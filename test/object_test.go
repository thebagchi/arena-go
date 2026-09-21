package arena_test

import (
	"testing"
	"unsafe"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
)

// TestAppend covers append.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAppend(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test basic append
	slice := arena.MakeSlice[int](a, 2, 4)
	slice[0] = 1
	slice[1] = 2

	// Append elements
	slice = arena.Append(a, slice, 3, 4)
	if len(slice) != 4 {
		t.Errorf("Expected length 4, got %d", len(slice))
	}

	if cap(slice) != 4 {
		t.Errorf("Expected capacity 4, got %d", cap(slice))
	}

	expected := []int{1, 2, 3, 4}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected slice[%d] = %d, got %d", i, v, slice[i])
		}
	}

	// Test append that requires growing
	slice = arena.Append(a, slice, 5, 6, 7) // This should grow the slice
	if len(slice) != 7 {
		t.Errorf("Expected length 7 after growth, got %d", len(slice))
	}

	if slice[4] != 5 || slice[5] != 6 || slice[6] != 7 {
		t.Errorf("Growth append failed: got %v", slice)
	}

	// Test append empty
	slice = arena.Append(a, slice) // Should not change anything
	if len(slice) != 7 {
		t.Errorf("Empty append changed length to %d", len(slice))
	}
}

// TestAppendStrings covers append strings.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAppendStrings(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[string](a, 1, 2)
	slice[0] = "hello"

	slice = arena.Append(a, slice, "world", "arena")
	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}

	expected := []string{"hello", "world", "arena"}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected slice[%d] = %q, got %q", i, v, slice[i])
		}
	}
}

// TestOwns covers owns.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestOwns(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test nil pointer
	if a.Owns(nil) {
		t.Error("a.Owns(nil) should return false")
	}

	// Test heap pointer
	heapPtr := new(int)
	if arena.OwnsPtr(a, heapPtr) {
		t.Error("OwnsPtr should return false for heap pointers")
	}

	// Test arena pointers
	obj := arena.MakeObject[int](a)
	if !arena.OwnsPtr(a, obj) {
		t.Error("OwnsPtr should return true for arena pointers")
	}

	slice := arena.MakeSlice[int](a, 5, 10)
	if !arena.OwnsSlice(a, slice) {
		t.Error("OwnsSlice should return true for arena slice pointers")
	}

	str := a.MakeString("test")
	if !arena.OwnsString(a, str) {
		t.Error("OwnsString should return true for arena string pointers")
	}

	// Test with nil pointer
	invalidPtr := unsafe.Pointer(nil)
	if a.Owns(invalidPtr) {
		t.Error("Owns should return false for nil pointer")
	}
}

// TestPtr covers ptr.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestPtr(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test with int
	value := 42

	ptr := arena.Ptr(a, value)
	if ptr == nil {
		t.Fatal("Ptr returned nil")
	}

	if *ptr != 42 {
		t.Errorf("Expected *ptr = 42, got %d", *ptr)
	}

	if !arena.OwnsPtr(a, ptr) {
		t.Error("Ptr should allocate in arena")
	}

	// Test modification
	*ptr = 100
	if *ptr != 100 {
		t.Errorf("Expected *ptr = 100 after modification, got %d", *ptr)
	}
	// Original value should be unchanged
	if value != 42 {
		t.Error("Original value should not be modified")
	}

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Alice", Age: 30}

	personPtr := arena.Ptr(a, person)
	if personPtr == nil {
		t.Fatal("Ptr returned nil for struct")
	}

	if personPtr.Name != "Alice" || personPtr.Age != 30 {
		t.Errorf("Expected Person{Alice, 30}, got %+v", *personPtr)
	}

	if !arena.OwnsPtr(a, personPtr) {
		t.Error("Ptr should allocate struct in arena")
	}

	// Test with string
	str := "hello"

	strPtr := arena.Ptr(a, str)
	if *strPtr != "hello" {
		t.Errorf("Expected string 'hello', got '%s'", *strPtr)
	}

	// Test with slice (copies the slice header, not the backing array)
	slice := []int{1, 2, 3}

	slicePtr := arena.Ptr(a, slice)
	if len(*slicePtr) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(*slicePtr))
	}

	if (*slicePtr)[0] != 1 || (*slicePtr)[1] != 2 || (*slicePtr)[2] != 3 {
		t.Errorf("Expected slice [1 2 3], got %v", *slicePtr)
	}
}

