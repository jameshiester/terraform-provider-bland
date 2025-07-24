// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ datasource.DataSource = &KnowledgeBaseDataSource{}

type KnowledgeBaseDataSource struct {
	utils.TypeInfo
	KnowledgeBaseClient *KnowledgeBaseClient
}

func NewKnowledgeBaseDataSource() datasource.DataSource {
	return &KnowledgeBaseDataSource{
		TypeInfo: utils.TypeInfo{
			TypeName: "knowledge_base",
		},
	}
}

func (d *KnowledgeBaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *KnowledgeBaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.KnowledgeBaseClient = NewKnowledgeBaseClient(client.Api)
}

func (d *KnowledgeBaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve a knowledge base by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique knowledge base id",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the knowledge base",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the knowledge base",
				Computed:            true,
			},
			"text": schema.StringAttribute{
				MarkdownDescription: "Extracted text from the knowledge base",
				Computed:            true,
				Sensitive:           true,
			},
			"extracted_text": schema.StringAttribute{
				MarkdownDescription: "Extracted text from the knowledge base",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *KnowledgeBaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var config KnowledgeBaseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	read, err := d.KnowledgeBaseClient.ReadKnowledgeBase(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}

	model := ConvertFromKnowledgeBaseDtoToDataSource(*read)
	resp.State.Set(ctx, model)
}
