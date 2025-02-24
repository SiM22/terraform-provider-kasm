//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-kasm/testutils"
)

const kasmGroupClient = "kasm_client"

// And add local providerConfig function if needed
func providerConfig() string {
	return testutils.ProviderConfig()
}

func TestAccGroupImport_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("tf-%d", time.Now().UnixNano())
	groupName := fmt.Sprintf("test-group-%s", uniqueIdentifier)
	resourceName := "kasm_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create the resource
			{
				Config: testAccGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
				),
			},
			// Then test importing it
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("name:%s", groupName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name        = "%s"
    description = "Test Group"
    priority    = 50
}
`, testutils.ProviderConfig(), name)
}

func TestAccGroupImport_existing(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("tf-%d", time.Now().UnixNano())
	groupName := fmt.Sprintf("import-test-group-%s", uniqueIdentifier)
	resourceName := "kasm_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create the resource
			{
				Config: testAccGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
				),
			},
			// Import by name
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("name:%s", groupName),
				ImportStateVerify: true,
			},
		},
	})
}

// Helper function to check if group exists
func testAccCheckGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Group ID is set")
		}

		client := testutils.GetTestClient(nil)
		_, err := client.GetGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error fetching group with ID %s: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

// Empty config for import testing
func testAccGroupConfig_empty() string {
	return fmt.Sprintf(`
%s
`, testutils.ProviderConfig())
}
