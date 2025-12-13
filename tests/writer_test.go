package arena_test

import (
	"testing"

	"github.com/thebagchi/arena-go"
)

func TestWriter(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	w := arena.NewWriter(a)

	// Test Write
	n, err := w.Write([]byte("hello"))
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != 5 {
		t.Errorf("Write: expected 5 bytes written, got %d", n)
	}

	// Test WriteString
	n, err = w.WriteString(" world")
	if err != nil {
		t.Errorf("WriteString failed: %v", err)
	}
	if n != 6 {
		t.Errorf("WriteString: expected 6 bytes written, got %d", n)
	}

	// Test WriteByte
	err = w.WriteByte('!')
	if err != nil {
		t.Errorf("WriteByte failed: %v", err)
	}

	// Test Bytes
	expected := "hello world!"
	actual := string(w.Bytes())
	if actual != expected {
		t.Errorf("Bytes: expected %q, got %q", expected, actual)
	}

	// Test Len
	if w.Len() != 12 {
		t.Errorf("Len: expected 12, got %d", w.Len())
	}

	// Test Cap
	if w.Cap() < w.Len() {
		t.Errorf("Cap: capacity %d should be >= len %d", w.Cap(), w.Len())
	}

	// Test Reset
	w.Reset()
	if w.Len() != 0 {
		t.Errorf("Reset: expected len 0, got %d", w.Len())
	}
	if w.Cap() == 0 {
		t.Errorf("Reset: capacity should remain after reset")
	}

	// Test growth
	w.Reset()
	largeData := make([]byte, 1000)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	n, err = w.Write(largeData)
	if err != nil {
		t.Errorf("Write large data failed: %v", err)
	}
	if n != 1000 {
		t.Errorf("Write large data: expected 1000 bytes written, got %d", n)
	}
	if w.Len() != 1000 {
		t.Errorf("Write large data: expected len 1000, got %d", w.Len())
	}
	if len(w.Bytes()) != 1000 {
		t.Errorf("Write large data: expected bytes len 1000, got %d", len(w.Bytes()))
	}
}
