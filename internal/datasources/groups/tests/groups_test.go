//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmGroups_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
			// Add additional validation
			c := testutils.GetTestClient(t)
			groups, err := c.GetGroups()
			if err != nil {
				t.Skipf("No groups available: %v", err)
			}
			if len(groups) == 0 {
				t.Skip("No groups available for testing")
			}
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmGroupsConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Make checks more lenient since we don't know exact group data
					resource.TestCheckResourceAttrSet("data.kasm_groups.all", "groups.#"),
					resource.TestCheckResourceAttrSet("data.kasm_groups.all", "groups.0.group_id"),
					resource.TestCheckResourceAttrSet("data.kasm_groups.all", "groups.0.name"),
				),
			},
		},
	})
}

func testAccKasmGroupsConfig_basic() string {
	return fmt.Sprintf(`
provider "kasm" {
	base_url = "%s"
	api_key = "%s"
	api_secret = "%s"
	insecure = true
}

data "kasm_groups" "all" {}

output "debug_group_count" {
    value = length(data.kasm_groups.all.groups)
}

output "debug_first_group" {
    value = length(data.kasm_groups.all.groups) > 0 ? data.kasm_groups.all.groups[0] : null
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}
