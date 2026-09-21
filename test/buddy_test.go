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

// TestBuddy_ReallocOrder covers buddy realloc order.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_ReallocOrder(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 128
	// First allocation
	ptrs := make([]*int64, count)
	for i := range count {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("first allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
	}

	// Delete all pointers
	for _, ptr := range ptrs {
		arena.DeleteObject(a, ptr)
	}

	// Second allocation and verify in one loop
	for i := range count {
		ptr := arena.Alloc[int64](a)
		if ptr == nil {
			t.Fatalf("second allocation %d failed", i)
		}

		*ptr = int64(i + 1000) // Different value to distinguish

		// Verify address matches first allocation
		addr := uintptr(unsafe.Pointer(ptr))
		if addr != uintptr(unsafe.Pointer(ptrs[i])) {
			temp := uintptr(unsafe.Pointer(ptrs[i]))
			t.Errorf("order mismatch at %d: want %#x, got %#x", i, temp, addr)
		}

		if *ptr != int64(i+1000) {
			t.Errorf(
				"value mismatch at index %d: expected %d, got %d",
				i,
				i+1000,
				*ptr,
			)
		}
	}

	t.Logf("Successfully verified reallocation order for %d int64 values", count)
}

// TestBuddy_FragmentationRecovery covers buddy fragmentation recovery.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_FragmentationRecovery(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 256
	// First allocation of 256 pointers
	ptrs := make([]*int64, count)
	for i := range count {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
	}

	// Free even pointers
	for i := range count {
		if i%2 == 0 {
			arena.DeleteObject(a, ptrs[i])
		}
	}

	// Allocate 128 new pointers
	for i := range 128 {
		ptr := arena.Alloc[int64](a)
		if ptr == nil {
			t.Fatalf("second allocation %d failed", i)
		}

		*ptr = int64(i + 2000)

		// Verify address matches the freed even pointer
		addr := uintptr(unsafe.Pointer(ptr))
		if addr != uintptr(unsafe.Pointer(ptrs[i*2])) {
			temp := uintptr(unsafe.Pointer(ptrs[i*2]))
			t.Errorf("order mismatch at %d: want %#x, got %#x", i, temp, addr)
		}

		if *ptr != int64(i+2000) {
			t.Errorf(
				"value mismatch at index %d: expected %d, got %d",
				i,
				i+2000,
				*ptr,
			)
		}
	}

	t.Logf("Successfully verified allocation order after freeing even pointers")

	// Free odd pointers
	for i := range count {
		if i%2 == 1 {
			arena.DeleteObject(a, ptrs[i])
		}
	}

	// Allocate 128 more pointers
	for i := range 128 {
		ptr := arena.Alloc[int64](a)
		if ptr == nil {
			t.Fatalf("third allocation %d failed", i)
		}

		*ptr = int64(i + 3000)

		// Verify address matches the freed odd pointer
		addr := uintptr(unsafe.Pointer(ptr))
		if addr != uintptr(unsafe.Pointer(ptrs[i*2+1])) {
			temp := uintptr(unsafe.Pointer(ptrs[i*2+1]))
			t.Errorf("order mismatch at %d: want %#x, got %#x", i, temp, addr)
		}

		if *ptr != int64(i+3000) {
			t.Errorf(
				"value mismatch at index %d: expected %d, got %d",
				i,
				i+3000,
				*ptr,
			)
		}
	}

	t.Logf("Successfully verified allocation order after freeing odd pointers")
}

// TestBuddy_100KInt64 covers buddy 100k int 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KInt64(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_1MInt64 covers buddy 1m int 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_1MInt64(t *testing.T) {
	//t.Skip("Skipping 1M test in buddy")
	a := arena.New(alloc.NewBuddyAllocator(
		alloc.WithGrowthStrategy(alloc.FIXED),
		alloc.WithSize(1_000_000*16),
	))
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

