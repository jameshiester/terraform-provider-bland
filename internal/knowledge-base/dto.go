// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

type KnowledgeBaseDto struct {
	ID            string  `json:"id,omitempty"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Text          string  `json:"text"`
	ExtractedText *string `json:"-"` // not included in response
	File          *[]byte `json:"-"` // Binary data, not serialized to JSON
}

type readKnowledgeBaseResponseDataDto struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ExtractedText *string `json:"text"`
}

type readKnowledgeBaseResponseDto struct {
	Data readKnowledgeBaseResponseDataDto `json:"data"`
}

type CreateKnowledgeBaseDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Text        *string `json:"text,omitempty"`
	File        *[]byte `json:"-"` // Binary data for multipart form
}

type createKnowledgeBaseUploadResponseDataDto struct {
	ID string `json:"vector_id"`
}

type createKnowledgeBaseUploadResponseDto struct {
	Data createKnowledgeBaseUploadResponseDataDto `json:"data"`
}

type UpdateKnowledgeBaseDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Text        *string `json:"text,omitempty"`
	File        *[]byte `json:"-"` // Binary data for multipart form
}
