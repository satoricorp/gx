//go:build e2e

package e2e

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

type demuxFixPayload struct {
	Proposal struct {
		ID        string `json:"id"`
		Revisions []struct {
			ID          string   `json:"id"`
			TargetStack string   `json:"target_stack"`
			BaseStack   string   `json:"base_stack"`
			DependsOn   []string `json:"depends_on"`
		} `json:"revisions"`
	} `json:"proposal"`
	Review struct {
		Valid       bool `json:"valid"`
		RepairHints []struct {
			Kind string `json:"kind"`
		} `json:"repair_hints"`
	} `json:"review"`
	State   string `json:"state"`
	Model   string `json:"model"`
	Updated bool   `json:"updated"`
}

type demuxReviewPayload struct {
	Valid       bool                 `json:"valid"`
	Proposal    demuxProposalPayload `json:"proposal"`
	Errors      []string             `json:"errors"`
	RepairHints []struct {
		Kind       string `json:"kind"`
		RevisionID string `json:"revision_id"`
		DependsOn  string `json:"depends_on"`
	} `json:"repair_hints"`
}

type demuxWorkflowPayload struct {
	Action        string               `json:"action"`
	State         string               `json:"state"`
	NextTool      string               `json:"next_tool"`
	FinalTool     string               `json:"final_tool"`
	Workflow      []string             `json:"workflow"`
	Proposal      demuxProposalPayload `json:"proposal"`
	Review        demuxReviewPayload   `json:"review"`
	RevisionShape map[string]any       `json:"revision_shape"`
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
	alphaBookmark := "gx/feat-alpha"
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

func TestGXDemuxProposesAndAppliesFileLevelRevisions(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("session-alpha")
	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "messy work")

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "compose", "--plan", "--json", "--intent", "demux e2e")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
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

	list := h.gx("compose", "list")
	for _, want := range []string{
		"Compose proposals",
		"* d1",
		proposal.ID,
		"pending / 2 revisions / 2 files",
		"demux e2e: alpha.txt",
		"default for revisions: gx compose show <revision-id>",
	} {
		if !strings.Contains(list, want) {
			t.Fatalf("gx compose list missing %q in:\n%s", want, list)
		}
	}

	showProposal := h.gx("compose", "show", "d1")
	for _, want := range []string{
		"Compose proposal",
		"demux e2e: alpha.txt",
		"JSON gx compose --json",
	} {
		if !strings.Contains(showProposal, want) {
			t.Fatalf("gx compose show d1 missing %q in:\n%s", want, showProposal)
		}
	}

	checkProposal := h.gx("compose", "review", "d1")
	for _, want := range []string{
		"Compose proposal",
		"Proposes 2 revisions",
		"Diagnostics",
		"hidden; use --raw or --json",
		"gx compose",
	} {
		if !strings.Contains(checkProposal, want) {
			t.Fatalf("gx compose review d1 missing %q in:\n%s", want, checkProposal)
		}
	}
	for _, unwanted := range []string{
		"Feasibility warnings",
		"hunk h1 in alpha.txt is not mapped to an enclosing symbol",
	} {
		if strings.Contains(checkProposal, unwanted) {
			t.Fatalf("gx compose review d1 should hide %q by default:\n%s", unwanted, checkProposal)
		}
	}

	show := h.gx("compose", "show", "u1")
	for _, want := range []string{
		"Compose revision",
		"Proposal " + proposal.ID,
		"Revision u1",
		"Intent demux e2e: alpha.txt",
		"Files",
		"alpha.txt",
		"Diff",
		"diff --git a/alpha.txt b/alpha.txt",
		"+alpha",
	} {
		if !strings.Contains(show, want) {
			t.Fatalf("gx compose show u1 missing %q in:\n%s", want, show)
		}
	}

	rawApply := h.gx("compose", "apply", proposal.ID, "--json")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{
		"demux e2e: alpha.txt": {"session-alpha"},
		"demux e2e: beta.txt":  {"session-alpha"},
	})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "u1", Intent: "demux e2e: alpha.txt", HunkIDsJSON: `["h1"]`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "u2", Intent: "demux e2e: beta.txt", HunkIDsJSON: `["h2"]`, ProvenanceStatus: "explicit"},
	})
}

