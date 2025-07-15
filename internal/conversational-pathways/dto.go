// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

type createPathwayDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updatePathwayDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type pathwayDto struct {
	ID          string
	Name        string
	Description string
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
