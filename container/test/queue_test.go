package container_test

import (
	"testing"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestQueue_Basic covers queue basic.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_Basic(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Test empty queue
	if !q.IsEmpty() {
		t.Error("Expected empty queue")
	}

	if q.Len() != 0 {
		t.Errorf("Expected length 0, got %d", q.Len())
	}

	// Test enqueue
	q.Enqueue(42)

	if q.IsEmpty() {
		t.Error("Expected non-empty queue")
	}

	if q.Len() != 1 {
		t.Errorf("Expected length 1, got %d", q.Len())
	}

	// Test peek
	val, ok := q.Peek()
	if !ok {
		t.Error("Expected peek to succeed")
	}

	if val != 42 {
		t.Errorf("Expected peeked value 42, got %d", val)
	}

	if q.Len() != 1 {
		t.Errorf("Expected length still 1 after peek, got %d", q.Len())
	}

	// Test dequeue
	val, ok = q.Dequeue()
	if !ok {
		t.Error("Expected dequeue to succeed")
	}

	if val != 42 {
		t.Errorf("Expected dequeued value 42, got %d", val)
	}

	if q.Len() != 0 {
		t.Errorf("Expected length 0 after dequeue, got %d", q.Len())
	}

	if !q.IsEmpty() {
		t.Error("Expected empty queue after dequeue")
	}

	// Test dequeue from empty queue
	val, ok = q.Dequeue()
	if ok {
		t.Error("Expected dequeue from empty queue to fail")
	}

	if val != 0 {
		t.Errorf("Expected zero value, got %d", val)
	}
}

// TestQueue_FIFO covers queue fifo.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_FIFO(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Enqueue multiple values
	for i := 1; i <= 5; i++ {
		q.Enqueue(i)
	}

	if q.Len() != 5 {
		t.Errorf("Expected length 5, got %d", q.Len())
	}

	// Dequeue and verify FIFO order
	for i := 1; i <= 5; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Errorf("Expected dequeue at index %d to succeed", i)
		}

		if val != i {
			t.Errorf("Expected dequeued value %d, got %d", i, val)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Expected empty queue, got length %d", q.Len())
	}
}

// TestQueue_Peek_Empty covers queue peek empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_Peek_Empty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Peek on empty queue
	val, ok := q.Peek()
	if ok {
		t.Error("Expected peek on empty queue to fail")
	}

	if val != 0 {
		t.Errorf("Expected zero value, got %d", val)
	}
}

// TestQueue_Clear covers queue clear.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_Clear(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Add some elements
	for i := 1; i <= 5; i++ {
		q.Enqueue(i)
	}

	if q.Len() != 5 {
		t.Errorf("Expected length 5, got %d", q.Len())
	}

	// Clear
	q.Clear()

	if q.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", q.Len())
	}

	if !q.IsEmpty() {
		t.Error("Expected empty queue after clear")
	}
}

// TestQueue_Compaction covers queue compaction.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_Compaction(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Add and remove many elements to trigger compaction
	for i := 0; i < 100; i++ {
		q.Enqueue(i)
	}

	// Dequeue most of them
	for i := 0; i < 80; i++ {
		_, ok := q.Dequeue()
		if !ok {
			t.Errorf("Expected dequeue at index %d to succeed", i)
		}
	}

	// Should have 20 elements left
	if q.Len() != 20 {
		t.Errorf("Expected length 20, got %d", q.Len())
	}

	// Dequeue remaining
	for i := 80; i < 100; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Errorf("Expected dequeue at index %d to succeed", i)
		}

		if val != i {
			t.Errorf("Expected value %d, got %d", i, val)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Expected empty queue, got length %d", q.Len())
	}
}

// TestQueue_String covers queue string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_String(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[string](a)

	q.Enqueue("hello")
	q.Enqueue("world")

	val, ok := q.Dequeue()
	if !ok || val != "hello" {
		t.Errorf("Expected 'hello', got %q", val)
	}

	val, ok = q.Dequeue()
	if !ok || val != "world" {
		t.Errorf("Expected 'world', got %q", val)
	}
}

