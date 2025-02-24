package staging

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &stagingResource{}
	_ resource.ResourceWithImportState = &stagingResource{}
)

// stagingResource is the resource implementation
type stagingResource struct {
	client *client.Client
}

// StagingResourceModel maps the resource schema data
type StagingResourceModel struct {
	ID                     types.String  `tfsdk:"id"`
	ZoneID                 types.String  `tfsdk:"zone_id"`
	ImageID                types.String  `tfsdk:"image_id"`
	NumSessions            types.Int64   `tfsdk:"num_sessions"`
	Expiration             types.Float64 `tfsdk:"expiration"`
	AllowKasmAudio         types.Bool    `tfsdk:"allow_kasm_audio"`
	AllowKasmUploads       types.Bool    `tfsdk:"allow_kasm_uploads"`
	AllowKasmDownloads     types.Bool    `tfsdk:"allow_kasm_downloads"`
	AllowKasmClipboardDown types.Bool    `tfsdk:"allow_kasm_clipboard_down"`
	AllowKasmClipboardUp   types.Bool    `tfsdk:"allow_kasm_clipboard_up"`
	AllowKasmMicrophone    types.Bool    `tfsdk:"allow_kasm_microphone"`
}

// New creates a new staging config resource
func New() resource.Resource {
	return &stagingResource{}
}

// Metadata returns the resource type name
func (r *stagingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_config"
}

// Schema defines the schema for the resource
func (r *stagingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kasm staging configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "The Zone ID for the staging config",
			},
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "The Image ID to use for the staging config",
			},
			"num_sessions": schema.Int64Attribute{
				Required:    true,
				Description: "The number of sessions to maintain staged",
			},
			"expiration": schema.Float64Attribute{
				Required:    true,
				Description: "The expiration time in hours for staged sessions",
			},
			"allow_kasm_audio": schema.BoolAttribute{
				Required:    true,
				Description: "Allow audio streaming from the session",
			},
			"allow_kasm_uploads": schema.BoolAttribute{
				Required:    true,
				Description: "Allow file uploads to the session",
			},
			"allow_kasm_downloads": schema.BoolAttribute{
				Required:    true,
				Description: "Allow file downloads from the session",
			},
			"allow_kasm_clipboard_down": schema.BoolAttribute{
				Required:    true,
				Description: "Allow clipboard copy from session to local",
			},
			"allow_kasm_clipboard_up": schema.BoolAttribute{
				Required:    true,
				Description: "Allow clipboard copy from local to session",
			},
			"allow_kasm_microphone": schema.BoolAttribute{
				Required:    true,
				Description: "Allow microphone access in the session",
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *stagingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state
func (r *stagingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StagingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateStagingConfigRequest{
		ZoneID:                 plan.ZoneID.ValueString(),
		ImageID:                plan.ImageID.ValueString(),
		NumSessions:            int(plan.NumSessions.ValueInt64()),
		Expiration:             plan.Expiration.ValueFloat64(),
		AllowKasmAudio:         plan.AllowKasmAudio.ValueBool(),
		AllowKasmUploads:       plan.AllowKasmUploads.ValueBool(),
		AllowKasmDownloads:     plan.AllowKasmDownloads.ValueBool(),
		AllowKasmClipboardDown: plan.AllowKasmClipboardDown.ValueBool(),
		AllowKasmClipboardUp:   plan.AllowKasmClipboardUp.ValueBool(),
		AllowKasmMicrophone:    plan.AllowKasmMicrophone.ValueBool(),
	}

	stagingConfig, err := r.client.CreateStagingConfig(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating staging config",
			fmt.Sprintf("Could not create staging config: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(stagingConfig.StagingConfigID)

	tflog.Info(ctx, fmt.Sprintf("Created staging config with ID: %s", plan.ID.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *stagingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StagingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stagingConfig, err := r.client.GetStagingConfig(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Staging Config",
			fmt.Sprintf("Could not read staging config ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.ZoneID = types.StringValue(stagingConfig.ZoneID)
	state.ImageID = types.StringValue(stagingConfig.ImageID)
	state.NumSessions = types.Int64Value(int64(stagingConfig.NumSessions))
	state.Expiration = types.Float64Value(stagingConfig.Expiration)
	state.AllowKasmAudio = types.BoolValue(stagingConfig.AllowKasmAudio)
	state.AllowKasmUploads = types.BoolValue(stagingConfig.AllowKasmUploads)
	state.AllowKasmDownloads = types.BoolValue(stagingConfig.AllowKasmDownloads)
	state.AllowKasmClipboardDown = types.BoolValue(stagingConfig.AllowKasmClipboardDown)
	state.AllowKasmClipboardUp = types.BoolValue(stagingConfig.AllowKasmClipboardUp)
	state.AllowKasmMicrophone = types.BoolValue(stagingConfig.AllowKasmMicrophone)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *stagingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StagingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateStagingConfigRequest{
		StagingConfigID:        plan.ID.ValueString(),
		ZoneID:                 plan.ZoneID.ValueString(),
		ImageID:                plan.ImageID.ValueString(),
		NumSessions:            new(int),
		Expiration:             new(float64),
		AllowKasmAudio:         new(bool),
		AllowKasmUploads:       new(bool),
		AllowKasmDownloads:     new(bool),
		AllowKasmClipboardDown: new(bool),
		AllowKasmClipboardUp:   new(bool),
		AllowKasmMicrophone:    new(bool),
	}

	*updateReq.NumSessions = int(plan.NumSessions.ValueInt64())
	*updateReq.Expiration = plan.Expiration.ValueFloat64()
	*updateReq.AllowKasmAudio = plan.AllowKasmAudio.ValueBool()
	*updateReq.AllowKasmUploads = plan.AllowKasmUploads.ValueBool()
	*updateReq.AllowKasmDownloads = plan.AllowKasmDownloads.ValueBool()
	*updateReq.AllowKasmClipboardDown = plan.AllowKasmClipboardDown.ValueBool()
	*updateReq.AllowKasmClipboardUp = plan.AllowKasmClipboardUp.ValueBool()
	*updateReq.AllowKasmMicrophone = plan.AllowKasmMicrophone.ValueBool()

	stagingConfig, err := r.client.UpdateStagingConfig(updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kasm Staging Config",
			fmt.Sprintf("Could not update staging config ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(stagingConfig.StagingConfigID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *stagingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StagingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStagingConfig(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kasm Staging Config",
			fmt.Sprintf("Could not delete staging config ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *stagingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
