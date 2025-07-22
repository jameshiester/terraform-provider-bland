package pathways

import (
	"sort"
	"testing"
)

type testVersion struct {
	versionNumber   int
	revisionNumber  int
	isPrevPublished *bool
}

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
