package clitui

import (
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/vcs"
	"github.com/mattn/go-runewidth"
)

func TestRenderStatusSummaryMatchesDesignShape(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	snapshot := vcs.StatusSnapshot{
		RepoLabel: "acme/console",
		BaseRef:   "main",
		Bookmarks: []vcs.BookmarkSnapshot{
			{
				Stack:         vcs.StackInfo{Name: "waitlist + gx-pr", Alias: "s1", BaseRef: "main"},
				Current:       true,
				ChangeCount:   5,
				ApprovedCount: 2,
				FileCount:     18,
			},
			{
				Stack:         vcs.StackInfo{Name: "onboarding repo picker", Alias: "s2", BaseRef: "main"},
				ChangeCount:   3,
				ApprovedCount: 3,
			},
		},
	}

	got := RenderStatusSummary(snapshot)
	for _, want := range []string{
		"$ gx status",
		"2 bookmarks · acme/console",
		"waitlist + gx-pr",
		"← current",
		"main",
		"2/5 approved",
		"onboarding repo picker",
		"base · ⌂ local · ⇡ on origin",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RenderStatusSummary() missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderModifyResultIncludesNextStep(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	snapshot := vcs.StatusSnapshot{
		BaseRef: "main",
		Bookmarks: []vcs.BookmarkSnapshot{{
			Stack:   vcs.StackInfo{Name: "waitlist + gx-pr", Alias: "s1", BaseRef: "main"},
			Current: true,
		}},
	}
	result := vcs.ModifyResult{
		CurrentChange: vcs.ChangeInfo{
			ChangeID:    "qpvkpnrmtxypchange",
			Description: "Validate demux hunk coverage",
		},
	}
	got := RenderModifyResult(snapshot, result, "r2")
	for _, want := range []string{
		"gx edit r2",
		"revision",
		"r2",
		"Validate demux hunk coverage",
		"bookmark",
		"waitlist + gx-pr",
		"next",
		`gx commit -m "Validate demux hunk coverage"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RenderModifyResult() missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderModifyInteractiveUsesCompactEditPicker(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	snapshot := vcs.StatusSnapshot{
		Bookmarks: []vcs.BookmarkSnapshot{{
			Stack:   vcs.StackInfo{Name: "update terminal theme", Alias: "s1", BaseRef: "main"},
			Current: true,
		}},
	}
	revisions := []vcs.RevisionSnapshot{
		{Index: 1, ChangeID: "pqkkqwzustsx", CommitID: "31316432abcdef", Description: "(no message)", ShortID: "31316432", Published: true, SyncNote: "local,origin"},
		{Index: 2, ChangeID: "tqxssqomnyzn", CommitID: "5cb53d15abcdef", Description: "update terminal theme", ShortID: "5cb53d15", Working: true},
	}

	got := RenderModifyInteractive(snapshot, revisions, 1)
	for _, want := range []string{
		"$ gx edit",
		"update terminal theme · 2 revs",
		"j/k move · enter edit · q quit",
		"↓ 31316432  (no message)  local+origin",
		"›  ↑ 5cb53d15  update terminal theme  @",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RenderModifyInteractive() missing %q in:\n%s", want, got)
		}
	}
	assertAlignedPickerRows(t, got, "31316432", "5cb53d15")
	if strings.Contains(got, "working target") || strings.Contains(got, "local only") {
		t.Fatalf("RenderModifyInteractive() should use compact tags:\n%s", got)
	}
}

func assertAlignedPickerRows(t *testing.T, text, first, second string) {
	t.Helper()
	firstColumn := -1
	secondColumn := -1
	for _, line := range strings.Split(text, "\n") {
		if idx := strings.Index(line, first); idx >= 0 {
			firstColumn = runewidth.StringWidth(line[:idx])
		}
		if idx := strings.Index(line, second); idx >= 0 {
			secondColumn = runewidth.StringWidth(line[:idx])
		}
	}
	if firstColumn < 0 || secondColumn < 0 || firstColumn != secondColumn {
		t.Fatalf("picker rows not aligned for %q/%q: columns %d/%d\n%s", first, second, firstColumn, secondColumn, text)
	}
}
