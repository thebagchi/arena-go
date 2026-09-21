package alloc

import (
	"fmt"
	"sort"
	"sync"
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

const (
	// MIN_CLASS_SIZE is the smallest size class, and the smallest object that
	// can hold the free-list link stored inside a free slot.
	MIN_CLASS_SIZE = 16
	// CLASS_COUNT is how many power-of-two classes run from MIN_CLASS_SIZE.
	// The largest is 16<<11, or 32 KiB; anything bigger gets its own mapping.
	CLASS_COUNT = 12
	// OBJECTS_PER_SLAB is the least number of objects a slab is sized to hold.
	//
	// Slabs used to be sized to one page, so every class at or above the page
	// size got one object per mmap call: 100 objects of 8 KiB cost 100 syscalls.
	// Sizing by object count instead makes that one syscall per 16 objects.
	OBJECTS_PER_SLAB = 16
	// LARGE_CLASS marks a slab that holds a single over-sized object.
	LARGE_CLASS = -1
	// BITS_PER_WORD is the width of one word of a slab's allocation bitmap.
	BITS_PER_WORD = 64
	// NOT_LISTED marks a slab that is in no bin's partial list.
	NOT_LISTED = -1
)

// SlabAllocator serves fixed-size classes from slabs of objects and supports
// individual free.
//
// One mutex guards everything. The previous design had a second mutex per bin
// and took the two in one order when allocating and the other when freeing,
// which deadlocked, and it mutated the shared page pool while holding only the
// bin lock, which raced. Neither is expressible with a single lock.
//
// Every slab carries an allocation bitmap beside its free list. The bitmap is
// what makes a double free and an interior pointer detectable: a free list
// alone has no memory of which slots are already on it, so freeing a pointer
// twice linked it to itself and handed the same address out twice.
type SlabAllocator struct {
	mtx     sync.Mutex
	table   *res.PageTable
	free    map[int][]*res.Page
	bins    [CLASS_COUNT]*_Bin
	slabs   []*_Slab
	deleted bool
	sorted  bool
}

// _Bin is one size class: the object size and the slabs that still have a free
// slot, most recently used last.
type _Bin struct {
	partial []*_Slab
	size    int
}

// _Slab is one page carved into equal slots.
//
// Slots are handed out in two stages. While next is below capacity the slab is
// still walking up untouched memory, which on a freshly mapped page is already
// zero and so costs nothing to hand over. Once something is freed it goes on the
// list headed by freed, and a slot taken from there is zeroed on the way out.
//
// bitmap records which slots are live, which is what makes a double free and an
// interior pointer detectable. The free list alone has no memory of that.
type _Slab struct {
	page     *res.Page
	bitmap   []uint64
	freed    unsafe.Pointer
	size     int
	capacity int
	used     int
	next     int
	class    int
	listed   int
	zeroed   bool
}

// NewSlabAllocator returns a slab allocator that has mapped nothing yet.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func NewSlabAllocator() *SlabAllocator {
	a := &SlabAllocator{
		table: res.NewPageTable(),
		free:  make(map[int][]*res.Page),
	}

	for i := range a.bins {
		a.bins[i] = &_Bin{size: MIN_CLASS_SIZE << uint(i)}
	}

	return a
}

// Alloc returns zeroed memory, or panics with res.ErrOutOfMemory.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Alloc(size, align uint64) unsafe.Pointer {
	ptr, ok := a.TryAlloc(size, align)
	if !ok {
		panic(res.ErrOutOfMemory)
	}

	return ptr
}

// TryAlloc is Alloc, reporting exhaustion instead of panicking.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) TryAlloc(size, align uint64) (unsafe.Pointer, bool) {
	size, align = res.Request(size, align)
	want := int(max(size, align))

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		panic(res.ErrDeleted)
	}

	slab, ok := a._SlabFor(want)
	if !ok {
		return nil, false
	}

	ptr, dirty := slab._Take()
	if slab._Full() {
		a._Unlist(slab)
	}

	// Untouched memory on a freshly mapped page is already zero: the kernel
	// guarantees it. Only a reused slot, or a recycled page, has to be cleared.
	if dirty {
		clear(unsafe.Slice((*byte)(ptr), slab.size))
	}

	return ptr, true
}

// Reset frees every object while keeping the pages for reuse.
//
// The pages go back to this allocator's own pool. They used to be left in the
// page table with nothing referencing them, so each Reset cycle mapped a fresh
// set and the footprint grew without bound.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Reset() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	for _, slab := range a.slabs {
		size := slab.page.Size()
		a.free[size] = append(a.free[size], slab.page)
	}

	a.slabs, a.sorted = nil, true

	for _, bin := range a.bins {
		bin.partial = nil
	}
}

