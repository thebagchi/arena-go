// Package arena provides high-performance, zero-GC slice implementation using arena memory.
// ArenaSlice offers appendable slices with small slice optimization (SSO) for minimal memory overhead.
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
	"unsafe"
)

// ArenaSlice[T] – the ultimate appendable slice in arena memory
// • Small slices → inline buffer (SSO-style)
// • Large slices → growable arena memory
// • Append/Push never touches the Go heap
// • 30+ methods for comprehensive slice operations
//
// Core operations: Append, Push, Pop, Get, Set, Insert, Remove
// Bulk operations: AppendSlice, Resize, Clear, Reset
// Algorithms: Sort, SortStable, SortBy, Reverse, Contains, IndexOf
// Conversion: Clone (heap), CloneSlice (arena), ToSlice
// Iteration: All, All2, Keys, Iter (pull-based), range loops
//
// Usage:
//
//	a := New(1024, BUMP) // Create arena
//	defer a.Delete()
//
//	// Create empty slice
//	slice := MakeArenaSlice[int](a)
//
//	// Append elements (zero heap allocations)
//	slice.Append(1)
//	slice.Append(2)
//	slice.Append(3)
//
//	// Append multiple elements
//	slice.AppendSlice([]int{4, 5, 6})
//
//	// Access elements
//	fmt.Println(slice.Slice()) // [1 2 3 4 5 6]
//
//	// Iterate using modern iterators (Go 1.23+)
//	for v := range slice.All() {
//		fmt.Println(v)
//	}
//
//	// Iterate with indices
//	for i, v := range slice.All2() {
//		fmt.Printf("index %d: %v\n", i, v)
//	}
//
//	// Traditional range loop
//	for i, v := range slice.Slice() {
//		fmt.Printf("index %d: %v\n", i, v)
//	}
//
//	// Pull-based iteration
//	iter := slice.Iter()
//	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
//		fmt.Println(v)
//	}
type ArenaSlice[T any] struct {
	arena    *Arena
	ptr      unsafe.Pointer
	length   int
	capacity int
	data     [16]T   // inline buffer – 16 elements of T (adjustable)
	sso      bool    // false = inline, true = arena-backed
	_        [7]byte // padding to 64 bytes (cache-line friendly)
}

// Len returns current length
func (s *ArenaSlice[T]) Len() int {
	return s.length
}

// Cap returns current capacity
func (s *ArenaSlice[T]) Cap() int {
	if !s.sso {
		return len(s.data)
	}
	return s.capacity
}

// Slice returns the current slice (zero-copy)
// This provides access to the underlying data as a standard Go slice.
// The returned slice shares memory with the ArenaSlice and remains valid
// until the arena is deleted or reset.
func (s *ArenaSlice[T]) Slice() []T {
	if !s.sso {
		return s.data[:s.length]
	}
	return unsafe.Slice((*T)(s.ptr), s.capacity)[:s.length]
}

// Append one element
// This operation never allocates on the heap - small slices use the inline buffer,
// large slices grow within arena memory.
//
// Example:
//
//	slice := MakeArenaSlice[int](a)
//	slice.Append(42)
//	slice.Append(24)
//	fmt.Println(slice.Len()) // 2
func (s *ArenaSlice[T]) Append(v T) {
	s.ensureCapacity(s.length + 1)
	if !s.sso {
		s.data[s.length] = v
	} else {
		var zero T
		elemSize := unsafe.Sizeof(zero)
		*(*T)(unsafe.Add(s.ptr, elemSize*uintptr(s.length))) = v
	}
	s.length++
}

