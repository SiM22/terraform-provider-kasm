package images

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &imageResource{}
	_ resource.ResourceWithImportState = &imageResource{}
)

// imageResource is the resource implementation
type imageResource struct {
	client *client.Client
}

// ImageResourceModel maps the resource schema data
type ImageResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	FriendlyName        types.String   `tfsdk:"friendly_name"`
	Desc                types.String   `tfsdk:"description"`
	Categories          []types.String `tfsdk:"categories"`
	Memory              types.Int64    `tfsdk:"memory"`
	Cores               types.Float64  `tfsdk:"cores"`
	CPUAllocationMethod types.String   `tfsdk:"cpu_allocation_method"`
	ImageSrc            types.String   `tfsdk:"image_src"`
	DockerRegistry      types.String   `tfsdk:"docker_registry"`
	Available           types.Bool     `tfsdk:"available"`       // Add this field
	VolumeMappings      types.Map      `tfsdk:"volume_mappings"` // Added
	NetworkName         types.String   `tfsdk:"network_name"`    // Added
	DockerUser          types.String   `tfsdk:"docker_user"`
	DockerPassword      types.String   `tfsdk:"docker_password"`
	UncompressedSizeMB  types.Int64    `tfsdk:"uncompressed_size_mb"`
	ImageType           types.String   `tfsdk:"image_type"`
	Enabled             types.Bool     `tfsdk:"enabled"`
	GPUCount            types.Int64    `tfsdk:"gpu_count"`
	RequireGPU          types.Bool     `tfsdk:"require_gpu"`
	RestrictToNetwork   types.Bool     `tfsdk:"restrict_to_network"`
	RestrictToServer    types.Bool     `tfsdk:"restrict_to_server"`
	RestrictToZone      types.Bool     `tfsdk:"restrict_to_zone"`
	ServerID            types.String   `tfsdk:"server_id"`
	ZoneID              types.String   `tfsdk:"zone_id"`
	Hidden              types.Bool     `tfsdk:"hidden"`
	RunConfig           types.Map      `tfsdk:"run_config"`
	ExecConfig          types.Map      `tfsdk:"exec_config"`
	KasmAudioDefaultOn  types.Bool     `tfsdk:"kasm_audio_default_on"`
}

// NewImageResource creates a new resource instance
func New() resource.Resource {
	return &imageResource{}
}

