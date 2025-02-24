# Exec Resource

Manages command execution capabilities for Kasm sessions. This resource allows you to configure and control command execution within workspace sessions.

## Example Usage

### Basic Exec Configuration
```hcl
resource "kasm_exec" "example" {
  session_id = kasm_session.workspace.id
  enabled    = true
}
```

### Restricted Exec Configuration
```hcl
resource "kasm_exec" "restricted" {
  session_id = kasm_session.workspace.id
  enabled    = true

  settings = {
    allow_privileged = false
    allowed_commands = ["ls", "cat", "echo"]
    working_dir     = "/home/user"
    timeout        = 30
  }
}
```

## Argument Reference

* `session_id` - (Required) The ID of the session to configure exec for.
* `enabled` - (Required) Whether command execution is enabled.
* `settings` - (Optional) Execution settings:
  * `allow_privileged` - Allow privileged commands
  * `allowed_commands` - List of allowed commands
  * `working_dir` - Default working directory
  * `timeout` - Command timeout in seconds

## Attribute Reference

* `id` - The unique identifier for the exec configuration.
* `last_command` - Details of the last executed command.
* `execution_status` - Current execution status.

## Import

Import an exec configuration:

```bash
terraform import kasm_exec.example session_id:exec_config_id
```

## Notes

1. Security:
   - Command restrictions
   - Privilege controls
   - Environment isolation

2. Execution:
   - Command validation
   - Output capture
   - Error handling

3. Usage:
   - Automation scripts
   - System maintenance
   - Configuration management

4. Limitations:
   - Security constraints
   - Resource limitations
   - Command restrictions
