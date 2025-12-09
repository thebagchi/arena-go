// Package arena provides memory allocation utilities for arena-based allocators.
// This package handles low-level memory operations using system calls for efficient
// memory management outside of Go's garbage collector.
package arena

import "syscall"

var pagesize int

func init() {
	pagesize = syscall.Getpagesize()
}

// MakePages allocates memory pages using mmap.
// It rounds up the requested size to the nearest page boundary to ensure
// proper alignment and prevent partial page allocations.
//
// Parameters:
//   - size: The minimum number of bytes to allocate. Will be rounded up to page size.
//
// Returns:
//   - []byte: A byte slice backed by the allocated memory pages.
//
// Panics:
//   - If mmap fails to allocate the requested memory.
//
// Note: The allocated memory is not managed by Go's GC and must be explicitly
// released using ReleasePages to avoid memory leaks.
func MakePages(size int) []byte {
	size = ((size + pagesize - 1) / pagesize) * pagesize
	data, err := syscall.Mmap(-1, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		panic(err)
	}
	return data
}

// ReleasePages frees memory pages allocated with MakePages.
// This function must be called to release memory allocated by MakePages,
// otherwise the memory will leak as it's not managed by Go's garbage collector.
//
// Parameters:
//   - data: The byte slice returned by MakePages. Must be the exact slice
//     returned by MakePages, not a subslice.
//
// Note: After calling ReleasePages, the data slice becomes invalid and
// should not be used. Attempting to access it may cause undefined behavior.
func ReleasePages(data []byte) {
	syscall.Munmap(data)
}
