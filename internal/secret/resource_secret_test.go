// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package secret_test

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

func TestUnitSecretResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.bland.ai/v1/secrets",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/secret/Validate_Create/post_secret.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bland.ai/v1/secrets/secret_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/secret/Validate_Create/update_secret.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://api.bland.ai/v1/secrets/secret_123",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/secrets/secret_123`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/resource/secret/Validate_Create/get_secret.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "bland_secret" "test" {
						name                              = "test_secret"
						value                       = "example secret value"
						static = true
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bland_secret.test", "name", "test_secret"),
					resource.TestCheckResourceAttr("bland_secret.test", "value", "example secret value"),
					resource.TestCheckResourceAttr("bland_secret.test", "id", "secret_123"),
				),
			},
		},
	})
}
