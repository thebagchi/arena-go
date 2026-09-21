package container_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestMap_BasicOperations covers map basic operations.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_BasicOperations(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	// Test empty map
	if m.Len() != 0 {
		t.Errorf("Expected length 0, got %d", m.Len())
	}

	// Test Set and Get
	m.Set("key1", 100)
	m.Set("key2", 200)

	if val, found := m.Get("key1"); !found || val != 100 {
		t.Errorf("Expected key1=100, got %v, %v", val, found)
	}

	if val, found := m.Get("key2"); !found || val != 200 {
		t.Errorf("Expected key2=200, got %v, %v", val, found)
	}

	if _, found := m.Get("nonexistent"); found {
		t.Error("Expected nonexistent key to not be found")
	}

	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}
}

// TestMap_Update covers map update.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Update(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	m.Set("key", 100)

	if val, _ := m.Get("key"); val != 100 {
		t.Errorf("Expected 100, got %v", val)
	}

	m.Set("key", 200) // Update

	if val, _ := m.Get("key"); val != 200 {
		t.Errorf("Expected 200 after update, got %v", val)
	}

	if m.Len() != 1 {
		t.Errorf("Expected length 1 after update, got %d", m.Len())
	}
}

// TestMap_Delete covers map delete.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Delete(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	m.Set("key1", 100)
	m.Set("key2", 200)
	m.Set("key3", 300)

	if m.Len() != 3 {
		t.Errorf("Expected length 3, got %d", m.Len())
	}

	m.Delete("key2")

	if _, found := m.Get("key2"); found {
		t.Error("Expected key2 to be deleted")
	}

	if m.Len() != 2 {
		t.Errorf("Expected length 2 after delete, got %d", m.Len())
	}

	// Delete nonexistent key (should not panic)
	m.Delete("nonexistent")

	if m.Len() != 2 {
		t.Errorf("Expected length 2 after deleting nonexistent key, got %d", m.Len())
	}
}

// TestMap_Range covers map range.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Range(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	entries := map[string]int{
		"key1": 100,
		"key2": 200,
		"key3": 300,
	}

	for k, v := range entries {
		m.Set(k, v)
	}

	collected := make(map[string]int)

	m.Range(func(k string, v int) bool {
		collected[k] = v
		return true
	})

	if len(collected) != len(entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(collected))
	}

	for k, v := range entries {
		if collected[k] != v {
			t.Errorf("Expected %s=%d, got %d", k, v, collected[k])
		}
	}
}

// TestMap_Range_StopEarly covers map range stop early.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Range_StopEarly(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	m.Set("key1", 100)
	m.Set("key2", 200)
	m.Set("key3", 300)

	count := 0

	m.Range(func(k string, v int) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})

	if count != 2 {
		t.Errorf("Expected to stop after 2 iterations, got %d", count)
	}
}

// TestMap_Growth covers map growth.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Growth(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[int, string](a)

	// Add enough entries to trigger growth
	for i := 0; i < 20; i++ {
		m.Set(i, "value"+string(rune(i+48)))
	}

	if m.Len() != 20 {
		t.Errorf("Expected length 20, got %d", m.Len())
	}

	// Verify all entries are still accessible
	for i := 0; i < 20; i++ {
		expected := "value" + string(rune(i+48))
		if val, found := m.Get(i); !found || val != expected {
			t.Errorf("Expected %s for key %d, got %v, %v", expected, i, val, found)
		}
	}
}

// TestMap_Reset covers map reset.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Reset(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	m.Set("key1", 100)
	m.Set("key2", 200)

	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	m.Reset()

	if m.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", m.Len())
	}

	if _, found := m.Get("key1"); found {
		t.Error("Expected key1 to not exist after reset")
	}
}

// TestMap_Clone covers map clone.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Clone(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	m.Set("key1", 100)
	m.Set("key2", 200)

	cloned := m.Clone()

	if len(cloned) != 2 {
		t.Errorf("Expected cloned map length 2, got %d", len(cloned))
	}

	if cloned["key1"] != 100 || cloned["key2"] != 200 {
		t.Errorf("Cloned map contents incorrect: %v", cloned)
	}

	// Modify original (should not affect clone)
	m.Set("key1", 999)

	if cloned["key1"] != 100 {
		t.Error("Clone should be independent of original")
	}
}

// TestMap_EmptyClone covers map empty clone.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_EmptyClone(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	cloned := m.Clone()

	if cloned != nil {
		t.Error("Expected nil for empty map clone")
	}
}

