#!/bin/sh
# This script is a wrapper for the t-walker program. It captures the output of the program and allows us to parse
# any commands that the program would like to execute on the calling shell.
# For example, hitting the 'c' key will cause t-walker to exit and switch the shell to the directory
# via this script.

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
