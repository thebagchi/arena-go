package arena_test

const (
	// FIXTURE_ARRAY is the array length inside MixedStruct. It is what makes the
	// fixture large enough that allocating one is a realistic unit of work.
	FIXTURE_ARRAY = 100
)

// MixedStruct is the shared fixture for tests and benchmarks that need a large
// allocation whose layout is varied rather than a flat block of bytes.
//
// Its fields are the specification: the arrays make it large, and the string
// and array-of-string fields make it carry pointers, which is what exercises
// alignment and the rule that arena memory is not scanned by the collector.
// Nothing reads them, because the layout rather than the contents is the point,
// which is also why a benchmark must not write to them: that would measure the
// writes instead of the allocation.
//
// It was seven near-identical local declarations across this package before it
// was one.
type MixedStruct struct {
	F01 [FIXTURE_ARRAY]int
	F02 [FIXTURE_ARRAY]int8
	F03 [FIXTURE_ARRAY]int16
	F04 [FIXTURE_ARRAY]int32
	F05 [FIXTURE_ARRAY]int64
	F06 [FIXTURE_ARRAY]uint
	F07 [FIXTURE_ARRAY]float32
	F08 [FIXTURE_ARRAY]float64
	F09 [FIXTURE_ARRAY]bool
	F10 [FIXTURE_ARRAY]string
	F11 [FIXTURE_ARRAY]byte
	F12 int
	F13 int8
	F14 int16
	F15 int32
	F16 int64
	F17 uint
	F18 float32
	F19 float64
	F20 bool
	F21 string
	F22 byte
}
