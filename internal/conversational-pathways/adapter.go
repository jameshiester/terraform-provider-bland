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

func ConvertFromPathwayNodeDataDto(data pathwayNodeDataDto) (*ConversationalPathwayNodeDataModel, error) {
	model := ConversationalPathwayNodeDataModel{
		GlobalPrompt: types.StringPointerValue(data.GlobalPrompt),
		GlobalLabel:  types.StringPointerValue(data.GlobalLabel),
		Method:       types.StringPointerValue(data.Method),
		IsStart:      types.BoolPointerValue(data.IsStart),
		Name:         types.StringValue(data.Name),
		Prompt:       types.StringPointerValue(data.Prompt),
		Text:         types.StringPointerValue(data.Text),
		URL:          types.StringPointerValue(data.URL),
	}
	for _, variable := range *data.ExtractVars {
		varModel, err := ConvertFromPathwayNodeDataExtractVars(variable)
		if err != nil {
			return nil, err
		}
		model.ExtractVars = append(model.ExtractVars, varModel)
	}
	for _, responseData := range *data.ResponseData {
		responseModel := ConvertFromPathwayNodeDataResponseData(responseData)
		model.ResponseData = append(model.ResponseData, responseModel)
	}
	for _, responsePathwayData := range *data.ResponsePathways {
		responsePathway, err := ConvertFromPathwayNodeDataResponsePathway(responsePathwayData)
		if err != nil {
			return nil, err
		}
		model.ResponsePathways = append(model.ResponsePathways, *responsePathway)
	}
	return &model, nil
}

func ConvertFromPathwayNodeDataModel(data ConversationalPathwayNodeDataModel) pathwayNodeDataDto {
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

	return pathwayNodeDataDto{
		Name:             data.Name.ValueString(),
		GlobalPrompt:     data.GlobalPrompt.ValueStringPointer(),
		Prompt:           data.Prompt.ValueStringPointer(),
		Text:             data.Text.ValueStringPointer(),
		IsStart:          data.IsStart.ValueBoolPointer(),
		Method:           data.Method.ValueStringPointer(),
		URL:              data.URL.ValueStringPointer(),
		GlobalLabel:      data.GlobalLabel.ValueStringPointer(),
		ExtractVars:      extractVars,
		ResponseData:     responseData,
		ResponsePathways: responsePathways,
	}
}

func ConvertFromPathwayNodeDto(node pathwayNodeDto) (*ConversationalPathwayNodeModel, error) {
	data, err := ConvertFromPathwayNodeDataDto(node.Data)
	if err != nil {
		return nil, err
	}
	return &ConversationalPathwayNodeModel{
		ID:   types.StringValue(node.ID),
		Type: types.StringValue(node.Type),
		Data: *data,
	}, nil
}

func ConvertFromPathwayNodeModel(node ConversationalPathwayNodeModel) pathwayNodeDto {
	return pathwayNodeDto{
		ID:   node.ID.ValueString(),
		Type: node.Type.ValueString(),
		Data: ConvertFromPathwayNodeDataModel(node.Data),
	}
}

func ConvertFromPathwayDto(pathway pathwayDto) (*ConversationalPathwayDataSourceModel, error) {

	path := ConversationalPathwayDataSourceModel{
		ID:          types.StringValue(pathway.ID),
		Name:        types.StringValue(pathway.Name),
		Description: types.StringValue(pathway.Description),
	}
	for _, node := range pathway.Nodes {
		nodeModel, err := ConvertFromPathwayNodeDto(node)
		if err != nil {
			return nil, err
		}
		path.Nodes = append(path.Nodes, *nodeModel)
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

	return path
}
