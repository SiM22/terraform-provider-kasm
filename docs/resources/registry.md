# Registry Resource

Manages Docker registry configurations in Kasm. This resource allows you to configure and manage Docker registry access for workspace images.

## Example Usage

### Basic Registry
```hcl
resource "kasm_registry" "dockerhub" {
  name     = "DockerHub"
  server   = "registry.hub.docker.com"
  username = "dockerhub_user"
  password = "dockerhub_password"
}
```

### Private Registry with Custom Configuration
```hcl
resource "kasm_registry" "private" {
  name        = "Private Registry"
  server      = "registry.company.com"
  username    = "registry_user"
  password    = "registry_password"
  insecure    = false
  channel     = "stable"
}
```

## Argument Reference

* `name` - (Required) The name of the registry configuration.
* `server` - (Required) The registry server URL.
* `username` - (Required) Username for registry authentication.
* `password` - (Required) Password for registry authentication.
* `insecure` - (Optional) Whether to skip TLS verification. Defaults to false.
* `channel` - (Optional) The channel to use for images. Defaults to "stable".

## Attribute Reference

* `id` - The unique identifier for the registry configuration.

## Import

Import a registry configuration:

```bash
terraform import kasm_registry.example registry_id_here
```

## Notes

1. Authentication:
   - Credentials are stored securely
   - Support for token-based authentication
   - Automatic credential rotation not supported

2. Security:
   - TLS verification enabled by default
   - Insecure mode available for testing
   - Credentials encrypted at rest

3. Channels:
   - Stable: Production-ready images
   - Testing: Pre-release images
   - Development: Latest builds

4. Usage:
   - Used for pulling workspace images
   - Supports private repositories
   - Multiple registries can be configured
