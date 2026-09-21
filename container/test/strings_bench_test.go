package container_test

import (
	"strings"
	"testing"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
	"github.com/thebagchi/arena-go/res"
)

var (
	benchStr       = "  hello world from arena allocator  "
	benchLongStr   = strings.Repeat("hello world ", 100)
	benchSubstr    = "world"
	benchCutset    = " "
	benchOld       = "world"
	benchNew       = "arena"
	benchSplitStr  = "a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z"
	benchFieldsStr = "hello world from arena allocator with many fields here"
)

// ToBytes/ToString Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdStringToBytes(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = []byte(benchStr)
	}
}

// BenchmarkString_ZeroCopyToBytes measures string zero copy to bytes.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyToBytes(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = res.UnsafeBytes(benchStr)
	}
}

// BenchmarkString_StdBytesToString measures string std bytes to string.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdBytesToString(b *testing.B) {
	data := []byte(benchStr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}

// BenchmarkString_ZeroCopyToString measures string zero copy to string.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyToString(b *testing.B) {
	data := []byte(benchStr)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = res.UnsafeString(data)
	}
}

// TrimSpace Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimSpace(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimSpace(benchStr)
	}
}

// BenchmarkString_ZeroCopyTrimSpace measures string zero copy trim space.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrimSpace(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimSpace(benchStr)
	}
}

// Contains Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdContains(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Contains(benchStr, benchSubstr)
	}
}

// BenchmarkString_ZeroCopyContains measures string zero copy contains.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyContains(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Contains(benchStr, benchSubstr)
	}
}

// HasPrefix Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdHasPrefix(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.HasPrefix(benchStr, "  hello")
	}
}

// BenchmarkString_ZeroCopyHasPrefix measures string zero copy has prefix.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyHasPrefix(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.HasPrefix(benchStr, "  hello")
	}
}

// HasSuffix Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdHasSuffix(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.HasSuffix(benchStr, "  ")
	}
}

// BenchmarkString_ZeroCopyHasSuffix measures string zero copy has suffix.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyHasSuffix(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.HasSuffix(benchStr, "  ")
	}
}

// Index Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdIndex(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Index(benchStr, benchSubstr)
	}
}

// BenchmarkString_ZeroCopyIndex measures string zero copy index.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyIndex(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Index(benchStr, benchSubstr)
	}
}

// LastIndex Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdLastIndex(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.LastIndex(benchStr, benchSubstr)
	}
}

// BenchmarkString_ZeroCopyLastIndex measures string zero copy last index.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyLastIndex(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.LastIndex(benchStr, benchSubstr)
	}
}

// Trim Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrim(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Trim(benchStr, benchCutset)
	}
}

// BenchmarkString_ZeroCopyTrim measures string zero copy trim.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrim(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Trim(benchStr, benchCutset)
	}
}

// TrimLeft Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimLeft(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimLeft(benchStr, benchCutset)
	}
}

// BenchmarkString_ZeroCopyTrimLeft measures string zero copy trim left.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrimLeft(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimLeft(benchStr, benchCutset)
	}
}

// TrimRight Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimRight(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimRight(benchStr, benchCutset)
	}
}

// BenchmarkString_ZeroCopyTrimRight measures string zero copy trim right.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrimRight(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimRight(benchStr, benchCutset)
	}
}

// EqualFold Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdEqualFold(b *testing.B) {
	str1 := "Hello World"
	str2 := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.EqualFold(str1, str2)
	}
}

// BenchmarkString_ZeroCopyEqualFold measures string zero copy equal fold.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyEqualFold(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	str1 := "Hello World"
	str2 := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.EqualFold(str1, str2)
	}
}

// Compare Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdCompare(b *testing.B) {
	str1 := "apple"
	str2 := "banana"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Compare(str1, str2)
	}
}

// BenchmarkString_ZeroCopyCompare measures string zero copy compare.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyCompare(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	str1 := "apple"
	str2 := "banana"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Compare(str1, str2)
	}
}

// ToLower Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdToLower(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ToLower(benchStr)
	}
}

// BenchmarkString_ZeroCopyToLower measures string zero copy to lower.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyToLower(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ToLower(benchStr)
	}
}

// ToUpper Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdToUpper(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ToUpper(benchStr)
	}
}

// BenchmarkString_ZeroCopyToUpper measures string zero copy to upper.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyToUpper(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ToUpper(benchStr)
	}
}

