package cli

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/satoricorp/gx/internal/authoring"
)

func TestDemuxInteractiveModelAppliesSelectedStack(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", TargetStack: "feature/demux-routing"},
		},
	})

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: 's', Text: "s"}))
	model = next.(demuxInteractiveModel)

	if model.action.Kind != "apply_stack" || model.action.Target != "feature/demux-routing" {
		t.Fatalf("action = %#v, want apply demux routing stack", model.action)
	}
}

func TestDemuxInteractiveModelAppliesAllOnlyFromShiftA(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management"},
		},
	})

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(demuxInteractiveModel)
	if model.mode != demuxInteractiveRevisions {
		t.Fatalf("enter mode = %v, want revisions", model.mode)
	}
	if model.action.Kind != "" {
		t.Fatalf("enter action = %#v, want no action", model.action)
	}

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: 'a', Text: "a"}))
	model = next.(demuxInteractiveModel)
	if model.action.Kind != "" {
		t.Fatalf("a action = %#v, want no action", model.action)
	}

	model = newDemuxInteractiveModel(model.proposal)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: 'A', Text: "A"}))
	model = next.(demuxInteractiveModel)
	if model.action.Kind != "apply_all" {
		t.Fatalf("A action = %#v, want apply_all", model.action)
	}
}

func TestDemuxInteractiveViewShowsAllStacksAndAcceptCallout(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", Files: []string{"storage.go"}, TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", Files: []string{"demux.go"}, TargetStack: "feature/demux-routing"},
		},
	})
	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(demuxInteractiveModel)

	text := model.View().Content
	for _, want := range []string{
		"Stacks",
		"s accept stack",
		"Shift+A accept all",
		"feature/stack-management",
		"r1 stack storage",
		"feature/demux-routing",
		"r2 route planner",
		"j/k stack",
		"● selected",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("View() missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Action required") || strings.Contains(text, "Apply all stacks") || strings.Contains(text, "apply selected") {
		t.Fatalf("View() contains old compose copy:\n%s", text)
	}
}

func TestDemuxViewportKeepsHeaderFooterAndShowsScrollIndicators(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	proposal := authoring.DemuxProposal{ID: "demux-scroll"}
	for i := 0; i < 8; i++ {
		proposal.Revisions = append(proposal.Revisions, authoring.RevisionProposal{
			ID:          fmt.Sprintf("r%d", i),
			Intent:      fmt.Sprintf("revision %d", i),
			TargetStack: fmt.Sprintf("feature/stack-%d", i),
		})
	}
	model := newDemuxInteractiveModel(proposal)
	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
	model = next.(demuxInteractiveModel)

	text := model.View().Content
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if lines[0] != "$ gx compose" {
		t.Fatalf("top line = %q, want gx compose header in:\n%s", lines[0], text)
	}
	if !strings.Contains(lines[len(lines)-1], "j/k stack") {
		t.Fatalf("last line should be legend, got %q in:\n%s", lines[len(lines)-1], text)
	}
	if !strings.Contains(text, "... more below") {
		t.Fatalf("viewport missing lower scroll indicator:\n%s", text)
	}

	for i := 0; i < 6; i++ {
		next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
		model = next.(demuxInteractiveModel)
	}
	text = model.View().Content
	if !strings.Contains(text, "... more above") {
		t.Fatalf("viewport missing upper scroll indicator after moving down:\n%s", text)
	}
	if !strings.Contains(text, "feature/stack-6") {
		t.Fatalf("viewport should keep selected stack visible:\n%s", text)
	}
	lines = strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if !strings.Contains(lines[len(lines)-1], "j/k stack") {
		t.Fatalf("last line should remain legend after scroll, got %q in:\n%s", lines[len(lines)-1], text)
	}
}

func TestDemuxInteractiveModelNavigatesRevisionsAndEscapesToStacks(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management"},
		},
	})

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(demuxInteractiveModel)
	if model.mode != demuxInteractiveRevisions {
		t.Fatalf("mode = %v, want revisions", model.mode)
	}

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc}))
	model = next.(demuxInteractiveModel)
	if model.mode != demuxInteractiveStacks {
		t.Fatalf("mode = %v, want stacks", model.mode)
	}
}

func TestDemuxInteractiveModelSelectsAcrossStacks(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", TargetStack: "feature/demux-routing"},
		},
	})

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: ' ', Text: " "}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	model = next.(demuxInteractiveModel)
	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: ' ', Text: " "}))
	model = next.(demuxInteractiveModel)

	if len(model.selected) != 2 {
		t.Fatalf("selected = %#v, want 2 revisions", model.selected)
	}
}

