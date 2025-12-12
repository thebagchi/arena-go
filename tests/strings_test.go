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
			got := arena.ToBytes(tt.s)
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
			got := arena.ToString(tt.b)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimSpace(t *testing.T) {
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
			got := arena.TrimSpace(tt.s)
			if got != tt.want {
				t.Errorf("TrimSpace(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
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
			got := arena.IsEmpty(tt.s)
			if got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
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
			got := arena.Contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
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
			got := arena.HasPrefix(tt.s, tt.prefix)
			if got != tt.want {
				t.Errorf("HasPrefix(%q, %q) = %v, want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
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
			got := arena.HasSuffix(tt.s, tt.suffix)
			if got != tt.want {
				t.Errorf("HasSuffix(%q, %q) = %v, want %v", tt.s, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
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
			got := arena.Index(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Index(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestLastIndex(t *testing.T) {
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
			got := arena.LastIndex(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("LastIndex(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
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
			got := arena.Trim(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("Trim(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestTrimLeft(t *testing.T) {
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
			got := arena.TrimLeft(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("TrimLeft(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestTrimRight(t *testing.T) {
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
			got := arena.TrimRight(tt.s, tt.cutset)
			if got != tt.want {
				t.Errorf("TrimRight(%q, %q) = %q, want %q", tt.s, tt.cutset, got, tt.want)
			}
		})
	}
}

func TestEqualFold(t *testing.T) {
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
			got := arena.EqualFold(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("EqualFold(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
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
			got := arena.Compare(tt.s, tt.t)
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
			got := arena.ToLower(tt.s)
			if got != tt.want {
				t.Errorf("ToLower(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
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
			got := arena.ToUpper(tt.s)
			if got != tt.want {
				t.Errorf("ToUpper(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestTitle(t *testing.T) {
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
			got := arena.Title(tt.s)
			if got != tt.want {
				t.Errorf("Title(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
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
			got := arena.Split(a, tt.s, tt.sep)
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
			got := arena.Join(a, tt.elems, tt.sep)
			if got != tt.want {
				t.Errorf("Join(%v, %q) = %q, want %q", tt.elems, tt.sep, got, tt.want)
			}
		})
	}
}

func TestFields(t *testing.T) {
	a := arena.New(1024, arena.BUMP)
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
			got := arena.Fields(a, tt.s)
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
			parts := arena.Split(a, tt.s, tt.sep)
			got := arena.Join(a, parts, tt.sep)
			if got != tt.s {
				t.Errorf("Split then Join roundtrip failed: got %q, want %q", got, tt.s)
			}
		})
	}
}
