# User Resource

Manages a Kasm user. This resource allows you to create, update, and delete users in your Kasm workspace.

## Example Usage

### Basic User
```hcl
resource "kasm_user" "example" {
  username     = "john.doe@example.com"
  password     = "SecurePass123!"
  first_name   = "John"
  last_name    = "Doe"
}
```

### User with Group Assignment
```hcl
resource "kasm_user" "developer" {
  username     = "dev@company.com"
  password     = "SecurePass123!"
  first_name   = "John"
  last_name    = "Developer"
  organization = "Development"
  groups       = [kasm_group.developers.name]
}
```

### Disabled User
```hcl
resource "kasm_user" "disabled_user" {
  username  = "inactive@company.com"
  password  = "SecurePass123!"
  disabled  = true
}
```

## Argument Reference

* `username` - (Required) The email address of the user.
* `password` - (Required) The user's password. Must meet complexity requirements.
* `first_name` - (Optional) The user's first name.
* `last_name` - (Optional) The user's last name.
* `organization` - (Optional) The organization the user belongs to.
* `phone` - (Optional) The user's phone number.
* `locked` - (Optional) Whether the user account is locked. Defaults to false.
* `disabled` - (Optional) Whether the user account is disabled. Defaults to false.
* `groups` - (Optional) List of group names to assign the user to.

## Attribute Reference

* `id` - The unique identifier for the user.

## Import

Import a user using either their ID or username:

```bash
# Import by ID
terraform import kasm_user.example user_id_here

# Import by username
terraform import kasm_user.example username:john.doe@example.com
```

## Notes

1. Password Requirements:
   - Must contain at least one number
   - Should be sufficiently complex for security

2. Group Management:
   - Groups are managed by adding/removing users
   - The "All Users" group is automatically assigned
   - Group names must be valid and existing

3. User States:
   - Users can be locked (temporarily prevented from logging in)
   - Users can be disabled (permanently deactivated)
   - State changes trigger appropriate cleanup actions

4. Organization:
   - Organization field is used for grouping and filtering
   - Does not affect permissions or access control