func TestGXComposeApplyAllCreatesEveryProposedStackAndPublishPublishesAll(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")
	h.writeMultiStackComposeFiles()

	rawProposal := h.gx("compose", "--json", "--intent", "multi compose")
	var packet struct {
		State    string               `json:"state"`
		Proposal demuxProposalPayload `json:"proposal"`
	}
	if err := json.Unmarshal([]byte(rawProposal), &packet); err != nil {
		t.Fatalf("decode compose packet: %v\n%s", err, rawProposal)
	}
	if packet.State != "ready_to_apply" || packet.Proposal.ID == "" {
		t.Fatalf("compose packet = state %q proposal %#v", packet.State, packet.Proposal)
	}
	wantStacks := map[string]bool{
		"feature/demux-routing":    false,
		"feature/stack-management": false,
		"test/e2e-tests":           false,
	}
	for _, revision := range packet.Proposal.Revisions {
		if _, ok := wantStacks[revision.TargetStack]; ok {
			wantStacks[revision.TargetStack] = true
		}
	}
	for bookmark, seen := range wantStacks {
		if !seen {
			t.Fatalf("compose proposal missing stack %q in revisions %#v\n%s", bookmark, packet.Proposal.Revisions, rawProposal)
		}
	}

	rawApply := h.gx("compose", "apply", packet.Proposal.ID, "--json")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode compose apply: %v\n%s", err, rawApply)
	}
	if applied.RemainingChanges || applied.AcceptedSubset || applied.NextAction != "" {
		t.Fatalf("compose apply partial fields = remaining:%t subset:%t next:%q\n%s", applied.RemainingChanges, applied.AcceptedSubset, applied.NextAction, rawApply)
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

	output := h.gx("publish", "--github")
	if !strings.Contains(output, "Published") || !strings.Contains(output, "3 stacks") {
		t.Fatalf("publish output = %q, want published 3 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
	for bookmark := range wantStacks {
		assertRemoteBranchExists(t, h, bookmark)
		assertStack(t, h, bookmark, "published", h.bookmarkTargets()[bookmark], "refs/heads/"+bookmark)
	}
}

func TestGXComposeSelectedStackAppearsInStacksBeforeRerun(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")
	h.writeMultiStackComposeFiles()

	rawProposal := h.gx("compose", "--json", "--intent", "multi compose")
	var packet struct {
		State    string               `json:"state"`
		Proposal demuxProposalPayload `json:"proposal"`
	}
	if err := json.Unmarshal([]byte(rawProposal), &packet); err != nil {
		t.Fatalf("decode compose packet: %v\n%s", err, rawProposal)
	}
	selectedStack := "feature/stack-management"
	planFile := h.writeDemuxPlanForStack(packet.Proposal, selectedStack)

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode compose apply-plan: %v\n%s", err, rawApply)
	}
	if !applied.RemainingChanges || !applied.AcceptedSubset || applied.NextAction != "gx compose" {
		t.Fatalf("selected-stack apply fields = remaining:%t subset:%t next:%q\n%s", applied.RemainingChanges, applied.AcceptedSubset, applied.NextAction, rawApply)
	}
	assertCurrentBranch(t, h, "main")
	assertStack(t, h, selectedStack, "draft", h.bookmarkTargets()[selectedStack], "")
	assertStackInStacksJSON(t, h, selectedStack, 1)

	rawNext := h.gx("compose", "--json", "--intent", "multi compose followup")
	var next struct {
		State    string               `json:"state"`
		Proposal demuxProposalPayload `json:"proposal"`
	}
	if err := json.Unmarshal([]byte(rawNext), &next); err != nil {
		t.Fatalf("decode followup compose packet: %v\n%s", err, rawNext)
	}
	if next.State != "ready_to_apply" || len(next.Proposal.Revisions) != 2 {
		t.Fatalf("followup compose = state %q revisions %#v\n%s", next.State, next.Proposal.Revisions, rawNext)
	}
	h.gx("compose", "apply", next.Proposal.ID, "--json")
	assertCurrentBranch(t, h, "main")
	assertStackInStacksJSON(t, h, "feature/demux-routing", 1)
	assertStackInStacksJSON(t, h, "test/e2e-tests", 1)
}

func TestGXComposeApplyPlanAppendsToExistingJJBookmark(t *testing.T) {
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

	h.writeMultiStackComposeFiles()
	rawProposal := h.gx("compose", "--json", "--intent", "append existing compose")
	var packet struct {
		State    string               `json:"state"`
		Proposal demuxProposalPayload `json:"proposal"`
	}
	if err := json.Unmarshal([]byte(rawProposal), &packet); err != nil {
		t.Fatalf("decode compose packet: %v\n%s", err, rawProposal)
	}
	planFile := h.writeDemuxPlanForStack(packet.Proposal, existingStack)

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode compose apply-plan: %v\n%s", err, rawApply)
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

	rawProposal := source.gx("compose", "--json", "--intent", "update throwaway sync markdown")
	var packet struct {
		State    string               `json:"state"`
		Proposal demuxProposalPayload `json:"proposal"`
	}
	if err := json.Unmarshal([]byte(rawProposal), &packet); err != nil {
		t.Fatalf("decode compose packet: %v\n%s", err, rawProposal)
	}
	if packet.State != "ready_to_apply" || len(packet.Proposal.Revisions) != 1 {
		t.Fatalf("compose packet = state %q revisions %#v\n%s", packet.State, packet.Proposal.Revisions, rawProposal)
	}
	branch := packet.Proposal.Revisions[0].TargetStack
	if branch == "" || strings.HasPrefix(branch, "gx/") {
		t.Fatalf("compose target_stack = %q, want conventional public branch", branch)
	}

	source.gx("compose", "apply", packet.Proposal.ID, "--json")
	assertCurrentBranch(t, source, "main")
	assertStackInStacksJSON(t, source, branch, 1)
	source.gx("publish", branch)
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
	if !strings.Contains(output, "Cloud catch-up 1 bookmark(s) fetched from remote") {
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

func TestGXDemuxWorkflowHappyPathAppliesReviewedProposal(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("session-alpha")
	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "messy workflow")

	rawPacket := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "compose", "--plan", "--json", "--intent", "workflow")
	var packet demuxWorkflowPayload
	if err := json.Unmarshal([]byte(rawPacket), &packet); err != nil {
		t.Fatalf("decode demux packet: %v\n%s", err, rawPacket)
	}
	if packet.Action != "demux_changes" || packet.State != "ready_to_apply" || packet.NextTool != "gx_apply_revision_plan" || packet.FinalTool != "gx_apply_revision_plan" {
		t.Fatalf("demux state = %#v, want ready-to-apply packet", packet)
	}
	if len(packet.Workflow) == 0 || packet.RevisionShape["intent"] == nil {
		t.Fatalf("demux packet missing workflow guidance: %#v", packet)
	}
	if !packet.Review.Valid || len(packet.Review.Errors) != 0 {
		t.Fatalf("demux changes review = %#v, want valid", packet.Review)
	}
	if len(packet.Review.Proposal.Revisions) != 2 {
		t.Fatalf("reviewed proposal revisions = %#v, want 2", packet.Review.Proposal.Revisions)
	}

	planData, err := json.Marshal(packet.Review.Proposal)
	if err != nil {
		t.Fatalf("marshal reviewed proposal: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "reviewed-plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		t.Fatalf("write reviewed proposal: %v", err)
	}

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply-plan: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{
		"workflow: alpha.txt": {"session-alpha"},
		"workflow: beta.txt":  {"session-alpha"},
	})
	assertDemuxEvidence(t, h, packet.Proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "u1", Intent: "workflow: alpha.txt", HunkIDsJSON: `["h1"]`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "u2", Intent: "workflow: beta.txt", HunkIDsJSON: `["h2"]`, ProvenanceStatus: "explicit"},
	})
}

