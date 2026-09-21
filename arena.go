// Package arena provides zero-GC memory allocators with several strategies, and
// the generic helpers that allocate Go values inside them.
//
// # The rule that governs everything else
//
// Arena memory comes from mmap and is invisible to Go's garbage collector. A Go
// pointer written into arena memory does not keep its target alive, so a value
// whose only reference lives in an arena will be collected and its storage
// reused while the arena still points at it. This is a use-after-free that the
// race detector cannot see and that only appears under memory pressure.
//
// What that rules out: storing a heap string, slice, map, channel, func,
// interface or pointer in arena memory. What it permits: values with no pointers
// at all, and pointers that lead back into the same arena.
//
// Bring outside data in by copying it:
//
//	name := a.MakeString(userInput) // copied into the arena, safe to store
//	node := arena.Ptr(a, Node{})    // lives in the arena, safe to point at
//
// The containers in the container package follow the same rule. Map copies
// string keys for exactly this reason; a key or value of some other
// pointer-bearing type is the caller's responsibility.
//
// # Allocator contract
//
// Every allocator here promises the same things, and the Allocator interface
// states them once:
//
//   - Alloc returns memory that is zeroed, or panics with ErrOutOfMemory. It
//     never returns nil, because every caller of a nil result dereferences it.
//   - TryAlloc reports exhaustion instead of panicking.
//   - align must be a power of two; zero means the natural default. Alignment up
//     to the system page size is honoured.
//   - Reset invalidates every pointer and keeps the memory. Delete releases it,
//     after which any allocation panics with ErrDeleted.
//   - Remove frees one allocation where the strategy supports it and is a no-op
//     where it does not. It must be given the pointer that was handed out.
//
// # Thread safety
//
// All three allocators are safe for concurrent use. So are Map, SkipList and
// Pool. Vec, Queue, Stack, Buffer and Str are not: they are the containers whose
// operations return interior pointers, which no lock inside them could protect.
// Give each goroutine its own arena, or guard a shared container yourself.
//
// # Strategies
//
//   - Bump is fastest and frees only in whole arenas.
//   - Slab suits fixed-size objects with high turnover.
//   - Buddy suits varied sizes that must be freed individually.
package arena

import (
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

var (
	// ErrOutOfMemory is what Alloc panics with when it cannot satisfy a request.
	ErrOutOfMemory = res.ErrOutOfMemory
	// ErrDeleted is what any operation panics with after Delete.
	ErrDeleted = res.ErrDeleted
)

// Allocator is the strategy an Arena is built on. The package documentation
// states the contract every implementation keeps.
type Allocator interface {
	Alloc(size, align uint64) unsafe.Pointer
	TryAlloc(size, align uint64) (unsafe.Pointer, bool)
	Reset()
	Delete()
	Remove(ptr unsafe.Pointer)
	Owns(ptr unsafe.Pointer) bool
}

// Arena is the multi-type facade over one allocator.
type Arena struct {
	Allocator
}

// New returns an arena backed by the given allocator.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func New(allocator Allocator) *Arena {
	return &Arena{Allocator: allocator}
}

// MakeString copies s into the arena and returns a string sharing that copy.
//
// Copying is the point: a string built on the Go heap and stored in arena memory
// would not be kept alive by that reference. The result is safe to store.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func (a *Arena) MakeString(s string) string {
	if len(s) == 0 {
		return ""
	}

	ptr := a.Alloc(uint64(len(s)), 1)
	copy(unsafe.Slice((*byte)(ptr), len(s)), s)

	return unsafe.String((*byte)(ptr), len(s))
}

// Alloc returns a pointer to a zeroed T in the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func Alloc[T any](a *Arena) *T {
	var zero T

	size, align := _Shape(zero)

	return (*T)(a.Alloc(size, align))
}

// MakeObject returns a pointer to a zeroed T in the arena. It is Alloc under the
// name that reads better for structs.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func MakeObject[T any](a *Arena) *T {
	return Alloc[T](a)
}

// Ptr copies value into the arena and returns a pointer to the copy.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func Ptr[T any](a *Arena, value T) *T {
	ptr := Alloc[T](a)
	*ptr = value

	return ptr
}

