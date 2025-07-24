// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package knowledgebase

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertFromKnowledgeBaseDto(dto KnowledgeBaseDto) KnowledgeBaseModel {
	return KnowledgeBaseModel{
		ID:            types.StringValue(dto.ID),
		Name:          types.StringValue(dto.Name),
		Description:   types.StringValue(dto.Description),
		Text:          types.StringValue(dto.Text),
		ExtractedText: types.StringPointerValue(dto.ExtractedText),
		FilePath:      types.StringNull(), // Not returned from API
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

func ConvertToCreateKnowledgeBaseDto(model KnowledgeBaseModel) (CreateKnowledgeBaseDto, error) {
	var fileBytesPtr *[]byte
	if !model.FilePath.IsNull() && model.FilePath.ValueString() != "" {
		fileBytes, err := os.ReadFile(model.FilePath.ValueString())
		if err != nil {
			return CreateKnowledgeBaseDto{}, err
		}
		fileBytesPtr = &fileBytes
	} else {
		fileBytesPtr = nil
	}
	return CreateKnowledgeBaseDto{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		File:        fileBytesPtr,
		Text:        model.Text.ValueStringPointer(),
	}, nil
}

func ConvertToUpdateKnowledgeBaseDto(model KnowledgeBaseModel) (UpdateKnowledgeBaseDto, error) {
	var fileBytesPtr *[]byte
	if !model.FilePath.IsNull() && model.FilePath.ValueString() != "" {
		fileBytes, err := os.ReadFile(model.FilePath.ValueString())
		if err != nil {
			return UpdateKnowledgeBaseDto{}, err
		}
		fileBytesPtr = &fileBytes
	} else {
		fileBytesPtr = nil
	}
	return UpdateKnowledgeBaseDto{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		File:        fileBytesPtr,
		Text:        model.Text.ValueStringPointer(),
	}, nil
}
