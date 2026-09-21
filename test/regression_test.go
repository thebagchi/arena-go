// Regression tests for the defects found reviewing this package.
//
// Each one failed before the fix it guards. They were written as reproductions
// first, behind a build tag, and moved here as each defect was closed. Each names
// the defect it covers, so a test that starts failing says what has come back.
package arena_test

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
	"unsafe"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
	arenaio "github.com/thebagchi/arena-go/io"
)

const (
	// DEADLOCK_TIMEOUT is how long the concurrency tests wait before calling a
	// stall a deadlock. A working allocator finishes these in under a second.
	DEADLOCK_TIMEOUT = 20 * time.Second
	// SENTINEL is a recognisable value written to memory that is then freed, so
	// that a later allocation handing it back unzeroed is obvious.
	SENTINEL = 0xDEADBEEF
	// MIN_BLOCK_BYTES is the buddy allocator's smallest block, used to size an
	// arena that must hold a known number of allocations without growing.
	MIN_BLOCK_BYTES = 16
)

// TestSlabConcurrentAllocFree covers a lock-order inversion: allocation took the bin lock
// then the allocator lock, and freeing took them the other way round, so the two
// deadlocked as soon as they ran together.
//
// Revisions:
//   - 2026-09-21 21:05: initial creation
func TestSlabConcurrentAllocFree(t *testing.T) {
	slab := alloc.NewSlabAllocator()
	defer slab.Delete()

	done := make(chan struct{})

	go func() {
		defer close(done)

		var wg sync.WaitGroup

		for range 8 {
			wg.Add(1)

			go func() {
				defer wg.Done()

				held := make([]unsafe.Pointer, 0, 512)

				for range 20000 {
					held = append(held, slab.Alloc(24, 8))
					if len(held) < cap(held) {
						continue
					}

					for _, ptr := range held {
						slab.Remove(ptr)
					}

					held = held[:0]
				}
			}()
		}

		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(DEADLOCK_TIMEOUT):
		t.Fatal("Alloc and Remove deadlocked: the two lock orders disagree")
	}
}

// TestSlabConcurrentClassesRace covers a data race: freeing appended to the
// shared page pool while holding only a bin lock, which raced with allocation
// reading that pool under the allocator lock. Run under -race to see it.
//
// Revisions:
//   - 2026-09-21 21:10: initial creation
func TestSlabConcurrentClassesRace(t *testing.T) {
	slab := alloc.NewSlabAllocator()
	defer slab.Delete()

	var wg sync.WaitGroup

	for _, size := range []uint64{64, 128, 256, 512} {
		wg.Add(1)

		go func(size uint64) {
			defer wg.Done()

			held := make([]unsafe.Pointer, 200)

			for range 50 {
				for i := range held {
					held[i] = slab.Alloc(size, 8)
				}

				for _, ptr := range held {
					slab.Remove(ptr)
				}
			}
		}(size)
	}

	wg.Wait()
}

// TestSlabRejectsDoubleFree covers a double free: a second free of the same
// pointer linked a slot to itself, so the next two allocations returned the same
// address.
//
// Revisions:
//   - 2026-09-21 21:14: initial creation
func TestSlabRejectsDoubleFree(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	keep := a.Alloc(64, 8)
	freed := a.Alloc(64, 8)

	a.Remove(freed)
	a.Remove(freed)

	first, second := a.Alloc(64, 8), a.Alloc(64, 8)
	if first == second {
		t.Errorf("a double free handed address %p out twice", first)
	}

	if first == keep || second == keep {
		t.Error("an allocation returned a block that is still live")
	}
}

// TestSlabRejectsInteriorPointer covers an interior pointer: a pointer into the middle
// of an object was accepted and threaded the free list through live data.
//
// Revisions:
//   - 2026-09-21 21:18: initial creation
func TestSlabRejectsInteriorPointer(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const offset = 8

	held := a.Alloc(64, 8)
	a.Remove(unsafe.Add(held, offset))

	next := a.Alloc(64, 8)
	if next == held {
		t.Error("freeing an interior pointer released the whole object")
	}
}

