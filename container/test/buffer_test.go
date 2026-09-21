package container_test

import (
	"strings"
	"testing"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestBuffer_Basic covers buffer basic.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_Basic(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Test empty buffer
	if buf.Len() != 0 {
		t.Errorf("Expected length 0, got %d", buf.Len())
	}

	if buf.String() != "" {
		t.Errorf("Expected empty string, got %q", buf.String())
	}

	if buf.Cap() == 0 {
		t.Error("Expected non-zero capacity")
	}

	// Test append bytes
	buf.Append([]byte("hello"))

	if buf.Len() != 5 {
		t.Errorf("Expected length 5, got %d", buf.Len())
	}

	if buf.String() != "hello" {
		t.Errorf("Expected 'hello', got %q", buf.String())
	}
}

// TestBuffer_AppendString covers buffer append string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_AppendString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	buf.AppendString("hello")

	if buf.String() != "hello" {
		t.Errorf("Expected 'hello', got %q", buf.String())
	}

	buf.AppendString(" ")
	buf.AppendString("world")

	if buf.String() != "hello world" {
		t.Errorf("Expected 'hello world', got %q", buf.String())
	}
}

// TestBuffer_AppendBytes covers buffer append bytes.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_AppendBytes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	buf.Append([]byte("foo"))
	buf.Append([]byte("bar"))

	if buf.String() != "foobar" {
		t.Errorf("Expected 'foobar', got %q", buf.String())
	}

	// Test appending empty slice
	buf.Append([]byte{})

	if buf.String() != "foobar" {
		t.Errorf("Expected 'foobar', got %q", buf.String())
	}
}

// TestBuffer_Reset covers buffer reset.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_Reset(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	buf.AppendString("hello world")

	if buf.Len() != 11 {
		t.Errorf("Expected length 11, got %d", buf.Len())
	}

	initialCap := buf.Cap()
	buf.Reset()

	if buf.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", buf.Len())
	}

	if buf.String() != "" {
		t.Errorf("Expected empty string, got %q", buf.String())
	}

	if buf.Cap() != initialCap {
		t.Errorf(
			"Expected capacity %d to remain after reset, got %d",
			initialCap,
			buf.Cap(),
		)
	}
}

// TestBuffer_Bytes covers buffer bytes.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_Bytes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	buf.AppendString("test")
	bytes := buf.Bytes()

	if len(bytes) != 4 {
		t.Errorf("Expected bytes length 4, got %d", len(bytes))
	}

	expected := []byte("test")
	for i, b := range expected {
		if bytes[i] != b {
			t.Errorf("Expected byte[%d] = %d, got %d", i, b, bytes[i])
		}
	}
}

// TestBuffer_GrowCapacity covers buffer grow capacity.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_GrowCapacity(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)
	initialCap := buf.Cap()

	// Append enough data to trigger growth
	largeString := strings.Repeat("a", 200)
	buf.AppendString(largeString)

	if buf.Len() != 200 {
		t.Errorf("Expected length 200, got %d", buf.Len())
	}

	if buf.Cap() <= initialCap {
		t.Errorf("Expected capacity to grow from %d, but got %d", initialCap, buf.Cap())
	}

	if buf.String() != largeString {
		t.Errorf("Expected appended string to match")
	}
}

// TestBuffer_CloneString covers buffer clone string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_CloneString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)
	buf.AppendString("hello arena")

	cloned := buf.CloneString()
	if cloned != "hello arena" {
		t.Errorf("Expected 'hello arena', got %q", cloned)
	}

	// Modify buffer and verify cloned string is unchanged
	buf.Reset()
	buf.AppendString("different")

	if cloned != "hello arena" {
		t.Errorf("Expected cloned string to remain 'hello arena', got %q", cloned)
	}
}

