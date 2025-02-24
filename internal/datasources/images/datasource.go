package images

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementations satisfy the expected interfaces.
var (
	_ datasource.DataSource = &imagesDataSource{}
	_ datasource.DataSource = &sessionRecordingDataSource{}
	_ datasource.DataSource = &sessionsRecordingsDataSource{}
)

// imagesDataSource implements the data source
type imagesDataSource struct {
	client *client.Client
}

// imagesDataSourceModel maps the data source schema data
type imagesDataSourceModel struct {
	Images []imageModel `tfsdk:"images"`
}

// imageModel for the images data source
type imageModel struct {
	ID                  types.String  `tfsdk:"id"`
	Name                types.String  `tfsdk:"name"`
	FriendlyName        types.String  `tfsdk:"friendly_name"`
	Description         types.String  `tfsdk:"description"`
	Memory              types.Int64   `tfsdk:"memory"`
	Cores               types.Float64 `tfsdk:"cores"`
	CPUAllocationMethod types.String  `tfsdk:"cpu_allocation_method"`
	// Add other fields as needed
}

// NewImagesDataSource creates a new images data source
func New() datasource.DataSource {
	return &imagesDataSource{}
}

// Metadata returns the data source type name
func (d *imagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

// Schema defines the schema for the data source
func (d *imagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of all available Kasm images.",
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Description: "List of images",
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
						"cpu_allocation_method": schema.StringAttribute{
							Description: "CPU allocation method",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source
func (d *imagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *imagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Starting Read method for images data source")

	var state imagesDataSourceModel

	images, err := d.client.GetImages()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Images",
			fmt.Sprintf("Could not read Kasm images: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Got %d images from API", len(images)))

	// Map response body to model
	state.Images = make([]imageModel, 0, len(images))
	for _, image := range images {
		tflog.Debug(ctx, fmt.Sprintf("Processing image: %s (%s)", image.Name, image.ImageID))

		imageState := imageModel{
			ID:                  types.StringValue(image.ImageID),
			Name:                types.StringValue(image.Name),
			FriendlyName:        types.StringValue(image.FriendlyName),
			Description:         types.StringValue(image.Description),
			Memory:              types.Int64Value(image.Memory),
			Cores:               types.Float64Value(image.Cores),
			CPUAllocationMethod: types.StringValue(image.CPUAllocationMethod),
		}
		state.Images = append(state.Images, imageState)
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapped %d images to state", len(state.Images)))

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error setting state", map[string]interface{}{
			"error": resp.Diagnostics.Errors(),
		})
		return
	}
}

type sessionRecordingDataSource struct {
	client *client.Client
}

type sessionRecordingModel struct {
	KasmID              types.String     `tfsdk:"kasm_id"`
	PreauthDownloadLink types.Bool       `tfsdk:"preauth_download_link"`
	Recordings          []recordingModel `tfsdk:"recordings"`
}

type recordingModel struct {
	RecordingID                 types.String `tfsdk:"recording_id"`
	AccountID                   types.String `tfsdk:"account_id"`
	SessionRecordingURL         types.String `tfsdk:"session_recording_url"`
	SessionRecordingMetadata    types.Map    `tfsdk:"session_recording_metadata"`
	SessionRecordingDownloadURL types.String `tfsdk:"session_recording_download_url"`
}

type sessionsRecordingsDataSource struct {
	client *client.Client
}

type sessionsRecordingsModel struct {
	KasmIDs             types.List                  `tfsdk:"kasm_ids"` // Changed from []types.String
	PreauthDownloadLink types.Bool                  `tfsdk:"preauth_download_link"`
	Sessions            map[string][]recordingModel `tfsdk:"sessions"`
}

// Implement Metadata, Schema, Configure, and Read methods for each data source
func (d *sessionRecordingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_recording"
}

func (d *sessionsRecordingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sessions_recordings"
}

// Schema implementations for the session recording data sources
func (d *sessionRecordingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches recordings for a specific Kasm session.",
		Attributes: map[string]schema.Attribute{
			"kasm_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Kasm session",
			},
			"preauth_download_link": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to include pre-authorized download links",
			},
			"recordings": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"recording_id": schema.StringAttribute{
							Computed: true,
						},
						"account_id": schema.StringAttribute{
							Computed: true,
						},
						"session_recording_url": schema.StringAttribute{
							Computed: true,
						},
						"session_recording_metadata": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"session_recording_download_url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *sessionsRecordingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches recordings for multiple Kasm sessions.",
		Attributes: map[string]schema.Attribute{
			"kasm_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of Kasm session IDs",
			},
			"preauth_download_link": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to include pre-authorized download links",
			},
			"sessions": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"recording_id": schema.StringAttribute{
							Computed: true,
						},
						"account_id": schema.StringAttribute{
							Computed: true,
						},
						"session_recording_url": schema.StringAttribute{
							Computed: true,
						},
						"session_recording_metadata": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"session_recording_download_url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read method for sessionRecordingDataSource
func (d *sessionRecordingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionRecordingModel

	// Get the configuration
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get recordings from the API
	recordings, err := d.client.GetSessionRecordings(
		state.KasmID.ValueString(),
		state.PreauthDownloadLink.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Session Recordings",
			fmt.Sprintf("Could not read recordings for session %s: %s",
				state.KasmID.ValueString(), err),
		)
		return
	}

	// Map the recordings to our model
	state.Recordings = make([]recordingModel, 0, len(recordings))
	for _, rec := range recordings {
		metadata, diags := types.MapValueFrom(ctx, types.StringType, rec.SessionRecordingMetadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		recordingState := recordingModel{
			RecordingID:                 types.StringValue(rec.RecordingID),
			AccountID:                   types.StringValue(rec.AccountID),
			SessionRecordingURL:         types.StringValue(rec.SessionRecordingURL),
			SessionRecordingMetadata:    metadata,
			SessionRecordingDownloadURL: types.StringValue(rec.SessionRecordingDownloadURL),
		}
		state.Recordings = append(state.Recordings, recordingState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read method for sessionsRecordingsDataSource
func (d *sessionsRecordingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionsRecordingsModel

	// Get configuration
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert types.List to []string for the API call
	var kasmIDs []string
	diags = state.KasmIDs.ElementsAs(ctx, &kasmIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get recordings from the API
	recordings, err := d.client.GetSessionsRecordings(
		kasmIDs,
		state.PreauthDownloadLink.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Sessions Recordings",
			fmt.Sprintf("Could not read recordings for sessions: %s", err),
		)
		return
	}

	// Initialize the sessions map
	state.Sessions = make(map[string][]recordingModel)

	// Map the recordings to our model
	for kasmID, sessionRecordings := range recordings {
		recordingModels := make([]recordingModel, 0, len(sessionRecordings))
		for _, rec := range sessionRecordings {
			metadata, diags := types.MapValueFrom(ctx, types.StringType, rec.SessionRecordingMetadata)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			recordingState := recordingModel{
				RecordingID:                 types.StringValue(rec.RecordingID),
				AccountID:                   types.StringValue(rec.AccountID),
				SessionRecordingURL:         types.StringValue(rec.SessionRecordingURL),
				SessionRecordingMetadata:    metadata,
				SessionRecordingDownloadURL: types.StringValue(rec.SessionRecordingDownloadURL),
			}
			recordingModels = append(recordingModels, recordingState)
		}
		state.Sessions[kasmID] = recordingModels
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure methods for both data sources
func (d *sessionRecordingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sessionsRecordingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
