package codereview

import (
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/cloud"
)

// An exact-match history row that carries its rule prints a "Rule:" line, so
// the reader can tell "this rule fired here before" apart from "something
// similar was said". Rows recorded before rules existed print no such line.
func TestReviewHistorySnippetsShowRuleWhenPayloadHasOne(t *testing.T) {
	result := cloud.CodeReviewHistorySearchResult{
		ExactMatches: []cloud.CodeReviewHistoryFindingMatch{
			{
				ID:          "row-1",
				Fingerprint: "fp-1",
				Outcome:     "valid",
				Category:    "review",
				FilePath:    "pkg/loop.go",
				Title:       "Off-by-one",
				Summary:     "Loop stops short.",
				Payload:     map[string]any{"rule_id": " review/wrong-logic "},
			},
			{
				ID:          "row-2",
				Fingerprint: "fp-2",
				Outcome:     "valid",
				Category:    "review",
				FilePath:    "pkg/other.go",
				Title:       "Old row",
				Summary:     "Predates rules.",
			},
		},
	}
	snippets := reviewHistorySnippets(result, 10)
	if len(snippets) != 2 {
		t.Fatalf("expected 2 snippets, got %d", len(snippets))
	}
	if !strings.Contains(snippets[0].Text, "\nRule: review/wrong-logic\n") {
		t.Errorf("first snippet should carry a trimmed Rule: line:\n%s", snippets[0].Text)
	}
	if strings.Contains(snippets[1].Text, "Rule:") {
		t.Errorf("row without rule_id must not print a Rule: line:\n%s", snippets[1].Text)
	}
	// The rest of the snippet is unchanged.
	for _, want := range []string{"Prior code review finding.", "Outcome: valid", "Category: review", "File: pkg/loop.go", "Title: Off-by-one", "Summary: Loop stops short."} {
		if !strings.Contains(snippets[0].Text, want) {
			t.Errorf("snippet missing %q:\n%s", want, snippets[0].Text)
		}
	}
	if snippets[0].Ref != "fp-1" || snippets[0].File != "pkg/loop.go" {
		t.Errorf("snippet identity changed: %+v", snippets[0])
	}
}
