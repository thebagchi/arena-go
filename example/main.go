package main

import (
	"fmt"

	"github.com/thebagchi/arena-go"
)

// Example struct for pointer demonstration
type Person struct {
	Name string
	Age  int
}

func main() {
	// Create an arena with 4KB memory
	a := arena.New(4096, arena.BUMP)
	defer a.Delete()

	fmt.Println("=== ArenaSlice Examples ===")

	// 1. Integer slice
	fmt.Println("\n1. Integer Slice:")
	intSlice := arena.MakeArenaSlice[int](a)
	intSlice.AppendSlice([]int{1, 2, 3, 4, 5})
	intSlice.Append(6, 7, 8) // Append multiple at once

	fmt.Printf("Length: %d, Capacity: %d\n", intSlice.Len(), intSlice.Cap())
	fmt.Printf("Contents: %v\n", intSlice.Slice())

	// Demonstrate zero-GC append
	for i := 9; i <= 15; i++ {
		intSlice.AppendOne(i)
	}
	fmt.Printf("After appending more: %v\n", intSlice.Slice())

	// 2. String slice
	fmt.Println("\n2. String Slice:")
	stringSlice := arena.MakeArenaSlice[string](a, "hello", "world")
	stringSlice.AppendSlice([]string{"arena", "memory"})
	stringSlice.Push("allocation") // Using Push alias

	fmt.Printf("String slice: %v\n", stringSlice.Slice())

	// Demonstrate Contains and IndexOf
	if stringSlice.Contains("arena") {
		fmt.Printf("'arena' found at index: %d\n", stringSlice.IndexOf("arena"))
	}

	// 3. Pointer to struct slice
	fmt.Println("\n3. Pointer to Struct Slice:")

	// Create some Person structs in the arena
	person1 := arena.Ptr(a, Person{Name: "Alice", Age: 30})
	person2 := arena.Ptr(a, Person{Name: "Bob", Age: 25})
	person3 := arena.Ptr(a, Person{Name: "Charlie", Age: 35})

	// Create a slice of pointers to Person
	pointerSlice := arena.MakeArenaSlice[*Person](a)
	pointerSlice.Append(person1, person2, person3)

	fmt.Printf("Number of people: %d\n", pointerSlice.Len())

	// Iterate and print
	for i, personPtr := range pointerSlice.All2() {
		fmt.Printf("Person %d: %s is %d years old\n", i+1, personPtr.Name, personPtr.Age)
	}

	// Demonstrate sorting (by age)
	fmt.Println("\nSorting people by age:")
	pointerSlice.Sort(func(a, b *Person) bool {
		return a.Age < b.Age
	})

	for i, personPtr := range pointerSlice.All2() {
		fmt.Printf("Person %d: %s is %d years old\n", i+1, personPtr.Name, personPtr.Age)
	}

	// 4. Demonstrate Clone (heap escape)
	fmt.Println("\n4. Clone to heap:")
	clonedInts := intSlice.Clone()
	fmt.Printf("Cloned integers: %v\n", clonedInts)

	// Arena is automatically cleaned up when main exits (defer a.Delete())
	fmt.Println("\n=== Example completed successfully! ===")
}
