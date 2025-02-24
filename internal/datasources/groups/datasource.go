package groups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource = &groupsDataSource{}
)

// groupsDataSource is the data source implementation
type groupsDataSource struct {
	client *client.Client
}

// groupsDataSourceModel maps the data source schema data
type groupsDataSourceModel struct {
	Groups []groupModel `tfsdk:"groups"`
}

// groupModel maps group schema data
type groupModel struct {
	GroupID     types.String   `tfsdk:"group_id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Priority    types.Int64    `tfsdk:"priority"`
	Permissions []types.String `tfsdk:"permissions"`
}

// New creates a new groups data source
func New() datasource.DataSource {
	return &groupsDataSource{}
}

// Metadata returns the data source type name
func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

// Schema defines the schema for the data source
func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of all available Kasm groups.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "List of groups",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Description: "Group ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Group name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Group description",
							Computed:    true,
						},
						"priority": schema.Int64Attribute{
							Description: "Group priority",
							Computed:    true,
						},
						"permissions": schema.ListAttribute{
							Description: "List of group permissions",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data
func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read groups data source")
	var state groupsDataSourceModel

	groups, err := d.client.GetGroups()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Groups",
			fmt.Sprintf("Could not read Kasm groups: %s", err),
		)
		return
	}

	// Map response body to model
	for _, group := range groups {
		groupState := groupModel{
			GroupID:     types.StringValue(group.GroupID),
			Name:        types.StringValue(group.Name),
			Description: types.StringValue(group.Description),
			Priority:    types.Int64Value(int64(group.Priority)),
		}

		// Convert permissions to []types.String
		permissions := make([]types.String, 0, len(group.Permissions))
		for _, perm := range group.Permissions {
			permissions = append(permissions, types.StringValue(perm))
		}
		groupState.Permissions = permissions

		state.Groups = append(state.Groups, groupState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Finished reading groups data source")
}
