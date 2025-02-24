package tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-kasm/testutils"
)

func TestAccGroupMembership_import(t *testing.T) {
	t.Parallel()

	// Generate unique identifiers with both timestamp and random number
	rand.Seed(time.Now().UnixNano())
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Intn(10000))
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resources
			{
				Config: testAccGroupMembershipConfig_forImport(username, groupname),
				Check: resource.ComposeTestCheckFunc(
					testutils.TestCheckResourceExists("kasm_user.test"),
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_membership.test"),
				),
			},
			// Test successful import using group_id:user_id format
			{
				ResourceName:            "kasm_group_membership.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["kasm_group_membership.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "kasm_group_membership.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["user_id"]), nil
				},
			},
			// Test import with invalid ID format
			{
				ResourceName:  "kasm_group_membership.test",
				ImportState:   true,
				ImportStateId: "invalid-id-format",
				ExpectError:   regexp.MustCompile(`Import ID must be in the format group_id:user_id`),
			},
			// Test import with non-existent IDs
			{
				ResourceName:  "kasm_group_membership.test",
				ImportState:   true,
				ImportStateId: "non-existent-group:non-existent-user",
				ExpectError:   regexp.MustCompile(`Error Reading User`),
			},
		},
	})
}

func testAccGroupMembershipConfig_forImport(username string, groupname string) string {
	return fmt.Sprintf(`
%s

# Create user
resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
}

# Create group
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Test group for acceptance tests"
}

# Create membership
resource "kasm_group_membership" "test" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test.id
}
`, testutils.ProviderConfig(), username, groupname)
}