// AppendSlice appends multiple elements
// Efficiently appends a slice of elements with a single capacity check.
// Uses copy() for optimal performance.
//
// Example:
//
//	slice := MakeArenaSlice[string](a)
//	slice.AppendSlice([]string{"hello", "world"})
//	slice.AppendSlice([]string{"foo", "bar"})
//	fmt.Println(slice.Slice()) // [hello world foo bar]
func (s *ArenaSlice[T]) AppendSlice(src []T) {
	if len(src) == 0 {
		return
	}
	s.ensureCapacity(s.length + len(src))
	if !s.sso {
		copy(s.data[s.length:], src)
	} else {
		dst := unsafe.Slice((*T)(s.ptr), s.capacity)
		copy(dst[s.length:], src)
	}
	s.length += len(src)
}

// ensureCapacity grows if needed
func (s *ArenaSlice[T]) ensureCapacity(needed int) {
	if needed <= s.Cap() {
		return
	}

	// Migrate to arena if still inline
	if !s.sso {
		s.migrateToArena(needed)
		return
	}

	// Grow arena-backed buffer
	newCap := max(max(s.capacity*2, needed), 64)

	var zero T
	elemSize := unsafe.Sizeof(zero)
	if elemSize == 0 {
		elemSize = 1
	}

	newPtr := s.arena.raw.Alloc(uint64(newCap)*uint64(elemSize), 16)
	if s.ptr != nil {
		copy(unsafe.Slice((*T)(newPtr), newCap), unsafe.Slice((*T)(s.ptr), s.capacity))
	}
	s.ptr = newPtr
	s.capacity = newCap
}

// migrateToArena moves inline data to arena
func (s *ArenaSlice[T]) migrateToArena(needed int) {
	newCap := max(max(len(s.data)*2, needed), 64)

	var zero T
	elemSize := unsafe.Sizeof(zero)
	if elemSize == 0 {
		elemSize = 1
	}

	s.ptr = s.arena.raw.Alloc(uint64(newCap)*uint64(elemSize), 16)
	copy(unsafe.Slice((*T)(s.ptr), newCap), s.data[:])
	s.capacity = newCap
	s.sso = true
}

// Reset keeps capacity, clears length
// This allows reusing the allocated memory for new data without deallocation.
// The capacity remains the same, making subsequent appends more efficient.
//
// Example:
//
//	slice := MakeArenaSlice[int](a)
//	slice.AppendSlice([]int{1, 2, 3})
//	fmt.Println(slice.Len()) // 3
//	slice.Reset()
//	fmt.Println(slice.Len()) // 0
//	fmt.Println(slice.Cap()) // still has capacity
func (s *ArenaSlice[T]) Reset() {
	s.length = 0
	// Keep ptr/cap for reuse
}

// Clone returns a heap-allocated copy of the slice that escapes the arena.
// The returned slice is independent of the arena lifecycle and can be safely
// used after the arena is deleted. Use this when you need to preserve slice
// data beyond the arena's lifetime.
//
// Example:
//
//	arenaSlice := MakeArenaSlice[int](a)
//	arenaSlice.AppendSlice([]int{1, 2, 3})
//
//	heapSlice := arenaSlice.Clone() // heap allocation here
//	a.Delete() // arena is gone, but heapSlice is still valid
//
//	fmt.Println(heapSlice) // [1 2 3]
func (s *ArenaSlice[T]) Clone() []T {
	if s.length == 0 {
		return nil
	}
	// Create a new heap-allocated slice
	result := make([]T, s.length)
	if !s.sso {
		// SSO path - copy from inline buffer
		copy(result, s.data[:s.length])
	} else {
		// Arena path - copy from arena memory
		copy(result, unsafe.Slice((*T)(s.ptr), s.length))
	}
	return result
}

// MakeArenaSlice creates a new ArenaSlice from initial data
// Automatically chooses between inline buffer (SSO) and arena allocation
// based on the initial data size.
//
// Example:
//
//	a := New(1024, BUMP)
//
//	// Small slice - uses inline buffer
//	small := MakeArenaSlice[int](a, 1, 2, 3)
//
//	// Large slice - uses arena memory
//	large := MakeArenaSlice[int](a)
//	for i := 0; i < 100; i++ {
//		large.Append(i)
//	}
func MakeArenaSlice[T any](a *Arena, initial ...T) ArenaSlice[T] {
	var as ArenaSlice[T]
	as.arena = a
	if len(initial) <= len(as.data) {
		copy(as.data[:], initial)
		as.length = len(initial)
		as.sso = false
	} else {
		as.AppendSlice(initial)
	}
	return as
}

