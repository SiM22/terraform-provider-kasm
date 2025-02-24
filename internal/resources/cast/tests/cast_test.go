//go:build acceptance
// +build acceptance

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
	"os"
	"regexp"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"
	"testing"
	"time"
)

var createdImageID string // Track the created image ID for cleanup

// cleanupTestImage deletes the test image if it was created
func cleanupTestImage(t *testing.T, c *client.Client) {
	if createdImageID != "" {
		if err := c.DeleteImage(createdImageID); err != nil {
			t.Logf("Warning: Failed to delete test image %s: %v", createdImageID, err)
		} else {
			if os.Getenv("KASM_DEBUG") != "" {
				log.Printf("[DEBUG] Successfully deleted test image: %s", createdImageID)
			}
			createdImageID = "" // Reset after successful deletion
		}
	}
}

func TestAccKasmCastConfig_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	c := testutils.GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
	}

	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] Starting TestAccKasmCastConfig_basic")
	}

	imageID := ensureWorkspaceImage(t, c)
	if imageID == "" {
		t.Fatal("Failed to get test image ID")
	}

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, c)
	})

	uniqueIdentifier := fmt.Sprintf("%d", time.Now().UnixNano())
	rName := fmt.Sprintf("tf-test-cast-%s", uniqueIdentifier)
	uniqueKey := fmt.Sprintf("key-%s", uniqueIdentifier)
	resourceName := "kasm_cast_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmCastConfigConfig_basic(rName, uniqueKey, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "key", uniqueKey),
					resource.TestCheckResourceAttr(resourceName, "image_id", imageID),
					resource.TestCheckResourceAttr(resourceName, "enable_sharing", "false"),
					resource.TestCheckResourceAttr(resourceName, "disable_control_panel", "false"),
					resource.TestCheckResourceAttr(resourceName, "allowed_referrers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "allowed_referrers.0", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "allowed_referrers.1", "test.com"),
					resource.TestCheckResourceAttr(resourceName, "limit_sessions", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_remaining", "10"),
					resource.TestCheckResourceAttr(resourceName, "error_url", "https://error.example.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] Completed TestAccKasmCastConfig_basic")
	}
}

func TestAccKasmCastConfig_update(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	c := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, c)

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, c)
	})

	uniqueIdentifier := fmt.Sprintf("%d", time.Now().UnixNano())
	rName := fmt.Sprintf("tf-test-cast-%s", uniqueIdentifier)
	updatedName := fmt.Sprintf("%s-updated", rName)
	uniqueKey := fmt.Sprintf("key-%s", uniqueIdentifier)
	resourceName := "kasm_cast_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmCastConfigConfig_basic(rName, uniqueKey, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccKasmCastConfigConfig_update(updatedName, uniqueKey, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "enable_sharing", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_anonymous", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_remaining", "20"),
					resource.TestCheckResourceAttr(resourceName, "allowed_referrers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ip_request_limit", "10"),
					resource.TestCheckResourceAttr(resourceName, "ip_request_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "error_url", "https://error-updated.example.com"),
					resource.TestCheckResourceAttr(resourceName, "allow_kasm_audio", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_kasm_uploads", "true"),
				),
			},
		},
	})
}

func TestAccKasmCastConfig_validation(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	c := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, c)

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, c)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                    provider "kasm" {
                        base_url = "%s"
                        api_key = "%s"
                        api_secret = "%s"
                        insecure = true
                    }
                    resource "kasm_cast_config" "test" {
                        name     = "test-invalid-key"
                        image_id = "%s"
                        key      = "invalid@key$"
                    }`,
					os.Getenv("KASM_BASE_URL"),
					os.Getenv("KASM_API_KEY"),
					os.Getenv("KASM_API_SECRET"),
					imageID,
				),
				ExpectError: regexp.MustCompile(`key must only contain alphanumeric characters`),
			},
			{
				Config: fmt.Sprintf(`
                    provider "kasm" {
                        base_url = "%s"
                        api_key = "%s"
                        api_secret = "%s"
                        insecure = true
                    }
                    resource "kasm_cast_config" "test" {
                        name      = "test-invalid-url"
                        image_id  = "%s"
                        key       = "validkey123"
                        error_url = "not-a-url"
                    }`,
					os.Getenv("KASM_BASE_URL"),
					os.Getenv("KASM_API_KEY"),
					os.Getenv("KASM_API_SECRET"),
					imageID,
				),
				ExpectError: regexp.MustCompile(`URL must start with http:// or https://`),
			},
		},
	})
}

