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
//
// listed records every pointer currently on the free list. The list alone
// cannot tell a second free of the same object from a first one, and the
// second one used to link the object to itself.
type Pool[T any] struct {
	mtx    sync.Mutex
	arena  *arena.Arena
	freed  *Vec[*T]
	listed map[uintptr]struct{}
	size   uintptr
}

// NewPool returns a pool that allocates T from the arena.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
//   - 2026-10-13 16:25: record which objects are on the free list
func NewPool[T any](a *arena.Arena) *Pool[T] {
	size := unsafe.Sizeof(*new(T))
	if size == 0 {
		size = 1
	}

	return &Pool[T]{
		arena:  a,
		freed:  NewVec[*T](a),
		listed: make(map[uintptr]struct{}),
		size:   (size + POOL_ALIGN - 1) &^ (POOL_ALIGN - 1),
	}
}

// Alloc returns a zeroed object, reusing a freed one when there is one.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
//   - 2026-10-13 16:25: drop the object from the free-list set before reuse
func (p *Pool[T]) Alloc() *T {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if ptr, ok := p.freed.Pop(); ok {
		delete(p.listed, uintptr(unsafe.Pointer(ptr)))

		// A pooled object never went back to the allocator, so nothing has
		// zeroed it since the caller last used it.
		clear(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), p.size))

		return ptr
	}

	return (*T)(p.arena.Alloc(uint64(p.size), POOL_ALIGN))
}

// Free returns an object to the pool for reuse. Freeing nil, a pointer from
// another arena, or an object already on the free list does nothing.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
//   - 2026-10-13 16:25: reject a second free of an object already waiting
func (p *Pool[T]) Free(obj *T) {
	if obj == nil {
		return
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	addr := uintptr(unsafe.Pointer(obj))
	if _, found := p.listed[addr]; found {
		return
	}

	if !p.arena.Owns(unsafe.Pointer(obj)) {
		return
	}

	p.listed[addr] = struct{}{}
	p.freed.Push(obj)
}

// Reset drops the free list. The objects stay in the arena until it is reset.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
//   - 2026-10-13 16:25: forget which objects were on the free list
func (p *Pool[T]) Reset() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.freed.Clear()
	clear(p.listed)
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
