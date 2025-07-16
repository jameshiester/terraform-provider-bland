// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

type createPathwayDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updatePathwayDto struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       []pathwayNodeDto `json:"nodes"`
}

type pathwayDto struct {
	ID          string           `json:"pathway_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       []pathwayNodeDto `json:"nodes"`
}

type pathwayNodeDto struct {
	ID   string             `json:"id"`
	Type string             `json:"type"`
	Data pathwayNodeDataDto `json:"data"`
}

type pathwayNodeDataDto struct {
	Name         string `json:"name"`
	Text         string `json:"text"`
	GlobalPrompt string `json:"global_prompt"`
	Prompt       string `json:"prompt"`
	IsStart      bool   `json:"isStart"`
}

type createPathwayResponseData struct {
	ID string `json:"pathway_id"`
}

type errorDto struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type createPathwayResponseDto struct {
	Errors *[]errorDto                `json:"errors,omitempty"`
	Data   *createPathwayResponseData `json:"data"`
}

type updatePathwayResponseDto struct {
	Errors *[]errorDto       `json:"errors,omitempty"`
	Data   *updatePathwayDto `json:"pathway_data"`
}
