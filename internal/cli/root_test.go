package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-runewidth"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/publication"
)

var errTestComposeRepair = errors.New("compose repair unavailable")

func TestPrintAddSummaryIncludesHashesSplitAndEditCommands(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	result := authoring.CheckpointResult{
		Change: authoring.ChangeInfo{
			ChangeID:    "zzzzzzchange",
			CommitID:    "abcdef123456",
			Description: "split generated work",
		},
	}
	var out bytes.Buffer

	printAddSummary(&out, result, true)

	text := out.String()
	for _, want := range []string{
		"Revision recorded",
		"Message split generated work",
		"Revision zzzzzzchange",
		"Commit abcdef12",
		"Split recorded the selected changes; remaining edits stay in the current revision",
		"Edit gx edit zzzzzzchange",
		"Next gx add -m \"split generated work\"",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printAddSummary() missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintStacksSummaryUsesCompactBookmarkDesignWithoutDroppingDetails(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	remote := "git@github.com:acme/console.git"
	current := authoring.StackInfo{
		Name:         "waitlist + gx-pr",
		Alias:        "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	stack := authoring.StackSummary{
		Repo:           authoring.RepoInfo{RootPath: "/tmp/console", RemoteURL: &remote},
		Stack:          &current,
		PublishedCount: 1,
		Stacks: []authoring.StackInfo{
			current,
			{
				Name:           "onboarding repo picker",
				Alias:          "onboarding",
				BookmarkName:   "feature/onboarding",
				BaseRef:        "main",
				Status:         "published",
				RevisionCount:  1,
				PublishedCount: 1,
				Revisions: []authoring.RevisionSummary{
					{Index: 1, ChangeID: "onboardabcde", CommitID: "fedcba4321", Description: "onboarding empty state", Published: true},
				},
			},
		},
		Revisions: []authoring.RevisionSummary{
			{Index: 1, ChangeID: "kmpqzsxabcde", CommitID: "1234567890", Description: "waitlist form component", Published: true},
			{Index: 2, ChangeID: "a3f8c12abcde", CommitID: "abcdef1234", Description: "gx-pr payload sync", Active: true},
		},
	}
	unrecorded := authoring.ChangeInfo{ChangeID: "dirtychange", CommitID: "dirtycommit", Description: "dirty scratch", Files: []string{"main.go"}}
	var out bytes.Buffer

	printStatusSummary(&out, stack, &unrecorded)

	text := out.String()
	for _, want := range []string{
		"$ gx stacks",
		"acme/console",
		"● waitlist + gx-pr",
		"waitlist · main · draft · ↑1 · ✓1",
		"waitlist form component",
		"gx-pr payload sync",
		"dirty scratch",
		"Published",
		"onboarding repo picker",
		"onboarding · main · ✓",
		"onboarding empty state",
		"j/k revision · e edit · d diff · Shift+D delete · esc stacks · q quit",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printStatusSummary() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{
		"j/k move stack",
		"j/k select",
		"repo acme/console",
		"stacks 2",
		"ref feature/waitlist",
		"base main",
		"state draft",
		"12345678",
		"abcdef12",
	} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printStatusSummary() should not include %q:\n%s", unwanted, text)
		}
	}
	if strings.Index(text, "dirty scratch") > strings.Index(text, "gx-pr payload sync") {
		t.Fatalf("latest revision should render before older revisions:\n%s", text)
	}
}

func TestRenderStacksSummaryUsesDisplayFallbackWhenOnlyRevisionsAreKnown(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	branch := "feature/change-kxwqpvuo"
	main := "main"
	stack := authoring.StackSummary{
		Repo: authoring.RepoInfo{RootPath: "/tmp/console", BranchName: &branch, DefaultBranch: &main},
		Revisions: []authoring.RevisionSummary{
			{Index: 1, ChangeID: "a3f8c12abcde", CommitID: "abcdef1234", Description: "gx-pr payload sync", Active: true},
		},
	}

	text := renderStacksSummary(stack, nil, 0, false, 0)
	for _, want := range []string{
		"$ gx stacks",
		"console",
		"● feature/change-kxwqpvuo",
		"feature/change-kxwqpvuo  main · draft · ↑1",
		"gx-pr payload sync",
		"j/k revision · e edit · d diff · Shift+D delete · esc stacks · q quit",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("renderStacksSummary() missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "change-kxwqpvuo · main") {
		t.Fatalf("renderStacksSummary() should not repeat the stack name in metadata:\n%s", text)
	}
}

func TestRenderSelectableRevisionLinesKeepsColumnsAligned(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	stack := authoring.StackSummary{
		Revisions: []authoring.RevisionSummary{
			{ChangeID: "olderchange", Description: "older"},
			{ChangeID: "latestchange", Description: "latest"},
		},
	}
	selected := renderSelectableRevisionLines(stack, nil, 0)
	unselected := renderSelectableRevisionLines(stack, nil, 1)
	if len(selected) != 2 || len(unselected) != 2 {
		t.Fatalf("renderSelectableRevisionLines() produced selected=%d unselected=%d lines, want 2 each", len(selected), len(unselected))
	}
	selectedColumn := displayColumn(selected[0], "latest")
	unselectedColumn := displayColumn(unselected[0], "latest")
	if selectedColumn < 0 || unselectedColumn < 0 || selectedColumn != unselectedColumn {
		t.Fatalf("selected revision id moved columns:\n%s\n%s", selected[0], unselected[0])
	}
}

func displayColumn(line string, needle string) int {
	index := strings.Index(line, needle)
	if index < 0 {
		return -1
	}
	return runewidth.StringWidth(line[:index])
}

func TestStacksModelDShowsDiffForSelectedRevision(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	model := newStacksModel(authoring.StackSummary{
		Stack:  &current,
		Stacks: []authoring.StackInfo{current},
		Revisions: []authoring.RevisionSummary{
			{ChangeID: "oldchange", Description: "old"},
			{ChangeID: "latestchange", Description: "latest", Active: true},
		},
	}, nil)
	model.mode = stacksModeRevisions

	if model.revCursor != 0 {
		t.Fatalf("newStacksModel().revCursor = %d, want latest revision at top", model.revCursor)
	}
	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'd', Text: "d"}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if updated.mode != stacksModeDiff {
		t.Fatalf("Update(d).mode = %v, want diff", updated.mode)
	}
	if updated.action != (stacksAction{}) {
		t.Fatalf("Update(d).action = %#v, want no exit action", updated.action)
	}
}

func TestStacksModelEnterShowsSelectedStackRevisions(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	other := authoring.StackInfo{
		Name:         "docs",
		BookmarkName: "feature/docs",
		BaseRef:      "main",
		Status:       "draft",
		Revisions: []authoring.RevisionSummary{
			{ChangeID: "old-docs-change", Description: "old docs"},
			{ChangeID: "latest-docs-change", Description: "latest docs"},
		},
	}
	model := newStacksModel(authoring.StackSummary{
		Stack:  &current,
		Stacks: []authoring.StackInfo{current, other},
	}, nil)
	model.mode = stacksModeStacks
	model.stackCursor = 1

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if cmd != nil {
		t.Fatalf("Update(enter) returned command for selected stack")
	}
	if updated.mode != stacksModeRevisions {
		t.Fatalf("Update(enter).mode = %v, want revisions", updated.mode)
	}
	if got := updated.selectedRevision(); got != "latest-docs-change" {
		t.Fatalf("selectedRevision() = %q, want latest selected stack revision", got)
	}
}

func TestStacksModelEnterUsesCurrentStackEntryRevisions(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
		Revisions: []authoring.RevisionSummary{
			{ChangeID: "entry-change", Description: "entry revision"},
		},
	}
	model := newStacksModel(authoring.StackSummary{
		Stack:  &current,
		Stacks: []authoring.StackInfo{current},
	}, nil)
	model.mode = stacksModeStacks

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if updated.mode != stacksModeRevisions {
		t.Fatalf("Update(enter).mode = %v, want revisions", updated.mode)
	}
	if got := updated.selectedRevision(); got != "entry-change" {
		t.Fatalf("selectedRevision() = %q, want current stack entry revision", got)
	}
}

