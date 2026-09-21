# arena-go

A high-performance, zero-GC memory allocator library for Go with generic containers. Allocate memory outside the garbage collector using custom allocators (Bump, Slab, Buddy) and work with type-safe generic containers.

## Features

- **Zero-GC Allocations**: Allocate memory outside Go's garbage collector using arena-based memory management
- **Multiple Allocators**: Choose from Bump (fastest), Slab (fixed-size), or Buddy (flexible) allocation strategies
- **Generic Containers**: Type-safe generic containers including Vec (dynamic arrays), Map, SkipList, Pool, and Str
- **Thread-Safe Allocators**: All three allocators are safe for concurrent use, as are `Map`, `SkipList` and `Pool`. See [Thread Safety](#thread-safety) for the containers that are not
- **Zeroed Memory**: Every allocation comes back zeroed, without the allocation path paying for it
- **Deterministic Memory Management**: Explicit control over memory lifecycle with Reset() and Delete()

## The one rule

Arena memory comes from `mmap` and **is invisible to Go's garbage collector**. A Go
pointer whose only reference lives in arena memory does not keep its target alive: it
will be collected and its storage reused while the arena still points at it. The race
detector cannot see this, and it only shows up under memory pressure.

What that rules out: storing a Go heap string, slice, map, channel, func, interface or
pointer in arena memory. What it permits: values with no pointers at all, and pointers
that lead back into the same arena.

Bring outside data in by copying it:

```go
name := a.MakeString(userInput)   // copied into the arena, safe to store
node := arena.Ptr(a, Node{})      // lives in the arena, safe to point at
```

```go
// Wrong: the only reference to each string is arena memory.
for i := range n {
    vec.AppendOne(fmt.Sprintf("item%d", i))
}

// Right: the arena owns the bytes.
for i := range n {
    vec.AppendOne(a.MakeString(fmt.Sprintf("item%d", i)))
}
```

`Map` copies string keys and string values for you. Anything else is yours to copy.

## Installation

```bash
go get github.com/thebagchi/arena-go
```

## Quick Start

### Creating an Arena

```go
package main

import (
    "github.com/thebagchi/arena-go"
    "github.com/thebagchi/arena-go/alloc"
)

func main() {
    // Create an arena with 4KB of memory using Bump allocator
    a := arena.New(alloc.NewBumpAllocator(4096))
    defer a.Delete()
    
    // Use the arena...
}
```

## Allocators

### Bump Allocator
Fastest allocator, best for batch allocations or when the arena is reset frequently.

```go
a := arena.New(alloc.NewBumpAllocator(1024 * 4096)) // 4MB
```

**Example use case:**
```go
// Perfect for processing a batch of requests
a := arena.New(alloc.NewBumpAllocator(1024 * 1024)) // 1MB per request
defer a.Delete()

// Process request...
users := container.NewVec[User](a)
for _, id := range requestUserIDs {
    user := arena.Alloc[User](a)
    user.ID = id
    users.Append(user)
}

// Reset for next batch
a.Reset()
```

### Slab Allocator
Best for fixed-size objects with high allocation/free turnover. Features a multi-tiered design with 17 size classes (16B to 1MB), exhausted/available slab tracking, and in-object free lists for efficient memory management.

```go
a := arena.New(alloc.NewSlabAllocator())
```

**Architecture highlights:**
- 17 size classes for objects from 16 bytes to 1MB
- Per-size-class bins with exhausted and available slab lists
- In-object free lists for O(1) allocation/deallocation
- Cross-bin free page pool for efficient memory recycling

**Example use case:**
```go
// Perfect for a long-lived server handling many connections
a := arena.New(alloc.NewSlabAllocator())
defer a.Delete()

pool := container.NewPool[Connection](a)

// Allocate on connect
conn := pool.Alloc()
conn.ID = nextID()

// Free on disconnect (reusable for next connection)
pool.Free(conn)
```

### Buddy Allocator
Most flexible, good for varied-size allocations with power-of-2 sizes.

```go
a := arena.New(alloc.NewBuddyAllocator(1024 * 4096))
```

**Example use case:**
```go
// Perfect for variable-size allocations
a := arena.New(alloc.NewBuddyAllocator(2 * 1024 * 1024)) // 2MB
defer a.Delete()

// Allocate varied sizes
smallVec := container.NewVec[byte](a)  // small data
largeVec := container.NewVec[int64](a) // larger data
map1 := container.NewMap[string, int](a)

// Flexible management of different sizes
smallVec.Append(1, 2, 3)
largeVec.Append(100, 200, 300)
```

### Allocator Comparison

Measured on 2026-09-21, Go 1.26.2, linux/amd64, five runs. `benchmarks.md` has the full
set and how to reproduce it.

| | Bump | Slab | Buddy |
| --- | --- | --- | --- |
| `Alloc` | 12.4 ns/op | 23.0 ns/op | 54.1 ns/op |
| Go heap allocations | 0 | 0 | 0 |
| Frees individually | no | yes | yes |
| Free strategy | reset the whole arena | per object | per block |
| Allocation pattern | linear | fixed-size classes | power-of-2 blocks |
| Fragmentation | none | minimal | low, blocks merge back |
| Grows | yes, chunks double | yes | yes, unless fixed |

**Choose Bump if** allocations are short-lived, arrive in batches, or the arena is
reset often. It is the fastest and frees nothing individually.

**Choose Slab if** the program allocates and frees objects of similar sizes over a
long life, such as a connection or node pool.

**Choose Buddy if** sizes vary and blocks must be freed one at a time. It costs more
per allocation because it splits and merges blocks, which is what keeps
fragmentation low.

## Core Operations

### Allocating Objects

```go
import "github.com/thebagchi/arena-go"

// Allocate a single integer
ptr := arena.Alloc[int](a)
*ptr = 42

// Allocate and initialize a struct
type Person struct {
    Name string
    Age  int
}

person := arena.Ptr(a, Person{Name: "Alice", Age: 30})

// Create a string in the arena
str := a.MakeString("hello world")
```

### Allocating Slices

```go
// Allocate a slice with length 10, capacity 20
slice := arena.MakeSlice[int](a, 10, 20)
slice[0] = 100
```

## Generic Containers

### Vec - Dynamic Array

```go
import "github.com/thebagchi/arena-go/container"

// Create a dynamic array
vec := container.NewVec[int](a)

// Append elements
vec.Append(1, 2, 3)
vec.AppendOne(4)
vec.AppendSlice([]int{5, 6, 7})

// Access elements
fmt.Println(vec.Len())      // 7
fmt.Println(vec.Cap())      // capacity
fmt.Println(vec.Slice())    // []int{1, 2, 3, 4, 5, 6, 7}

// Pop elements
if val, ok := vec.Pop(); ok {
    fmt.Println(val) // 7
}

// Search. These are functions rather than methods, because finding a value
// needs the element type to be comparable and a Vec holds any type.
idx := container.IndexOf(vec, 3)
fmt.Println(idx) // 2
fmt.Println(container.Contains(vec, 3)) // true

// Clear
vec.Clear()
fmt.Println(vec.Len()) // 0
```

### Map - Type-Safe Hash Map

```go
// Create a map
m := container.NewMap[string, int](a)

// Set values
m.Set("alice", 30)
m.Set("bob", 25)

// Get values
if age, found := m.Get("alice"); found {
    fmt.Println(age) // 30
}

// Check existence
if m.Contains("charlie") {
    fmt.Println("Found charlie")
}

// Iterate
m.Range(func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)
    return true
})

// Delete
m.Delete("bob")

// Get length
fmt.Println(m.Len()) // 1
```

### SkipList - Ordered Map

```go
// Create a skip list (ordered by key)
sl := container.NewSkipList[int, string](a)

// Insert elements
sl.Insert(10, "ten")
sl.Insert(5, "five")
sl.Insert(15, "fifteen")

// Search
if val, found := sl.Search(5); found {
    fmt.Println(val) // "five"
}

// Iterate in order
sl.Range(func(key int, value string) bool {
    fmt.Printf("%d: %s\n", key, value)
    return true
})

// Delete
sl.Delete(5)
```

### Pool - Object Pool

```go
// Create an object pool for efficient allocation/deallocation
pool := container.NewPool[Person](a)

// Allocate from pool
p := pool.Alloc()
p.Name = "Alice"
p.Age = 30

// Return to pool for reuse
pool.Free(p)

// Next allocation will reuse the same memory (zeroed)
p2 := pool.Alloc()
fmt.Println(p2.Name) // "" (zeroed)
fmt.Println(p2.Age)  // 0 (zeroed)

// Check pool stats
fmt.Println(pool.Len()) // number in free list
fmt.Println(pool.Cap()) // total capacity
```

### Str - String Utilities

```go
// Create a string builder
str := container.NewStr(a)

// Build strings
str.WriteString("Hello, ")
str.WriteString("World!")

// Get result
result := str.String()
fmt.Println(result)

// Or as bytes
bytes := str.Bytes()
```

## I/O Operations

### Writer

```go
import arenaio "github.com/thebagchi/arena-go/io"

// Create a writer backed by arena memory
writer := arenaio.NewWriter(a)

// Write data
writer.Write([]byte("hello"))
writer.WriteString(" world")
writer.WriteByte('!')

// Get result
fmt.Println(string(writer.Bytes())) // "hello world!"

// Check stats
fmt.Println(writer.Len()) // bytes written
fmt.Println(writer.Cap()) // capacity

// Reset for reuse
writer.Reset()
```

### Reader

```go
// Create a reader
data := []byte("hello world")
reader := arenaio.NewReader(a, data)

// Read data
buf := make([]byte, 5)
n, err := reader.Read(buf)
fmt.Println(string(buf[:n])) // "hello"

// Check position
fmt.Println(reader.Len())  // bytes remaining
fmt.Println(reader.Size()) // total size

// Reset to beginning
reader.Reset()
```

## Memory Management

### Resetting an Arena

Clears allocations but retains underlying memory pages:

```go
a.Reset()
```

### Deleting an Arena

Frees all memory back to the OS:

```go
a.Delete()
```

### Checking Memory Ownership

```go
// Check if a pointer was allocated by this arena
if arena.OwnsPtr(a, ptr) {
    fmt.Println("Pointer belongs to this arena")
}

// Check if a slice was allocated by this arena
if arena.OwnsSlice(a, slice) {
    fmt.Println("Slice belongs to this arena")
}
```

## Performance Tips

1. **Choose the right allocator**: Use Bump for batch allocations, Slab for fixed-size objects, Buddy for mixed sizes
2. **Use pools for frequent allocations**: Object pools reduce allocation overhead for frequently created/destroyed objects
3. **Pre-allocate containers**: Specify capacity when creating containers if you know the size
4. **Reset vs Delete**: Use Reset() for batch operations, Delete() only when completely done
5. **Arena per operation**: Create separate arenas for independent operations and delete them when done

## Thread Safety

**Safe for concurrent use:** all three allocators, and `Map`, `SkipList` and `Pool`.

**Not safe for concurrent use:** `Vec`, `Queue`, `Stack`, `Buffer` and `Str`. These hand
back interior references — `Slice()`, `Bytes()`, and the strings `Str` returns — which
stay valid after any lock inside them would have been released, so a lock would promise
a safety it could not deliver. Give each goroutine its own arena, or guard the container
yourself.

```go
// Safe: the allocator serialises, and Map takes its own lock.
a := arena.New(alloc.NewSlabAllocator())
defer a.Delete()

m := container.NewMap[string, int](a)

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        m.Set(fmt.Sprintf("key%d", n), n)
    }(i)
}
wg.Wait()
```

```go
// Safe: one arena per goroutine, which is the usual shape and the fastest.
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()

        a := arena.New(alloc.NewBumpAllocator(4096))
        defer a.Delete()

        vec := container.NewVec[int](a)
        vec.AppendOne(n)
    }(i)
}
wg.Wait()
```

## Examples

### Example 1: Processing User Requests with Bump Allocator

```go
package main

import (
    "fmt"
    "github.com/thebagchi/arena-go"
    "github.com/thebagchi/arena-go/alloc"
    "github.com/thebagchi/arena-go/container"
)

type User struct {
    ID   int
    Name string
}

// Each request gets its own arena, reset after processing
func processRequest(userIDs []int) {
    a := arena.New(alloc.NewBumpAllocator(1024 * 100)) // 100KB per request
    defer a.Delete()
    
    users := container.NewVec[*User](a)
    
    // Allocate users from the arena
    for _, id := range userIDs {
        user := arena.Alloc[User](a)
        user.ID = id
        user.Name = fmt.Sprintf("User_%d", id)
        users.Append(user)
    }
    
    fmt.Printf("Processed %d users\n", users.Len())
    // All memory freed when a.Delete() is called
}

func main() {
    processRequest([]int{1, 2, 3, 4, 5})
}
```

### Example 2: Connection Pool with Slab Allocator

```go
type Connection struct {
    ID        int
    Buffer    *container.Str
    Timestamp int64
}

func runServer() {
    // Single arena for the entire server lifetime
    a := arena.New(alloc.NewSlabAllocator())
    defer a.Delete()
    
    connPool := container.NewPool[Connection](a)
    activeConnections := container.NewMap[int, *Connection](a)
    
    // Handle new connection
    onConnect := func(connID int) {
        // Get connection from pool (reuses memory if available)
        conn := connPool.Alloc()
        conn.ID = connID
        conn.Buffer = container.NewStr(a)
        conn.Timestamp = 0
        
        activeConnections.Set(connID, conn)
    }
    
    // Handle disconnection
    onDisconnect := func(connID int) {
        if conn, found := activeConnections.Get(connID); found {
            activeConnections.Delete(connID)
            // Memory is reused for next connection
            connPool.Free(conn)
        }
    }
    
    // Simulate connections
    onConnect(1)
    onConnect(2)
    onDisconnect(1)
    onConnect(3) // Reuses memory from connection 1
    
    fmt.Printf("Active connections: %d\n", activeConnections.Len())
}
```

### Example 3: Data Analysis with Buddy Allocator

```go
type DataPoint struct {
    Value     float64
    Timestamp int64
}

func analyzeData() {
    a := arena.New(alloc.NewBuddyAllocator(2 * 1024 * 1024)) // 2MB
    defer a.Delete()
    
    // Store different data types in the same arena
    measurements := container.NewVec[float64](a)
    metadata := container.NewMap[string, string](a)
    timeSeries := container.NewVec[DataPoint](a)
    
    // Populate data
    for i := 0; i < 1000; i++ {
        measurements.Append(float64(i) * 3.14)
        point := arena.Alloc[DataPoint](a)
        point.Value = float64(i)
        point.Timestamp = int64(i * 1000)
        timeSeries.Append(point)
    }
    
    metadata.Set("source", "sensor_1")
    metadata.Set("location", "lab_a")
    
    fmt.Printf("Measurements: %d, Time Series: %d\n", 
        measurements.Len(), timeSeries.Len())
}
```

### Example 4: String Building with Multiple Types

```go
func buildJSON(users []map[string]interface{}) string {
    a := arena.New(alloc.NewBumpAllocator(1024 * 10))
    defer a.Delete()
    
    builder := container.NewStr(a)
    builder.WriteString("[")
    
    for i, user := range users {
        if i > 0 {
            builder.WriteString(",")
        }
        builder.WriteString(`{"name":"`)
        builder.WriteString(user["name"].(string))
        builder.WriteString(`","age":`)
        builder.WriteString(fmt.Sprintf("%v", user["age"]))
        builder.WriteString("}")
    }
    
    builder.WriteString("]")
    return string(builder.Bytes())
}
```

### Example 5: Using Clone for Heap-Independent Data

```go
func processWithSnapshot() {
    a := arena.New(alloc.NewBumpAllocator(1024 * 4))
    defer a.Delete()
    
    // Create data in arena
    arenaSlice := arena.MakeSlice[int](a, 5, 10)
    arenaSlice[0] = 1
    arenaSlice[1] = 2
    arenaSlice[2] = 3
    
    // Clone to heap for passing outside arena scope
    heapSlice := arena.CloneSlice(arenaSlice[:3])
    
    // Now safe to use heapSlice after arena is deleted
    return heapSlice // Safe - not owned by arena
}
```

### Example 6: Multiple Arenas for Parallel Operations

```go
func processParallel(batches [][]int) []int {
    results := container.NewVec[int](nil) // Use standard allocation
    
    var wg sync.WaitGroup
    for _, batch := range batches {
        wg.Add(1)
        go func(b []int) {
            defer wg.Done()
            
            // Each goroutine gets its own arena
            a := arena.New(alloc.NewBumpAllocator(1024))
            defer a.Delete()
            
            vec := container.NewVec[int](a)
            for _, v := range b {
                vec.Append(v * 2)
            }
            
            // Sum and store result
            sum := 0
            vec.Range(func(v int) bool {
                sum += v
                return true
            })
            
            results.Append(sum)
        }(batch)
    }
    
    wg.Wait()
    return results.Slice()
}
```

See the `example/main.go` file for more comprehensive examples covering all features.

Run the example:

```bash
go run ./example/main.go
```

## Package Structure

- `arena.go` - the Arena type, the Allocator contract, and the generic helpers
- `limits.go` - the bounds those helpers enforce
- `alloc/` - allocator implementations
  - `bump.go` - bump allocator, a cursor over growing chunks
  - `slab.go` - slab allocator, size classes with individual free
  - `buddy.go` - buddy allocator, splitting and merging power-of-2 blocks
  - `chunk.go` - one buddy chunk and its free-block tree
- `res/` - the raw memory everything is built on
  - `mem.go` - mmap and munmap
  - `page.go` - one mapping, with its bounds
  - `table.go` - a set of pages, and ownership lookup
  - `bump.go` - a cursor over a page table
  - `stats.go`, `errors.go`, `utilities.go`
- `container/` - Vec, Map, SkipList, Pool, Queue, Stack, Buffer, Str
- `io/` - Reader and Writer over arena memory
- `test/`, `container/test/`, `io/test/` - the suites

## License

See LICENSE file for details.