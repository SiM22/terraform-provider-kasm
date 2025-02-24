# Images Data Source

Provides information about available workspace images in Kasm. This data source is commonly used to find available images for group authorization and session creation.

## Example Usage

### Basic Query
```hcl
data "kasm_images" "all" {}

output "available_images" {
  value = data.kasm_images.all.images
}
```

### Image Authorization Workflow
```hcl
# Get available images
data "kasm_images" "available" {}

# Create a group
resource "kasm_group" "workspace_users" {
  name        = "WorkspaceUsers"
  priority    = 1
  description = "Users with access to workspace images"
}

# Authorize first available image for the group
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
  groups       = []

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

# Create a session with the authorized image
resource "kasm_session" "developer_workspace" {
  depends_on = [kasm_group_membership.developer_access, kasm_group_image.workspace_access]
  image_id   = data.kasm_images.available.images[0].id
  user_id    = kasm_user.developer.id
}
```

### Output Image Details
```hcl
data "kasm_images" "available" {}

output "image_details" {
  value = [
    for img in data.kasm_images.available.images : {
      id            = img.id
      name          = img.name
      friendly_name = img.friendly_name
      description   = img.description
      cores         = img.cores
      memory        = img.memory
    }
  ]
}
```

## Argument Reference

The data source doesn't require any arguments, but accepts the following:

* `filter` - (Optional) Filter criteria for images:
  * `name` - Filter by image name
  * `enabled` - Filter by enabled status

## Attribute Reference

* `images` - List of images matching the filter criteria. Each image contains:
  * `id` - Image identifier
  * `name` - Image name (e.g., "kasmweb/terminal:1.16.0")
  * `friendly_name` - User-friendly name (e.g., "Terminal")
  * `description` - Image description
  * `cores` - Number of CPU cores
  * `memory` - Memory allocation in bytes
  * `cpu_allocation_method` - CPU allocation method
  * `image_src` - Source path of the image
  * `enabled` - Whether the image is enabled
  * `available` - Whether the image is available for use

## Notes

1. Image Authorization:
   - Images must be authorized for groups before users can create sessions
   - Use with `kasm_group_image` to authorize images for groups
   - Users must be members of groups with authorized images

2. Session Creation:
   - Use image IDs from this data source when creating sessions
   - Ensure users have proper authorization through group membership
   - Session creation will fail if the user's groups don't have the image authorized

3. Best Practices:
   - Query available images before creating group authorizations
   - Use image attributes to select appropriate images for workloads
   - Consider memory and CPU requirements when selecting images

4. Integration:
   - Core component of the image authorization system
   - Used in conjunction with groups, users, and sessions
   - Essential for workspace provisioning and management
