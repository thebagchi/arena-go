package arena_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"unsafe"

	"github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
	"github.com/thebagchi/arena-go/res"
)

// TestSlab_100KInt64 covers slab 100k int 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KInt64(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*int64, count)

	for i := range count {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
	}

	seen := make(map[uintptr]bool)
	pageCounts := make(map[uintptr]int)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}

		page := addr / uintptr(res.PAGE_SIZE)
		pageCounts[page]++
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}

// TestSlab_1MInt64 covers slab 1m int 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_1MInt64(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 1_000_000

	ptrs := make([]*int64, count)

	for i := range count {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
	}

	seen := make(map[uintptr]bool)
	pageCounts := make(map[uintptr]int)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}

		page := addr / uintptr(res.PAGE_SIZE)
		pageCounts[page]++
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}

// TestSlab_100KInt32 covers slab 100k int 32.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KInt32(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*int32, count)

	for i := range count {
		ptrs[i] = arena.Alloc[int32](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int32(i)
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int32(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int32 values", count)
}

// TestSlab_100KInt16 covers slab 100k int 16.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KInt16(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*int16, count)

	for i := range count {
		ptrs[i] = arena.Alloc[int16](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int16(i)
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int16(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int16 values", count)
}

// TestSlab_100KInt8 covers slab 100k int 8.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KInt8(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*int8, count)

	for i := range count {
		ptrs[i] = arena.Alloc[int8](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int8(i)
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int8(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int8 values", count)
}

// TestSlab_100KEmpty covers slab 100k empty.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KEmpty(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type Empty struct{}

	const count = 100_000

	ptrs := make([]*Empty, count)

	for i := range count {
		ptrs[i] = arena.Alloc[Empty](a)
		if ptrs[i] == nil {
			t.Fatalf("zero-sized allocation %d failed", i)
		}
	}

	for i, ptr := range ptrs {
		if ptr == nil {
			t.Errorf("zero-sized allocation at index %d is nil", i)
		}
	}

	t.Logf("Successfully allocated and verified %d zero-sized values", count)
}

// TestSlab_100KByte100 covers slab 100k byte 100.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KByte100(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*[100]byte, count)

	for i := range count {
		ptrs[i] = arena.Alloc[[100]byte](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		copy(
			(*ptrs[i])[:],
			[]byte{byte(i % 256), byte((i / 256) % 256), byte((i / 65536) % 256)},
		)
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true
		expected := [3]byte{byte(i % 256), byte((i / 256) % 256), byte((i / 65536) % 256)}

		actual := [3]byte{(*ptrs[i])[0], (*ptrs[i])[1], (*ptrs[i])[2]}
		if actual != expected {
			t.Errorf("index %d: expected %v, got %v", i, expected, actual)
		}
	}

	t.Logf("Successfully allocated and verified %d [100]byte values", count)
}

// TestSlab_100KFloat32 covers slab 100k float 32.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KFloat32(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*float32, count)

	for i := range count {
		ptrs[i] = arena.Alloc[float32](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = float32(i) + 0.5
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		expected := float32(i) + 0.5
		if *ptrs[i] != expected {
			t.Errorf("index %d: expected %f, got %f", i, expected, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d float32 values", count)
}

// TestSlab_100KFloat64 covers slab 100k float 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KFloat64(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*float64, count)

	for i := range count {
		ptrs[i] = arena.Alloc[float64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = float64(i) + 0.5
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		expected := float64(i) + 0.5
		if *ptrs[i] != expected {
			t.Errorf("index %d: expected %f, got %f", i, expected, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d float64 values", count)
}

// TestSlab_100KBool covers slab 100k bool.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KBool(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*bool, count)

	for i := range count {
		ptrs[i] = arena.Alloc[bool](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i%2 == 0
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		expected := i%2 == 0
		if *ptrs[i] != expected {
			t.Errorf("index %d: expected %t, got %t", i, expected, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d bool values", count)
}

// TestSlab_100KTestStruct covers slab 100k test struct.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KTestStruct(t *testing.T) {
	type Struct struct {
		f1 int8
		f2 int16
		f3 int32
		f4 int64
		f5 bool
		f6 float32
		f7 float64
	}

	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*Struct, count)

	for i := range count {
		ptrs[i] = arena.Alloc[Struct](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		ptrs[i].f1 = int8(i % 128)
		ptrs[i].f2 = int16(i % 32768)
		ptrs[i].f3 = int32(i)
		ptrs[i].f4 = int64(i)
		ptrs[i].f5 = i%2 == 0
		ptrs[i].f6 = float32(i) + 0.5
		ptrs[i].f7 = float64(i) + 0.5
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		expected := Struct{
			f1: int8(i % 128),
			f2: int16(i % 32768),
			f3: int32(i),
			f4: int64(i),
			f5: i%2 == 0,
			f6: float32(i) + 0.5,
			f7: float64(i) + 0.5,
		}

		actual := *ptrs[i]
		if actual != expected {
			t.Errorf("index %d: expected %+v, got %+v", i, expected, actual)
		}
	}

	t.Logf("Successfully allocated and verified %d TestStruct values", count)
}

// TestSlab_100KStrings covers slab 100k strings.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KStrings(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	strs := make([]string, count)

	for i := range count {
		s := fmt.Sprintf("string value %d", i)

		strs[i] = a.MakeString(s)
		if strs[i] != s {
			t.Fatalf("string %d mismatch: expected %q, got %q", i, s, strs[i])
		}
	}

	pageCounts := make(map[uintptr]int)

	for i := range count {
		s := fmt.Sprintf("string value %d", i)
		if strs[i] != s {
			t.Errorf(
				"string %d verification failed: expected %q, got %q",
				i,
				s,
				strs[i],
			)
		}
		// Get address of string data
		if len(strs[i]) > 0 {
			addr := (*[2]uintptr)(unsafe.Pointer(&strs[i]))[1]
			page := addr / uintptr(res.PAGE_SIZE)
			pageCounts[page]++
		}
	}

	t.Logf("Successfully allocated and verified %d strings", count)
}

// TestSlab_100KTypesAlignment covers slab 100k types alignment.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KTypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	for i := range count {
		// int8 (1 byte alignment)
		_ = arena.Alloc[int8](a)

		// int16 (2 byte alignment)
		p16 := arena.Alloc[int16](a)
		if uintptr(unsafe.Pointer(p16))%2 != 0 {
			t.Errorf("int16 %d: addr %p not aligned to 2 bytes", i, p16)
		}

		// int32 (4 byte alignment)
		p32 := arena.Alloc[int32](a)
		if uintptr(unsafe.Pointer(p32))%4 != 0 {
			t.Errorf("int32 %d: addr %p not aligned to 4 bytes", i, p32)
		}

		// int64 (8 byte alignment)
		p64 := arena.Alloc[int64](a)
		if uintptr(unsafe.Pointer(p64))%8 != 0 {
			t.Errorf("int64 %d: addr %p not aligned to 8 bytes", i, p64)
		}
	}

	t.Logf("Successfully verified alignment for %d allocations of various types", count)
}

// TestSlab_100KArray1000TypesAlignment covers slab 100k array 1000types alignment.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_100KArray1000TypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100_000

	for i := range count {
		// [1000]int8 (1 byte alignment)
		_ = arena.Alloc[[1000]int8](a)

		// [1000]int16 (2 byte alignment)
		p16 := arena.Alloc[[1000]int16](a)
		if uintptr(unsafe.Pointer(p16))%2 != 0 {
			t.Errorf("[1000]int16 %d: addr %p not aligned to 2 bytes", i, p16)
		}

		// [1000]int32 (4 byte alignment)
		p32 := arena.Alloc[[1000]int32](a)
		if uintptr(unsafe.Pointer(p32))%4 != 0 {
			t.Errorf("[1000]int32 %d: addr %p not aligned to 4 bytes", i, p32)
		}

		// [1000]int64 (8 byte alignment)
		p64 := arena.Alloc[[1000]int64](a)
		if uintptr(unsafe.Pointer(p64))%8 != 0 {
			t.Errorf("[1000]int64 %d: addr %p not aligned to 8 bytes", i, p64)
		}
	}

	t.Logf("Successfully verified alignment for %d allocations of various array types", count)
}

// TestSlab_AppendSlice covers slab append slice.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_AppendSlice(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 0, 10)
	for i := 0; i < 100; i++ {
		slice = arena.Append(a, slice, i)
	}

	if len(slice) != 100 {
		t.Errorf("Expected length 100, got %d", len(slice))
	}

	for i, v := range slice {
		if v != i {
			t.Errorf("index %d: expected %d, got %d", i, i, v)
		}
	}

	t.Logf("Successfully verified Append for slice")
}

// TestSlab_RandomTypesLambda covers slab random types lambda.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_RandomTypesLambda(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const iterations = 100_000

	counter := 0

	allocInt8 := func() {
		ptr := arena.Alloc[int8](a)
		if ptr == nil {
			t.Fatalf("int8 allocation failed at counter %d", counter)
		}

		*ptr = int8(counter % 128)
		counter++
	}

	allocInt16 := func() {
		ptr := arena.Alloc[int16](a)
		if ptr == nil {
			t.Fatalf("int16 allocation failed at counter %d", counter)
		}

		*ptr = int16(counter % 32768)

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%2 != 0 {
			t.Errorf(
				"int16 counter %d: addr %#x not aligned to 2 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocInt32 := func() {
		ptr := arena.Alloc[int32](a)
		if ptr == nil {
			t.Fatalf("int32 allocation failed at counter %d", counter)
		}

		*ptr = int32(counter)

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%4 != 0 {
			t.Errorf(
				"int32 counter %d: addr %#x not aligned to 4 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocInt64 := func() {
		ptr := arena.Alloc[int64](a)
		if ptr == nil {
			t.Fatalf("int64 allocation failed at counter %d", counter)
		}

		*ptr = int64(counter)

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%8 != 0 {
			t.Errorf(
				"int64 counter %d: addr %#x not aligned to 8 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocFloat32 := func() {
		ptr := arena.Alloc[float32](a)
		if ptr == nil {
			t.Fatalf("float32 allocation failed at counter %d", counter)
		}

		*ptr = float32(counter) + 0.5

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%4 != 0 {
			t.Errorf(
				"float32 counter %d: addr %#x not aligned to 4 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocFloat64 := func() {
		ptr := arena.Alloc[float64](a)
		if ptr == nil {
			t.Fatalf("float64 allocation failed at counter %d", counter)
		}

		*ptr = float64(counter) + 0.5

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%8 != 0 {
			t.Errorf(
				"float64 counter %d: addr %#x not aligned to 8 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocators := []func(){
		allocInt8,
		allocInt16,
		allocInt32,
		allocInt64,
		allocFloat32,
		allocFloat64,
	}

	for range iterations {
		rand.Shuffle(len(allocators), func(i, j int) {
			allocators[i], allocators[j] = allocators[j], allocators[i]
		})

		for _, allocFunc := range allocators {
			allocFunc()
		}
	}

	t.Logf("verified alignment over %d random type allocations", iterations)
}

// TestSlab_ExpandCases covers slab expand cases.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_ExpandCases(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	t.Run("EXPANDCASE1_SINGLEELEMENTNOGROWTH", func(t *testing.T) {
		slice := arena.MakeSlice[int](a, 2, 4)
		slice[0] = 1
		slice[1] = 2

		originalCap := cap(slice)
		originalPtr := &slice[0]

		slice = arena.Append(a, slice, 3)

		if len(slice) != 3 {
			t.Errorf("Expected length 3, got %d", len(slice))
		}

		if cap(slice) != originalCap {
			t.Errorf(
				"Capacity changed unexpectedly: was %d, now %d",
				originalCap,
				cap(slice),
			)
		}

		if &slice[0] != originalPtr {
			t.Errorf("Slice backing changed unexpectedly")
		}

		if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
			t.Errorf("Slice values incorrect: %v", slice)
		}
	})

	t.Run("EXPANDCASE2_SINGLEELEMENTWITHGROWTH", func(t *testing.T) {
		slice := arena.MakeSlice[int](a, 2, 2)
		slice[0] = 1
		slice[1] = 2

		originalCap := cap(slice)
		originalPtr := &slice[0]

		slice = arena.Append(a, slice, 3)

		if len(slice) != 3 {
			t.Errorf("Expected length 3, got %d", len(slice))
		}

		if cap(slice) <= originalCap {
			t.Errorf("Capacity did not grow: was %d, now %d", originalCap, cap(slice))
		}

		if &slice[0] == originalPtr {
			t.Errorf("Slice backing should have changed")
		}

		if slice[0] != 1 || slice[1] != 2 || slice[2] != 3 {
			t.Errorf("Slice values incorrect: %v", slice)
		}
	})

	t.Run("EXPANDCASE3_MULTIELEMENTNOGROWTH", func(t *testing.T) {
		slice := arena.MakeSlice[int](a, 2, 6)
		slice[0] = 1
		slice[1] = 2

		originalCap := cap(slice)
		originalPtr := &slice[0]

		slice = arena.Append(a, slice, 3, 4)

		if len(slice) != 4 {
			t.Errorf("Expected length 4, got %d", len(slice))
		}

		if cap(slice) != originalCap {
			t.Errorf(
				"Capacity changed unexpectedly: was %d, now %d",
				originalCap,
				cap(slice),
			)
		}

		if &slice[0] != originalPtr {
			t.Errorf("Slice backing changed unexpectedly")
		}

		expected := []int{1, 2, 3, 4}
		for i, v := range expected {
			if slice[i] != v {
				t.Errorf(
					"Slice[%d] incorrect: expected %d, got %d",
					i,
					v,
					slice[i],
				)
			}
		}
	})

	t.Run("EXPANDCASE4_MULTIELEMENTWITHGROWTH", func(t *testing.T) {
		slice := arena.MakeSlice[int](a, 2, 3)
		slice[0] = 1
		slice[1] = 2

		originalCap := cap(slice)
		originalPtr := &slice[0]

		slice = arena.Append(a, slice, 3, 4, 5)

		if len(slice) != 5 {
			t.Errorf("Expected length 5, got %d", len(slice))
		}

		if cap(slice) <= originalCap {
			t.Errorf("Capacity did not grow: was %d, now %d", originalCap, cap(slice))
		}

		if &slice[0] == originalPtr {
			t.Errorf("Slice backing should have changed")
		}

		expected := []int{1, 2, 3, 4, 5}
		for i, v := range expected {
			if slice[i] != v {
				t.Errorf(
					"Slice[%d] incorrect: expected %d, got %d",
					i,
					v,
					slice[i],
				)
			}
		}
	})

	t.Run("EXPANDCASE5_APPENDTOEMPTYSLICE", func(t *testing.T) {
		slice := arena.MakeSlice[int](a, 0, 2)

		slice = arena.Append(a, slice, 10)

		if len(slice) != 1 {
			t.Errorf("Expected length 1, got %d", len(slice))
		}

		if cap(slice) < 1 {
			t.Errorf("Capacity should be at least 1, got %d", cap(slice))
		}

		if slice[0] != 10 {
			t.Errorf("Slice[0] incorrect: expected 10, got %d", slice[0])
		}
	})
}

// TestSlab_StressTest covers slab stress test.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_StressTest(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const (
		totalIterations = 1_000_000
		resetInterval   = 10_000
	)

	counter := 0

	allocInt32 := func() {
		ptr := arena.Alloc[int32](a)
		if ptr == nil {
			t.Fatalf("int32 allocation failed at counter %d", counter)
		}

		*ptr = int32(counter)
		counter++
	}

	allocInt64 := func() {
		ptr := arena.Alloc[int64](a)
		if ptr == nil {
			t.Fatalf("int64 allocation failed at counter %d", counter)
		}

		*ptr = int64(counter)
		counter++
	}

	allocFloat64 := func() {
		ptr := arena.Alloc[float64](a)
		if ptr == nil {
			t.Fatalf("float64 allocation failed at counter %d", counter)
		}

		*ptr = float64(counter) + 0.5
		counter++
	}

	allocSlice := func() {
		slice := arena.MakeSlice[int](a, 0, 10)
		for i := 0; i < 5; i++ {
			slice = arena.Append(a, slice, counter+i)
		}

		counter += 5
	}

	allocators := []func(){allocInt32, allocInt64, allocFloat64, allocSlice}

	for i := range totalIterations {
		// Randomly select an allocator
		idx := rand.Intn(len(allocators))
		allocators[idx]()

		// Reset arena periodically
		if (i+1)%resetInterval == 0 {
			a.Reset()

			counter = 0 // Reset counter to avoid large numbers
		}
	}

	t.Logf("Stress test completed: %d iterations with periodic resets", totalIterations)
}

// TestSlab_ConcurrentAllocations covers slab concurrent allocations.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_ConcurrentAllocations(t *testing.T) {
	const (
		numGoroutines           = 10
		allocationsPerGoroutine = 1000
	)

	var wg sync.WaitGroup

	errors := make(chan error, numGoroutines)

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			a := arena.New(alloc.NewSlabAllocator())
			defer a.Delete()

			// Allocate and verify various types
			for j := 0; j < allocationsPerGoroutine; j++ {
				intPtr := arena.Alloc[int](a)
				if intPtr == nil {
					errors <- fmt.Errorf(
						"goroutine %d: failed to alloc int %d",
						id,
						j,
					)

					return
				}

				*intPtr = j

				int64Ptr := arena.Alloc[int64](a)
				if int64Ptr == nil {
					errors <- fmt.Errorf(
						"goroutine %d: failed to alloc int64 %d",
						id,
						j,
					)

					return
				}

				*int64Ptr = int64(j)

				float64Ptr := arena.Alloc[float64](a)
				if float64Ptr == nil {
					errors <- fmt.Errorf(
						"goroutine %d: failed to alloc float64 %d",
						id,
						j,
					)

					return
				}

				*float64Ptr = float64(j) + 0.5

				boolPtr := arena.Alloc[bool](a)
				if boolPtr == nil {
					errors <- fmt.Errorf(
						"goroutine %d: failed to alloc bool %d",
						id,
						j,
					)

					return
				}

				*boolPtr = j%2 == 0

				str := a.MakeString(fmt.Sprintf("goroutine %d string %d", id, j))
				if str == "" {
					errors <- fmt.Errorf(
						"goroutine %d: failed to make string %d",
						id,
						j,
					)

					return
				}

				arrayPtr := arena.Alloc[[10]int](a)
				if arrayPtr == nil {
					errors <- fmt.Errorf(
						"goroutine %d: failed to alloc array %d",
						id,
						j,
					)

					return
				}

				for k := range *arrayPtr {
					(*arrayPtr)[k] = j + k
				}

				// Verify immediately
				if *intPtr != j {
					errors <- fmt.Errorf(
						"goroutine %d: int verification failed for %d",
						id,
						j,
					)

					return
				}

				if *int64Ptr != int64(j) {
					errors <- fmt.Errorf(
						"goroutine %d: int64 verification failed for %d",
						id,
						j,
					)

					return
				}

				if *float64Ptr != float64(j)+0.5 {
					errors <- fmt.Errorf(
						"goroutine %d: float64 verification failed for %d",
						id,
						j,
					)

					return
				}

				if *boolPtr != (j%2 == 0) {
					errors <- fmt.Errorf(
						"goroutine %d: bool verification failed for %d",
						id,
						j,
					)

					return
				}

				expectedStr := fmt.Sprintf("goroutine %d string %d", id, j)
				if str != expectedStr {
					errors <- fmt.Errorf(
						"goroutine %d: string %d: got %q, want %q",
						id,
						j,
						str,
						expectedStr,
					)

					return
				}

				for k := range *arrayPtr {
					if (*arrayPtr)[k] != j+k {
						errors <- fmt.Errorf(
							"goroutine %d: array %d index %d mismatch",
							id,
							j,
							k,
						)

						return
					}
				}
			}

			a.Reset()
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}

	t.Logf("%d goroutines x %d allocations completed", numGoroutines, allocationsPerGoroutine)
}

// TestSlab_InvalidInputs covers slab invalid inputs.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_InvalidInputs(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	t.Run("MAKESLICE_NEGATIVELEN", func(t *testing.T) {
		// Rejected the way the built-in make rejects it. This used to return an
		// empty slice, because the capacity was checked before the length.
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative length in MakeSlice")
			}
		}()

		arena.MakeSlice[int](a, -1, 0)
	})

	t.Run("MAKESLICE_NEGATIVECAP", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative capacity in MakeSlice")
			}
		}()

		arena.MakeSlice[int](a, 0, -1)
	})

	t.Run("MAKESLICE_LENGREATERTHANCAP", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for len > cap in MakeSlice")
			}
		}()

		arena.MakeSlice[int](a, 5, 3)
	})

	t.Run("MAKESTRING_EMPTY", func(t *testing.T) {
		s := a.MakeString("")
		if s != "" {
			t.Errorf("Expected empty string, got %q", s)
		}
	})

	t.Run("ALLOC_ZEROSIZE", func(t *testing.T) {
		ptr := arena.Alloc[struct{}](a)
		if ptr == nil {
			t.Error("Failed to allocate zero-sized struct")
		}
	})

	t.Run("RESET_AFTERDELETE", func(t *testing.T) {
		a2 := arena.New(alloc.NewSlabAllocator())
		a2.Delete()
		a2.Reset()
	})

	t.Run("ALLOC_AFTERDELETE", func(t *testing.T) {
		a3 := arena.New(alloc.NewSlabAllocator())
		a3.Delete()

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for allocation after delete")
			}
		}()

		arena.Alloc[int](a3)
	})

	t.Logf("Invalid input tests completed")
}

