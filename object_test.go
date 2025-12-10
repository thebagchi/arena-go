package arena

import (
	"testing"
	"unsafe"
)

func TestAppend(t *testing.T) {
	a := New(1024, BUMP)
	defer a.Delete()

	// Test basic append
	slice := MakeSlice[int](a, 2, 4)
	slice[0] = 1
	slice[1] = 2

	// Append elements
	slice = Append(a, slice, 3, 4)
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
	slice = Append(a, slice, 5, 6, 7) // This should grow the slice
	if len(slice) != 7 {
		t.Errorf("Expected length 7 after growth, got %d", len(slice))
	}
	if slice[4] != 5 || slice[5] != 6 || slice[6] != 7 {
		t.Errorf("Growth append failed: got %v", slice)
	}

	// Test append empty
	slice = Append(a, slice) // Should not change anything
	if len(slice) != 7 {
		t.Errorf("Empty append changed length to %d", len(slice))
	}
}

func TestAppendStrings(t *testing.T) {
	a := New(1024, BUMP)
	defer a.Delete()

	slice := MakeSlice[string](a, 1, 2)
	slice[0] = "hello"

	slice = Append(a, slice, "world", "arena")
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

func TestOwns(t *testing.T) {
	a := New(1024, BUMP)
	defer a.Delete()

	// Test nil pointer
	if a.Owns(nil) {
		t.Error("Owns(nil) should return false")
	}

	// Test heap pointer
	heapPtr := new(int)
	if OwnsPtr(a, heapPtr) {
		t.Error("OwnsPtr should return false for heap pointers")
	}

	// Test arena pointers
	obj := MakeObject[int](a)
	if !OwnsPtr(a, obj) {
		t.Error("OwnsPtr should return true for arena pointers")
	}

	slice := MakeSlice[int](a, 5, 10)
	if !OwnsSlice(a, slice) {
		t.Error("OwnsSlice should return true for arena slice pointers")
	}

	str := a.MakeString("test")
	if !OwnsString(a, str) {
		t.Error("OwnsString should return true for arena string pointers")
	}

	// Test with nil pointer
	invalidPtr := unsafe.Pointer(nil)
	if a.Owns(invalidPtr) {
		t.Error("Owns should return false for nil pointer")
	}
}

func TestPtr(t *testing.T) {
	a := New(1024, BUMP)
	defer a.Delete()

	// Test with int
	value := 42
	ptr := Ptr(a, value)
	if ptr == nil {
		t.Fatal("Ptr returned nil")
	}
	if *ptr != 42 {
		t.Errorf("Expected *ptr = 42, got %d", *ptr)
	}
	if !OwnsPtr(a, ptr) {
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
	personPtr := Ptr(a, person)
	if personPtr == nil {
		t.Fatal("Ptr returned nil for struct")
	}
	if personPtr.Name != "Alice" || personPtr.Age != 30 {
		t.Errorf("Expected Person{Alice, 30}, got %+v", *personPtr)
	}
	if !OwnsPtr(a, personPtr) {
		t.Error("Ptr should allocate struct in arena")
	}

	// Test with string
	str := "hello"
	strPtr := Ptr(a, str)
	if *strPtr != "hello" {
		t.Errorf("Expected string 'hello', got '%s'", *strPtr)
	}

	// Test with slice (copies the slice header, not the backing array)
	slice := []int{1, 2, 3}
	slicePtr := Ptr(a, slice)
	if len(*slicePtr) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(*slicePtr))
	}
	if (*slicePtr)[0] != 1 || (*slicePtr)[1] != 2 || (*slicePtr)[2] != 3 {
		t.Errorf("Expected slice [1 2 3], got %v", *slicePtr)
	}
}
