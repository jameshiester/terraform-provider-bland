// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
)

type KnowledgeBaseClient struct {
	Api *api.Client
}

func NewKnowledgeBaseClient(apiClient *api.Client) *KnowledgeBaseClient {
	return &KnowledgeBaseClient{Api: apiClient}
}

func (c *KnowledgeBaseClient) CreateKnowledgeBase(ctx context.Context, kbModel KnowledgeBaseModel) (*string, error) {
	createDto, err := ConvertToCreateKnowledgeBaseDto(kbModel)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for knowledge base: %w", err)
	}

	var created createKnowledgeBaseUploadResponseDto
	if createDto.File != nil {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   c.Api.Config.BaseURL,
			Path:   "/v1/knowledgebases/upload",
		}

		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Add text fields
		err := writer.WriteField("name", createDto.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		err = writer.WriteField("description", createDto.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		// Add file
		if len(*createDto.File) > 0 {
			filename := filepath.Base(kbModel.FilePath.ValueString())
			part, err := writer.CreateFormFile("file", filename)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			_, err = part.Write(*createDto.File)
			if err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}

		writer.Close()

		// Set content type header
		headers := http.Header{}
		headers.Set("Content-Type", writer.FormDataContentType())
		reqBody := bytes.NewReader(buf.Bytes())
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		_, err = c.Api.ExecuteMultipart(ctxWithTimeout, "POST", apiUrl.String(), headers, reqBody, []int{http.StatusOK}, &created)
		if err != nil {
			return nil, fmt.Errorf("failed to create knowledge base: %w", err)
		}
	} else {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   c.Api.Config.BaseURL,
			Path:   "/v1/knowledgebases",
		}
		_, err := c.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, createDto, []int{http.StatusOK}, &created)
		if err != nil {
			return nil, fmt.Errorf("failed to create knowledge base: %w", err)
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
	// Set content type header
	headers := http.Header{}
	headers.Set("Include-Text", "true")
	var kb readKnowledgeBaseResponseDto
	_, err := c.Api.Execute(ctx, nil, "GET", apiUrl.String(), headers, nil, []int{http.StatusOK}, &kb)
	if err != nil {
		return nil, fmt.Errorf("failed to read knowledge base: %w", err)
	}
	result := KnowledgeBaseDto{
		ID:            id,
		Name:          kb.Data.Name,
		Description:   kb.Data.Description,
		ExtractedText: kb.Data.ExtractedText,
	}
	return &result, nil
}

func (c *KnowledgeBaseClient) UpdateKnowledgeBase(ctx context.Context, id string, kbModel KnowledgeBaseModel) (*KnowledgeBaseDto, error) {
	updateDto, err := ConvertToUpdateKnowledgeBaseDto(kbModel)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for knowledge base: %w", err)
	}
	var updated createKnowledgeBaseUploadResponseDto
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/knowledgebases/%s", id),
	}
	_, err = c.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, updateDto, []int{http.StatusOK}, &updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update knowledge base: %w", err)
	}
	updatedDto := KnowledgeBaseDto{
		ID:          id,
		Name:        kbModel.Name.ValueString(),
		Description: kbModel.Description.ValueString(),
		File:        updateDto.File,
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
