package container_test

import (
	"testing"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestStack_Basic covers stack basic.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_Basic(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Test empty stack
	if !s.IsEmpty() {
		t.Error("Expected empty stack")
	}

	if s.Len() != 0 {
		t.Errorf("Expected length 0, got %d", s.Len())
	}

	// Test push
	s.Push(42)

	if s.IsEmpty() {
		t.Error("Expected non-empty stack")
	}

	if s.Len() != 1 {
		t.Errorf("Expected length 1, got %d", s.Len())
	}

	// Test peek
	val, ok := s.Peek()
	if !ok {
		t.Error("Expected peek to succeed")
	}

	if val != 42 {
		t.Errorf("Expected peeked value 42, got %d", val)
	}

	if s.Len() != 1 {
		t.Errorf("Expected length still 1 after peek, got %d", s.Len())
	}

	// Test pop
	val, ok = s.Pop()
	if !ok {
		t.Error("Expected pop to succeed")
	}

	if val != 42 {
		t.Errorf("Expected popped value 42, got %d", val)
	}

	if s.Len() != 0 {
		t.Errorf("Expected length 0 after pop, got %d", s.Len())
	}

	if !s.IsEmpty() {
		t.Error("Expected empty stack after pop")
	}

	// Test pop from empty stack
	val, ok = s.Pop()
	if ok {
		t.Error("Expected pop from empty stack to fail")
	}

	if val != 0 {
		t.Errorf("Expected zero value, got %d", val)
	}
}

// TestStack_LIFO covers stack lifo.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_LIFO(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Push multiple values
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}

	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}

	// Pop and verify LIFO order (5, 4, 3, 2, 1)
	expectedOrder := []int{5, 4, 3, 2, 1}
	for i, expected := range expectedOrder {
		val, ok := s.Pop()
		if !ok {
			t.Errorf("Expected pop at index %d to succeed", i)
		}

		if val != expected {
			t.Errorf("Expected popped value %d, got %d", expected, val)
		}
	}

	if s.Len() != 0 {
		t.Errorf("Expected empty stack, got length %d", s.Len())
	}
}

// TestStack_Peek_Empty covers stack peek empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_Peek_Empty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Peek on empty stack
	val, ok := s.Peek()
	if ok {
		t.Error("Expected peek on empty stack to fail")
	}

	if val != 0 {
		t.Errorf("Expected zero value, got %d", val)
	}
}

// TestStack_Clear covers stack clear.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_Clear(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Add some elements
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}

	if s.Len() != 5 {
		t.Errorf("Expected length 5, got %d", s.Len())
	}

	// Clear
	s.Clear()

	if s.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", s.Len())
	}

	if !s.IsEmpty() {
		t.Error("Expected empty stack after clear")
	}
}

// TestStack_MultipleTypes covers stack multiple types.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_MultipleTypes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	// Test with strings
	ss := container.NewStack[string](a)
	ss.Push("hello")
	ss.Push("world")

	val, ok := ss.Pop()
	if !ok || val != "world" {
		t.Errorf("Expected 'world', got %q", val)
	}

	val, ok = ss.Pop()
	if !ok || val != "hello" {
		t.Errorf("Expected 'hello', got %q", val)
	}

	// Test with floats
	sf := container.NewStack[float64](a)
	sf.Push(3.14)
	sf.Push(2.71)

	fval, ok := sf.Pop()
	if !ok || fval != 2.71 {
		t.Errorf("Expected 2.71, got %f", fval)
	}

	fval, ok = sf.Pop()
	if !ok || fval != 3.14 {
		t.Errorf("Expected 3.14, got %f", fval)
	}
}

// TestStack_LargeScale covers stack large scale.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_LargeScale(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Push 1000 elements
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	if s.Len() != 1000 {
		t.Errorf("Expected length 1000, got %d", s.Len())
	}

	// Pop all in reverse order
	for i := 999; i >= 0; i-- {
		val, ok := s.Pop()
		if !ok {
			t.Errorf("Expected pop at index %d to succeed", i)
		}

		if val != i {
			t.Errorf("Expected value %d, got %d", i, val)
		}
	}

	if s.Len() != 0 {
		t.Errorf("Expected empty stack, got length %d", s.Len())
	}
}

// TestStack_Cap covers stack cap.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_Cap(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	initialCap := s.Cap()
	if initialCap == 0 {
		t.Error("Expected non-zero initial capacity")
	}

	// Add elements and check capacity
	for i := 0; i < 50; i++ {
		s.Push(i)
	}

	if s.Cap() < 50 {
		t.Errorf("Expected capacity >= 50, got %d", s.Cap())
	}
}

