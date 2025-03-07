//go:build acceptance
// +build acceptance

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

// TestAccSessionsDataSource tests the sessions data source against a real Kasm instance.
func TestAccSessionsDataSource(t *testing.T) {
	// Skip if not running acceptance tests
	testutils.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSessionsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_sessions.test", "current_time"),
					// We can't check for specific sessions as they may vary, but we can check that the data source works
					resource.TestCheckResourceAttrSet("data.kasm_sessions.test", "id"),
				),
			},
		},
	})
}

// testAccSessionsDataSourceConfig returns the Terraform configuration for testing the sessions data source.
func testAccSessionsDataSourceConfig() string {
	return testutils.ProviderConfig() + `
	data "kasm_sessions" "test" {}
	`
}
