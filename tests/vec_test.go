package arena_test

import (
	"reflect"
	"testing"

	"github.com/thebagchi/arena-go"
)

func TestVecBasic(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	// Test empty slice
	slice := arena.NewVec[int](a)
	if slice.Len() != 0 {
		t.Errorf("Expected length 0, got %d", slice.Len())
	}
	if slice.Cap() < 16 { // Should have inline capacity
		t.Errorf("Expected capacity >= 16, got %d", slice.Cap())
	}

	// Test append
	slice.Append(42)
	if slice.Len() != 1 {
		t.Errorf("Expected length 1, got %d", slice.Len())
	}
	if slice.Slice()[0] != 42 {
		t.Errorf("Expected first element 42, got %d", slice.Slice()[0])
	}
}

func TestVecAppendSlice(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	slice := arena.NewVec[int](a)

	// Test appending empty slice
	slice.AppendSlice([]int{})
	if slice.Len() != 0 {
		t.Errorf("Expected length 0 after empty append, got %d", slice.Len())
	}

	// Test appending data
	data := []int{1, 2, 3, 4, 5}
	slice.AppendSlice(data)
	if slice.Len() != 5 {
		t.Errorf("Expected length 5, got %d", slice.Len())
	}

	result := slice.Slice()
	for i, v := range data {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}

	// Test appending more data
	slice.AppendSlice([]int{6, 7, 8})
	if slice.Len() != 8 {
		t.Errorf("Expected length 8, got %d", slice.Len())
	}
}

func TestVecSSO(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	// Test small slice stays in SSO
	small := arena.NewVec[int](a, 1, 2, 3)
	if small.Len() != 3 {
		t.Errorf("Expected length 3, got %d", small.Len())
	}

	// Force migration to arena by appending many elements
	for i := 0; i < 20; i++ {
		small.Append(i + 10)
	}

	if small.Len() != 23 {
		t.Errorf("Expected length 23, got %d", small.Len())
	}

	// Verify data integrity
	slice := small.Slice()
	expected := []int{1, 2, 3}
	for i := 0; i < 3; i++ {
		if slice[i] != expected[i] {
			t.Errorf("Expected slice[%d] = %d, got %d", i, expected[i], slice[i])
		}
	}
	for i := 3; i < 23; i++ {
		if slice[i] != i+7 { // 3 + (i-3) + 10 - 3 = i + 7
			t.Errorf("Expected slice[%d] = %d, got %d", i, i+7, slice[i])
		}
	}
}

func TestVecReset(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	slice := arena.NewVec[int](a)
	slice.AppendSlice([]int{1, 2, 3, 4, 5})

	if slice.Len() != 5 {
		t.Errorf("Expected length 5, got %d", slice.Len())
	}

	capBefore := slice.Cap()
	slice.Reset()

	if slice.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", slice.Len())
	}
	if slice.Cap() != capBefore {
		t.Errorf("Expected capacity to remain %d, got %d", capBefore, slice.Cap())
	}

	// Test reuse after reset
	slice.Append(99)
	if slice.Len() != 1 {
		t.Errorf("Expected length 1 after append, got %d", slice.Len())
	}
	if slice.Slice()[0] != 99 {
		t.Errorf("Expected first element 99, got %d", slice.Slice()[0])
	}
}

func TestVecClone(t *testing.T) {
	a := arena.New(1024, arena.BUMP)

	slice := arena.NewVec[string](a)
	slice.AppendSlice([]string{"hello", "world", "arena"})

	cloned := slice.Clone()
	a.Delete() // Arena is gone, but cloned should still work

	if len(cloned) != 3 {
		t.Errorf("Expected cloned length 3, got %d", len(cloned))
	}

	expected := []string{"hello", "world", "arena"}
	for i, v := range expected {
		if cloned[i] != v {
			t.Errorf("Expected cloned[%d] = %q, got %q", i, v, cloned[i])
		}
	}
}

