package arena_test

import (
	"strings"
	"testing"

	arena "github.com/thebagchi/arena-go"
)

var (
	benchStr       = "  hello world from arena allocator  "
	benchLongStr   = strings.Repeat("hello world ", 100)
	benchSubstr    = "world"
	benchSep       = " "
	benchCutset    = " "
	benchOld       = "world"
	benchNew       = "arena"
	benchSplitStr  = "a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z"
	benchFieldsStr = "hello world from arena allocator with many fields here"
)

// ToBytes/ToString Benchmarks
func BenchmarkStdStringToBytes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(benchStr)
	}
}

func BenchmarkZeroCopyToBytes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.UnsafeBytes(benchStr)
	}
}

func BenchmarkStdBytesToString(b *testing.B) {
	data := []byte(benchStr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

func BenchmarkZeroCopyToString(b *testing.B) {
	data := []byte(benchStr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.UnsafeString(data)
	}
}

// TrimSpace Benchmarks
func BenchmarkStdTrimSpace(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimSpace(benchStr)
	}
}

func BenchmarkZeroCopyTrimSpace(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimSpace(benchStr)
	}
}

// Contains Benchmarks
func BenchmarkStdContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Contains(benchStr, benchSubstr)
	}
}

func BenchmarkZeroCopyContains(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Contains(benchStr, benchSubstr)
	}
}

// HasPrefix Benchmarks
func BenchmarkStdHasPrefix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.HasPrefix(benchStr, "  hello")
	}
}

func BenchmarkZeroCopyHasPrefix(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.HasPrefix(benchStr, "  hello")
	}
}

// HasSuffix Benchmarks
func BenchmarkStdHasSuffix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.HasSuffix(benchStr, "  ")
	}
}

func BenchmarkZeroCopyHasSuffix(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.HasSuffix(benchStr, "  ")
	}
}

// Index Benchmarks
func BenchmarkStdIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Index(benchStr, benchSubstr)
	}
}

func BenchmarkZeroCopyIndex(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Index(benchStr, benchSubstr)
	}
}

// LastIndex Benchmarks
func BenchmarkStdLastIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.LastIndex(benchStr, benchSubstr)
	}
}

func BenchmarkZeroCopyLastIndex(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.LastIndex(benchStr, benchSubstr)
	}
}

// Trim Benchmarks
func BenchmarkStdTrim(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Trim(benchStr, benchCutset)
	}
}

func BenchmarkZeroCopyTrim(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Trim(benchStr, benchCutset)
	}
}

// TrimLeft Benchmarks
func BenchmarkStdTrimLeft(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimLeft(benchStr, benchCutset)
	}
}

func BenchmarkZeroCopyTrimLeft(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimLeft(benchStr, benchCutset)
	}
}

// TrimRight Benchmarks
func BenchmarkStdTrimRight(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimRight(benchStr, benchCutset)
	}
}

func BenchmarkZeroCopyTrimRight(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimRight(benchStr, benchCutset)
	}
}

// EqualFold Benchmarks
func BenchmarkStdEqualFold(b *testing.B) {
	str1 := "Hello World"
	str2 := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.EqualFold(str1, str2)
	}
}

func BenchmarkZeroCopyEqualFold(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	str1 := "Hello World"
	str2 := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.EqualFold(str1, str2)
	}
}

// Compare Benchmarks
func BenchmarkStdCompare(b *testing.B) {
	str1 := "apple"
	str2 := "banana"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Compare(str1, str2)
	}
}

func BenchmarkZeroCopyCompare(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	str1 := "apple"
	str2 := "banana"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Compare(str1, str2)
	}
}

// ToLower Benchmarks
func BenchmarkStdToLower(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ToLower(benchStr)
	}
}

func BenchmarkZeroCopyToLower(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ToLower(benchStr)
	}
}

// ToUpper Benchmarks
func BenchmarkStdToUpper(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ToUpper(benchStr)
	}
}

