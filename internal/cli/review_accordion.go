package cli

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/termstyle"
)

// The accordion: the interactive shape of the review report.
//
// When gx review has a terminal on both ends, the linear render becomes a
// menu. The header and the run ledger print once and stay; the verdict prints
// as its own node so the reader sees NO-SHIP before choosing what to read;
// then a list of sections. Enter opens a section in place under its own
// marker, and the menu re-renders beneath with the cursor advanced to the next
// unvisited row, so the reader walks the report top to bottom without ever
// losing the list. Esc collapses. j prints the JSON. q (or Done) quits with
// the same exit code the linear render would have produced.
//
// It is the clack rail-and-diamond grammar — ◆ for an open group, ● for the
// current row, ◉ for a visited one, ○ for an unvisited one, │ down the left —
// drawn with the bubbletea and lipgloss already in the tree; no new dependency.
//
// Two rules keep it honest. Agents and CI never see it: the model is only
// entered when useInteractiveTerminal says both fds are ttys, and every other
// path prints RenderReviewText, so the exit code and the Verdict/Next seam are
// identical in both modes. And it renders from the same Report the linear
// render reads — the accordion is presentation over a stable struct, never a
// second source of truth.

type accordionSection int

const (
	sectionBlocking accordionSection = iota
	sectionAdvisory
	sectionFixPlan
	sectionStory
	sectionRunDetails
	sectionDone
	sectionCount
)

type accordionModel struct {
	report   codereview.Report
	color    bool
	cursor   accordionSection
	open     accordionSection // sectionCount when nothing is open
	visited  [sectionCount]bool
	showJSON bool
	quit     bool
	width    int

	// Precomputed once: the lanes and the plan, so View() never re-derives.
	blocking []codereview.Finding
	advisory []codereview.Finding
	plan     []codereview.FixStep
}

func newAccordionModel(report codereview.Report, color bool) accordionModel {
	blocking, advisory := codereview.SplitFindingsByLane(report.Findings)
	m := accordionModel{
		report:   report,
		color:    color,
		open:     sectionCount,
		blocking: blocking,
		advisory: advisory,
		plan:     codereview.BuildFixPlan(report.Findings),
		width:    100,
	}
	// Start on the first section that has content, so enter does something.
	m.cursor = m.firstNonEmpty()
	return m
}

func (m accordionModel) Init() tea.Cmd { return nil }

func (m accordionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "up", "k":
			m.cursor = m.prevRow(m.cursor)
		case "down", "j":
			m.cursor = m.nextRow(m.cursor)
		case "enter", " ":
			if m.cursor == sectionDone {
				m.quit = true
				return m, tea.Quit
			}
			if m.open == m.cursor {
				m.open = sectionCount
			} else {
				m.open = m.cursor
				m.visited[m.cursor] = true
				m.showJSON = false
				// Advance to the next unvisited row so the reader walks down.
				if next := m.nextUnvisited(m.cursor); next != m.cursor {
					m.cursor = next
				}
			}
		case "esc":
			if m.showJSON {
				m.showJSON = false
			} else {
				m.open = sectionCount
			}
		case "J":
			m.showJSON = !m.showJSON
			m.open = sectionCount
		}
	}
	return m, nil
}

// nextRow / prevRow move the cursor, skipping sections that would render
// nothing — an empty ADVISORY row is still listed (as "0 findings") but the
// cursor does not stop on it, so enter never opens an empty section.
func (m accordionModel) nextRow(from accordionSection) accordionSection {
	for s := from + 1; s < sectionCount; s++ {
		if m.selectable(s) {
			return s
		}
	}
	return from
}

func (m accordionModel) prevRow(from accordionSection) accordionSection {
	for s := from - 1; s >= 0; s-- {
		if m.selectable(s) {
			return s
		}
	}
	return from
}

func (m accordionModel) nextUnvisited(from accordionSection) accordionSection {
	for s := from + 1; s < sectionCount; s++ {
		if m.selectable(s) && !m.visited[s] {
			return s
		}
	}
	return from
}

func (m accordionModel) firstNonEmpty() accordionSection {
	for s := accordionSection(0); s < sectionCount; s++ {
		if m.selectable(s) {
			return s
		}
	}
	return sectionDone
}