// ─────────────────────────────────────────────────────────────────────────────
// Extended Methods — Super User-Friendly!
// ─────────────────────────────────────────────────────────────────────────────

// Push = Append (alias — very common)
func (s *ArenaSlice[T]) Push(v T) { s.Append(v) }

// Pop removes and returns last element
func (s *ArenaSlice[T]) Pop() (T, bool) {
	if s.length == 0 {
		var zero T
		return zero, false
	}
	s.length--
	return s.At(s.length), true
}

// Get returns element at index (safe)
func (s *ArenaSlice[T]) Get(i int) (T, bool) {
	if i < 0 || i >= s.length {
		var zero T
		return zero, false
	}
	return s.At(i), true
}

// Set replaces element at index
func (s *ArenaSlice[T]) Set(i int, v T) bool {
	if i < 0 || i >= s.length {
		return false
	}
	if s.sso {
		s.data[i] = v
	} else {
		unsafe.Slice((*T)(s.ptr), s.capacity)[i] = v
	}
	return true
}

// Insert at index (shifts elements)
func (s *ArenaSlice[T]) Insert(i int, v T) bool {
	if i < 0 || i > s.length {
		return false
	}
	s.ensureCapacity(s.length + 1)
	if s.sso {
		copy(s.data[i+1:], s.data[i:s.length])
		s.data[i] = v
	} else {
		slice := unsafe.Slice((*T)(s.ptr), s.capacity)
		copy(slice[i+1:], slice[i:s.length])
		slice[i] = v
	}
	s.length++
	return true
}

// Remove at index (shifts elements)
func (s *ArenaSlice[T]) Remove(i int) bool {
	if i < 0 || i >= s.length {
		return false
	}
	if s.sso {
		copy(s.data[i:], s.data[i+1:s.length])
	} else {
		slice := unsafe.Slice((*T)(s.ptr), s.capacity)
		copy(slice[i:], slice[i+1:s.length])
	}
	s.length--
	return true
}

// Clear keeps capacity
func (s *ArenaSlice[T]) Clear() { s.length = 0 }

// Resize to exact length (zero-fill if growing)
func (s *ArenaSlice[T]) Resize(n int) {
	if n <= s.length {
		s.length = n
		return
	}
	s.ensureCapacity(n)
	if s.sso {
		for i := s.length; i < n; i++ {
			s.data[i] = *new(T)
		}
	} else {
		slice := unsafe.Slice((*T)(s.ptr), s.capacity)
		for i := s.length; i < n; i++ {
			slice[i] = *new(T)
		}
	}
	s.length = n
}

// Truncate shrinks length
func (s *ArenaSlice[T]) Truncate(n int) bool {
	if n < 0 || n > s.length {
		return false
	}
	s.length = n
	return true
}

