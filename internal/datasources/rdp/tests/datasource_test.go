//go:build unit
// +build unit

package tests

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
	"terraform-provider-kasm/internal/datasources/rdp"
)

// Test the data source schema
func TestRDPClientConnectionInfoDataSource_Schema(t *testing.T) {
	t.Parallel()

	// Create the data source
	ds := rdp.NewRDPClientConnectionInfoDataSource()

	// Verify it implements the expected interface
	_, ok := ds.(datasource.DataSource)
	assert.True(t, ok, "Data source doesn't implement datasource.DataSource")
}

// Test the data source metadata
func TestRDPClientConnectionInfoDataSource_Metadata(t *testing.T) {
	t.Parallel()

	// Create the data source
	ds := rdp.NewRDPClientConnectionInfoDataSource()

	// Get the data source as the correct interface
	dataSource, ok := ds.(datasource.DataSource)
	assert.True(t, ok, "Data source doesn't implement datasource.DataSource")

	// Create a metadata request and response
	req := datasource.MetadataRequest{
		ProviderTypeName: "kasm",
	}
	resp := &datasource.MetadataResponse{}

	// Call the Metadata method
	dataSource.Metadata(context.Background(), req, resp)

	// Verify the type name
	assert.Equal(t, "kasm_rdp_client_connection_info", resp.TypeName, "Unexpected type name")
}

// Test the data source schema method
func TestRDPClientConnectionInfoDataSource_SchemaMethod(t *testing.T) {
	t.Parallel()

	// Create the data source
	ds := rdp.NewRDPClientConnectionInfoDataSource()

	// Get the data source as the correct interface
	dataSource, ok := ds.(datasource.DataSource)
	assert.True(t, ok, "Data source doesn't implement datasource.DataSource")

	// Create a schema request and response
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	// Call the Schema method
	dataSource.Schema(context.Background(), req, resp)

	// Verify schema attributes
	assert.NotNil(t, resp.Schema, "Schema should not be nil")
	assert.NotEmpty(t, resp.Schema.Description, "Schema should have a description")

	// Check required attributes
	assert.Contains(t, resp.Schema.Attributes, "id", "Schema should have id attribute")
	assert.Contains(t, resp.Schema.Attributes, "user_id", "Schema should have user_id attribute")
	assert.Contains(t, resp.Schema.Attributes, "kasm_id", "Schema should have kasm_id attribute")
	assert.Contains(t, resp.Schema.Attributes, "connection_type", "Schema should have connection_type attribute")
	assert.Contains(t, resp.Schema.Attributes, "file", "Schema should have file attribute")
	assert.Contains(t, resp.Schema.Attributes, "url", "Schema should have url attribute")
}