// Delete unmaps everything. A failed unmap means this allocator's record of its
// pages disagrees with the kernel, which is a bug here rather than a condition
// a caller can handle.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Delete() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	a.slabs, a.free, a.deleted, a.sorted = nil, nil, true, true

	for _, bin := range a.bins {
		bin.partial = nil
	}

	if err := a.table.Delete(); err != nil {
		panic(fmt.Errorf("alloc: slab delete: %w", err))
	}
}

// Remove frees one object.
//
// It ignores a pointer this allocator did not hand out, a pointer into the
// middle of an object, and a second free of the same object. Each of those was
// previously accepted, and the last one linked a slot to itself so that the next
// two allocations returned the same address.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Remove(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	slab, ok := a._FindSlab(ptr)
	if !ok {
		return
	}

	offset := int(uintptr(ptr) - slab.page.Start())
	if offset%slab.size != 0 {
		return
	}

	slot := offset / slab.size
	if slot >= slab.capacity || !slab._Marked(slot) {
		return
	}

	slab._Give(ptr, slot)
	a._Relist(slab)
}

// Owns reports whether ptr came from this allocator.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Owns(ptr unsafe.Pointer) bool {
	return a.table.Owns(ptr)
}

// Stats reports how much memory this allocator has mapped.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) Stats() res.Stats {
	return res.Stats{Mapped: a.table.Size(), Mappings: a.table.Len()}
}

// _SlabFor returns a slab of the right class with a free slot, creating one when
// no existing slab has room.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _SlabFor(want int) (*_Slab, bool) {
	class := _ClassOf(want)
	if class == LARGE_CLASS {
		return a._NewSlab(LARGE_CLASS, want)
	}

	bin := a.bins[class]
	if n := len(bin.partial); n > 0 {
		return bin.partial[n-1], true
	}

	return a._NewSlab(class, bin.size)
}

// _NewSlab maps or reuses a page and carves it into slots of the given size.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _NewSlab(class, size int) (*_Slab, bool) {
	want := res.RoundUp(size*OBJECTS_PER_SLAB, res.PAGE_SIZE)
	if class == LARGE_CLASS {
		want = res.RoundUp(size, res.PAGE_SIZE)
	}

	page, fresh, ok := a._Page(want)
	if !ok {
		return nil, false
	}

	slab := _NewSlabOf(page, class, size, fresh)
	a._Track(slab)

	if class != LARGE_CLASS {
		bin := a.bins[class]
		slab.listed = len(bin.partial)
		bin.partial = append(bin.partial, slab)
	}

	return slab, true
}

// _Page takes a page of exactly size bytes from the recycled pool, or maps one.
//
// The pool is keyed by size rather than scanned for a page that is big enough,
// so reuse is exact and costs no search.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Page(size int) (*res.Page, bool, bool) {
	if pool := a.free[size]; len(pool) > 0 {
		page := pool[len(pool)-1]
		a.free[size] = pool[:len(pool)-1]

		return page, false, true
	}

	page, err := a.table.New(size)
	if err != nil {
		return nil, false, false
	}

	return page, true, true
}

// _Track records a slab in address order, so _FindSlab is a binary search.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Track(slab *_Slab) {
	a.slabs = append(a.slabs, slab)
	a.sorted = false
}

// _Sort puts the slabs in address order, which _FindSlab needs and an append
// does not maintain.
//
// Sorting on insert instead is quadratic, because mmap hands out addresses
// downward and every new slab therefore sorts before all the others.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Sort() {
	if a.sorted {
		return
	}

	sort.Slice(a.slabs, func(i, j int) bool {
		return a.slabs[i].page.Start() < a.slabs[j].page.Start()
	})

	a.sorted = true
}

// _FindSlab returns the slab containing ptr.
//
// Freeing used to walk every bin and every slab in it, twice, under as many as
// thirty-two lock acquisitions. This is one binary search over all slabs.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _FindSlab(ptr unsafe.Pointer) (*_Slab, bool) {
	a._Sort()

	addr := uintptr(ptr)

	idx := sort.Search(len(a.slabs), func(i int) bool {
		return a.slabs[i].page.End() > addr
	})
	if idx < len(a.slabs) && a.slabs[idx].page.Contains(addr) {
		return a.slabs[idx], true
	}

	return nil, false
}

