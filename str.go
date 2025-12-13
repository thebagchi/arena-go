package arena

import (
	"bytes"
	"iter"
	"unicode"
	"unicode/utf8"
)

// Str is a string utility struct that holds an arena reference for memory management.
type Str struct {
	arena *Arena
}

// NewStr creates a new Str instance with the given arena.
func NewStr(a *Arena) *Str {
	return &Str{arena: a}
}

// TrimSpace trims whitespace from the string using bytes operations and returns an unsafe string sharing memory.
func (s *Str) TrimSpace(str string) string {
	return UnsafeString(bytes.TrimSpace(UnsafeBytes(str)))
}

// IsEmpty checks if the string is empty or contains only whitespace.
func (s *Str) IsEmpty(str string) bool {
	return len(bytes.TrimSpace(UnsafeBytes(str))) == 0
}

// Contains checks if the string contains the substring without copying.
func (s *Str) Contains(str, substr string) bool {
	return bytes.Contains(UnsafeBytes(str), UnsafeBytes(substr))
}

// HasPrefix checks if the string starts with the prefix without copying.
func (s *Str) HasPrefix(str, prefix string) bool {
	return bytes.HasPrefix(UnsafeBytes(str), UnsafeBytes(prefix))
}

// HasSuffix checks if the string ends with the suffix without copying.
func (s *Str) HasSuffix(str, suffix string) bool {
	return bytes.HasSuffix(UnsafeBytes(str), UnsafeBytes(suffix))
}

// Index returns the index of the first occurrence of substr in str, or -1 if not found, without copying.
func (s *Str) Index(str, substr string) int {
	return bytes.Index(UnsafeBytes(str), UnsafeBytes(substr))
}

// LastIndex returns the index of the last occurrence of substr in str, or -1 if not found, without copying.
func (s *Str) LastIndex(str, substr string) int {
	return bytes.LastIndex(UnsafeBytes(str), UnsafeBytes(substr))
}

// Trim trims characters from the cutset from both ends of the string without copying.
func (s *Str) Trim(str string, cutset string) string {
	return UnsafeString(bytes.Trim(UnsafeBytes(str), cutset))
}

// TrimLeft trims characters from the cutset from the left end of the string without copying.
func (s *Str) TrimLeft(str string, cutset string) string {
	return UnsafeString(bytes.TrimLeft(UnsafeBytes(str), cutset))
}

// TrimRight trims characters from the cutset from the right end of the string without copying.
func (s *Str) TrimRight(str string, cutset string) string {
	return UnsafeString(bytes.TrimRight(UnsafeBytes(str), cutset))
}

// EqualFold performs case-insensitive comparison of two strings without copying.
func (s *Str) EqualFold(str, t string) bool {
	return bytes.EqualFold(UnsafeBytes(str), UnsafeBytes(t))
}

// Compare performs lexicographical comparison of two strings without copying.
func (s *Str) Compare(str, t string) int {
	return bytes.Compare(UnsafeBytes(str), UnsafeBytes(t))
}

// ToLower converts the string to lowercase.
// Returns the original string without allocation if already lowercase.
func (s *Str) ToLower(str string) string {
	// Fast path: check if already lowercase
	var (
		b       = UnsafeBytes(str)
		convert = false
	)
	for _, c := range b {
		if c >= 'A' && c <= 'Z' {
			convert = true
			break
		}
	}
	if !convert {
		return str
	}

	// Convert using Buffer
	buf := NewBuffer(s.arena)
	for _, c := range b {
		if c >= 'A' && c <= 'Z' {
			buf.Append([]byte{c + 32})
		} else {
			buf.Append([]byte{c})
		}
	}
	return buf.String()
}

