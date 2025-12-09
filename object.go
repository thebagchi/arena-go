package arena

import (
	"unsafe"
)

// Alloc allocates and returns a pointer to a new instance of type T in the arena.
// The object is zero-initialized. This is useful for creating instances without
// heap allocation. The pointer remains valid until the arena is deleted or reset.
//
// Example:
//
//	ptr := arena.Alloc[int](a)
//	*ptr = 42
func Alloc[T any](a *Arena) *T {
	var zero T
	size := unsafe.Sizeof(zero)
	if size == 0 {
		size = 1
	}
	ptr := a.raw.Alloc(uint64(size), 16)
	return (*T)(ptr)
}

// MakeObject allocates and returns a pointer to a new instance of type T in the arena.
// The object is zero-initialized. This is useful for creating struct instances without
// heap allocation. The pointer remains valid until the arena is deleted or reset.
//
// Example:
//
//	type Node struct { Value int; Next *Node }
//	node := arena.MakeObject[Node](a)
//	node.Value = 42
func MakeObject[T any](a *Arena) *T {
	var zero T
	var (
		size  uintptr = unsafe.Sizeof(zero)
		align uintptr = unsafe.Alignof(zero)
	)
	if size == 0 {
		size = 1
	}
	ptr := a.raw.Alloc(uint64(size), uint64(align))
	return (*T)(ptr)
}

// CloneObject returns a heap-allocated copy of an arena-allocated object.
// The returned object is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve object
// data beyond the arena's lifetime.
//
// Example:
//
//	type Node struct { Value int; Next *Node }
//	arenaNode := arena.MakeObject[Node](a)
//	arenaNode.Value = 42
//	heapNode := arena.CloneObject(arenaNode)
//	a.Delete() // heapNode is still valid
func CloneObject[T any](obj *T) *T {
	if obj == nil {
		return nil
	}
	result := new(T)
	*result = *obj
	return result
}

// MakeSlice allocates and returns a slice of type T with the specified length and capacity in the arena.
// The slice elements are zero-initialized. This is useful for creating slices without
// heap allocation. The slice remains valid until the arena is deleted or reset.
//
// Example:
//
//	slice := arena.MakeSlice[int](a, 10, 20)
//	slice[0] = 42
func MakeSlice[T any](a *Arena, length, capacity int) []T {
	if capacity == 0 {
		return nil
	}
	var (
		zero T
		size uintptr = unsafe.Sizeof(zero)
	)
	if size == 0 {
		size = 1
	}
	// Check for overflow
	if uint64(capacity) > (1<<63)/uint64(size) {
		panic("arena: slice allocation size overflow")
	}
	var (
		ptr   = a.raw.Alloc(uint64(capacity)*uint64(size), 16)
		slice = unsafe.Slice((*T)(ptr), capacity)
	)
	return slice[:length]
}

// CloneSlice returns a heap-allocated copy of an arena-backed slice.
// The returned slice is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve slice
// data beyond the arena's lifetime.
func CloneSlice[T any](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

// MakeString allocates and returns a string with the specified content in the arena.
// The string is zero-copy, meaning it shares the underlying bytes with the input string.
// This is useful for creating strings without heap allocation. The string remains valid until the arena is deleted or reset.
//
// Example:
//
//	str := arena.MakeString("hello world")
//	fmt.Println(str) // prints "hello world"
func (a *Arena) MakeString(s string) string {
	if len(s) == 0 {
		return ""
	}
	ptr := a.raw.Alloc(uint64(len(s)), 1)
	copy((*[1 << 30]byte)(ptr)[:len(s):len(s)], s)
	return unsafe.String((*byte)(ptr), len(s))
}

// CloneString returns a heap-allocated copy of an arena-backed string.
// The returned string is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve string
// data beyond the arena's lifetime.
func CloneString(s string) string {
	if len(s) == 0 {
		return ""
	}
	// Force allocation on heap by creating a new string
	return string([]byte(s))
}

// DeleteObject marks an arena-allocated object for deletion.
func DeleteObject[T any](a *Arena, obj *T) {
	a.raw.Remove(unsafe.Pointer(obj))
}

// DeleteSlice marks an arena-allocated slice for deletion.
func DeleteSlice[T any](a *Arena, slice []T) {
	if len(slice) > 0 {
		a.raw.Remove(unsafe.Pointer(&slice[0]))
	}
}

// DeleteString marks an arena-allocated string for deletion.
func DeleteString(a *Arena, s string) {
	if len(s) > 0 {
		a.raw.Remove(unsafe.Pointer(unsafe.StringData(s)))
	}
}
