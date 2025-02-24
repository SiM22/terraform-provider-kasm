package group_image

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-kasm/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &groupImageResource{}
	_ resource.ResourceWithImportState = &groupImageResource{}
)

// groupImageResource is the resource implementation
type groupImageResource struct {
	client *client.Client
}

// GroupImageResourceModel maps the resource schema data
type GroupImageResourceModel struct {
	ID                types.String `tfsdk:"id"`
	GroupID           types.String `tfsdk:"group_id"`
	ImageID           types.String `tfsdk:"image_id"`
	GroupImageID      types.String `tfsdk:"group_image_id"`
	ImageName         types.String `tfsdk:"image_name"`
	GroupName         types.String `tfsdk:"group_name"`
	ImageFriendlyName types.String `tfsdk:"image_friendly_name"`
	ImageSrc          types.String `tfsdk:"image_src"`
}

// New creates a new group image resource
func New() resource.Resource {
	return &groupImageResource{}
}

// Metadata returns the resource type name
func (r *groupImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_image"
}

// Schema defines the schema for the resource
func (r *groupImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages image authorization for a group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the group",
			},
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the image to authorize",
			},
			"group_image_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the group image authorization",
			},
			"image_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the image",
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the group",
			},
			"image_friendly_name": schema.StringAttribute{
				Computed:    true,
				Description: "The friendly name of the image",
			},
			"image_src": schema.StringAttribute{
				Computed:    true,
				Description: "The source path of the image",
			},
		},
	}
}

// Configure adds the provider configured client to the resource
func (r *groupImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// waitForImage waits for the image to be available in the group's authorized images
func (r *groupImageResource) waitForImage(groupID string, imageID string) (*client.GroupImage, error) {
	maxAttempts := 10
	delay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		images, err := r.client.GetGroupImages(groupID)
		if err != nil {
			return nil, fmt.Errorf("error reading group images: %v", err)
		}

		for _, img := range images {
			if img.ImageID == imageID {
				return &img, nil
			}
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("image not found in group's authorized images after %d attempts", maxAttempts)
}

// Create creates the resource and sets the initial Terraform state
func (r *groupImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Verify the image exists first by getting all images
	images, err := r.client.GetImages()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting images",
			fmt.Sprintf("Could not get images: %v", err),
		)
		return
	}

	var imageExists bool
	for _, img := range images {
		if img.ImageID == plan.ImageID.ValueString() {
			imageExists = true
			break
		}
	}

	if !imageExists {
		resp.Diagnostics.AddError(
			"Image not found",
			fmt.Sprintf("Image with ID %s does not exist", plan.ImageID.ValueString()),
		)
		return
	}

	err = r.client.AddGroupImage(plan.GroupID.ValueString(), plan.ImageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error authorizing image for group",
			fmt.Sprintf("Could not authorize image: %v", err),
		)
		return
	}

	// Wait for the image to be available
	groupImage, err := r.waitForImage(plan.GroupID.ValueString(), plan.ImageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error finding authorized image",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.GroupID.ValueString(), plan.ImageID.ValueString()))
	plan.GroupImageID = types.StringValue(groupImage.GroupImageID)
	plan.ImageName = types.StringValue(groupImage.ImageName)
	plan.GroupName = types.StringValue(groupImage.GroupName)
	plan.ImageFriendlyName = types.StringValue(groupImage.ImageFriendlyName)
	plan.ImageSrc = types.StringValue(groupImage.ImageSrc)

	tflog.Info(ctx, fmt.Sprintf("Created group image authorization with ID: %s", plan.ID.ValueString()))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data
func (r *groupImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	images, err := r.client.GetGroupImages(state.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Group Images",
			fmt.Sprintf("Could not read group images: %s", err),
		)
		return
	}

	// Find the image in the group's authorized images
	var found bool
	for _, img := range images {
		if img.ImageID == state.ImageID.ValueString() {
			state.GroupImageID = types.StringValue(img.GroupImageID)
			state.ImageName = types.StringValue(img.ImageName)
			state.GroupName = types.StringValue(img.GroupName)
			state.ImageFriendlyName = types.StringValue(img.ImageFriendlyName)
			state.ImageSrc = types.StringValue(img.ImageSrc)
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

// Update updates the resource and sets the updated Terraform state on success
func (r *groupImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Group image authorizations cannot be updated, only created or deleted
	var plan GroupImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *groupImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveGroupImage(state.GroupID.ValueString(), state.ImageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing Group Image Authorization",
			fmt.Sprintf("Could not remove authorization: %s", err),
		)
		return
	}
}

// ImportState imports the resource into Terraform state
func (r *groupImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format group_id:image_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("image_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
