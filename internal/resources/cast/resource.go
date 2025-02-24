package cast

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &castConfigResource{}
	_ resource.ResourceWithImportState = &castConfigResource{}
)

// castConfigResource is the resource implementation
type castConfigResource struct {
	client *client.Client
}

// CastConfigResourceModel maps the resource schema data
type CastConfigResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	ImageID               types.String `tfsdk:"image_id"`
	AllowedReferrers      types.List   `tfsdk:"allowed_referrers"`
	LimitSessions         types.Bool   `tfsdk:"limit_sessions"`
	SessionRemaining      types.Int64  `tfsdk:"session_remaining"`
	LimitIPs              types.Bool   `tfsdk:"limit_ips"`
	IPRequestLimit        types.Int64  `tfsdk:"ip_request_limit"`
	IPRequestSeconds      types.Int64  `tfsdk:"ip_request_seconds"`
	ErrorURL              types.String `tfsdk:"error_url"`
	EnableSharing         types.Bool   `tfsdk:"enable_sharing"`
	DisableControlPanel   types.Bool   `tfsdk:"disable_control_panel"`
	DisableTips           types.Bool   `tfsdk:"disable_tips"`
	DisableFixedRes       types.Bool   `tfsdk:"disable_fixed_res"`
	Key                   types.String `tfsdk:"key"`
	AllowAnonymous        types.Bool   `tfsdk:"allow_anonymous"`
	GroupID               types.String `tfsdk:"group_id"`
	RequireRecaptcha      types.Bool   `tfsdk:"require_recaptcha"`
	KasmURL               types.String `tfsdk:"kasm_url"`
	DynamicKasmURL        types.Bool   `tfsdk:"dynamic_kasm_url"`
	DynamicDockerNetwork  types.Bool   `tfsdk:"dynamic_docker_network"`
	AllowResume           types.Bool   `tfsdk:"allow_resume"`
	EnforceClientSettings types.Bool   `tfsdk:"enforce_client_settings"`
	AllowKasmAudio        types.Bool   `tfsdk:"allow_kasm_audio"`
	AllowKasmUploads      types.Bool   `tfsdk:"allow_kasm_uploads"`
	AllowKasmDownloads    types.Bool   `tfsdk:"allow_kasm_downloads"`
	AllowClipboardDown    types.Bool   `tfsdk:"allow_clipboard_down"`
	AllowClipboardUp      types.Bool   `tfsdk:"allow_clipboard_up"`
	AllowMicrophone       types.Bool   `tfsdk:"allow_microphone"`
	ValidUntil            types.String `tfsdk:"valid_until"`
	AllowSharing          types.Bool   `tfsdk:"allow_sharing"`
	AudioDefaultOn        types.Bool   `tfsdk:"audio_default_on"`
	IMEModeDefaultOn      types.Bool   `tfsdk:"ime_mode_default_on"`
}

// Validator functions
func validateKey() validator.String {
	return validators.StringValidator{
		Desc: "key must only contain alphanumeric characters, underscores, or hyphens",
		ValidateFn: func(v string) bool {
			matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, v)
			return matched
		},
		ErrMessage: "key must only contain alphanumeric characters, underscores, or hyphens",
	}
}

func validateImageID() validator.String {
	return validators.StringValidator{
		Desc: "must be a 32-character hexadecimal string",
		ValidateFn: func(v string) bool {
			matched, _ := regexp.MatchString(`^[a-fA-F0-9]{32}$`, v)
			return matched
		},
		ErrMessage: "image_id must be a 32-character hexadecimal string",
	}
}

// Cast-specific helper function for checking if a time string is in the correct format
func validateTimeFormat(timeStr string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", timeStr)
	return err == nil
}

// Cast-specific helper to validate session limits
func validateSessionLimits(limitSessions bool, sessionRemaining int64) bool {
	if limitSessions && sessionRemaining <= 0 {
		return false
	}
	return true
}

