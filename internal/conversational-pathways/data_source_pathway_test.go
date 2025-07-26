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

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/pathway/123`,
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
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.kb", "test-kb"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.kb_tool", "test-kb-tool"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.#", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.0.chosen_pathway", "The user has asked about something to do with rosters or rostering."),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.0.conversation_history.basic_history", "i want to talk to the rostering team"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.chosen_pathway", "User responded"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.#", "3"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.0.role", "user"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.0.content", "Hello?"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.1.role", "assistant"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.1.content", "Hello, this is YLDP. How can I help you today?"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.2.role", "user"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.pathway_examples.1.conversation_history.advanced_history.2.content", "I broke a door in  my house"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.headers.#", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.headers.0.name", "a"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.headers.0.value", "val"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.headers.1.name", "b"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.headers.1.value", "val2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.auth.type", "Bearer"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.auth.token", "124"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.auth.encode", "false"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.body", "test body"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.#", "1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.conditions.#", "1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.conditions.0.field", "expected_annual_salary"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.conditions.0.value", "500000"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.conditions.0.is_group", "false"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.conditions.0.operator", "less than"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.routes.0.target_node_id", "78136d68-d3d7-4d91-917e-26853c830d09"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.0.data.fallback_node_id", "fallback-node-id"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.id", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.type", "Default"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.data.name", "Edge 2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "nodes.1.data.prompt", "Second edge information"),

					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.id", "Edge1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.source", "1"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.target", "2"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.data.label", "Edge Label"),
					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "edges.0.data.always_pick", "true"),

					resource.TestCheckResourceAttr("data.bland_conversational_pathway.pathway", "global_config.global_prompt", "Example global prompt"),
				),
			},
		},
	})
}
