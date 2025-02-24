package group

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
	"strings"
	"terraform-provider-kasm/internal/client"
)

var _ resource.Resource = &groupResource{}
var _ resource.ResourceWithImportState = &groupResource{}

type groupResource struct {
	client *client.Client
}

type GroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Priority    types.Int64  `tfsdk:"priority"`
	Description types.String `tfsdk:"description"`
	Permissions types.List   `tfsdk:"permissions"`
}

func New() resource.Resource {
	return &groupResource{}
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"priority": schema.Int64Attribute{
				Required: true,
			},
			"permissions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert permissions from types.List to []string
	permissions := make([]string, 0)
	if !plan.Permissions.IsNull() {
		var permissionValues []string
		diags = plan.Permissions.ElementsAs(ctx, &permissionValues, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		permissions = permissionValues
	}

	group := &client.Group{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    int(plan.Priority.ValueInt64()),
		Permissions: permissions,
	}

	createdGroup, err := r.client.CreateGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			fmt.Sprintf("Could not create group, unexpected error: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(createdGroup.GroupID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading group", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	group, err := r.client.GetGroup(state.ID.ValueString())
	if err != nil {
		if client.IsGroupNotFoundError(err) {
			tflog.Debug(ctx, "Group not found, removing from state", map[string]interface{}{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading group",
			fmt.Sprintf("Could not read group ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	state.Name = types.StringValue(group.Name)
	state.Description = types.StringValue(group.Description)
	state.Priority = types.Int64Value(int64(group.Priority))

	// Convert permissions from []string to types.List
	if len(group.Permissions) > 0 {
		permissionValues := make([]types.String, len(group.Permissions))
		for i, p := range group.Permissions {
			permissionValues[i] = types.StringValue(p)
		}
		state.Permissions, diags = types.ListValueFrom(ctx, types.StringType, permissionValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert permissions from types.List to []string
	permissions := make([]string, 0)
	if !plan.Permissions.IsNull() {
		var permissionValues []string
		diags = plan.Permissions.ElementsAs(ctx, &permissionValues, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		permissions = permissionValues
	}

	group := &client.Group{
		GroupID:     plan.ID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    int(plan.Priority.ValueInt64()),
		Permissions: permissions,
	}

	updatedGroup, err := r.client.UpdateGroup(group)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group",
			fmt.Sprintf("Could not update group, unexpected error: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(updatedGroup.GroupID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting group",
			fmt.Sprintf("Could not delete group, unexpected error: %v", err),
		)
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	var groupID string
	switch len(idParts) {
	case 1:
		// Import by group ID
		groupID = idParts[0]
	case 2:
		if idParts[0] != "name" {
			resp.Diagnostics.AddError(
				"Invalid Import Format",
				`To import by name use format: name:group-name`,
			)
			return
		}
		// Find group by name
		groups, err := r.client.GetGroups()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing group",
				fmt.Sprintf("Could not list groups: %s", err),
			)
			return
		}

		for _, group := range groups {
			if group.Name == idParts[1] {
				groupID = group.GroupID
				break
			}
		}

		if groupID == "" {
			resp.Diagnostics.AddError(
				"Group Not Found",
				fmt.Sprintf("Could not find group with name: %s", idParts[1]),
			)
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Invalid Import Format",
			`Use format: "group_id" or "name:group-name"`,
		)
		return
	}

	// Set the imported ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), groupID)...)
}
