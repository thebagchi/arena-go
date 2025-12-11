package arena_test

import (
	"fmt"
	"testing"

	"github.com/thebagchi/arena-go"
)

func TestMap_Keys(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)
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

func TestMap_Values(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)
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

func TestMap_All_Iterator(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)
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

func TestMap_All_EarlyTermination(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)
	for i := 0; i < 10; i++ {
		m.Set(fmt.Sprintf("key%d", i), i)
	}

	count := 0
	for _, _ = range m.All() {
		count++
		if count >= 5 {
			break
		}
	}

	if count != 5 {
		t.Errorf("Expected to stop at 5 iterations, got %d", count)
	}
}

func TestMap_Iter(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)
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

func TestMap_IterEmpty(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

	iter := m.Iter()
	_, _, ok := iter.Next()
	if ok {
		t.Error("Expected no entries in empty map")
	}
}

func TestMap_KeysEmpty(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	m := arena.NewMap[string, int](a)

	count := 0
	for range m.Keys() {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 keys in empty map, got %d", count)
	}
}
