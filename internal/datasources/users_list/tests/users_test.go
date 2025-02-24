//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmUsersDataSource_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create a user to ensure we have data
			{
				Config: testAccKasmUserConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),
					resource.TestCheckResourceAttr("kasm_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("kasm_user.test", "last_name", "User"),
					resource.TestCheckResourceAttr("kasm_user.test", "organization", "Test Org"),
				),
			},
			// Then test the users data source
			{
				Config: testAccKasmUserConfig(username) + testAccKasmUsersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Check that we have users
					resource.TestCheckResourceAttrSet("data.kasm_users.test", "users.#"),
					// Check that our test user exists and has the correct attributes
					resource.TestMatchTypeSetElemNestedAttrs("data.kasm_users.test", "users.*", map[string]*regexp.Regexp{
						"username": regexp.MustCompile(username),
					}),
					// Check that each user has the required fields
					resource.TestCheckResourceAttrSet("data.kasm_users.test", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.kasm_users.test", "users.0.username"),
					// Verify we have at least one user
					resource.TestMatchResourceAttr("data.kasm_users.test", "users.#", regexp.MustCompile(`[1-9][0-9]*`)),
				),
			},
		},
	})
}

// Create a test user
func testAccKasmUserConfig(username string) string {
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
	groups = []
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), username)
}

// Basic config to list all users
func testAccKasmUsersDataSourceConfig() string {
	return `
data "kasm_users" "test" {}
`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
