package arena_test

import (
	"testing"
	"unsafe"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
)

const (
	// LARGE_OBJECT_BYTES is the size the large-object benchmarks allocate.
	LARGE_OBJECT_BYTES = 1024
)

// Benchmark basic allocation performance
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_Alloc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = arena.Alloc[int](a)
	}
}

// Benchmark allocation with reset
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_AllocWithReset(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = arena.Alloc[int](a)
		if i%1000 == 0 {
			a.Reset()
		}
	}
}

// Benchmark allocation and deallocation (bump doesn't truly free, but test for consistency)
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_AllocFree(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := arena.Alloc[int](a)
		arena.DeleteObject(a, p)
	}
}

// Benchmark large allocations
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_AllocLarge(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	// A named array rather than a struct wrapping one: the size is the point,
	// and a field nothing reads is a field nothing can justify.
	type Large [LARGE_OBJECT_BYTES]byte

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = arena.Alloc[Large](a)
	}
}

// Benchmark mixed workload
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_MixedWorkload(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = arena.Alloc[int](a)
		_ = arena.Alloc[int64](a)
		_ = arena.Alloc[float64](a)
	}
}

// Benchmark Owns operation
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_Owns(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	p := arena.Alloc[int](a)
	*p = 42

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = a.Owns(unsafe.Pointer(p))
	}
}

// Benchmark make slice
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_MakeSlice(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1000 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		slice := arena.MakeSlice[int](a, 10, 10)
		slice[0] = i
	}
}

// Benchmark make string
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_MakeString(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1000 * 4096))
	defer a.Delete()

	str := "benchmark string for allocation"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := a.MakeString(str)
		_ = s
	}
}

// Benchmark comparison with standard allocator
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_VsStandard_Alloc(b *testing.B) {
	b.Run("Bump", func(b *testing.B) {
		a := arena.New(alloc.NewBumpAllocator(100 * 4096))
		defer a.Delete()

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = arena.Alloc[int](a)
		}
	})

	b.Run("Standard", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = new(int)
		}
	})
}

// Benchmark batch allocations
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_BatchAlloc(b *testing.B) {
	var sizes = []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			a := arena.New(alloc.NewBumpAllocator(1000 * 4096))
			defer a.Delete()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := 0; j < size; j++ {
					_ = arena.Alloc[int](a)
				}

				a.Reset()
			}
		})
	}
}

// Benchmark reset overhead
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_Reset(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	// Pre-allocate some objects
	for i := 0; i < 1000; i++ {
		_ = arena.Alloc[int](a)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.Reset()
		// Re-allocate to keep state consistent
		for j := 0; j < 1000; j++ {
			_ = arena.Alloc[int](a)
		}
	}
}

// Benchmark struct allocation
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_AllocStruct(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	type TestStruct struct {
		A int
		B string
		C [10]int
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := arena.Alloc[TestStruct](a)
		_ = p
	}
}

// Benchmark array allocation
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_AllocArray(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = arena.Alloc[[100]int](a)
	}
}

// Benchmark stress test
//
// Revisions:
//   - 2025-12-31 16:13: initial creation
func BenchmarkBump_StressTest(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1000 * 4096))
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_ = arena.Alloc[int](a)
		}

		a.Reset()
	}
}
