package container

import (
	"iter"
	"math/bits"
	"math/rand/v2"
	"sync"

	arena "github.com/thebagchi/arena-go"
)

const (
	// MAX_LEVEL is the tallest tower a node can have, which supports on the
	// order of two to the sixteenth entries before search degrades.
	MAX_LEVEL = 16
	// NODE_ALIGN is the alignment given to a node.
	NODE_ALIGN = 16
)

// Ordered is any type a skip list can order by, which is every type Go's own
// comparison operators accept.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// Pair is one key and value from a skip list.
type Pair[K Ordered, V any] struct {
	Key   K
	Value V
}

// SkipList is an ordered key-value store whose nodes live in an arena.
//
// Safe for concurrent use: reads share the lock, writes take it exclusively.
type SkipList[K Ordered, V any] struct {
	mtx   sync.RWMutex
	arena *arena.Arena
	head  *_Node[K, V]
	count int
	level int
}

// _Node is one tower in the list. forward holds one pointer per level.
type _Node[K Ordered, V any] struct {
	forward []*_Node[K, V]
	key     K
	value   V
}

// NewSkipList returns an empty skip list backed by the arena.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func NewSkipList[K Ordered, V any](a *arena.Arena) *SkipList[K, V] {
	head := arena.Alloc[_Node[K, V]](a)
	head.forward = arena.MakeSlice[*_Node[K, V]](a, MAX_LEVEL+1, MAX_LEVEL+1)

	return &SkipList[K, V]{arena: a, head: head}
}

// Search returns the value stored under key, reporting whether there was one.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Search(key K) (V, bool) {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	node := sl._Find(key)
	if node != nil && node.key == key {
		return node.value, true
	}

	var zero V

	return zero, false
}

// Contains reports whether key is present.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Contains(key K) bool {
	_, found := sl.Search(key)

	return found
}

// Insert stores value under key, replacing any previous value.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Insert(key K, value V) {
	sl.mtx.Lock()
	defer sl.mtx.Unlock()

	// A fixed-size array rather than a slice, so the search path costs no heap
	// allocation on a structure whose whole point is to avoid the heap.
	var update [MAX_LEVEL + 1]*_Node[K, V]

	node := sl.head

	for i := sl.level; i >= 0; i-- {
		for node.forward[i] != nil && node.forward[i].key < key {
			node = node.forward[i]
		}

		update[i] = node
	}

	if next := node.forward[0]; next != nil && next.key == key {
		next.value = value

		return
	}

	level := _RandomLevel()
	if level > sl.level {
		for i := sl.level + 1; i <= level; i++ {
			update[i] = sl.head
		}

		sl.level = level
	}

	fresh := arena.Alloc[_Node[K, V]](sl.arena)
	fresh.key, fresh.value = key, value
	fresh.forward = arena.MakeSlice[*_Node[K, V]](sl.arena, level+1, level+1)

	for i := 0; i <= level; i++ {
		fresh.forward[i] = update[i].forward[i]
		update[i].forward[i] = fresh
	}

	sl.count = sl.count + 1
}

// Delete removes key, reporting whether it was present.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Delete(key K) bool {
	sl.mtx.Lock()
	defer sl.mtx.Unlock()

	var update [MAX_LEVEL + 1]*_Node[K, V]

	node := sl.head

	for i := sl.level; i >= 0; i-- {
		for node.forward[i] != nil && node.forward[i].key < key {
			node = node.forward[i]
		}

		update[i] = node
	}

	target := node.forward[0]
	if target == nil || target.key != key {
		return false
	}

	for i := 0; i <= sl.level; i++ {
		if update[i].forward[i] != target {
			break
		}

		update[i].forward[i] = target.forward[i]
	}

	for sl.level > 0 && sl.head.forward[sl.level] == nil {
		sl.level = sl.level - 1
	}

	sl.count = sl.count - 1

	return true
}

// Len returns the number of entries.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Len() int {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	return sl.count
}

// Min returns the smallest key and its value.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Min() (K, V, bool) {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	if first := sl.head.forward[0]; first != nil {
		return first.key, first.value, true
	}

	var (
		zeroKey K
		zeroVal V
	)

	return zeroKey, zeroVal, false
}

// Max returns the largest key and its value.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Max() (K, V, bool) {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	node := sl.head

	for i := sl.level; i >= 0; i-- {
		for node.forward[i] != nil {
			node = node.forward[i]
		}
	}

	if node != sl.head {
		return node.key, node.value, true
	}

	var (
		zeroKey K
		zeroVal V
	)

	return zeroKey, zeroVal, false
}

// Range calls visit for each entry in key order until it returns false.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Range(visit func(K, V) bool) {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	for node := sl.head.forward[0]; node != nil; node = node.forward[0] {
		if !visit(node.key, node.value) {
			return
		}
	}
}

// All returns an iterator over key and value pairs in key order.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		sl.Range(yield)
	}
}

// Keys returns an iterator over the keys in order.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		sl.Range(func(key K, _ V) bool {
			return yield(key)
		})
	}
}

// Values returns an iterator over the values in key order.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		sl.Range(func(_ K, val V) bool {
			return yield(val)
		})
	}
}

// Reset empties the list, keeping the memory for reuse.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Reset() {
	sl.mtx.Lock()
	defer sl.mtx.Unlock()

	clear(sl.head.forward)

	sl.level, sl.count = 0, 0
}

// Clone returns a heap map of every entry, which outlives the arena. It does not
// preserve order.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) Clone() map[K]V {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	if sl.count == 0 {
		return nil
	}

	result := make(map[K]V, sl.count)
	for node := sl.head.forward[0]; node != nil; node = node.forward[0] {
		result[node.key] = node.value
	}

	return result
}

// CloneSlice returns a heap slice of every entry in key order, which outlives
// the arena.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) CloneSlice() []Pair[K, V] {
	sl.mtx.RLock()
	defer sl.mtx.RUnlock()

	if sl.count == 0 {
		return nil
	}

	result := make([]Pair[K, V], 0, sl.count)
	for node := sl.head.forward[0]; node != nil; node = node.forward[0] {
		result = append(result, Pair[K, V]{Key: node.key, Value: node.value})
	}

	return result
}

// _Find returns the first node at or after key.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func (sl *SkipList[K, V]) _Find(key K) *_Node[K, V] {
	node := sl.head

	for i := sl.level; i >= 0; i-- {
		for node.forward[i] != nil && node.forward[i].key < key {
			node = node.forward[i]
		}
	}

	return node.forward[0]
}

// _RandomLevel returns a tower height, halving in probability per level.
//
// It counts trailing ones of one random word rather than drawing a float per
// level, so a whole tower height costs one random number.
//
// Revisions:
//   - 2025-12-13 00:05: initial creation
func _RandomLevel() int {
	return min(bits.TrailingZeros64(^rand.Uint64()), MAX_LEVEL)
}
