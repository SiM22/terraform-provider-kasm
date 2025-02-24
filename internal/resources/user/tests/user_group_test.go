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

func TestAccKasmGroupAndUser_advanced(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	initialGroupName := fmt.Sprintf("InitialGroup_%s", uniqueIdentifier)
	updatedGroupName := fmt.Sprintf("UpdatedGroup_%s", uniqueIdentifier)
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial configuration
			{
				Config: testutils.TestAccUserGroupConfig(username, initialGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_group.test", "name", initialGroupName),
					resource.TestCheckResourceAttr("kasm_group.test", "priority", "50"),
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),
					resource.TestCheckResourceAttr("kasm_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("kasm_user.test", "last_name", "Import"),
				),
			},
			// Updated configuration
			{
				Config: testutils.TestAccUserGroupConfig(username, updatedGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_group.test", "name", updatedGroupName),
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),
				),
			},
		},
	})
}

func testAccKasmGroupAndUserConfig_initial(groupName, username string) string {
	return fmt.Sprintf(`
provider "kasm" {
  base_url = "%s"
  api_key = "%s"
  api_secret = "%s"
  insecure = true
}

resource "kasm_group" "test_group" {
  name = "%s"
  priority = 50
  description = "Initial test group"
}

resource "kasm_user" "test_user" {
  username = "%s"
  first_name = "Test"
  last_name = "User"
  password = "initialPassword123!"
  locked = false
  disabled = false
  groups = [kasm_group.test_group.name]

}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), groupName, username)
}

func testAccKasmGroupAndUserConfig_updated(groupName, username string) string {
	return fmt.Sprintf(`
provider "kasm" {
  base_url = "%s"
  api_key = "%s"
  api_secret = "%s"
  insecure = true
}

resource "kasm_group" "test_group" {
  name = "%s"
  priority = 75
  description = "Updated test group"
}

resource "kasm_user" "test_user" {
  username = "%s"
  first_name = "Updated"
  last_name = "TestUser"
  password = "updatedPassword123!"
  locked = false
  disabled = false
  groups = [kasm_group.test_group.name]

  depends_on = [kasm_group.test_group]
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), groupName, username)
}
