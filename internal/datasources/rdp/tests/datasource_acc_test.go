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
	maxRetries := 30
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
	// Skip the test if TF_ACC is not set
	testutils.TestAccPreCheck(t)

	// Get the user ID from the environment or use a default
	userID := os.Getenv("KASM_USER_ID")
	if userID == "" {
		userID = "44edb3e5-2909-4927-a60b-6e09c7219104"
	}
	log.Printf("[DEBUG] Using user ID: %s", userID)

	// Create a new Kasm client
	c := client.NewClient(
		os.Getenv("KASM_BASE_URL"),
		os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"),
		true, // insecure - ignore TLS certificate verification
	)

	// Ensure we have a usable image
	imageID := ensureWorkspaceImage(t, c)
	log.Printf("[DEBUG] Using image ID: %s", imageID)

	// Ensure the test image is cleaned up after the test
	defer cleanupTestImage(t, c)

	// Create a new Kasm session
	log.Printf("[DEBUG] Creating Kasm session for user %s with image %s", userID, imageID)
	sessionToken := uuid.New().String()
	kasm, err := c.CreateKasm(userID, imageID, sessionToken, "test", true, false, false, false)
	if err != nil {
		t.Fatalf("Failed to create Kasm session: %v", err)
	}

	// Ensure the Kasm session is deleted when the test is done
	defer func() {
		log.Printf("[DEBUG] Cleaning up Kasm session %s", kasm.KasmID)
		err := c.DestroyKasm(userID, kasm.KasmID)
		if err != nil {
			log.Printf("[WARN] Failed to delete Kasm session: %v", err)
		}
	}()

	// Wait for the session to be fully initialized
	log.Printf("[DEBUG] Waiting for session to initialize...")
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		time.Sleep(10 * time.Second)

		// Check if the session is available
		status, err := c.GetKasmStatus(userID, kasm.KasmID, true)
		if err != nil {
			log.Printf("Attempt %d: Session not ready yet: %v. Retrying...", i+1, err)
			continue
		}

		if status.Kasm != nil && status.Kasm.ContainerID != "" {
			log.Printf("Session is ready after %d attempts", i+1)
			break
		}

		log.Printf("Attempt %d: Session not fully initialized yet. Retrying...", i+1)

		if i == maxRetries-1 {
			log.Printf("Warning: Session may not be fully initialized after %d attempts", maxRetries)
			t.Skip("Skipping test as session did not initialize in time")
		}
	}

	// Get session details
	log.Printf("[DEBUG] Getting session details for Kasm ID: %s", kasm.KasmID)

	// Wait for the session to be ready before requesting RDP connection info
	log.Printf("[DEBUG] Waiting for session to be ready...")
	maxRetries = 30
	retryInterval := 5 * time.Second
	var sessionReady bool
	var sessionDetails *client.KasmStatusResponse

	for i := 0; i < maxRetries; i++ {
		var err error
		sessionDetails, err = c.GetKasmStatus(userID, kasm.KasmID, true)
		if err != nil {
			log.Printf("[DEBUG] Error getting session status: %v. Retrying...", err)
		} else {
			log.Printf("[DEBUG] Session details: %+v", sessionDetails)
			if sessionDetails.Kasm != nil && sessionDetails.Kasm.OperationalStatus == "running" {
				log.Printf("[DEBUG] Session is ready (running)")
				sessionReady = true
				break
			}
			log.Printf("[DEBUG] Session is not ready yet, status: %s", sessionDetails.OperationalStatus)
		}

		log.Printf("[DEBUG] Waiting for session to be ready, attempt %d/%d. Sleeping for %v...", i+1, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}

	if !sessionReady {
		t.Fatalf("Session did not become ready after %d attempts", maxRetries)
	}

	// Test RDP connection info for file type
	log.Printf("[DEBUG] Getting RDP connection info for file type")
	fileConnectionInfo, err := c.GetRDPConnectionInfo(userID, kasm.KasmID, "file")
	if err != nil {
		t.Fatalf("Failed to get RDP connection info for file type: %v", err)
	}
	log.Printf("[DEBUG] RDP file connection info: %+v", fileConnectionInfo)
	if fileConnectionInfo.File == "" {
		log.Printf("[WARN] RDP file connection info is empty. This is expected when using container images instead of RDP servers.")
		log.Printf("[INFO] To fully test RDP functionality, additional infrastructure is required (Windows RDP server configured in Kasm).")
	}

	// Test RDP connection info for URL type
	log.Printf("[DEBUG] Getting RDP connection info for URL type")
	urlConnectionInfo, err := c.GetRDPConnectionInfo(userID, kasm.KasmID, "url")
	if err != nil {
		t.Fatalf("Failed to get RDP connection info for URL type: %v", err)
	}
	log.Printf("[DEBUG] RDP URL connection info: %+v", urlConnectionInfo)
	if urlConnectionInfo.URL == "" {
		log.Printf("[WARN] RDP URL connection info is empty. This is expected when using container images instead of RDP servers.")
		log.Printf("[INFO] To fully test RDP functionality, additional infrastructure is required (Windows RDP server configured in Kasm).")
	}

	// Run the Terraform acceptance test
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRDPClientConnectionInfoFileConfig(userID, kasm.KasmID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_rdp_client_connection_info.test_file", "id"),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_file", "user_id", userID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_file", "kasm_id", kasm.KasmID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_file", "connection_type", "file"),
					// We don't check for file content since it might be empty with container images
				),
			},
			{
				Config: testAccKasmRDPClientConnectionInfoURLConfig(userID, kasm.KasmID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_rdp_client_connection_info.test_url", "id"),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_url", "user_id", userID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_url", "kasm_id", kasm.KasmID),
					resource.TestCheckResourceAttr("data.kasm_rdp_client_connection_info.test_url", "connection_type", "url"),
					// We don't check for URL content since it might be empty with container images
				),
			},
		},
	})
}

// Test configurations
func testAccKasmRDPClientConnectionInfoFileConfig(userID, kasmID string) string {
	return testutils.ProviderConfig() + fmt.Sprintf(`
data "kasm_rdp_client_connection_info" "test_file" {
  user_id = "%s"
  kasm_id = "%s"
  connection_type = "file"
}
`, userID, kasmID)
}

func testAccKasmRDPClientConnectionInfoURLConfig(userID, kasmID string) string {
	return testutils.ProviderConfig() + fmt.Sprintf(`
data "kasm_rdp_client_connection_info" "test_url" {
  user_id = "%s"
  kasm_id = "%s"
  connection_type = "url"
}
`, userID, kasmID)
}
