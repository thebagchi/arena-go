// Package alloc holds the three allocation strategies an Arena can be built on.
//
// All three satisfy the same contract, stated once on the Allocator interface in
// the arena package: Alloc returns zeroed memory or panics with
// res.ErrOutOfMemory, TryAlloc reports exhaustion instead, alignment must be a
// power of two, and any use after Delete panics with res.ErrDeleted.
//
// None of them is safe for concurrent use on its own. Wrap one in
// alloc.Synchronized when an arena is shared between goroutines; a program that
// gives each goroutine its own arena, which is the usual shape, pays nothing.
package alloc

import (
	"fmt"
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

// BumpAllocator hands out memory by advancing a cursor and reclaims it all at
// once. It is the fastest strategy here and supports no individual free.
//
// It is a thin shell over res.Bump, which holds the single lock on the path.
// The shell used to add a second mutex of its own, so every allocation paid for
// two uncontended lock pairs to protect one cursor.
type BumpAllocator struct {
	bump *res.Bump
}

// NewBumpAllocator returns a bump allocator whose first chunk holds at least
// size bytes. Later chunks double in size, so a working set much larger than
// size costs a handful of mappings rather than one per chunk's worth of bytes.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func NewBumpAllocator(size int) *BumpAllocator {
	return &BumpAllocator{bump: res.NewBump(size)}
}

// Alloc returns zeroed memory, or panics with res.ErrOutOfMemory.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Alloc(size, align uint64) unsafe.Pointer {
	return b.bump.Alloc(size, align)
}

// TryAlloc is Alloc, reporting exhaustion instead of panicking.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) TryAlloc(size, align uint64) (unsafe.Pointer, bool) {
	return b.bump.TryAlloc(size, align)
}

// Reset reclaims every allocation and zeroes what was handed out, keeping the
// mapped pages for reuse. All previously returned pointers become invalid.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Reset() {
	b.bump.Reset()
}

// Delete unmaps everything. A failed unmap means this allocator's own record of
// its pages disagrees with the kernel, which is a bug here rather than a
// condition a caller can handle, so it panics rather than returning.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Delete() {
	if err := b.bump.Delete(); err != nil {
		panic(fmt.Errorf("alloc: bump delete: %w", err))
	}
}

// Remove is a no-op: a bump allocator frees only in whole arenas.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Remove(ptr unsafe.Pointer) {
	// Empty
}

// Owns reports whether ptr came from this allocator.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Owns(ptr unsafe.Pointer) bool {
	return b.bump.Owns(ptr)
}

// Stats reports how much memory this allocator has mapped.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func (b *BumpAllocator) Stats() res.Stats {
	return res.Stats{Mapped: b.bump.Size(), Mappings: b.bump.Mappings()}
}
