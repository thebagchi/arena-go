package arena_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/thebagchi/arena-go"
)

func TestMap_BasicOperations(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Update(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Delete(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Range(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Range_StopEarly(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Growth(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[int, string](a)

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

func TestMap_Reset(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_Clone(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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

func TestMap_EmptyClone(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

	cloned := m.Clone()

	if cloned != nil {
		t.Error("Expected nil for empty map clone")
	}
}

func TestMap_ConcurrentAccess(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[int, int](a)

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

func TestMap_DifferentTypes(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	// Test with different key/value types
	intToString := arena.NewMap[int, string](a)
	intToString.Set(42, "answer")

	if val, found := intToString.Get(42); !found || val != "answer" {
		t.Errorf("Expected 'answer', got %v, %v", val, found)
	}

	// Test with struct values
	type Point struct{ X, Y int }
	pointMap := arena.NewMap[string, Point](a)
	pointMap.Set("origin", Point{0, 0})

	if val, found := pointMap.Get("origin"); !found || val.X != 0 || val.Y != 0 {
		t.Errorf("Expected Point{0,0}, got %v, %v", val, found)
	}
}

func BenchmarkMap_Set(b *testing.B) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[int, int](a)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i*2)
	}
}

func BenchmarkMap_Get(b *testing.B) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[int, int](a)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		m.Set(i, i*2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

func BenchmarkMap_Range(b *testing.B) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[int, int](a)

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

func TestMap_GetAllocations(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

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
