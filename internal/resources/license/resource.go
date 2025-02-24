package license

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-kasm/internal/client"
)

var _ resource.Resource = &licenseResource{}
var _ resource.ResourceWithConfigure = &licenseResource{}

type licenseResource struct {
	client *client.Client
}

type licenseResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ActivationKey types.String `tfsdk:"activation_key"`
	Seats         types.Int64  `tfsdk:"seats"`
	IssuedTo      types.String `tfsdk:"issued_to"`

	// Computed fields from response
	Expiration  types.String `tfsdk:"expiration"`
	IssuedAt    types.String `tfsdk:"issued_at"`
	Limit       types.Int64  `tfsdk:"limit"`
	IsVerified  types.Bool   `tfsdk:"is_verified"`
	LicenseType types.String `tfsdk:"license_type"`
	SKU         types.String `tfsdk:"sku"`

	// Features
	AutoScaling       types.Bool `tfsdk:"auto_scaling"`
	Branding          types.Bool `tfsdk:"branding"`
	SessionStaging    types.Bool `tfsdk:"session_staging"`
	SessionCasting    types.Bool `tfsdk:"session_casting"`
	LogForwarding     types.Bool `tfsdk:"log_forwarding"`
	DeveloperAPI      types.Bool `tfsdk:"developer_api"`
	InjectSSHKeys     types.Bool `tfsdk:"inject_ssh_keys"`
	SAML              types.Bool `tfsdk:"saml"`
	LDAP              types.Bool `tfsdk:"ldap"`
	SessionSharing    types.Bool `tfsdk:"session_sharing"`
	LoginBanner       types.Bool `tfsdk:"login_banner"`
	URLCategorization types.Bool `tfsdk:"url_categorization"`
	UsageLimit        types.Bool `tfsdk:"usage_limit"`
}

func New() resource.Resource {
	return &licenseResource{}
}

func (r *licenseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

func (r *licenseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Kasm license activation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "License identifier.",
				Computed:    true,
			},
			"activation_key": schema.StringAttribute{
				Description: "The activation key provided by Kasm Technologies.",
				Required:    true,
				Sensitive:   true,
			},
			"seats": schema.Int64Attribute{
				Description: "The desired number of seats to license.",
				Optional:    true,
			},
			"issued_to": schema.StringAttribute{
				Description: "Organization the deployment is licensed for.",
				Optional:    true,
			},

			// Computed fields
			"expiration": schema.StringAttribute{
				Description: "License expiration date.",
				Computed:    true,
			},
			"issued_at": schema.StringAttribute{
				Description: "License issue date.",
				Computed:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Licensed seat limit.",
				Computed:    true,
			},
			"is_verified": schema.BoolAttribute{
				Description: "Whether the license is verified.",
				Computed:    true,
			},
			"license_type": schema.StringAttribute{
				Description: "Type of license.",
				Computed:    true,
			},
			"sku": schema.StringAttribute{
				Description: "License SKU.",
				Computed:    true,
			},

			// Feature flags
			"auto_scaling": schema.BoolAttribute{
				Description: "Auto scaling feature flag.",
				Computed:    true,
			},
			"branding": schema.BoolAttribute{
				Description: "Branding feature flag.",
				Computed:    true,
			},
			"session_staging": schema.BoolAttribute{
				Description: "Session staging feature flag.",
				Computed:    true,
			},
			"session_casting": schema.BoolAttribute{
				Description: "Session casting feature flag.",
				Computed:    true,
			},
			"log_forwarding": schema.BoolAttribute{
				Description: "Log forwarding feature flag.",
				Computed:    true,
			},
			"developer_api": schema.BoolAttribute{
				Description: "Developer API feature flag.",
				Computed:    true,
			},
			"inject_ssh_keys": schema.BoolAttribute{
				Description: "SSH key injection feature flag.",
				Computed:    true,
			},
			"saml": schema.BoolAttribute{
				Description: "SAML feature flag.",
				Computed:    true,
			},
			"ldap": schema.BoolAttribute{
				Description: "LDAP feature flag.",
				Computed:    true,
			},
			"session_sharing": schema.BoolAttribute{
				Description: "Session sharing feature flag.",
				Computed:    true,
			},
			"login_banner": schema.BoolAttribute{
				Description: "Login banner feature flag.",
				Computed:    true,
			},
			"url_categorization": schema.BoolAttribute{
				Description: "URL categorization feature flag.",
				Computed:    true,
			},
			"usage_limit": schema.BoolAttribute{
				Description: "Usage limit feature flag.",
				Computed:    true,
			},
		},
	}
}

