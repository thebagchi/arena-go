// arena/arena.go
//
// Package arena provides high-performance, zero-GC memory allocators with multiple strategies.
//
// Thread Safety:
//   - All allocators (Bump, Slab, Buddy) are thread-safe and can be used concurrently
//   - Alloc() operations are serialized with mutexes to prevent data races
//   - Reset() and Delete() should NOT be called concurrently with Alloc() or with each other
//   - Multiple Arena instances are completely independent and require no synchronization
//
// Memory Model:
//   - All memory is allocated via mmap and lives outside Go's garbage collector
//   - Memory is never returned to the OS until Delete() is called
//   - Reset() clears allocations but retains underlying memory pages
//
// Allocator Strategies:
//   - BUMP: Fastest, best for batch allocations or when arena is reset frequently
//   - SLAB: Best for fixed-size objects with high allocation/free turnover
//   - BUDDY: Most flexible, good for varied-size allocations with power-of-2 sizes
package arena

import (
	"syscall"
	"unsafe"
)

// ---------------------------------------------------------------
// Public API – one arena for all types
// ---------------------------------------------------------------

type Type int

const (
	BUMP Type = iota
	SLAB
	BUDDY
)

// Arena is the beautiful multi-type facade.
// Thread-safe: Multiple goroutines can safely call Alloc concurrently.
// The underlying allocator handles synchronization internally.
type Arena struct {
	Allocator
}

// New creates an arena. pages == 0 → 1 page (4 KiB default)
func New(pages int, alloc Type) *Arena {
	if pages <= 0 {
		pages = 1 // ← your request: treat 0 as 1
	}
	size := pages * syscall.Getpagesize()

	var raw Allocator
	switch alloc {
	case BUMP:
		raw = NewBumpAllocator(size)
	case SLAB:
		raw = NewSlabAllocator(256, size) // configurable block size
	case BUDDY:
		raw = NewBuddyAllocator(syscall.Getpagesize(), pages)
	default:
		raw = NewBumpAllocator(size)
	}
	return &Arena{Allocator: raw}
}

func (a *Arena) Reset() {
	a.Allocator.Reset()
}
func (a *Arena) Delete() {
	a.Allocator.Delete()
}

// Owns checks if the given pointer belongs to memory managed by this arena.
// Returns true if the pointer was allocated by this arena and is still valid.
// Returns false for nil pointers or pointers not managed by this arena.
func (a *Arena) Owns(ptr unsafe.Pointer) bool {
	return a.Allocator.Owns(ptr)
}

// ---------------------------------------------------------------
// Internal raw allocators (all support growing)
// ---------------------------------------------------------------

type Allocator interface {
	Alloc(size, align uint64) unsafe.Pointer
	Reset()
	Delete()
	Remove(ptr unsafe.Pointer)
	Owns(ptr unsafe.Pointer) bool
}
