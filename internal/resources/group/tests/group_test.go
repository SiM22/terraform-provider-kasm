//go:build acceptance
// +build acceptance

// #nosec G101 -- Test credentials
package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmGroup_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	groupName := fmt.Sprintf("test-group-%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_group.test", "name", groupName),
					resource.TestCheckResourceAttr("kasm_group.test", "priority", "50"),
					resource.TestCheckResourceAttr("kasm_group.test", "description", "Test group"),
				),
			},
		},
	})
}

func testAccKasmGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_group" "test" {
    name = "%s"
    priority = 50
    description = "Test group"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), name)
}
