package arena

import (
	"sort"
	"sync"
	"unsafe"
)

type BumpAllocator struct {
	chunks  [][]byte
	current int
	offset  int
	mtx     sync.Mutex
}

// NewBumpAllocator creates a new bump allocator with an initial chunk of the given size.
func NewBumpAllocator(size int) *BumpAllocator {
	return &BumpAllocator{
		chunks: [][]byte{MakePages(size)},
	}
}

// Alloc allocates memory of the specified size and alignment.
// It uses a bump allocation strategy, growing the heap as needed.
// Note: Pointers returned by Alloc become invalid after Reset() or Delete() and should not be used.
func (b *BumpAllocator) Alloc(size, align uint64) unsafe.Pointer {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	// log.Println("Allocating: ", size, align)
	// log.Println("current: ", b.current, "offset: ", b.offset)
	// log.Println("chunks: ", len(b.chunks))
	aligned := (b.offset + int(align-1)) &^ int(align-1)
	// log.Println("aligned: ", aligned)
	// log.Println("current chunk size: ", len(b.chunks[b.current]))
	if aligned+int(size) > len(b.chunks[b.current]) {
		// grow
		if b.current+1 >= len(b.chunks) {
			sz := max(int(size), len(b.chunks[0]))
			// log.Println("creating page with size: ", sz)
			b.chunks = append(b.chunks, MakePages(sz))
		}
		b.current++
		b.offset = 0
		aligned = 0
	}
	ptr := unsafe.Pointer(&b.chunks[b.current][aligned])
	b.offset = aligned + int(size)
	return ptr
}

// Reset resets the allocator to its initial state, allowing reuse of allocated memory.
// Note: All previously allocated pointers become invalid and should not be used.
func (b *BumpAllocator) Reset() {
	b.mtx.Lock()
	b.current, b.offset = 0, 0
	b.mtx.Unlock()
}

// Delete frees all memory allocated by the allocator.
// Note: All previously allocated pointers become invalid and should not be used.
func (b *BumpAllocator) Delete() {
	b.mtx.Lock()
	for _, c := range b.chunks {
		ReleasePages(c)
	}
	b.chunks = nil
	b.mtx.Unlock()
}

// Remove is a no-op for bump allocator, as individual deallocations are not supported.
// Note: This does not invalidate any pointers.
func (b *BumpAllocator) Remove(ptr unsafe.Pointer) {
	// no op for bump allocator
}

// Owns checks if the given pointer belongs to memory managed by this allocator.
func (b *BumpAllocator) Owns(ptr unsafe.Pointer) bool {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if ptr == nil {
		return false
	}

	ptrAddr := uintptr(ptr)

	// Binary search to find the chunk that might contain the pointer
	// Chunks are allocated sequentially, so they should be in increasing memory order
	idx := sort.Search(len(b.chunks), func(i int) bool {
		chunkStart := uintptr(unsafe.Pointer(&b.chunks[i][0]))
		return ptrAddr < chunkStart
	})

	// If idx == 0, pointer is before first chunk
	if idx == 0 {
		return false
	}

	// Check the previous chunk (idx - 1) since sort.Search returns insertion point
	chunk := b.chunks[idx-1]
	chunkStart := uintptr(unsafe.Pointer(&chunk[0]))
	chunkEnd := chunkStart + uintptr(len(chunk))

	return ptrAddr >= chunkStart && ptrAddr < chunkEnd
}