func TestAccKasmCastConfig_fullConfig(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	c := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, c)

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, c)
	})

	uniqueIdentifier := fmt.Sprintf("%d", time.Now().UnixNano())
	rName := fmt.Sprintf("tf-test-%s", uniqueIdentifier)
	uniqueKey := fmt.Sprintf("key-%s", uniqueIdentifier)
	resourceName := "kasm_cast_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmCastConfigConfig_full(rName, uniqueKey, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "limit_sessions", "true"),
					resource.TestCheckResourceAttr(resourceName, "session_remaining", "10"),
					resource.TestCheckResourceAttr(resourceName, "limit_ips", "true"),
					resource.TestCheckResourceAttr(resourceName, "ip_request_limit", "5"),
					resource.TestCheckResourceAttr(resourceName, "ip_request_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceName, "enforce_client_settings", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_kasm_audio", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_kasm_uploads", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_kasm_downloads", "false"),
				),
			},
		},
	})
}

func ensureWorkspaceImage(t *testing.T, c *client.Client) string {
	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] Starting ensureWorkspaceImage")
	}

	// If we already created an image in this test run, return its ID
	if createdImageID != "" {
		return createdImageID
	}

	// Try to find an existing image first
	images, err := c.GetImages()
	if err != nil {
		if os.Getenv("KASM_DEBUG") != "" {
			log.Printf("[DEBUG] Error getting images: %v", err)
			log.Printf("[DEBUG] Will attempt to create new image")
		}
	} else {
		if os.Getenv("KASM_DEBUG") != "" {
			log.Printf("[DEBUG] Found %d existing images", len(images))
			for i, img := range images {
				log.Printf("[DEBUG] Image %d: ID=%s Name=%s Available=%v", i, img.ImageID, img.Name, img.Available)
			}
		}
		// Look for any usable image
		for _, img := range images {
			if img.ImageID != "" && img.Available {
				if os.Getenv("KASM_DEBUG") != "" {
					log.Printf("[DEBUG] Found existing usable image: %s (%s)", img.Name, img.ImageID)
				}
				return img.ImageID
			}
		}
		if os.Getenv("KASM_DEBUG") != "" {
			log.Printf("[DEBUG] No usable image found, will create new one")
		}
	}

	runConfig := map[string]interface{}{
		"hostname":       "kasm-test",
		"container_name": "kasm_test_container",
		"network":        "kasm-network",
		"environment": map[string]string{
			"KASM_TEST": "true",
		},
	}
	runConfigJSON, _ := json.Marshal(runConfig)

	execConfig := map[string]interface{}{
		"shell": "/bin/bash",
		"first_launch": map[string]interface{}{
			"command": "google-chrome --start-maximized",
			"environment": map[string]string{
				"CHROME_TEST": "true",
			},
		},
	}
	execConfigJSON, _ := json.Marshal(execConfig)

	image := &client.CreateImageRequest{
		ImageSrc:           "img/thumbnails/chrome.png",
		Categories:         "Testing",
		RunConfig:          string(runConfigJSON),
		Description:        "Test image for acceptance tests",
		FriendlyName:       "Chrome_Test",
		DockerRegistry:     "https://index.docker.io/v1/",
		Name:               "kasmweb/chrome:1.16.0",
		UncompressedSizeMB: 2000,
		ImageType:          "Container",
		Enabled:            true,
		Memory:             2048000000,
		Cores:              2,
		GPUCount:           0,
		RequireGPU:         nil,
		ExecConfig:         string(execConfigJSON),
	}

	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] Creating test image with config: %+v", image)
	}

	var createdImage *client.Image
	var lastErr error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if os.Getenv("KASM_DEBUG") != "" {
			log.Printf("[DEBUG] Attempt %d/%d to create image", i+1, maxRetries)
		}

		createdImage, lastErr = c.AddWorkspaceImage(image)

		if lastErr != nil {
			if os.Getenv("KASM_DEBUG") != "" {
				log.Printf("[DEBUG] Attempt %d failed with error: %v", i+1, lastErr)
			}
			if i < maxRetries-1 {
				sleepDuration := time.Second * time.Duration(i+1)
				if os.Getenv("KASM_DEBUG") != "" {
					log.Printf("[DEBUG] Waiting %v before next attempt", sleepDuration)
				}
				time.Sleep(sleepDuration)
			}
			continue
		}

		if createdImage == nil {
			lastErr = fmt.Errorf("API returned success but image is nil")
			if os.Getenv("KASM_DEBUG") != "" {
				log.Printf("[DEBUG] Attempt %d failed: %v", i+1, lastErr)
			}
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		if createdImage.ImageID == "" {
			lastErr = fmt.Errorf("API returned image but ImageID is empty")
			if os.Getenv("KASM_DEBUG") != "" {
				log.Printf("[DEBUG] Attempt %d failed: %v", i+1, lastErr)
			}
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		// Success! Store the ID for cleanup
		createdImageID = createdImage.ImageID
		if os.Getenv("KASM_DEBUG") != "" {
			log.Printf("[DEBUG] Successfully created test image with ID: %s", createdImageID)
			log.Printf("[DEBUG] Image details: %+v", createdImage)
		}
		return createdImageID
	}

	errorMsg := fmt.Sprintf("Failed to create workspace image after %d attempts. Last error: %v", maxRetries, lastErr)
	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] %s", errorMsg)
	}
	t.Fatal(errorMsg)
	return "" // unreachable
}

