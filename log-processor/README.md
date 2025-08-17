# Log Processor Example

Processes log files to extract error and warning entries, outputting timestamp and log level in CSV format.

## Running

**Shell version:**
```bash
./process-logs.sh
```

**yupsh Go version:**
```bash
go run main.go
```

Both produce identical `results.csv` output.

## Learning

The code files are heavily commented to show the direct translation between shell and Go:
- `process-logs.sh` - Shell script with comments showing the yupsh equivalent
- `main.go` - Go program with comments showing the shell equivalent

Read both side-by-side to understand the patterns.
