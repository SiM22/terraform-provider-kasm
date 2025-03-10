#!/bin/bash

# Default directory to analyze
DIRECTORY="."

# Function to display usage information
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -d, --directory DIR   Directory to analyze (default: current directory)"
    echo "  -c, --clear           Clear existing data in the vector database"
    echo "  -h, --help            Display this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -d internal/provider         Analyze the internal/provider directory"
    echo "  $0 -d internal/resources        Analyze the internal/resources directory"
    echo "  $0 -d internal/datasources      Analyze the internal/datasources directory"
    echo "  $0 -c -d internal/client        Clear database and analyze internal/client"
    exit 1
}

# Parse command line arguments
CLEAR_DB=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--directory)
            DIRECTORY="$2"
            shift 2
            ;;
        -c|--clear)
            CLEAR_DB=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Check if the directory exists
if [ ! -d "$DIRECTORY" ]; then
    echo "Error: Directory '$DIRECTORY' does not exist"
    exit 1
fi

# Create a temporary directory for the code to analyze
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Copy the specified directory to the temporary directory
echo "Copying code from '$DIRECTORY' to temporary directory..."
if [ "$DIRECTORY" = "." ]; then
    # For the root directory, copy specific directories that contain Go code
    mkdir -p "$TMP_DIR/internal"
    cp -r "internal/provider" "$TMP_DIR/internal/" 2>/dev/null || true
    cp -r "internal/resources" "$TMP_DIR/internal/" 2>/dev/null || true
    cp -r "internal/datasources" "$TMP_DIR/internal/" 2>/dev/null || true
    cp -r "internal/client" "$TMP_DIR/internal/" 2>/dev/null || true
    cp -r "internal/utils" "$TMP_DIR/internal/" 2>/dev/null || true
    # Copy any Go files in the root directory
    cp -r *.go "$TMP_DIR/" 2>/dev/null || true
else
    # For specific directories, copy as before
    cp -r "$DIRECTORY/"* "$TMP_DIR/" 2>/dev/null || true
fi

# Check if there are any Go files in the directory
GO_FILES=$(find "$TMP_DIR" -name "*.go" | wc -l)
if [ "$GO_FILES" -eq 0 ]; then
    echo "Warning: No Go files found in '$DIRECTORY'"
    # Create a sample Go file for testing if no Go files found
    echo "Creating a sample Go file for testing..."
    cat > "$TMP_DIR/sample.go" << 'EOL'
package main

import (
	"fmt"
)

func main() {
	fmt.Println("This is a sample Go file")
}
EOL
fi

# Run the analyzer on the directory
echo "Running analyzer on directory: $DIRECTORY"
docker run --rm -v "$TMP_DIR":/app/code terraform-provider-kasm-analyzer bash -c "cd /app/code && /app/analyzer" > provider-analysis.json

# Check if the analyzer produced output
if [ ! -s provider-analysis.json ]; then
    echo "Error: Analyzer did not produce any output"
    exit 1
fi

echo "Analyzer output saved to provider-analysis.json"

# Wait for the vector-db to be ready
echo "Waiting for vector-db to be ready..."
sleep 5

# Clear the database if requested
if [ "$CLEAR_DB" = true ]; then
    echo "Clearing existing data in the vector database..."
    curl -X POST "http://localhost:8000/api/clear"
    sleep 2
fi

# Send the analyzer output to the vector database
echo "Sending data to vector database for directory: $DIRECTORY"
SOURCE_NAME="$DIRECTORY"
if [ "$DIRECTORY" = "." ]; then
    SOURCE_NAME="entire-provider"
fi
curl -X POST "http://localhost:8000/api/add" \
    -H "Content-Type: application/json" \
    -d "{\"structures\": $(cat provider-analysis.json), \"source\": \"$SOURCE_NAME\"}"

# Check the count of documents in the database
echo "Checking document count..."
curl "http://localhost:8000/api/count"

echo "Done! You can now access the UI at http://localhost:8000"
echo "Analyzed directory: $DIRECTORY"
