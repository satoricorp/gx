package cli

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/satoricorp/gx/internal/codereview"
)

// accordionFixture is the checkout-retry shape: one blocking finding both
// graders agreed on, one demoted advisory, a story item, and a fix plan.
func accordionFixture() codereview.Report {
	twoLegs := []string{"Bedrock A", "Bedrock B"}
	findings := codereview.AssignLanes([]codereview.Finding{
		{
			ID: "bedrock-a.ai.review.1", RuleID: "gx:recommended/no-secrets-in-logs",
			Title: "Secrets must not reach logs", Summary: "req embeds PaymentToken.",
			Recommendation: "log req.ID instead of req.", Strength: "Strong",
			Corroboration: twoLegs, JudgeVerdict: "confirmed", JudgeConfidence: 0.91,
			File: "internal/checkout/session.go", Line: 88,
		},
		{
			ID: "bedrock-b.ai.review.2", RuleID: "gx:recommended/comments-match-code",
			Title: "Comments must describe what the code now does", Summary: "stale comment.",
			Recommendation: "update the comment.", Strength: "Strong",
			Corroboration: []string{"Bedrock B"}, JudgeVerdict: "confirmed", JudgeConfidence: 0.77,
			File: "internal/checkout/session.go", Line: 71,
		},
	})
	return codereview.Report{
		Reviewed: true, ReviewMode: codereview.ReviewModeRange,
		ReviewRange: "origin/main...feature/checkout-retry", ReviewBase: "origin/main",
		ChangedFiles: []string{"a.go", "b.go"}, Reviewer: "heuristic+ai",
		ReviewModels: []string{"claude-a", "claude-b"}, ReviewTransport: "gx Cloud",
		Findings: findings,
		Story: []codereview.StoryItem{{
			Headline: "Checkout retries on 5xx where it used to fail fast", Category: codereview.StoryBehaviorDelta,
			Materiality: "high", File: "internal/checkout/session.go", Consequence: "up to 4 attempts over ~15s.",
		}},
	}
}

func key(s string) tea.KeyPressMsg {
	switch s {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "tab":
		return tea.KeyPressMsg{Code: tea.KeyTab}
	case "shift+tab":
		return tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}
	}
	r := []rune(s)[0]
	return tea.KeyPressMsg{Code: r, Text: s}
}

func press(t *testing.T, m accordionModel, keys ...string) accordionModel {
	t.Helper()
	for _, k := range keys {
		next, _ := m.Update(key(k))
		var ok bool
		m, ok = next.(accordionModel)
		if !ok {
			t.Fatalf("Update returned %T, want accordionModel", next)
		}
	}
	return m
}

func view(m accordionModel) string { return m.View().Content }

func TestAccordionOpensWithHeaderVerdictAndMenu(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	out := view(m)
	for _, want := range []string{
		"◆  gx review — origin/main...feature/checkout-retry",
		"◇  Verdict: NO-SHIP", // verdict node before the menu, always
		"◆  Sections",
		"❯ 1 ●  BLOCKING", // cursor starts on the first non-empty section; 1 is its shortcut
		"1 finding",
		"no-secrets-in-logs", // rule name in the hint column
		"2 ○  ADVISORY",
		"3 ○  FIX PLAN",
		"4 ○  WORTH KNOWING",
		"5 ○  RUN DETAILS",
		"○  Done",
		"exit 3 (no-ship)",
		"└  gx:recommended v",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("initial view missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("color=false view must carry no ANSI")
	}
	// Nothing is open yet: no finding body on screen.
	if strings.Contains(out, "Why  req embeds") {
		t.Fatalf("no section should be open before enter:\n%s", out)
	}
}

func TestAccordionEnterOpensSectionAndCursorStays(t *testing.T) {
	m := press(t, newAccordionModel(accordionFixture(), false), "enter")
	out := view(m)
	// BLOCKING opened under its own node with the finding body, rendered by
	// the same writer the linear render uses.
	for _, want := range []string{
		"◆  BLOCKING  1 finding",
		"Secrets must not reach logs",
		"no-secrets-in-logs  gx:recommended  ·  internal/checkout/session.go:88",
		"Why  req embeds PaymentToken.",
		"Fix  log req.ID instead of req.",
		"both graders agreed · judge confirmed",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("after enter, missing %q:\n%s", want, out)
		}
	}
	// The menu re-renders beneath with BLOCKING open (●) and the cursor still
	// on it: enter is "show me this", not "and move on". ↓ moves.
	if !strings.Contains(out, "❯ 1 ●  BLOCKING") {
		t.Fatalf("menu should show the cursor on the open BLOCKING row:\n%s", out)
	}
	if m.open != sectionBlocking || m.cursor != sectionBlocking {
		t.Fatalf("open=%d cursor=%d, want both blocking", m.open, m.cursor)
	}
	// Enter again closes it.
	m = press(t, m, "enter")
	if m.open != sectionCount {
		t.Fatalf("second enter should close the section, open=%d", m.open)
	}
}