func TestJjDiffArgsPreserveColorForInteractiveDiffs(t *testing.T) {
	args := jjDiffArgs("change123")
	want := []string{"diff", "-r", "change123", "--color=always"}
	if strings.Join(args, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("jjDiffArgs() = %#v, want %#v", args, want)
	}
}

func TestColorizeDiffLeavesAnsiColoredOutputUntouched(t *testing.T) {
	diff := "\x1b[38;5;3mModified regular file main.go:\x1b[39m\n"
	if got := colorizeDiff(diff); got != diff {
		t.Fatalf("colorizeDiff() changed ANSI-colored diff:\n%q\nwant:\n%q", got, diff)
	}
}

func TestStacksModelShiftDDeletesActiveRevision(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	model := newStacksModel(authoring.StackSummary{
		Stack:  &current,
		Stacks: []authoring.StackInfo{current},
		Revisions: []authoring.RevisionSummary{
			{ChangeID: "oldchange", Description: "old"},
			{ChangeID: "activechange", Description: "active", Active: true},
		},
	}, nil)
	model.mode = stacksModeRevisions

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'D', Text: "D"}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if updated.action != (stacksAction{Kind: "delete-revision", Target: "activechange"}) {
		t.Fatalf("Update(D).action = %#v, want active revision delete", updated.action)
	}
}

func TestStacksModelShiftDDeletesSelectedStack(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	model := newStacksModel(authoring.StackSummary{
		Stack: &current,
		Stacks: []authoring.StackInfo{
			current,
			{Name: "docs", BookmarkName: "feature/docs", BaseRef: "main", Status: "draft"},
		},
	}, nil)
	model.mode = stacksModeStacks
	model.stackCursor = 1

	next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: 'D', Text: "D"}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if updated.action != (stacksAction{Kind: "delete-stack", Target: "feature/docs"}) {
		t.Fatalf("Update(D).action = %#v, want selected stack delete", updated.action)
	}
}

func TestStacksModelEnterDoesNotSwitchCurrentStack(t *testing.T) {
	current := authoring.StackInfo{
		Name:         "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	model := newStacksModel(authoring.StackSummary{
		Stack:  &current,
		Stacks: []authoring.StackInfo{current},
	}, nil)
	model.mode = stacksModeStacks

	next, cmd := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	updated, ok := next.(stacksModel)
	if !ok {
		t.Fatalf("Update() = %T, want stacksModel", next)
	}
	if cmd != nil {
		t.Fatalf("Update(enter) returned command for current stack")
	}
	if updated.action != (stacksAction{}) {
		t.Fatalf("Update(enter).action = %#v, want no-op for current stack", updated.action)
	}
}

func TestPrintCurrentStatusHumanExplainsMessageAssignment(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	main := "main"
	status := currentStatus{
		Repo: authoring.RepoInfo{DefaultBranch: &main},
		Stack: &authoring.StackInfo{
			Name:    "change-kxwqpvuo",
			BaseRef: "main",
		},
		Current: authoring.ChangeInfo{
			ChangeID: "rzlmkskwokkzwnqpxttqstlwrpvlmttl",
			CommitID: "73f3245b3523",
			Files:    []string{"internal/cli/root.go", "internal/cli/stacks_tui.go"},
		},
		NeedsMessage:  true,
		Files:         []string{"internal/cli/root.go", "internal/cli/stacks_tui.go"},
		Next:          []string{"gx compose -a", "gx stacks"},
		GitStatusNote: "gx stores new changes in revisions, so `git status` may be clean.",
	}
	var out bytes.Buffer

	printCurrentStatusHuman(&out, status)
	text := out.String()
	for _, want := range []string{
		"2 files are currently waiting to be assigned.",
		"internal/cli/root.go",
		"internal/cli/stacks_tui.go",
		"gx stores new changes in revisions, so `git status` may be clean.",
		"Next",
		"gx compose -a",
		"gx stacks",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printCurrentStatusHuman() missing %q in:\n%s", want, text)
		}
	}

	status.Current.Description = "update stack navigation"
	status.NeedsMessage = false
	out.Reset()
	printCurrentStatusHuman(&out, status)
	text = out.String()
	for _, unwanted := range []string{"Commit ", "Parent ", "Revision ", "Message:"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printCurrentStatusHuman() should not include %q:\n%s", unwanted, text)
		}
	}
}

