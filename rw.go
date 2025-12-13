package arena

import "io"

// Writer provides a way to write bytes to an arena-allocated buffer
// without the byte array escaping to the heap.
type Writer struct {
	arena  *Arena
	buffer []byte
	offset int
}

// NewWriter creates a new Writer with an arena-allocated buffer.
func NewWriter(a *Arena) *Writer {
	buf := MakeSlice[byte](a, 0, 32)
	buf = buf[:cap(buf)] // set len to cap to allow writing
	return &Writer{
		arena:  a,
		buffer: buf,
		offset: 0,
	}
}

// Write writes p to the buffer, growing it as needed.
// The buffer is reallocated in the arena if necessary.
func (w *Writer) Write(p []byte) (n int, err error) {
	needed := w.offset + len(p)
	if needed > cap(w.buffer) {
		w.grow(needed)
	}
	copy(w.buffer[w.offset:], p)
	w.offset = w.offset + len(p)
	return len(p), nil
}

// WriteString writes s to the buffer, growing it as needed.
func (w *Writer) WriteString(s string) (n int, err error) {
	needed := w.offset + len(s)
	if needed > cap(w.buffer) {
		w.grow(needed)
	}
	copy(w.buffer[w.offset:], s)
	w.offset = w.offset + len(s)
	return len(s), nil
}

// WriteByte writes a single byte to the buffer, growing it as needed.
func (w *Writer) WriteByte(c byte) error {
	if w.offset >= cap(w.buffer) {
		w.grow(w.offset + 1)
	}
	w.buffer[w.offset] = c
	w.offset = w.offset + 1
	return nil
}

// Bytes returns the written bytes as a slice.
// The underlying array is arena-allocated and does not escape to the heap.
func (w *Writer) Bytes() []byte {
	return w.buffer[:w.offset]
}

// Len returns the number of bytes written.
func (w *Writer) Len() int {
	return w.offset
}

// Cap returns the capacity of the buffer.
func (w *Writer) Cap() int {
	return cap(w.buffer)
}

// Reset resets the writer to be empty but retains the underlying buffer.
func (w *Writer) Reset() {
	w.offset = 0
}

// grow ensures the buffer has at least the given capacity.
func (w *Writer) grow(size int) {
	var capacity int = cap(w.buffer) * 2
	if capacity < size {
		capacity = size
	}
	if capacity < 64 {
		capacity = 64
	}
	temp := MakeSlice[byte](w.arena, 0, capacity)
	temp = temp[:cap(temp)]
	copy(temp, w.buffer[:w.offset])
	DeleteSlice(w.arena, w.buffer)
	w.buffer = temp
}

// Reader provides a way to read bytes from an arena-allocated buffer
// without the byte array escaping to the heap.
type Reader struct {
	arena  *Arena
	buffer []byte
	offset int
}

// NewReader creates a new Reader with an arena-allocated buffer.
func NewReader(a *Arena, data []byte) *Reader {
	return &Reader{
		arena:  a,
		buffer: data,
		offset: 0,
	}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.offset >= len(r.buffer) {
		return 0, io.EOF
	}
	n = copy(p, r.buffer[r.offset:])
	r.offset = r.offset + n
	return n, nil
}

// Len returns the number of bytes remaining to be read.
func (r *Reader) Len() int {
	return len(r.buffer) - r.offset
}

// Size returns the original length of the buffer.
func (r *Reader) Size() int {
	return len(r.buffer)
}

// Reset resets the reader to the beginning of the buffer.
func (r *Reader) Reset() {
	r.offset = 0
}
