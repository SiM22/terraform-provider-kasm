package registry

import (
    "context"
    "fmt"
    "time"
    "regexp"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "terraform-provider-kasm/internal/client"
    "terraform-provider-kasm/internal/validators"
)

// Define interfaces implementations
var (
    _ resource.Resource = &registryResource{}
    _ resource.ResourceWithImportState = &registryResource{}
)

// Resource types
type registryResource struct {
    client *client.Client
}

type RegistryResourceModel struct {
    ID             types.String `tfsdk:"id"`
    URL            types.String `tfsdk:"url"`
    OverrideSchema types.String `tfsdk:"override_schema"`
    Channel        types.String `tfsdk:"channel"`
}

func New() resource.Resource {
    return &registryResource{}
}

func (r *registryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_registry"
}

func (r *registryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "url": schema.StringAttribute{
                Required: true,
                Validators: []validator.String{
                    validators.ValidateURL(),
                },
            },
            "channel": schema.StringAttribute{
                Required: true,
                Validators: []validator.String{
                    validateKasmVersion(),
                },
            },
            "override_schema": schema.StringAttribute{
                Optional:    true,
                Description: "Schema override for the registry",
            },
        },
    }
}

func (r *registryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *registryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan RegistryResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    createReq := &client.CreateRegistryRequest{
        Registry:       plan.URL.ValueString(),
        OverrideSchema: plan.OverrideSchema.ValueString(),
        Channel:        plan.Channel.ValueString(),
    }

    err := r.client.CreateRegistry(createReq)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating registry",
            fmt.Sprintf("Could not create registry: %v", err),
        )
        return
    }

    tflog.Info(ctx, "Registry created successfully, attempting to retrieve")

    // Find the created registry
    var createdRegistry *client.Registry
    maxRetries := 2
    for i := 0; i < maxRetries; i++ {
        registries, err := r.client.GetRegistries()
        if err != nil {
            resp.Diagnostics.AddError(
                "Error retrieving registries",
                fmt.Sprintf("Could not get registries: %v", err),
            )
            return
        }
        for _, reg := range registries {
            if reg.RegistryURL == plan.URL.ValueString() {
                createdRegistry = &reg
                break
            }
        }
        if createdRegistry != nil {
            break
        }
        if i < maxRetries-1 {
            tflog.Info(ctx, fmt.Sprintf("Registry not found, retrying in 5 seconds (attempt %d/%d)", i+1, maxRetries))
            time.Sleep(5 * time.Second)
        }
    }

    if createdRegistry == nil {
        resp.Diagnostics.AddError(
            "Error finding created registry",
            "Could not find the registry after creation",
        )
        return
    }

    // Keep the originally planned channel value
    plan.ID = types.StringValue(createdRegistry.RegistryID)
    // Don't update the channel from the API response
    // plan.Channel = types.StringValue(createdRegistry.Channel)

    if createdRegistry.Channel != plan.Channel.ValueString() {
        tflog.Warn(ctx, fmt.Sprintf(
            "API normalized channel %s to version %s. This is expected behavior.",
            plan.Channel.ValueString(),
            createdRegistry.Channel,
        ))
    }

    tflog.Info(ctx, fmt.Sprintf("Created registry with ID: %s, URL: %s, Channel: %s",
        plan.ID.ValueString(), plan.URL.ValueString(), plan.Channel.ValueString()))

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
}

func (r *registryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state RegistryResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    registries, err := r.client.GetRegistries()
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading Kasm Registry",
            fmt.Sprintf("Could not read registry ID %s: %s", state.ID.ValueString(), err),
        )
        return
    }

    // Find the specific registry
    var found bool
    for _, registry := range registries {
        if registry.RegistryID == state.ID.ValueString() {
            state.URL = types.StringValue(registry.RegistryURL)
            // Keep the original channel value instead of using the one from the API
            // state.Channel = types.StringValue(registry.Channel)
            found = true
            break
        }
    }

    if !found {
        resp.State.RemoveResource(ctx)
        return
    }

    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

// Add helper function for version validation
func isValidKasmVersion(version string) bool {
    matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, version)
    return matched
}

func (r *registryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan, state RegistryResourceModel

    // Get plan and state
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Important: Keep the existing channel if it hasn't changed in the plan
    if plan.Channel == state.Channel {
        tflog.Info(ctx, "Channel unchanged in plan, preserving existing value")
        plan.Channel = state.Channel
    }

    // Perform update operations...

    // Keep the plan's channel value rather than using the API response
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
}

