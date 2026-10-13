package res

import (
	"sync"
	"unsafe"
)

const (
	// DEFAULT_ALIGN is the alignment used when a caller asks for none.
	DEFAULT_ALIGN = 8
	// CHUNK_GROWTH multiplies each new chunk's size by the last one's.
	//
	// Chunks used to be a fixed size, so a large working set on a small arena
	// meant one mmap per chunk's worth of bytes and an ownership scan over all
	// of them. Doubling bounds both at the cost of one partly used chunk.
	CHUNK_GROWTH = 2
	// MAX_CHUNK_SIZE caps that doubling, so a long-lived arena does not reserve
	// hundreds of megabytes to satisfy one more small allocation.
	MAX_CHUNK_SIZE = 64 << 20
)

// Bump hands out bytes by advancing a cursor through a growing list of pages.
// It is the fastest allocator here and the one with the least bookkeeping:
// nothing is freed individually, and Reset returns the whole arena at once.
//
// Memory handed out is always zero. Fresh mappings are zeroed by the kernel, and
// Reset clears exactly the bytes that were handed out since the last one, so the
// allocation path itself never pays for zeroing.
//
// The mutex guards every field below it, and is the only lock on the path: the
// allocator wrapping this type adds none of its own.
type Bump struct {
	mtx     sync.Mutex
	chunks  []*_Chunk
	table   *PageTable
	current int
	next    int
	deleted bool
}

// _Chunk is one page plus the two offsets Bump keeps for it: how far the cursor
// has advanced, and how far it ever advanced since the last Reset.
type _Chunk struct {
	page  *Page
	used  int
	dirty int
}

// NewBump returns a bump allocator whose first chunk holds at least size bytes.
//
// Revisions:
//   - 2026-09-21 11:10: initial creation
func NewBump(size int) *Bump {
	first := max(size, PAGE_SIZE)

	return &Bump{
		table: NewPageTable(),
		next:  first,
	}
}

// Alloc returns zeroed memory of the given size and alignment.
//
// It panics with ErrOutOfMemory rather than returning nil, because every caller
// of a nil result here dereferences it: the generic helpers in the arena package
// would crash somewhere further on, naming neither the arena nor the cause. Use
// TryAlloc where exhaustion is an expected outcome.
//
// Revisions:
//   - 2026-09-21 11:12: initial creation
func (b *Bump) Alloc(size, align uint64) unsafe.Pointer {
	ptr, ok := b.TryAlloc(size, align)
	if !ok {
		panic(ErrOutOfMemory)
	}

	return ptr
}

// TryAlloc is Alloc, reporting exhaustion instead of panicking.
//
// Revisions:
//   - 2026-09-21 11:14: initial creation
//   - 2026-10-13 16:25: a size that does not fit in an int reports exhaustion
func (b *Bump) TryAlloc(size, align uint64) (unsafe.Pointer, bool) {
	size, align, fits := Request(size, align)

	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.deleted {
		panic(ErrDeleted)
	}

	if !fits {
		return nil, false
	}

	if len(b.chunks) == 0 && !b._Grow(int(size+align)) {
		return nil, false
	}

	for {
		if ptr, ok := b._FromChunk(b.chunks[b.current], size, align); ok {
			return ptr, true
		}

		if b.current+1 >= len(b.chunks) && !b._Grow(int(size+align)) {
			return nil, false
		}

		b.current++
	}
}

// Reset rewinds the cursor and zeroes the bytes handed out since the last Reset,
// so the next allocation sees zeroed memory without the allocation path paying
// for it. Every pointer handed out before the call becomes invalid.
//
// Revisions:
//   - 2026-09-21 11:17: initial creation
//   - 2026-10-13 16:25: panic with ErrDeleted after Delete, as Alloc does
func (b *Bump) Reset() {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.deleted {
		panic(ErrDeleted)
	}

	for _, chunk := range b.chunks {
		if chunk.dirty > 0 {
			clear(chunk.page.Base()[:chunk.dirty])
			chunk.dirty = 0
		}

		chunk.used = 0
	}

	b.current = 0
}

// Delete unmaps every page. Any later allocation panics with ErrDeleted.
//
// Revisions:
//   - 2026-09-21 11:19: initial creation
func (b *Bump) Delete() error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.deleted {
		return nil
	}

	b.chunks, b.current, b.deleted = nil, 0, true

	return b.table.Delete()
}

// Remove is a no-op apart from the deleted check. A bump allocator frees only
// in whole arenas, and a call after Delete is the same lifetime bug as Alloc.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func (b *Bump) Remove(ptr unsafe.Pointer) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.deleted {
		panic(ErrDeleted)
	}
}

// Owns reports whether ptr came from this allocator.
//
// Revisions:
//   - 2026-09-21 11:20: initial creation
func (b *Bump) Owns(ptr unsafe.Pointer) bool {
	return b.table.Owns(ptr)
}

// Size returns the total bytes mapped.
//
// Revisions:
//   - 2026-09-21 11:21: initial creation
func (b *Bump) Size() int {
	return b.table.Size()
}

// Mappings returns the number of distinct mappings held.
//
// Revisions:
//   - 2026-09-21 11:42: initial creation
func (b *Bump) Mappings() int {
	return b.table.Len()
}

// _FromChunk advances the cursor in one chunk, reporting whether it fitted.
//
// Revisions:
//   - 2026-09-21 11:22: initial creation
//   - 2026-10-13 16:25: reject a request by the room left, so a size that would
//     wrap the end address is not treated as fitting
func (b *Bump) _FromChunk(chunk *_Chunk, size, align uint64) (unsafe.Pointer, bool) {
	var (
		base    = chunk.page.Start()
		aligned = AlignUp(base+uintptr(chunk.used), uintptr(align))
	)

	if aligned < base || aligned >= chunk.page.End() {
		return nil, false
	}

	room := uint64(chunk.page.End() - aligned)
	if size > room {
		return nil, false
	}

	offset := int(aligned - base)
	chunk.used = offset + int(size)

	if chunk.used > chunk.dirty {
		chunk.dirty = chunk.used
	}

	return unsafe.Add(chunk.page.Ptr(), offset), true
}

// _Grow appends a chunk large enough for need bytes, doubling the chunk size
// until it reaches MAX_CHUNK_SIZE.
//
// Revisions:
//   - 2026-09-21 11:24: initial creation
func (b *Bump) _Grow(need int) bool {
	size := max(b.next, need)

	page, err := b.table.New(size)
	if err != nil {
		return false
	}

	b.chunks = append(b.chunks, &_Chunk{page: page})
	b.next = min(b.next*CHUNK_GROWTH, MAX_CHUNK_SIZE)

	return true
}
