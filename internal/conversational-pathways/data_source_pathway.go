// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConversationalPathwayDataSource{}

func NewConversationalPathwayDataSource() datasource.DataSource {
	return &ConversationalPathwayDataSource{
		TypeInfo: utils.TypeInfo{
			TypeName: "conversational_pathway",
		},
	}
}

func (r *ConversationalPathwayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *ConversationalPathwayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source to retrieve a specific conversational pathway by `id`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the conversational pathway for which you want to retrieve detailed information.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the conversational pathway.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the conversational pathway.",
				Computed:            true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "Data about all the nodes in the pathway.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier for the node.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the node.",
							Computed:            true,
						},
						"data": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the node.",
									Computed:            true,
								},
								"text": schema.StringAttribute{
									MarkdownDescription: "Text for the node.",
									Computed:            true,
								},
								"global_prompt": schema.StringAttribute{
									MarkdownDescription: "Prompt for a global node.",
									Computed:            true,
								},
								"prompt": schema.StringAttribute{
									MarkdownDescription: "Prompt for a knowledge base node.",
									Computed:            true,
								},
								"is_start": schema.BoolAttribute{
									MarkdownDescription: "Defines if this is the start node of the pathway.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *ConversationalPathwayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	d.ApplicationClient = newPathwayClient(client.Api)
}

func (d *ConversationalPathwayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONVERSATIONAL PATHWAYS START: %s", d.FullTypeName()))
	if state.ID.ValueString() == "" {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), "ID is missing from state")
		return
	}

	pathway, err := d.ApplicationClient.GetPathway(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	model, err := ConvertFromPathwayDto(*pathway)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error when parsing %s", d.FullTypeName()), err.Error())
		return
	}
	state.Name = model.Name
	state.Description = model.Description
	state.Nodes = model.Nodes
	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONVERSATIONAL PATHWAYS END: %s", d.FullTypeName()))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
