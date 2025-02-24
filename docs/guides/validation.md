# Validation Guide

This guide explains the validation rules and constraints applied to resources in the Kasm provider.

## Resource Validation

### User Resource

1. Username Validation
```hcl
resource "kasm_user" "example" {
  username = "user@example.com"  # Must be valid email format
}
```

2. Password Requirements
- Minimum length: 8 characters
- Must contain at least one number
- Must contain at least one special character

### Group Resource

1. Name Validation
- Must be unique
- Cannot contain special characters
- Length between 3 and 64 characters

2. Priority Validation
- Must be positive integer
- Range: 1-100

### Session Resource

1. Image ID Validation
- Must be valid UUID format
- Must reference existing image

2. User ID Validation
- Must be valid UUID format
- Must reference existing user

## Common Validations

The provider includes a centralized validation package (`internal/validators`) that implements common validation functions used across resources:

### String Validators
- Email format validation
- UUID format validation
- Name format validation (alphanumeric, dashes, underscores)
- Length constraints
- Character set restrictions

### Numeric Validators
- Range validation
- Integer validation
- Positive number validation
- Port number validation

### Resource Validators
- Resource name validation
- Resource reference validation
- Resource state validation

### Configuration Validators
- URL format validation
- Path format validation
- Version format validation

## Using Validators

Resources use the common validators to ensure consistency:

```hcl
# Email validation is consistent across resources
resource "kasm_user" "example" {
  username = "user@example.com"  # Uses common email validator
}

# Name validation is consistent across resources
resource "kasm_group" "example" {
  name = "developers"  # Uses common name validator
}
```

## Error Messages

Validation errors include:
1. Detailed error message
2. Field reference
3. Expected format/value
4. Current invalid value

Example error:
```
Error: Invalid username format
Field: username
Expected: valid email address
Got: "invalid@email@example.com"
```

## Best Practices

1. Resource Naming
- Use descriptive names
- Follow naming conventions
- Avoid special characters
- Use validated formats

2. Configuration Values
- Use appropriate types
- Follow format requirements
- Consider dependencies
- Validate against schemas

3. Resource References
- Verify existence
- Use correct reference format
- Handle dependencies
- Validate UUIDs

4. State Management
- Valid state transitions
- Proper cleanup
- Resource dependencies
- Validate state changes

## Extending Validation

When adding new resources:
1. Use existing validators from the validators package
2. Add new validators to the package if needed
3. Follow consistent error message formats
4. Document validation rules

## Common Validation Patterns

### Required Fields
```hcl
resource "kasm_example" "test" {
  required_field = "value"  # Will fail validation if empty
}
```

### Optional Fields with Validation
```hcl
resource "kasm_example" "test" {
  optional_field = "valid-name"  # If provided, must match format
}
```

### Dependent Field Validation
```hcl
resource "kasm_example" "test" {
  feature_enabled = true
  feature_config  = "required-when-enabled"  # Validated based on feature_enabled
}
```

## Testing Validation

The provider includes comprehensive validation tests:
1. Unit tests for individual validators
2. Integration tests for validation chains
3. Acceptance tests for full resource validation
4. Error message format tests
