# Arena-Go Examples

This directory contains examples demonstrating the usage of the arena-go library.

## Running the Example

To run the main example that showcases ArenaSlice with different types:

```bash
cd example
go run main.go
```

## What the Example Demonstrates

The `main.go` file shows how to use `ArenaSlice` with:

1. **Integer slices**: Creating, appending, and manipulating int slices
2. **String slices**: Working with strings, including search operations
3. **Pointer to struct slices**: Allocating structs in the arena and working with slices of pointers

Key features demonstrated:
- Zero-GC memory allocation
- Various append methods (AppendOne, Append, AppendSlice, Push)
- Search operations (Contains, IndexOf)
- Sorting with custom comparers
- Iteration using modern Go iterators
- Cloning to heap memory when needed
- Automatic arena cleanup

## Output

When run, the example will output:
```
=== ArenaSlice Examples ===

1. Integer Slice:
Length: 8, Capacity: 16
Contents: [1 2 3 4 5 6 7 8]
After appending more: [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]

2. String Slice:
String slice: [hello world arena memory allocation]
'arena' found at index: 2

3. Pointer to Struct Slice:
Number of people: 3
Person 1: Alice is 30 years old
Person 2: Bob is 25 years old
Person 3: Charlie is 35 years old

Sorting people by age:
Person 1: Bob is 25 years old
Person 2: Alice is 30 years old
Person 3: Charlie is 35 years old

4. Clone to heap:
Cloned integers: [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15]

=== Example completed successfully! ===
```