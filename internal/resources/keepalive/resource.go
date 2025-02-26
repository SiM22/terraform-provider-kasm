package keepalive

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &keepaliveResource{}
)

// NewKeepaliveResource is a helper function to simplify the provider implementation.
func NewKeepaliveResource() resource.Resource {
	return &keepaliveResource{}
}

// keepaliveResource is the resource implementation.
type keepaliveResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *keepaliveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keepalive"
}

// Schema defines the schema for the resource.
func (r *keepaliveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"kasm_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *keepaliveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *keepaliveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KeepaliveResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set a unique ID before making the API call
	plan.ID = types.StringValue(plan.KasmID.ValueString())

	// Make the keepalive API call
	_, err := r.client.Keepalive(plan.KasmID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error sending keepalive",
			fmt.Sprintf("Could not send keepalive: %v", err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *keepaliveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KeepaliveResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *keepaliveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KeepaliveResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make the keepalive API call
	_, err := r.client.Keepalive(plan.KasmID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error sending keepalive",
			fmt.Sprintf("Could not send keepalive: %v", err),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *keepaliveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KeepaliveResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No API call needed for deletion
}

// KeepaliveResourceModel maps the resource schema data.
type KeepaliveResourceModel struct {
	ID     types.String `tfsdk:"id"`
	KasmID types.String `tfsdk:"kasm_id"`
}
