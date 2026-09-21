package arena_test

import (
	"errors"
	"io"
	"testing"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	iopackage "github.com/thebagchi/arena-go/io"
)

// TestReader covers reader.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	data := []byte("hello world")
	reader := iopackage.NewReader(a, data)

	// Test Read
	buf := make([]byte, 5)

	n, err := reader.Read(buf)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	if n != 5 {
		t.Errorf("Read: expected 5 bytes read, got %d", n)
	}

	if string(buf) != "hello" {
		t.Errorf("Read: expected 'hello', got '%s'", string(buf))
	}

	// Test Len
	if reader.Len() != 6 {
		t.Errorf("Len: expected 6, got %d", reader.Len())
	}

	// Test Size
	if reader.Size() != 11 {
		t.Errorf("Size: expected 11, got %d", reader.Size())
	}

	// Read remaining
	buf2 := make([]byte, 10)

	n, err = reader.Read(buf2)
	if err != nil {
		t.Errorf("Read remaining failed: %v", err)
	}

	if n != 6 {
		t.Errorf("Read remaining: expected 6 bytes read, got %d", n)
	}

	if string(buf2[:n]) != " world" {
		t.Errorf("Read remaining: expected ' world', got '%s'", string(buf2[:n]))
	}

	// Test EOF
	n, err = reader.Read(buf)
	if n != 0 {
		t.Errorf("Read after EOF: expected 0 bytes, got %d", n)
	}

	if !errors.Is(err, io.EOF) {
		t.Errorf("Read after EOF: expected EOF, got %v", err)
	}

	// Test Reset
	reader.Reset()

	if reader.Len() != 11 {
		t.Errorf("Reset: expected len 11, got %d", reader.Len())
	}
}

// TestReader_EmptyData covers reader empty data.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_EmptyData(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	reader := iopackage.NewReader(a, []byte{})

	// Test reading from empty data
	buf := make([]byte, 5)

	n, err := reader.Read(buf)
	if n != 0 || !errors.Is(err, io.EOF) {
		t.Errorf("Read empty: expected 0 bytes and EOF, got %d bytes and %v", n, err)
	}

	if reader.Len() != 0 {
		t.Errorf("Empty Len: expected 0, got %d", reader.Len())
	}

	if reader.Size() != 0 {
		t.Errorf("Empty Size: expected 0, got %d", reader.Size())
	}
}

// TestReader_PartialReads covers reader partial reads.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_PartialReads(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	data := []byte("0123456789")
	reader := iopackage.NewReader(a, data)

	// Read 2 bytes at a time
	expected := []string{"01", "23", "45", "67", "89"}
	for i, exp := range expected {
		buf := make([]byte, 2)

		n, err := reader.Read(buf)
		if n != 2 || string(buf) != exp {
			t.Errorf(
				"Partial read %d: expected %q got %q (n=%d)",
				i,
				exp,
				string(buf),
				n,
			)
		}

		if err != nil && !errors.Is(err, io.EOF) {
			t.Errorf("Partial read %d: unexpected error %v", i, err)
		}
	}

	// Next read should return EOF
	buf := make([]byte, 5)

	n, err := reader.Read(buf)
	if n != 0 || !errors.Is(err, io.EOF) {
		t.Errorf("Expected EOF, got %d bytes and %v", n, err)
	}
}

// TestReader_LargeBuffer covers reader large buffer.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_LargeBuffer(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	// Create large data
	largeData := make([]byte, 5000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	reader := iopackage.NewReader(a, largeData)

	// Test Size and Len
	if reader.Size() != 5000 {
		t.Errorf("Size: expected 5000, got %d", reader.Size())
	}

	if reader.Len() != 5000 {
		t.Errorf("Len: expected 5000, got %d", reader.Len())
	}

	// Read in chunks
	buf := make([]byte, 1000)
	for i := 0; i < 5; i++ {
		n, err := reader.Read(buf)
		if n != 1000 {
			t.Errorf("Chunk %d: expected 1000 bytes, got %d", i, n)
		}

		if err != nil {
			t.Errorf("Chunk %d: unexpected error %v", i, err)
		}

		if reader.Len() != 5000-(i+1)*1000 {
			t.Errorf(
				"Chunk %d: Len expected %d, got %d",
				i,
				5000-(i+1)*1000,
				reader.Len(),
			)
		}
	}
}

// TestReader_ReadMoreThanAvailable covers reader read more than available.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_ReadMoreThanAvailable(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	data := []byte("hello")
	reader := iopackage.NewReader(a, data)

	// Request more bytes than available
	buf := make([]byte, 100)

	n, err := reader.Read(buf)
	if n != 5 {
		t.Errorf("Expected 5 bytes, got %d", n)
	}

	if string(buf[:n]) != "hello" {
		t.Errorf("Expected 'hello', got %q", string(buf[:n]))
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestReader_MultipleResets covers reader multiple resets.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_MultipleResets(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	data := []byte("test")
	reader := iopackage.NewReader(a, data)

	// Read, reset, read again
	buf := make([]byte, 4)
	_MustRead(t, reader, buf)

	if reader.Len() != 0 {
		t.Errorf("After read: Len expected 0, got %d", reader.Len())
	}

	reader.Reset()

	if reader.Len() != 4 {
		t.Errorf("After reset: Len expected 4, got %d", reader.Len())
	}

	n, _ := reader.Read(buf)
	if n != 4 || string(buf) != "test" {
		t.Errorf(
			"Second read: expected 'test' with 4 bytes, got %q with %d",
			string(buf),
			n,
		)
	}
}

// TestReader_SequentialPartialReads covers reader sequential partial reads.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestReader_SequentialPartialReads(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	defer a.Delete()

	data := []byte("abcdefghij")
	reader := iopackage.NewReader(a, data)

	// Read 3 bytes
	buf := make([]byte, 3)

	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("First read: %v", err)
	}

	if n != len(buf) || string(buf) != "abc" || reader.Len() != 7 {
		t.Errorf(
			"First read: expected 'abc' with 7 remaining, got %q with %d",
			string(buf),
			reader.Len(),
		)
	}

	// Read 2 bytes
	buf2 := make([]byte, 2)

	n, err = reader.Read(buf2)
	if err != nil {
		t.Fatalf("Second read: %v", err)
	}

	if n != len(buf2) || string(buf2) != "de" || reader.Len() != 5 {
		t.Errorf(
			"Second read: expected 'de' with 5 remaining, got %q with %d",
			string(buf2),
			reader.Len(),
		)
	}

	// Read remaining
	buf3 := make([]byte, 100)

	n, _ = reader.Read(buf3)
	if string(buf3[:n]) != "fghij" || reader.Len() != 0 {
		t.Errorf(
			"Final read: expected 'fghij' with 0 remaining, got %q with %d",
			string(buf3[:n]),
			reader.Len(),
		)
	}
}

// _MustRead reads into p and fails the test if the reader reports an error.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func _MustRead(t *testing.T, r *iopackage.Reader, p []byte) int {
	t.Helper()

	n, err := r.Read(p)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	return n
}
