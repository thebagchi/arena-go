package arena

import "testing"

func TestBumpAllocator(t *testing.T) {
	a := New(1, BUMP)
	p1 := Alloc[int](a)
	if p1 == nil {
		t.Fatal("alloc failed")
	}
	p2 := Alloc[int](a)
	if p2 == nil {
		t.Fatal("alloc failed")
	}
	if p1 == p2 {
		t.Fatal("same pointer")
	}
	a.Reset()
	p3 := Alloc[int](a)
	if p3 != p1 {
		t.Fatal("not reset")
	}
}

func TestBumpAllocatorVariousSizes(t *testing.T) {
	a := New(10, BUMP) // 10 pages for larger allocations

	// Test allocating different basic types
	intPtr := Alloc[int](a)
	if intPtr == nil {
		t.Fatal("failed to alloc int")
	}
	*intPtr = 42

	int64Ptr := Alloc[int64](a)
	if int64Ptr == nil {
		t.Fatal("failed to alloc int64")
	}
	*int64Ptr = 123456789

	float64Ptr := Alloc[float64](a)
	if float64Ptr == nil {
		t.Fatal("failed to alloc float64")
	}
	*float64Ptr = 3.14159

	// Test allocating a struct
	type TestStruct struct {
		A int
		B string
		C [10]int
	}
	structPtr := Alloc[TestStruct](a)
	if structPtr == nil {
		t.Fatal("failed to alloc struct")
	}
	structPtr.A = 1
	structPtr.B = "test"

	// Test allocating slices of different sizes
	slice1 := MakeSlice[int](a, 5, 10)
	if len(slice1) != 5 || cap(slice1) != 10 {
		t.Fatalf("slice1: len=%d cap=%d, expected len=5 cap=10", len(slice1), cap(slice1))
	}
	for i := range slice1 {
		slice1[i] = i
	}

	slice2 := MakeSlice[string](a, 3, 5)
	if len(slice2) != 3 || cap(slice2) != 5 {
		t.Fatalf("slice2: len=%d cap=%d, expected len=3 cap=5", len(slice2), cap(slice2))
	}
	slice2[0] = "hello"
	slice2[1] = "world"

	// Test allocating a larger slice
	largeSlice := MakeSlice[byte](a, 1000, 2000)
	if len(largeSlice) != 1000 || cap(largeSlice) != 2000 {
		t.Fatalf("largeSlice: len=%d cap=%d, expected len=1000 cap=2000", len(largeSlice), cap(largeSlice))
	}
	for i := range largeSlice {
		largeSlice[i] = byte(i % 256)
	}

	// Test string allocation
	str := a.MakeString("various sizes test")
	if str != "various sizes test" {
		t.Fatalf("string alloc failed: got %s", str)
	}

	// Verify values are set correctly
	if *intPtr != 42 {
		t.Fatalf("intPtr value: got %d, expected 42", *intPtr)
	}
	if *int64Ptr != 123456789 {
		t.Fatalf("int64Ptr value: got %d, expected 123456789", *int64Ptr)
	}
	if *float64Ptr != 3.14159 {
		t.Fatalf("float64Ptr value: got %f, expected 3.14159", *float64Ptr)
	}
	if structPtr.A != 1 || structPtr.B != "test" {
		t.Fatalf("structPtr values: A=%d B=%s, expected A=1 B=test", structPtr.A, structPtr.B)
	}
	for i, v := range slice1 {
		if v != i {
			t.Fatalf("slice1[%d]: got %d, expected %d", i, v, i)
		}
	}
	if slice2[0] != "hello" || slice2[1] != "world" {
		t.Fatalf("slice2 values: %v, expected [hello world]", slice2)
	}
}

func TestBumpAllocatorGrow(t *testing.T) {
	// Create a small arena with only 1 page (typically 4096 bytes)
	a := New(1, BUMP)

	// Allocate a large slice that exceeds the initial arena size
	// 1000 * 8 bytes (sizeof(int)) = 8000 bytes > 4096 bytes, forcing growth
	largeSlice := MakeSlice[int](a, 1000, 1000)
	if len(largeSlice) != 1000 || cap(largeSlice) != 1000 {
		t.Fatalf("largeSlice: len=%d cap=%d, expected 1000", len(largeSlice), cap(largeSlice))
	}

	// Fill the slice
	for i := range largeSlice {
		largeSlice[i] = i * 2
	}

	// Allocate another item to ensure growth worked
	anotherPtr := Alloc[int64](a)
	if anotherPtr == nil {
		t.Fatal("failed to alloc after growth")
	}
	*anotherPtr = 999999

	// Verify the large slice values
	for i, v := range largeSlice {
		expected := i * 2
		if v != expected {
			t.Fatalf("largeSlice[%d]: got %d, expected %d", i, v, expected)
		}
	}

	// Verify the additional allocation
	if *anotherPtr != 999999 {
		t.Fatalf("anotherPtr: got %d, expected 999999", *anotherPtr)
	}
}
