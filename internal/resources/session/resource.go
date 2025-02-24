package session

import (
	"context"
	"fmt"

	"terraform-provider-kasm/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &sessionTokenResource{}
	_ resource.ResourceWithImportState = &sessionTokenResource{}
)

// sessionTokenResource is the resource implementation
type sessionTokenResource struct {
	client *client.Client
}

// SessionTokenResourceModel maps the resource schema data
type SessionTokenResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	UserID                types.String `tfsdk:"user_id"`
	SessionToken          types.String `tfsdk:"session_token"`
	SessionTokenDate      types.String `tfsdk:"session_token_date"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	SessionJWT            types.String `tfsdk:"session_jwt"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	AllowResume           types.Bool   `tfsdk:"allow_resume"`
	SessionAuthentication types.Bool   `tfsdk:"session_authentication"`
}

// New creates a new session token resource
func New() resource.Resource {
	return &sessionTokenResource{}
}

// Metadata returns the resource type name
func (r *sessionTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_token"
}

// Schema defines the schema for the resource
func (r *sessionTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kasm session token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the user to create the session token for",
			},
			"session_token": schema.StringAttribute{
				Computed:    true,
				Description: "The session token value",
			},
			"session_token_date": schema.StringAttribute{
				Computed:    true,
				Description: "The time the token was created or last promoted",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "The time the token will expire",
			},
			"session_jwt": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The JWT token used for authentication",
			},
			"persistent": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the session should be persistent",
			},
			"allow_resume": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the session can be resumed",
			},
			"session_authentication": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether session authentication is required",
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *sessionTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *sessionTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SessionTokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateSessionTokenRequest{}
	createReq.TargetUser.UserID = plan.UserID.ValueString()

	sessionToken, err := r.client.CreateSessionToken(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating session token",
			fmt.Sprintf("Could not create session token: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(sessionToken.SessionToken)
	plan.SessionToken = types.StringValue(sessionToken.SessionToken)
	plan.SessionTokenDate = types.StringValue(sessionToken.SessionTokenDate)
	plan.ExpiresAt = types.StringValue(sessionToken.ExpiresAt)
	plan.SessionJWT = types.StringValue(sessionToken.SessionJWT)

	tflog.Info(ctx, fmt.Sprintf("Created session token with ID: %s", plan.ID.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *sessionTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SessionTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getReq := &client.GetSessionTokenRequest{}
	getReq.TargetSessionToken.SessionToken = state.SessionToken.ValueString()

	sessionToken, err := r.client.GetSessionToken(getReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Session Token",
			fmt.Sprintf("Could not read session token ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	state.SessionTokenDate = types.StringValue(sessionToken.SessionTokenDate)
	state.ExpiresAt = types.StringValue(sessionToken.ExpiresAt)
	state.SessionJWT = types.StringValue(sessionToken.SessionJWT)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *sessionTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SessionTokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateSessionTokenRequest{}
	updateReq.TargetSessionToken.SessionToken = plan.SessionToken.ValueString()

	sessionToken, err := r.client.UpdateSessionToken(updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kasm Session Token",
			fmt.Sprintf("Could not update session token ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	plan.SessionTokenDate = types.StringValue(sessionToken.SessionTokenDate)
	plan.ExpiresAt = types.StringValue(sessionToken.ExpiresAt)
	plan.SessionJWT = types.StringValue(sessionToken.SessionJWT)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *sessionTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SessionTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteReq := &client.DeleteSessionTokenRequest{}
	deleteReq.TargetSessionToken.SessionToken = state.SessionToken.ValueString()

	err := r.client.DeleteSessionToken(deleteReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kasm Session Token",
			fmt.Sprintf("Could not delete session token ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *sessionTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("session_token"), req, resp)
}
