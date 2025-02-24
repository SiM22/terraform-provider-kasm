# Data Source: kasm_zones

Use this data source to get information about Kasm deployment zones.

## Example Usage

```hcl
# Get all zones with full details
data "kasm_zones" "all" {
  brief = false
}

# Get zones with limited information
data "kasm_zones" "brief" {
  brief = true
}

# Access individual zone information
output "default_zone_id" {
  value = [
    for zone in data.kasm_zones.all.zones :
    zone.id
    if zone.name == "default"
  ][0]
}
```

## Argument Reference

* `brief` - (Optional) Limit the information returned for each zone. Defaults to false.

## Attributes Reference

* `id` - The ID of this resource.
* `zones` - A list of zones. Each zone contains the following attributes:
  * `id` - The ID of the zone.
  * `name` - The name of the zone.
  * `auto_scaling_enabled` - Whether auto-scaling is enabled for the zone.
  * `aws_enabled` - Whether AWS integration is enabled for the zone.
  * `aws_region` - The AWS region configured for the zone.
  * `aws_access_key_id` - The AWS access key ID configured for the zone.
  * `aws_secret_access_key` - The AWS secret access key configured for the zone (sensitive).
  * `ec2_agent_ami_id` - The EC2 agent AMI ID configured for the zone.