// MakeSlice returns a zeroed slice of the given length and capacity, backed by
// arena memory.
//
// It panics on the argument combinations the built-in make rejects, and on an
// arena that cannot satisfy the request.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func MakeSlice[T any](a *Arena, length, capacity int) []T {
	// The same arguments the built-in make rejects, rejected the same way: a
	// silent nil would be a slice the caller writes to and a crash somewhere
	// else.
	if length < 0 || capacity < 0 || length > capacity {
		panic("arena: invalid slice length or capacity")
	}

	if capacity == 0 {
		return nil
	}

	var zero T

	size, align := _Shape(zero)
	if uint64(capacity) > MAX_ALLOC_BYTES/size {
		panic("arena: slice allocation size overflow")
	}

	ptr := a.Alloc(uint64(capacity)*size, align)

	return unsafe.Slice((*T)(ptr), capacity)[:length]
}

// Append adds elements to an arena-backed slice, moving it to a larger block
// when it runs out of capacity and releasing the old block.
//
// The slice must be one MakeSlice or a previous Append returned, not a re-slice
// of one: growing frees the old backing by the address of its first element, and
// for a re-sliced slice that address is inside the block rather than its start.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func Append[T any](a *Arena, slice []T, elems ...T) []T {
	if len(elems) == 0 {
		return slice
	}

	length := len(slice) + len(elems)
	if length <= cap(slice) {
		grown := slice[:length]
		copy(grown[len(slice):], elems)

		return grown
	}

	capacity := max(cap(slice)*GROWTH_FACTOR, length, MIN_SLICE_CAPACITY)
	grown := MakeSlice[T](a, length, capacity)

	copy(grown, slice)
	copy(grown[len(slice):], elems)

	if len(slice) > 0 {
		a.Remove(unsafe.Pointer(&slice[0]))
	}

	return grown
}

// OwnsPtr reports whether ptr points into the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func OwnsPtr[T any](a *Arena, ptr *T) bool {
	return a.Owns(unsafe.Pointer(ptr))
}

// OwnsSlice reports whether a slice's backing array is in the arena. An empty
// slice has no backing array, so it belongs to no arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func OwnsSlice[T any](a *Arena, slice []T) bool {
	if len(slice) == 0 {
		return false
	}

	return a.Owns(unsafe.Pointer(unsafe.SliceData(slice)))
}

// OwnsString reports whether a string's bytes are in the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func OwnsString(a *Arena, s string) bool {
	if len(s) == 0 {
		return false
	}

	return a.Owns(unsafe.Pointer(unsafe.StringData(s)))
}

// CloneObject returns a heap copy of an arena object, which outlives the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func CloneObject[T any](obj *T) *T {
	if obj == nil {
		return nil
	}

	result := new(T)
	*result = *obj

	return result
}

// CloneSlice returns a heap copy of an arena slice, which outlives the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func CloneSlice[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}

	result := make([]T, len(slice))
	copy(result, slice)

	return result
}

// CloneString returns a heap copy of an arena string, which outlives the arena.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func CloneString(s string) string {
	if len(s) == 0 {
		return ""
	}

	return string(res.UnsafeBytes(s))
}

// DeleteObject frees one arena object where the allocator supports it.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func DeleteObject[T any](a *Arena, obj *T) {
	if obj != nil {
		a.Remove(unsafe.Pointer(obj))
	}
}

// DeleteSlice frees an arena slice's backing block. It must be given the slice
// MakeSlice returned rather than a re-slice of it.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func DeleteSlice[T any](a *Arena, slice []T) {
	if len(slice) > 0 {
		a.Remove(unsafe.Pointer(&slice[0]))
	}
}

// DeleteString frees an arena string's bytes.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func DeleteString(a *Arena, s string) {
	if len(s) > 0 {
		a.Remove(unsafe.Pointer(unsafe.StringData(s)))
	}
}

// _Shape returns the allocation size and alignment for a value.
//
// A zero-sized type still gets one byte, so that two allocations of it have
// distinct addresses and an ownership test can tell them apart.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func _Shape[T any](zero T) (uint64, uint64) {
	size := unsafe.Sizeof(zero)
	if size == 0 {
		size = 1
	}

	return uint64(size), uint64(unsafe.Alignof(zero))
}
