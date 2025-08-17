#!/bin/bash
set -e

# Analyze files in a directory and generate statistics
# yupsh equivalent: See main.go

# Get directory from command line, default to current directory
# yupsh: dir := "."; if len(os.Args) > 1 { dir = os.Args[1] }
DIR=${1:-.}
echo "Analyzing files in: ${DIR}"

# === File Count by Type ===
echo "=== File Count by Type ==="
# Find all files with extensions, extract extension, count occurrences
# yupsh: Pipeline with find, While(extractExtension), sort, uniq, sort
find "${DIR}" -type f -name "*.*" \
| while read file; do
    # Extract extension (everything after last dot)
    # yupsh: filepath.Ext(filename) then strings.TrimPrefix(ext, ".")
    echo "${file##*.}"
done \
| sort \              # yupsh: sort.Sort()
| uniq -c \           # yupsh: uniq.Uniq(uniq.Count) - counts and shows count
| sort -nr            # yupsh: sort.Sort(sort.Numeric, sort.Reverse)

echo ""
# === Largest Files ===
echo "=== Largest Files ==="
# Find all files, get size info, sort by size, show top 10
# yupsh: Pipeline with find, While(getFileSize), sort, head
find "${DIR}" -type f -exec ls -la {} \; \
| awk '{print $5 "\t" $9}' \    # Extract size (field 5) and name (field 9)
                                # yupsh: os.Stat(filename) then fmt.Sprintf("%d\t%s", info.Size(), filename)
| sort -nr \                    # Sort numerically, descending
                                # yupsh: sort.Sort(sort.Numeric, sort.Reverse)
| head -10                      # Take top 10
                                # yupsh: head.Head(head.LineCount(10))

echo ""
# === Total Size ===
echo "=== Total Size ==="
# Find all files, extract sizes, sum them up
# yupsh: Pipeline with find, While(getFileSizeOnly), awk.Awk(&totalSizeProgram{})
find "${DIR}" -type f -exec ls -la {} \; \
| awk '{sum += $5} END {print "Total: " sum " bytes"}'
# yupsh: Custom awk program that accumulates sizes:
#   func (p *totalSizeProgram) Action(ctx *awk.Context) {
#     var size int64
#     fmt.Sscanf(ctx.Field(1), "%d", &size)
#     p.sum += size
#   }
#   func (p *totalSizeProgram) End(ctx *awk.Context) (string, error) {
#     return fmt.Sprintf("Total: %d bytes", p.sum), nil
#   }
