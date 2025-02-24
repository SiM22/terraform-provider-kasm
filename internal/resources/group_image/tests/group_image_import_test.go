package tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"terraform-provider-kasm/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupImage_import(t *testing.T) {
	// Remove parallel execution since we're dealing with shared resources
	// t.Parallel()

	// Get client and ensure test image exists
	client := testutils.GetTestClient(t)
	imageID := ensureWorkspaceImage(t, client)
	if imageID == "" {
		t.Fatal("Failed to get a valid image ID")
	}

	// Ensure cleanup runs after the test
	t.Cleanup(func() {
		if imageID != "" {
			cleanupTestImage(t, client)
		}
	})

	// Generate unique identifiers with both timestamp and random number
	rand.Seed(time.Now().UnixNano())
	uniqueIdentifier := fmt.Sprintf("%d_%d", time.Now().Unix(), rand.Intn(10000))
	groupname := fmt.Sprintf("testgroup_%s", uniqueIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resources
			{
				Config: testAccGroupImageConfig_forImport(groupname, imageID),
				Check: resource.ComposeTestCheckFunc(
					testutils.TestCheckResourceExists("kasm_group.test"),
					testutils.TestCheckResourceExists("kasm_group_image.test"),
					// Verify all attributes are set
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "group_image_id"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "group_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_friendly_name"),
					resource.TestCheckResourceAttrSet("kasm_group_image.test", "image_src"),
				),
			},
			// Test successful import using group_id:image_id format
			{
				ResourceName:            "kasm_group_image.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"id"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["kasm_group_image.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "kasm_group_image.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["image_id"]), nil
				},
			},
			// Test import with invalid ID format
			{
				ResourceName:  "kasm_group_image.test",
				ImportState:   true,
				ImportStateId: "invalid-id-format",
				ExpectError:   regexp.MustCompile(`Import ID must be in the format group_id:image_id`),
			},
			// Test import with non-existent IDs
			{
				ResourceName:  "kasm_group_image.test",
				ImportState:   true,
				ImportStateId: "non-existent-group:non-existent-image",
				ExpectError:   regexp.MustCompile(`Error Reading Group Images`),
			},
		},
	})
}

func testAccGroupImageConfig_forImport(groupname, imageID string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name = "%s"
    priority = 50
    description = "Test Import Group"
}

resource "kasm_group_image" "test" {
    group_id = kasm_group.test.id
    image_id = "%s"
}
`, testutils.ProviderConfig(), groupname, imageID)
}
