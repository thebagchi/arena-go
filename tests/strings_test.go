package arena_test

import (
	"testing"

	arena "github.com/thebagchi/arena-go"
)

func TestToBytes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{"empty", "", []byte{}},
		{"simple", "hello", []byte("hello")},
		{"unicode", "世界", []byte("世界")},
		{"whitespace", "  test  ", []byte("  test  ")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := arena.UnsafeBytes(tt.s)
			if len(got) != len(tt.want) {
				t.Errorf("ToBytes() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ToBytes()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{"simple", []byte("hello"), "hello"},
		{"unicode", []byte("世界"), "世界"},
		{"whitespace", []byte("  test  "), "  test  "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := arena.UnsafeString(tt.b)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimSpace(t *testing.T) {
	a := arena.New(1, arena.BUMP); str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"no space", "hello", "hello"},
		{"leading", "  hello", "hello"},
		{"trailing", "hello  ", "hello"},
		{"both", "  hello  ", "hello"},
		{"tabs", "\t\thello\t\t", "hello"},
		{"newlines", "\n\nhello\n\n", "hello"},
		{"mixed", " \t\n hello \t\n ", "hello"},
		{"empty", "", ""},
		{"only spaces", "   ", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.TrimSpace(tt.s)
			if got != tt.want {
				t.Errorf("TrimSpace(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", true},
		{"spaces", "   ", true},
		{"tabs", "\t\t", true},
		{"newlines", "\n\n", true},
		{"mixed whitespace", " \t\n ", true},
		{"with content", "hello", false},
		{"content with spaces", "  hello  ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.IsEmpty(tt.s)
			if got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"contains", "hello world", "world", true},
		{"not contains", "hello world", "golang", false},
		{"empty substr", "hello", "", true},
		{"empty string", "", "test", false},
		{"both empty", "", "", true},
		{"case sensitive", "Hello", "hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{"has prefix", "hello world", "hello", true},
		{"no prefix", "hello world", "world", false},
		{"empty prefix", "hello", "", true},
		{"longer prefix", "hi", "hello", false},
		{"case sensitive", "Hello", "hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.HasPrefix(tt.s, tt.prefix)
			if got != tt.want {
				t.Errorf("HasPrefix(%q, %q) = %v, want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		suffix string
		want   bool
	}{
		{"has suffix", "hello world", "world", true},
		{"no suffix", "hello world", "hello", false},
		{"empty suffix", "hello", "", true},
		{"longer suffix", "hi", "hello", false},
		{"case sensitive", "World", "world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.HasSuffix(tt.s, tt.suffix)
			if got != tt.want {
				t.Errorf("HasSuffix(%q, %q) = %v, want %v", tt.s, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{"found at start", "hello world", "hello", 0},
		{"found in middle", "hello world", "o w", 4},
		{"found at end", "hello world", "world", 6},
		{"not found", "hello world", "golang", -1},
		{"empty substr", "hello", "", 0},
		{"multiple occurrences", "hello hello", "hello", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Index(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Index(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestLastIndex(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{"found at end", "hello world", "world", 6},
		{"multiple occurrences", "hello hello", "hello", 6},
		{"single occurrence", "hello world", "o w", 4},
		{"not found", "hello world", "golang", -1},
		{"empty substr", "hello", "", 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.LastIndex(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("LastIndex(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"trim spaces", "  hello  ", " ", "hello"},
		{"trim multiple chars", "xxxhelloxxx", "x", "hello"},
		{"trim mixed", "abchelloabc", "abc", "hello"},
		{"no trim needed", "hello", "x", "hello"},
		{"trim all", "xxxxx", "x", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Trim(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("Trim(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestTrimLeft(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"trim left spaces", "  hello  ", " ", "hello  "},
		{"trim left chars", "xxxhelloxxx", "x", "helloxxx"},
		{"no trim needed", "hello", "x", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.TrimLeft(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("TrimLeft(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestTrimRight(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"trim right spaces", "  hello  ", " ", "  hello"},
		{"trim right chars", "xxxhelloxxx", "x", "xxxhello"},
		{"no trim needed", "hello", "x", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.TrimRight(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("TrimRight(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestEqualFold(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		t    string
		want bool
	}{
		{"equal", "hello", "hello", true},
		{"case insensitive", "Hello", "hello", true},
		{"mixed case", "HeLLo", "hEllO", true},
		{"not equal", "hello", "world", false},
		{"different length", "hello", "hell", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.EqualFold(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("EqualFold(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		t    string
		want int
	}{
		{"equal", "hello", "hello", 0},
		{"less than", "apple", "banana", -1},
		{"greater than", "zebra", "apple", 1},
		{"empty strings", "", "", 0},
		{"empty vs non-empty", "", "hello", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Compare(tt.s, tt.t)
			// Normalize to -1, 0, 1
			if got < 0 {
				got = -1
			} else if got > 0 {
				got = 1
			}
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"all upper", "HELLO", "hello"},
		{"mixed", "HeLLo", "hello"},
		{"already lower", "hello", "hello"},
		{"with numbers", "Hello123", "hello123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.ToLower(tt.s)
			if got != tt.want {
				t.Errorf("ToLower(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"all lower", "hello", "HELLO"},
		{"mixed", "HeLLo", "HELLO"},
		{"already upper", "HELLO", "HELLO"},
		{"with numbers", "Hello123", "HELLO123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.ToUpper(tt.s)
			if got != tt.want {
				t.Errorf("ToUpper(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"single word", "hello", "Hello"},
		{"multiple words", "hello world", "Hello World"},
		{"already title", "Hello World", "Hello World"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Title(tt.s)
			if got != tt.want {
				t.Errorf("Title(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()

	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{"simple", "a,b,c", ",", []string{"a", "b", "c"}},
		{"empty sep", "hello", "", []string{"h", "e", "l", "l", "o"}},
		{"no match", "hello", ",", []string{"hello"}},
		{"trailing sep", "a,b,", ",", []string{"a", "b", ""}},
		{"leading sep", ",a,b", ",", []string{"", "a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Split(tt.s, tt.sep)
			if len(got) != len(tt.want) {
				t.Errorf("Split(%q, %q) length = %v, want %v", tt.s, tt.sep, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Split(%q, %q)[%d] = %q, want %q", tt.s, tt.sep, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestJoin(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()

	tests := []struct {
		name  string
		elems []string
		sep   string
		want  string
	}{
		{"simple", []string{"a", "b", "c"}, ",", "a,b,c"},
		{"empty sep", []string{"a", "b", "c"}, "", "abc"},
		{"single elem", []string{"hello"}, ",", "hello"},
		{"empty slice", []string{}, ",", ""},
		{"with spaces", []string{"hello", "world"}, " ", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Join(tt.elems, tt.sep)
			if got != tt.want {
				t.Errorf("Join(%v, %q) = %q, want %q", tt.elems, tt.sep, got, tt.want)
			}
		})
	}
}

func TestFields(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()

	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"simple", "hello world", []string{"hello", "world"}},
		{"multiple spaces", "hello  world", []string{"hello", "world"}},
		{"tabs", "hello\tworld", []string{"hello", "world"}},
		{"mixed whitespace", "hello \t\n world", []string{"hello", "world"}},
		{"leading trailing", "  hello world  ", []string{"hello", "world"}},
		{"single word", "hello", []string{"hello"}},
		{"empty", "", []string{}},
		{"only whitespace", "   ", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Fields(tt.s)
			if len(got) != len(tt.want) {
				t.Errorf("Fields(%q) length = %v, want %v", tt.s, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Fields(%q)[%d] = %q, want %q", tt.s, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitJoinRoundtrip(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()

	tests := []struct {
		name string
		s    string
		sep  string
	}{
		{"simple", "a,b,c", ","},
		{"spaces", "hello world test", " "},
		{"pipes", "one|two|three", "|"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := str.Split(tt.s, tt.sep)
			got := str.Join(parts, tt.sep)
			if got != tt.s {
				t.Errorf("Split then Join roundtrip failed: got %q, want %q", got, tt.s)
			}
		})
	}
}

func TestLines(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{"empty", "", []string{}},
		{"single line no newline", "hello", []string{"hello"}},
		{"single line with newline", "hello\n", []string{"hello\n"}},
		{"multiple lines", "line1\nline2\nline3", []string{"line1\n", "line2\n", "line3"}},
		{"multiple lines no final newline", "line1\nline2\nline3", []string{"line1\n", "line2\n", "line3"}},
		{"empty lines", "\n\n", []string{"\n", "\n"}},
		{"mixed", "line1\n\nline3\n", []string{"line1\n", "\n", "line3\n"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for line := range str.Lines(tt.s) {
				got = append(got, line)
			}
			if len(got) != len(tt.want) {
				t.Errorf("Lines() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Lines()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestClone(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Reset()

	tests := []struct {
		name string
		s    string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"unicode", "世界"},
		{"long", "this is a longer string with multiple words"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.Clone(tt.s)
			if got != tt.s {
				t.Errorf("Clone() = %q, want %q", got, tt.s)
			}
		})
	}
}

func TestFieldsFunc(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Reset()

	tests := []struct {
		name string
		s    string
		f    func(rune) bool
		want []string
	}{
		{"spaces", "hello world test", func(r rune) bool { return r == ' ' }, []string{"hello", "world", "test"}},
		{"commas", "a,b,c", func(r rune) bool { return r == ',' }, []string{"a", "b", "c"}},
		{"empty", "", func(r rune) bool { return r == ' ' }, []string{}},
		{"no fields", "   ", func(r rune) bool { return r == ' ' }, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.FieldsFunc(tt.s, tt.f)
			if len(got) != len(tt.want) {
				t.Errorf("FieldsFunc() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("FieldsFunc()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestContainsFunc(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		f    func(rune) bool
		want bool
	}{
		{"contains space", "hello world", func(r rune) bool { return r == ' ' }, true},
		{"no space", "helloworld", func(r rune) bool { return r == ' ' }, false},
		{"contains digit", "abc123", func(r rune) bool { return r >= '0' && r <= '9' }, true},
		{"no digit", "abc", func(r rune) bool { return r >= '0' && r <= '9' }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.ContainsFunc(tt.s, tt.f)
			if got != tt.want {
				t.Errorf("ContainsFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexFunc(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		f    func(rune) bool
		want int
	}{
		{"first space", "hello world", func(r rune) bool { return r == ' ' }, 5},
		{"no space", "helloworld", func(r rune) bool { return r == ' ' }, -1},
		{"first digit", "abc123", func(r rune) bool { return r >= '0' && r <= '9' }, 3},
		{"no digit", "abc", func(r rune) bool { return r >= '0' && r <= '9' }, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.IndexFunc(tt.s, tt.f)
			if got != tt.want {
				t.Errorf("IndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndexFunc(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		f    func(rune) bool
		want int
	}{
		{"last space", "hello world test", func(r rune) bool { return r == ' ' }, 11},
		{"no space", "helloworld", func(r rune) bool { return r == ' ' }, -1},
		{"last digit", "abc123def456", func(r rune) bool { return r >= '0' && r <= '9' }, 11},
		{"no digit", "abc", func(r rune) bool { return r >= '0' && r <= '9' }, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.LastIndexFunc(tt.s, tt.f)
			if got != tt.want {
				t.Errorf("LastIndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapASCII(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name    string
		s       string
		mapping func(byte) int
		want    string
	}{
		{"to upper", "hello", func(c byte) int { return int(c - 32) }, "HELLO"},
		{"drop spaces", "hello world", func(c byte) int {
			if c == ' ' {
				return -1
			}
			return int(c)
		}, "helloworld"},
		{"identity", "hello", func(c byte) int { return int(c) }, "hello"},
		{"digits only", "a1b2c3", func(c byte) int {
			if c >= '0' && c <= '9' {
				return int(c)
			}
			return -1
		}, "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.MapASCII(tt.mapping, tt.s)
			if got != tt.want {
				t.Errorf("MapASCII() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapUTF8(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name    string
		s       string
		mapping func(rune) rune
		want    string
	}{
		{"to upper", "hello", func(r rune) rune { return r - 32 }, "HELLO"},
		{"drop spaces", "hello world", func(r rune) rune {
			if r == ' ' {
				return -1
			}
			return r
		}, "helloworld"},
		{"identity", "hello", func(r rune) rune { return r }, "hello"},
		{"unicode", "café", func(r rune) rune {
			if r == 'é' {
				return 'e'
			}
			return r
		}, "cafe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.MapUTF8(tt.mapping, tt.s)
			if got != tt.want {
				t.Errorf("MapUTF8() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapString(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name    string
		s       string
		mapping func(rune) rune
		want    string
	}{
		{"to upper", "hello", func(r rune) rune { return r - 32 }, "HELLO"},
		{"drop spaces", "hello world", func(r rune) rune {
			if r == ' ' {
				return -1
			}
			return r
		}, "helloworld"},
		{"identity", "hello", func(r rune) rune { return r }, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.MapString(tt.mapping, tt.s)
			if got != tt.want {
				t.Errorf("MapString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Reset()

	tests := []struct {
		name string
		s    string
		want string
	}{
		{"simple", "hello world", "HELLO WORLD"},
		{"mixed", "hello WORLD", "HELLO WORLD"},
		{"unicode", "héllo wörld", "HÉLLO WÖRLD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.ToTitle(tt.s)
			if got != tt.want {
				t.Errorf("ToTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToValidUTF8(t *testing.T) {
	a := arena.New(4096, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Reset()

	tests := []struct {
		name        string
		s           string
		replacement string
		want        string
	}{
		{"valid", "hello", "?", "hello"},
		{"invalid bytes", "hello\xffworld", "?", "hello?world"},
		{"empty replacement", "hello\xffworld", "", "helloworld"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.ToValidUTF8(tt.s, tt.replacement)
			if got != tt.want {
				t.Errorf("ToValidUTF8() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimFunc(t *testing.T) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	tests := []struct {
		name string
		s    string
		f    func(rune) bool
		want string
	}{
		{"trim spaces", "  hello  ", func(r rune) bool { return r == ' ' }, "hello"},
		{"trim digits", "123hello456", func(r rune) bool { return r >= '0' && r <= '9' }, "hello"},
		{"no trim", "hello", func(r rune) bool { return r == ' ' }, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := str.TrimFunc(tt.s, tt.f)
			if got != tt.want {
				t.Errorf("TrimFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}
