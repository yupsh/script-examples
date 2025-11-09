package main

import (
	"fmt"
	"os"

	gloo "github.com/gloo-foo/framework"
	grep "github.com/yupsh/grep"
	head "github.com/yupsh/head"
	pipe "github.com/gloo-foo/pipe"
	seq "github.com/yupsh/seq"
	tail "github.com/yupsh/tail"
	yes "github.com/yupsh/yes"
)

// Demonstrates how pipes are closed when downstream commands end
//
// Key Concept: When a command like head or tail has read enough data,
// it closes its input pipe. This signals upstream commands to stop
// generating data, preventing wasted work.
//
// This is crucial for:
// - Infinite streams (yes, seq with large ranges)
// - Expensive data generation
// - Resource conservation
//
// Shell equivalent: See pipe-closure.sh

func main() {
	fmt.Println("=== Example 1: head closes pipe after 3 lines ===")
	fmt.Println("Generating 10 lines, but head will only read 3...")
	runExample(pipe.Pipeline(
		seq.Seq("1", "10"),
		head.Head(head.LineCount(3)),
	))
	fmt.Println()

	fmt.Println("=== Example 2: tail closes pipe after reading for 5 lines ===")
	fmt.Println("Generating 100 lines, but tail will only keep the last 5...")
	runExample(pipe.Pipeline(
		seq.Seq("1", "100"),
		tail.Tail(tail.LineCount(5)),
	))
	fmt.Println()

	fmt.Println("=== Example 3: Pipeline with head in middle ===")
	fmt.Println("Generate 20 lines -> head keeps 5 -> grep filters for '3'...")
	runExample(pipe.Pipeline(
		seq.Seq("1", "20"),
		head.Head(head.LineCount(5)),
		grep.Grep("3"),
	))
	fmt.Println()

	fmt.Println("=== Example 4: Yes command with head (infinite stream) ===")
	fmt.Println("yes generates infinite output, but head stops it after 3 lines...")
	runExample(pipe.Pipeline(
		yes.Yes("hello"),
		head.Head(head.LineCount(3)),
	))
	fmt.Println()

	fmt.Println("=== Example 5: Multiple heads in sequence ===")
	fmt.Println("Generate 100 -> head 50 -> head 10 -> head 3...")
	runExample(pipe.Pipeline(
		seq.Seq("1", "100"),
		head.Head(head.LineCount(50)),
		head.Head(head.LineCount(10)),
		head.Head(head.LineCount(3)),
	))
	fmt.Println()

	fmt.Println("=== Example 6: Large data generator with early termination ===")
	fmt.Println("Simulating expensive data generation that stops early...")
	// In real scenarios, pipe closure prevents wasted computation
	// Note: We use a smaller number (10000) to keep the demo fast
	runExample(pipe.Pipeline(
		seq.Seq("1", "10000"),
		head.Head(head.LineCount(5)),
	))
	fmt.Println()

	fmt.Println("Done! Notice how pipe closure prevents unnecessary work.")
}

func runExample(cmd gloo.Command) {
	if err := gloo.Run(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