// ===== ALLOC TESTS =====

// TestAlloc covers alloc.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAlloc(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Test allocation of int
	intPtr := arena.Alloc[int](a)
	if intPtr == nil {
		t.Fatal("Alloc returned nil for int")
	}

	if *intPtr != 0 {
		t.Errorf("Expected zero-initialized value, got %d", *intPtr)
	}

	*intPtr = 99
	if *intPtr != 99 {
		t.Errorf("Expected 99, got %d", *intPtr)
	}

	if !arena.OwnsPtr(a, intPtr) {
		t.Error("Alloc should allocate in arena")
	}

	// Test allocation of struct
	type Point struct {
		X, Y int
	}

	pointPtr := arena.Alloc[Point](a)
	if pointPtr == nil {
		t.Fatal("Alloc returned nil for struct")
	}

	if pointPtr.X != 0 || pointPtr.Y != 0 {
		t.Errorf("Expected {0, 0}, got {%d, %d}", pointPtr.X, pointPtr.Y)
	}

	pointPtr.X = 10

	pointPtr.Y = 20
	if pointPtr.X != 10 || pointPtr.Y != 20 {
		t.Errorf("Expected {10, 20}, got {%d, %d}", pointPtr.X, pointPtr.Y)
	}
}

// TestAllocMultiple covers alloc multiple.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAllocMultiple(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Allocate multiple objects
	ptr1 := arena.Alloc[int](a)
	ptr2 := arena.Alloc[int](a)
	ptr3 := arena.Alloc[int](a)

	*ptr1 = 1
	*ptr2 = 2
	*ptr3 = 3

	if *ptr1 != 1 || *ptr2 != 2 || *ptr3 != 3 {
		t.Errorf("Expected {1, 2, 3}, got {%d, %d, %d}", *ptr1, *ptr2, *ptr3)
	}

	// All should be owned by arena
	if !arena.OwnsPtr(a, ptr1) || !arena.OwnsPtr(a, ptr2) || !arena.OwnsPtr(a, ptr3) {
		t.Error("All allocated objects should be owned by arena")
	}
}

// ===== MAKE SLICE TESTS =====

// TestMakeSlice covers make slice.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeSlice(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 5, 10)
	if len(slice) != 5 {
		t.Errorf("Expected length 5, got %d", len(slice))
	}

	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}

	// Test zero initialization
	for i := 0; i < len(slice); i++ {
		if slice[i] != 0 {
			t.Errorf("Expected zero-initialized slice[%d], got %d", i, slice[i])
		}
	}

	// Test writing
	for i := 0; i < len(slice); i++ {
		slice[i] = i * 10
	}

	for i := 0; i < len(slice); i++ {
		if slice[i] != i*10 {
			t.Errorf("Expected slice[%d] = %d, got %d", i, i*10, slice[i])
		}
	}
}

// TestMakeSliceCapacity covers make slice capacity.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeSliceCapacity(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Create slice with length < capacity
	slice := arena.MakeSlice[string](a, 2, 10)
	if len(slice) != 2 {
		t.Errorf("Expected length 2, got %d", len(slice))
	}

	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}

	slice[0] = "a"
	slice[1] = "b"

	if slice[0] != "a" || slice[1] != "b" {
		t.Errorf("Expected {a, b}, got {%s, %s}", slice[0], slice[1])
	}
}

