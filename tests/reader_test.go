package arena_test

import (
	"io"
	"testing"

	arena "github.com/thebagchi/arena-go"
)

func TestReader(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	data := []byte("hello world")
	reader := arena.NewReader(a, data)

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
	if err != io.EOF {
		t.Errorf("Read after EOF: expected EOF, got %v", err)
	}

	// Test Reset
	reader.Reset()
	if reader.Len() != 11 {
		t.Errorf("Reset: expected len 11, got %d", reader.Len())
	}
}