// Reverse in place
func (s *ArenaSlice[T]) Reverse() {
	slice := s.Slice()
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Sort (for ordered types)
func (s *ArenaSlice[T]) Sort(less func(a, b T) bool) {
	slice := s.Slice()
	sort.Slice(slice, func(i, j int) bool { return less(slice[i], slice[j]) })
}

// SortStable
func (s *ArenaSlice[T]) SortStable(less func(a, b T) bool) {
	slice := s.Slice()
	sort.SliceStable(slice, func(i, j int) bool { return less(slice[i], slice[j]) })
}

// SortBy (for cmp.Ordered)
func (s *ArenaSlice[T]) SortBy(cmpFn func(a, b T) int) {
	if cmpFn == nil {
		// For basic ordered types, this will panic if T is not ordered
		// Users should provide their own comparison function
		panic("SortBy requires a comparison function for non-ordered types")
	}
	s.Sort(func(a, b T) bool { return cmpFn(a, b) < 0 })
}

// Contains
func (s *ArenaSlice[T]) Contains(v T) bool {
	for _, x := range s.Slice() {
		if any(x) == any(v) {
			return true
		}
	}
	return false
}

// IndexOf
func (s *ArenaSlice[T]) IndexOf(v T) int {
	for i, x := range s.Slice() {
		if any(x) == any(v) {
			return i
		}
	}
	return -1
}

// CloneSlice returns a deep copy as new ArenaSlice
func (s *ArenaSlice[T]) CloneSlice() ArenaSlice[T] {
	clone := MakeArenaSlice[T](s.arena)
	clone.AppendSlice(s.Slice())
	return clone
}

// ToSlice returns as normal []T (copy to heap)
func (s *ArenaSlice[T]) ToSlice() []T {
	dst := make([]T, s.length)
	copy(dst, s.Slice())
	return dst
}

// Keys returns an iterator over indices
func (s *ArenaSlice[T]) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < s.length; i++ {
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
func (s *ArenaSlice[T]) LenForRange() int {
	return s.length
}

// At returns element at index for range loops
// Used internally by Go's range loop implementation.
// Zero-allocation access to elements.
func (s *ArenaSlice[T]) At(i int) T {
	if !s.sso {
		return s.data[i]
	}
	return unsafe.Slice((*T)(s.ptr), s.capacity)[i]
}

// All returns an iterator over values (Go 1.23+ iter.Seq)
// Push-style iteration with early termination support.
//
// Example:
//
//	slice := MakeArenaSlice[int](a)
//	slice.AppendSlice([]int{1, 2, 3, 4, 5})
//
//	// Iterate all values
//	for v := range slice.All() {
//		fmt.Println(v)
//	}
//
//	// Early termination
//	for v := range slice.All() {
//		if v > 3 {
//			break // stops iteration
//		}
//		fmt.Println(v)
//	}
func (s *ArenaSlice[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if !s.sso {
			for i := 0; i < s.length; i++ {
				if !yield(s.data[i]) {
					return
				}
			}
			return
		}
		slice := unsafe.Slice((*T)(s.ptr), s.capacity)[:s.length]
		for _, v := range slice {
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
//	slice := MakeArenaSlice[string](a)
//	slice.AppendSlice([]string{"apple", "banana", "cherry"})
//
//	for i, fruit := range slice.All2() {
//		fmt.Printf("Index %d: %s\n", i, fruit)
//	}
func (s *ArenaSlice[T]) All2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		if !s.sso {
			for i := 0; i < s.length; i++ {
				if !yield(i, s.data[i]) {
					return
				}
			}
			return
		}
		slice := unsafe.Slice((*T)(s.ptr), s.capacity)[:s.length]
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

// ArenaSliceIter provides pull-based iteration
// Similar to channels or iterators in other languages.
type ArenaSliceIter[T any] struct {
	s     *ArenaSlice[T]
	index int
}

// Iter returns a pull-based iterator
// Use Next() to pull values one by one.
//
// Example:
//
//	slice := MakeArenaSlice[int](a)
//	slice.AppendSlice([]int{10, 20, 30})
//
//	iter := slice.Iter()
//	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
//		fmt.Println(v) // prints 10, 20, 30
//	}
func (s *ArenaSlice[T]) Iter() ArenaSliceIter[T] {
	return ArenaSliceIter[T]{s: s, index: 0}
}

// Next returns the next element and whether it exists
// Returns (zero_value, false) when iteration is complete.
func (it *ArenaSliceIter[T]) Next() (T, bool) {
	if it.index >= it.s.Len() {
		var zero T
		return zero, false
	}
	val := it.s.At(it.index)
	it.index++
	return val, true
}
