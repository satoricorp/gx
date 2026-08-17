package codereview

import (
	"strings"
	"testing"
)

func fixPromptReport() Report {
	return Report{
		Reviewed: true,
		Prompt:   "make checkout survive gateway blips",
		Findings: []Finding{
			// Advisory, but listed first: the plan must still lead with the
			// blocking one below.
			{
				ID: "ai.review.1", Title: "Name the retry budget",
				Lane: LaneAdvisory, RuleID: "gx:recommended/name-the-budget",
				Summary: "The literal 3 appears twice.", Recommendation: "Extract maxRetries.",
				File: "internal/checkout/retry.go", Line: 12,
			},
			{
				ID: "ai.review.2", Title: "Retry can double-charge the customer",
				Lane: LaneBlocking, RuleID: "gx:recommended/no-double-charge",
				Summary:        "The retry re-POSTs the charge without an idempotency key.",
				Benefit:        "A gateway blip bills the customer twice.",
				Recommendation: "Send the order ID as the Idempotency-Key header.",
				Example:        "-  post(chargeURL, body)\n+  post(chargeURL, body, idempotencyKey(order.ID))",
				DiffHunk:       "+\tresp, err := post(chargeURL, body)",
				File:           "internal/checkout/charge.go", Line: 88,
				Corroboration: []string{"bedrock-a", "bedrock-b"},
			},
			// No recommendation: not actionable, so never the top fix.
			{ID: "ai.review.3", Title: "Consider a comment", Lane: LaneAdvisory},
		},
	}
}

// The whole contract in one test: enhance leads with the blocking finding even
// when an advisory one is listed first, because it reads the review's own plan
// rather than ranking again.
func TestTopFixLeadsWithTheBlockingFinding(t *testing.T) {
	fix, ok := TopFix(fixPromptReport())
	if !ok {
		t.Fatal("TopFix() found nothing, want the blocking finding")
	}
	if fix.Finding.ID != "ai.review.2" {
		t.Fatalf("top fix = %q, want the blocking ai.review.2", fix.Finding.ID)
	}
	if !fix.Step.FlipsVerdict {
		t.Error("the blocking fix should be marked as flipping the verdict")
	}
	// One other actionable finding remains (the third has no recommendation
	// and is not a step at all), and none of the rest block.
	if fix.Remaining != 1 || fix.RemainingBlocking != 0 {
		t.Errorf("Remaining = %d (%d blocking), want 1 (0 blocking)", fix.Remaining, fix.RemainingBlocking)
	}
}

func TestTopFixReportsNothingWhenThereIsNothingActionable(t *testing.T) {
	// Findings with no recommendation are not fixes; a review of them has
	// nothing to hand a model.
	r := Report{Reviewed: true, Findings: []Finding{
		{ID: "ai.review.1", Title: "Looks fine", Lane: LaneAdvisory},
	}}
	if _, ok := TopFix(r); ok {
		t.Fatal("TopFix() reported a fix for a report with no actionable finding")
	}
	if _, ok := TopFix(Report{}); ok {
		t.Fatal("TopFix() reported a fix for an empty report")
	}
}

func TestRenderFixPromptCarriesEverythingAModelNeeds(t *testing.T) {
	report := fixPromptReport()
	fix, ok := TopFix(report)
	if !ok {
		t.Fatal("TopFix() found nothing")
	}
	out := RenderFixPrompt(report, fix)

	for _, want := range []string{
		"# Retry can double-charge the customer",            // the problem, as the title
		"internal/checkout/charge.go:88",                    // where
		"gx:recommended/no-double-charge",                   // which rule
		"blocking — this is what fails the review",          // what it costs
		"2 reviewers agreed",                                // how much to trust it
		"make checkout survive gateway blips",               // what the change was for
		"## What's wrong",                                   //
		"re-POSTs the charge without an idempotency",        //
		"## Why it matters",                                 //
		"bills the customer twice",                          //
		"## The code",                                       //
		"```diff",                                           // a diff hunk is fenced as diff
		"resp, err := post(chargeURL, body)",                //
		"## Change to make",                                 //
		"Send the order ID as the Idempotency-Key",          //
		"idempotencyKey(order.ID)",                          // the model's own sketch
		"## Done when",                                      //
		"no longer trips `gx:recommended/no-double-charge`", // a concrete finish line
		"1 more fix left after this one",                    // honest about what it hides
	} {
		if !strings.Contains(out, want) {
			t.Errorf("prompt missing %q\n---\n%s", want, out)
		}
	}
	// It is for pasting into a model: no escape codes, ever.
	if strings.Contains(out, "\x1b[") {
		t.Error("the fix prompt must never contain ANSI")
	}
	// It shows ONE fix: the other findings' titles must not appear.
	if strings.Contains(out, "Name the retry budget") {
		t.Error("the prompt should carry one fix, not the whole review")
	}
}

