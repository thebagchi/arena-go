package container

import (
	"encoding/binary"
	"iter"
	"strings"
	"unicode"
	"unicode/utf8"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/res"
)

const (
	// CASE_SHIFT is the distance between an ASCII letter's two cases.
	CASE_SHIFT = 'a' - 'A'
	// ALL_MATCHES asks a replace for every occurrence rather than a count.
	ALL_MATCHES = -1
	// WORD_BYTES is how many bytes the word-at-a-time scans handle per step.
	WORD_BYTES = 8
	// HIGH_BITS selects the top bit of every byte in a word, which is set for
	// exactly the bytes that are not ASCII.
	HIGH_BITS = 0x8080808080808080
	// ASCII_LIMIT is the first byte value that is not ASCII.
	ASCII_LIMIT = 0x80
	// ONE_PER_BYTE has the low bit of every byte set, for broadcasting a byte
	// across a word.
	ONE_PER_BYTE = 0x0101010101010101
	// ASCII_SPACE is a bit per whitespace byte: tab, newline, vertical tab, form
	// feed, carriage return and space. All six are below 64, so the whole set
	// fits in one word and the test needs no memory at all.
	ASCII_SPACE = 1<<'\t' | 1<<'\n' | 1<<'\v' | 1<<'\f' | 1<<'\r' | 1<<' '
	// SPACE_LIMIT is one past the highest byte ASCII_SPACE can answer for.
	SPACE_LIMIT = 64
	// FIELDS_GUESS is the capacity a single-pass split starts with, before it
	// grows in the arena.
	FIELDS_GUESS = 8
)

// Str holds an arena and offers the strings package's operations over it.
//
// Operations that only narrow their input return a string sharing the input's
// memory and allocate nothing. Operations that build a new string allocate it in
// the arena, so the result lives exactly as long as the arena does.
//
// Splitting returns parts that point into the input rather than copies of it. If
// the input is a Go heap string and the parts are stored in arena memory, the
// caller has to keep the input alive, because arena memory does not. Passing a
// string the arena already owns, from MakeString or Clone, removes the question.
//
// It is not safe for concurrent use, because the strings it returns share arena
// memory that a later call may reuse.
type Str struct {
	arena *arena.Arena
}

// NewStr returns a string helper backed by the arena.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func NewStr(a *arena.Arena) *Str {
	return &Str{arena: a}
}

// Clone copies s into the arena.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Clone(str string) string {
	return s.arena.MakeString(str)
}

// IsEmpty reports whether str is empty or only whitespace.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// TrimSpace removes leading and trailing whitespace without copying.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// Trim removes any leading and trailing characters in cutset.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Trim(str, cutset string) string {
	return strings.Trim(str, cutset)
}

// TrimLeft removes any leading characters in cutset.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimLeft(str, cutset string) string {
	return strings.TrimLeft(str, cutset)
}

// TrimRight removes any trailing characters in cutset.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimRight(str, cutset string) string {
	return strings.TrimRight(str, cutset)
}

// TrimPrefix removes prefix from the front of str when it is there.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimPrefix(str, prefix string) string {
	if s.HasPrefix(str, prefix) {
		return str[len(prefix):]
	}

	return str
}

// TrimSuffix removes suffix from the end of str when it is there.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimSuffix(str, suffix string) string {
	if s.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	}

	return str
}

// TrimFunc removes leading and trailing runes the predicate accepts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimFunc(str string, match func(rune) bool) string {
	return s.TrimRightFunc(s.TrimLeftFunc(str, match), match)
}

// TrimLeftFunc removes leading runes the predicate accepts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimLeftFunc(str string, match func(rune) bool) string {
	for i, r := range str {
		if !match(r) {
			return str[i:]
		}
	}

	return ""
}

// TrimRightFunc removes trailing runes the predicate accepts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) TrimRightFunc(str string, match func(rune) bool) string {
	for len(str) > 0 {
		r, size := utf8.DecodeLastRuneInString(str)
		if !match(r) {
			return str
		}

		str = str[:len(str)-size]
	}

	return ""
}

