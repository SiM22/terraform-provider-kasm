# Stats Resource

Manages statistics collection for Kasm sessions. This resource allows you to configure and manage performance and usage statistics for workspace sessions.

## Example Usage

### Basic Stats Configuration
```hcl
resource "kasm_stats" "example" {
  session_id = kasm_session.workspace.id
  enabled    = true
}
```

### Detailed Stats Configuration
```hcl
resource "kasm_stats" "detailed" {
  session_id = kasm_session.workspace.id
  enabled    = true

  collection_settings = {
    interval     = 60
    cpu_stats    = true
    memory_stats = true
    network_stats = true
    disk_stats   = true
  }
}
```

## Argument Reference

* `session_id` - (Required) The ID of the session to collect stats for.
* `enabled` - (Required) Whether stats collection is enabled.
* `collection_settings` - (Optional) Statistics collection settings:
  * `interval` - Collection interval in seconds
  * `cpu_stats` - Collect CPU statistics
  * `memory_stats` - Collect memory statistics
  * `network_stats` - Collect network statistics
  * `disk_stats` - Collect disk statistics

## Attribute Reference

* `id` - The unique identifier for the stats configuration.
* `current_stats` - Current statistics snapshot.
* `collection_status` - Status of stats collection.

## Import

Import a stats configuration:

```bash
terraform import kasm_stats.example session_id:stats_config_id
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
