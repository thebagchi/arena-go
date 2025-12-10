package arena

import (
	"unsafe"
)

// AsUnsafePointer converts a pointer to unsafe.Pointer.
// This is a generic helper that eliminates the need for explicit unsafe.Pointer casts.
//
// Example:
//
//	ptr := MakeObject[int](a)
//	unsafePtr := AsUnsafePointer(ptr)
func AsUnsafePointer[T any](ptr *T) unsafe.Pointer {
	return unsafe.Pointer(ptr)
}

// AsUnsafePointerSlice converts a slice to unsafe.Pointer pointing to its underlying array.
// Returns nil for empty slices.
//
// Example:
//
//	slice := MakeSlice[int](a, 10, 20)
//	slicePtr := AsUnsafePointerSlice(slice)
func AsUnsafePointerSlice[T any](slice []T) unsafe.Pointer {
	if len(slice) == 0 {
		return nil
	}
	return unsafe.Pointer(unsafe.SliceData(slice))
}

// AsUnsafePointerString converts a string to unsafe.Pointer pointing to its underlying data.
// Returns nil for empty strings.
//
// Example:
//
//	str := "hello"
//	strPtr := AsUnsafePointerString(str)
func AsUnsafePointerString(s string) unsafe.Pointer {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Pointer(unsafe.StringData(s))
}

// OwnsPtr checks if the given pointer to a value belongs to memory managed by this arena.
// This is a convenience wrapper around Owns that eliminates the need for unsafe.Pointer casts.
func OwnsPtr[T any](a *Arena, ptr *T) bool {
	return a.Allocator.Owns(unsafe.Pointer(ptr))
}

// OwnsSlice checks if the underlying array of the given slice belongs to memory managed by this arena.
// Returns false for nil or empty slices.
func OwnsSlice[T any](a *Arena, slice []T) bool {
	if len(slice) == 0 {
		return false
	}
	return a.Owns(unsafe.Pointer(unsafe.SliceData(slice)))
}

// OwnsString checks if the underlying data of the given string belongs to memory managed by this arena.
// Returns false for empty strings.
func OwnsString(a *Arena, s string) bool {
	if len(s) == 0 {
		return false
	}
	return a.Allocator.Owns(unsafe.Pointer(unsafe.StringData(s)))
}
