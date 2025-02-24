# Image Resource

Manages workspace images in Kasm. This resource allows you to configure and manage the images available for workspace sessions.

## Example Usage

### Basic Image
```hcl
resource "kasm_image" "ubuntu" {
  name           = "Ubuntu Workspace"
  image_src      = "ubuntu:focal"
  friendly_name  = "Ubuntu 20.04"
  description    = "Basic Ubuntu workspace"
  memory         = 2048
  cores          = 2
  enabled        = true
}
```

### Custom Image with Registry
```hcl
resource "kasm_image" "development" {
  name              = "Development Environment"
  image_src         = "registry.company.com/dev-workspace:latest"
  friendly_name     = "Development Workspace"
  description       = "Full development environment with tools"
  docker_registry   = "registry.company.com"
  memory            = 4096
  cores             = 4
  enabled           = true
}
```

## Argument Reference

* `name` - (Required) The name of the image.
* `image_src` - (Required) The source of the Docker image.
* `friendly_name` - (Required) A user-friendly name for the image.
* `description` - (Optional) A description of the image.
* `docker_registry` - (Optional) The Docker registry where the image is hosted.
* `memory` - (Optional) The amount of memory in MB to allocate to containers using this image.
* `cores` - (Optional) The number of CPU cores to allocate to containers using this image.
* `enabled` - (Optional) Whether the image is available for use. Defaults to true.
* `uncompressed_size_mb` - (Optional) The uncompressed size of the image in MB.
* `image_type` - (Optional) The type of the image.
* `run_config` - (Optional) Configuration for running the container.
* `exec_config` - (Optional) Configuration for executing commands in the container.
* `volume_mapping` - (Optional) Volume mappings for the container.
* `restrict_to_network` - (Optional) Whether to restrict the image to a specific network.
* `restrict_to_server` - (Optional) Whether to restrict the image to a specific server.
* `restrict_to_zone` - (Optional) Whether to restrict the image to a specific zone.
* `server_id` - (Optional) The ID of the server to restrict the image to.
* `zone_id` - (Optional) The ID of the zone to restrict the image to.
* `network_name` - (Optional) The name of the network to restrict the image to.

## Attribute Reference

* `id` - The unique identifier for the image.
* `image_id` - The ID used to reference this image in other resources.
* `available` - Whether the image is currently available for use.

## Import

Images can be imported using their ID:

```bash
terraform import kasm_image.example <image_id>
```

## Notes

1. Image Sources:
   - Public Docker Hub images
   - Private registry images
   - Custom built images

2. Registry Integration:
   - Links to configured registries
   - Authentication handled automatically
   - Support for private repositories

3. Usage:
   - Referenced in session configurations
   - Used as base for workspaces
   - Can be enabled/disabled as needed

4. Management:
   - Version control through image tags
   - Automatic updates possible
   - State tracking for running sessions

5. Resource Operations:
   - Create: Creates a new workspace image
   - Read: Retrieves current image configuration
   - Update: Modifies existing image settings
   - Delete: Removes the image from Kasm
