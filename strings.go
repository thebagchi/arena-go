package arena

import (
	"bytes"
	"unsafe"
)

// ToBytes converts a string to a byte slice without copying (unsafe).
// Warning: Do not modify the returned slice, as it shares memory with the string.
func ToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ToString converts a byte slice to a string without copying (unsafe).
// Warning: Do not modify the original slice after conversion, as it shares memory with the string.
func ToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// TrimSpace trims whitespace from the string using bytes operations and returns an unsafe string sharing memory.
func TrimSpace(s string) string {
	b := ToBytes(s)
	trimmed := bytes.TrimSpace(b)
	return ToString(trimmed)
}

// IsEmpty checks if the string is empty or contains only whitespace.
func IsEmpty(s string) bool {
	b := ToBytes(s)
	trimmed := bytes.TrimSpace(b)
	return len(trimmed) == 0
}

// Contains checks if the string contains the substring without copying.
func Contains(s, substr string) bool {
	return bytes.Contains(ToBytes(s), ToBytes(substr))
}

// HasPrefix checks if the string starts with the prefix without copying.
func HasPrefix(s, prefix string) bool {
	return bytes.HasPrefix(ToBytes(s), ToBytes(prefix))
}

// HasSuffix checks if the string ends with the suffix without copying.
func HasSuffix(s, suffix string) bool {
	return bytes.HasSuffix(ToBytes(s), ToBytes(suffix))
}

// Index returns the index of the first occurrence of substr in s, or -1 if not found, without copying.
func Index(s, substr string) int {
	return bytes.Index(ToBytes(s), ToBytes(substr))
}

// LastIndex returns the index of the last occurrence of substr in s, or -1 if not found, without copying.
func LastIndex(s, substr string) int {
	return bytes.LastIndex(ToBytes(s), ToBytes(substr))
}

// Trim trims characters from the cutset from both ends of the string without copying.
func Trim(s string, cutset string) string {
	b := ToBytes(s)
	trimmed := bytes.Trim(b, cutset)
	return ToString(trimmed)
}

// TrimLeft trims characters from the cutset from the left end of the string without copying.
func TrimLeft(s string, cutset string) string {
	b := ToBytes(s)
	trimmed := bytes.TrimLeft(b, cutset)
	return ToString(trimmed)
}

// TrimRight trims characters from the cutset from the right end of the string without copying.
func TrimRight(s string, cutset string) string {
	b := ToBytes(s)
	trimmed := bytes.TrimRight(b, cutset)
	return ToString(trimmed)
}

// EqualFold performs case-insensitive comparison of two strings without copying.
func EqualFold(s, t string) bool {
	return bytes.EqualFold(ToBytes(s), ToBytes(t))
}

// Compare performs lexicographical comparison of two strings without copying.
func Compare(s, t string) int {
	return bytes.Compare(ToBytes(s), ToBytes(t))
}

// ToLower converts the string to lowercase.
// Note: This allocates a new byte slice for the result.
func ToLower(s string) string {
	b := ToBytes(s)
	lower := bytes.ToLower(b)
	return ToString(lower)
}

// ToUpper converts the string to uppercase.
// Note: This allocates a new byte slice for the result.
func ToUpper(s string) string {
	b := ToBytes(s)
	upper := bytes.ToUpper(b)
	return ToString(upper)
}

// Title capitalizes the first letter of each word.
// Note: This allocates a new byte slice for the result.
func Title(s string) string {
	b := ToBytes(s)
	title := bytes.Title(b)
	return ToString(title)
}

// Split splits the string by separator and allocates the result in the arena.
func Split(a *Arena, s, sep string) []string {
	if sep == "" {
		// Split into individual runes
		n := len(s)
		slice := MakeSlice[string](a, 0, n)
		for i := 0; i < n; i++ {
			slice = Append(a, slice, s[i:i+1])
		}
		return slice
	}

	// Count occurrences to pre-allocate with exact capacity
	n := bytes.Count(ToBytes(s), ToBytes(sep)) + 1
	slice := MakeSlice[string](a, 0, n)

	// Manually split to avoid intermediate allocations
	start := 0
	length := len(sep)
	for {
		idx := Index(s[start:], sep)
		if idx < 0 {
			// Add remaining part
			slice = Append(a, slice, s[start:])
			break
		}
		// Add part before separator
		slice = Append(a, slice, s[start:start+idx])
		start += idx + length
	}
	return slice
}

// Join joins the elements with separator and allocates the result in the arena.
func Join(a *Arena, elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	// Calculate total length
	length := len(sep) * (len(elems) - 1)
	for _, e := range elems {
		length += len(e)
	}
	// Allocate in arena
	var (
		data = MakeSlice[byte](a, length, length)
		pos  = 0
	)
	for i, e := range elems {
		if i > 0 {
			copy(data[pos:], sep)
			pos += len(sep)
		}
		copy(data[pos:], e)
		pos = pos + len(e)
	}
	return unsafe.String(&data[0], length)
}

