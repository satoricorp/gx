package codereview

import (
	"strings"
	"testing"
)

func renderTestReport() ConstraintsReport {
	return ConstraintsReport{
		RepoRoot:    "/tmp/repo",
		Reviewed:    true,
		ReviewMode:  ReviewModeRange,
		ReviewRange: "main...HEAD",
		DiffStats:   ConstraintsDiffStats{Files: 2, AddedLines: 40, RemovedLines: 3},
		Verdict:     VerdictNoShip,
		Gates: []GateResult{
			{Gate: GateCorrectness, Title: "Correctness", Status: GatePass, Summary: "go test: ok", Files: []string{"a.go"}},
			{Gate: GateSecurity, Title: "Security", Status: GateFail, Summary: "1 possible secret(s) in added lines",
				Findings: []Finding{{Title: "Possible AWS access key in the diff", File: "a.go", Line: 3, Recommendation: "Rotate it."}},
				Files:    []string{"a.go"}},
			{Gate: GateAccessibility, Title: "Accessibility", Status: GateSkipped, SkipReason: "no UI files changed"},
		},
	}
}

func TestRenderConstraintsTextPlain(t *testing.T) {
	report := renderTestReport()
	report.Color = false
	out := RenderConstraintsText(report)
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("plain render carries ANSI escapes:\n%s", out)
	}
	for _, want := range []string{
		"Constraints check — main...HEAD (2 files, +40/−3)",
		"1. Correctness",
		"PASS",
		"FAIL",
		"SKIPPED  no UI files changed",
		"a.go:3 — Possible AWS access key in the diff — Rotate it.",
		"Verdict: NO-SHIP — 1 gate(s) failed (security)",
		"Next: address the findings above, then rerun `gx constraints`.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("render missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "files: a.go") {
		t.Fatalf("files listed without --verbose:\n%s", out)
	}
}

func TestRenderConstraintsTextVerboseListsFiles(t *testing.T) {
	report := renderTestReport()
	report.Verbose = true
	out := RenderConstraintsText(report)
	if !strings.Contains(out, "files: a.go") {
		t.Fatalf("verbose render missing the files line:\n%s", out)
	}
}

func TestRenderConstraintsTextShipVariant(t *testing.T) {
	report := renderTestReport()
	report.Verdict = VerdictShip
	report.Gates[1].Status = GatePass
	report.Gates[1].Findings = nil
	out := RenderConstraintsText(report)
	if !strings.Contains(out, "Verdict: SHIP — all gates clear (2 passed, 1 skipped)") {
		t.Fatalf("ship verdict line wrong:\n%s", out)
	}
	if !strings.Contains(out, "Next: ship it — open the PR.") {
		t.Fatalf("ship next line wrong:\n%s", out)
	}
}

func TestRenderConstraintsTextNothingToCheck(t *testing.T) {
	report := ConstraintsReport{Verdict: VerdictNothingToCheck, ReviewTarget: "the working tree"}
	out := RenderConstraintsText(report)
	if !strings.Contains(out, "Verdict: NOTHING-TO-CHECK") || !strings.Contains(out, "the working tree") {
		t.Fatalf("nothing-to-check render wrong:\n%s", out)
	}
}

func TestRenderConstraintsMarkdownTable(t *testing.T) {
	report := renderTestReport()
	out := RenderConstraintsMarkdown(report)
	for _, want := range []string{
		"## Constraints check — main...HEAD (2 files, +40/−3)",
		"| # | Gate | Status | Evidence |",
		"| 1 | Correctness | ✅ PASS | go test: ok |",
		"| 2 | Security | ❌ FAIL |",
		"| 3 | Accessibility | ⏭️ SKIPPED | no UI files changed |",
		"**Findings**",
		"- **Security** — a.go:3 — Possible AWS access key in the diff — Rotate it.",
		"**Verdict: NO-SHIP — 1 gate(s) failed (security)**",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("markdown missing %q:\n%s", want, out)
		}
	}
}

func TestRenderConstraintsDegradedWarning(t *testing.T) {
	report := renderTestReport()
	report.Verdict = VerdictDegraded
	report.Gates[1].Status = GatePass
	report.Gates[1].Findings = nil
	report.DegradedReasons = []string{"not signed in to gx Cloud"}
	text := RenderConstraintsText(report)
	if !strings.Contains(text, "Warning: not signed in to gx Cloud") {
		t.Fatalf("text render missing the degraded warning:\n%s", text)
	}
	if !strings.Contains(text, "Verdict: DEGRADED") {
		t.Fatalf("text render missing the degraded verdict:\n%s", text)
	}
	markdown := RenderConstraintsMarkdown(report)
	if !strings.Contains(markdown, "> Warning: not signed in to gx Cloud") {
		t.Fatalf("markdown render missing the degraded warning:\n%s", markdown)
	}
}
