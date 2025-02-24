package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmSessionPermission_userSpecific(t *testing.T) {
	t.Parallel()

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	username1 := fmt.Sprintf("testuser1_%s", uniqueIdentifier)
	username2 := fmt.Sprintf("testuser2_%s", uniqueIdentifier)
	groupname := fmt.Sprintf("SessionGroup_%s", uniqueIdentifier)

	// Get and configure a test image
	imageID, available := ensureImageAvailable(t)
	if !available {
		t.Skip("No suitable test images available")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Basic Creation Test
			{
				Config: testAccKasmSessionPermissionConfig_userSpecific(username1, username2, groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					// Check resources exist
					testutils.TestCheckResourceExists("kasm_session.test"),
					testutils.TestCheckResourceExists("kasm_session_permission.test"),
					testutils.TestCheckResourceExists("kasm_join.test"),

					// Check session attributes with retries
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["kasm_session.test"]
						if !ok {
							return fmt.Errorf("kasm_session.test not found")
						}

						maxRetries := 5
						for i := 0; i < maxRetries; i++ {
							shareID := rs.Primary.Attributes["share_id"]
							if shareID != "" {
								// Verify other attributes
								if rs.Primary.ID == "" {
									return fmt.Errorf("no ID is set")
								}
								if rs.Primary.Attributes["operational_status"] == "" {
									return fmt.Errorf("operational_status is empty")
								}
								if rs.Primary.Attributes["share"] != "true" {
									return fmt.Errorf("share is not true")
								}
								if rs.Primary.Attributes["enable_sharing"] != "true" {
									return fmt.Errorf("enable_sharing is not true")
								}
								return nil
							}

							// Wait before retrying
							time.Sleep(2 * time.Second)

							// Get fresh state
							rs, ok = s.RootModule().Resources["kasm_session.test"]
							if !ok {
								return fmt.Errorf("kasm_session.test not found during retry")
							}
						}

						return fmt.Errorf("share_id is still empty after %d retries", maxRetries)
					},

					// Check permissions
					resource.TestCheckResourceAttrSet("kasm_session_permission.test", "kasm_id"),
					resource.TestCheckResourceAttr("kasm_session_permission.test", "user_permissions.#", "1"),

					// Check join
					resource.TestCheckResourceAttrSet("kasm_join.test", "kasm_url"),

					// Custom check for share_id consistency
					func(s *terraform.State) error {
						sessionRS, ok := s.RootModule().Resources["kasm_session.test"]
						if !ok {
							return fmt.Errorf("kasm_session.test not found")
						}

						joinRS, ok := s.RootModule().Resources["kasm_join.test"]
						if !ok {
							return fmt.Errorf("kasm_join.test not found")
						}

						sessionShareID := sessionRS.Primary.Attributes["share_id"]
						joinShareID := joinRS.Primary.Attributes["share_id"]

						if sessionShareID == "" {
							return fmt.Errorf("session share_id is empty")
						}

						if sessionShareID != joinShareID {
							return fmt.Errorf("share_id mismatch: session=%s, join=%s",
								sessionShareID, joinShareID)
						}

						return nil
					},
				),
			},
			// Import Test
			{
				ResourceName:                         "kasm_session_permission.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "kasm_id",
				ImportStateVerifyIgnore: []string{
					"id",
					"user_permissions.#",
					"user_permissions.0.access",
					"user_permissions.0.user_id",
				},
			},
		},
	})
}

func testAccKasmSessionPermissionConfig_userSpecific(username1, username2, groupname, imageID string) string {
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

# Create group with sharing permissions
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Group for session testing"
    permissions = ["share_sessions", "allow_kasm_sharing", "shared_session_full_control"]
}

# Add first user to group
resource "kasm_group_membership" "test1" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test1.id
}

# Add second user to group
resource "kasm_group_membership" "test2" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test2.id
}

# Add image to group
resource "kasm_group_image" "test" {
    group_id = kasm_group.test.id
    image_id = "%s"
}

# Create session
resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test1]
    image_id = "%s"
    user_id = kasm_user.test1.id
    share = true
    enable_sharing = true
}

# Set session permissions for second user
resource "kasm_session_permission" "test" {
    depends_on = [kasm_session.test]
    kasm_id = kasm_session.test.id
    user_permissions = [
        {
            user_id = kasm_user.test2.id
            access = "rw"
        }
    ]
}

# Have second user join the session
resource "kasm_join" "test" {
    depends_on = [kasm_session_permission.test]
    share_id = kasm_session.test.share_id
    user_id = kasm_user.test2.id
}
`, testutils.ProviderConfig(), username1, username2, groupname, imageID, imageID)
}