// _Unlist drops a slab from its bin's partial list once it is full.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Unlist(slab *_Slab) {
	if slab.class == LARGE_CLASS || slab.listed == NOT_LISTED {
		return
	}

	var (
		bin  = a.bins[slab.class]
		last = len(bin.partial) - 1
		move = bin.partial[last]
	)

	bin.partial[slab.listed] = move
	move.listed = slab.listed
	bin.partial = bin.partial[:last]
	slab.listed = NOT_LISTED
}

// _Relist puts a slab back in its bin after a free, and recycles it once it is
// empty and some other slab in the bin can serve the next request.
//
// Keeping the last slab of a class avoids unmapping and remapping a page on
// every alloc/free pair at the class boundary.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Relist(slab *_Slab) {
	if slab.class == LARGE_CLASS {
		if slab.used == 0 {
			a._Recycle(slab)
		}

		return
	}

	bin := a.bins[slab.class]

	if slab.listed == NOT_LISTED {
		slab.listed = len(bin.partial)
		bin.partial = append(bin.partial, slab)
	}

	if slab.used == 0 && len(bin.partial) > 1 {
		a._Unlist(slab)
		a._Recycle(slab)
	}
}

// _Recycle returns an empty slab's page to the pool for reuse.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *SlabAllocator) _Recycle(slab *_Slab) {
	a._Sort()

	idx := sort.Search(len(a.slabs), func(i int) bool {
		return a.slabs[i].page.Start() >= slab.page.Start()
	})
	if idx < len(a.slabs) && a.slabs[idx] == slab {
		a.slabs = append(a.slabs[:idx], a.slabs[idx+1:]...)
	}

	size := slab.page.Size()
	a.free[size] = append(a.free[size], slab.page)
}

// _NewSlabOf carves a page into slots and threads a free list through them.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func _NewSlabOf(page *res.Page, class, size int, fresh bool) *_Slab {
	slot := size
	if class == LARGE_CLASS {
		slot = page.Size()
	}

	capacity := page.Size() / slot

	// No free list is threaded through the slots here. Writing one would touch
	// every slot on the page, which costs a pass over the whole slab and dirties
	// memory the kernel had already zeroed. Slots are handed out by walking up
	// instead, and the list only ever holds slots that were freed.
	return &_Slab{
		page:     page,
		bitmap:   make([]uint64, (capacity+BITS_PER_WORD-1)/BITS_PER_WORD),
		size:     slot,
		capacity: capacity,
		class:    class,
		listed:   NOT_LISTED,
		zeroed:   fresh,
	}
}

// _ClassOf returns the size class index for a request, or LARGE_CLASS when the
// request is bigger than the largest class.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func _ClassOf(size int) int {
	rounded := res.RoundPow2(uint64(max(size, MIN_CLASS_SIZE)))

	class := int(res.Log2(rounded) - res.Log2(MIN_CLASS_SIZE))
	if class >= CLASS_COUNT {
		return LARGE_CLASS
	}

	return class
}

// _Take hands out a slot, reporting whether it holds stale bytes that the
// caller has to clear.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (s *_Slab) _Take() (unsafe.Pointer, bool) {
	var (
		ptr   unsafe.Pointer
		slot  int
		dirty bool
	)

	if s.freed != nil {
		ptr = s.freed
		s.freed = *(*unsafe.Pointer)(ptr)
		slot = int(uintptr(ptr)-s.page.Start()) / s.size
		dirty = true
	} else {
		slot = s.next
		ptr = unsafe.Add(s.page.Ptr(), slot*s.size)
		s.next = s.next + 1
		dirty = !s.zeroed
	}

	s.used = s.used + 1
	s.bitmap[slot/BITS_PER_WORD] |= 1 << uint(slot%BITS_PER_WORD)

	return ptr, dirty
}

// _Full reports whether the slab has no slot left to hand out.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (s *_Slab) _Full() bool {
	return s.freed == nil && s.next >= s.capacity
}

// _Give returns a slot to the free list and marks it free.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (s *_Slab) _Give(ptr unsafe.Pointer, slot int) {
	s.bitmap[slot/BITS_PER_WORD] &^= 1 << uint(slot%BITS_PER_WORD)
	*(*unsafe.Pointer)(ptr) = s.freed
	s.freed = ptr
	s.used = s.used - 1
}

// _Marked reports whether a slot is currently handed out.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (s *_Slab) _Marked(slot int) bool {
	return s.bitmap[slot/BITS_PER_WORD]&(1<<uint(slot%BITS_PER_WORD)) != 0
}
