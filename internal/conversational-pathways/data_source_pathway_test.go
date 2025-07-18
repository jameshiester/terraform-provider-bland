// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jameshiester/terraform-provider-bland/internal/mocks"
	"github.com/jarcoal/httpmock"
)

func TestUnitConversationalPathwayDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/pathway/123?api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/datasource/Validate_Read/get_pathway.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "bland_conversational_pathway" "pathway" {
					id = "123"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "id", "123"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "name", "TestPathwayName"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "description", "TestPathwayDescription"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.#", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.id", "1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.type", "Default"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.name", "Start"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.text", "Hey there, how are you doing today?"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.is_start", "true"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.id", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.type", "Default"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.data.name", "Edge 2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.data.prompt", "Second edge information"),

					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.id", "Edge1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.source", "1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.target", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.data.label", "Edge Label"),

					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "global_config.global_prompt", "Example global prompt"),
				),
			},
		},
	})
}