// TestMakeSliceZeroLength covers make slice zero length.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeSliceZeroLength(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 0, 10)
	if len(slice) != 0 {
		t.Errorf("Expected length 0, got %d", len(slice))
	}

	if cap(slice) != 10 {
		t.Errorf("Expected capacity 10, got %d", cap(slice))
	}

	// Should be able to append to it
	slice = arena.Append(a, slice, 1, 2, 3)
	if len(slice) != 3 {
		t.Errorf("Expected length 3 after append, got %d", len(slice))
	}
}

// TestMakeSliceLarge covers make slice large.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeSliceLarge(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 1000, 2000)
	if len(slice) != 1000 {
		t.Errorf("Expected length 1000, got %d", len(slice))
	}

	if cap(slice) != 2000 {
		t.Errorf("Expected capacity 2000, got %d", cap(slice))
	}

	// Test some values
	slice[0] = 42

	slice[999] = 999
	if slice[0] != 42 || slice[999] != 999 {
		t.Errorf("Expected {42, ..., 999}, got {%d, ..., %d}", slice[0], slice[999])
	}
}

// TestSliceOwnership covers slice ownership.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestSliceOwnership(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 10, 20)
	if !arena.OwnsSlice(a, slice) {
		t.Error("Slice should be owned by arena")
	}

	// Heap slice should not be owned
	heapSlice := make([]int, 10)
	if arena.OwnsSlice(a, heapSlice) {
		t.Error("Heap slice should not be owned by arena")
	}
}

// ===== MAKE STRING TESTS =====

// TestMakeString covers make string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	str := a.MakeString("hello world")
	if str != "hello world" {
		t.Errorf("Expected 'hello world', got %q", str)
	}

	if !arena.OwnsString(a, str) {
		t.Error("String should be owned by arena")
	}
}

// TestMakeStringEmpty covers make string empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeStringEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	str := a.MakeString("")
	if str != "" {
		t.Errorf("Expected empty string, got %q", str)
	}

	if arena.OwnsString(a, str) {
		t.Error("Empty string should not be owned by arena")
	}
}

// TestMakeStringUnicode covers make string unicode.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeStringUnicode(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	str := a.MakeString("Hello 世界 🚀")
	if str != "Hello 世界 🚀" {
		t.Errorf("Expected 'Hello 世界 🚀', got %q", str)
	}

	if !arena.OwnsString(a, str) {
		t.Error("Unicode string should be owned by arena")
	}
}

// TestMakeStringLarge covers make string large.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMakeStringLarge(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	largeStr := "x"
	for i := 0; i < 1000; i++ {
		largeStr += "x"
	}

	str := a.MakeString(largeStr)
	if str != largeStr {
		t.Error("Large string content mismatch")
	}

	if len(str) != 1001 {
		t.Errorf("Expected length 1001, got %d", len(str))
	}

	if !arena.OwnsString(a, str) {
		t.Error("Large string should be owned by arena")
	}
}

// ===== CLONE TESTS =====

// TestCloneObject covers clone object.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneObject(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	type Data struct {
		Value int
		Name  string
	}

	// Create object in arena
	arenaObj := arena.Alloc[Data](a)
	arenaObj.Value = 42
	arenaObj.Name = "test"

	// Clone to heap
	heapObj := arena.CloneObject(arenaObj)

	if heapObj == nil {
		t.Fatal("CloneObject returned nil")
	}

	if heapObj == arenaObj {
		t.Error("Clone should create a different object")
	}

	if heapObj.Value != 42 || heapObj.Name != "test" {
		t.Errorf("Expected {42, test}, got {%d, %s}", heapObj.Value, heapObj.Name)
	}

	// Verify they're independent
	heapObj.Value = 100
	if arenaObj.Value != 42 {
		t.Error("Modifying clone should not affect original")
	}
}

// TestCloneObjectNil covers clone object nil.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneObjectNil(t *testing.T) {
	var nilPtr *int

	cloned := arena.CloneObject(nilPtr)
	if cloned != nil {
		t.Error("CloneObject of nil should return nil")
	}
}

