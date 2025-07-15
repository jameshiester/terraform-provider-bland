// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways_test

import (
	"net/http"
	"regexp"
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

	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`^https://api.bland.ai/v1/pathway/create?api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("./tests/resource/pathway/Validate_Create/post_pathway.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "bland_conversational_pathway" "path" {
						name                              = "Test Provider Name"
						description                       = "Test Provider Description"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "name", "Test Provider Name"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "description", "Test Provider Description"),
					resource.TestCheckResourceAttr("bland_conversational_pathway.path", "id", "123"),
				),
			},
		},
	})
}
