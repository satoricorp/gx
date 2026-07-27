package codereview

import "testing"

func TestFindingMentionsChangedFile(t *testing.T) {
	tests := []struct {
		name    string
		finding Finding
		changed []string
		want    bool
	}{
		{
			name:    "exact full path",
			finding: Finding{Summary: "The change in `internal/app/app.go` needs review."},
			changed: []string{"internal/app/app.go"},
			want:    true,
		},
		{
			name:    "bare basename matches changed basename",
			finding: Finding{Recommendation: "Update utils.go with a guard."},
			changed: []string{"internal/utils.go"},
			want:    true,
		},
		{
			name:    "substring filename does not match (regression)",
			finding: Finding{Summary: "Refactor myutils.go for clarity."},
			changed: []string{"utils.go"},
			want:    false,
		},
		{
			name:    "extension-collision substring does not match (regression)",
			finding: Finding{Summary: "Document app.golang somewhere."},
			changed: []string{"app.go"},
			want:    false,
		},
		{
			name:    "explicit different directory does not match",
			finding: Finding{Summary: "Look at src/utils.go for the helper."},
			changed: []string{"internal/utils.go"},
			want:    false,
		},
		{
			name:    "no file mention",
			finding: Finding{Summary: "General advice with no files."},
			changed: []string{"internal/app/app.go"},
			want:    false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := findingMentionsChangedFile(tc.finding, tc.changed); got != tc.want {
				t.Fatalf("findingMentionsChangedFile() = %v, want %v", got, tc.want)
			}
		})
	}
}
