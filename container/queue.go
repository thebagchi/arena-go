package container

import (
	arena "github.com/thebagchi/arena-go"
)

const (
	// COMPACT_THRESHOLD is the number of dequeued slots that must accumulate
	// before the queue shifts the remainder down, so that a queue used as a
	// short pipeline never pays for the shift.
	COMPACT_THRESHOLD = 16
	// COMPACT_DIVISOR compacts once the dead prefix is this fraction of the
	// backing block, here one half.
	COMPACT_DIVISOR = 2
)

// Queue is a first-in first-out queue backed by a Vec in arena memory. It is not
// safe for concurrent use.
type Queue[T any] struct {
	data *Vec[T]
	head int
}

// NewQueue returns an empty queue backed by the arena.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func NewQueue[T any](a *arena.Arena) *Queue[T] {
	return &Queue[T]{data: NewVec[T](a)}
}

// Enqueue adds an element to the back.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Enqueue(value T) {
	q.data.AppendOne(value)
}

// Dequeue removes and returns the front element, reporting whether there was
// one.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.head >= q.data.Len() {
		var zero T

		return zero, false
	}

	value := q.data.At(q.head)
	q.head = q.head + 1

	q._Compact()

	return value, true
}

// Peek returns the front element without removing it.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Peek() (T, bool) {
	return q.data.Get(q.head)
}

// Len returns the number of elements still queued.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Len() int {
	return q.data.Len() - q.head
}

// Cap returns the capacity of the backing block.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Cap() int {
	return q.data.Cap()
}

// IsEmpty reports whether the queue holds nothing.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) IsEmpty() bool {
	return q.head >= q.data.Len()
}

// Clear drops every element and keeps the capacity.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) Clear() {
	q.data.Clear()
	q.head = 0
}

// _Compact shifts the live elements down once the dequeued prefix is worth
// reclaiming, so a long-lived queue does not grow without bound.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (q *Queue[T]) _Compact() {
	if q.head <= COMPACT_THRESHOLD || q.head <= q.data.Len()/COMPACT_DIVISOR {
		return
	}

	data := q.data.Slice()
	live := copy(data, data[q.head:])

	q.data.Truncate(live)
	q.head = 0
}
