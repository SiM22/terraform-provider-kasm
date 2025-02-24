//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	// "terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"
)

const kasmUserClient = "kasm_client"

func TestAccUserImport_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("tf-%d", time.Now().UnixNano())
	username := fmt.Sprintf("test-%s@kasm.local", uniqueIdentifier)
	resourceName := "kasm_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("username:%s", username),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserConfig_basic(username string) string {
	return fmt.Sprintf(`
%s

resource "kasm_user" "test" {
    username     = "%s"
    password     = "Password123!"
    first_name   = "Test"
    last_name    = "User"
    organization = "Test Org"
    locked       = false
    disabled     = false
    groups       = []   # Explicitly set as empty list
}
`, testutils.ProviderConfig(), username)
}

func testAccUserConfig_withGroups(username string, groups []string) string {
	// Convert groups array to HCL list string
	groupsStr := "[]"
	if len(groups) > 0 {
		groupsStr = fmt.Sprintf("[\"%s\"]", strings.Join(groups, "\", \""))
	}

	return fmt.Sprintf(`
%s

resource "kasm_user" "test" {
    username     = "%s"
    password     = "Password123!"
    first_name   = "Test"
    last_name    = "User"
    organization = "Test Org"
    locked       = false
    disabled     = false
    groups       = %s
}
`, testutils.ProviderConfig(), username, groupsStr)
}

func TestAccUserImport_existing(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("tf-%d", time.Now().UnixNano())
	username := fmt.Sprintf("test-%s@kasm.local", uniqueIdentifier)
	resourceName := "kasm_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("username:%s", username),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// Add this function for import testing
func testAccUserConfig_forImport(username, uniqueIdentifier string) string {
	return fmt.Sprintf(`
%s

resource "kasm_user" "test" {
    username     = "%s"
    password     = "TestImport123!%s"
    first_name   = "PreExisting"
    last_name    = "User"
    organization = "Import Test Org"
}
`, testutils.ProviderConfig(), username, uniqueIdentifier)
}

// Helper function to check if user exists
func testAccCheckUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		client := testutils.GetTestClient(nil)
		_, err := client.GetUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error fetching user with ID %s: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

// Replace plancheck.CheckResourceDestroy with a custom destroy check
func testAccCheckUserDestroy(s *terraform.State) error {
	client := testutils.GetTestClient(nil)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kasm_user" {
			continue
		}

		_, err := client.GetUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Empty config for import testing
func testAccUserConfig_empty() string {
	return fmt.Sprintf(`
%s
`, testutils.ProviderConfig())
}
