---
page_title: "kasm_sessions Data Source - terraform-provider-kasm"
subcategory: ""
description: |-
  Retrieves a list of all active Kasm sessions.
---

# kasm_sessions (Data Source)

Retrieves a list of all active Kasm sessions.

## Example Usage

```terraform
data "kasm_sessions" "all" {}

# Access a specific session by index
output "first_session_id" {
  value = length(data.kasm_sessions.all.sessions) > 0 ? data.kasm_sessions.all.sessions[0].kasm_id : "No sessions"
}

# Access a specific session by kasm_id
output "specific_session_status" {
  value = lookup(data.kasm_sessions.all.sessions_map, "your-kasm-id", null) != null ? data.kasm_sessions.all.sessions_map["your-kasm-id"].operational_status : "Session not found"
}
```

## Argument Reference

This data source has no arguments.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `current_time` - Current time as reported by the Kasm API.
* `sessions` - List of all active Kasm sessions.
  * `expiration_date` - Date and time when the session will expire.
  * `container_ip` - IP address of the container.
  * `start_date` - Date and time when the session was started.
  * `token` - Session token.
  * `image_id` - ID of the image used for the session.
  * `view_only_token` - View-only token for the session.
  * `cores` - Number of CPU cores allocated to the session.
  * `hostname` - Hostname of the session container.
  * `kasm_id` - Unique identifier for the Kasm session.
  * `port_map` - Map of port mappings for the session.
  * `image_name` - Name of the image used for the session.
  * `image_friendly_name` - Friendly name of the image used for the session.
  * `image_src` - Source URL of the image used for the session.
  * `is_persistent_profile` - Whether the session uses a persistent profile.
  * `memory` - Amount of memory allocated to the session in MB.
  * `operational_status` - Current operational status of the session.
  * `container_id` - ID of the container running the session.
  * `port` - Port number used by the session.
  * `keepalive_date` - Date and time of the last keepalive signal.
  * `user_id` - ID of the user who owns the session.
  * `share_id` - Share ID for the session, if shared.
  * `host` - Host where the session is running.
  * `server_id` - ID of the server running the session.
* `sessions_map` - Map of all active Kasm sessions, indexed by kasm_id. Contains the same attributes as the sessions list.
