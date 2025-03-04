package testutils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/joho/godotenv"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/provider"
)

// LoadEnvFile loads environment variables from .env file
func LoadEnvFile(t *testing.T) {
	// Try to find .env file in current directory and parent directories
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			err = godotenv.Load(envFile)
			if err != nil {
				t.Fatalf("Error loading .env file: %v", err)
			}
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Log("No .env file found, using existing environment variables")
}

// TestAccPreCheck validates required environment variables exist
func TestAccPreCheck(t *testing.T) {
	LoadEnvFile(t)

	if v := os.Getenv("KASM_BASE_URL"); v == "" {
		t.Fatal("KASM_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("KASM_API_KEY"); v == "" {
		t.Fatal("KASM_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("KASM_API_SECRET"); v == "" {
		t.Fatal("KASM_API_SECRET must be set for acceptance tests")
	}
}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance tests
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kasm": providerserver.NewProtocol6WithError(provider.New()),
}

// GetTestClient returns a configured Kasm client for testing
func GetTestClient(t testing.TB) *client.Client {
	return client.NewClient(
		os.Getenv("KASM_BASE_URL"),
		os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"),
		true,
	)
}

// providerConfig (lowercase) for internal use
func providerConfig() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

// ProviderConfig (uppercase) for external use
func ProviderConfig() string {
	return providerConfig()
}

// GenerateUniqueUsername generates a unique username for testing
func GenerateUniqueUsername() string {
	// Use timestamp and random number to ensure uniqueness
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("testuser_%d_%d", timestamp, randomNum)
}

// CleanupExistingSessions cleans up any existing sessions
func CleanupExistingSessions(t testing.TB) {
	// Get the client
	c := GetTestClient(t)
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

// Global variable to track created image ID for cleanup
var createdTestImageID string

// CleanupTestImage deletes the test image if it was created during testing
func CleanupTestImage(t testing.TB) {
	if createdTestImageID == "" {
		return
	}

	c := GetTestClient(t)
	if c == nil {
		t.Logf("Warning: Failed to get client for test image cleanup")
		return
	}

	if err := c.DeleteImage(createdTestImageID); err != nil {
		t.Logf("Warning: Failed to delete test image %s: %v", createdTestImageID, err)
	} else {
		t.Logf("Successfully deleted test image: %s", createdTestImageID)
		createdTestImageID = "" // Reset after successful deletion
	}
}

// EnsureImageAvailable ensures a valid image is available for testing
// If no suitable image exists, it will create one
func EnsureImageAvailable(t testing.TB) (string, bool) {
	// If we already created an image in this test run, return its ID
	if createdTestImageID != "" {
		t.Logf("Using previously created test image: %s", createdTestImageID)
		return createdTestImageID, true
	}

	maxRetries := 10
	retryDelay := 5 * time.Second

	// Get the client
	c := GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
		return "", false
	}

	// Try to get images with retries
	var images []client.Image
	for i := 0; i < maxRetries; i++ {
		resp, err := c.GetImages()
		if err != nil {
			t.Logf("Warning: Failed to get images (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		images = resp
		break
	}

	// Find a suitable image (prefer Ubuntu or similar if available)
	preferredNames := []string{"ubuntu", "debian", "centos", "fedora", "alpine", "chrome", "firefox", "filezilla"}

	// First try to find a preferred image
	for _, name := range preferredNames {
		for _, img := range images {
			if img.Enabled && img.Available && (containsIgnoreCase(img.Name, name) || containsIgnoreCase(img.FriendlyName, name)) {
				t.Logf("Found preferred image: %s (%s)", img.FriendlyName, img.ImageID)
				return img.ImageID, true
			}
		}
	}

	// If no preferred image, use the first enabled and available image
	for _, img := range images {
		if img.Enabled && img.Available {
			t.Logf("Using image: %s (%s)", img.FriendlyName, img.ImageID)
			return img.ImageID, true
		}
	}

	// No suitable image found, create one
	t.Log("No suitable images found, creating a test image")
	return createTestImage(t, c)
}

// createTestImage creates a test image for testing
func createTestImage(t testing.TB, c *client.Client) (string, bool) {
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
		Description:        "Test image for Terraform Provider Kasm tests",
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

	t.Logf("Creating test image with config: %+v", image)

	var createdImage *client.Image
	var lastErr error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		t.Logf("Attempt %d/%d to create image", i+1, maxRetries)

		createdImage, lastErr = c.AddWorkspaceImage(image)

		if lastErr != nil {
			t.Logf("Attempt %d failed with error: %v", i+1, lastErr)
			if i < maxRetries-1 {
				sleepDuration := time.Second * time.Duration(i+1)
				t.Logf("Waiting %v before next attempt", sleepDuration)
				time.Sleep(sleepDuration)
			}
			continue
		}

		if createdImage == nil {
			lastErr = fmt.Errorf("API returned success but image is nil")
			t.Logf("Attempt %d failed: %v", i+1, lastErr)
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		if createdImage.ImageID == "" {
			lastErr = fmt.Errorf("API returned image but ImageID is empty")
			t.Logf("Attempt %d failed: %v", i+1, lastErr)
			if i < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(i+1))
			}
			continue
		}

		// Success! Store the ID for cleanup
		createdTestImageID = createdImage.ImageID
		t.Logf("Successfully created test image with ID: %s", createdTestImageID)

		// Wait for the image to be downloaded and available
		t.Logf("Waiting for image to be downloaded and available...")
		waitForImageAvailable(t, c, createdTestImageID)

		return createdTestImageID, true
	}

	errorMsg := fmt.Sprintf("Failed to create workspace image after %d attempts. Last error: %v", maxRetries, lastErr)
	t.Logf("%s", errorMsg)
	return "", false
}

// waitForImageAvailable waits for an image to be downloaded and available
func waitForImageAvailable(t testing.TB, c *client.Client, imageID string) {
	maxRetries := 30
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		images, err := c.GetImages()
		if err != nil {
			t.Logf("Warning: Failed to get images while waiting (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		for _, img := range images {
			if img.ImageID == imageID {
				if img.Available {
					t.Logf("Image %s is now available", imageID)
					return
				}
				t.Logf("Image %s is not yet available (attempt %d/%d)", imageID, i+1, maxRetries)
				break
			}
		}

		time.Sleep(retryDelay)
	}

	t.Logf("Warning: Timed out waiting for image %s to become available", imageID)
}

// Helper function to check if a string contains another string (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// User test configurations
func TestAccUserConfig(username string) string {
	return fmt.Sprintf(`
%s

resource "kasm_user" "test" {
    username     = "%s"
    password     = "TestImport123!"
    first_name   = "Test"
    last_name    = "Import"
    organization = "Import Testing"
    phone        = "1234567890"
    groups       = []
}`, ProviderConfig(), username)
}

// Group test configurations
func TestAccGroupConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name        = "%s"
    priority    = 50
    description = "Test Import Group"
}`, ProviderConfig(), name)
}

// Combined user and group test configurations
func TestAccUserGroupConfig(username, groupName string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name        = "%s"
    priority    = 50
    description = "Test group for user import"
}

resource "kasm_user" "test" {
    username     = "%s"
    password     = "TestImport123!"
    first_name   = "Test"
    last_name    = "Import"
    organization = "Import Testing"
    groups       = [kasm_group.test.name]
}`, ProviderConfig(), groupName, username)
}

// Empty config for import testing
func TestAccEmptyConfig() string {
	return ProviderConfig()
}

// TestCheckResourceExists verifies that a resource exists in state by checking an attribute
func TestCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"), // Change "id" to an actual attribute name to confirm existence
	)
}
