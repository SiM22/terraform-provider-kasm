package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionsDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionsDataSource{}
)

// New is a helper function to simplify the provider implementation.
func New() datasource.DataSource {
	return &sessionsDataSource{}
}

// sessionsDataSource is the data source implementation.
type sessionsDataSource struct {
	client *client.Client
}

// sessionsDataSourceModel maps the data source schema data.
type sessionsDataSourceModel struct {
	ID          types.String            `tfsdk:"id"`
	CurrentTime types.String            `tfsdk:"current_time"`
	Sessions    []sessionModel          `tfsdk:"sessions"`
	SessionsMap map[string]sessionModel `tfsdk:"sessions_map"`
}

// sessionModel maps session schema data.
type sessionModel struct {
	ExpirationDate      types.String            `tfsdk:"expiration_date"`
	ContainerIP         types.String            `tfsdk:"container_ip"`
	StartDate           types.String            `tfsdk:"start_date"`
	Token               types.String            `tfsdk:"token"`
	ImageID             types.String            `tfsdk:"image_id"`
	ViewOnlyToken       types.String            `tfsdk:"view_only_token"`
	Cores               types.Float64           `tfsdk:"cores"`
	Hostname            types.String            `tfsdk:"hostname"`
	KasmID              types.String            `tfsdk:"kasm_id"`
	PortMap             map[string]types.String `tfsdk:"port_map"`
	ImageName           types.String            `tfsdk:"image_name"`
	ImageFriendlyName   types.String            `tfsdk:"image_friendly_name"`
	ImageSrc            types.String            `tfsdk:"image_src"`
	IsPersistentProfile types.Bool              `tfsdk:"is_persistent_profile"`
	Memory              types.Int64             `tfsdk:"memory"`
	OperationalStatus   types.String            `tfsdk:"operational_status"`
	ContainerID         types.String            `tfsdk:"container_id"`
	Port                types.Int64             `tfsdk:"port"`
	KeepaliveDate       types.String            `tfsdk:"keepalive_date"`
	UserID              types.String            `tfsdk:"user_id"`
	ShareID             types.String            `tfsdk:"share_id"`
	Host                types.String            `tfsdk:"host"`
	ServerID            types.String            `tfsdk:"server_id"`
}

// Metadata returns the data source type name.
func (d *sessionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sessions"
}

