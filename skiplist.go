// skiplist.go — The Ultimate Arena-Backed Skip List
package arena

import (
	"iter"
	"math/rand"
	"sync"
	"unsafe"
)

type signedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type floatingPoint interface {
	~float32 | ~float64
}

type stringType interface {
	~string
}

type ordered interface {
	signedInteger | unsignedInteger | floatingPoint | stringType
}

const (
	DEFAULT_MAX_LEVEL   = 16
	DEFAULT_PROBABILITY = 0.5
)

func RandomLevel() int {
	level := 0
	for level < DEFAULT_MAX_LEVEL && rand.Float64() < DEFAULT_PROBABILITY {
		level++
	}
	return level
}

// SkipList is a thread-safe, ordered key-value store using skip list algorithm.
// All operations (Search, Insert, Delete, Range) are protected by RWMutex.
// Memory is allocated entirely from the arena, avoiding GC pressure.
type SkipList[K ordered, V any] struct {
	arena *Arena
	head  *node[K, V]
	level int
	lock  sync.RWMutex
}

type Pair[K ordered, V any] struct {
	Key   K
	Value V
}

type node[K ordered, V any] struct {
	key     K
	value   V
	level   int
	forward []*node[K, V]
}

func NewSkipList[K ordered, V any](a *Arena) *SkipList[K, V] {
	// Allocate head node
	head := (*node[K, V])(a.Allocator.Alloc(uint64(unsafe.Sizeof(node[K, V]{})), 16))
	head.level = DEFAULT_MAX_LEVEL
	head.forward = MakeSlice[*node[K, V]](a, DEFAULT_MAX_LEVEL+1, DEFAULT_MAX_LEVEL+1)

	return &SkipList[K, V]{
		arena: a,
		head:  head,
		level: 0,
	}
}

// Search finds a value by key
func (sl *SkipList[K, V]) Search(key K) (V, bool) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()

	x := sl.head
	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.key == key {
		return x.value, true
	}
	return *new(V), false
}

// Insert adds or updates a key-value pair
func (sl *SkipList[K, V]) Insert(key K, value V) {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	update := make([]*node[K, V], DEFAULT_MAX_LEVEL+1)
	x := sl.head

	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]
	if x != nil && x.key == key {
		x.value = value
		return
	}

	level := RandomLevel()
	if level > sl.level {
		for i := sl.level + 1; i <= level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}

	// Allocate new node
	n := (*node[K, V])(sl.arena.Allocator.Alloc(uint64(unsafe.Sizeof(node[K, V]{})), 16))
	n.key = key
	n.value = value
	n.level = level
	n.forward = MakeSlice[*node[K, V]](sl.arena, level+1, level+1)

	for i := 0; i <= level; i++ {
		n.forward[i] = update[i].forward[i]
		update[i].forward[i] = n
	}
}

// Delete removes a key-value pair
func (sl *SkipList[K, V]) Delete(key K) bool {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	update := make([]*node[K, V], DEFAULT_MAX_LEVEL+1)
	x := sl.head

	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]
	if x == nil || x.key != key {
		return false
	}

	for i := 0; i <= sl.level; i++ {
		if update[i].forward[i] != x {
			break
		}
		update[i].forward[i] = x.forward[i]
	}

	for sl.level > 0 && sl.head.forward[sl.level] == nil {
		sl.level--
	}
	return true
}

// Range iterates over all key-value pairs in sorted order
func (sl *SkipList[K, V]) Range(f func(K, V) bool) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	x := sl.head.forward[0]
	for x != nil {
		if !f(x.key, x.value) {
			return
		}
		x = x.forward[0]
	}
}

// All returns an iterator over all key-value pairs in sorted order.
// This can be used with Go 1.23+ range-over-func:
//
//	for key, val := range skiplist.All() {
//	    // process key, val
//	}
func (sl *SkipList[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		sl.lock.RLock()
		defer sl.lock.RUnlock()
		x := sl.head.forward[0]
		for x != nil {
			if !yield(x.key, x.value) {
				return
			}
			x = x.forward[0]
		}
	}
}

// Keys returns an iterator over all keys in sorted order.
// This can be used with Go 1.23+ range-over-func:
//
//	for key := range skiplist.Keys() {
//	    // process key
//	}
func (sl *SkipList[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		sl.lock.RLock()
		defer sl.lock.RUnlock()
		x := sl.head.forward[0]
		for x != nil {
			if !yield(x.key) {
				return
			}
			x = x.forward[0]
		}
	}
}

// Values returns an iterator over all values in key-sorted order.
// This can be used with Go 1.23+ range-over-func:
//
//	for val := range skiplist.Values() {
//	    // process val
//	}
func (sl *SkipList[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		sl.lock.RLock()
		defer sl.lock.RUnlock()
		x := sl.head.forward[0]
		for x != nil {
			if !yield(x.value) {
				return
			}
			x = x.forward[0]
		}
	}
}

// Len returns the number of elements in the skip list
func (sl *SkipList[K, V]) Len() int {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	count := 0
	x := sl.head.forward[0]
	for x != nil {
		count++
		x = x.forward[0]
	}
	return count
}

// Reset clears all elements from the skip list
func (sl *SkipList[K, V]) Reset() {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	for i := range sl.head.forward {
		sl.head.forward[i] = nil
	}
	sl.level = 0
}

// Contains checks if a key exists
func (sl *SkipList[K, V]) Contains(key K) bool {
	_, ok := sl.Search(key)
	return ok
}

// Min returns the minimum key-value pair
func (sl *SkipList[K, V]) Min() (K, V, bool) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	if x := sl.head.forward[0]; x != nil {
		return x.key, x.value, true
	}
	return *new(K), *new(V), false
}

// Max returns the maximum key-value pair
func (sl *SkipList[K, V]) Max() (K, V, bool) {
	sl.lock.RLock()
	defer sl.lock.RUnlock()
	x := sl.head
	for i := sl.level; i >= 0; i-- {
		for x.forward[i] != nil {
			x = x.forward[i]
		}
	}
	if x != sl.head {
		return x.key, x.value, true
	}
	return *new(K), *new(V), false
}

// Clone returns a heap-allocated standard Go map with all entries from the skip list.
// The returned map is independent of the arena lifecycle and can be safely used
// after the arena is deleted. Use this when you need to preserve skip list data
// beyond the arena's lifetime. Note: The returned map does not preserve order.
func (sl *SkipList[K, V]) Clone() map[K]V {
	sl.lock.RLock()
	defer sl.lock.RUnlock()

	count := sl.Len()
	if count == 0 {
		return nil
	}

	result := make(map[K]V, count)
	x := sl.head.forward[0]
	for x != nil {
		result[x.key] = x.value
		x = x.forward[0]
	}
	return result
}

// CloneSlice returns a heap-allocated slice of key-value pairs in sorted order.
// The returned slice is independent of the arena lifecycle and can be safely used
// after the arena is deleted. Use this when you need to preserve skip list data
// with ordering beyond the arena's lifetime.
func (sl *SkipList[K, V]) CloneSlice() []Pair[K, V] {
	sl.lock.RLock()
	defer sl.lock.RUnlock()

	count := sl.Len()
	if count == 0 {
		return nil
	}

	result := make([]Pair[K, V], 0, count)
	x := sl.head.forward[0]
	for x != nil {
		result = append(result, Pair[K, V]{Key: x.key, Value: x.value})
		x = x.forward[0]
	}
	return result
}