func TestGXDemuxAttachesRepoLocalSessionWithoutEnv(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("cursor-session")
	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "messy cursor work")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "cursor demux")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	if len(proposal.Revisions) != 1 {
		t.Fatalf("proposal revisions = %#v, want 1", proposal.Revisions)
	}
	if proposal.Revisions[0].ProvenanceStatus != "repo_local" || !reflect.DeepEqual(proposal.Revisions[0].SessionIDs, []string{"cursor-session"}) {
		t.Fatalf("proposal provenance = %#v, want repo-local cursor session", proposal.Revisions[0])
	}

	h.gx("compose", "apply", proposal.ID, "--json")
	assertChangeSessions(t, h, map[string][]string{
		"cursor demux: alpha.txt": {"cursor-session"},
	})
}

func TestGXDemuxReviewPlanReportsRepairableErrorsWithoutApplying(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", numberedLines(20, nil))
	h.run("jj", "describe", "-m", "seed alpha")
	h.gx("add", "-m", "seed alpha")
	h.gx("base", "--set", "gx/seed-alpha")

	h.writeTrackedFile("alpha.txt", numberedLines(20, map[int]string{
		2:  "line 02 changed",
		18: "line 18 changed",
	}))
	h.run("jj", "describe", "-m", "messy alpha")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "review alpha")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	if len(proposal.Hunks) < 2 {
		t.Fatalf("proposal hunks = %#v, want at least 2", proposal.Hunks)
	}

	plan := map[string]any{
		"id":                 proposal.ID,
		"repo_root":          proposal.RepoRoot,
		"proposed_change_id": proposal.ProposedChangeID,
		"proposed_commit_id": proposal.ProposedCommitID,
		"status":             "pending",
		"hunks":              proposal.Hunks,
		"revisions": []map[string]any{{
			"id":        "r1",
			"intent":    "only first hunk",
			"use_hunks": true,
			"hunk_ids":  []string{proposal.Hunks[0].ID},
		}},
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal invalid plan: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		t.Fatalf("write invalid plan: %v", err)
	}

	rawReview := h.gx("compose", "review-plan", "--json", "--plan-file", planFile)
	var reviewed demuxReviewPayload
	if err := json.Unmarshal([]byte(rawReview), &reviewed); err != nil {
		t.Fatalf("decode demux review-plan: %v\n%s", err, rawReview)
	}
	if reviewed.Valid || len(reviewed.Errors) != 1 || !strings.Contains(reviewed.Errors[0], "is not assigned to any revision") {
		t.Fatalf("review payload = %#v, want repairable unassigned hunk error", reviewed)
	}
	assertChangeSessions(t, h, nil)
	assertDemuxEvidence(t, h, proposal.ID, nil)
}