// Contains reports whether substr is inside str.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ContainsAny reports whether any rune of chars is inside str.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ContainsAny(str, chars string) bool {
	return strings.ContainsAny(str, chars)
}

// ContainsRune reports whether r is inside str.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ContainsRune(str string, r rune) bool {
	return strings.ContainsRune(str, r)
}

// ContainsFunc reports whether any rune of str satisfies the predicate.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ContainsFunc(str string, match func(rune) bool) bool {
	return s.IndexFunc(str, match) != NOT_FOUND
}

// HasPrefix reports whether str starts with prefix.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasSuffix reports whether str ends with suffix.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// Index returns the first position of substr, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Index(str, substr string) int {
	return strings.Index(str, substr)
}

// LastIndex returns the last position of substr, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) LastIndex(str, substr string) int {
	return strings.LastIndex(str, substr)
}

// IndexByte returns the first position of c, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) IndexByte(str string, c byte) int {
	return strings.IndexByte(str, c)
}

// LastIndexByte returns the last position of c, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) LastIndexByte(str string, c byte) int {
	return strings.LastIndexByte(str, c)
}

// IndexAny returns the first position of any rune in chars, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) IndexAny(str, chars string) int {
	return strings.IndexAny(str, chars)
}

// LastIndexAny returns the last position of any rune in chars, or NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) LastIndexAny(str, chars string) int {
	return strings.LastIndexAny(str, chars)
}

// IndexFunc returns the position of the first rune the predicate accepts, or
// NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) IndexFunc(str string, match func(rune) bool) int {
	for i, r := range str {
		if match(r) {
			return i
		}
	}

	return NOT_FOUND
}

// LastIndexFunc returns the position of the last rune the predicate accepts, or
// NOT_FOUND.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) LastIndexFunc(str string, match func(rune) bool) int {
	for len(str) > 0 {
		r, size := utf8.DecodeLastRuneInString(str)
		if match(r) {
			return len(str) - size
		}

		str = str[:len(str)-size]
	}

	return NOT_FOUND
}

// Count returns how many non-overlapping copies of substr are in str.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Count(str, substr string) int {
	return strings.Count(str, substr)
}

// Compare orders two strings lexicographically.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Compare(str, other string) int {
	return strings.Compare(str, other)
}

// EqualFold compares two strings ignoring case.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) EqualFold(str, other string) bool {
	return strings.EqualFold(str, other)
}

// Cut splits str around the first copy of sep.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Cut(str, sep string) (string, string, bool) {
	if i := s.Index(str, sep); i >= 0 {
		return str[:i], str[i+len(sep):], true
	}

	return str, "", false
}

// CutPrefix removes prefix, reporting whether it was there.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) CutPrefix(str, prefix string) (string, bool) {
	if s.HasPrefix(str, prefix) {
		return str[len(prefix):], true
	}

	return str, false
}

// CutSuffix removes suffix, reporting whether it was there.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) CutSuffix(str, suffix string) (string, bool) {
	if s.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)], true
	}

	return str, false
}

// ToLower returns str with ASCII letters lowercased, allocating in the arena
// only when there is something to change.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
//   - 2026-09-21 17:20: take the range rather than a mapping function, which
//     removes an indirect call per byte from the scan
func (s *Str) ToLower(str string) string {
	return s._MapCase(str, 'A', 'Z', CASE_SHIFT)
}

// ToUpper returns str with ASCII letters uppercased, allocating in the arena
// only when there is something to change.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
//   - 2026-09-21 17:21: take the range rather than a mapping function, which
//     removes an indirect call per byte from the scan
func (s *Str) ToUpper(str string) string {
	return s._MapCase(str, 'a', 'z', -CASE_SHIFT)
}

// Title uppercases the first letter of each word.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Title(str string) string {
	if !_HasNonASCII(str) {
		return s._TitleASCII(str)
	}

	start := true

	return s.MapUTF8(func(r rune) rune {
		mapped := r
		if start && unicode.IsLetter(r) {
			mapped = unicode.ToTitle(r)
		}

		start = unicode.IsSpace(r)

		return mapped
	}, str)
}