// TestAllocAfterDeletePanics covers use after delete: one allocator panicked with
// an index out of range, another with a nil dereference, and the third with a
// message of its own. All three now report the same thing.
//
// Revisions:
//   - 2026-09-21 21:22: initial creation
func TestAllocAfterDeletePanics(t *testing.T) {
	cases := []struct {
		name string
		make func() arena.Allocator
	}{
		{
			name: "bump",
			make: func() arena.Allocator { return alloc.NewBumpAllocator(4096) },
		},
		{name: "slab", make: func() arena.Allocator { return alloc.NewSlabAllocator() }},
		{name: "buddy", make: func() arena.Allocator { return alloc.NewBuddyAllocator() }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := arena.New(tc.make())
			a.Delete()

			defer func() {
				recovered, _ := recover().(error)
				if !errors.Is(recovered, arena.ErrDeleted) {
					t.Errorf(
						"Alloc after Delete panicked with %v, want %v",
						recovered,
						arena.ErrDeleted,
					)
				}
			}()

			_ = a.Alloc(8, 8)
		})
	}
}

// TestBuddyFixedExhaustion covers exhaustion: a full fixed-size allocator
// returned nil, and every generic helper dereferenced it. Under -race the crash
// was a fatal error that no recover could catch.
//
// Revisions:
//   - 2026-09-21 21:27: initial creation
func TestBuddyFixedExhaustion(t *testing.T) {
	const oversized = 1 << 20

	a := arena.New(alloc.NewBuddyAllocator(alloc.WithGrowthStrategy(alloc.FIXED)))
	defer a.Delete()

	_ = a.Alloc(16, 8)

	if _, ok := a.TryAlloc(oversized, 16); ok {
		t.Fatal("TryAlloc succeeded on a full fixed allocator")
	}

	defer func() {
		recovered, _ := recover().(error)
		if !errors.Is(recovered, arena.ErrOutOfMemory) {
			t.Errorf("MakeSlice on a full arena panicked with %v, want %v",
				recovered, arena.ErrOutOfMemory)
		}
	}()

	_ = arena.MakeSlice[byte](a, oversized, oversized)
}

// TestSlabResetReusesPages covers a mapping leak: Reset left every slab page
// mapped and unreferenced, so each cycle mapped a fresh set.
//
// Revisions:
//   - 2026-09-21 21:32: initial creation
func TestSlabResetReusesPages(t *testing.T) {
	slab := alloc.NewSlabAllocator()
	defer slab.Delete()

	mapped := make([]int, 0, 5)

	for range cap(mapped) {
		for range 10000 {
			_ = slab.Alloc(64, 8)
		}

		slab.Reset()

		mapped = append(mapped, slab.Stats().Mapped)
	}

	first, last := mapped[0], mapped[len(mapped)-1]
	if last > first {
		t.Errorf("mapped bytes grew across Reset cycles: %v", mapped)
	}
}

// TestAllocIsZeroed covers unzeroed memory: no allocator zeroed anything, so a
// slab handed back its own free-list pointer in the first word of every fresh
// object, and bump and buddy handed back the previous occupant's data.
//
// Revisions:
//   - 2026-09-21 21:37: initial creation
func TestAllocIsZeroed(t *testing.T) {
	cases := []struct {
		name  string
		make  func() arena.Allocator
		reuse func(*arena.Arena, *int64)
	}{
		{
			name:  "slab reuses a freed object",
			make:  func() arena.Allocator { return alloc.NewSlabAllocator() },
			reuse: func(a *arena.Arena, ptr *int64) { arena.DeleteObject(a, ptr) },
		},
		{
			name:  "buddy reuses a freed block",
			make:  func() arena.Allocator { return alloc.NewBuddyAllocator() },
			reuse: func(a *arena.Arena, ptr *int64) { arena.DeleteObject(a, ptr) },
		},
		{
			name:  "bump reuses the whole arena",
			make:  func() arena.Allocator { return alloc.NewBumpAllocator(4096) },
			reuse: func(a *arena.Arena, _ *int64) { a.Reset() },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := arena.New(tc.make())
			defer a.Delete()

			first := arena.Alloc[int64](a)
			if *first != 0 {
				t.Errorf("a fresh allocation held %#x, want 0", *first)
			}

			*first = SENTINEL
			tc.reuse(a, first)

			if again := arena.Alloc[int64](a); *again != 0 {
				t.Errorf("a reused allocation held %#x, want 0", *again)
			}

			for i, value := range arena.MakeSlice[int64](a, 4, 4) {
				if value != 0 {
					t.Errorf("MakeSlice element %d held %#x, want 0", i, value)
				}
			}
		})
	}
}