// Schema defines the schema for the data source.
func (d *sessionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of active Kasm sessions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for the data source.",
				Computed:    true,
			},
			"current_time": schema.StringAttribute{
				Description: "Current server time when the sessions were fetched.",
				Computed:    true,
			},
			"sessions": schema.ListNestedAttribute{
				Description: "List of active Kasm sessions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"expiration_date": schema.StringAttribute{
							Description: "Date when the session will expire.",
							Computed:    true,
						},
						"container_ip": schema.StringAttribute{
							Description: "IP address of the container.",
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Date when the session was started.",
							Computed:    true,
						},
						"token": schema.StringAttribute{
							Description: "Session token.",
							Computed:    true,
							Sensitive:   true,
						},
						"image_id": schema.StringAttribute{
							Description: "ID of the image used for the session.",
							Computed:    true,
						},
						"view_only_token": schema.StringAttribute{
							Description: "Token for view-only access.",
							Computed:    true,
							Sensitive:   true,
						},
						"cores": schema.Float64Attribute{
							Description: "Number of CPU cores allocated to the session.",
							Computed:    true,
						},
						"hostname": schema.StringAttribute{
							Description: "Hostname of the session container.",
							Computed:    true,
						},
						"kasm_id": schema.StringAttribute{
							Description: "Unique identifier for the Kasm session.",
							Computed:    true,
						},
						"port_map": schema.MapAttribute{
							Description: "Mapping of service names to ports.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"image_name": schema.StringAttribute{
							Description: "Name of the image used for the session.",
							Computed:    true,
						},
						"image_friendly_name": schema.StringAttribute{
							Description: "User-friendly name of the image.",
							Computed:    true,
						},
						"image_src": schema.StringAttribute{
							Description: "Source path for the image thumbnail.",
							Computed:    true,
						},
						"is_persistent_profile": schema.BoolAttribute{
							Description: "Whether the session has a persistent profile.",
							Computed:    true,
						},
						"memory": schema.Int64Attribute{
							Description: "Amount of memory allocated to the session in bytes.",
							Computed:    true,
						},
						"operational_status": schema.StringAttribute{
							Description: "Current operational status of the session.",
							Computed:    true,
						},
						"container_id": schema.StringAttribute{
							Description: "ID of the container running the session.",
							Computed:    true,
						},
						"port": schema.Int64Attribute{
							Description: "Main port used by the session.",
							Computed:    true,
						},
						"keepalive_date": schema.StringAttribute{
							Description: "Date of the last keepalive signal.",
							Computed:    true,
						},
						"user_id": schema.StringAttribute{
							Description: "ID of the user who owns the session.",
							Computed:    true,
						},
						"share_id": schema.StringAttribute{
							Description: "ID used for sharing the session.",
							Computed:    true,
						},
						"host": schema.StringAttribute{
							Description: "Host where the session is running.",
							Computed:    true,
						},
						"server_id": schema.StringAttribute{
							Description: "ID of the server running the session.",
							Computed:    true,
						},
					},
				},
			},
			"sessions_map": schema.MapNestedAttribute{
				Description: "Map of active Kasm sessions keyed by kasm_id.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"expiration_date": schema.StringAttribute{
							Description: "Date when the session will expire.",
							Computed:    true,
						},
						"container_ip": schema.StringAttribute{
							Description: "IP address of the container.",
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Date when the session was started.",
							Computed:    true,
						},
						"token": schema.StringAttribute{
							Description: "Session token.",
							Computed:    true,
							Sensitive:   true,
						},
						"image_id": schema.StringAttribute{
							Description: "ID of the image used for the session.",
							Computed:    true,
						},
						"view_only_token": schema.StringAttribute{
							Description: "Token for view-only access.",
							Computed:    true,
							Sensitive:   true,
						},
						"cores": schema.Float64Attribute{
							Description: "Number of CPU cores allocated to the session.",
							Computed:    true,
						},
						"hostname": schema.StringAttribute{
							Description: "Hostname of the session container.",
							Computed:    true,
						},
						"kasm_id": schema.StringAttribute{
							Description: "Unique identifier for the Kasm session.",
							Computed:    true,
						},
						"port_map": schema.MapAttribute{
							Description: "Mapping of service names to ports.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"image_name": schema.StringAttribute{
							Description: "Name of the image used for the session.",
							Computed:    true,
						},
						"image_friendly_name": schema.StringAttribute{
							Description: "User-friendly name of the image.",
							Computed:    true,
						},
						"image_src": schema.StringAttribute{
							Description: "Source path for the image thumbnail.",
							Computed:    true,
						},
						"is_persistent_profile": schema.BoolAttribute{
							Description: "Whether the session has a persistent profile.",
							Computed:    true,
						},
						"memory": schema.Int64Attribute{
							Description: "Amount of memory allocated to the session in bytes.",
							Computed:    true,
						},
						"operational_status": schema.StringAttribute{
							Description: "Current operational status of the session.",
							Computed:    true,
						},
						"container_id": schema.StringAttribute{
							Description: "ID of the container running the session.",
							Computed:    true,
						},
						"port": schema.Int64Attribute{
							Description: "Main port used by the session.",
							Computed:    true,
						},
						"keepalive_date": schema.StringAttribute{
							Description: "Date of the last keepalive signal.",
							Computed:    true,
						},
						"user_id": schema.StringAttribute{
							Description: "ID of the user who owns the session.",
							Computed:    true,
						},
						"share_id": schema.StringAttribute{
							Description: "ID used for sharing the session.",
							Computed:    true,
						},
						"host": schema.StringAttribute{
							Description: "Host where the session is running.",
							Computed:    true,
						},
						"server_id": schema.StringAttribute{
							Description: "ID of the server running the session.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *sessionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *sessionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionsDataSourceModel

	// Get sessions from the API
	tflog.Debug(ctx, "Getting sessions from Kasm API")
	sessionsResp, err := d.client.GetKasms()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kasm Sessions",
			fmt.Sprintf("Could not read Kasm sessions: %s", err),
		)
		return
	}

	// Set state values
	state.ID = types.StringValue("kasm_sessions")

	// Set current_time - if it's empty from the API, use the current time
	if sessionsResp.CurrentTime != "" {
		state.CurrentTime = types.StringValue(sessionsResp.CurrentTime)
	} else {
		// Use current time as fallback
		state.CurrentTime = types.StringValue(time.Now().UTC().Format(time.RFC3339))
		tflog.Debug(ctx, "CurrentTime not provided by API, using current time")
	}

	state.Sessions = make([]sessionModel, 0, len(sessionsResp.Kasms))
	state.SessionsMap = make(map[string]sessionModel, len(sessionsResp.Kasms))

	// Map response body to model
	for _, kasm := range sessionsResp.Kasms {
		sessionState := sessionModel{
			ExpirationDate:      types.StringValue(kasm.ExpirationDate),
			ContainerIP:         types.StringValue(kasm.ContainerIP),
			StartDate:           types.StringValue(kasm.StartDate),
			Token:               types.StringValue(kasm.Token),
			ImageID:             types.StringValue(kasm.ImageID),
			ViewOnlyToken:       types.StringValue(kasm.ViewOnlyToken),
			Cores:               types.Float64Value(kasm.Cores),
			Hostname:            types.StringValue(kasm.Hostname),
			KasmID:              types.StringValue(kasm.KasmID),
			IsPersistentProfile: types.BoolValue(kasm.IsPersistentProfile),
			Memory:              types.Int64Value(kasm.Memory),
			OperationalStatus:   types.StringValue(kasm.OperationalStatus),
			ContainerID:         types.StringValue(kasm.ContainerID),
			Port:                types.Int64Value(int64(kasm.Port)),
			KeepaliveDate:       types.StringValue(kasm.KeepaliveDate),
			UserID:              types.StringValue(kasm.UserID),
			ShareID:             types.StringValue(kasm.ShareID),
			Host:                types.StringValue(kasm.Host),
			ServerID:            types.StringValue(kasm.ServerID),
		}

		// Add image information if available
		if kasm.Image.ImageID != "" {
			sessionState.ImageName = types.StringValue(kasm.Image.Name)
			sessionState.ImageFriendlyName = types.StringValue(kasm.Image.FriendlyName)
			sessionState.ImageSrc = types.StringValue(kasm.Image.ImageSrc)
		}

		// Add port map if available
		portMap := make(map[string]types.String)

		// Convert the port map to a string representation
		if kasm.PortMap.Audio.Port != 0 {
			portMap["audio"] = types.StringValue(fmt.Sprintf("%d", kasm.PortMap.Audio.Port))
		}
		if kasm.PortMap.VNC.Port != 0 {
			portMap["vnc"] = types.StringValue(fmt.Sprintf("%d", kasm.PortMap.VNC.Port))
		}
		if kasm.PortMap.AudioInput.Port != 0 {
			portMap["audio_input"] = types.StringValue(fmt.Sprintf("%d", kasm.PortMap.AudioInput.Port))
		}

		sessionState.PortMap = portMap

		state.Sessions = append(state.Sessions, sessionState)
		state.SessionsMap[kasm.KasmID] = sessionState
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