func BenchmarkZeroCopyToUpper(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ToUpper(benchStr)
	}
}

// Title Benchmarks
func BenchmarkStdTitle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Title(benchStr)
	}
}

func BenchmarkZeroCopyTitle(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Title(benchStr)
	}
}

// Split Benchmarks
func BenchmarkStdSplit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Split(benchSplitStr, ",")
	}
}

func BenchmarkArenaSplit(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Split(benchSplitStr, ",")
		a.Reset()
	}
}

// Join Benchmarks
func BenchmarkStdJoin(b *testing.B) {
	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Join(parts, ",")
	}
}

func BenchmarkArenaJoin(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Join(parts, ",")
		a.Reset()
	}
}

// Fields Benchmarks
func BenchmarkStdFields(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Fields(benchFieldsStr)
	}
}

func BenchmarkArenaFields(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Fields(benchFieldsStr)
		a.Reset()
	}
}

// Count Benchmarks
func BenchmarkStdCount(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Count(benchLongStr, benchSubstr)
	}
}

func BenchmarkZeroCopyCount(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Count(benchLongStr, benchSubstr)
	}
}

// Replace Benchmarks
func BenchmarkStdReplace(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Replace(benchLongStr, benchOld, benchNew, -1)
	}
}

func BenchmarkArenaReplace(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Replace(benchLongStr, benchOld, benchNew, -1)
		a.Reset()
	}
}

// ReplaceAll Benchmarks
func BenchmarkStdReplaceAll(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ReplaceAll(benchLongStr, benchOld, benchNew)
	}
}

func BenchmarkArenaReplaceAll(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ReplaceAll(benchLongStr, benchOld, benchNew)
		a.Reset()
	}
}

// Repeat Benchmarks
func BenchmarkStdRepeat(b *testing.B) {
	testStr := "hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Repeat(testStr, 10)
	}
}

func BenchmarkArenaRepeat(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	testStr := "hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Repeat(testStr, 10)
		a.Reset()
	}
}

// TrimPrefix Benchmarks
func BenchmarkStdTrimPrefix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimPrefix(benchStr, "  hello")
	}
}

func BenchmarkZeroCopyTrimPrefix(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimPrefix(benchStr, "  hello")
	}
}

// TrimSuffix Benchmarks
func BenchmarkStdTrimSuffix(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimSuffix(benchStr, "  ")
	}
}

func BenchmarkZeroCopyTrimSuffix(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimSuffix(benchStr, "  ")
	}
}

// Cut Benchmarks
func BenchmarkStdCut(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = strings.Cut(benchStr, benchSubstr)
	}
}

func BenchmarkZeroCopyCut(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = str.Cut(benchStr, benchSubstr)
	}
}

// IndexByte Benchmarks
func BenchmarkStdIndexByte(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.IndexByte(benchStr, 'w')
	}
}

func BenchmarkZeroCopyIndexByte(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.IndexByte(benchStr, 'w')
	}
}

// ContainsAny Benchmarks
func BenchmarkStdContainsAny(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ContainsAny(benchStr, "xyz")
	}
}

func BenchmarkZeroCopyContainsAny(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ContainsAny(benchStr, "xyz")
	}
}

// Long String Benchmarks (to test scalability)
func BenchmarkStdSplitLong(b *testing.B) {
	longStr := strings.Repeat("word,", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Split(longStr, ",")
	}
}

func BenchmarkArenaSplitLong(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	longStr := strings.Repeat("word,", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Split(longStr, ",")
		a.Reset()
	}
}

// Memory allocation comparison benchmarks
func BenchmarkStdSplitAllocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = strings.Split(benchSplitStr, ",")
	}
}

func BenchmarkArenaSplitAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.Split(benchSplitStr, ",")
		a.Reset()
	}
}

func BenchmarkStdFieldsAllocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = strings.Fields(benchFieldsStr)
	}
}

func BenchmarkArenaFieldsAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.Fields(benchFieldsStr)
		a.Reset()
	}
}