// The core of the new contract: arrows move between sections AND show the one
// you land on, so ↓ ↓ ↓ reads the report. Tab does the same. Both wrap.
func TestAccordionArrowsMoveAndOpen(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	m = press(t, m, "down")
	if m.cursor != sectionAdvisory || m.open != sectionAdvisory {
		t.Fatalf("↓ should move to ADVISORY and open it: cursor=%d open=%d", m.cursor, m.open)
	}
	if !strings.Contains(view(m), "◆  ADVISORY") {
		t.Fatalf("ADVISORY body should be on screen after ↓:\n%s", view(m))
	}
	m = press(t, m, "up")
	if m.cursor != sectionBlocking || m.open != sectionBlocking {
		t.Fatalf("↑ should move back to BLOCKING and open it: cursor=%d open=%d", m.cursor, m.open)
	}
	// Wrap: ↑ from the first row lands on Done, which is highlighted but not
	// "opened" — the section that was showing stays showing.
	m = press(t, m, "up")
	if m.cursor != sectionDone {
		t.Fatalf("↑ from the first row should wrap to Done, got %d", m.cursor)
	}
	if m.open != sectionBlocking {
		t.Fatalf("landing on Done must not close what was open, open=%d", m.open)
	}
	// And ↓ from Done wraps to the first section.
	m = press(t, m, "down")
	if m.cursor != sectionBlocking {
		t.Fatalf("↓ from Done should wrap to BLOCKING, got %d", m.cursor)
	}
	// tab / shift+tab are the same moves.
	m = press(t, m, "tab")
	if m.cursor != sectionAdvisory || m.open != sectionAdvisory {
		t.Fatalf("tab should behave like ↓: cursor=%d open=%d", m.cursor, m.open)
	}
	m = press(t, m, "shift+tab")
	if m.cursor != sectionBlocking {
		t.Fatalf("shift+tab should behave like ↑, got %d", m.cursor)
	}
}

func TestAccordionDigitsJump(t *testing.T) {
	m := press(t, newAccordionModel(accordionFixture(), false), "4")
	if m.cursor != sectionStory || m.open != sectionStory {
		t.Fatalf("4 should jump to WORTH KNOWING and open it: cursor=%d open=%d", m.cursor, m.open)
	}
	if !strings.Contains(view(m), "◆  WORTH KNOWING") {
		t.Fatalf("story body should be on screen after 4")
	}
	// A digit for a section with nothing in it is ignored.
	r := accordionFixture()
	r.Story = nil
	m = press(t, newAccordionModel(r, false), "4")
	if m.cursor != sectionBlocking {
		t.Fatalf("4 on an empty story section should do nothing, cursor=%d", m.cursor)
	}
}

func TestAccordionEscCollapsesAndVisitedMarkPersists(t *testing.T) {
	// Open BLOCKING, move the cursor off it (↓ opens ADVISORY), esc closes.
	m := press(t, newAccordionModel(accordionFixture(), false), "enter", "down", "esc")
	out := view(m)
	if strings.Contains(out, "◆  ADVISORY") {
		t.Fatalf("esc should collapse the open section:\n%s", out)
	}
	if !strings.Contains(out, "1 ◉  BLOCKING") {
		t.Fatalf("a visited section should keep its ◉ mark after the cursor leaves it:\n%s", out)
	}
	if m.open != sectionCount {
		t.Fatalf("open should be reset after esc, got %d", m.open)
	}
	// ← closes too; esc with nothing open quits.
	m = press(t, m, "enter")
	if m.open != sectionAdvisory {
		t.Fatalf("enter should reopen the cursor's section")
	}
	m = press(t, m, "left")
	if m.open != sectionCount {
		t.Fatalf("← should close the open section")
	}
	if next, cmd := m.Update(key("esc")); cmd == nil || !next.(accordionModel).quit {
		t.Fatalf("esc with nothing open should quit")
	}
}

func TestAccordionWalksSectionsAndDelegatesToLinearRenderers(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	// enter (blocking) → ↓ (advisory) → ↓ (fix plan) → ↓ (story)
	m = press(t, m, "enter", "down", "down", "down")
	out := view(m)
	if !strings.Contains(out, "◆  WORTH KNOWING") || !strings.Contains(out, "Checkout retries on 5xx") {
		t.Fatalf("third ↓ should open the story:\n%s", out)
	}
	if !strings.Contains(out, "What changes for you: up to 4 attempts") {
		t.Fatalf("story body should come from the linear renderer:\n%s", out)
	}
	for _, s := range []accordionSection{sectionBlocking, sectionAdvisory, sectionFixPlan, sectionStory} {
		if !m.visited[s] {
			t.Fatalf("section %d should be visited", s)
		}
	}
	if m.cursor != sectionStory {
		t.Fatalf("cursor should be on the story it opened, got %d", m.cursor)
	}
}

