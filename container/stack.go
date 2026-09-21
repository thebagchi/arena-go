package container

import (
	arena "github.com/thebagchi/arena-go"
)

// Stack is a last-in first-out stack backed by a Vec in arena memory. It is not
// safe for concurrent use.
type Stack[T any] struct {
	data *Vec[T]
}

// NewStack returns an empty stack backed by the arena.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func NewStack[T any](a *arena.Arena) *Stack[T] {
	return &Stack[T]{data: NewVec[T](a)}
}

// Push adds an element to the top.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Push(value T) {
	s.data.AppendOne(value)
}

// Pop removes and returns the top element, reporting whether there was one.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Pop() (T, bool) {
	return s.data.Pop()
}

// Peek returns the top element without removing it.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Peek() (T, bool) {
	return s.data.Get(s.data.Len() - 1)
}

// Len returns the number of elements held.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Len() int {
	return s.data.Len()
}

// Cap returns the capacity of the backing block.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Cap() int {
	return s.data.Cap()
}

// IsEmpty reports whether the stack holds nothing.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) IsEmpty() bool {
	return s.data.Len() == 0
}

// Clear drops every element and keeps the capacity.
//
// Revisions:
//   - 2025-12-19 13:20: initial creation
func (s *Stack[T]) Clear() {
	s.data.Clear()
}
