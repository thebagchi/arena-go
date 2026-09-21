package container

import (
	"sync"
	"unsafe"

	arena "github.com/thebagchi/arena-go"
)

const (
	// POOL_ALIGN is the alignment every pooled object gets, one cache-friendly
	// step up from the natural alignment of most structs.
	POOL_ALIGN = 16
)

// Pool hands out objects of one type from an arena and keeps freed ones for
// reuse, which suits nodes that are created and discarded in large numbers.
//
// Safe for concurrent use.
type Pool[T any] struct {
	mtx   sync.Mutex
	arena *arena.Arena
	freed *Vec[*T]
	size  uintptr
}

// NewPool returns a pool that allocates T from the arena.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func NewPool[T any](a *arena.Arena) *Pool[T] {
	size := unsafe.Sizeof(*new(T))
	if size == 0 {
		size = 1
	}

	return &Pool[T]{
		arena: a,
		freed: NewVec[*T](a),
		size:  (size + POOL_ALIGN - 1) &^ (POOL_ALIGN - 1),
	}
}

// Alloc returns a zeroed object, reusing a freed one when there is one.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (p *Pool[T]) Alloc() *T {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if ptr, ok := p.freed.Pop(); ok {
		// A pooled object never went back to the allocator, so nothing has
		// zeroed it since the caller last used it.
		clear(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), p.size))

		return ptr
	}

	return (*T)(p.arena.Alloc(uint64(p.size), POOL_ALIGN))
}

// Free returns an object to the pool for reuse. Freeing nil does nothing.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (p *Pool[T]) Free(obj *T) {
	if obj == nil {
		return
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.freed.Push(obj)
}

// Reset drops the free list. The objects stay in the arena until it is reset.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (p *Pool[T]) Reset() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.freed.Clear()
}

// Len returns how many freed objects are waiting for reuse.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (p *Pool[T]) Len() int {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.freed.Len()
}

// Cap returns the capacity of the free list.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (p *Pool[T]) Cap() int {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.freed.Cap()
}