// Lines Benchmarks
func BenchmarkLines(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	linesStr := strings.Repeat("This is a line of text\n", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for _ = range str.Lines(linesStr) {
			count++
		}
		_ = count
	}
}

func BenchmarkStdLines(b *testing.B) {
	linesStr := strings.Repeat("This is a line of text\n", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Split(strings.TrimSuffix(linesStr, "\n"), "\n")
	}
}

// Clone Benchmarks
func BenchmarkArenaClone(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.Clone(benchStr)
		a.Reset()
	}
}

// FieldsFunc Benchmarks
func BenchmarkStdFieldsFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.FieldsFunc(benchFieldsStr, isSpace)
	}
}

func BenchmarkArenaFieldsFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.FieldsFunc(benchFieldsStr, isSpace)
		a.Reset()
	}
}

// ContainsFunc Benchmarks
func BenchmarkContainsFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ContainsFunc(benchStr, isDigit)
	}
}

// IndexFunc Benchmarks
func BenchmarkIndexFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.IndexFunc(benchStr, isSpace)
	}
}

// LastIndexFunc Benchmarks
func BenchmarkLastIndexFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.LastIndexFunc(benchStr, isSpace)
	}
}

// MapString Benchmarks
func BenchmarkStdMap(b *testing.B) {
	toUpper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Map(toUpper, benchStr)
	}
}

func BenchmarkArenaMapASCII(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	toUpper := func(c byte) int {
		if c >= 'a' && c <= 'z' {
			return int(c - 32)
		}
		return int(c)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.MapASCII(toUpper, benchStr)
	}
}

func BenchmarkArenaMapUTF8(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	toUpper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.MapUTF8(toUpper, benchStr)
	}
}

func BenchmarkArenaMapString(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	toUpper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.MapString(toUpper, benchStr)
	}
}

// ToTitle Benchmarks
func BenchmarkStdToTitle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ToTitle(benchStr)
	}
}

func BenchmarkArenaToTitle(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ToTitle(benchStr)
		a.Reset()
	}
}

// ToValidUTF8 Benchmarks
func BenchmarkStdToValidUTF8(b *testing.B) {
	invalidStr := "hello\xffworld\xfe"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.ToValidUTF8(invalidStr, "?")
	}
}

func BenchmarkArenaToValidUTF8(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	invalidStr := "hello\xffworld\xfe"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.ToValidUTF8(invalidStr, "?")
		a.Reset()
	}
}

// TrimFunc Benchmarks
func BenchmarkStdTrimFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimFunc(benchStr, isSpace)
	}
}

func BenchmarkTrimFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimFunc(benchStr, isSpace)
	}
}

// TrimLeftFunc Benchmarks
func BenchmarkStdTrimLeftFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimLeftFunc(benchStr, isSpace)
	}
}

func BenchmarkTrimLeftFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimLeftFunc(benchStr, isSpace)
	}
}

// TrimRightFunc Benchmarks
func BenchmarkStdTrimRightFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.TrimRightFunc(benchStr, isSpace)
	}
}

func BenchmarkTrimRightFunc(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = str.TrimRightFunc(benchStr, isSpace)
	}
}

// Allocation comparison benchmarks for new functions
func BenchmarkArenaCloneAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.Clone(benchStr)
		a.Reset()
	}
}

func BenchmarkArenaFieldsFuncAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.FieldsFunc(benchFieldsStr, isSpace)
		a.Reset()
	}
}

func BenchmarkArenaMapStringAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	toUpper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.MapString(toUpper, benchStr)
	}
}

func BenchmarkArenaToTitleAllocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.ToTitle(benchStr)
		a.Reset()
	}
}

func BenchmarkArenaToValidUTF8Allocs(b *testing.B) {
	a := arena.New(1, arena.BUMP)
	str := arena.NewStr(a)
	defer a.Delete()
	invalidStr := "hello\xffworld\xfe"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = str.ToValidUTF8(invalidStr, "?")
		a.Reset()
	}
}
