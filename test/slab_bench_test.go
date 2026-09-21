package arena_test

import (
	"testing"
	"unsafe"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
)

// Benchmark basic allocation performance
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_Alloc(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[int](a)
	}
}

// Benchmark allocation with reset
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocWithReset(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[int](a)
		if i%1000 == 0 {
			a.Reset()
		}
	}
}

// Benchmark allocation and deallocation
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocFree(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		p := arena.Alloc[int](a)
		arena.DeleteObject(a, p)
	}
}

// Benchmark small object allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocSmall(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		p := arena.Alloc[byte](a)
		_ = p
	}
}

// Benchmark medium object allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocMedium(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type MediumStruct struct {
		A [32]int64
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[MediumStruct](a)
	}
}

// Benchmark large object allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocLarge(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type LargeStruct struct {
		A [512]int64
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[LargeStruct](a)
	}
}

// Benchmark xlarge object allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocXLarge(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type XLargeStruct struct {
		A [4096]int64
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[XLargeStruct](a)
	}
}

// Benchmark mixed size allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocMixed(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		switch i % 4 {
		case 0:
			_ = arena.Alloc[int](a)
		case 1:
			_ = arena.Alloc[int32](a)
		case 2:
			_ = arena.Alloc[int64](a)
		case 3:
			_ = arena.Alloc[float64](a)
		}
	}
}

// Benchmark many small allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocMany(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		for j := 0; j < 100; j++ {
			_ = arena.Alloc[int](a)
		}
	}
}

// Benchmark alloc/free pattern
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocFreePattern(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	ptrs := make([]*int, 100)

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		// Allocate 100 objects
		for j := 0; j < 100; j++ {
			ptrs[j] = arena.Alloc[int](a)
		}
		// Free odd indices
		for j := 1; j < 100; j += 2 {
			arena.DeleteObject(a, ptrs[j])
		}
		// Free even indices
		for j := 0; j < 100; j += 2 {
			arena.DeleteObject(a, ptrs[j])
		}
	}
}

// Benchmark sequential large allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_SequentialLarge(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type LargeType struct {
		A [256]int64
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		for j := 0; j < 10; j++ {
			_ = arena.Alloc[LargeType](a)
		}
	}
}

// Benchmark reset overhead
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_Reset(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Pre-allocate some objects
	for i := 0; i < 1000; i++ {
		_ = arena.Alloc[int](a)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		a.Reset()
		// Re-allocate to keep state consistent
		for j := 0; j < 1000; j++ {
			_ = arena.Alloc[int](a)
		}
	}
}

// Benchmark owns check
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_Owns(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	p := arena.Alloc[int](a)

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = a.Owns(unsafe.Pointer(p))
	}
}

// Benchmark struct allocation
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocStruct(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type TestStruct struct {
		A int64
		B string
		C float64
		D [16]byte
	}

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		p := arena.Alloc[TestStruct](a)
		_ = p
	}
}

// Benchmark array allocation
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_AllocArray(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[[100]int](a)
	}
}

// Benchmark mixed type allocations
//
// Revisions:
//   - 2025-12-19 13:53: initial creation
func BenchmarkSlab_MixedTypes(b *testing.B) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[int](a)
		_ = arena.Alloc[string](a)
		_ = arena.Alloc[float64](a)
		_ = arena.Alloc[bool](a)
	}
}