func TestGXDemuxApplyPlanResolvesHunkIDs(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", numberedLines(20, nil))
	h.run("jj", "describe", "-m", "seed alpha")
	h.gx("add", "-m", "seed alpha")
	h.gx("base", "--set", "gx/seed-alpha")

	h.insertSession("session-alpha")
	h.writeTrackedFile("alpha.txt", numberedLines(20, map[int]string{
		2:  "line 02 changed",
		18: "line 18 changed",
	}))
	h.run("jj", "describe", "-m", "messy alpha")

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "compose", "--plan", "--json", "--intent", "split alpha")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	if len(proposal.Hunks) < 2 {
		t.Fatalf("proposal hunks = %#v, want at least 2", proposal.Hunks)
	}

	plan := map[string]any{
		"id":                 proposal.ID,
		"repo_root":          proposal.RepoRoot,
		"proposed_change_id": proposal.ProposedChangeID,
		"proposed_commit_id": proposal.ProposedCommitID,
		"status":             "pending",
		"hunks":              proposal.Hunks,
		"revisions": []map[string]any{{
			"id":                "r1",
			"intent":            "split alpha hunks",
			"use_hunks":         true,
			"hunk_ids":          []string{proposal.Hunks[0].ID, proposal.Hunks[1].ID},
			"provenance_status": "explicit",
			"session_ids":       []string{"session-alpha"},
			"confidence":        0.8,
		}},
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal hunk plan: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		t.Fatalf("write hunk plan: %v", err)
	}

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply-plan: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 1 {
		t.Fatalf("applied revisions = %#v, want 1", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{"split alpha hunks": {"session-alpha"}})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "r1", Intent: "split alpha hunks", HunkIDsJSON: `["` + proposal.Hunks[0].ID + `","` + proposal.Hunks[1].ID + `"]`, ProvenanceStatus: "explicit"},
	})
}

func TestGXDemuxApplyPlanRoutesRevisionToExistingStack(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("target-anchor.txt", "target\n")
	h.run("jj", "describe", "-m", "target stack")
	h.gx("add", "-m", "target stack")
	targetBookmark := "gx/target-stack"
	targetBefore := h.bookmarkTargets()[targetBookmark]
	if targetBefore == "" {
		t.Fatalf("missing target bookmark before routed demux: %#v", h.bookmarkTargets())
	}

	h.run("git", "switch", "main")
	h.writeTrackedFile("source-anchor.txt", "source\n")
	h.run("jj", "describe", "-m", "source stack")
	h.gx("add", "-m", "source stack")
	sourceBookmark := "gx/source-stack"
	h.gx("base", "--set", sourceBookmark)

	h.insertSession("session-route")
	h.writeTrackedFile("routed.txt", "routed\n")
	h.writeTrackedFile("local.txt", "local\n")
	h.run("jj", "describe", "-m", "messy routed source")

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-route"}, "compose", "--plan", "--json", "--intent", "route e2e")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}

	plan := map[string]any{
		"id":                 proposal.ID,
		"repo_root":          proposal.RepoRoot,
		"proposed_change_id": proposal.ProposedChangeID,
		"proposed_commit_id": proposal.ProposedCommitID,
		"status":             "pending",
		"revisions": []map[string]any{
			{
				"id":                "r1",
				"intent":            "route e2e routed file",
				"files":             []string{"routed.txt"},
				"target_stack":      targetBookmark,
				"route_source":      "user",
				"route_reason":      "e2e routes this file to an existing stack",
				"route_confidence":  1.0,
				"provenance_status": "explicit",
				"session_ids":       []string{"session-route"},
				"confidence":        0.9,
			},
			{
				"id":                "r2",
				"intent":            "route e2e local file",
				"files":             []string{"local.txt"},
				"provenance_status": "explicit",
				"session_ids":       []string{"session-route"},
				"confidence":        0.9,
			},
		},
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal routed plan: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "routed-plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		t.Fatalf("write routed plan: %v", err)
	}

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply-plan: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	targetAfter := h.bookmarkTargets()[targetBookmark]
	sourceAfter := h.bookmarkTargets()[sourceBookmark]
	if targetAfter == "" || targetAfter == targetBefore {
		t.Fatalf("target bookmark after routed demux = %q, before %q", targetAfter, targetBefore)
	}
	if sourceAfter == "" {
		t.Fatalf("missing source bookmark after routed demux: %#v", h.bookmarkTargets())
	}
	assertCurrentBranch(t, h, "main")
	assertChangeSessions(t, h, map[string][]string{
		"route e2e local file":  {"session-route"},
		"route e2e routed file": {"session-route"},
	})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "r1", Intent: "route e2e routed file", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "r2", Intent: "route e2e local file", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
	})
}

