// Command example demonstrates the arena API: vectors, maps, strings and the
// io adapters, all allocated outside Go's garbage collector.
package main

import (
	"encoding/json"
	"fmt"

	arena "github.com/thebagchi/arena-go"
	"github.com/thebagchi/arena-go/alloc"
	"github.com/thebagchi/arena-go/container"
	arenaio "github.com/thebagchi/arena-go/io"
)

const (
	// ARENA_SIZE is the first chunk the demonstration arena maps.
	ARENA_SIZE = 64 << 10
	// PEOPLE is how many entries the map demonstration inserts.
	PEOPLE = 20
	// PREVIEW is how many entries the map demonstration prints.
	PREVIEW = 3
	// AGE_STEP spaces the generated ages apart.
	AGE_STEP = 10
	// ALICE_AGE, BOB_AGE and JOHN_AGE are the ages the sorting and encoding
	// demonstrations use.
	ALICE_AGE = 30
	BOB_AGE   = 25
	JOHN_AGE  = 28
	// ALICE_NEW_AGE is the value the map demonstration overwrites with.
	ALICE_NEW_AGE = 31
)

var (
	// FIRST_NUMBERS and MORE_NUMBERS are the two batches the vector
	// demonstration appends, kept as data rather than inline literals.
	FIRST_NUMBERS = []int{1, 2, 3, 4, 5}
	MORE_NUMBERS  = []int{6, 7, 8}
)

// Person is the struct the examples store in the arena.
type Person struct {
	Name string
	Age  int
}

// main runs each demonstration against one arena and releases it at the end.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func main() {
	a := arena.New(alloc.NewBumpAllocator(ARENA_SIZE))
	defer a.Delete()

	showVec(a)
	showMap(a)
	showStrings(a)
	showIO(a)
}

// showVec appends to an arena-backed vector and sorts it.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func showVec(a *arena.Arena) {
	fmt.Println("== vector")

	numbers := container.NewVec[int](a)
	numbers.AppendSlice(FIRST_NUMBERS)
	numbers.Append(MORE_NUMBERS...)

	fmt.Printf(
		"  contents: %v (len %d, cap %d)\n",
		numbers.Slice(),
		numbers.Len(),
		numbers.Cap(),
	)

	words := container.NewVec[string](a)
	words.Append(a.MakeString("arena"), a.MakeString("memory"))

	fmt.Printf("  index of arena: %d\n", container.IndexOf(words, "arena"))

	// Pointers into the arena are safe to store in arena memory; a pointer to a
	// Go heap value would not be, because nothing here is a garbage collector
	// root.
	people := container.NewVec[*Person](a)
	people.Append(
		arena.Ptr(a, Person{Name: a.MakeString("Alice"), Age: ALICE_AGE}),
		arena.Ptr(a, Person{Name: a.MakeString("Bob"), Age: BOB_AGE}),
	)

	people.Sort(func(l, r *Person) bool {
		return l.Age < r.Age
	})

	for i, person := range people.All2() {
		fmt.Printf("  person %d: %s is %d\n", i+1, person.Name, person.Age)
	}
}

// showMap fills an arena-backed map and reads it back.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func showMap(a *arena.Arena) {
	fmt.Println("== map")

	ages := container.NewMap[string, int](a)
	for i := range PEOPLE {
		// Map copies string keys into the arena, so a key built on the Go heap
		// is safe to hand over.
		ages.Set(fmt.Sprintf("person%d", i), i*AGE_STEP)
	}

	ages.Set("alice", ALICE_NEW_AGE)
	ages.Delete("person0")

	if age, found := ages.Get("alice"); found {
		fmt.Printf("  alice: %d\n", age)
	}

	fmt.Printf("  entries: %d\n", ages.Len())

	shown := 0

	for key, age := range ages.All() {
		if shown >= PREVIEW {
			break
		}

		fmt.Printf("  %s: %d\n", key, age)

		shown = shown + 1
	}
}

// showStrings runs a few arena-allocated string operations.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func showStrings(a *arena.Arena) {
	fmt.Println("== strings")

	str := container.NewStr(a)
	text := a.MakeString("  the quick brown fox  ")

	fmt.Printf("  trimmed: %q\n", str.TrimSpace(text))
	fmt.Printf("  title:   %q\n", str.Title(str.TrimSpace(text)))
	fmt.Printf("  fields:  %v\n", str.Fields(text))
	fmt.Printf("  joined:  %q\n", str.Join(str.Fields(text), "-"))
}

// showIO serialises a struct through the arena-backed writer and reads bytes
// back through the reader.
//
// Revisions:
//   - 2025-12-10 20:59: initial creation
func showIO(a *arena.Arena) {
	fmt.Println("== io")

	writer := arenaio.NewWriter(a)
	person := arena.Ptr(a, Person{Name: a.MakeString("John Doe"), Age: JOHN_AGE})

	if err := json.NewEncoder(writer).Encode(person); err != nil {
		fmt.Printf("  encode failed: %v\n", err)

		return
	}

	fmt.Printf("  json: %s", writer.String())
	fmt.Printf("  arena-backed: %t\n", arena.OwnsSlice(a, writer.Bytes()))

	reader := arenaio.NewReader(a, writer.Bytes())
	buf := make([]byte, PEOPLE)

	n, err := reader.Read(buf)
	if err != nil {
		fmt.Printf("  read failed: %v\n", err)

		return
	}

	fmt.Printf("  read %d bytes, %d left\n", n, reader.Len())
}
