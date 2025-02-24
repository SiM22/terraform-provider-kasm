package session_permission

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &sessionPermissionResource{}
	_ resource.ResourceWithImportState = &sessionPermissionResource{}
)

// sessionPermissionResource is the resource implementation
type sessionPermissionResource struct {
	client *client.Client
}

// UserPermissionModel maps the user permission schema data
type UserPermissionModel struct {
	UserID      types.String `tfsdk:"user_id"`
	Access      types.String `tfsdk:"access"`
	Username    types.String `tfsdk:"username"`
	VNCUsername types.String `tfsdk:"vnc_username"`
}

// SessionPermissionModel maps the resource schema data
type SessionPermissionModel struct {
	ID              types.String          `tfsdk:"id"`
	KasmID          types.String          `tfsdk:"kasm_id"`
	GlobalAccess    types.String          `tfsdk:"global_access"`
	UserPermissions []UserPermissionModel `tfsdk:"user_permissions"`
}

// New creates a new session permission resource
func New() resource.Resource {
	return &sessionPermissionResource{}
}

// Metadata returns the resource type name
func (r *sessionPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_permission"
}

// Schema defines the schema for the resource
func (r *sessionPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Kasm session permissions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kasm_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Kasm session",
			},
			"global_access": schema.StringAttribute{
				Optional:    true,
				Description: "The global access level (r, rw)",
			},
			"user_permissions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of user-specific permissions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user_id": schema.StringAttribute{
							Required:    true,
							Description: "The ID of the user",
						},
						"access": schema.StringAttribute{
							Required:    true,
							Description: "The access level (r, rw)",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "The username of the user",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"vnc_username": schema.StringAttribute{
							Computed:    true,
							Description: "The VNC username for the user",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *sessionPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *sessionPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SessionPermissionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID to the Kasm ID
	plan.ID = plan.KasmID

	// Convert the plan to API request
	permissions := []client.SessionPermissionAccess{}
	for _, p := range plan.UserPermissions {
		permissions = append(permissions, client.SessionPermissionAccess{
			UserID: p.UserID.ValueString(),
			Access: p.Access.ValueString(),
		})
	}

	request := &client.SetSessionPermissionsRequest{
		TargetSessionPermissions: client.TargetSessionPermissions{
			KasmID:             plan.KasmID.ValueString(),
			Access:             plan.GlobalAccess.ValueString(),
			SessionPermissions: permissions,
		},
	}

	// Create the session permission
	perms, err := r.client.SetSessionPermissions(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating session permission",
			fmt.Sprintf("Could not create session permission: %v", err),
		)
		return
	}

	// Update the plan with the response
	if len(perms) > 0 {
		plan.UserPermissions = make([]UserPermissionModel, len(perms))
		for i, p := range perms {
			plan.UserPermissions[i] = UserPermissionModel{
				UserID:      types.StringValue(p.UserID),
				Access:      types.StringValue(p.Access),
				Username:    types.StringValue(p.Username),
				VNCUsername: types.StringValue(p.VNCUsername),
			}
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *sessionPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SessionPermissionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to preserve access levels
	var currentState SessionPermissionModel
	diags = req.State.Get(ctx, &currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a map of user IDs to their access levels from current state
	userAccessLevels := make(map[string]string)
	for _, p := range currentState.UserPermissions {
		userAccessLevels[p.UserID.ValueString()] = p.Access.ValueString()
	}

	// Get the permissions
	request := &client.GetSessionPermissionsRequest{
		TargetSessionPermissions: struct {
			KasmID string `json:"kasm_id"`
		}{
			KasmID: state.KasmID.ValueString(),
		},
	}

	perms, err := r.client.GetSessionPermissions(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading session permission",
			fmt.Sprintf("Could not read session permission: %v", err),
		)
		return
	}

	// Update state with the response
	if len(perms) > 0 {
		state.UserPermissions = make([]UserPermissionModel, len(perms))
		for i, p := range perms {
			// Use access level from current state if available
			access := p.Access
			if currentAccess, ok := userAccessLevels[p.UserID]; ok {
				access = currentAccess
			}

			state.UserPermissions[i] = UserPermissionModel{
				UserID:      types.StringValue(p.UserID),
				Access:      types.StringValue(access),
				Username:    types.StringValue(p.Username),
				VNCUsername: types.StringValue(p.VNCUsername),
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *sessionPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SessionPermissionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to preserve computed fields
	var state SessionPermissionModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a map of user IDs to their computed fields
	userComputedFields := make(map[string]struct {
		username    string
		vncUsername string
	})
	for _, p := range state.UserPermissions {
		userComputedFields[p.UserID.ValueString()] = struct {
			username    string
			vncUsername string
		}{
			username:    p.Username.ValueString(),
			vncUsername: p.VNCUsername.ValueString(),
		}
	}

	// Convert the plan to API request
	permissions := []client.SessionPermissionAccess{}
	for _, p := range plan.UserPermissions {
		permissions = append(permissions, client.SessionPermissionAccess{
			UserID: p.UserID.ValueString(),
			Access: p.Access.ValueString(),
		})
	}

	request := &client.SetSessionPermissionsRequest{
		TargetSessionPermissions: client.TargetSessionPermissions{
			KasmID:             plan.KasmID.ValueString(),
			Access:             plan.GlobalAccess.ValueString(),
			SessionPermissions: permissions,
		},
	}

	// Update the session permission
	perms, err := r.client.SetSessionPermissions(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating session permission",
			fmt.Sprintf("Could not update session permission: %v", err),
		)
		return
	}

	// Update the plan with the response
	if len(perms) > 0 {
		plan.UserPermissions = make([]UserPermissionModel, len(perms))
		for i, p := range perms {
			// Use computed fields from state if available
			username := p.Username
			vncUsername := p.VNCUsername
			if computed, ok := userComputedFields[p.UserID]; ok {
				username = computed.username
				vncUsername = computed.vncUsername
			}

			plan.UserPermissions[i] = UserPermissionModel{
				UserID:      types.StringValue(p.UserID),
				Access:      types.StringValue(p.Access),
				Username:    types.StringValue(username),
				VNCUsername: types.StringValue(vncUsername),
			}
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *sessionPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SessionPermissionModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the session permission
	request := &client.DeleteAllSessionPermissionsRequest{
		TargetSessionPermissions: struct {
			KasmID string `json:"kasm_id"`
		}{
			KasmID: state.KasmID.ValueString(),
		},
	}

	err := r.client.DeleteAllSessionPermissions(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting session permission",
			fmt.Sprintf("Could not delete session permission: %v", err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *sessionPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("kasm_id"), req, resp)
}
