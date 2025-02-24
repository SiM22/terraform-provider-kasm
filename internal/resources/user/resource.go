package user

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

var _ resource.Resource = &userResource{}
var _ resource.ResourceWithImportState = &userResource{}

type userResource struct {
	client *client.Client
}

type UserResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	FirstName        types.String `tfsdk:"first_name"`
	LastName         types.String `tfsdk:"last_name"`
	Locked           types.Bool   `tfsdk:"locked"`
	Disabled         types.Bool   `tfsdk:"disabled"`
	Organization     types.String `tfsdk:"organization"`
	Phone            types.String `tfsdk:"phone"`
	Groups           types.List   `tfsdk:"groups"`
	Attributes       types.Map    `tfsdk:"attributes"`
	AuthorizedImages types.List   `tfsdk:"authorized_images"`
}

// New/Constructor function
func New() resource.Resource {
	return &userResource{}
}

// Metadata implementation
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema implementation
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"locked": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"organization": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phone": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"groups": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of group names to assign the user to",
			},
			"attributes": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "User-specific attributes",
			},
			"authorized_images": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of image IDs the user is authorized to use",
			},
		},
	}
}

// Configure implementation
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Helper functions
func (r *userResource) retryOperation(ctx context.Context, operation func() error) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		if i < maxRetries-1 {
			tflog.Warn(ctx, fmt.Sprintf("Operation failed, retrying (%d/%d)", i+1, maxRetries))
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return err
	}
	return nil
}

// Helper functions
func resourceUserToState(ctx context.Context, user *client.User, state *UserResourceModel) {
	state.ID = types.StringValue(user.UserID)
	state.Username = types.StringValue(user.Username)
	state.FirstName = types.StringValue(user.FirstName)
	state.LastName = types.StringValue(user.LastName)
	state.Locked = types.BoolValue(user.Locked)
	state.Disabled = types.BoolValue(user.Disabled)

	// Handle optional string fields
	if user.Organization == "" {
		state.Organization = types.StringNull()
	} else {
		state.Organization = types.StringValue(user.Organization)
	}

	if user.Phone == "" {
		state.Phone = types.StringNull()
	} else {
		state.Phone = types.StringValue(user.Phone)
	}

	// Password is handled separately and only updated when changed
	// If password is null, it will be handled by the update function

	// Always ensure groups is a non-null list
	groupNames := []string{}
	for _, group := range user.Groups {
		if !group.IsSystem {
			groupNames = append(groupNames, group.Name)
		}
	}

	groupsList, diags := types.ListValueFrom(ctx, types.StringType, groupNames)
	if diags.HasError() {
		groupsList, _ = types.ListValueFrom(ctx, types.StringType, []string{})
	}
	state.Groups = groupsList

	// Ensure authorized_images is always a non-null list
	if user.AuthorizedImages != nil {
		imageList, diags := types.ListValueFrom(ctx, types.StringType, user.AuthorizedImages)
		if diags.HasError() {
			imageList, _ = types.ListValueFrom(ctx, types.StringType, []string{})
		}
		state.AuthorizedImages = imageList
	} else {
		emptyList, _ := types.ListValueFrom(ctx, types.StringType, []string{})
		state.AuthorizedImages = emptyList
	}
}

func (r *userResource) handleGroupUpdates(ctx context.Context, userID string, groups types.List) error {
	if groups.IsNull() || groups.IsUnknown() {
		// Handle null or unknown values by setting empty list
		return r.client.UpdateUserGroupsByName(userID, []string{})
	}

	var groupStrings []types.String
	diags := groups.ElementsAs(ctx, &groupStrings, false)
	if diags.HasError() {
		return fmt.Errorf("invalid group format: %v", diags)
	}

	groupNames := make([]string, 0, len(groupStrings))
	for _, g := range groupStrings {
		if !g.IsNull() && !g.IsUnknown() {
			groupNames = append(groupNames, g.ValueString())
		}
	}

	return r.client.UpdateUserGroupsByName(userID, groupNames)
}

