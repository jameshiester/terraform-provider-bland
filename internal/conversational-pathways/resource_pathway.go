// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

var _ resource.Resource = &ConversationalPathwayResource{}
var _ resource.ResourceWithImportState = &ConversationalPathwayResource{}

type ConversationalPathwayResource struct {
	utils.TypeInfo
	PathwayClient client
}

func NewConversationalPathwayResource() resource.Resource {
	return &ConversationalPathwayResource{
		TypeInfo: utils.TypeInfo{
			TypeName: "conversational_pathway",
		},
	}
}

func (r *ConversationalPathwayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ConversationalPathwayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [Conversational Pathway](https://docs.bland.ai/tutorials/pathways).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique pathway id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the pathway",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the pathway",
				Required:            true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "Data about all the nodes in the pathway.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier for the node.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the node.",
							Required:            true,
						},
						"data": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the node.",
									Required:            true,
								},
								"text": schema.StringAttribute{
									MarkdownDescription: "Text for the node.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("prompt")),
									},
								},
								"global_prompt": schema.StringAttribute{
									MarkdownDescription: "Prompt for a global node.",
									Optional:            true,
								},
								"global_label": schema.StringAttribute{
									MarkdownDescription: "Label for a global node.",
									Optional:            true,
								},
								"method": schema.StringAttribute{
									MarkdownDescription: "Method for the node.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|TRACE|CONNECT)$`),
											"must be a valid HTTP method in uppercase (e.g., GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT)",
										),
									},
								},
								"is_start": schema.BoolAttribute{
									MarkdownDescription: "Defines if this is the start node of the pathway.",
									Optional:            true,
								},
								"is_global": schema.BoolAttribute{
									MarkdownDescription: "Defines if this is a global node.",
									Optional:            true,
								},
								"prompt": schema.StringAttribute{
									MarkdownDescription: "Prompt for a knowledge base node.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("text")),
									},
								},
								"url": schema.StringAttribute{
									MarkdownDescription: "URL for the node.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^https?://[^\s]+$`),
											"must be a valid URL starting with http:// or https:// (e.g., http://example.com)",
										),
									},
								},
								"condition": schema.StringAttribute{
									MarkdownDescription: "Condition for the node.",
									Optional:            true,
								},
								"kb": schema.StringAttribute{
									MarkdownDescription: "Knowledge base for the node.",
									Optional:            true,
								},
								"transfer_number": schema.StringAttribute{
									MarkdownDescription: "Transfer number for the node.",
									Optional:            true,
								},
								"extract_vars": schema.ListNestedAttribute{
									MarkdownDescription: "Variables to extract from the node.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "Name of the variable.",
												Required:            true,
											},
											"type": schema.StringAttribute{
												MarkdownDescription: "Type of the variable.",
												Required:            true,
											},
											"description": schema.StringAttribute{
												MarkdownDescription: "Description of the variable.",
												Required:            true,
											},
										},
									},
								},
								"response_data": schema.ListNestedAttribute{
									MarkdownDescription: "Response data for the node.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"data": schema.StringAttribute{
												MarkdownDescription: "Data value.",
												Required:            true,
											},
											"name": schema.StringAttribute{
												MarkdownDescription: "Name of the response data.",
												Required:            true,
											},
											"context": schema.StringAttribute{
												MarkdownDescription: "Context for the response data.",
												Required:            true,
											},
										},
									},
								},
								"response_pathways": schema.ListNestedAttribute{
									MarkdownDescription: "Response pathways for the node.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"condition": schema.SingleNestedAttribute{
												Required: true,
												Attributes: map[string]schema.Attribute{
													"variable": schema.StringAttribute{
														MarkdownDescription: "Condition variable.",
														Required:            true,
													},
													"condition": schema.StringAttribute{
														MarkdownDescription: "Condition operator.",
														Required:            true,
													},
													"value": schema.StringAttribute{
														MarkdownDescription: "Condition value.",
														Required:            true,
													},
												},
											},
											"outcome": schema.SingleNestedAttribute{
												Required: true,
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														MarkdownDescription: "Outcome node id.",
														Required:            true,
													},
													"node_name": schema.StringAttribute{
														MarkdownDescription: "Outcome node name.",
														Required:            true,
													},
												},
											},
										},
									},
								},
								"model_options": schema.SingleNestedAttribute{
									MarkdownDescription: "Model options for the node.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"model_type": schema.StringAttribute{
											MarkdownDescription: "Type of the model.",
											Required:            true,
										},
										"interruption_threshold": schema.StringAttribute{
											MarkdownDescription: "Interruption threshold for the model.",
											Optional:            true,
										},
										"temperature": schema.Float32Attribute{
											MarkdownDescription: "Temperature setting for the model.",
											Optional:            true,
										},
										"skip_user_response": schema.BoolAttribute{
											MarkdownDescription: "Whether to skip user response.",
											Optional:            true,
										},
										"block_interruptions": schema.BoolAttribute{
											MarkdownDescription: "Whether to block interruptions.",
											Optional:            true,
										},
									},
								},
								"pathway_examples": schema.ListNestedAttribute{
									MarkdownDescription: "Example conversations and chosen pathways for this node.",
									Optional:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"chosen_pathway": schema.StringAttribute{
												MarkdownDescription: "The chosen pathway for the example.",
												Required:            true,
											},
											"conversation_history": schema.SingleNestedAttribute{
												MarkdownDescription: "The conversation history for the example.",
												Required:            true,
												Attributes: map[string]schema.Attribute{
													"basic_history": schema.StringAttribute{
														MarkdownDescription: "Conversation history as a string.",
														Optional:            true,
														Validators: []validator.String{
															stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("advanced_history")),
														},
													},
													"advanced_history": schema.ListNestedAttribute{
														MarkdownDescription: "Conversation history as a list of messages.",
														Optional:            true,
														Validators: []validator.List{
															listvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic_history")),
														},
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"role": schema.StringAttribute{
																	MarkdownDescription: "Role of the message (user or assistant).",
																	Required:            true,
																},
																"content": schema.StringAttribute{
																	MarkdownDescription: "Content of the message.",
																	Required:            true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"edges": schema.ListNestedAttribute{
				MarkdownDescription: "Data about all the edges in the pathway.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier for the edge.",
							Required:            true,
						},
						"source": schema.StringAttribute{
							MarkdownDescription: "Source node ID for the edge.",
							Required:            true,
						},
						"target": schema.StringAttribute{
							MarkdownDescription: "Target node ID for the edge.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the edge.",
							Required:            true,
						},
						"data": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"label": schema.StringAttribute{
									MarkdownDescription: "Label for the edge.",
									Required:            true,
								},
								"is_highlighted": schema.BoolAttribute{
									MarkdownDescription: "Whether the edge is highlighted.",
									Optional:            true,
									Computed:            true,
									Default:             booldefault.StaticBool(false),
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description of the edge.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"global_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Global configuration for the pathway.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"global_prompt": schema.StringAttribute{
						MarkdownDescription: "Global prompt for the pathway.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *ConversationalPathwayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	_, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
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
	r.PathwayClient = newPathwayClient(client.Api)
}

func (r *ConversationalPathwayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan ConversationalPathwayDataSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dto := ConvertFromPathwayModel(plan)

	modelToCreate := createPathwayDto{
		Name:        dto.Name,
		Description: dto.Description,
		Nodes:       dto.Nodes,
		Edges:       dto.Edges,
	}

	connection, err := r.PathwayClient.CreatePathway(ctx, modelToCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create pathway", err.Error())
		return
	}

	responseModel, err := ConvertFromPathwayDto(*connection)
	if err != nil {
		resp.Diagnostics.AddError("Error occurred when parsing create pathway response", err.Error())
		return
	}
	plan.ID = responseModel.ID
	plan.Description = responseModel.Description
	plan.Name = responseModel.Name
	plan.Nodes = responseModel.Nodes
	plan.Edges = responseModel.Edges
	plan.GlobalConfig = responseModel.GlobalConfig
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConversationalPathwayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pathway, err := r.PathwayClient.GetPathway(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	model, err := ConvertFromPathwayDto(*pathway)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error when converting %s", r.FullTypeName()), err.Error())
		return
	}
	state.Name = model.Name
	state.Description = model.Description
	state.Nodes = model.Nodes
	state.Edges = model.Edges
	state.GlobalConfig = model.GlobalConfig
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ConversationalPathwayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ConversationalPathwayDataSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dto := ConvertFromPathwayModel(plan)

	updateParams := updatePathwayDto{
		Name:        dto.Name,
		Description: dto.Description,
		Nodes:       dto.Nodes,
		Edges:       dto.Edges,
	}

	updateReponse, err := r.PathwayClient.UpdatePathway(ctx, plan.ID.ValueString(), updateParams)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}

	modelState, err := ConvertFromPathwayDto(*updateReponse)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error when converting updated %s", r.FullTypeName()), err.Error())
		return
	}
	plan.Name = modelState.Name
	plan.Description = modelState.Description
	plan.Nodes = modelState.Nodes
	plan.Edges = modelState.Edges
	plan.GlobalConfig = modelState.GlobalConfig
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConversationalPathwayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ConversationalPathwayDataSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.PathwayClient.DeletePathway(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *ConversationalPathwayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := utils.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Custom validator to ensure 'text' and 'prompt' are mutually exclusive
