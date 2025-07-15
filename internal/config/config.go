// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

type ProviderConfig struct {
	APIKey           string
	TerraformVersion string
	TestMode         bool
	BaseURL          string
}
