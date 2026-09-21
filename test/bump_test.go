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

// TestBump_100KInt64 covers bump 100k int 64.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KInt64(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}

// TestBump_1MInt64 covers bump 1m int 64.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_1MInt64(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(2000 * res.PAGE_SIZE))
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

	for i := range count {
		addr := uintptr(unsafe.Pointer(ptrs[i]))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}

// TestBump_100KInt32 covers bump 100k int 32.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KInt32(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KInt16 covers bump 100k int 16.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KInt16(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KInt8 covers bump 100k int 8.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KInt8(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KEmpty covers bump 100k empty.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KEmpty(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	type Empty struct{}

	t.Logf(
		"Empty struct: Sizeof=%d, Alignof=%d",
		unsafe.Sizeof(Empty{}),
		unsafe.Alignof(Empty{}),
	)

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

// TestBump_100KByte100 covers bump 100k byte 100.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KByte100(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KFloat32 covers bump 100k float 32.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KFloat32(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KFloat64 covers bump 100k float 64.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KFloat64(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KBool covers bump 100k bool.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KBool(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_100KMixedStruct covers bump 100k mixed struct.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KMixedStruct(t *testing.T) {
	type Struct struct {
		f1 int8
		f2 int16
		f3 int32
		f4 int64
		f5 bool
		f6 float32
		f7 float64
	}

	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	t.Logf(
		"MixedStruct: Sizeof=%d, Alignof=%d",
		unsafe.Sizeof(Struct{}),
		unsafe.Alignof(Struct{}),
	)

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

	t.Logf("Successfully allocated and verified %d MixedStruct values", count)
}

// TestBump_100KStrings covers bump 100k strings.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KStrings(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	const count = 100_000

	strings := make([]string, count)

	for i := range count {
		strings[i] = a.MakeString(fmt.Sprintf("string %d", i))
		if strings[i] == "" {
			t.Fatalf("allocation %d failed", i)
		}
	}

	for i := range count {
		addr := uintptr(unsafe.Pointer(unsafe.StringData(strings[i])))
		t.Logf("index %d: str=%q, addr=%#x", i, strings[i], addr)
	}

	seen := make(map[uintptr]bool)

	for i := range count {
		addr := uintptr(unsafe.Pointer(unsafe.StringData(strings[i])))
		if seen[addr] {
			t.Errorf("duplicate address at index %d: %#x", i, addr)
		}

		seen[addr] = true

		expected := fmt.Sprintf("string %d", i)
		if strings[i] != expected {
			t.Errorf("index %d: expected %s, got %s", i, expected, strings[i])
		}
	}

	t.Logf("Successfully allocated and verified %d strings", count)
}

// TestBump_100KTypesAlignment covers bump 100k types alignment.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KTypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

	allocators := []func(){allocInt8, allocInt16, allocInt32, allocInt64}

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

