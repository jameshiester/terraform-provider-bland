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
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Nodes           []pathwayNodeDto `json:"nodes"`
	Edges           []pathwayEdgeDto `json:"edges"`
	Revision        int              `json:"revision_number"`
	Version         int              `json:"version_number"`
	PostCallActions []string         `json:"post_call_actions"`
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
	Label         string             `json:"label"`
	IsHighlighted bool               `json:"isHighlighted"`
	Description   *string            `json:"description,omitempty"`
	AlwaysPick    *bool              `json:"alwaysPick,omitempty"`
	Condition     []EdgeConditionDto `json:"condition,omitempty"`
}

type pathwayNodeDataResponseDataDto struct {
	Data    string  `json:"data"`
	Name    string  `json:"name"`
	Context *string `json:"context,omitempty"`
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
	ExtractVars      *[][]interface{}                  `json:"extractVars,omitempty"`
	ResponseData     *[]pathwayNodeDataResponseDataDto `json:"responseData,omitempty"`
	ResponsePathways *[][]interface{}                  `json:"responsePathways,omitempty"`
	Condition        *string                           `json:"condition,omitempty"`
	KnowledgeBase    *string                           `json:"kb,omitempty"`
	KbTool           *string                           `json:"kbTool,omitempty"`
	TransferNumber   *string                           `json:"transferNumber,omitempty"`
	ModelOptions     *modelOptionDto                   `json:"modelOptions,omitempty"`
	PathwayExamples  *[]pathwayExampleDto              `json:"pathway_examples,omitempty"`
	Headers          *[][]string                       `json:"headers,omitempty"`
	Auth             *AuthDto                          `json:"auth,omitempty"`
	Body             *string                           `json:"body,omitempty"`
	Routes           *[]RouteDto                       `json:"routes,omitempty"`
	FallbackNodeId   *string                           `json:"fallbackNodeId,omitempty"`
}

type AuthDto struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	Encode bool   `json:"encode"`
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

type updatePathwayDataDto struct {
	Message string `json:"message"`
}

type updatePathwayResponseDto struct {
	Errors *[]errorDto           `json:"errors,omitempty"`
	Data   *updatePathwayDataDto `json:"pathway_data"`
}

type pathwayVersionDto struct {
	VersionNumber       int    `json:"version_number"`
	RevisionNumber      int    `json:"revision_number"`
	CreatedAt           string `json:"created_at"`
	Name                string `json:"name"`
	SourceVersionNumber *int   `json:"source_version_number"`
	IsStaging           *bool  `json:"is_staging,omitempty"`
	IsProduction        *bool  `json:"is_production,omitempty"`
	IsPrevPublished     *bool  `json:"is_prev_published,omitempty"`
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

// Custom unmarshaller for Conversation History.
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

type RouteConditionDto struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	IsGroup  bool   `json:"isGroup"`
	Operator string `json:"operator"`
}

type RouteDto struct {
	Conditions   []RouteConditionDto `json:"conditions"`
	TargetNodeId string              `json:"targetNodeId"`
}

type EdgeConditionDto struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	IsGroup  bool   `json:"isGroup"`
	Operator string `json:"operator"`
}
