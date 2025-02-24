package join

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

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &joinResource{}
	_ resource.ResourceWithImportState = &joinResource{}
)

// joinResource is the resource implementation
type joinResource struct {
	client *client.Client
}

// JoinResourceModel maps the resource schema data
type JoinResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ShareID      types.String `tfsdk:"share_id"`
	UserID       types.String `tfsdk:"user_id"`
	KasmID       types.String `tfsdk:"kasm_id"`
	SessionToken types.String `tfsdk:"session_token"`
	KasmURL      types.String `tfsdk:"kasm_url"`
}

// New creates a new join resource
func New() resource.Resource {
	return &joinResource{}
}

// Metadata returns the resource type name
func (r *joinResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_join"
}

// Schema defines the schema for the resource
func (r *joinResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Join a shared Kasm session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"share_id": schema.StringAttribute{
				Required:    true,
				Description: "The share ID of the session to join",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the user joining the session",
			},
			"kasm_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the joined Kasm session",
			},
			"session_token": schema.StringAttribute{
				Computed:    true,
				Description: "The session token for the joined session",
			},
			"kasm_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL to access the joined session",
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *joinResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *joinResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan JoinResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set a unique ID before making the API call
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ShareID.ValueString(), plan.UserID.ValueString()))

	// Try to join the session with retries
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		joinResp, err := r.client.JoinKasm(plan.ShareID.ValueString(), plan.UserID.ValueString())
		if err != nil {
			if i == maxRetries-1 {
				resp.Diagnostics.AddError(
					"Error joining Kasm session",
					fmt.Sprintf("Could not join session after %d retries: %v", maxRetries, err),
				)
				return
			}
			tflog.Info(ctx, fmt.Sprintf("Join attempt %d failed, retrying in 2 seconds...", i+1))
			time.Sleep(2 * time.Second)
			continue
		}

		if joinResp == nil {
			if i == maxRetries-1 {
				resp.Diagnostics.AddError(
					"Error joining Kasm session",
					"Received nil response from join operation",
				)
				return
			}
			tflog.Info(ctx, fmt.Sprintf("Join attempt %d returned nil response, retrying in 2 seconds...", i+1))
			time.Sleep(2 * time.Second)
			continue
		}

		if joinResp.ErrorMessage != "" {
			if i == maxRetries-1 {
				resp.Diagnostics.AddError(
					"Error joining Kasm session",
					fmt.Sprintf("API returned error: %s", joinResp.ErrorMessage),
				)
				return
			}
			tflog.Info(ctx, fmt.Sprintf("Join attempt %d failed with error: %s, retrying in 2 seconds...", i+1, joinResp.ErrorMessage))
			time.Sleep(2 * time.Second)
			continue
		}

		// Success! Set the values from the response
		plan.KasmID = types.StringValue(joinResp.Kasm.KasmID)
		plan.SessionToken = types.StringValue(joinResp.SessionToken)
		plan.KasmURL = types.StringValue(joinResp.KasmURL)

		tflog.Info(ctx, fmt.Sprintf("Successfully joined Kasm session on attempt %d - ID: %s, KasmID: %s",
			i+1,
			plan.ID.ValueString(),
			plan.KasmID.ValueString()))

		// Set full state with all fields
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		return
	}
}

// Read refreshes the Terraform state with the latest data
func (r *joinResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state JoinResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Maintain the ID format
	if state.ID.IsNull() || state.ID.ValueString() == "" {
		state.ID = types.StringValue(fmt.Sprintf("%s:%s", state.ShareID.ValueString(), state.UserID.ValueString()))
	}

	// Try to join the session again to get fresh details
	joinResp, err := r.client.JoinKasm(state.ShareID.ValueString(), state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Session",
			fmt.Sprintf("Could not read session ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	if joinResp != nil && joinResp.ErrorMessage == "" {
		state.KasmID = types.StringValue(joinResp.Kasm.KasmID)
		state.SessionToken = types.StringValue(joinResp.SessionToken)
		state.KasmURL = types.StringValue(joinResp.KasmURL)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *joinResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan JoinResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Maintain the ID format
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.ShareID.ValueString(), plan.UserID.ValueString()))

	// Re-join the session with the new parameters
	joinResp, err := r.client.JoinKasm(plan.ShareID.ValueString(), plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error joining Kasm session",
			fmt.Sprintf("Could not join session: %v", err),
		)
		return
	}

	if joinResp != nil && joinResp.ErrorMessage == "" {
		plan.KasmID = types.StringValue(joinResp.Kasm.KasmID)
		plan.SessionToken = types.StringValue(joinResp.SessionToken)
		plan.KasmURL = types.StringValue(joinResp.KasmURL)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *joinResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No need to do anything as we're just joining an existing session
}

// ImportState imports the resource into Terraform state
func (r *joinResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("share_id"), req, resp)
}
