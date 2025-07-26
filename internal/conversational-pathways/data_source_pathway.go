// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
								"global_label": schema.StringAttribute{
									MarkdownDescription: "Label for a global node.",
									Computed:            true,
								},
								"method": schema.StringAttribute{
									MarkdownDescription: "Method for the node.",
									Computed:            true,
								},
								"is_start": schema.BoolAttribute{
									MarkdownDescription: "Defines if this is the start node of the pathway.",
									Computed:            true,
								},
								"is_global": schema.BoolAttribute{
									MarkdownDescription: "Defines if this is a global node.",
									Computed:            true,
								},
								"prompt": schema.StringAttribute{
									MarkdownDescription: "Prompt for a knowledge base node.",
									Computed:            true,
								},
								"url": schema.StringAttribute{
									MarkdownDescription: "URL for the node.",
									Computed:            true,
								},
								"condition": schema.StringAttribute{
									MarkdownDescription: "Condition for the node.",
									Computed:            true,
								},
								"kb": schema.StringAttribute{
									MarkdownDescription: "Knowledge base for the node.",
									Computed:            true,
								},
								"kb_tool": schema.StringAttribute{
									MarkdownDescription: "Knowledge base tool for the node.",
									Computed:            true,
								},
								"transfer_number": schema.StringAttribute{
									MarkdownDescription: "Transfer number for the node.",
									Computed:            true,
								},
								"extract_vars": schema.ListNestedAttribute{
									MarkdownDescription: "Variables to extract from the node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "Name of the variable.",
												Computed:            true,
											},
											"type": schema.StringAttribute{
												MarkdownDescription: "Type of the variable.",
												Computed:            true,
											},
											"description": schema.StringAttribute{
												MarkdownDescription: "Description of the variable.",
												Computed:            true,
											},
											"increase_spelling_precision": schema.BoolAttribute{
												MarkdownDescription: "Indicates if model uses increased spelling precision",
												Computed:            true,
											},
										},
									},
								},
								"response_data": schema.ListNestedAttribute{
									MarkdownDescription: "Response data for the node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"data": schema.StringAttribute{
												MarkdownDescription: "Data value.",
												Computed:            true,
											},
											"name": schema.StringAttribute{
												MarkdownDescription: "Name of the response data.",
												Computed:            true,
											},
											"context": schema.StringAttribute{
												MarkdownDescription: "Context for the response data.",
												Computed:            true,
											},
										},
									},
								},
								"response_pathways": schema.ListNestedAttribute{
									MarkdownDescription: "Response pathways for the node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"condition": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"variable": schema.StringAttribute{
														MarkdownDescription: "Condition variable.",
														Computed:            true,
													},
													"condition": schema.StringAttribute{
														MarkdownDescription: "Condition operator.",
														Computed:            true,
													},
													"value": schema.StringAttribute{
														MarkdownDescription: "Condition value.",
														Computed:            true,
													},
												},
											},
											"outcome": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														MarkdownDescription: "Outcome node id.",
														Computed:            true,
													},
													"node_name": schema.StringAttribute{
														MarkdownDescription: "Outcome node name.",
														Computed:            true,
													},
												},
											},
										},
									},
								},
								"model_options": schema.SingleNestedAttribute{
									MarkdownDescription: "Model options for the node.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"model_type": schema.StringAttribute{
											MarkdownDescription: "Type of the model.",
											Computed:            true,
										},
										"interruption_threshold": schema.StringAttribute{
											MarkdownDescription: "Interruption threshold for the model.",
											Computed:            true,
										},
										"temperature": schema.Float32Attribute{
											MarkdownDescription: "Temperature setting for the model.",
											Computed:            true,
										},
										"skip_user_response": schema.BoolAttribute{
											MarkdownDescription: "Whether to skip user response.",
											Computed:            true,
										},
										"block_interruptions": schema.BoolAttribute{
											MarkdownDescription: "Whether to block interruptions.",
											Computed:            true,
										},
									},
								},
								"pathway_examples": schema.ListNestedAttribute{
									MarkdownDescription: "Example conversations and chosen pathways for this node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"chosen_pathway": schema.StringAttribute{
												MarkdownDescription: "The chosen pathway for the example.",
												Computed:            true,
											},
											"conversation_history": schema.SingleNestedAttribute{
												MarkdownDescription: "The conversation history for the example.",
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"basic_history": schema.StringAttribute{
														MarkdownDescription: "Conversation history as a string.",
														Computed:            true,
														Validators: []validator.String{
															stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("advanced_history")),
														},
													},
													"advanced_history": schema.ListNestedAttribute{
														MarkdownDescription: "Conversation history as a list of messages.",
														Computed:            true,
														Validators: []validator.List{
															listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic_history")),
														},
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"role": schema.StringAttribute{
																	MarkdownDescription: "Role of the message (user or assistant).",
																	Computed:            true,
																},
																"content": schema.StringAttribute{
																	MarkdownDescription: "Content of the message.",
																	Computed:            true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"auth": schema.SingleNestedAttribute{
									MarkdownDescription: "Authentication for the node.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "Auth type (e.g., Bearer).",
											Computed:            true,
										},
										"token": schema.StringAttribute{
											MarkdownDescription: "Auth token.",
											Computed:            true,
										},
										"encode": schema.BoolAttribute{
											MarkdownDescription: "Whether to encode the token.",
											Computed:            true,
										},
									},
								},
								"body": schema.StringAttribute{
									MarkdownDescription: "Body for the node.",
									Computed:            true,
								},
								"headers": schema.ListNestedAttribute{
									MarkdownDescription: "Headers for the node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "Header name.",
												Computed:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "Header value.",
												Computed:            true,
											},
										},
									},
								},
								"routes": schema.ListNestedAttribute{
									MarkdownDescription: "Routes for the node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"conditions": schema.ListNestedAttribute{
												MarkdownDescription: "Conditions for the route.",
												Computed:            true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"field": schema.StringAttribute{
															MarkdownDescription: "Field name.",
															Computed:            true,
														},
														"value": schema.StringAttribute{
															MarkdownDescription: "Field value.",
															Computed:            true,
														},
														"is_group": schema.BoolAttribute{
															MarkdownDescription: "Whether this is a group condition.",
															Computed:            true,
														},
														"operator": schema.StringAttribute{
															MarkdownDescription: "Condition operator.",
															Computed:            true,
														},
													},
												},
											},
											"target_node_id": schema.StringAttribute{
												MarkdownDescription: "Target node ID.",
												Computed:            true,
											},
										},
									},
								},
								"fallback_node_id": schema.StringAttribute{
									MarkdownDescription: "Fallback node ID.",
									Computed:            true,
								},
								"timeout_value": schema.Int64Attribute{
									MarkdownDescription: "Timeout value for the node.",
									Computed:            true,
								},
								"max_retries": schema.Int64Attribute{
									MarkdownDescription: "Maximum number of retries for the node.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"edges": schema.ListNestedAttribute{
				MarkdownDescription: "Data about all the edges in the pathway.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier for the edge.",
							Computed:            true,
						},
						"source": schema.StringAttribute{
							MarkdownDescription: "Source node ID for the edge.",
							Computed:            true,
						},
						"target": schema.StringAttribute{
							MarkdownDescription: "Target node ID for the edge.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the edge.",
							Computed:            true,
						},
						"data": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"label": schema.StringAttribute{
									MarkdownDescription: "Label for the edge.",
									Computed:            true,
								},
								"is_highlighted": schema.BoolAttribute{
									MarkdownDescription: "Whether the edge is highlighted.",
									Computed:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description of the edge.",
									Computed:            true,
								},
								"always_pick": schema.BoolAttribute{
									MarkdownDescription: "Whether this edge should always be picked.",
									Computed:            true,
								},
								"conditions": schema.ListNestedAttribute{
									MarkdownDescription: "Conditions for the edge.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"field": schema.StringAttribute{
												MarkdownDescription: "Field name.",
												Computed:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "Field value.",
												Computed:            true,
											},
											"is_group": schema.BoolAttribute{
												MarkdownDescription: "Whether this is a group condition.",
												Computed:            true,
											},
											"operator": schema.StringAttribute{
												MarkdownDescription: "Condition operator.",
												Computed:            true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"global_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Global configuration for the pathway.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"global_prompt": schema.StringAttribute{
						MarkdownDescription: "Global prompt for the pathway.",
						Computed:            true,
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
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s %s", d.FullTypeName(), state.Name.ValueString()), "ID is missing from state")
		return
	}
	state.Name = types.StringValue(state.Name.ValueString())
	state.Description = types.StringValue(state.Description.ValueString())

	pathway, err := d.ApplicationClient.GetPathway(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	model, err := ConvertFromPathwayDto(*pathway)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error when converting %s", d.FullTypeName()), err.Error())
		return
	}

	state.Name = types.StringValue(model.Name.ValueString())
	state.Description = types.StringValue(model.Description.ValueString())
	state.Nodes = model.Nodes
	state.Edges = model.Edges
	state.GlobalConfig = model.GlobalConfig
	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONVERSATIONAL PATHWAYS END: %s", d.FullTypeName()))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