func TestPrintCurrentStatusHumanPointsDirtyEditModeToAdd(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	main := "main"
	status := currentStatus{
		Repo: authoring.RepoInfo{DefaultBranch: &main},
		Stack: &authoring.StackInfo{
			Name:         "change-kxwqpvuo",
			BookmarkName: "feature/change-kxwqpvuo",
			BaseRef:      "main",
		},
		Refs: currentStatusRefs{
			GXBaseRef:      "main",
			GXStackRef:     "feature/change-kxwqpvuo",
			GitCheckoutRef: "feature/change-kxwqpvuo",
		},
		Current: authoring.ChangeInfo{
			ChangeID: "rzlmkskwokkzwnqpxttqstlwrpvlmttl",
			CommitID: "73f3245b3523",
			Files:    []string{"internal/cli/root.go", "internal/cli/style.go"},
		},
		NeedsMessage: true,
		Files:        []string{"internal/cli/root.go", "internal/cli/style.go"},
		Next:         []string{`gx add -m "describe this revision"`, "gx stacks"},
	}
	var out bytes.Buffer

	printCurrentStatusHuman(&out, status)
	text := out.String()
	for _, want := range []string{
		"gx add -m \"describe this revision\"",
		"gx stacks",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("dirty edit status missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "pending compose proposal") {
		t.Fatalf("dirty edit status should not point at compose:\n%s", text)
	}
}

func TestPrintCurrentStatusHumanShowsReportPromptForUploadError(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	status := currentStatus{
		PublishUploads: publication.QueueStatus{LastError: `POST "https://api.gx.run/v1/publish": tls: failed to verify certificate`},
		Files:          []string{"internal/github/client_test.go"},
		Next:           []string{"gx compose -a", "gx stacks"},
		GitStatusNote:  "gx stores new changes in revisions, so `git status` may be clean.",
	}
	var out bytes.Buffer

	printCurrentStatusHuman(&out, status)
	text := out.String()
	for _, want := range []string{
		`ERROR: POST "https://api.gx.run/v1/publish": tls: failed to verify certificate`,
		"Run `gx report` to report this issue.",
		"internal/github/client_test.go",
		"gx compose -a",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printCurrentStatusHuman() missing %q in:\n%s", want, text)
		}
	}
}

func TestRenderStacksInteractiveShowsNavigationHintAndCursor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{Name: "waitlist + gx-pr", Alias: "waitlist", BookmarkName: "feature/waitlist", BaseRef: "main", Status: "draft"}
	stack := authoring.StackSummary{
		Repo:  authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack: &current,
		Stacks: []authoring.StackInfo{
			current,
			{
				Name:         "device auth polish",
				Alias:        "device-auth",
				BookmarkName: "feature/device-auth",
				BaseRef:      "main",
				Status:       "draft",
				Revisions: []authoring.RevisionSummary{
					{ChangeID: "deviceauthlatest", Description: "device auth latest"},
					{ChangeID: "deviceautholder", Description: "device auth older"},
				},
			},
		},
		Revisions: []authoring.RevisionSummary{
			{Index: 1, ChangeID: "a3f8c12abcde", CommitID: "abcdef1234", Description: "gx-pr payload sync", Active: true},
		},
	}

	first := renderStacksSummary(stack, nil, 0, false, 0)
	for _, want := range []string{"j/k revision · e edit · d diff · Shift+D delete · esc stacks · q quit", "● waitlist + gx-pr", "gx-pr payload sync", "latest"} {
		if !strings.Contains(first, want) {
			t.Fatalf("renderStacksSummary(cursor 0) missing %q in:\n%s", want, first)
		}
	}

	second := renderStacksSummary(stack, nil, 1, true, 0)
	for _, want := range []string{"● device auth polish", "device-auth · main · draft", "○ waitlist + gx-pr"} {
		if !strings.Contains(second, want) {
			t.Fatalf("renderStacksSummary(cursor 1) missing %q in:\n%s", want, second)
		}
	}

	third := renderStacksSummary(stack, nil, 1, false, 0)
	for _, want := range []string{"● device auth polish", "›   ↑ devicea", "device auth older"} {
		if !strings.Contains(third, want) {
			t.Fatalf("renderStacksSummary(non-current revision mode) missing %q in:\n%s", want, third)
		}
	}
}

func TestStacksViewportKeepsHeaderFooterAndShowsScrollIndicators(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{
		Name:         "stack 0",
		Alias:        "s0",
		BookmarkName: "feature/stack-0",
		BaseRef:      "main",
		Status:       "draft",
	}
	summary := authoring.StackSummary{
		Repo:      authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack:     &current,
		Stacks:    []authoring.StackInfo{current},
		Revisions: []authoring.RevisionSummary{{ChangeID: "change0", Description: "revision 0"}},
	}
	for i := 1; i < 8; i++ {
		summary.Stacks = append(summary.Stacks, authoring.StackInfo{
			Name:         fmt.Sprintf("stack %d", i),
			Alias:        fmt.Sprintf("s%d", i),
			BookmarkName: fmt.Sprintf("feature/stack-%d", i),
			BaseRef:      "main",
			Status:       "draft",
			Revisions:    []authoring.RevisionSummary{{ChangeID: fmt.Sprintf("change%d", i), Description: fmt.Sprintf("revision %d", i)}},
		})
	}
	model := newStacksModel(summary, nil)
	next, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
	model = next.(stacksModel)

	text := model.View().Content
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if lines[0] != "$ gx stacks" {
		t.Fatalf("top line = %q, want gx stacks header in:\n%s", lines[0], text)
	}
	if !strings.Contains(lines[len(lines)-1], "j/k stack") {
		t.Fatalf("last line should be legend, got %q in:\n%s", lines[len(lines)-1], text)
	}
	if !strings.Contains(text, "... more below") {
		t.Fatalf("viewport missing lower scroll indicator:\n%s", text)
	}

	for i := 0; i < 6; i++ {
		next, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
		model = next.(stacksModel)
	}
	text = model.View().Content
	if !strings.Contains(text, "... more above") {
		t.Fatalf("viewport missing upper scroll indicator after moving down:\n%s", text)
	}
	if !strings.Contains(text, "stack 6") {
		t.Fatalf("viewport should keep selected stack visible:\n%s", text)
	}
	lines = strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if !strings.Contains(lines[len(lines)-1], "j/k stack") {
		t.Fatalf("last line should remain legend after scroll, got %q in:\n%s", lines[len(lines)-1], text)
	}
}

func TestStacksDisplayHidesMergedAndKeepsPublishedAtBottomByDefault(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{Name: "active work", Alias: "active", BookmarkName: "feature/active", BaseRef: "main", Status: "draft"}
	stack := authoring.StackSummary{
		Repo:  authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack: &current,
		Stacks: []authoring.StackInfo{
			{Name: "published work", Alias: "published", BookmarkName: "feature/published", BaseRef: "main", Status: "published", RevisionCount: 1, PublishedCount: 1, Revisions: []authoring.RevisionSummary{{ChangeID: "publishedchange", Description: "published"}}},
			current,
			{Name: "merged work", Alias: "merged", BookmarkName: "feature/merged", BaseRef: "main", Status: "merged"},
			{Name: "next work", Alias: "next", BookmarkName: "feature/next", BaseRef: "main", Status: "draft", RevisionCount: 1, Revisions: []authoring.RevisionSummary{{ChangeID: "nextchange", Description: "next"}}},
		},
		Revisions: []authoring.RevisionSummary{{ChangeID: "activechange", Description: "active"}},
	}

	display := stackSummaryForStacksDisplay(stack, stackDisplayOptions{})
	text := renderStacksSummary(display, nil, 0, true, 0)

	activeIndex := strings.Index(text, "active work")
	nextIndex := strings.Index(text, "next work")
	publishedIndex := strings.Index(text, "published work")
	if activeIndex == -1 || nextIndex == -1 || publishedIndex == -1 {
		t.Fatalf("renderStacksSummary() missing expected stacks:\n%s", text)
	}
	if strings.Contains(text, "merged work") {
		t.Fatalf("renderStacksSummary() should hide merged stacks by default:\n%s", text)
	}
	if publishedIndex < activeIndex || publishedIndex < nextIndex {
		t.Fatalf("published stack should render after draft stacks:\n%s", text)
	}
}

