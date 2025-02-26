//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
	"os"
	"terraform-provider-kasm/testutils"
	"testing"
)

func TestAccKasmKeepalive_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	// Setup test client
	c := testutils.GetTestClient(t)
	if c == nil {
		t.Fatal("Failed to get test client")
	}

	// Get a valid user ID
	users, err := c.GetUsers()
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}
	if len(users) == 0 {
		t.Fatal("No users found in the system")
	}
	userID := users[0].UserID
	log.Printf("[DEBUG] Using user ID: %s", userID)

	// Get a valid image ID
	images, err := c.GetImages()
	if err != nil {
		t.Fatalf("Failed to get images: %v", err)
	}
	if len(images) == 0 {
		t.Fatal("No images found in the system")
	}

	// Find an available image
	var imageID string
	for _, img := range images {
		if img.Available {
			imageID = img.ImageID
			break
		}
	}
	if imageID == "" {
		t.Fatal("No available images found")
	}
	log.Printf("[DEBUG] Using image ID: %s", imageID)

	// Create a test kasm
	kasm, err := c.CreateKasm(
		userID,
		imageID,
		"", // empty session token, will be created automatically
		users[0].Username,
		false, // share
		false, // persistent
		false, // allowResume
		false, // sessionAuthentication
	)
	if err != nil {
		t.Fatalf("Failed to create test kasm: %v", err)
	}
	log.Printf("[DEBUG] Created test kasm with ID: %s", kasm.KasmID)

	// Ensure kasm cleanup
	t.Cleanup(func() {
		if err := c.DestroyKasm(userID, kasm.KasmID); err != nil {
			t.Logf("Warning: Failed to delete test kasm %s: %v", kasm.KasmID, err)
		} else {
			log.Printf("[DEBUG] Successfully deleted test kasm: %s", kasm.KasmID)
		}
	})

	resourceName := "kasm_keepalive.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmKeepaliveConfig_basic(kasm.KasmID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", kasm.KasmID),
					resource.TestCheckResourceAttr(resourceName, "kasm_id", kasm.KasmID),
				),
			},
		},
	})
}

func testAccKasmKeepaliveConfig_basic(kasmID string) string {
	return fmt.Sprintf(`
provider "kasm" {
  base_url   = "%s"
  api_key    = "%s"
  api_secret = "%s"
  insecure   = true
}

resource "kasm_keepalive" "test" {
	kasm_id = "%s"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), kasmID)
}