// Test Configuration Functions

// testAccKasmCastConfigConfig_basic creates a basic cast configuration
// Parameters:
// - rName: resource name
// - uniqueKey: unique identifier for the configuration
// - imageID: ID of the workspace image to use
func testAccKasmCastConfigConfig_basic(rName, uniqueKey, imageID string) string {
	baseURL := os.Getenv("KASM_BASE_URL")
	apiKey := os.Getenv("KASM_API_KEY")
	apiSecret := os.Getenv("KASM_API_SECRET")

	return fmt.Sprintf(`
provider "kasm" {
    base_url   = "%s"
    api_key    = "%s"
    api_secret = "%s"
    insecure   = true
}

resource "kasm_cast_config" "test" {
    # Basic identification
    name                  = "%s"
    image_id             = "%s"
    key                  = "%s"

    # Boolean configuration flags - all explicitly set
    enable_sharing       = false
    disable_control_panel = false
    disable_tips         = false
    disable_fixed_res    = false
    allow_anonymous      = false
    require_recaptcha    = false
    dynamic_kasm_url     = false
    dynamic_docker_network = false
    allow_resume         = false
    enforce_client_settings = true
    allow_kasm_audio     = true
    allow_kasm_uploads   = false
    allow_kasm_downloads = false
    allow_clipboard_down = false
    allow_clipboard_up   = false
    allow_microphone     = false
    allow_sharing        = false
    audio_default_on     = true
    ime_mode_default_on  = true

    # Numeric configuration values
    limit_sessions      = true
    session_remaining   = 10
    limit_ips          = true
    ip_request_limit   = 5
    ip_request_seconds = 60

    # String configuration values
    error_url          = "https://error.example.com"
    kasm_url           = "https://start.example.com"
    valid_until        = "2024-12-31 23:59:59"

    # List configuration values
    allowed_referrers  = ["example.com", "test.com"]
}
`, baseURL, apiKey, apiSecret, rName, imageID, uniqueKey)
}