// Metadata returns the resource type name
func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Schema defines the schema for the resource
func (r *imageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kasm workspace image.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Docker image name and tag",
			},
			"friendly_name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the image",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Image description",
			},
			"categories": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Image categories",
			},
			"memory": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory allocation in bytes",
			},
			"cores": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of CPU cores",
			},
			"cpu_allocation_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU allocation method (e.g., Inherit, Quotas)",
			},
			"image_src": schema.StringAttribute{
				Optional:    true,
				Description: "Path to image icon/thumbnail",
			},
			"docker_registry": schema.StringAttribute{
				Required:    true,
				Description: "Docker registry URL",
			},
			"docker_user": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Docker registry username",
			},
			"docker_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Docker registry password",
			},
			"uncompressed_size_mb": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Uncompressed image size in MB",
			},
			"image_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of image",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the image is enabled",
			},
			"gpu_count": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of GPUs required",
			},
			"require_gpu": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether GPU is required",
			},
			"restrict_to_network": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to restrict to specific network",
			},
			"restrict_to_server": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to restrict to specific server",
			},
			"restrict_to_zone": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to restrict to specific zone",
			},
			"server_id": schema.StringAttribute{
				Optional:    true,
				Description: "Server ID if restricted to server",
			},
			"zone_id": schema.StringAttribute{
				Optional:    true,
				Description: "Zone ID if restricted to zone",
			},
			"hidden": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the image is hidden",
			},
			"run_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Docker run configuration",
			},
			"exec_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Docker exec configuration",
			},
			"volume_mappings": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Volume mapping configuration",
			},
			"available": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the image is available",
			},
			"network_name": schema.StringAttribute{
				Optional:    true,
				Description: "Network name for the image",
			},
			"kasm_audio_default_on": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether Kasm audio is enabled by default",
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *imageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeMapping := make(map[string]interface{})
	if !plan.VolumeMappings.IsNull() {
		var tempMap map[string]string
		diags = plan.VolumeMappings.ElementsAs(ctx, &tempMap, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		for k, v := range tempMap {
			volumeMapping[k] = v
		}
	}

	runConfig := make(map[string]string)
	if !plan.RunConfig.IsNull() {
		diags = plan.RunConfig.ElementsAs(ctx, &runConfig, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	execConfig := make(map[string]string)
	if !plan.ExecConfig.IsNull() {
		diags = plan.ExecConfig.ElementsAs(ctx, &execConfig, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	image := &client.Image{
		Name:                plan.Name.ValueString(),
		FriendlyName:        plan.FriendlyName.ValueString(),
		Description:         plan.Desc.ValueString(),
		Memory:              plan.Memory.ValueInt64(),
		Cores:               plan.Cores.ValueFloat64(),
		CPUAllocationMethod: plan.CPUAllocationMethod.ValueString(),
		DockerRegistry:      plan.DockerRegistry.ValueString(),
		DockerUser:          plan.DockerUser.ValueString(),
		DockerPassword:      plan.DockerPassword.ValueString(),
		ImageSrc:            plan.ImageSrc.ValueString(),
		Enabled:             plan.Enabled.ValueBool(),
		Available:           plan.Available.ValueBool(),
		RestrictToZone:      plan.RestrictToZone.ValueBool(),
		RestrictToServer:    plan.RestrictToServer.ValueBool(),
		RestrictToNetwork:   plan.RestrictToNetwork.ValueBool(),
		ServerID:            plan.ServerID.ValueString(),
		ZoneID:              plan.ZoneID.ValueString(),
		NetworkName:         plan.NetworkName.ValueString(),
		VolumeMappings:      volumeMapping,
		RunConfig:           stringMapToInterface(runConfig),
		ExecConfig:          stringMapToInterface(execConfig),
	}

	createdImage, err := r.client.CreateImage(image)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating image",
			fmt.Sprintf("Could not create image: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(createdImage.ImageID)
	plan.CPUAllocationMethod = types.StringValue("Inherit")
	plan.ImageType = types.StringValue("Container")
	plan.UncompressedSizeMB = types.Int64Value(0)
	plan.KasmAudioDefaultOn = types.BoolValue(false)

	if createdImage.CPUAllocationMethod != "" {
		plan.CPUAllocationMethod = types.StringValue(createdImage.CPUAllocationMethod)
	}
	if createdImage.ImageType != "" {
		plan.ImageType = types.StringValue(createdImage.ImageType)
	}
	if createdImage.UncompressedSizeMB > 0 {
		plan.UncompressedSizeMB = types.Int64Value(int64(createdImage.UncompressedSizeMB))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	image, err := r.client.GetImage(state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading image",
			fmt.Sprintf("Could not read image ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	state.Name = types.StringValue(image.Name)
	state.FriendlyName = types.StringValue(image.FriendlyName)
	state.Desc = types.StringValue(image.Description)
	state.Memory = types.Int64Value(image.Memory)
	state.Cores = types.Float64Value(image.Cores)
	state.CPUAllocationMethod = types.StringValue(image.CPUAllocationMethod)
	state.ImageSrc = types.StringValue(image.ImageSrc)
	state.DockerRegistry = types.StringValue(image.DockerRegistry)
	state.Available = types.BoolValue(image.Available)
	state.NetworkName = types.StringValue(image.NetworkName)
	state.DockerUser = types.StringValue(image.DockerUser)
	state.DockerPassword = types.StringValue(image.DockerPassword)
	state.Enabled = types.BoolValue(image.Enabled)
	state.RestrictToNetwork = types.BoolValue(image.RestrictToNetwork)
	state.RestrictToServer = types.BoolValue(image.RestrictToServer)
	state.RestrictToZone = types.BoolValue(image.RestrictToZone)
	state.ServerID = types.StringValue(image.ServerID)
	state.ZoneID = types.StringValue(image.ZoneID)

	volumeMapping, diags := types.MapValueFrom(ctx, types.StringType, image.VolumeMappings)
	resp.Diagnostics.Append(diags...)
	state.VolumeMappings = volumeMapping

	runConfig, diags := types.MapValueFrom(ctx, types.StringType, image.RunConfig)
	resp.Diagnostics.Append(diags...)
	state.RunConfig = runConfig

	execConfig, diags := types.MapValueFrom(ctx, types.StringType, image.ExecConfig)
	resp.Diagnostics.Append(diags...)
	state.ExecConfig = execConfig

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeMapping := make(map[string]interface{})
	if !plan.VolumeMappings.IsNull() {
		var tempMap map[string]string
		diags = plan.VolumeMappings.ElementsAs(ctx, &tempMap, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		for k, v := range tempMap {
			volumeMapping[k] = v
		}
	}

	runConfig := make(map[string]string)
	if !plan.RunConfig.IsNull() {
		diags = plan.RunConfig.ElementsAs(ctx, &runConfig, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	execConfig := make(map[string]string)
	if !plan.ExecConfig.IsNull() {
		diags = plan.ExecConfig.ElementsAs(ctx, &execConfig, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	image := &client.Image{
		ImageID:             plan.ID.ValueString(),
		Name:                plan.Name.ValueString(),
		FriendlyName:        plan.FriendlyName.ValueString(),
		Description:         plan.Desc.ValueString(),
		Memory:              plan.Memory.ValueInt64(),
		Cores:               plan.Cores.ValueFloat64(),
		CPUAllocationMethod: plan.CPUAllocationMethod.ValueString(),
		DockerRegistry:      plan.DockerRegistry.ValueString(),
		DockerUser:          plan.DockerUser.ValueString(),
		DockerPassword:      plan.DockerPassword.ValueString(),
		ImageSrc:            plan.ImageSrc.ValueString(),
		Enabled:             plan.Enabled.ValueBool(),
		Available:           plan.Available.ValueBool(),
		RestrictToZone:      plan.RestrictToZone.ValueBool(),
		RestrictToServer:    plan.RestrictToServer.ValueBool(),
		RestrictToNetwork:   plan.RestrictToNetwork.ValueBool(),
		ServerID:            plan.ServerID.ValueString(),
		ZoneID:              plan.ZoneID.ValueString(),
		NetworkName:         plan.NetworkName.ValueString(),
		VolumeMappings:      volumeMapping,
		RunConfig:           stringMapToInterface(runConfig),
		ExecConfig:          stringMapToInterface(execConfig),
	}

	updatedImage, err := r.client.UpdateImage(image)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating image",
			fmt.Sprintf("Could not update image ID %s: %v", plan.ID.ValueString(), err),
		)
		return
	}

	plan.ID = types.StringValue(updatedImage.ImageID)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func stringMapToInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteImage(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Image Deletion Not Supported",
			"The image deletion API is not officially documented. Please use Kasm Web UI for image deletion.",
		)
	}
}

func (r *imageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import image using ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

// Docker image validator
func validateDockerImage() validator.String {
	return validators.StringValidator{
		Desc: "must be a valid docker image name",
		ValidateFn: func(val string) bool {
			matched, _ := regexp.MatchString(`^(?:([a-zA-Z0-9-_.]+(?::[0-9]+)?)/)?([a-zA-Z0-9-_.]+(?:/[a-zA-Z0-9-_.]+)*):?([a-zA-Z0-9-_.]+)?$`, val)
			return matched
		},
		ErrMessage: "invalid docker image name format",
	}
}

func validateCPUAllocationMethod() validator.String {
	return validators.StringValidator{
		Desc: "must be one of: Inherit, Quotas",
		ValidateFn: func(val string) bool {
			validMethods := map[string]bool{
				"Inherit": true,
				"Quotas":  true,
			}
			return validMethods[val]
		},
		ErrMessage: "cpu allocation method must be either 'Inherit' or 'Quotas'",
	}
}

func validateMemory() validator.Int64 {
	return validators.Int64Validator{
		Desc: "memory must be a positive value in bytes",
		ValidateFn: func(val int64) bool {
			return val > 0
		},
		ErrMessage: "memory must be greater than 0 bytes",
	}
}

func validateCores() validator.Float64 {
	return validators.Float64Validator{
		Desc: "cores must be a positive value",
		ValidateFn: func(val float64) bool {
			return val > 0 && val <= 128
		},
		ErrMessage: "cores must be greater than 0 and less than or equal to 128",
	}
}

// Helper function for checking sensitive changes
func sensitiveValuesChanged(ctx context.Context, plan, state ImageResourceModel) bool {
	if !plan.DockerUser.Equal(state.DockerUser) ||
		!plan.DockerPassword.Equal(state.DockerPassword) {
		return true
	}
	return false
}

// Helper function to check if any restriction flags are enabled but their corresponding IDs are empty
func validateRestrictions(image *ImageResourceModel) error {
	if image.RestrictToZone.ValueBool() && image.ZoneID.ValueString() == "" {
		return fmt.Errorf("zone_id is required when restrict_to_zone is true")
	}
	if image.RestrictToServer.ValueBool() && image.ServerID.ValueString() == "" {
		return fmt.Errorf("server_id is required when restrict_to_server is true")
	}
	if image.RestrictToNetwork.ValueBool() && image.NetworkName.ValueString() == "" {
		return fmt.Errorf("network_name is required when restrict_to_network is true")
	}
	return nil
}