func TestVecIterators(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	slice := arena.NewVec[int](a)
	data := []int{10, 20, 30, 40, 50}
	slice.AppendSlice(data)

	// Test All() iterator
	var collected []int
	for v := range slice.All() {
		collected = append(collected, v)
	}
	if !reflect.DeepEqual(collected, data) {
		t.Errorf("All() iterator failed: expected %v, got %v", data, collected)
	}

	// Test All2() iterator
	collected = collected[:0]
	for i, v := range slice.All2() {
		if i >= len(data) || v != data[i] {
			t.Errorf("All2() iterator failed at index %d: expected %d, got %d", i, data[i], v)
		}
		collected = append(collected, v)
	}
	if !reflect.DeepEqual(collected, data) {
		t.Errorf("All2() iterator failed: expected %v, got %v", data, collected)
	}

	// Test pull-based iterator
	iter := slice.Iter()
	collected = collected[:0]
	for v, ok := iter.Next(); ok; v, ok = iter.Next() {
		collected = append(collected, v)
	}
	if !reflect.DeepEqual(collected, data) {
		t.Errorf("Pull iterator failed: expected %v, got %v", data, collected)
	}

	// Test iterator on empty slice
	empty := arena.NewVec[int](a)
	iter2 := empty.Iter()
	if v, ok := iter2.Next(); ok {
		t.Errorf("Empty iterator should return false, got value %d", v)
	}
}

func TestVecRangeLoop(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	slice := arena.NewVec[string](a)
	slice.AppendSlice([]string{"apple", "banana", "cherry"})

	// Test range loop
	result := slice.Slice()

	expected := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Range loop failed: expected %v, got %v", expected, result)
	}
}

func TestVecLargeData(t *testing.T) {
	a := arena.New(100*1024, arena.BUMP) // 100KB arena
	defer a.Delete()

	slice := arena.NewVec[int](a)

	// Add a lot of data to force arena allocation
	for i := 0; i < 1000; i++ {
		slice.Append(i)
	}

	if slice.Len() != 1000 {
		t.Errorf("Expected length 1000, got %d", slice.Len())
	}

	// Verify data integrity
	for i := 0; i < 1000; i++ {
		if slice.Slice()[i] != i {
			t.Errorf("Data corruption at index %d: expected %d, got %d", i, i, slice.Slice()[i])
		}
	}
}

func TestVecGenerics(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	// Test with different types
	intSlice := arena.NewVec[int](a, 1, 2, 3)
	stringSlice := arena.NewVec[string](a, "a", "b", "c")
	boolSlice := arena.NewVec[bool](a, true, false, true)

	if intSlice.Len() != 3 || stringSlice.Len() != 3 || boolSlice.Len() != 3 {
		t.Error("Generic type slices failed")
	}

	// Test struct types
	type Point struct {
		X, Y int
	}
	structSlice := arena.NewVec[Point](a)
	structSlice.Append(Point{1, 2})
	structSlice.Append(Point{3, 4})

	if structSlice.Len() != 2 {
		t.Errorf("Struct slice failed: expected length 2, got %d", structSlice.Len())
	}
}

func TestVecEdgeCases(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	defer a.Delete()

	// Test zero-sized types
	slice := arena.NewVec[struct{}](a)
	for i := 0; i < 10; i++ {
		slice.Append(struct{}{})
	}
	if slice.Len() != 10 {
		t.Errorf("Zero-sized type slice failed: expected length 10, got %d", slice.Len())
	}

	// Test Clone on empty slice
	empty := arena.NewVec[int](a)
	cloned := empty.Clone()
	if cloned != nil {
		t.Errorf("Clone of empty slice should return nil, got %v", cloned)
	}

	// Test Reset on empty slice
	empty.Reset()
	if empty.Len() != 0 {
		t.Errorf("Reset on empty slice should keep length 0, got %d", empty.Len())
	}
}

func BenchmarkVecAppend(b *testing.B) {
	a := arena.New(1024*1024, arena.BUMP) // 1MB arena
	defer a.Delete()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := arena.NewVec[int](a)
		for j := 0; j < 100; j++ {
			slice.Append(j)
		}
		if i%100 == 0 { // Reset arena periodically
			a.Reset()
		}
	}
}

func BenchmarkVecAppendSlice(b *testing.B) {
	a := arena.New(1024*1024, arena.BUMP)
	defer a.Delete()

	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := arena.NewVec[int](a)
		slice.AppendSlice(data)
		if i%100 == 0 {
			a.Reset()
		}
	}
}

func BenchmarkVecIterate(b *testing.B) {
	a := arena.New(1024*1024, arena.BUMP)
	defer a.Delete()

	slice := arena.NewVec[int](a)
	for i := 0; i < 1000; i++ {
		slice.Append(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for v := range slice.All() {
			sum += v
		}
	}
}

func BenchmarkStandardSlice(b *testing.B) {
	for b.Loop() {
		slice := make([]int, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, j)
		}
		_ = slice // Use the slice to avoid SA4010
	}
}