// TestQueue_Cap covers queue cap.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_Cap(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	initialCap := q.Cap()
	if initialCap == 0 {
		t.Error("Expected non-zero initial capacity")
	}

	// Add elements and check capacity grows
	for i := 0; i < 50; i++ {
		q.Enqueue(i)
	}

	if q.Cap() < 50 {
		t.Errorf("Expected capacity >= 50, got %d", q.Cap())
	}
}

// TestQueue_AlternatingOps covers queue alternating ops.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_AlternatingOps(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Alternate between enqueue and dequeue
	q.Enqueue(1)
	q.Enqueue(2)

	val, _ := q.Dequeue()
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	q.Enqueue(3)

	val, _ = q.Dequeue()
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	val, _ = q.Dequeue()
	if val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}
}

// TestQueue_PeekSequence covers queue peek sequence.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_PeekSequence(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	for i := 10; i < 15; i++ {
		q.Enqueue(i)
	}

	// Peek should always return the same element
	for j := 0; j < 5; j++ {
		val, ok := q.Peek()
		if !ok || val != 10 {
			t.Errorf("Iteration %d: Expected 10, got %d", j, val)
		}
	}

	// Dequeue should return what peek returned
	val, _ := q.Dequeue()
	if val != 10 {
		t.Errorf("Expected 10 from dequeue, got %d", val)
	}
}

// TestQueue_SingleElement covers queue single element.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_SingleElement(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	q.Enqueue(99)

	if q.Len() != 1 {
		t.Errorf("Expected length 1, got %d", q.Len())
	}

	val, ok := q.Peek()
	if !ok || val != 99 {
		t.Errorf("Expected 99 from peek, got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok || val != 99 {
		t.Errorf("Expected 99 from dequeue, got %d", val)
	}

	if !q.IsEmpty() {
		t.Error("Expected empty queue after dequeue")
	}
}

// TestQueue_StressTest covers queue stress test.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_StressTest(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	const n = 5000

	// Enqueue n elements
	for i := 0; i < n; i++ {
		q.Enqueue(i)
	}

	if q.Len() != n {
		t.Errorf("Expected length %d, got %d", n, q.Len())
	}

	// Dequeue all
	for i := 0; i < n; i++ {
		val, ok := q.Dequeue()
		if !ok {
			t.Errorf("Dequeue failed at iteration %d", i)
		}

		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	if !q.IsEmpty() {
		t.Error("Expected empty queue after all dequeues")
	}
}

// TestQueue_MultipleTypes covers queue multiple types.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_MultipleTypes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	// Test with structs
	type Person struct {
		name string
		age  int
	}

	qp := container.NewQueue[Person](a)
	qp.Enqueue(Person{"Alice", 30})
	qp.Enqueue(Person{"Bob", 25})

	p, ok := qp.Dequeue()
	if !ok || p.name != "Alice" || p.age != 30 {
		t.Errorf("Expected Alice 30, got %+v", p)
	}

	p, ok = qp.Dequeue()
	if !ok || p.name != "Bob" || p.age != 25 {
		t.Errorf("Expected Bob 25, got %+v", p)
	}
}

// TestQueue_EnqueueDequeuePattern covers queue enqueue dequeue pattern.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestQueue_EnqueueDequeuePattern(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	q := container.NewQueue[int](a)

	// Pattern: E, E, D, E, D, E, E, D, D
	q.Enqueue(1)
	q.Enqueue(2)
	v1, _ := q.Dequeue()
	q.Enqueue(3)
	v2, _ := q.Dequeue()
	q.Enqueue(4)
	q.Enqueue(5)
	v3, _ := q.Dequeue()
	v4, _ := q.Dequeue()

	if v1 != 1 || v2 != 2 || v3 != 3 || v4 != 4 {
		t.Errorf("Pattern mismatch: %d, %d, %d, %d", v1, v2, v3, v4)
	}

	peek, ok := q.Peek()
	if !ok || peek != 5 {
		t.Errorf("Expected remaining element 5, got %d", peek)
	}
}