// TestStack_PushAndPeek covers stack push and peek.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_PushAndPeek(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	for i := 1; i <= 10; i++ {
		s.Push(i)

		val, ok := s.Peek()
		if !ok {
			t.Errorf("Expected peek at iteration %d to succeed", i)
		}

		if val != i {
			t.Errorf("Expected peeked value %d, got %d", i, val)
		}
	}

	if s.Len() != 10 {
		t.Errorf("Expected length 10, got %d", s.Len())
	}
}

// TestStack_AlternatingOps covers stack alternating ops.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_AlternatingOps(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Alternate between push and pop
	s.Push(1)
	s.Push(2)

	val, _ := s.Pop()
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	s.Push(3)

	val, _ = s.Pop()
	if val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}

	val, _ = s.Pop()
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
}

// TestStack_PeekSequence covers stack peek sequence.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_PeekSequence(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	for i := 10; i < 15; i++ {
		s.Push(i)
	}

	// Peek should always return the top element
	for j := 0; j < 5; j++ {
		val, ok := s.Peek()
		if !ok || val != 14 {
			t.Errorf("Iteration %d: Expected 14, got %d", j, val)
		}
	}

	// Pop should return what peek returned
	val, _ := s.Pop()
	if val != 14 {
		t.Errorf("Expected 14 from pop, got %d", val)
	}
}

// TestStack_SingleElement covers stack single element.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_SingleElement(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	s.Push(99)

	if s.Len() != 1 {
		t.Errorf("Expected length 1, got %d", s.Len())
	}

	val, ok := s.Peek()
	if !ok || val != 99 {
		t.Errorf("Expected 99 from peek, got %d", val)
	}

	val, ok = s.Pop()
	if !ok || val != 99 {
		t.Errorf("Expected 99 from pop, got %d", val)
	}

	if !s.IsEmpty() {
		t.Error("Expected empty stack after pop")
	}
}

// TestStack_StressTest covers stack stress test.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_StressTest(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	const n = 5000

	// Push n elements
	for i := 0; i < n; i++ {
		s.Push(i)
	}

	if s.Len() != n {
		t.Errorf("Expected length %d, got %d", n, s.Len())
	}

	// Pop all in LIFO order
	for i := n - 1; i >= 0; i-- {
		val, ok := s.Pop()
		if !ok {
			t.Errorf("Pop failed at iteration %d", i)
		}

		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	if !s.IsEmpty() {
		t.Error("Expected empty stack after all pops")
	}
}

// TestStack_InterleavedPushPop covers stack interleaved push pop.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_InterleavedPushPop(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Push-pop-push-pop pattern
	for i := 0; i < 10; i++ {
		s.Push(i)

		val, ok := s.Pop()
		if !ok || val != i {
			t.Errorf("Iteration %d: Expected %d, got %d", i, i, val)
		}
	}

	if !s.IsEmpty() {
		t.Error("Expected empty stack")
	}
}

// TestStack_PushPopSequence covers stack push pop sequence.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_PushPopSequence(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	s := container.NewStack[int](a)

	// Push multiple, pop some, push more, pop all
	for i := 1; i <= 5; i++ {
		s.Push(i)
	}

	// Pop first 2
	val, _ := s.Pop()
	if val != 5 {
		t.Errorf("Expected 5, got %d", val)
	}

	val, _ = s.Pop()
	if val != 4 {
		t.Errorf("Expected 4, got %d", val)
	}

	// Push more
	for i := 6; i <= 8; i++ {
		s.Push(i)
	}

	// Verify remaining and new
	expected := []int{8, 7, 6, 3, 2, 1}
	for _, exp := range expected {
		val, ok := s.Pop()
		if !ok || val != exp {
			t.Errorf("Expected %d, got %d", exp, val)
		}
	}
}

// TestStack_StructElements covers stack struct elements.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestStack_StructElements(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	type Task struct {
		id    int
		title string
	}

	st := container.NewStack[Task](a)
	st.Push(Task{1, "Task 1"})
	st.Push(Task{2, "Task 2"})
	st.Push(Task{3, "Task 3"})

	for i := 3; i >= 1; i-- {
		task, ok := st.Pop()
		if !ok || task.id != i {
			t.Errorf("Expected id %d, got %d", i, task.id)
		}
	}
}
