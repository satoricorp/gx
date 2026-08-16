package codereview

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/termstyle"
)

// The terminal render for gx review.
//
// gx review has printed markdown to stdout since the day it shipped — the
// same string it posts as the PR comment. This is the render a terminal
// actually deserves: a run ledger that says what ran, findings in two lanes
// with the code under discussion and a one-line fix, a fix plan the agent can
// work top to bottom, a story of what changed that the reader now owns, and
// then the two-line Verdict/Next seam that the MCP tool and the slash commands
// already relay — kept last, kept the same shape.
//
// Two audiences read one report. The agent reads the lanes and the fix plan
// and converges to clean. The human reads the story and the verdict. Nothing
// in the story is wrong; every item there ends in ownership, not action.
//
// Color rides on report.Color and termstyle.Enabled(), so --json, --md, and
// piped output stay plain — the same discipline the constraints render keeps.

const (
	reviewRenderIndent     = "  "
	reviewRenderCodeIndent = "      "
	reviewRenderRuleWidth  = 74
)

// RenderReviewText is the terminal render. Sections appear in a fixed order
// and a section with nothing in it is omitted rather than printed empty, with
// one exception: the verdict seam always prints, because agents key on it.
func RenderReviewText(report Report) string {
	var b strings.Builder
	color := report.Color && termstyle.Enabled()

	writeReviewLedger(&b, report, color)

	blocking, advisory := splitFindingsByLane(report.Findings)
	if len(blocking) > 0 {
		b.WriteString("\n")
		writeReviewLane(&b, report, color, "BLOCKING", termstyle.Danger, blocking)
	}
	if len(advisory) > 0 {
		b.WriteString("\n")
		writeReviewLane(&b, report, color, "ADVISORY", termstyle.Warning, advisory)
	}
	if plan := BuildFixPlan(report.Findings); len(plan) > 0 {
		b.WriteString("\n")
		writeReviewFixPlan(&b, report, color, plan)
	}
	if len(report.Story) > 0 {
		b.WriteString("\n")
		writeReviewStory(&b, report, color, report.Story)
	}
	if !report.Reviewed || len(report.Findings) == 0 && len(report.Story) == 0 {
		writeReviewNothingToShow(&b, report, color)
	}

	b.WriteString("\n")
	verdict, next := ReviewVerdictLines(report)
	paint := termstyle.Success
	if len(blocking) > 0 || !report.Reviewed {
		paint = termstyle.Danger
	} else if len(report.DegradedReasons) > 0 {
		paint = termstyle.Warning
	}
	b.WriteString(colorize(color, paint, verdict))
	b.WriteString("\n")
	b.WriteString(next)
	b.WriteString("\n")
	return b.String()
}

