package kasm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-kasm/internal/client"
)

var (
	_ resource.Resource                = &kasmSessionResource{}
	_ resource.ResourceWithConfigure   = &kasmSessionResource{}
	_ resource.ResourceWithImportState = &kasmSessionResource{}
)

type kasmSessionResource struct {
	client *client.Client
}

type kasmSessionResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	ImageID               types.String `tfsdk:"image_id"`
	UserID                types.String `tfsdk:"user_id"`
	Share                 types.Bool   `tfsdk:"share"`
	EnableSharing         types.Bool   `tfsdk:"enable_sharing"`
	ShareID               types.String `tfsdk:"share_id"`
	RDPEnabled            types.Bool   `tfsdk:"rdp_enabled"`
	RDPConnectionFile     types.String `tfsdk:"rdp_connection_file"`
	EnableStats           types.Bool   `tfsdk:"enable_stats"`
	AllowExec             types.Bool   `tfsdk:"allow_exec"`
	OperationalStatus     types.String `tfsdk:"operational_status"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	AllowResume           types.Bool   `tfsdk:"allow_resume"`
	SessionAuthentication types.Bool   `tfsdk:"session_authentication"`
}

func NewKasmSessionResource() resource.Resource {
	tflog.Info(context.Background(), "Creating new kasm session resource")
	return &kasmSessionResource{
		client: nil,
	}
}

func (r *kasmSessionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(context.Background(), "Configuring kasm session resource")

	if req.ProviderData == nil {
		// During validation, provider data might be nil, which is okay
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
	tflog.Info(context.Background(), "Successfully configured kasm session resource")
}

func (r *kasmSessionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session"
}

func (r *kasmSessionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kasm session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the Kasm session.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Description: "The ID of the image to use for the session.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user to create the session for.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"share": schema.BoolAttribute{
				Description: "Whether to enable session sharing.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"enable_sharing": schema.BoolAttribute{
				Description: "Whether to enable sharing features.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"share_id": schema.StringAttribute{
				Description: "The share ID for the session when sharing is enabled.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rdp_enabled": schema.BoolAttribute{
				Description: "Whether to enable RDP for the session.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"rdp_connection_file": schema.StringAttribute{
				Description: "The RDP connection file content when RDP is enabled.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_stats": schema.BoolAttribute{
				Description: "Whether to enable session statistics.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_exec": schema.BoolAttribute{
				Description: "Whether to allow command execution in the session.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"operational_status": schema.StringAttribute{
				Description: "The operational status of the session.",
				Computed:    true,
			},
			"persistent": schema.BoolAttribute{
				Description: "Whether the session should be persistent.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_resume": schema.BoolAttribute{
				Description: "Whether the session can be resumed.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"session_authentication": schema.BoolAttribute{
				Description: "Whether session authentication is required.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *kasmSessionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating kasm session")

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. This is a bug in the provider that should be reported to the provider developers.",
		)
		return
	}

	var plan kasmSessionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If share is true, ensure enable_sharing is also true
	if plan.Share.ValueBool() {
		plan.EnableSharing = types.BoolValue(true)
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating session with user_id: %s, image_id: %s",
		plan.UserID.ValueString(),
		plan.ImageID.ValueString()))

	// Create the session
	status, err := r.client.CreateKasm(
		plan.UserID.ValueString(),
		plan.ImageID.ValueString(),
		"",
		"",
		plan.Share.ValueBool(),
		plan.Persistent.ValueBool(),
		plan.AllowResume.ValueBool(),
		plan.SessionAuthentication.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Kasm session",
			fmt.Sprintf("Unable to create session: %v", err),
		)
		return
	}

	// Debug log the response
	tflog.Debug(ctx, fmt.Sprintf("CreateKasm response: %+v", status))

	// Set ID and ShareID from initial response
	plan.ID = types.StringValue(status.KasmID)
	plan.OperationalStatus = types.StringValue(status.Status)

	// Set share_id from initial response if available
	if status.ShareID != "" {
		plan.ShareID = types.StringValue(status.ShareID)
		tflog.Debug(ctx, fmt.Sprintf("Got share_id from initial response: %s", status.ShareID))
	} else {
		plan.ShareID = types.StringValue("")
	}

	plan.RDPConnectionFile = types.StringValue("") // Initialize empty

	// Get full session details
	tflog.Debug(ctx, "Getting session details")
	sessionInfo, err := r.client.GetKasmStatus(
		plan.UserID.ValueString(),
		status.KasmID,
		true,
	)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Error getting session details: %v", err))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Session details: %+v", sessionInfo))

		if sessionInfo.Kasm != nil {
			// Update operational status from fresh data
			plan.OperationalStatus = types.StringValue(sessionInfo.Kasm.OperationalStatus)

			// Set share_id if available and not already set
			if sessionInfo.Kasm.ShareID != "" && plan.ShareID.ValueString() == "" {
				plan.ShareID = types.StringValue(sessionInfo.Kasm.ShareID)
				tflog.Debug(ctx, fmt.Sprintf("Got share_id from session details: %s", sessionInfo.Kasm.ShareID))
			}
		}
	}

	// If we don't have a share_id yet but sharing is enabled, wait and retry
	if plan.Share.ValueBool() && plan.ShareID.ValueString() == "" {
		tflog.Debug(ctx, "Waiting for share_id to be available")
		maxRetries := 10 // Increase retries
		for i := 0; i < maxRetries; i++ {
			time.Sleep(2 * time.Second)

			sessionInfo, err := r.client.GetKasmStatus(
				plan.UserID.ValueString(),
				status.KasmID,
				true,
			)
			if err == nil && sessionInfo.Kasm != nil && sessionInfo.Kasm.ShareID != "" {
				plan.ShareID = types.StringValue(sessionInfo.Kasm.ShareID)
				tflog.Debug(ctx, fmt.Sprintf("Got share_id after retry: %s", sessionInfo.Kasm.ShareID))
				break
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Retry %d/%d: No share_id yet", i+1, maxRetries))
			}
		}

		if plan.ShareID.ValueString() == "" {
			tflog.Warn(ctx, "Failed to get share_id after retries")
		}
	}

	// Handle RDP if enabled
	if plan.RDPEnabled.ValueBool() {
		tflog.Debug(ctx, "Getting RDP connection info")
		rdpResp, err := r.client.GetRDPConnectionInfo(
			plan.UserID.ValueString(),
			status.KasmID,
			client.RDPConnectionTypeFile,
		)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Unable to get RDP connection info: %v", err))
			plan.RDPConnectionFile = types.StringValue("")
		} else if rdpResp != nil {
			plan.RDPConnectionFile = types.StringValue(rdpResp.File)
			tflog.Debug(ctx, "Successfully got RDP connection file")
		} else {
			plan.RDPConnectionFile = types.StringValue("")
			tflog.Warn(ctx, "RDP response was nil")
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Created session with ID: %s, ShareID: %s, Status: %s",
		plan.ID.ValueString(),
		plan.ShareID.ValueString(),
		plan.OperationalStatus.ValueString()))

	// Set the final state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *kasmSessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. This is a bug in the provider that should be reported to the provider developers.",
		)
		return
	}

	var state kasmSessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get session status
	status, err := r.client.GetKasmStatus(
		state.UserID.ValueString(),
		state.ID.ValueString(),
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Kasm session",
			fmt.Sprintf("Unable to read session: %v", err),
		)
		return
	}

	// Update state with status from response
	if status.Kasm != nil {
		state.OperationalStatus = types.StringValue(status.Kasm.OperationalStatus)
		if status.Kasm.ShareID != "" {
			state.ShareID = types.StringValue(status.Kasm.ShareID)
		}
	} else {
		// If Kasm is nil, use the top-level status
		state.OperationalStatus = types.StringValue(status.OperationalStatus)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *kasmSessionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. This is a bug in the provider that should be reported to the provider developers.",
		)
		return
	}

	var plan kasmSessionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state kasmSessionResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update RDP status if changed
	if !plan.RDPEnabled.Equal(state.RDPEnabled) && plan.RDPEnabled.ValueBool() {
		rdpResp, err := r.client.GetRDPConnectionInfo(
			plan.UserID.ValueString(),
			state.ID.ValueString(),
			client.RDPConnectionTypeFile,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating RDP connection info",
				fmt.Sprintf("Unable to update RDP connection info: %v", err),
			)
			return
		}
		plan.RDPConnectionFile = types.StringValue(rdpResp.File)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *kasmSessionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client Not Configured",
			"The provider client has not been configured. This is a bug in the provider that should be reported to the provider developers.",
		)
		return
	}

	var state kasmSessionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Destroy the session
	err := r.client.DestroyKasm(state.UserID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error destroying Kasm session",
			fmt.Sprintf("Unable to destroy session: %v", err),
		)
		return
	}
}

func (r *kasmSessionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using the format: user_id:kasm_id
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: user_id:kasm_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
