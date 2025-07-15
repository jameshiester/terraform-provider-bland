// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	test "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	pathways "github.com/jameshiester/terraform-provider-bland/internal/conversational-pathways"
	"github.com/jameshiester/terraform-provider-bland/internal/mocks"
	"github.com/jameshiester/terraform-provider-bland/internal/provider"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitBlandProviderHasChildDataSources_Basic(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		pathways.NewConversationalPathwayDataSource(),
	}
	providerInstance := provider.NewBlandProvider(context.Background())()
	datasources := providerInstance.DataSources(context.Background())

	require.Equalf(t, len(expectedDataSources), len(datasources), "Expected %d data sources, got %d", len(expectedDataSources), len(datasources))
	for _, d := range datasources {
		require.Containsf(t, expectedDataSources, d(), "Data source %+v was not expected", d())
	}
}

func TestUnitBlandProviderHasChildResources_Basic(t *testing.T) {
	expectedResources := []resource.Resource{
		pathways.NewConversationalPathwayResource(),
	}
	providerInstance := provider.NewBlandProvider(context.Background())()
	resources := providerInstance.Resources(context.Background())

	require.Equalf(t, len(expectedResources), len(resources), "Expected %d resources, got %d", len(expectedResources), len(resources))
	for _, r := range resources {
		require.Containsf(t, expectedResources, r(), "Resource %+v was not expected", r())
	}
}

func TestBlandProvider_Validate_Telementry_Optout_Is_False(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	test.Test(t, test.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []test.TestStep{
			{
				Config: `provider "bland" {
					api_key = "123"
				}`,
			},
		},
	})
}
