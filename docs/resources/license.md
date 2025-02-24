# License Resource

Manages Kasm license activation. This resource allows you to activate and manage licenses for your Kasm deployment.

## Example Usage

### Basic License Activation
```hcl
resource "kasm_license" "enterprise" {
  activation_key = "-----BEGIN ACTIVATION KEY-----eyJ0eXAiOiJ...-----END ACTIVATION KEY-----"
}
```

### Custom License Configuration
```hcl
resource "kasm_license" "enterprise" {
  activation_key = "-----BEGIN ACTIVATION KEY-----eyJ0eXAiOiJ...-----END ACTIVATION KEY-----"
  seats         = 10
  issued_to     = "ACME Corp"
}
```

## Argument Reference

* `activation_key` - (Required) The activation key provided by Kasm Technologies.
* `seats` - (Optional) The desired number of seats to license the deployment for.
* `issued_to` - (Optional) Organization the deployment is licensed for.

## Attribute Reference

* `id` - The license identifier.
* `expiration` - License expiration date.
* `issued_at` - License issue date.
* `limit` - Licensed seat limit.
* `is_verified` - Whether the license is verified.
* `license_type` - Type of license.
* `sku` - License SKU.

### Feature Flags

* `auto_scaling` - Auto scaling feature flag.
* `branding` - Branding feature flag.
* `session_staging` - Session staging feature flag.
* `session_casting` - Session casting feature flag.
* `log_forwarding` - Log forwarding feature flag.
* `developer_api` - Developer API feature flag.
* `inject_ssh_keys` - SSH key injection feature flag.
* `saml` - SAML feature flag.
* `ldap` - LDAP feature flag.
* `session_sharing` - Session sharing feature flag.
* `login_banner` - Login banner feature flag.
* `url_categorization` - URL categorization feature flag.
* `usage_limit` - Usage limit feature flag.

## Notes

1. License Activation:
   - Activation keys are sensitive and should be handled securely
   - Keys are validated during activation
   - Features are determined by the license type

2. Seat Management:
   - Optional seat specification
   - Maximum seats determined by entitlement
   - Seat changes require reactivation

3. License Types:
   - Per Concurrent Kasm
   - Features vary by SKU
   - Enterprise features available

4. Security:
   - Secure activation key storage
   - License verification
   - Feature enforcement