// ToUpper converts the string to uppercase.
// Returns the original string without allocation if already uppercase.
func (s *Str) ToUpper(str string) string {
	// Fast path: check if already uppercase
	var (
		b       = UnsafeBytes(str)
		convert = false
	)
	for _, c := range b {
		if c >= 'a' && c <= 'z' {
			convert = true
			break
		}
	}
	if !convert {
		return str
	}

	// Convert using Buffer
	buf := NewBuffer(s.arena)
	for _, c := range b {
		if c >= 'a' && c <= 'z' {
			buf.Append([]byte{c - 32})
		} else {
			buf.Append([]byte{c})
		}
	}
	return buf.String()
}

// Title capitalizes the first letter of each word.
// Returns the original string without allocation if already title case.
func (s *Str) Title(str string) string {
	// Fast path: check if already title case (simplified check)
	// This is a basic check - if no lowercase letters at word starts, might be title case
	var (
		prevWasSpace    = true
		needsConversion = false
	)
	for _, r := range str {
		if prevWasSpace && r >= 'a' && r <= 'z' {
			needsConversion = true
			break
		}
		prevWasSpace = r == ' ' || r == '\t' || r == '\n' || r == '\r'
	}
	if !needsConversion {
		return str
	}

	// Convert with proper title casing using Buffer
	var (
		buf         = NewBuffer(s.arena)
		isWordStart = true
		runeBuf     [utf8.UTFMax]byte
	)
	for _, r := range str {
		if isWordStart && unicode.IsLetter(r) {
			n := utf8.EncodeRune(runeBuf[:], unicode.ToTitle(r))
			isWordStart = false
			buf.Append(runeBuf[:n])
		} else {
			n := utf8.EncodeRune(runeBuf[:], r)
			buf.Append(runeBuf[:n])
			if unicode.IsSpace(r) {
				isWordStart = true
			} else {
				isWordStart = false
			}
		}
	}
	return buf.String()
}

// Split splits the string by separator and allocates the result in the arena.
func (s *Str) Split(str, sep string) []string {
	if sep == "" {
		// Split into individual runes
		var (
			n     = utf8.RuneCountInString(str)
			slice = MakeSlice[string](s.arena, 0, n)
		)
		for _, r := range str {
			slice = Append(s.arena, slice, string(r))
		}
		return slice
	}

	// Count occurrences to pre-allocate with exact capacity
	var (
		n     = bytes.Count(UnsafeBytes(str), UnsafeBytes(sep)) + 1
		slice = MakeSlice[string](s.arena, 0, n)
	)

	// Manually split to avoid intermediate allocations
	var (
		start  = 0
		length = len(sep)
	)
	for {
		idx := s.Index(str[start:], sep)
		if idx < 0 {
			// Add remaining part
			slice = Append(s.arena, slice, str[start:])
			break
		}
		// Add part before separator
		slice = Append(s.arena, slice, str[start:start+idx])
		start = start + (idx + length)
	}
	return slice
}

// Join joins the elements with separator and allocates the result in the arena.
func (s *Str) Join(elems []string, sep string) string {
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
		data = MakeSlice[byte](s.arena, length, length)
		pos  = 0
	)
	for i, e := range elems {
		if i > 0 {
			copy(data[pos:], sep)
			pos = pos + len(sep)
		}
		copy(data[pos:], e)
		pos = pos + len(e)
	}
	return UnsafeString(data)
}

// Fields splits the string on whitespace and allocates the result in the arena.
func (s *Str) Fields(str string) []string {
	// Fast path for empty string
	if len(str) == 0 {
		return nil
	}

	// Count fields first to pre-allocate with exact capacity
	var (
		n       = 0
		inField = false
	)
	for _, r := range str {
		wasInField := inField
		inField = r != ' ' && r != '\t' && r != '\n' && r != '\r'
		if inField && !wasInField {
			n = n + 1
		}
	}

	if n == 0 {
		return nil
	}

	var (
		slice = MakeSlice[string](s.arena, 0, n)
		start = -1
	)
	for i, r := range str {
		isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'
		if start < 0 {
			if !isSpace {
				start = i
			}
		} else if isSpace {
			slice = Append(s.arena, slice, str[start:i])
			start = -1
		}
	}
	if start >= 0 {
		slice = Append(s.arena, slice, str[start:])
	}

	return slice
}