// TestSlab_Vec_NativeTypes covers slab vec native types.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_Vec_NativeTypes(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Test Vec with int
	vecInt := container.NewVec[int](a)
	vecInt.AppendOne(1)
	vecInt.AppendOne(2)
	vecInt.AppendOne(3)

	if vecInt.Len() != 3 {
		t.Errorf("Vec[int] length: expected 3, got %d", vecInt.Len())
	}

	if vecInt.Slice()[0] != 1 || vecInt.Slice()[1] != 2 || vecInt.Slice()[2] != 3 {
		t.Errorf("Vec[int] values incorrect: %v", vecInt.Slice())
	}

	// Test Vec with float64
	vecFloat := container.NewVec[float64](a)
	vecFloat.Append(1.1, 2.2, 3.3)

	if vecFloat.Len() != 3 {
		t.Errorf("Vec[float64] length: expected 3, got %d", vecFloat.Len())
	}

	expectedFloat := []float64{1.1, 2.2, 3.3}
	for i, v := range vecFloat.Slice() {
		if v != expectedFloat[i] {
			t.Errorf(
				"Vec[float64] index %d: expected %f, got %f",
				i,
				expectedFloat[i],
				v,
			)
		}
	}

	// Test Vec with string
	vecString := container.NewVec[string](a)
	vecString.AppendSlice([]string{"hello", "world"})

	if vecString.Len() != 2 {
		t.Errorf("Vec[string] length: expected 2, got %d", vecString.Len())
	}

	if vecString.Slice()[0] != "hello" || vecString.Slice()[1] != "world" {
		t.Errorf("Vec[string] values incorrect: %v", vecString.Slice())
	}

	// Test Vec with bool
	vecBool := container.NewVec[bool](a)
	vecBool.Push(true)
	vecBool.Push(false)

	if vecBool.Len() != 2 {
		t.Errorf("Vec[bool] length: expected 2, got %d", vecBool.Len())
	}

	if vecBool.Slice()[0] != true || vecBool.Slice()[1] != false {
		t.Errorf("Vec[bool] values incorrect: %v", vecBool.Slice())
	}

	// Test Vec with byte
	vecByte := container.NewVec[byte](a)
	vecByte.AppendOne('A')
	vecByte.AppendOne('B')

	if vecByte.Len() != 2 {
		t.Errorf("Vec[byte] length: expected 2, got %d", vecByte.Len())
	}

	if vecByte.Slice()[0] != 'A' || vecByte.Slice()[1] != 'B' {
		t.Errorf("Vec[byte] values incorrect: %v", vecByte.Slice())
	}

	// Test Vec with 100K items for each type
	const count = 100_000

	// 100K int
	vecIntLarge := container.NewVec[int](a)
	for i := range count {
		vecIntLarge.AppendOne(i)
	}

	if vecIntLarge.Len() != count {
		t.Errorf("Vec[int] 100K length: expected %d, got %d", count, vecIntLarge.Len())
	}

	if vecIntLarge.Slice()[0] != 0 || vecIntLarge.Slice()[count-1] != count-1 {
		t.Errorf(
			"Vec[int] 100K first/last incorrect: first=%d, last=%d",
			vecIntLarge.Slice()[0],
			vecIntLarge.Slice()[count-1],
		)
	}

	// 100K float64
	vecFloatLarge := container.NewVec[float64](a)
	for i := range count {
		vecFloatLarge.AppendOne(float64(i) + 0.5)
	}

	if vecFloatLarge.Len() != count {
		t.Errorf(
			"Vec[float64] 100K length: expected %d, got %d",
			count,
			vecFloatLarge.Len(),
		)
	}

	floats := vecFloatLarge.Slice()
	if floats[0] != 0.5 || floats[count-1] != float64(count-1)+0.5 {
		t.Errorf(
			"Vec[float64] 100K first/last incorrect: first=%f, last=%f",
			floats[0],
			floats[count-1],
		)
	}

	// 100K string
	vecStringLarge := container.NewVec[string](a)

	// Copied into the arena first. A Go heap string whose only reference is
	// arena memory is not kept alive by it, so storing the Sprintf result
	// directly lets the collector reuse the bytes mid-loop.
	for i := range count {
		vecStringLarge.AppendOne(a.MakeString(fmt.Sprintf("item%d", i)))
	}

	if vecStringLarge.Len() != count {
		t.Errorf(
			"Vec[string] 100K length: expected %d, got %d",
			count,
			vecStringLarge.Len(),
		)
	}

	strs := vecStringLarge.Slice()
	if strs[0] != "item0" || strs[count-1] != fmt.Sprintf("item%d", count-1) {
		t.Errorf(
			"Vec[string] 100K first/last incorrect: first=%s, last=%s",
			strs[0],
			strs[count-1],
		)
	}

	// 100K bool
	vecBoolLarge := container.NewVec[bool](a)
	for i := range count {
		vecBoolLarge.AppendOne(i%2 == 0)
	}

	if vecBoolLarge.Len() != count {
		t.Errorf("Vec[bool] 100K length: expected %d, got %d", count, vecBoolLarge.Len())
	}

	if vecBoolLarge.Slice()[0] != true || vecBoolLarge.Slice()[count-1] != ((count-1)%2 == 0) {
		t.Errorf(
			"Vec[bool] 100K first/last incorrect: first=%t, last=%t",
			vecBoolLarge.Slice()[0],
			vecBoolLarge.Slice()[count-1],
		)
	}

	// 100K byte
	vecByteLarge := container.NewVec[byte](a)
	for i := range count {
		vecByteLarge.AppendOne(byte((i % 256)))
	}

	if vecByteLarge.Len() != count {
		t.Errorf("Vec[byte] 100K length: expected %d, got %d", count, vecByteLarge.Len())
	}

	if vecByteLarge.Slice()[0] != 0 || vecByteLarge.Slice()[count-1] != byte((count-1)%256) {
		t.Errorf(
			"Vec[byte] 100K first/last incorrect: first=%d, last=%d",
			vecByteLarge.Slice()[0],
			vecByteLarge.Slice()[count-1],
		)
	}

	t.Logf("Vec native types test completed successfully")
}

