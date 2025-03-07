package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
	groupsds "terraform-provider-kasm/internal/datasources/groups"
	imageds "terraform-provider-kasm/internal/datasources/images"
	rdpds "terraform-provider-kasm/internal/datasources/rdp"
	registryds "terraform-provider-kasm/internal/datasources/registries"
	registryimageds "terraform-provider-kasm/internal/datasources/registry_images"
	sessionstatusds "terraform-provider-kasm/internal/datasources/session_status"
	sessionsds "terraform-provider-kasm/internal/datasources/sessions"
	usersds "terraform-provider-kasm/internal/datasources/users_list"
	zonesds "terraform-provider-kasm/internal/datasources/zones"
	"terraform-provider-kasm/internal/resources/cast"
	"terraform-provider-kasm/internal/resources/group"
	"terraform-provider-kasm/internal/resources/group_image"
	"terraform-provider-kasm/internal/resources/group_membership"
	imageres "terraform-provider-kasm/internal/resources/image"
	"terraform-provider-kasm/internal/resources/join"
	"terraform-provider-kasm/internal/resources/kasm"
	"terraform-provider-kasm/internal/resources/keepalive"
	"terraform-provider-kasm/internal/resources/license"
	"terraform-provider-kasm/internal/resources/login"
	"terraform-provider-kasm/internal/resources/registry"
	"terraform-provider-kasm/internal/resources/session"
	"terraform-provider-kasm/internal/resources/session_permission"
	"terraform-provider-kasm/internal/resources/staging"
	"terraform-provider-kasm/internal/resources/stats"
	"terraform-provider-kasm/internal/resources/user"
)

var _ provider.Provider = &kasmProvider{}

type kasmProvider struct {
	// version string
	client *client.Client
	// Custom server URL for testing
	testServerURL string
}

type kasmProviderModel struct {
	BaseURL   types.String `tfsdk:"base_url"`
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	Insecure  types.Bool   `tfsdk:"insecure"`
}

func New(opts ...string) provider.Provider {
	p := &kasmProvider{}
	if len(opts) > 0 && opts[0] != "" {
		p.testServerURL = opts[0]
	}
	return p
}

func (p *kasmProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kasm"
}

func (p *kasmProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "The base URL of the Kasm API",
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The API key for Kasm",
			},
			"api_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The API secret for Kasm",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip TLS verification",
			},
		},
	}
}

func (p *kasmProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Kasm provider")

	// Retrieve provider data from configuration
	var config kasmProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var baseURL, apiKey, apiSecret string
	var hasConfigErrors bool

	// First check configuration values
	if config.BaseURL.IsNull() {
		// If we have a test server URL, use it
		if p.testServerURL != "" {
			baseURL = p.testServerURL
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("base_url"),
				"Missing Kasm API Base URL",
				"The provider requires a base_url value to be set in the configuration.",
			)
			hasConfigErrors = true
		}
	} else {
		baseURL = config.BaseURL.ValueString()
		// Validate URL format
		if _, err := url.ParseRequestURI(baseURL); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("base_url"),
				"Invalid Base URL Format",
				fmt.Sprintf("The base_url value must be a valid URL: %v", err),
			)
			hasConfigErrors = true
		}
	}

	if config.APIKey.IsNull() {
		// If we have a test server URL, use a dummy API key
		if p.testServerURL != "" {
			apiKey = "test-api-key"
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_key"),
				"Missing Kasm API Key",
				"The provider requires an api_key value to be set in the configuration.",
			)
			hasConfigErrors = true
		}
	} else {
		apiKey = config.APIKey.ValueString()
	}

	if config.APISecret.IsNull() {
		// If we have a test server URL, use a dummy API secret
		if p.testServerURL != "" {
			apiSecret = "test-api-secret"
		} else {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_secret"),
				"Missing Kasm API Secret",
				"The provider requires an api_secret value to be set in the configuration.",
			)
			hasConfigErrors = true
		}
	} else {
		apiSecret = config.APISecret.ValueString()
	}

	// If there are configuration errors, return early
	if hasConfigErrors {
		return
	}

	// Only fall back to environment variables if no configuration was provided
	if baseURL == "" {
		baseURL = os.Getenv("KASM_BASE_URL")
	}
	if apiKey == "" {
		apiKey = os.Getenv("KASM_API_KEY")
	}
	if apiSecret == "" {
		apiSecret = os.Getenv("KASM_API_SECRET")
	}

	tflog.Debug(ctx, fmt.Sprintf("Configuration values - Base URL: %s, API Key: %s", baseURL, apiKey))

	// Create client with provided configuration
	insecure := false
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	tflog.Info(ctx, "Creating Kasm client")
	client := client.NewClient(baseURL, apiKey, apiSecret, insecure)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unable to Create Client",
			fmt.Sprintf("Failed to create Kasm API client with base_url=%s, api_key=%s, api_secret=%s, insecure=%v", baseURL, apiKey, apiSecret, insecure),
		)
		return
	}

	// Store the client in the provider
	p.client = client
	tflog.Info(ctx, "Successfully stored client in provider")

	// Make the client available during DataSource and Resource Configure methods
	if p.client == nil {
		resp.Diagnostics.AddError(
			"Provider Client Not Set",
			"Failed to store client in provider. This is an error in the provider that should be reported to the provider developers.",
		)
		return
	}

	// Set the provider data for resources and data sources
	resp.DataSourceData = p.client
	resp.ResourceData = p.client
	tflog.Info(ctx, "Successfully configured Kasm provider")
}

func (p *kasmProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		user.New,
		group.New,
		cast.New,
		imageres.New,
		registry.New,
		kasm.NewKasmSessionResource,
		session.New,
		login.New,
		license.New,
		staging.New,
		session_permission.New,
		group_image.New,
		group_membership.New,
		join.New,
		stats.NewStatsResource,
		keepalive.NewKeepaliveResource,
	}
}

func (p *kasmProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		imageds.New,
		registryds.New,
		zonesds.New,
		registryimageds.New,
		groupsds.New,
		usersds.New,
		rdpds.NewRDPClientConnectionInfoDataSource,
		sessionsds.New,
		sessionstatusds.New,
	}
}