// TrimPrefix removes the prefix from the string if present, without copying.
func (s *Str) TrimPrefix(str, prefix string) string {
	if s.HasPrefix(str, prefix) {
		return str[len(prefix):]
	}
	return str
}

// TrimSuffix removes the suffix from the string if present, without copying.
func (s *Str) TrimSuffix(str, suffix string) string {
	if s.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	}
	return str
}

// Count counts the number of non-overlapping instances of substr in str.
func (s *Str) Count(str, substr string) int {
	return bytes.Count(UnsafeBytes(str), UnsafeBytes(substr))
}

// IndexByte returns the index of the first instance of byte c in str, or -1 if not found.
func (s *Str) IndexByte(str string, c byte) int {
	return bytes.IndexByte(UnsafeBytes(str), c)
}

// LastIndexByte returns the index of the last instance of byte c in str, or -1 if not found.
func (s *Str) LastIndexByte(str string, c byte) int {
	return bytes.LastIndexByte(UnsafeBytes(str), c)
}

// IndexAny returns the index of the first instance of any character from chars in str, or -1 if not found.
func (s *Str) IndexAny(str, chars string) int {
	return bytes.IndexAny(UnsafeBytes(str), chars)
}

// LastIndexAny returns the index of the last instance of any character from chars in str, or -1 if not found.
func (s *Str) LastIndexAny(str, chars string) int {
	return bytes.LastIndexAny(UnsafeBytes(str), chars)
}

// ContainsAny checks if the string contains any character from chars without copying.
func (s *Str) ContainsAny(str, chars string) bool {
	return bytes.ContainsAny(UnsafeBytes(str), chars)
}

// ContainsRune checks if the string contains the rune without copying.
func (s *Str) ContainsRune(str string, r rune) bool {
	return bytes.ContainsRune(UnsafeBytes(str), r)
}

// Replace replaces the first n occurrences of old with new and allocates the result in the arena.
// If n < 0, all occurrences are replaced.
func (s *Str) Replace(str, old, new string, n int) string {
	if n == 0 || old == "" {
		return str
	}

	buf := NewBuffer(s.arena)
	var (
		start = 0
		count = 0
	)

	for {
		idx := s.Index(str[start:], old)
		if idx < 0 {
			buf.AppendString(str[start:])
			break
		}
		buf.AppendString(str[start : start+idx])
		buf.AppendString(new)
		start = start + (idx + len(old))
		count = count + 1
		if n > 0 && count >= n {
			buf.AppendString(str[start:])
			break
		}
	}
	return buf.String()
}

// ReplaceAll replaces all occurrences of old with new and allocates the result in the arena.
func (s *Str) ReplaceAll(str, old, new string) string {
	return s.Replace(str, old, new, -1)
}

// Repeat returns a new string consisting of count copies of str, allocated in the arena.
func (s *Str) Repeat(str string, count int) string {
	if count <= 0 {
		return ""
	}
	if count == 1 {
		return str
	}

	buf := NewBuffer(s.arena)
	for range count {
		buf.AppendString(str)
	}
	return buf.String()
}

// Cut cuts str around the first instance of sep, returning the text before and after sep.
// The found result reports whether sep appears in str.
// If sep does not appear in str, cut returns str, "", false.
func (s *Str) Cut(str, sep string) (before, after string, found bool) {
	i := s.Index(str, sep)
	if i < 0 {
		return str, "", false
	}
	return str[:i], str[i+len(sep):], true
}

// CutPrefix returns str without the provided leading prefix string and reports whether it found the prefix.
// If str doesn't start with prefix, CutPrefix returns str, false.
func (s *Str) CutPrefix(str, prefix string) (after string, found bool) {
	if s.HasPrefix(str, prefix) {
		return str[len(prefix):], true
	}
	return str, false
}

