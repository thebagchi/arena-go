package container_test

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
)

// BOUNDARY_BYTES holds every byte where the word-at-a-time scans change their
// answer, plus the six ASCII whitespace characters and some non-ASCII. Random
// strings drawn from this alphabet exercise the lane arithmetic far harder than
// prose would, because almost every byte sits on an edge.
var (
	BOUNDARY_BYTES = []byte{
		0, 1,
		'@', 'A', 'B', 'Y', 'Z', '[',
		'`', 'a', 'b', 'y', 'z', '{',
		'0', '9',
		' ', '\t', '\n', '\v', '\f', '\r',
		0x7e, 0x7f, 0x80, 0xc3, 0xa9, 0xff,
	}
)

// TestStr_ASCIIMatchesStdlib compares the ASCII fast paths against the standard
// library over random inputs.
//
// ToLower, ToUpper and Fields read eight bytes at a time and decide from lane
// arithmetic rather than from a comparison per byte. That is easy to get subtly
// wrong at a boundary or at a length that is not a multiple of the word, and a
// unit test with a friendly input would not notice. Lengths from zero to just
// past a word are covered explicitly for the same reason.
//
// Revisions:
//   - 2026-09-21 18:05: initial creation
func TestStr_ASCIIMatchesStdlib(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 << 20))
	defer a.Delete()

	str := container.NewStr(a)

	for length := range 40 {
		for range 200 {
			checkAll(t, str, randomString(length))
		}
	}
}

// TestStr_ASCIIEdgeLengths covers the lengths either side of a word boundary
// with inputs that are entirely one byte value, so a tail that is handled
// wrongly shows up as a changed length rather than a changed character.
//
// Revisions:
//   - 2026-09-21 18:08: initial creation
func TestStr_ASCIIEdgeLengths(t *testing.T) {
	a := arena.New(alloc.NewBumpAllocator(1 << 20))
	defer a.Delete()

	str := container.NewStr(a)

	for _, c := range []byte{'A', 'Z', 'a', 'z', ' ', '@', '[', 0x80} {
		for _, length := range []int{0, 1, 7, 8, 9, 15, 16, 17, 23, 24, 25} {
			checkAll(t, str, strings.Repeat(string([]byte{c}), length))
		}
	}
}

// randomString draws length bytes from BOUNDARY_BYTES.
//
// Revisions:
//   - 2026-09-21 18:10: initial creation
func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = BOUNDARY_BYTES[rand.IntN(len(BOUNDARY_BYTES))]
	}

	return string(b)
}

// checkAll compares every ASCII fast path against the standard library.
//
// Revisions:
//   - 2026-09-21 18:20: initial creation
func checkAll(t *testing.T, str *container.Str, in string) {
	t.Helper()

	got, want := str.ToLower(in), strings.ToLower(in)
	if !asciiEqual(got, want, in) {
		t.Fatalf("ToLower(%q) = %q, want %q", in, got, want)
	}

	got, want = str.ToUpper(in), strings.ToUpper(in)
	if !asciiEqual(got, want, in) {
		t.Fatalf("ToUpper(%q) = %q, want %q", in, got, want)
	}

	got, want = str.ToTitle(in), strings.ToTitle(in)
	if !asciiEqual(got, want, in) {
		t.Fatalf("ToTitle(%q) = %q, want %q", in, got, want)
	}

	// Title has no standard-library counterpart that is not deprecated, so the
	// reference is the rune-by-rune definition the ASCII path replaced.
	if got, want = str.Title(in), referenceTitle(in); got != want {
		t.Fatalf("Title(%q) = %q, want %q", in, got, want)
	}

	if got, want = str.ToValidUTF8(in, "?"), strings.ToValidUTF8(in, "?"); got != want {
		t.Fatalf("ToValidUTF8(%q) = %q, want %q", in, got, want)
	}

	checkFields(t, str, in)
}

// referenceTitle is the obvious rune-by-rune implementation, kept as the thing
// the ASCII path has to agree with.
//
// Revisions:
//   - 2026-09-21 18:52: initial creation
func referenceTitle(in string) string {
	var (
		out   strings.Builder
		start = true
	)

	for _, r := range in {
		mapped := r
		if start && unicode.IsLetter(r) {
			mapped = unicode.ToTitle(r)
		}

		start = unicode.IsSpace(r)

		out.WriteRune(mapped)
	}

	return out.String()
}

// checkFields compares Fields against the standard library, which it must match
// exactly: the ASCII path is taken only when the input has no byte that
// unicode.IsSpace could disagree about.
//
// Revisions:
//   - 2026-09-21 18:12: initial creation
func checkFields(t *testing.T, str *container.Str, in string) {
	t.Helper()

	var (
		got  = str.Fields(in)
		want = strings.Fields(in)
	)

	if len(got) != len(want) {
		t.Fatalf("Fields(%q) = %q, want %q", in, got, want)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("Fields(%q)[%d] = %q, want %q", in, i, got[i], want[i])
		}
	}
}

// asciiEqual reports whether two case conversions agree, allowing the documented
// difference: this package maps ASCII only, where the standard library maps the
// whole of Unicode.
//
// Revisions:
//   - 2026-09-21 18:14: initial creation
func asciiEqual(got, want, in string) bool {
	if got == want {
		return true
	}

	for i := 0; i < len(in); i++ {
		if in[i] >= 0x80 {
			return true
		}
	}

	return false
}