// ReviewVerdictLines is the two-line machine seam: a "Verdict: ..." line and a
// "Next: ..." line, the last two lines of the render. It mirrors the
// constraints seam so anything that already relays one relays the other.
func ReviewVerdictLines(report Report) (verdict, next string) {
	if !report.Reviewed {
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return "Verdict: NOTHING-TO-REVIEW — no change found (looked at " + target + ")",
			"Next: make a change, then rerun `gx review`."
	}
	blocking, advisory := splitFindingsByLane(report.Findings)
	if len(blocking) > 0 {
		names := make([]string, 0, len(blocking))
		seen := map[string]struct{}{}
		for _, f := range blocking {
			name := reviewRuleShortName(f)
			if _, dup := seen[name]; dup || name == "" {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
		detail := ""
		if len(names) > 0 {
			detail = " (" + strings.Join(names, ", ") + ")"
		}
		return fmt.Sprintf("Verdict: NO-SHIP — %d blocking finding(s)%s", len(blocking), detail),
			"Next: work the fix plan, then rerun `gx review`."
	}
	if len(report.DegradedReasons) > 0 {
		return "Verdict: DEGRADED — no blocking finding, but the review saw less than a healthy one would (" + strings.Join(report.DegradedReasons, "; ") + ")",
			"Next: fix the degradation above, then rerun `gx review` for a verdict worth shipping on."
	}
	if len(advisory) > 0 {
		return fmt.Sprintf("Verdict: SHIP — no blocking finding (%d advisory)", len(advisory)),
			"Next: ship it — the advisory items are yours to take or leave."
	}
	return "Verdict: SHIP — no findings",
		"Next: ship it."
}

// ---- ledger ---------------------------------------------------------------

func writeReviewLedger(b *strings.Builder, report Report, color bool) {
	writeReviewLedgerHeader(b, report, color)
	writeReviewLedgerRows(b, report, color)
}

func writeReviewLedgerHeader(b *strings.Builder, report Report, color bool) {
	target := strings.TrimSpace(report.ReviewRange)
	if target == "" {
		target = strings.TrimSpace(report.ReviewTarget)
	}
	if target == "" {
		target = "the current change"
	}
	head := colorize(color, reviewMint, "gx review")
	b.WriteString(head)
	b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf(" — %s · %d file(s)", target, len(report.ChangedFiles))))
	b.WriteString("\n")
	if prompt := strings.TrimSpace(report.Prompt); prompt != "" {
		b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf("intent: %q", prompt)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func writeReviewLedgerRows(b *strings.Builder, report Report, color bool) {
	n := 0
	row := func(label, value string) {
		n++
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf("%d  %-11s", n, label)))
		b.WriteString(value)
		b.WriteString("\n")
	}
	// The scope row states, in words, what was reviewed — "Reviewed the
	// repository", "Reviewed `origin/main...HEAD`" — because that sentence is
	// the contract a reader (and the tests) hold the report to. The file count
	// rides along as detail.
	scope := "Reviewed " + strings.TrimSpace(report.ReviewTarget)
	if strings.TrimSpace(report.ReviewTarget) == "" {
		scope = "Reviewed the current change"
	}
	if n := len(report.ChangedFiles); n > 0 {
		scope += colorize(color, termstyle.Muted, fmt.Sprintf(" · %d file(s)", n))
	}
	if !report.Reviewed {
		scope = colorize(color, termstyle.Warning, "Nothing to review") +
			colorize(color, termstyle.Muted, ": no changes found in "+strings.TrimSpace(report.ReviewTarget)+". No code was inspected, so this is not a clean review.")
	}
	row("scope", scope)

	models := "no model ran"
	if report.aiReviewRan() {
		models = strings.Join(report.ReviewModels, " · ")
		if models == "" {
			models = "AI panel"
		}
		if report.ReviewTransport != "" {
			models += colorize(color, termstyle.Muted, " via "+report.ReviewTransport)
		}
	}
	row("reviewers", models)

	if len(report.Evidence) > 0 {
		var have, missing []string
		for _, e := range report.Evidence {
			label := strings.TrimSpace(e.Source)
			if label == "" {
				continue
			}
			switch e.State {
			case EvidenceOK:
				have = append(have, label)
			case EvidenceEmpty, EvidenceDisabled, EvidenceSkipped:
				// Configured and answered, or deliberately not asked: not a
				// gap the reader needs to chase.
			default:
				missing = append(missing, label)
			}
		}
		line := ""
		if len(have) > 0 {
			line = colorize(color, termstyle.Success, strings.Join(have, " · "))
		}
		if len(missing) > 0 {
			if line != "" {
				line += "   "
			}
			line += colorize(color, termstyle.Muted, "missing: "+strings.Join(missing, ", "))
		}
		if line != "" {
			row("evidence", line)
		}
	}

	blocking, advisory := splitFindingsByLane(report.Findings)
	demoted := 0
	for _, f := range report.Findings {
		if f.DemotedFrom != "" {
			demoted++
		}
	}
	verify := fmt.Sprintf("%d blocking · %d advisory", len(blocking), len(advisory))
	if demoted > 0 {
		verify += fmt.Sprintf(" · %d demoted", demoted)
	}
	row("findings", verify)

	if report.Coverage.Partial() {
		if s := strings.TrimSpace(report.Coverage.Statement()); s != "" {
			row("coverage", colorize(color, termstyle.Warning, s))
		}
	}
	if len(report.DegradedReasons) > 0 {
		b.WriteString("\n")
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Warning, "Warning: "+strings.Join(report.DegradedReasons, "; ")))
		b.WriteString("\n")
	}
}

