//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

const e2eCloudToken = "gx_e2eabcdefghijklmnopqrstuvwxyz"

type harness struct {
	t      *testing.T
	bin    string
	repo   string
	home   string
	gxHome string
	remote string
}

type stackRow struct {
	BookmarkName string
	Status       string
	HeadChangeID sql.NullString
	RemoteRef    sql.NullString
}

type demuxProposalPayload struct {
	ID               string `json:"id"`
	RepoRoot         string `json:"repo_root"`
	ProposedChangeID string `json:"proposed_change_id"`
	ProposedCommitID string `json:"proposed_commit_id"`
	Hunks            []struct {
		ID    string `json:"id"`
		File  string `json:"file"`
		Patch string `json:"patch"`
	} `json:"hunks"`
	Revisions []struct {
		ID               string   `json:"id"`
		Intent           string   `json:"intent"`
		Files            []string `json:"files"`
		TargetStack      string   `json:"target_stack"`
		UseHunks         bool     `json:"use_hunks"`
		HunkIDs          []string `json:"hunk_ids"`
		ProvenanceStatus string   `json:"provenance_status"`
		SessionIDs       []string `json:"session_ids"`
	} `json:"revisions"`
	FeasibilityWarnings []struct {
		RevisionID string `json:"revision_id"`
		Severity   string `json:"severity"`
		Source     string `json:"source"`
		Message    string `json:"message"`
	} `json:"feasibility_warnings"`
	Warnings []string `json:"warnings"`
}

func (p *demuxProposalPayload) UnmarshalJSON(data []byte) error {
	type proposalAlias demuxProposalPayload
	var direct proposalAlias
	if err := json.Unmarshal(data, &direct); err != nil {
		return err
	}
	if direct.ID != "" {
		*p = demuxProposalPayload(direct)
		return nil
	}
	var wrapped struct {
		Proposal proposalAlias `json:"proposal"`
	}
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	*p = demuxProposalPayload(wrapped.Proposal)
	return nil
}

type demuxApplyPayload struct {
	Proposal         demuxProposalPayload `json:"proposal"`
	RemainingChanges bool                 `json:"remaining_changes"`
	AcceptedSubset   bool                 `json:"accepted_subset"`
	NextAction       string               `json:"next_action"`
	Revisions        []struct {
		Change struct {
			Description string `json:"description"`
		} `json:"change"`
		Stack struct {
			BookmarkName string `json:"BookmarkName"`
		} `json:"Stack"`
	} `json:"revisions"`
}

type demuxEvidenceRow struct {
	RevisionProposalID string
	Intent             string
	HunkIDsJSON        string
	ProvenanceStatus   string
}

func TestGXAddMaintainsBookmarkStateAndAttachesExplicitSession(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)

	if _, err := os.Stat(filepath.Join(h.repo, ".jj")); !os.IsNotExist(err) {
		t.Fatalf("before gx init .jj exists or stat failed: %v", err)
	}
	assertCurrentBranch(t, h, "main")
	assertNoGXBookmarks(t, h)

	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	if _, err := os.Stat(filepath.Join(h.repo, ".jj")); err != nil {
		t.Fatalf("after gx init .jj missing: %v", err)
	}
	assertCurrentBranch(t, h, "main")
	assertNoGXBookmarks(t, h)

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.insertSession("session-alpha")
	h.insertSessionRequest("session-alpha", "implement alpha")

	assertCurrentBranch(t, h, "main")
	assertNoGXBookmarks(t, h)
	assertChangeSessions(t, h, nil)

	h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "add", "-m", "feat alpha")

	bookmarks := h.bookmarkTargets()
	alphaBookmark := "feature/feat-alpha"
	alphaTarget := bookmarks[alphaBookmark]
	if alphaTarget == "" {
		t.Fatalf("after gx add missing %s bookmark: %#v", alphaBookmark, bookmarks)
	}
	assertCurrentBranch(t, h, alphaBookmark)
	assertStack(t, h, alphaBookmark, "draft", alphaTarget, "")
	assertChangeSessions(t, h, map[string][]string{"feat alpha": {"session-alpha"}})
}

