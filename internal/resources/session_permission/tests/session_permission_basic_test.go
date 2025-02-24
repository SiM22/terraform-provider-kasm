package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmSessionPermission_basic(t *testing.T) {
	t.Parallel()

	// Generate unique identifiers with both timestamp and random number
	rand.Seed(time.Now().UnixNano())
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Intn(10000))
	username := fmt.Sprintf("testuser_%s", uniqueIdentifier)
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
				Config: testAccKasmSessionPermissionConfig_basic(username, groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					// Check that resources exist
					testutils.TestCheckResourceExists("kasm_user.test"),
					testutils.TestCheckResourceExists("kasm_session.test"),
					testutils.TestCheckResourceExists("kasm_session_permission.test"),

					// Check user attributes
					resource.TestCheckResourceAttr("kasm_user.test", "username", username),

					// Check session attributes
					resource.TestCheckResourceAttrSet("kasm_session.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "share_id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
					resource.TestCheckResourceAttr("kasm_session.test", "share", "true"),
					resource.TestCheckResourceAttr("kasm_session.test", "enable_sharing", "true"),

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
					resource.TestCheckResourceAttrSet("kasm_session.test", "share_id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
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

# Create session
resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test]
    image_id = "%s"
    user_id = kasm_user.test.id
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

# Create session
resource "kasm_session" "test" {
    depends_on = [kasm_group_image.test, kasm_group_membership.test]
    image_id = "%s"
    user_id = kasm_user.test.id
    share = true
    enable_sharing = true
}

# Set session permissions
resource "kasm_session_permission" "test" {
    depends_on = [kasm_session.test]
    kasm_id = kasm_session.test.id
    global_access = "rw"
}
`, testutils.ProviderConfig(), username, groupname, imageID, imageID)
}

func ensureImageAvailable(t *testing.T) (string, bool) {
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		c := testutils.GetTestClient(t)
		if c == nil {
			t.Logf("Attempt %d: Failed to get test client", i+1)
			return "", false
		}

		images, err := c.GetImages()
		if err != nil {
			t.Logf("Attempt %d: Error getting images: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return "", false
		}

		// Look for any Chrome test image
		for _, img := range images {
			// Check if it's a Chrome test image
			if img.Name == "kasmweb/chrome:1.16.0" && img.Enabled {
				// Keep all existing fields but update run_config
				runConfig := map[string]interface{}{
					"name":     "kasm_test_container",
					"hostname": "kasm-test",
					"network":  "kasm-network",
					"environment": map[string]string{
						"KASM_TEST": "true",
					},
				}

				// Create a copy of the image with all fields
				updatedImage := &client.Image{
					ImageID:             img.ImageID,
					Name:                img.Name,
					FriendlyName:        img.FriendlyName,
					Description:         img.Description,
					Categories:          img.Categories,
					Memory:              img.Memory,
					Cores:               img.Cores,
					CPUAllocationMethod: img.CPUAllocationMethod,
					DockerRegistry:      img.DockerRegistry,
					UncompressedSizeMB:  img.UncompressedSizeMB,
					ImageType:           img.ImageType,
					Enabled:             img.Enabled,
					Available:           img.Available,
					ImageSrc:            img.ImageSrc,
					ExecConfig:          img.ExecConfig,
					RunConfig:           runConfig,
					// Set volume mappings to nil since it's optional
					VolumeMappings: nil,
				}

				// Update the image
				_, err := c.UpdateImage(updatedImage)
				if err != nil {
					t.Logf("Failed to update image run_config: %v", err)
					return "", false
				}

				t.Logf("Found and updated test image: %s", img.ImageID)
				return img.ImageID, true
			}
		}

		t.Logf("Attempt %d: No suitable Chrome test image found", i+1)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	t.Logf("Failed to find suitable test image after %d attempts", maxRetries)
	return "", false
}
