package rdp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &rdpClientConnectionInfoDataSource{}
)

// NewRDPClientConnectionInfoDataSource is a helper function to simplify the provider implementation.
func NewRDPClientConnectionInfoDataSource() datasource.DataSource {
	return &rdpClientConnectionInfoDataSource{
		client: &client.Client{},
	}
}

// rdpClientConnectionInfoDataSource is the data source implementation.
type rdpClientConnectionInfoDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *rdpClientConnectionInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rdp_client_connection_info"
}

// Schema defines the schema for the data source.
func (d *rdpClientConnectionInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves RDP client connection information for a Kasm session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the RDP client connection info.",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the user for which to retrieve the RDP connection info.",
			},
			"kasm_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Kasm session for which to retrieve the RDP connection info.",
			},
			"connection_type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of connection to retrieve (url or file)",
				Validators: []validator.String{
					validators.StringOneOf("url", "file"),
				},
			},
			"file": schema.StringAttribute{
				Computed:    true,
				Description: "The RDP file content (if connection_type is 'file' or not specified).",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The RDP URL (if connection_type is 'url').",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *rdpClientConnectionInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *rdpClientConnectionInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state rdpClientConnectionInfoModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the connection info from the API
	connectionType := client.RDPConnectionTypeFile
	if !state.ConnectionType.IsNull() && state.ConnectionType.ValueString() == "url" {
		connectionType = client.RDPConnectionTypeURL
	}

	tflog.Debug(ctx, fmt.Sprintf("Getting RDP connection info for user %s, kasm %s, type %s", state.UserID.ValueString(), state.KasmID.ValueString(), connectionType))
	connectionInfo, err := d.client.GetRDPConnectionInfo(state.UserID.ValueString(), state.KasmID.ValueString(), connectionType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting RDP connection info",
			fmt.Sprintf("Could not get RDP connection info: %v", err),
		)
		return
	}

	// Set the ID
	state.ID = types.StringValue(fmt.Sprintf("%s-%s-%s", state.UserID.ValueString(), state.KasmID.ValueString(), connectionType))

	// Set the connection info
	state.File = types.StringValue(connectionInfo.File)
	state.URL = types.StringValue(connectionInfo.URL)

	// Save the data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// rdpClientConnectionInfoModel maps the data source schema data.
type rdpClientConnectionInfoModel struct {
	ID             types.String `tfsdk:"id"`
	UserID         types.String `tfsdk:"user_id"`
	KasmID         types.String `tfsdk:"kasm_id"`
	ConnectionType types.String `tfsdk:"connection_type"`
	File           types.String `tfsdk:"file"`
	URL            types.String `tfsdk:"url"`
}
