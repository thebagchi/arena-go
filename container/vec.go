// Package container holds data structures whose storage lives in an arena
// rather than on the Go heap.
//
// The rule from the arena package applies to every one of them: arena memory is
// invisible to the garbage collector, so a Go pointer, string, slice, map,
// channel, func or interface stored in one of these containers is not kept alive
// by being stored there. Element types that contain no pointers are always safe,
// and so are pointers back into the same arena. Map copies string keys into the
// arena for this reason; for any other pointer-bearing type the caller copies.
//
// Map, SkipList and Pool are safe for concurrent use. Vec, Queue, Stack, Buffer
// and Str are not: their operations hand back interior pointers and slices that
// stay valid after any lock inside them would have been released, so a lock
// would promise a safety it could not deliver.
package container

import (
	"iter"
	"sort"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

const (
	// SSO_THRESHOLD is the capacity a small vector starts with, so that the
	// first handful of appends do not move the backing block.
	SSO_THRESHOLD = 16
	// LARGE_CAPACITY is the starting capacity once the first request is already
	// past the small-vector threshold.
	LARGE_CAPACITY = 64
	// GROWTH_FACTOR is how much capacity grows when a vector is full.
	GROWTH_FACTOR = 2
	// NOT_FOUND is what a search returns when there is no match.
	NOT_FOUND = -1
)

// Vec is an appendable slice whose backing array lives in an arena.
//
// It is not safe for concurrent use. Slice returns the backing array itself, so
// a caller can hold a reference that outlives any lock this type could take.
type Vec[T any] struct {
	arena *arena.Arena
	data  []T
}

// NewVec returns a vector holding the given initial elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func NewVec[T any](a *arena.Arena, initial ...T) *Vec[T] {
	v := &Vec[T]{arena: a}

	if len(initial) > 0 {
		v.AppendSlice(initial)

		return v
	}

	v.data = arena.MakeSlice[T](a, 0, SSO_THRESHOLD)

	return v
}

// Len returns the number of elements held.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Len() int {
	return len(v.data)
}

// Cap returns the capacity of the backing array.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Cap() int {
	return cap(v.data)
}

// Slice returns the backing array as a plain slice, without copying. It stays
// valid until the arena is reset or deleted.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Slice() []T {
	return v.data
}

// AppendOne adds one element.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) AppendOne(value T) {
	v.Ensure(len(v.data) + 1)
	v.data = v.data[:len(v.data)+1]
	v.data[len(v.data)-1] = value
}

// Append adds several elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Append(elems ...T) {
	v.AppendSlice(elems)
}

// AppendSlice adds every element of src.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) AppendSlice(src []T) {
	if len(src) == 0 {
		return
	}

	old := len(v.data)

	v.Ensure(old + len(src))
	v.data = v.data[:old+len(src)]
	copy(v.data[old:], src)
}

// Ensure grows the backing array so it can hold needed elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Ensure(needed int) {
	if needed <= cap(v.data) {
		return
	}

	capacity := max(cap(v.data)*GROWTH_FACTOR, needed)
	if cap(v.data) == 0 {
		capacity = SSO_THRESHOLD
		if needed > SSO_THRESHOLD {
			capacity = max(needed, LARGE_CAPACITY)
		}
	}

	grown := arena.MakeSlice[T](v.arena, len(v.data), capacity)
	copy(grown, v.data)
	v.arena.Remove(res.SlicePtr(v.data))
	v.data = grown
}

// Push adds one element to the end.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Push(value T) {
	v.AppendOne(value)
}

// Pop removes and returns the last element, reporting whether there was one.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Pop() (T, bool) {
	if len(v.data) == 0 {
		var zero T

		return zero, false
	}

	value := v.data[len(v.data)-1]
	v.data = v.data[:len(v.data)-1]

	return value, true
}

// Get returns the element at idx, reporting whether idx was in range.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Get(idx int) (T, bool) {
	if idx < 0 || idx >= len(v.data) {
		var zero T

		return zero, false
	}

	return v.data[idx], true
}

// Set replaces the element at idx, reporting whether idx was in range.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Set(idx int, value T) bool {
	if idx < 0 || idx >= len(v.data) {
		return false
	}

	v.data[idx] = value

	return true
}

// At returns the element at idx without a range check.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) At(idx int) T {
	return v.data[idx]
}

// Insert places value at idx, shifting later elements up.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Insert(idx int, value T) bool {
	if idx < 0 || idx > len(v.data) {
		return false
	}

	v.Ensure(len(v.data) + 1)
	v.data = v.data[:len(v.data)+1]
	copy(v.data[idx+1:], v.data[idx:len(v.data)-1])
	v.data[idx] = value

	return true
}

// Remove deletes the element at idx, shifting later elements down.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Remove(idx int) bool {
	if idx < 0 || idx >= len(v.data) {
		return false
	}

	copy(v.data[idx:], v.data[idx+1:])
	v.data = v.data[:len(v.data)-1]

	return true
}

