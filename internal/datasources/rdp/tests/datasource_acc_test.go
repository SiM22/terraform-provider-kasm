//go:build acceptance
// +build acceptance

package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"
)

var createdImageID string // Track the created image ID for cleanup

// cleanupTestImage deletes the test image if it was created
func cleanupTestImage(t *testing.T, c *client.Client) {
	if createdImageID != "" {
		if err := c.DeleteImage(createdImageID); err != nil {
			t.Logf("Warning: Failed to delete test image %s: %v", createdImageID, err)
		} else {
			log.Printf("[DEBUG] Successfully deleted test image: %s", createdImageID)
			createdImageID = "" // Reset after successful deletion
		}
	}
}

// ensureWorkspaceImage ensures a test image exists and returns its ID
func ensureWorkspaceImage(t *testing.T, c *client.Client) string {
	log.Printf("[DEBUG] Starting ensureWorkspaceImage")

	// If we already created an image in this test run, return its ID
	if createdImageID != "" {
		return createdImageID
	}

	// Try to find an existing image first
	images, err := c.GetImages()
	if err != nil {
		log.Printf("[DEBUG] Error getting images: %v", err)
		log.Printf("[DEBUG] Will attempt to create new image")
	} else {
		log.Printf("[DEBUG] Found %d existing images", len(images))
		for i, img := range images {
			log.Printf("[DEBUG] Image %d: ID=%s Name=%s Available=%v", i, img.ImageID, img.Name, img.Available)
		}
		// Look for any usable image
		for _, img := range images {
			if img.ImageID != "" && img.Available {
				log.Printf("[DEBUG] Found existing usable image: %s (%s)", img.Name, img.ImageID)
				return img.ImageID
			}
		}
		log.Printf("[DEBUG] No usable image found, will create new one")
	}

	// Create a small, lightweight image for testing
	runConfig := map[string]interface{}{
		"hostname":       "kasm-test",
		"container_name": "kasm_test_container",
		"network":        "kasm-network",
		"environment": map[string]string{
			"KASM_TEST": "true",
		},
	}
	runConfigJSON, _ := json.Marshal(runConfig)

	execConfig := map[string]interface{}{}
	execConfigJSON, _ := json.Marshal(execConfig)

	volumeMappings := map[string]interface{}{}
	volumeMappingsJSON, _ := json.Marshal(volumeMappings)

	// Using a small image for faster download
	image := &client.CreateImageRequest{
		ImageSrc:           "img/thumbnails/filezilla.png",
		Categories:         "Testing",
		RunConfig:          string(runConfigJSON),
		Description:        "Test image for RDP connection tests",
		FriendlyName:       "FileZilla_Test",
		DockerRegistry:     "https://index.docker.io/v1/",
		Name:               "kasmweb/filezilla:1.16.1",
		UncompressedSizeMB: 1000,
		ImageType:          "Container",
		Enabled:            true,
		Memory:             1024000000,
		Cores:              1,
		GPUCount:           0,
		RequireGPU:         nil,
		ExecConfig:         string(execConfigJSON),
		VolumeMappings:     string(volumeMappingsJSON),
	}

	log.Printf("[DEBUG] Creating test image with config: %+v", image)

	var createdImage *client.Image
	var lastErr error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		log.Printf("[DEBUG] Attempt %d/%d to create image", i+1, maxRetries)

		createdImage, lastErr = c.AddWorkspaceImage(image)

		if lastErr != nil {
			log.Printf("[DEBUG] Attempt %d failed with error: %v", i+1, lastErr)
			if i < maxRetries-1 {
				sleepDuration := time.Second * time.Duration(i+1)
				log.Printf("[DEBUG] Waiting %v before next attempt", sleepDuration)
				time.Sleep(sleepDuration)
			}
			continue
		}

		if createdImage == nil {
			lastErr = fmt.Errorf("API returned success but image is nil")
			log.Printf("[DEBUG] Attempt %d failed: %v", i+1, lastErr)
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		if createdImage.ImageID == "" {
			lastErr = fmt.Errorf("API returned image but ImageID is empty")
			log.Printf("[DEBUG] Attempt %d failed: %v", i+1, lastErr)
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		// Success! Store the ID for cleanup
		createdImageID = createdImage.ImageID
		log.Printf("[DEBUG] Successfully created test image with ID: %s", createdImageID)
		log.Printf("[DEBUG] Image details: %+v", createdImage)

		// Wait for the image to be downloaded and available
		log.Printf("[DEBUG] Waiting for image to be downloaded and available...")
		waitForImageAvailable(t, c, createdImageID)

		return createdImageID
	}

	errorMsg := fmt.Sprintf("Failed to create workspace image after %d attempts. Last error: %v", maxRetries, lastErr)
	log.Printf("[DEBUG] %s", errorMsg)
	t.Fatal(errorMsg)
	return "" // unreachable
}

// waitForImageAvailable waits for an image to be downloaded and available
func waitForImageAvailable(t *testing.T, c *client.Client, imageID string) {
	maxRetries := 5
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		images, err := c.GetImages()
		if err != nil {
			log.Printf("[DEBUG] Error getting images: %v", err)
		} else {
			for _, img := range images {
				if img.ImageID == imageID {
					if img.Available {
						log.Printf("[DEBUG] Image %s is now available", imageID)
						return
					}
					log.Printf("[DEBUG] Image %s is not yet available, status: %v", imageID, img.Available)
					break
				}
			}
		}

		log.Printf("[DEBUG] Waiting for image to be available, attempt %d/%d. Sleeping for %v...", i+1, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}

	log.Printf("[WARN] Image %s did not become available after %d attempts", imageID, maxRetries)
}

// TestAccKasmRDPClientConnectionInfo tests the RDP client connection info data source
func TestAccKasmRDPClientConnectionInfo(t *testing.T) {
	t.Skip("Skipping test as the RDP client connection info API endpoint is not working as expected")
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	// Get a test client
	c := testutils.GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
	}

	// Ensure we have an image for testing
	imageID, available := testutils.EnsureImageAvailable(t)
	if !available {
		t.Skip("Skipping test as no suitable test images are available")
	}

	// Create a test user
	username := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])
	user := &client.User{
		Username:     username,
		Password:     "Test@123",
		FirstName:    "Test",
		LastName:     "User",
		Organization: "test@example.com",
		Locked:       false,
		Disabled:     false,
	}

	createdUser, err := c.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	userID := createdUser.UserID
	defer func() {
		if err := c.DeleteUser(userID); err != nil {
			t.Logf("Warning: Failed to delete test user: %v", err)
		}
	}()

	// Create a test session using the existing CreateKasm method
	sessionToken := uuid.New().String()
	session, err := c.CreateKasm(userID, imageID, sessionToken, username, false, true, false, false)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	kasmID := session.KasmID
	defer func() {
		if err := c.DestroyKasm(userID, kasmID); err != nil {
			t.Logf("Warning: Failed to destroy test session: %v", err)
		}
		// Clean up test image if it was created
		testutils.CleanupTestImage(t)
	}()

	// Wait for session to be ready
	maxRetries := 10
	retryDelay := 2 * time.Second
	var sessionReady bool

	for i := 0; i < maxRetries; i++ {
		status, err := c.GetKasmStatus(userID, kasmID, true)
		if err != nil {
			t.Logf("Warning: Failed to get session status (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		// Check both the status.OperationalStatus and status.Kasm.OperationalStatus
		operationalStatus := status.OperationalStatus
		if status.Kasm != nil && status.Kasm.OperationalStatus != "" {
			operationalStatus = status.Kasm.OperationalStatus
		}

		t.Logf("Session status: %s (attempt %d/%d)", operationalStatus, i+1, maxRetries)

		if operationalStatus == "running" {
			sessionReady = true
			break
		}

		time.Sleep(retryDelay)
	}

	if !sessionReady {
		t.Fatalf("Timed out waiting for session to be ready")
	}

	// Run the tests
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRDPClientConnectionInfoFileConfig(userID, kasmID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "user_id", userID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "kasm_id", kasmID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "connection_type", "file"),
					resource.TestCheckResourceAttrSet("data.kasm_rdp_client_connection_info.test", "file"),
				),
			},
			{
				Config: testAccKasmRDPClientConnectionInfoURLConfig(userID, kasmID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "user_id", userID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "kasm_id", kasmID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test", "connection_type", "url"),
					resource.TestCheckResourceAttrSet("data.kasm_rdp_client_connection_info.test", "url"),
				),
			},
		},
	})
}

// Test configurations
func testAccKasmRDPClientConnectionInfoFileConfig(userID, kasmID string) string {
	return testutils.ProviderConfig() + fmt.Sprintf(`
data "kasm_rdp_client_connection_info" "test" {
  user_id = "%s"
  kasm_id = "%s"
  connection_type = "file"
}
`, userID, kasmID)
}

func testAccKasmRDPClientConnectionInfoURLConfig(userID, kasmID string) string {
	return testutils.ProviderConfig() + fmt.Sprintf(`
data "kasm_rdp_client_connection_info" "test" {
  user_id = "%s"
  kasm_id = "%s"
  connection_type = "url"
}
`, userID, kasmID)
}
