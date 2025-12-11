package arena_test

import (
	"runtime"
	"testing"

	"github.com/thebagchi/arena-go"
)

func BenchmarkBumpAlloc(b *testing.B) {
	a := arena.New(100, arena.BUMP)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := arena.Alloc[int](a)
		*ptr = i
	}
}

func BenchmarkBumpMakeSlice(b *testing.B) {
	a := arena.New(1000, arena.BUMP)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slice := arena.MakeSlice[int](a, 10, 10)
		slice[0] = i
	}
}

func BenchmarkBumpMakeObject(b *testing.B) {
	a := arena.New(100, arena.BUMP)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := arena.MakeObject[struct{ A, B int }](a)
		obj.A = i
		obj.B = i * 2
	}
}

func BenchmarkBumpMakeString(b *testing.B) {
	a := arena.New(1000, arena.BUMP)
	str := "benchmark string for allocation"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := a.MakeString(str)
		_ = s
	}
}

func BenchmarkHeapAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ptr := new(int)
		*ptr = i
	}
}

func BenchmarkHeapMakeSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := make([]int, 10)
		slice[0] = i
	}
}

func BenchmarkHeapMakeString(b *testing.B) {
	str := "benchmark string for allocation"
	for i := 0; i < b.N; i++ {
		s := string([]byte(str))
		_ = s
	}
}

func BenchmarkBumpAllocatorStress(b *testing.B) {
	a := arena.New(100, arena.BUMP)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate TestBumpAllocatorVariousSizes logic
		intPtr := arena.Alloc[int](a)
		*intPtr = 42

		int64Ptr := arena.Alloc[int64](a)
		*int64Ptr = 123456789

		float64Ptr := arena.Alloc[float64](a)
		*float64Ptr = 3.14159

		type TestStruct struct {
			A int
			B string
			C [10]int
		}
		structPtr := arena.Alloc[TestStruct](a)
		structPtr.A = 1
		structPtr.B = "test"

		slice1 := arena.MakeSlice[int](a, 5, 10)
		for j := range slice1 {
			slice1[j] = j
		}

		slice2 := arena.MakeSlice[string](a, 3, 5)
		slice2[0] = "hello"
		slice2[1] = "world"

		largeSlice := arena.MakeSlice[byte](a, 100, 200)
		for j := range largeSlice {
			largeSlice[j] = byte(j % 256)
		}

		str := a.MakeString("various sizes test")
		_ = str

		// Reset after each iteration
		a.Reset()
	}
}

func BenchmarkHeapAllocatorStress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Simulate similar allocations on heap
		intPtr := new(int)
		*intPtr = 42

		int64Ptr := new(int64)
		*int64Ptr = 123456789

		float64Ptr := new(float64)
		*float64Ptr = 3.14159

		type TestStruct struct {
			A int
			B string
			C [10]int
		}
		structPtr := new(TestStruct)
		structPtr.A = 1
		structPtr.B = "test"

		slice1 := make([]int, 5, 10)
		for j := range slice1 {
			slice1[j] = j
		}

		slice2 := make([]string, 3, 5)
		slice2[0] = "hello"
		slice2[1] = "world"

		largeSlice := make([]byte, 100, 200)
		for j := range largeSlice {
			largeSlice[j] = byte(j % 256)
		}

		str := string([]byte("various sizes test"))
		_ = str

		// Run GC after each iteration
		runtime.GC()
	}
}
