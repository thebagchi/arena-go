package container

import (
	"hash/maphash"
	"iter"
	"reflect"
	"sync"
	"unsafe"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

const (
	// INITIAL_BUCKET_COUNT is the bucket array size a new map starts with.
	INITIAL_BUCKET_COUNT = 16
	// LOAD_NUMERATOR is the numerator of the load factor at which the bucket
	// array doubles.
	LOAD_NUMERATOR = 3
	// LOAD_DENOMINATOR is its denominator, so the two together are three
	// quarters full.
	LOAD_DENOMINATOR = 4
)

// Map is a hash map whose entries and bucket array live in an arena, using
// separate chaining.
//
// String keys and string values are copied into the arena as they are stored.
// Arena memory is not scanned by the garbage collector, so a heap string kept
// only by a map entry would be collected and its bytes reused underneath the
// map; copying is what makes the common case safe. A key or value of any other
// pointer-bearing type is the caller's to copy, as the package documentation
// says.
//
// Safe for concurrent use: reads share the lock, writes take it exclusively.
type Map[K comparable, V any] struct {
	mtx     sync.RWMutex
	arena   *arena.Arena
	buckets *Vec[*_Entry[K, V]]
	seed    maphash.Seed
	count   int
	mask    uint64
	keyStr  bool
	valStr  bool
}

// _Entry is one link in a bucket's chain.
type _Entry[K comparable, V any] struct {
	next *_Entry[K, V]
	hash uint64
	key  K
	val  V
}

// NewMap returns an empty map backed by the arena.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func NewMap[K comparable, V any](a *arena.Arena) *Map[K, V] {
	buckets := NewVec[*_Entry[K, V]](a)
	buckets.Resize(INITIAL_BUCKET_COUNT)

	return &Map[K, V]{
		arena:   a,
		buckets: buckets,
		seed:    maphash.MakeSeed(),
		mask:    INITIAL_BUCKET_COUNT - 1,
		keyStr:  _IsString[K](),
		valStr:  _IsString[V](),
	}
}

// Set inserts or replaces the value stored under key.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Set(key K, value V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if m.count > m._Buckets()*LOAD_NUMERATOR/LOAD_DENOMINATOR {
		m._Grow()
	}

	if m.valStr {
		m._Intern(unsafe.Pointer(&value))
	}

	var (
		hash  = maphash.Comparable(m.seed, key)
		index = int(hash & m.mask)
		head  = m.buckets.At(index)
	)

	for entry := head; entry != nil; entry = entry.next {
		if entry.hash == hash && entry.key == key {
			entry.val = value

			return
		}
	}

	if m.keyStr {
		m._Intern(unsafe.Pointer(&key))
	}

	entry := arena.Alloc[_Entry[K, V]](m.arena)
	entry.hash, entry.key, entry.val, entry.next = hash, key, value, head

	m.buckets.Set(index, entry)
	m.count = m.count + 1
}

// Get returns the value stored under key, reporting whether there was one.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	var (
		hash  = maphash.Comparable(m.seed, key)
		index = int(hash & m.mask)
	)

	for entry := m.buckets.At(index); entry != nil; entry = entry.next {
		if entry.hash == hash && entry.key == key {
			return entry.val, true
		}
	}

	var zero V

	return zero, false
}

// Contains reports whether key is present.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Contains(key K) bool {
	_, found := m.Get(key)

	return found
}

// Delete removes key and frees its entry.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Delete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var (
		hash  = maphash.Comparable(m.seed, key)
		index = int(hash & m.mask)
		prev  *_Entry[K, V]
	)

	for entry := m.buckets.At(index); entry != nil; entry = entry.next {
		if entry.hash != hash || entry.key != key {
			prev = entry

			continue
		}

		if prev == nil {
			m.buckets.Set(index, entry.next)
		} else {
			prev.next = entry.next
		}

		m.arena.Remove(unsafe.Pointer(entry))
		m.count = m.count - 1

		return
	}
}

