package keepalive

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

func TestKeepaliveResource_Metadata(t *testing.T) {
	r := &keepaliveResource{}
	req := resource.MetadataRequest{
		ProviderTypeName: "kasm",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "kasm_keepalive", resp.TypeName)
}

func TestKeepaliveResource_Schema(t *testing.T) {
	r := &keepaliveResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	assert.NotNil(t, resp.Schema)
	assert.NotNil(t, resp.Schema.Attributes["id"])
	assert.NotNil(t, resp.Schema.Attributes["kasm_id"])
}
