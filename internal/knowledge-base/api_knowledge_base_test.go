// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/config"
	knowledgebase "github.com/jameshiester/terraform-provider-bland/internal/knowledge-base"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestKnowledgeBaseClient_CreateKnowledgeBase_WithFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	filePath := "./tests/example.txt"

	// Mock the file upload endpoint
	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/knowledgebases/upload",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `{"data":{"vector_id":"kb_123"}}`), nil
		})

	providerConfig := &config.ProviderConfig{
		BaseURL: "api.bland.ai",
		APIKey:  "123",
	}
	apiClient := api.NewApiClientBase(providerConfig, api.NewAuthBase(providerConfig))
	client := knowledgebase.NewKnowledgeBaseClient(apiClient)

	model := knowledgebase.KnowledgeBaseModel{
		Name:        types.StringValue("Test KB"),
		Description: types.StringValue("Test Description"),
		FilePath:    types.StringValue(filePath),
	}

	result, err := client.CreateKnowledgeBase(context.Background(), model)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "kb_123", *result)
}

func TestKnowledgeBaseClient_CreateKnowledgeBase_TextOnly(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the text-only create endpoint
	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/knowledgebases",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `{"data":{"vector_id":"kb_456"}}`), nil
		})

	providerConfig := &config.ProviderConfig{
		BaseURL: "api.bland.ai",
		APIKey:  "123",
	}
	apiClient := api.NewApiClientBase(providerConfig, api.NewAuthBase(providerConfig))
	client := knowledgebase.NewKnowledgeBaseClient(apiClient)

	model := knowledgebase.KnowledgeBaseModel{
		Name:        types.StringValue("Test KB"),
		Description: types.StringValue("Test Description"),
		Text:        types.StringValue("test text content"),
	}

	result, err := client.CreateKnowledgeBase(context.Background(), model)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "kb_456", *result)
}

func TestKnowledgeBaseClient_ReadKnowledgeBase(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the read endpoint
	httpmock.RegisterResponder("GET", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `{
				"data": {
					"name": "Test Knowledge Base",
					"description": "Test Description",
					"text": "Extracted text content"
				}
			}`), nil
		})

	providerConfig := &config.ProviderConfig{
		BaseURL: "api.bland.ai",
		APIKey:  "123",
	}
	apiClient := api.NewApiClientBase(providerConfig, api.NewAuthBase(providerConfig))
	client := knowledgebase.NewKnowledgeBaseClient(apiClient)

	result, err := client.ReadKnowledgeBase(context.Background(), "kb_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "kb_123" {
		t.Errorf("Expected ID 'kb_123', got '%s'", result.ID)
	}

	if result.Name != "Test Knowledge Base" {
		t.Errorf("Expected name 'Test Knowledge Base', got '%s'", result.Name)
	}

	if result.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", result.Description)
	}

	if *result.ExtractedText != "Extracted text content" {
		t.Errorf("Expected text 'Extracted text content', got '%s'", result.Text)
	}
}

func TestKnowledgeBaseClient_UpdateKnowledgeBase(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	filePath := "./tests/example.txt"

	// Mock the update endpoint
	httpmock.RegisterResponder("PATCH", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `{"data":{"vector_id":"kb_123"}}`), nil
		})

	providerConfig := &config.ProviderConfig{
		BaseURL: "api.bland.ai",
		APIKey:  "123",
	}
	apiClient := api.NewApiClientBase(providerConfig, api.NewAuthBase(providerConfig))
	client := knowledgebase.NewKnowledgeBaseClient(apiClient)

	model := knowledgebase.KnowledgeBaseModel{
		Name:        types.StringValue("Updated KB"),
		Description: types.StringValue("Updated Description"),
		FilePath:    types.StringValue(filePath),
	}

	result, err := client.UpdateKnowledgeBase(context.Background(), "kb_123", model)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "kb_123", result.ID)
}

func TestKnowledgeBaseClient_DeleteKnowledgeBase(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the delete endpoint
	httpmock.RegisterResponder("DELETE", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	providerConfig := &config.ProviderConfig{
		BaseURL: "api.bland.ai",
		APIKey:  "123",
	}
	apiClient := api.NewApiClientBase(providerConfig, api.NewAuthBase(providerConfig))
	client := knowledgebase.NewKnowledgeBaseClient(apiClient)

	err := client.DeleteKnowledgeBase(context.Background(), "kb_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
