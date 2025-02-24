package registry_images

import (
	"context"
	"fmt"
	"regexp"

	"terraform-provider-kasm/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &registryImageDataSource{}

type registryImageDataSource struct {
	client *client.Client
}

type registryImageDataSourceModel struct {
	ID         types.String                   `tfsdk:"id"`
	RegistryID types.String                   `tfsdk:"registry_id"`
	Images     []registryImageDataSourceImage `tfsdk:"images"`
}

type registryImageDataSourceImage struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	FriendlyName types.String  `tfsdk:"friendly_name"`
	Description  types.String  `tfsdk:"description"`
	Memory       types.Int64   `tfsdk:"memory"`
	Cores        types.Float64 `tfsdk:"cores"`
}

func New() datasource.DataSource {
	return &registryImageDataSource{}
}

func (d *registryImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_images"
}

func (d *registryImageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of Kasm registry images.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Data source identifier",
				Computed:    true,
			},
			"registry_id": schema.StringAttribute{
				Description: "Registry ID to filter images by",
				Optional:    true,
			},
			"images": schema.ListNestedAttribute{
				Description: "List of registry images",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Image ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Image name",
							Computed:    true,
						},
						"friendly_name": schema.StringAttribute{
							Description: "User-friendly name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Image description",
							Computed:    true,
						},
						"memory": schema.Int64Attribute{
							Description: "Memory in bytes",
							Computed:    true,
						},
						"cores": schema.Float64Attribute{
							Description: "CPU cores",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *registryImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *registryImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading registry images data source")

	var config registryImageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all images from the client
	images, err := d.client.ListRegistryImages()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Images",
			fmt.Sprintf("Could not read Kasm images: %s", err),
		)
		return
	}

	// Debug: Print all available images
	tflog.Debug(ctx, "Available images:")
	for _, img := range images {
		tflog.Debug(ctx, fmt.Sprintf("  - Image: Name='%s', Registry='%s', ID='%s'",
			img.Name, img.DockerRegistry, img.ImageID))
	}

	// Filter images by registry ID if provided
	var filteredImages []client.RegistryImage
	if config.RegistryID.IsNull() {
		filteredImages = images
	} else {
		registryID := config.RegistryID.ValueString()

		// Validate registry ID format (must be a UUID)
		if !isValidUUID(registryID) {
			resp.Diagnostics.AddError(
				"Invalid Registry ID",
				fmt.Sprintf("Registry ID '%s' is invalid", registryID),
			)
			return
		}

		// Get registries to validate the registry ID
		registries, err := d.client.GetRegistries()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Kasm Registries",
				fmt.Sprintf("Could not read Kasm registries: %s", err),
			)
			return
		}

		// Check if the registry ID exists
		registryExists := false
		for _, reg := range registries {
			if reg.RegistryID == registryID {
				registryExists = true
				break
			}
		}

		if !registryExists {
			// Return empty list for non-existent registry
			filteredImages = []client.RegistryImage{}
		} else {
			// Filter images by registry ID
			for _, img := range images {
				if img.DockerRegistry == registryID {
					filteredImages = append(filteredImages, img)
				}
			}
		}
	}

	// Convert filtered images to data source model
	var dataSourceImages []registryImageDataSourceImage
	for _, img := range filteredImages {
		dataSourceImage := registryImageDataSourceImage{
			ID:           types.StringValue(img.ImageID),
			Name:         types.StringValue(img.Name),
			FriendlyName: types.StringValue(img.FriendlyName),
			Description:  types.StringValue(img.Description),
			Memory:       types.Int64Value(img.Memory),
			Cores:        types.Float64Value(img.Cores),
		}
		dataSourceImages = append(dataSourceImages, dataSourceImage)
	}

	// Set data source model
	config.ID = types.StringValue("kasm_registry_images")
	config.Images = dataSourceImages
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(id string) bool {
	// Basic UUID format check (8-4-4-4-12)
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidPattern.MatchString(id)
}
