# Terraform Provider for Kasm Workspaces

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kasm/terraform-provider-kasm)](https://github.com/kasm/terraform-provider-kasm/releases)
[![Build Status](https://github.com/SiM22/terraform-provider-kasm/actions/workflows/test.yml/badge.svg)](https://github.com/SiM22/terraform-provider-kasm/actions/workflows/test.yml)

> ⚠️ **Work in Progress**: This provider is currently under active development and may not be fully stable. Some features might not work as expected, and breaking changes could occur. Use with caution in production environments.

This Terraform provider allows you to manage resources in [Kasm Workspaces](https://kasmweb.com/) through Infrastructure as Code.

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/kasm/kasm/latest/docs) and in the [Kasm Dev API Docs](docs/Kasm%20Dev%20API%20Docs).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Go >= 1.20 (to build the provider plugin)

## Quick Start

```hcl
provider "kasm" {
  base_url = "https://your.kasm.instance"
  api_key = "your-api-key"
  api_secret = "your-api-secret"
  insecure = true # For self-signed certificates
}

resource "kasm_user" "example" {
  username = "testuser"
  password = "TestPassword123!"
  first_name = "Test"
  last_name = "User"
}
```

## Features

- **User Management**
  - Create, update, and delete users
  - Manage user attributes and permissions
  - Group assignments and role management

- **Session Management**
  - Create and manage workspace sessions
  - Configure RDP/VNC access
  - Manage session permissions

- **Group Management**
  - Create and manage groups
  - Configure group policies
  - Manage group memberships

## Contributing

We welcome contributions! Please see our [Contribution Guidelines](docs/CONTRIBUTING.md) for details.
