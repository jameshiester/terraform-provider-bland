// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

import "github.com/hashicorp/terraform-plugin-framework/types"

type SecretConfigModel struct {
	URL             types.String `tfsdk:"url"`
	Method          types.String `tfsdk:"method"`
	Response        types.String `tfsdk:"response"`
	Body            types.String `tfsdk:"body"`
	RefreshInterval types.Int32  `tfsdk:"refresh_interval"`
	Headers         types.Map    `tfsdk:"headers"`
}

type SecretModel struct {
	ID     types.String       `tfsdk:"id"`
	Name   types.String       `tfsdk:"name"`
	Static types.Bool         `tfsdk:"static"`
	Value  types.String       `tfsdk:"value"`
	Config *SecretConfigModel `tfsdk:"config"`
}

type SecretDataSourceModel struct {
	ID     types.String       `tfsdk:"id"`
	Name   types.String       `tfsdk:"name"`
	Static types.Bool         `tfsdk:"static"`
	Value  types.String       `tfsdk:"value"`
	Config *SecretConfigModel `tfsdk:"config"`
}
