// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/config"
	knowledgebase "github.com/jameshiester/terraform-provider-bland/internal/knowledge-base"
	"github.com/jarcoal/httpmock"
)

func TestKnowledgeBaseClient_CreateKnowledgeBase_WithFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	fileData := []byte("test file content")
	dto := knowledgebase.CreateKnowledgeBaseDto{
		Name:        "Test KB",
		Description: "Test Description",
		File:        &fileData,
	}

	result, err := client.CreateKnowledgeBase(context.Background(), dto)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if *result != "kb_123" {
		t.Errorf("Expected vector_id 'kb_123', got '%s'", *result)
	}
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

	text := "test text content"
	dto := knowledgebase.CreateKnowledgeBaseDto{
		Name:        "Test KB",
		Description: "Test Description",
		Text:        &text,
	}

	result, err := client.CreateKnowledgeBase(context.Background(), dto)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if *result != "kb_456" {
		t.Errorf("Expected vector_id 'kb_456', got '%s'", *result)
	}
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

	if result.Text != "Extracted text content" {
		t.Errorf("Expected text 'Extracted text content', got '%s'", result.Text)
	}
}

func TestKnowledgeBaseClient_UpdateKnowledgeBase(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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

	fileData := []byte("updated file content")
	dto := knowledgebase.UpdateKnowledgeBaseDto{
		Name:        "Updated KB",
		Description: "Updated Description",
		File:        &fileData,
	}

	result, err := client.UpdateKnowledgeBase(context.Background(), "kb_123", dto)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "kb_123" {
		t.Errorf("Expected ID 'kb_123', got '%s'", result.ID)
	}
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
