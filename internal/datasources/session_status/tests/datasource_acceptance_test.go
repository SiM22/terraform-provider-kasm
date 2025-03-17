//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/testutils"
)

// TestAccSessionStatusDataSource tests the session_status data source against a real Kasm instance.
func TestAccSessionStatusDataSource(t *testing.T) {
	// Skip if not running acceptance tests
	testutils.TestAccPreCheck(t)

	// Clean up any existing sessions
	testutils.CleanupExistingSessions(t)

	// Get the client
	c := testutils.GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
	}

	// Generate a unique username for the test
	username := testutils.GenerateUniqueUsername()

	// Ensure we have a valid image
	imageID, found := testutils.EnsureImageAvailable(t)
	if !found {
		t.Skip("No valid image found for testing")
	}

	// Create a user for the test
	user := &client.User{
		Username:  username,
		Password:  "Password123!",
		FirstName: "Test",
		LastName:  "User",
		Phone:     "",
		Locked:    false,
		Disabled:  false,
	}

	createdUser, err := c.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	userID := createdUser.UserID
	defer func() {
		// Clean up the user
		err := c.DeleteUser(userID)
		if err != nil {
			t.Logf("Warning: Failed to delete test user: %v", err)
		}
	}()

	// Create a session for the user
	// First get a session token
	// For testing purposes, we'll use the CreateKasm method directly
	sessionResp, err := c.CreateKasm(
		userID,
		imageID,
		"", // No session token needed for API key auth
		username,
		false, // share
		false, // persistent
		true,  // allow resume
		false, // session authentication
	)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}
	kasmID := sessionResp.KasmID

	// Log detailed session creation response
	t.Logf("[DEBUG] Session created successfully:")
	t.Logf("[DEBUG] - Kasm ID: %s", kasmID)
	t.Logf("[DEBUG] - User ID: %s", userID)
	t.Logf("[DEBUG] - Username: %s", username)
	t.Logf("[DEBUG] - Image ID: %s", imageID)
	t.Logf("[DEBUG] - Session Token: %s", sessionResp.SessionToken)
	t.Logf("[DEBUG] - Kasm URL: %s", sessionResp.KasmURL)

	defer func() {
		// Clean up the session
		t.Logf("[DEBUG] Cleaning up session: %s", kasmID)
		err := c.DestroyKasm(userID, kasmID)
		if err != nil {
			t.Logf("Warning: Failed to destroy test session: %v", err)
		}
	}()

	// Wait for the session to be fully initialized
	log.Printf("[DEBUG] Waiting for session to initialize... (kasm_id: %s, user_id: %s)", kasmID, userID)
	maxRetries := 15 // Increased from 10 to 15
	var sessionReady bool

	// Add an initial longer wait to give the session more time to initialize
	log.Printf("[DEBUG] Initial wait of 60 seconds for session to initialize...")
	time.Sleep(10 * time.Second)

	for i := 0; i < maxRetries; i++ {
		// Check if the session is available
		log.Printf("[DEBUG] Checking session status (attempt %d/%d)...", i+1, maxRetries)
		status, err := c.GetKasmStatus(userID, kasmID, true)
		if err != nil {
			t.Logf("Attempt %d: Session not ready yet: %v. Retrying...", i+1, err)
			// Exponential backoff: 10s, 15s, 20s, 25s, etc.
			waitTime := 10 + (i * 5)
			t.Logf("[DEBUG] Waiting %d seconds before next retry...", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			continue
		}

		if status.Kasm != nil && status.Kasm.ContainerID != "" {
			t.Logf("Session is ready after %d attempts", i+1)
			t.Logf("[DEBUG] Session details:")
			t.Logf("[DEBUG] - Container ID: %s", status.Kasm.ContainerID)
			t.Logf("[DEBUG] - Container IP: %s", status.Kasm.ContainerIP)
			t.Logf("[DEBUG] - Host: %s", status.Kasm.Host)
			t.Logf("[DEBUG] - Port: %d", status.Kasm.Port)
			sessionReady = true
			break
		}

		t.Logf("Attempt %d: Session not fully initialized yet (container_id empty). Retrying...", i+1)
		// Exponential backoff: 10s, 15s, 20s, 25s, etc.
		waitTime := 10 + (i * 5)
		t.Logf("[DEBUG] Waiting %d seconds before next retry...", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)
	}

	if !sessionReady {
		t.Logf("[WARN] Session may not be fully initialized after %d attempts. Proceeding with test anyway.", maxRetries)
		t.Logf("[DEBUG] Session details - kasm_id: %s, user_id: %s", kasmID, userID)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionStatusDataSourceConfig(userID, kasmID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_session_status.test", "status"),
					resource.TestCheckResourceAttrSet("data.kasm_session_status.test", "operational_status"),
					resource.TestCheckResourceAttrSet("data.kasm_session_status.test", "image_id"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "kasm_id", kasmID),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "user_id", userID),
				),
			},
		},
	})
}

// testAccSessionStatusDataSourceConfig returns the Terraform configuration for testing the session_status data source.
func testAccSessionStatusDataSourceConfig(userID, kasmID string) string {
	return testutils.ProviderConfig() + fmt.Sprintf(`
	data "kasm_session_status" "test" {
		user_id = "%s"
		kasm_id = "%s"
		skip_agent_check = true
	}
	`, userID, kasmID)
}