func TestDemuxInteractiveModelCombinesSelectedRevisions(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", Files: []string{"storage.go"}, HunkIDs: []string{"h1"}, Hunks: []authoring.HunkRange{{ID: "h1", File: "storage.go"}}, TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", Files: []string{"demux.go"}, HunkIDs: []string{"h2"}, Hunks: []authoring.HunkRange{{ID: "h2", File: "demux.go"}}, TargetStack: "feature/stack-management"},
		},
	})
	model.mode = demuxInteractiveRevisions
	model.selected = map[string]struct{}{"r1": {}, "r2": {}}

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'c', Text: "c"}))
	model = next.(demuxInteractiveModel)

	if len(model.proposal.Revisions) != 1 {
		t.Fatalf("revisions = %#v, want 1 combined revision", model.proposal.Revisions)
	}
	combined := model.proposal.Revisions[0]
	if combined.ID != "r1" || !strings.Contains(combined.Intent, "route planner") {
		t.Fatalf("combined = %#v, want r1 with merged intent", combined)
	}
	if len(combined.Files) != 2 || len(combined.Hunks) != 2 {
		t.Fatalf("combined files/hunks = %#v %#v, want merged", combined.Files, combined.Hunks)
	}
}

func TestDemuxInteractiveModelReleasesSelectedRevisions(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", TargetStack: "feature/demux-routing"},
		},
	})
	model.mode = demuxInteractiveRevisions
	model.selected = map[string]struct{}{"r1": {}}

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'r', Text: "r"}))
	model = next.(demuxInteractiveModel)

	if len(model.proposal.Revisions) != 1 || model.proposal.Revisions[0].ID != "r2" {
		t.Fatalf("revisions = %#v, want r2 only", model.proposal.Revisions)
	}
	if len(model.selected) != 0 {
		t.Fatalf("selected = %#v, want cleared", model.selected)
	}
}

func TestDemuxInteractiveModelOpensDiffInPlace(t *testing.T) {
	model := newDemuxInteractiveModel(authoring.DemuxProposal{
		ID: "demux-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", TargetStack: "feature/stack-management", Hunks: []authoring.HunkRange{{ID: "h1", File: "storage.go", Patch: "@@ -1 +1 @@\n-old\n+new\n"}}},
		},
	})
	model.mode = demuxInteractiveRevisions

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'd', Text: "d"}))
	model = next.(demuxInteractiveModel)
	if model.mode != demuxInteractiveDiff {
		t.Fatalf("mode = %v, want diff", model.mode)
	}
	if !strings.Contains(model.View().Content, "@@ -1 +1 @@") {
		t.Fatalf("diff view missing hunk:\n%s", model.View().Content)
	}

	next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc}))
	model = next.(demuxInteractiveModel)
	if model.mode != demuxInteractiveRevisions {
		t.Fatalf("mode = %v, want revisions", model.mode)
	}
}

func TestFilterDemuxProposalForStackKeepsOnlySelectedStack(t *testing.T) {
	proposal := authoring.DemuxProposal{
		ID:               "demux-1",
		ProposedChangeID: "change-1",
		ProposedCommitID: "commit-1",
		Hunks: []authoring.HunkRange{
			{ID: "h1", File: "storage.go"},
			{ID: "h2", File: "demux.go"},
		},
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", Files: []string{"storage.go"}, TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", Files: []string{"demux.go"}, TargetStack: "feature/demux-routing"},
		},
		FeasibilityWarnings: []authoring.FeasibilityWarning{
			{RevisionID: "r1", Severity: "warning", Message: "selected"},
			{RevisionID: "r2", Severity: "warning", Message: "other"},
		},
	}

	filtered := filterDemuxProposalForStack(proposal, "feature/stack-management")

	if filtered.ID != "" {
		t.Fatalf("filtered ID = %q, want new proposal", filtered.ID)
	}
	if filtered.ProposedChangeID != "change-1" || filtered.ProposedCommitID != "commit-1" {
		t.Fatalf("filtered proposal identity = %s/%s, want source change/commit", filtered.ProposedChangeID, filtered.ProposedCommitID)
	}
	if len(filtered.Revisions) != 1 || filtered.Revisions[0].ID != "r1" {
		t.Fatalf("filtered revisions = %#v, want r1 only", filtered.Revisions)
	}
	if len(filtered.Hunks) != 1 || filtered.Hunks[0].ID != "h1" {
		t.Fatalf("filtered hunks = %#v, want h1 only", filtered.Hunks)
	}
	if len(filtered.FeasibilityWarnings) != 1 || filtered.FeasibilityWarnings[0].RevisionID != "r1" {
		t.Fatalf("filtered warnings = %#v, want r1 only", filtered.FeasibilityWarnings)
	}
}
