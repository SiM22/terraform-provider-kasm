package tests

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func generateUniqueUsername() string {
	// Use timestamp and random number to ensure uniqueness
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("testuser_%d_%d", timestamp, randomNum)
}

func cleanupExistingSessions(t *testing.T) {
	// Get the client
	c := testutils.GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
	}

	// Get all sessions
	kasms, err := c.GetKasms()
	if err != nil {
		t.Logf("Warning: Failed to get existing sessions: %v", err)
		return
	}

	// Destroy each session
	for _, kasm := range kasms.Kasms {
		err := c.DestroyKasm(kasm.UserID, kasm.KasmID)
		if err != nil {
			t.Logf("Warning: Failed to destroy session %s: %v", kasm.KasmID, err)
		}
	}
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

func waitForGroupMembership(t *testing.T, userID string, groupID string) error {
	maxRetries := 30
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		c := testutils.GetTestClient(t)
		if c == nil {
			return fmt.Errorf("failed to get test client")
		}

		user, err := c.GetUser(userID)
		if err != nil {
			t.Logf("Attempt %d: Error getting user: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to get user after %d attempts: %v", maxRetries, err)
		}

		// Check if the target group is in the list
		for _, group := range user.Groups {
			if group.GroupID == groupID {
				t.Logf("Successfully found group membership after %d attempts", i+1)
				return nil
			}
		}

		t.Logf("Attempt %d: Group %s not found in user's groups. Waiting...", i+1, groupID)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("group membership not found after %d attempts (userID: %s, groupID: %s)", maxRetries, userID, groupID)
}

func TestAccKasmSession_Basic(t *testing.T) {
	// Clean up any existing sessions before starting the test
	cleanupExistingSessions(t)

	// Check for available images
	imageID, available := ensureImageAvailable(t)
	if !available {
		t.Skip("Skipping test as no suitable Chrome test images are available")
	}

	username := generateUniqueUsername()
	groupname := fmt.Sprintf("testgroup_%d", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
			// Clean up any existing sessions before each test run
			cleanupExistingSessions(t)
			// Double check images are still available
			if _, available := ensureImageAvailable(t); !available {
				t.Skip("Images are no longer available")
			}
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// Wait for group membership to propagate before creating the session
					time.Sleep(30 * time.Second)
					// Verify images are available before running the test step
					if _, available := ensureImageAvailable(t); !available {
						t.Fatal("No images available before running test step")
					}
				},
				Config: testAccKasmSessionConfig_basic(username, groupname, imageID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("kasm_session.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_session.test", "operational_status"),
					func(state *terraform.State) error {
						// Get the resource
						rs, ok := state.RootModule().Resources["kasm_session.test"]
						if !ok {
							return fmt.Errorf("not found: %s", "kasm_session.test")
						}

						if rs.Primary.ID == "" {
							return fmt.Errorf("no ID is set")
						}

						return nil
					},
				),
			},
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
    permissions = ["allow_all_images"]
}

resource "kasm_group_image" "test" {
    depends_on = [kasm_group.test]
    group_id = kasm_group.test.id
    image_id = "%s"  # Chrome test image
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
    share = true
    enable_sharing = true
    persistent = false
    allow_resume = false
    session_authentication = false
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), groupname, imageID, username)
}