// TestCloneSlice covers clone slice.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneSlice(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	arenaSlice := arena.MakeSlice[int](a, 5, 10)
	for i := 0; i < len(arenaSlice); i++ {
		arenaSlice[i] = i * 10
	}

	heapSlice := arena.CloneSlice(arenaSlice)

	if len(heapSlice) != 5 {
		t.Errorf("Expected length 5, got %d", len(heapSlice))
	}

	// Verify content
	for i := 0; i < len(heapSlice); i++ {
		if heapSlice[i] != arenaSlice[i] {
			t.Errorf(
				"Expected heapSlice[%d] = %d, got %d",
				i,
				arenaSlice[i],
				heapSlice[i],
			)
		}
	}

	// Verify independence
	heapSlice[0] = 999

	if arenaSlice[0] != 0 {
		t.Error("Modifying clone should not affect original")
	}

	// Heap slice should not be owned by arena
	if arena.OwnsSlice(a, heapSlice) {
		t.Error("Cloned slice should not be owned by arena")
	}
}

// TestCloneSliceEmpty covers clone slice empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneSliceEmpty(t *testing.T) {
	emptySlice := []int{}

	cloned := arena.CloneSlice(emptySlice)
	if cloned != nil {
		t.Error("CloneSlice of empty slice should return nil")
	}
}

// TestCloneString covers clone string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	arenaStr := a.MakeString("hello arena")
	heapStr := arena.CloneString(arenaStr)

	if heapStr != "hello arena" {
		t.Errorf("Expected 'hello arena', got %q", heapStr)
	}

	// Heap string should not be owned by arena
	if arena.OwnsString(a, heapStr) {
		t.Error("Cloned string should not be owned by arena")
	}
}

// TestCloneStringEmpty covers clone string empty.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestCloneStringEmpty(t *testing.T) {
	cloned := arena.CloneString("")
	if cloned != "" {
		t.Errorf("Expected empty string, got %q", cloned)
	}
}

// ===== DELETE TESTS =====

// TestDeleteObject covers delete object.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestDeleteObject(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	obj := arena.Alloc[int](a)
	*obj = 42

	// DeleteObject should not panic
	arena.DeleteObject(a, obj)
}

// TestDeleteSlice covers delete slice.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestDeleteSlice(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 10, 20)
	for i := 0; i < len(slice); i++ {
		slice[i] = i
	}

	// DeleteSlice should not panic
	arena.DeleteSlice(a, slice)
}

// TestDeleteString covers delete string.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestDeleteString(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	str := a.MakeString("test string")

	// DeleteString should not panic
	arena.DeleteString(a, str)
}

// ===== APPEND EDGE CASES =====

// TestAppendGrowth covers append growth.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAppendGrowth(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 0, 2)

	// Single element append requiring growth
	slice = arena.Append(a, slice, 1)
	if len(slice) != 1 || slice[0] != 1 {
		t.Errorf("Expected [1], got %v", slice)
	}

	// Multiple element append requiring growth
	slice = arena.Append(a, slice, 2, 3, 4, 5)
	if len(slice) != 5 {
		t.Errorf("Expected length 5, got %d", len(slice))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range expected {
		if slice[i] != v {
			t.Errorf("Expected slice[%d] = %d, got %d", i, v, slice[i])
		}
	}
}

// TestAppendMultipleTypes covers append multiple types.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAppendMultipleTypes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	type Item struct {
		ID    int
		Value string
	}

	slice := arena.MakeSlice[Item](a, 0, 5)
	slice = arena.Append(a, slice, Item{1, "one"}, Item{2, "two"})

	if len(slice) != 2 {
		t.Errorf("Expected length 2, got %d", len(slice))
	}

	if slice[0].ID != 1 || slice[0].Value != "one" {
		t.Errorf("Expected {1, one}, got {%d, %s}", slice[0].ID, slice[0].Value)
	}

	if slice[1].ID != 2 || slice[1].Value != "two" {
		t.Errorf("Expected {2, two}, got {%d, %s}", slice[1].ID, slice[1].Value)
	}
}