// TestMap_ConcurrentAccess covers map concurrent access.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_ConcurrentAccess(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[int, int](a)

	var wg sync.WaitGroup

	// Start multiple goroutines doing concurrent operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			// Each goroutine does some sets and gets
			for j := 0; j < 100; j++ {
				key := start*100 + j
				m.Set(key, key*2)

				if val, found := m.Get(key); !found || val != key*2 {
					t.Errorf("Concurrent access failed for key %d", key)
				}
			}
		}(i)
	}

	wg.Wait()

	if m.Len() != 1000 {
		t.Errorf("Expected length 1000, got %d", m.Len())
	}
}

// TestMap_DifferentTypes covers map different types.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_DifferentTypes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	// Test with different key/value types
	intToString := container.NewMap[int, string](a)
	intToString.Set(42, "answer")

	if val, found := intToString.Get(42); !found || val != "answer" {
		t.Errorf("Expected 'answer', got %v, %v", val, found)
	}

	// Test with struct values
	type Point struct{ X, Y int }

	pointMap := container.NewMap[string, Point](a)
	pointMap.Set("origin", Point{0, 0})

	if val, found := pointMap.Get("origin"); !found || val.X != 0 || val.Y != 0 {
		t.Errorf("Expected Point{0,0}, got %v, %v", val, found)
	}
}

// BenchmarkMap_Set measures map set.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func BenchmarkMap_Set(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[int, int](a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i, i*2)
	}
}

// BenchmarkMap_Get measures map get.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func BenchmarkMap_Get(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[int, int](a)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		m.Set(i, i*2)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

// BenchmarkMap_Range measures map range.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func BenchmarkMap_Range(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[int, int](a)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		m.Set(i, i*2)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(k, v int) bool {
			return true
		})
	}
}

// TestMap_GetAllocations covers map get allocations.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_GetAllocations(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	// Add some entries to trigger growth
	for i := 0; i < 20; i++ {
		m.Set(fmt.Sprintf("key%d", i), i)
	}

	// Verify all entries are present
	if m.Len() != 20 {
		t.Errorf("Expected 20 entries, got %d", m.Len())
	}

	// Verify we can read all entries
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key%d", i)

		val, ok := m.Get(key)
		if !ok || val != i {
			t.Errorf("Failed to get key%d: got %d, ok=%v", i, val, ok)
		}
	}

	// Reset and verify map is empty
	m.Reset()

	if m.Len() != 0 {
		t.Errorf("Expected 0 entries after reset, got %d", m.Len())
	}

	// Clone should work correctly
	m.Set("test", 42)

	clone := m.Clone()
	if clone["test"] != 42 {
		t.Errorf("Clone failed: expected 42, got %d", clone["test"])
	}
}

// Iterator tests (from map_iter_test.go)

// TestMap_Keys covers map keys.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Keys(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	keys := make(map[string]bool)
	for key := range m.Keys() {
		keys[key] = true
	}

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	for _, k := range []string{"a", "b", "c"} {
		if !keys[k] {
			t.Errorf("Missing key: %s", k)
		}
	}
}

// TestMap_Values covers map values.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Values(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	values := make(map[int]bool)
	for val := range m.Values() {
		values[val] = true
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	for _, v := range []int{1, 2, 3} {
		if !values[v] {
			t.Errorf("Missing value: %d", v)
		}
	}
}

// TestMap_All_Iterator covers map all iterator.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_All_Iterator(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	entries := make(map[string]int)
	for key, val := range m.All() {
		entries[key] = val
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	expected := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range expected {
		if entries[k] != v {
			t.Errorf("Expected %s=%d, got %d", k, v, entries[k])
		}
	}
}

// TestMap_All_EarlyTermination covers map all early termination.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_All_EarlyTermination(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	for i := 0; i < 10; i++ {
		m.Set(fmt.Sprintf("key%d", i), i)
	}

	count := 0
	for range m.All() {
		count++
		if count >= 5 {
			break
		}
	}

	if count != 5 {
		t.Errorf("Expected to stop at 5 iterations, got %d", count)
	}
}

// TestMap_Iter covers map iter.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_Iter(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	entries := make(map[string]int)

	iter := m.Iter()
	for key, val, ok := iter.Next(); ok; key, val, ok = iter.Next() {
		entries[key] = val
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	expected := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range expected {
		if entries[k] != v {
			t.Errorf("Expected %s=%d, got %d", k, v, entries[k])
		}
	}
}

// TestMap_IterEmpty covers map iter empty.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_IterEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	iter := m.Iter()

	_, _, ok := iter.Next()
	if ok {
		t.Error("Expected no entries in empty map")
	}
}

// TestMap_KeysEmpty covers map keys empty.
//
// Revisions:
//   - 2025-12-11 23:51: initial creation
func TestMap_KeysEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(4096 * 4096))
	defer a.Delete()

	m := container.NewMap[string, int](a)

	count := 0
	for range m.Keys() {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 keys in empty map, got %d", count)
	}
}
