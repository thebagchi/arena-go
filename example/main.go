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
	intSlice := arena.NewVec[int](a)
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
	stringSlice := arena.NewVec[string](a, "hello", "world")
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
	pointerSlice := arena.NewVec[*Person](a)
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

	fmt.Println("\n=== ArenaMap Examples ===")

	// 5. String to int map
	fmt.Println("\n5. String to Int Map:")
	stringMap := arena.NewMap[string, int](a)
	stringMap.Set("alice", 30)
	stringMap.Set("bob", 25)
	stringMap.Set("charlie", 35)

	// Add more entries to trigger growth
	for i := 0; i < 20; i++ {
		stringMap.Set(fmt.Sprintf("person%d", i), i*10)
	}

	fmt.Printf("Map length: %d\n", stringMap.Len())

	// Get values
	if age, found := stringMap.Get("alice"); found {
		fmt.Printf("Alice's age: %d\n", age)
	}

	// Range over map
	fmt.Println("All entries:")
	stringMap.Range(func(key string, value int) bool {
		fmt.Printf("  %s: %d\n", key, value)
		return true
	})

	// Update and delete
	stringMap.Set("alice", 31) // Update
	stringMap.Delete("bob")    // Delete

	fmt.Printf("After update/delete, length: %d\n", stringMap.Len())

	// Demonstrate Clone (heap escape)
	fmt.Println("\n6. Clone Map to heap:")
	clonedMap := stringMap.Clone()
	fmt.Printf("Cloned map: %v\n", clonedMap)

	// Demonstrate iterators
	fmt.Println("\n7. Iterator Examples:")

	// Keys iterator
	fmt.Print("Keys: ")
	count := 0
	for key := range stringMap.Keys() {
		if count < 5 {
			fmt.Printf("%s ", key)
		}
		count++
	}
	fmt.Printf("... (total: %d)\n", count)

	// All iterator (key-value pairs)
	fmt.Println("First 3 entries using All():")
	count = 0
	for key, val := range stringMap.All() {
		if count < 3 {
			fmt.Printf("  %s: %d\n", key, val)
		}
		count++
		if count >= 3 {
			break
		}
	}

	// Pull-based iterator
	fmt.Println("First 3 entries using Iter():")
	iter := stringMap.Iter()
	for i := 0; i < 3; i++ {
		key, val, ok := iter.Next()
		if !ok {
			break
		}
		fmt.Printf("  %s: %d\n", key, val)
	}

	// Arena is automatically cleaned up when main exits (defer a.Delete())
	fmt.Println("\n=== Example completed successfully! ===")
}
