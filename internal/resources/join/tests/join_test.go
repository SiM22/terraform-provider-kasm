//go:build acceptance
// +build acceptance

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/testutils"
)

func TestAccKasmJoin_basic(t *testing.T) {
	// Skip this test for now as it requires additional setup
	t.Skip("Skipping join test until sharing functionality is fully implemented")

	testutils.LoadEnvVars(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.PreCheck(t) },
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "kasm_join" "test" {
						share_id = "test_share_id"
						user_id  = "test_user_id"
					}
				`,
				ExpectError: nil,
			},
		},
	})
}
