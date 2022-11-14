---
sidebar_position: 7
id: provider
title: Provider
---

The `provider.go`  type is the main entry point for your plugin, when
Terraform instantiates the plugin it calls the `New` method. This should
return an instance of your provider. Terraform then uses this to fetch the
resource and data types that it will use when processing the Terraform
configuration.

```go
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
```

Like the schema resource we can begin by modifying the `provider.go` file
that the template provides for you.

First replace all references for `Scaffolding` with `Minecraft`, then we can
start to look at how a provider can handle configuration.

## Creating the Provider Model

Like a resource Terraform will process the HCL configuration which can then
be deserialized into a data model. The data model for the provider is going
to be far simpler than the `schema` resource. You only need two attributes
`endpoint` that will alow the configuration of the API endpoint and the
`api_key` which is used to secure access.

To create the model, modify your `MinecraftProviderModel` so
that it is the same as the block below.

```go
type MinecraftProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}
```

Next you can define the schema for the provider.

## Defining the Schema

In the same way that a resource has attributes, a provider has the same.
Configuration of these attributes is done in in the `GetSchema` method.

Often you want to make an attribute optional so that it can be configured 
through alternate sources such as environment variables or configuration 
files that are read from a default location.

The Minecraft API client needs to be configured with an endpoint and an API key
let's add these as provider attributes so that they can be confiugured. Update
the `GetShema` block so that it is the same as the following example. 

```go
func (p *MinecraftProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}
```

## Configuring the Provider

Now that you have defined the configuration you can configure the provider,
this is where you will create the client that the schema resource used.

When instantiating a new provider Terraform will call the `Configure`
method allowing you to deserialize any attributes included in the provider
stanza and create anything needed by resources and data types.

Data is deserialized in a very similar way to how you did it when creating the 
resource.

```go
var data ScaffoldingProviderModel

resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

if resp.Diagnostics.HasError() {
	return
}
```

Next you need to fetch the configuration values from either the config
or optionally an environment variable.

```go
endpoint := ""
apiKey := ""

// set endpoint and apiKey from config
if !data.Endpoint.IsNull() {
	endpoint = data.Endpoint.ValueString()
}

if !data.APIKey.IsNull() {
	endpoint = data.APIKey.ValueString()
}

// override from environment variables if set
if ep := os.Getenv("MINECRAFT_ENDPOINT"); ep != "" {
	endpoint = ep
}

if apk := os.Getenv("MINECRAFT_APIKEY"); apk != "" {
	apiKey = apk
}
```

If either the `endpoint` or `apiKey` is empty then return
an error message to the user to let them know what they need to do.

```go
if endpoint == "" {
	resp.Diagnostics.AddError(
		"Configuration Error",
		"Unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_ENDPOINT'",
	)
	return
}

if apiKey == "" && data.APIKey.IsNull() {
	resp.Diagnostics.AddError(
		"Configuration Error",
		"unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_APIKEY'",
	)
	return
}
```

Finally you can instantate the Minecraft client with the `newClient` function
passing it the two values you have just obtained.

```go
// Example client configuration for data sources and resources
client := newClient(endpoint, apiKey)
resp.DataSourceData = client
resp.ResourceData = client	
```

The whole method looks like this, replace your `Configure` method with this
block.

```go
func (p *MinecraftProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MinecraftProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := ""
	apiKey := ""

	// set endpoint and apiKey from config
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.APIKey.IsNull() {
		endpoint = data.APIKey.ValueString()
	}

	// override from environment variables if set
	if ep := os.Getenv("MINECRAFT_ENDPOINT"); ep != "" {
		endpoint = ep
	}

	if apk := os.Getenv("MINECRAFT_APIKEY"); apk != "" {
		apiKey = apk
	}

	// if config is not set, return an error
	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Configuration Error",
			"Unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_ENDPOINT'",
		)
		return
	}

	if apiKey == "" && data.APIKey.IsNull() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			"unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_APIKEY'",
		)
		return
	}

	// Example client configuration for data sources and resources
	client := newClient(endpoint, apiKey)
	resp.DataSourceData = client
	resp.ResourceData = client
}
```

## Configuring the Provider Metadata

So that Terraform knows that a resouce called `minecraft_schema` is handled
by your plugin you need to configure the metadata to set the `TypeName` to the
prefix for your resources.

This is completed in the `Metadata` method, set the `resp.TypeName` to the
value of `minecraft` as shown in the following example.

```go
func (p *MinecraftProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "minecraft"
	resp.Version = p.version
}
```
	
## Compiling the Provider

Your completed provider code should look like the following example.

<details>
  <summary>provider.go</summary>

```go
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &MinecraftProvider{}
var _ provider.ProviderWithMetadata = &MinecraftProvider{}

// ScaffoldingProvider defines the provider implementation.
type MinecraftProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type MinecraftProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *MinecraftProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *MinecraftProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *MinecraftProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MinecraftProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := ""
	apiKey := ""

	// set endpoint and apiKey from config
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.APIKey.IsNull() {
		endpoint = data.APIKey.ValueString()
	}

	// override from environment variables if set
	if ep := os.Getenv("MINECRAFT_ENDPOINT"); ep != "" {
		endpoint = ep
	}

	if apk := os.Getenv("MINECRAFT_APIKEY"); apk != "" {
		apiKey = apk
	}

	// if config is not set, return an error
	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Configuration Error",
			"Unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_ENDPOINT'",
		)
		return
	}

	if apiKey == "" && data.APIKey.IsNull() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			"unable to set endpoint, please set either the endpoint property in the provider or the environment variable 'MINECRAFT_APIKEY'",
		)
		return
	}

	// Example client configuration for data sources and resources
	client := newClient(endpoint, apiKey)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MinecraftProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
	}
}

func (p *MinecraftProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MinecraftProvider{
			version: version,
		}
	}
}

```
</details>
	
Let's do one final check to make sure that everything is ok

```shell
make build
```

```shell
go build -o bin/terraform-provider-example_v0.1.0
```

Let's now test the provider

	
	
	
	
	
	
	