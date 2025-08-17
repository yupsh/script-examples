# yupsh Script Examples

This directory contains real-world examples demonstrating how to convert shell scripts into yupsh Go programs.

## Examples

### 📊 [log-processor](./log-processor/)
Processes log files to extract errors and warnings, demonstrating:
- File iteration patterns
- Text filtering and field extraction
- Output redirection with `tee.Tee()`
- The `while.While()` command for line-by-line processing

```bash
cd log-processor
go run main.go
```

### 📈 [file-stats](./file-stats/)
Analyzes files in a directory to generate statistics, demonstrating:
- File system operations with `find.Find()`
- Custom data processing functions
- Using `awk.Awk()` for calculations
- Combining multiple pipelines

```bash
cd file-stats
go run main.go [directory]
```

## Common Patterns

### 1. Shell Loop → yupsh Pipeline with `While()`
```bash
# Shell
while read line; do
    # process line
done
```

```go
// yupsh
import . "github.com/yupsh/while"

pipe.Pipeline(
    cat.Cat("input.txt"),
    While(func(args ...any) gloo.Command {
        line := args[0].(string)
        return echo.Echo("processed: " + line)
    }),
)
```

### 2. File Iteration
```bash
# Shell
for file in logs/*.log; do
    grep "ERROR" "$file"
done
```

```go
// yupsh
pipe.Pipeline(
    ls.Ls("logs/*.log"),
    While(func(args ...any) gloo.Command {
        filename := args[0].(string)
        return pipe.Pipeline(
            cat.Cat(filename),
            grep.Grep("ERROR"),
        )
    }),
)
```

### 3. Field Extraction with Custom Functions
```bash
# Shell
find . -name "*.go" | while read file; do
    echo "${file##*.}"  # Extract extension
done
```

```go
// yupsh
pipe.Pipeline(
    find.Find(find.Dir("."), find.Name("*.go")),
    While(func(args ...any) gloo.Command {
        filename := args[0].(string)
        ext := filepath.Ext(filename)
        return echo.Echo(strings.TrimPrefix(ext, "."))
    }),
)
```

## Getting Started

Each example directory contains:
- `main.go` - The yupsh Go program
- `*.sh` - The equivalent shell script
- `README.md` - Detailed explanation of patterns used

Run any example to see yupsh in action, then compare with the shell script to understand the translation patterns.

