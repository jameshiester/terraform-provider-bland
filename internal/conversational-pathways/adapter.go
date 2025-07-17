// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExtractVars      *[][]string                       `json:"extractVars,omitempty"`
// GlobalLabel      *string                           `json:"globalLabel,omitempty"`
// GlobalPrompt     *string                           `json:"global_prompt,omitempty"`
// IsStart          *bool                             `json:"isStart,omitempty"`
// Method           *string                           `json:"method,omitempty"`
// Name             string                            `json:"name"`
// Prompt           *string                           `json:"prompt,omitempty"`
// ResponseData     *[]pathwayNodeDataResponseDataDto `json:"responseData,omitempty"`
// ResponsePathways *[]interface{}                    `json:"responsePathways,omitempty"`
// Text             *string                           `json:"text,omitempty"`
// URL              *string                           `json:"url,omitempty"`

func ConvertFromPathwayNodeDataExtractVars(vals []string) (ConversationalPathwayNodeDataExtractVariableModel, error) {
	model := ConversationalPathwayNodeDataExtractVariableModel{}
	if len(vals) != 3 {
		return model, fmt.Errorf("failed to get pathway: %d", len(vals))
	}
	model.Name = types.StringValue(vals[0])
	model.Type = types.StringValue(vals[1])
	model.Description = types.StringValue(vals[2])
	return model, nil
}

func ConvertFromPathwayNodeDataResponseData(data pathwayNodeDataResponseDataDto) ConversationalPathwayNodeDataResponseDataModel {
	return ConversationalPathwayNodeDataResponseDataModel{
		Name:    types.StringValue(data.Name),
		Data:    types.StringValue(data.Data),
		Context: types.StringValue(data.Context),
	}
}

func ConvertFromPathwayNodeDataResponsePathway(
	item []interface{},
) (*ConversationalPathwayNodeDataResponsePathwayModel, error) {
	if len(item) != 4 {
		return nil, fmt.Errorf("expected 4 elements got %d", len(item))
	}

	// Extract and type assert the first three elements
	variable, ok1 := item[0].(string)
	condition, ok2 := item[1].(string)
	value, ok3 := item[2].(string)
	outcomeMap, ok4 := item[3].(map[string]interface{})
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, fmt.Errorf("invalid types")
	}

	// Extract outcome fields
	id, okId := outcomeMap["id"].(string)
	name, okName := outcomeMap["name"].(string)
	if !okId || !okName {
		return nil, fmt.Errorf("invalid outcome object")
	}

	return &ConversationalPathwayNodeDataResponsePathwayModel{
		Condition: ConversationalPathwayNodeDataReponsePathwayConditionModel{
			Variable:  types.StringValue(variable),
			Condition: types.StringValue(condition),
			Value:     types.StringValue(value),
		},
		Outcome: ConversationalPathwayNodeDataReponsePathwayOutcomeModel{
			ID:   types.StringValue(id),
			Name: types.StringValue(name),
		},
	}, nil
}

func ConvertFromPathwayNodeDataDto(data *pathwayNodeDataDto) (*ConversationalPathwayNodeDataModel, error) {
	if data == nil {
		return nil, nil
	}
	model := ConversationalPathwayNodeDataModel{
		GlobalPrompt:   types.StringPointerValue(data.GlobalPrompt),
		GlobalLabel:    types.StringPointerValue(data.GlobalLabel),
		Method:         types.StringPointerValue(data.Method),
		IsStart:        types.BoolPointerValue(data.IsStart),
		IsGlobal:       types.BoolPointerValue(data.IsGlobal),
		Name:           types.StringValue(data.Name),
		Prompt:         types.StringPointerValue(data.Prompt),
		Text:           types.StringPointerValue(data.Text),
		URL:            types.StringPointerValue(data.URL),
		Condition:      types.StringPointerValue(data.Condition),
		KnowledgeBase:  types.StringPointerValue(data.KnowledgeBase),
		TransferNumber: types.StringPointerValue(data.TransferNumber),
	}

	if data.ModelOptions != nil {
		model.ModelOptions = &ConversationalPathwayNodeDataModelOptionModel{
			Type:                  types.StringValue(data.ModelOptions.Type),
			InterruptionThreshold: types.StringPointerValue(data.ModelOptions.InterruptionThreshold),
			Temperature:           types.Float32PointerValue(data.ModelOptions.Temperature),
			SkipUserResponse:      types.BoolPointerValue(data.ModelOptions.SkipUserResponse),
			BlockInterruptions:    types.BoolPointerValue(data.ModelOptions.BlockInterruptions),
		}
	}

	if data.ExtractVars != nil {
		for _, variable := range *data.ExtractVars {
			varModel, err := ConvertFromPathwayNodeDataExtractVars(variable)
			if err != nil {
				return nil, err
			}
			model.ExtractVars = append(model.ExtractVars, varModel)
		}
	}
	if data.ResponseData != nil {
		for _, responseData := range *data.ResponseData {
			responseModel := ConvertFromPathwayNodeDataResponseData(responseData)
			model.ResponseData = append(model.ResponseData, responseModel)
		}
	}
	if data.ResponsePathways != nil {
		for _, responsePathwayData := range *data.ResponsePathways {
			responsePathway, err := ConvertFromPathwayNodeDataResponsePathway(responsePathwayData)
			if err != nil {
				return nil, err
			}
			model.ResponsePathways = append(model.ResponsePathways, *responsePathway)
		}
	}
	return &model, nil
}

