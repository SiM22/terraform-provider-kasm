# Managing Sessions with the Kasm Provider

This guide provides instructions on how to manage sessions in Kasm Workspaces using the Terraform provider.

## Creating a Session

To create a session in Kasm, use the following resource configuration:

```hcl
resource "kasm_session" "new_session" {
  user_id    = kasm_user.example.id
  image_id   = kasm_image.example.id
  share      = true
  rdp_enabled = true
}
```

## Updating a Session

To update an existing session, you can modify the attributes in the resource configuration:

```hcl
resource "kasm_session" "existing_session" {
  user_id    = kasm_user.example.id
  image_id   = kasm_image.example.id
  share      = false
  rdp_enabled = true
}
```

## Deleting a Session

To delete a session, you can use the following configuration:

```hcl
resource "kasm_session" "delete_session" {
  user_id = kasm_user.example.id
}
```

## Session Attributes

When creating or updating a session, you can specify the following attributes:
- `user_id`: The ID of the user for the session (required).
- `image_id`: The ID of the image to use for the session (required).
- `share`: Whether to allow sharing of the session (optional).
- `rdp_enabled`: Whether RDP access is enabled (optional).

## Example Usage

Here’s a complete example of managing a session:

```hcl
resource "kasm_session" "example" {
  user_id    = kasm_user.example.id
  image_id   = kasm_image.example.id
  share      = true
  rdp_enabled = true
}

resource "kasm_session" "updated_session" {
  user_id    = kasm_user.example.id
  image_id   = kasm_image.example.id
  share      = false
  rdp_enabled = true
}

resource "kasm_session" "deleted_session" {
  user_id = kasm_user.example.id
}
```

## Next Steps

After managing sessions, you can refer to the following guides for additional management tasks:
- [Managing Users](managing_users.md)
- [Managing Groups](managing_groups.md)

For more details, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs).