func (m accordionModel) selectable(s accordionSection) bool {
	switch s {
	case sectionBlocking:
		return len(m.blocking) > 0
	case sectionAdvisory:
		return len(m.advisory) > 0
	case sectionFixPlan:
		return len(m.plan) > 0
	case sectionStory:
		return len(m.report.Story) > 0
	case sectionRunDetails, sectionDone:
		return true
	}
	return false
}

// ---- view -----------------------------------------------------------------

func (m accordionModel) View() tea.View {
	var b strings.Builder
	c := func(paint func(string) string, s string) string {
		if !m.color || s == "" {
			return s
		}
		return paint(s)
	}
	mint := func(s string) string { return c(codereview.ReviewMint, s) }
	mute := func(s string) string { return c(termstyle.Muted, s) }
	rail := func(s string) string {
		if s == "" {
			return mute("│")
		}
		return mute("│") + "  " + s
	}

	// Header node.
	target := strings.TrimSpace(m.report.ReviewRange)
	if target == "" {
		target = strings.TrimSpace(m.report.ReviewTarget)
	}
	b.WriteString(mint("◆") + "  " + mint("gx review") + mute(fmt.Sprintf(" — %s · %d file(s)", target, len(m.report.ChangedFiles))) + "\n")
	if p := strings.TrimSpace(m.report.Prompt); p != "" {
		b.WriteString(rail(mute(fmt.Sprintf("intent: %q", p))) + "\n")
	}
	b.WriteString(rail("") + "\n")
	for _, line := range accordionLedgerRows(m.report, m.color, m.blocking, m.advisory) {
		b.WriteString(rail(line) + "\n")
	}
	b.WriteString(rail("") + "\n")

	// Verdict node — before the menu, always.
	verdict, _ := codereview.ReviewVerdictLines(m.report)
	vp := termstyle.Success
	if len(m.blocking) > 0 || !m.report.Reviewed {
		vp = termstyle.Danger
	} else if len(m.report.DegradedReasons) > 0 {
		vp = termstyle.Warning
	}
	head, _, _ := strings.Cut(verdict, " — ")
	tail := strings.TrimPrefix(verdict, head)
	b.WriteString(c(vp, "◇") + "  " + c(vp, head) + mute(tail) + "\n")
	b.WriteString(rail("") + "\n")

	// The open section, if any, under its own node.
	if m.showJSON {
		b.WriteString(mint("◆") + "  " + c(termstyle.Section, "JSON") + mute("  gx.review/2 · same object --json prints") + "\n")
		b.WriteString(rail("") + "\n")
		for _, line := range strings.Split(strings.TrimRight(accordionJSON(m.report), "\n"), "\n") {
			b.WriteString(rail(mute(line)) + "\n")
		}
		b.WriteString(rail("") + "\n")
	} else if m.open < sectionCount {
		title, paint := m.sectionLabel(m.open)
		b.WriteString(mint("◆") + "  " + c(paint, title) + mute("  "+m.sectionCount(m.open)) + "\n")
		b.WriteString(rail("") + "\n")
		for _, line := range strings.Split(strings.TrimRight(m.sectionBody(m.open), "\n"), "\n") {
			b.WriteString(rail(line) + "\n")
		}
		b.WriteString(rail(mute("  [esc] collapse   [↑↓] move   [J] json   [q] quit")) + "\n")
		b.WriteString(rail("") + "\n")
	}

	// The menu.
	b.WriteString(mint("◆") + "  " + c(termstyle.Section, "Open a section") + mute("   ↑↓ move · enter open · esc collapse · J json · q quit") + "\n")
	b.WriteString(rail("") + "\n")
	for s := accordionSection(0); s < sectionCount; s++ {
		title, paint := m.sectionLabel(s)
		var mark string
		switch {
		case s == m.open || s == m.cursor:
			mark = mint("●")
		case m.visited[s]:
			mark = mute("◉")
		default:
			mark = mute("○")
		}
		label := c(paint, title)
		if !m.selectable(s) {
			label = mute(title)
		}
		pad := strings.Repeat(" ", max(1, 16-len(title)))
		row := mark + "  " + label + pad + fmt.Sprintf("%-11s", m.sectionCount(s)) + mute(m.sectionHint(s))
		if s == m.cursor {
			row = mint("❯") + " " + row
		} else {
			row = "  " + row
		}
		b.WriteString(rail(row) + "\n")
	}
	b.WriteString(rail("") + "\n")
	b.WriteString(mint("└") + "  " + mute(accordionFooter(m.report)) + "\n")
	return tea.NewView(b.String())
}

