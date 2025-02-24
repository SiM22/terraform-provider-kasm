package registries

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-kasm/internal/client"
)

var _ datasource.DataSource = &RegistriesDataSource{}

type RegistriesDataSource struct {
	client *client.Client
}

type RegistriesDataSourceModel struct {
	Registries []RegistryModel `tfsdk:"registries"`
}

type RegistryModel struct {
	ID            types.String `tfsdk:"id"`
	URL           types.String `tfsdk:"url"`
	DoAutoUpdate  types.Bool   `tfsdk:"do_auto_update"`
	SchemaVersion types.String `tfsdk:"schema_version"`
	IsVerified    types.Bool   `tfsdk:"is_verified"`
	Channel       types.String `tfsdk:"channel"`
}

func New() datasource.DataSource {
	return &RegistriesDataSource{}
}

func (d *RegistriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RegistriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registries"
}

func (d *RegistriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all available Kasm registries.",
		Attributes: map[string]schema.Attribute{
			"registries": schema.ListNestedAttribute{
				Description: "List of registries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Registry identifier.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "Registry URL.",
							Computed:    true,
						},
						"do_auto_update": schema.BoolAttribute{
							Description: "Whether auto-update is enabled.",
							Computed:    true,
						},
						"schema_version": schema.StringAttribute{
							Description: "Schema version.",
							Computed:    true,
						},
						"is_verified": schema.BoolAttribute{
							Description: "Whether the registry is verified.",
							Computed:    true,
						},
						"channel": schema.StringAttribute{
							Description: "Registry channel.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *RegistriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegistriesDataSourceModel

	registries, err := d.client.GetRegistries()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading registries",
			fmt.Sprintf("Unable to read registries: %v", err),
		)
		return
	}

	for _, registry := range registries {
		state.Registries = append(state.Registries, RegistryModel{
			ID:            types.StringValue(registry.RegistryID),
			URL:           types.StringValue(registry.RegistryURL),
			DoAutoUpdate:  types.BoolValue(registry.DoAutoUpdate),
			SchemaVersion: types.StringValue(registry.SchemaVersion),
			IsVerified:    types.BoolValue(registry.IsVerified),
			Channel:       types.StringValue(registry.Channel),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
