# Terraform Provider Kasm - Architecture Documentation

## Overview
The Terraform Provider for Kasm is structured to provide a clean separation of concerns between the API client, resource implementations, data sources, and testing infrastructure. This document outlines the code organization, design principles, and testing strategy.

## Directory Structure

```
tf-provider-kasm/
├── docs/                    # Provider documentation
│   ├── data-sources/       # Data source documentation
│   ├── guides/             # User guides
│   ├── resources/          # Resource documentation
│   └── testing.md          # Testing documentation
├── internal/               # Internal provider code
│   ├── client/            # API client implementation
│   │   ├── *_ops.go       # API operations
│   │   └── *_types.go     # API types
│   ├── datasources/       # Data source implementations
│   │   └── */             # Each data source in its own package
│   ├── resources/         # Resource implementations
│   │   └── */             # Each resource in its own package
│   │       ├── resource.go # Resource implementation
│   │       └── tests/     # Resource-specific tests
│   ├── provider/          # Provider configuration
│   └── validators/        # Common validation functions
└── testutils/             # Testing utilities
    ├── debugger/          # Debug logging and error capture
    └── aicapture/         # AI-assisted debugging tools
```

## Component Design

### 1. Client Layer (`internal/client/`)
The client layer handles all API communication with Kasm.

#### Key Components:
- `client.go`: Core client implementation with configuration and HTTP client setup.
- `*_types.go`: Type definitions for API requests/responses.
- `*_ops.go`: API operation implementations.
- `backoff.go`: Retry logic and backoff strategies.
- `errors.go`: Error type definitions and handling.
- `http.go`: HTTP client configuration and middleware.

Example:
```go
// client/image_types.go - Type definitions

// Image represents a Kasm image

type Image struct {
    ImageID string `json:"image_id"`
    Name    string `json:"name"`
    // ... other fields
}
```

#### Test Patterns
- Tests for the client layer are organized by functionality, ensuring that all API interactions are covered.
- Utilize mocking to simulate API responses for unit tests, allowing for comprehensive coverage without relying on live API calls.

### 2. Resource Layer (`internal/resources/`)
- Each resource has its own package, containing the resource implementation and associated tests.
- Resources are designed to follow Terraform's resource lifecycle, including create, read, update, and delete (CRUD) operations.

### 3. Data Source Layer (`internal/datasources/`)
- Similar to the resource layer, each data source has its own package.
- Data sources are responsible for fetching and returning information from the Kasm API.

### 4. Testing Strategy
- The provider uses a combination of unit tests and integration tests to ensure functionality.
- Unit tests focus on individual components, while integration tests validate the interactions between components and the Kasm API.

### Design Principles
- **Separation of Concerns**: Each layer of the provider is responsible for a distinct aspect of functionality, promoting maintainability and clarity.
- **Modularity**: The architecture allows for easy addition of new features and resources without affecting existing functionality.

### Conclusion
This architecture documentation serves as a guide for understanding the structure and design of the Terraform Provider for Kasm. It is essential for maintaining and extending the provider as new features and updates are introduced.
