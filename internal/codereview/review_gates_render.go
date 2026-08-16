package codereview

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

// Renders for the review report. The last two lines are a contract: a
// single "Verdict: ..." line and a "Next: ..." line, which is what the MCP
// tool and the /review slash command tell the host agent to relay.

const reviewMaxRenderedFiles = 8

func reviewColorize(report ReviewReport, paint func(string) string, text string) string {
	return colorize(report.Color, paint, text)
}

func reviewAccent(report ReviewReport, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func reviewStatusPaint(status GateStatus) func(string) string {
	switch status {
	case GatePass:
		return termstyle.Success
	case GateFail:
		return termstyle.Danger
	default:
		return termstyle.Muted
	}
}

func reviewHeaderLine(report ReviewReport) string {
	target := strings.TrimSpace(report.ReviewRange)
	if target == "" {
		target = strings.TrimSpace(report.ReviewTarget)
	}
	if target == "" {
		target = "the current change"
	}
	return fmt.Sprintf("Review check — %s (%d files, +%d/−%d)",
		target, report.DiffStats.Files, report.DiffStats.AddedLines, report.DiffStats.RemovedLines)
}

func reviewGateStatusCounts(report ReviewReport) (passed, failed, skipped int) {
	for _, gate := range report.Gates {
		switch gate.Status {
		case GatePass:
			passed++
		case GateFail:
			failed++
		default:
			skipped++
		}
	}
	return passed, failed, skipped
}

// reviewVerdictLines is the shared two-line seam, unstyled.
func reviewVerdictLines(report ReviewReport) (verdict, next string) {
	passed, failed, skipped := reviewGateStatusCounts(report)
	switch report.Verdict {
	case VerdictNothingToCheck:
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return "Verdict: NOTHING-TO-CHECK — no change found (looked at " + target + ")",
			"Next: make a change, then rerun `gx review`."
	case VerdictNoShip:
		names := make([]string, 0, failed)
		for _, id := range report.FailedGates() {
			names = append(names, string(id))
		}
		return fmt.Sprintf("Verdict: NO-SHIP — %d gate(s) failed (%s)", failed, strings.Join(names, ", ")),
			"Next: address the findings above, then rerun `gx review`."
	case VerdictDegraded:
		reason := strings.Join(report.DegradedReasons, "; ")
		return "Verdict: DEGRADED — no gate failed, but the run saw less than a healthy one would (" + reason + ")",
			"Next: fix the degradation above, then rerun `gx review` for a verdict worth shipping on."
	default:
		return fmt.Sprintf("Verdict: SHIP — all gates clear (%d passed, %d skipped)", passed, skipped),
			"Next: ship it — open the PR."
	}
}

// reviewFindingLine is the compact one-liner under a gate: location and
// claim only. The fix lives in the How to resolve section, once.
func reviewFindingLine(finding Finding) string {
	return reviewFindingLocation(finding) + strings.TrimSpace(finding.Title)
}

func reviewFindingLocation(finding Finding) string {
	switch {
	case finding.File != "" && finding.Line > 0:
		return fmt.Sprintf("%s:%d — ", finding.File, finding.Line)
	case finding.File != "":
		return finding.File + " — "
	default:
		return ""
	}
}

// reviewResolveStep is one entry of the How to resolve section: a finding
// with the gate it came from, or a degradation to repair.
type reviewResolveStep struct {
	gateTitle string
	finding   Finding
}

func reviewResolveSteps(report ReviewReport) []reviewResolveStep {
	var steps []reviewResolveStep
	// Failed gates first: the steps that flip the verdict outrank the warnings.
	for _, failing := range []bool{true, false} {
		for _, gate := range report.Gates {
			if (gate.Status == GateFail) != failing {
				continue
			}
			for _, finding := range gate.Findings {
				steps = append(steps, reviewResolveStep{gateTitle: gate.Title, finding: finding})
			}
		}
	}
	for _, reason := range report.DegradedReasons {
		steps = append(steps, reviewResolveStep{
			gateTitle: "Degraded run",
			finding: Finding{
				Title:          reason,
				Recommendation: "Restore this, then rerun `gx review` for a verdict worth shipping on.",
			},
		})
	}
	return steps
}

func reviewFilesLine(files []string) string {
	if len(files) == 0 {
		return ""
	}
	shown := files
	extra := ""
	if len(shown) > reviewMaxRenderedFiles {
		extra = fmt.Sprintf(" (+%d more)", len(shown)-reviewMaxRenderedFiles)
		shown = shown[:reviewMaxRenderedFiles]
	}
	return "files: " + strings.Join(shown, ", ") + extra
}

// RenderReviewText is the terminal render: one line per gate, findings
// indented beneath, the Verdict/Next seam at the end. Colors ride on
// report.Color and termstyle.Enabled, so --json and piped output stay plain.
func RenderReviewText(report ReviewReport) string {
	var b strings.Builder
	fmt.Fprintln(&b, reviewAccent(report, reviewHeaderLine(report)))
	fmt.Fprintln(&b)
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := reviewVerdictLines(report)
		fmt.Fprintln(&b, reviewColorize(report, termstyle.Muted, verdict))
		fmt.Fprintln(&b, next)
		return b.String()
	}
	for index, gate := range report.Gates {
		title := fmt.Sprintf("%-14s", gate.Title)
		status := fmt.Sprintf("%-8s", string(gate.Status))
		detail := gate.Summary
		if gate.Status == GateSkipped {
			detail = gate.SkipReason
		}
		fmt.Fprintf(&b, "  %d. %s %s %s\n",
			index+1,
			reviewAccent(report, title),
			reviewColorize(report, reviewStatusPaint(gate.Status), status),
			detail)
		for _, finding := range gate.Findings {
			fmt.Fprintf(&b, "%s- %s\n", strings.Repeat(" ", 26), reviewFindingLine(finding))
		}
		if report.Verbose {
			if line := reviewFilesLine(gate.Files); line != "" {
				fmt.Fprintf(&b, "%s%s\n", strings.Repeat(" ", 26), reviewColorize(report, termstyle.Muted, line))
			}
		}
	}
	if len(report.DegradedReasons) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, reviewColorize(report, termstyle.Danger, "Warning: "+strings.Join(report.DegradedReasons, "; ")))
	}
	if steps := reviewResolveSteps(report); len(steps) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, reviewAccent(report, "How to resolve"))
		fmt.Fprintln(&b)
		for index, step := range steps {
			header := fmt.Sprintf("%d. %s — %s", index+1, step.gateTitle, reviewFindingLine(step.finding))
			fmt.Fprintf(&b, "  %s\n", header)
			for _, line := range renderCodeLines(report.Color, step.finding, "     ") {
				fmt.Fprintln(&b, line)
			}
			if fix := strings.TrimSpace(step.finding.Recommendation); fix != "" {
				fmt.Fprintf(&b, "     %s %s\n", reviewColorize(report, termstyle.Success, "Fix:"), fix)
			}
		}
	}
	fmt.Fprintln(&b)
	verdict, next := reviewVerdictLines(report)
	paint := termstyle.Success
	if report.Verdict != VerdictShip {
		paint = termstyle.Danger
	}
	fmt.Fprintln(&b, reviewColorize(report, paint, verdict))
	fmt.Fprintln(&b, next)
	return b.String()
}

