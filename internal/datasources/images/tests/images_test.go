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

func TestAccKasmImagesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmImagesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Check that we have images
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.#"),
					// Check first image has all required attributes
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.id"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.friendly_name"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.description"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.memory"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.cores"),
					resource.TestCheckResourceAttrSet("data.kasm_images.test", "images.0.cpu_allocation_method"),
					// Verify we have at least one image
					resource.TestMatchResourceAttr("data.kasm_images.test", "images.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

// Basic config to list all images
func testAccKasmImagesDataSourceConfig() string {
	return fmt.Sprintf(`
	provider "kasm" {
		base_url = "%s"
		api_key = "%s"
		api_secret = "%s"
		insecure = true
	}

	data "kasm_images" "test" {}
	`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}