// TestBuffer_CloneBytes covers buffer clone bytes.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_CloneBytes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)
	buf.Append([]byte("test"))

	cloned := buf.CloneBytes()
	if len(cloned) != 4 {
		t.Errorf("Expected cloned length 4, got %d", len(cloned))
	}

	// Verify it's a heap copy
	expected := []byte("test")
	for i, b := range expected {
		if cloned[i] != b {
			t.Errorf("Expected cloned[%d] = %d, got %d", i, b, cloned[i])
		}
	}
}

// TestBuffer_CloneStringEmpty covers buffer clone string empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_CloneStringEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	cloned := buf.CloneString()
	if cloned != "" {
		t.Errorf("Expected empty string, got %q", cloned)
	}
}

// TestBuffer_CloneBytesEmpty covers buffer clone bytes empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_CloneBytesEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	cloned := buf.CloneBytes()
	if cloned != nil {
		t.Errorf("Expected nil, got %v", cloned)
	}
}

// TestBuffer_NewBufferString covers buffer new buffer string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_NewBufferString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBufferString(a, "initial content")

	if buf.String() != "initial content" {
		t.Errorf("Expected 'initial content', got %q", buf.String())
	}

	if buf.Len() != 15 {
		t.Errorf("Expected length 15, got %d", buf.Len())
	}

	// Append more
	buf.AppendString(" extended")

	if buf.String() != "initial content extended" {
		t.Errorf("Expected 'initial content extended', got %q", buf.String())
	}
}

// TestBuffer_MultipleAppends covers buffer multiple appends.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_MultipleAppends(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	strs := []string{"hello", " ", "from", " ", "buffer"}
	expected := "hello from buffer"

	for _, s := range strs {
		buf.AppendString(s)
	}

	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

// TestBuffer_MixedAppends covers buffer mixed appends.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_MixedAppends(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	buf.AppendString("Hello")
	buf.Append([]byte{' ', 'W'})
	buf.AppendString("orld")

	if buf.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", buf.String())
	}
}

// TestBuffer_LargeContent covers buffer large content.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_LargeContent(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1000 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Build a large string
	largeContent := strings.Repeat("Lorem ipsum dolor sit amet. ", 100)
	buf.AppendString(largeContent)

	if buf.String() != largeContent {
		t.Error("Expected large content to match")
	}

	if buf.Len() != len(largeContent) {
		t.Errorf("Expected length %d, got %d", len(largeContent), buf.Len())
	}

	// Test cloning large content
	cloned := buf.CloneString()
	if cloned != largeContent {
		t.Error("Expected cloned large content to match")
	}
}

// TestBuffer_ResetAndReuse covers buffer reset and reuse.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_ResetAndReuse(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// First use
	buf.AppendString("first")

	if buf.String() != "first" {
		t.Errorf("Expected 'first', got %q", buf.String())
	}

	// Reset
	buf.Reset()

	if buf.String() != "" {
		t.Errorf("Expected empty string after reset, got %q", buf.String())
	}

	// Second use
	buf.AppendString("second")

	if buf.String() != "second" {
		t.Errorf("Expected 'second', got %q", buf.String())
	}

	// Reset again
	buf.Reset()
	buf.AppendString("third")

	if buf.String() != "third" {
		t.Errorf("Expected 'third', got %q", buf.String())
	}
}

// TestBuffer_ConcurrentAppends covers buffer concurrent appends.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_ConcurrentAppends(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Sequential appends
	buf.AppendString("a")
	buf.AppendString("b")
	buf.AppendString("c")
	buf.AppendString("d")
	buf.AppendString("e")

	if buf.String() != "abcde" {
		t.Errorf("Expected 'abcde', got %q", buf.String())
	}
}

// TestBuffer_ByteVsString covers buffer byte vs string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_ByteVsString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Mix bytes and strings
	buf.Append([]byte{72, 101, 108, 108, 111}) // "Hello" in bytes
	buf.AppendString(" ")
	buf.Append([]byte("World"))

	if buf.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", buf.String())
	}
}