// Title Benchmarks
//
// There is no standard-library counterpart to benchmark against: strings.Title
// was deprecated in Go 1.18 for mishandling Unicode word boundaries, and the
// replacement lives in golang.org/x/text, which this module does not depend on.

// BenchmarkString_ZeroCopyTitle measures string zero copy title.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTitle(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Title(benchStr)
	}
}

// Split Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdSplit(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Split(benchSplitStr, ",")
	}
}

// BenchmarkString_ArenaSplit measures string arena split.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaSplit(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Split(benchSplitStr, ",")

		a.Reset()
	}
}

// Join Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdJoin(b *testing.B) {
	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Join(parts, ",")
	}
}

// BenchmarkString_ArenaJoin measures string arena join.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaJoin(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	parts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Join(parts, ",")

		a.Reset()
	}
}

// Fields Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdFields(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Fields(benchFieldsStr)
	}
}

// BenchmarkString_ArenaFields measures string arena fields.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaFields(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Fields(benchFieldsStr)

		a.Reset()
	}
}

// Count Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdCount(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Count(benchLongStr, benchSubstr)
	}
}

// BenchmarkString_ZeroCopyCount measures string zero copy count.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyCount(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Count(benchLongStr, benchSubstr)
	}
}

// Replace Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdReplace(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ReplaceAll(benchLongStr, benchOld, benchNew)
	}
}

// BenchmarkString_ArenaReplace measures string arena replace.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaReplace(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Replace(benchLongStr, benchOld, benchNew, -1)

		a.Reset()
	}
}

// ReplaceAll Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdReplaceAll(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ReplaceAll(benchLongStr, benchOld, benchNew)
	}
}

// BenchmarkString_ArenaReplaceAll measures string arena replace all.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaReplaceAll(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ReplaceAll(benchLongStr, benchOld, benchNew)

		a.Reset()
	}
}

// Repeat Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdRepeat(b *testing.B) {
	testStr := "hello"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Repeat(testStr, 10)
	}
}

// BenchmarkString_ArenaRepeat measures string arena repeat.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaRepeat(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	testStr := "hello"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Repeat(testStr, 10)

		a.Reset()
	}
}

// TrimPrefix Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimPrefix(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimPrefix(benchStr, "  hello")
	}
}

// BenchmarkString_ZeroCopyTrimPrefix measures string zero copy trim prefix.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrimPrefix(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimPrefix(benchStr, "  hello")
	}
}

// TrimSuffix Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimSuffix(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimSuffix(benchStr, "  ")
	}
}

// BenchmarkString_ZeroCopyTrimSuffix measures string zero copy trim suffix.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyTrimSuffix(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimSuffix(benchStr, "  ")
	}
}

// Cut Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdCut(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = strings.Cut(benchStr, benchSubstr)
	}
}

// BenchmarkString_ZeroCopyCut measures string zero copy cut.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyCut(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = str.Cut(benchStr, benchSubstr)
	}
}

// IndexByte Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdIndexByte(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.IndexByte(benchStr, 'w')
	}
}

// BenchmarkString_ZeroCopyIndexByte measures string zero copy index byte.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyIndexByte(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.IndexByte(benchStr, 'w')
	}
}

// ContainsAny Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdContainsAny(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ContainsAny(benchStr, "xyz")
	}
}

// BenchmarkString_ZeroCopyContainsAny measures string zero copy contains any.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ZeroCopyContainsAny(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ContainsAny(benchStr, "xyz")
	}
}

// Long String Benchmarks (to test scalability)
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdSplitLong(b *testing.B) {
	longStr := strings.Repeat("word,", 1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Split(longStr, ",")
	}
}

// BenchmarkString_ArenaSplitLong measures string arena split long.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaSplitLong(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	longStr := strings.Repeat("word,", 1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Split(longStr, ",")

		a.Reset()
	}
}

// Memory allocation comparison benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdSplitAllocs(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = strings.Split(benchSplitStr, ",")
	}
}

// BenchmarkString_ArenaSplitAllocs measures string arena split allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaSplitAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.Split(benchSplitStr, ",")

		a.Reset()
	}
}

// BenchmarkString_StdFieldsAllocs measures string std fields allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdFieldsAllocs(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = strings.Fields(benchFieldsStr)
	}
}