func TestStacksDisplayShowAllStillHidesMerged(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{Name: "active work", Alias: "active", BookmarkName: "feature/active", BaseRef: "main", Status: "draft"}
	stack := authoring.StackSummary{
		Repo:  authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack: &current,
		Stacks: []authoring.StackInfo{
			{Name: "merged work", Alias: "merged", BookmarkName: "feature/merged", BaseRef: "main", Status: "merged", RevisionCount: 2, PublishedCount: 2},
			{Name: "published work", Alias: "published", BookmarkName: "feature/published", BaseRef: "main", Status: "published", RevisionCount: 1, PublishedCount: 1, Revisions: []authoring.RevisionSummary{{ChangeID: "publishedchange", Description: "published"}}},
			current,
		},
		Revisions: []authoring.RevisionSummary{{ChangeID: "activechange", Description: "active"}},
	}

	display := stackSummaryForStacksDisplay(stack, stackDisplayOptions{ShowAll: true})
	text := renderStacksSummary(display, nil, 0, true, 0)

	activeIndex := strings.Index(text, "active work")
	publishedIndex := strings.Index(text, "published work")
	if activeIndex == -1 || publishedIndex == -1 {
		t.Fatalf("renderStacksSummary() missing expected stacks:\n%s", text)
	}
	if strings.Contains(text, "merged work") {
		t.Fatalf("show-all should still hide merged stacks:\n%s", text)
	}
	if publishedIndex < activeIndex {
		t.Fatalf("published stack should render after draft stacks:\n%s", text)
	}
}

func TestStacksDisplayHidesEmptyStacksAndShowsListNotice(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{Name: "active work", Alias: "active", BookmarkName: "feature/active", BaseRef: "main", Status: "draft"}
	stack := authoring.StackSummary{
		Repo:  authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack: &current,
		Stacks: []authoring.StackInfo{
			current,
			{Name: "empty docs", Alias: "docs", BookmarkName: "feature/docs", BaseRef: "main", Status: "draft"},
			{Name: "empty cli", Alias: "cli", BookmarkName: "feature/cli", BaseRef: "main", Status: "draft"},
		},
		Revisions: []authoring.RevisionSummary{{ChangeID: "activechange", Description: "active"}},
	}

	display := stackDisplaySummaryForStacksDisplay(stack, stackDisplayOptions{})
	text := renderStacksSummaryWithHidden(display.Stack, nil, 0, true, 0, display.HiddenEmpty)

	if display.HiddenEmpty != 2 {
		t.Fatalf("HiddenEmpty = %d, want 2", display.HiddenEmpty)
	}
	for _, want := range []string{
		"2 branches without revisions. Run 'gx stacks list' to see a full list of stacks.",
		"active work",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("renderStacksSummaryWithHidden() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"empty docs", "empty cli"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("renderStacksSummaryWithHidden() should hide %q:\n%s", unwanted, text)
		}
	}

	full := stackSummaryForStacksDisplay(stack, stackDisplayOptions{ShowEmpty: true})
	fullText := renderStacksSummary(full, nil, 0, true, 0)
	for _, want := range []string{"empty docs", "empty cli"} {
		if !strings.Contains(fullText, want) {
			t.Fatalf("gx stacks list display missing %q in:\n%s", want, fullText)
		}
	}
}

func TestPrintStatusAgentListsStacksAndTargetDetails(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{
		Name:         "waitlist + gx-pr",
		Alias:        "waitlist",
		BookmarkName: "feature/waitlist",
		BaseRef:      "main",
		Status:       "draft",
	}
	stack := authoring.StackSummary{
		Repo:           authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack:          &current,
		PublishedCount: 0,
		Stacks: []authoring.StackInfo{
			current,
			{Name: "device auth polish", Alias: "device-auth", BookmarkName: "feature/device-auth", BaseRef: "main", Status: "draft"},
		},
		Revisions: []authoring.RevisionSummary{
			{Index: 1, ChangeID: "a3f8c12abcde", CommitID: "abcdef1234", Description: "gx-pr payload sync", Active: true},
		},
	}
	var out bytes.Buffer

	printStatusAgent(&out, stack, "")
	text := out.String()
	for _, want := range []string{
		"repo=console stacks=2 current=\"waitlist + gx-pr\"",
		"stack name=\"waitlist + gx-pr\" alias=waitlist ref=feature/waitlist base=main status=draft changes=1 published=0/1 sync=local current=true",
		"revision id=a3f8c12abcde message=\"gx-pr payload sync\" status=draft commit=abcdef12 active=true",
		"stack name=\"device auth polish\" alias=device-auth ref=feature/device-auth base=main status=draft changes=0 published=0/0 sync=local current=false",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printStatusAgent() missing %q in:\n%s", want, text)
		}
	}

	out.Reset()
	printStatusAgent(&out, stack, "waitlist")
	text = out.String()
	for _, want := range []string{
		"stack=\"waitlist + gx-pr\" repo=console",
		"stack alias=waitlist ref=feature/waitlist base=main status=draft changes=1 published=0/1 sync=local current=true",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printStatusAgent(target) missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintModifySummaryKeepsTargetAndNextCommand(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	result := authoring.ModifyResult{
		Output:        "Working copy now at abc",
		CurrentChange: authoring.ChangeInfo{ChangeID: "qpvkpnrmtxyp", CommitID: "abcdef1234", Description: "Review demux hunk coverage"},
		Stack:         &authoring.StackInfo{Name: "waitlist + gx-pr", Alias: "waitlist", BookmarkName: "feature/waitlist", BaseRef: "main", Status: "draft"},
	}
	var out bytes.Buffer

	printModifySummary(&out, result)

	text := out.String()
	for _, want := range []string{
		"Working copy now at abc",
		"revision qpvkpnrmtxyp",
		"commit abcdef12",
		"message Review demux hunk coverage",
		"stack waitlist + gx-pr",
		"next gx compose",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printModifySummary() missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintDemuxProposalHidesDiagnosticsByDefault(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	proposal := authoring.DemuxProposal{
		ID:               "demux-1",
		ProposedChangeID: "change-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "split generated work", Files: []string{"alpha.go"}},
		},
		FeasibilityWarnings: []authoring.FeasibilityWarning{
			{
				RevisionID: "r1",
				Severity:   "warning",
				Source:     "structural_dependency",
				Message:    "alpha.go references Beta from beta.go, but r2 is proposed after it",
			},
		},
		Warnings: []string{"no structural dependency edges found"},
	}
	var out bytes.Buffer

	printDemuxProposal(&out, proposal, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"Compose proposal",
		"Found 1 file",
		"Proposes 1 revisions",
		"Status needs review",
		"Revisions",
		"r1 split generated work",
		"Diagnostics 2 hidden; use --raw or --json",
		"JSON gx compose --json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxProposal() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{
		"[structural_dependency / r1]",
		"alpha.go references Beta from beta.go, but r2 is proposed after it",
		"no structural dependency edges found",
	} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printDemuxProposal() should hide %q by default:\n%s", unwanted, text)
		}
	}
}

