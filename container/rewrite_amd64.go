//go:build amd64

//go:generate go run ./gen -out case_amd64.s

package container

var (
	// _USE_AVX2 is decided once. CPUID is far more expensive than the rewrite of a
	// short string, so the check does not belong on the path.
	_USE_AVX2 = _HasAVX2()
)

const (
	// SPEC_HI_SHIFT is where _CaseSpec stores the top of the range.
	SPEC_HI_SHIFT = 8
	// SPEC_DELTA_SHIFT is where _CaseSpec stores the size of the shift.
	SPEC_DELTA_SHIFT = 16
	// SPEC_SIGN_SHIFT is where _CaseSpec stores the sign of the shift.
	SPEC_SIGN_SHIFT = 24
	// SPEC_BYTE selects one byte of the packed spec.
	SPEC_BYTE = 0xFF
	// SPEC_SIGN is the sign bit inside the packed spec.
	SPEC_SIGN = uint64(1) << SPEC_SIGN_SHIFT
)

// _CaseSpec packs the span into one word so the kernels can take it by value.
// A pointer would escape into the assembly, and every conversion would allocate.
//
// Revisions:
//   - 2026-10-13 17:40: initial creation
func _CaseSpec(span *_CaseSpan) uint64 {
	spec := uint64(span.lo) |
		uint64(span.hi)<<SPEC_HI_SHIFT |
		_DeltaMagnitude(span.delta)<<SPEC_DELTA_SHIFT
	if span.delta < 0 {
		spec |= SPEC_SIGN
	}

	return spec
}

// _RewriteWide rewrites whole vectors of an ASCII tail and returns how many
// bytes it wrote. The rest, including any vector that holds a non-ASCII byte,
// is left for the word path.
//
// Revisions:
//   - 2026-10-13 17:25: initial creation
//   - 2026-10-13 17:40: the span is passed by value, so a short conversion does
//     not allocate
func _RewriteWide(dst []byte, src []byte, span *_CaseSpan) int {
	if len(src) < SSE_MIN || !_ShiftFits(span) || span.lo == 0 {
		return 0
	}

	spec := _CaseSpec(span)
	if _USE_AVX2 && len(src) >= AVX_MIN {
		return _RewriteAVX2(dst, src, spec)
	}

	return _RewriteSSE2(dst, src, spec)
}

// _RewriteSSE2 rewrites 16-byte chunks. It is generated.
//
// Revisions:
//   - 2026-10-13 17:25: initial creation
//   - 2026-10-13 17:40: take the span as a packed word
func _RewriteSSE2(dst []byte, src []byte, spec uint64) int

// _RewriteAVX2 rewrites 32-byte chunks. It is generated.
//
// Revisions:
//   - 2026-10-13 17:25: initial creation
//   - 2026-10-13 17:40: take the span as a packed word
func _RewriteAVX2(dst []byte, src []byte, spec uint64) int

// _HasAVX2 reports whether the CPU and the OS both accept AVX2.
//
// Revisions:
//   - 2026-10-13 17:25: initial creation
func _HasAVX2() bool
