package codereview

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

// Renders for the gates report. The last two lines are a contract: a
// single "Verdict: ..." line and a "Next: ..." line, which is what the MCP
// tool and the /gates slash command tell the host agent to relay.

const gatesMaxRenderedFiles = 8

func gatesColorize(report GatesReport, paint func(string) string, text string) string {
	return colorize(report.Color, paint, text)
}

func gatesAccent(report GatesReport, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func gatesStatusPaint(status GateStatus) func(string) string {
	switch status {
	case GatePass:
		return termstyle.Success
	case GateFail:
		return termstyle.Danger
	default:
		return termstyle.Muted
	}
}

func gatesHeaderLine(report GatesReport) string {
	target := strings.TrimSpace(report.ReviewRange)
	if target == "" {
		target = strings.TrimSpace(report.ReviewTarget)
	}
	if target == "" {
		target = "the current change"
	}
	return fmt.Sprintf("Gates check — %s (%d files, +%d/−%d)",
		target, report.DiffStats.Files, report.DiffStats.AddedLines, report.DiffStats.RemovedLines)
}

func gateStatusCounts(report GatesReport) (passed, failed, skipped int) {
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

// gatesVerdictLines is the shared two-line seam, unstyled.
func gatesVerdictLines(report GatesReport) (verdict, next string) {
	passed, failed, skipped := gateStatusCounts(report)
	switch report.Verdict {
	case VerdictNothingToCheck:
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return "Verdict: NOTHING-TO-CHECK — no change found (looked at " + target + ")",
			"Next: make a change, then rerun `gx gates`."
	case VerdictNoShip:
		names := make([]string, 0, failed)
		for _, id := range report.FailedGates() {
			names = append(names, string(id))
		}
		return fmt.Sprintf("Verdict: NO-SHIP — %d gate(s) failed (%s)", failed, strings.Join(names, ", ")),
			"Next: address the findings above, then rerun `gx gates`."
	case VerdictDegraded:
		reason := strings.Join(report.DegradedReasons, "; ")
		return "Verdict: DEGRADED — no gate failed, but the run saw less than a healthy one would (" + reason + ")",
			"Next: fix the degradation above, then rerun `gx gates` for a verdict worth shipping on."
	default:
		return fmt.Sprintf("Verdict: SHIP — all gates clear (%d passed, %d skipped)", passed, skipped),
			"Next: ship it — open the PR."
	}
}

// gatesFindingLine is the compact one-liner under a gate: location and
// claim only. The fix lives in the How to resolve section, once.
func gatesFindingLine(finding Finding) string {
	return gatesFindingLocation(finding) + strings.TrimSpace(finding.Title)
}

func gatesFindingLocation(finding Finding) string {
	switch {
	case finding.File != "" && finding.Line > 0:
		return fmt.Sprintf("%s:%d — ", finding.File, finding.Line)
	case finding.File != "":
		return finding.File + " — "
	default:
		return ""
	}
}

// gatesResolveStep is one entry of the How to resolve section: a finding
// with the gate it came from, or a degradation to repair.
type gatesResolveStep struct {
	gateTitle string
	finding   Finding
}

func gatesResolveSteps(report GatesReport) []gatesResolveStep {
	var steps []gatesResolveStep
	// Failed gates first: the steps that flip the verdict outrank the warnings.
	for _, failing := range []bool{true, false} {
		for _, gate := range report.Gates {
			if (gate.Status == GateFail) != failing {
				continue
			}
			for _, finding := range gate.Findings {
				steps = append(steps, gatesResolveStep{gateTitle: gate.Title, finding: finding})
			}
		}
	}
	for _, reason := range report.DegradedReasons {
		steps = append(steps, gatesResolveStep{
			gateTitle: "Degraded run",
			finding: Finding{
				Title:          reason,
				Recommendation: "Restore this, then rerun `gx gates` for a verdict worth shipping on.",
			},
		})
	}
	return steps
}

func gatesFilesLine(files []string) string {
	if len(files) == 0 {
		return ""
	}
	shown := files
	extra := ""
	if len(shown) > gatesMaxRenderedFiles {
		extra = fmt.Sprintf(" (+%d more)", len(shown)-gatesMaxRenderedFiles)
		shown = shown[:gatesMaxRenderedFiles]
	}
	return "files: " + strings.Join(shown, ", ") + extra
}

// RenderGatesText is the terminal render: one line per gate, findings
// indented beneath, the Verdict/Next seam at the end. Colors ride on
// report.Color and termstyle.Enabled, so --json and piped output stay plain.
func RenderGatesText(report GatesReport) string {
	var b strings.Builder
	fmt.Fprintln(&b, gatesAccent(report, gatesHeaderLine(report)))
	fmt.Fprintln(&b)
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := gatesVerdictLines(report)
		fmt.Fprintln(&b, gatesColorize(report, termstyle.Muted, verdict))
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
			gatesAccent(report, title),
			gatesColorize(report, gatesStatusPaint(gate.Status), status),
			detail)
		for _, finding := range gate.Findings {
			fmt.Fprintf(&b, "%s- %s\n", strings.Repeat(" ", 26), gatesFindingLine(finding))
		}
		if report.Verbose {
			if line := gatesFilesLine(gate.Files); line != "" {
				fmt.Fprintf(&b, "%s%s\n", strings.Repeat(" ", 26), gatesColorize(report, termstyle.Muted, line))
			}
		}
	}
	if len(report.DegradedReasons) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, gatesColorize(report, termstyle.Danger, "Warning: "+strings.Join(report.DegradedReasons, "; ")))
	}
	if steps := gatesResolveSteps(report); len(steps) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, gatesAccent(report, "How to resolve"))
		fmt.Fprintln(&b)
		for index, step := range steps {
			header := fmt.Sprintf("%d. %s — %s", index+1, step.gateTitle, gatesFindingLine(step.finding))
			fmt.Fprintf(&b, "  %s\n", header)
			for _, line := range renderCodeLines(report.Color, step.finding, "     ") {
				fmt.Fprintln(&b, line)
			}
			if fix := strings.TrimSpace(step.finding.Recommendation); fix != "" {
				fmt.Fprintf(&b, "     %s %s\n", gatesColorize(report, termstyle.Success, "Fix:"), fix)
			}
		}
	}
	fmt.Fprintln(&b)
	verdict, next := gatesVerdictLines(report)
	paint := termstyle.Success
	if report.Verdict != VerdictShip {
		paint = termstyle.Danger
	}
	fmt.Fprintln(&b, gatesColorize(report, paint, verdict))
	fmt.Fprintln(&b, next)
	return b.String()
}

func gatesStatusEmoji(status GateStatus) string {
	switch status {
	case GatePass:
		return "✅ PASS"
	case GateFail:
		return "❌ FAIL"
	default:
		return "⏭️ SKIPPED"
	}
}

// RenderGatesMarkdown is what agents relay and what a PR body can carry:
// a status table, findings as bullets, and the same Verdict/Next seam in bold.
func RenderGatesMarkdown(report GatesReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", gatesHeaderLine(report))
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := gatesVerdictLines(report)
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
			index+1, gate.Title, gatesStatusEmoji(gate.Status), markdownTableCell(detail))
	}
	if steps := gatesResolveSteps(report); len(steps) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "### How to resolve")
		for index, step := range steps {
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "**%d. %s — %s**\n", index+1, step.gateTitle, gatesFindingLine(step.finding))
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
			if line := gatesFilesLine(gate.Files); line != "" {
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
	verdict, next := gatesVerdictLines(report)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "**%s**\n%s\n", verdict, next)
	return b.String()
}

func markdownTableCell(text string) string {
	text = strings.ReplaceAll(text, "|", "\\|")
	return strings.ReplaceAll(text, "\n", " ")
}
