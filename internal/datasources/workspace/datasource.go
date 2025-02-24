package workspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKasmWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKasmWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Add other attributes as needed
		},
	}
}

func dataSourceKasmWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	// Implement read logic
	return nil
}
