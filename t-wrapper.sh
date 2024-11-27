#!/bin/sh

# Create a temporary file to store the output
tempfile=$(mktemp)

# Execute the t-wrapper.sh-walker program, capture its output, and display the TUI
script -q "$tempfile" t-walker

# Read the captured output into a variable
output=$(cat "$tempfile")

## Clean up the temporary file
rm "$tempfile"

parsedOutput=$(echo "$output" | sed -n 's/.*Executing command: \x1b\[32m \(.*\) \x1b\[0m.*/\1/p' | sed 's/^ *//')
eval "$parsedOutput"
