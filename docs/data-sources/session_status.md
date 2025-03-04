---
page_title: "kasm_session_status Data Source - terraform-provider-kasm"
subcategory: ""
description: |-
  Retrieves the status of a specific Kasm session.
---

# kasm_session_status (Data Source)

Retrieves the status of a specific Kasm session.

## Example Usage

```terraform
data "kasm_session_status" "example" {
  kasm_id = "your-kasm-id"
  user_id = "your-user-id"
  skip_agent_check = true
}

output "session_operational_status" {
  value = data.kasm_session_status.example.operational_status
}
```

## Argument Reference

The following arguments are supported:

* `kasm_id` - (Required) The ID of the Kasm session to check.
* `user_id` - (Required) The ID of the user who owns the session.
* `skip_agent_check` - (Optional) Whether to skip checking the agent status. Defaults to `false`.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `status` - The status of the Kasm session (e.g., "running").
* `operational_status` - The operational status of the Kasm session.
* `operational_message` - A message describing the current operational status.
* `error_message` - Error message if any.
* `kasm_url` - URL to access the Kasm session.
* `container_ip` - The IP address of the container running the session.
* `port` - The port number used by the session.
* `container_id` - The ID of the container running the session.
* `server_id` - The ID of the server running the session.
* `host` - The host where the session is running.
* `hostname` - The hostname of the session container.
* `image_id` - The ID of the image used for the session.
* `image_name` - The name of the image used for the session.
