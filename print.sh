#!/bin/bash

# Check if a directory argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi

# Directory to search in
DIR="$1"

# Find all .go files recursively in the specified directory and cat their contents
find "$DIR" -type f -name '*.go' -exec cat {} +