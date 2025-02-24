# kasm_group_image (Resource)

Manages image authorization for a Kasm group. This resource allows you to authorize specific images for use by members of a group. Image authorization through groups is required before users can create sessions with specific images.

## Example Usage

### Basic Image Authorization
```hcl
# Get available images
data "kasm_images" "available" {
}

# Create a group
resource "kasm_group" "example" {
  name        = "DevelopmentTeam"
  priority    = 1
  description = "Development team group"
}

# Authorize an image for the group
resource "kasm_group_image" "example" {
  group_id = kasm_group.example.id
  image_id = data.kasm_images.available.images[0].id
}
```

### Complete Session Setup with Image Authorization
```hcl
# Get available images
data "kasm_images" "available" {
}

# Create a group
resource "kasm_group" "workspace_users" {
  name        = "WorkspaceUsers"
  priority    = 1
  description = "Users with access to workspace images"
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

# Add user to group
resource "kasm_group_membership" "developer_access" {
  group_id = kasm_group.workspace_users.id
  user_id  = kasm_user.developer.id
}

# Create a session
resource "kasm_session" "developer_workspace" {
  depends_on = [kasm_group_membership.developer_access, kasm_group_image.workspace_access]
  image_id   = data.kasm_images.available.images[0].id
  user_id    = kasm_user.developer.id
}
```

## Argument Reference

* `group_id` - (Required) The ID of the group to authorize the image for.
* `image_id` - (Required) The ID of the image to authorize. This can be obtained from the `kasm_images` data source.

## Attribute Reference

* `id` - The ID of the group image authorization.
* `group_image_id` - The ID of the group image authorization.
* `image_name` - The name of the authorized image.
* `group_name` - The name of the group.
* `image_friendly_name` - The friendly name of the image.
* `image_src` - The source path of the image.

## Import

Group image authorizations can be imported using the group_image_id:

```shell
terraform import kasm_group_image.example {group_image_id}
```

## Notes

1. Image Authorization System:
   - Users can only create sessions with images authorized for their groups
   - A user must be a member of a group (via `kasm_group_membership`) that has the image authorized (via `kasm_group_image`)
   - Use the `kasm_images` data source to get available image IDs
   - Session creation will fail if the user's groups don't have the required image authorization

2. Dependencies:
   - When creating sessions, use `depends_on` to ensure image authorization is set up first
   - Both group membership and image authorization must be in place before creating sessions

3. Multiple Groups:
   - Users can be members of multiple groups
   - Users have access to all images authorized for any of their groups
   - Image authorization is additive across groups