// ToTitle returns str with every rune in title case.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ToTitle(str string) string {
	if _HasNonASCII(str) {
		return s.MapUTF8(unicode.ToTitle, str)
	}

	// For ASCII, unicode.ToTitle and unicode.ToUpper agree on every rune.
	return s._MapCase(str, 'a', 'z', -CASE_SHIFT)
}

// MapASCII rewrites each byte through the mapping, dropping bytes it maps below
// zero.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) MapASCII(mapping func(byte) int, str string) string {
	var (
		src = res.UnsafeBytes(str)
		dst = arena.MakeSlice[byte](s.arena, 0, len(src))
	)

	for _, c := range src {
		if mapped := mapping(c); mapped >= 0 {
			dst = append(dst, byte(mapped))
		}
	}

	return res.UnsafeString(dst)
}

// MapUTF8 rewrites each rune through the mapping, dropping runes it maps below
// zero.
//
// It builds through a Buffer rather than into a slice sized from the input: a
// mapped rune can encode longer than the one it replaces, and an append past a
// preallocated arena block silently moves the result to the Go heap.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) MapUTF8(mapping func(rune) rune, str string) string {
	var (
		dst  = NewBuffer(s.arena)
		temp [utf8.UTFMax]byte
	)

	for _, r := range str {
		mapped := mapping(r)
		if mapped < 0 {
			continue
		}

		n := utf8.EncodeRune(temp[:], mapped)
		dst.Append(temp[:n])
	}

	return dst.String()
}

// MapString rewrites each rune through the mapping.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) MapString(mapping func(rune) rune, str string) string {
	return s.MapUTF8(mapping, str)
}

// Repeat returns count copies of str, joined, in the arena.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Repeat(str string, count int) string {
	if count <= 0 || len(str) == 0 {
		return ""
	}

	dst := arena.MakeSlice[byte](s.arena, 0, len(str)*count)
	for range count {
		dst = append(dst, str...)
	}

	return res.UnsafeString(dst)
}

// Replace replaces the first n copies of old with fresh, or every copy when n
// is negative.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Replace(str, old, fresh string, n int) string {
	if n == 0 || old == "" {
		return str
	}

	count := s.Count(str, old)
	if count == 0 {
		return str
	}

	if n > 0 && n < count {
		count = n
	}

	var (
		size = len(str) + count*(len(fresh)-len(old))
		dst  = arena.MakeSlice[byte](s.arena, 0, max(size, len(str)))
		done = 0
	)

	for done < count {
		idx := s.Index(str, old)
		if idx < 0 {
			break
		}

		dst = append(dst, str[:idx]...)
		dst = append(dst, fresh...)

		str = str[idx+len(old):]
		done = done + 1
	}

	dst = append(dst, str...)

	return res.UnsafeString(dst)
}

// ReplaceAll replaces every copy of old with fresh.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ReplaceAll(str, old, fresh string) string {
	return s.Replace(str, old, fresh, ALL_MATCHES)
}

// Split breaks str at each copy of sep. An empty separator splits into runes.
//
// The parts share str's memory; only the slice holding them is allocated.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Split(str, sep string) []string {
	if sep == "" {
		parts := arena.MakeSlice[string](s.arena, 0, utf8.RuneCountInString(str))
		for i, r := range str {
			parts = append(parts, str[i:i+utf8.RuneLen(r)])
		}

		return parts
	}

	parts := arena.MakeSlice[string](s.arena, 0, s.Count(str, sep)+1)

	for {
		idx := s.Index(str, sep)
		if idx < 0 {
			return append(parts, str)
		}

		parts = append(parts, str[:idx])
		str = str[idx+len(sep):]
	}
}

// SplitN breaks str at each copy of sep, into at most n parts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) SplitN(str, sep string, n int) []string {
	return s._Wrap(strings.SplitN(str, sep, n))
}

// SplitAfter breaks str after each copy of sep.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) SplitAfter(str, sep string) []string {
	return s._Wrap(strings.SplitAfter(str, sep))
}

// SplitAfterN breaks str after each copy of sep, into at most n parts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) SplitAfterN(str, sep string, n int) []string {
	return s._Wrap(strings.SplitAfterN(str, sep, n))
}

