package codereview

import (
	"strings"
	"testing"
)

// A healthy review that found nothing used to be labelled "heuristic fallback",
// because the label was set from the finding count rather than from whether a
// model ran. That produced the same "results are from deterministic checks
// only" banner as a real AI outage — which teaches readers to discount the
// banner in exactly the case it matters.
func TestDegradedBannerDistinguishesRanFromUnavailable(t *testing.T) {
	unavailable := Report{
		Reviewer:        "heuristic fallback",
		DegradedReasons: []string{"AI review unavailable (no credentials)"},
	}
	if got := RenderMarkdown(unavailable); !strings.Contains(got, "AI review unavailable") {
		t.Fatalf("a review no model ran must say so, got:\n%s", got)
	}

	// The AI ran, one shard failed, and it found nothing. "No issues found" is
	// least trustworthy here, so the banner must not claim the review was clean
	// nor presuppose findings that are absent.
	ranButThin := Report{
		Reviewer:        "heuristic+ai",
		DegradedReasons: []string{"1 of 4 parallel reviews failed"},
	}
	got := RenderMarkdown(ranButThin)
	if strings.Contains(got, "AI review unavailable") {
		t.Fatalf("a review the AI did run must not claim it was unavailable, got:\n%s", got)
	}
	if !strings.Contains(got, "ran degraded") {
		t.Fatalf("a degraded run must say it was degraded, got:\n%s", got)
	}
	if strings.Contains(got, "the findings below are real") {
		t.Fatalf("must not promise findings when there are none, got:\n%s", got)
	}

	// With findings present the original sentence is the right one.
	ranWithFindings := Report{
		Reviewer:        "heuristic+ai",
		DegradedReasons: []string{"1 of 4 parallel reviews failed"},
		Findings:        []Finding{{Title: "A real problem", Strength: "Blocking"}},
	}
	if got := RenderMarkdown(ranWithFindings); !strings.Contains(got, "the findings below are real") {
		t.Fatalf("a degraded run with findings should vouch for them, got:\n%s", got)
	}
}
