# Getting Started with the Kasm Provider

This guide will help you get started with the Kasm Terraform provider, enabling you to manage your Kasm Workspaces effectively.

## Prerequisites

Before you begin, ensure you have the following:
- A Kasm Workspaces account
- API credentials (API key and secret)
- Terraform installed (version >= 1.0)

## Installation

To install the Kasm provider, add the following configuration to your Terraform configuration file:

```hcl
tf
terraform {
  required_providers {
    kasm = {
      source = "kasmtech/kasm"
      version = ">= 1.0.0"
    }
  }
}
```

## Provider Configuration

Configure the provider in your Terraform configuration:

```hcl
provider "kasm" {
  base_url   = "https://kasm.example.com"
  api_key    = "your-api-key"
  api_secret = "your-api-secret"
  insecure   = false  # Optional: Skip TLS verification
}
```

## Example Usage

Here’s a simple example of how to create a user in Kasm:

```hcl
resource "kasm_user" "example" {
  username     = "testuser"
  password     = "TestPassword123!"
  first_name   = "Test"
  last_name    = "User"
}
```

## Next Steps

Once you have the provider configured, you can start managing your Kasm resources. Refer to the following guides for more information:
- [Managing Users](managing_users.md)
- [Managing Groups](managing_groups.md)
- [Session Management](session_management.md)

For further assistance, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs) or the [TerraformProviderDocumentation](TerraformProviderDocumentation).
