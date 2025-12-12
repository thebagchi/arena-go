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
	ptr := a.Allocator.Alloc(uint64(size), 16)
	return (*T)(ptr)
}

// Ptr allocates memory for a value in the arena and returns a pointer to it.
// The value is copied into arena memory, making it independent of the original.
//
// Example:
//
//	a := New(1024, BUMP)
//	defer a.Delete()
//
//	value := 42
//	ptr := Ptr(a, value)  // allocates int in arena, returns *int
//	*ptr = 100            // modify the arena-backed value
func Ptr[T any](a *Arena, value T) *T {
	ptr := Alloc[T](a)
	*ptr = value
	return ptr
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
	ptr := a.Allocator.Alloc(uint64(size), uint64(align))
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
		ptr   = a.Allocator.Alloc(uint64(capacity)*uint64(size), 16)
		slice = unsafe.Slice((*T)(ptr), capacity)
	)
	return slice[:length]
}

// Append appends elements to an arena-backed slice, growing it if necessary.
// This function ensures that appended elements stay within arena memory and
// don't cause heap allocations. When growing is required, the old slice backing
// is automatically marked for deletion in the arena. Use this instead of the
// built-in append function when working with arena-backed slices.
//
// Parameters:
//   - a: The arena that backs the slice
//   - slice: The arena-backed slice to append to
//   - elems: Elements to append
//
// Returns:
//   - A new slice that includes the original elements plus the appended ones
//
// Example:
//
//	slice := arena.MakeSlice[int](a, 2, 4) // []int with cap 4
//	slice[0] = 1
//	slice[1] = 2
//
//	// Append more elements
//	slice = arena.Append(a, slice, 3, 4, 5)
//	fmt.Println(slice) // [1 2 3 4 5]
func Append[T any](a *Arena, slice []T, elems ...T) []T {
	if len(elems) == 0 {
		return slice
	}

	// Fast path for single element append (most common case)
	if len(elems) == 1 {
		length := len(slice) + 1
		if length <= cap(slice) {
			// Have capacity, direct assignment
			slice = slice[:length]
			slice[length-1] = elems[0]
			return slice
		}
		// Need to grow
		capacity := max(cap(slice)*2, 4)
		temp := MakeSlice[T](a, length, capacity)
		copy(temp, slice)
		temp[length-1] = elems[0]
		if len(slice) > 0 {
			a.Allocator.Remove(unsafe.Pointer(&slice[0]))
		}
		return temp
	}

	// Multi-element append
	length := len(slice) + len(elems)
	if length > cap(slice) {
		// Need to allocate new backing
		capacity := max(max(cap(slice)*2, length), 4)
		temp := MakeSlice[T](a, length, capacity)
		copy(temp[:len(slice)], slice)
		copy(temp[len(slice):], elems)
		if len(slice) > 0 {
			a.Allocator.Remove(unsafe.Pointer(&slice[0]))
		}
		return temp
	}
	// Enough capacity, just append in place
	copy(slice[len(slice):length], elems)
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
	ptr := a.Allocator.Alloc(uint64(len(s)), 1)
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
// This function should be used with allocators that support individual object deletion.
// Note that not all allocator types support individual deletions.
//
// Example:
//
//	obj := arena.MakeObject[MyStruct](a)
//	// ... use obj ...
//	arena.DeleteObject(a, obj)
func DeleteObject[T any](a *Arena, obj *T) {
	a.Allocator.Remove(unsafe.Pointer(obj))
}

// DeleteSlice marks an arena-allocated slice for deletion.
// This function should be used with allocators that support individual slice deletion.
// Note that not all allocator types support individual deletions.
//
// Example:
//
//	slice := arena.MakeSlice[int](a, 10, 20)
//	// ... use slice ...
//	arena.DeleteSlice(a, slice)
func DeleteSlice[T any](a *Arena, slice []T) {
	if len(slice) > 0 {
		a.Allocator.Remove(unsafe.Pointer(&slice[0]))
	}
}

// DeleteString marks an arena-allocated string for deletion.
// This function should be used with allocators that support individual string deletion.
// Note that not all allocator types support individual deletions.
//
// Example:
//
//	str := a.MakeString("hello world")
//	// ... use str ...
//	arena.DeleteString(a, str)
func DeleteString(a *Arena, s string) {
	if len(s) > 0 {
		a.Allocator.Remove(unsafe.Pointer(unsafe.StringData(s)))
	}
}
