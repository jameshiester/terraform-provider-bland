// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
)

type SecretClient struct {
	Api *api.Client
}

func newSecretClient(apiClient *api.Client) *SecretClient {
	return &SecretClient{Api: apiClient}
}

func (c *SecretClient) CreateSecret(ctx context.Context, secret createSecretDto) (*secretDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   "/v1/secrets",
	}
	var created createSecretResponseDto
	_, err := c.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, secret, []int{http.StatusCreated}, &created)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	createdSecret := secretDto{
		ID:    created.Data.ID,
		Name:  secret.Name,
		Value: secret.Value,
	}
	return &createdSecret, nil
}

func (c *SecretClient) ReadSecret(ctx context.Context, id string) (*secretDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/secrets/%s", id),
	}
	var secret secretDto
	_, err := c.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	return &secret, nil
}

func (c *SecretClient) UpdateSecret(ctx context.Context, secretID string, secret updateSecretDto) (*secretDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/secrets/%s", secretID),
	}
	var updated secretDto
	_, err := c.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, secret, []int{http.StatusOK}, &updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", err)
	}
	return &updated, nil
}

func (c *SecretClient) DeleteSecret(ctx context.Context, id string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   c.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/secrets/%s", id),
	}
	_, err := c.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}
