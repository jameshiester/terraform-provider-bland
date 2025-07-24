// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertFromKnowledgeBaseDto(dto KnowledgeBaseDto) KnowledgeBaseModel {
	var fileValue types.String
	if dto.File == nil {
		fileValue = types.StringNull()
	} else {
		fileValue = types.StringValue(base64.StdEncoding.EncodeToString(*dto.File))
	}
	return KnowledgeBaseModel{
		ID:            types.StringValue(dto.ID),
		Name:          types.StringValue(dto.Name),
		Description:   types.StringValue(dto.Description),
		Text:          types.StringValue(dto.Text),
		File:          fileValue,
		ExtractedText: types.StringPointerValue(dto.ExtractedText),
	}
}

func ConvertFromKnowledgeBaseDtoToDataSource(dto KnowledgeBaseDto) KnowledgeBaseDataSourceModel {
	return KnowledgeBaseDataSourceModel{
		ID:            types.StringValue(dto.ID),
		Name:          types.StringValue(dto.Name),
		Description:   types.StringValue(dto.Description),
		ExtractedText: types.StringPointerValue(dto.ExtractedText),
	}
}

func ConvertToCreateKnowledgeBaseDto(model KnowledgeBaseModel) CreateKnowledgeBaseDto {
	var fileBytesPtr *[]byte
	if !model.File.IsNull() && model.File.ValueString() != "" {
		fileBytes, _ := base64.StdEncoding.DecodeString(model.File.ValueString())
		fileBytesPtr = &fileBytes
	} else {
		fileBytesPtr = nil
	}
	return CreateKnowledgeBaseDto{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		File:        fileBytesPtr,
		Text:        model.Text.ValueStringPointer(),
	}
}

func ConvertToUpdateKnowledgeBaseDto(model KnowledgeBaseModel) UpdateKnowledgeBaseDto {
	var fileBytesPtr *[]byte
	if !model.File.IsNull() && model.File.ValueString() != "" {
		fileBytes, _ := base64.StdEncoding.DecodeString(model.File.ValueString())
		fileBytesPtr = &fileBytes
	} else {
		fileBytesPtr = nil
	}
	return UpdateKnowledgeBaseDto{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		File:        fileBytesPtr,
		Text:        model.Text.ValueStringPointer(),
	}
}
