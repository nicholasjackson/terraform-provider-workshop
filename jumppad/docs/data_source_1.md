---
sidebar_position: 9
id: data_source_1
title: Data Source - Creating a Data Source
---

A data source allows you to provide a read only view for resources in your
API. For our Minecraft API we can retrieve the details of a block using
the api path `v1/block/{x}/{y}/{z}`.

[http://minecraft.container.shipyard.run:9090/redoc#tag/Block/operation/getSingleBlock](http://minecraft.container.shipyard.run:9090/redoc#tag/Block/operation/getSingleBlock)

Let's see how to build a data source, like the provider and the resource
you can begin by modifying the existing `example_data_source.go` the template
created for you.

First rename the file `block_data_source.go`, then you can start renaming
the default objects.

## Renaming References

In the `block_data_source.go` file rename the references `ExampleDataSource`
to `BlockDataSource`. Next let's configure the model source that the configuration
will be deserialzied to.

## Creating the BlockDataSourceModel

The `BlockDataSourceModel` looks like the following example. Only the the x,
y, and z attributes will be user configurable. The `id` attribute will be returned
from the API.

Modify the BlockDataSourceModel in the `block_data_source.go` file so that
it looks like the following example.

```go
type BlockDataSourceModel struct {
	X                     types.Number `tfsdk:"x"`
	Y                     types.Number `tfsdk:"y"`
	Z                     types.Number `tfsdk:"z"`
  Material              types.String `tfsdk:"material"`
	Id                    types.String `tfsdk:"id"`
}
```

Next you will define the schema that Terraform will use to populate your model.

## Defining the Schema

The schema is configured in the `GetSchema` method, you define attributes that
are either user configurable or computed in exactly the same way that you 
did when creating the schema resource.

Go ahead and change the `GetSchema` method so that it looks like the following
example:

```go
func (d *BlockDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Block data source",

		Attributes: map[string]tfsdk.Attribute{
			"x": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
			},
			"y": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
			},
			"z": {
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Type:                types.NumberType,
			},
            "material": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}
```

## Configuring the Data Source

In the same way that resources work, any client or properties that are needed
by the datasource are injected by the provider. It does this by calling the
`Configure` method with the `ProviderData` payload.

First change the client on the `BlockDataSource` struct to the following
so that you can use the same Minecraft client as you did with the resources.

```go
type BlockDataSource struct {
	minecraftClient *client
}
```

Next update the `Configure` method to save the client as a local reference. 

```go
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
```

Once all that is done, we need to change the reference in `provider.go` and
then we can compile the provider to check everything is ok.

## Configuring Data Sources in the Provider

Like resources you need to configure which data sources are available in your
provider. This is done using the `DataSources` function defined in the `provider.go` file. Update `DataSources` so that it reflects the new name for your
data source.

```go
func (p *MinecraftProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBlockDataSource,
	}
}
```

Once that is done, give the compiler a quick compile.

```shell
make build
```

And if everything is ok, let now look at how you can implement the `Read` 
method.
