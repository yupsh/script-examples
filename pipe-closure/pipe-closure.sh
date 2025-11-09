#!/bin/bash

# pipe-closure.sh
# Demonstrates how pipes are closed when downstream commands end
#
# Key concept: Commands like head/tail read only N lines then close the pipe.
# This causes upstream commands to receive SIGPIPE when they try to write more.

echo "=== Example 1: head closes pipe after 3 lines ==="
echo "Generating 10 lines, but head will only read 3..."
seq 1 10 | head -n 3
echo

echo "=== Example 2: tail closes pipe after reading for 5 lines ==="
echo "Generating 100 lines, but tail will only keep the last 5..."
seq 1 100 | tail -n 5
echo

echo "=== Example 3: Pipeline with head in middle ==="
echo "Generate 20 lines -> head keeps 5 -> grep filters for '3'..."
seq 1 20 | head -n 5 | grep 3
echo

echo "=== Example 4: Yes command with head (infinite stream) ==="
echo "yes generates infinite output, but head stops it after 3 lines..."
yes "hello" | head -n 3
echo

echo "=== Example 5: Multiple heads in sequence ==="
echo "Generate 100 -> head 50 -> head 10 -> head 3..."
seq 1 100 | head -n 50 | head -n 10 | head -n 3
echo

echo "=== Example 6: Large data generator with early termination ==="
echo "Simulating expensive data generation that stops early..."
# In real scenarios, this prevents wasted computation
(for i in {1..1000000}; do echo "Processing expensive item $i"; done) | head -n 5
echo

echo "Done! Notice how pipe closure prevents unnecessary work."

