# RDP Resource

Manages RDP (Remote Desktop Protocol) configurations for Kasm sessions. This resource allows you to configure and manage RDP access to workspace sessions.

## Example Usage

### Basic RDP Configuration
```hcl
resource "kasm_rdp" "example" {
  session_id = kasm_session.windows.id
  enabled    = true
}
```

### Advanced RDP Configuration
```hcl
resource "kasm_rdp" "custom" {
  session_id = kasm_session.windows.id
  enabled    = true

  connection_settings = {
    port           = 3389
    authentication = "nla"
    quality        = "high"
    compression    = true
  }
}
```

## Argument Reference

* `session_id` - (Required) The ID of the session to enable RDP for.
* `enabled` - (Required) Whether RDP access is enabled.
* `connection_settings` - (Optional) RDP connection settings:
  * `port` - RDP port number
  * `authentication` - Authentication level
  * `quality` - Connection quality setting
  * `compression` - Enable compression

## Attribute Reference

* `id` - The unique identifier for the RDP configuration.
* `connection_file` - The RDP connection file content.
* `connection_url` - The URL for RDP access.

## Import

Import an RDP configuration:

```bash
terraform import kasm_rdp.example session_id:rdp_config_id
```

## Notes

1. Security:
   - Network Level Authentication (NLA) support
   - Encrypted connections
   - Session-specific credentials

2. Connection Settings:
   - Configurable port numbers
   - Quality/performance options
   - Compression settings

3. Access:
   - Connection file generation
   - URL-based access
   - Credential management

4. Requirements:
   - Windows-based workspace image
   - Proper network configuration
   - Supported authentication methods
