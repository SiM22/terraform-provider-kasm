# Stats Resource

Manages frame statistics for Kasm sessions. This resource allows you to retrieve performance and usage statistics for workspace sessions.

## Example Usage

### Basic Stats Configuration
```hcl
resource "kasm_stats" "example" {
  kasm_id = kasm_session.workspace.id
}
```

### Stats with User ID
```hcl
resource "kasm_stats" "detailed" {
  kasm_id = kasm_session.workspace.id
  user_id = kasm_user.example.id
}
```

## Argument Reference

* `kasm_id` - (Required) The ID of the Kasm session to retrieve stats for.
* `user_id` - (Optional) The ID of the user who owns the session.

## Attribute Reference

* `id` - The unique identifier for the stats resource.
* `res_x` - The horizontal resolution of the session.
* `res_y` - The vertical resolution of the session.
* `changed` - The number of changed pixels.
* `server_time` - The server processing time in milliseconds.
* `client_count` - The number of connected clients.
* `analysis` - The time spent on frame analysis in milliseconds.
* `screenshot` - The time spent on screenshot processing in milliseconds.
* `encoding_time` - The total encoding time in milliseconds.
* `last_updated` - Timestamp of the last refresh of the stats.

## Import

Stats can be imported using the Kasm ID:

```bash
terraform import kasm_stats.example <kasm_id>
```

## Notes

1. Performance Impact:
   - Collection interval affects overhead
   - Selective metric collection
   - Resource usage considerations

2. Data Collection:
   - CPU usage and load
   - Memory utilization
   - Network traffic
   - Disk operations

3. Storage:
   - Time-series data
   - Aggregation options
   - Retention settings

4. Usage:
   - Performance monitoring
   - Resource planning
   - Usage analysis
   - Billing information