// ---- lanes ----------------------------------------------------------------

func splitFindingsByLane(findings []Finding) (blocking, advisory []Finding) {
	for _, f := range findings {
		if LaneOf(f) == LaneBlocking {
			blocking = append(blocking, f)
		} else {
			advisory = append(advisory, f)
		}
	}
	return blocking, advisory
}

func writeReviewLane(b *strings.Builder, report Report, color bool, label string, paint func(string) string, findings []Finding) {
	b.WriteString(reviewDivider(color, label, paint, fmt.Sprintf("%d", len(findings))))
	b.WriteString("\n")
	for _, f := range findings {
		b.WriteString("\n")
		writeReviewFinding(b, report, color, f)
	}
}

func reviewDivider(color bool, label string, paint func(string) string, right string) string {
	fill := reviewRenderRuleWidth - len(label) - len(right)
	if fill < 4 {
		fill = 4
	}
	return colorize(color, paint, label) + " " +
		colorize(color, termstyle.Muted, strings.Repeat("─", fill)+" "+right+" ──")
}

func writeReviewFinding(b *strings.Builder, report Report, color bool, f Finding) {
	// Requirement line: the rule's plain-English name when there is one,
	// otherwise the finding's own title.
	title := strings.TrimSpace(f.Title)
	b.WriteString(reviewRenderIndent)
	b.WriteString(colorize(color, termstyle.Section, title))
	b.WriteString("\n")

	// Identity line: namespace/rule · file:line · how it was decided.
	b.WriteString(reviewRenderIndent)
	if f.RuleID != "" {
		ns, name, _ := strings.Cut(f.RuleID, "/")
		b.WriteString(colorize(color, reviewMint, ns))
		b.WriteString(colorize(color, termstyle.Muted, "/"))
		b.WriteString(colorize(color, termstyle.Section, name))
	} else {
		b.WriteString(colorize(color, termstyle.Muted, f.ID))
	}
	if loc := findingLocation(f); loc != "" {
		b.WriteString(colorize(color, termstyle.Muted, "  ·  "+loc))
	}
	b.WriteString(colorize(color, termstyle.Muted, "  ["+reviewDecidedBy(f)+"]"))
	b.WriteString("\n")

	if lines := renderCodeLines(color, f, reviewRenderCodeIndent); len(lines) > 0 {
		b.WriteString("\n")
		for _, line := range lines {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	if why := strings.TrimSpace(f.Summary); why != "" {
		b.WriteString("\n")
		writeReviewLabeled(b, color, termstyle.Warning, "Why", why)
	}
	if fix := strings.TrimSpace(f.Recommendation); fix != "" {
		writeReviewLabeled(b, color, termstyle.Success, "Fix", fix)
	}

	// Evidence footer: what the trust machinery said.
	if ev := reviewEvidenceLine(f); ev != "" {
		b.WriteString("\n")
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Muted, ev))
		b.WriteString("\n")
	}
	if f.RuleID != "" && f.File != "" {
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Muted, "Suppress ▸ paste into REVIEW.md ## Exceptions:"))
		b.WriteString("\n")
		b.WriteString(reviewRenderIndent + reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Muted, SuppressLine(f)))
		b.WriteString("\n")
	}
}

