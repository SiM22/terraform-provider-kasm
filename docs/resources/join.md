# kasm_join (Resource)

Join a shared Kasm session.

## Example Usage

```hcl
# Create a session
resource "kasm_session" "example" {
  image_id = "a5273e40ab5d47909aa7f047a2413da6"
  user_id = kasm_user.owner.id
  share = true
}

# Join the session as another user
resource "kasm_join" "example" {
  share_id = kasm_session.example.share_id
  user_id = kasm_user.viewer.id
}
```

## Argument Reference

* `share_id` - (Required) The share ID of the session to join.
* `user_id` - (Required) The ID of the user joining the session.

## Attribute Reference

* `id` - The ID of the join resource.
* `kasm_id` - The ID of the joined Kasm session.
* `session_token` - The session token for the joined session.
* `kasm_url` - The URL to access the joined session.

## Import

Join resources can be imported using the share ID:

```shell
terraform import kasm_join.example {share_id}