func TestPrintDemuxProposalGroupsRevisionsByStack(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	proposal := authoring.DemuxProposal{
		ID:               "demux-1",
		ProposedChangeID: "change-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "stack storage", Files: []string{"internal/storage/schema.sql"}, TargetStack: "feature/stack-management"},
			{ID: "r2", Intent: "route planner", Files: []string{"internal/authoring/demux_routing.go"}, TargetStack: "feature/demux-routing"},
			{ID: "r3", Intent: "stack CLI", Files: []string{"internal/cli/root.go"}, TargetStack: "feature/stack-management"},
		},
	}
	var out bytes.Buffer

	printDemuxProposal(&out, proposal, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"Stacks",
		"● feature/stack-management  2 revisions",
		"r1 stack storage",
		"r3 stack CLI",
		"● feature/demux-routing  1 revision",
		"r2 route planner",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxProposal() missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Revisions") {
		t.Fatalf("stacked proposal should use Stacks section:\n%s", text)
	}
}

func TestPrintDemuxProposalRawIncludesDiagnostics(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	proposal := authoring.DemuxProposal{
		ID:               "demux-1",
		ProposedChangeID: "change-1",
		Revisions: []authoring.RevisionProposal{
			{ID: "r1", Intent: "split generated work", Files: []string{"alpha.go"}},
		},
		FeasibilityWarnings: []authoring.FeasibilityWarning{
			{
				RevisionID: "r1",
				Severity:   "warning",
				Source:     "structural_dependency",
				Message:    "alpha.go references Beta from beta.go, but r2 is proposed after it",
			},
		},
		Warnings: []string{"no structural dependency edges found"},
	}
	var out bytes.Buffer

	printDemuxProposal(&out, proposal, demuxPrintOptions{Raw: true})

	text := out.String()
	for _, want := range []string{
		"Warnings",
		"alpha.go references Beta from beta.go, but r2 is proposed after it",
		"[structural_dependency / r1]",
		"no structural dependency edges found",
		"JSON gx compose --json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxProposal() missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "hidden; use --raw or --json") {
		t.Fatalf("raw print should not include hidden diagnostics summary:\n%s", text)
	}
}

func TestPrintDemuxChangesPacketIncludesWorkflowHeader(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	packet := authoring.DemuxPlanPacket{
		State:    authoring.DemuxWorkflowReadyToApply,
		NextTool: "gx_apply_revision_plan",
		Proposal: authoring.DemuxProposal{
			ID:               "demux-1",
			ProposedChangeID: "change-1",
			Revisions: []authoring.RevisionProposal{
				{ID: "r1", Intent: "split generated work", Files: []string{"alpha.go"}},
			},
		},
	}
	var out bytes.Buffer

	printDemuxChangesPacket(&out, packet, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"$ gx compose",
		"Compose proposal",
		"Found 1 file",
		"Revisions",
		"Accept gx compose apply demux-1",
		"JSON gx compose --json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxChangesPacket() missing %q in:\n%s", want, text)
		}
	}
}

func TestDemuxHunksForRevisionsIncludesWholeFileHunksWithHunkIDs(t *testing.T) {
	proposal := authoring.DemuxProposal{
		Hunks: []authoring.HunkRange{
			{ID: "h1", File: ".env.example", Patch: "env patch"},
			{ID: "h2", File: "server/src/routes/openai.ts", Patch: "openai patch"},
		},
	}
	revisions := []authoring.RevisionProposal{
		{ID: "r1", Intent: "env hunk", UseHunks: true, HunkIDs: []string{"h1"}},
		{ID: "u61", Intent: "update api route", Files: []string{"server/src/routes/openai.ts"}, TargetStack: "feature/server-src"},
	}

	hunks := demuxHunksForRevisions(proposal, revisions)

	if len(hunks) != 2 {
		t.Fatalf("demuxHunksForRevisions() returned %d hunks, want 2: %#v", len(hunks), hunks)
	}
	ids := map[string]bool{}
	for _, hunk := range hunks {
		ids[hunk.ID] = true
	}
	if !ids["h1"] || !ids["h2"] {
		t.Fatalf("demuxHunksForRevisions() ids = %#v, want h1 and h2", ids)
	}
}

func TestDemuxErrorHasInteractiveProposalRequiresUsableHumanProposal(t *testing.T) {
	packet := authoring.DemuxPlanPacket{
		Proposal: authoring.DemuxProposal{ID: "demux-1"},
	}

	if !demuxErrorHasInteractiveProposal(packet, false, false) {
		t.Fatalf("demuxErrorHasInteractiveProposal() = false, want true")
	}
	if demuxErrorHasInteractiveProposal(packet, true, false) {
		t.Fatalf("raw output should not enter interactive compose")
	}
	if demuxErrorHasInteractiveProposal(packet, false, true) {
		t.Fatalf("plan-only output should not enter interactive compose")
	}
	if demuxErrorHasInteractiveProposal(authoring.DemuxPlanPacket{}, false, false) {
		t.Fatalf("missing proposal should not enter interactive compose")
	}
}

func TestDemuxAutoAcceptBlockedReasonRequiresCleanReadyProposal(t *testing.T) {
	ready := authoring.DemuxPlanPacket{
		State: authoring.DemuxWorkflowReadyToApply,
		Proposal: authoring.DemuxProposal{
			ID: "demux-1",
			Revisions: []authoring.RevisionProposal{
				{ID: "r1", Intent: "ready"},
			},
		},
	}
	if reason := demuxAutoAcceptBlockedReason(ready); reason != "" {
		t.Fatalf("ready proposal blocked: %s", reason)
	}

	withWarning := ready
	withWarning.Proposal.FeasibilityWarnings = []authoring.FeasibilityWarning{{Severity: "warning", Message: "blocked"}}
	if reason := demuxAutoAcceptBlockedReason(withWarning); !strings.Contains(reason, "blocking warnings") {
		t.Fatalf("warning reason = %q, want blocking warnings", reason)
	}

	withRepairHint := ready
	withRepairHint.Review.RepairHints = []authoring.RepairHint{{Kind: "split"}}
	if reason := demuxAutoAcceptBlockedReason(withRepairHint); !strings.Contains(reason, "repair hints") {
		t.Fatalf("repair reason = %q, want repair hints", reason)
	}
}

