package arena

import (
	"unsafe"
)

type SlabAllocator struct {
	blockSize uintptr
}

func NewSlabAllocator(blockSize, totalBytes int) *SlabAllocator {
	if blockSize < 16 {
		blockSize = 16
	}
	blockSize = (blockSize + 15) &^ 15
	s := &SlabAllocator{blockSize: uintptr(blockSize)}
	// dummy implementation, no actual allocation
	return s
}

func (s *SlabAllocator) Alloc(size, align uint64) unsafe.Pointer {
	// dummy
	return nil
}

func (s *SlabAllocator) Reset() {
	// dummy
}

func (s *SlabAllocator) Delete() {
	// dummy
}

func (s *SlabAllocator) Remove(ptr unsafe.Pointer) {
	// no op for slab allocator
}

func (s *SlabAllocator) Owns(ptr unsafe.Pointer) bool {
	// TODO: implement when slab allocator is fully implemented
	return false
}
