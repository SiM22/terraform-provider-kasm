package zones

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-kasm/internal/client"
)

var (
	_ datasource.DataSource = &zonesDataSource{}
)

// zonesDataSource is the data source implementation.
type zonesDataSource struct {
	client *client.Client
}

// zoneModel maps zone data
type zoneModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	AutoScalingEnabled types.Bool   `tfsdk:"auto_scaling_enabled"`
	AWSEnabled         types.Bool   `tfsdk:"aws_enabled"`
	AWSRegion          types.String `tfsdk:"aws_region"`
	AWSAccessKeyID     types.String `tfsdk:"aws_access_key_id"`
	AWSSecretAccessKey types.String `tfsdk:"aws_secret_access_key"`
	EC2AgentAMIID      types.String `tfsdk:"ec2_agent_ami_id"`
}

// zonesDataSourceModel maps the data source schema data
type zonesDataSourceModel struct {
	Brief types.Bool   `tfsdk:"brief"`
	Zones []zoneModel  `tfsdk:"zones"`
	ID    types.String `tfsdk:"id"`
}

// New creates a new zones data source
func New() datasource.DataSource {
	return &zonesDataSource{}
}

// Metadata returns the data source type name
func (d *zonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zones"
}

// Schema defines the schema for the data source
func (d *zonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Kasm deployment zones.",
		Attributes: map[string]schema.Attribute{
			"brief": schema.BoolAttribute{
				Optional:    true,
				Description: "Limit the information returned for each zone",
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"zones": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the zone",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the zone",
						},
						"auto_scaling_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether auto-scaling is enabled",
						},
						"aws_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether AWS integration is enabled",
						},
						"aws_region": schema.StringAttribute{
							Computed:    true,
							Description: "The AWS region",
						},
						"aws_access_key_id": schema.StringAttribute{
							Computed:    true,
							Description: "The AWS access key ID",
						},
						"aws_secret_access_key": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The AWS secret access key",
						},
						"ec2_agent_ami_id": schema.StringAttribute{
							Computed:    true,
							Description: "The EC2 agent AMI ID",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *zonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *zonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zonesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	brief := false
	if !state.Brief.IsNull() {
		brief = state.Brief.ValueBool()
	}

	zones, err := d.client.GetZones(brief)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Zones",
			fmt.Sprintf("Could not read zones: %v", err),
		)
		return
	}

	// Map response body to model
	for _, zone := range zones {
		zoneState := zoneModel{
			ID:                 types.StringValue(zone.ZoneID),
			Name:               types.StringValue(zone.ZoneName),
			AutoScalingEnabled: types.BoolValue(zone.AutoScalingEnabled),
			AWSEnabled:         types.BoolValue(zone.AWSEnabled),
			AWSRegion:          types.StringValue(zone.AWSRegion),
			AWSAccessKeyID:     types.StringValue(zone.AWSAccessKeyID),
			AWSSecretAccessKey: types.StringValue(zone.AWSSecretAccessKey),
			EC2AgentAMIID:      types.StringValue(zone.EC2AgentAMIID),
		}
		state.Zones = append(state.Zones, zoneState)
	}

	// Set ID based on the timestamp
	state.ID = types.StringValue("zones")

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
