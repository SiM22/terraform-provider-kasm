package users

import (
	"context"
	"fmt"
	"terraform-provider-kasm/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &userDataSource{}

type userDataSource struct {
	client *client.Client
}

type userDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	UserID           types.String `tfsdk:"user_id"`
	Username         types.String `tfsdk:"username"`
	FirstName        types.String `tfsdk:"first_name"`
	LastName         types.String `tfsdk:"last_name"`
	Organization     types.String `tfsdk:"organization"`
	Phone            types.String `tfsdk:"phone"`
	Groups           types.List   `tfsdk:"groups"`
	AuthorizedImages types.List   `tfsdk:"authorized_images"`
}

func New() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch information about a Kasm user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "ID of the user.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the user.",
				Computed:    true,
			},
			"first_name": schema.StringAttribute{
				Description: "First name of the user.",
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "Last name of the user.",
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "Organization of the user.",
				Computed:    true,
			},
			"phone": schema.StringAttribute{
				Description: "Phone number of the user.",
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of groups the user belongs to.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"authorized_images": schema.ListAttribute{
				Description: "List of image IDs the user is authorized to use.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm User",
			fmt.Sprintf("Could not read user ID %s: %s", state.UserID.ValueString(), err),
		)
		return
	}

	// Get authorized images
	authorizedImages, err := d.client.GetUserAuthorizedImages(user.UserID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user authorized images",
			fmt.Sprintf("Could not read user authorized images: %v", err),
		)
		return
	}

	// Convert groups to list of names
	groupNames := []string{}
	for _, group := range user.Groups {
		if !group.IsSystem {
			groupNames = append(groupNames, group.Name)
		}
	}

	// Create lists
	groupList, diags := types.ListValueFrom(ctx, types.StringType, groupNames)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageList, diags := types.ListValueFrom(ctx, types.StringType, authorizedImages)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	state.ID = types.StringValue(user.UserID)
	state.Username = types.StringValue(user.Username)
	state.FirstName = types.StringValue(user.FirstName)
	state.LastName = types.StringValue(user.LastName)
	state.Organization = types.StringValue(user.Organization)
	state.Phone = types.StringValue(user.Phone)
	state.Groups = groupList
	state.AuthorizedImages = imageList

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