// testAccKasmCastConfigConfig_full creates a comprehensive cast configuration
// with all available options configured
// Parameters:
// - rName: resource name
// - uniqueKey: unique identifier for the configuration
// - imageID: ID of the workspace image to use
func testAccKasmCastConfigConfig_full(rName, uniqueKey, imageID string) string {
	baseURL := os.Getenv("KASM_BASE_URL")
	apiKey := os.Getenv("KASM_API_KEY")
	apiSecret := os.Getenv("KASM_API_SECRET")

	return fmt.Sprintf(`
provider "kasm" {
    base_url   = "%s"
    api_key    = "%s"
    api_secret = "%s"
    insecure   = true
}

resource "kasm_cast_config" "test" {
    # Basic settings
    name                  = "%s"
    image_id             = "%s"
    key                  = "%s"

    # Comprehensive configuration
    enable_sharing        = true
    disable_control_panel = false
    disable_tips         = false
    disable_fixed_res    = false
    allow_anonymous      = true
    require_recaptcha    = false
    limit_sessions       = true
    session_remaining    = 10
    limit_ips           = true
    ip_request_limit    = 5
    ip_request_seconds  = 3600
    error_url           = "https://error.example.com"
    kasm_url            = "https://start.example.com"
    dynamic_kasm_url    = true
    dynamic_docker_network = false
    allow_resume        = true
    enforce_client_settings = true
    allow_kasm_audio    = true
    allow_kasm_uploads  = false
    allow_kasm_downloads = false
    allow_clipboard_down = false
    allow_clipboard_up  = false
    allow_microphone    = false
    allowed_referrers   = ["example.com", "test.com"]
}
`, baseURL, apiKey, apiSecret, rName, imageID, uniqueKey)
}

// testAccKasmCastConfigConfig_update creates an updated cast configuration
// for testing configuration changes
// Parameters:
// - rName: resource name
// - uniqueKey: unique identifier for the configuration
// - imageID: ID of the workspace image to use
func testAccKasmCastConfigConfig_update(rName, uniqueKey, imageID string) string {
	baseURL := os.Getenv("KASM_BASE_URL")
	apiKey := os.Getenv("KASM_API_KEY")
	apiSecret := os.Getenv("KASM_API_SECRET")

	return fmt.Sprintf(`
provider "kasm" {
    base_url   = "%s"
    api_key    = "%s"
    api_secret = "%s"
    insecure   = true
}

resource "kasm_cast_config" "test" {
    # Basic settings with updated name
    name                  = "%s"
    image_id             = "%s"
    key                  = "%s"

    # All boolean fields updated to true
    enable_sharing       = true
    disable_control_panel = true
    disable_tips         = true
    disable_fixed_res    = true
    allow_anonymous      = true
    require_recaptcha    = true
    dynamic_kasm_url     = true
    dynamic_docker_network = true
    allow_resume         = true
    enforce_client_settings = true
    allow_kasm_audio     = true
    allow_kasm_uploads   = true
    allow_kasm_downloads = true
    allow_clipboard_down = true
    allow_clipboard_up   = true
    allow_microphone     = true
    allow_sharing        = true
    audio_default_on     = true
    ime_mode_default_on  = true

    # Updated numeric values
    limit_sessions      = true
    session_remaining   = 20
    limit_ips          = true
    ip_request_limit   = 10
    ip_request_seconds = 120

    # Updated string values
    error_url          = "https://error-updated.example.com"
    kasm_url           = "https://start-updated.example.com"
    valid_until        = "2025-12-31 23:59:59"

    # Updated list values
    allowed_referrers  = ["example.com", "test.com", "updated.com"]
}
`, baseURL, apiKey, apiSecret, rName, imageID, uniqueKey)
}