func TestGXBaseShowsAndSetsAuthoringBase(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	base := h.gx("base")
	for _, want := range []string{"GX base", "Base main", "Current main", "On base yes"} {
		if !strings.Contains(base, want) {
			t.Fatalf("gx base missing %q in:\n%s", want, base)
		}
	}

	h.run("git", "checkout", "-b", "develop")
	updated := h.gx("base", "--set", "develop")
	for _, want := range []string{"Base updated", "Base develop", "Current develop", "On base yes"} {
		if !strings.Contains(updated, want) {
			t.Fatalf("gx base --set missing %q in:\n%s", want, updated)
		}
	}
	assertCurrentBranch(t, h, "develop")

	h.run("git", "switch", "main")
	status := h.gx("base")
	for _, want := range []string{"Base develop", "Current main", "On base yes"} {
		if !strings.Contains(status, want) {
			t.Fatalf("gx base after switch missing %q in:\n%s", want, status)
		}
	}

	reset := h.gx("base", "--set", "develop")
	for _, want := range []string{"Base updated", "Base develop", "Current develop", "On base yes"} {
		if !strings.Contains(reset, want) {
			t.Fatalf("gx base --set from clean checkout missing %q in:\n%s", want, reset)
		}
	}

	internal := runCommandAllowError(h.t, h.repo, h.env(), h.bin, "base", "--set", "gx/develop")
	if internal.exitCode == 0 || !strings.Contains(internal.output, "Revision `gx/develop` doesn't exist") {
		t.Fatalf("gx base --set gx/develop = exit %d:\n%s", internal.exitCode, internal.output)
	}
}

func TestGXBaseSetMainUsesGXCheckout(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")
	assertCurrentBranch(t, h, "main")

	updated := h.gx("base", "--set", "main")
	for _, want := range []string{"Base updated", "Base main", "Current main", "On base yes"} {
		if !strings.Contains(updated, want) {
			t.Fatalf("gx base --set main missing %q in:\n%s", want, updated)
		}
	}
	assertCurrentBranch(t, h, "main")
}

func TestGXGenerateCreatesFileLevelRevisionsAndEvidence(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("session-alpha")
	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "messy work")

	rawApply := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "generate", "--json", "--intent", "demux e2e")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	proposal := applied.Proposal
	if proposal.ID == "" || proposal.ProposedChangeID == "" || proposal.ProposedCommitID == "" {
		t.Fatalf("proposal missing identity fields: %#v", proposal)
	}
	if len(proposal.Revisions) != 2 {
		t.Fatalf("proposal revisions = %#v, want 2", proposal.Revisions)
	}
	if proposal.Revisions[0].Intent != "demux e2e: alpha.txt" || !reflect.DeepEqual(proposal.Revisions[0].Files, []string{"alpha.txt"}) {
		t.Fatalf("first proposed revision = %#v", proposal.Revisions[0])
	}
	if proposal.Revisions[1].Intent != "demux e2e: beta.txt" || !reflect.DeepEqual(proposal.Revisions[1].Files, []string{"beta.txt"}) {
		t.Fatalf("second proposed revision = %#v", proposal.Revisions[1])
	}
	if proposal.Revisions[0].ProvenanceStatus != "explicit" || !reflect.DeepEqual(proposal.Revisions[0].SessionIDs, []string{"session-alpha"}) {
		t.Fatalf("proposal provenance = %#v", proposal.Revisions[0])
	}
	if len(proposal.FeasibilityWarnings) < 2 {
		t.Fatalf("proposal feasibility warnings = %#v, want at least unmapped hunk warnings", proposal.FeasibilityWarnings)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	assertCurrentBranch(t, h, "main")
	assertChangeSessions(t, h, map[string][]string{
		"demux e2e: alpha.txt": {"session-alpha"},
		"demux e2e: beta.txt":  {"session-alpha"},
	})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "u1", Intent: "demux e2e: alpha.txt", HunkIDsJSON: `["h1"]`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "u2", Intent: "demux e2e: beta.txt", HunkIDsJSON: `["h2"]`, ProvenanceStatus: "explicit"},
	})
}

