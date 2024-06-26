// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SchemaResource{}
var _ resource.ResourceWithImportState = &SchemaResource{}

func NewSchemaResource() resource.Resource {
	return &SchemaResource{}
}

// SchemaResource defines the resource implementation.
type SchemaResource struct {
	minecraftClient *client
}

// SchemaResourceModel describes the resource data model.
type SchemaResourceModel struct {
	X          types.Number `tfsdk:"x"`
	Y          types.Number `tfsdk:"y"`
	Z          types.Number `tfsdk:"z"`
	Rotation   types.Number `tfsdk:"rotation"`
	Schema     types.String `tfsdk:"schema"`
	SchemaHash types.String `tfsdk:"schema_hash"`
	Id         types.String `tfsdk:"id"`
}

func (r *SchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"x": schema.NumberAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.Number{
					numberplanmodifier.RequiresReplace(),
				},
			},
			"y": schema.NumberAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.Number{
					numberplanmodifier.RequiresReplace(),
				},
			},
			"z": schema.NumberAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.Number{
					numberplanmodifier.RequiresReplace(),
				},
			},
			"rotation": schema.NumberAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.Number{
					numberplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema_hash": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					&schemaPlanModifier{},
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.minecraftClient = client
}

func (r *SchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SchemaResourceModel

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

	data.Id = types.StringValue(id)

	hash, err := calculateHashFromFile(data.Schema.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to generate hash for file", err.Error())
	}
	data.SchemaHash = types.StringValue(hash)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SchemaResourceModel

	// Read Terraform prior state data into the model
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SchemaResourceModel

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

func (r *SchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func calculateHashFromFile(path string) (string, error) {
	// generate a hash of the file so that we can track changes
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to generate hash for schema file: %s", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("unable to generate hash for schema file: %s", err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

type schemaPlanModifier struct{}

func (s schemaPlanModifier) Description(ctx context.Context) string {
	return "checks if the file represented by schema has changed"
}

func (s schemaPlanModifier) MarkdownDescription(ctx context.Context) string {
	return s.Description(ctx)
}

func (s schemaPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// generate a hash from the file
	attrPath := path.Empty().AtName("schema")
	schema := ""
	req.Plan.GetAttribute(ctx, attrPath, &schema)

	attrPath = path.Empty().AtName("schema_hash")
	schemaHash := ""
	req.Plan.GetAttribute(ctx, attrPath, &schemaHash)

	// if empty probably first apply
	if schemaHash == "" {
		return
	}

	// caclculate and compare the hash
	newHash, _ := calculateHashFromFile(schema)
	if newHash != schemaHash {
		// set the new value and set the requires replace, Terraform will force the resource to be re-created
		resp.PlanValue = types.StringValue(newHash)
		resp.RequiresReplace = true
		resp.Diagnostics.AddWarning(
			"Schema File Changed",
			fmt.Sprintf("The file %s has changed from when the resource was originally created, this forces the destruction of the resource. Old file hash: %s, New file hash: %s", schema, schemaHash, newHash),
		)
	}
}
