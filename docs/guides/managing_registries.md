# Managing Registries with the Kasm Provider

This guide provides instructions on how to manage Docker registries in Kasm Workspaces using the Terraform provider.

## Creating a Registry

To create a registry in Kasm, use the following resource configuration:

```hcl
resource "kasm_registry" "new_registry" {
  name     = "My Private Registry"
  server   = "registry.example.com"
  username = "user"
  password = "password"
}
```

## Updating a Registry

To update an existing registry, you can modify the attributes in the resource configuration:

```hcl
resource "kasm_registry" "existing_registry" {
  name     = "My Private Registry"
  server   = "registry.example.com"
  username = "new_user"
  password = "new_password"
}
```

## Deleting a Registry

To delete a registry, you can use the following configuration:

```hcl
resource "kasm_registry" "delete_registry" {
  name = "Registry to Delete"
}
```

## Registry Attributes

When creating or updating a registry, you can specify the following attributes:
- `name`: The name of the registry (required).
- `server`: The server URL of the registry (required).
- `username`: The username for authentication (optional).
- `password`: The password for authentication (optional).

## Example Usage

Here’s a complete example of managing a registry:

```hcl
resource "kasm_registry" "example" {
  name     = "My Private Registry"
  server   = "registry.example.com"
  username = "user"
  password = "password"
}

resource "kasm_registry" "updated_registry" {
  name     = "My Private Registry"
  server   = "registry.example.com"
  username = "new_user"
  password = "new_password"
}

resource "kasm_registry" "deleted_registry" {
  name = "Registry to Delete"
}
```

## Next Steps

After managing registries, you can refer to the following guides for additional management tasks:
- [Managing Users](managing_users.md)
- [Managing Groups](managing_groups.md)

For more details, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs).
