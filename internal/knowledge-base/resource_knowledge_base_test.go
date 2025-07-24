// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jameshiester/terraform-provider-bland/internal/mocks"
	"github.com/jarcoal/httpmock"
)

func TestAccKnowledgeBaseResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "bland_knowledge_base" "kb" {
						name        = "TestKnowledgeBase"
						description = "Test knowledge base description"
						file_path        = "./tests/example.txt"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "name", "TestKnowledgeBase"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "description", "Test knowledge base description"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "extracted_text", "test file content"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "file_path", "./tests/example.txt"),
				),
			},
		},
	})
}

func TestUnitKnowledgeBaseResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Use a static test file for file_path
	filePath := "./tests/example.txt"

	// Mock file upload endpoint
	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/knowledgebases/upload",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/knowledge_base/Validate_Create/post_knowledge_base.json").String()), nil
		})

	// Mock read endpoint
	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/knowledgebases/kb_123`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/knowledge_base/Validate_Create/get_knowledge_base.json").String()), nil
		})

	// Mock update endpoint
	httpmock.RegisterResponder("PATCH", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/knowledge_base/Validate_Create/update_knowledge_base.json").String()), nil
		})

	// Mock delete endpoint
	httpmock.RegisterResponder("DELETE", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "bland_knowledge_base" "kb" {
						name        = "TestKnowledgeBase"
						description = "Test knowledge base description"
						file_path   = "%s"
					}
				`, filePath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "name", "TestKnowledgeBase"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "description", "Test knowledge base description"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "id", "kb_123"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "extracted_text", "This is the extracted text from the knowledge base file."),
				),
			},
		},
	})
}

func TestUnitKnowledgeBaseResource_Validate_Create_TextOnly(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock text-only create endpoint
	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/knowledgebases",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/knowledge_base/Validate_Create/post_knowledge_base.json").String()), nil
		})

	// Mock read endpoint
	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/knowledgebases/kb_123`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/knowledge_base/Validate_Create/get_knowledge_base.json").String()), nil
		})

	// Mock delete endpoint
	httpmock.RegisterResponder("DELETE", "https://api.bland.ai/v1/knowledgebases/kb_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "bland_knowledge_base" "kb" {
						name        = "TestKnowledgeBase"
						description = "Test knowledge base description"
						text        = "This is the extracted text from the knowledge base file."
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "name", "TestKnowledgeBase"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "description", "Test knowledge base description"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "id", "kb_123"),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "text", "This is the extracted text from the knowledge base file."),
					resource.TestCheckResourceAttr("bland_knowledge_base.kb", "extracted_text", "This is the extracted text from the knowledge base file."),
				),
			},
		},
	})
}
