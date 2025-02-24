# Screenshot Resource

Manages screenshot capabilities for Kasm sessions. This resource allows you to configure and manage screenshot functionality for workspace sessions.

## Example Usage

### Basic Screenshot Configuration
```hcl
resource "kasm_screenshot" "example" {
  session_id = kasm_session.workspace.id
  enabled    = true
}
```

### Custom Screenshot Configuration
```hcl
resource "kasm_screenshot" "custom" {
  session_id = kasm_session.workspace.id
  enabled    = true

  settings = {
    width      = 1920
    height     = 1080
    quality    = "high"
    format     = "jpeg"
    auto_clean = true
  }
}
```

## Argument Reference

* `session_id` - (Required) The ID of the session to enable screenshots for.
* `enabled` - (Required) Whether screenshot functionality is enabled.
* `settings` - (Optional) Screenshot settings:
  * `width` - Screenshot width in pixels
  * `height` - Screenshot height in pixels
  * `quality` - Image quality setting
  * `format` - Image format (jpeg/png)
  * `auto_clean` - Automatically clean old screenshots

## Attribute Reference

* `id` - The unique identifier for the screenshot configuration.
* `latest_screenshot` - URL or path to the latest screenshot.
* `screenshot_count` - Number of screenshots taken.

## Import

Import a screenshot configuration:

```bash
terraform import kasm_screenshot.example session_id:screenshot_config_id
```

## Notes

1. Performance:
   - Resolution affects performance
   - Quality settings impact size
   - Format selection trade-offs

2. Storage:
   - Automatic cleanup options
   - Storage location configuration
   - Retention policies

3. Usage:
   - Session monitoring
   - Automated testing
   - Documentation generation

4. Limitations:
   - Resource intensive
   - Storage space requirements
   - Network bandwidth impact
