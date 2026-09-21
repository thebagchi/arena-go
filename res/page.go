package res

import (
	"unsafe"
)

// Page is one mapped region, held together with its address bounds so that an
// ownership test is two integer comparisons rather than a slice dereference.
type Page struct {
	base  []byte
	start uintptr
	end   uintptr
}

// NewPage maps a region of at least size bytes, rounded up to a page boundary.
//
// Revisions:
//   - 2026-09-21 10:48: initial creation
func NewPage(size int) (*Page, error) {
	base, err := MakePages(size)
	if err != nil {
		return nil, err
	}

	return _NewPageOf(base), nil
}

// Base returns the page's bytes. The slice is the whole mapping, so it is also
// what ReleasePages must be given.
//
// Revisions:
//   - 2026-09-21 10:49: initial creation
func (p *Page) Base() []byte {
	return p.base
}

// Ptr returns a pointer to the first byte in the page.
//
// Pointers into a page are derived from this with unsafe.Add. Rebuilding one by
// converting a uintptr back is what go vet calls a possible misuse: the integer
// is not a reference, so nothing stops the value moving between the two casts.
//
// Revisions:
//   - 2026-09-21 11:52: initial creation
func (p *Page) Ptr() unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(p.base))
}

// Start returns the address of the first byte in the page.
//
// Revisions:
//   - 2026-09-21 10:49: initial creation
func (p *Page) Start() uintptr {
	return p.start
}

// End returns the address one past the last byte in the page.
//
// Revisions:
//   - 2026-09-21 10:50: initial creation
func (p *Page) End() uintptr {
	return p.end
}

// Size returns the number of bytes mapped, which is the requested size rounded
// up to a page boundary.
//
// Revisions:
//   - 2026-09-21 10:50: initial creation
func (p *Page) Size() int {
	return len(p.base)
}

// Contains reports whether addr falls inside this page.
//
// Revisions:
//   - 2026-09-21 10:51: initial creation
func (p *Page) Contains(addr uintptr) bool {
	return addr >= p.start && addr < p.end
}

// Release unmaps the page and leaves it empty. A released page reports that it
// contains nothing, so a stale reference cannot claim a pointer.
//
// Revisions:
//   - 2026-09-21 10:52: initial creation
func (p *Page) Release() error {
	if p.base == nil {
		return nil
	}

	err := ReleasePages(p.base)
	p.base, p.start, p.end = nil, 0, 0

	return err
}

// _NewPageOf wraps an existing mapping, caching its bounds once.
//
// Revisions:
//   - 2026-09-21 10:53: initial creation
func _NewPageOf(base []byte) *Page {
	start := uintptr(unsafe.Pointer(unsafe.SliceData(base)))

	return &Page{
		base:  base,
		start: start,
		end:   start + uintptr(len(base)),
	}
}