func TestGXGenerateCreatesEveryProposedStackAndPushPublishesAll(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")
	h.writeMultiStackGenerateFiles()

	rawApply := h.gx("generate", "--json", "--intent", "multi generate")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	if applied.Proposal.ID == "" {
		t.Fatalf("generate result missing proposal: %#v", applied)
	}
	wantStacks := map[string]bool{
		"feature/demux-routing":    false,
		"feature/stack-management": false,
		"test/e2e-tests":           false,
	}
	for _, revision := range applied.Proposal.Revisions {
		if _, ok := wantStacks[revision.TargetStack]; ok {
			wantStacks[revision.TargetStack] = true
		}
	}
	for bookmark, seen := range wantStacks {
		if !seen {
			t.Fatalf("generate proposal missing stack %q in revisions %#v\n%s", bookmark, applied.Proposal.Revisions, rawApply)
		}
	}
	if applied.RemainingChanges || applied.AcceptedSubset || applied.NextAction != "" {
		t.Fatalf("generate apply partial fields = remaining:%t subset:%t next:%q\n%s", applied.RemainingChanges, applied.AcceptedSubset, applied.NextAction, rawApply)
	}
	assertCurrentBranch(t, h, "main")
	if len(applied.Revisions) != 3 {
		t.Fatalf("applied revisions = %#v, want 3", applied.Revisions)
	}
	for _, revision := range applied.Revisions {
		if _, ok := wantStacks[revision.Stack.BookmarkName]; !ok {
			t.Fatalf("applied unexpected stack %q in %#v", revision.Stack.BookmarkName, applied.Revisions)
		}
	}
	for bookmark := range wantStacks {
		assertStack(t, h, bookmark, "draft", h.bookmarkTargets()[bookmark], "")
		assertStackInStacksJSON(t, h, bookmark, 1)
		assertRemoteBranchMissing(t, h, bookmark)
	}

	output := h.gx("push")
	if !strings.Contains(output, "Pushed") || !strings.Contains(output, "3 stacks") {
		t.Fatalf("push output = %q, want pushed 3 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
	for bookmark := range wantStacks {
		assertRemoteBranchExists(t, h, bookmark)
		assertStack(t, h, bookmark, "published", h.bookmarkTargets()[bookmark], "refs/heads/"+bookmark)
	}
}

func TestGXGenerateAppendsToExistingJJBookmark(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	existingStack := "feature/stack-management"
	h.writeTrackedFile("existing-stack.txt", "existing\n")
	h.run("jj", "describe", "-m", "existing stack-management work")
	h.run("jj", "bookmark", "set", existingStack, "-r", "@")
	existingTarget := h.bookmarkTargets()[existingStack]
	if existingTarget == "" {
		t.Fatalf("missing raw existing bookmark %q in %#v", existingStack, h.bookmarkTargets())
	}
	h.run("jj", "new", "main")
	assertCurrentBranch(t, h, "main")

	h.writeMultiStackGenerateFiles()
	rawApply := h.gx("generate", "--json", "--intent", "append existing generate", "internal/storage")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 1 || applied.Revisions[0].Stack.BookmarkName != existingStack {
		t.Fatalf("applied revisions = %#v, want one revision appended to %s", applied.Revisions, existingStack)
	}
	if got := h.bookmarkTargets()[existingStack]; got == "" || got == existingTarget {
		t.Fatalf("existing bookmark target after append = %q, before %q", got, existingTarget)
	}
	assertCurrentBranch(t, h, "main")
	assertStackInStacksJSON(t, h, existingStack, 2)
}

func TestGXSyncPullsCloudBookmarkForThrowawayMarkdownChange(t *testing.T) {
	source := newHarness(t)
	source.initGitRepo(true)
	source.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	const fileName = "throwaway-sync.md"
	const fileContents = "sync pulled this markdown change\n"
	source.writeTrackedFile(fileName, fileContents)
	source.run("jj", "describe", "-m", "throwaway sync markdown")

	rawApply := source.gx("generate", "--json", "--intent", "update throwaway sync markdown")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	if len(applied.Proposal.Revisions) != 1 || len(applied.Revisions) != 1 {
		t.Fatalf("generate result revisions = proposal:%#v applied:%#v\n%s", applied.Proposal.Revisions, applied.Revisions, rawApply)
	}
	branch := applied.Proposal.Revisions[0].TargetStack
	if branch == "" || strings.HasPrefix(branch, "gx/") {
		t.Fatalf("generate target_stack = %q, want conventional public branch", branch)
	}

	assertCurrentBranch(t, source, "main")
	assertStackInStacksJSON(t, source, branch, 1)
	source.gx("push", branch)
	assertRemoteBranchExists(t, source, branch)
	remoteHead := strings.TrimSpace(source.run("git", "--git-dir", source.remote, "rev-parse", "refs/heads/"+branch))
	if remoteHead == "" {
		t.Fatalf("remote branch %s has empty head", branch)
	}

	target := newHarness(t)
	const githubRemoteURL = "git@github.com:e2e/gx-sync.git"
	target.cloneFromRemoteAsGitHub(source.remote, githubRemoteURL)
	target.gx("init", "--name", "Joe Example", "--email", "joe@example.com")
	target.saveCloudCredentials()
	assertCurrentBranch(t, target, "main")
	if _, ok := target.bookmarkTargets()[branch]; ok {
		t.Fatalf("target unexpectedly had bookmark %s before gx sync: %#v", branch, target.bookmarkTargets())
	}

	apiURL := startBookmarkCloudAPI(t, branch, remoteHead)
	output := target.gxWithEnv([]string{"GX_CLOUD_URL=" + apiURL + "/gx/pr"}, "sync")
	if !strings.Contains(output, "Remote catch-up 1 bookmark(s) fetched from remote") {
		t.Fatalf("gx sync output missing cloud catch-up:\n%s", output)
	}
	assertCurrentBranch(t, target, "main")
	if target.bookmarkTargets()[branch] == "" {
		t.Fatalf("target missing synced bookmark %s after gx sync: %#v", branch, target.bookmarkTargets())
	}
	got := target.run("jj", "file", "show", "-r", branch, fileName)
	if got != fileContents {
		t.Fatalf("synced file contents = %q, want %q", got, fileContents)
	}
}

func TestGXGenerateAttachesRepoLocalSessionWithoutEnv(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("cursor-session")
	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "messy cursor work")

	rawApply := h.gx("generate", "--json", "--intent", "cursor generate")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	proposal := applied.Proposal
	if len(proposal.Revisions) != 1 {
		t.Fatalf("proposal revisions = %#v, want 1", proposal.Revisions)
	}
	if proposal.Revisions[0].ProvenanceStatus != "repo_local" || !reflect.DeepEqual(proposal.Revisions[0].SessionIDs, []string{"cursor-session"}) {
		t.Fatalf("proposal provenance = %#v, want repo-local cursor session", proposal.Revisions[0])
	}
	if len(applied.Revisions) != 1 {
		t.Fatalf("applied revisions = %#v, want 1", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{
		"cursor generate: alpha.txt": {"cursor-session"},
	})
}

func TestGXGeneratePreservesSymbolLevelHunks(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("app.go", symbolFixture("one", "two"))
	h.run("jj", "describe", "-m", "seed symbols")
	h.gx("add", "-m", "seed symbols")
	h.gx("base", "--set", "feature/seed-symbols")

	h.insertSession("session-alpha")
	h.writeTrackedFile("app.go", symbolFixture("ONE", "TWO"))
	h.run("jj", "describe", "-m", "messy symbols")

	rawApply := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "generate", "--json", "--intent", "split symbols")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode generate apply: %v\n%s", err, rawApply)
	}
	proposal := applied.Proposal
	if len(proposal.Revisions) != 1 {
		t.Fatalf("proposal revisions = %#v, want one hunk-level symbol revision", proposal.Revisions)
	}
	if !proposal.Revisions[0].UseHunks || len(proposal.Revisions[0].HunkIDs) != 2 {
		t.Fatalf("symbol revision = %#v, want both symbol hunks preserved", proposal.Revisions[0])
	}
	if len(applied.Revisions) != 1 {
		t.Fatalf("applied revisions = %#v, want 1", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{proposal.Revisions[0].Intent: {"session-alpha"}})
}

func TestGXPushAllPublishesEveryStack(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.insertSession("session-alpha")
	h.insertSessionRequest("session-alpha", "implement alpha")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "add", "-m", "feat alpha")
	alphaBookmark := "feature/feat-alpha"
	alphaTarget := h.bookmarkTargets()[alphaBookmark]

	h.run("git", "switch", "main")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "beta")
	h.insertSession("session-beta")
	h.insertSessionRequest("session-beta", "implement beta")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-beta"}, "add", "-m", "feat beta")
	betaBookmark := "feature/feat-beta"
	betaTarget := h.bookmarkTargets()[betaBookmark]

	assertRemoteBranchMissing(t, h, alphaBookmark)
	assertRemoteBranchMissing(t, h, betaBookmark)

	output := h.gx("push")
	if !strings.Contains(output, "Pushed") || !strings.Contains(output, "2 stacks") {
		t.Fatalf("push output = %q, want pushed 2 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
	assertRemoteBranchExists(t, h, alphaBookmark)
	assertRemoteBranchExists(t, h, betaBookmark)
	assertStack(t, h, alphaBookmark, "published", alphaTarget, "refs/heads/"+alphaBookmark)
	assertStack(t, h, betaBookmark, "published", betaTarget, "refs/heads/"+betaBookmark)
	assertChangeBookmark(t, h, "feat alpha", alphaBookmark)
	assertChangeBookmark(t, h, "feat beta", betaBookmark)
}

func TestGXPushNoStacksIsNoopOnMain(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	output := h.gx("push")

	if !strings.Contains(output, "Pushed") || !strings.Contains(output, "0 stacks") {
		t.Fatalf("push output = %q, want pushed 0 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
}

func TestGXStatusHidesStackMergedIntoBase(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.gx("add", "-m", "feat alpha")
	alphaBookmark := "feature/feat-alpha"
	assertStackInStacksJSON(t, h, alphaBookmark, 1)

	h.run("jj", "bookmark", "set", "main", "-r", alphaBookmark, "--allow-backwards")
	h.run("git", "switch", "main")

	assertCurrentBranch(t, h, "main")
	assertStackNotInStacksJSON(t, h, alphaBookmark)
	output := h.gx("push")
	if !strings.Contains(output, "Pushed") || !strings.Contains(output, "0 stacks") {
		t.Fatalf("push output = %q, want pushed 0 stacks", output)
	}
}

func TestGXBaseSetRefusesDirtyEditCheckout(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.gx("add", "-m", "feat alpha")
	alphaBookmark := strings.TrimSpace(h.run("git", "branch", "--show-current"))
	assertCurrentBranch(t, h, alphaBookmark)
	status := h.gx("status")
	if strings.Contains(status, "Active working change is not yet recorded") {
		t.Fatalf("status reported empty JJ working change as unrecorded:\n%s", status)
	}

	h.writeTrackedFile("scratch.txt", "scratch\n")
	denied := runCommandAllowError(h.t, h.repo, h.env(), h.bin, "base", "--set", "main")
	if denied.exitCode == 0 {
		t.Fatalf("gx base --set main with dirty edit checkout succeeded:\n%s", denied.output)
	}
	assertCurrentBranch(t, h, alphaBookmark)
	if _, err := os.Stat(filepath.Join(h.repo, "scratch.txt")); err != nil {
		t.Fatalf("scratch file missing after refused base set: %v", err)
	}
}

func TestGXStatusShowsImplicitStackAliases(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.insertSession("session-alpha")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "add", "-m", "feat alpha")

	h.run("git", "switch", "main")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "beta")
	h.insertSession("session-beta")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-beta"}, "add", "-m", "feat beta")

	status := h.gx("status")
	if !strings.Contains(status, "feat beta  s1") {
		t.Fatalf("status missing beta alias:\n%s", status)
	}
	if !strings.Contains(status, "feat alpha  s2") {
		t.Fatalf("status missing alpha alias:\n%s", status)
	}
}

func symbolFixture(first, second string) string {
	return strings.Join([]string{
		"package main",
		"",
		"func First() string {",
		"\treturn " + strconv.Quote(first),
		"}",
		"",
		"// gap 01",
		"// gap 02",
		"// gap 03",
		"// gap 04",
		"// gap 05",
		"// gap 06",
		"// gap 07",
		"// gap 08",
		"",
		"func Second() string {",
		"\treturn " + strconv.Quote(second),
		"}",
		"",
	}, "\n")
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	requireTool(t, "git")
	requireTool(t, "jj")
	bin := filepath.Join(t.TempDir(), "gx")
	runHost(t, e2eSourceRoot(t), "go", "build", "-o", bin, "./cmd/gx")

	h := &harness{
		t:      t,
		bin:    bin,
		repo:   filepath.Join(t.TempDir(), "repo"),
		home:   filepath.Join(t.TempDir(), "home"),
		gxHome: filepath.Join(t.TempDir(), "gx-home"),
	}
	if err := os.MkdirAll(h.repo, 0o755); err != nil {
		t.Fatalf("create repo dir: %v", err)
	}
	if resolved, err := filepath.EvalSymlinks(h.repo); err == nil {
		h.repo = resolved
	}
	t.Setenv("GX_HOME", h.gxHome)
	return h
}

func e2eSourceRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve e2e source root: runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("resolve e2e source root: %v", err)
	}
	return root
}