// CutSuffix returns str without the provided trailing suffix string and reports whether it found the suffix.
// If str doesn't end with suffix, CutSuffix returns str, false.
func (s *Str) CutSuffix(str, suffix string) (before string, found bool) {
	if s.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)], true
	}
	return str, false
}

// SplitN splits the string by separator with a maximum of n parts and allocates the result in the arena.
// If n < 0, there is no limit on the number of parts.
func (s *Str) SplitN(str, sep string, n int) []string {
	var (
		parts = bytes.SplitN(UnsafeBytes(str), UnsafeBytes(sep), n)
		slice = MakeSlice[string](s.arena, 0, len(parts))
	)
	for _, p := range parts {
		slice = Append(s.arena, slice, UnsafeString(p))
	}
	return slice
}

// SplitAfter splits the string after each instance of sep and allocates the result in the arena.
func (s *Str) SplitAfter(str, sep string) []string {
	var (
		parts = bytes.SplitAfter(UnsafeBytes(str), UnsafeBytes(sep))
		slice = MakeSlice[string](s.arena, 0, len(parts))
	)
	for _, p := range parts {
		slice = Append(s.arena, slice, UnsafeString(p))
	}
	return slice
}

// SplitAfterN splits the string after each instance of sep with a maximum of n parts and allocates the result in the arena.
func (s *Str) SplitAfterN(str, sep string, n int) []string {
	var (
		parts = bytes.SplitAfterN(UnsafeBytes(str), UnsafeBytes(sep), n)
		slice = MakeSlice[string](s.arena, 0, len(parts))
	)
	for _, p := range parts {
		slice = Append(s.arena, slice, UnsafeString(p))
	}
	return slice
}

// Lines returns an iterator over the newline-terminated lines in the string str.
// The lines yielded by the iterator include their terminating newlines.
// If str is empty, the iterator yields no lines at all.
// If str does not end in a newline, the final yielded line will not end in a newline.
func (s *Str) Lines(str string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if len(str) == 0 {
			return
		}
		start := 0
		for i, r := range str {
			if r == '\n' {
				if !yield(str[start : i+1]) {
					return
				}
				start = i + 1
			}
		}
		if start < len(str) {
			yield(str[start:])
		}
	}
}

// Clone returns a copy of the string, allocated in the arena.
func (s *Str) Clone(str string) string {
	return s.arena.MakeString(str)
}

// FieldsFunc splits the string str at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of str allocated in the arena.
// If all code points in str satisfy f(c) or the string is empty, an empty slice is returned.
func (s *Str) FieldsFunc(str string, f func(rune) bool) []string {
	// Count fields first
	var (
		n       = 0
		inField = false
	)
	for _, r := range str {
		wasInField := inField
		inField = !f(r)
		if inField && !wasInField {
			n = n + 1
		}
	}

	if n == 0 {
		return nil
	}

	var (
		slice = MakeSlice[string](s.arena, 0, n)
		start = -1
	)
	for i, r := range str {
		if start < 0 {
			if !f(r) {
				start = i
			}
		} else if f(r) {
			slice = Append(s.arena, slice, str[start:i])
			start = -1
		}
	}
	if start >= 0 {
		slice = Append(s.arena, slice, str[start:])
	}

	return slice
}

// ContainsFunc reports whether any Unicode code point in str satisfies f(r).
func (s *Str) ContainsFunc(str string, f func(rune) bool) bool {
	for _, r := range str {
		if f(r) {
			return true
		}
	}
	return false
}