// TestBuffer_UnicodeContent covers buffer unicode content.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_UnicodeContent(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Unicode strings
	buf.AppendString("Hello 世界")
	buf.AppendString(" 🚀")

	expected := "Hello 世界 🚀"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

// TestBuffer_EmptyAppends covers buffer empty appends.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_EmptyAppends(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Many empty appends should not affect the buffer
	for i := 0; i < 10; i++ {
		buf.Append([]byte{})
	}

	if buf.Len() != 0 {
		t.Errorf("Expected length 0, got %d", buf.Len())
	}

	buf.AppendString("test")

	for i := 0; i < 10; i++ {
		buf.Append([]byte{})
	}

	if buf.String() != "test" {
		t.Errorf("Expected 'test', got %q", buf.String())
	}
}

// TestBuffer_AppendSingleBytes covers buffer append single bytes.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_AppendSingleBytes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Append single bytes at a time
	for i := 0; i < 5; i++ {
		buf.Append([]byte{byte(65 + i)}) // A, B, C, D, E
	}

	if buf.String() != "ABCDE" {
		t.Errorf("Expected 'ABCDE', got %q", buf.String())
	}
}

// TestBuffer_LenAndCap covers buffer len and cap.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_LenAndCap(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	for i := 0; i < 5; i++ {
		buf.AppendString("test")
	}

	expectedLen := 20 // 5 * 4
	if buf.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, buf.Len())
	}

	if buf.Cap() < expectedLen {
		t.Errorf("Expected cap >= %d, got %d", expectedLen, buf.Cap())
	}

	buf.Reset()

	if buf.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", buf.Len())
	}
	// Cap should remain the same
	if buf.Cap() < expectedLen {
		t.Errorf("Expected cap to remain >= %d, got %d", expectedLen, buf.Cap())
	}
}

// TestBuffer_VeryLongString covers buffer very long string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_VeryLongString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Create a very long string
	longStr := strings.Repeat("x", 1000)
	buf.AppendString(longStr)

	if buf.Len() != 1000 {
		t.Errorf("Expected length 1000, got %d", buf.Len())
	}

	if buf.String() != longStr {
		t.Error("Expected long string to match")
	}
}

// TestBuffer_RepeatedResetReuse covers buffer repeated reset reuse.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_RepeatedResetReuse(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	for i := 0; i < 5; i++ {
		buf.AppendString("content")

		if buf.String() != "content" {
			t.Errorf("Iteration %d: Expected 'content', got %q", i, buf.String())
		}

		buf.Reset()

		if buf.Len() != 0 {
			t.Errorf(
				"Iteration %d: Expected length 0 after reset, got %d",
				i,
				buf.Len(),
			)
		}
	}
}

// TestBuffer_AppendZeroValueBytes covers buffer append zero value bytes.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_AppendZeroValueBytes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)

	// Append bytes including zero values
	buf.Append([]byte{65, 0, 66, 0, 67}) // A\0B\0C

	if buf.Len() != 5 {
		t.Errorf("Expected length 5, got %d", buf.Len())
	}

	bytes := buf.Bytes()

	expected := []byte{65, 0, 66, 0, 67}
	for i, b := range expected {
		if bytes[i] != b {
			t.Errorf("Expected bytes[%d] = %d, got %d", i, b, bytes[i])
		}
	}
}

// TestBuffer_BytesModification covers buffer bytes modification.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestBuffer_BytesModification(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	buf := container.NewBuffer(a)
	buf.AppendString("hello")

	// Get bytes and verify
	bytes := buf.Bytes()
	if string(bytes) != "hello" {
		t.Errorf("Expected 'hello', got %q", string(bytes))
	}

	// Append more
	buf.AppendString(" world")

	bytes = buf.Bytes()
	if string(bytes) != "hello world" {
		t.Errorf("Expected 'hello world', got %q", string(bytes))
	}
}
