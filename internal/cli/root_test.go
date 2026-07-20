package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/vcs"
)

var errTestComposeRepair = errors.New("compose repair unavailable")

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
		"$ gx status",
		"acme/console",
		"● waitlist + gx-pr",
		"waitlist · main · draft · ↑1 · ✓1",
		"waitlist form component",
		"gx-pr payload sync",
		"dirty scratch",
		"Remote",
		"onboarding repo picker",
		"onboarding · feature/onboarding · main · remote · ✓1",
		"onboarding empty state",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printStatusSummary() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{
		"j/k up/down",
		"j/k revision",
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

func TestPrintStacksSummaryShowsUnstagedFilesWithoutStacks(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	stack := authoring.StackSummary{
		Repo: authoring.RepoInfo{RootPath: "/tmp/console"},
	}
	unrecorded := authoring.ChangeInfo{
		ChangeID:    "dirtychange",
		CommitID:    "dirtycommit",
		Description: "(no description set)",
		Files:       []string{"README.md", "internal/cli/root.go"},
	}
	var out bytes.Buffer

	printStacksSummary(&out, stack, &unrecorded, 0)

	text := out.String()
	for _, want := range []string{
		"$ gx status",
		"Unstaged",
		"README.md",
		"internal/cli/root.go",
		"unrecorded",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("printStacksSummary() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{
		"No GX revisions recorded yet.",
		"j/k up/down",
	} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("printStacksSummary() should not include %q:\n%s", unwanted, text)
		}
	}
}

func TestStatusNextHints(t *testing.T) {
	clean := vcs.GitWorkingStatus{}
	dirty := vcs.GitWorkingStatus{Unstaged: []string{"a.go"}}
	staged := vcs.GitWorkingStatus{Staged: []string{"a.go"}}
	unpushed := authoring.StackSummary{
		Revisions: []authoring.RevisionSummary{{Published: false}},
	}
	prURL := "https://github.com/acme/gx/pull/7"
	pushed := authoring.StackSummary{
		Stack:     &authoring.StackInfo{GitHubPRURL: &prURL},
		Revisions: []authoring.RevisionSummary{{Published: true}},
	}

	if got := statusNextHints(staged, authoring.StackSummary{}); len(got) != 1 || !strings.Contains(got[0], "gx commit") {
		t.Fatalf("staged hints = %v, want gx commit", got)
	}
	if got := statusNextHints(dirty, authoring.StackSummary{}); len(got) != 1 || got[0] != "git add" {
		t.Fatalf("dirty hints = %v, want git add", got)
	}
	if got := statusNextHints(clean, unpushed); len(got) != 1 || got[0] != vcs.HintPush {
		t.Fatalf("unpushed hints = %v, want %q", got, vcs.HintPush)
	}
	if got := statusNextHints(clean, pushed); len(got) != 1 || got[0] != prURL {
		t.Fatalf("pushed hints = %v, want PR URL", got)
	}
	if got := statusNextHints(clean, authoring.StackSummary{}); got != nil {
		t.Fatalf("clean hints = %v, want nil", got)
	}
}

