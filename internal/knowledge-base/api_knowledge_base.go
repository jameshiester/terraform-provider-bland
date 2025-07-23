// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
)

type KnowledgeBaseClient struct {
	Api *api.Client
}

func newKnowledgeBaseClient(apiClient *api.Client) *KnowledgeBaseClient {
	return &KnowledgeBaseClient{Api: apiClient}
}

func (c *KnowledgeBaseClient) CreateKnowledgeBase(ctx context.Context, kb CreateKnowledgeBaseDto) (*string, error) {

	var created createKnowledgeBaseUploadResponseDto
	if kb.File != nil {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   c.Api.Config.BaseURL,
			Path:   "/v1/knowledgebases/upload",
		}

		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Add text fields
		writer.WriteField("name", kb.Name)
		writer.WriteField("description", kb.Description)

		// Add file
		if len(*kb.File) > 0 {
			part, err := writer.CreateFormFile("file", "knowledge_base_file")
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			_, err = part.Write(*kb.File)
			if err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}

		writer.Close()

		// Set content type header
		headers := http.Header{}
		headers.Set("Content-Type", writer.FormDataContentType())
		_, err := c.Api.Execute(ctx, nil, "POST", apiUrl.String(), headers, io.NopCloser(&buf), []int{http.StatusCreated}, &created)
		if err != nil {
			return nil, fmt.Errorf("failed to create knowledge base: %w", err)
		}
	} else {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   c.Api.Config.BaseURL,
			Path:   "/v1/knowledgebases",
		}
		_, err := c.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, kb, []int{http.StatusOK}, &created)
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %w", err)
		}
	}

	return &created.Data.ID, nil
}

func (c *KnowledgeBaseClient) ReadKnowledgeBase(ctx context.Context, id string) (*KnowledgeBaseDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/knowledgebases/%s", id),
	}
	var kb readKnowledgeBaseResponseDto
	_, err := c.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &kb)
	if err != nil {
		return nil, fmt.Errorf("failed to read knowledge base: %w", err)
	}
	result := KnowledgeBaseDto{
		ID:          id,
		Name:        kb.Data.Name,
		Description: kb.Data.Description,
		Text:        kb.Data.Text,
	}
	return &result, nil
}

func (c *KnowledgeBaseClient) UpdateKnowledgeBase(ctx context.Context, id string, kb UpdateKnowledgeBaseDto) (*KnowledgeBaseDto, error) {
	var updated createKnowledgeBaseUploadResponseDto
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/knowledgebases/%s", id),
	}
	_, err := c.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, kb, []int{http.StatusOK}, &updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	updatedDto := KnowledgeBaseDto{
		ID:          id,
		Name:        kb.Name,
		Description: kb.Description,
		File:        kb.File,
	}
	return &updatedDto, nil
}

func (c *KnowledgeBaseClient) DeleteKnowledgeBase(ctx context.Context, id string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/knowledgebases/%s", id),
	}
	_, err := c.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge base: %w", err)
	}
	return nil
}
