package io

import (
	"io"

	arena "github.com/thebagchi/arena-go"
)

// Reader reads from a byte slice, usually one an arena owns.
//
// It is not safe for concurrent use.
type Reader struct {
	arena  *arena.Arena
	buffer []byte
	offset int
}

// NewReader returns a reader over data.
//
// Revisions:
//   - 2025-12-16 17:26: initial creation
func NewReader(a *arena.Arena, data []byte) *Reader {
	return &Reader{arena: a, buffer: data}
}

// Read copies the next bytes into p, returning io.EOF once nothing is left.
//
// Revisions:
//   - 2025-12-16 17:26: initial creation
func (r *Reader) Read(p []byte) (int, error) {
	if r.offset >= len(r.buffer) {
		return 0, io.EOF
	}

	n := copy(p, r.buffer[r.offset:])
	r.offset = r.offset + n

	return n, nil
}

// Len returns how many bytes are left to read.
//
// Revisions:
//   - 2025-12-16 17:26: initial creation
func (r *Reader) Len() int {
	return len(r.buffer) - r.offset
}

// Size returns the length of the whole buffer.
//
// Revisions:
//   - 2025-12-16 17:26: initial creation
func (r *Reader) Size() int {
	return len(r.buffer)
}

// Reset moves back to the start of the buffer.
//
// Revisions:
//   - 2025-12-16 17:26: initial creation
func (r *Reader) Reset() {
	r.offset = 0
}