// Fields breaks str at runs of whitespace.
//
// An all-ASCII input takes a path that reads bytes rather than decoding runes,
// and answers "is this whitespace" from a table rather than through a function
// value. Anything else falls back to FieldsFunc, because unicode.IsSpace counts
// two characters outside ASCII that the table cannot know about.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
//   - 2026-09-21 17:23: add the ASCII path, which was worth three quarters of
//     the call
func (s *Str) Fields(str string) []string {
	if _HasNonASCII(str) {
		return s.FieldsFunc(str, unicode.IsSpace)
	}

	return s._FieldsASCII(str)
}

// FieldsFunc breaks str at runs of runes the predicate accepts.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
//   - 2026-09-21 17:25: call the predicate once per rune rather than twice
func (s *Str) FieldsFunc(str string, sep func(rune) bool) []string {
	var (
		parts = arena.MakeSlice[string](s.arena, 0, FIELDS_GUESS)
		start = NOT_FOUND
	)

	// One pass. Counting first meant decoding every rune and calling the
	// predicate for it twice over, to save a few reallocations of a small
	// pointer array.
	for i, r := range str {
		if sep(r) {
			if start >= 0 {
				parts = arena.Append(s.arena, parts, str[start:i])
				start = NOT_FOUND
			}

			continue
		}

		if start < 0 {
			start = i
		}
	}

	if start >= 0 {
		parts = arena.Append(s.arena, parts, str[start:])
	}

	if len(parts) == 0 {
		return nil
	}

	return parts
}

// Join concatenates elems with sep between them, in the arena.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Join(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}

	size := len(sep) * (len(elems) - 1)
	for _, e := range elems {
		size = size + len(e)
	}

	dst := arena.MakeSlice[byte](s.arena, 0, size)

	for i, e := range elems {
		if i > 0 {
			dst = append(dst, sep...)
		}

		dst = append(dst, e...)
	}

	return res.UnsafeString(dst)
}

// ToValidUTF8 replaces each run of invalid bytes with the replacement.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) ToValidUTF8(str, replacement string) string {
	if utf8.ValidString(str) {
		return str
	}

	var (
		dst   = NewBuffer(s.arena)
		valid = 0
	)

	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		if r != utf8.RuneError || size != 1 {
			i = i + size

			continue
		}

		// Hand over everything since the last bad byte in one go, rather than a
		// call to the buffer for every rune that was already fine.
		dst.AppendString(str[valid:i])
		dst.AppendString(replacement)

		for i < len(str) {
			r, size = utf8.DecodeRuneInString(str[i:])
			if r != utf8.RuneError || size != 1 {
				break
			}

			i = i + 1
		}

		valid = i
	}

	dst.AppendString(str[valid:])

	return dst.String()
}

// Lines returns an iterator over the newline-terminated lines of str.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
func (s *Str) Lines(str string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(str) > 0 {
			idx := s.IndexByte(str, '\n')
			if idx < 0 {
				yield(str)

				return
			}

			if !yield(str[:idx+1]) {
				return
			}

			str = str[idx+1:]
		}
	}
}

// _HasNonASCII reports whether any byte of str is outside ASCII.
//
// Eight bytes are tested per step: a word has a non-ASCII byte exactly when one
// of its lanes has its high bit set.
//
// Revisions:
//   - 2026-09-21 17:28: initial creation
func _HasNonASCII(str string) bool {
	b := res.UnsafeBytes(str)

	i := 0
	for ; i+WORD_BYTES <= len(b); i += WORD_BYTES {
		if binary.LittleEndian.Uint64(b[i:])&HIGH_BITS != 0 {
			return true
		}
	}

	for ; i < len(b); i++ {
		if b[i] >= ASCII_LIMIT {
			return true
		}
	}

	return false
}

// _IsASCIISpace reports whether c is one of the six ASCII whitespace bytes.
//
// Revisions:
//   - 2026-09-21 17:44: initial creation
func _IsASCIISpace(c byte) bool {
	return c < SPACE_LIMIT && ASCII_SPACE&(1<<c) != 0
}

