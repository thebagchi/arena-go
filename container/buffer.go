package container

import (
	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

const (
	// BUFFER_INITIAL_CAPACITY is what a new buffer starts with.
	BUFFER_INITIAL_CAPACITY = 32
	// BUFFER_MIN_GROWTH is the smallest capacity a grown buffer is given.
	BUFFER_MIN_GROWTH = 64
)

// Buffer builds a byte string in arena memory, like bytes.Buffer without the
// heap. It is not safe for concurrent use.
type Buffer struct {
	arena *arena.Arena
	buf   []byte
}

// NewBuffer returns an empty buffer backed by the arena.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func NewBuffer(a *arena.Arena) *Buffer {
	return &Buffer{
		arena: a,
		buf:   arena.MakeSlice[byte](a, 0, BUFFER_INITIAL_CAPACITY),
	}
}

// NewBufferString returns a buffer holding a copy of s.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func NewBufferString(a *arena.Arena, s string) *Buffer {
	capacity := max(len(s)*GROWTH_FACTOR, BUFFER_INITIAL_CAPACITY)
	buffer := &Buffer{
		arena: a,
		buf:   arena.MakeSlice[byte](a, 0, capacity),
	}

	buffer.AppendString(s)

	return buffer
}

// String returns the contents as a string sharing the buffer's memory. It stays
// valid until the buffer grows again or the arena is reset.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) String() string {
	return res.UnsafeString(b.buf)
}

// Bytes returns the contents without copying.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Len returns the number of bytes held.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Cap returns the capacity of the backing block.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Cap() int {
	return cap(b.buf)
}

// Append adds bytes to the buffer.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Append(bytes []byte) {
	if len(bytes) == 0 {
		return
	}

	b.Grow(len(bytes))
	b.buf = append(b.buf, bytes...)
}

// AppendString adds a string's bytes to the buffer.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) AppendString(s string) {
	if len(s) == 0 {
		return
	}

	b.Grow(len(s))
	b.buf = append(b.buf, s...)
}

// AppendByte adds one byte to the buffer.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) AppendByte(c byte) {
	b.Grow(1)
	b.buf = append(b.buf, c)
}

// Grow makes room for at least needed more bytes.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Grow(needed int) {
	if len(b.buf)+needed <= cap(b.buf) {
		return
	}

	var (
		wanted   = len(b.buf) + needed
		capacity = max(cap(b.buf)*GROWTH_FACTOR, wanted, BUFFER_MIN_GROWTH)
		grown    = arena.MakeSlice[byte](b.arena, len(b.buf), capacity)
	)

	copy(grown, b.buf)
	b.arena.Remove(res.SlicePtr(b.buf))
	b.buf = grown
}

// Reset empties the buffer and keeps its capacity.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
}

// CloneString returns a heap copy of the contents, which outlives the arena.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) CloneString() string {
	return string(b.buf)
}

// CloneBytes returns a heap copy of the contents, which outlives the arena.
//
// Revisions:
//   - 2025-12-12 20:18: initial creation
func (b *Buffer) CloneBytes() []byte {
	return arena.CloneSlice(b.buf)
}
