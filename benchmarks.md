# Allocator benchmarks

The numbers: where the time goes, why each allocator is as fast as it is, and what is
still on the table.

Measured on 2026-09-21, Go 1.26.2, linux/amd64, 12th Gen Intel Core i7-12800H,
five runs of each benchmark. Reproduce with:

```bash
go test -run '^$' -bench 'BenchmarkBump_' -benchmem -count 5 ./test/
```

Every figure here comes from a fresh run of the command above. A single run of a
benchmark proves nothing; the comparison against the previous implementation
further down was taken with `benchstat` over six runs of each side, then repeated with
the two sides swapped, because an unswapped comparison reports thermal drift as a
regression in whichever side ran second.

## Bump allocator

| Operation | ns/op | B/op | allocs/op |
| --- | --- | --- | --- |
| `Alloc` | 12.4 | 0 | 0 |
| `Alloc` with periodic `Reset` | 12.4 | 0 | 0 |
| `Owns` | 11.5 | 0 | 0 |
| `MakeString` | 21.2 | 0 | 0 |
| `MakeSlice` | 37.9 | 0 | 0 |
| Mixed workload | 35.6 | 0 | 0 |

`B/op` and `allocs/op` count the **Go heap**. Zero across the board is the point
of the library: the bytes come from `mmap` and the collector never sees them.

## Buddy allocator

| Operation | ns/op | B/op | allocs/op |
| --- | --- | --- | --- |
| `Alloc` | 54.1 | 0 | 0 |
| `Alloc` with periodic `Reset` | 40.6 | 0 | 0 |

`Reset` is cheaper than a plain allocation because resetting the chunk hands the whole
tree back at once, so the allocations after it never have to split anything.

## The three allocators compared

| Allocator | `Alloc` ns/op | Frees individually | Best for |
| --- | --- | --- | --- |
| Bump | 12.4 | no, whole arena only | batches, request handlers, anything reset often |
| Slab | 23.0 | yes | fixed-size objects with high turnover |
| Buddy | 54.1 | yes | varied sizes that must be freed one at a time |

Buddy costs more per allocation because it walks a tree to find and split a
block, and merges on the way back out. That is what buys individual free of any
size with low fragmentation.

## Against the previous implementation

Commit `acb966c`, same machine, `benchstat` over six runs:

| Benchmark | Before | After | Change |
| --- | --- | --- | --- |
| `BenchmarkBuddy_Alloc` | 35 972 ns/op | 54.1 ns/op | −99.85% |
| `BenchmarkSlab_Alloc` | 39.1 ns/op | 23.0 ns/op | −41.1% |
| `BenchmarkBump_Alloc` | 21.6 ns/op | 12.4 ns/op | −42.6% |

Where the time went:

- **Buddy** searched a one-bit-per-node bitmap that could not tell a full
  subtree from a partly used one, so every allocation re-walked every live
  block. Per-operation cost therefore grew with the number of live blocks; the
  35.9 µs above is that cost at the benchmark's working set, not a constant. It
  is now one byte per node holding the largest free block in that subtree, which
  makes both allocation and free O(log n).
- **Slab** kept its page index sorted by inserting in place. Linux hands out
  `mmap` addresses downward, so every new page sorted before all the others and
  the insert moved the whole slice. It now appends and sorts only when a lookup
  needs the order.
- **Bump** took two mutexes per allocation, one in the allocator and one in the
  cursor underneath it, to protect a single offset. It takes one.

## What the numbers do not say

These are single-goroutine benchmarks of one operation in a loop. They do not
measure contention, they do not measure what a real allocation pattern does to
fragmentation, and they say nothing about the cost the Go allocator avoids by
not having to scan this memory. Use them to choose between the three
strategies, not to predict a program's throughput.

One trap worth naming, because it hid the buddy allocator's behaviour for a long
time: a benchmark that frees or resets between iterations never builds up a live
set, so an allocator whose cost grows with the live set looks constant. If you
add a benchmark here, consider whether it needs a companion that does not reset.
`TestBuddySequentialAllocScales` is the test that checks for exactly that.

## Running the rest

```bash
# every benchmark, with memory statistics
make bench

# one allocator
make bench BENCH=BenchmarkSlab_

# compare the working tree against a commit
git worktree add --detach /tmp/base HEAD
(cd /tmp/base && go test -run '^$' -bench 'BenchmarkBump_' -benchmem -count 6 ./test/) >old.txt
go test -run '^$' -bench 'BenchmarkBump_' -benchmem -count 6 ./test/ >new.txt
go tool benchstat old.txt new.txt
git worktree remove /tmp/base
```
