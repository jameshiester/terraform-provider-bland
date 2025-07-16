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

// func TestAccConversationalPathwayResource_Validate_Create(t *testing.T) {
// 	resource.Test(t, resource.TestCase{

// 		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 					resource "bland_conversational_pathway" "path" {
// 						name                              = "Test Provider Name"
// 						description                       = "Test Provider Description"
// 					}
// 					`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "name", "Test Provider Name"),
// 					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "description", "Test Provider Description"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitConversationalPathwayResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/pathway/create?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("./tests/resource/pathway/Validate_Create/post_pathway.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/pathway/123?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/pathway/Validate_Create/update_pathway.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.bland.ai/v1/pathway/123?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/pathway/123?api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/pathway/Validate_Create/get_pathway.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "bland_conversational_pathway" "path" {
						name                              = "TestPathwayName"
						description                       = "TestPathwayDescription"
						nodes = [
							{
								id = "1"
								type = "Default"
								data = {
              						name = "Start"
              						text = "Hey there, how are you doing today?"
              						is_start = true
           						}
							}
						]
						global_config = {
							global_prompt = "Example global prompt"
						}
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "name", "TestPathwayName"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "description", "TestPathwayDescription"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "id", "123"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "nodes.0.id", "1"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "global_config.global_prompt", "Example global prompt"),
				),
			},
		},
	})
}
