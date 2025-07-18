// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"encoding/json"
	"fmt"
)

type createPathwayDto struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       []pathwayNodeDto `json:"nodes"`
	Edges       []pathwayEdgeDto `json:"edges"`
}

type updatePathwayDto struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       []pathwayNodeDto `json:"nodes"`
	Edges       []pathwayEdgeDto `json:"edges"`
}

type pathwayDto struct {
	ID          string           `json:"pathway_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       []pathwayNodeDto `json:"nodes"`
	Edges       []pathwayEdgeDto `json:"edges"`
}

type getPathwayDto struct {
	ID          string      `json:"pathway_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Nodes       NodesOrBool `json:"nodes"`
	Edges       EdgesOrBool `json:"edges"`
}

type pathwayGlobalConfigDto struct {
	GlobalPrompt string `json:"globalPrompt"`
}

type pathwayNodeDto struct {
	ID           *string                 `json:"id"`
	Type         *string                 `json:"type"`
	GlobalConfig *pathwayGlobalConfigDto `json:"globalConfig,omitempty"`
	Data         *pathwayNodeDataDto     `json:"data"`
}

type pathwayEdgeDto struct {
	ID     string             `json:"id"`
	Source string             `json:"source"`
	Target string             `json:"target"`
	Type   string             `json:"type"`
	Data   pathwayEdgeDataDto `json:"data"`
}

type pathwayEdgeDataDto struct {
	Label         string  `json:"label"`
	IsHighlighted bool    `json:"isHighlighted"`
	Description   *string `json:"description,omitempty"`
}

type pathwayNodeDataResponseDataDto struct {
	Data    string `json:"data"`
	Name    string `json:"name"`
	Context string `json:"context"`
}

type pathwayNodeDataDto struct {
	Name             string                            `json:"name"`
	Text             *string                           `json:"text,omitempty"`
	GlobalPrompt     *string                           `json:"global_prompt,omitempty"`
	Prompt           *string                           `json:"prompt,omitempty"`
	IsStart          *bool                             `json:"isStart,omitempty"`
	IsGlobal         *bool                             `json:"isGlobal,omitempty"`
	GlobalLabel      *string                           `json:"globalLabel,omitempty"`
	URL              *string                           `json:"url,omitempty"`
	Method           *string                           `json:"method,omitempty"`
	ExtractVars      *[][]string                       `json:"extractVars,omitempty"`
	ResponseData     *[]pathwayNodeDataResponseDataDto `json:"responseData,omitempty"`
	ResponsePathways *[][]interface{}                  `json:"responsePathways,omitempty"`
	Condition        *string                           `json:"condition,omitempty"`
	KnowledgeBase    *string                           `json:"kb,omitempty"`
	TransferNumber   *string                           `json:"transferNumber,omitempty"`
	ModelOptions     *modelOptionDto                   `json:"modelOptions,omitempty"`
	PathwayExamples  *[]pathwayExampleDto              `json:"pathway_examples,omitempty"`
}

type modelOptionDto struct {
	Type                  string   `json:"modelType"`
	InterruptionThreshold *string  `json:"interruption_threshold,omitempty"`
	Temperature           *float32 `json:"temperature,omitempty"`
	SkipUserResponse      *bool    `json:"skipUserResponse,omitempty"`
	BlockInterruptions    *bool    `json:"block_interruptions,omitempty"`
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

// Custom type for nodes that can be a boolean or an array
// If boolean, will be nil. If array, will be the array.
type NodesOrBool []pathwayNodeDto

func (n *NodesOrBool) UnmarshalJSON(data []byte) error {
	// Check for boolean
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*n = nil
		return nil
	}
	// Otherwise, try as array
	var arr []pathwayNodeDto
	if err := json.Unmarshal(data, &arr); err == nil {
		*n = arr
		return nil
	}
	return fmt.Errorf("nodes is neither bool nor array")
}

// Custom type for edges that can be a boolean or an array
// If boolean, will be nil. If array, will be the array.
type EdgesOrBool []pathwayEdgeDto

func (e *EdgesOrBool) UnmarshalJSON(data []byte) error {
	// Check for boolean
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*e = nil
		return nil
	}
	// Otherwise, try as array
	var arr []pathwayEdgeDto
	if err := json.Unmarshal(data, &arr); err == nil {
		*e = arr
		return nil
	}
	return fmt.Errorf("edges is neither bool nor array")
}

// pathwayExampleDto supports both string and array of messages for Conversation History

type pathwayExampleDto struct {
	ChosenPathway       string                `json:"Chosen Pathway"`
	ConversationHistory pathwayExampleHistory `json:"Conversation History"`
}

type pathwayExampleHistory struct {
	BasicHistory    *string
	AdvancedHistory *[]pathwayExampleMessageDto
}

type pathwayExampleMessageDto struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Custom unmarshaller for Conversation History
func (h *pathwayExampleHistory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		h.BasicHistory = &s
		return nil
	}
	var arr []pathwayExampleMessageDto
	if err := json.Unmarshal(data, &arr); err == nil {
		h.AdvancedHistory = &arr
		return nil
	}
	return fmt.Errorf("conversation History is neither string nor array of messages")
}
