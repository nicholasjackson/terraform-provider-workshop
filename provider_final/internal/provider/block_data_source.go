// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &BlockDataSource{}

func NewBlockDataSource() datasource.DataSource {
	return &BlockDataSource{}
}

// BlockDataSource defines the data source implementation.
type BlockDataSource struct {
	minecraftClient *client
}

// BlockDataSourceModel describes the data source data model.
type BlockDataSourceModel struct {
	X        types.Number `tfsdk:"x"`
	Y        types.Number `tfsdk:"y"`
	Z        types.Number `tfsdk:"z"`
	Material types.String `tfsdk:"material"`
	Id       types.String `tfsdk:"id"`
}

func (d *BlockDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block"
}

func (d *BlockDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Block data source",

		Attributes: map[string]schema.Attribute{
			"x": schema.NumberAttribute{
				MarkdownDescription: "Block configurable attribute",
				Required:            true,
			},
			"y": schema.NumberAttribute{
				MarkdownDescription: "Block configurable attribute",
				Required:            true,
			},
			"z": schema.NumberAttribute{
				MarkdownDescription: "Block configurable attribute",
				Required:            true,
			},
			"material": schema.StringAttribute{
				MarkdownDescription: "Block configurable attribute",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Block identifier",
				Computed:            true,
			},
		},
	}
}

func (d *BlockDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	minecraftClient, ok := req.ProviderData.(*client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.minecraftClient = minecraftClient
}

func (d *BlockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BlockDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	x, _ := data.X.ValueBigFloat().Int64()
	y, _ := data.Y.ValueBigFloat().Int64()
	z, _ := data.Z.ValueBigFloat().Int64()

	block, err := d.minecraftClient.getBlock(int(x), int(y), int(z))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve block",
			fmt.Sprintf("Unable to get block, got error: %s", err),
		)
		return
	}

	data.Material = types.StringValue(block.Material)
	data.Id = types.StringValue(block.ID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