// TestSlab_AllocationsAfterReset covers slab allocations after reset.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_AllocationsAfterReset(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 10_000

	// First allocation batch
	ptrs1 := make([]*int64, count)
	for i := range count {
		ptrs1[i] = arena.Alloc[int64](a)
		if ptrs1[i] == nil {
			t.Fatalf("first allocation batch %d failed", i)
		}

		*ptrs1[i] = int64(i)
	}

	// Verify first batch
	for i := range count {
		if *ptrs1[i] != int64(i) {
			t.Errorf("first batch index %d: expected %d, got %d", i, i, *ptrs1[i])
		}
	}

	// Reset the arena
	a.Reset()

	// Second allocation batch after reset
	ptrs2 := make([]*int64, count)
	for i := range count {
		ptrs2[i] = arena.Alloc[int64](a)
		if ptrs2[i] == nil {
			t.Fatalf("second allocation batch %d failed after reset", i)
		}

		*ptrs2[i] = int64(i * 2)
	}

	// Verify second batch
	for i := range count {
		if *ptrs2[i] != int64(i*2) {
			t.Errorf("second batch index %d: expected %d, got %d", i, i*2, *ptrs2[i])
		}
	}

	// Test with different types after reset
	a.Reset()

	const count2 = 5000

	// Float64 allocations
	floatPtrs := make([]*float64, count2)
	for i := range count2 {
		floatPtrs[i] = arena.Alloc[float64](a)
		if floatPtrs[i] == nil {
			t.Fatalf("float64 allocation %d failed after reset", i)
		}

		*floatPtrs[i] = float64(i) + 0.5
	}

	// String allocations
	strs := make([]string, count2)
	for i := range count2 {
		strs[i] = a.MakeString(fmt.Sprintf("reset_test_%d", i))
	}

	// Verify float64 and string allocations
	for i := range count2 {
		if *floatPtrs[i] != float64(i)+0.5 {
			t.Errorf(
				"float64 index %d: expected %f, got %f",
				i,
				float64(i)+0.5,
				*floatPtrs[i],
			)
		}

		expected := fmt.Sprintf("reset_test_%d", i)
		if strs[i] != expected {
			t.Errorf("string index %d: expected %q, got %q", i, expected, strs[i])
		}
	}

	t.Logf("Successfully verified allocations after Reset")
}