// TestBuddy_100KStrings covers buddy 100k strings.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KStrings(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KInt32 covers buddy 100k int 32.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KInt32(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KInt16 covers buddy 100k int 16.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KInt16(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KInt8 covers buddy 100k int 8.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KInt8(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KEmpty covers buddy 100k empty.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KEmpty(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KByte100 covers buddy 100k byte 100.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KByte100(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KFloat32 covers buddy 100k float 32.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KFloat32(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KFloat64 covers buddy 100k float 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KFloat64(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KBool covers buddy 100k bool.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KBool(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KTestStruct covers buddy 100k test struct.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KTestStruct(t *testing.T) {
	type Struct struct {
		f1 int8
		f2 int16
		f3 int32
		f4 int64
		f5 bool
		f6 float32
		f7 float64
	}

	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KTypesAlignment covers buddy 100k types alignment.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KTypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_100KArray1000TypesAlignment covers buddy 100k array 1000types alignment.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KArray1000TypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_AppendSlice covers buddy append slice.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_AppendSlice(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_RandomTypesLambda covers buddy random types lambda.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_RandomTypesLambda(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_ExpandCases covers buddy expand cases.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_ExpandCases(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_StressTest covers buddy stress test.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_StressTest(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_ConcurrentAllocations covers buddy concurrent allocations.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_ConcurrentAllocations(t *testing.T) {
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

			a := arena.New(alloc.NewBuddyAllocator())
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

// TestBuddy_InvalidInputs covers buddy invalid inputs.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_InvalidInputs(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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
		a2 := arena.New(alloc.NewBuddyAllocator())
		a2.Delete()
		a2.Reset()
	})

	t.Run("ALLOC_AFTERDELETE", func(t *testing.T) {
		a3 := arena.New(alloc.NewBuddyAllocator())
		a3.Delete()

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when allocating after delete")
			}
		}()

		arena.Alloc[int](a3)
	})

	t.Logf("Invalid input tests completed")
}

// TestBuddy_Vec_NativeTypes covers buddy vec native types.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_Vec_NativeTypes(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
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
		vecByteLarge.AppendOne(byte(i % 256))
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

	t.Logf("Vec tests with native types completed successfully")
}

// TestBuddy_BasicAllocation covers buddy basic allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_BasicAllocation(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate various sizes
	p1 := arena.Alloc[int](a)
	if p1 == nil {
		t.Fatal("allocation failed")
	}

	*p1 = 42

	p2 := arena.Alloc[uint64](a)
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

// TestBuddy_RemoveAndReuse covers buddy remove and reuse.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_RemoveAndReuse(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate and free
	p1 := arena.Alloc[int](a)
	*p1 = 100
	addr1 := unsafe.Pointer(p1)

	arena.DeleteObject(a, p1)

	// Allocate again - should reuse the freed block
	p2 := arena.Alloc[int](a)
	addr2 := unsafe.Pointer(p2)

	// Due to buddy coalescing, addresses might be the same or different
	// but allocation should succeed
	if p2 == nil {
		t.Fatal("reallocation failed")
	}

	*p2 = 200

	if *p2 != 200 {
		t.Errorf("expected 200, got %d", *p2)
	}

	t.Logf("addr1=%p addr2=%p", addr1, addr2)
}

// TestBuddy_MultipleAllocations covers buddy multiple allocations.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_MultipleAllocations(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 100

	ptrs := make([]*int, count)

	// Allocate many objects
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i * 10
	}

	// Verify all pointers are unique
	seen := make(map[uintptr]bool)

	for i := 0; i < count; i++ {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true
		// Verify value is still what we set
		if *ptrs[i] != i*10 {
			t.Errorf("index %d: expected %d, got %d", i, i*10, *ptrs[i])
		}
	}
}

// TestBuddy_Coalescing covers buddy coalescing.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_Coalescing(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate multiple small blocks
	p1 := arena.Alloc[int](a)
	p2 := arena.Alloc[int](a)
	p3 := arena.Alloc[int](a)
	p4 := arena.Alloc[int](a)

	*p1 = 1
	*p2 = 2
	*p3 = 3
	*p4 = 4

	// Free some blocks - should coalesce
	arena.DeleteObject(a, p1)
	arena.DeleteObject(a, p2)

	// Allocate larger block - might reuse coalesced space
	p5 := arena.Alloc[uint64](a)
	if p5 == nil {
		t.Fatal("allocation after coalescing failed")
	}

	*p5 = 123456789

	// Remaining allocations should still be valid
	if *p3 != 3 {
		t.Errorf("p3: expected 3, got %d", *p3)
	}

	if *p4 != 4 {
		t.Errorf("p4: expected 4, got %d", *p4)
	}
}

