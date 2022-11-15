---
sidebar_position: 3
id: schema_1
title: Schema Resource - Create Method
---

Providers are Go types that implement the `resource.Resource` interface.
This interface requires that your type has specific methods that Terraform
will call to process Terraform configuration and convert it into API calls.

The full interface is listed below:

```go
type Resource interface {
	// Metadata should return the full name of the resource, such as
	// examplecloud_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// GetSchema returns the schema for this resource.
	GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics)

	// Create is called when the provider must create a new resource. Config
	// and planned state values should be read from the
	// CreateRequest and new state values set on the CreateResponse.
	Create(context.Context, CreateRequest, *CreateResponse)

	// Read is called when the provider must read resource values in order
	// to update state. Planned state values should be read from the
	// ReadRequest and new state values set on the ReadResponse.
	Read(context.Context, ReadRequest, *ReadResponse)

	// Update is called to update the state of the resource. Config, planned
	// state, and prior state values should be read from the
	// UpdateRequest and new state values set on the UpdateResponse.
	Update(context.Context, UpdateRequest, *UpdateResponse)

	// Delete is called when the provider must delete the resource. Config
	// values may be read from the DeleteRequest.
	//
	// If execution completes without error, the framework will automatically
	// call DeleteResponse.State.RemoveResource(), so it can be omitted
	// from provider logic.
	Delete(context.Context, DeleteRequest, *DeleteResponse)
}
```

You are going to build a resource that interacts with the schema capabilities
of the Minecraft API.