// TestSlab_BasicAllocation covers slab basic allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_BasicAllocation(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate small objects
	p1 := arena.Alloc[int](a)
	if p1 == nil {
		t.Fatal("allocation failed")
	}

	*p1 = 42

	p2 := arena.Alloc[int64](a)
	if p2 == nil {
		t.Fatal("allocation failed")
	}

	*p2 = 99

	if *p1 != 42 {
		t.Errorf("expected 42, got %d", *p1)
	}

	if *p2 != 99 {
		t.Errorf("expected 99, got %d", *p2)
	}
}

// TestSlab_SizeClasses covers slab size classes.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_SizeClasses(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Test allocations of various sizes using MakeSlice
	type TestCase struct {
		name string
		size int
	}

	cases := []TestCase{
		{"tiny", 1},
		{"small", 16},
		{"medium", 256},
		{"large", 4096},
		{"xlarge", 262144},
	}

	for _, tc := range cases {
		p := arena.MakeSlice[byte](a, tc.size, tc.size)
		if len(p) == 0 {
			t.Fatalf("%s: allocation failed", tc.name)
		}
		// Write to allocated memory to verify it's accessible
		p[0] = 42
	}
}

// TestSlab_MultipleObjects covers slab multiple objects.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_MultipleObjects(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	const count = 100

	ptrs := make([]*int, count)

	// Allocate many small objects
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i * 10
	}

	// Verify all values
	for i := 0; i < count; i++ {
		if *ptrs[i] != i*10 {
			t.Errorf("index %d: expected %d, got %d", i, i*10, *ptrs[i])
		}
	}
}

