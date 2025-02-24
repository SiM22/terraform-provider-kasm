# Managing Groups with the Kasm Provider

This guide provides instructions on how to manage groups in Kasm Workspaces using the Terraform provider.

## Creating a Group

To create a group in Kasm, use the following resource configuration:

```hcl
resource "kasm_group" "new_group" {
  name        = "Development Team"
  priority    = 50
  description = "Group for development team members"
}
```

## Updating a Group

To update an existing group, you can modify the attributes in the resource configuration:

```hcl
resource "kasm_group" "existing_group" {
  name        = "Development Team"
  priority    = 40
  description = "Updated description for the development team"
}
```

## Deleting a Group

To delete a group, you can use the following configuration:

```hcl
resource "kasm_group" "delete_group" {
  name = "Group to Delete"
}
```

## Group Attributes

When creating or updating a group, you can specify the following attributes:
- `name`: The name of the group (required).
- `priority`: The priority level for the group (optional).
- `description`: A description of the group (optional).

## Example Usage

Here’s a complete example of managing a group:

```hcl
resource "kasm_group" "example" {
  name        = "Development Team"
  priority    = 50
  description = "Group for development team members"
}

resource "kasm_group" "updated_group" {
  name        = "Development Team"
  priority    = 40
  description = "Updated description for the development team"
}

resource "kasm_group" "deleted_group" {
  name = "Group to Delete"
}
```

## Next Steps

After managing groups, you can refer to the following guides for additional management tasks:
- [Managing Users](managing_users.md)
- [Session Management](session_management.md)

For more details, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs).