// Cast-specific helper to validate IP rate limiting configuration
func validateIPRateLimits(limitIPs bool, requestLimit, requestSeconds int64) bool {
	if limitIPs {
		if requestLimit <= 0 || requestSeconds <= 0 {
			return false
		}
	}
	return true
}

// NewCastConfigResource creates a new resource instance
func New() resource.Resource {
	return &castConfigResource{}
}

func (r *castConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Metadata returns the resource type name
func (r *castConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cast_config"
}

// Schema defines the schema for the resource
func (r *castConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kasm Session Casting configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the casting configuration.",
				Validators: []validator.String{
					&validators.StringValidator{
						Desc: "name must not be empty",
						ValidateFn: func(v string) bool {
							return len(v) > 0
						},
						ErrMessage: "name cannot be empty",
					},
				},
			},
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the image to use for the casting configuration.",
				Validators: []validator.String{
					validateImageID(),
				},
			},
			"allowed_referrers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of allowed referrer domains.",
			},
			"limit_sessions": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to limit the total number of sessions.",
			},
			"session_remaining": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of sessions allowed to be spawned.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"limit_ips": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to enable IP-based rate limiting.",
			},
			"ip_request_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of sessions allowed per IP within the time window.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"ip_request_seconds": schema.Int64Attribute{
				Optional:    true,
				Computed:    true, // This allows the provider to set a default value
				Description: "Time window in seconds for IP rate limiting. Defaults to 0.",
				Default:     int64default.StaticInt64(0), // Set default value to 0
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647), // Validates range from 0 to max int32
				},
			},
			"error_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL to redirect to on errors.",
				Validators: []validator.String{
					validators.ValidateURL(),
				},
			},
			"enable_sharing": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to enable session sharing.",
			},
			"disable_control_panel": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to disable the control panel.",
			},
			"disable_tips": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to disable the tips dialogue.",
			},
			"disable_fixed_res": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to disable fixed resolution in sharing mode.",
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the casting URL",
				Validators: []validator.String{
					validateKey(),
				},
			},
			"allow_anonymous": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow anonymous access.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "The group ID for anonymous users.",
			},
			"require_recaptcha": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to require reCAPTCHA verification.",
			},
			"kasm_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL to load in browser-based sessions.",
				Validators: []validator.String{
					validators.ValidateURL(),
				},
			},
			"dynamic_kasm_url": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow dynamic KASM_URL via query parameter.",
			},
			"dynamic_docker_network": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow dynamic docker network selection.",
			},
			"allow_resume": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow session resuming.",
			},
			"enforce_client_settings": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to enforce client settings.",
			},
			"allow_kasm_audio": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow audio streaming.",
			},
			"allow_kasm_uploads": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow file uploads.",
			},
			"allow_kasm_downloads": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow file downloads.",
			},
			"allow_clipboard_down": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow clipboard download.",
			},
			"allow_clipboard_up": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow clipboard upload.",
			},
			"allow_microphone": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow microphone access.",
			},
			"valid_until": schema.StringAttribute{
				Optional:    true,
				Description: "Expiration timestamp in UTC.",
			},
			"allow_sharing": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to allow session sharing.",
			},
			"audio_default_on": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether audio should be enabled by default.",
			},
			"ime_mode_default_on": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether IME mode should be enabled by default.",
			},
		},
	}
}