func TestGXDemuxFixPlanRepairsStackRoutes(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("target-anchor.txt", "target\n")
	h.run("jj", "describe", "-m", "target stack")
	h.gx("add", "-m", "target stack")
	targetBookmark := "gx/target-stack"

	h.run("git", "switch", "main")
	h.writeTrackedFile("source-anchor.txt", "source\n")
	h.run("jj", "describe", "-m", "source stack")
	h.gx("add", "-m", "source stack")
	sourceBookmark := "gx/source-stack"
	h.gx("base", "--set", sourceBookmark)

	h.writeTrackedFile("routed.txt", "routed\n")
	h.writeTrackedFile("source-local.txt", "source local\n")
	h.run("jj", "describe", "-m", "messy repair route")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "repair route e2e")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	planFile := h.writeDemuxPlan(proposal, []map[string]any{
		{
			"id":                "r1",
			"intent":            "repair route e2e routed file",
			"files":             []string{"routed.txt"},
			"target_stack":      targetBookmark,
			"route_source":      "user",
			"route_reason":      "route to target stack",
			"route_confidence":  1.0,
			"provenance_status": "absent",
			"confidence":        0.5,
		},
		{
			"id":                "r2",
			"intent":            "repair route e2e source file",
			"files":             []string{"source-local.txt"},
			"depends_on":        []string{"r1"},
			"target_stack":      "s1",
			"route_source":      "user",
			"route_reason":      "alias route missing base stack",
			"route_confidence":  1.0,
			"provenance_status": "absent",
			"confidence":        0.5,
		},
	})

	result := runCommandAllowError(h.t, h.repo, h.env(), h.bin, "compose", "apply-plan", "--json", "--plan-file", planFile)
	if result.exitCode == 0 || !strings.Contains(result.output, "invalid route") {
		t.Fatalf("bad routed apply-plan = exit %d:\n%s", result.exitCode, result.output)
	}

	rawFix := h.gx("compose", "fix", proposal.ID, "--plan", "--json")
	var fixed demuxFixPayload
	if err := json.Unmarshal([]byte(rawFix), &fixed); err != nil {
		t.Fatalf("decode demux fix: %v\n%s", err, rawFix)
	}
	if !fixed.Updated || fixed.Model != "deterministic" || fixed.State != "ready_to_apply" || !fixed.Review.Valid {
		t.Fatalf("demux fix = %#v, want deterministic ready repair", fixed)
	}
	if len(fixed.Review.RepairHints) != 0 {
		t.Fatalf("demux fix repair hints = %#v, want none", fixed.Review.RepairHints)
	}
	if len(fixed.Proposal.Revisions) != 2 {
		t.Fatalf("fixed revisions = %#v, want 2", fixed.Proposal.Revisions)
	}
	if fixed.Proposal.Revisions[0].TargetStack != targetBookmark {
		t.Fatalf("r1 target_stack = %q, want %q", fixed.Proposal.Revisions[0].TargetStack, targetBookmark)
	}
	if fixed.Proposal.Revisions[1].TargetStack != sourceBookmark {
		t.Fatalf("r2 target_stack = %q, want %q", fixed.Proposal.Revisions[1].TargetStack, sourceBookmark)
	}
	if fixed.Proposal.Revisions[1].BaseStack != targetBookmark {
		t.Fatalf("r2 base_stack = %q, want %q", fixed.Proposal.Revisions[1].BaseStack, targetBookmark)
	}
}

func TestGXDemuxApplyPlanRoutesRevisionToNewStack(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("source-anchor.txt", "source\n")
	h.run("jj", "describe", "-m", "source stack")
	h.gx("add", "-m", "source stack")
	sourceBookmark := "gx/source-stack"
	h.gx("base", "--set", sourceBookmark)

	h.insertSession("session-new-route")
	h.writeTrackedFile("new-stack.txt", "new stack\n")
	h.writeTrackedFile("source-local.txt", "source local\n")
	h.run("jj", "describe", "-m", "messy new stack route")

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-new-route"}, "compose", "--plan", "--json", "--intent", "new route e2e")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}

	newBookmark := "gx/new-routed-stack"
	plan := map[string]any{
		"id":                 proposal.ID,
		"repo_root":          proposal.RepoRoot,
		"proposed_change_id": proposal.ProposedChangeID,
		"proposed_commit_id": proposal.ProposedCommitID,
		"status":             "pending",
		"revisions": []map[string]any{
			{
				"id":                "r1",
				"intent":            "new route e2e new stack file",
				"files":             []string{"new-stack.txt"},
				"target_stack":      newBookmark,
				"base_stack":        "main",
				"route_source":      "user",
				"route_reason":      "e2e creates a new routed stack",
				"route_confidence":  1.0,
				"provenance_status": "explicit",
				"session_ids":       []string{"session-new-route"},
				"confidence":        0.9,
			},
			{
				"id":                "r2",
				"intent":            "new route e2e source file",
				"files":             []string{"source-local.txt"},
				"provenance_status": "explicit",
				"session_ids":       []string{"session-new-route"},
				"confidence":        0.9,
			},
		},
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal new route plan: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "new-route-plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		t.Fatalf("write new route plan: %v", err)
	}

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply-plan: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	bookmarks := h.bookmarkTargets()
	if bookmarks[newBookmark] == "" {
		t.Fatalf("missing new routed bookmark after demux: %#v", bookmarks)
	}
	if bookmarks[sourceBookmark] == "" {
		t.Fatalf("missing source bookmark after demux: %#v", bookmarks)
	}
	assertCurrentBranch(t, h, "main")
	assertChangeSessions(t, h, map[string][]string{
		"new route e2e new stack file": {"session-new-route"},
		"new route e2e source file":    {"session-new-route"},
	})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "r1", Intent: "new route e2e new stack file", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "r2", Intent: "new route e2e source file", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
	})
}

