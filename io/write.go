// Package io adapts arena memory to the standard library's Reader and Writer
// interfaces, so that encoders and decoders can work without touching the Go
// heap.
package io

import (
	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

const (
	// WRITER_INITIAL_CAPACITY is the block a new writer starts with.
	WRITER_INITIAL_CAPACITY = 32
	// WRITER_MIN_GROWTH is the smallest block a grown writer moves to.
	WRITER_MIN_GROWTH = 64
	// GROWTH_FACTOR is how much the block grows when the writer fills it.
	GROWTH_FACTOR = 2
)

// Writer collects bytes into arena memory.
//
// It is not safe for concurrent use, and Bytes returns the live block rather
// than a copy.
type Writer struct {
	arena  *arena.Arena
	buffer []byte
	offset int
}

// NewWriter returns a writer backed by the arena.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func NewWriter(a *arena.Arena) *Writer {
	return &Writer{
		arena:  a,
		buffer: arena.MakeSlice[byte](a, WRITER_INITIAL_CAPACITY, WRITER_INITIAL_CAPACITY),
	}
}

// Write appends p and returns how many bytes it took, which is always all of
// them.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Write(p []byte) (int, error) {
	w.Grow(len(p))
	copy(w.buffer[w.offset:], p)

	w.offset = w.offset + len(p)

	return len(p), nil
}

// WriteString appends s and returns how many bytes it took.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) WriteString(s string) (int, error) {
	w.Grow(len(s))
	copy(w.buffer[w.offset:], s)

	w.offset = w.offset + len(s)

	return len(s), nil
}

// WriteByte appends one byte.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) WriteByte(c byte) error {
	w.Grow(1)

	w.buffer[w.offset] = c
	w.offset = w.offset + 1

	return nil
}

// Grow makes room for at least needed more bytes, in the arena.
//
// The replacement block comes from the arena and the old one is handed back.
// Growth used to call make, so a writer that outgrew its first block moved to
// the Go heap and quietly stopped being arena-backed at all.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Grow(needed int) {
	wanted := w.offset + needed
	if wanted <= cap(w.buffer) {
		w.buffer = w.buffer[:cap(w.buffer)]

		return
	}

	var (
		capacity = max(cap(w.buffer)*GROWTH_FACTOR, wanted, WRITER_MIN_GROWTH)
		grown    = arena.MakeSlice[byte](w.arena, capacity, capacity)
	)

	copy(grown, w.buffer[:w.offset])
	w.arena.Remove(res.SlicePtr(w.buffer))
	w.buffer = grown
}

// Bytes returns what has been written, sharing the writer's memory.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Bytes() []byte {
	return w.buffer[:w.offset]
}

// String returns what has been written, sharing the writer's memory.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) String() string {
	return res.UnsafeString(w.buffer[:w.offset])
}

// Len returns how many bytes have been written.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Len() int {
	return w.offset
}

// Cap returns the capacity of the backing block.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Cap() int {
	return cap(w.buffer)
}

// Reset empties the writer and keeps its block.
//
// Revisions:
//   - 2025-12-14 00:33: initial creation
func (w *Writer) Reset() {
	w.offset = 0
}