// writeReviewLabeled writes "  Why  text" with continuation lines aligned
// under the text, wrapping at a comfortable width.
func writeReviewLabeled(b *strings.Builder, color bool, paint func(string) string, label, text string) {
	const width = 72
	pad := strings.Repeat(" ", len(reviewRenderIndent)+len(label)+2)
	first := true
	for _, line := range wrapText(text, width) {
		if first {
			b.WriteString(reviewRenderIndent)
			b.WriteString(colorize(color, paint, label))
			b.WriteString("  ")
			first = false
		} else {
			b.WriteString(pad)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
}

func reviewEvidenceLine(f Finding) string {
	var parts []string
	switch n := len(f.Corroboration); {
	case n >= 2:
		parts = append(parts, "both graders agreed")
	case n == 1:
		parts = append(parts, "one grader flagged")
	}
	verdict := strings.ToLower(strings.TrimSpace(f.JudgeVerdict))
	if verdict == "confirmed" {
		if f.JudgeConfidence > 0 {
			parts = append(parts, fmt.Sprintf("judge confirmed %.2f", f.JudgeConfidence))
		} else {
			parts = append(parts, "judge confirmed")
		}
	} else if verdict != "" {
		parts = append(parts, "judge "+verdict)
	}
	if f.DemotedFrom != "" {
		parts = append(parts, "demoted from "+f.DemotedFrom)
	}
	return strings.Join(parts, " · ")
}

// reviewDecidedBy names how a finding was decided: a rule with no reviewer leg
// behind it came from code; anything the panel raised was graded.
func reviewDecidedBy(f Finding) string {
	if len(f.Corroboration) == 0 && f.JudgeVerdict == "" {
		return "det"
	}
	return "graded"
}

func reviewRuleShortName(f Finding) string {
	if f.RuleID == "" {
		return ""
	}
	_, name, ok := strings.Cut(f.RuleID, "/")
	if !ok {
		return f.RuleID
	}
	return name
}

func findingLocation(f Finding) string {
	switch {
	case f.File != "" && f.Line > 0:
		return fmt.Sprintf("%s:%d", f.File, f.Line)
	case f.File != "":
		return f.File
	default:
		return ""
	}
}

// SuppressLine is the exact, paste-ready exception a reader adds to REVIEW.md
// to silence this rule where the finding landed. Friction belongs in the
// decision to suppress, not in the mechanics of doing it.
func SuppressLine(f Finding) string {
	name := reviewRuleShortName(f)
	scope := f.File
	if scope == "" {
		scope = "<path>"
	}
	return fmt.Sprintf("- %q doesn't apply in %s — <reason>.", name, scope)
}

// ---- fix plan -------------------------------------------------------------

func writeReviewFixPlan(b *strings.Builder, report Report, color bool, plan []FixStep) {
	blocking, advisory := 0, 0
	for _, s := range plan {
		if s.Lane == LaneBlocking {
			blocking++
		} else {
			advisory++
		}
	}
	b.WriteString(colorize(color, reviewMint, "FIX PLAN"))
	b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf(" ──── %d blocking · %d advisory ── fix, then rerun ", blocking, advisory)))
	b.WriteString(colorize(color, termstyle.Command, "gx review"))
	b.WriteString(colorize(color, termstyle.Muted, " ────"))
	b.WriteString("\n\n")
	writeReviewFixPlanBody(b, report, color, plan)
}

func writeReviewFixPlanBody(b *strings.Builder, report Report, color bool, plan []FixStep) {
	for _, s := range plan {
		paint := termstyle.Warning
		if s.Lane == LaneBlocking {
			paint = termstyle.Danger
		}
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, paint, fmt.Sprintf("%d.", s.Order)))
		b.WriteString(" ")
		b.WriteString(s.Action)
		b.WriteString("\n")
		meta := s.RuleID
		if meta == "" {
			meta = s.FindingID
		}
		if _, name, ok := strings.Cut(meta, "/"); ok {
			meta = name
		}
		if s.File != "" {
			meta += " · " + s.File
			if s.Line > 0 {
				meta += fmt.Sprintf(":%d", s.Line)
			}
		}
		b.WriteString(reviewRenderIndent + "   ")
		b.WriteString(colorize(color, termstyle.Muted, meta))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(reviewRenderIndent)
	b.WriteString(colorize(color, termstyle.Muted, "re-runs re-grade only the hunks you touch — the loop is cheap"))
	b.WriteString("\n")
}

// ---- story ----------------------------------------------------------------

func writeReviewStory(b *strings.Builder, report Report, color bool, items []StoryItem) {
	b.WriteString(colorize(color, termstyle.Command, "WORTH KNOWING"))
	b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf(" ────── %d change(s) that passed, and that you now own ──", len(items))))
	b.WriteString("\n")
	writeReviewStoryBody(b, report, color, items)
}

