# Kasm Provider

The Kasm provider allows you to manage resources in [Kasm Workspaces](https://kasmweb.com/), enabling you to interact with various components of your Kasm environment, including users, groups, sessions, and configurations.

## Example Usage

```hcl
terraform {
  required_providers {
    kasm = {
      source = "kasmtech/kasm"
      version = ">= 1.0.0"
    }
  }
}

provider "kasm" {
  base_url   = "https://kasm.example.com"
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
}
```

## Authentication

The provider requires API credentials, which can be provided in multiple ways:

### Environment Variables
```bash
export KASM_BASE_URL="https://kasm.example.com"
export KASM_API_KEY="your-api-key"
export KASM_API_SECRET="your-api-secret"
```

### Provider Configuration
```hcl
provider "kasm" {
  base_url   = "https://kasm.example.com"
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
  insecure   = false  # Optional: Skip TLS verification
}
```

## Argument Reference

- `base_url` - (Required) The base URL of your Kasm instance. Can also be provided via `KASM_BASE_URL` environment variable.
- `api_key` - (Required) API key for authentication. Can also be provided via `KASM_API_KEY` environment variable.
- `api_secret` - (Required) API secret for authentication. Can also be provided via `KASM_API_SECRET` environment variable.
- `insecure` - (Optional) Skip TLS verification. Defaults to false.

## Resource Types

- `kasm_user` - Manage Kasm users.
- `kasm_group` - Manage Kasm groups.
- `kasm_session` - Manage Kasm sessions.
- `kasm_login` - Generates login URLs for users
- `kasm_rdp` - Configures RDP access
- `kasm_screenshot` - Manages session screenshots
- `kasm_stats` - Handles session statistics
- `kasm_keepalive` - Manages session keepalive settings
- `kasm_exec` - Manages command execution
- `kasm_cast` - Manages session casting configurations
- `kasm_registry` - Manages Docker registry configurations
- `kasm_image` - Manages workspace images
- `kasm_license` - Manages Kasm license activation

## Data Sources

- `kasm_images` - Query available workspace images
- `kasm_registries` - Query available registries
- `kasm_workspace` - Query workspace information

## Guides

- [Getting Started](guides/getting_started.md)
- [Managing Users](guides/managing_users.md)
- [Session Management](guides/session_management.md)

## Contributing

We welcome contributions! Please see our [Contribution Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
