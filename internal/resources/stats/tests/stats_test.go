//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func generateUniqueUsername() string {
	// Use timestamp and random number to ensure uniqueness
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("testuser_%d_%d", timestamp, randomNum)
}

func TestAccKasmStats_FrameStats(t *testing.T) {
	username := generateUniqueUsername()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmStatsConfig_frameStats(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("kasm_stats.test", "id"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "kasm_id"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "res_x"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "res_y"),
					resource.TestCheckResourceAttrSet("kasm_stats.test", "last_updated"),
				),
			},
		},
	})
}

func testAccKasmStatsConfig_frameStats(username string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_user" "test" {
    username = "%s"
    password = "TestPassword123!"
    first_name = "Test"
    last_name = "User"
    organization = "Test Org"
    locked = false
    disabled = false
    groups = []
}

data "kasm_images" "available" {
    depends_on = [kasm_user.test]
}

resource "kasm_session" "test" {
    image_id = data.kasm_images.available.images[0].image_id
    user_id = kasm_user.test.id
    enable_stats = true
}

resource "kasm_stats" "test" {
    kasm_id = kasm_session.test.id
    user_id = kasm_user.test.id
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), username)
}