func TestRenderDefaultStatusSummaryShowsGitSectionsAndNext(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	current := authoring.StackInfo{Name: "waitlist", Alias: "waitlist", BookmarkName: "feature/waitlist", BaseRef: "main", Status: "draft"}
	stack := authoring.StackSummary{
		Repo:  authoring.RepoInfo{RootPath: "/tmp/console"},
		Stack: &current,
		Stacks: []authoring.StackInfo{
			current,
			{Name: "other", Alias: "other", BookmarkName: "feature/other", BaseRef: "main", Status: "draft", RevisionCount: 1},
		},
		Revisions: []authoring.RevisionSummary{{ChangeID: "abc", Description: "payload sync"}},
		GitWorking: vcs.GitWorkingStatus{
			Staged:    []string{"internal/cli/root.go"},
			Untracked: []string{"scratch.txt"},
		},
		Next: []string{`gx commit -m "describe this revision"`},
	}

	text := renderDefaultStatusSummary(stack, nil, 0, false)
	for _, want := range []string{
		"$ gx status",
		"Staged",
		"internal/cli/root.go",
		"Untracked",
		"scratch.txt",
		"waitlist",
		"1 other stacks — run gx status list",
		"Next",
		`gx commit -m "describe this revision"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("renderDefaultStatusSummary() missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"Unstaged", "other · feature/other"} {
		if strings.Contains(text, unwanted) {
			t.Fatalf("renderDefaultStatusSummary() should not include %q:\n%s", unwanted, text)
		}
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
		"$ gx status",
		"console",
		"● feature/change-kxwqpvuo",
		"feature/change-kxwqpvuo  main · draft · ↑1",
		"gx-pr payload sync",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("renderStacksSummary() missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "j/k up/down") || strings.Contains(text, "j/k revision") {
		t.Fatalf("renderStacksSummary() should not include interactive legend:\n%s", text)
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
	if got := updated.selectedRevision(); got != "old-docs-change" {
		t.Fatalf("selectedRevision() = %q, want first selected stack revision", got)
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

func TestColorizeDiffLeavesAnsiColoredOutputUntouched(t *testing.T) {
	diff := "\x1b[38;5;3mModified regular file main.go:\x1b[39m\n"
	if got := colorizeDiff(diff); got != diff {
		t.Fatalf("colorizeDiff() changed ANSI-colored diff:\n%q\nwant:\n%q", got, diff)
	}
}

func TestStacksModelShiftDDoesNotDeleteActiveRevision(t *testing.T) {
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
	if updated.action != (stacksAction{}) {
		t.Fatalf("Update(D).action = %#v, want no action", updated.action)
	}
}

func TestStacksModelShiftDDoesNotDeleteSelectedStack(t *testing.T) {
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
	if updated.action != (stacksAction{}) {
		t.Fatalf("Update(D).action = %#v, want no action", updated.action)
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
		Next:          []string{`git add <files>`, `gx commit -m "describe this revision"`, "gx status"},
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
		"gx commit -m \"describe this revision\"",
		"gx status",
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
		Next:         []string{`git add <files>`, `gx commit -m "describe this revision"`, "gx status"},
	}
	var out bytes.Buffer

	printCurrentStatusHuman(&out, status)
	text := out.String()
	for _, want := range []string{
		"gx commit -m \"describe this revision\"",
		"gx status",
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
		Next:           []string{`git add <files>`, `gx commit -m "describe this revision"`, "gx status"},
		GitStatusNote:  "gx stores new changes in revisions, so `git status` may be clean.",
	}
	var out bytes.Buffer

	printCurrentStatusHuman(&out, status)
	text := out.String()
	for _, want := range []string{
		`ERROR: POST "https://api.gx.run/v1/publish": tls: failed to verify certificate`,
		"Run `gx report` to report this issue.",
		"internal/github/client_test.go",
		"gx commit -m \"describe this revision\"",
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

	first := renderInteractiveStacksSummaryWithHidden(stack, nil, 0, false, 0, 0)
	for _, want := range []string{"↑/↓ navigate · enter open revisions", "● waitlist + gx-pr", "gx-pr payload sync", "latest"} {
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

	text := model.View().Content
	if !strings.Contains(text, "$ gx status") {
		t.Fatalf("missing gx status header in:\n%s", text)
	}
	if !strings.Contains(text, "↑/↓ navigate · enter open revisions · q quit") {
		t.Fatalf("missing legend in:\n%s", text)
	}
	for _, want := range []string{"stack 0", "stack 7", "change7"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}

	for i := 0; i < 6; i++ {
		next, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
		model = next.(stacksModel)
	}
	if model.stackCursor != 6 {
		t.Fatalf("stackCursor = %d, want 6 after moving down", model.stackCursor)
	}
}

func TestStacksDisplayHidesMergedAndKeepsRemoteAtBottomByDefault(t *testing.T) {
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
		"2 branches without revisions. Run 'gx status list' to see a full list of features.",
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
			t.Fatalf("gx status list display missing %q in:\n%s", want, fullText)
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
	for _, want := range []string{"review (gxr)", "status (gxs)", "commit"} {
		if !strings.Contains(text, want) {
			t.Fatalf("root help missing alias %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "generate (gxg)") {
		t.Fatalf("root help should hide generate:\n%s", text)
	}
	if strings.Contains(text, "Shortcuts:") {
		t.Fatalf("root help should render aliases inline instead of a shortcut section:\n%s", text)
	}
}

func TestReviewAliasResolves(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"gxr"})
	if err != nil || cmd == nil || cmd.Name() != "review" {
		t.Fatalf("Find(gxr) = cmd=%v err=%v, want review command", cmd, err)
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
	t.Chdir(t.TempDir())
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
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_JUDGE", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_REVIEW_RESOURCES", "0")
	t.Setenv("GX_REVIEW_INDEXED_CONTEXT", "0")

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
		"- No material issues found in this change.",
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

func writeTestGXConfig(t *testing.T, gxHome, name, email string) {
	t.Helper()
	if err := os.MkdirAll(gxHome, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := fmt.Sprintf(`{"user":{"name":%q,"email":%q}}`, name, email)
	if err := os.WriteFile(filepath.Join(gxHome, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func TestInitYesAcceptsDefaultsAndSuppressesOutput(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	t.Chdir(root)
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "-y"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("gx init -y error = %v\n%s", err, out.String())
	}
	if out.String() != "" {
		t.Fatalf("gx init -y output = %q, want empty", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		t.Fatalf("after gx init -y .git missing: %v", err)
	}
}

func TestReviewCommandAcceptsScopeFlag(t *testing.T) {
	root := initGitRepo(t)
	t.Chdir(root)
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_REVIEW_RESOURCES", "0")
	t.Setenv("GX_REVIEW_INDEXED_CONTEXT", "0")

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

func TestPostReviewSummaryCommentUpsertsGitHubPRComment(t *testing.T) {
	remote := "https://github.com/acme/gx.git"
	branch := "feature/demo"
	var gotCommentBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/pulls":
			if r.URL.Query().Get("head") != "acme:"+branch {
				t.Fatalf("head query = %q", r.URL.Query().Get("head"))
			}
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/gx/pull/7"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode comment body: %v", err)
			}
			gotCommentBody = payload["body"]
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/gx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_REVIEW_RESOURCES", "0")
	t.Setenv("GX_REVIEW_INDEXED_CONTEXT", "0")

	report := codereview.Report{
		Findings: []codereview.Finding{{
			ID:               "docs.missing-readme",
			Title:            "Missing README",
			Summary:          "The repo lacks the standard entrypoint document new maintainers expect first.",
			Benefit:          "Improves onboarding speed by giving humans and agents one place to find setup, purpose, and common commands.",
			Recommendation:   "Add a concise README.",
			Strength:         "Strong",
			SourcePublishers: []string{"Go project"},
		}},
		Sources: []codereview.Source{{
			ID:        "go-code-review-comments",
			Title:     "Go Code Review Comments",
			URL:       "https://go.dev/wiki/CodeReviewComments",
			Publisher: "Go project",
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)

	if gotCommentBody == "" {
		t.Fatal("postReviewSummaryComment() did not send a comment")
	}
	for _, want := range []string{"<!-- gx review summary -->", "## Recommendations", "**Informed by:** Go project"} {
		if !strings.Contains(gotCommentBody, want) {
			t.Fatalf("comment body missing %q:\n%s", want, gotCommentBody)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("postReviewSummaryComment() wrote warnings:\n%s", stderr.String())
	}
}

func TestPostReviewSummaryCommentAttemptsInlineCommentForValidAnchor(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test User")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeTestFile(t, root, "main.go", "package main\nfunc main() {}\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "initial")
	writeTestFile(t, root, "main.go", "package main\nfunc main() { println(\"hi\") }\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "change")

	remote := "https://github.com/acme/gx.git"
	branch := "feature/demo"
	var inlineAttempted bool
	var summaryAttempted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/pulls":
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/gx/pull/7"}]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/gx/pulls/7/comments":
			inlineAttempted = true
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode inline comment body: %v", err)
			}
			if payload["path"] != "main.go" || int(payload["line"].(float64)) != 2 {
				t.Fatalf("inline payload = %#v", payload)
			}
			_, _ = w.Write([]byte(`{"id":99}`))
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			summaryAttempted = true
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/gx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")

	report := codereview.Report{
		RepoRoot: root,
		Findings: []codereview.Finding{{
			ID:             "ai.review.1",
			Title:          "Inline finding",
			Summary:        "main.go changed.",
			Benefit:        "Keeps comments anchored.",
			Recommendation: "Fix main.go.",
			Anchors:        []codereview.FindingAnchor{{File: "main.go", Line: 2}},
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RootPath:   root,
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)
	if !inlineAttempted {
		t.Fatal("postReviewSummaryComment() did not attempt inline comment")
	}
	if !summaryAttempted {
		t.Fatal("postReviewSummaryComment() did not post summary fallback")
	}
	if stderr.Len() != 0 {
		t.Fatalf("postReviewSummaryComment() wrote warnings:\n%s", stderr.String())
	}
}

func TestPostReviewSummaryCommentFallsBackWhenInlineCommentFails(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test User")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeTestFile(t, root, "main.go", "package main\nfunc main() {}\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "initial")
	writeTestFile(t, root, "main.go", "package main\nfunc main() { println(\"hi\") }\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "change")

	remote := "https://github.com/acme/gx.git"
	branch := "feature/demo"
	var summaryAttempted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/pulls":
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/gx/pull/7"}]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/gx/pulls/7/comments":
			http.Error(w, "line cannot be commented", http.StatusUnprocessableEntity)
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/gx/issues/7/comments":
			summaryAttempted = true
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/gx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")

	report := codereview.Report{
		RepoRoot: root,
		Findings: []codereview.Finding{{
			ID:             "ai.review.1",
			Title:          "Inline finding",
			Summary:        "main.go changed.",
			Benefit:        "Keeps comments anchored.",
			Recommendation: "Fix main.go.",
			Anchors:        []codereview.FindingAnchor{{File: "main.go", Line: 2}},
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RootPath:   root,
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)
	if !summaryAttempted {
		t.Fatal("postReviewSummaryComment() did not post summary after inline failure")
	}
	if !strings.Contains(stderr.String(), "Could not post GX inline review comment") {
		t.Fatalf("postReviewSummaryComment() warning = %q, want inline failure warning", stderr.String())
	}
}

func TestAutoReportFailurePostsCommandError(t *testing.T) {
	root := initGitRepo(t)
	t.Chdir(root)
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	var gotReport cloud.ReportLogRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/reported-logs" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReport); err != nil {
			t.Fatalf("decode report body: %v", err)
		}
		_, _ = w.Write([]byte(`{"id":"report-1","url":"https://gx.run/reports/report-1"}`))
	}))
	defer server.Close()
	t.Setenv("GX_CLOUD_URL", server.URL)

	autoReportFailure(context.Background(), authoring.NewEngine(), fmt.Errorf("push exploded"), "gx push")
	if !strings.Contains(gotReport.Error, "gx push: push exploded") {
		t.Fatalf("report error = %q, want command error", gotReport.Error)
	}
	if gotReport.GXVersion == "" || gotReport.OS == "" || gotReport.Arch == "" {
		t.Fatalf("report metadata incomplete: %#v", gotReport)
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
	cmd := exec.Command("git", "init", "-b", "main")
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
	runGitTest(t, root, args...)
}

func runGitTest(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s error = %v\n%s", strings.Join(args, " "), err, out)
	}
}
