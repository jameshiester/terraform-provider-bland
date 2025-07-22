// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package pathways

import (
	"context"
	"net/http"
	"sort"
	"testing"

	"github.com/jameshiester/terraform-provider-bland/internal/api"
	"github.com/jameshiester/terraform-provider-bland/internal/config"
	"github.com/jarcoal/httpmock"
)

func boolPtr(b bool) *bool { return &b }

func TestFindLatestUnpublishedVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   []pathwayVersionDto
		expectV int
		expectR int
		found   bool
	}{
		{
			name:    "no versions",
			input:   []pathwayVersionDto{},
			expectV: 0, expectR: 0, found: false,
		},
		{
			name: "all published",
			input: []pathwayVersionDto{
				{VersionNumber: 2, RevisionNumber: 1, IsPrevPublished: boolPtr(true)},
				{VersionNumber: 1, RevisionNumber: 2, IsPrevPublished: boolPtr(true)},
			},
			expectV: 0, expectR: 0, found: false,
		},
		{
			name: "one unpublished",
			input: []pathwayVersionDto{
				{VersionNumber: 2, RevisionNumber: 1, IsPrevPublished: boolPtr(false)},
				{VersionNumber: 1, RevisionNumber: 2, IsPrevPublished: boolPtr(true)},
			},
			expectV: 2, expectR: 1, found: true,
		},
		{
			name: "multiple unpublished, pick latest",
			input: []pathwayVersionDto{
				{VersionNumber: 3, RevisionNumber: 1, IsPrevPublished: boolPtr(false)},
				{VersionNumber: 2, RevisionNumber: 5, IsPrevPublished: boolPtr(false)},
				{VersionNumber: 2, RevisionNumber: 4, IsPrevPublished: boolPtr(false)},
			},
			expectV: 3, expectR: 1, found: true,
		},
		{
			name: "nil isPrevPublished treated as published",
			input: []pathwayVersionDto{
				{VersionNumber: 2, RevisionNumber: 1, IsPrevPublished: nil},
				{VersionNumber: 1, RevisionNumber: 2, IsPrevPublished: boolPtr(false)},
			},
			expectV: 1, expectR: 2, found: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, r, found := FindLatestUnpublishedVersion(tc.input)
			if v != tc.expectV || r != tc.expectR || found != tc.found {
				t.Errorf("expected (%d, %d, %v), got (%d, %d, %v)", tc.expectV, tc.expectR, tc.found, v, r, found)
			}
		})
	}
}

func TestGetPathwayVersions_Sorting(t *testing.T) {
	versions := []pathwayVersionDto{
		{VersionNumber: 1, RevisionNumber: 2},
		{VersionNumber: 2, RevisionNumber: 1},
		{VersionNumber: 1, RevisionNumber: 3},
		{VersionNumber: 2, RevisionNumber: 2},
	}
	// Simulate what GetPathwayVersions does
	sort.Slice(versions, func(i, j int) bool {
		if versions[i].VersionNumber != versions[j].VersionNumber {
			return versions[i].VersionNumber > versions[j].VersionNumber
		}
		return versions[i].RevisionNumber > versions[j].RevisionNumber
	})
	// Expect order: (2,2), (2,1), (1,3), (1,2)
	expect := []struct{ v, r int }{{2, 2}, {2, 1}, {1, 3}, {1, 2}}
	for i, e := range expect {
		if versions[i].VersionNumber != e.v || versions[i].RevisionNumber != e.r {
			t.Errorf("at %d: expected (%d,%d), got (%d,%d)", i, e.v, e.r, versions[i].VersionNumber, versions[i].RevisionNumber)
		}
	}
}

func TestGetPathwayVersions_HTTPMock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResp := `[
		{"version_number": 1, "revision_number": 2, "created_at": "2025-07-22T00:16:28.052Z", "name": "Version 1", "source_version_number": null, "is_prev_published": true},
		{"version_number": 2, "revision_number": 1, "created_at": "2025-07-23T00:16:28.052Z", "name": "Version 2", "source_version_number": 1, "is_prev_published": false},
		{"version_number": 1, "revision_number": 3, "created_at": "2025-07-24T00:16:28.052Z", "name": "Version 1.3", "source_version_number": null, "is_prev_published": false}
	]`

	pathwayID := "abc123"
	url := "https://api.bland.ai/v1/pathway/" + pathwayID + "/versions"
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, mockResp), nil
		},
	)

	client := client{Api: &api.Client{Config: &config.ProviderConfig{BaseURL: "api.bland.ai"}}}
	ctx := context.Background()
	versions, err := client.GetPathwayVersions(ctx, pathwayID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(versions))
	}
	// Should be sorted: version 2,1; version 1,3; version 1,2
	if versions[0].VersionNumber != 2 || versions[0].RevisionNumber != 1 {
		t.Errorf("expected first version to be (2,1), got (%d,%d)", versions[0].VersionNumber, versions[0].RevisionNumber)
	}
	if versions[1].VersionNumber != 1 || versions[1].RevisionNumber != 3 {
		t.Errorf("expected second version to be (1,3), got (%d,%d)", versions[1].VersionNumber, versions[1].RevisionNumber)
	}
	if versions[2].VersionNumber != 1 || versions[2].RevisionNumber != 2 {
		t.Errorf("expected third version to be (1,2), got (%d,%d)", versions[2].VersionNumber, versions[2].RevisionNumber)
	}
	// Check is_prev_published field
	if versions[0].IsPrevPublished == nil || *versions[0].IsPrevPublished != false {
		t.Errorf("expected is_prev_published for first version to be false")
	}
	if versions[1].IsPrevPublished == nil || *versions[1].IsPrevPublished != false {
		t.Errorf("expected is_prev_published for second version to be false")
	}
	if versions[2].IsPrevPublished == nil || *versions[2].IsPrevPublished != true {
		t.Errorf("expected is_prev_published for third version to be true")
	}
}