// TestAppendToZeroCapacity covers append to zero capacity.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAppendToZeroCapacity(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 0, 0)
	slice = arena.Append(a, slice, 1, 2, 3)

	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}

	if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", slice)
	}
}

// ===== INTEGRATION TESTS =====

// TestAllocationLifecycle covers allocation lifecycle.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestAllocationLifecycle(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	type Person struct {
		Name string
		Age  int
	}

	// Allocate
	person := arena.Alloc[Person](a)
	person.Name = "Alice"
	person.Age = 30

	// Verify ownership
	if !arena.OwnsPtr(a, person) {
		t.Error("Object should be owned by arena")
	}

	// Clone to heap
	heapPerson := arena.CloneObject(person)
	if arena.OwnsPtr(a, heapPerson) {
		t.Error("Cloned object should not be owned by arena")
	}

	// Verify independence
	person.Age = 31
	if heapPerson.Age != 30 {
		t.Error("Heap copy should not be affected by arena changes")
	}

	// Delete from arena
	arena.DeleteObject(a, person)
}

// TestMixedAllocations covers mixed allocations.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestMixedAllocations(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Allocate different types
	intPtr := arena.Alloc[int](a)
	strSlice := arena.MakeSlice[string](a, 3, 5)
	arenaStr := a.MakeString("test")
	intSlice := arena.MakeSlice[int](a, 10, 20)

	*intPtr = 42
	strSlice[0] = "hello"
	strSlice[1] = "world"

	// Verify all are owned by arena
	if !arena.OwnsPtr(a, intPtr) {
		t.Error("intPtr should be owned by arena")
	}

	if !arena.OwnsSlice(a, strSlice) {
		t.Error("strSlice should be owned by arena")
	}

	if !arena.OwnsString(a, arenaStr) {
		t.Error("arenaStr should be owned by arena")
	}

	if !arena.OwnsSlice(a, intSlice) {
		t.Error("intSlice should be owned by arena")
	}
}

// TestArenaReset covers arena reset.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestArenaReset(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))

	// Allocate some objects
	slice := arena.MakeSlice[int](a, 10, 20)
	for i := 0; i < len(slice); i++ {
		slice[i] = i
	}

	obj := arena.Alloc[int](a)
	*obj = 42

	// Reset should not panic
	a.Reset()

	// After reset, we should be able to allocate again
	newSlice := arena.MakeSlice[int](a, 5, 10)
	if len(newSlice) != 5 {
		t.Errorf("After reset, expected length 5, got %d", len(newSlice))
	}

	a.Delete()
}

// TestInterleavedOperations covers interleaved operations.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestInterleavedOperations(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	// Interleave allocations, appends, and deletions
	slice := arena.MakeSlice[int](a, 0, 5)
	slice = arena.Append(a, slice, 1)
	obj := arena.Alloc[int](a)
	slice = arena.Append(a, slice, 2, 3)
	*obj = 99
	str := a.MakeString("test")
	slice = arena.Append(a, slice, 4)
	arena.DeleteObject(a, obj)

	if len(slice) != 4 {
		t.Errorf("Expected length 4, got %d", len(slice))
	}

	if str != "test" {
		t.Errorf("Expected 'test', got %q", str)
	}
}

// TestLargeStructAllocation covers large struct allocation.
//
// Revisions:
//   - 2026-01-01 00:03: initial creation
func TestLargeStructAllocation(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	type LargeStruct struct {
		Data [1000]int
		Name string
	}

	obj := arena.Alloc[LargeStruct](a)
	obj.Data[0] = 42
	obj.Data[999] = 999
	obj.Name = "large"

	if obj.Data[0] != 42 || obj.Data[999] != 999 {
		t.Error("Large struct data corruption")
	}

	if !arena.OwnsPtr(a, obj) {
		t.Error("Large struct should be owned by arena")
	}
}
