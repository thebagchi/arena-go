package alloc

import (
	"math/bits"
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

const (
	// MIN_ORDER is the order of the smallest block, so blocks run 16, 32, 64...
	MIN_ORDER = 4
	// MIN_BLOCK_SIZE is the smallest block a buddy chunk hands out.
	MIN_BLOCK_SIZE = 1 << MIN_ORDER
	// SUBTREE_FREE is added to a node's order to encode "this whole subtree is
	// free", leaving zero free to mean "nothing in this subtree is available".
	SUBTREE_FREE = 1
	// TREE_NODES is the size of a complete binary tree with 1<<order leaves,
	// indexed from one so that a node's children are 2i and 2i+1.
	TREE_NODES = 2
	// FILL_GROWTH is how much of a level Reset fills per copy: each one doubles
	// the span already written.
	FILL_GROWTH = 2
)

// Chunk is one page managed as a buddy system: a complete binary tree whose
// leaves are MIN_BLOCK_SIZE blocks and whose internal nodes are the merges.
//
// avail[i] records the largest free block inside subtree i, as its order plus
// one, with zero meaning nothing there is free. That single number is what makes
// every operation O(log n): allocation walks down following any child that can
// still fit the request, and freeing walks up merging while both halves are
// free.
//
// The previous encoding was one bit per node, where a set bit meant "allocated
// or partly used". Those two states are indistinguishable once a subtree fills
// up, so the search could not skip a full subtree and re-walked every allocated
// block on every call, which made allocation cost grow with the number of live
// blocks. The cost of telling them apart is one byte per node instead of one
// bit: two bytes of bookkeeping per sixteen-byte block.
//
// avail lives on the Go heap rather than at the front of the page. Keeping it
// out of the payload means the payload starts page-aligned at offset zero, and
// that a write running off the end of one block cannot corrupt the allocator's
// own records.
type Chunk struct {
	page  *res.Page
	avail []uint8
	size  int
	order int
}

// NewChunk maps a chunk whose payload is size bytes rounded up to a power of
// two, and marks the whole of it free.
//
// Revisions:
//   - 2026-09-21 13:05: initial creation
func NewChunk(table *res.PageTable, size int) (*Chunk, error) {
	// Floored at a page: a chunk smaller than one would still map a whole page
	// and then refuse to hand out the rest of it.
	payload := int(res.RoundPow2(uint64(max(size, res.PAGE_SIZE))))

	page, err := table.New(payload)
	if err != nil {
		return nil, err
	}

	order := int(res.Log2(uint64(payload / MIN_BLOCK_SIZE)))
	chunk := &Chunk{
		page:  page,
		avail: make([]uint8, TREE_NODES<<uint(order)),
		size:  payload,
		order: order,
	}

	chunk.Reset()

	return chunk, nil
}

// Allocate returns a block of exactly blockSize bytes, which must be a power of
// two, or reports that the chunk cannot fit one.
//
// Revisions:
//   - 2026-09-21 13:09: initial creation
func (c *Chunk) Allocate(blockSize int) (unsafe.Pointer, bool) {
	need := _OrderOf(blockSize)
	if need > c.order || int(c.avail[1]) < need+SUBTREE_FREE {
		return nil, false
	}

	idx := 1

	for ord := c.order; ord > need; ord-- {
		idx = c._Descend(idx, need)
	}

	c.avail[idx] = 0
	c._Update(idx)

	return unsafe.Add(c.page.Ptr(), c._Offset(idx)), true
}

// Free releases the block containing ptr and merges it with its buddy for as
// long as both halves are free. It reports whether ptr named a live block.
//
// A pointer into the middle of a block frees the whole block, which is what lets
// a caller hand back a re-sliced slice.
//
// Revisions:
//   - 2026-09-21 13:13: initial creation
func (c *Chunk) Free(ptr unsafe.Pointer) bool {
	addr := uintptr(ptr)
	if !c.page.Contains(addr) {
		return false
	}

	offset := int(addr - c.page.Start())
	if offset >= c.size {
		return false
	}

	leaf := (1 << uint(c.order)) + offset/MIN_BLOCK_SIZE
	idx := 1

	for ord := c.order; ; ord-- {
		if c.avail[idx] == 0 && c._Whole(idx, ord) {
			break
		}

		if ord == 0 {
			return false
		}

		idx = leaf >> uint(ord-1)
	}

	c.avail[idx] = uint8(c._Ord(idx) + SUBTREE_FREE)
	c._Update(idx)

	return true
}

// Reset marks every block free.
//
// Each level is filled by writing its first entry and then doubling the filled
// span with copy, which is runtime.memmove and therefore vector assembly. The
// obvious byte loop is ten times slower on a chunk this size: Go rewrites a
// zero fill into memclr, but has no such rewrite for a non-zero value.
//
// Revisions:
//   - 2026-09-21 13:16: initial creation
//   - 2026-09-21 16:31: fill each level through copy rather than a byte loop,
//     which takes a 256 KiB chunk's metadata from 4.05us to 0.41us
func (c *Chunk) Reset() {
	for depth := 0; depth <= c.order; depth++ {
		var (
			value = uint8(c.order - depth + SUBTREE_FREE)
			from  = 1 << uint(depth)
			to    = min(from<<1, len(c.avail))
		)

		if from >= to {
			continue
		}

		c.avail[from] = value

		for n := 1; from+n < to; n = n * FILL_GROWTH {
			copy(c.avail[from+n:to], c.avail[from:from+n])
		}
	}
}

// Owns reports whether ptr falls inside this chunk.
//
// Revisions:
//   - 2026-09-21 13:18: initial creation
func (c *Chunk) Owns(ptr unsafe.Pointer) bool {
	return c.page.Contains(uintptr(ptr))
}

// Contains reports whether an address falls inside this chunk. It takes the
// address rather than a pointer so a caller that already holds one does not
// have to convert it back, which go vet reads as a possible misuse.
//
// Revisions:
//   - 2026-09-21 14:16: initial creation
func (c *Chunk) Contains(addr uintptr) bool {
	return c.page.Contains(addr)
}

// Size returns the payload size in bytes.
//
// Revisions:
//   - 2026-09-21 13:19: initial creation
func (c *Chunk) Size() int {
	return c.size
}

// Order returns the height of the tree, so the chunk holds 1<<Order leaves.
//
// Revisions:
//   - 2026-09-21 13:20: initial creation
func (c *Chunk) Order() int {
	return c.order
}

// Largest returns the biggest block this chunk can still hand out.
//
// Revisions:
//   - 2026-09-21 13:21: initial creation
func (c *Chunk) Largest() int {
	if c.avail[1] == 0 {
		return 0
	}

	return MIN_BLOCK_SIZE << uint(c.avail[1]-SUBTREE_FREE)
}

// Start returns the address of the chunk's first payload byte.
//
// Revisions:
//   - 2026-09-21 13:22: initial creation
func (c *Chunk) Start() uintptr {
	return c.page.Start()
}

// _Descend picks the child to continue the search in, preferring the tighter fit
// so that a large free block is not broken up while a smaller one would serve.
//
// Revisions:
//   - 2026-09-21 13:24: initial creation
func (c *Chunk) _Descend(idx, need int) int {
	var (
		left  = 2 * idx
		right = left + 1
		fits  = uint8(need + SUBTREE_FREE)
	)

	if c.avail[left] >= fits && (c.avail[right] < fits || c.avail[left] <= c.avail[right]) {
		return left
	}

	return right
}

// _Update recomputes the ancestors of idx. A parent whose halves are both
// entirely free is itself entirely free, which is where merging happens.
//
// It stops as soon as a recomputed value matches the one already stored. A
// node's value is a pure function of its two children, so if a parent did not
// change, nothing above it can have changed either.
//
// Revisions:
//   - 2026-09-21 13:27: initial creation
//   - 2026-09-21 16:33: stop once a value is unchanged instead of always walking
//     to the root, worth about 8% of an allocation
func (c *Chunk) _Update(idx int) {
	for idx > 1 {
		var (
			parent = idx / 2
			ord    = c._Ord(parent)
			left   = c.avail[2*parent]
			right  = c.avail[2*parent+1]
			want   = max(left, right)
		)

		if int(left) == ord && int(right) == ord {
			want = uint8(ord + SUBTREE_FREE)
		}

		if c.avail[parent] == want {
			return
		}

		c.avail[parent] = want
		idx = parent
	}
}

// _Whole reports whether node idx is allocated as one block rather than split.
//
// A node with nothing available is either one allocated block or two allocated
// halves. Its children tell the two apart: the children of a block allocated
// whole were never touched, so they still read entirely free.
//
// Revisions:
//   - 2026-09-21 13:31: initial creation
func (c *Chunk) _Whole(idx, ord int) bool {
	if ord == 0 {
		return true
	}

	return int(c.avail[2*idx]) == ord && int(c.avail[2*idx+1]) == ord
}

// _Ord returns the order of the block node idx stands for.
//
// Revisions:
//   - 2026-09-21 13:33: initial creation
func (c *Chunk) _Ord(idx int) int {
	return c.order - (bits.Len(uint(idx)) - 1)
}

// _Offset returns the payload offset of node idx's leftmost leaf.
//
// Revisions:
//   - 2026-09-21 13:34: initial creation
func (c *Chunk) _Offset(idx int) int {
	leaf := idx << uint(c._Ord(idx))

	return (leaf - (1 << uint(c.order))) * MIN_BLOCK_SIZE
}

// _OrderOf returns the order of a block of the given size.
//
// Revisions:
//   - 2026-09-21 13:36: initial creation
func _OrderOf(blockSize int) int {
	rounded := res.RoundPow2(uint64(max(blockSize, MIN_BLOCK_SIZE)))

	return int(res.Log2(rounded) - MIN_ORDER)
}