[http://localhost:9090/redoc#tag/Schema](http://localhost:9090/redoc#tag/Schema)

It allows you to place multiple blocks defined in a zipped JSON file.

## Creating a Schema using the API

Let's test this by manually curling the API. Open your Minecraft client either 
in the browser or with the desktop client and log in.

All blocks in Minecraft are placed in 3D space, x, y, and z. Where `x and z` are
horizontal coordinates and `y` is a vertical coordinate.

You can determine your position at any time by using the command:

```
/position
```

![](./images/schema/position.jpg)

Let's create a car at your current location. The schema for a car is located 
in the ./schema/car.zip file. Using `curl` you are simulating what Terraform
will do when a `terraform apply` runs containing one of your custom resources.

```shell
➜ curl -H 'X-API-Key: supertopsecret' http://minecraft.container.shipyard.run:9090/v1/schema/-1260/24/288/0 --data-binary @./schema/car.zip
```

After executing the API request will return an ID of the schema that can be used 
to undo the operation and retrieve the details for it.

```
1668343502799
```

![](./images/schema/create.jpg)

## Retrieving the Details of a Schema

You can retrieve the details of the schema you have just applied using
the following command; remember to replace the id with your id. This will
be the API called when Terraform reads information for the resource.

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

The car you created earlier will now have been removed.

![](./images/schema/remove.jpg)

Let's now start codifying this as a provider resource.

## Creating the Examples

Before you write any code, first, let's define the resources as HCL examples.

Rather than creating new files, you can modify the files created by the
template. First, let's rename the folder 

`./examples/resources/scaffolding_example`

to 

`./examples/resources/minecraft_schema`

### Adding the provider stanza

Then let's modify the `resource.tf` file inside it. You can remove all the
existing content, and then let's add a `provider block.

The provider stanza defines the providers and their versions that the config will use. 
The following block represents the configuration your provider.

:::note 
Source is prefixed with `local` this differentiates it from a provider
that is pulled from the Terraform registry. This corresponds to the `local`
folder where the plugin is installed.

`~/.terraform.d/plugins/local/`
:::

Next, let's define the provider stanza; the API needs a key for authentication. 
You will add that and the location of the API as custom parameters in the
provider.

```javascript
provider "minecraft" {
  endpoint = "http://minecraft.container.shipyard.run:9090"
  api_key = "supertopsecret"
}
```

Finally, create the `minecraft_schema` resource; once completed, you can 
begin writing the code that will interact with these.

We will use these when testing our provider, but you first need to write the 
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
Terraform will serialize the HCL configurations. It uses struct tags
similar to defining tags for JSON serialization and deserialzation.

The model must have all the fields used when Terraform processes your config
or serializes it to the state. For the `schema` example, it will look like the following.

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
The `SchemaResourceModel` does not use standard Go types; it uses SDK types.
When Terraform deserializes your configuration, it needs to differentiate a missing 
value from an empty value. In Go, a defined type has a base value. It is not nil 
unless it is a reference.
:::

Next, let's modify the `Metadata`. The Metadata method is called by Terraform
when it creates your resource. It tells Terraform how to link your resource
type to the configuration.

```go
func (r *SchemaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}
```

## Defining the schema

Next, you need to define the schema that allows you to specify the validation 
rules, like optional or computed attributes. Computed attributes are set by 
your code; the end user does not set them.

### Required attribute

A basic example of a required attribute looks like the following. This is
a number type that, when changed in the state, will force Terraform to
perform a destroy operation on the provider before it creates a new resource.

In the Schema API example, it is impossible to mutate a schema; any elements that
change need to force a replacement. This is set by using the `PlanModifiers` and
the `resource.RequiresReplace` function.

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

An example of a computed attribute is the `schema_hash``; this attribute is 
calculated by the provider and is a hash of the file used to define the schema.
Storing the hash is vital, as should the file that defines the schema 
change you need to replace the resource. Keeping only the file name is not enough, 
as the file could be replaced with one of the same name.

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

To interact with the API, you need to know the endpoint and the APIKey. Ideally, 
each resource will not create the client itself. The client should be injected 
into the resource, and the SDK has a mechanism for this with the Configure method; 
the example code already has a client defined in the file `client.go`. This wraps the 
RESTFul API with Go types and methods. You will learn how this is configured
when we look at the Provider type. 

For now, rename `client` in the `SchemaResource` struct fields to `minecraftClient`
and change the type to `*client`, which is the go client for the API
from this package. You can then update the configure block to look like the 
following example.

```go
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
```

## Creating Resources

When Terraform wants to create a resource, it calls the `Create` method on
your provider.  In this method, you make any necessary API calls to create 
the things your resource defines. In our example, this will be to call 
the `v1/schema` API on the Minecraft server with the correct
details.

If you look at the `Create` method, you will see the following code at the
top.

```go
var data *SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
```

What this is doing is deserializing the Terraform resources written in HCL
and converting them into your data model.

The `Get` method returns a `diag.Diagnostics` type, this is an SDK type that
is used for rich error messages and logging. When building a provider, you 
will use this type instead of logging to StdOut or returning errors from methods. 

`diag.`Diagnostics` ensures that any log output or error is correctly formatted
so that the Terraform application can display it consistently for all providers.

To apply a schema, you will use the `createSchema` method on the 
`client`. This method accepts parameters which are `schemaRequest`.

```go
type schemaRequest struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Z        int    `json:"z"`
	Rotation int    `json:"rotation"`
	Schema   string `json:"schema"`
}
```

The types for `schemaRequest` are plain go types so to use the fields from
`data` you will need to convert them. To get the go type from a SDK type
you can use the sdk.Type methods as follows:

#### converting to int64

```go
data.X.ValueBigFloat().Int64()
```

#### converting to string

```go
data.Schema.ValueString(),
```

:::note
The reason the SDK uses sdk.Type rather than Go types is to differentate
between not set and the Go types default value.

If `data.X` has not been set such as when it is an optional attribute then
the method call `VauleBigFloat()` will panic. You should always code 
defensively against nil values. In our example this should be safe
as all the attributes are `required`.
:::

Add the following code block to your `Create` method beneath 
`if resp.Diagnositcs.HasError() {}`.

```go
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
```

Now you can call the `createSchema` method to make the API call, this method
returns an error so you need to halt all further processing and return
should you get an error back.

Add the following beneath your `schemaRequest`.

```go
id, err := r.minecraftClient.createSchema(sr)
if err != nil {
	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	return
}
```

To handle changes to the underlying schema file, you need to store a hash of
the file in the Terraform state. The following code block is a pretty standard
Go way to calculate a hash from a file. Add this to the bottom of your `schema_resource.go`
file.

```go
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
```

You can now call this method and set the returned value to the `data` model that
will be persisted to the state. Add this code to your `Create` method beneath 
the last block you just added.

```go
	strHash, _ := calculateHashFromFile(data.Schema.ValueString())

	data.SchemaHash = types.StringValue(strHash)
  data.Id = types.StringValue(id)
```

:::note
To convert a Go type to a `sdk.Type` you can use the methods

`types.StringValue("mystring")`

and

`types.Float64(23)`
:::

That completes the `Create` method; if you look at the final line from your
`Create` method, you will see the following.

```go
resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
```

This code is saving your `SchemaResourceModel` to the Terraform state
it is essential that you save your model with the unique id, or Terraform
will not know how to delete or update it.

## Custom Plan Modifier

During this example, we have mentioned that you need to keep track of the 
file hash for the schema so that Terraform can detect changes and recreate
the resource. So far, you have implemented the tracking of the hash, but there 
are no checks.

Let's now add this check, you are going to add the check as a plan modifier
a block of code that executes when Terraform runs a plan. At this point
you can check the existing stored hash with a new computed value from the 
schema file in the new plan.

To create a plan modifier, you need to create a type that implements the 
`AttributePlanModifier` interface. This interface requires you to implement
three methods.

```go
// AttributePlanModifier represents a modifier for an attribute at plan time.
// An AttributePlanModifier can only modify the planned value for the attribute
// on which it is defined. For plan-time modifications that modify the values of
// several attributes at once, please instead use the ResourceWithModifyPlan
// interface by defining a ModifyPlan function on the resource.
type AttributePlanModifier interface {
	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this modifier is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description(context.Context) string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this modifier is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription(context.Context) string

	// Modify is called when the provider has an opportunity to modify
	// the plan: once during the plan phase when Terraform is determining
	// the diff that should be shown to the user for approval, and once
	// during the apply phase with any unknown values from configuration
	// filled in with their final values.
	//
	// The Modify function has access to the config, state, and plan for
	// both the attribute in question and the entire resource, but it can
	// only modify the value of the one attribute.
	//
	// Any returned errors will stop further execution of plan modifications
	// for this Attribute and any nested Attribute. Other Attribute at the same
	// or higher levels of the Schema will still execute any plan modifications
	// to ensure all warnings and errors across all root Attribute are
	// captured.
	//
	// Please see the documentation for ResourceWithModifyPlan#ModifyPlan
	// for further details.
	Modify(context.Context, ModifyAttributePlanRequest, *ModifyAttributePlanResponse)
}
```

Our plan modifier to check the hash looks like the following.

```go
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
    resp.AttributePlan = types.StringValue(newHash)
    resp.RequiresReplace = true
    resp.Diagnostics.AddWarning(
    	"Schema File Changed",
      fmt.Sprintf("The file %s has changed from when the resource was originally created, this forces the destruction of the resource. Old file hash: %s, New file hash: %s", schema, schemaHash, newHash),
    )
  }
}
```

Breaking this down, the following code allows you to get an attribute 
value from the Terraform plan and set it to a Go type.

```go
attrPath := path.Empty().AtName("schema")
schema := ""
req.Plan.GetAttribute(ctx, attrPath, &schema)
```

Then you check if the hash has been set; if you do not do this check, then 
you will return an incorrect error the first time a plan runs since there will
not be an existing hash to compare.

```go
if schemaHash == "" {
  return
}
```

Finally, compare the hash and return a warning if things have changed. There is
no need to return an error as it is legitimate that the schema file
could change. However, since the API does not accept mutations of schema,
Terraform needs to destroy the resource and re-create it.

```go
newHash, _ := calculateHashFromFile(schema)
if newHash != schemaHash {
  // set the new value and set the requires replace, Terraform will force the resource to be re-created
  resp.AttributePlan = types.StringValue(newHash)
  resp.RequiresReplace = true
  resp.Diagnostics.AddWarning(
  	"Schema File Changed",
    fmt.Sprintf("The file %s has changed from when the resource was originally created, this forces the destruction of the resource. Old file hash: %s, New file hash: %s", schema, schemaHash, newHash),
  )
}
```

Add the `schemaPlanModifier` in the previous code block to the bottom of your 
`schema_resouce.go` file.

You can then update the `schema_hash` attribute to add this modifier to the
exising collection of `tfsdk.AttributePlanModifier`s.

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

:::note
There is already an existing plan modifier `resource.RequiresReplace()`, this 
tells Terraform to replace the entire resource should the value of this attribute
change from an existing value in the state.
:::

## Defining Provider Resources
This completes the `Create` method before you move on to the `Read`, method
let's check that your code compiles. Before you do, there is one small change
that needs to be made in the `provider.go` file.

All resources you create for your provider must be registered so that Terraform
understands which resources the provider has.

This is done in the `Resources` method for the provider that is defined in
the file `provider.go`. The method returns a collection of available resources
it is important to remember that when you create a new resource type, you add it 
to this collection.

Update the Resources method in the provider so that it looks like the following.

```go
func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
	}
}
```

You can then check your code compiles by running `make build` in the terminal.

```shell
make build
```

```shell
go build -o bin/terraform-provider-example_v0.1.0
```

Assuming all went well, you should have no errors, and you can now move on
to implementing the `Read` method for the `schema` resource.