func (h *harness) saveCloudCredentials() {
	h.t.Helper()
	credentialsPath := filepath.Join(h.gxHome, "credentials.json")
	if err := os.MkdirAll(filepath.Dir(credentialsPath), 0o755); err != nil {
		h.t.Fatalf("create gx home: %v", err)
	}
	credentials := map[string]any{
		"cloud": map[string]any{
			"token":                  e2eCloudToken,
			"cli_session_token":      e2eCloudToken,
			"cli_session_expires_at": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano),
			"user_id":                "e2e-user",
			"login":                  "e2e",
			"session_id":             "e2e-session",
			"machine_id":             "e2e-machine",
			"machine_name":           "e2e-machine",
			"obtained_at":            time.Now().UTC().Format(time.RFC3339Nano),
		},
	}
	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		h.t.Fatalf("marshal cloud credentials: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(credentialsPath, data, 0o600); err != nil {
		h.t.Fatalf("write cloud credentials: %v", err)
	}
}

func (h *harness) initGitRepo(withRemote bool) {
	h.t.Helper()
	h.run("git", "init", "-b", "main")
	h.writeFile("README.md", "# e2e\n")
	h.run("git", "add", "README.md")
	h.run("git", "-c", "user.name=Joe Example", "-c", "user.email=joe@example.com", "commit", "-m", "initial")
	if !withRemote {
		return
	}
	h.remote = filepath.Join(h.t.TempDir(), "remote.git")
	runHost(h.t, h.t.TempDir(), "git", "init", "--bare", h.remote)
	h.run("git", "remote", "add", "origin", h.remote)
	h.run("git", "push", "origin", "main")
}