func (m accordionModel) sectionLabel(s accordionSection) (string, func(string) string) {
	switch s {
	case sectionBlocking:
		return "BLOCKING", termstyle.Danger
	case sectionAdvisory:
		return "ADVISORY", termstyle.Warning
	case sectionFixPlan:
		return "FIX PLAN", codereview.ReviewMint
	case sectionStory:
		return "WORTH KNOWING", termstyle.Command
	case sectionRunDetails:
		return "RUN DETAILS", termstyle.Muted
	default:
		return "Done", termstyle.Value
	}
}

func (m accordionModel) sectionCount(s accordionSection) string {
	switch s {
	case sectionBlocking:
		return plural(len(m.blocking), "finding")
	case sectionAdvisory:
		return plural(len(m.advisory), "finding")
	case sectionFixPlan:
		return plural(len(m.plan), "step")
	case sectionStory:
		return plural(len(m.report.Story), "change")
	}
	return ""
}

func (m accordionModel) sectionHint(s accordionSection) string {
	names := func(fs []codereview.Finding) string {
		var out []string
		seen := map[string]bool{}
		for _, f := range fs {
			n := codereview.RuleShortName(f)
			if n == "" || seen[n] {
				continue
			}
			seen[n] = true
			out = append(out, n)
			if len(out) == 3 {
				break
			}
		}
		return strings.Join(out, " · ")
	}
	switch s {
	case sectionBlocking:
		return names(m.blocking)
	case sectionAdvisory:
		return names(m.advisory)
	case sectionFixPlan:
		return "the agent's list — fix, then rerun"
	case sectionStory:
		high := 0
		for _, it := range m.report.Story {
			if it.Materiality == "high" {
				high++
			}
		}
		if high > 0 {
			return fmt.Sprintf("passed, and yours now · %d high materiality", high)
		}
		return "passed, and yours now"
	case sectionRunDetails:
		return "reviewers · evidence · coverage"
	default:
		verdict, _ := codereview.ReviewVerdictLines(m.report)
		code := 0
		if len(m.blocking) > 0 {
			code = reviewFindingsExitCode
		}
		head, _, _ := strings.Cut(strings.TrimPrefix(verdict, "Verdict: "), " — ")
		return fmt.Sprintf("exit %d (%s)", code, strings.ToLower(head))
	}
}

// sectionBody renders one section by delegating to the linear renderer's
// section writers, so the accordion and RenderReviewText cannot drift.
func (m accordionModel) sectionBody(s accordionSection) string {
	switch s {
	case sectionBlocking:
		return codereview.RenderReviewLane(m.report, m.color, "BLOCKING", m.blocking)
	case sectionAdvisory:
		return codereview.RenderReviewLane(m.report, m.color, "ADVISORY", m.advisory)
	case sectionFixPlan:
		return codereview.RenderReviewFixPlan(m.report, m.color)
	case sectionStory:
		return codereview.RenderReviewStory(m.report, m.color)
	case sectionRunDetails:
		return codereview.RenderReviewDetails(m.report, m.color)
	}
	return ""
}

func accordionLedgerRows(report codereview.Report, color bool, blocking, advisory []codereview.Finding) []string {
	return strings.Split(strings.TrimRight(codereview.RenderReviewLedgerRows(report, color), "\n"), "\n")
}

func accordionFooter(report codereview.Report) string {
	footer := "gx:recommended v" + codereview.RecommendedPackVersion
	if report.RepoRoot != "" {
		if policy := codereview.LoadReviewPolicy(report.RepoRoot); len(policy.Rules) > 0 {
			footer += fmt.Sprintf(" + %d rule(s) from REVIEW.md", len(policy.Rules))
		}
	}
	return footer
}

func accordionJSON(report codereview.Report) string {
	var b strings.Builder
	if err := writeReviewJSON(&b, report); err != nil {
		return "(could not encode report: " + err.Error() + ")"
	}
	return b.String()
}

func plural(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}
