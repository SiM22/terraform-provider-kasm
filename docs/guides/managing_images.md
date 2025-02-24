# Managing Images with the Kasm Provider

This guide provides instructions on how to manage images in Kasm Workspaces using the Terraform provider.

## Creating an Image

To create an image in Kasm, use the following resource configuration:

```hcl
resource "kasm_image" "new_image" {
  name        = "Ubuntu 20.04"
  description = "A base image for development"
}
```

## Updating an Image

To update an existing image, you can modify the attributes in the resource configuration:

```hcl
resource "kasm_image" "existing_image" {
  name        = "Ubuntu 20.04"
  description = "Updated description for the image"
}
```

## Deleting an Image

To delete an image, you can use the following configuration:

```hcl
resource "kasm_image" "delete_image" {
  name = "Image to Delete"
}
```

## Image Attributes

When creating or updating an image, you can specify the following attributes:
- `name`: The name of the image (required).
- `description`: A description of the image (optional).

## Example Usage

Here’s a complete example of managing an image:

```hcl
resource "kasm_image" "example" {
  name        = "Ubuntu 20.04"
  description = "A base image for development"
}

resource "kasm_image" "updated_image" {
  name        = "Ubuntu 20.04"
  description = "Updated description for the image"
}

resource "kasm_image" "deleted_image" {
  name = "Image to Delete"
}
```

## Next Steps

After managing images, you can refer to the following guides for additional management tasks:
- [Managing Users](managing_users.md)
- [Managing Groups](managing_groups.md)

For more details, consult the [Kasm Dev API Docs](Kasm%20Dev%20API%20Docs).
