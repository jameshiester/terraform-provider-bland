// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import "github.com/hashicorp/terraform-plugin-framework/types"

type KnowledgeBaseModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	File          types.String `tfsdk:"file"`
	Text          types.String `tfsdk:"text"`
	ExtractedText types.String `tfsdk:"extracted_text"`
}

type KnowledgeBaseDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Text          types.String `tfsdk:"text"`
	ExtractedText types.String `tfsdk:"extracted_text"`
}