func TestGXDemuxApplyPlanRoutesBaseRevisionToNewStacks(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.insertSession("session-base-route")
	h.writeTrackedFile("first-stack.txt", "first stack\n")
	h.writeTrackedFile("second-stack.txt", "second stack\n")
	h.run("jj", "describe", "-m", "messy base route")

	rawStatus := h.gx("status", "--agent")
	if !strings.Contains(rawStatus, `stack=""`) {
		t.Fatalf("status before routed apply = %s, want anonymous base source", rawStatus)
	}

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-base-route"}, "compose", "--plan", "--json", "--intent", "base route e2e")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}

	firstBookmark := "gx/base-routed-first"
	secondBookmark := "gx/base-routed-second"
	planFile := h.writeDemuxPlan(proposal, []map[string]any{
		{
			"id":                "r1",
			"intent":            "base route e2e first stack",
			"files":             []string{"first-stack.txt"},
			"target_stack":      firstBookmark,
			"base_stack":        "main",
			"route_source":      "user",
			"route_reason":      "e2e creates first stack from base source",
			"route_confidence":  1.0,
			"provenance_status": "explicit",
			"session_ids":       []string{"session-base-route"},
			"confidence":        0.9,
		},
		{
			"id":                "r2",
			"intent":            "base route e2e second stack",
			"files":             []string{"second-stack.txt"},
			"target_stack":      secondBookmark,
			"base_stack":        "main",
			"route_source":      "user",
			"route_reason":      "e2e creates second stack from base source",
			"route_confidence":  1.0,
			"provenance_status": "explicit",
			"session_ids":       []string{"session-base-route"},
			"confidence":        0.9,
		},
	})

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply-plan: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 2 {
		t.Fatalf("applied revisions = %#v, want 2", applied.Revisions)
	}
	bookmarks := h.bookmarkTargets()
	if bookmarks[firstBookmark] == "" || bookmarks[secondBookmark] == "" {
		t.Fatalf("missing routed bookmarks after base demux: %#v", bookmarks)
	}
	status := h.gx("status", "--agent")
	if !strings.Contains(status, "files=0") {
		t.Fatalf("source status after base routed apply = %s, want no remaining files", status)
	}
	assertChangeSessions(t, h, map[string][]string{
		"base route e2e first stack":  {"session-base-route"},
		"base route e2e second stack": {"session-base-route"},
	})
	assertDemuxEvidence(t, h, proposal.ID, []demuxEvidenceRow{
		{RevisionProposalID: "r1", Intent: "base route e2e first stack", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
		{RevisionProposalID: "r2", Intent: "base route e2e second stack", HunkIDsJSON: `null`, ProvenanceStatus: "explicit"},
	})
}

func TestGXDemuxRoutedApplyCopiesFilesAbsentFromTargetTree(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("target-anchor.txt", "target\n")
	h.run("jj", "describe", "-m", "target stack")
	h.gx("add", "-m", "target stack")
	targetBookmark := "gx/target-stack"

	h.run("git", "switch", "main")
	h.writeTrackedFile("source-anchor.txt", "source\n")
	h.run("jj", "describe", "-m", "source stack")
	h.gx("add", "-m", "source stack")
	sourceBookmark := "gx/source-stack"
	h.gx("base", "--set", sourceBookmark)

	if err := os.MkdirAll(filepath.Join(h.repo, "internal", "hooks"), 0o755); err != nil {
		t.Fatalf("mkdir internal/hooks: %v", err)
	}
	h.writeTrackedFile("internal/hooks/run.go", "package hooks\n")
	h.writeTrackedFile("internal/hooks/run_test.go", "package hooks\n")
	h.run("jj", "describe", "-m", "messy hook route")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "hook route")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	planFile := h.writeDemuxPlan(proposal, []map[string]any{{
		"id":                "r1",
		"intent":            "hook route files",
		"files":             []string{"internal/hooks/run.go", "internal/hooks/run_test.go"},
		"target_stack":      targetBookmark,
		"route_source":      "user",
		"route_reason":      "target stack lacks hook files",
		"route_confidence":  1.0,
		"provenance_status": "absent",
		"confidence":        0.5,
	}})

	rawApply := h.gx("compose", "apply-plan", "--json", "--plan-file", planFile)
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode routed apply: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 1 || applied.Revisions[0].Change.Description != "hook route files" {
		t.Fatalf("applied revisions = %#v, want routed hook revision", applied.Revisions)
	}
	status := h.gx("status")
	if strings.Contains(status, "internal/hooks/run.go") || strings.Contains(status, "internal/hooks/run_test.go") {
		t.Fatalf("source status after routed apply still includes hook files:\n%s", status)
	}
}

func TestGXComposeRoutedApplyRestoresSourceStackOnTargetRecordFailure(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("target-anchor.txt", "target\n")
	h.run("jj", "describe", "-m", "target stack")
	h.gx("add", "-m", "target stack")
	targetBookmark := "gx/target-stack"
	targetBefore := h.bookmarkTargets()[targetBookmark]

	h.run("git", "switch", "main")
	h.writeTrackedFile("source-anchor.txt", "source\n")
	h.run("jj", "describe", "-m", "source stack")
	h.gx("add", "-m", "source stack")
	sourceBookmark := "gx/source-stack"
	h.gx("base", "--set", sourceBookmark)

	h.writeTrackedFile("record-fail.txt", "record fail\n")
	h.run("jj", "describe", "-m", "messy record failure route")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "record failure route")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode compose proposal: %v\n%s", err, rawProposal)
	}
	planFile := h.writeDemuxPlan(proposal, []map[string]any{{
		"id":                "r1",
		"intent":            "(no description set)",
		"files":             []string{"record-fail.txt"},
		"target_stack":      targetBookmark,
		"route_source":      "user",
		"route_reason":      "force target record validation failure",
		"route_confidence":  1.0,
		"provenance_status": "absent",
		"confidence":        0.5,
	}})

	result := runCommandAllowError(h.t, h.repo, h.env(), h.bin, "compose", "apply-plan", "--json", "--plan-file", planFile)
	if result.exitCode == 0 || !strings.Contains(result.output, "(no description set)") {
		t.Fatalf("routed apply result = exit %d:\n%s", result.exitCode, result.output)
	}
	if got := h.bookmarkTargets()[targetBookmark]; got != targetBefore {
		t.Fatalf("target bookmark moved after failed routed record: got %q want %q", got, targetBefore)
	}
	assertCurrentBranch(t, h, "main")
	status := h.gx("status")
	if !strings.Contains(status, "record-fail.txt") {
		t.Fatalf("source status after rollback missing record-fail.txt:\n%s", status)
	}
}

