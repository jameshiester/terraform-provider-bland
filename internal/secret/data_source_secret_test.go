// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jameshiester/terraform-provider-bland/internal/mocks"
	"github.com/jarcoal/httpmock"
)

func TestUnitSecretDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bland.ai/v1/secrets/secret123`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("./tests/datasource/Validate_Read/get_secret.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "bland_secret" "secret" {
					id = "secret123"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.bland_secret.secret", "id", "secret123"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "name", "TestSecret"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "value", "secret-value"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.method", "GET"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.url", "https://api.example.com/secret"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.refresh_interval", "300"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.response", "token"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.headers.%", "2"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.headers.Authorization", "Bearer token"),
					resource.TestCheckResourceAttr("data.bland_secret.secret", "config.headers.Content-Type", "application/json"),
				),
			},
		},
	})
}
