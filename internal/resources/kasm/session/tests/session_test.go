//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-kasm/testutils"
)

func generateUniqueUsername() string {
	// Use timestamp and random number to ensure uniqueness
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("testuser_%d_%d", timestamp, randomNum)
}

func TestAccKasmSession_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	// Clean up any existing sessions before starting the test
	testutils.CleanupExistingSessions(t)

	// Check for available images
	imageID, available := testutils.EnsureImageAvailable(t)
	if !available {
		t.Skip("Skipping test as no suitable test images are available")
	}

	username := generateUniqueUsername()
	groupname := fmt.Sprintf("testgroup_%d", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
			// Clean up any existing sessions before each test run
			testutils.CleanupExistingSessions(t)
			// Double check images are still available
			if _, available := testutils.EnsureImageAvailable(t); !available {
				t.Skip("Images are no longer available")
			}
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmSessionConfig_basic(username, groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("kasm_session.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "user_id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
					resource.TestCheckResourceAttr("kasm_session.test", "share", "false"),
					resource.TestCheckResourceAttr("kasm_session.test", "enable_sharing", "false"),
					resource.TestCheckResourceAttr("kasm_session.test", "persistent", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			// Clean up test image if it was created
			testutils.CleanupTestImage(t)
			return nil
		},
	})
}

func testAccKasmSessionConfig_basic(username, groupname, imageID string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_group" "test" {
    name = "%s"
    description = "Test group for session testing"
    priority = 100
    permissions = ["allow_all_images", "allow_sharing"]
}

resource "kasm_group_image" "test" {
    depends_on = [kasm_group.test]
    group_id = kasm_group.test.id
    image_id = "%s"  # Test image
}

resource "kasm_user" "test" {
    depends_on = [kasm_group.test]
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
    organization = "Test Org"
    locked = false
    disabled = false
    groups = []
}

resource "kasm_group_membership" "test" {
    depends_on = [kasm_group.test, kasm_user.test]
    group_id = kasm_group.test.id
    user_id = kasm_user.test.id
}

resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test]
    image_id = kasm_group_image.test.image_id
    user_id = kasm_user.test.id
    share = false
    enable_sharing = false
    persistent = true
    allow_resume = true
    session_authentication = false
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), groupname, imageID, username)
}
