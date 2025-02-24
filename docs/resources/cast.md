# Cast Resource

Manages session casting configurations in Kasm. This resource allows you to configure how workspace sessions can be shared and viewed by other users.

## Example Usage

### Basic Cast Configuration
```hcl
resource "kasm_cast" "example" {
  name           = "Development Cast"
  image_id       = data.kasm_images.ubuntu.id
  enable_sharing = true
}
```

### Advanced Cast Configuration
```hcl
resource "kasm_cast" "restricted" {
  name           = "Restricted Cast"
  image_id       = data.kasm_images.windows.id
  enable_sharing = true
  key            = "dev_cast"

  client_settings = {
    allow_kasm_audio        = true
    idle_disconnect        = 3600
    allow_kasm_clipboard   = true
    allow_kasm_downloads   = false
    allow_kasm_uploads     = false
    allow_kasm_microphone  = false
  }
}
```

## Argument Reference

* `name` - (Required) The name of the cast configuration.
* `image_id` - (Required) The ID of the workspace image to use.
* `enable_sharing` - (Required) Whether to enable session sharing.
* `key` - (Optional) A unique key for the cast configuration.
* `client_settings` - (Optional) A map of client settings:
  * `allow_kasm_audio` - Allow audio streaming
  * `idle_disconnect` - Idle timeout in seconds
  * `allow_kasm_clipboard` - Allow clipboard sharing
  * `allow_kasm_downloads` - Allow file downloads
  * `allow_kasm_uploads` - Allow file uploads
  * `allow_kasm_microphone` - Allow microphone access

## Attribute Reference

* `id` - The unique identifier for the cast configuration.

## Import

Import a cast configuration:

```bash
terraform import kasm_cast.example cast_id_here
```

## Notes

1. Sharing Settings:
   - Control what features are available during sharing
   - Configure security restrictions
   - Manage resource access

2. Client Settings:
   - Audio and video controls
   - File transfer permissions
   - Input device access
   - Session timeouts

3. Security:
   - Feature-level access control
   - Session isolation
   - Resource restrictions

4. Usage:
   - Used for configuring shared sessions
   - Applies to all sessions using this configuration
   - Can be updated while sessions are active
