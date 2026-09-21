package res

import "errors"

var (
	// ErrOutOfMemory is the value Alloc panics with when it cannot satisfy a
	// request: a fixed-size allocator is full, or the kernel refused a mapping.
	// Callers that expect exhaustion use TryAlloc instead of recovering.
	ErrOutOfMemory = errors.New("arena: out of memory")

	// ErrDeleted is the value every operation panics with after Delete. Reusing
	// a deleted arena is a lifetime bug in the caller, so it is reported rather
	// than absorbed; the previous behaviour was an index-out-of-range or a nil
	// dereference, which named neither the arena nor the mistake.
	ErrDeleted = errors.New("arena: used after Delete")
)
