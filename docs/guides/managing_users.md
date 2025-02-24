# Managing Users with the Kasm Provider

This guide provides instructions on how to manage users in Kasm Workspaces using the Terraform provider.

## Creating a User

To create a user in Kasm, use the following resource configuration:

```hcl
resource "kasm_user" "new_user" {
  username     = "newuser"
  password     = "SecurePassword123!"
  first_name   = "New"
  last_name    = "User"
}
```

## Updating a User

To update an existing user, you can modify the attributes in the resource configuration:

```hcl
resource "kasm_user" "existing_user" {
  username     = "existinguser"
  password     = "NewSecurePassword456!"
  first_name   = "Updated"
  last_name    = "User"
}
```

## Deleting a User

To delete a user, you can use the following configuration:

```hcl
resource "kasm_user" "delete_user" {
  username = "user_to_delete"
}
```

## User Attributes

When creating or updating a user, you can specify the following attributes:
- `username`: The username for the user (required).
- `password`: The password for the user (required).
- `first_name`: The first name of the user (optional).
- `last_name`: The last name of the user (optional).

## Example Usage

Here’s a complete example of managing a user:

```hcl
resource "kasm_user" "example" {
  username     = "testuser"
  password     = "TestPassword123!"
  first_name   = "Test"
  last_name    = "User"
}

resource "kasm_user" "updated_user" {
  username     = "testuser"
  password     = "NewPassword456!"
  first_name   = "Updated"
  last_name    = "User"
}

resource "kasm_user" "deleted_user" {
  username = "user_to_delete"
}
```

## Next Steps

After managing users, you can refer to the following guides for additional management tasks:
- [Managing Groups](managing_groups.md)
- [Session Management](session_management.md)

For more details, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs).
