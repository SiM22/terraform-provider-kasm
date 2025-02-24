package users_list

import (
	"context"
	"fmt"
	"terraform-provider-kasm/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &usersDataSource{}

type usersDataSource struct {
	client *client.Client
}

type userModel struct {
	ID               types.String `tfsdk:"id"`
	Username         types.String `tfsdk:"username"`
	FirstName        types.String `tfsdk:"first_name"`
	LastName         types.String `tfsdk:"last_name"`
	Organization     types.String `tfsdk:"organization"`
	Phone            types.String `tfsdk:"phone"`
	Groups           types.List   `tfsdk:"groups"`
	AuthorizedImages types.List   `tfsdk:"authorized_images"`
}

type usersDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Users []userModel  `tfsdk:"users"`
}

func New() datasource.DataSource {
	return &usersDataSource{}
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of all Kasm users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "List of users",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User ID",
							Computed:    true,
						},
						"username": schema.StringAttribute{
							Description: "Username",
							Computed:    true,
						},
						"first_name": schema.StringAttribute{
							Description: "First name",
							Computed:    true,
						},
						"last_name": schema.StringAttribute{
							Description: "Last name",
							Computed:    true,
						},
						"organization": schema.StringAttribute{
							Description: "Organization",
							Computed:    true,
						},
						"phone": schema.StringAttribute{
							Description: "Phone number",
							Computed:    true,
						},
						"groups": schema.ListAttribute{
							Description: "List of groups the user belongs to",
							Computed:    true,
							ElementType: types.StringType,
						},
						"authorized_images": schema.ListAttribute{
							Description: "List of authorized image IDs",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading users data source")

	var state usersDataSourceModel

	// Get list of users from API
	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d users", len(users)))

	// Map response body to model
	state.Users = make([]userModel, 0)
	for _, user := range users {
		var userState userModel
		userState.ID = types.StringValue(user.UserID)
		userState.Username = types.StringValue(user.Username)
		userState.FirstName = types.StringValue(user.FirstName)
		userState.LastName = types.StringValue(user.LastName)
		userState.Organization = types.StringValue(user.Organization)
		userState.Phone = types.StringValue(user.Phone)

		// Convert groups to types.List
		groups := make([]string, 0)
		for _, group := range user.Groups {
			groups = append(groups, group.Name)
		}
		groupsList, diags := types.ListValueFrom(ctx, types.StringType, groups)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		userState.Groups = groupsList

		// Convert authorized images to types.List
		authorizedImages := make([]string, 0)
		for _, image := range user.AuthorizedImages {
			authorizedImages = append(authorizedImages, image)
		}
		imagesList, diags := types.ListValueFrom(ctx, types.StringType, authorizedImages)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		userState.AuthorizedImages = imagesList

		state.Users = append(state.Users, userState)
	}

	// Set id
	state.ID = types.StringValue("users")

	tflog.Debug(ctx, "Setting state")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Finished reading users data source")
}
