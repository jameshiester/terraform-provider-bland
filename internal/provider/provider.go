// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/config"
	pathways "github.com/jameshiester/terraform-provider-bland/internal/conversational-pathways"
)

// Ensure BlandProvider satisfies various provider interfaces.
var _ provider.Provider = &BlandProvider{}

// BlandProvider defines the provider implementation.
type BlandProvider struct {
	Config  *config.ProviderConfig
	Api     *api.Client
	version string
}

// BlandProviderModel describes the provider data model.
type BlandProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *BlandProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bland"
	resp.Version = p.version
}

func (p *BlandProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key.  Can also be sourced from the `BLAND_API_KEY` environment variable",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *BlandProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BlandProviderModel
	// Check environment variables
	apiToken := os.Getenv("BLAND_API_KEY")
	baseUrl := os.Getenv("BLAND_BASE_URL")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.APIKey.ValueString() != "" {
		apiToken = data.APIKey.ValueString()
	}
	if baseUrl == "" {
		baseUrl = "api.bland.ai"
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Key Configuration",
			"While configuring the provider, the API key was not found in "+
				"the BLAND_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}
	p.Config.APIKey = apiToken
	p.Config.BaseURL = baseUrl
	p.Config.TerraformVersion = req.TerraformVersion
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	providerClient := api.ProviderClient{
		Config: p.Config,
		Api:    p.Api,
	}
	resp.DataSourceData = &providerClient
	resp.ResourceData = &providerClient
}

func (p *BlandProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return pathways.NewConversationalPathwayResource() },
	}
}

func (p *BlandProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return pathways.NewConversationalPathwayDataSource() },
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BlandProvider{
			version: version,
		}
	}
}

func NewBlandProvider(ctx context.Context, testModeEnabled ...bool) func() provider.Provider {
	providerConfig := config.ProviderConfig{
		TerraformVersion: "unknown",
	}

	if len(testModeEnabled) > 0 && testModeEnabled[0] {
		tflog.Warn(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
		providerConfig.TestMode = true
	}

	return func() provider.Provider {
		p := &BlandProvider{
			Config: &providerConfig,
			Api:    api.NewApiClientBase(&providerConfig, api.NewAuthBase(&providerConfig)),
		}
		return p
	}
}
