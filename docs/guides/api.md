# API Usage Guide

This guide explains how the Kasm provider interacts with the Kasm API and how to troubleshoot API-related issues.

## Authentication

### API Keys
```hcl
provider "kasm" {
  base_url   = "https://kasm.example.com"
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
}
```

### Environment Variables
```bash
export KASM_BASE_URL="https://kasm.example.com"
export KASM_API_KEY="your-api-key"
export KASM_API_SECRET="your-api-secret"
```

## API Client Organization

The provider's API client (`internal/client`) is organized into focused operation files:

### Core Operations
- `client.go`: Core client implementation and configuration
- `http.go`: HTTP request handling and response processing
- `errors.go`: Error definitions and handling
- `ratelimit.go`: Rate limiting implementation
- `backoff.go`: Retry and backoff logic

### Resource Operations
- `user_ops.go`: User management operations
- `group_ops.go`: Group management operations
- `image_ops.go`: Image management operations
- `registry_ops.go`: Registry management operations
- `cast_ops.go`: Cast configuration operations
- `license_ops.go`: License management operations
- `login_ops.go`: Login URL generation operations
- `session_ops.go`: Session management operations

### Type Definitions
- `user_types.go`: User-related types
- `group_types.go`: Group-related types
- `image_types.go`: Image-related types
- `registry_types.go`: Registry-related types
- `cast_types.go`: Cast-related types
- `license_types.go`: License-related types
- `login_types.go`: Login-related types
- `session_types.go`: Session-related types

## API Endpoints

### User Management
```
POST /api/public/create_user
POST /api/public/get_user
POST /api/public/update_user
POST /api/public/delete_user
```

### Group Management
```
POST /api/public/create_group
POST /api/public/get_groups
POST /api/public/update_group
POST /api/public/delete_group
```

### Session Management
```
POST /api/public/create_session
POST /api/public/get_session
POST /api/public/update_session
POST /api/public/delete_session
```

### Login Management
```
POST /api/public/get_login      # Generate login URL for user
```

## Rate Limiting and Retries

The provider implements sophisticated request handling:

### Rate Limiting (`ratelimit.go`)
1. Default: 100 requests per second
2. Configurable through settings
3. Automatic rate limiting
4. Queue management

### Retries (`backoff.go`)
1. Exponential backoff
2. Configurable retry limits
3. Error-specific retry logic
4. Jitter implementation

## Error Handling

### Error Types (`errors.go`)
1. Authentication Errors (401)
   - Invalid credentials
   - Expired tokens
   - Missing authentication

2. Permission Errors (403)
   - Insufficient permissions
   - Resource access denied
   - Policy violations

3. Not Found Errors (404)
   - Resource not found
   - Invalid endpoints
   - Deleted resources

4. Rate Limit Errors (429)
   - Too many requests
   - Quota exceeded
   - Throttling applied

5. Server Errors (500)
   - Internal server errors
   - Service unavailable
   - Database errors

### Error Response Format
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {}
  }
}
```

## Best Practices

1. API Usage
- Use appropriate client operations
- Leverage built-in rate limiting
- Utilize retry mechanisms
- Validate responses
- Handle errors appropriately

2. Error Handling
- Use defined error types
- Implement proper backoff
- Log detailed errors
- Handle retries correctly
- Check error contexts

3. Performance
- Minimize API calls
- Use bulk operations where available
- Cache responses when appropriate
- Handle timeouts properly
- Monitor rate limits

4. Security
- Secure credential storage
- Use TLS verification
- Implement key rotation
- Audit API access
- Follow least privilege

## Debugging

Enable debug logging to see detailed API interaction:
```bash
export KASM_DEBUG=1
export KASM_LOG_LEVEL=debug
```

This will show:
1. API requests and payloads
2. Response data
3. Rate limiting status
4. Retry attempts
5. Error details