// Fields splits the string on whitespace and allocates the result in the arena.
func Fields(a *Arena, s string) []string {
	// Fast path for empty string
	if len(s) == 0 {
		return nil
	}

	// Count fields first to pre-allocate with exact capacity
	n := 0
	inField := false
	for _, r := range s {
		wasInField := inField
		inField = r != ' ' && r != '\t' && r != '\n' && r != '\r'
		if inField && !wasInField {
			n++
		}
	}

	if n == 0 {
		return nil
	}

	slice := MakeSlice[string](a, 0, n)

	// Extract fields without intermediate allocations
	start := -1
	for i, r := range s {
		isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'
		if start < 0 {
			if !isSpace {
				start = i
			}
		} else if isSpace {
			slice = Append(a, slice, s[start:i])
			start = -1
		}
	}
	if start >= 0 {
		slice = Append(a, slice, s[start:])
	}

	return slice
}

// TrimPrefix removes the prefix from the string if present, without copying.
func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// TrimSuffix removes the suffix from the string if present, without copying.
func TrimSuffix(s, suffix string) string {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// Count counts the number of non-overlapping instances of substr in s.
func Count(s, substr string) int {
	return bytes.Count(ToBytes(s), ToBytes(substr))
}

// IndexByte returns the index of the first instance of byte c in s, or -1 if not found.
func IndexByte(s string, c byte) int {
	return bytes.IndexByte(ToBytes(s), c)
}

// LastIndexByte returns the index of the last instance of byte c in s, or -1 if not found.
func LastIndexByte(s string, c byte) int {
	return bytes.LastIndexByte(ToBytes(s), c)
}

// IndexAny returns the index of the first instance of any character from chars in s, or -1 if not found.
func IndexAny(s, chars string) int {
	return bytes.IndexAny(ToBytes(s), chars)
}

// LastIndexAny returns the index of the last instance of any character from chars in s, or -1 if not found.
func LastIndexAny(s, chars string) int {
	return bytes.LastIndexAny(ToBytes(s), chars)
}

// ContainsAny checks if the string contains any character from chars without copying.
func ContainsAny(s, chars string) bool {
	return bytes.ContainsAny(ToBytes(s), chars)
}

// ContainsRune checks if the string contains the rune without copying.
func ContainsRune(s string, r rune) bool {
	return bytes.ContainsRune(ToBytes(s), r)
}

// Replace replaces the first n occurrences of old with new and allocates the result in the arena.
// If n < 0, all occurrences are replaced.
func Replace(a *Arena, s, old, new string, n int) string {
	result := bytes.Replace(ToBytes(s), ToBytes(old), ToBytes(new), n)
	return a.MakeString(ToString(result))
}

// ReplaceAll replaces all occurrences of old with new and allocates the result in the arena.
func ReplaceAll(a *Arena, s, old, new string) string {
	result := bytes.ReplaceAll(ToBytes(s), ToBytes(old), ToBytes(new))
	return a.MakeString(ToString(result))
}

// Repeat returns a new string consisting of count copies of s, allocated in the arena.
func Repeat(a *Arena, s string, count int) string {
	result := bytes.Repeat(ToBytes(s), count)
	return a.MakeString(ToString(result))
}

// Cut cuts s around the first instance of sep, returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func Cut(s, sep string) (before, after string, found bool) {
	i := Index(s, sep)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}

// CutPrefix returns s without the provided leading prefix string and reports whether it found the prefix.
// If s doesn't start with prefix, CutPrefix returns s, false.
func CutPrefix(s, prefix string) (after string, found bool) {
	if HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return s, false
}

// CutSuffix returns s without the provided trailing suffix string and reports whether it found the suffix.
// If s doesn't end with suffix, CutSuffix returns s, false.
func CutSuffix(s, suffix string) (before string, found bool) {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)], true
	}
	return s, false
}

// SplitN splits the string by separator with a maximum of n parts and allocates the result in the arena.
// If n < 0, there is no limit on the number of parts.
func SplitN(a *Arena, s, sep string, n int) []string {
	parts := bytes.SplitN(ToBytes(s), ToBytes(sep), n)
	slice := MakeSlice[string](a, 0, len(parts))
	for _, p := range parts {
		slice = Append(a, slice, a.MakeString(ToString(p)))
	}
	return slice
}

// SplitAfter splits the string after each instance of sep and allocates the result in the arena.
func SplitAfter(a *Arena, s, sep string) []string {
	parts := bytes.SplitAfter(ToBytes(s), ToBytes(sep))
	slice := MakeSlice[string](a, 0, len(parts))
	for _, p := range parts {
		slice = Append(a, slice, a.MakeString(ToString(p)))
	}
	return slice
}

// SplitAfterN splits the string after each instance of sep with a maximum of n parts and allocates the result in the arena.
func SplitAfterN(a *Arena, s, sep string, n int) []string {
	parts := bytes.SplitAfterN(ToBytes(s), ToBytes(sep), n)
	slice := MakeSlice[string](a, 0, len(parts))
	for _, p := range parts {
		slice = Append(a, slice, a.MakeString(ToString(p)))
	}
	return slice
}
