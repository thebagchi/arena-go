package arena

import (
	"hash/maphash"
	"iter"
	"sync"
	"unsafe"
)

const INITIAL_BUCKET_COUNT = 16 // Initial number of buckets in the hash map

// Map is a high-performance, zero-GC hash map that lives entirely in arena memory.
// Uses separate chaining for collision resolution, eliminating clustering issues.
// Thread-safe: All operations (Get, Set, Delete, Range) are protected by an RWMutex.
// Multiple goroutines can safely call Get concurrently, while Set/Delete operations are serialized.
type Map[K comparable, V any] struct {
	mu      sync.RWMutex
	arena   *Arena
	buckets *Vec[*entry[K, V]] // arena-backed bucket array (array of pointers to chain heads)
	count   int
	cap     int
	mask    uint64
	seed    maphash.Seed
}

// entry is a node in the hash chain (linked list)
type entry[K comparable, V any] struct {
	hash uint64
	key  K
	val  V
	next *entry[K, V]
}

// NewMap creates a new Map with separate chaining for collision resolution
func NewMap[K comparable, V any](a *Arena) *Map[K, V] {
	// Create arena-backed vec for buckets
	buckets := NewVec[*entry[K, V]](a)

	// Initialize with nil pointers
	for i := 0; i < INITIAL_BUCKET_COUNT; i++ {
		buckets.AppendOne(nil)
	}

	m := &Map[K, V]{
		arena:   a,
		buckets: buckets,
		cap:     INITIAL_BUCKET_COUNT,
		mask:    uint64(INITIAL_BUCKET_COUNT - 1),
		seed:    maphash.MakeSeed(),
	}
	return m
}

// hash function using maphash for better performance and security
func (m *Map[K, V]) hash(key K) uint64 {
	var h maphash.Hash
	h.SetSeed(m.seed)

	// Write key data to hasher
	switch v := any(key).(type) {
	case string:
		h.WriteString(v)
	case int:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case int8:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case int16:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case int32:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case int64:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uint:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uint8:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uint16:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uint32:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uint64:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	case uintptr:
		writeBytes(&h, unsafe.Pointer(&v), unsafe.Sizeof(v))
	default:
		// For other comparable types, use their memory representation
		writeBytes(&h, unsafe.Pointer(&key), unsafe.Sizeof(key))
	}

	return h.Sum64()
}

// writeBytes writes raw bytes to the hasher
func writeBytes(h *maphash.Hash, ptr unsafe.Pointer, size uintptr) {
	data := unsafe.Slice((*byte)(ptr), size)
	h.Write(data)
}

// Set inserts or updates a key-value pair using separate chaining
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Grow when load factor > 0.75
	if m.count > m.cap*3/4 {
		m.grow()
	}

	hash := m.hash(key)
	index := hash & m.mask
	head, ok := m.buckets.Get(int(index))
	if !ok {
		panic("arena map: bucket index out of bounds")
	}

	// Check if key exists in chain and update
	e := head
	for e != nil {
		if e.hash == hash && e.key == key {
			e.val = value
			return
		}
		e = e.next
	}

	// Key not found, allocate new entry and prepend to chain
	// Note: entries are freed immediately on Delete/Reset via arena.Remove()
	item := (*entry[K, V])(m.arena.Alloc(uint64(unsafe.Sizeof(entry[K, V]{})), 8))

	*item = entry[K, V]{
		hash: hash,
		key:  key,
		val:  value,
		next: head,
	}

	m.buckets.Set(int(index), item)
	m.count++
}

// Get returns value and true if found
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cap == 0 {
		var zero V
		return zero, false
	}

	hash := m.hash(key)
	index := hash & m.mask
	e, ok := m.buckets.Get(int(index))
	if !ok {
		panic("arena map: bucket index out of bounds")
	}

	// Walk the chain
	for e != nil {
		if e.hash == hash && e.key == key {
			return e.val, true
		}
		e = e.next
	}

	var zero V
	return zero, false
}

// Delete removes a key from the chain and frees the entry memory
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cap == 0 {
		return
	}

	hash := m.hash(key)
	index := hash & m.mask

	// Walk the chain and remove the matching entry
	var prev *entry[K, V]
	curr, ok := m.buckets.Get(int(index))
	if !ok {
		panic("arena map: bucket index out of bounds")
	}

	for curr != nil {
		if curr.hash == hash && curr.key == key {
			// Found it - unlink from chain
			if prev == nil {
				// Removing head of chain
				m.buckets.Set(int(index), curr.next)
			} else {
				// Removing from middle/end of chain
				prev.next = curr.next
			}
			// Free the entry memory via arena
			m.arena.Remove(unsafe.Pointer(curr))
			m.count--
			return
		}
		prev = curr
		curr = curr.next
	}
}

// Range calls f for each entry in all chains
func (m *Map[K, V]) Range(f func(K, V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for i := 0; i < m.cap; i++ {
		e, ok := m.buckets.Get(i)
		if !ok {
			panic("arena map: bucket index out of bounds")
		}
		// Walk the chain at this bucket
		for e != nil {
			if !f(e.key, e.val) {
				return
			}
			e = e.next
		}
	}
}

// Len returns number of entries
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.count
}