func TestPrintDemuxChangesPacketShowsPartialAcceptAndFollowup(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	packet := authoring.DemuxPlanPacket{
		State: authoring.DemuxWorkflowReadyToApply,
		Proposal: authoring.DemuxProposal{
			ID: "demux-partial",
			Revisions: []authoring.RevisionProposal{
				{ID: "r1", Intent: "alpha", Files: []string{"alpha.go"}},
			},
			Warnings: []string{demuxPartialComposeWarning},
		},
	}
	var out bytes.Buffer

	printDemuxChangesPacket(&out, packet, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"Partial proposal",
		demuxPartialComposeWarning,
		"Accept gx compose apply demux-partial",
		"Then gx compose",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxChangesPacket() missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintDemuxReviewPacketSummarizesWarningsAndRepairHintsByDefault(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	packet := authoring.DemuxPlanPacket{
		State: authoring.DemuxWorkflowRepairRequired,
		Proposal: authoring.DemuxProposal{
			ID: "demux-1",
			Revisions: []authoring.RevisionProposal{
				{ID: "u1", Intent: "app", Files: []string{"app.go"}},
				{ID: "u2", Intent: "helper", Files: []string{"helper.go"}},
			},
			FeasibilityWarnings: []authoring.FeasibilityWarning{{
				RevisionID: "u1",
				Severity:   "warning",
				Source:     "structural_dependency",
				Message:    "app.go depends on helper.go",
			}},
			Warnings: []string{"generic diagnostic"},
		},
		Review: authoring.ReviewDemuxResult{
			Valid: true,
			RepairHints: []authoring.RepairHint{{
				Kind:       "reorder_dependency",
				RevisionID: "u1",
				DependsOn:  "u2",
				Suggestion: "move u2 before u1",
			}},
		},
	}
	var out bytes.Buffer

	printDemuxReviewPacket(&out, packet, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"Compose proposal",
		"Found 2 files",
		"Proposes 2 revisions",
		"Status needs review",
		"Repair hints 1 hidden; use --raw or --json",
		"Blocking warnings",
		"u1  app.go depends on helper.go",
		"Next",
		"gx compose fix demux-1",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxReviewPacket() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{
		"reorder_dependency - u1 - depends on u2 - move u2 before u1",
		"generic diagnostic",
	} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printDemuxReviewPacket() should hide %q by default:\n%s", unwanted, text)
		}
	}

	out.Reset()
	printDemuxReviewPacket(&out, packet, demuxPrintOptions{Raw: true})
	text = out.String()
	for _, want := range []string{
		"Repair hints",
		"reorder_dependency - u1 - depends on u2 - move u2 before u1",
		"Blocking warnings",
		"app.go depends on helper.go",
		"[structural_dependency / u1]",
		"Diagnostics",
		"generic diagnostic",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxReviewPacket(raw) missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "hidden; use --raw or --json") {
		t.Fatalf("printDemuxReviewPacket(raw) should not hide diagnostics:\n%s", text)
	}
}

func TestPrintDemuxReviewPacketShowsInfoWarningsInReviewScreen(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	packet := authoring.DemuxPlanPacket{
		State: authoring.DemuxWorkflowReadyToApply,
		Proposal: authoring.DemuxProposal{
			ID:        "demux-1",
			Revisions: []authoring.RevisionProposal{{ID: "u1", Intent: "alpha"}},
			FeasibilityWarnings: []authoring.FeasibilityWarning{{
				RevisionID: "u1",
				Severity:   "info",
				Source:     "changed_symbol",
				Message:    "hunk h1 is not mapped to an enclosing symbol",
			}},
		},
		Review: authoring.ReviewDemuxResult{Valid: true},
	}
	var out bytes.Buffer

	printDemuxReviewPacket(&out, packet, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{"Compose proposal", "Diagnostics 1 hidden; use --raw or --json", "Next", "gx compose"} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxReviewPacket() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"Status needs review", "Blocking warnings", "hunk h1 is not mapped to an enclosing symbol"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printDemuxReviewPacket() should not show %q for info diagnostics:\n%s", unwanted, text)
		}
	}

	out.Reset()
	printDemuxReviewPacket(&out, packet, demuxPrintOptions{Raw: true})
	text = out.String()
	for _, want := range []string{"Diagnostics", "not mapped to an enclosing symbol", "[changed_symbol / u1]"} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxReviewPacket(raw) missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Blocking warnings") {
		t.Fatalf("printDemuxReviewPacket(raw) should not show info diagnostics as blocking:\n%s", text)
	}
}

func TestPrintDemuxProposalListShowsLatestPendingMarker(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	proposals := []authoring.DemuxProposalSummary{
		{
			ID:                    "demux-new",
			Status:                authoring.ProposalPending,
			CreatedAt:             1700000000000,
			RevisionCount:         2,
			Files:                 []string{"alpha.go", "beta.go"},
			HiddenDiagnosticCount: 3,
			FirstRevisionIntent:   "update alpha",
			LatestPendingForShow:  true,
			Alias:                 "d1",
		},
		{
			ID:            "demux-old",
			Status:        authoring.ProposalPending,
			CreatedAt:     1699990000000,
			RevisionCount: 1,
			Files:         []string{"gamma.go"},
			Alias:         "d2",
		},
	}
	var out bytes.Buffer

	printDemuxProposalList(&out, proposals)

	text := out.String()
	for _, want := range []string{
		"Compose proposals",
		"* d1",
		"demux-new",
		"pending / 2 revisions / 2 files / 3 diagnostics",
		"update alpha",
		"default for revisions: gx compose show <revision-id>",
		"  d2",
		"demux-old",
		"pending / 1 revisions / 1 files",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxProposalList() missing %q in:\n%s", want, text)
		}
	}
}

func TestDemuxShowTargetIsProposal(t *testing.T) {
	for _, target := range []string{"d1", "D12", "demux-abc123"} {
		if !demuxShowTargetIsProposal(target) {
			t.Fatalf("demuxShowTargetIsProposal(%q) = false, want true", target)
		}
	}
	for _, target := range []string{"u1", "r1", "docs", "d", "d0x"} {
		if demuxShowTargetIsProposal(target) {
			t.Fatalf("demuxShowTargetIsProposal(%q) = true, want false", target)
		}
	}
}

