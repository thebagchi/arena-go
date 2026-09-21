package alloc

import (
	"fmt"
	"sort"
	"sync"
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

const (
	// DEFAULT_CHUNK_SIZE is the payload of a chunk when none is configured.
	DEFAULT_CHUNK_SIZE = 256 << 10
	// NO_CHUNK marks that no chunk has served an allocation yet.
	NO_CHUNK = -1
)

// Growth decides what a buddy allocator does when every chunk is full.
type Growth int

const (
	// FIXED keeps one chunk and fails once it is full.
	FIXED Growth = iota
	// ADDITIVE maps another chunk.
	ADDITIVE
)

// BuddyAllocator serves power-of-two blocks from chunks that split and merge,
// which is what lets it free individual blocks of any size and still hand the
// space back as one piece.
//
// Chunks are kept sorted by address so a free is a binary search, and the chunk
// that served the last allocation is tried first.
type BuddyAllocator struct {
	mtx     sync.Mutex
	table   *res.PageTable
	chunks  []*Chunk
	size    int
	last    int
	growth  Growth
	deleted bool
	sorted  bool
}

// BuddyOption configures a BuddyAllocator at construction.
type BuddyOption func(*BuddyAllocator)

// NewBuddyAllocator returns a buddy allocator that has mapped nothing yet.
//
// Nothing is mapped until the first allocation. It used to reserve a chunk's
// worth of pages that no allocation ever reached, because chunks were mapped
// separately from it.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func NewBuddyAllocator(opts ...BuddyOption) *BuddyAllocator {
	a := &BuddyAllocator{
		table:  res.NewPageTable(),
		size:   DEFAULT_CHUNK_SIZE,
		last:   NO_CHUNK,
		growth: ADDITIVE,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// WithGrowthStrategy sets what happens when every chunk is full.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func WithGrowthStrategy(growth Growth) BuddyOption {
	return func(a *BuddyAllocator) {
		a.growth = growth
	}
}

// WithSize sets the payload of each chunk, rounded up to a power of two.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func WithSize(size int) BuddyOption {
	return func(a *BuddyAllocator) {
		if size > 0 {
			a.size = int(res.RoundPow2(uint64(size)))
		}
	}
}

// Alloc returns zeroed memory, or panics with res.ErrOutOfMemory.
//
// Alignment up to the system page size is guaranteed: a block of a given size
// sits at a multiple of that size from a page-aligned base.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Alloc(size, align uint64) unsafe.Pointer {
	ptr, ok := a.TryAlloc(size, align)
	if !ok {
		panic(res.ErrOutOfMemory)
	}

	return ptr
}

// TryAlloc is Alloc, reporting exhaustion instead of panicking. A fixed-growth
// allocator reaches that state in normal use, which is why it is reachable
// without recovering from a panic.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) TryAlloc(size, align uint64) (unsafe.Pointer, bool) {
	size, align = res.Request(size, align)

	// The block has to cover the alignment as well as the size, or the block's
	// own alignment says nothing about the pointer handed back. Passing only the
	// size is what made Alloc(8, 64) return a 16-byte-aligned pointer.
	block := int(max(res.RoundPow2(size), align, MIN_BLOCK_SIZE))

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		panic(res.ErrDeleted)
	}

	ptr, ok := a._FromChunks(block)
	if !ok {
		if ptr, ok = a._Grow(block); !ok {
			return nil, false
		}
	}

	clear(unsafe.Slice((*byte)(ptr), size))

	return ptr, true
}

// Reset frees every block in every chunk, keeping the chunks mapped.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Reset() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	for _, chunk := range a.chunks {
		chunk.Reset()
	}

	a.last = NO_CHUNK
}

// Delete unmaps every chunk. A failed unmap means this allocator's record of its
// pages disagrees with the kernel, which is a bug here rather than a condition a
// caller can handle.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Delete() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	a.chunks, a.last, a.deleted, a.sorted = nil, NO_CHUNK, true, true

	if err := a.table.Delete(); err != nil {
		panic(fmt.Errorf("alloc: buddy delete: %w", err))
	}
}

// Remove frees the block containing ptr, ignoring a pointer from elsewhere.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Remove(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.deleted {
		return
	}

	if chunk, ok := a._FindChunk(uintptr(ptr)); ok {
		chunk.Free(ptr)
	}
}

// Owns reports whether ptr came from this allocator.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Owns(ptr unsafe.Pointer) bool {
	return a.table.Owns(ptr)
}

// Stats reports how much memory this allocator has mapped.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) Stats() res.Stats {
	return res.Stats{Mapped: a.table.Size(), Mappings: a.table.Len()}
}

// _FromChunks tries the chunk that served the last allocation, then the rest.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) _FromChunks(block int) (unsafe.Pointer, bool) {
	if a.last >= 0 && a.last < len(a.chunks) {
		if ptr, ok := a.chunks[a.last].Allocate(block); ok {
			return ptr, true
		}
	}

	for i, chunk := range a.chunks {
		if chunk.Largest() < block {
			continue
		}

		if ptr, ok := chunk.Allocate(block); ok {
			a.last = i

			return ptr, true
		}
	}

	return nil, false
}

// _Grow maps another chunk and allocates from it.
//
// The new chunk is the configured size, or just big enough for an over-sized
// request. Growth used to ask for twice the rounded request, so a one-megabyte
// allocation mapped two megabytes.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) _Grow(block int) (unsafe.Pointer, bool) {
	if a.growth == FIXED && len(a.chunks) > 0 {
		return nil, false
	}

	chunk, err := NewChunk(a.table, max(a.size, block))
	if err != nil {
		return nil, false
	}

	a.chunks = append(a.chunks, chunk)
	a.last = len(a.chunks) - 1
	a.sorted = false

	return chunk.Allocate(block)
}

// _FindChunk returns the chunk containing addr.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) _FindChunk(addr uintptr) (*Chunk, bool) {
	a._Sort()

	idx := sort.Search(len(a.chunks), func(i int) bool {
		return a.chunks[i].Start()+uintptr(a.chunks[i].Size()) > addr
	})
	if idx < len(a.chunks) && a.chunks[idx].Contains(addr) {
		return a.chunks[idx], true
	}

	return nil, false
}

// _Sort puts the chunks in address order, which _FindChunk needs and an append
// does not maintain.
//
// Sorting on insert instead is quadratic, because mmap hands out addresses
// downward and every new chunk therefore sorts before all the others. It also
// invalidates the cached index, so that is reset with it.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (a *BuddyAllocator) _Sort() {
	if a.sorted {
		return
	}

	sort.Slice(a.chunks, func(i, j int) bool {
		return a.chunks[i].Start() < a.chunks[j].Start()
	})

	a.last, a.sorted = NO_CHUNK, true
}
