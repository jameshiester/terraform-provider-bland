// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertFromPathwayNodeDataExtractVars(vals []interface{}) (ConversationalPathwayNodeDataExtractVariableModel, error) {
	model := ConversationalPathwayNodeDataExtractVariableModel{}
	if len(vals) != 3 && len(vals) != 4 {
		return model, fmt.Errorf("failed to get pathway: %d", len(vals))
	}

	// Convert first three values to strings
	if name, ok := vals[0].(string); ok {
		model.Name = types.StringValue(name)
	} else {
		return model, fmt.Errorf("name must be a string")
	}

	if varType, ok := vals[1].(string); ok {
		model.Type = types.StringValue(varType)
	} else {
		return model, fmt.Errorf("type must be a string")
	}

	if description, ok := vals[2].(string); ok {
		model.Description = types.StringValue(description)
	} else {
		return model, fmt.Errorf("description must be a string")
	}

	// Handle optional fourth value (boolean)
	if len(vals) == 4 {
		if increaseSpellingPrecision, ok := vals[3].(bool); ok {
			model.IncreaseSpellingPrecision = types.BoolValue(increaseSpellingPrecision)
		} else {
			return model, fmt.Errorf("increase_spelling_precision must be a boolean")
		}
	}

	return model, nil
}

func ConvertFromPathwayNodeDataResponseData(data pathwayNodeDataResponseDataDto) ConversationalPathwayNodeDataResponseDataModel {
	return ConversationalPathwayNodeDataResponseDataModel{
		Name:    types.StringValue(data.Name),
		Data:    types.StringValue(data.Data),
		Context: types.StringPointerValue(data.Context),
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
		KbTool:         types.StringPointerValue(data.KbTool),
		TransferNumber: types.StringPointerValue(data.TransferNumber),
		Body:           types.StringPointerValue(data.Body),
		FallbackNodeId: types.StringPointerValue(data.FallbackNodeId),
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
	if data.Auth != nil {
		model.Auth = &ConversationalPathwayAuthModel{
			Type:   types.StringValue(data.Auth.Type),
			Token:  types.StringValue(data.Auth.Token),
			Encode: types.BoolValue(data.Auth.Encode),
		}
	}
	if data.Headers != nil {
		for _, h := range *data.Headers {
			if len(h) == 2 {
				model.Headers = append(model.Headers, ConversationalPathwayHeaderModel{
					Name:  types.StringValue(h[0]),
					Value: types.StringValue(h[1]),
				})
			}
		}
	}
	if data.Routes != nil {
		for _, route := range *data.Routes {
			routeModel := ConversationalPathwayRouteModel{
				TargetNodeId: types.StringValue(route.TargetNodeId),
			}
			for _, condition := range route.Conditions {
				routeModel.Conditions = append(routeModel.Conditions, ConversationalPathwayRouteConditionModel{
					Field:    types.StringValue(condition.Field),
					Value:    types.StringValue(condition.Value),
					IsGroup:  types.BoolValue(condition.IsGroup),
					Operator: types.StringValue(condition.Operator),
				})
			}
			model.Routes = append(model.Routes, routeModel)
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
	if data.PathwayExamples != nil {
		for _, ex := range *data.PathwayExamples {
			model.PathwayExamples = append(model.PathwayExamples, ConvertFromPathwayExampleDto(ex))
		}
	}
	return &model, nil
}

func ConvertFromPathwayExampleDto(dto pathwayExampleDto) ConversationalPathwayExampleModel {
	model := ConversationalPathwayExampleModel{
		ChosenPathway: types.StringValue(dto.ChosenPathway),
	}
	if dto.ConversationHistory.BasicHistory != nil {
		model.ConversationHistory.BasicHistory = types.StringValue(*dto.ConversationHistory.BasicHistory)
	}
	if dto.ConversationHistory.AdvancedHistory != nil {
		for _, msg := range *dto.ConversationHistory.AdvancedHistory {
			model.ConversationHistory.AdvancedHistory = append(model.ConversationHistory.AdvancedHistory, ConvertFromPathwayExampleMessageDto(msg))
		}
	}
	return model
}

func ConvertFromPathwayExampleMessageDto(dto pathwayExampleMessageDto) ConversationalPathwayExampleMessageModel {
	return ConversationalPathwayExampleMessageModel{
		Role:    types.StringValue(dto.Role),
		Content: types.StringValue(dto.Content),
	}
}

func ConvertFromPathwayNodeDataModel(data ConversationalPathwayNodeDataModel) *pathwayNodeDataDto {
	var extractVars *[][]interface{}
	if len(data.ExtractVars) > 0 {
		tmp := make([][]interface{}, 0, len(data.ExtractVars))
		for _, v := range data.ExtractVars {
			if !v.IncreaseSpellingPrecision.IsNull() {
				tmp = append(tmp, []interface{}{
					v.Name.ValueString(),
					v.Type.ValueString(),
					v.Description.ValueString(),
					v.IncreaseSpellingPrecision.ValueBool(),
				})
			} else {
				tmp = append(tmp, []interface{}{
					v.Name.ValueString(),
					v.Type.ValueString(),
					v.Description.ValueString(),
				})
			}
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
				Context: v.Context.ValueStringPointer(),
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

	var pathwayExamples *[]pathwayExampleDto
	if len(data.PathwayExamples) > 0 {
		tmp := make([]pathwayExampleDto, 0, len(data.PathwayExamples))
		for _, ex := range data.PathwayExamples {
			tmp = append(tmp, ConvertFromPathwayExampleModel(ex))
		}
		pathwayExamples = &tmp
	} else {
		pathwayExamples = nil
	}

	var headers *[][]string
	if len(data.Headers) > 0 {
		tmp := make([][]string, 0, len(data.Headers))
		for _, h := range data.Headers {
			tmp = append(tmp, []string{h.Name.ValueString(), h.Value.ValueString()})
		}
		headers = &tmp
	} else {
		headers = nil
	}

	var auth *AuthDto
	if data.Auth != nil {
		auth = &AuthDto{
			Type:   data.Auth.Type.ValueString(),
			Token:  data.Auth.Token.ValueString(),
			Encode: data.Auth.Encode.ValueBool(),
		}
	} else {
		auth = nil
	}
	var body *string
	if !data.Body.IsNull() && !data.Body.IsUnknown() {
		b := data.Body.ValueString()
		body = &b
	}

	var routes *[]RouteDto
	if len(data.Routes) > 0 {
		tmp := make([]RouteDto, 0, len(data.Routes))
		for _, route := range data.Routes {
			routeDto := RouteDto{
				TargetNodeId: route.TargetNodeId.ValueString(),
			}
			for _, condition := range route.Conditions {
				routeDto.Conditions = append(routeDto.Conditions, RouteConditionDto{
					Field:    condition.Field.ValueString(),
					Value:    condition.Value.ValueString(),
					IsGroup:  condition.IsGroup.ValueBool(),
					Operator: condition.Operator.ValueString(),
				})
			}
			tmp = append(tmp, routeDto)
		}
		routes = &tmp
	} else {
		routes = nil
	}
	var fallbackNodeId *string
	if !data.FallbackNodeId.IsNull() && !data.FallbackNodeId.IsUnknown() {
		f := data.FallbackNodeId.ValueString()
		fallbackNodeId = &f
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
		KbTool:           data.KbTool.ValueStringPointer(),
		TransferNumber:   data.TransferNumber.ValueStringPointer(),
		ExtractVars:      extractVars,
		ResponseData:     responseData,
		ResponsePathways: responsePathways,
		ModelOptions:     modelOptions,
		PathwayExamples:  pathwayExamples,
		Headers:          headers,
		Auth:             auth,
		Body:             body,
		Routes:           routes,
		FallbackNodeId:   fallbackNodeId,
	}
}

func ConvertFromPathwayExampleModel(model ConversationalPathwayExampleModel) pathwayExampleDto {
	dto := pathwayExampleDto{
		ChosenPathway: model.ChosenPathway.ValueString(),
	}
	if !model.ConversationHistory.BasicHistory.IsNull() && !model.ConversationHistory.BasicHistory.IsUnknown() {
		str := model.ConversationHistory.BasicHistory.ValueString()
		dto.ConversationHistory.BasicHistory = &str
	}
	if len(model.ConversationHistory.AdvancedHistory) > 0 {
		arr := make([]pathwayExampleMessageDto, 0, len(model.ConversationHistory.AdvancedHistory))
		for _, msg := range model.ConversationHistory.AdvancedHistory {
			arr = append(arr, ConvertFromPathwayExampleMessageModel(msg))
		}
		dto.ConversationHistory.AdvancedHistory = &arr
	}
	return dto
}

func ConvertFromPathwayExampleMessageModel(model ConversationalPathwayExampleMessageModel) pathwayExampleMessageDto {
	return pathwayExampleMessageDto{
		Role:    model.Role.ValueString(),
		Content: model.Content.ValueString(),
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
	nodeType := edge.Type
	if nodeType == "" {
		nodeType = "custom"
	}
	model := ConversationalPathwayEdgeModel{
		ID:     types.StringValue(edge.ID),
		Source: types.StringValue(edge.Source),
		Target: types.StringValue(edge.Target),
		Type:   types.StringValue(nodeType),
		Data: ConversationalPathwayEdgeDataModel{
			Label:         types.StringValue(edge.Data.Label),
			IsHighlighted: types.BoolValue(edge.Data.IsHighlighted),
			Description:   types.StringPointerValue(edge.Data.Description),
			AlwaysPick:    types.BoolPointerValue(edge.Data.AlwaysPick),
		},
	}
	if edge.Data.Condition != nil {
		for _, condition := range *edge.Data.Condition {
			model.Data.Conditions = &[]ConversationalPathwayEdgeConditionModel{}
			*model.Data.Conditions = append(*model.Data.Conditions, ConversationalPathwayEdgeConditionModel{
				Field:    types.StringValue(condition.Field),
				Value:    types.StringValue(condition.Value),
				IsGroup:  types.BoolValue(condition.IsGroup),
				Operator: types.StringValue(condition.Operator),
			})
		}
	}
	return model
}

func ConvertFromPathwayEdgeModel(edge ConversationalPathwayEdgeModel) pathwayEdgeDto {
	edgeData := pathwayEdgeDataDto{
		Label:         edge.Data.Label.ValueString(),
		IsHighlighted: edge.Data.IsHighlighted.ValueBool(),
		Description:   edge.Data.Description.ValueStringPointer(),
		AlwaysPick:    edge.Data.AlwaysPick.ValueBoolPointer(),
	}
	if edge.Data.Conditions != nil {
		for _, condition := range *edge.Data.Conditions {
			if edgeData.Condition == nil {
				edgeData.Condition = &[]EdgeConditionDto{}
			}
			*edgeData.Condition = append(*edgeData.Condition, EdgeConditionDto{
				Field:    condition.Field.ValueString(),
				Value:    condition.Value.ValueString(),
				IsGroup:  condition.IsGroup.ValueBool(),
				Operator: condition.Operator.ValueString(),
			})
		}
	}
	return pathwayEdgeDto{
		ID:     edge.ID.ValueString(),
		Source: edge.Source.ValueString(),
		Target: edge.Target.ValueString(),
		Type:   edge.Type.ValueString(),
		Data:   edgeData,
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
