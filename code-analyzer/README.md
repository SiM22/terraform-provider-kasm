# Code Analysis System

This system provides semantic search capabilities for your Go codebase using Chroma vector database.

## Components

1. **Go Parser**: Extracts code structure using Go AST
2. **Vector Database**: Stores and queries code embeddings

## Usage

1. Build and run the system:
```bash
docker-compose up --build
```

2. Query the codebase:
```bash
curl -X POST http://localhost:8000/query -H "Content-Type: application/json" -d '{"text": "your query here"}'
```