func ConvertFromPathwayNodeDataModel(data ConversationalPathwayNodeDataModel) *pathwayNodeDataDto {
	var extractVars *[][]string
	if len(data.ExtractVars) > 0 {
		tmp := make([][]string, 0, len(data.ExtractVars))
		for _, v := range data.ExtractVars {
			tmp = append(tmp, []string{
				v.Name.ValueString(),
				v.Type.ValueString(),
				v.Description.ValueString(),
			})
		}
		extractVars = &tmp
	} else {
		extractVars = nil
	}

	var responseData *[]pathwayNodeDataResponseDataDto
	if len(data.ResponseData) > 0 {
		tmp := make([]pathwayNodeDataResponseDataDto, 0, len(data.ResponseData))
		for _, v := range data.ResponseData {
			tmp = append(tmp, pathwayNodeDataResponseDataDto{
				Name:    v.Name.ValueString(),
				Data:    v.Data.ValueString(),
				Context: v.Context.ValueString(),
			})
		}
		responseData = &tmp
	} else {
		responseData = nil
	}

	var responsePathways *[][]interface{}
	if len(data.ResponsePathways) > 0 {
		tmp := make([][]interface{}, 0, len(data.ResponsePathways))
		for _, v := range data.ResponsePathways {
			row := []interface{}{
				v.Condition.Variable.ValueString(),
				v.Condition.Condition.ValueString(),
				v.Condition.Value.ValueString(),
				map[string]interface{}{
					"id":   v.Outcome.ID.ValueString(),
					"name": v.Outcome.Name.ValueString(),
				},
			}
			tmp = append(tmp, row)
		}
		responsePathways = &tmp
	} else {
		responsePathways = nil
	}

	var modelOptions *modelOptionDto
	if data.ModelOptions != nil {
		modelOptions = &modelOptionDto{
			Type:                  data.ModelOptions.Type.ValueString(),
			InterruptionThreshold: data.ModelOptions.InterruptionThreshold.ValueStringPointer(),
			Temperature:           data.ModelOptions.Temperature.ValueFloat32Pointer(),
			SkipUserResponse:      data.ModelOptions.SkipUserResponse.ValueBoolPointer(),
			BlockInterruptions:    data.ModelOptions.BlockInterruptions.ValueBoolPointer(),
		}
	} else {
		modelOptions = nil
	}

	return &pathwayNodeDataDto{
		Name:             data.Name.ValueString(),
		GlobalPrompt:     data.GlobalPrompt.ValueStringPointer(),
		Prompt:           data.Prompt.ValueStringPointer(),
		Text:             data.Text.ValueStringPointer(),
		IsStart:          data.IsStart.ValueBoolPointer(),
		IsGlobal:         data.IsGlobal.ValueBoolPointer(),
		Method:           data.Method.ValueStringPointer(),
		URL:              data.URL.ValueStringPointer(),
		GlobalLabel:      data.GlobalLabel.ValueStringPointer(),
		Condition:        data.Condition.ValueStringPointer(),
		KnowledgeBase:    data.KnowledgeBase.ValueStringPointer(),
		TransferNumber:   data.TransferNumber.ValueStringPointer(),
		ExtractVars:      extractVars,
		ResponseData:     responseData,
		ResponsePathways: responsePathways,
		ModelOptions:     modelOptions,
	}
}

