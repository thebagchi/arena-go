package arena_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/thebagchi/arena-go/res"
)

// TestRoundPow2 checks the rounding helper, including the zero that used to
// underflow and return zero rather than one.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestRoundPow2(t *testing.T) {
	cases := []struct {
		name string
		in   uint64
		want uint64
	}{
		{name: "zero", in: 0, want: 1},
		{name: "one", in: 1, want: 1},
		{name: "two", in: 2, want: 2},
		{name: "three", in: 3, want: 4},
		{name: "already a power of two", in: 4096, want: 4096},
		{name: "one over", in: 4097, want: 8192},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := res.RoundPow2(tc.in); got != tc.want {
				t.Errorf("RoundPow2(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestLog2 checks the log helper, including the zero that used to wrap to the
// maximum uint64.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestLog2(t *testing.T) {
	cases := []struct {
		name string
		in   uint64
		want uint64
	}{
		{name: "zero", in: 0, want: 0},
		{name: "one", in: 1, want: 0},
		{name: "sixteen", in: 16, want: 4},
		{name: "not a power of two", in: 17, want: 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := res.Log2(tc.in); got != tc.want {
				t.Errorf("Log2(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestIsPow2 checks the alignment predicate every masking calculation relies on.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestIsPow2(t *testing.T) {
	cases := []struct {
		name string
		in   uint64
		want bool
	}{
		{name: "zero", in: 0, want: false},
		{name: "one", in: 1, want: true},
		{name: "eight", in: 8, want: true},
		{name: "twelve", in: 12, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := res.IsPow2(tc.in); got != tc.want {
				t.Errorf("IsPow2(%d) = %t, want %t", tc.in, got, tc.want)
			}
		})
	}
}

// TestPage checks that a page reports the bounds it was mapped with and stops
// claiming pointers once released.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestPage(t *testing.T) {
	page, err := res.NewPage(res.PAGE_SIZE)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}

	if page.Size() < res.PAGE_SIZE {
		t.Errorf("Size = %d, want at least %d", page.Size(), res.PAGE_SIZE)
	}

	if !page.Contains(page.Start()) {
		t.Error("page does not contain its own start")
	}

	if page.Contains(page.End()) {
		t.Error("page contains the address one past its end")
	}

	start := page.Start()
	if err := page.Release(); err != nil {
		t.Fatalf("Release: %v", err)
	}

	if page.Contains(start) {
		t.Error("released page still claims its old start")
	}
}

// TestPageTableFind checks that a table finds the page holding a pointer and
// rejects one from elsewhere.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestPageTableFind(t *testing.T) {
	table := res.NewPageTable()

	defer func() {
		if err := table.Delete(); err != nil {
			t.Errorf("Delete: %v", err)
		}
	}()

	const pages = 8

	for range pages {
		if _, err := table.New(res.PAGE_SIZE); err != nil {
			t.Fatalf("New: %v", err)
		}
	}

	if table.Len() != pages {
		t.Errorf("Len = %d, want %d", table.Len(), pages)
	}

	if table.Size() < pages*res.PAGE_SIZE {
		t.Errorf("Size = %d, want at least %d", table.Size(), pages*res.PAGE_SIZE)
	}

	page, err := table.New(res.PAGE_SIZE)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if !table.Owns(page.Ptr()) {
		t.Error("table does not own a page it mapped")
	}

	outside := new(int)
	if table.Owns(unsafe.Pointer(outside)) {
		t.Error("table claims a heap pointer")
	}
}

// TestPageTableDeleteEmpties checks that Delete releases everything and leaves
// the table claiming nothing.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestPageTableDeleteEmpties(t *testing.T) {
	table := res.NewPageTable()

	page, err := table.New(res.PAGE_SIZE)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	start := page.Start()

	if err := table.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if table.Len() != 0 || table.Size() != 0 {
		t.Errorf(
			"after Delete: Len = %d, Size = %d, want 0 and 0",
			table.Len(),
			table.Size(),
		)
	}

	if page.Contains(start) {
		t.Error("page released by Delete still claims its old start")
	}
}

// TestBumpAllocAlignment checks that every alignment a caller can ask for is
// honoured, including one larger than the value being stored.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestBumpAllocAlignment(t *testing.T) {
	bump := res.NewBump(res.PAGE_SIZE)

	defer func() {
		if err := bump.Delete(); err != nil {
			t.Errorf("Delete: %v", err)
		}
	}()

	for _, align := range []uint64{1, 2, 4, 8, 16, 32, 64, 128} {
		t.Run(fmt.Sprintf("align %d", align), func(t *testing.T) {
			// One byte at a time in between, so the cursor is left at an odd
			// offset and the alignment has real work to do.
			_ = bump.Alloc(1, 1)

			ptr := bump.Alloc(8, align)
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

// TestBumpAllocDistinct checks that separate allocations never overlap.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestBumpAllocDistinct(t *testing.T) {
	bump := res.NewBump(res.PAGE_SIZE)

	defer func() {
		if err := bump.Delete(); err != nil {
			t.Errorf("Delete: %v", err)
		}
	}()

	const count = 4096

	seen := make(map[uintptr]bool, count)

	for i := range count {
		ptr := bump.Alloc(64, 8)
		if seen[uintptr(ptr)] {
			t.Fatalf("allocation %d repeated address %p", i, ptr)
		}

		seen[uintptr(ptr)] = true

		if !bump.Owns(ptr) {
			t.Fatalf("allocation %d is not owned by the allocator that made it", i)
		}
	}
}

// TestBumpGrowsGeometrically checks that a working set much larger than the
// first chunk costs a handful of mappings rather than one per chunk's worth.
//
// Revisions:
//   - 2025-12-31 18:33: initial creation
func TestBumpGrowsGeometrically(t *testing.T) {
	bump := res.NewBump(res.PAGE_SIZE)

	defer func() {
		if err := bump.Delete(); err != nil {
			t.Errorf("Delete: %v", err)
		}
	}()

	const (
		total   = 4 << 20
		perCall = 1024
		limit   = 16
	)

	for range total / perCall {
		_ = bump.Alloc(perCall, 8)
	}

	if bump.Mappings() > limit {
		t.Errorf(
			"%d bytes took %d mappings, want at most %d",
			total,
			bump.Mappings(),
			limit,
		)
	}
}