// RemoveBy deletes elements the predicate accepts, at most limit of them, or all
// of them when limit is not positive. It returns how many it removed.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) RemoveBy(limit int, match func(int, T) bool) int {
	removed := 0

	for i := len(v.data) - 1; i >= 0; i-- {
		if !match(i, v.data[i]) {
			continue
		}

		v.Remove(i)

		removed = removed + 1
		if limit > 0 && removed >= limit {
			return removed
		}
	}

	return removed
}

// Clear drops every element and keeps the capacity.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Clear() {
	v.data = v.data[:0]
}

// Reset drops every element and keeps the capacity.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Reset() {
	v.data = v.data[:0]
}

// Resize sets the length, zero-filling any new elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Resize(n int) {
	if n <= len(v.data) {
		v.data = v.data[:max(n, 0)]

		return
	}

	old := len(v.data)

	v.Ensure(n)
	v.data = v.data[:n]

	var zero T

	for i := old; i < n; i++ {
		v.data[i] = zero
	}
}

// Truncate shortens the vector to n elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Truncate(n int) bool {
	if n < 0 || n > len(v.data) {
		return false
	}

	v.data = v.data[:n]

	return true
}

// Reverse reverses the elements in place.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Reverse() {
	for i, j := 0, len(v.data)-1; i < j; i, j = i+1, j-1 {
		v.data[i], v.data[j] = v.data[j], v.data[i]
	}
}

// Sort orders the elements by the given less function.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Sort(less func(T, T) bool) {
	data := v.data
	sort.Slice(data, func(i, j int) bool {
		return less(data[i], data[j])
	})
}

// SortStable orders the elements by the given less function, keeping equal
// elements in their original order.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) SortStable(less func(T, T) bool) {
	data := v.data
	sort.SliceStable(data, func(i, j int) bool {
		return less(data[i], data[j])
	})
}

// SortBy orders the elements by a three-way comparison.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) SortBy(compare func(T, T) int) {
	v.Sort(func(a, b T) bool {
		return compare(a, b) < 0
	})
}

// IndexFunc returns the first index the predicate accepts, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) IndexFunc(match func(T) bool) int {
	for i, value := range v.data {
		if match(value) {
			return i
		}
	}

	return NOT_FOUND
}

// ContainsFunc reports whether any element satisfies the predicate.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) ContainsFunc(match func(T) bool) bool {
	return v.IndexFunc(match) != NOT_FOUND
}

// Clone returns a heap copy of the elements, which outlives the arena.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Clone() []T {
	return arena.CloneSlice(v.data)
}

// ToSlice returns a heap copy of the elements, which outlives the arena.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) ToSlice() []T {
	return arena.CloneSlice(v.data)
}

// CloneSlice returns a second vector in the same arena holding the same
// elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) CloneSlice() *Vec[T] {
	clone := NewVec[T](v.arena)
	clone.AppendSlice(v.data)

	return clone
}

// Keys returns an iterator over the indices.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range v.data {
			if !yield(i) {
				return
			}
		}
	}
}

// All returns an iterator over the elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range v.data {
			if !yield(value) {
				return
			}
		}
	}
}

// All2 returns an iterator over index and element pairs.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) All2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, value := range v.data {
			if !yield(i, value) {
				return
			}
		}
	}
}

// Iter returns a pull-based iterator over the elements.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (v *Vec[T]) Iter() VecIter[T] {
	return VecIter[T]{vec: v}
}

// VecIter walks a vector one element at a time.
type VecIter[T any] struct {
	vec   *Vec[T]
	index int
}

// Next returns the next element, reporting whether there was one.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func (it *VecIter[T]) Next() (T, bool) {
	if it.index >= it.vec.Len() {
		var zero T

		return zero, false
	}

	value := it.vec.At(it.index)
	it.index = it.index + 1

	return value, true
}

// IndexOf returns the first index holding value, or NOT_FOUND.
//
// It is a function rather than a method because finding a value needs T to be
// comparable, and a Vec holds any type. The method form compared through the
// empty interface, which allocated on every element and panicked outright on an
// element type that is not comparable at all.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func IndexOf[T comparable](v *Vec[T], value T) int {
	for i, held := range v.data {
		if held == value {
			return i
		}
	}

	return NOT_FOUND
}

// LastIndexOf returns the last index holding value, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func LastIndexOf[T comparable](v *Vec[T], value T) int {
	for i := len(v.data) - 1; i >= 0; i-- {
		if v.data[i] == value {
			return i
		}
	}

	return NOT_FOUND
}

// Contains reports whether the vector holds value.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func Contains[T comparable](v *Vec[T], value T) bool {
	return IndexOf(v, value) != NOT_FOUND
}
