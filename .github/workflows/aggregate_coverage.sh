#!/bin/bash
# aggregate_coverage.sh

# Initialize associative array
declare -A cov_lines
declare -A cov_hits

# Read coverage data
while IFS=':' read -r filename pct; do
    dir=$(dirname "$filename")
    cov_lines["$dir"]+=$(echo "$pct" | awk '{print $1}')
    cov_hits["$dir"]+=$(echo "$pct" | awk '{print $2}')
done < "$1"

# Calculate and print coverage per directory
for dir in "${!cov_lines[@]}"; do
    lines=${cov_lines["$dir"]}
    hits=${cov_hits["$dir"]}
    coverage=$(awk "BEGIN {printf \"%.2f\", ($hits/$lines)*100}")
    echo "$dir coverage: $coverage%"
done