func TestAccordionCursorSkipsEmptySections(t *testing.T) {
	r := accordionFixture()
	r.Story = nil // no story → WORTH KNOWING not selectable
	m := newAccordionModel(r, false)
	m = press(t, m, "tab", "tab", "tab") // blocking → advisory → fix plan → (skip story) run details
	if m.cursor != sectionRunDetails {
		t.Fatalf("cursor should skip the empty story section, got %d", m.cursor)
	}
	// And the menu still lists it, muted, with its zero count.
	if !strings.Contains(view(m), "WORTH KNOWING") {
		t.Fatalf("empty section should still be listed")
	}
}

func TestAccordionJSONToggle(t *testing.T) {
	m := press(t, newAccordionModel(accordionFixture(), false), "J")
	out := view(m)
	if !strings.Contains(out, "◆  JSON") || !strings.Contains(out, `"schema": "gx.review/2"`) {
		t.Fatalf("J should show the JSON envelope:\n%s", out)
	}
	m = press(t, m, "esc")
	if strings.Contains(view(m), "◆  JSON") {
		t.Fatalf("esc should hide the JSON")
	}
}

func TestAccordionQuitPaths(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	if next, cmd := m.Update(key("q")); cmd == nil || !next.(accordionModel).quit {
		t.Fatalf("q should quit")
	}
	// Enter on Done quits too.
	m = press(t, m, "tab", "tab", "tab", "tab", "tab") // to Done
	if m.cursor != sectionDone {
		t.Fatalf("expected cursor on Done, got %d", m.cursor)
	}
	if next, cmd := m.Update(key("enter")); cmd == nil || !next.(accordionModel).quit {
		t.Fatalf("enter on Done should quit")
	}
}

func TestAccordionNothingToReviewShowsVerdictAndOnlyDetailsAndDone(t *testing.T) {
	r := codereview.Report{Reviewed: false, ReviewTarget: "the working tree"}
	m := newAccordionModel(r, false)
	out := view(m)
	if !strings.Contains(out, "◇  Verdict: NOTHING-TO-REVIEW") {
		t.Fatalf("verdict node missing:\n%s", out)
	}
	if m.cursor != sectionRunDetails {
		t.Fatalf("with no findings the first selectable row is RUN DETAILS, got %d", m.cursor)
	}
}

func TestBrowseReviewInteractivelyFallsBackWhenNotATTY(t *testing.T) {
	// A bytes.Buffer is not an *os.File, so useInteractiveTerminal says no and
	// the caller must print the linear render. This is the agent/CI path.
	var in, out strings.Builder
	if browseReviewInteractively(strings.NewReader(in.String()), &out, accordionFixture()) {
		t.Fatalf("non-tty must not open the accordion")
	}
	if out.Len() != 0 {
		t.Fatalf("non-tty path must write nothing itself, got %q", out.String())
	}
}

// The bug this guards: opening a long section used to leave the reader at
// its END (inline repaint from the bottom), scrolling up to find the start.
// The viewport must put the section header at the top and scroll down.
func TestAccordionOpenSectionStartsAtItsHeader(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	// Give it a small window so the content overflows and scrolling matters.
	next, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 24})
	m = next.(accordionModel)
	m = press(t, m, "enter") // open BLOCKING
	if !m.ready {
		t.Fatalf("model should be sized after WindowSizeMsg")
	}
	// The section's ◆ header must be at the top of the viewport — or, when
	// the content is too short for the header to reach row 0 (the viewport
	// clamps at the bottom), the header must at least be visible and nothing
	// of the section may sit above the fold. Either way: the reader sees the
	// START, never the end.
	top := m.vp.YOffset()
	header := m.openSectionLine()
	maxOff := m.vp.TotalLineCount() - m.vp.VisibleLineCount()
	if header < 0 {
		t.Fatalf("no open section header found")
	}
	if top != min(header, maxOff) {
		t.Fatalf("viewport top = %d, want the section header line %d (or the clamp %d)", top, header, maxOff)
	}
	if header < top {
		t.Fatalf("section header (%d) scrolled above the fold (%d)", header, top)
	}
	// And the view is on the alt screen with the menu beneath the content.
	v := m.View()
	if !v.AltScreen {
		t.Fatalf("accordion must take the alt screen so scrollback stays intact")
	}
	if !strings.Contains(v.Content, "◆  Sections") {
		t.Fatalf("menu must be present beneath the viewport:\n%s", v.Content)
	}
	// PgDn scrolls the open section; the cursor stays put.
	before := m.cursor
	m = press(t, m, "pgdown")
	if m.cursor != before {
		t.Fatalf("PgDn must scroll, not move the menu cursor")
	}
	// Scrolled DOWN from the header (clamped at the bottom), never up above it.
	if m.vp.YOffset() < top {
		t.Fatalf("scrolling after open must not move above the section header")
	}
	if maxOff > top && m.vp.YOffset() == top {
		t.Fatalf("PgDn should have moved the viewport down from %d (max %d)", top, maxOff)
	}
	// The mouse wheel scrolls too.
	off := m.vp.YOffset()
	next, _ = m.Update(tea.MouseWheelMsg{Button: tea.MouseWheelUp})
	m = next.(accordionModel)
	if off > 0 && m.vp.YOffset() >= off {
		t.Fatalf("wheel up should scroll up: %d → %d", off, m.vp.YOffset())
	}
}