// TestPoolAllocIsZeroed covers unzeroed memory for the pool, whose objects never
// go back to the allocator and so are zeroed by the pool itself.
//
// Revisions:
//   - 2026-09-21 21:44: initial creation
func TestPoolAllocIsZeroed(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	pool := container.NewPool[[4]int64](a)

	first := pool.Alloc()
	if *first != ([4]int64{}) {
		t.Errorf("a fresh pool object held %v, want zeroes", *first)
	}

	*first = [4]int64{SENTINEL, SENTINEL, SENTINEL, SENTINEL}
	pool.Free(first)

	if again := pool.Alloc(); *again != ([4]int64{}) {
		t.Errorf("a reused pool object held %v, want zeroes", *again)
	}
}

// TestBuddyHonoursAlign covers a broken alignment: the buddy allocator worked out the
// block size the alignment needed and then asked for a block sized only by the
// request, so Alloc(8, 64) came back 16-byte aligned.
//
// Revisions:
//   - 2026-09-21 21:48: initial creation
func TestBuddyHonoursAlign(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	for _, align := range []uint64{16, 32, 64, 128, 256} {
		t.Run(fmt.Sprintf("align %d", align), func(t *testing.T) {
			// An odd block first, so the next one cannot be aligned by accident.
			_ = a.Alloc(16, 16)

			ptr := a.Alloc(8, align)
			if uintptr(ptr)%uintptr(align) != 0 {
				t.Errorf(
					"Alloc(8, %d) returned %p, which is not %d-aligned",
					align,
					ptr,
					align,
				)
			}
		})
	}
}

// TestMapStructKeyWithStringField covers a wrong hash: keys were hashed from
// their raw memory, so two equal keys whose string fields sat at different
// addresses hashed differently and the second one was never found.
//
// Revisions:
//   - 2026-09-21 21:53: initial creation
func TestMapStructKeyWithStringField(t *testing.T) {
	type key struct {
		Name string
		Rank int
	}

	a := arena.New(alloc.NewBumpAllocator(1 << 20))
	defer a.Delete()

	m := container.NewMap[key, int](a)

	stored := key{Name: string([]byte("hello")), Rank: 1}
	lookup := key{Name: string([]byte("hello")), Rank: 1}

	if stored != lookup {
		t.Fatal("the two keys must compare equal for this test to mean anything")
	}

	m.Set(stored, 42)

	if got, found := m.Get(lookup); !found || got != 42 {
		t.Errorf("Get with an equal key returned (%d, %t), want (42, true)", got, found)
	}
}

