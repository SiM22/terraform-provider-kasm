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

func TestAccKasmUser_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmUserConfig_basic(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),
					resource.TestCheckResourceAttr("kasm_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("kasm_user.test", "last_name", "User"),
					resource.TestCheckResourceAttr("kasm_user.test", "organization", "Test Org"),
				),
			},
		},
	})
}

func testAccKasmUserConfig_basic(username string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
    organization = "Test Org"
    locked = false
    disabled = false

    # Empty list in HCL format
    groups = []
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), username)
}
