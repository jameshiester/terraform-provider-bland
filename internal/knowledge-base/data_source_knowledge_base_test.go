// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jameshiester/terraform-provider-bland/internal/mocks"
	"github.com/jarcoal/httpmock"
)

func TestUnitKnowledgeBaseDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/knowledgebases/kb_123`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/datasource/Validate_Read/get_knowledge_base.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "bland_knowledge_base" "kb" {
					id = "kb_123"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bland_knowledge_base.kb", "id", "kb_123"),
					resource.TestCheckResourceAttr("data.bland_knowledge_base.kb", "name", "TestKnowledgeBase"),
					resource.TestCheckResourceAttr("data.bland_knowledge_base.kb", "description", "Test knowledge base description"),
					resource.TestCheckResourceAttr("data.bland_knowledge_base.kb", "text", "This is the extracted text from the knowledge base file."),
				),
			},
		},
	})
}