func TestComposeHelpHidesAgentOnlyCommandsAndRemovesOldAliases(t *testing.T) {
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"compose", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	for _, want := range []string{"list", "review", "fix", "show"} {
		if !strings.Contains(text, want) {
			t.Fatalf("compose help missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"  apply ", "  check ", "  proposal ", "apply-plan", "review-plan", "validate", "changes"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("compose help should hide %q:\n%s", unwanted, text)
		}
	}

	for _, args := range [][]string{{"compose", "review-plan"}, {"compose", "apply-plan"}} {
		cmd, _, err := root.Find(args)
		if err != nil || cmd == nil || !cmd.Hidden {
			t.Fatalf("Find(%v) = cmd=%v err=%v, want hidden command", args, cmd, err)
		}
	}
	for _, args := range [][]string{{"demux"}, {"gxa"}, {"gxt"}, {"pr"}, {"modify"}, {"stack", "--new", "demo"}} {
		cmd, _, err := root.Find(args)
		if err == nil && cmd != nil && cmd.Name() == args[0] {
			t.Fatalf("Find(%v) resolved removed command %q", args, cmd.Name())
		}
	}
}

func TestAddHelpShowsMessageInExamples(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"add", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"A non-empty commit message is required; pass it with -m.",
		"Examples:",
		`gx add -m "describe this revision"`,
		`gx add -i internal/termstyle/theme.go -m "update terminal theme"`,
		`gx add --hunk --patch-file /tmp/selected.patch -m "record selected hunks"`,
		"-m, --message string      commit message (required)",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("add help missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "gx add -i internal/termstyle/theme.go\n") {
		t.Fatalf("add help includes interactive example without message:\n%s", text)
	}
}

func TestPrintDemuxAIReviewResult(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	result := authoring.DemuxAIReviewResult{
		Proposal: authoring.DemuxProposal{
			ID:        "demux-1",
			Revisions: []authoring.RevisionProposal{{ID: "u1"}, {ID: "u2"}},
		},
		Review:          authoring.ReviewDemuxResult{Valid: true},
		State:           authoring.DemuxWorkflowReadyToApply,
		Model:           "test-model",
		Updated:         true,
		Notes:           []string{"moved helper before caller"},
		WarningsSent:    3,
		WarningsTotal:   8,
		RepairHintSent:  2,
		RepairHintTotal: 5,
	}
	var out bytes.Buffer

	printDemuxAIReviewResult(&out, result)

	text := out.String()
	for _, want := range []string{
		"Compose fix",
		"Proposal demux-1",
		"Model test-model",
		"Updated yes",
		"State ready_to_apply",
		"Revisions 2",
		"Blocking issues 0",
		"Sent 3/8 warnings, 2/5 repair hints",
		"moved helper before caller",
		"Accept gx compose",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxAIReviewResult() missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintDemuxAIReviewResultShowsBlockingWarnings(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	result := authoring.DemuxAIReviewResult{
		Proposal: authoring.DemuxProposal{
			ID: "demux-1",
			FeasibilityWarnings: []authoring.FeasibilityWarning{
				{Severity: "info"},
				{Severity: "warning"},
			},
		},
		Review: authoring.ReviewDemuxResult{
			Valid:       true,
			RepairHints: []authoring.RepairHint{{Kind: "reorder_dependency"}},
		},
		State: authoring.DemuxWorkflowRepairRequired,
		Model: "test-model",
	}
	var out bytes.Buffer

	printDemuxAIReviewResult(&out, result)

	text := out.String()
	for _, want := range []string{
		"Updated no",
		"Blocking issues 1",
		"Diagnostics 3 hidden; use --raw or --json",
		"Review gx compose review demux-1 --raw",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxAIReviewResult() missing %q in:\n%s", want, text)
		}
	}
}

func TestPrintDemuxRevisionViewShowsMetadataAndDiff(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	view := authoring.DemuxRevisionView{
		ProposalID:       "demux-1",
		ProposalStatus:   authoring.ProposalPending,
		ProposedChangeID: "change-1",
		Revision: authoring.RevisionProposal{
			ID:               "u1",
			Intent:           "split generated work",
			Files:            []string{"alpha.go"},
			UseHunks:         true,
			HunkIDs:          []string{"h1"},
			ProvenanceStatus: "explicit",
			SessionIDs:       []string{"session-alpha"},
			Confidence:       0.75,
		},
		Hunks: []authoring.HunkRange{{
			ID:    "h1",
			File:  "alpha.go",
			Patch: "diff --git a/alpha.go b/alpha.go\n+++ b/alpha.go\n@@ -1 +1 @@\n-alpha\n+beta\n",
		}},
		FeasibilityWarnings: []authoring.FeasibilityWarning{{
			RevisionID: "u1",
			Severity:   "info",
			Source:     "changed_symbol",
			Message:    "coarser grouping",
		}},
	}
	var out bytes.Buffer

	printDemuxRevisionView(&out, view, demuxPrintOptions{})

	text := out.String()
	for _, want := range []string{
		"Compose revision",
		"Proposal demux-1",
		"Revision u1",
		"Intent split generated work",
		"Mode hunk (1)",
		"Confidence 0.75",
		"Files",
		"alpha.go",
		"Diagnostics 1 hidden; use --raw or --json",
		"Diff",
		"diff --git a/alpha.go b/alpha.go",
		"-alpha",
		"+beta",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxRevisionView() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"coarser grouping", "Provenance explicit session-alpha"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printDemuxRevisionView() should hide %q by default:\n%s", unwanted, text)
		}
	}
}

func TestPrintDemuxRevisionViewRawShowsDiagnostics(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	view := authoring.DemuxRevisionView{
		ProposalID:       "demux-1",
		ProposalStatus:   authoring.ProposalPending,
		ProposedChangeID: "change-1",
		Revision: authoring.RevisionProposal{
			ID:               "u1",
			Intent:           "split generated work",
			Files:            []string{"alpha.go"},
			ProvenanceStatus: "explicit",
			SessionIDs:       []string{"session-alpha"},
		},
		FeasibilityWarnings: []authoring.FeasibilityWarning{{
			RevisionID: "u1",
			Severity:   "warning",
			Source:     "structural_dependency",
			Message:    "alpha depends on beta",
		}},
		StructuralDeps: []authoring.StructuralDependency{{
			FromFile: "alpha.go",
			ToFile:   "beta.go",
			Symbol:   "Beta",
		}},
	}
	var out bytes.Buffer

	printDemuxRevisionView(&out, view, demuxPrintOptions{Raw: true})

	text := out.String()
	for _, want := range []string{
		"Feasibility warnings",
		"Provenance explicit session-alpha",
		"alpha depends on beta",
		"[structural_dependency / u1]",
		"Structural dependencies",
		"alpha.go -> beta.go",
		"Beta",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printDemuxRevisionView(raw) missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "hidden; use --raw or --json") {
		t.Fatalf("raw revision view should not hide diagnostics:\n%s", text)
	}
}

func TestVersionCommandPrintsLabeledVersion(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := strings.TrimSpace(out.String())
	if !strings.HasPrefix(text, "version ") {
		t.Fatalf("version output = %q, want compact version line", text)
	}
	if strings.Contains(text, "$ gx version") {
		t.Fatalf("version output should not echo command:\n%s", text)
	}
}

func TestVersionCommandPrintsJSON(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"version", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var got struct {
		Version  string `json:"version"`
		Release  string `json:"release_version"`
		Revision string `json:"revision"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("version --json output is not JSON: %v\n%s", err, out.String())
	}
	if got.Version == "" {
		t.Fatalf("version --json = %#v, want non-empty version", got)
	}
}

func TestRootHelpPrintsAsciiLogoAtTop(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GX_HOME", t.TempDir())
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.HasPrefix(text, gxLogoRaw+"\n") {
		t.Fatalf("root help should start with logo:\n%s", text)
	}
	if !strings.Contains(text, "$ gx help") {
		t.Fatalf("root help missing gx help invocation:\n%s", text)
	}
	if !strings.Contains(text, "Not signed in  gx auth login") {
		t.Fatalf("root help missing signed-out auth line:\n%s", text)
	}
}

func TestRootHelpShowsSignedInUser(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GX_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "octocat", GitHubAccessToken: "gho_saved", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	if text := out.String(); !strings.Contains(text, "Signed in as octocat") {
		t.Fatalf("root help missing signed-in auth line:\n%s", text)
	}
}

func TestRootHelpIgnoresLegacyLoginWithoutGitHubToken(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GX_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "api-key", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if strings.Contains(text, "Signed in as api-key") {
		t.Fatalf("root help showed stale legacy login:\n%s", text)
	}
	if !strings.Contains(text, "Not signed in  gx auth login") {
		t.Fatalf("root help missing signed-out auth line:\n%s", text)
	}
}

func TestRootHelpShowsStoredLoginWithCloudEnvPresent(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_CLOUD_URL", "http://localhost:3200/gx/pr")
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "joelachance", GitHubAccessToken: "gho_saved", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "Signed in as joelachance") {
		t.Fatalf("root help missing stored login auth line:\n%s", text)
	}
}

func TestOpsCommandPrintsCompactMenu(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"ops"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"Advanced commands for capture, diagnostics, and ingest",
		"  capture     Record a gx session for replay",
		"  diagnose    Inspect local gx state and config",
		"  ingest      Import external change metadata",
		`Use "gx ops [command] --help" for details.`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("gx ops output missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Usage:") {
		t.Fatalf("gx ops output should not render generic help:\n%s", text)
	}
}

func TestOpsDiagnoseAliasStillWorks(t *testing.T) {
	root := NewRoot(context.Background())
	for _, args := range [][]string{{"ops", "diagnose", "doctor"}, {"ops", "diag", "doctor"}} {
		cmd, _, err := root.Find(args)
		if err != nil || cmd == nil || cmd.Name() != "doctor" {
			t.Fatalf("Find(%v) = cmd=%v err=%v, want doctor", args, cmd, err)
		}
	}
}

func TestPublishHelpDoesNotPublish(t *testing.T) {
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"publish", "--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "Publish accepted GX stacks for review") {
		t.Fatalf("help output missing publish summary:\n%s", text)
	}
	if !strings.Contains(text, "gx publish [stack]") {
		t.Fatalf("help output missing stack usage:\n%s", text)
	}
	if !strings.Contains(text, "--all") {
		t.Fatalf("help output missing --all flag:\n%s", text)
	}
	if strings.Contains(text, "--github") {
		t.Fatalf("help output should hide compatibility --github flag:\n%s", text)
	}
	if strings.Contains(text, "--no-github") {
		t.Fatalf("help output should not include removed --no-github flag:\n%s", text)
	}
	if strings.Contains(text, "Pushing ") || strings.Contains(text, "Creating draft PR") {
		t.Fatalf("publish help entered publish path:\n%s", text)
	}
}

func TestPublishGitHubCompatibilityFlagRejectsFalseAndNoGitHubIsRemoved(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"publish"})
	if err != nil {
		t.Fatalf("Find(publish) error = %v", err)
	}
	githubFlag := cmd.Flags().Lookup("github")
	if githubFlag == nil {
		t.Fatal("publish command missing --github flag")
	}
	if githubFlag.DefValue != "true" {
		t.Fatalf("--github default = %q, want true", githubFlag.DefValue)
	}
	noGitHubFlag := cmd.Flags().Lookup("no-github")
	if noGitHubFlag != nil {
		t.Fatal("publish command still registers removed --no-github flag")
	}
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"publish", "--github=false"})
	err = root.Execute()
	if err == nil || !strings.Contains(err.Error(), "--github=false is no longer supported") {
		t.Fatalf("publish --github=false error = %v, want unsupported flag value", err)
	}
}

func TestPRAliasIsRemoved(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"pr"})
	if err == nil && cmd != nil && cmd.Name() == "pr" {
		t.Fatalf("Find(pr) resolved removed alias")
	}
}

func TestReviewCommandUsesDefaults(t *testing.T) {
	root := initGitRepo(t)
	writeTestFile(t, root, "README.md", "# repo\n")
	writeTestFile(t, root, "AGENTS.md", "# agents\n")
	writeTestFile(t, root, "go.mod", "module example.com/repo\n")
	writeTestFile(t, root, "main_test.go", "package main\n")
	gitAddTestFiles(t, root, "README.md", "AGENTS.md", "go.mod", "main_test.go")
	t.Chdir(root)
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"review"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("gx review error = %v\n%s", err, out.String())
	}
	text := out.String()
	for _, want := range []string{
		"## Recommendations",
		"- No recommendations yet.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("gx review output missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"brief", "PR Summary"} {
		if strings.Contains(strings.ToLower(text), unwanted) {
			t.Fatalf("gx review output should not include %q:\n%s", unwanted, text)
		}
	}
}

func TestReviewCommandAcceptsScopeFlag(t *testing.T) {
	root := initGitRepo(t)
	t.Chdir(root)
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"review", "--scope", "architecture"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gx review --scope error = %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "## Recommendations") {
		t.Fatalf("gx review --scope output missing recommendations:\n%s", out.String())
	}
}

func TestReviewCommandRejectsMultiplePrompts(t *testing.T) {
	cmd := NewRoot(context.Background())
	cmd.SetArgs([]string{"review", "one prompt", "second prompt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("gx review accepted multiple positional prompts")
	}
	if !strings.Contains(err.Error(), "accepts at most 1 arg") {
		t.Fatalf("gx review error = %v, want maximum arg error", err)
	}
}

func initGitRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init error = %v\n%s", err, out)
	}
	return root
}

func writeTestFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func gitAddTestFiles(t *testing.T, root string, files ...string) {
	t.Helper()
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add error = %v\n%s", err, out)
	}
}
