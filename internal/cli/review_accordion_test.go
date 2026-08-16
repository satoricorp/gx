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
		"◆  Open a section",
		"❯ ●  BLOCKING", // cursor starts on the first non-empty section
		"1 finding",
		"no-secrets-in-logs", // rule name in the hint column
		"○  ADVISORY",
		"○  FIX PLAN",
		"○  WORTH KNOWING",
		"○  RUN DETAILS",
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

func TestAccordionEnterOpensSectionAndAdvancesCursor(t *testing.T) {
	m := press(t, newAccordionModel(accordionFixture(), false), "enter")
	out := view(m)
	// BLOCKING opened under its own node with the finding body, rendered by
	// the same writer the linear render uses.
	for _, want := range []string{
		"◆  BLOCKING  1 finding",
		"Secrets must not reach logs",
		"no-secrets-in-logs  gx:recommended  ·  internal/checkout/session.go:88  [graded]",
		"Why  req embeds PaymentToken.",
		"Fix  log req.ID instead of req.",
		"both graders agreed · judge confirmed 0.91",
		"[esc] collapse",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("after enter, missing %q:\n%s", want, out)
		}
	}
	// The menu re-renders beneath, BLOCKING marked open (●), and the cursor
	// advanced to the next unvisited row.
	if !strings.Contains(out, "●  BLOCKING") || !strings.Contains(out, "❯ ●  ADVISORY") {
		t.Fatalf("menu should show BLOCKING open and cursor on ADVISORY:\n%s", out)
	}
	if m.open != sectionBlocking || m.cursor != sectionAdvisory {
		t.Fatalf("open=%d cursor=%d, want open=blocking cursor=advisory", m.open, m.cursor)
	}
}

func TestAccordionEscCollapsesAndVisitedMarkPersists(t *testing.T) {
	m := press(t, newAccordionModel(accordionFixture(), false), "enter", "esc")
	out := view(m)
	if strings.Contains(out, "◆  BLOCKING") {
		t.Fatalf("esc should collapse the open section:\n%s", out)
	}
	if !strings.Contains(out, "◉  BLOCKING") {
		t.Fatalf("a visited section should keep its ◉ mark after collapse:\n%s", out)
	}
	if m.open != sectionCount {
		t.Fatalf("open should be reset after esc, got %d", m.open)
	}
}

func TestAccordionWalksSectionsAndDelegatesToLinearRenderers(t *testing.T) {
	m := newAccordionModel(accordionFixture(), false)
	// enter (blocking) → enter (advisory) → enter (fix plan) → enter (story)
	m = press(t, m, "enter", "enter", "enter", "enter")
	out := view(m)
	if !strings.Contains(out, "◆  WORTH KNOWING") || !strings.Contains(out, "Checkout retries on 5xx") {
		t.Fatalf("fourth enter should open the story:\n%s", out)
	}
	if !strings.Contains(out, "What changes for you: up to 4 attempts") {
		t.Fatalf("story body should come from the linear renderer:\n%s", out)
	}
	for _, s := range []accordionSection{sectionBlocking, sectionAdvisory, sectionFixPlan, sectionStory} {
		if !m.visited[s] {
			t.Fatalf("section %d should be visited", s)
		}
	}
	if m.cursor != sectionRunDetails {
		t.Fatalf("cursor should have advanced to RUN DETAILS, got %d", m.cursor)
	}
}

func TestAccordionCursorSkipsEmptySections(t *testing.T) {
	r := accordionFixture()
	r.Story = nil // no story → WORTH KNOWING not selectable
	m := newAccordionModel(r, false)
	m = press(t, m, "down", "down", "down") // blocking → advisory → fix plan → (skip story) run details
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
	m = press(t, m, "down", "down", "down", "down", "down") // to Done
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