// TestBuddy_Reset covers buddy reset.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_Reset(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate some memory
	for i := 0; i < 10; i++ {
		p := arena.Alloc[int](a)
		*p = i
	}

	// Reset should clear allocations
	a.Reset()

	// Should be able to allocate again
	p := arena.Alloc[int](a)
	if p == nil {
		t.Fatal("allocation after reset failed")
	}

	*p = 999
	if *p != 999 {
		t.Errorf("expected 999, got %d", *p)
	}
}

// TestBuddy_LargeAllocation covers buddy large allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_LargeAllocation(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate large structure
	type Large struct {
		data [1024]byte
	}

	p := arena.Alloc[Large](a)
	if p == nil {
		t.Fatal("large allocation failed")
	}

	p.data[0] = 0xFF
	p.data[1023] = 0xAA

	if p.data[0] != 0xFF || p.data[1023] != 0xAA {
		t.Error("large allocation data corruption")
	}
}

// TestBuddy_Owns covers buddy owns.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_Owns(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	p := arena.Alloc[int](a)
	*p = 42

	if !a.Owns(unsafe.Pointer(p)) {
		t.Error("Owns should return true for allocated pointer")
	}

	external := new(int)
	if a.Owns(unsafe.Pointer(external)) {
		t.Error("Owns should return false for external pointer")
	}

	if a.Owns(nil) {
		t.Error("Owns should return false for nil pointer")
	}
}

// TestBuddy_ZeroSizedType covers buddy zero sized type.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_ZeroSizedType(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	type Empty struct{}

	p := arena.Alloc[Empty](a)
	if p == nil {
		t.Fatal("zero-sized allocation failed")
	}

	// Should not crash
	_ = *p
}

// TestBuddy_AlignedAllocation covers buddy aligned allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_AlignedAllocation(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	p := arena.Alloc[uint64](a)
	if p == nil {
		t.Fatal("allocation failed")
	}

	// Check alignment
	addr := uintptr(unsafe.Pointer(p))
	if addr%8 != 0 {
		t.Errorf("uint64 not aligned: address=%#x", addr)
	}
}

// TestBuddy_InterleavedAllocFree covers buddy interleaved alloc free.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_InterleavedAllocFree(t *testing.T) {
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Test alloc/free/alloc cycle with Reset instead of Remove
	// since buddy doesn't track allocation sizes for individual frees
	for cycle := 0; cycle < 3; cycle++ {
		ptrs := make([]*int, 10)

		// Allocate
		for i := 0; i < 10; i++ {
			ptrs[i] = arena.Alloc[int](a)
			if ptrs[i] == nil {
				t.Fatalf("cycle %d: allocation %d failed", cycle, i)
			}

			*ptrs[i] = cycle*100 + i
		}

		// Verify
		for i := 0; i < 10; i++ {
			if *ptrs[i] != cycle*100+i {
				t.Errorf("cycle %d, index %d: expected %d, got %d",
					cycle, i, cycle*100+i, *ptrs[i])
			}
		}

		// Reset for next cycle
		a.Reset()
	}
}

// TestBuddy_MultipleChunks covers buddy multiple chunks.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_MultipleChunks(t *testing.T) {
	// Start with multiple chunks
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 200

	ptrs := make([]*int, count)

	// Allocate across multiple chunks
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = i * 3
	}

	// Verify all pointers are unique
	seen := make(map[uintptr]bool)

	for i := 0; i < count; i++ {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true
		// Verify value
		if *ptrs[i] != i*3 {
			t.Errorf("index %d: expected %d, got %d", i, i*3, *ptrs[i])
		}
	}
}