func (r *registryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state RegistryResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    err := r.client.DeleteRegistry(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Deleting Kasm Registry",
            fmt.Sprintf("Could not delete registry ID %s: %s", state.ID.ValueString(), err),
        )
        return
    }
}
func (r *registryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    // Import registry using ID
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to validate registry URL format
func validateRegistryURL(url string) error {
    if url == "" {
        return fmt.Errorf("registry URL cannot be empty")
    }
    // Add additional URL validation if needed
    return nil
}

// Helper function to find a registry by ID
func (r *registryResource) findRegistryByID(registryID string) (*client.Registry, error) {
    registries, err := r.client.GetRegistries()
    if err != nil {
        return nil, err
    }

    for _, registry := range registries {
        if registry.RegistryID == registryID {
            return &registry, nil
        }
    }

    return nil, fmt.Errorf("registry with ID %s not found", registryID)
}

// Helper function to find a registry by URL
func (r *registryResource) findRegistryByURL(url string) (*client.Registry, error) {
    registries, err := r.client.GetRegistries()
    if err != nil {
        return nil, err
    }

    for _, registry := range registries {
        if registry.RegistryURL == url {
            return &registry, nil
        }
    }

    return nil, fmt.Errorf("registry with URL %s not found", url)
}
// RegistriesDataSource implementation
type registriesDataSource struct {
    client *client.Client
}

type registriesDataSourceModel struct {
    Registries []registryModel `tfsdk:"registries"`
}

type registryModel struct {
    ID            types.String   `tfsdk:"id"`
    URL           types.String   `tfsdk:"url"`
    AutoUpdate    types.Bool     `tfsdk:"auto_update"`
    SchemaVersion types.String   `tfsdk:"schema_version"`
    IsVerified    types.Bool     `tfsdk:"is_verified"`
    Channel       types.String   `tfsdk:"channel"`
}

func (d *registriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_registries"
}

func (d *registriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = dataschema.Schema{
        Description: "Fetch information about Kasm registries",
        Attributes: map[string]dataschema.Attribute{
            "registries": dataschema.ListNestedAttribute{
                Description: "List of registries",
                Computed:    true,
                NestedObject: dataschema.NestedAttributeObject{
                    Attributes: map[string]dataschema.Attribute{
                        "id": dataschema.StringAttribute{
                            Description: "Registry identifier",
                            Computed:    true,
                        },
                        "url": dataschema.StringAttribute{
                            Description: "Registry URL",
                            Computed:    true,
                        },
                        "auto_update": dataschema.BoolAttribute{
                            Description: "Whether auto-update is enabled",
                            Computed:    true,
                        },
                        "schema_version": dataschema.StringAttribute{
                            Description: "Registry schema version",
                            Computed:    true,
                        },
                        "is_verified": dataschema.BoolAttribute{
                            Description: "Whether the registry is verified",
                            Computed:    true,
                        },
                        "channel": dataschema.StringAttribute{
                            Description: "Registry channel",
                            Computed:    true,
                        },
                    },
                },
            },
        },
    }
}

func (d *registriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state registriesDataSourceModel

    registries, err := d.client.GetRegistries()
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading Kasm Registries",
            fmt.Sprintf("Could not read registries: %v", err),
        )
        return
    }

    // Map response body to model
    for _, registry := range registries {
        registryState := registryModel{
            ID:            types.StringValue(registry.RegistryID),
            URL:           types.StringValue(registry.RegistryURL),
            AutoUpdate:    types.BoolValue(registry.DoAutoUpdate),
            SchemaVersion: types.StringValue(registry.SchemaVersion),
            IsVerified:    types.BoolValue(registry.IsVerified),
            Channel:       types.StringValue(registry.Channel),
        }
        state.Registries = append(state.Registries, registryState)
    }

    // Set state
    diags := resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

// Add a version validator
func validateKasmVersion() validator.String {
    return validators.StringValidator{
        Desc: "must be a valid Kasm version (x.y.z) or channel name (stable, beta, develop)",
        ValidateFn: func(v string) bool {
            // Check if it's a version number format (e.g., 1.16.0)
            versionMatch, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, v)
            if versionMatch {
                return true
            }

            // Check if it's a valid channel name
            validChannels := map[string]bool{
                "stable":  true,
                "beta":    true,
                "develop": true,
            }
            return validChannels[v]
        },
        ErrMessage: "Channel must be a valid Kasm version (x.y.z) or a valid channel name (stable, beta, develop)",
    }
}
