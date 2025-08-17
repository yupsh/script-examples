package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	awk `github.com/yupsh/awk`
	echo `github.com/yupsh/echo`
	find `github.com/yupsh/find`
	gloo `github.com/gloo-foo/framework`
	head `github.com/yupsh/head`
	pipe `github.com/gloo-foo/pipe`
	sort `github.com/yupsh/sort`
	uniq `github.com/yupsh/uniq`
	. `github.com/yupsh/while`
)

// Analyze files in a directory and generate statistics
// Shell equivalent: See analyze-files.sh
//
// This demonstrates three common yupsh patterns:
// 1. Using While() with custom functions to process filenames
// 2. Combining find, sort, uniq, head commands in pipelines
// 3. Using awk.Awk() with custom programs for aggregations
//
// Key advantage: Native Go file operations (os.Stat) instead of parsing ls output
func main() {
	// Get directory from command line, default to current directory
	// Shell: DIR=${1:-.}
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fmt.Fprintf(os.Stderr, "Analyzing files in: %s\n", dir)

	// === File Count by Type ===
	// Shell: find | while read | sort | uniq -c | sort -nr
	fmt.Fprintf(os.Stderr, "\n=== File Count by Type ===\n")
	err := gloo.Run(pipe.Pipeline(
		// Find all files with extensions
		// Shell: find "${DIR}" -type f -name "*.*"
		find.Find(find.Dir(dir), find.FileType, find.Name("*.*")),

		// Extract extension from each filename
		// Shell: while read file; do echo "${file##*.}"; done
		// yupsh: While() calls extractExtension() for each line
		//        extractExtension uses filepath.Ext() and strings.TrimPrefix()
		While(extractExtension),

		// Sort extensions alphabetically (prepares for uniq)
		// Shell: sort
		sort.Sort(),

		// Count occurrences of each unique extension
		// Shell: uniq -c
		// Output format: "  5 go" (count followed by value)
		uniq.Uniq(uniq.Count),

		// Sort by count in descending order (most common first)
		// Shell: sort -nr (numeric, reverse)
		sort.Sort(sort.Numeric, sort.Reverse),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// === Largest Files ===
	// Shell: find | ls -la | awk | sort -nr | head -10
	fmt.Fprintf(os.Stderr, "\n=== Largest Files ===\n")
	err = gloo.Run(pipe.Pipeline(
		// Find all files (no name filter this time)
		// Shell: find "${DIR}" -type f
		find.Find(find.Dir(dir), find.FileType),

		// Get size and name for each file
		// Shell: find -exec ls -la {} \; | awk '{print $5 "\t" $9}'
		// yupsh: While() calls getFileSize() which uses os.Stat()
		//        This is better than parsing ls output - native Go!
		//        Output format: "12485\t./README.md" (size TAB filename)
		While(getFileSize),

		// Sort by size in descending order (largest first)
		// Shell: sort -nr (numeric, reverse)
		// Note: sort.Numeric tells sort to compare numbers, not strings
		sort.Sort(sort.Numeric, sort.Reverse),

		// Take only the top 10 largest files
		// Shell: head -10
		head.Head(head.LineCount(10)),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// === Total Size ===
	// Shell: find | ls -la | awk '{sum += $5} END {print "Total: " sum " bytes"}'
	fmt.Fprintf(os.Stderr, "\n=== Total Size ===\n")
	err = gloo.Run(pipe.Pipeline(
		// Find all files
		// Shell: find "${DIR}" -type f
		find.Find(find.Dir(dir), find.FileType),

		// Get just the file size (not the name)
		// Shell: (implicit in awk processing)
		// yupsh: While() calls getFileSizeOnly() which uses os.Stat()
		//        Output format: "12485" (just the size in bytes)
		While(getFileSizeOnly),

		// Sum all sizes using a custom awk program
		// Shell: awk '{sum += $5} END {print "Total: " sum " bytes"}'
		// yupsh: Custom totalSizeProgram that accumulates and formats output
		//        Action() accumulates each size, End() prints the total
		awk.Awk(&totalSizeProgram{}),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// extractExtension extracts the file extension from a filepath
//
// Shell equivalent:
//   echo "${file##*.}"
//
// This is called by While() for each filename from find.Find().
// It demonstrates using native Go functions instead of shell string manipulation.
//
// Shell pattern: ${file##*.} removes everything up to last dot
// yupsh pattern: filepath.Ext() gets extension, strings.TrimPrefix() removes dot
func extractExtension(args ...any) gloo.Command {
	// args[0] is the filename from find.Find() output
	filename := args[0].(string)

	// Get the extension (everything after last dot)
	// filepath.Ext() returns ".go", ".txt", etc.
	ext := filepath.Ext(filename)
	if ext == "" {
		// Skip files without extensions
		// Shell: (would still output empty line)
		// yupsh: return nil to skip this line entirely
		return nil
	}

	// Remove the leading dot: ".go" -> "go"
	// Shell: ${file##*.} doesn't include the dot
	ext = strings.TrimPrefix(ext, ".")

	// Output just the extension
	// Shell: echo "${file##*.}"
	return echo.Echo(ext)
}

// getFileSize gets the file size and name formatted as "size\tname"
//
// Shell equivalent:
//   find -exec ls -la {} \; | awk '{print $5 "\t" $9}'
//
// This demonstrates a key yupsh advantage: instead of spawning ls and parsing
// its output with awk, we use native Go's os.Stat() to get file information
// directly. This is faster, more reliable, and cross-platform.
//
// Shell approach: find -exec ls -la {} \; spawns ls for each file
// yupsh approach: os.Stat() is a single function call, no subprocess
func getFileSize(args ...any) gloo.Command {
	// args[0] is the filename from find.Find() output
	filename := args[0].(string)

	// Get file info using native Go
	// Shell: ls -la outputs multiple fields, awk extracts field 5 (size)
	// yupsh: os.Stat() gives us structured data directly
	info, err := os.Stat(filename)
	if err != nil {
		// Skip files we can't access (permissions, deleted, etc.)
		// Shell: ls would print error to stderr but continue
		// yupsh: return nil to skip this line
		return nil
	}

	// Format as "size\tname" for sorting
	// Shell: awk '{print $5 "\t" $9}'
	// Output example: "12485\t./README.md"
	return echo.Echo(fmt.Sprintf("%d\t%s", info.Size(), filename))
}

// getFileSizeOnly gets just the file size (for summing)
//
// Shell equivalent:
//   find -exec ls -la {} \; | awk '{print $5}'
//
// This is similar to getFileSize but only outputs the size,
// not the filename. Used for the total size calculation.
func getFileSizeOnly(args ...any) gloo.Command {
	// args[0] is the filename from find.Find() output
	filename := args[0].(string)

	// Get file info using native Go
	info, err := os.Stat(filename)
	if err != nil {
		return nil // Skip files we can't access
	}

	// Return just the size in bytes
	// Shell: awk '{print $5}' (field 5 from ls -la)
	// Output example: "12485"
	return echo.Echo(fmt.Sprintf("%d", info.Size()))
}

// totalSizeProgram is a custom awk program that sums all input numbers
//
// Shell equivalent:
//   awk '{sum += $1} END {print "Total: " sum " bytes"}'
//
// This demonstrates how to use yupsh's awk.Awk() command with a custom
// program. The shell's awk has three sections: BEGIN, Action, and END.
// We only need Action (process each line) and END (output final result).
//
// Shell awk pattern:
//   {sum += $1}                         - Action: add field 1 to sum
//   END {print "Total: " sum " bytes"}  - End: print final total
//
// yupsh pattern:
//   Action() - called for each input line
//   End() - called once at the end
type totalSizeProgram struct {
	awk.SimpleProgram // Provides basic awk program structure
	sum int64         // Accumulator for total size
}

// Action is called for each input line
// Shell: {sum += $1}
func (p *totalSizeProgram) Action(ctx *awk.Context) (string, bool) {
	// Parse the size from field 1 (the only field in our input)
	// Shell: $1 (automatic in awk)
	// yupsh: ctx.Field(1) (explicit field access)
	var size int64
	fmt.Sscanf(ctx.Field(1), "%d", &size)

	// Add to running total
	// Shell: sum += $1
	p.sum += size

	// Don't emit anything during processing (only at the end)
	// Shell: (no print statement in action, so nothing output)
	// yupsh: return "", false (empty string, don't emit)
	return "", false
}

// End is called once after all lines are processed
// Shell: END {print "Total: " sum " bytes"}
func (p *totalSizeProgram) End(ctx *awk.Context) (string, error) {
	// Format and return the total
	// Shell: print "Total: " sum " bytes"
	return fmt.Sprintf("Total: %d bytes", p.sum), nil
}

