package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SchemaResource{}

//var _ resource.ResourceWithImportState = &BlockResource{}

func NewSchemaResource() resource.Resource {
	return &SchemaResource{}
}

// ExampleResource defines the resource implementation.
type SchemaResource struct {
	minecraftClient *client
}

// ExampleResourceModel describes the resource data model.
type SchemaResourceModel struct {
	X          types.Number `tfsdk:"x"`
	Y          types.Number `tfsdk:"y"`
	Z          types.Number `tfsdk:"z"`
	Rotation   types.Number `tfsdk:"rotation"`
	Schema     types.String `tfsdk:"schema"`
	SchemaHash types.String `tfsdk:"schema_hash"`
	Id         types.String `tfsdk:"id"`
}

type schemaPlanModifier struct{}

func (s *schemaPlanModifier) Description(ctx context.Context) string {
	return "checks if the file represented by schema has changed"
}

func (s *schemaPlanModifier) MarkdownDescription(ctx context.Context) string {
	return s.Description(ctx)
}

func (s *schemaPlanModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// generate a hash from the file
	attrPath := path.Empty().AtName("schema")
	schema := ""
	req.Plan.GetAttribute(ctx, attrPath, &schema)

	// caclculate and compare the hash
	newHash, _ := calculateHashFromFile(schema)

	// set the new value and set the requires replace, Terraform will force the resource to be re-created
	resp.AttributePlan = types.StringValue(newHash)
	resp.RequiresReplace = true
	resp.Diagnostics.AddWarning(
		"Schema File Changed",
		fmt.Sprintf("The file %s has changed from when the resource was originally created, this forces the destruction of the resource.", schema),
	)
}

func (r *SchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema resource",

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
			"rotation": {
				MarkdownDescription: "Example configurable attribute",
				Required:            false,
				Optional:            true,
				Type:                types.NumberType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"schema": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"schema_hash": {
				Computed:            true,
				MarkdownDescription: "Example configurable attribute",
				Type:                types.StringType,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					&schemaPlanModifier{},
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

func (r *SchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	x, _ := data.X.ValueBigFloat().Int64()
	y, _ := data.Y.ValueBigFloat().Int64()
	z, _ := data.Z.ValueBigFloat().Int64()
	rotation, _ := data.Rotation.ValueBigFloat().Int64()

	sr := schemaRequest{
		X:        int(x),
		Y:        int(y),
		Z:        int(z),
		Rotation: int(rotation),
		Schema:   data.Schema.ValueString(),
	}

	id, err := r.minecraftClient.createSchema(sr)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	strHash, _ := calculateHashFromFile(data.Schema.ValueString())

	data.SchemaHash = types.StringValue(strHash)

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SchemaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		return
	}

	_, err := r.minecraftClient.getSchemaDetails(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read schema, got error: %s", err))
		return
	}
}

func (r *SchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *SchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.minecraftClient.undoSchema(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete schema, got error: %s", err))
		return
	}
}

func calculateHashFromFile(path string) (string, error) {
	// generate a hash of the file so that we can track changes
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Unable to generate hash for schema file: %s", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("Unable to generate hash for schema file: %s", err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

//func (r *EResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
//}
