# RDP Client Connection Info Data Source

Retrieves RDP client connection information for a Kasm session.

> **Note:** This functionality requires a Windows RDP server configured in Kasm. For container-based sessions, the RDP connection info will be empty. See the [Kasm documentation on Fixed Infrastructure](https://kasmweb.com/docs/latest/how_to/fixed_infrastructure.html#fixed-infrastructure-rdp-vnc-ssh-kasmvnc) for more details on setting up RDP servers.

## Example Usage

### Basic Usage
```hcl
data "kasm_rdp_client_connection_info" "example_file" {
  user_id = "44edb3e5-2909-4927-a60b-6e09c7219104"
  kasm_id = "898813d7-a677-4c60-8999-c9ea346a3e21"
  connection_type = "file"
}

output "rdp_file" {
  value = data.kasm_rdp_client_connection_info.example_file.file
}

data "kasm_rdp_client_connection_info" "example_url" {
  user_id = "44edb3e5-2909-4927-a60b-6e09c7219104"
  kasm_id = "898813d7-a677-4c60-8999-c9ea346a3e21"
  connection_type = "url"
}

output "rdp_url" {
  value = data.kasm_rdp_client_connection_info.example_url.url
}
```

### Advanced Usage
```hcl
data "kasm_rdp_client_connection_info" "custom" {
  user_id         = kasm_user.example.id
  kasm_id         = kasm_session.example.id
  connection_type = "file"
}

output "rdp_file" {
  value     = data.kasm_rdp_client_connection_info.custom.file
  sensitive = true
}
```

## Argument Reference

* `user_id` - (Required) The ID of the user requesting the connection.
* `kasm_id` - (Required) The ID of the Kasm session.
* `connection_type` - (Optional) The type of connection to retrieve ("url" or "file"). Defaults to "url".

## Attribute Reference

* `id` - The unique identifier for the connection info.
* `file` - The RDP connection file content (if connection_type is "file").
* `url` - The URL for RDP access (if connection_type is "url").
