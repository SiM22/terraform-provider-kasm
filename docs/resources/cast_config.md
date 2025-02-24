# Resource: kasm_cast_config

Manages a Kasm casting configuration. Casting configurations allow you to create and manage session casting links that can be used to launch Kasm sessions.

## Example Usage

```hcl
# Create a basic casting configuration
resource "kasm_cast_config" "example" {
  name    = "Chrome Configuration"
  image_id = "90c473864949488c902aa6e4adbf89a6"
  key     = "abc123"

  allowed_referrers = [
    "acme.com",
    "contoso.com"
  ]

  # Session limits
  limit_sessions    = true
  session_remaining = 100

  # IP rate limiting
  limit_ips           = true
  ip_request_limit    = 1
  ip_request_seconds  = 60

  # Error handling
  error_url = "https://error.example.com"

  # Session settings
  enable_sharing         = true
  disable_control_panel = false
  disable_tips         = false
  disable_fixed_res    = false

  # Authentication
  allow_anonymous    = true
  group_id          = "68d557ac4cac42cca9f31c7c853de0f3"
  require_recaptcha = true

  # URL handling
  kasm_url              = "https://google.com"
  dynamic_kasm_url      = true
  dynamic_docker_network = false
  allow_resume          = true

  # Client settings
  enforce_client_settings   = true
  allow_kasm_audio         = true
  allow_kasm_uploads       = false
  allow_kasm_downloads     = false
  allow_kasm_clipboard_down = false
  allow_kasm_clipboard_up   = false
  allow_kasm_microphone    = false
  allow_kasm_sharing       = false
  kasm_audio_default_on    = true
  kasm_ime_mode_default_on = true

  # Expiration
  valid_until = "2024-12-31 23:59:59"
}

# Output the casting URL
output "casting_url" {
  value = "https://my.kasm.server/#/cast/${kasm_cast_config.example.key}"
}
```

## Argument Reference

* `name` - (Required) The configuration name.
* `image_id` - (Required) The Image ID to use for the casting config.
* `key` - (Required) The unique identifier for the Casting URL. Users will launch sessions via `https://my.kasm.server/#/cast/<key>`.
* `allowed_referrers` - (Optional) A list of domains allowed as referrers when a casting link is visited.
* `limit_sessions` - (Optional) When enabled, the total number of sessions for this config will be limited.
* `session_remaining` - (Optional) The number of sessions that are allowed to be spawned from this casting link.
* `limit_ips` - (Optional) When enabled, the system will limit requests based on source IP.
* `ip_request_limit` - (Optional) The total number of sessions allowed for the given time period.
* `ip_request_seconds` - (Optional) The timeframe in seconds for IP rate limiting.
* `error_url` - (Optional) URL to redirect to when an error occurs.
* `enable_sharing` - (Optional) When enabled, this session will automatically have sharing activated.
* `disable_control_panel` - (Optional) When enabled, the Control Panel widget is not shown.
* `disable_tips` - (Optional) When enabled, the Tips dialogue is not shown.
* `disable_fixed_res` - (Optional) When enabled and in sharing mode, the resolution will be dynamic.
* `allow_anonymous` - (Optional) If enabled, requests will not require authentication.
* `group_id` - (Optional) The group ID for anonymous users.
* `require_recaptcha` - (Optional) When enabled, requests will be validated by Google reCAPTCHA.
* `kasm_url` - (Optional) The URL to populate as KASM_URL environment variable.
* `dynamic_kasm_url` - (Optional) Allow kasm_url query parameter in cast URL.
* `dynamic_docker_network` - (Optional) Allow docker_network query parameter in cast URL.
* `allow_resume` - (Optional) Allow session resumption for authenticated users.
* `enforce_client_settings` - (Optional) Enforce client settings on the session.
* `allow_kasm_audio` - (Optional) Allow audio streaming from the session.
* `allow_kasm_uploads` - (Optional) Allow file uploads to the session.
* `allow_kasm_downloads` - (Optional) Allow file downloads from the session.
* `allow_kasm_clipboard_down` - (Optional) Allow clipboard copy from session to local.
* `allow_kasm_clipboard_up` - (Optional) Allow clipboard copy from local to session.
* `allow_kasm_microphone` - (Optional) Allow microphone access in the session.
* `valid_until` - (Optional) The time until which the casting link is valid (UTC).
* `allow_kasm_sharing` - (Optional) Allow the user to place their session in sharing mode.
* `kasm_audio_default_on` - (Optional) Enable audio by default.
* `kasm_ime_mode_default_on` - (Optional) Enable IME mode by default.

## Attribute Reference

* `id` - The ID of the cast configuration.
* `image_friendly_name` - The friendly name of the selected image.
* `group_name` - The name of the selected group.

## Import

Cast configurations can be imported using their ID:

```shell
terraform import kasm_cast_config.example <cast_config_id>
```

## Notes

* The casting URL is constructed by combining your Kasm server URL with the cast key: `https://my.kasm.server/#/cast/<key>`
* When `allow_anonymous` is enabled, the system will create new user accounts for each request.
* Anonymous users are automatically added to the All Users Group and the group specified in `group_id`.
* To use reCAPTCHA validation, the Google reCAPTCHA Private Key and Site Key must be set in Server Settings.
* Client settings (`allow_kasm_*`) only apply when `enforce_client_settings` is enabled.
* When `dynamic_kasm_url` is enabled, users can append `?kasm_url=example.com` to the casting URL.
* When `dynamic_docker_network` is enabled, users can append `?docker_network=example_network` to the casting URL.
