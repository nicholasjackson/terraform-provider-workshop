---
id: schema
title: Schema Resource
---

The schema resource interacts with the schema capabilities of the Minecraft
API.

[http://localhost:9090/redoc#tag/Schema](http://localhost:9090/redoc#tag/Schema)

If allows you to place a multiple blocks defined in a zipped JSON file.

## Creating a Schema using the API

Let's test this by curling the API manually, open your Minecraft client either
in the browser or with the desktop client and log in.

All blocks in Minecraft are placed in 3D space, x, y, and z. Where `x and z` are
horizontal coordinates and `y` is a vertical coordinate.

You can determine your position at any time by using the command 

```
/position
```

![](./images/schema/position.jpg)

Let's create a car at the car at our current location. The schema for a car
can be found in the ./schema/car.zip file. This will be the API that 
Terraform uses when creating resources.

```shell
➜ curl -H 'X-API-Key: supertopsecret' http://minecraft.container.shipyard.run:9090/v1/schema/-1260/24/288/0 --data-binary @./schema/car.zip
```

The API request will return an ID of the schema that can be used to undo the 
operation and to retrieve the details for it.

```
1668343502799
```

![](./images/schema/create.jpg)

## Retrieving the Details of a Schema

You can retrieve the details of the schema that you have just applied using
the following command, remember to replace the id with your own id. This will
be the API that will be called when Terraform reads information for the resource.

Enter the following details in your terminal, replacing the `[id]` with the
id returned from the previous request.

```shell
curl -H 'X-API-Key: supertopsecret' http://minecraft.container.shipyard.run:9090/v1/schema/details/[id]
```

```json
{"startX":-1260,"startY":24,"startZ":288,"endX":-1252,"endY":34,"endZ":292}
```

## Removing a Schema

The API includes a method `v1/schema/undo/[schema id]` that allows you to undo
the application of the schema created in a previous step. This will be the
API that your provider will use when destroying a schema resource.

Execute the following command in your terminal, replacing the `[id]` with the
id returned from the previous request.

```shell
curl -H 'X-API-Key: supertopsecret' -XDELETE http://minecraft.container.shipyard.run:9090/v1/schema/undo/[id]
```

The car you created earlier will have been removed.

![](./images/schema/remove.jpg)

Let's now start codifying this as a provider resource.

## Creating the Examples

Before you write any code first, let's define the resources as HCL examples.

Rather than creating new files you can modify the files created by the
template. First let's rename the folder 

`./examples/resources/scaffolding_example`

to 

`./examples/resources/minecraft_schema`

### Adding the provider stanza

Then let's modify the `resource.tf` file inside it. You can remove all the
existing content, and then let's add a `provider block.

The provider stanza defines the providers and their versions that will be
used by the config. The following block defines a provider requirement
for the provider that is built by this example. 

:::note 
Source is prefixed with `local` this differentiates it from a provider
that is pulled from the Terraform registry. This corresponds to the `local`
folder where the plugin is installed.

`~/.terraform.d/plugins/local/`
:::

```javascript
terraform {
  required_providers {
    minecraft = {
      source  = "local/hashicraft/minecraft"
      version = "0.1.0"
    }
  }
}

```

Next let's define the provider stanza, the API needs a key for authentication
you will add that and the location of the API as custom parameters in the
provider.

```javascript
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
```

Finally, create the `minecraft_schema` resource, once this is done, you can
begin writing the code that will interact with these.

```javascript
resource "minecraft_schema" "bus" {
  x = -1278
  y = 24
  z = 288
  rotation = 270
  schema = "../../../example_schemas/car.zip"
}
```

We will use these when testing our provider but, first you need to write the
code. Let's do that now.

## Creating the Schema Resource

Like you did with the examples, let's modify an existing file for creating 
the schema resource. First, rename the file

`internal/provider/example_resource.go`

to 

`internal/provider/schema_resource.go`

Then let's rename all the `ExampleResource` references to `SchemaResource`
you will need to change all of the references.

Next, you need to define the `SchemaResourceModel`, this struct is where
Terraform will serialize the HCL configurations to. It uses struct tags
similar to defining tags for JSON serialization and deserialzation.

The model needs to have all the fields that are going to be used when
Terraform processes your config or serializes it to the state.
For the `schema` example, it will look like the following.

```go
type SchemaResourceModel struct {
	X          types.Number `tfsdk:"x"`
	Y          types.Number `tfsdk:"y"`
	Z          types.Number `tfsdk:"z"`
	Rotation   types.Number `tfsdk:"rotation"`
	Schema     types.String `tfsdk:"schema"`
	SchemaHash types.String `tfsdk:"schema_hash"`
	Id         types.String `tfsdk:"id"`
}
```

:::note
The `SchemaResourceModel` does not use standard Go types, it uses SDK types.
When Terraform deserializes your configuration it needs to differentiate
from a missing value from an empty value. In Go a defined type has a base
value it is not nil unless it is a reference.
:::

Next, let's modify the `Metadata` the Metadata method is called by Terraform
when it creates your resource. It tells Terraform how to link your resource
type to the configuration.

```go
func (r *SchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}
```

## Defining the schema

Next, you need to define the schema that allows you to define the validation
rules like optional or computed attributes. Computed attributes are set by
the provider; they are not set by the end user.

### Required attribute

A basic example of a required attribute looks like the following. This is
a number type that when changed in the state will force Terraform to
perform a destroy operation on the provider before it creates a new resouce.

In the Schema example it is not possible to mutate a schema, any elements that
change need to force a replace.

```go
"x": {
		MarkdownDescription: "Example configurable attribute",
		Required:            true,
		Type:                types.NumberType,
		PlanModifiers: []tfsdk.AttributePlanModifier{
			resource.RequiresReplace(),
		},
	},
```

### Computed attribute

An example of a computed attribute will be the `schema_hash`, this attribute
is computed by the provider and is a hash of the file used to define the schema.
Storing the hash is important as should the file that defines the schema 
change you need to replace the resource. Storing only the file name is not
enough as the file could be replaced with one of the same name.

```go
"schema_hash": {
	Computed:            true,
	MarkdownDescription: "Example configurable attribute",
	Type:                types.StringType,
	PlanModifiers: []tfsdk.AttributePlanModifier{
		&schemaPlanModifier{},
		resource.RequiresReplace(),
	},
},
```

Replace the GetSchema method with the following example.

```go
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
```

## Configuring the API Client

To interact with the API you need to know the endpoint and the APIKey, ideally,
each resource will not create the client itself. It is better that the client
is injected into the resource.

Terraform has a mechanism for this with the Configure method, the example 
code already has a client defined in the file `client.go`. This wraps the 
RESTFul API with Go types and methods. You will learn how this is configured
when we look at the Provider type. 

```go
func (r *SchemaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
```