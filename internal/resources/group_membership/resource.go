package group_membership

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &groupMembershipResource{}
	_ resource.ResourceWithImportState = &groupMembershipResource{}
)

// groupMembershipResource is the resource implementation
type groupMembershipResource struct {
	client *client.Client
}

// GroupMembershipResourceModel maps the resource schema data
type GroupMembershipResourceModel struct {
	ID      types.String `tfsdk:"id"`
	GroupID types.String `tfsdk:"group_id"`
	UserID  types.String `tfsdk:"user_id"`
}

// New creates a new group membership resource
func New() resource.Resource {
	return &groupMembershipResource{}
}

// Metadata returns the resource type name
func (r *groupMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_membership"
}

// Schema defines the schema for the resource
func (r *groupMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages user membership in a group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *groupMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *groupMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupMembershipResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AddUserToGroup(plan.UserID.ValueString(), plan.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error adding user to group",
			fmt.Sprintf("Could not add user to group: %v", err),
		)
		return
	}

	// Debug: Get user details after adding to group
	user, err := r.client.GetUser(plan.UserID.ValueString())
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not get user details after adding to group: %v", err))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("User details after adding to group: %+v", user))
		if user.Groups != nil {
			groupNames := []string{}
			for _, group := range user.Groups {
				groupNames = append(groupNames, group.Name)
			}
			tflog.Debug(ctx, fmt.Sprintf("User groups after adding to group: %+v", groupNames))
		}
		if user.AuthorizedImages != nil {
			tflog.Debug(ctx, fmt.Sprintf("User authorized images after adding to group: %+v", user.AuthorizedImages))
		}

		// NOTE: The Kasm API automatically updates the user's authorized_images and groups
		// when a user is added to a group. This can cause state drift in the user resource.
		// In a real-world scenario, we would need to coordinate with the user resource
		// to update its state. For testing purposes, we'll just log this information.
	}

	// Set a unique ID for the membership
	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.GroupID.ValueString(), plan.UserID.ValueString()))

	tflog.Info(ctx, fmt.Sprintf("Created group membership with ID: %s", plan.ID.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *groupMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupMembershipResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the user's groups
	user, err := r.client.GetUser(state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading User",
			fmt.Sprintf("Could not read user: %s", err),
		)
		return
	}

	// Check if the user is still in the specific group
	var found bool
	for _, group := range user.Groups {
		if group.GroupID == state.GroupID.ValueString() {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *groupMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Group memberships cannot be updated, only created or deleted
	var plan GroupMembershipResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *groupMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupMembershipResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveUserFromGroup(state.UserID.ValueString(), state.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing User from Group",
			fmt.Sprintf("Could not remove user from group: %s", err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *groupMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format group_id:user_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
