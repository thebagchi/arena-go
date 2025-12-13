// Package arena provides high-performance, zero-GC slice implementation using arena memory.
// Slice offers appendable slices with small slice optimization (SSO) for minimal memory overhead.
//
// Features:
// • Zero heap allocations for append operations
// • Small Slice Optimization (SSO) for small slices
// • Comprehensive API with 30+ methods
// • Full iterator support (Go 1.23+)
// • Range loop compatibility
// • Sorting, searching, and manipulation operations
package arena

import (
	"iter"
	"sort"
)

// Vec[T] – the ultimate appendable slice in arena memory
// • All data allocated from arena memory
// • Append/Push never touches the Go heap
// • 30+ methods for comprehensive slice operations
//
// Core operations: AppendOne, Push, Pop, Get, Set, Insert, Remove
// Bulk operations: AppendSlice, Append, Resize, Clear, Reset
// Algorithms: Sort, SortStable, SortBy, Reverse, Contains, IndexOf
// Conversion: Clone (heap), CloneSlice (arena), ToSlice
// Iteration: All, All2, Keys, Iter (pull-based), range loops
//
// Usage:
//
// a := New(1024, BUMP) // Create arena
// defer a.Delete()
//
// // Create empty slice
// slice := NewVec[int](a)
//
// // Append elements (zero heap allocations)
// slice.AppendOne(1)
// slice.AppendOne(2)
// slice.AppendOne(3)
//
// // Append multiple elements
// slice.AppendSlice([]int{4, 5, 6})
//
// // Access elements
// fmt.Println(slice.Slice()) // [1 2 3 4 5 6]
//
// // Iterate using modern iterators (Go 1.23+)
// for v := range slice.All() {
// fmt.Println(v)
// }
//
// // Iterate with indices
// for i, v := range slice.All2() {
// fmt.Printf("index %d: %v\n", i, v)
// }
//
// // Traditional range loop
// for i, v := range slice.Slice() {
// fmt.Printf("index %d: %v\n", i, v)
// }
//
// // Pull-based iteration
// iter := slice.Iter()
// for v, ok := iter.Next(); ok; v, ok = iter.Next() {
// fmt.Println(v)
// }
type Vec[T any] struct {
	arena *Arena
	data  []T
}

const SSO_THRESHOLD = 16 // SSO for slices up to 16 elements

// Len returns current length
func (s *Vec[T]) Len() int {
	return len(s.data)
}

// Cap returns current capacity
func (s *Vec[T]) Cap() int {
	return cap(s.data)
}

// Slice returns the current slice (zero-copy)
// This provides access to the underlying data as a standard Go slice.
// The returned slice shares memory with the ArenaSlice and remains valid
// until the arena is deleted or reset.
// ⚠️ CAUTION: Storing the returned slice in a long-lived variable may cause heap escape.
func (s *Vec[T]) Slice() []T {
	return s.data
}

// AppendOne appends one element
// This operation never allocates on the heap - all data is stored in arena memory.
// Small slices (up to ssoThreshold elements) get small initial capacity.
//
// Example:
//
// slice := NewVec[int](a)
// slice.AppendOne(42)
// slice.AppendOne(24)
// fmt.Println(slice.Len()) // 2
func (s *Vec[T]) AppendOne(v T) {
	s.ensure(len(s.data) + 1)
	s.data = s.data[:len(s.data)+1]
	s.data[len(s.data)-1] = v
}

// Append adds multiple elements to the slice
// Similar to Go's built-in append function but for ArenaSlice.
// This method takes any number of elements and appends them efficiently.
//
// Example:
//
// slice := NewVec[int](a)
// slice.Append(1, 2, 3)  // append multiple elements at once
// slice.Append(4)         // append single element
// fmt.Println(slice.Slice()) // [1 2 3 4]
func (s *Vec[T]) Append(elems ...T) {
	s.AppendSlice(elems)
}

// AppendSlice appends multiple elements
// Efficiently appends a slice of elements with a single capacity check.
// Uses copy() for optimal performance.
//
// Example:
//
// slice := NewVec[string](a)
// slice.AppendSlice([]string{"hello", "world"})
// slice.AppendSlice([]string{"foo", "bar"})
// fmt.Println(slice.Slice()) // [hello world foo bar]
func (s *Vec[T]) AppendSlice(src []T) {
	if len(src) == 0 {
		return
	}
	s.ensure(len(s.data) + len(src))
	oldLen := len(s.data)
	s.data = s.data[:oldLen+len(src)]
	copy(s.data[oldLen:], src)
}

// ensure grows if needed
func (s *Vec[T]) ensure(needed int) {
	if needed <= cap(s.data) {
		return
	}

	// Determine new capacity with SSO awareness
	var capacity int
	if cap(s.data) == 0 {
		// Initial allocation - use SSO threshold for small slices
		if needed <= SSO_THRESHOLD {
			capacity = SSO_THRESHOLD
		} else {
			capacity = max(needed, 64)
		}
	} else {
		// Growth - double capacity or fit needed
		capacity = max(cap(s.data)*2, needed)
	}

	// Use MakeSlice from object.go to allocate from arena
	temp := MakeSlice[T](s.arena, len(s.data), capacity)
	copy(temp, s.data)
	s.arena.Remove(AsUnsafePointerSlice(s.data))
	s.data = temp
}

