package res

import (
	"math/bits"
	"unsafe"
)

// RoundPow2 returns the smallest power of two greater than or equal to n.
//
// Zero and one both round to one, so callers never receive a block size of
// zero; the previous formula underflowed on n == 0 and returned it.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func RoundPow2(n uint64) uint64 {
	if n <= 1 {
		return 1
	}

	return uint64(1) << uint(bits.Len64(n-1))
}

// Log2 returns the index of n's highest set bit, which is log2(n) for a power of
// two. Zero returns zero rather than wrapping to the maximum uint64.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func Log2(n uint64) uint64 {
	if n == 0 {
		return 0
	}

	return uint64(bits.Len64(n) - 1)
}

// IsPow2 reports whether n is a power of two. Alignment arguments must satisfy
// it, because every alignment calculation here masks with n-1.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func IsPow2(n uint64) bool {
	return n != 0 && n&(n-1) == 0
}

// RoundUp returns size rounded up to the next multiple of unit.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func RoundUp(size, unit int) int {
	if unit <= 1 {
		return size
	}

	return ((size + unit - 1) / unit) * unit
}

// AlignUp returns addr rounded up to the next multiple of align, which must be a
// power of two.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func AlignUp(addr, align uintptr) uintptr {
	return (addr + align - 1) &^ (align - 1)
}

// SlicePtr returns a pointer to a slice's backing array, or nil when the slice
// is empty.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func SlicePtr[T any](slice []T) unsafe.Pointer {
	if len(slice) == 0 {
		return nil
	}

	return unsafe.Pointer(unsafe.SliceData(slice))
}

// UnsafeBytes views a string's bytes as a slice without copying. The result
// aliases the string and must never be written to.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func UnsafeBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}

	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString views a byte slice as a string without copying. The caller keeps
// the promise that the bytes will not change for the string's lifetime.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func UnsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return unsafe.String(&b[0], len(b))
}

// Request normalises an allocation request: a zero size still needs a distinct
// address, a zero alignment means the default, and every alignment calculation
// here masks with align-1, which is only correct for a power of two.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func Request(size, align uint64) (uint64, uint64) {
	if size == 0 {
		size = 1
	}

	if align == 0 {
		align = DEFAULT_ALIGN
	}

	if !IsPow2(align) {
		align = RoundPow2(align)
	}

	return size, align
}