// IndexFunc returns the index into str of the first Unicode code point satisfying f(c),
// or -1 if none do.
func (s *Str) IndexFunc(str string, f func(rune) bool) int {
	for i, r := range str {
		if f(r) {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the index into str of the last Unicode code point satisfying f(c),
// or -1 if none do.
func (s *Str) LastIndexFunc(str string, f func(rune) bool) int {
	for i := len(str) - 1; i >= 0; i-- {
		r, size := utf8.DecodeLastRuneInString(str[:i+1])
		if f(r) {
			return i - (size - 1)
		}
		i = i - (size - 1)
	}
	return -1
}

// MapASCII returns a copy of the string str with all its bytes modified according to the mapping function.
// This is optimized for ASCII-only strings. If mapping returns a negative value, the byte is dropped.
// For full Unicode support, use MapUTF8.
func (s *Str) MapASCII(mapping func(byte) int, str string) string {
	var (
		b   = UnsafeBytes(str)
		buf = NewBuffer(s.arena)
	)

	for _, c := range b {
		mapped := mapping(c)
		if mapped >= 0 {
			buf.Append([]byte{byte(mapped)})
		}
	}
	return buf.String()
}

// MapUTF8 returns a copy of the string str with all its characters modified according to the mapping function.
// If mapping returns a negative value, the character is dropped from the string with no replacement.
// This handles full Unicode correctly. For ASCII-only strings, MapASCII is faster.
func (s *Str) MapUTF8(mapping func(rune) rune, str string) string {
	var (
		buf  = NewBuffer(s.arena)
		temp [utf8.UTFMax]byte
	)

	for _, r := range str {
		mapped := mapping(r)
		if mapped >= 0 {
			n := utf8.EncodeRune(temp[:], mapped)
			buf.Append(temp[:n])
		}
	}
	return buf.String()
}

// MapString is an alias for MapUTF8 for backward compatibility.
func (s *Str) MapString(mapping func(rune) rune, str string) string {
	return s.MapUTF8(mapping, str)
}

// ToTitle returns a copy of the string str with all Unicode letters mapped to their Unicode title case.
// The result is allocated in the arena.
func (s *Str) ToTitle(str string) string {
	result := s.MapUTF8(unicode.ToTitle, str)
	return s.arena.MakeString(result)
}

// ToValidUTF8 returns a copy of the string str with each run of invalid UTF-8 byte sequences
// replaced by the replacement string, which may be empty.
// The result is allocated in the arena.
func (s *Str) ToValidUTF8(str, replacement string) string {
	// Fast path: check if string is already valid UTF-8
	if utf8.ValidString(str) {
		return str
	}

	buf := NewBuffer(s.arena)
	for i := 0; i < len(str); {
		r1, size1 := utf8.DecodeRuneInString(str[i:])
		if r1 == utf8.RuneError && size1 == 1 {
			// Invalid UTF-8 sequence
			buf.AppendString(replacement)
			i = i + 1
			// Skip consecutive invalid bytes
			for i < len(str) {
				r2, size2 := utf8.DecodeRuneInString(str[i:])
				if r2 != utf8.RuneError || size2 != 1 {
					break
				}
				i = i + 1
			}
		} else {
			// Valid UTF-8 rune
			var (
				temp [utf8.UTFMax]byte
				n    = utf8.EncodeRune(temp[:], r1)
			)
			buf.Append(temp[:n])
			i = i + size1
		}
	}
	return buf.String()
}

// TrimFunc returns a slice of the string str with all leading and trailing
// Unicode code points c satisfying f(c) removed.
func (s *Str) TrimFunc(str string, f func(rune) bool) string {
	return s.TrimRightFunc(s.TrimLeftFunc(str, f), f)
}

// TrimLeftFunc returns a slice of the string str with all leading
// Unicode code points c satisfying f(c) removed.
func (s *Str) TrimLeftFunc(str string, f func(rune) bool) string {
	for i, r := range str {
		if !f(r) {
			return str[i:]
		}
	}
	return ""
}

// TrimRightFunc returns a slice of the string str with all trailing
// Unicode code points c satisfying f(c) removed.
func (s *Str) TrimRightFunc(str string, f func(rune) bool) string {
	for i := len(str) - 1; i >= 0; i-- {
		r, size := utf8.DecodeLastRuneInString(str[:i+1])
		if !f(r) {
			return str[:i+1]
		}
		i = i - (size - 1)
	}
	return ""
}