func writeReviewStoryBody(b *strings.Builder, report Report, color bool, items []StoryItem) {
	for i, item := range items {
		b.WriteString("\n")
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Command, fmt.Sprintf("%d", i+1)))
		b.WriteString("  ")
		b.WriteString(colorize(color, termstyle.Section, strings.TrimSpace(item.Headline)))
		if m := strings.TrimSpace(item.Materiality); m != "" {
			b.WriteString(colorize(color, termstyle.Muted, "   materiality: "+m))
		}
		b.WriteString("\n")
		if loc := storyLocation(item); loc != "" {
			b.WriteString(reviewRenderIndent + "   ")
			b.WriteString(colorize(color, termstyle.Muted, loc))
			b.WriteString("\n")
		}
		if hunk := strings.TrimSpace(item.DiffHunk); hunk != "" {
			b.WriteString("\n")
			for _, line := range renderHunkLines(color, hunk, reviewRenderCodeIndent) {
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
		if c := strings.TrimSpace(item.Consequence); c != "" {
			b.WriteString("\n")
			pad := reviewRenderIndent + "   "
			label := "What changes for you: "
			// Wrap the whole sentence with the label counted in, so the first
			// line does not overrun the width the rest of the report keeps.
			wrapped := wrapText(label+c, 72)
			for i, line := range wrapped {
				b.WriteString(pad)
				if i == 0 {
					line = colorize(color, termstyle.Section, label) + strings.TrimPrefix(line, label)
				}
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
		if item.PassedRule != "" {
			b.WriteString(reviewRenderIndent + "   ")
			b.WriteString(colorize(color, reviewMint, reviewRuleShortNameOf(item.PassedRule)+" passed on this"))
			b.WriteString("\n")
		}
		if item.Asked != "" || item.Chose != "" || len(item.Watch) > 0 {
			b.WriteString("\n")
		}
		if a := strings.TrimSpace(item.Asked); a != "" {
			writeStoryKV(b, color, "asked", fmt.Sprintf("%q", a))
		}
		if c := strings.TrimSpace(item.Chose); c != "" {
			writeStoryKV(b, color, "chose", c)
		}
		if len(item.Watch) > 0 {
			writeStoryKV(b, color, "watch", strings.Join(item.Watch, " · "))
		}
	}
}

func writeStoryKV(b *strings.Builder, color bool, key, value string) {
	b.WriteString(reviewRenderIndent + "   ")
	b.WriteString(colorize(color, termstyle.Muted, key))
	b.WriteString("  ")
	b.WriteString(value)
	b.WriteString("\n")
}

func storyLocation(item StoryItem) string {
	if item.File == "" {
		return ""
	}
	if len(item.Files) > 1 {
		return fmt.Sprintf("%s (+%d more)", item.File, len(item.Files)-1)
	}
	return item.File
}

func reviewRuleShortNameOf(ruleID string) string {
	if _, name, ok := strings.Cut(ruleID, "/"); ok {
		return name
	}
	return ruleID
}

// ---- empty state ----------------------------------------------------------

func writeReviewNothingToShow(b *strings.Builder, report Report, color bool) {
	b.WriteString("\n")
	b.WriteString(reviewRenderIndent)
	switch {
	case !report.Reviewed:
		b.WriteString(colorize(color, termstyle.Muted, "Nothing to review — no code was inspected."))
	case strings.TrimSpace(report.NoFindingsMessage) != "":
		b.WriteString(colorize(color, termstyle.Muted, strings.TrimSpace(report.NoFindingsMessage)))
	default:
		b.WriteString(colorize(color, termstyle.Muted, "No findings."))
	}
	b.WriteString("\n")
}

// ---- small utilities ------------------------------------------------------

// reviewMint paints with the review accent, in the same function shape as the
// termstyle painters so it can be handed to colorize.
func reviewMint(text string) string {
	return reviewMintANSI + text + reviewResetANSI
}

// wrapText breaks text at spaces to fit width, never splitting a word.
func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) > width {
			lines = append(lines, cur)
			cur = w
			continue
		}
		cur += " " + w
	}
	return append(lines, cur)
}

