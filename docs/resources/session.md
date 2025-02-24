# Session Resource

Manages a Kasm workspace session. This resource allows you to create, configure, and manage Kasm workspace sessions.

## Example Usage

### Basic Session with Required Group Authorization
```hcl
# Create a group
resource "kasm_group" "example" {
  name        = "workspace-users"
  description = "Group for workspace access"
  priority    = 100
}

# Get available images
data "kasm_images" "available" {
}

# Authorize the image for the group
resource "kasm_group_image" "example" {
  group_id = kasm_group.example.id
  image_id = data.kasm_images.available.images[0].id
}

# Create a user
resource "kasm_user" "example" {
  username     = "example-user"
  password     = "securepassword123!"
  first_name   = "Example"
  last_name    = "User"
  organization = "Example Org"
  locked       = false
  disabled     = false
  groups       = []  # Groups managed by kasm_group_membership

  lifecycle {
    ignore_changes = [
      groups,
      authorized_images,
    ]
  }
}

# Add user to group
resource "kasm_group_membership" "example" {
  group_id = kasm_group.example.id
  user_id  = kasm_user.example.id
}

# Create session
resource "kasm_session" "example" {
  depends_on = [kasm_group_membership.example, kasm_group_image.example]
  image_id   = data.kasm_images.available.images[0].id
  user_id    = kasm_user.example.id
}
```

### Session with RDP Access
```hcl
resource "kasm_session" "rdp_enabled" {
  depends_on   = [kasm_group_membership.example, kasm_group_image.example]
  image_id     = data.kasm_images.available.images[0].id
  user_id      = kasm_user.example.id
  rdp_enabled  = true
}
```

### Shared Session with Stats
```hcl
resource "kasm_session" "shared" {
  depends_on    = [kasm_group_membership.example, kasm_group_image.example]
  image_id      = data.kasm_images.available.images[0].id
  user_id       = kasm_user.example.id
  share         = true
  enable_stats  = true
}
```

## Argument Reference

* `image_id` - (Required) The ID of the workspace image to use for the session. The user must be authorized to use this image through group membership.
* `user_id` - (Required) The ID of the user to create the session for.
* `share` - (Optional) Whether to enable session sharing. Defaults to false.
* `enable_sharing` - (Optional) Whether to enable sharing features. Automatically set to true if share is true.
* `rdp_enabled` - (Optional) Whether to enable RDP for the session. Defaults to false.
* `enable_stats` - (Optional) Whether to enable session statistics. Defaults to false.
* `allow_exec` - (Optional) Whether to allow command execution in the session. Defaults to false.

## Attribute Reference

* `id` - The ID of the Kasm session.
* `share_id` - The share ID for the session when sharing is enabled.
* `rdp_connection_file` - The RDP connection file content when RDP is enabled.
* `operational_status` - The current status of the session.

## Import

Import a session using the format `user_id:session_id`:

```bash
terraform import kasm_session.example user_123:session_456
```

## Notes

1. Image Authorization:
   - Users must be authorized to use an image through group membership
   - Use `kasm_group_image` to authorize images for groups
   - Use `kasm_group_membership` to add users to groups
   - The session creation will fail if the user is not authorized for the image

2. Session States:
   - Creating: Initial session setup
   - Running: Active session
   - Stopped: Session terminated
   - Failed: Session creation/operation failed

3. RDP Access:
   - RDP configuration is generated when enabled
   - Connection file contains secure credentials
   - RDP port is automatically assigned

4. Session Sharing:
   - Share ID is generated when sharing is enabled
   - Share ID required for other users to join
   - Sharing can be toggled after creation
   - When `share` is set to true, `enable_sharing` is automatically set to true

5. Statistics:
   - Performance metrics when enabled
   - Resource usage tracking
   - Session analytics
