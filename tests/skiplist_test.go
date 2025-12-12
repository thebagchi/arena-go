package arena_test

import (
	"testing"

	"github.com/thebagchi/arena-go"
)

func TestSkipListInsertSearch(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	// Insert some values
	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")
	sl.Insert(3, "three")
	sl.Insert(20, "twenty")

	// Search for existing values
	tests := []struct {
		key      int
		expected string
		found    bool
	}{
		{10, "ten", true},
		{5, "five", true},
		{15, "fifteen", true},
		{3, "three", true},
		{20, "twenty", true},
		{100, "", false},
		{1, "", false},
	}

	for _, tt := range tests {
		val, found := sl.Search(tt.key)
		if found != tt.found {
			t.Errorf("Search(%d): expected found=%v, got %v", tt.key, tt.found, found)
		}
		if found && val != tt.expected {
			t.Errorf("Search(%d): expected %s, got %s", tt.key, tt.expected, val)
		}
	}
}

func TestSkipListUpdate(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	// Insert and update
	sl.Insert(10, "ten")
	sl.Insert(10, "TEN")

	val, found := sl.Search(10)
	if !found {
		t.Fatal("Expected to find key 10")
	}
	if val != "TEN" {
		t.Errorf("Expected 'TEN', got %s", val)
	}

	if sl.Len() != 1 {
		t.Errorf("Expected length 1, got %d", sl.Len())
	}
}

func TestSkipListDelete(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	// Insert values
	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")

	// Delete existing
	if !sl.Delete(10) {
		t.Error("Expected Delete(10) to return true")
	}

	if _, found := sl.Search(10); found {
		t.Error("Key 10 should have been deleted")
	}

	// Delete non-existing
	if sl.Delete(100) {
		t.Error("Expected Delete(100) to return false")
	}

	// Verify remaining elements
	if sl.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sl.Len())
	}
}

func TestSkipListContains(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")

	if !sl.Contains(10) {
		t.Error("Expected Contains(10) to return true")
	}
	if !sl.Contains(5) {
		t.Error("Expected Contains(5) to return true")
	}
	if sl.Contains(100) {
		t.Error("Expected Contains(100) to return false")
	}
}

func TestSkipListMinMax(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	// Empty list
	_, _, found := sl.Min()
	if found {
		t.Error("Expected Min() to return false for empty list")
	}

	_, _, found = sl.Max()
	if found {
		t.Error("Expected Max() to return false for empty list")
	}

	// Insert values
	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")
	sl.Insert(3, "three")
	sl.Insert(20, "twenty")

	// Check min
	minKey, minVal, found := sl.Min()
	if !found {
		t.Fatal("Expected Min() to find a value")
	}
	if minKey != 3 {
		t.Errorf("Expected min key 3, got %d", minKey)
	}
	if minVal != "three" {
		t.Errorf("Expected min value 'three', got %s", minVal)
	}

	// Check max
	maxKey, maxVal, found := sl.Max()
	if !found {
		t.Fatal("Expected Max() to find a value")
	}
	if maxKey != 20 {
		t.Errorf("Expected max key 20, got %d", maxKey)
	}
	if maxVal != "twenty" {
		t.Errorf("Expected max value 'twenty', got %s", maxVal)
	}
}

func TestSkipListRange(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")
	sl.Insert(3, "three")

	// Collect all elements
	var keys []int
	var values []string
	sl.Range(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	// Should be in sorted order
	expected := []int{3, 5, 10, 15}
	if len(keys) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("Expected key[%d] = %d, got %d", i, expected[i], k)
		}
	}

	// Test early termination
	count := 0
	sl.Range(func(k int, v string) bool {
		count++
		return count < 2
	})
	if count != 2 {
		t.Errorf("Expected to iterate 2 times, got %d", count)
	}
}

func TestSkipListLen(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	if sl.Len() != 0 {
		t.Errorf("Expected empty list to have length 0, got %d", sl.Len())
	}

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")

	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}

	sl.Delete(10)

	if sl.Len() != 2 {
		t.Errorf("Expected length 2 after delete, got %d", sl.Len())
	}
}

func TestSkipListReset(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")

	sl.Reset()

	if sl.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", sl.Len())
	}

	if _, found := sl.Search(10); found {
		t.Error("Expected no elements after reset")
	}
}

func TestSkipListClone(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")

	// Clone to map
	m := sl.Clone()
	if m == nil {
		t.Fatal("Expected non-nil map")
	}
	if len(m) != 3 {
		t.Errorf("Expected map length 3, got %d", len(m))
	}
	if m[10] != "ten" || m[5] != "five" || m[15] != "fifteen" {
		t.Error("Clone map has incorrect values")
	}

	// Empty skip list
	sl2 := arena.NewSkipList[int, string](a)
	m2 := sl2.Clone()
	if m2 != nil {
		t.Error("Expected nil map for empty skip list")
	}
}

