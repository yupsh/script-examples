package main

import (
	"fmt"
	"os"

	cat `github.com/yupsh/cat`
	echo `github.com/yupsh/echo`
	gloo `github.com/gloo-foo/framework`
	grep `github.com/yupsh/grep`
	ls `github.com/yupsh/ls`
	pipe `github.com/gloo-foo/pipe`
	tee `github.com/yupsh/tee`
	. `github.com/yupsh/while`
)

// Process log files to extract errors and warnings
// Shell equivalent: See process-logs.sh
//
// This demonstrates the yupsh pattern for converting shell pipelines with
// nested while loops into Go programs.
//
// Key pattern: Shell's "while read" loops become While() commands in yupsh.
// Each While() receives a callback function that processes one line (or set
// of fields) at a time.
func main() {
	// Main pipeline: List log files and process each one
	// Shell: ls -1 logs/*.log | while read -r file; do ... done
	err := gloo.Run(pipe.Pipeline(
		// List all .log files in logs/ directory
		// Shell: ls -1 logs/*.log
		ls.Ls("logs/*.log"),

		// For each filename (one per line), call processLogFile()
		// Shell: while read -r file; do ... done
		// The While() command reads each line and passes it as args[0]
		While(processLogFile),
	))

	if err != nil {
		fmt.Fprintf(os.Stderr, "log-processor: %v\n", err)
		os.Exit(1)
	}
}

// processLogLine extracts timestamp and level from each matching log line
//
// Shell equivalent:
//   timestamp=$(echo "$line" | cut -d' ' -f1)
//   level=$(echo "$line" | cut -d' ' -f2)
//   echo "${timestamp},${level}" >> results.csv
//
// yupsh pattern: When FieldSeparator(" ") is specified in While(), each line
// is automatically split on whitespace. The fields are passed as separate
// arguments to this function:
//   args[0] = first field (timestamp)
//   args[1] = second field (level)
//   args[2..n] = remaining fields (if any)
//
// This eliminates the need for manual field extraction with cut/awk.
func processLogLine(args ...any) gloo.Command {
	// Extract the fields we need
	if len(args) < 2 {
		return nil // Skip malformed lines (safety check)
	}
	timestamp := args[0].(string)
	level := args[1].(string)

	return pipe.Pipeline(
		// Format as CSV: timestamp,level
		// Shell: echo "${timestamp},${level}"
		echo.Echo(fmt.Sprintf("%s,%s", timestamp, level)),

		// Append to results.csv
		// Shell: >> results.csv
		// tee.Append makes it append instead of overwrite
		tee.Tee("results.csv", tee.Append),
	)
}

// processLogFile reads a log file and extracts errors/warnings
//
// Shell equivalent:
//   grep -i "error\|warning" "${file}" | while read -r line; do
//     timestamp=$(echo "$line" | cut -d' ' -f1)
//     level=$(echo "$line" | cut -d' ' -f2)
//     echo "${timestamp},${level}" >> results.csv
//   done
//
// This function is called once per filename from the outer While() loop.
// It creates a nested pipeline to process each file.
func processLogFile(args ...any) gloo.Command {
	// args[0] is the filename from ls.Ls() output
	// Shell: while read -r file; do ... "${file}" ... done
	filename := args[0].(string)
	filepath := "logs/" + filename

	// Print progress to stderr (won't interfere with pipeline)
	// Shell: echo "Processing ${file}"
	fmt.Fprintf(os.Stderr, "Processing %s\n", filepath)

	return pipe.Pipeline(
		// Read the file contents
		// Shell: (implicit - grep reads the file)
		cat.Cat(filepath),

		// Filter for lines containing "error" or "warning" (case insensitive)
		// Shell: grep -i "error\|warning" "${file}"
		// Note: yupsh uses "|" for regex alternation instead of "\|"
		grep.Grep("error|warning", grep.IgnoreCase),

		// For each matching line, split on whitespace and extract fields
		// Shell: while read -r line; do
		//          timestamp=$(echo "$line" | cut -d' ' -f1)
		//          level=$(echo "$line" | cut -d' ' -f2)
		//        done
		//
		// yupsh: FieldSeparator(" ") automatically splits each line on spaces
		// and passes the fields as separate args to processLogLine()
		While(processLogLine, FieldSeparator(" ")),
	)
}

