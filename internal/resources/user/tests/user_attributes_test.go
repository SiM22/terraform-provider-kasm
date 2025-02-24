//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmUserAttributes(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create a user with attributes
			{
				Config: testAccKasmUserWithAttributesConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),
					resource.TestCheckResourceAttr("kasm_user.test", "first_name", "Test"),
					resource.TestCheckResourceAttr("kasm_user.test", "last_name", "User"),
					resource.TestCheckResourceAttr("kasm_user.test", "organization", "Test Org"),
					resource.TestCheckResourceAttr("kasm_user.test", "attributes.theme", "dark"),
					resource.TestCheckResourceAttr("kasm_user.test", "attributes.language", "en"),
				),
			},
		},
	})
}

func testAccKasmUserWithAttributesConfig(username string) string {
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

    attributes = {
        theme = "dark"
        language = "en"
    }

    groups = []
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), username)
}
