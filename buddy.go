package arena

import (
	"math/bits"
	"unsafe"
)

type BuddyAllocator struct {
	chunkSize uint64
	free      [][]int
	order     int
}

func NewBuddyAllocator(chunkSize, numChunks int) *BuddyAllocator {
	if chunkSize&(chunkSize-1) != 0 {
		panic("chunkSize must be power of 2")
	}
	order := bits.Len(uint(chunkSize)) - 1
	b := &BuddyAllocator{
		chunkSize: uint64(chunkSize),
		order:     order,
		free:      make([][]int, order+1),
	}
	// dummy, no chunks added
	return b
}

func (b *BuddyAllocator) Alloc(size, align uint64) unsafe.Pointer {
	// dummy
	return nil
}

func (b *BuddyAllocator) Reset() {
	// dummy
}

func (b *BuddyAllocator) Delete() {
	// dummy
}

func (b *BuddyAllocator) Remove(ptr unsafe.Pointer) {
	// no op for buddy allocator
}
