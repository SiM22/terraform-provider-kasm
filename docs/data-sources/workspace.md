# Workspace Data Source

Provides information about Kasm workspaces and their configurations.

## Example Usage

### Basic Query
```hcl
data "kasm_workspace" "current" {}

output "workspace_info" {
  value = data.kasm_workspace.current
}
```

### Filtered Query
```hcl
data "kasm_workspace" "dev" {
  filter = {
    name = "development"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter criteria for workspaces:
  * `name` - Filter by workspace name
  * `status` - Filter by status
  * `type` - Filter by workspace type

## Attribute Reference

* `workspaces` - List of workspaces matching the filter criteria:
  * `id` - Workspace identifier
  * `name` - Workspace name
  * `status` - Current status
  * `type` - Workspace type
  * `created_at` - Creation timestamp
  * `updated_at` - Last update timestamp
  * `configuration` - Workspace configuration details

## Notes

1. Usage:
   - Environment information
   - Configuration details
   - Status monitoring

2. Filtering:
   - Name-based filtering
   - Status filtering
   - Type selection

3. Performance:
   - Real-time data
   - Cached information
   - Efficient queries

4. Integration:
   - Resource management
   - Configuration planning
   - Status monitoring