func (r *castConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CastConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start with required fields
	config := &client.CastConfig{
		CastingConfigName: plan.Name.ValueString(),
		ImageID:           plan.ImageID.ValueString(),
		Key:               plan.Key.ValueString(),
	}

	// Handle optional boolean fields
	if !plan.LimitSessions.IsNull() {
		config.LimitSessions = plan.LimitSessions.ValueBool()
	}
	if !plan.LimitIPs.IsNull() {
		config.LimitIPs = plan.LimitIPs.ValueBool()
	}
	if !plan.EnableSharing.IsNull() {
		config.EnableSharing = plan.EnableSharing.ValueBool()
	}
	if !plan.DisableControlPanel.IsNull() {
		config.DisableControlPanel = plan.DisableControlPanel.ValueBool()
	}
	if !plan.DisableTips.IsNull() {
		config.DisableTips = plan.DisableTips.ValueBool()
	}
	if !plan.DisableFixedRes.IsNull() {
		config.DisableFixedRes = plan.DisableFixedRes.ValueBool()
	}
	if !plan.AllowAnonymous.IsNull() {
		config.AllowAnonymous = plan.AllowAnonymous.ValueBool()
	}
	if !plan.RequireRecaptcha.IsNull() {
		config.RequireRecaptcha = plan.RequireRecaptcha.ValueBool()
	}
	if !plan.DynamicKasmURL.IsNull() {
		config.DynamicKasmURL = plan.DynamicKasmURL.ValueBool()
	}
	if !plan.DynamicDockerNetwork.IsNull() {
		config.DynamicDockerNetwork = plan.DynamicDockerNetwork.ValueBool()
	}
	if !plan.AllowResume.IsNull() {
		config.AllowResume = plan.AllowResume.ValueBool()
	}
	if !plan.EnforceClientSettings.IsNull() {
		config.EnforceClientSettings = plan.EnforceClientSettings.ValueBool()
	}
	if !plan.AllowKasmAudio.IsNull() {
		config.AllowKasmAudio = plan.AllowKasmAudio.ValueBool()
	}
	if !plan.AllowKasmUploads.IsNull() {
		config.AllowKasmUploads = plan.AllowKasmUploads.ValueBool()
	}
	if !plan.AllowKasmDownloads.IsNull() {
		config.AllowKasmDownloads = plan.AllowKasmDownloads.ValueBool()
	}
	if !plan.AllowClipboardDown.IsNull() {
		config.AllowClipboardDown = plan.AllowClipboardDown.ValueBool()
	}
	if !plan.AllowClipboardUp.IsNull() {
		config.AllowClipboardUp = plan.AllowClipboardUp.ValueBool()
	}
	if !plan.AllowMicrophone.IsNull() {
		config.AllowMicrophone = plan.AllowMicrophone.ValueBool()
	}
	if !plan.AllowSharing.IsNull() {
		config.AllowSharing = plan.AllowSharing.ValueBool()
	}
	if !plan.AudioDefaultOn.IsNull() {
		config.AudioDefaultOn = plan.AudioDefaultOn.ValueBool()
	}
	if !plan.IMEModeDefaultOn.IsNull() {
		config.IMEModeDefaultOn = plan.IMEModeDefaultOn.ValueBool()
	}

	// Handle optional numeric fields
	if !plan.SessionRemaining.IsNull() {
		config.SessionRemaining = int(plan.SessionRemaining.ValueInt64())
	}
	if !plan.IPRequestLimit.IsNull() {
		config.IPRequestLimit = int(plan.IPRequestLimit.ValueInt64())
	}
	if !plan.IPRequestSeconds.IsNull() {
		config.IPRequestSeconds = int(plan.IPRequestSeconds.ValueInt64())
	}

	// Handle optional string fields
	if !plan.ErrorURL.IsNull() {
		config.ErrorURL = plan.ErrorURL.ValueString()
	}
	if !plan.KasmURL.IsNull() {
		config.KasmURL = plan.KasmURL.ValueString()
	}
	if !plan.ValidUntil.IsNull() {
		config.ValidUntil = plan.ValidUntil.ValueString()
	}
	if !plan.GroupID.IsNull() {
		config.GroupID = plan.GroupID.ValueString()
	}

	// Handle AllowedReferrers list
	if !plan.AllowedReferrers.IsNull() {
		var referrers []string
		diags = plan.AllowedReferrers.ElementsAs(ctx, &referrers, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		config.AllowedReferrers = referrers
	}

	// Create the resource
	createdConfig, err := r.client.CreateCastConfig(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cast config",
			fmt.Sprintf("Could not create cast config, unexpected error: %v", err),
		)
		return
	}

	// Convert response to model
	state := r.toResourceModel(createdConfig)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *castConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CastConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get cast config from API
	castConfig, err := r.client.GetCastConfig(state.ID.ValueString())
	if err != nil {
		var notFoundErr *client.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading cast config",
			fmt.Sprintf("Could not read cast config ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	// Map response to state
	updatedState := r.toResourceModel(castConfig)

	// Set refreshed state
	diags = resp.State.Set(ctx, updatedState)
	resp.Diagnostics.Append(diags...)
}

func (r *castConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan and state
	var plan, state CastConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start with current state
	config := &client.CastConfig{
		ID:                state.ID.ValueString(),
		CastingConfigName: plan.Name.ValueString(),
		ImageID:           plan.ImageID.ValueString(),
		Key:               plan.Key.ValueString(),

		// Always set boolean values with defaults
		LimitSessions:         plan.LimitSessions.ValueBool(),
		LimitIPs:              plan.LimitIPs.ValueBool(),
		EnableSharing:         plan.EnableSharing.ValueBool(),
		DisableControlPanel:   plan.DisableControlPanel.ValueBool(),
		DisableTips:           plan.DisableTips.ValueBool(),
		DisableFixedRes:       plan.DisableFixedRes.ValueBool(),
		AllowAnonymous:        plan.AllowAnonymous.ValueBool(),
		RequireRecaptcha:      plan.RequireRecaptcha.ValueBool(),
		DynamicKasmURL:        plan.DynamicKasmURL.ValueBool(),
		DynamicDockerNetwork:  plan.DynamicDockerNetwork.ValueBool(),
		AllowResume:           plan.AllowResume.ValueBool(),
		EnforceClientSettings: plan.EnforceClientSettings.ValueBool(),
		AllowKasmAudio:        plan.AllowKasmAudio.ValueBool(),
		AllowKasmUploads:      plan.AllowKasmUploads.ValueBool(),
		AllowKasmDownloads:    plan.AllowKasmDownloads.ValueBool(),
		AllowClipboardDown:    plan.AllowClipboardDown.ValueBool(),
		AllowClipboardUp:      plan.AllowClipboardUp.ValueBool(),
		AllowMicrophone:       plan.AllowMicrophone.ValueBool(),
		AllowSharing:          plan.AllowSharing.ValueBool(),
		AudioDefaultOn:        plan.AudioDefaultOn.ValueBool(),
		IMEModeDefaultOn:      plan.IMEModeDefaultOn.ValueBool(),
	}

	// Handle numeric fields
	if !plan.SessionRemaining.IsNull() {
		config.SessionRemaining = int(plan.SessionRemaining.ValueInt64())
	}
	if !plan.IPRequestLimit.IsNull() {
		config.IPRequestLimit = int(plan.IPRequestLimit.ValueInt64())
	}
	if !plan.IPRequestSeconds.IsNull() {
		config.IPRequestSeconds = int(plan.IPRequestSeconds.ValueInt64())
	}

	// Handle string fields
	if !plan.ErrorURL.IsNull() {
		config.ErrorURL = plan.ErrorURL.ValueString()
	}
	if !plan.KasmURL.IsNull() {
		config.KasmURL = plan.KasmURL.ValueString()
	}
	if !plan.ValidUntil.IsNull() {
		config.ValidUntil = plan.ValidUntil.ValueString()
	}
	if !plan.GroupID.IsNull() {
		config.GroupID = plan.GroupID.ValueString()
	}

	// Handle AllowedReferrers list
	if !plan.AllowedReferrers.IsNull() {
		var referrers []string
		diags = plan.AllowedReferrers.ElementsAs(ctx, &referrers, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		config.AllowedReferrers = referrers
	}

	// Update the resource
	updatedConfig, err := r.client.UpdateCastConfig(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cast config",
			fmt.Sprintf("Could not update cast config, unexpected error: %v", err),
		)
		return
	}

	// Convert response to model
	newState := r.toResourceModel(updatedConfig)

	// Set state
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (r *castConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CastConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCastConfig(state.ID.ValueString())
	if err != nil {
		var notFoundErr *client.NotFoundError
		var unauthorizedErr *client.UnauthorizedError

		if errors.As(err, &notFoundErr) {
			resp.Diagnostics.AddWarning(
				"Cast config not found",
				fmt.Sprintf("Cast config with ID %s was not found, it may have been deleted externally", state.ID.ValueString()),
			)
			return
		}

		if errors.As(err, &unauthorizedErr) {
			resp.Diagnostics.AddError(
				"Unauthorized to delete cast config",
				"Please check your credentials and permissions",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting cast config",
			fmt.Sprintf("Could not delete cast config, unexpected error: %v", err),
		)
	}
}

func (r *castConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import cast config using ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper functions for model conversion
func (r *castConfigResource) toAPIModel(plan *CastConfigResourceModel) *client.CastConfig {
	config := &client.CastConfig{
		CastingConfigName:     plan.Name.ValueString(),
		ImageID:               plan.ImageID.ValueString(),
		Key:                   plan.Key.ValueString(),
		LimitSessions:         plan.LimitSessions.ValueBool(),
		SessionRemaining:      int(plan.SessionRemaining.ValueInt64()),
		LimitIPs:              plan.LimitIPs.ValueBool(),
		IPRequestLimit:        int(plan.IPRequestLimit.ValueInt64()),
		IPRequestSeconds:      int(plan.IPRequestSeconds.ValueInt64()),
		ErrorURL:              plan.ErrorURL.ValueString(),
		EnableSharing:         plan.EnableSharing.ValueBool(),
		DisableControlPanel:   plan.DisableControlPanel.ValueBool(),
		DisableTips:           plan.DisableTips.ValueBool(),
		DisableFixedRes:       plan.DisableFixedRes.ValueBool(),
		AllowAnonymous:        plan.AllowAnonymous.ValueBool(),
		GroupID:               plan.GroupID.ValueString(),
		RequireRecaptcha:      plan.RequireRecaptcha.ValueBool(),
		KasmURL:               plan.KasmURL.ValueString(),
		DynamicKasmURL:        plan.DynamicKasmURL.ValueBool(),
		DynamicDockerNetwork:  plan.DynamicDockerNetwork.ValueBool(),
		AllowResume:           plan.AllowResume.ValueBool(),
		EnforceClientSettings: plan.EnforceClientSettings.ValueBool(),
		AllowKasmAudio:        plan.AllowKasmAudio.ValueBool(),
		AllowKasmUploads:      plan.AllowKasmUploads.ValueBool(),
		AllowKasmDownloads:    plan.AllowKasmDownloads.ValueBool(),
		AllowClipboardDown:    plan.AllowClipboardDown.ValueBool(),
		AllowClipboardUp:      plan.AllowClipboardUp.ValueBool(),
		AllowMicrophone:       plan.AllowMicrophone.ValueBool(),
		ValidUntil:            plan.ValidUntil.ValueString(),
		AllowSharing:          plan.AllowSharing.ValueBool(),
		AudioDefaultOn:        plan.AudioDefaultOn.ValueBool(),
		IMEModeDefaultOn:      plan.IMEModeDefaultOn.ValueBool(),
	}

	// Handle optional arrays properly
	if !plan.AllowedReferrers.IsNull() {
		var referrers []string
		plan.AllowedReferrers.ElementsAs(context.Background(), &referrers, false)
		config.AllowedReferrers = referrers
	} else {
		config.AllowedReferrers = []string{} // Always set an empty array if null
	}

	return config
}

func (r *castConfigResource) toResourceModel(api *client.CastConfig) *CastConfigResourceModel {
	model := &CastConfigResourceModel{
		// Required fields always get set
		ID:      types.StringValue(api.ID),
		Name:    types.StringValue(api.CastingConfigName),
		ImageID: types.StringValue(api.ImageID),
		Key:     types.StringValue(api.Key),

		// Always set all boolean values with defaults
		LimitSessions:         types.BoolValue(api.LimitSessions),
		LimitIPs:              types.BoolValue(api.LimitIPs),
		EnableSharing:         types.BoolValue(api.EnableSharing),
		DisableControlPanel:   types.BoolValue(api.DisableControlPanel),
		DisableTips:           types.BoolValue(api.DisableTips),
		DisableFixedRes:       types.BoolValue(api.DisableFixedRes),
		AllowAnonymous:        types.BoolValue(api.AllowAnonymous),
		RequireRecaptcha:      types.BoolValue(api.RequireRecaptcha),
		DynamicKasmURL:        types.BoolValue(api.DynamicKasmURL),
		DynamicDockerNetwork:  types.BoolValue(api.DynamicDockerNetwork),
		AllowResume:           types.BoolValue(api.AllowResume),
		EnforceClientSettings: types.BoolValue(api.EnforceClientSettings),
		AllowKasmAudio:        types.BoolValue(api.AllowKasmAudio),
		AllowKasmUploads:      types.BoolValue(api.AllowKasmUploads),
		AllowKasmDownloads:    types.BoolValue(api.AllowKasmDownloads),
		AllowClipboardDown:    types.BoolValue(api.AllowClipboardDown),
		AllowClipboardUp:      types.BoolValue(api.AllowClipboardUp),
		AllowMicrophone:       types.BoolValue(api.AllowMicrophone),
		AllowSharing:          types.BoolValue(api.AllowSharing),
		AudioDefaultOn:        types.BoolValue(api.AudioDefaultOn),
		IMEModeDefaultOn:      types.BoolValue(api.IMEModeDefaultOn),
	}

	// Handle list field
	if api.AllowedReferrers != nil {
		referrers, diags := types.ListValueFrom(context.Background(), types.StringType, api.AllowedReferrers)
		if !diags.HasError() {
			model.AllowedReferrers = referrers
		} else {
			model.AllowedReferrers = types.ListNull(types.StringType)
		}
	} else {
		model.AllowedReferrers = types.ListNull(types.StringType)
	}

	// Handle numeric fields - set to null if zero
	if api.SessionRemaining != 0 {
		model.SessionRemaining = types.Int64Value(int64(api.SessionRemaining))
	} else {
		model.SessionRemaining = types.Int64Null()
	}

	if api.IPRequestLimit != 0 {
		model.IPRequestLimit = types.Int64Value(int64(api.IPRequestLimit))
	} else {
		model.IPRequestLimit = types.Int64Null()
	}

	if api.IPRequestSeconds != 0 {
		model.IPRequestSeconds = types.Int64Value(int64(api.IPRequestSeconds))
	} else {
		model.IPRequestSeconds = types.Int64Null()
	}

	// Handle string fields - set to null if empty
	if api.ErrorURL != "" {
		model.ErrorURL = types.StringValue(api.ErrorURL)
	} else {
		model.ErrorURL = types.StringNull()
	}

	if api.KasmURL != "" {
		model.KasmURL = types.StringValue(api.KasmURL)
	} else {
		model.KasmURL = types.StringNull()
	}

	if api.ValidUntil != "" {
		model.ValidUntil = types.StringValue(api.ValidUntil)
	} else {
		model.ValidUntil = types.StringNull()
	}

	if api.GroupID != "" {
		model.GroupID = types.StringValue(api.GroupID)
	} else {
		model.GroupID = types.StringNull()
	}

	return model
}
