# File Statistics Example

Analyzes files in a directory and generates three types of statistics:
1. **File count by type** - Groups files by extension and counts them
2. **Largest files** - Shows the 10 largest files by size
3. **Total size** - Calculates total size of all files

## Running

**Shell version:**
```bash
./analyze-files.sh [directory]
```

**yupsh Go version:**
```bash
go run main.go [directory]
```

Both produce identical output.

## Learning

The code files are heavily commented to show the direct translation between shell and Go:
- `analyze-files.sh` - Shell script with comments showing the yupsh equivalent
- `main.go` - Go program with comments showing the shell equivalent

Read both side-by-side to understand the patterns.
