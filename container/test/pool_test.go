package container_test

import (
	"sync"
	"testing"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// TestPool_Basic covers pool basic.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_Basic(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Allocate first object
	x := p.Alloc()
	if x == nil {
		t.Fatal("alloc failed")
	}

	*x = 42

	// Free it
	p.Free(x)

	if p.Len() != 1 {
		t.Errorf("Expected free list len 1, got %d", p.Len())
	}

	// Allocate again - should reuse
	y := p.Alloc()
	if y == nil {
		t.Fatal("alloc failed")
	}

	if *y != 0 {
		t.Errorf("Expected zeroed value, got %d", *y)
	}

	if y != x {
		t.Error("Expected to reuse same pointer")
	}

	// Free list should be empty now
	if p.Len() != 0 {
		t.Errorf("Expected free list len 0, got %d", p.Len())
	}
}

// TestPool_MultipleObjects covers pool multiple objects.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_MultipleObjects(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Allocate multiple objects
	ptrs := make([]*int, 100)
	for i := range ptrs {
		ptrs[i] = p.Alloc()
		*ptrs[i] = i
	}

	// Free them all
	for _, ptr := range ptrs {
		p.Free(ptr)
	}

	if p.Len() != 100 {
		t.Errorf("Expected free list len 100, got %d", p.Len())
	}

	// Reuse them all
	for i := range ptrs {
		ptr := p.Alloc()
		if *ptr != 0 {
			t.Errorf("Expected zeroed value at %d, got %d", i, *ptr)
		}
	}

	if p.Len() != 0 {
		t.Errorf("Expected free list len 0, got %d", p.Len())
	}
}

// TestPool_Reset covers pool reset.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_Reset(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Allocate all objects first, then free them
	ptrs := make([]*int, 10)
	for i := 0; i < 10; i++ {
		ptrs[i] = p.Alloc()
		*ptrs[i] = i
	}

	for _, ptr := range ptrs {
		p.Free(ptr)
	}

	if p.Len() != 10 {
		t.Errorf("Expected free list len 10, got %d", p.Len())
	}

	// Reset should clear free list
	p.Reset()

	if p.Len() != 0 {
		t.Errorf("Expected free list len 0 after reset, got %d", p.Len())
	}

	// Allocating after reset should get new memory
	ptr := p.Alloc()
	if ptr == nil {
		t.Fatal("alloc failed after reset")
	}
}

// TestPool_FreeNil covers pool free nil.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_FreeNil(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Should not panic
	p.Free(nil)

	if p.Len() != 0 {
		t.Errorf("Expected free list len 0, got %d", p.Len())
	}
}

// TestPool_StructTypes covers pool struct types.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_StructTypes(t *testing.T) {
	type Node struct {
		Value int
		Left  *Node
		Right *Node
	}

	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	defer a.Delete()

	p := container.NewPool[Node](a)

	// Allocate a node
	n1 := p.Alloc()
	n1.Value = 42

	// Free and reallocate
	p.Free(n1)
	n2 := p.Alloc()

	// Should be same pointer and zeroed
	if n2 != n1 {
		t.Error("Expected to reuse same pointer")
	}

	if n2.Value != 0 || n2.Left != nil || n2.Right != nil {
		t.Error("Expected zeroed struct")
	}
}

// TestPool_ThreadSafety covers pool thread safety.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_ThreadSafety(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	var wg sync.WaitGroup

	workers := 10
	iterations := 1000

	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				ptr := p.Alloc()
				if ptr == nil {
					t.Errorf("alloc failed")
					return
				}

				*ptr = i
				p.Free(ptr)
			}
		}()
	}

	wg.Wait()

	// All objects should be in free list
	freeCount := p.Len()
	if freeCount > workers*iterations {
		t.Errorf(
			"Free list too large: %d (max expected %d)",
			freeCount,
			workers*iterations,
		)
	}
}

// TestPool_ArenaLifecycle covers pool arena lifecycle.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_ArenaLifecycle(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096))
	p := container.NewPool[int](a)

	// Allocate some objects
	ptrs := make([]*int, 10)
	for i := range ptrs {
		ptrs[i] = p.Alloc()
		*ptrs[i] = i
	}

	// Delete arena - this should free all pool memory
	a.Delete()

	// Pool should still be safe to use for free list operations
	p.Reset()

	if p.Len() != 0 {
		t.Errorf("Expected free list len 0, got %d", p.Len())
	}
}

// TestPool_DifferentAllocators covers pool different allocators.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_DifferentAllocators(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(100 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Allocate, use, free, reuse
	x := p.Alloc()
	*x = 42
	p.Free(x)

	y := p.Alloc()
	if *y != 0 {
		t.Errorf("Expected zeroed value, got %d", *y)
	}

	if y != x {
		t.Error("Expected to reuse same pointer")
	}
}

// TestPool_LargeStructs covers pool large structs.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func TestPool_LargeStructs(t *testing.T) {
	type LargeStruct struct {
		Data [1024]byte
		ID   int
	}

	a := arena.New(alloc.NewBumpAllocator(1024 * 4096))
	defer a.Delete()

	p := container.NewPool[LargeStruct](a)

	// Allocate and verify zeroing
	s := p.Alloc()

	s.ID = 123
	for i := range s.Data {
		s.Data[i] = byte(i)
	}

	// Free and reallocate
	p.Free(s)
	s2 := p.Alloc()

	// Should be zeroed
	if s2.ID != 0 {
		t.Errorf("Expected zeroed ID, got %d", s2.ID)
	}

	for i := range s2.Data {
		if s2.Data[i] != 0 {
			t.Errorf("Expected zeroed data at %d, got %d", i, s2.Data[i])
		}
	}
}

// Benchmarks
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkPool_Alloc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(10240 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = p.Alloc()
	}
}

// BenchmarkPool_AllocFree measures pool alloc free.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkPool_AllocFree(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(10240 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr := p.Alloc()
		p.Free(ptr)
	}
}

// BenchmarkPool_AllocFreeReuse measures pool alloc free reuse.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkPool_AllocFreeReuse(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(10240 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	// Pre-populate free list
	for i := 0; i < 100; i++ {
		p.Free(p.Alloc())
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr := p.Alloc()
		p.Free(ptr)
	}
}

// BenchmarkPool_Parallel measures pool parallel.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkPool_Parallel(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(10240 * 4096))
	defer a.Delete()

	p := container.NewPool[int](a)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ptr := p.Alloc()
			p.Free(ptr)
		}
	})
}

// BenchmarkStdAlloc_AllocFree measures std alloc alloc free.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkStdAlloc_AllocFree(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr := new(int)
		_ = ptr
	}
}

type BenchNode struct {
	Value int
	Left  *BenchNode
	Right *BenchNode
}

// BenchmarkPool_Struct measures pool struct.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkPool_Struct(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(10240 * 4096))
	defer a.Delete()

	p := container.NewPool[BenchNode](a)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr := p.Alloc()
		p.Free(ptr)
	}
}

// BenchmarkStdAlloc_Struct measures std alloc struct.
//
// Revisions:
//   - 2025-12-16 08:35: initial creation
func BenchmarkStdAlloc_Struct(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ptr := new(BenchNode)
		_ = ptr
	}
}
