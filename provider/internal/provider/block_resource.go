package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &BlockResource{}

//var _ resource.ResourceWithImportState = &BlockResource{}

func NewBlockResource() resource.Resource {
	return &BlockResource{}
}

// ExampleResource defines the resource implementation.
type BlockResource struct {
	minecraftClient *client
}

// ExampleResourceModel describes the resource data model.
type BlockResourceModel struct {
	X        types.Number `tfsdk:"x"`
	Y        types.Number `tfsdk:"y"`
	Z        types.Number `tfsdk:"z"`
	Material types.String `tfsdk:"material"`
	Id       types.String `tfsdk:"id"`
}

func (r *BlockResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block"
}

func (r *BlockResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Block resource",

		Attributes: map[string]tfsdk.Attribute{
			"x": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"y": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"z": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"material": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r *BlockResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	minecraftClient, ok := req.ProviderData.(*client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.minecraftClient = minecraftClient
}

func (r *BlockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	x, _ := data.X.ValueBigFloat().Int64()
	y, _ := data.Y.ValueBigFloat().Int64()
	z, _ := data.Z.ValueBigFloat().Int64()

	br := blockRequest{
		X:        int(x),
		Y:        int(y),
		Z:        int(z),
		Material: data.Material.ValueString(),
	}

	err := r.minecraftClient.createBlock(br)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(fmt.Sprintf("%d_%d_%d", x, y, z))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BlockResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	x, _ := data.X.ValueBigFloat().Int64()
	y, _ := data.Y.ValueBigFloat().Int64()
	z, _ := data.Z.ValueBigFloat().Int64()

	err := r.minecraftClient.getBlock(int(x), int(y), int(z))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *BlockResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *BlockResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	x, _ := data.X.ValueBigFloat().Int64()
	y, _ := data.Y.ValueBigFloat().Int64()
	z, _ := data.Z.ValueBigFloat().Int64()

	br := blockRequest{
		X: int(x),
		Y: int(y),
		Z: int(z),
	}

	err := r.minecraftClient.deleteBlock(br)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}
}

//func (r *EResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
//}