// TestBump_100KArray1000TypesAlignment covers bump 100k array 1000types alignment.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KArray1000TypesAlignment(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	const iterations = 10_000

	counter := 0

	allocInt8 := func() {
		ptr := arena.Alloc[[1000]int8](a)
		if ptr == nil {
			t.Fatalf("int8 array allocation failed at counter %d", counter)
		}

		for i := range *ptr {
			(*ptr)[i] = int8((counter + i) % 128)
		}

		counter++
	}

	allocInt16 := func() {
		ptr := arena.Alloc[[1000]int16](a)
		if ptr == nil {
			t.Fatalf("int16 array allocation failed at counter %d", counter)
		}

		for i := range *ptr {
			(*ptr)[i] = int16((counter + i) % 32768)
		}

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%2 != 0 {
			t.Errorf(
				"int16 array counter %d: addr %#x not aligned to 2 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocInt32 := func() {
		ptr := arena.Alloc[[1000]int32](a)
		if ptr == nil {
			t.Fatalf("int32 array allocation failed at counter %d", counter)
		}

		for i := range *ptr {
			(*ptr)[i] = int32(counter + i)
		}

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%4 != 0 {
			t.Errorf(
				"int32 array counter %d: addr %#x not aligned to 4 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocInt64 := func() {
		ptr := arena.Alloc[[1000]int64](a)
		if ptr == nil {
			t.Fatalf("int64 array allocation failed at counter %d", counter)
		}

		for i := range *ptr {
			(*ptr)[i] = int64(counter + i)
		}

		addr := uintptr(unsafe.Pointer(ptr))
		if addr%8 != 0 {
			t.Errorf(
				"int64 array counter %d: addr %#x not aligned to 8 bytes",
				counter,
				addr,
			)
		}

		counter++
	}

	allocators := []func(){allocInt8, allocInt16, allocInt32, allocInt64}

	for range iterations {
		rand.Shuffle(len(allocators), func(i, j int) {
			allocators[i], allocators[j] = allocators[j], allocators[i]
		})

		for _, allocFunc := range allocators {
			allocFunc()
		}
	}

	t.Logf("verified alignment over %d random array allocations", iterations)
}

// TestBump_AppendSlice covers bump append slice.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_AppendSlice(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	slice := arena.MakeSlice[int](a, 0, 4)

	const count = 100_000
	for i := range count {
		slice = arena.Append(a, slice, i)
	}

	if len(slice) != count {
		t.Errorf("Expected length %d, got %d", count, len(slice))
	}

	for i := range count {
		if slice[i] != i {
			t.Errorf("Expected slice[%d] = %d, got %d", i, i, slice[i])
		}
	}

	t.Logf(
		"Successfully appended %d elements to slice: final length %d, capacity %d",
		count,
		len(slice),
		cap(slice),
	)
}

// TestBump_RandomTypesLambda covers bump random types lambda.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_RandomTypesLambda(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	const iterations = 100_000

	counter := 0

	allocInt8 := func() {
		ptr := arena.Alloc[int8](a)
		if ptr == nil {
			t.Fatalf("int8 allocation failed at counter %d", counter)
		}

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("int8 allocation not owned by arena at counter %d", counter)
		}

		*ptr = int8(counter % 128)
		counter++
	}

	allocInt16 := func() {
		ptr := arena.Alloc[int16](a)
		if ptr == nil {
			t.Fatalf("int16 allocation failed at counter %d", counter)
		}

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("int16 allocation not owned by arena at counter %d", counter)
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

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("int32 allocation not owned by arena at counter %d", counter)
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

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("int64 allocation not owned by arena at counter %d", counter)
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

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("float32 allocation not owned by arena at counter %d", counter)
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

		if !a.Owns(unsafe.Pointer(ptr)) {
			t.Errorf("float64 allocation not owned by arena at counter %d", counter)
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

// TestBump_ExpandCases covers bump expand cases.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_ExpandCases(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_StressTest covers bump stress test.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_StressTest(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	const (
		totalIterations = 10_000_000
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

// TestBump_ConcurrentAllocations covers bump concurrent allocations.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_ConcurrentAllocations(t *testing.T) {
	const (
		numGoroutines           = 10
		allocationsPerGoroutine = 10000
	)

	var wg sync.WaitGroup

	errors := make(chan error, numGoroutines)

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// TestBump_InvalidInputs covers bump invalid inputs.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_InvalidInputs(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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
		a2 := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
		a2.Delete()
		// Reset after delete may be safe, check if it doesn't crash
		a2.Reset() // If it panics, the test will fail
	})

	t.Run("ALLOC_AFTERDELETE", func(t *testing.T) {
		a3 := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// FuzzBumpAlloc covers bump alloc.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func FuzzBumpAlloc(f *testing.F) {
	// Seed with some initial values
	f.Add(10, 20)   // len, cap for slice
	f.Add(0, 0)     // empty slice
	f.Add(100, 200) // larger slice
	f.Add(1, 1)     // minimal

	f.Fuzz(func(t *testing.T, sliceLen, sliceCap int) {
		// Skip invalid inputs to avoid expected panics
		if sliceLen < 0 || sliceCap < 0 || sliceLen > sliceCap {
			return
		}

		a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
		defer a.Delete()

		// Allocate a slice with fuzzed len/cap
		slice := arena.MakeSlice[int](a, sliceLen, sliceCap)
		if len(slice) != sliceLen {
			t.Errorf("Slice len mismatch: got %d, expected %d", len(slice), sliceLen)
		}

		if cap(slice) < sliceCap {
			t.Errorf(
				"Slice cap too small: got %d, expected >= %d",
				cap(slice),
				sliceCap,
			)
		}

		// Fill the slice to test allocation
		for i := range slice {
			slice[i] = i
		}

		// Allocate some primitives
		intPtr := arena.Alloc[int](a)
		if intPtr == nil {
			t.Error("Failed to allocate int")
		} else {
			*intPtr = 42
		}

		floatPtr := arena.Alloc[float64](a)
		if floatPtr == nil {
			t.Error("Failed to allocate float64")
		} else {
			*floatPtr = 3.14
		}

		boolPtr := arena.Alloc[bool](a)
		if boolPtr == nil {
			t.Error("Failed to allocate bool")
		} else {
			*boolPtr = true
		}

		// Allocate a string
		str := a.MakeString("fuzz test string")
		if str != "fuzz test string" {
			t.Errorf("String allocation failed: got %q", str)
		}

		// Allocate an array
		arrayPtr := arena.Alloc[[5]int](a)
		if arrayPtr == nil {
			t.Error("Failed to allocate array")
		} else {
			for i := range *arrayPtr {
				(*arrayPtr)[i] = i * 2
			}
		}

		// A struct of every basic kind, written through arena memory and read
		// back whole. Every field is part of the shape being checked, so every
		// field is set and compared rather than left as padding.
		type RoundTripStruct struct {
			f01 [5]int
			f02 [5]int8
			f03 [5]int16
			f04 [5]int32
			f05 [5]int64
			f06 [5]uint
			f07 [5]float32
			f08 [5]float64
			f09 [5]bool
			f10 [5]string
			f11 [5]byte
			f12 int
			f13 int8
			f14 int16
			f15 int32
			f16 int64
			f17 uint
			f18 float32
			f19 float64
			f20 bool
			f21 string
			f22 byte
		}

		want := RoundTripStruct{
			f01: [5]int{1, 2, 3, 4, 5},
			f02: [5]int8{1, 2, 3, 4, 5},
			f03: [5]int16{1, 2, 3, 4, 5},
			f04: [5]int32{1, 2, 3, 4, 5},
			f05: [5]int64{1, 2, 3, 4, 5},
			f06: [5]uint{1, 2, 3, 4, 5},
			f07: [5]float32{1.5, 2.5, 3.5, 4.5, 5.5},
			f08: [5]float64{1.5, 2.5, 3.5, 4.5, 5.5},
			f09: [5]bool{true, false, true, false, true},
			f10: [5]string{"a", "b", "c", "d", "e"},
			f11: [5]byte{1, 2, 3, 4, 5},
			f12: 123,
			f13: 12,
			f14: 1234,
			f15: 123456,
			f16: 1234567890,
			f17: 123,
			f18: 1.5,
			f19: 2.5,
			f20: true,
			f21: "fuzz",
			f22: 'x',
		}

		structPtr := arena.Alloc[RoundTripStruct](a)
		if structPtr == nil {
			t.Error("Failed to allocate complex struct")
		} else {
			*structPtr = want
			if *structPtr != want {
				t.Error("struct did not survive a round trip through the arena")
			}
		}

		// Reset to ensure cleanup
		a.Reset()
	})
}

// TestBump_Vec_NativeTypes covers bump vec native types.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_Vec_NativeTypes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
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

// BenchmarkBump_Allocations measures bump allocations.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func BenchmarkBump_Allocations(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	b.ReportAllocs()

	for b.Loop() {
		// Allocate native types
		arena.Alloc[int](a)
		arena.Alloc[int8](a)
		arena.Alloc[int16](a)
		arena.Alloc[int32](a)
		arena.Alloc[int64](a)
		arena.Alloc[uint](a)
		arena.Alloc[float32](a)
		arena.Alloc[float64](a)
		arena.Alloc[bool](a)

		// Allocate arrays
		arena.Alloc[[100]int](a)
		arena.Alloc[[100]byte](a)
		arena.Alloc[[100]int8](a)
		arena.Alloc[[100]int16](a)
		arena.Alloc[[100]int32](a)
		arena.Alloc[[100]int64](a)
		arena.Alloc[[100]uint](a)
		arena.Alloc[[100]float32](a)
		arena.Alloc[[100]float64](a)
		arena.Alloc[[100]bool](a)

		// Allocate slices
		arena.MakeSlice[int](a, 100, 200)
		arena.MakeSlice[string](a, 100, 200)
		arena.MakeSlice[int8](a, 100, 200)
		arena.MakeSlice[int16](a, 100, 200)
		arena.MakeSlice[int32](a, 100, 200)
		arena.MakeSlice[int64](a, 100, 200)
		arena.MakeSlice[uint](a, 100, 200)
		arena.MakeSlice[float32](a, 100, 200)
		arena.MakeSlice[float64](a, 100, 200)
		arena.MakeSlice[bool](a, 100, 200)

		// Allocate strings
		a.MakeString("benchmark string")

		// Allocate structs
		arena.Alloc[MixedStruct](a)

		// Allocate empty struct
		type Empty struct{}

		arena.Alloc[Empty](a)

		a.Reset()
	}
}

// BenchmarkBump_StructAlloc measures bump struct alloc.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func BenchmarkBump_StructAlloc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(res.PAGE_SIZE))
	defer a.Delete()

	b.ReportAllocs()

	for b.Loop() {
		vec := container.NewVec[*MixedStruct](a)

		const count = 100_000
		for i := range count {
			ptr := arena.Alloc[MixedStruct](a)
			if ptr == nil {
				b.Fatal("Failed to allocate MixedStruct")
			}

			ptr.F12 = i
			vec.AppendOne(ptr)
		}

		if vec.Len() != count {
			b.Fatalf("Vec length mismatch: expected %d, got %d", count, vec.Len())
		}

		a.Reset()
	}
}

