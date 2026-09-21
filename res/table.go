package res

import (
	"sort"
	"sync"
	"unsafe"
)

// PageTable owns a set of mapped pages and answers ownership questions about
// them.
//
// Find is a binary search over the pages in address order. The allocators call
// it on every free, where the linear scan it replaces cost O(pages) and, in the
// slab allocator, was performed twice per call.
//
// Pages are appended as they are mapped and put in order only when a lookup
// needs it. Inserting in place instead looks tidier and is quadratic: Linux
// hands out mmap addresses downward, so every new page sorts before all the
// others and the insert moves the whole slice.
//
// The mutex guards every field below it. A PageTable is safe for concurrent use;
// the allocators built on it are not, unless they say so.
type PageTable struct {
	mtx    sync.Mutex
	pages  []*Page
	mapped int
	sorted bool
}

// NewPageTable returns an empty table that has mapped nothing yet.
//
// Revisions:
//   - 2026-09-21 10:56: initial creation
func NewPageTable() *PageTable {
	return new(PageTable)
}

// New maps a page of at least size bytes and records it in address order.
//
// Revisions:
//   - 2026-09-21 10:57: initial creation
func (t *PageTable) New(size int) (*Page, error) {
	page, err := NewPage(size)
	if err != nil {
		return nil, err
	}

	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.pages = append(t.pages, page)
	t.mapped = t.mapped + page.Size()
	t.sorted = false

	return page, nil
}

// Find returns the page containing ptr, or nil and false when no page does.
//
// Revisions:
//   - 2026-09-21 10:59: initial creation
func (t *PageTable) Find(ptr unsafe.Pointer) (*Page, bool) {
	if ptr == nil {
		return nil, false
	}

	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t._Find(uintptr(ptr))
}

// Owns reports whether ptr falls inside any page this table holds.
//
// Revisions:
//   - 2026-09-21 11:00: initial creation
func (t *PageTable) Owns(ptr unsafe.Pointer) bool {
	_, found := t.Find(ptr)

	return found
}

// Size returns the total bytes mapped, which is what a test asserting that a
// Reset cycle does not leak reads.
//
// Revisions:
//   - 2026-09-21 11:01: initial creation
func (t *PageTable) Size() int {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.mapped
}

// Len returns the number of pages held.
//
// Revisions:
//   - 2026-09-21 11:02: initial creation
func (t *PageTable) Len() int {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return len(t.pages)
}

// Delete unmaps every page and empties the table. It keeps going after a failed
// unmap so one bad page cannot strand the rest, and reports the first error.
//
// Revisions:
//   - 2026-09-21 11:03: initial creation
func (t *PageTable) Delete() error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	var first error

	for _, page := range t.pages {
		if err := page.Release(); err != nil && first == nil {
			first = err
		}
	}

	t.pages, t.mapped = nil, 0

	return first
}

// _Find is Find without the lock, for callers that already hold it.
//
// Revisions:
//   - 2026-09-21 11:04: initial creation
func (t *PageTable) _Find(addr uintptr) (*Page, bool) {
	t._Sort()

	idx := sort.Search(len(t.pages), func(i int) bool {
		return t.pages[i].end > addr
	})
	if idx < len(t.pages) && t.pages[idx].Contains(addr) {
		return t.pages[idx], true
	}

	return nil, false
}

// _Sort puts the pages in address order, which a binary search needs and an
// append does not maintain.
//
// Revisions:
//   - 2026-09-22 00:20: initial creation
func (t *PageTable) _Sort() {
	if t.sorted {
		return
	}

	sort.Slice(t.pages, func(i, j int) bool {
		return t.pages[i].start < t.pages[j].start
	})

	t.sorted = true
}
