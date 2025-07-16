// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/constants"
)

func newPathwayClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) CreatePathway(ctx context.Context, pathwayToCreate createPathwayDto) (*pathwayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().BaseURL,
		Path:   "/v1/pathway/create",
	}
	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	response := createPathwayResponseDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, pathwayToCreate, []int{http.StatusCreated}, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create pathway: %w", err)
	}

	if response.Errors != nil {
		messages := make([]string, 0, len(*response.Errors))
		for _, err := range *response.Errors {
			messages = append(messages, err.Message)
		}
		message := strings.Join(messages, ". ")
		return nil, fmt.Errorf("failed to create pathway with incorrect status: %s", message)
	}
	if response.Data == nil || response.Data.ID == "" {
		return nil, fmt.Errorf("failed to create pathway: %s", "invalid data in response")
	}
	pathway := pathwayDto{
		ID:          response.Data.ID,
		Name:        pathwayToCreate.Name,
		Description: pathwayToCreate.Description,
	}

	return &pathway, nil
}

func (client *client) UpdatePathway(ctx context.Context, pathwayID string, pathwayToUpdate updatePathwayDto) (*pathwayDto, error) {
	_, err := client.GetPathway(ctx, pathwayID)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s", pathwayID),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	updateResponse := updatePathwayResponseDto{}
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, pathwayToUpdate, []int{http.StatusOK}, &updateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to update connection: %w", err)
	}
	updatedPathway := pathwayDto{}
	updatedPathway.ID = pathwayID
	updatedPathway.Name = updateResponse.Data.Name
	updatedPathway.Description = updateResponse.Data.Description
	return &updatedPathway, nil
}

func (client *client) GetPathway(ctx context.Context, pathwayID string) (*pathwayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s", pathwayID),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	pathway := pathwayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &pathway)
	if err != nil {
		if strings.Contains(err.Error(), "PathwayNotFound") {
			return nil, api.WrapIntoProviderError(err, api.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("Pathway '%s' not found", pathwayID))
		}
		return nil, fmt.Errorf("failed to get pathway: %w", err)
	}
	return &pathway, nil
}

func (client *client) DeletePathway(ctx context.Context, pathwayID string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s", pathwayID),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete pathway: %w", err)
	}
	return nil
}