// Create implementation
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating user", map[string]interface{}{
		"username": plan.Username.ValueString(),
	})

	// Initialize empty lists if they are null or unknown
	if plan.Groups.IsNull() || plan.Groups.IsUnknown() {
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Groups = emptyList
	}

	if plan.AuthorizedImages.IsNull() || plan.AuthorizedImages.IsUnknown() {
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.AuthorizedImages = emptyList
	}

	// Create the user first
	user := &client.User{
		Username:     plan.Username.ValueString(),
		Password:     plan.Password.ValueString(),
		FirstName:    plan.FirstName.ValueString(),
		LastName:     plan.LastName.ValueString(),
		Locked:       plan.Locked.ValueBool(),
		Disabled:     plan.Disabled.ValueBool(),
		Organization: plan.Organization.ValueString(),
		Phone:        plan.Phone.ValueString(),
	}

	// Add authorized images if specified
	if !plan.AuthorizedImages.IsNull() {
		var imageIDs []string
		diags = plan.AuthorizedImages.ElementsAs(ctx, &imageIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		user.AuthorizedImages = imageIDs
	}

	var createdUser *client.User
	err := r.retryOperation(ctx, func() error {
		var err error
		createdUser, err = r.client.CreateUser(user)
		return err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			fmt.Sprintf("Could not create user: %v", err),
		)
		return
	}

	tflog.Debug(ctx, "Created user successfully", map[string]interface{}{
		"user_id": createdUser.UserID,
	})

	// Handle attributes if specified
	if !plan.Attributes.IsNull() {
		var attributes map[string]string
		diags = plan.Attributes.ElementsAs(ctx, &attributes, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Convert string map to interface map
		attributesInterface := make(map[string]interface{})
		for k, v := range attributes {
			attributesInterface[k] = v
		}

		err = r.client.UpdateUserAttributes(createdUser.UserID, attributesInterface)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting user attributes",
				fmt.Sprintf("Could not set user attributes: %v", err),
			)
			return
		}
	}

	// Handle group assignments
	if !plan.Groups.IsNull() {
		err = r.handleGroupUpdates(ctx, createdUser.UserID, plan.Groups)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error assigning user to groups",
				err.Error(),
			)
			return
		}
	}

	// Read the user back to get the final state
	user, err = r.client.GetUser(createdUser.UserID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created user",
			fmt.Sprintf("Could not read created user: %v", err),
		)
		return
	}

	// Set the state
	plan.ID = types.StringValue(createdUser.UserID)
	resourceUserToState(ctx, user, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, "Finished creating user", map[string]interface{}{
		"user_id": user.UserID,
	})
}

// Read implementation
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Starting Read method for user resource")

	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize empty lists if null
	if state.Groups.IsNull() {
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Groups = emptyList
	}

	if state.AuthorizedImages.IsNull() {
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.AuthorizedImages = emptyList
	}

	var user *client.User
	err := r.retryOperation(ctx, func() error {
		var err error
		user, err = r.client.GetUser(state.ID.ValueString())
		return err
	})

	if err != nil {
		if client.IsUserNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	// Update state with user data
	resourceUserToState(ctx, user, &state)

	// Ensure consistent handling of groups
	if state.Groups.IsNull() || len(user.Groups) == 0 {
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Groups = emptyList
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implementation
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create update payload
	user := &client.User{
		UserID:       plan.ID.ValueString(),
		Username:     plan.Username.ValueString(),
		FirstName:    plan.FirstName.ValueString(),
		LastName:     plan.LastName.ValueString(),
		Organization: plan.Organization.ValueString(),
		Phone:        plan.Phone.ValueString(),
		Locked:       plan.Locked.ValueBool(),
		Disabled:     plan.Disabled.ValueBool(),
	}

	// Only include password if changed
	if !plan.Password.Equal(state.Password) {
		user.Password = plan.Password.ValueString()
	}

	var updatedUser *client.User
	err := r.retryOperation(ctx, func() error {
		var err error
		updatedUser, err = r.client.UpdateUser(user)
		return err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			fmt.Sprintf("Could not update user: %v", err),
		)
		return
	}

	// Update groups first if specified
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		var groupNames []string
		diags := plan.Groups.ElementsAs(ctx, &groupNames, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		err = r.client.UpdateUserGroupsByName(updatedUser.UserID, groupNames)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating user groups",
				fmt.Sprintf("Could not update user groups: %v", err),
			)
			return
		}

		// Convert group names to List
		groupList, diags := types.ListValueFrom(ctx, types.StringType, groupNames)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.Groups = groupList
	} else {
		// Ensure we always have a valid list
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Groups = emptyList
	}

	// Read back the final state
	readUser, err := r.client.GetUser(updatedUser.UserID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated user",
			fmt.Sprintf("Could not read updated user: %v", err),
		)
		return
	}

	// Update state with the latest information
	resourceUserToState(ctx, readUser, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implementation
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting user", map[string]interface{}{
		"user_id": state.ID.ValueString(),
	})

	// Remove from groups first
	if !state.Groups.IsNull() && !state.Groups.IsUnknown() {
		err := r.retryOperation(ctx, func() error {
			return r.client.UpdateUserGroupsByName(state.ID.ValueString(), []string{})
		})
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Error removing user from groups",
				fmt.Sprintf("Could not remove user from groups: %v", err),
			)
		}
	}

	// Delete the user
	err := r.retryOperation(ctx, func() error {
		return r.client.DeleteUser(state.ID.ValueString())
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			fmt.Sprintf("Could not delete user: %v", err),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted user", map[string]interface{}{
		"user_id": state.ID.ValueString(),
	})
}

// ImportState implementation
func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	var userID string
	switch len(idParts) {
	case 1:
		// Import by user ID
		userID = idParts[0]
	case 2:
		if idParts[0] != "username" {
			resp.Diagnostics.AddError(
				"Invalid Import Format",
				`To import by username use format: username:example@domain.com`,
			)
			return
		}
		// Find user by username
		users, err := r.client.GetUsers()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing user",
				fmt.Sprintf("Could not list users: %s", err),
			)
			return
		}

		for _, user := range users {
			if user.Username == idParts[1] {
				userID = user.UserID
				break
			}
		}

		if userID == "" {
			resp.Diagnostics.AddError(
				"User Not Found",
				fmt.Sprintf("Could not find user with username: %s", idParts[1]),
			)
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Invalid Import Format",
			`Use format: "user_id" or "username:example@domain.com"`,
		)
		return
	}

	// Set the imported ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), userID)...)
}
