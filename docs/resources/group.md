# Group Resource

Manages a Kasm group. Groups are used to organize users and manage permissions in your Kasm workspace.

## Example Usage

### Basic Group
```hcl
resource "kasm_group" "developers" {
  name        = "Developers"
  description = "Development team group"
  priority    = 50
}
```

### System Group
```hcl
resource "kasm_group" "admins" {
  name        = "Administrators"
  description = "System administrators group"
  priority    = 100
  is_system   = true
}
```

## Argument Reference

* `name` - (Required) The name of the group.
* `description` - (Optional) A description of the group's purpose.
* `priority` - (Required) The group's priority level. Higher numbers indicate higher priority.
* `is_system` - (Optional) Whether this is a system group. Defaults to false.

## Attribute Reference

* `group_id` - The unique identifier for the group.

## Import

Import a group using either the ID or name:

```bash
# Import by ID
terraform import kasm_group.example group_id_here

# Import by name
terraform import kasm_group.example name:group_name_here
```

## Notes

1. Priority Levels:
   - Higher priority groups take precedence in permission conflicts
   - Priority must be a positive integer
   - Consider using standardized priority ranges (e.g., 1-100)

2. System Groups:
   - System groups have special privileges
   - Cannot be deleted through normal means
   - Used for core system functionality

3. Group Names:
   - Must be unique across the workspace
   - Case-sensitive
   - Used as identifiers in user assignments

4. Usage:
   - Groups are referenced in user configurations
   - Can be used for access control
   - Support organizational structure
