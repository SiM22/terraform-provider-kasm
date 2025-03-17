package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"terraform-provider-kasm/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupMembership_basic(t *testing.T) {
	t.Parallel()

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Test
			{
				Config: testAccGroupMembershipConfig_basic(username, groupname),
				Check: resource.ComposeTestCheckFunc(
					// Check resources exist
					testutils.TestCheckResourceExists("kasm_user.test"),
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_membership.test"),

					// Check attributes
					resource.TestCheckResourceAttrSet("kasm_group_membership.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test", "group_id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test", "user_id"),
				),
			},
			// Import Test
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
		},
	})
}

func TestAccGroupMembership_multiple(t *testing.T) {
	t.Parallel()

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	username1 := fmt.Sprintf("testuser1_%s", uniqueIdentifier)
	username2 := fmt.Sprintf("testuser2_%s", uniqueIdentifier)
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Test Multiple Users in Group
			{
				Config: testAccGroupMembershipConfig_multiple(username1, username2, groupname),
				Check: resource.ComposeTestCheckFunc(
					// Check resources exist
					testutils.TestCheckResourceExists("kasm_user.test1"),
					testutils.TestCheckResourceExists("kasm_user.test2"),
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_membership.test1"),
					testutils.TestCheckResourceExists("kasm_group_membership.test2"),

					// Check attributes for first membership
					resource.TestCheckResourceAttrSet("kasm_group_membership.test1", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test1", "group_id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test1", "user_id"),

					// Check attributes for second membership
					resource.TestCheckResourceAttrSet("kasm_group_membership.test2", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test2", "group_id"),
					resource.TestCheckResourceAttrSet("kasm_group_membership.test2", "user_id"),
				),
			},
		},
	})
}

func testAccGroupMembershipConfig_basic(username, groupname string) string {
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

func testAccGroupMembershipConfig_multiple(username1, username2, groupname string) string {
	return fmt.Sprintf(`
%s

# Create first user
resource "kasm_user" "test1" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User1"
}

# Create second user
resource "kasm_user" "test2" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User2"
}

# Create group
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Test group for acceptance tests"
}

# Create first membership
resource "kasm_group_membership" "test1" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test1.id
}

# Create second membership
resource "kasm_group_membership" "test2" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test2.id
}
`, testutils.ProviderConfig(), username1, username2, groupname)
}