func TestGXComposeRoutedApplyRestoresMainOnTargetRecordFailure(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("target-anchor.txt", "target\n")
	h.run("jj", "describe", "-m", "target stack")
	h.gx("add", "-m", "target stack")
	targetBookmark := "gx/target-stack"
	targetBefore := h.bookmarkTargets()[targetBookmark]

	h.run("git", "switch", "main")
	assertCurrentBranch(t, h, "main")
	h.writeTrackedFile("record-fail-main.txt", "record fail\n")
	h.run("jj", "describe", "-m", "messy main record failure route")

	rawProposal := h.gx("compose", "--plan", "--json", "--intent", "main record failure route")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode compose proposal: %v\n%s", err, rawProposal)
	}
	planFile := h.writeDemuxPlan(proposal, []map[string]any{{
		"id":                "r1",
		"intent":            "(no description set)",
		"files":             []string{"record-fail-main.txt"},
		"target_stack":      targetBookmark,
		"route_source":      "user",
		"route_reason":      "force target record validation failure",
		"route_confidence":  1.0,
		"provenance_status": "absent",
		"confidence":        0.5,
	}})

	result := runCommandAllowError(h.t, h.repo, h.env(), h.bin, "compose", "apply-plan", "--json", "--plan-file", planFile)
	if result.exitCode == 0 || !strings.Contains(result.output, "(no description set)") {
		t.Fatalf("routed apply result = exit %d:\n%s", result.exitCode, result.output)
	}
	if got := h.bookmarkTargets()[targetBookmark]; got != targetBefore {
		t.Fatalf("target bookmark moved after failed routed record: got %q want %q", got, targetBefore)
	}
	assertCurrentBranch(t, h, "main")
	status := h.gx("status")
	if !strings.Contains(status, "record-fail-main.txt") {
		t.Fatalf("source status after rollback missing record-fail-main.txt:\n%s", status)
	}
}

func TestGXDemuxProposesSymbolLevelRevisions(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(false)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("app.go", symbolFixture("one", "two"))
	h.run("jj", "describe", "-m", "seed symbols")
	h.gx("add", "-m", "seed symbols")
	h.gx("base", "--set", "gx/seed-symbols")

	h.insertSession("session-alpha")
	h.writeTrackedFile("app.go", symbolFixture("ONE", "TWO"))
	h.run("jj", "describe", "-m", "messy symbols")

	rawProposal := h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "compose", "--plan", "--json", "--intent", "split symbols")
	var proposal demuxProposalPayload
	if err := json.Unmarshal([]byte(rawProposal), &proposal); err != nil {
		t.Fatalf("decode demux proposal: %v\n%s", err, rawProposal)
	}
	if len(proposal.Revisions) != 1 {
		t.Fatalf("proposal revisions = %#v, want one hunk-level symbol revision", proposal.Revisions)
	}
	if !proposal.Revisions[0].UseHunks || len(proposal.Revisions[0].HunkIDs) != 2 {
		t.Fatalf("symbol revision = %#v, want both symbol hunks preserved", proposal.Revisions[0])
	}

	rawApply := h.gx("compose", "apply", proposal.ID, "--json")
	var applied demuxApplyPayload
	if err := json.Unmarshal([]byte(rawApply), &applied); err != nil {
		t.Fatalf("decode demux apply: %v\n%s", err, rawApply)
	}
	if len(applied.Revisions) != 1 {
		t.Fatalf("applied revisions = %#v, want 1", applied.Revisions)
	}
	assertChangeSessions(t, h, map[string][]string{proposal.Revisions[0].Intent: {"session-alpha"}})
}