// TestSlab_RemoveAndReuse covers slab remove and reuse.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_RemoveAndReuse(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate and free
	p1 := arena.Alloc[int](a)
	*p1 = 100
	addr1 := unsafe.Pointer(p1)

	arena.DeleteObject(a, p1)

	// Allocate again - should potentially reuse the freed block
	p2 := arena.Alloc[int](a)
	addr2 := unsafe.Pointer(p2)

	if p2 == nil {
		t.Fatal("reallocation failed")
	}

	*p2 = 200

	if *p2 != 200 {
		t.Errorf("expected 200, got %d", *p2)
	}

	t.Logf("addr1=%p addr2=%p (may or may not reuse)", addr1, addr2)
}

// TestSlab_InterleavedAllocFree covers slab interleaved alloc free.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_InterleavedAllocFree(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate, free, allocate pattern
	ptrs := make([]*int, 100)

	for i := 0; i < 100; i++ {
		ptrs[i] = arena.Alloc[int](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i
	}

	// Free every other one
	for i := 0; i < 100; i += 2 {
		arena.DeleteObject(a, ptrs[i])
		ptrs[i] = nil
	}

	// Allocate more - should fill freed slots
	for i := 0; i < 50; i++ {
		p := arena.Alloc[int](a)
		if p == nil {
			t.Fatalf("reallocation %d failed", i)
		}

		*p = i + 1000
	}

	// Verify remaining originals
	for i := 1; i < 100; i += 2 {
		if *ptrs[i] != i {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}
}

// TestSlab_DifferentTypes covers slab different types.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_DifferentTypes(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate different types
	pi := arena.Alloc[int](a)
	*pi = 42

	ps := arena.Alloc[string](a)
	*ps = "hello"

	pf := arena.Alloc[float64](a)
	*pf = 3.14

	pb := arena.Alloc[bool](a)
	*pb = true

	// Verify
	if *pi != 42 {
		t.Errorf("int: expected 42, got %d", *pi)
	}

	if *ps != "hello" {
		t.Errorf("string: expected 'hello', got %q", *ps)
	}

	if *pf != 3.14 {
		t.Errorf("float64: expected 3.14, got %f", *pf)
	}

	if *pb != true {
		t.Errorf("bool: expected true, got %v", *pb)
	}
}

