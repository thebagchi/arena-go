// Package res provides the raw memory every arena allocator is built on, and the
// two ways this module divides it up: a PageTable that owns mapped pages and
// answers ownership questions, and a Bump cursor that hands out bytes from them.
//
// Memory here comes from mmap and lives outside Go's garbage collector. A Go
// pointer written into it does not keep its target alive, so anything stored in
// this memory must either contain no pointers or point back into the same arena.
// The arena package documentation states the rule in full.
package res

import (
	"fmt"
	"syscall"
)

const (
	// MMAP_NO_FILE is the descriptor for an anonymous mapping.
	MMAP_NO_FILE = -1
	// MMAP_NO_OFFSET is the file offset for an anonymous mapping.
	MMAP_NO_OFFSET = 0
	// MMAP_PROT is read/write, which is all arena memory ever needs.
	MMAP_PROT = syscall.PROT_READ | syscall.PROT_WRITE
	// MMAP_FLAGS keeps the mapping private to this process and file-free.
	MMAP_FLAGS = syscall.MAP_PRIVATE | syscall.MAP_ANONYMOUS
)

var (
	// PAGE_SIZE is the operating system page size, read once at process start.
	PAGE_SIZE = syscall.Getpagesize()
)

// MakePages maps size bytes of anonymous memory, rounded up to a page boundary.
//
// The kernel zeroes a fresh mapping, which is what lets the bump allocator skip
// zeroing memory it has never handed out. The result is invisible to Go's
// garbage collector and must be handed back with ReleasePages.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func MakePages(size int) ([]byte, error) {
	rounded := RoundUp(size, PAGE_SIZE)

	data, err := syscall.Mmap(MMAP_NO_FILE, MMAP_NO_OFFSET, rounded, MMAP_PROT, MMAP_FLAGS)
	if err != nil {
		return nil, fmt.Errorf("res: mmap %d bytes: %w", rounded, err)
	}

	return data, nil
}

// ReleasePages unmaps pages obtained from MakePages.
//
// It returns the unmap error rather than discarding it: a failed munmap leaks
// the mapping for the life of the process, which is worth reporting even though
// no caller can repair it.
//
// Revisions:
//   - 2025-12-09 22:38: initial creation
func ReleasePages(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if err := syscall.Munmap(data); err != nil {
		return fmt.Errorf("res: munmap %d bytes: %w", len(data), err)
	}

	return nil
}
