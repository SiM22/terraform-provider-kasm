# Keepalive Resource

Manages session keepalive settings for Kasm sessions. This resource allows you to configure automatic session maintenance to prevent timeouts.

## Example Usage

### Basic Keepalive Configuration
```hcl
resource "kasm_keepalive" "example" {
  session_id = kasm_session.workspace.id
  enabled    = true
}
```

### Custom Keepalive Configuration
```hcl
resource "kasm_keepalive" "custom" {
  session_id = kasm_session.workspace.id
  enabled    = true

  settings = {
    interval     = 300
    max_duration = 86400
    idle_timeout = 3600
  }
}
```

## Argument Reference

* `session_id` - (Required) The ID of the session to configure keepalive for.
* `enabled` - (Required) Whether keepalive is enabled.
* `settings` - (Optional) Keepalive settings:
  * `interval` - Keepalive interval in seconds
  * `max_duration` - Maximum session duration in seconds
  * `idle_timeout` - Idle timeout in seconds

## Attribute Reference

* `id` - The unique identifier for the keepalive configuration.
* `last_keepalive` - Timestamp of last keepalive signal.
* `session_duration` - Current session duration.

## Import

Import a keepalive configuration:

```bash
terraform import kasm_keepalive.example session_id:keepalive_config_id
```

## Notes

1. Timing:
   - Regular interval checks
   - Maximum session limits
   - Idle detection

2. Resource Management:
   - Automatic cleanup
   - Resource release
   - Session termination

3. Usage:
   - Long-running sessions
   - Automated workflows
   - Continuous operations

4. Limitations:
   - Network requirements
   - Resource constraints
   - Policy compliance
