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
	resp.TypeName = "minecraft"
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
		apiKey = data.APIKey.ValueString()
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
		NewBlockResource,
	}
}

func (p *MinecraftProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBlockDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MinecraftProvider{
			version: version,
		}
	}
}