func TestGXPublishAllPublishesEveryStack(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.insertSession("session-alpha")
	h.insertSessionRequest("session-alpha", "implement alpha")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-alpha"}, "add", "-m", "feat alpha")
	alphaBookmark := "gx/feat-alpha"
	alphaTarget := h.bookmarkTargets()[alphaBookmark]

	h.run("git", "switch", "main")
	h.writeTrackedFile("beta.txt", "beta\n")
	h.run("jj", "describe", "-m", "beta")
	h.insertSession("session-beta")
	h.insertSessionRequest("session-beta", "implement beta")
	h.gxWithEnv([]string{"GX_SESSION_ID=session-beta"}, "add", "-m", "feat beta")
	betaBookmark := "gx/feat-beta"
	betaTarget := h.bookmarkTargets()[betaBookmark]

	assertRemoteBranchMissing(t, h, alphaBookmark)
	assertRemoteBranchMissing(t, h, betaBookmark)

	output := h.gx("publish")
	if !strings.Contains(output, "Published") || !strings.Contains(output, "2 stacks") {
		t.Fatalf("publish output = %q, want published 2 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
	assertRemoteBranchExists(t, h, alphaBookmark)
	assertRemoteBranchExists(t, h, betaBookmark)
	assertStack(t, h, alphaBookmark, "published", alphaTarget, "refs/heads/"+alphaBookmark)
	assertStack(t, h, betaBookmark, "published", betaTarget, "refs/heads/"+betaBookmark)
	assertChangeBookmark(t, h, "feat alpha", alphaBookmark)
	assertChangeBookmark(t, h, "feat beta", betaBookmark)
}

func TestGXPublishNoStacksIsNoopOnMain(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	output := h.gx("publish")

	if !strings.Contains(output, "Published") || !strings.Contains(output, "0 stacks") {
		t.Fatalf("publish output = %q, want published 0 stacks", output)
	}
	assertCurrentBranch(t, h, "main")
}

func TestGXStacksHidesStackMergedIntoBase(t *testing.T) {
	h := newHarness(t)
	h.initGitRepo(true)
	h.gx("init", "--name", "Joe Example", "--email", "joe@example.com")

	h.writeTrackedFile("alpha.txt", "alpha\n")
	h.run("jj", "describe", "-m", "alpha")
	h.gx("add", "-m", "feat alpha")
	alphaBookmark := "gx/feat-alpha"
	assertStackInStacksJSON(t, h, alphaBookmark, 1)

	h.run("jj", "bookmark", "set", "main", "-r", alphaBookmark, "--allow-backwards")
	h.run("git", "switch", "main")

	assertCurrentBranch(t, h, "main")
	assertStackNotInStacksJSON(t, h, alphaBookmark)
	output := h.gx("publish")
	if !strings.Contains(output, "Published") || !strings.Contains(output, "0 stacks") {
		t.Fatalf("publish output = %q, want published 0 stacks", output)
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

func TestGXStacksShowsImplicitStackAliases(t *testing.T) {
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

	stacks := h.gx("stacks")
	if !strings.Contains(stacks, "feat beta  s1") {
		t.Fatalf("stacks missing beta alias:\n%s", stacks)
	}
	if !strings.Contains(stacks, "feat alpha  s2") {
		t.Fatalf("stacks missing alpha alias:\n%s", stacks)
	}
}

func numberedLines(count int, replacements map[int]string) string {
	var b strings.Builder
	for i := 1; i <= count; i++ {
		line := fmt.Sprintf("line %02d", i)
		if replacement, ok := replacements[i]; ok {
			line = replacement
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
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
			"token":        e2eCloudToken,
			"user_id":      "e2e-user",
			"login":        "e2e",
			"session_id":   "e2e-session",
			"machine_id":   "e2e-machine",
			"machine_name": "e2e-machine",
			"obtained_at":  time.Now().UTC().Format(time.RFC3339Nano),
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

func (h *harness) gxWithInput(input string, args ...string) commandResult {
	h.t.Helper()
	return runCommandAllowErrorWithInput(h.t, h.repo, h.env(), input, h.bin, args...)
}

func (h *harness) writeDemuxPlan(proposal demuxProposalPayload, revisions []map[string]any) string {
	h.t.Helper()
	plan := map[string]any{
		"id":                 proposal.ID,
		"repo_root":          proposal.RepoRoot,
		"proposed_change_id": proposal.ProposedChangeID,
		"proposed_commit_id": proposal.ProposedCommitID,
		"status":             "pending",
		"revisions":          revisions,
	}
	planData, err := json.Marshal(plan)
	if err != nil {
		h.t.Fatalf("marshal demux plan: %v", err)
	}
	planFile := filepath.Join(h.t.TempDir(), "demux-plan.json")
	if err := os.WriteFile(planFile, planData, 0o644); err != nil {
		h.t.Fatalf("write demux plan: %v", err)
	}
	return planFile
}

func (h *harness) writeDemuxPlanForStack(proposal demuxProposalPayload, targetStack string) string {
	h.t.Helper()
	var revisions []map[string]any
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != targetStack {
			continue
		}
		revisions = append(revisions, map[string]any{
			"id":                revision.ID,
			"intent":            revision.Intent,
			"files":             revision.Files,
			"target_stack":      revision.TargetStack,
			"route_source":      "user",
			"route_reason":      "selected stack from compose",
			"route_confidence":  1.0,
			"provenance_status": revision.ProvenanceStatus,
			"confidence":        1.0,
		})
	}
	if len(revisions) == 0 {
		h.t.Fatalf("proposal has no revisions for stack %q: %#v", targetStack, proposal.Revisions)
	}
	return h.writeDemuxPlan(proposal, revisions)
}

func (h *harness) writeMultiStackComposeFiles() {
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
	h.run("jj", "describe", "-m", "multi compose source")
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
	rawStacks := h.gx("stacks", "--json")
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
	rawStacks := h.gx("stacks", "--json")
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
