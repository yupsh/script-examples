# Pipe Closure Example

This example demonstrates a critical concept in Unix pipelines: **how pipes are closed when downstream commands finish reading**.

## The Problem This Solves

When you have a pipeline like:

```bash
seq 1 1000000 | head -n 5
```

You don't want `seq` to generate all 1,000,000 numbers. Instead, once `head` has read 5 lines, it should **close the pipe**, signaling `seq` to stop generating data.

This is crucial for:
- **Performance**: Avoid wasted computation
- **Resource conservation**: Don't generate unnecessary data
- **Infinite streams**: Commands like `yes` would run forever without pipe closure
- **Memory efficiency**: Large datasets don't need to be fully materialized

## How It Works

### Shell Behavior

In Unix shells, when a downstream command (like `head`) closes its stdin:
1. The pipe between commands is closed
2. Upstream commands receive `SIGPIPE` when they try to write
3. Proper programs handle this gracefully and exit
4. The pipeline completes efficiently

### yupsh Behavior

The yupsh framework handles this automatically:
1. Commands like `head.Head()` and `tail.Tail()` stop reading after their quota
2. The framework properly closes pipes between commands
3. Upstream commands detect the closed pipe and stop generating data
4. No error is reported (this is expected behavior)

## Examples in This Demo

### 1. Basic Head Usage
```go
pipe.Pipeline(
    seq.Seq("1", "10"),
    head.Head(head.LineCount(3)),
)
```
**Shell equivalent:**
```bash
seq 1 10 | head -n 3
```

Generates 10 numbers, but `head` only reads 3. The remaining 7 are never generated after the pipe closes.

### 2. Tail Usage
```go
pipe.Pipeline(
    seq.Seq("1", "100"),
    tail.Tail(tail.LineCount(5)),
)
```
**Shell equivalent:**
```bash
seq 1 100 | tail -n 5
```

`tail` must read all input to find the last N lines, but still demonstrates proper pipe handling.

### 3. Head in the Middle
```go
pipe.Pipeline(
    seq.Seq("1", "20"),
    head.Head(head.LineCount(5)),
    grep.Grep("3"),
)
```
**Shell equivalent:**
```bash
seq 1 20 | head -n 5 | grep 3
```

Shows how `head` limits data flow to downstream commands.

### 4. Infinite Streams
```go
pipe.Pipeline(
    yes.Yes("hello"),
    head.Head(head.LineCount(3)),
)
```
**Shell equivalent:**
```bash
yes "hello" | head -n 3
```

**Critical example!** `yes` would run forever, but `head` stops it after 3 lines. Without proper pipe closure, this would never terminate.

### 5. Cascading Heads
```go
pipe.Pipeline(
    seq.Seq("1", "100"),
    head.Head(head.LineCount(50)),
    head.Head(head.LineCount(10)),
    head.Head(head.LineCount(3)),
)
```
**Shell equivalent:**
```bash
seq 1 100 | head -n 50 | head -n 10 | head -n 3
```

Each `head` further limits the data, demonstrating multiple closure points.

### 6. Expensive Generation
```go
pipe.Pipeline(
    seq.Seq("1", "10000"),
    head.Head(head.LineCount(5)),
)
```
**Shell equivalent:**
```bash
for i in {1..1000000}; do echo "Processing $i"; done | head -n 5
```

In real applications, this pattern prevents expensive computations from running unnecessarily.

## Key Patterns

### Pattern 1: Early Termination of Expensive Operations

**Problem:** You have an expensive data generation pipeline but only need a few results.

**Solution:** Use `head` to limit output and trigger pipe closure.

```go
pipe.Pipeline(
    // Expensive operation: scan large directory tree
    find.Find(find.Dir("/large/directory")),
    // Only need first 10 results
    head.Head(head.LineCount(10)),
)
```

### Pattern 2: Sampling Infinite Streams

**Problem:** You want to test or sample from an infinite data source.

**Solution:** Combine infinite generators with `head`.

```go
pipe.Pipeline(
    yes.Yes("test"),  // or any infinite generator
    head.Head(head.LineCount(100)),
)
```

### Pattern 3: Preview Large Files

**Problem:** You want to quickly peek at a large file's contents.

**Solution:** Use `head` to read just the beginning.

```go
pipe.Pipeline(
    cat.Cat("huge-file.txt"),
    head.Head(head.LineCount(20)),
)
```

### Pattern 4: Cascading Filters

**Problem:** You have multiple filtering stages and want to stop as soon as possible.

**Solution:** Place `head` strategically in the pipeline.

```go
pipe.Pipeline(
    cat.Cat("data.txt"),
    grep.Grep("ERROR"),
    head.Head(head.LineCount(5)),  // Stop after 5 errors found
    // Further processing...
)
```

## Technical Details

### SIGPIPE Handling

In Unix, writing to a closed pipe generates a `SIGPIPE` signal. Proper programs handle this by:
1. Catching the signal (or ignoring it)
2. Detecting the write error
3. Exiting gracefully

The yupsh framework handles this automatically for all commands.

### Context Cancellation

The yupsh implementation uses Go's context package to propagate cancellation:
1. When a downstream command closes, it cancels its context
2. Upstream commands check the context
3. The entire pipeline shuts down efficiently

### Performance Benefits

Pipe closure enables:
- **O(N)** instead of **O(M)** where M >> N
- Early resource cleanup
- Reduced memory usage
- Faster overall execution

## Running the Examples

```bash
cd pipe-closure

# Run the shell version
bash pipe-closure.sh

# Run the yupsh version
go run main.go
```

Both versions should produce similar output, demonstrating proper pipe closure behavior.

## Common Pitfalls

### ❌ Not Handling Pipe Closure

```go
// BAD: Custom generator that doesn't check for context cancellation
While(func(args ...any) gloo.Command {
    for i := 0; i < 1000000; i++ {
        // This will keep running even if downstream closed
        echo.Echo(fmt.Sprintf("%d", i))
    }
    return nil
})
```

### ✅ Proper Context Checking

```go
// GOOD: Check context or let framework handle it
pipe.Pipeline(
    seq.Seq(1, 1000000),  // Built-in commands handle this
    head.Head(head.Lines(10)),
)
```

## Related Concepts

- **Backpressure**: Downstream commands can slow upstream generation
- **Buffering**: Pipes have limited buffers (typically 64KB in Unix)
- **SIGPIPE**: Signal sent when writing to closed pipe
- **Context cancellation**: Go's mechanism for propagating stop signals

## Further Reading

- Unix pipe(2) man page
- Go context package documentation
- "The Art of Unix Programming" by Eric S. Raymond
- yupsh framework documentation on pipeline execution

