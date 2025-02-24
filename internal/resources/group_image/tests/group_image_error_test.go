package tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"terraform-provider-kasm/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupImage_errors(t *testing.T) {
	// Remove parallel execution since we're dealing with shared resources
	// t.Parallel()

	// Get client and ensure test image exists
	client := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, client)
	if imageID == "" {
		t.Fatal("Failed to get a valid image ID")
	}

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		if imageID != "" {
			cleanupTestImage(t, client)
		}
	})

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccGroupImageConfig_nonExistentImage(groupname),
				ExpectError: regexp.MustCompile(`Image with ID non-existent-image-id does not exist`),
			},
			{
				Config:      testAccGroupImageConfig_nonExistentGroup(imageID),
				ExpectError: regexp.MustCompile(`Error authorizing image for group`),
			},
			{
				Config: testAccGroupImageConfig_restrictedGroup(groupname, username, imageID),
				ExpectError: regexp.MustCompile(
					`Error: Error assigning user to groups\s+with kasm_user\.test,\s+on terraform_plugin_test\.tf line \d+, in resource "kasm_user" "test":\s+\d+: resource "kasm_user" "test" {\s+group not found: [a-f0-9]+`),
			},
		},
	})
}

func testAccGroupImageConfig_nonExistentImage(groupname string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Test group for acceptance tests"
    permissions = ["allow_all_images"]
}

resource "kasm_group_image" "test" {
    depends_on = [kasm_group.test]
    group_id = kasm_group.test.id
    image_id = "non-existent-image-id"
}
`, testutils.ProviderConfig(), groupname)
}

func testAccGroupImageConfig_nonExistentGroup(imageID string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group_image" "test" {
    group_id = "non-existent-group-id"
    image_id = "%s"
}
`, testutils.ProviderConfig(), imageID)
}

func testAccGroupImageConfig_restrictedGroup(groupname, username, imageID string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name = "%s"
    priority = 1
    description = "Test group for acceptance tests"
}

resource "kasm_user" "test" {
    depends_on = [kasm_group.test]
    username = "%s"
    password = "testpassword123!"
    first_name = "Test"
    last_name = "User"
    groups = [kasm_group.test.id]
}

resource "kasm_group_image" "test" {
    depends_on = [kasm_group.test, kasm_user.test]
    group_id = kasm_group.test.id
    image_id = "%s"
}

resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test]
    image_id = kasm_group_image.test.image_id
    user_id = kasm_user.test.id
}
`, testutils.ProviderConfig(), groupname, username, imageID)
}
