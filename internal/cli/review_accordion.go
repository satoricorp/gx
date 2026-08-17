package cli

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/termstyle"
)

// The accordion: the interactive shape of the review report.
//
// When gx review has a terminal on both ends, the linear render becomes a
// menu. The header and the run ledger print once and stay; the verdict prints
// as its own node so the reader sees NO-SHIP before choosing what to read;
// then a list of sections. The keys are the ones every list UI trains:
//
//	↑ ↓ (k j)   move between sections — and open the one you land on, so
//	            browsing is reading; no second keypress to see it
//	tab ⇧tab    the same, wrapping at either end
//	1‥5         jump straight to that section
//	enter → l   open / close the section under the cursor (enter on Done quits)
//	← h esc     close the open section (esc with nothing open quits)
//	PgDn PgUp   scroll the open section; space / b and ctrl+d / ctrl+u too,
//	            and the mouse wheel; g / G for top / bottom
//	J           toggle the --json object
//	q ctrl+c    quit with the same exit code the linear render produces
//
// The first cut used tab to move and ↑↓ to scroll, pager-style. In practice
// nobody reached for tab, ↑↓ "did nothing" (they scrolled a report that was
// not yet open), tab did not wrap so it read as one-way, and ⇧tab is dropped
// by some terminals — so the menu felt broken. Arrows move; the wheel scrolls.
//
// It is the clack rail-and-diamond grammar — ◆ for an open group, ● for the
// current row, ◉ for a visited one, ○ for an unvisited one, │ down the left —
// drawn with the bubbletea and lipgloss already in the tree; no new dependency.
//
// Screen model. The program takes the alternate screen, like less or a
// pager, so it owns the whole window and leaves the shell scrollback intact
// on exit. The report — header, ledger, verdict, and the open section — lives
// in a viewport anchored at the TOP: opening a section shows its first line,
// and the reader scrolls DOWN through it. The menu is a fixed footer beneath
// the viewport, so it never scrolls out of view. The first cut rendered
// everything inline and let the terminal repaint from the bottom, which put
// the reader at the end of a long section and made them scroll up to find
// its start — the wrong way round.
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
	height   int
	// vp holds everything above the menu and scrolls it; the menu itself is
	// drawn beneath vp with a fixed height so it is always on screen.
	vp    viewport.Model
	ready bool

	// Precomputed once: the lanes, the plan, and the footer, so View() never
	// re-derives — View() runs on every keypress, and the footer reads
	// REVIEW.md from disk. (gx review flagged this itself, on this branch.)
	blocking []codereview.Finding
	advisory []codereview.Finding
	plan     []codereview.FixStep
	footer   string
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
		footer:   accordionFooter(report),
		width:    100,
		height:   40,
		vp:       viewport.New(),
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
		if msg.Height > 0 {
			m.height = msg.Height
		}
		m.resize()
		m.ready = true
		return m, nil
	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			m.vp.ScrollUp(3)
		case tea.MouseWheelDown:
			m.vp.ScrollDown(3)
		}
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		// Move — and show. Landing on a section opens it, so ↓ ↓ ↓ reads the
		// report top to bottom. Wraps, so tab from the last row is not a dead
		// key. Landing on Done leaves whatever was open on screen.
		case "down", "j", "tab":
			m.moveTo(m.nextRow(m.cursor))
		case "up", "k", "shift+tab":
			m.moveTo(m.prevRow(m.cursor))
		case "1", "2", "3", "4", "5":
			if s := accordionSection(msg.String()[0] - '1'); m.selectable(s) {
				m.moveTo(s)
			}
		// Open / close the section under the cursor. Enter on Done quits.
		case "enter", "right", "l":
			if m.cursor == sectionDone {
				m.quit = true
				return m, tea.Quit
			}
			if m.open == m.cursor {
				m.close()
			} else {
				m.show(m.cursor)
			}
		case "left", "h":
			m.close()
		case "esc":
			if !m.showJSON && m.open >= sectionCount {
				m.quit = true
				return m, tea.Quit
			}
			m.close()
		// Scroll the open section: pager keys, and the wheel above.
		case "pgup", "b", "ctrl+u":
			m.vp.HalfPageUp()
		case "pgdown", "f", "ctrl+d", " ":
			m.vp.HalfPageDown()
		case "g", "home":
			m.vp.GotoTop()
		case "G", "end":
			m.vp.GotoBottom()
		case "J":
			m.showJSON = !m.showJSON
			m.open = sectionCount
			m.refresh()
			m.scrollToOpenSection()
		}
	}
	return m, nil
}

// moveTo puts the cursor on s and, when s is a readable section, opens it.
// Done is only highlighted: it has nothing to show, and closing what the
// reader was looking at just because they moved past it would be a surprise.
func (m *accordionModel) moveTo(s accordionSection) {
	m.cursor = s
	if s == sectionDone || !m.selectable(s) {
		return
	}
	m.show(s)
}

// show opens s and puts its header at the top of the viewport.
func (m *accordionModel) show(s accordionSection) {
	m.open = s
	m.visited[s] = true
	m.showJSON = false
	m.refresh()
	m.scrollToOpenSection()
}

// close collapses whatever is open — a section or the JSON — and returns the
// viewport to the top, where the header and verdict are.
func (m *accordionModel) close() {
	m.showJSON = false
	m.open = sectionCount
	m.refresh()
	m.vp.GotoTop()
}

// resize fits the viewport under the fixed footer (the menu) and refreshes.
func (m *accordionModel) resize() {
	footer := len(strings.Split(m.menuView(), "\n"))
	h := m.height - footer
	if h < 5 {
		h = 5
	}
	m.vp.SetWidth(m.width)
	m.vp.SetHeight(h)
	m.refresh()
}

