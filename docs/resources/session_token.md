# Resource: kasm_session_token

Manages a Kasm session token. Session tokens authenticate user's requests to access functionality within the system.

## Example Usage

```hcl
# Create a session token for a user
resource "kasm_session_token" "example" {
  user_id = "009c3779-4fa0-4af8-9722-8daf195718c0"
}

# Output the JWT token for use in other configurations
output "session_jwt" {
  value     = kasm_session_token.example.session_jwt
  sensitive = true
}
```

## Argument Reference

* `user_id` - (Required) The ID of the user to create the session token for.

## Attribute Reference

* `id` - The ID of the session token (same as `session_token`).
* `session_token` - The value of the session token.
* `session_token_date` - The time the token was created or last promoted.
* `expires_at` - The time the token will no longer be valid. This is session_token_date + the global setting "Session Lifetime".
* `session_jwt` - The JWT token used by clients for authentication. This value is marked as sensitive.

## Import

Session tokens can be imported using their token value:

```shell
terraform import kasm_session_token.example <session_token>
```

## Notes

* Session tokens have a limited lifetime determined by the global "Session Lifetime" setting.
* The token can be promoted (refreshed) by updating the resource, which will extend its lifetime.
* When the resource is destroyed, the session token is invalidated.
* Multiple session tokens can be created for the same user.
* The session JWT is sensitive information and should be handled securely.
