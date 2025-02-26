package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kasm": providerserver.NewProtocol6WithError(New()),
}

func testAccPreCheck(t *testing.T) {
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

func TestProvider(t *testing.T) {
	t.Parallel()

	p := New()
	if p == nil {
		t.Fatal("provider is nil")
	}
}

func TestProvider_Configure(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		values      map[string]tftypes.Value
		expectError bool
	}{
		"valid": {
			values: map[string]tftypes.Value{
				"base_url":   tftypes.NewValue(tftypes.String, "https://example.com"),
				"api_key":    tftypes.NewValue(tftypes.String, "test-key"),
				"api_secret": tftypes.NewValue(tftypes.String, "test-secret"),
				"insecure":   tftypes.NewValue(tftypes.Bool, false),
			},
			expectError: false,
		},
		"missing_base_url": {
			values: map[string]tftypes.Value{
				"base_url":   tftypes.NewValue(tftypes.String, nil),
				"api_key":    tftypes.NewValue(tftypes.String, "test-key"),
				"api_secret": tftypes.NewValue(tftypes.String, "test-secret"),
				"insecure":   tftypes.NewValue(tftypes.Bool, false),
			},
			expectError: true,
		},
		"missing_api_key": {
			values: map[string]tftypes.Value{
				"base_url":   tftypes.NewValue(tftypes.String, "https://example.com"),
				"api_key":    tftypes.NewValue(tftypes.String, nil),
				"api_secret": tftypes.NewValue(tftypes.String, "test-secret"),
				"insecure":   tftypes.NewValue(tftypes.Bool, false),
			},
			expectError: true,
		},
		"missing_api_secret": {
			values: map[string]tftypes.Value{
				"base_url":   tftypes.NewValue(tftypes.String, "https://example.com"),
				"api_key":    tftypes.NewValue(tftypes.String, "test-key"),
				"api_secret": tftypes.NewValue(tftypes.String, nil),
				"insecure":   tftypes.NewValue(tftypes.Bool, false),
			},
			expectError: true,
		},
		"invalid_base_url": {
			values: map[string]tftypes.Value{
				"base_url":   tftypes.NewValue(tftypes.String, "not-a-url"),
				"api_key":    tftypes.NewValue(tftypes.String, "test-key"),
				"api_secret": tftypes.NewValue(tftypes.String, "test-secret"),
				"insecure":   tftypes.NewValue(tftypes.Bool, false),
			},
			expectError: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			p := &kasmProvider{}
			var diags diag.Diagnostics

			schemaType := tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"base_url":   tftypes.String,
					"api_key":    tftypes.String,
					"api_secret": tftypes.String,
					"insecure":   tftypes.Bool,
				},
			}

			configValue := tftypes.NewValue(schemaType, tc.values)

			resp := &provider.ConfigureResponse{
				DataSourceData: nil,
				ResourceData:   nil,
				Diagnostics:    diags,
			}

			req := provider.ConfigureRequest{
				Config: tfsdk.Config{
					Raw: configValue,
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"base_url": schema.StringAttribute{
								Required: true,
							},
							"api_key": schema.StringAttribute{
								Required:  true,
								Sensitive: true,
							},
							"api_secret": schema.StringAttribute{
								Required:  true,
								Sensitive: true,
							},
							"insecure": schema.BoolAttribute{
								Optional: true,
							},
						},
					},
				},
			}

			ctx := context.Background()
			p.Configure(ctx, req, resp)

			if tc.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tc.expectError && resp.Diagnostics.HasError() {
				t.Errorf("expected no error but got: %v", resp.Diagnostics)
			}
		})
	}
}

func TestProvider_Schema(t *testing.T) {
	t.Parallel()

	p := New()

	// Test schema request
	var resp provider.SchemaResponse
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	// Verify schema is set
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("Schema attributes are empty")
	}

	// Verify required attributes
	requiredAttrs := []string{"base_url", "api_key", "api_secret"}
	for _, attrName := range requiredAttrs {
		attr := resp.Schema.Attributes[attrName]
		if attr == nil {
			t.Errorf("expected %s to exist", attrName)
			continue
		}
		if !attr.IsRequired() {
			t.Errorf("expected %s to be required", attrName)
		}
	}

	// Verify optional attributes
	optionalAttrs := []string{"insecure"}
	for _, attrName := range optionalAttrs {
		attr := resp.Schema.Attributes[attrName]
		if attr == nil {
			t.Errorf("expected %s to exist", attrName)
			continue
		}
		if attr.IsRequired() {
			t.Errorf("expected %s to be optional", attrName)
		}
	}

	// Verify sensitive attributes
	sensitiveAttrs := []string{"api_key", "api_secret"}
	for _, attrName := range sensitiveAttrs {
		attr := resp.Schema.Attributes[attrName]
		if attr == nil {
			t.Errorf("expected %s to exist", attrName)
			continue
		}
		if !attr.IsSensitive() {
			t.Errorf("expected %s to be sensitive", attrName)
		}
	}
}

func TestProvider_Metadata(t *testing.T) {
	t.Parallel()

	p := New()
	var resp provider.MetadataResponse
	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	if resp.TypeName != "kasm" {
		t.Errorf("expected type name to be 'kasm', got %s", resp.TypeName)
	}
}

func TestProvider_Resources(t *testing.T) {
	t.Parallel()

	p := New()
	resources := p.Resources(context.Background())

	// Verify that we have resources defined
	if len(resources) == 0 {
		t.Error("expected resources to be non-empty")
	}
}

func TestProvider_DataSources(t *testing.T) {
	t.Parallel()

	p := New()
	dataSources := p.DataSources(context.Background())

	// Verify that we have data sources defined
	if len(dataSources) == 0 {
		t.Error("expected data sources to be non-empty")
	}
}