func (h *harness) cloneFromRemoteAsGitHub(remotePath, remoteURL string) {
	h.t.Helper()
	if err := os.RemoveAll(h.repo); err != nil {
		h.t.Fatalf("remove empty repo before clone: %v", err)
	}
	runHost(h.t, h.t.TempDir(), "git", "clone", remotePath, h.repo)
	h.remote = remotePath
	h.run("git", "remote", "set-url", "origin", remoteURL)
	h.run("git", "config", "url.file://"+remotePath+".insteadOf", remoteURL)
}

func (h *harness) writeFile(name, contents string) {
	h.t.Helper()
	path := filepath.Join(h.repo, name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		h.t.Fatalf("write %s: %v", name, err)
	}
}

func (h *harness) writeTrackedFile(name, contents string) {
	h.t.Helper()
	h.writeFile(name, contents)
	h.run("jj", "file", "track", name)
}

func (h *harness) insertSession(id string) {
	h.t.Helper()
	db, err := storage.Open(context.Background())
	if err != nil {
		h.t.Fatalf("open gx db: %v", err)
	}
	defer db.Close()
	now := time.Now().UnixMilli()
	if _, err := db.ExecContext(context.Background(), `
		INSERT INTO sessions (id, created_at, command, cwd, gx_version, repo_root)
		VALUES (?, ?, ?, ?, ?, ?)
	`, id, now, "cursor agent", h.repo, "e2e", h.repo); err != nil {
		h.t.Fatalf("insert session: %v", err)
	}
}

