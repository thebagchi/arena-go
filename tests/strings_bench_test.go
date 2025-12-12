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
		_ = arena.ToBytes(benchStr)
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
		_ = arena.ToString(data)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.TrimSpace(benchStr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Contains(benchStr, benchSubstr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.HasPrefix(benchStr, "  hello")
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.HasSuffix(benchStr, "  ")
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Index(benchStr, benchSubstr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.LastIndex(benchStr, benchSubstr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Trim(benchStr, benchCutset)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.TrimLeft(benchStr, benchCutset)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.TrimRight(benchStr, benchCutset)
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
	str1 := "Hello World"
	str2 := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.EqualFold(str1, str2)
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
	str1 := "apple"
	str2 := "banana"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Compare(str1, str2)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.ToLower(benchStr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.ToUpper(benchStr)
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
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Split(a, benchSplitStr, ",")
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
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Join(a, parts, ",")
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
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Fields(a, benchFieldsStr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Count(benchLongStr, benchSubstr)
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
	a := arena.New(8192, arena.BUMP)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Replace(a, benchLongStr, benchOld, benchNew, -1)
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
	a := arena.New(8192, arena.BUMP)
	defer a.Delete()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.ReplaceAll(a, benchLongStr, benchOld, benchNew)
		a.Reset()
	}
}

// Repeat Benchmarks
func BenchmarkStdRepeat(b *testing.B) {
	str := "hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = strings.Repeat(str, 10)
	}
}

func BenchmarkArenaRepeat(b *testing.B) {
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	str := "hello"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Repeat(a, str, 10)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.TrimPrefix(benchStr, "  hello")
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.TrimSuffix(benchStr, "  ")
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = arena.Cut(benchStr, benchSubstr)
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.IndexByte(benchStr, 'w')
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.ContainsAny(benchStr, "xyz")
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
	a := arena.New(65536, arena.BUMP)
	defer a.Delete()
	longStr := strings.Repeat("word,", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arena.Split(a, longStr, ",")
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
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = arena.Split(a, benchSplitStr, ",")
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
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = arena.Fields(a, benchFieldsStr)
		a.Reset()
	}
}