func reviewStatusEmoji(status GateStatus) string {
	switch status {
	case GatePass:
		return "✅ PASS"
	case GateFail:
		return "❌ FAIL"
	default:
		return "⏭️ SKIPPED"
	}
}

// RenderReviewMarkdown is what agents relay and what a PR body can carry:
// a status table, findings as bullets, and the same Verdict/Next seam in bold.
func RenderReviewMarkdown(report ReviewReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", reviewHeaderLine(report))
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := reviewVerdictLines(report)
		fmt.Fprintf(&b, "**%s**\n%s\n", verdict, next)
		return b.String()
	}
	fmt.Fprintln(&b, "| # | Gate | Status | Evidence |")
	fmt.Fprintln(&b, "|---|------|--------|----------|")
	for index, gate := range report.Gates {
		detail := gate.Summary
		if gate.Status == GateSkipped {
			detail = gate.SkipReason
		}
		fmt.Fprintf(&b, "| %d | %s | %s | %s |\n",
			index+1, gate.Title, reviewStatusEmoji(gate.Status), markdownTableCell(detail))
	}
	if steps := reviewResolveSteps(report); len(steps) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "### How to resolve")
		for index, step := range steps {
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "**%d. %s — %s**\n", index+1, step.gateTitle, reviewFindingLine(step.finding))
			switch {
			case strings.TrimSpace(step.finding.DiffHunk) != "":
				// Raw hunk text in a diff fence: agents and GitHub render the
				// +/− coloring themselves.
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "```diff\n%s\n```\n", step.finding.DiffHunk)
			case strings.TrimSpace(step.finding.CodeExcerpt) != "":
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "```%s\n%s\n```\n", fenceLanguage(step.finding.File), step.finding.CodeExcerpt)
			}
			if fix := strings.TrimSpace(step.finding.Recommendation); fix != "" {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "Fix: %s\n", fix)
			}
		}
	}
	if report.Verbose {
		var fileLines []string
		for _, gate := range report.Gates {
			if line := reviewFilesLine(gate.Files); line != "" {
				fileLines = append(fileLines, fmt.Sprintf("- %s — %s", gate.Title, line))
			}
		}
		if len(fileLines) > 0 {
			fmt.Fprintln(&b)
			fmt.Fprintln(&b, "**Files by gate**")
			for _, line := range fileLines {
				fmt.Fprintln(&b, line)
			}
		}
	}
	if len(report.DegradedReasons) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "> Warning: %s\n", strings.Join(report.DegradedReasons, "; "))
	}
	verdict, next := reviewVerdictLines(report)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "**%s**\n%s\n", verdict, next)
	return b.String()
}

func markdownTableCell(text string) string {
	text = strings.ReplaceAll(text, "|", "\\|")
	return strings.ReplaceAll(text, "\n", " ")
}