func (h *harness) insertSessionRequest(sessionID, text string) {
	h.t.Helper()
	db, err := storage.Open(context.Background())
	if err != nil {
		h.t.Fatalf("open gx db: %v", err)
	}
	defer db.Close()
	now := time.Now().UnixMilli()
	if _, err := db.ExecContext(context.Background(), `
		INSERT INTO requests (id, session_id, created_at, provider, endpoint, method, model, request_body, request_headers)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, "request-"+sessionID, sessionID, now, "openai", "/v1/responses", "POST", "gpt-e2e", []byte(`{"input":`+strconv.Quote(text)+`}`), "{}"); err != nil {
		h.t.Fatalf("insert session request: %v", err)
	}
}

func (h *harness) gx(args ...string) string {
	h.t.Helper()
	return h.gxWithEnv(nil, args...)
}

func (h *harness) gxWithEnv(extra []string, args ...string) string {
	h.t.Helper()
	return runCommand(h.t, h.repo, h.env(extra...), h.bin, args...)
}

func (h *harness) writeMultiStackGenerateFiles() {
	h.t.Helper()
	for _, dir := range []string{
		filepath.Join(h.repo, "internal", "authoring"),
		filepath.Join(h.repo, "internal", "storage"),
		filepath.Join(h.repo, "test", "e2e"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			h.t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	h.writeTrackedFile("internal/authoring/demux.go", "package authoring\n")
	h.writeTrackedFile("internal/storage/schema.sql", "CREATE TABLE compose_route(id INTEGER);\n")
	h.writeTrackedFile("test/e2e/gx_e2e_test.go", "package e2e\n")
	h.run("jj", "describe", "-m", "multi generate source")
}

func (h *harness) run(name string, args ...string) string {
	h.t.Helper()
	return runCommand(h.t, h.repo, h.env(), name, args...)
}

func (h *harness) env(extra ...string) []string {
	blocked := []string{
		"GX_HOME=",
		"GX_SESSION_ID=",
		"GX_SESSION_IDS=",
		"GX_CLOUD_URL=",
		"GX_SEMANTIC_INDEX=",
		"GX_TPUF_NAMESPACE=",
		"GX_TPUF_BASE_URL=",
		"GX_OPENAI_BASE_URL=",
		"GX_OPENAI_EMBEDDING_MODEL=",
		"GX_EMBEDDING_DIMENSIONS=",
		"GX_SEMANTIC_BATCH_SIZE=",
		"GX_SEMANTIC_MAX_CHUNK_BYTES=",
		"TURBOPUFFER_API_KEY=",
		"OPENAI_API_KEY=",
		"GX_POSTLIST_URL=",
		"GX_POSTLIST_API_KEY=",
		"HOME=",
		"NO_COLOR=",
		"TERM=",
		"GIT_CONFIG_GLOBAL=",
		"GIT_CONFIG_NOSYSTEM=",
	}
	var env []string
	for _, value := range os.Environ() {
		skip := false
		for _, prefix := range blocked {
			if strings.HasPrefix(value, prefix) {
				skip = true
				break
			}
		}
		if !skip {
			env = append(env, value)
		}
	}
	env = append(env,
		"GX_HOME="+h.gxHome,
		"HOME="+h.home,
		"NO_COLOR=1",
		"TERM=dumb",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_NOSYSTEM=1",
	)
	env = append(env, extra...)
	return env
}

func (h *harness) bookmarkTargets() map[string]string {
	h.t.Helper()
	result := runCommandAllowError(h.t, h.repo, h.env(), "jj", "bookmark", "list", "-T", `name ++ "|" ++ normal_target.change_id() ++ "\n"`)
	if result.exitCode != 0 && strings.Contains(result.output, "There is no jj repo") {
		return map[string]string{}
	}
	if result.exitCode != 0 {
		h.t.Fatalf("jj bookmark list failed with exit %d:\n%s", result.exitCode, result.output)
	}
	out := result.output
	bookmarks := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		bookmarks[parts[0]] = parts[1]
	}
	return bookmarks
}

func startBookmarkCloudAPI(t *testing.T, branchName, remoteHead string) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/bookmarks" {
			t.Errorf("cloud request = %s %s, want GET /bookmarks", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+e2eCloudToken {
			t.Errorf("cloud auth = %q, want bearer token", got)
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		payload := []map[string]any{{
			"id":              "bookmark-sync-e2e",
			"repo_full_name":  "e2e/gx-sync",
			"branch_name":     branchName,
			"title":           "throwaway sync markdown",
			"revision":        1,
			"head_commit_id":  remoteHead,
			"remote_head_sha": remoteHead,
			"merge_status":    "open",
			"updated_at_ms":   time.Now().UnixMilli(),
		}}
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("encode bookmark payload: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func assertNoGXBookmarks(t *testing.T, h *harness) {
	t.Helper()
	for name := range h.bookmarkTargets() {
		if strings.HasPrefix(name, "gx/") {
			t.Fatalf("unexpected gx bookmark %q in %#v", name, h.bookmarkTargets())
		}
	}
}

func assertCurrentBranch(t *testing.T, h *harness, want string) {
	t.Helper()
	got := strings.TrimSpace(h.run("git", "branch", "--show-current"))
	if got != want {
		t.Fatalf("current branch = %q, want %q", got, want)
	}
}

func assertStack(t *testing.T, h *harness, bookmark, status, headChangeID, remoteRef string) {
	t.Helper()
	stacks := readStacks(t, h)
	stack, ok := stacks[bookmark]
	if !ok {
		t.Fatalf("stack %q missing from %#v", bookmark, stacks)
	}
	if stack.Status != status {
		t.Fatalf("stack %q status = %q, want %q", bookmark, stack.Status, status)
	}
	if headChangeID != "" && (!stack.HeadChangeID.Valid || stack.HeadChangeID.String != headChangeID) {
		t.Fatalf("stack %q head_change_id = %v, want %q", bookmark, stack.HeadChangeID, headChangeID)
	}
	if remoteRef == "" {
		if stack.RemoteRef.Valid && stack.RemoteRef.String != "" {
			t.Fatalf("stack %q remote_ref = %q, want empty", bookmark, stack.RemoteRef.String)
		}
		return
	}
	if !stack.RemoteRef.Valid || stack.RemoteRef.String != remoteRef {
		t.Fatalf("stack %q remote_ref = %v, want %q", bookmark, stack.RemoteRef, remoteRef)
	}
}

func assertStackInStacksJSON(t *testing.T, h *harness, bookmark string, revisions int) {
	t.Helper()
	rawStacks := h.gx("status", "--json")
	var summary struct {
		Stacks []struct {
			BookmarkName  string `json:"BookmarkName"`
			RevisionCount int    `json:"RevisionCount"`
		} `json:"Stacks"`
	}
	if err := json.Unmarshal([]byte(rawStacks), &summary); err != nil {
		t.Fatalf("decode stacks: %v\n%s", err, rawStacks)
	}
	for _, stack := range summary.Stacks {
		if stack.BookmarkName != bookmark {
			continue
		}
		if stack.RevisionCount != revisions {
			t.Fatalf("stack %q revision count = %d, want %d\n%s", bookmark, stack.RevisionCount, revisions, rawStacks)
		}
		return
	}
	t.Fatalf("stack %q missing from gx stacks:\n%s", bookmark, rawStacks)
}

func assertStackNotInStacksJSON(t *testing.T, h *harness, bookmark string) {
	t.Helper()
	rawStacks := h.gx("status", "--json")
	var summary struct {
		Stacks []struct {
			BookmarkName string `json:"BookmarkName"`
		} `json:"Stacks"`
	}
	if err := json.Unmarshal([]byte(rawStacks), &summary); err != nil {
		t.Fatalf("decode stacks: %v\n%s", err, rawStacks)
	}
	for _, stack := range summary.Stacks {
		if stack.BookmarkName == bookmark {
			t.Fatalf("stack %q should not be visible in gx stacks after merge:\n%s", bookmark, rawStacks)
		}
	}
}

func readStacks(t *testing.T, h *harness) map[string]stackRow {
	t.Helper()
	db := openDB(t, h)
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), `
		SELECT bookmark_name, status, head_change_id, remote_ref
		FROM stacks
		ORDER BY bookmark_name
	`)
	if err != nil {
		t.Fatalf("query stacks: %v", err)
	}
	defer rows.Close()
	stacks := map[string]stackRow{}
	for rows.Next() {
		var stack stackRow
		if err := rows.Scan(&stack.BookmarkName, &stack.Status, &stack.HeadChangeID, &stack.RemoteRef); err != nil {
			t.Fatalf("scan stack: %v", err)
		}
		stacks[stack.BookmarkName] = stack
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate stacks: %v", err)
	}
	return stacks
}

func assertChangeSessions(t *testing.T, h *harness, want map[string][]string) {
	t.Helper()
	got := map[string][]string{}
	db := openDB(t, h)
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), `
		SELECT c.description, s.id
		FROM change_sessions cs
		JOIN changes c ON c.id = cs.change_id
		JOIN sessions s ON s.id = cs.session_id
		ORDER BY c.description, s.id
	`)
	if err != nil {
		t.Fatalf("query change sessions: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var description, sessionID string
		if err := rows.Scan(&description, &sessionID); err != nil {
			t.Fatalf("scan change session: %v", err)
		}
		got[description] = append(got[description], sessionID)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate change sessions: %v", err)
	}
	if len(want) == 0 {
		want = nil
	}
	if len(got) == 0 {
		got = nil
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("change sessions = %#v, want %#v", got, want)
	}
}

func assertDemuxEvidence(t *testing.T, h *harness, proposalID string, want []demuxEvidenceRow) {
	t.Helper()
	db := openDB(t, h)
	defer db.Close()
	rows, err := db.QueryContext(context.Background(), `
		SELECT revision_proposal_id, intent, hunk_ids_json, provenance_status
		FROM change_demux_evidence
		WHERE demux_proposal_id = ?
		ORDER BY id
	`, proposalID)
	if err != nil {
		t.Fatalf("query demux evidence: %v", err)
	}
	defer rows.Close()
	var got []demuxEvidenceRow
	for rows.Next() {
		var row demuxEvidenceRow
		if err := rows.Scan(&row.RevisionProposalID, &row.Intent, &row.HunkIDsJSON, &row.ProvenanceStatus); err != nil {
			t.Fatalf("scan demux evidence: %v", err)
		}
		got = append(got, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate demux evidence: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("demux evidence = %#v, want %#v", got, want)
	}
}

func assertChangeBookmark(t *testing.T, h *harness, description, bookmark string) {
	t.Helper()
	count := changeBookmarkCount(t, h, description, bookmark)
	if count != 1 {
		t.Fatalf("change bookmark count for %q/%q = %d, want 1", description, bookmark, count)
	}
}

func assertNoChangeBookmark(t *testing.T, h *harness, description, bookmark string) {
	t.Helper()
	count := changeBookmarkCount(t, h, description, bookmark)
	if count != 0 {
		t.Fatalf("change bookmark count for %q/%q = %d, want 0", description, bookmark, count)
	}
}

func changeBookmarkCount(t *testing.T, h *harness, description, bookmark string) int {
	t.Helper()
	db := openDB(t, h)
	defer db.Close()
	var count int
	if err := db.QueryRowContext(context.Background(), `
		SELECT count(*)
		FROM change_bookmarks cb
		JOIN changes c ON c.id = cb.change_id
		WHERE c.description = ? AND cb.bookmark_name = ?
	`, description, bookmark).Scan(&count); err != nil {
		t.Fatalf("query change bookmark: %v", err)
	}
	return count
}

func assertRemoteBranchMissing(t *testing.T, h *harness, branch string) {
	t.Helper()
	out := runCommandAllowError(t, h.repo, h.env(), "git", "--git-dir", h.remote, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if out.exitCode == 0 {
		t.Fatalf("remote branch %q exists before push", branch)
	}
}

func assertRemoteBranchExists(t *testing.T, h *harness, branch string) {
	t.Helper()
	out := runCommandAllowError(t, h.repo, h.env(), "git", "--git-dir", h.remote, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if out.exitCode != 0 {
		t.Fatalf("remote branch %q missing after push:\n%s", branch, out.output)
	}
}

func openDB(t *testing.T, h *harness) *sql.DB {
	t.Helper()
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("open gx db: %v", err)
	}
	return db
}

func requireTool(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s is required for e2e tests: %v", name, err)
	}
}

func runHost(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	return runCommand(t, dir, os.Environ(), name, args...)
}

type commandResult struct {
	output   string
	exitCode int
}

func runCommand(t *testing.T, dir string, env []string, name string, args ...string) string {
	t.Helper()
	result := runCommandAllowError(t, dir, env, name, args...)
	if result.exitCode != 0 {
		t.Fatalf("%s %s failed with exit %d:\n%s", name, strings.Join(args, " "), result.exitCode, result.output)
	}
	return result.output
}

func runCommandAllowError(t *testing.T, dir string, env []string, name string, args ...string) commandResult {
	t.Helper()
	return runCommandAllowErrorWithInput(t, dir, env, "", name, args...)
}

func runCommandAllowErrorWithInput(t *testing.T, dir string, env []string, input string, name string, args ...string) commandResult {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = env
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	out, err := cmd.CombinedOutput()
	output := string(out)
	if ctx.Err() == context.DeadlineExceeded {
		return commandResult{output: output + "\ncommand timed out", exitCode: -1}
	}
	if err == nil {
		return commandResult{output: output, exitCode: 0}
	}
	if exit, ok := err.(*exec.ExitError); ok {
		return commandResult{output: output, exitCode: exit.ExitCode()}
	}
	return commandResult{output: output + "\n" + err.Error(), exitCode: -1}
}
