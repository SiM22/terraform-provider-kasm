package login

import (
	"context"
	"fmt"

	"terraform-provider-kasm/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &loginResource{}
var _ resource.ResourceWithConfigure = &loginResource{}

type loginResource struct {
	client *client.Client
}

type loginResourceModel struct {
	ID       types.String `tfsdk:"id"`
	UserID   types.String `tfsdk:"user_id"`
	LoginURL types.String `tfsdk:"login_url"`
}

func New() resource.Resource {
	return &loginResource{}
}

func (r *loginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_login"
}

func (r *loginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a login URL for a user to access Kasm without credentials.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource identifier.",
				Computed:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The ID of the user to generate a login URL for.",
				Required:    true,
			},
			"login_url": schema.StringAttribute{
				Description: "The generated login URL.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *loginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *loginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	loginResp, err := r.client.GetLoginURL(plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating login URL",
			fmt.Sprintf("Unable to generate login URL: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(plan.UserID.ValueString())
	plan.LoginURL = types.StringValue(loginResp.URL)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *loginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	loginResp, err := r.client.GetLoginURL(state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading login URL",
			fmt.Sprintf("Unable to read login URL: %v", err),
		)
		return
	}

	state.LoginURL = types.StringValue(loginResp.URL)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *loginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan loginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	loginResp, err := r.client.GetLoginURL(plan.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating login URL",
			fmt.Sprintf("Unable to update login URL: %v", err),
		)
		return
	}

	plan.LoginURL = types.StringValue(loginResp.URL)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *loginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No action needed for deletion as the URL is ephemeral
}