// Len returns the number of entries.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.count
}

// Range calls visit for each entry until it returns false.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Range(visit func(K, V) bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for _, head := range m.buckets.Slice() {
		for entry := head; entry != nil; entry = entry.next {
			if !visit(entry.key, entry.val) {
				return
			}
		}
	}
}

// Reset frees every entry and empties the map, keeping the bucket array.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	buckets := m.buckets.Slice()
	for i, head := range buckets {
		for entry := head; entry != nil; {
			next := entry.next
			m.arena.Remove(unsafe.Pointer(entry))
			entry = next
		}

		buckets[i] = nil
	}

	m.count = 0
}

// Clone returns a heap map holding the same entries, which outlives the arena.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Clone() map[K]V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if m.count == 0 {
		return nil
	}

	result := make(map[K]V, m.count)

	for _, head := range m.buckets.Slice() {
		for entry := head; entry != nil; entry = entry.next {
			result[entry.key] = entry.val
		}
	}

	return result
}

// Keys returns an iterator over the keys.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.Range(func(key K, _ V) bool {
			return yield(key)
		})
	}
}

// Values returns an iterator over the values.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.Range(func(_ K, val V) bool {
			return yield(val)
		})
	}
}

// All returns an iterator over key and value pairs.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(yield)
	}
}

// Iter returns a pull-based iterator over the entries.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) Iter() *MapIter[K, V] {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	it := &MapIter[K, V]{owner: m}
	it._Advance()

	return it
}

// _Buckets returns the bucket count. The caller holds the lock.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) _Buckets() int {
	return m.buckets.Len()
}

// _Grow doubles the bucket array and rehashes every chain into it.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) _Grow() {
	var (
		old    = m.buckets
		wanted = old.Len() * GROWTH_FACTOR
		grown  = NewVec[*_Entry[K, V]](m.arena)
	)

	grown.Resize(wanted)

	m.buckets = grown
	m.mask = uint64(wanted - 1)

	buckets := grown.Slice()

	for _, head := range old.Slice() {
		for entry := head; entry != nil; {
			next := entry.next
			index := int(entry.hash & m.mask)
			entry.next = buckets[index]
			buckets[index] = entry
			entry = next
		}
	}

	m.arena.Remove(res.SlicePtr(old.Slice()))
}

// _Intern replaces a string in place with a copy that lives in the arena.
//
// The pointer names a value whose type has string layout, which is what lets a
// named string type be copied as well. Doing it through the empty interface
// would box the string and allocate on the Go heap on every insert.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (m *Map[K, V]) _Intern(at unsafe.Pointer) {
	held := (*string)(at)
	*held = m.arena.MakeString(*held)
}

// MapIter walks a map one entry at a time.
type MapIter[K comparable, V any] struct {
	owner   *Map[K, V]
	current *_Entry[K, V]
	index   int
}

// Next returns the next key and value, reporting whether there was one.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (it *MapIter[K, V]) Next() (K, V, bool) {
	it.owner.mtx.RLock()
	defer it.owner.mtx.RUnlock()

	if it.current == nil {
		var (
			zeroKey K
			zeroVal V
		)

		return zeroKey, zeroVal, false
	}

	key, val := it.current.key, it.current.val

	it.current = it.current.next
	if it.current == nil {
		it.index = it.index + 1
		it._Advance()
	}

	return key, val, true
}

// _Advance moves to the first non-empty bucket at or after index.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func (it *MapIter[K, V]) _Advance() {
	buckets := it.owner.buckets.Slice()

	for it.index < len(buckets) {
		if buckets[it.index] != nil {
			it.current = buckets[it.index]

			return
		}

		it.index = it.index + 1
	}
}

// _IsString reports whether T's underlying type is a string, so that values of
// it can be copied into the arena.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func _IsString[T any]() bool {
	var zero T

	return reflect.TypeOf(&zero).Elem().Kind() == reflect.String
}