// Reset keeps capacity, clears length
// This allows reusing the allocated memory for new data without deallocation.
// The capacity remains the same, making subsequent appends more efficient.
//
// Example:
//
// slice := NewVec[int](a)
// slice.AppendSlice([]int{1, 2, 3})
// fmt.Println(slice.Len()) // 3
// slice.Reset()
// fmt.Println(slice.Len()) // 0
// fmt.Println(slice.Cap()) // still has capacity
func (s *Vec[T]) Reset() {
	s.data = s.data[:0]
}

// Clone returns a heap-allocated copy of the slice that escapes the arena.
// ⚠️ HEAP ESCAPE: This function allocates on the heap.
// The returned slice is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve slice
// data beyond the arena's lifetime.
//
// Example:
//
// arenaSlice := NewVec[int](a)
// arenaSlice.AppendSlice([]int{1, 2, 3})
//
// heapSlice := arenaSlice.Clone() // heap allocation here
// a.Delete() // arena is gone, but heapSlice is still valid
//
// fmt.Println(heapSlice) // [1 2 3]
func (s *Vec[T]) Clone() []T {
	if len(s.data) == 0 {
		return nil
	}
	result := make([]T, len(s.data))
	copy(result, s.data)
	return result
}

// NewSlice creates a new Slice from initial data
// All data is allocated from arena memory. Small slices benefit from SSO threshold.
//
// Example:
//
// a := New(1024, BUMP)
//
// // Small slice - efficient SSO allocation
// small := NewVec[int](a, 1, 2, 3)
//
// // Large slice - arena memory
// large := NewVec[int](a)
// for i := 0; i < 100; i++ {
// large.AppendOne(i)
// }
func NewVec[T any](a *Arena, initial ...T) *Vec[T] {
	as := &Vec[T]{arena: a}
	if len(initial) > 0 {
		as.AppendSlice(initial)
	} else {
		// Pre-allocate SSO capacity for empty slices
		as.data = MakeSlice[T](a, 0, SSO_THRESHOLD)
	}
	return as
}

// ─────────────────────────────────────────────────────────────────────────────
// Extended Methods — Super User-Friendly!
// ─────────────────────────────────────────────────────────────────────────────

// Push = AppendOne (alias — very common)
func (s *Vec[T]) Push(v T) {
	s.AppendOne(v)
}