// _TitleASCII uppercases the first letter of each whitespace-separated word.
//
// Revisions:
//   - 2026-09-21 18:40: initial creation
func (s *Str) _TitleASCII(str string) string {
	dst := arena.MakeSlice[byte](s.arena, len(str), len(str))
	start := true

	for i := 0; i < len(str); i++ {
		c := str[i]
		if start && c >= 'a' && c <= 'z' {
			c = c - CASE_SHIFT
		}

		start = _IsASCIISpace(str[i])
		dst[i] = c
	}

	return res.UnsafeString(dst)
}

// _FieldsASCII splits str on ASCII whitespace.
//
// Both passes cost a shift and a mask per byte. The rune-based path they replace
// decoded UTF-8 and made an indirect call for every character.
//
// Revisions:
//   - 2026-09-21 17:30: initial creation
func (s *Str) _FieldsASCII(str string) []string {
	count, inside := 0, false

	for i := 0; i < len(str); i++ {
		if _IsASCIISpace(str[i]) {
			inside = false

			continue
		}

		if !inside {
			count = count + 1
		}

		inside = true
	}

	if count == 0 {
		return nil
	}

	parts := arena.MakeSlice[string](s.arena, 0, count)
	start := NOT_FOUND

	for i := 0; i < len(str); i++ {
		if _IsASCIISpace(str[i]) {
			if start >= 0 {
				parts = append(parts, str[start:i])
				start = NOT_FOUND
			}

			continue
		}

		if start < 0 {
			start = i
		}
	}

	if start >= 0 {
		parts = append(parts, str[start:])
	}

	return parts
}

// _MapCase shifts every byte in [lo,hi] by delta, returning str untouched when
// there is nothing to change.
//
// The range arrives as two bytes rather than as a func(byte) byte, which is what
// this took before. A function value cannot be inlined, so the scan spent an
// indirect call on every byte of its input.
//
// Revisions:
//   - 2026-09-21 17:32: initial creation
func (s *Str) _MapCase(str string, lo, hi byte, delta int) string {
	src := res.UnsafeBytes(str)

	at := _FirstInRange(src, lo, hi)
	if at < 0 {
		return str
	}

	dst := arena.MakeSlice[byte](s.arena, len(src), len(src))
	copy(dst, src[:at])

	for i := at; i < len(src); i++ {
		c := src[i]
		if c >= lo && c <= hi {
			c = byte(int(c) + delta)
		}

		dst[i] = c
	}

	return res.UnsafeString(dst)
}

// _FirstInRange returns the first index holding a byte in [lo,hi], or NOT_FOUND.
//
// Eight bytes are tested per step while the input stays ASCII. Adding 0x80-lo to
// a lane sets its high bit when the byte is at least lo, and adding 0x80-hi-1
// sets it when the byte is past hi; a lane inside the range is therefore one
// where the first is set and the second is not. A word carrying any non-ASCII
// byte leaves the loop, because those adds would carry between lanes.
//
// Revisions:
//   - 2026-09-21 17:35: initial creation
func _FirstInRange(src []byte, lo, hi byte) int {
	var (
		atLeastLo = ONE_PER_BYTE * uint64(ASCII_LIMIT-lo)
		pastHi    = ONE_PER_BYTE * uint64(ASCII_LIMIT-hi-1)
	)

	i := 0
	for ; i+WORD_BYTES <= len(src); i += WORD_BYTES {
		w := binary.LittleEndian.Uint64(src[i:])
		if w&HIGH_BITS != 0 {
			break
		}

		if (w+atLeastLo)&^(w+pastHi)&HIGH_BITS != 0 {
			break
		}
	}

	for ; i < len(src); i++ {
		if c := src[i]; c >= lo && c <= hi {
			return i
		}
	}

	return NOT_FOUND
}

// _Wrap copies parts into an arena-backed slice. The parts themselves are not
// copied: each one points into the string that was split.
//
// Revisions:
//   - 2025-12-13 23:50: initial creation
//   - 2026-09-21 17:52: take strings rather than byte slices, now that the
//     splitting goes through the strings package
func (s *Str) _Wrap(parts []string) []string {
	wrapped := arena.MakeSlice[string](s.arena, 0, len(parts))
	wrapped = append(wrapped, parts...)

	return wrapped
}