// refresh re-renders the scrollable content into the viewport, preserving
// the reader's scroll position when the content did not shrink past it.
func (m *accordionModel) refresh() {
	off := m.vp.YOffset()
	m.vp.SetContent(m.contentView())
	if off < m.vp.TotalLineCount() {
		m.vp.SetYOffset(off)
	}
}

// scrollToOpenSection puts the open section's header at the top of the
// viewport — the reader sees the start and scrolls down, never up.
func (m *accordionModel) scrollToOpenSection() {
	if m.open >= sectionCount && !m.showJSON {
		m.vp.GotoTop()
		return
	}
	line := m.openSectionLine()
	if line < 0 {
		m.vp.GotoTop()
		return
	}
	m.vp.SetYOffset(line)
}

// nextRow / prevRow move the cursor, skipping sections that would render
// nothing — an empty ADVISORY row is still listed (as "0 findings") but the
// cursor does not stop on it, so enter never opens an empty section. Both
// wrap: past Done comes the first section, before the first comes Done. A
// menu whose last row eats keypresses reads as broken, not as finished.
func (m accordionModel) nextRow(from accordionSection) accordionSection {
	for i := accordionSection(1); i <= sectionCount; i++ {
		s := (from + i) % sectionCount
		if m.selectable(s) {
			return s
		}
	}
	return from
}

func (m accordionModel) prevRow(from accordionSection) accordionSection {
	for i := accordionSection(1); i <= sectionCount; i++ {
		s := (from + sectionCount - i) % sectionCount
		if m.selectable(s) {
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

// paint helpers shared by the content and menu views.
func (m accordionModel) painters() (c func(func(string) string, string) string, mint, mute, rail func(string) string) {
	c = func(paint func(string) string, s string) string {
		if !m.color || s == "" {
			return s
		}
		return paint(s)
	}
	mint = func(s string) string { return c(codereview.ReviewMint, s) }
	mute = func(s string) string { return c(termstyle.Muted, s) }
	rail = func(s string) string {
		if s == "" {
			return mute("│")
		}
		return mute("│") + "  " + s
	}
	return c, mint, mute, rail
}

// contentView is everything that scrolls: header, ledger, verdict, and the
// open section. It is what the viewport holds.
func (m accordionModel) contentView() string {
	var b strings.Builder
	c, mint, mute, rail := m.painters()

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
		b.WriteString(rail("") + "\n")
	}
	return b.String()
}

// openSectionLine is the 0-based line in contentView where the open section's
// ◆ header sits, or -1 when nothing is open. The viewport scrolls to it.
func (m accordionModel) openSectionLine() int {
	if m.open >= sectionCount && !m.showJSON {
		return -1
	}
	lines := strings.Split(m.contentView(), "\n")
	// The section header is the second ◆ line (the first is the gx review
	// header). Find it by scanning for a line that starts with the ◆ glyph
	// after the first.
	seen := 0
	for i, l := range lines {
		plain := stripANSI(l)
		if strings.HasPrefix(plain, "◆") {
			seen++
			if seen == 2 {
				return i
			}
		}
	}
	return -1
}

// menuView is the fixed footer: the section list and the pack line.
func (m accordionModel) menuView() string {
	var b strings.Builder
	c, mint, mute, rail := m.painters()
	b.WriteString(mint("◆") + "  " + c(termstyle.Section, "Sections") + mute("   ↑↓ move · enter open/close · PgUp PgDn scroll · 1-5 jump · J json · q quit") + "\n")
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
		// The digit is the shortcut: 1‥5 jumps to that section. Done has none.
		digit := " "
		if s < sectionDone {
			digit = mute(fmt.Sprintf("%d", int(s)+1))
		}
		row := digit + " " + mark + "  " + label + pad + fmt.Sprintf("%-11s", m.sectionCount(s)) + mute(m.sectionHint(s))
		if s == m.cursor {
			row = mint("❯") + " " + row
		} else {
			row = "  " + row
		}
		b.WriteString(rail(row) + "\n")
	}
	b.WriteString(rail("") + "\n")
	b.WriteString(mint("└") + "  " + mute(m.footer))
	// A scroll indicator when the content is taller than the viewport, so the
	// reader knows there is more above or below.
	if m.ready && m.vp.TotalLineCount() > m.vp.VisibleLineCount() {
		pct := 100
		if m.vp.TotalLineCount() > 0 {
			pct = int(float64(m.vp.YOffset()+m.vp.VisibleLineCount()) / float64(m.vp.TotalLineCount()) * 100)
			if pct > 100 {
				pct = 100
			}
		}
		b.WriteString(mute(fmt.Sprintf("  · %d%%", pct)))
	}
	b.WriteString("\n")
	return b.String()
}

// View composes the scrolling content over the fixed menu, on the alternate
// screen so the shell's scrollback is untouched.
func (m accordionModel) View() tea.View {
	var content string
	if m.ready {
		content = m.vp.View()
	} else {
		// Before the first WindowSizeMsg the viewport has no size; show the
		// content unclipped so tests and a headless run still see it.
		content = strings.TrimRight(m.contentView(), "\n")
	}
	v := tea.NewView(content + "\n" + m.menuView())
	v.AltScreen = true
	// Cell-motion mouse reporting is what makes the wheel scroll the report;
	// it does not capture drag-select the way all-motion would.
	v.MouseMode = tea.MouseModeCellMotion
	return v
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
		if len(m.report.DegradedReasons) > 0 {
			return "RUN DETAILS", termstyle.Warning
		}
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
		if n := len(m.report.DegradedReasons); n > 0 {
			return fmt.Sprintf("⚠ degraded run — %s", codereview.SummarizeDegradedReasons(m.report.DegradedReasons))
		}
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