// TestBump_Allocator covers bump allocator.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_Allocator(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	p1 := arena.Alloc[int](a)
	if p1 == nil {
		t.Fatal("alloc failed")
	}

	p2 := arena.Alloc[int](a)
	if p2 == nil {
		t.Fatal("alloc failed")
	}

	if p1 == p2 {
		t.Fatal("same pointer")
	}

	a.Reset()

	p3 := arena.Alloc[int](a)
	if p3 != p1 {
		t.Fatal("not reset")
	}
}

// TestBump_AllocatorVariousSizes covers bump allocator various sizes.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_AllocatorVariousSizes(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(10 * 4096)) // 10 pages for larger allocations

	// Test allocating different basic types
	intPtr := arena.Alloc[int](a)
	if intPtr == nil {
		t.Fatal("failed to alloc int")
	}

	*intPtr = 42

	int64Ptr := arena.Alloc[int64](a)
	if int64Ptr == nil {
		t.Fatal("failed to alloc int64")
	}

	*int64Ptr = 123456789

	float64Ptr := arena.Alloc[float64](a)
	if float64Ptr == nil {
		t.Fatal("failed to alloc float64")
	}

	*float64Ptr = 3.14159

	// Test allocating a struct
	type SmallStruct struct {
		A int
		B string
		C [10]int
	}

	structPtr := arena.Alloc[SmallStruct](a)
	if structPtr == nil {
		t.Fatal("failed to alloc struct")
	}

	structPtr.A = 1
	structPtr.B = "test"

	// Test allocating slices of different sizes
	slice1 := arena.MakeSlice[int](a, 5, 10)
	if len(slice1) != 5 || cap(slice1) != 10 {
		t.Fatalf("slice1: len=%d cap=%d, expected len=5 cap=10", len(slice1), cap(slice1))
	}

	for i := range slice1 {
		slice1[i] = i
	}

	slice2 := arena.MakeSlice[string](a, 3, 5)
	if len(slice2) != 3 || cap(slice2) != 5 {
		t.Fatalf("slice2: len=%d cap=%d, expected len=3 cap=5", len(slice2), cap(slice2))
	}

	slice2[0] = "hello"
	slice2[1] = "world"

	// Test allocating a larger slice
	largeSlice := arena.MakeSlice[byte](a, 1000, 2000)
	if len(largeSlice) != 1000 || cap(largeSlice) != 2000 {
		t.Fatalf(
			"largeSlice: len=%d cap=%d, expected len=1000 cap=2000",
			len(largeSlice),
			cap(largeSlice),
		)
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
		t.Fatalf(
			"structPtr values: A=%d B=%s, expected A=1 B=test",
			structPtr.A,
			structPtr.B,
		)
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

// TestBump_AllocatorGrow covers bump allocator grow.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_AllocatorGrow(t *testing.T) {
	// Create a small arena with only 1 page (typically 4096 bytes)
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	// Allocate a large slice that exceeds the initial arena size
	// 1000 * 8 bytes (sizeof(int)) = 8000 bytes > 4096 bytes, forcing growth
	largeSlice := arena.MakeSlice[int](a, 1000, 1000)
	if len(largeSlice) != 1000 || cap(largeSlice) != 1000 {
		t.Fatalf(
			"largeSlice: len=%d cap=%d, expected 1000",
			len(largeSlice),
			cap(largeSlice),
		)
	}

	// Fill the slice
	for i := range largeSlice {
		largeSlice[i] = i * 2
	}

	// Allocate another item to ensure growth worked
	anotherPtr := arena.Alloc[int64](a)
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

// TestBump_100KInt64_FromTests covers bump 100k int 64 from tests.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_100KInt64_FromTests(t *testing.T) {
	// Allocate 100K int64 values using bump allocator
	a := arena.New(alloc.NewBumpAllocator(100 * 4096)) // Start with 100 pages
	defer a.Delete()

	const count = 100_000

	ptrs := make([]*int64, count)

	// Allocate 100K int64
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
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
		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}

// TestBump_1MInt64_FromTests covers bump 1m int 64 from tests.
//
// Revisions:
//   - 2025-12-20 23:37: initial creation
func TestBump_1MInt64_FromTests(t *testing.T) {
	// Allocate 1M int64 values using bump allocator
	a := arena.New(alloc.NewBumpAllocator(2000 * 4096)) // Start with 2000 pages (~8MB)
	defer a.Delete()

	const count = 1_000_000

	ptrs := make([]*int64, count)

	// Allocate 1M int64
	for i := 0; i < count; i++ {
		ptrs[i] = arena.Alloc[int64](a)
		if ptrs[i] == nil {
			t.Fatalf("allocation %d failed", i)
		}

		*ptrs[i] = int64(i)
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
		if *ptrs[i] != int64(i) {
			t.Errorf("index %d: expected %d, got %d", i, i, *ptrs[i])
		}
	}

	t.Logf("Successfully allocated and verified %d int64 values", count)
}
