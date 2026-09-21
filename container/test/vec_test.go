package container_test

import (
	"reflect"
	"testing"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestVec_Basic covers vec basic.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_Basic(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test empty slice
	slice := container.NewVec[int](a)
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

// TestVec_AppendSlice covers vec append slice.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_AppendSlice(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := container.NewVec[int](a)

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

// TestVec_SSO covers vec sso.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_SSO(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test small slice stays in SSO
	small := container.NewVec[int](a, 1, 2, 3)
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

// TestVec_Reset covers vec reset.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_Reset(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := container.NewVec[int](a)
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

// TestVec_Clone covers vec clone.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_Clone(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))

	slice := container.NewVec[string](a)
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

// TestVec_Iterators covers vec iterators.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_Iterators(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := container.NewVec[int](a)
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
			t.Errorf(
				"All2() iterator failed at index %d: expected %d, got %d",
				i,
				data[i],
				v,
			)
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
	empty := container.NewVec[int](a)

	iter2 := empty.Iter()
	if v, ok := iter2.Next(); ok {
		t.Errorf("Empty iterator should return false, got value %d", v)
	}
}

// TestVec_RangeLoop covers vec range loop.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_RangeLoop(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := container.NewVec[string](a)
	slice.AppendSlice([]string{"apple", "banana", "cherry"})

	// Test range loop
	result := slice.Slice()

	expected := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Range loop failed: expected %v, got %v", expected, result)
	}
}

// TestVec_LargeData covers vec large data.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_LargeData(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 1024 * 4096)) // 100KB arena
	defer a.Delete()

	slice := container.NewVec[int](a)

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
			t.Errorf(
				"Data corruption at index %d: expected %d, got %d",
				i,
				i,
				slice.Slice()[i],
			)
		}
	}
}

// TestVec_Generics covers vec generics.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_Generics(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test with different types
	intSlice := container.NewVec[int](a, 1, 2, 3)
	stringSlice := container.NewVec[string](a, "a", "b", "c")
	boolSlice := container.NewVec[bool](a, true, false, true)

	if intSlice.Len() != 3 || stringSlice.Len() != 3 || boolSlice.Len() != 3 {
		t.Error("Generic type slices failed")
	}

	// Test struct types
	type Point struct {
		X, Y int
	}

	structSlice := container.NewVec[Point](a)
	structSlice.Append(Point{1, 2})
	structSlice.Append(Point{3, 4})

	if structSlice.Len() != 2 {
		t.Errorf("Struct slice failed: expected length 2, got %d", structSlice.Len())
	}
}

// TestVec_EdgeCases covers vec edge cases.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func TestVec_EdgeCases(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test zero-sized types
	slice := container.NewVec[struct{}](a)
	for i := 0; i < 10; i++ {
		slice.Append(struct{}{})
	}

	if slice.Len() != 10 {
		t.Errorf("Zero-sized type slice failed: expected length 10, got %d", slice.Len())
	}

	// Test Clone on empty slice
	empty := container.NewVec[int](a)

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

// BenchmarkVecAppend measures vec append.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func BenchmarkVecAppend(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 1024 * 4096)) // 1MB arena
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		slice := container.NewVec[int](a)
		for j := 0; j < 100; j++ {
			slice.Append(j)
		}

		if i%100 == 0 { // Reset arena periodically
			a.Reset()
		}
	}
}

// BenchmarkVecAppendSlice measures vec append slice.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func BenchmarkVecAppendSlice(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 1024 * 4096))
	defer a.Delete()

	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		slice := container.NewVec[int](a)
		slice.AppendSlice(data)

		if i%100 == 0 {
			a.Reset()
		}
	}
}

// BenchmarkVecIterate measures vec iterate.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func BenchmarkVecIterate(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 1024 * 4096))
	defer a.Delete()

	slice := container.NewVec[int](a)
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

// BenchmarkStandardSlice measures standard slice.
//
// Revisions:
//   - 2025-12-09 23:15: initial creation
func BenchmarkStandardSlice(b *testing.B) {
	for b.Loop() {
		slice := make([]int, 0, 100)
		for j := 0; j < 100; j++ {
			slice = append(slice, j)
		}

		_ = slice // Use the slice to avoid SA4010
	}
}
