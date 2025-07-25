// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

// ConversationalPathwayDataSource defines the data source implementation.
type ConversationalPathwayDataSource struct {
	utils.TypeInfo
	ApplicationClient client
}

// ConversationalPathwayNodeModel describes the node model.
type ConversationalPathwayNodeModel struct {
	Type types.String                       `tfsdk:"type"`
	ID   types.String                       `tfsdk:"id"`
	Data ConversationalPathwayNodeDataModel `tfsdk:"data"`
}

type ConversationalPathwayNodeDataExtractVariableModel struct {
	Name                      types.String `tfsdk:"name"`
	Type                      types.String `tfsdk:"type"`
	Description               types.String `tfsdk:"description"`
	IncreaseSpellingPrecision types.Bool   `tfsdk:"increase_spelling_precision"`
}

type ConversationalPathwayNodeDataResponseDataModel struct {
	Data    types.String `tfsdk:"data"`
	Name    types.String `tfsdk:"name"`
	Context types.String `tfsdk:"context"`
}

type ConversationalPathwayNodeDataReponsePathwayConditionModel struct {
	Variable  types.String `tfsdk:"variable"`
	Condition types.String `tfsdk:"condition"`
	Value     types.String `tfsdk:"value"`
}

type ConversationalPathwayNodeDataReponsePathwayOutcomeModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"node_name"`
}

type ConversationalPathwayNodeDataModelOptionModel struct {
	Type                  types.String  `tfsdk:"model_type"`
	InterruptionThreshold types.String  `tfsdk:"interruption_threshold"`
	Temperature           types.Float32 `tfsdk:"temperature"`
	SkipUserResponse      types.Bool    `tfsdk:"skip_user_response"`
	BlockInterruptions    types.Bool    `tfsdk:"block_interruptions"`
}

type ConversationalPathwayNodeDataResponsePathwayModel struct {
	Condition ConversationalPathwayNodeDataReponsePathwayConditionModel `tfsdk:"condition"`
	Outcome   ConversationalPathwayNodeDataReponsePathwayOutcomeModel   `tfsdk:"outcome"`
}

// ConversationalPathwayNodeDataModel describes the node data model.
type ConversationalPathwayNodeDataModel struct {
	ExtractVars      []ConversationalPathwayNodeDataExtractVariableModel `tfsdk:"extract_vars"`
	GlobalLabel      types.String                                        `tfsdk:"global_label"`
	GlobalPrompt     types.String                                        `tfsdk:"global_prompt"`
	IsGlobal         types.Bool                                          `tfsdk:"is_global"`
	Method           types.String                                        `tfsdk:"method"`
	IsStart          types.Bool                                          `tfsdk:"is_start"`
	Name             types.String                                        `tfsdk:"name"`
	Prompt           types.String                                        `tfsdk:"prompt"`
	ResponseData     []ConversationalPathwayNodeDataResponseDataModel    `tfsdk:"response_data"`
	ResponsePathways []ConversationalPathwayNodeDataResponsePathwayModel `tfsdk:"response_pathways"`
	Text             types.String                                        `tfsdk:"text"`
	URL              types.String                                        `tfsdk:"url"`
	Condition        types.String                                        `tfsdk:"condition"`
	KnowledgeBase    types.String                                        `tfsdk:"kb"`
	KbTool           types.String                                        `tfsdk:"kb_tool"`
	TransferNumber   types.String                                        `tfsdk:"transfer_number"`
	ModelOptions     *ConversationalPathwayNodeDataModelOptionModel      `tfsdk:"model_options"`
	PathwayExamples  []ConversationalPathwayExampleModel                 `tfsdk:"pathway_examples"`
	Headers          []ConversationalPathwayHeaderModel                  `tfsdk:"headers"`
	Auth             *ConversationalPathwayAuthModel                     `tfsdk:"auth"`
	Body             types.String                                        `tfsdk:"body"`
	Routes           []ConversationalPathwayRouteModel                   `tfsdk:"routes"`
	FallbackNodeId   types.String                                        `tfsdk:"fallback_node_id"`
	TimeoutValue     types.Int64                                         `tfsdk:"timeout_value"`
	MaxRetries       types.Int64                                         `tfsdk:"max_retries"`
}

type ConversationalPathwayAuthModel struct {
	Type   types.String `tfsdk:"type"`
	Token  types.String `tfsdk:"token"`
	Encode types.Bool   `tfsdk:"encode"`
}

type ConversationalPathwayHeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type ConversationalPathwayExampleModel struct {
	ChosenPathway       types.String                             `tfsdk:"chosen_pathway"`
	ConversationHistory ConversationalPathwayExampleHistoryModel `tfsdk:"conversation_history"`
}

type ConversationalPathwayExampleHistoryModel struct {
	BasicHistory    types.String                               `tfsdk:"basic_history"`
	AdvancedHistory []ConversationalPathwayExampleMessageModel `tfsdk:"advanced_history"`
}

type ConversationalPathwayExampleMessageModel struct {
	Role    types.String `tfsdk:"role"`
	Content types.String `tfsdk:"content"`
}

type ConversationalPathwayRouteConditionModel struct {
	Field    types.String `tfsdk:"field"`
	Value    types.String `tfsdk:"value"`
	IsGroup  types.Bool   `tfsdk:"is_group"`
	Operator types.String `tfsdk:"operator"`
}

type ConversationalPathwayRouteModel struct {
	Conditions   []ConversationalPathwayRouteConditionModel `tfsdk:"conditions"`
	TargetNodeId types.String                               `tfsdk:"target_node_id"`
}

// ConversationalPathwayDataSourceModel describes the data source data model.
type ConversationalPathwayDataSourceModel struct {
	Name         types.String                       `tfsdk:"name"`
	ID           types.String                       `tfsdk:"id"`
	Description  types.String                       `tfsdk:"description"`
	Nodes        []ConversationalPathwayNodeModel   `tfsdk:"nodes"`
	Edges        []ConversationalPathwayEdgeModel   `tfsdk:"edges"`
	GlobalConfig *ConversationalPathwayGlobalConfig `tfsdk:"global_config"`
}

type ConversationalPathwayGlobalConfig struct {
	GlobalPrompt types.String `tfsdk:"global_prompt"`
}

type ConversationalPathwayEdgeModel struct {
	ID     types.String                       `tfsdk:"id"`
	Source types.String                       `tfsdk:"source"`
	Target types.String                       `tfsdk:"target"`
	Type   types.String                       `tfsdk:"type"`
	Data   ConversationalPathwayEdgeDataModel `tfsdk:"data"`
}

type ConversationalPathwayEdgeConditionModel struct {
	Field    types.String `tfsdk:"field"`
	Value    types.String `tfsdk:"value"`
	IsGroup  types.Bool   `tfsdk:"is_group"`
	Operator types.String `tfsdk:"operator"`
}

type ConversationalPathwayEdgeDataModel struct {
	Label         types.String                               `tfsdk:"label"`
	IsHighlighted types.Bool                                 `tfsdk:"is_highlighted"`
	Description   types.String                               `tfsdk:"description"`
	AlwaysPick    types.Bool                                 `tfsdk:"always_pick"`
	Conditions    *[]ConversationalPathwayEdgeConditionModel `tfsdk:"conditions"`
}