// Pop removes and returns last element
func (s *Vec[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val, true
}

// Get returns element at index (safe)
func (s *Vec[T]) Get(i int) (T, bool) {
	if i < 0 || i >= len(s.data) {
		var zero T
		return zero, false
	}
	return s.data[i], true
}

// Set replaces element at index
func (s *Vec[T]) Set(i int, v T) bool {
	if i < 0 || i >= len(s.data) {
		return false
	}
	s.data[i] = v
	return true
}

// Insert at index (shifts elements)
func (s *Vec[T]) Insert(i int, v T) bool {
	if i < 0 || i > len(s.data) {
		return false
	}
	s.ensure(len(s.data) + 1)
	s.data = s.data[:len(s.data)+1]
	copy(s.data[i+1:], s.data[i:len(s.data)-1])
	s.data[i] = v
	return true
}

// Remove at index (shifts elements)
func (s *Vec[T]) Remove(i int) bool {
	if i < 0 || i >= len(s.data) {
		return false
	}
	copy(s.data[i:], s.data[i+1:])
	s.data = s.data[:len(s.data)-1]
	return true
}

// RemoveBy removes elements matching a condition with quantity control.
// The limit parameter controls maximum number of elements to remove (0 = unlimited).
// Returns the number of elements removed.
//
// Example:
//
//	slice := NewVec[int](a, 1, 2, 3, 4, 5, 5, 5)
//	removed := slice.RemoveBy(2, func(i int, v int) bool { return v == 5 })
//	// removed = 2, slice contains [1, 2, 3, 4, 5]
func (s *Vec[T]) RemoveBy(limit int, fn func(index int, v T) bool) int {
	var removed int
	for i := len(s.data) - 1; i >= 0; i-- {
		if fn(i, s.data[i]) {
			s.Remove(i)
			removed++
			if removed >= limit && limit > 0 {
				return removed
			}
		}
	}
	return removed
}

// Clear keeps capacity
func (s *Vec[T]) Clear() {
	s.data = s.data[:0]
}

// Resize to exact length (zero-fill if growing)
func (s *Vec[T]) Resize(n int) {
	if n <= len(s.data) {
		s.data = s.data[:n]
		return
	}
	s.ensure(n)
	oldLen := len(s.data)
	s.data = s.data[:n]
	for i := oldLen; i < n; i++ {
		s.data[i] = *new(T)
	}
}

// Truncate shrinks length
func (s *Vec[T]) Truncate(n int) bool {
	if n < 0 || n > len(s.data) {
		return false
	}
	s.data = s.data[:n]
	return true
}

// Reverse in place
func (s *Vec[T]) Reverse() {
	slice := s.Slice()
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Sort (for ordered types)
// ⚠️ CAUTION: The comparison function may cause closure allocations.
func (s *Vec[T]) Sort(less func(a, b T) bool) {
	slice := s.Slice()
	sort.Slice(slice, func(i, j int) bool { return less(slice[i], slice[j]) })
}

// SortStable
// ⚠️ CAUTION: The comparison function may cause closure allocations.
func (s *Vec[T]) SortStable(less func(a, b T) bool) {
	slice := s.Slice()
	sort.SliceStable(slice, func(i, j int) bool { return less(slice[i], slice[j]) })
}

// SortBy (for cmp.Ordered)
func (s *Vec[T]) SortBy(cmpFn func(a, b T) int) {
	if cmpFn == nil {
		// For basic ordered types, this will panic if T is not ordered
		// Users should provide their own comparison function
		panic("SortBy requires a comparison function for non-ordered types")
	}
	s.Sort(func(a, b T) bool { return cmpFn(a, b) < 0 })
}

// Contains
// ⚠️ CAUTION: Using any() for comparison may cause interface allocations.
func (s *Vec[T]) Contains(v T) bool {
	for _, x := range s.Slice() {
		if any(x) == any(v) {
			return true
		}
	}
	return false
}

// IndexOf finds the first occurrence of an element
// ⚠️ CAUTION: Using any() for comparison may cause interface allocations.
func (s *Vec[T]) IndexOf(v T) int {
	for i, x := range s.Slice() {
		if any(x) == any(v) {
			return i
		}
	}
	return -1
}

// LastIndexOf finds the last occurrence of an element
// Returns -1 if not found.
// ⚠️ CAUTION: Using any() for comparison may cause interface allocations.
func (s *Vec[T]) LastIndexOf(v T) int {
	for i := len(s.data) - 1; i >= 0; i-- {
		if any(s.data[i]) == any(v) {
			return i
		}
	}
	return -1
}

// CloneSlice returns a deep copy as new Slice
func (s *Vec[T]) CloneSlice() *Vec[T] {
	clone := NewVec[T](s.arena)
	clone.AppendSlice(s.Slice())
	return clone
}

// ToSlice returns as normal []T (copy to heap)
// ⚠️ HEAP ESCAPE: This function allocates on the heap.
func (s *Vec[T]) ToSlice() []T {
	dst := make([]T, len(s.data))
	copy(dst, s.data)
	return dst
}

// Keys returns an iterator over indices
func (s *Vec[T]) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range len(s.data) {
			if !yield(i) {
				return
			}
		}
	}
}

// -----------------------------
// Iterator support
// -----------------------------

// LenForRange returns length for range loops
func (s *Vec[T]) LenForRange() int {
	return len(s.data)
}

// At returns element at index for range loops
// Used internally by Go's range loop implementation.
// Zero-allocation access to elements.
func (s *Vec[T]) At(i int) T {
	return s.data[i]
}

// All returns an iterator over values (Go 1.23+ iter.Seq)
// Push-style iteration with early termination support.
//
// Example:
//
// slice := NewVec[int](a)
// slice.AppendSlice([]int{1, 2, 3, 4, 5})
//
// // Iterate all values
// for v := range slice.All() {
// fmt.Println(v)
// }
//
// // Early termination
// for v := range slice.All() {
// if v > 3 {
// break // stops iteration
// }
// fmt.Println(v)
// }
func (s *Vec[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.data {
			if !yield(v) {
				return
			}
		}
	}
}

// All2 returns an iterator over index-value pairs (Go 1.23+ iter.Seq2)
// Push-style iteration with indices and early termination.
//
// Example:
//
// slice := NewVec[string](a)
// slice.AppendSlice([]string{"apple", "banana", "cherry"})
//
// for i, fruit := range slice.All2() {
// fmt.Printf("Index %d: %s\n", i, fruit)
// }
func (s *Vec[T]) All2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, v := range s.data {
			if !yield(i, v) {
				return
			}
		}
	}
}

// SliceIter provides pull-based iteration
// Similar to channels or iterators in other languages.
type SliceIter[T any] struct {
	s     *Vec[T]
	index int
}

// Iter returns a pull-based iterator
// Use Next() to pull values one by one.
//
// Example:
//
// slice := NewVec[int](a)
// slice.AppendSlice([]int{10, 20, 30})
//
// iter := slice.Iter()
// for v, ok := iter.Next(); ok; v, ok = iter.Next() {
// fmt.Println(v) // prints 10, 20, 30
// }
func (s *Vec[T]) Iter() SliceIter[T] {
	return SliceIter[T]{s: s, index: 0}
}

// Next returns the next element and whether it exists
// Returns (zero_value, false) when iteration is complete.
func (it *SliceIter[T]) Next() (T, bool) {
	if it.index >= it.s.Len() {
		var zero T
		return zero, false
	}
	val := it.s.At(it.index)
	it.index++
	return val, true
}