// ---- section-level exports for the interactive shape ---------------------
//
// The accordion in internal/cli renders one section at a time under its own
// node. It calls these rather than re-implementing the sections, so the
// interactive and linear renders cannot drift: whatever RenderReviewText
// prints for a lane, the accordion prints for that lane.

// SplitFindingsByLane partitions findings into the blocking and advisory
// lanes, preserving report order within each.
func SplitFindingsByLane(findings []Finding) (blocking, advisory []Finding) {
	return splitFindingsByLane(findings)
}

// RuleShortName is the part of a rule ID after the namespace — the name a
// human reads and a suppression references.
func RuleShortName(f Finding) string {
	return reviewRuleShortName(f)
}

// ReviewMint paints with the review accent, in the same shape as the termstyle
// painters so a caller can pick it by severity alongside them.
func ReviewMint(text string) string {
	return reviewMint(text)
}

// RenderReviewLedgerRows returns the ledger rows — scope, reviewers, evidence,
// findings, coverage — one per line, without the header line above them.
func RenderReviewLedgerRows(report Report, color bool) string {
	var b strings.Builder
	writeReviewLedgerRows(&b, report, color)
	return b.String()
}

// RenderReviewLane renders one lane's findings, without the lane divider —
// the accordion draws its own section header.
func RenderReviewLane(report Report, color bool, label string, findings []Finding) string {
	var b strings.Builder
	for i, f := range findings {
		if i > 0 {
			b.WriteString("\n")
		}
		writeReviewFinding(&b, report, color, f)
	}
	return b.String()
}

// RenderReviewFixPlan renders the fix plan section body.
func RenderReviewFixPlan(report Report, color bool) string {
	plan := BuildFixPlan(report.Findings)
	if len(plan) == 0 {
		return ""
	}
	var b strings.Builder
	writeReviewFixPlanBody(&b, report, color, plan)
	return b.String()
}

// RenderReviewStory renders the story section body.
func RenderReviewStory(report Report, color bool) string {
	if len(report.Story) == 0 {
		return ""
	}
	var b strings.Builder
	writeReviewStoryBody(&b, report, color, report.Story)
	return b.String()
}

// RenderReviewDetails renders the run details a reader opens on demand:
// models, transport, every evidence source with its state, coverage, and any
// degradation — the ledger's long form.
func RenderReviewDetails(report Report, color bool) string {
	var b strings.Builder
	kv := func(k, v string) {
		b.WriteString(reviewRenderIndent)
		b.WriteString(colorize(color, termstyle.Muted, fmt.Sprintf("%-12s", k)))
		b.WriteString(v)
		b.WriteString("\n")
	}
	if report.aiReviewRan() {
		kv("reviewers", strings.Join(report.ReviewModels, " · "))
		if report.ReviewTransport != "" {
			kv("transport", report.ReviewTransport)
		}
	} else {
		kv("reviewers", "no model ran")
	}
	if report.ReviewMode != "" {
		kv("mode", report.ReviewMode)
	}
	if report.ReviewBase != "" {
		kv("base", report.ReviewBase)
	}
	if report.ContextSnippets > 0 {
		kv("context", fmt.Sprintf("%d snippet(s)", report.ContextSnippets))
	}
	for _, e := range report.Evidence {
		src := strings.TrimSpace(e.Source)
		if src == "" {
			continue
		}
		state := e.State
		if d := strings.TrimSpace(e.Detail); d != "" {
			state += " — " + d
		}
		paint := termstyle.Muted
		if e.State == EvidenceOK {
			paint = termstyle.Success
		} else if e.Degraded() {
			paint = termstyle.Warning
		}
		kv("evidence", colorize(color, paint, src+": "+state))
	}
	if s := strings.TrimSpace(report.Coverage.Statement()); s != "" {
		kv("coverage", s)
	}
	for _, r := range report.DegradedReasons {
		kv("degraded", colorize(color, termstyle.Warning, r))
	}
	if b.Len() == 0 {
		return reviewRenderIndent + colorize(color, termstyle.Muted, "no run details recorded") + "\n"
	}
	return b.String()
}
