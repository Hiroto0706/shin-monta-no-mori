#!/bin/bash

# Check if a coverage report file was passed as an argument
if [[ -z "$1" ]]; then
    echo "please this format -> $0 <coverage_report.txt>"
    exit 1
fi

# Aggregate coverage by directory
awk -F'[:\t]+' '{
    # Extract the directory path from the filename
    dir = $1;
    sub("/[^/]*$", "", dir);  # Remove the file name, keep the directory path

    # Parse coverage percentage
    coverage = $4;  # Assumes that coverage percentage is in the fourth field

    # Remove any extra spaces or percent sign
    gsub(/[% ]/, "", coverage);

    # Sum coverage percentages and count them for each directory
    total[dir] += coverage;
    count[dir]++;

    # Keep track of overall total and count for average calculation
    overall_total += coverage;
    overall_count++;
} END {
    # Print average coverage for each directory
    for (d in total) {
        if (d == "total") continue; # d = totalの時はスキップ
        average = total[d] / count[d];
        printf "%s: %.2f%%\n", d, average;
    }

    # Calculate and print the overall average coverage
    if (overall_count > 0) {
        overall_average = overall_total / overall_count;
        printf "Total average coverage: %.2f%%\n", overall_average;
    }
}' "$1"
