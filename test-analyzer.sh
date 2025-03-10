#!/bin/bash

# Run the analyzer and capture its output
echo "Running analyzer..."
ANALYZER_OUTPUT=$(docker-compose run --rm analyzer)

# Check if the analyzer produced output
if [ -z "$ANALYZER_OUTPUT" ]; then
    echo "Error: Analyzer did not produce any output"
    exit 1
fi

echo "Analyzer output:"
echo "$ANALYZER_OUTPUT"

# Wait for the vector-db to be ready
echo "Waiting for vector-db to be ready..."
sleep 5

# Send the analyzer output to the vector database
echo "Sending data to vector database..."
curl -X POST "http://localhost:8000/add" \
    -H "Content-Type: application/json" \
    -d "{\"structures\": $ANALYZER_OUTPUT}"

# Check the count of documents in the database
echo "Checking document count..."
curl "http://localhost:8000/count"

# Test a query
echo "Testing a query for 'main'..."
curl -X POST "http://localhost:8000/query" \
    -H "Content-Type: application/json" \
    -d '{"text": "main function"}'
