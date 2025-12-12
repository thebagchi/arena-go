package arena

import (
	"unsafe"
)

// Buffer is a string builder for arena allocators, similar to bytes.Buffer.
// All memory is allocated from the arena, never from the heap.
type Buffer struct {
	arena *Arena
	buf   []byte
}

// String returns the current string value
func (s *Buffer) String() string {
	if len(s.buf) == 0 {
		return ""
	}
	return unsafe.String(&s.buf[0], len(s.buf))
}

// Len returns current length
func (s *Buffer) Len() int {
	return len(s.buf)
}

// Cap returns current capacity
func (s *Buffer) Cap() int {
	return cap(s.buf)
}

// Append appends bytes – never touches the Go heap
func (s *Buffer) Append(bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	s.grow(len(bytes))
	s.buf = append(s.buf, bytes...)
}

// AppendString appends a string – convenience method
func (s *Buffer) AppendString(str string) {
	s.Append(unsafe.Slice(unsafe.StringData(str), len(str)))
}

// grow ensures capacity >= len + needed
func (s *Buffer) grow(needed int) {
	if len(s.buf)+needed <= cap(s.buf) {
		return
	}
	capacity := max(max(cap(s.buf)*2, len(s.buf)+needed), 64)

	buffer := MakeSlice[byte](s.arena, len(s.buf), capacity)
	copy(buffer, s.buf)

	// Remove old buffer from arena
	if len(s.buf) > 0 {
		s.arena.Allocator.Remove(unsafe.Pointer(&s.buf[0]))
	}
	s.buf = buffer
}

// Reset clears the string (keeps capacity)
func (s *Buffer) Reset() {
	s.buf = s.buf[:0]
}

// Bytes returns the inner byte slice backed by arena memory.
// Warning: Do not modify the returned slice, as it's shared with the buffer.
// The slice is only valid until the arena is deleted or reset.
func (s *Buffer) Bytes() []byte {
	return s.buf
}

// CloneString returns a heap-allocated copy of the string that escapes the arena.
// The returned string is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve string
// data beyond the arena's lifetime.
func (s *Buffer) CloneString() string {
	if len(s.buf) == 0 {
		return ""
	}
	return string(s.CloneBytes())
}

// CloneBytes returns a heap-allocated copy of the buffer content.
// The returned slice is independent of the arena and safe to use after arena deletion.
func (s *Buffer) CloneBytes() []byte {
	if len(s.buf) == 0 {
		return nil
	}
	b := make([]byte, len(s.buf))
	copy(b, s.buf)
	return b
}

// NewBuffer creates a new Buffer backed by the arena with initial 32-byte capacity
func NewBuffer(a *Arena) *Buffer {
	return &Buffer{
		arena: a,
		buf:   MakeSlice[byte](a, 0, 32),
	}
}

// NewBufferString creates a new Buffer with initial string content
func NewBufferString(a *Arena, s string) *Buffer {
	capacity := max(len(s)*2, 32)
	buf := &Buffer{
		arena: a,
		buf:   MakeSlice[byte](a, 0, capacity),
	}
	buf.AppendString(s)
	return buf
}
