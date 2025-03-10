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
		Description:        "Test image for frame stats tests",
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
	maxRetries := 30 // 5 minutes max wait time
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("[DEBUG] Checking if image %s is available (attempt %d/%d)", imageID, i+1, maxRetries)

		images, err := c.GetImages()
		if err != nil {
			log.Printf("[DEBUG] Error checking image status: %v", err)
			time.Sleep(retryInterval)
			continue
		}

		for _, img := range images {
			if img.ImageID == imageID {
				if img.Available {
					log.Printf("[DEBUG] Image %s is now available", imageID)
					return
				}
				log.Printf("[DEBUG] Image %s is not yet available, waiting...", imageID)
				break
			}
		}

		time.Sleep(retryInterval)
	}

	log.Printf("[WARN] Timed out waiting for image %s to become available", imageID)
}

// TestAccKasmStats_FrameStats is an acceptance test for the kasm_stats resource
// that tests the ability to get frame stats from a Kasm session.
// This test requires manual interaction to open the Kasm session URL in a browser
// unless KASM_SKIP_MANUAL_PROMPT=true is set.
func TestAccKasmStats_FrameStats(t *testing.T) {
	// Skip the test if TF_ACC is not set
	testutils.TestAccPreCheck(t)

	// Skip the test if KASM_SKIP_BROWSER_TEST is set to true
	// This allows skipping the test in CI environments where manual interaction isn't possible
	if os.Getenv("KASM_SKIP_BROWSER_TEST") == "true" {
		t.Skip("Skipping test that requires browser interaction as KASM_SKIP_BROWSER_TEST=true")
	}

	// Skip the test if CI is set to true
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test that requires browser interaction in CI environment")
	}

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
		if client.IsResourceUnavailableError(err) {
			t.Skip("Skipping test as no resources are available to create a Kasm session. An active session is required for this test to work.")
		}
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

	// Get the share ID from the response
	shareID := kasm.ShareID
	if shareID == "" {
		t.Fatalf("Share ID is empty")
	}
	log.Printf("[DEBUG] Using share ID: %s", shareID)

	// Get the join URL for the session using the share_id
	log.Printf("[DEBUG] Getting join URL for the session...")
	joinResp, err := c.JoinKasm(shareID, userID)
	if err != nil {
		t.Fatalf("Failed to get join URL: %v", err)
	} else {
		// Construct the full URL to join the session
		joinURL := fmt.Sprintf("%s%s", os.Getenv("KASM_BASE_URL"), joinResp.KasmURL)

		// Print the join URL prominently for manual connection
		fmt.Println("\n===========================================================")
		fmt.Println("MANUAL TESTING REQUIRED")
		fmt.Println("===========================================================")
		fmt.Println("PLEASE OPEN THIS URL IN YOUR BROWSER TO CONNECT TO THE SESSION:")
		fmt.Println(joinURL)
		fmt.Println("===========================================================")
		fmt.Println("IMPORTANT: After opening the URL, please:")
		fmt.Println("1. Wait for the session to fully load")
		fmt.Println("2. Interact with the session (click, type, move mouse)")
		fmt.Println("3. Keep the browser window open until the test completes")
		fmt.Println("===========================================================")

		// Check if we should skip the manual connection prompt
		skipManualPrompt := os.Getenv("KASM_SKIP_MANUAL_PROMPT") == "true"
		if !skipManualPrompt {
			// Wait for manual connection
			fmt.Println("Waiting 60 seconds for manual connection...")
			time.Sleep(60 * time.Second)
		} else {
			fmt.Println("Skipping manual connection prompt as KASM_SKIP_MANUAL_PROMPT=true")
			fmt.Println("WARNING: Test may fail without manual browser connection")
		}
	}

	// Try to get frame stats a few times before proceeding with the test
	log.Printf("[DEBUG] Attempting to get frame stats...")
	var frameStats *client.FrameStatsResponse
	maxStatsRetries := 10
	statsRetryInterval := 5 * time.Second
	for i := 0; i < maxStatsRetries; i++ {
		stats, err := c.GetFrameStats(kasm.KasmID, userID)
		if err == nil && stats != nil {
			frameStats = stats
			log.Printf("[DEBUG] Successfully retrieved frame stats on attempt %d", i+1)
			break
		}
		log.Printf("[DEBUG] Attempt %d: Failed to get frame stats: %v. Retrying in %s...", i+1, err, statsRetryInterval)
		time.Sleep(statsRetryInterval)
	}

	if frameStats == nil {
		log.Printf("[WARN] Could not get frame stats after %d attempts. Test may fail.", maxStatsRetries)
		// Skip the Terraform test if we couldn't get frame stats
		t.Skip("Skipping Terraform test as frame stats could not be retrieved. Please ensure you manually connected to the Kasm session.")
	} else {
		// Print the frame stats
		log.Printf("[DEBUG] Frame stats retrieved successfully")
		log.Printf("[DEBUG] Frame stats: %+v", frameStats)
	}

	// Run the Terraform acceptance test
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmStatsConfig(kasm.KasmID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("kasm_stats.test", "kasm_id", kasm.KasmID),
					resource.TestCheckResourceAttr("kasm_stats.test", "user_id", userID),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "fps"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "latency_avg"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "latency_max"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "latency_min"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "latency_mdev"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "tx_bps"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "rx_bps"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "tx_bytes"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "rx_bytes"),
				),
			},
		},
	})
}

// testAccKasmStatsConfig creates the Terraform configuration for testing the kasm_stats resource
func testAccKasmStatsConfig(kasmID, userID string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_stats" "test" {
    kasm_id = "%s"
    user_id = "%s"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), kasmID, userID)
}
