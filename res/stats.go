package res

// Stats reports what an allocator has taken from the operating system.
//
// Mappings matters separately from Mapped because an mmap call is a syscall
// costing microseconds, while the bytes it returns cost nothing until touched:
// an allocator that maps the right number of bytes in the wrong number of calls
// looks efficient by Mapped alone. It is also what lets a test assert a
// footprint without reaching into an allocator's unexported fields.
type Stats struct {
	Mapped   int
	Mappings int
}
