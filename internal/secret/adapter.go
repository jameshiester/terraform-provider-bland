// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ConvertFromSecretDto(dto secretDto) SecretModel {
	return SecretModel{
		ID:     types.StringValue(dto.ID),
		Name:   types.StringValue(dto.Name),
		Value:  types.StringPointerValue(dto.Value),
		Static: types.BoolPointerValue(dto.Static),
		Config: ConvertFromSecretConfigDto(dto.Config),
	}
}

func ConvertToSecretConfigDto(ctx context.Context, model *SecretConfigModel) (*secretConfigDto, error) {
	if model == nil {
		return nil, nil
	}
	dto := secretConfigDto{
		URL:             model.URL.ValueString(),
		Method:          model.Method.ValueString(),
		Response:        model.Response.ValueString(),
		Body:            model.Body.ValueStringPointer(),
		RefreshInterval: model.RefreshInterval.ValueInt32(),
	}
	var headers map[string]string
	if !model.Headers.IsNull() {
		elements := make(map[string]types.String, len(model.Headers.Elements()))
		diags := model.Headers.ElementsAs(ctx, &elements, false)
		if diags.HasError() {
			return nil, errors.New("headers are not valid")
		}
		dto.Headers = &headers
	}
	return &dto, nil
}

func ConvertFromSecretConfigDto(dto *secretConfigDto) *SecretConfigModel {
	if dto == nil {
		return nil
	}
	var headers types.Map
	if dto.Headers != nil {
		headerVals := make(map[string]attr.Value)
		for k, v := range *dto.Headers {
			headerVals[k] = types.StringValue(v)
		}
		headers, _ = types.MapValue(types.StringType, headerVals)
	} else {
		headers = basetypes.NewMapNull(types.StringType)
	}

	return &SecretConfigModel{
		URL:             types.StringValue(dto.URL),
		Method:          types.StringValue(dto.Method),
		Response:        types.StringValue(dto.Response),
		Body:            types.StringPointerValue(dto.Body),
		RefreshInterval: types.Int32Value(dto.RefreshInterval),
		Headers:         headers,
	}
}

func ConvertToSecretDto(ctx context.Context, model SecretModel) (*secretDto, error) {
	dto := secretDto{
		ID:    model.ID.ValueString(),
		Name:  model.Name.ValueString(),
		Value: model.Value.ValueStringPointer(),
	}
	if model.Config != nil {
		config, err := ConvertToSecretConfigDto(ctx, model.Config)
		if err != nil {
			return nil, err
		}
		dto.Config = config
	}
	return &dto, nil
}
