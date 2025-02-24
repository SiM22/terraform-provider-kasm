//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmRegistryImagesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryImagesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.#"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.id"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.friendly_name"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.description"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.memory"),
					resource.TestCheckResourceAttrSet("data.kasm_registry_images.test", "images.0.cores"),
				),
			},
		},
	})
}

func TestAccKasmRegistryImagesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryImagesDataSourceConfigEmpty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kasm_registry_images.test", "images.#", "0"),
				),
			},
		},
	})
}

func TestAccKasmRegistryImagesDataSource_invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccKasmRegistryImagesDataSourceConfigInvalid(),
				ExpectError: regexp.MustCompile(`Registry ID .* is invalid`),
			},
		},
	})
}

func testAccKasmRegistryImagesDataSourceConfig() string {
	return fmt.Sprintf(`
provider "kasm" {
	base_url = "%s"
	api_key = "%s"
	api_secret = "%s"
	insecure = true
}

data "kasm_registry_images" "test" {}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryImagesDataSourceConfigEmpty() string {
	return fmt.Sprintf(`
provider "kasm" {
	base_url = "%s"
	api_key = "%s"
	api_secret = "%s"
	insecure = true
}

data "kasm_registry_images" "test" {
	registry_id = "12345678-1234-1234-1234-123456789abc"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryImagesDataSourceConfigInvalid() string {
	return fmt.Sprintf(`
provider "kasm" {
	base_url = "%s"
	api_key = "%s"
	api_secret = "%s"
	insecure = true
}

data "kasm_registry_images" "test" {
	registry_id = "invalid-registry-id"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}
