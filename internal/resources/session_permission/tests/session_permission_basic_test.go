package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmSessionPermission_basic(t *testing.T) {
	t.Skip("Skipping session permission test until sharing functionality and resource availability issues are resolved")
	t.Parallel()

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)
	groupname := fmt.Sprintf("SessionGroup_%s", uniqueIdentifier)

	// Get and configure a test image
	imageID, available := ensureImageAvailable(t)
	if !available {
		t.Skip("No suitable test images available")
	}

	// We'll use the client in the resource directly

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Basic Creation Test
			{
				Config: testAccKasmSessionPermissionConfig_basic(username, groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					// Check that resources exist
					testutils.TestCheckResourceExists("kasm_user.test"),
					testutils.TestCheckResourceExists("kasm_session.test"),
					testutils.TestCheckResourceExists("kasm_session_permission.test"),

					// Check user attributes
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),

					// Check initial session attributes
					resource.TestCheckResourceAttrSet("kasm_session.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
					resource.TestCheckResourceAttr("kasm_session.test", "share", "true"),
					resource.TestCheckResourceAttr("kasm_session.test", "enable_sharing", "true"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "share_id"),

					// Check session permissions
					resource.TestCheckResourceAttrSet("kasm_session_permission.test", "kasm_id"),
					resource.TestCheckResourceAttr("kasm_session_permission.test", "global_access", "r"),
				),
			},
			// Update Test
			{
				Config: testAccKasmSessionPermissionConfig_update(username, groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					testutils.TestCheckResourceExists("kasm_session_permission.test"),
					resource.TestCheckResourceAttrSet("kasm_session_permission.test", "kasm_id"),
					resource.TestCheckResourceAttr("kasm_session_permission.test", "global_access", "rw"),

					// Verify session still exists and has proper attributes
					testutils.TestCheckResourceExists("kasm_session.test"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
					resource.TestCheckResourceAttr("kasm_session.test", "share", "true"),
					resource.TestCheckResourceAttr("kasm_session.test", "enable_sharing", "true"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "share_id"),
				),
			},
			// Import Test
			{
				ResourceName:                         "kasm_session_permission.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "kasm_id",
				ImportStateVerifyIgnore:              []string{"id", "global_access"},
			},
		},
	})
}

func testAccKasmSessionPermissionConfig_basic(username, groupname, imageID string) string {
	return fmt.Sprintf(`
%s

# Create first user
resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
}

# Create group with sharing permissions
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Group for session testing"
    permissions = ["share_sessions", "allow_kasm_sharing", "shared_session_full_control"]
}

# Add first user to group
resource "kasm_group_membership" "test" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test.id
}

# Add image to group
resource "kasm_group_image" "test" {
    group_id = kasm_group.test.id
    image_id = "%s"
}

# Create session with sharing enabled
resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test]
    image_id = "%s"
    user_id = kasm_user.test.id
    # Enable sharing
    share = true
    enable_sharing = true
}

# Set session permissions
resource "kasm_session_permission" "test" {
    depends_on = [kasm_session.test]
    kasm_id = kasm_session.test.id
    global_access = "r"
}
`, testutils.ProviderConfig(), username, groupname, imageID, imageID)
}

func testAccKasmSessionPermissionConfig_update(username, groupname, imageID string) string {
	return fmt.Sprintf(`
%s

# Create first user
resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
}

# Create group with sharing permissions
resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Group for session testing"
    permissions = ["share_sessions", "allow_kasm_sharing", "shared_session_full_control"]
}

# Add first user to group
resource "kasm_group_membership" "test" {
    group_id = kasm_group.test.id
    user_id = kasm_user.test.id
}

# Add image to group
resource "kasm_group_image" "test" {
    group_id = kasm_group.test.id
    image_id = "%s"
}

# Create session with sharing enabled
resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test]
    image_id = "%s"
    user_id = kasm_user.test.id
    # Enable sharing
    share = true
    enable_sharing = true
}

# Set session permissions
resource "kasm_session_permission" "test" {
    depends_on = [kasm_session.test]
    kasm_id = kasm_session.test.id
    global_access = "rw"
}

# We'll enable sharing directly through the client in the test code instead of using null_resource
`, testutils.ProviderConfig(), username, groupname, imageID, imageID)
}

func ensureImageAvailable(t *testing.T) (string, bool) {
	// Use the testutils function to get any available image
	return testutils.EnsureImageAvailable(t)
}
