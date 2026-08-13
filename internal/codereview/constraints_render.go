package codereview

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

// Renders for the constraints report. The last two lines are a contract: a
// single "Verdict: ..." line and a "Next: ..." line, which is what the MCP
// tool and the /constraints slash command tell the host agent to relay.

const constraintsMaxRenderedFiles = 8

func constraintsColorize(report ConstraintsReport, paint func(string) string, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return paint(text)
}

func constraintsAccent(report ConstraintsReport, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func constraintsStatusPaint(status GateStatus) func(string) string {
	switch status {
	case GatePass:
		return termstyle.Success
	case GateFail:
		return termstyle.Danger
	default:
		return termstyle.Muted
	}
}

func constraintsHeaderLine(report ConstraintsReport) string {
	target := strings.TrimSpace(report.ReviewRange)
	if target == "" {
		target = strings.TrimSpace(report.ReviewTarget)
	}
	if target == "" {
		target = "the current change"
	}
	return fmt.Sprintf("Constraints check — %s (%d files, +%d/−%d)",
		target, report.DiffStats.Files, report.DiffStats.AddedLines, report.DiffStats.RemovedLines)
}

func constraintsGateStatusCounts(report ConstraintsReport) (passed, failed, skipped int) {
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

// constraintsVerdictLines is the shared two-line seam, unstyled.
func constraintsVerdictLines(report ConstraintsReport) (verdict, next string) {
	passed, failed, skipped := constraintsGateStatusCounts(report)
	switch report.Verdict {
	case VerdictNothingToCheck:
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return "Verdict: NOTHING-TO-CHECK — no change found (looked at " + target + ")",
			"Next: make a change, then rerun `gx constraints`."
	case VerdictNoShip:
		names := make([]string, 0, failed)
		for _, id := range report.FailedGates() {
			names = append(names, string(id))
		}
		return fmt.Sprintf("Verdict: NO-SHIP — %d gate(s) failed (%s)", failed, strings.Join(names, ", ")),
			"Next: address the findings above, then rerun `gx constraints`."
	case VerdictDegraded:
		reason := strings.Join(report.DegradedReasons, "; ")
		return "Verdict: DEGRADED — no gate failed, but the run saw less than a healthy one would (" + reason + ")",
			"Next: fix the degradation above, then rerun `gx constraints` for a verdict worth shipping on."
	default:
		return fmt.Sprintf("Verdict: SHIP — all gates clear (%d passed, %d skipped)", passed, skipped),
			"Next: ship it — open the PR."
	}
}

func constraintsFindingLine(finding Finding) string {
	var location string
	switch {
	case finding.File != "" && finding.Line > 0:
		location = fmt.Sprintf("%s:%d — ", finding.File, finding.Line)
	case finding.File != "":
		location = finding.File + " — "
	}
	line := location + strings.TrimSpace(finding.Title)
	if recommendation := strings.TrimSpace(finding.Recommendation); recommendation != "" {
		line += " — " + recommendation
	}
	return line
}

func constraintsFilesLine(files []string) string {
	if len(files) == 0 {
		return ""
	}
	shown := files
	extra := ""
	if len(shown) > constraintsMaxRenderedFiles {
		extra = fmt.Sprintf(" (+%d more)", len(shown)-constraintsMaxRenderedFiles)
		shown = shown[:constraintsMaxRenderedFiles]
	}
	return "files: " + strings.Join(shown, ", ") + extra
}

// RenderConstraintsText is the terminal render: one line per gate, findings
// indented beneath, the Verdict/Next seam at the end. Colors ride on
// report.Color and termstyle.Enabled, so --json and piped output stay plain.
func RenderConstraintsText(report ConstraintsReport) string {
	var b strings.Builder
	fmt.Fprintln(&b, constraintsAccent(report, constraintsHeaderLine(report)))
	fmt.Fprintln(&b)
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := constraintsVerdictLines(report)
		fmt.Fprintln(&b, constraintsColorize(report, termstyle.Muted, verdict))
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
			constraintsAccent(report, title),
			constraintsColorize(report, constraintsStatusPaint(gate.Status), status),
			detail)
		for _, finding := range gate.Findings {
			fmt.Fprintf(&b, "%s- %s\n", strings.Repeat(" ", 26), constraintsFindingLine(finding))
		}
		if report.Verbose {
			if line := constraintsFilesLine(gate.Files); line != "" {
				fmt.Fprintf(&b, "%s%s\n", strings.Repeat(" ", 26), constraintsColorize(report, termstyle.Muted, line))
			}
		}
	}
	if len(report.DegradedReasons) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, constraintsColorize(report, termstyle.Danger, "Warning: "+strings.Join(report.DegradedReasons, "; ")))
	}
	fmt.Fprintln(&b)
	verdict, next := constraintsVerdictLines(report)
	paint := termstyle.Success
	if report.Verdict != VerdictShip {
		paint = termstyle.Danger
	}
	fmt.Fprintln(&b, constraintsColorize(report, paint, verdict))
	fmt.Fprintln(&b, next)
	return b.String()
}

func constraintsStatusEmoji(status GateStatus) string {
	switch status {
	case GatePass:
		return "✅ PASS"
	case GateFail:
		return "❌ FAIL"
	default:
		return "⏭️ SKIPPED"
	}
}

// RenderConstraintsMarkdown is what agents relay and what a PR body can carry:
// a status table, findings as bullets, and the same Verdict/Next seam in bold.
func RenderConstraintsMarkdown(report ConstraintsReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", constraintsHeaderLine(report))
	if report.Verdict == VerdictNothingToCheck {
		verdict, next := constraintsVerdictLines(report)
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
			index+1, gate.Title, constraintsStatusEmoji(gate.Status), markdownTableCell(detail))
	}
	var findingLines []string
	for _, gate := range report.Gates {
		for _, finding := range gate.Findings {
			findingLines = append(findingLines, fmt.Sprintf("- **%s** — %s", gate.Title, constraintsFindingLine(finding)))
		}
	}
	if len(findingLines) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "**Findings**")
		for _, line := range findingLines {
			fmt.Fprintln(&b, line)
		}
	}
	if report.Verbose {
		var fileLines []string
		for _, gate := range report.Gates {
			if line := constraintsFilesLine(gate.Files); line != "" {
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
	verdict, next := constraintsVerdictLines(report)
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "**%s**\n%s\n", verdict, next)
	return b.String()
}

func markdownTableCell(text string) string {
	text = strings.ReplaceAll(text, "|", "\\|")
	return strings.ReplaceAll(text, "\n", " ")
}
