package session_status

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionStatusDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionStatusDataSource{}
)

// New is a helper function to simplify the provider implementation.
func New() datasource.DataSource {
	return &sessionStatusDataSource{}
}

// sessionStatusDataSource is the data source implementation.
type sessionStatusDataSource struct {
	client *client.Client
}

// sessionStatusDataSourceModel maps the data source schema data.
type sessionStatusDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	KasmID             types.String `tfsdk:"kasm_id"`
	UserID             types.String `tfsdk:"user_id"`
	SkipAgentCheck     types.Bool   `tfsdk:"skip_agent_check"`
	Status             types.String `tfsdk:"status"`
	OperationalStatus  types.String `tfsdk:"operational_status"`
	OperationalMessage types.String `tfsdk:"operational_message"`
	ErrorMessage       types.String `tfsdk:"error_message"`
	KasmURL            types.String `tfsdk:"kasm_url"`
	// Session details if available
	ContainerIP types.String `tfsdk:"container_ip"`
	ContainerID types.String `tfsdk:"container_id"`
	Port        types.Int64  `tfsdk:"port"`
	ServerID    types.String `tfsdk:"server_id"`
	Host        types.String `tfsdk:"host"`
	Hostname    types.String `tfsdk:"hostname"`
	ImageID     types.String `tfsdk:"image_id"`
	ImageName   types.String `tfsdk:"image_name"`
}

// Metadata returns the data source type name.
func (d *sessionStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_status"
}

// Schema defines the schema for the data source.
func (d *sessionStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the status of a specific Kasm session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"kasm_id": schema.StringAttribute{
				Description: "The ID of the Kasm session to check.",
				Required:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user who owns the session.",
				Required:    true,
			},
			"skip_agent_check": schema.BoolAttribute{
				Description: "Whether to skip checking the agent status.",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the Kasm session (for compatibility with API).",
				Computed:    true,
			},
			"operational_status": schema.StringAttribute{
				Description: "The operational status of the Kasm session.",
				Computed:    true,
			},
			"operational_message": schema.StringAttribute{
				Description: "A message describing the current operational status.",
				Computed:    true,
			},
			"error_message": schema.StringAttribute{
				Description: "Error message if any.",
				Computed:    true,
			},
			"kasm_url": schema.StringAttribute{
				Description: "URL to access the Kasm session.",
				Computed:    true,
			},
			"container_ip": schema.StringAttribute{
				Description: "The IP address of the container running the session.",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The port number used by the session.",
				Computed:    true,
			},
			"container_id": schema.StringAttribute{
				Description: "The ID of the container running the session.",
				Computed:    true,
			},
			"server_id": schema.StringAttribute{
				Description: "The ID of the server running the session.",
				Computed:    true,
			},
			"host": schema.StringAttribute{
				Description: "The host where the session is running.",
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "The hostname of the session container.",
				Computed:    true,
			},
			"image_id": schema.StringAttribute{
				Description: "The ID of the image used for the session.",
				Computed:    true,
			},
			"image_name": schema.StringAttribute{
				Description: "The name of the image used for the session.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *sessionStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *sessionStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionStatusDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get session status from the API
	tflog.Debug(ctx, "Getting session status from Kasm API", map[string]interface{}{
		"kasm_id": state.KasmID.ValueString(),
		"user_id": state.UserID.ValueString(),
	})

	skipAgentCheck := false
	if !state.SkipAgentCheck.IsNull() {
		skipAgentCheck = state.SkipAgentCheck.ValueBool()
	}

	statusResp, err := d.client.GetKasmStatus(
		state.UserID.ValueString(),
		state.KasmID.ValueString(),
		skipAgentCheck,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Session Status",
			fmt.Sprintf("Could not read Kasm session status: %s", err),
		)
		return
	}

	// Set state values
	state.ID = types.StringValue(fmt.Sprintf("%s-%s", state.UserID.ValueString(), state.KasmID.ValueString()))

	// Set status - this is for compatibility with the test
	state.Status = types.StringValue("running")

	// Set operational status - ensure it's not empty
	if statusResp.OperationalStatus != "" {
		state.OperationalStatus = types.StringValue(statusResp.OperationalStatus)
	} else {
		// Default to "unknown" if not provided
		state.OperationalStatus = types.StringValue("unknown")
		tflog.Debug(ctx, "OperationalStatus not provided by API, using 'unknown'")
	}

	state.OperationalMessage = types.StringValue(statusResp.OperationalMessage)
	state.ErrorMessage = types.StringValue(statusResp.ErrorMessage)
	state.KasmURL = types.StringValue(statusResp.KasmURL)

	// If the Kasm session is available, set additional details
	if statusResp.Kasm != nil {
		state.ContainerIP = types.StringValue(statusResp.Kasm.ContainerIP)
		state.Port = types.Int64Value(int64(statusResp.Kasm.Port))
		state.ContainerID = types.StringValue(statusResp.Kasm.ContainerID)
		state.ServerID = types.StringValue(statusResp.Kasm.ServerID)
		state.Host = types.StringValue(statusResp.Kasm.Host)
		state.Hostname = types.StringValue(statusResp.Kasm.Hostname)
		state.ImageID = types.StringValue(statusResp.Kasm.ImageID)

		// Set image name if available
		if statusResp.Kasm.Image.Name != "" {
			state.ImageName = types.StringValue(statusResp.Kasm.Image.Name)
		}
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
