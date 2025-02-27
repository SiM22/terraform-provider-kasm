package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &statsResource{}
	_ resource.ResourceWithConfigure   = &statsResource{}
	_ resource.ResourceWithImportState = &statsResource{}
)

// NewStatsResource is a helper function to simplify the provider implementation.
func NewStatsResource() resource.Resource {
	return &statsResource{}
}

// statsResource is the resource implementation.
type statsResource struct {
	client *client.Client
}

// statsResourceModel maps the resource schema data.
type statsResourceModel struct {
	ID           types.String `tfsdk:"id"`
	KasmID       types.String `tfsdk:"kasm_id"`
	UserID       types.String `tfsdk:"user_id"`
	ResX         types.Int64  `tfsdk:"res_x"`
	ResY         types.Int64  `tfsdk:"res_y"`
	Changed      types.Int64  `tfsdk:"changed"`
	ServerTime   types.Int64  `tfsdk:"server_time"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	ClientCount  types.Int64  `tfsdk:"client_count"`
	Analysis     types.Int64  `tfsdk:"analysis"`
	Screenshot   types.Int64  `tfsdk:"screenshot"`
	EncodingTime types.Int64  `tfsdk:"encoding_time"`
}

// Metadata returns the resource type name.
func (r *statsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stats"
}

// Schema defines the schema for the resource.
func (r *statsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages frame statistics for a Kasm session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the stats resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kasm_id": schema.StringAttribute{
				Description: "The ID of the Kasm session to get stats for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user who owns the session.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"res_x": schema.Int64Attribute{
				Description: "The horizontal resolution of the session.",
				Computed:    true,
			},
			"res_y": schema.Int64Attribute{
				Description: "The vertical resolution of the session.",
				Computed:    true,
			},
			"changed": schema.Int64Attribute{
				Description: "The number of changed pixels.",
				Computed:    true,
			},
			"server_time": schema.Int64Attribute{
				Description: "The server processing time in milliseconds.",
				Computed:    true,
			},
			"client_count": schema.Int64Attribute{
				Description: "The number of connected clients.",
				Computed:    true,
			},
			"analysis": schema.Int64Attribute{
				Description: "The time spent on frame analysis in milliseconds.",
				Computed:    true,
			},
			"screenshot": schema.Int64Attribute{
				Description: "The time spent on screenshot processing in milliseconds.",
				Computed:    true,
			},
			"encoding_time": schema.Int64Attribute{
				Description: "The total encoding time in milliseconds.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last refresh of the stats.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *statsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *statsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan statsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate ID
	plan.ID = types.StringValue(fmt.Sprintf("%s-stats", plan.KasmID.ValueString()))

	// Get initial stats
	frameStats, err := r.client.GetFrameStats(plan.KasmID.ValueString(), plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Frame Stats",
			fmt.Sprintf("Could not get frame stats: %s", err),
		)
		return
	}

	// Set computed values
	plan.ResX = types.Int64Value(int64(frameStats.Frame.ResX))
	plan.ResY = types.Int64Value(int64(frameStats.Frame.ResY))
	plan.Changed = types.Int64Value(int64(frameStats.Frame.Changed))
	plan.ServerTime = types.Int64Value(int64(frameStats.Frame.ServerTime))
	plan.ClientCount = types.Int64Value(int64(len(frameStats.Frame.Clients)))
	plan.Analysis = types.Int64Value(int64(frameStats.Frame.Analysis))
	plan.Screenshot = types.Int64Value(int64(frameStats.Frame.Screenshot))
	plan.EncodingTime = types.Int64Value(int64(frameStats.Frame.EncodingTotal))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *statsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state statsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed stats
	frameStats, err := r.client.GetFrameStats(state.KasmID.ValueString(), state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Frame Stats",
			fmt.Sprintf("Could not read frame stats: %s", err),
		)
		return
	}

	// Update state with refreshed values
	state.ResX = types.Int64Value(int64(frameStats.Frame.ResX))
	state.ResY = types.Int64Value(int64(frameStats.Frame.ResY))
	state.Changed = types.Int64Value(int64(frameStats.Frame.Changed))
	state.ServerTime = types.Int64Value(int64(frameStats.Frame.ServerTime))
	state.ClientCount = types.Int64Value(int64(len(frameStats.Frame.Clients)))
	state.Analysis = types.Int64Value(int64(frameStats.Frame.Analysis))
	state.Screenshot = types.Int64Value(int64(frameStats.Frame.Screenshot))
	state.EncodingTime = types.Int64Value(int64(frameStats.Frame.EncodingTotal))
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *statsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan statsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed stats
	frameStats, err := r.client.GetFrameStats(plan.KasmID.ValueString(), plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Frame Stats",
			fmt.Sprintf("Could not update frame stats: %s", err),
		)
		return
	}

	// Update state with refreshed values
	plan.ResX = types.Int64Value(int64(frameStats.Frame.ResX))
	plan.ResY = types.Int64Value(int64(frameStats.Frame.ResY))
	plan.Changed = types.Int64Value(int64(frameStats.Frame.Changed))
	plan.ServerTime = types.Int64Value(int64(frameStats.Frame.ServerTime))
	plan.ClientCount = types.Int64Value(int64(len(frameStats.Frame.Clients)))
	plan.Analysis = types.Int64Value(int64(frameStats.Frame.Analysis))
	plan.Screenshot = types.Int64Value(int64(frameStats.Frame.Screenshot))
	plan.EncodingTime = types.Int64Value(int64(frameStats.Frame.EncodingTotal))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set updated state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *statsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Stats is a read-only resource, so we don't need to do anything for delete
	// Just log that we're "deleting" the resource
	var state statsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Stats resource delete complete", map[string]interface{}{
		"kasm_id": state.KasmID.ValueString(),
	})
}

// ImportState imports an existing resource into Terraform.
func (r *statsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using the kasm_id
	resource.ImportStatePassthroughID(ctx, path.Root("kasm_id"), req, resp)
}
