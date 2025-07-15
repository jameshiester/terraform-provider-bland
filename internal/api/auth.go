// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import "github.com/jameshiester/terraform-provider-bland/internal/config"

type Auth struct {
	config *config.ProviderConfig
}

func NewAuthBase(configValue *config.ProviderConfig) *Auth {
	return &Auth{
		config: configValue,
	}
}
