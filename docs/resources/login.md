# Login Resource

Generates a login URL that allows users to access Kasm without entering credentials.

## Example Usage

```hcl
resource "kasm_login" "example" {
  user_id = kasm_user.example.id
}

output "login_url" {
  value     = kasm_login.example.login_url
  sensitive = true
}
```

## Argument Reference

* `user_id` - (Required) The ID of the user to generate a login URL for.

## Attribute Reference

* `id` - The resource identifier (same as user_id).
* `login_url` - The generated login URL. This URL can be used to access Kasm without entering credentials.

## Notes

1. Security:
   - The login URL is sensitive and should be handled securely
   - URLs are typically valid for a limited time
   - Each generation creates a new unique URL

2. Usage:
   - Useful for automated access
   - Single sign-on integration
   - System automation

3. Limitations:
   - URLs are ephemeral
   - Subject to session timeout policies
   - Requires valid user ID
