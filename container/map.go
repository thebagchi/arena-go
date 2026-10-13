package container

import (
	"errors"
	"hash/maphash"
	"iter"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

var (
	// ErrIterStale is what MapIter.Next panics with when the map was written
	// after the iterator was created. Walking on would follow a chain Set may
	// have rehashed, or a node Delete has already freed.
	ErrIterStale = errors.New("container: map changed during iteration")
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
//
// Range copies the entries it will visit and releases the lock before the
// callback, so the callback may call back into the map. It sees the entries
// as they were when Range started. An iterator does not: Next panics with
// ErrIterStale if the map was written after Iter.
type Map[K comparable, V any] struct {
	mtx     sync.RWMutex
	arena   *arena.Arena
	buckets *Vec[*_Entry[K, V]]
	seed    maphash.Seed
	count   int
	gen     uint64
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
//   - 2026-10-13 16:25: free the arena copy of a string value that is replaced
func (m *Map[K, V]) Set(key K, value V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m._Touch()

	if m.count > m._Buckets()*LOAD_NUMERATOR/LOAD_DENOMINATOR {
		m._Grow()
	}

	if m.valStr {
		_CopyString(m.arena, unsafe.Pointer(&value))
	}

	var (
		hash  = maphash.Comparable(m.seed, key)
		index = int(hash & m.mask)
		head  = m.buckets.At(index)
	)

	for entry := head; entry != nil; entry = entry.next {
		if entry.hash == hash && entry.key == key {
			if m.valStr {
				_Release(
					m.arena,
					unsafe.Pointer(&entry.val),
					unsafe.Pointer(&value),
				)
			}

			entry.val = value

			return
		}
	}

	if m.keyStr {
		_CopyString(m.arena, unsafe.Pointer(&key))
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
//   - 2026-10-13 16:25: free the arena copies of a string key and value
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

		heldKey := entry.key
		heldVal := entry.val

		if prev == nil {
			m.buckets.Set(index, entry.next)
		} else {
			prev.next = entry.next
		}

		m._FreeStrings(heldKey, heldVal)
		m.arena.Remove(unsafe.Pointer(entry))
		m.count = m.count - 1
		m._Touch()

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
// The entries are copied under the lock and visit runs after it is released,
// so visit may call back into the map. String keys and values in that copy
// are heap copies: a callback that replaces or deletes the arena string must
// not pull the bytes out from under the copy still being walked.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
//   - 2026-10-13 16:25: release the lock before visit, so a callback can write
func (m *Map[K, V]) Range(visit func(K, V) bool) {
	pairs := m._Pairs()

	for i := range pairs {
		if !visit(pairs[i].key, pairs[i].val) {
			return
		}
	}
}

// Reset frees every entry and empties the map, keeping the bucket array.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
//   - 2026-10-13 16:25: free the arena copies of string keys and values
func (m *Map[K, V]) Reset() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	buckets := m.buckets.Slice()
	for i, head := range buckets {
		for entry := head; entry != nil; {
			var (
				next    = entry.next
				heldKey = entry.key
				heldVal = entry.val
			)

			m._FreeStrings(heldKey, heldVal)
			m.arena.Remove(unsafe.Pointer(entry))
			entry = next
		}

		buckets[i] = nil
	}

	m.count = 0
	m._Touch()
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
//   - 2026-10-13 16:25: record the map's generation so a later write is visible
func (m *Map[K, V]) Iter() *MapIter[K, V] {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	it := &MapIter[K, V]{owner: m, gen: m.gen}
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

// _CopyString replaces a string in place with a copy that lives in the arena.
//
// The pointer names a value whose type has string layout, which is what lets a
// named string type be copied as well. Doing it through the empty interface
// would box the string and allocate on the Go heap on every insert.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
//   - 2026-10-13 16:25: a function of the arena, so a skip list can copy too
func _CopyString(a *arena.Arena, at unsafe.Pointer) {
	held := (*string)(at)
	*held = a.MakeString(*held)
}

// _Release frees the arena copy of a string, unless keep names that same copy.
//
// keep is nil when the string is being dropped rather than replaced. An empty
// string was never allocated.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func _Release(a *arena.Arena, at, keep unsafe.Pointer) {
	held := *(*string)(at)
	if len(held) == 0 {
		return
	}

	if keep != nil && unsafe.StringData(held) == unsafe.StringData(*(*string)(keep)) {
		return
	}

	a.Remove(unsafe.Pointer(unsafe.StringData(held)))
}

// _FreeStrings frees the arena copies of a key and a value when those types
// are strings. The caller has already copied them out of the entry.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func (m *Map[K, V]) _FreeStrings(key K, value V) {
	if m.keyStr {
		_Release(m.arena, unsafe.Pointer(&key), nil)
	}

	if m.valStr {
		_Release(m.arena, unsafe.Pointer(&value), nil)
	}
}

// _Touch records that the map changed, so an iterator created before the
// change can tell. The caller holds the write lock.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func (m *Map[K, V]) _Touch() {
	m.gen = m.gen + 1
}

// _Pair is one entry copied out of the map for Range.
type _Pair[K comparable, V any] struct {
	key K
	val V
}

// _Pairs copies every entry. String keys and values are heap copies, so the
// result does not depend on arena bytes a later write might free. The caller
// does not hold the lock.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func (m *Map[K, V]) _Pairs() []_Pair[K, V] {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	pairs := make([]_Pair[K, V], 0, m.count)

	for _, head := range m.buckets.Slice() {
		for entry := head; entry != nil; entry = entry.next {
			key := entry.key
			val := entry.val

			if m.keyStr {
				key = _CloneString(key)
			}

			if m.valStr {
				val = _CloneString(val)
			}

			pairs = append(pairs, _Pair[K, V]{key: key, val: val})
		}
	}

	return pairs
}

// _CloneString returns a heap copy of a string-typed value. The copy is what
// lets Range free the arena string from inside the callback without the bytes
// under the walk disappearing.
//
// Revisions:
//   - 2026-10-13 16:25: initial creation
func _CloneString[T any](value T) T {
	var (
		held   = *(*string)(unsafe.Pointer(&value))
		cloned = strings.Clone(held)
	)

	*(*string)(unsafe.Pointer(&value)) = cloned

	return value
}

// MapIter walks a map one entry at a time.
type MapIter[K comparable, V any] struct {
	owner   *Map[K, V]
	current *_Entry[K, V]
	index   int
	gen     uint64
}

// Next returns the next key and value, reporting whether there was one.
//
// It panics with ErrIterStale when the map was written after Iter. The node
// it would walk may already have been freed.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
//   - 2026-10-13 16:25: panic when the map's generation has moved on
func (it *MapIter[K, V]) Next() (K, V, bool) {
	it.owner.mtx.RLock()
	defer it.owner.mtx.RUnlock()

	if it.gen != it.owner.gen {
		panic(ErrIterStale)
	}

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