func TestSkipListCloneSlice(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")
	sl.Insert(3, "three")

	// Clone to slice
	s := sl.CloneSlice()
	if s == nil {
		t.Fatal("Expected non-nil slice")
	}
	if len(s) != 4 {
		t.Errorf("Expected slice length 4, got %d", len(s))
	}

	// Should be in sorted order
	expected := []struct {
		key int
		val string
	}{
		{3, "three"},
		{5, "five"},
		{10, "ten"},
		{15, "fifteen"},
	}

	for i, pair := range s {
		if pair.Key != expected[i].key || pair.Value != expected[i].val {
			t.Errorf("Expected pair[%d] = {%d, %s}, got {%d, %s}",
				i, expected[i].key, expected[i].val, pair.Key, pair.Value)
		}
	}

	// Empty skip list
	sl2 := arena.NewSkipList[int, string](a)
	s2 := sl2.CloneSlice()
	if s2 != nil {
		t.Error("Expected nil slice for empty skip list")
	}
}

func TestSkipListIterators(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, string](a)

	sl.Insert(10, "ten")
	sl.Insert(5, "five")
	sl.Insert(15, "fifteen")
	sl.Insert(3, "three")

	// Test All()
	var keys []int
	var values []string
	for k, v := range sl.All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	expectedKeys := []int{3, 5, 10, 15}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("All(): expected %d pairs, got %d", len(expectedKeys), len(keys))
	}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("All(): expected key[%d] = %d, got %d", i, expectedKeys[i], k)
		}
	}

	// Test Keys()
	keys = nil
	for k := range sl.Keys() {
		keys = append(keys, k)
	}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Keys(): expected %d keys, got %d", len(expectedKeys), len(keys))
	}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Keys(): expected key[%d] = %d, got %d", i, expectedKeys[i], k)
		}
	}

	// Test Values()
	values = nil
	for v := range sl.Values() {
		values = append(values, v)
	}
	expectedValues := []string{"three", "five", "ten", "fifteen"}
	if len(values) != len(expectedValues) {
		t.Fatalf("Values(): expected %d values, got %d", len(expectedValues), len(values))
	}
	for i, v := range values {
		if v != expectedValues[i] {
			t.Errorf("Values(): expected value[%d] = %s, got %s", i, expectedValues[i], v)
		}
	}
}

func TestSkipListStringKeys(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[string, int](a)

	sl.Insert("banana", 2)
	sl.Insert("apple", 1)
	sl.Insert("cherry", 3)

	// Should be in sorted order
	var keys []string
	for k := range sl.Keys() {
		keys = append(keys, k)
	}

	expected := []string{"apple", "banana", "cherry"}
	if len(keys) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("Expected key[%d] = %s, got %s", i, expected[i], k)
		}
	}
}

func TestSkipListFloatKeys(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[float64, string](a)

	sl.Insert(3.14, "pi")
	sl.Insert(2.71, "e")
	sl.Insert(1.41, "sqrt2")

	minKey, minVal, found := sl.Min()
	if !found || minKey != 1.41 || minVal != "sqrt2" {
		t.Errorf("Min(): expected (1.41, sqrt2), got (%f, %s)", minKey, minVal)
	}

	maxKey, maxVal, found := sl.Max()
	if !found || maxKey != 3.14 || maxVal != "pi" {
		t.Errorf("Max(): expected (3.14, pi), got (%f, %s)", maxKey, maxVal)
	}
}

func TestSkipListManyElements(t *testing.T) {
	a := arena.New(10, arena.BUMP)
	defer a.Delete()

	sl := arena.NewSkipList[int, int](a)

	// Insert many elements
	n := 1000
	for i := 0; i < n; i++ {
		sl.Insert(i, i*10)
	}

	if sl.Len() != n {
		t.Errorf("Expected length %d, got %d", n, sl.Len())
	}

	// Verify all elements
	for i := 0; i < n; i++ {
		val, found := sl.Search(i)
		if !found {
			t.Errorf("Expected to find key %d", i)
		}
		if val != i*10 {
			t.Errorf("Expected value %d, got %d", i*10, val)
		}
	}

	// Verify order
	prev := -1
	for k := range sl.Keys() {
		if k <= prev {
			t.Errorf("Keys not in sorted order: %d after %d", k, prev)
		}
		prev = k
	}
}
