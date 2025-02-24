package tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccGroupMembership_errors(t *testing.T) {
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
			// Try to add user to non-existent group
			{
				Config:      testAccGroupMembershipConfig_nonExistentGroup(username),
				ExpectError: regexp.MustCompile(`Error adding user to group`),
			},
			// Try to add non-existent user to group
			{
				Config:      testAccGroupMembershipConfig_nonExistentUser(groupname),
				ExpectError: regexp.MustCompile(`Error adding user to group`),
			},
		},
	})
}

func testAccGroupMembershipConfig_nonExistentGroup(username string) string {
	return fmt.Sprintf(`
%s

# Create user
resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
}

# Create membership with non-existent group
resource "kasm_group_membership" "test" {
    group_id = "non-existent-group-id"
    user_id = kasm_user.test.id
}
`, testutils.ProviderConfig(), username)
}

func testAccGroupMembershipConfig_nonExistentUser(groupname string) string {
	return fmt.Sprintf(`
%s

# Create group
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Test group for acceptance tests"
}

# Create membership with non-existent user
resource "kasm_group_membership" "test" {
    group_id = kasm_group.test.id
    user_id = "non-existent-user-id"
}
`, testutils.ProviderConfig(), groupname)
}