func (r *licenseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *licenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan licenseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	activateReq := &client.ActivateRequest{
		ActivationKey: plan.ActivationKey.ValueString(),
	}

	if !plan.Seats.IsNull() {
		seats := int(plan.Seats.ValueInt64())
		activateReq.Seats = &seats
	}

	if !plan.IssuedTo.IsNull() {
		activateReq.IssuedTo = plan.IssuedTo.ValueString()
	}

	license, err := r.client.Activate(activateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error activating license",
			fmt.Sprintf("Unable to activate license: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(license.LicenseID)
	plan.Expiration = types.StringValue(license.Expiration)
	plan.IssuedAt = types.StringValue(license.IssuedAt)
	plan.Limit = types.Int64Value(int64(license.Limit))
	plan.IsVerified = types.BoolValue(license.IsVerified)
	plan.LicenseType = types.StringValue(license.LicenseType)
	plan.SKU = types.StringValue(license.SKU)

	// Set feature flags
	plan.AutoScaling = types.BoolValue(license.Features.AutoScaling)
	plan.Branding = types.BoolValue(license.Features.Branding)
	plan.SessionStaging = types.BoolValue(license.Features.SessionStaging)
	plan.SessionCasting = types.BoolValue(license.Features.SessionCasting)
	plan.LogForwarding = types.BoolValue(license.Features.LogForwarding)
	plan.DeveloperAPI = types.BoolValue(license.Features.DeveloperAPI)
	plan.InjectSSHKeys = types.BoolValue(license.Features.InjectSSHKeys)
	plan.SAML = types.BoolValue(license.Features.SAML)
	plan.LDAP = types.BoolValue(license.Features.LDAP)
	plan.SessionSharing = types.BoolValue(license.Features.SessionSharing)
	plan.LoginBanner = types.BoolValue(license.Features.LoginBanner)
	plan.URLCategorization = types.BoolValue(license.Features.URLCategorization)
	plan.UsageLimit = types.BoolValue(license.Features.UsageLimit)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *licenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// License information is only available during activation
	var state licenseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *licenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// License updates require reactivation
	var plan licenseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	activateReq := &client.ActivateRequest{
		ActivationKey: plan.ActivationKey.ValueString(),
	}

	if !plan.Seats.IsNull() {
		seats := int(plan.Seats.ValueInt64())
		activateReq.Seats = &seats
	}

	if !plan.IssuedTo.IsNull() {
		activateReq.IssuedTo = plan.IssuedTo.ValueString()
	}

	license, err := r.client.Activate(activateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error activating license",
			fmt.Sprintf("Unable to activate license: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(license.LicenseID)
	plan.Expiration = types.StringValue(license.Expiration)
	plan.IssuedAt = types.StringValue(license.IssuedAt)
	plan.Limit = types.Int64Value(int64(license.Limit))
	plan.IsVerified = types.BoolValue(license.IsVerified)
	plan.LicenseType = types.StringValue(license.LicenseType)
	plan.SKU = types.StringValue(license.SKU)

	// Set feature flags
	plan.AutoScaling = types.BoolValue(license.Features.AutoScaling)
	plan.Branding = types.BoolValue(license.Features.Branding)
	plan.SessionStaging = types.BoolValue(license.Features.SessionStaging)
	plan.SessionCasting = types.BoolValue(license.Features.SessionCasting)
	plan.LogForwarding = types.BoolValue(license.Features.LogForwarding)
	plan.DeveloperAPI = types.BoolValue(license.Features.DeveloperAPI)
	plan.InjectSSHKeys = types.BoolValue(license.Features.InjectSSHKeys)
	plan.SAML = types.BoolValue(license.Features.SAML)
	plan.LDAP = types.BoolValue(license.Features.LDAP)
	plan.SessionSharing = types.BoolValue(license.Features.SessionSharing)
	plan.LoginBanner = types.BoolValue(license.Features.LoginBanner)
	plan.URLCategorization = types.BoolValue(license.Features.URLCategorization)
	plan.UsageLimit = types.BoolValue(license.Features.UsageLimit)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *licenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// License deletion is handled by Kasm
}
