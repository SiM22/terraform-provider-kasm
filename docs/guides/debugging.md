# Debugging Guide

This guide explains how to debug issues with the Kasm Terraform provider.

## Debug Mode

The provider includes built-in debug functionality that can be enabled through environment variables.

### Enabling Debug Mode

```bash
export KASM_DEBUG=1
export KASM_LOG_LEVEL=debug
```

### Debug Output

When debug mode is enabled, the provider will output detailed information about:
- API requests and responses
- Resource operations
- State changes
- Error details

## Common Issues

### Authentication Issues

If you encounter authentication errors:
1. Verify your API credentials
2. Check the base URL configuration
3. Ensure TLS settings are correct
4. Review API logs with debug mode enabled

### Resource Creation Failures

When resources fail to create:
1. Enable debug mode
2. Check the API response
3. Verify resource configurations
4. Review dependent resources

### State Management

For state-related issues:
1. Use `terraform refresh` to sync state
2. Check resource IDs and references
3. Review import statements
4. Enable debug logging

## Logging

### Log Levels

The provider supports multiple log levels:
- error: Only error messages
- warn: Warnings and errors
- info: General information
- debug: Detailed debugging information

### Log Format

Debug logs include:
```
YYYY-MM-DD HH:MM:SS [LEVEL] Message
  Key1: Value1
  Key2: Value2
```

## Troubleshooting Tools

1. terraform console
```bash
$ terraform console
> var.kasm_base_url
```

2. terraform plan with debug
```bash
TF_LOG=DEBUG terraform plan
```

3. State inspection
```bash
terraform show
```

## Getting Help

1. Enable debug mode
2. Collect relevant logs
3. Create an issue with:
   - Debug output
   - Resource configurations
   - Error messages
   - Steps to reproduce