// A finding whose line the change did not touch has no diff hunk; the source
// excerpt stands in, and must not be fenced as a diff.
func TestRenderFixPromptFencesSourceExcerptsAsSource(t *testing.T) {
	report := Report{Reviewed: true, Findings: []Finding{{
		ID: "ai.review.1", Title: "Unbounded read", Lane: LaneBlocking,
		Recommendation: "Bound it.", File: "internal/io/read.go", Line: 4,
		CodeExcerpt: "data, _ := io.ReadAll(r)", CodeExcerptStart: 3,
	}}}
	fix, ok := TopFix(report)
	if !ok {
		t.Fatal("TopFix() found nothing")
	}
	out := RenderFixPrompt(report, fix)
	if !strings.Contains(out, "```go") {
		t.Errorf("a source excerpt should be fenced by language, got:\n%s", out)
	}
	if strings.Contains(out, "```diff") {
		t.Error("a source excerpt must not be fenced as a diff — a model would reproduce +/- markers")
	}
	// Nothing else remains, and the prompt says so rather than implying more.
	if !strings.Contains(out, "This is the only fix this review found.") {
		t.Errorf("a single-finding review should say so, got:\n%s", out)
	}
}

// The counts are the honesty check: a prompt showing one fix while blocking
// findings wait must say how many.
func TestRenderFixPromptCountsRemainingBlockingWork(t *testing.T) {
	report := Report{Reviewed: true, Findings: []Finding{
		{ID: "a", Title: "First", Lane: LaneBlocking, Recommendation: "Do a.", File: "a.go", Line: 1},
		{ID: "b", Title: "Second", Lane: LaneBlocking, Recommendation: "Do b.", File: "b.go", Line: 1},
		{ID: "c", Title: "Third", Lane: LaneBlocking, Recommendation: "Do c.", File: "c.go", Line: 1},
	}}
	fix, ok := TopFix(report)
	if !ok {
		t.Fatal("TopFix() found nothing")
	}
	if fix.RemainingBlocking != 2 {
		t.Fatalf("RemainingBlocking = %d, want 2", fix.RemainingBlocking)
	}
	out := RenderFixPrompt(report, fix)
	if !strings.Contains(out, "2 more fixes left after this one, 2 still blocking") {
		t.Errorf("prompt should count the blocking work it is not showing, got:\n%s", out)
	}
	// With blocking work left, fixing this one does not clear the gate, and
	// the acceptance line must not claim it does.
	if strings.Contains(out, "reports no blocking findings") {
		t.Error("acceptance must not promise a clean gate while other blocking findings remain")
	}
}

// A deterministic finding points at a failing tool, not at a line: it has no
// diff hunk and no excerpt, and its evidence is the only actionable thing in
// it. Measured on a real run before this: the prompt said "eslint failed" and
// gave the model nothing to act on.
func TestRenderFixPromptCarriesToolOutputWhenThereIsNoCode(t *testing.T) {
	report := Report{Reviewed: true, Findings: []Finding{{
		ID: "tools.static-failure", Title: "Fix static tool failures", Lane: LaneBlocking,
		Summary:        "`eslint .` failed.",
		Recommendation: "Fix the failing tool output, then rerun.",
		Evidence: []Evidence{
			{Label: "eslint", Value: "convex/schema.ts:14:3  error  'v' is defined but never used"},
		},
	}}}
	fix, ok := TopFix(report)
	if !ok {
		t.Fatal("TopFix() found nothing")
	}
	out := RenderFixPrompt(report, fix)
	for _, want := range []string{
		"## The failing output",
		"eslint:",
		"'v' is defined but never used",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("prompt missing %q\n---\n%s", want, out)
		}
	}
	// With no rule and no line, the acceptance line must not point at code
	// that was never shown.
	if strings.Contains(out, "the code above") {
		t.Errorf("acceptance should not reference code that is not in the prompt:\n%s", out)
	}
	if !strings.Contains(out, "- The problem described above is fixed.") {
		t.Errorf("want the location-free acceptance line, got:\n%s", out)
	}
}
