package arena

const (
	// MAX_ALLOC_BYTES bounds a single allocation, so that a length multiplied by
	// an element size cannot wrap around and ask for a tiny block to hold a huge
	// slice.
	MAX_ALLOC_BYTES = 1 << 63
	// GROWTH_FACTOR is how much a slice's capacity grows when it is appended to.
	GROWTH_FACTOR = 2
	// MIN_SLICE_CAPACITY is the smallest capacity a grown slice is given, so that
	// the first few appends do not each move the backing block.
	MIN_SLICE_CAPACITY = 4
)