// TestSlab_StructAllocation covers slab struct allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_StructAllocation(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	type TestStruct struct {
		A int
		B string
		C float64
	}

	p := arena.Alloc[TestStruct](a)
	if p == nil {
		t.Fatal("struct allocation failed")
	}

	p.A = 42
	p.B = "test"
	p.C = 3.14

	if p.A != 42 || p.B != "test" || p.C != 3.14 {
		t.Error("struct field values incorrect")
	}
}

// TestSlab_ArrayAllocation covers slab array allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_ArrayAllocation(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate array
	arr := arena.MakeSlice[int](a, 100, 100)
	if len(arr) == 0 {
		t.Fatal("array allocation failed")
	}

	// Fill with values
	for i := 0; i < 100; i++ {
		arr[i] = i * 2
	}

	// Verify
	for i := 0; i < 100; i++ {
		if arr[i] != i*2 {
			t.Errorf("index %d: expected %d, got %d", i, i*2, arr[i])
		}
	}
}

// TestSlab_Owns covers slab owns.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_Owns(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	p := arena.Alloc[int](a)
	if !a.Owns(unsafe.Pointer(p)) {
		t.Error("allocator should own allocated pointer")
	}

	// Test with unowned pointer
	var x int
	if a.Owns(unsafe.Pointer(&x)) {
		t.Error("allocator should not own stack pointer")
	}
}

