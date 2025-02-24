//go:build acceptance
// +build acceptance

package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// ensureWorkspaceImage ensures a test image exists and returns its ID
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

	volumeMappings := map[string]interface{}{}
	volumeMappingsJSON, _ := json.Marshal(volumeMappings)

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
		VolumeMappings:     string(volumeMappingsJSON),
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

func TestAccGroupImage_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	t.Parallel()

	// Get client and ensure test image exists
	client := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, client)

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, client)
	})

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Test
			{
				Config: testAccGroupImageConfig_basic(groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_image.test"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "group_image_id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "group_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_friendly_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_src"),
					resource.TestCheckResourceAttr("kasm_group_image.test", "group_name", groupname),
				),
			},
			// Import Test
			{
				ResourceName:            "kasm_group_image.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
				ImportStateIdFunc:       testAccGroupImageImportIDFunc,
			},
		},
	})
}

func TestAccGroupImage_multiple(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	t.Parallel()

	// Get client and ensure test images exist
	client := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, client)

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		cleanupTestImage(t, client)
	})

	// Initialize random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), r.Intn(10000))
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupImageConfig_multiple(groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_image.test1"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test1", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test1", "group_image_id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test1", "image_name"),
					resource.TestCheckResourceAttr("kasm_group_image.test1", "group_name", groupname),
				),
			},
		},
	})
}

func testAccGroupImageConfig_basic(name string, imageID string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name = "%s"
    priority = 50
    description = "Test group for image association"
}

resource "kasm_group_image" "test" {
    group_id = kasm_group.test.id
    image_id = "%s"
}
`, testutils.ProviderConfig(), name, imageID)
}

func testAccGroupImageConfig_multiple(name string, imageID string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name = "%s"
    priority = 50
    description = "Test group for image association"
}

resource "kasm_group_image" "test1" {
    group_id = kasm_group.test.id
    image_id = "%s"
}

resource "kasm_group_image" "test2" {
    group_id = kasm_group.test.id
    image_id = "%s"
}
`, testutils.ProviderConfig(), name, imageID, imageID)
}

func testAccGroupImageImportIDFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["kasm_group_image.test"]
	if !ok {
		return "", fmt.Errorf("Not found: %s", "kasm_group_image.test")
	}
	return fmt.Sprintf("%s:%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["image_id"]), nil
}
