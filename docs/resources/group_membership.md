# kasm_group_membership (Resource)

Manages user membership in a Kasm group. This resource allows you to add and remove users from groups. Group membership is used to control access to workspace images and other resources.

## Example Usage

### Basic Group Membership
```hcl
# Create a group
resource "kasm_group" "example" {
  name        = "DevelopmentTeam"
  priority    = 1
  description = "Development team group"
}

# Create a user
resource "kasm_user" "example" {
  username     = "developer"
  password     = "Password123!"
  first_name   = "John"
  last_name    = "Doe"
  groups       = []  # Groups managed by kasm_group_membership

  lifecycle {
    ignore_changes = [
      groups,
      authorized_images,
    ]
  }
}

# Add the user to the group
resource "kasm_group_membership" "example" {
  group_id = kasm_group.example.id
  user_id  = kasm_user.example.id
}
```

### Group Membership with Image Authorization
```hcl
# Create a group
resource "kasm_group" "workspace_users" {
  name        = "WorkspaceUsers"
  priority    = 1
  description = "Users with access to workspace images"
}

# Get available images
data "kasm_images" "available" {
}

# Authorize an image for the group
resource "kasm_group_image" "workspace_access" {
  group_id = kasm_group.workspace_users.id
  image_id = data.kasm_images.available.images[0].id
}

# Create a user
resource "kasm_user" "developer" {
  username     = "developer"
  password     = "Password123!"
  first_name   = "John"
  last_name    = "Doe"
  groups       = []  # Groups managed by kasm_group_membership

  lifecycle {
    ignore_changes = [
      groups,
      authorized_images,
    ]
  }
}

# Add the user to the group
resource "kasm_group_membership" "developer_access" {
  group_id = kasm_group.workspace_users.id
  user_id  = kasm_user.developer.id
}

# Create a session (user now has access to the image)
resource "kasm_session" "developer_workspace" {
  depends_on = [kasm_group_membership.developer_access, kasm_group_image.workspace_access]
  image_id   = data.kasm_images.available.images[0].id
  user_id    = kasm_user.developer.id
}
```

## Argument Reference

* `group_id` - (Required) The ID of the group to add the user to.
* `user_id` - (Required) The ID of the user to add to the group.

## Attribute Reference

* `id` - A unique identifier for the group membership in the format "group_id:user_id".

## Import

Group memberships can be imported using the ID in the format "group_id:user_id":

```shell
terraform import kasm_group_membership.example "group_id:user_id"
```

## Notes

1. Image Authorization:
   - Group membership is used to control which images a user can access
   - Use `kasm_group_image` to authorize images for a group
   - Users in the group will have access to all authorized images
   - Required for creating sessions with specific images

2. User Configuration:
   - When using `kasm_group_membership`, set the user's `groups` field to an empty list
   - Add a lifecycle block to ignore changes to `groups` and `authorized_images`
   - This prevents conflicts between direct group assignment and group membership resources

3. Dependencies:
   - When creating sessions, use `depends_on` to ensure group membership and image authorization are set up first
   - This prevents errors from trying to use unauthorized images