// TestBuddy_OversizedAllocation covers buddy oversized allocation.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_OversizedAllocation(t *testing.T) {
	// Create buddy allocator with small chunk size
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	// Allocate something larger than chunk size (4KB default)
	type Huge struct {
		data [8192]byte // 8KB, larger than default 4KB page
	}

	p := arena.Alloc[Huge](a)
	if p == nil {
		t.Fatal("oversized allocation failed")
	}

	// Write to verify it's accessible
	p.data[0] = 0xAA
	p.data[8191] = 0xBB

	if p.data[0] != 0xAA || p.data[8191] != 0xBB {
		t.Error("oversized allocation data corruption")
	}

	// Allocate another oversized block
	p2 := arena.Alloc[Huge](a)
	if p2 == nil {
		t.Fatal("second oversized allocation failed")
	}

	// Ensure they're different
	if unsafe.Pointer(p) == unsafe.Pointer(p2) {
		t.Error("oversized allocations returned same address")
	}

	p2.data[0] = 0xCC
	p2.data[8191] = 0xDD

	// Verify both still work
	if p.data[0] != 0xAA || p.data[8191] != 0xBB {
		t.Error("first oversized allocation corrupted")
	}

	if p2.data[0] != 0xCC || p2.data[8191] != 0xDD {
		t.Error("second oversized allocation corrupted")
	}
}

// TestBuddy_DetailedDebug covers buddy detailed debug.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_DetailedDebug(t *testing.T) {
	// Allocate just enough to see the 32766-32768 issue in detail
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 32800

	ptrs := make([]*uint64, count)
	seen := make(map[uintptr]bool)

	// Allocate up to the problem area
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[uint64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = uint64(i)

		// Check for duplicates using map instead of nested loop
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("DUPLICATE DETECTED: alloc %d ptr=%p", i, ptrs[i])
		}

		seen[addr] = true
	}

	t.Logf("Completed %d allocations", count)
}

// TestBuddy_100KUint64 covers buddy 100k uint 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_100KUint64(t *testing.T) {
	// Test with 100K to isolate the issue
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*uint64, count)

	// Allocate 100K uint64
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[uint64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = uint64(i)
	}

	// Verify all pointers are unique
	seen := make(map[uintptr]bool)

	for i := 0; i < count; i++ {
		addr := uintptr(unsafe.Pointer(ptrs[i]))

		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true
		// Verify value
		if *ptrs[i] != uint64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf(
		"Successfully allocated and verified %d uint64 values across multiple chunks",
		count,
	)
}

// TestBuddy_10KUint64 covers buddy 10k uint 64.
//
// Revisions:
//   - 2025-12-23 19:20: initial creation
func TestBuddy_10KUint64(t *testing.T) {
	// Test with 10K to isolate the issue
	a := arena.New(alloc.NewBuddyAllocator())
	defer a.Delete()

	const count = 10_000

	ptrs := make([]*uint64, count)

	// Allocate 10K uint64
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[uint64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = uint64(i)
	}

	// Verify all pointers are unique
	seen := make(map[uintptr]bool)

	for i := 0; i < count; i++ {
		addr := uintptr(unsafe.Pointer(ptrs[i]))

		// Figure out which chunk this belongs to (simplified without Res)
		chunkIdx := -1

		if seen[addr] {
			t.Errorf(
				"duplicate address at index %d (chunk %d): %#x",
				i,
				chunkIdx,
				addr,
			)
		}

		seen[addr] = true
		// Verify value
		if *ptrs[i] != uint64(i) {
			t.Errorf(
				"index %d (chunk %d): expected %d, got %d",
				i,
				chunkIdx,
				i,
				*ptrs[i],
			)

			if i > 32760 && i < 32770 {
				t.Logf(
					"near error: index=%d ptr=%p chunk=%d value=%d",
					i,
					ptrs[i],
					chunkIdx,
					*ptrs[i],
				)
			}
		}
	}

	t.Logf(
		"Successfully allocated and verified %d uint64 values across multiple chunks",
		count,
	)
}
