#!/bin/bash
set -e

# Process log files to extract errors and warnings
# yupsh equivalent: See main.go

# List all .log files in logs/ directory
# yupsh: ls.Ls("logs/*.log")
ls -1 logs/*.log \
| while read -r file; do
  # For each file, process it
  # yupsh: While(processLogFile)

  if ! [[ -f "${file}" ]]; then
    continue
  fi

  echo "Processing ${file}"

  # Read file and filter for errors/warnings (case insensitive)
  # yupsh: cat.Cat(filepath), grep.Grep("error|warning", grep.IgnoreCase)
  grep -i "error\|warning" "${file}" \
  | while read -r line; do
    # For each matching line, extract fields
    # yupsh: While(processLogLine, FieldSeparator(" "))
    #        This automatically splits the line on whitespace into args

    # Extract timestamp (1st field) and level (2nd field)
    # yupsh: args[0].(string) and args[1].(string)
    timestamp=$(echo "${line}" | cut -d' ' -f1)
    level=$(echo "${line}" | cut -d' ' -f2)

    # Write CSV output to file
    # yupsh: echo.Echo(fmt.Sprintf("%s,%s", timestamp, level))
    #        tee.Tee("results.csv", tee.Append)
    echo "${timestamp},${level}" >> results.csv
  done
done