// grow doubles the bucket array and rehashes all entries
func (m *Map[K, V]) grow() {
	obkt := m.buckets.Slice()
	ocap := m.cap

	ncap := ocap * 2
	if ncap < INITIAL_BUCKET_COUNT {
		ncap = INITIAL_BUCKET_COUNT
	}

	// Allocate new bucket array using Vec
	nbkt := NewVec[*entry[K, V]](m.arena)

	// Initialize with nil pointers
	for i := 0; i < ncap; i++ {
		nbkt.AppendOne(nil)
	}

	// Update map metadata
	m.buckets = nbkt
	m.cap = ncap
	m.mask = uint64(ncap - 1)
	ocount := m.count
	m.count = 0

	// Rehash all entries from old chains
	for i := 0; i < ocap; i++ {
		e := obkt[i]
		// Walk each chain
		for e != nil {
			next := e.next // Save next before we modify e.next

			// Reinsert entry into new bucket array
			index := e.hash & m.mask
			head, ok := nbkt.Get(int(index))
			if !ok {
				panic("arena map: bucket index out of bounds during grow")
			}
			e.next = head
			nbkt.Set(int(index), e)
			m.count++

			e = next
		}
	}

	// Sanity check
	if m.count != ocount {
		panic("arena map: lost entries during grow")
	}
}

// Reset frees all entries and clears the map while keeping capacity
func (m *Map[K, V]) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Free all entry nodes
	for i := 0; i < m.cap; i++ {
		e, ok := m.buckets.Get(i)
		if !ok {
			panic("arena map: bucket index out of bounds")
		}
		for e != nil {
			next := e.next
			m.arena.Remove(unsafe.Pointer(e))
			e = next
		}
		m.buckets.Set(i, nil)
	}
	m.count = 0
}

// Clone returns a heap-allocated standard Go map with all entries from the Map.
// The returned map is independent of the arena lifecycle and can be safely used
// after the arena is deleted. Use this when you need to preserve map data beyond
// the arena's lifetime.
func (m *Map[K, V]) Clone() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.count == 0 {
		return nil
	}

	result := make(map[K]V, m.count)
	for i := 0; i < m.cap; i++ {
		e, ok := m.buckets.Get(i)
		if !ok {
			panic("arena map: bucket index out of bounds")
		}
		// Walk the chain
		for e != nil {
			result[e.key] = e.val
			e = e.next
		}
	}
	return result
}

// -----------------------------
// Iterator support
// -----------------------------

// Keys returns an iterator over all keys in the map
// Example:
//
//	m := arena.NewMap[string, int](a)
//	m.Set("a", 1)
//	m.Set("b", 2)
//	for key := range m.Keys() {
//	    fmt.Println(key)
//	}
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for i := 0; i < m.cap; i++ {
			e, ok := m.buckets.Get(i)
			if !ok {
				panic("arena map: bucket index out of bounds")
			}
			for e != nil {
				if !yield(e.key) {
					return
				}
				e = e.next
			}
		}
	}
}

// Values returns an iterator over all values in the map
// Example:
//
//	m := arena.NewMap[string, int](a)
//	m.Set("a", 1)
//	m.Set("b", 2)
//	for val := range m.Values() {
//	    fmt.Println(val)
//	}
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for i := 0; i < m.cap; i++ {
			e, ok := m.buckets.Get(i)
			if !ok {
				panic("arena map: bucket index out of bounds")
			}
			for e != nil {
				if !yield(e.val) {
					return
				}
				e = e.next
			}
		}
	}
}

// All returns an iterator over all key-value pairs in the map
// Example:
//
//	m := arena.NewMap[string, int](a)
//	m.Set("a", 1)
//	m.Set("b", 2)
//	for key, val := range m.All() {
//	    fmt.Printf("%s: %d\n", key, val)
//	}
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for i := 0; i < m.cap; i++ {
			e, ok := m.buckets.Get(i)
			if !ok {
				panic("arena map: bucket index out of bounds")
			}
			for e != nil {
				if !yield(e.key, e.val) {
					return
				}
				e = e.next
			}
		}
	}
}

// MapIter provides pull-based iteration over map entries
type MapIter[K comparable, V any] struct {
	m       *Map[K, V]
	index   int
	current *entry[K, V]
}

// Iter returns a pull-based iterator for the map
// Use Next() to pull key-value pairs one by one.
//
// Example:
//
//	m := arena.NewMap[string, int](a)
//	m.Set("a", 1)
//	m.Set("b", 2)
//
//	iter := m.Iter()
//	for key, val, ok := iter.Next(); ok; key, val, ok = iter.Next() {
//	    fmt.Printf("%s: %d\n", key, val)
//	}
func (m *Map[K, V]) Iter() *MapIter[K, V] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	it := &MapIter[K, V]{
		m:       m,
		index:   0,
		current: nil,
	}

	// Find first non-empty bucket
	for it.index < m.cap {
		if e, ok := m.buckets.Get(it.index); ok && e != nil {
			it.current = e
			break
		}
		it.index++
	}

	return it
}

// Next returns the next key-value pair and whether it exists
// Returns (zero_key, zero_value, false) when iteration is complete.
func (it *MapIter[K, V]) Next() (K, V, bool) {
	it.m.mu.RLock()
	defer it.m.mu.RUnlock()

	if it.current == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	// Get current entry
	key := it.current.key
	val := it.current.val

	// Advance to next entry
	it.current = it.current.next

	// If current chain is exhausted, find next non-empty bucket
	if it.current == nil {
		it.index++
		for it.index < it.m.cap {
			if e, ok := it.m.buckets.Get(it.index); ok && e != nil {
				it.current = e
				break
			}
			it.index++
		}
	}

	return key, val, true
}
