// Copyright (c) James Hiester.
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"sort"

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
		Nodes:       pathwayToCreate.Nodes,
		Edges:       pathwayToCreate.Edges,
	}

	return &pathway, nil
}

func (client *client) GetPathwayVersions(ctx context.Context, pathwayID string) ([]pathwayVersionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s/versions", pathwayID),
	}

	var versions []pathwayVersionDto
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &versions)
	if err != nil {
		return nil, fmt.Errorf("failed to get pathway versions: %w", err)
	}

	sort.Slice(versions, func(i, j int) bool {
		if versions[i].VersionNumber != versions[j].VersionNumber {
			return versions[i].VersionNumber > versions[j].VersionNumber
		}
		return versions[i].RevisionNumber > versions[j].RevisionNumber
	})

	return versions, nil
}

func (client *client) UpdatePathway(ctx context.Context, pathwayID string, pathwayToUpdate updatePathwayDto) (*pathwayDto, error) {
	_, err := client.GetPathway(ctx, pathwayID)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   "/convo_pathway/update",
	}

	updateResponse := updatePathwayResponseDto{}
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, pathwayToUpdate, []int{http.StatusOK}, &updateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to update connection: %w", err)
	}
	updatedPathway := pathwayDto{}
	updatedPathway.ID = pathwayID
	updatedPathway.Name = pathwayToUpdate.Name
	updatedPathway.Description = pathwayToUpdate.Description
	updatedPathway.Nodes = pathwayToUpdate.Nodes
	updatedPathway.Edges = pathwayToUpdate.Edges
	return &updatedPathway, nil
}

func (client *client) GetPathway(ctx context.Context, pathwayID string) (*pathwayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s", pathwayID),
	}

	pathway := getPathwayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &pathway)
	if err != nil {
		if strings.Contains(err.Error(), "PathwayNotFound") {
			return nil, api.WrapIntoProviderError(err, api.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("Pathway '%s' not found", pathwayID))
		}
		return nil, fmt.Errorf("failed to get pathway: %w", err)
	}

	result := pathwayDto{
		ID:          pathwayID,
		Name:        pathway.Name,
		Description: pathway.Description,
		Nodes:       pathway.Nodes,
		Edges:       pathway.Edges,
	}

	return &result, nil
}

func (client *client) DeletePathway(ctx context.Context, pathwayID string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.Config.BaseURL,
		Path:   fmt.Sprintf("/v1/pathway/%s", pathwayID),
	}

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return fmt.Errorf("failed to delete pathway: %w", err)
	}
	return nil
}
