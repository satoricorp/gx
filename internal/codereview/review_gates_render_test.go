package codereview

import (
	"strings"
	"testing"
)

func renderTestReport() ReviewReport {
	return ReviewReport{
		RepoRoot:    "/tmp/repo",
		Reviewed:    true,
		ReviewMode:  ReviewModeRange,
		ReviewRange: "main...HEAD",
		DiffStats:   ReviewDiffStats{Files: 2, AddedLines: 40, RemovedLines: 3},
		Verdict:     VerdictNoShip,
		Gates: []GateResult{
			{Gate: GateCorrectness, Title: "Correctness", Status: GatePass, Summary: "go test: ok", Files: []string{"a.go"}},
			{Gate: GateSecurity, Title: "Security", Status: GateFail, Summary: "1 possible secret(s) in added lines",
				Findings: []Finding{{
					Title: "Possible AWS access key in the diff", File: "a.go", Line: 3,
					Recommendation:   "Rotate it.",
					CodeExcerpt:      "package app\n\nconst key = \"XXXX\"",
					CodeExcerptStart: 1,
				}},
				Files: []string{"a.go"}},
			{Gate: GateAccessibility, Title: "Accessibility", Status: GateSkipped, SkipReason: "no UI files changed"},
		},
	}
}

func TestRenderReviewTextPlain(t *testing.T) {
	report := renderTestReport()
	report.Color = false
	out := RenderReviewText(report)
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("plain render carries ANSI escapes:\n%s", out)
	}
	for _, want := range []string{
		"Review check — main...HEAD (2 files, +40/−3)",
		"1. Correctness",
		"PASS",
		"FAIL",
		"SKIPPED  no UI files changed",
		"a.go:3 — Possible AWS access key in the diff",
		"How to resolve",
		"1. Security — a.go:3 — Possible AWS access key in the diff",
		"→    3 | const key = \"XXXX\"",
		"   1 | package app",
		"Fix: Rotate it.",
		"Verdict: NO-SHIP — 1 gate(s) failed (security)",
		"Next: address the findings above, then rerun `gx review`.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("render missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "files: a.go") {
		t.Fatalf("files listed without --verbose:\n%s", out)
	}
}

func TestRenderReviewTextVerboseListsFiles(t *testing.T) {
	report := renderTestReport()
	report.Verbose = true
	out := RenderReviewText(report)
	if !strings.Contains(out, "files: a.go") {
		t.Fatalf("verbose render missing the files line:\n%s", out)
	}
}

func TestRenderReviewTextShipVariant(t *testing.T) {
	report := renderTestReport()
	report.Verdict = VerdictShip
	report.Gates[1].Status = GatePass
	report.Gates[1].Findings = nil
	out := RenderReviewText(report)
	if !strings.Contains(out, "Verdict: SHIP — all gates clear (2 passed, 1 skipped)") {
		t.Fatalf("ship verdict line wrong:\n%s", out)
	}
	if !strings.Contains(out, "Next: ship it — open the PR.") {
		t.Fatalf("ship next line wrong:\n%s", out)
	}
}

func TestRenderReviewTextNothingToCheck(t *testing.T) {
	report := ReviewReport{Verdict: VerdictNothingToCheck, ReviewTarget: "the working tree"}
	out := RenderReviewText(report)
	if !strings.Contains(out, "Verdict: NOTHING-TO-CHECK") || !strings.Contains(out, "the working tree") {
		t.Fatalf("nothing-to-check render wrong:\n%s", out)
	}
}

func TestRenderReviewMarkdownTable(t *testing.T) {
	report := renderTestReport()
	out := RenderReviewMarkdown(report)
	for _, want := range []string{
		"## Review check — main...HEAD (2 files, +40/−3)",
		"| # | Gate | Status | Evidence |",
		"| 1 | Correctness | ✅ PASS | go test: ok |",
		"| 2 | Security | ❌ FAIL |",
		"| 3 | Accessibility | ⏭️ SKIPPED | no UI files changed |",
		"### How to resolve",
		"**1. Security — a.go:3 — Possible AWS access key in the diff**",
		"```go\npackage app\n\nconst key = \"XXXX\"\n```",
		"Fix: Rotate it.",
		"**Verdict: NO-SHIP — 1 gate(s) failed (security)**",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("markdown missing %q:\n%s", want, out)
		}
	}
}

func TestRenderReviewShowsDiffHunks(t *testing.T) {
	report := renderTestReport()
	report.Gates[1].Findings[0].DiffHunk = "@@ -0,0 +1,3 @@\n+package app\n+\n+const key = \"XXXX\""
	report.Gates[1].Findings[0].CodeExcerpt = ""

	text := RenderReviewText(report)
	// The hunk body renders; the @@ header does not — it is addressing for
	// tools, and the finding already names file:line. The markdown render
	// still carries the raw hunk (header included) inside its diff fence.
	for _, want := range []string{"+package app", "+const key = \"XXXX\""} {
		if !strings.Contains(text, want) {
			t.Fatalf("text render missing hunk line %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "@@ -0,0 +1,3 @@") {
		t.Fatalf("text render should drop the @@ hunk header:\n%s", text)
	}
	if md := RenderReviewMarkdown(report); !strings.Contains(md, "@@ -0,0 +1,3 @@") {
		t.Fatalf("markdown render should keep the raw hunk header inside its diff fence:\n%s", md)
	}
	if strings.Contains(text, "   1 | package app") {
		t.Fatalf("text render fell back to the file excerpt despite a hunk:\n%s", text)
	}

	markdown := RenderReviewMarkdown(report)
	if !strings.Contains(markdown, "```diff\n@@ -0,0 +1,3 @@\n+package app\n+\n+const key = \"XXXX\"\n```") {
		t.Fatalf("markdown render missing the raw diff fence:\n%s", markdown)
	}
}

func TestRenderReviewDegradedWarning(t *testing.T) {
	report := renderTestReport()
	report.Verdict = VerdictDegraded
	report.Gates[1].Status = GatePass
	report.Gates[1].Findings = nil
	report.DegradedReasons = []string{"not signed in to gx Cloud"}
	text := RenderReviewText(report)
	if !strings.Contains(text, "Warning: not signed in to gx Cloud") {
		t.Fatalf("text render missing the degraded warning:\n%s", text)
	}
	if !strings.Contains(text, "Verdict: DEGRADED") {
		t.Fatalf("text render missing the degraded verdict:\n%s", text)
	}
	if !strings.Contains(text, "Degraded run — not signed in to gx Cloud") {
		t.Fatalf("text render missing the degraded resolve step:\n%s", text)
	}
	markdown := RenderReviewMarkdown(report)
	if !strings.Contains(markdown, "> Warning: not signed in to gx Cloud") {
		t.Fatalf("markdown render missing the degraded warning:\n%s", markdown)
	}
}
