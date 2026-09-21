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
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_Alloc(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[int](a)
	}
}

// Benchmark allocation with reset
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocWithReset(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
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
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocFree(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		p := arena.Alloc[int](a)
		arena.DeleteObject(a, p)
	}
}

// Benchmark large allocations
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocLarge(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// A named array rather than a struct wrapping one: the size is the point,
	// and a field nothing reads is a field nothing can justify.
	type Large [LARGE_OBJECT_BYTES]byte

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = arena.Alloc[Large](a)
	}
}

// Benchmark mixed workload
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_MixedWorkload(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	ptrs := make([]*int, 100)

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		var idx = i % len(ptrs)
		if ptrs[idx] != nil {
			arena.DeleteObject(a, ptrs[idx])
		}

		ptrs[idx] = arena.Alloc[int](a)
		*ptrs[idx] = i
	}
}

// Benchmark Owns operation
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_Owns(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	p := arena.Alloc[int](a)
	*p = 42

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		_ = a.Owns(unsafe.Pointer(p))
	}
}

// Benchmark parallel allocations
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocParallel(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := arena.Alloc[int](a)
			*p = 42
		}
	})
}

// Benchmark parallel alloc/free
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocFreeParallel(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := arena.Alloc[int](a)
			*p = 42
			arena.DeleteObject(a, p)
		}
	})
}

// Benchmark various allocation sizes
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_AllocSizes(b *testing.B) {
	var sizes = []struct {
		name string
		size int
	}{
		{"8B", 8},
		{"64B", 64},
		{"256B", 256},
		{"1KB", 1024},
		{"4KB", 4096},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			a := arena.New(alloc.NewBuddyAllocator())
			defer a.Delete()

			b.ResetTimer()

			for i := 0; i < b.N; i = i + 1 {
				_ = a.Alloc(uint64(s.size), 8)
			}
		})
	}
}

// Benchmark coalescing behavior
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_Coalescing(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		// Allocate 4 adjacent blocks
		p1 := arena.Alloc[int](a)
		p2 := arena.Alloc[int](a)
		p3 := arena.Alloc[int](a)
		p4 := arena.Alloc[int](a)

		*p1 = 1
		*p2 = 2
		*p3 = 3
		*p4 = 4

		// Free them - should trigger coalescing
		arena.DeleteObject(a, p1)
		arena.DeleteObject(a, p2)
		arena.DeleteObject(a, p3)
		arena.DeleteObject(a, p4)
	}
}

// Benchmark comparison with standard allocator
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_VsStandard_Alloc(b *testing.B) {
	b.Run("Buddy", func(b *testing.B) {
		a := arena.New(alloc.NewBuddyAllocator())
		defer a.Delete()

		b.ResetTimer()

		for i := 0; i < b.N; i = i + 1 {
			_ = arena.Alloc[int](a)
		}
	})

	b.Run("Standard", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i = i + 1 {
			_ = new(int)
		}
	})
}

// Benchmark comparison with BUMP allocator
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_VsBump_Alloc(b *testing.B) {
	b.Run("Buddy", func(b *testing.B) {
		a := arena.New(alloc.NewBuddyAllocator())
		defer a.Delete()

		b.ResetTimer()

		for i := 0; i < b.N; i = i + 1 {
			_ = arena.Alloc[int](a)
		}
	})

	b.Run("Bump", func(b *testing.B) {
		a := arena.New(alloc.NewBumpAllocator(100 * 4096))
		defer a.Delete()

		b.ResetTimer()

		for i := 0; i < b.N; i = i + 1 {
			_ = arena.Alloc[int](a)
		}
	})
}

// Benchmark fragmentation scenario
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_Fragmentation(b *testing.B) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 100

	var ptrs = make([]*int, count)

	b.ResetTimer()

	for i := 0; i < b.N; i = i + 1 {
		// Allocate
		for j := 0; j < count; j = j + 1 {
			ptrs[j] = arena.Alloc[int](a)
			*ptrs[j] = j
		}

		// Free every other allocation
		for j := 0; j < count; j = j + 2 {
			arena.DeleteObject(a, ptrs[j])
		}

		// Free remaining
		for j := 1; j < count; j = j + 2 {
			arena.DeleteObject(a, ptrs[j])
		}
	}
}

// Benchmark batch allocations
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkBuddy_BatchAlloc(b *testing.B) {
	var sizes = []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			a := arena.New(alloc.NewBuddyAllocator())
			defer a.Delete()

			b.ResetTimer()

			for i := 0; i < b.N; i = i + 1 {
				for j := 0; j < size; j = j + 1 {
					_ = arena.Alloc[int](a)
				}

				a.Reset()
			}
		})
	}
}
