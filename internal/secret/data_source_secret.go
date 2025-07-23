// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ datasource.DataSource = &SecretDataSource{}

type SecretDataSource struct {
	utils.TypeInfo
	SecretClient *SecretClient
}

func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{
		TypeInfo: utils.TypeInfo{
			TypeName: "secret",
		},
	}
}

func (d *SecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *SecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to retrieve a specific secret by `id`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the secret for which you want to retrieve detailed information.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the secret.",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the secret.",
				Computed:            true,
				Sensitive:           true,
			},
			"config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for refreshable secret.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						MarkdownDescription: "JSON body for the refresh request.",
						Computed:            true,
					},
					"method": schema.StringAttribute{
						MarkdownDescription: "HTTP method for the refresh request.",
						Computed:            true,
					},
					"refresh_interval": schema.Int32Attribute{
						MarkdownDescription: "Refresh interval for the refresh request.",
						Computed:            true,
					},
					"response": schema.StringAttribute{
						MarkdownDescription: "Value to extract from the refresh request response.",
						Computed:            true,
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "URL for the refresh request.",
						Computed:            true,
					},
					"headers": schema.MapAttribute{
						MarkdownDescription: "Headers for the refresh request.",
						ElementType:         types.StringType,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *SecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.SecretClient = newSecretClient(client.Api)
}

func (d *SecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var data struct {
		ID types.String `tfsdk:"id"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := d.SecretClient.ReadSecret(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading secret", err.Error())
		return
	}

	model := ConvertFromSecretDto(*secret)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