// TestArenaMemoryIsNotGCRoot covers the collector hazard: arena memory is not scanned
// by the garbage collector, so heap strings kept only by a map entry were
// collected and their bytes reused. Map now copies string keys into the arena.
//
// Revisions:
//   - 2026-09-21 21:58: initial creation
func TestArenaMemoryIsNotGCRoot(t *testing.T) {
	const entries = 20000

	a := arena.New(alloc.NewBumpAllocator(1 << 24))
	defer a.Delete()

	m := container.NewMap[string, int](a)
	key := func(i int) string {
		return fmt.Sprintf("person-%06d-%s", i, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}

	for i := range entries {
		m.Set(key(i), i)
	}

	_ChurnHeap(entries)

	missed := 0

	for i := range entries {
		if value, found := m.Get(key(i)); !found || value != i {
			missed = missed + 1
		}
	}

	if missed > 0 {
		t.Errorf("%d of %d keys did not survive a garbage collection", missed, entries)
	}
}

// TestBuddySequentialAllocScales covers a scaling defect: the one-bit-per-node
// bitmap could not tell a full subtree from a partly used one, so the search
// re-walked every live block and allocation cost grew with the live set.
//
// Four times the work should take about four times as long. It took sixteen.
//
// Revisions:
//   - 2026-09-21 22:04: initial creation
func TestBuddySequentialAllocScales(t *testing.T) {
	const (
		small = 2500
		large = 10000
		limit = 8.0
	)

	run := func(count int) time.Duration {
		a := arena.New(alloc.NewBuddyAllocator(
			alloc.WithGrowthStrategy(alloc.FIXED),
			alloc.WithSize(large*MIN_BLOCK_BYTES),
		))
		defer a.Delete()

		start := time.Now()

		for range count {
			_ = arena.Alloc[int64](a)
		}

		return time.Since(start)
	}

	run(small / 2)

	quick, slow := run(small), run(large)

	ratio := float64(slow) / float64(quick)
	t.Logf("%d allocations took %v, %d took %v, ratio %.1f", small, quick, large, slow, ratio)

	if ratio > limit {
		t.Errorf("four times the work took %.1fx: cost grows with the live set",
			ratio)
	}
}

// TestSlabLargeObjectsSharePages covers one mmap per object: slabs were sized to one
// page, so every object at or above the page size cost its own mmap call.
//
// Revisions:
//   - 2026-09-21 22:10: initial creation
func TestSlabLargeObjectsSharePages(t *testing.T) {
	const (
		size  = 8192
		count = 100
	)

	slab := alloc.NewSlabAllocator()
	defer slab.Delete()

	before := slab.Stats().Mappings

	for range count {
		_ = slab.Alloc(size, 8)
	}

	mappings := slab.Stats().Mappings - before
	t.Logf("%d objects of %d bytes took %d mappings", count, size, mappings)

	if mappings >= count {
		t.Errorf("%d objects took %d mappings: every object mapped its own page",
			count, mappings)
	}
}

// TestWriterStaysArenaBacked covers a silent fall back to the heap: the writer grew with make,
// once it outgrew its first block it was on the Go heap and had quietly stopped
// being arena-backed.
//
// Revisions:
//   - 2026-09-21 22:15: initial creation
func TestWriterStaysArenaBacked(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 << 20))
	defer a.Delete()

	writer := arenaio.NewWriter(a)

	for range 8 {
		if _, err := writer.Write(make([]byte, 100)); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}

	if !arena.OwnsSlice(a, writer.Bytes()) {
		t.Error("after growing, the writer's buffer is not arena memory")
	}
}

// TestMakeStringCopiesLargeInput covers a size limit: the copy went through a fixed-size
// array type, which panicked for anything at or above one gibibyte. This uses a
// far smaller string, because the bug was the cast rather than the size.
//
// Revisions:
//   - 2026-09-21 22:19: initial creation
func TestMakeStringCopiesLargeInput(t *testing.T) {
	const size = 1 << 20

	a := arena.New(alloc.NewBumpAllocator(1 << 20))
	defer a.Delete()

	source := string(make([]byte, size))

	copied := a.MakeString(source)
	if len(copied) != size {
		t.Errorf("MakeString returned %d bytes, want %d", len(copied), size)
	}

	if !arena.OwnsString(a, copied) {
		t.Error("MakeString returned a string that is not arena memory")
	}
}

// _ChurnHeap forces several collections with allocation in between, so that
// anything unreachable is not merely collected but has its storage reused.
//
// Revisions:
//   - 2026-09-21 22:23: initial creation
func _ChurnHeap(count int) {
	const junkSize = 48

	for range 5 {
		runtime.GC()
		debug.FreeOSMemory()

		junk := make([][]byte, 0, count)

		for range count {
			block := make([]byte, junkSize)
			for i := range block {
				block[i] = 'Z'
			}

			junk = append(junk, block)
		}

		runtime.KeepAlive(junk)
	}
}