// BenchmarkString_ArenaFieldsAllocs measures string arena fields allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaFieldsAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.Fields(benchFieldsStr)

		a.Reset()
	}
}

// Lines Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_Lines(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	linesStr := strings.Repeat("This is a line of text\n", 100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		for range str.Lines(linesStr) {
			count++
		}

		_ = count
	}
}

// BenchmarkString_StdLines measures string std lines.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdLines(b *testing.B) {
	linesStr := strings.Repeat("This is a line of text\n", 100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.Split(strings.TrimSuffix(linesStr, "\n"), "\n")
	}
}

// Clone Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaClone(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.Clone(benchStr)

		a.Reset()
	}
}

// FieldsFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdFieldsFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.FieldsFunc(benchFieldsStr, isSpace)
	}
}

// BenchmarkString_ArenaFieldsFunc measures string arena fields func.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaFieldsFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.FieldsFunc(benchFieldsStr, isSpace)

		a.Reset()
	}
}

// ContainsFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ContainsFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ContainsFunc(benchStr, isDigit)
	}
}

// IndexFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_IndexFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.IndexFunc(benchStr, isSpace)
	}
}

// LastIndexFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_LastIndexFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.LastIndexFunc(benchStr, isSpace)
	}
}

// MapString Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdMap(b *testing.B) {
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

// BenchmarkString_ArenaMapASCII measures string arena map ascii.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaMapASCII(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
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

// BenchmarkString_ArenaMapUTF8 measures string arena map utf 8.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaMapUTF8(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
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

// BenchmarkString_ArenaMapString measures string arena map string.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaMapString(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
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
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdToTitle(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ToTitle(benchStr)
	}
}

// BenchmarkString_ArenaToTitle measures string arena to title.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaToTitle(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ToTitle(benchStr)

		a.Reset()
	}
}

// ToValidUTF8 Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdToValidUTF8(b *testing.B) {
	invalidStr := "hello\xffworld\xfe"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.ToValidUTF8(invalidStr, "?")
	}
}

// BenchmarkString_ArenaToValidUTF8 measures string arena to valid utf 8.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaToValidUTF8(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	invalidStr := "hello\xffworld\xfe"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.ToValidUTF8(invalidStr, "?")

		a.Reset()
	}
}

// TrimFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimFunc(benchStr, isSpace)
	}
}

// BenchmarkString_TrimFunc measures string trim func.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_TrimFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimFunc(benchStr, isSpace)
	}
}

// TrimLeftFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimLeftFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimLeftFunc(benchStr, isSpace)
	}
}

// BenchmarkString_TrimLeftFunc measures string trim left func.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_TrimLeftFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimLeftFunc(benchStr, isSpace)
	}
}

// TrimRightFunc Benchmarks
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_StdTrimRightFunc(b *testing.B) {
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = strings.TrimRightFunc(benchStr, isSpace)
	}
}

// BenchmarkString_TrimRightFunc measures string trim right func.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_TrimRightFunc(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
	isSpace := func(r rune) bool { return r == ' ' || r == '\t' }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = str.TrimRightFunc(benchStr, isSpace)
	}
}

// Allocation comparison benchmarks for new functions
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaCloneAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.Clone(benchStr)

		a.Reset()
	}
}

// BenchmarkString_ArenaFieldsFuncAllocs measures string arena fields func allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaFieldsFuncAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	isSpace := func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.FieldsFunc(benchFieldsStr, isSpace)

		a.Reset()
	}
}

// BenchmarkString_ArenaMapStringAllocs measures string arena map string allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaMapStringAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))
	str := container.NewStr(a)
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

// BenchmarkString_ArenaToTitleAllocs measures string arena to title allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaToTitleAllocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.ToTitle(benchStr)

		a.Reset()
	}
}

// BenchmarkString_ArenaToValidUTF8Allocs measures string arena to valid utf 8allocs.
//
// Revisions:
//   - 2025-12-12 20:07: initial creation
func BenchmarkString_ArenaToValidUTF8Allocs(b *testing.B) {
	a := arena.New(alloc.NewBumpAllocator(1 * 4096))

	str := container.NewStr(a)
	defer a.Delete()

	invalidStr := "hello\xffworld\xfe"

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = str.ToValidUTF8(invalidStr, "?")

		a.Reset()
	}
}