func ConvertFromPathwayNodeDto(node pathwayNodeDto) (*ConversationalPathwayNodeModel, error) {
	data, err := ConvertFromPathwayNodeDataDto(node.Data)
	if err != nil {
		return nil, err
	}
	return &ConversationalPathwayNodeModel{
		ID:   types.StringPointerValue(node.ID),
		Type: types.StringPointerValue(node.Type),
		Data: *data,
	}, nil
}

func ConvertFromPathwayGlobalConfigNodeDto(node pathwayNodeDto) *ConversationalPathwayGlobalConfig {
	return &ConversationalPathwayGlobalConfig{
		GlobalPrompt: types.StringValue(node.GlobalConfig.GlobalPrompt),
	}
}

func ConvertFromPathwayNodeModel(node ConversationalPathwayNodeModel) pathwayNodeDto {
	return pathwayNodeDto{
		ID:   node.ID.ValueStringPointer(),
		Type: node.Type.ValueStringPointer(),
		Data: ConvertFromPathwayNodeDataModel(node.Data),
	}
}

func ConvertFromPathwayEdgeDto(edge pathwayEdgeDto) ConversationalPathwayEdgeModel {
	return ConversationalPathwayEdgeModel{
		ID:     types.StringValue(edge.ID),
		Source: types.StringValue(edge.Source),
		Target: types.StringValue(edge.Target),
		Type:   types.StringValue(edge.Type),
		Data: ConversationalPathwayEdgeDataModel{
			Label:         types.StringValue(edge.Data.Label),
			IsHighlighted: types.BoolValue(edge.Data.IsHighlighted),
			Description:   types.StringValue(edge.Data.Description),
		},
	}
}

func ConvertFromPathwayEdgeModel(edge ConversationalPathwayEdgeModel) pathwayEdgeDto {
	return pathwayEdgeDto{
		ID:     edge.ID.ValueString(),
		Source: edge.Source.ValueString(),
		Target: edge.Target.ValueString(),
		Type:   edge.Type.ValueString(),
		Data: pathwayEdgeDataDto{
			Label:         edge.Data.Label.ValueString(),
			IsHighlighted: edge.Data.IsHighlighted.ValueBool(),
			Description:   edge.Data.Description.ValueString(),
		},
	}
}

func ConvertFromPathwayGlobalConfigDto(config *pathwayGlobalConfigDto) ConversationalPathwayGlobalConfig {
	if config == nil {
		return ConversationalPathwayGlobalConfig{}
	}
	return ConversationalPathwayGlobalConfig{
		GlobalPrompt: types.StringValue(config.GlobalPrompt),
	}
}

func ConvertFromPathwayGlobalConfigModel(config ConversationalPathwayGlobalConfig) pathwayNodeDto {
	return pathwayNodeDto{
		GlobalConfig: &pathwayGlobalConfigDto{
			GlobalPrompt: config.GlobalPrompt.ValueString(),
		},
	}
}

func ConvertFromPathwayDto(pathway pathwayDto) (*ConversationalPathwayDataSourceModel, error) {

	path := ConversationalPathwayDataSourceModel{
		ID:          types.StringValue(pathway.ID),
		Name:        types.StringValue(pathway.Name),
		Description: types.StringValue(pathway.Description),
	}
	for _, node := range pathway.Nodes {
		if node.GlobalConfig == nil {
			nodeModel, err := ConvertFromPathwayNodeDto(node)
			if err != nil {
				return nil, err
			}
			path.Nodes = append(path.Nodes, *nodeModel)
		} else {
			path.GlobalConfig = ConvertFromPathwayGlobalConfigNodeDto(node)
		}
	}
	for _, edge := range pathway.Edges {
		edgeModel := ConvertFromPathwayEdgeDto(edge)
		path.Edges = append(path.Edges, edgeModel)
	}

	return &path, nil
}

func ConvertFromPathwayModel(pathway ConversationalPathwayDataSourceModel) pathwayDto {

	path := pathwayDto{
		ID:          pathway.ID.ValueString(),
		Name:        pathway.Name.ValueString(),
		Description: pathway.Description.ValueString(),
	}
	for _, node := range pathway.Nodes {
		nodeModel := ConvertFromPathwayNodeModel(node)
		path.Nodes = append(path.Nodes, nodeModel)
	}
	if pathway.GlobalConfig != nil {
		globalNode := ConvertFromPathwayGlobalConfigModel(*pathway.GlobalConfig)
		path.Nodes = append(path.Nodes, globalNode)
	}
	for _, edge := range pathway.Edges {
		edgeModel := ConvertFromPathwayEdgeModel(edge)
		path.Edges = append(path.Edges, edgeModel)
	}
	return path
}