// TestSlab_Reset covers slab reset.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_Reset(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate before reset
	p1 := arena.Alloc[int](a)
	*p1 = 42

	a.Reset()

	// Allocate after reset
	p2 := arena.Alloc[int](a)
	if p2 == nil {
		t.Fatal("allocation after reset failed")
	}

	*p2 = 99

	if *p2 != 99 {
		t.Errorf("expected 99, got %d", *p2)
	}
}

// TestSlab_ZeroSizedAllocation covers slab zero sized allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_ZeroSizedAllocation(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate zero-length array
	arr := arena.MakeSlice[int](a, 0, 0)
	if len(arr) != 0 {
		t.Error("zero-length allocation should be empty")
	}
}

// TestSlab_LargeAllocations covers slab large allocations.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_LargeAllocations(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate increasingly large arrays
	sizes := []int{1024, 4096, 16384, 65536, 262144}

	for _, size := range sizes {
		arr := arena.MakeSlice[byte](a, size, size)
		if len(arr) == 0 {
			t.Fatalf("allocation of size %d failed", size)
		}

		// Write and verify at boundaries
		arr[0] = 42
		arr[size-1] = 99

		if arr[0] != 42 {
			t.Errorf("size %d: first byte incorrect", size)
		}

		if arr[size-1] != 99 {
			t.Errorf("size %d: last byte incorrect", size)
		}
	}
}

// TestSlab_RapidAllocations covers slab rapid allocations.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_RapidAllocations(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Test rapid sequential allocations
	const count = 100

	ptrs := make([]*int, count)

	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i
	}

	// Verify
	for i := 0; i < count; i++ {
		if *ptrs[i] != i {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}
}

// TestSlab_MixedSizes covers slab mixed sizes.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestSlab_MixedSizes(t *testing.T) {
	a := arena.New(alloc.NewSlabAllocator())
	defer a.Delete()

	// Allocate mix of small and large objects
	for i := 0; i < 20; i++ {
		// Small object
		ps := arena.Alloc[int](a)
		if ps == nil {
			t.Fatalf("small allocation %d failed", i)
		}

		*ps = i

		// Medium object
		pm := arena.MakeSlice[byte](a, 256, 256)
		if len(pm) == 0 {
			t.Fatalf("medium allocation %d failed", i)
		}

		pm[0] = byte(i % 256)

		// Large object every 5 iterations
		if i%5 == 0 {
			pl := arena.MakeSlice[int](a, 100, 100)
			if len(pl) == 0 {
				t.Fatalf("large allocation %d failed", i)
			}

			pl[0] = i * 10
		}
	}
}
