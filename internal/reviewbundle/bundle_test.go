package reviewbundle

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/reviewsource"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func TestBuildPushIncludesDemuxEvidenceContext(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	changeID := seedChange(t, store, repoID, "change-1", "commit-1", "split alpha", []string{"alpha.txt"})
	if err := store.UpsertDemuxProposal(ctx, storage.DemuxProposal{
		ID:           "demux-test",
		RepoID:       repoID,
		BaseChangeID: "base-change",
		Status:       "applied",
		PayloadJSON:  `{"id":"demux-test"}`,
		CreatedAt:    10,
		UpdatedAt:    10,
	}); err != nil {
		t.Fatalf("UpsertDemuxProposal() error = %v", err)
	}
	if err := store.WriteChangeDemuxEvidence(ctx, storage.ChangeDemuxEvidence{
		ChangeID:           changeID,
		DemuxProposalID:    "demux-test",
		RevisionProposalID: "r1",
		Intent:             "split alpha",
		FilesJSON:          `["alpha.txt"]`,
		HunkIDsJSON:        `["h1"]`,
		UseHunks:           true,
		Confidence:         0.8,
		ProvenanceStatus:   "explicit",
		EvidenceJSON:       reviewContextFixture(),
		CreatedAt:          20,
	}); err != nil {
		t.Fatalf("WriteChangeDemuxEvidence() error = %v", err)
	}

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-1",
			RevisionID: "change-1",
			Message:    "split alpha",
			Files:      []string{"alpha.txt"},
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if bundle.SchemaVersion != SchemaVersion {
		t.Fatalf("schema version = %d, want %d", bundle.SchemaVersion, SchemaVersion)
	}
	if len(bundle.Revisions) != 1 {
		t.Fatalf("bundle revisions = %#v, want one entry", bundle.Revisions)
	}
	revision := bundle.Revisions[0]
	if revision.RevisionID != "change-1" || revision.CommitID != "commit-1" {
		t.Fatalf("revision identity = %#v", revision)
	}
	if revision.ReviewContext == nil {
		t.Fatal("revision review context = nil")
	}
	if revision.ReviewContext.ProvenanceStatus != "explicit" || revision.ReviewContext.StructuralStatus != "available" {
		t.Fatalf("review context status = %#v", revision.ReviewContext)
	}
	if revision.ReviewContext.Risk.Level != "high" || revision.ReviewContext.Risk.Score == 0 {
		t.Fatalf("review context risk = %#v, want high risk", revision.ReviewContext.Risk)
	}
	if len(revision.ReviewContext.StructuralFacts) != 1 || revision.ReviewContext.StructuralFacts[0].File != "alpha.txt" {
		t.Fatalf("review context structural facts = %#v", revision.ReviewContext.StructuralFacts)
	}
	if len(revision.ReviewContext.ChangedSymbols) != 1 || revision.ReviewContext.ChangedSymbols[0].Symbol != "Run" {
		t.Fatalf("review context changed symbols = %#v", revision.ReviewContext.ChangedSymbols)
	}
	if len(revision.ReviewContext.FeasibilityWarnings) != 1 || revision.ReviewContext.FeasibilityWarnings[0].Source != "structural_dependency" {
		t.Fatalf("review context warnings = %#v", revision.ReviewContext.FeasibilityWarnings)
	}
	if !hasEvidenceKind(revision.ReviewContext.Evidence, "provenance") ||
		!hasEvidenceKind(revision.ReviewContext.Evidence, "structural") ||
		!hasEvidenceKind(revision.ReviewContext.Evidence, "risk") {
		t.Fatalf("review context evidence = %#v, want typed provenance/structural/risk evidence", revision.ReviewContext.Evidence)
	}
}

func TestBuildPushIncludesPatchesAndDedupedSessions(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	alphaID := seedChange(t, store, repoID, "alpha-change", "alpha-commit", "alpha", []string{"alpha.txt"})
	betaID := seedChange(t, store, repoID, "beta-change", "beta-commit", "beta", []string{"beta.txt"})
	// UpsertObservedSession, not WriteSession: the latter has no production
	// callers, and seeding through it is what let a bundle test pass against a
	// `sessions` table production could never have written.
	if err := store.UpsertObservedSession(ctx, storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, alphaID, []string{"session-one"}, 2); err != nil {
		t.Fatalf("WriteChangeSessions(alpha) error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, betaID, []string{"session-one"}, 3); err != nil {
		t.Fatalf("WriteChangeSessions(beta) error = %v", err)
	}
	if err := store.WriteRequest(ctx, storage.Request{
		ID:             "request-one",
		SessionID:      "session-one",
		CreatedAt:      4,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		RequestBody:    []byte(`{"input":"why alpha?"}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest() error = %v", err)
	}
	if err := store.WriteResponse(ctx, storage.Response{
		ID:              "response-one",
		RequestID:       "request-one",
		CreatedAt:       5,
		CompletedAt:     6,
		StatusCode:      200,
		ResponseBody:    []byte(`{"output":"because alpha"}`),
		ResponseHeaders: "{}",
		DurationMS:      1,
	}); err != nil {
		t.Fatalf("WriteResponse() error = %v", err)
	}

	main := "main"
	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "beta-commit",
		Repo: vcs.RepoInfo{
			RootPath:      "/repo",
			Backend:       "git",
			DefaultBranch: &main,
		},
		GXStackRef: "feature/alpha",
		Commits: []vcs.PushedCommit{
			{
				CommitID:   "alpha-commit",
				RevisionID: "alpha-change",
				Message:    "alpha",
				Files:      []string{"alpha.txt"},
				Patch:      "diff --git a/alpha.txt b/alpha.txt\n",
			},
			{
				CommitID:   "beta-commit",
				RevisionID: "beta-change",
				Message:    "beta",
				Files:      []string{"beta.txt"},
				Patch:      "diff --git a/beta.txt b/beta.txt\n",
			},
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Revisions) != 2 {
		t.Fatalf("bundle revisions = %#v, want two entries", bundle.Revisions)
	}
	if bundle.Revisions[0].Patch == "" || bundle.Revisions[1].Patch == "" {
		t.Fatalf("bundle revision patches missing: %#v", bundle.Revisions)
	}
	if bundle.Revisions[0].BranchName != "feature/alpha" || bundle.Revisions[0].BaseBranchName != "main" {
		t.Fatalf("revision branches = %#v, want feature/alpha on main", bundle.Revisions[0])
	}
	if bundle.Revisions[0].ReviewContext == nil || bundle.Revisions[1].ReviewContext == nil {
		t.Fatalf("bundle revision review context missing: %#v", bundle.Revisions)
	}
	alphaContext := bundle.Revisions[0].ReviewContext
	if alphaContext.ProvenanceStatus != "linked" || alphaContext.LinkedSessionCount != 1 {
		t.Fatalf("alpha review context = %#v, want linked session provenance", alphaContext)
	}
	if len(alphaContext.ProvenanceSources) != 1 || alphaContext.ProvenanceSources[0].SessionID != "session-one" {
		t.Fatalf("alpha provenance sources = %#v, want session-one", alphaContext.ProvenanceSources)
	}
	if len(alphaContext.TranscriptSources) != 1 {
		t.Fatalf("alpha transcript sources = %#v, want one source", alphaContext.TranscriptSources)
	}
	transcriptSource := alphaContext.TranscriptSources[0]
	if transcriptSource.SessionID != "session-one" || transcriptSource.RequestID != "request-one" || transcriptSource.ResponseID == nil || *transcriptSource.ResponseID != "response-one" || transcriptSource.Status != "linked" {
		t.Fatalf("alpha transcript source = %#v, want request/response ids", transcriptSource)
	}
	if containsString(alphaContext.Risk.Signals, "missing_provenance") {
		t.Fatalf("alpha risk signals = %#v, did not expect missing provenance", alphaContext.Risk.Signals)
	}
	if len(bundle.Sessions) != 1 || bundle.Sessions[0].ID != "session-one" {
		t.Fatalf("bundle sessions = %#v, want deduped session-one", bundle.Sessions)
	}
}

func TestBuildPushIncludesReviewContextWithoutDemuxEvidence(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	seedChange(t, store, repoID, "change-1", "commit-1", "manual alpha", []string{"alpha.txt"})

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-1",
			RevisionID: "change-1",
			Message:    "manual alpha",
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Revisions) != 1 || bundle.Revisions[0].ReviewContext == nil {
		t.Fatalf("bundle revisions = %#v, want review context present", bundle.Revisions)
	}
	context := bundle.Revisions[0].ReviewContext
	if context.ProvenanceStatus != "absent" || context.StructuralStatus != "unavailable" {
		t.Fatalf("review context = %#v, want absent/unavailable", context)
	}
	if context.Risk.Level != "low" || !containsString(context.Risk.Signals, "missing_provenance") {
		t.Fatalf("risk = %#v, want low missing-provenance signal", context.Risk)
	}
	if !reflect.DeepEqual(bundle.Revisions[0].Files, []string{"alpha.txt"}) {
		t.Fatalf("revision files = %#v, want change-row fallback", bundle.Revisions[0].Files)
	}
}

func TestBuildPushIncludesAgentProvenance(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	changeID := seedChange(t, store, repoID, "change-1", "commit-1", "codex change", []string{"main.go"})
	source := "ambient"
	// UpsertObservedSession is the push path's writer and the only production
	// writer of `sessions` now that the cursor ingest is retired. It cannot
	// set process_name — nothing can any more — so the agent tool must resolve
	// from the command.
	if err := store.UpsertObservedSession(ctx, storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex exec",
		Cwd:       "/repo",
		GXVersion: "test",
		Source:    &source,
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	model := "gpt-5"
	if err := store.WriteRequest(ctx, storage.Request{
		ID:             "request-one",
		SessionID:      "session-one",
		CreatedAt:      2,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		Model:          &model,
		RequestBody:    []byte(`{"input":"change main"}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, changeID, []string{"session-one"}, 3); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-1",
			RevisionID: "change-1",
			Message:    "codex change",
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Revisions) != 1 || bundle.Revisions[0].ReviewContext == nil {
		t.Fatalf("bundle revision review context missing: %#v", bundle.Revisions)
	}
	got := bundle.Revisions[0].ReviewContext.AgentProvenance
	if len(got) != 1 {
		t.Fatalf("agent provenance = %#v, want one row", got)
	}
	if got[0].SessionID != "session-one" || got[0].AgentTool != "codex" || got[0].Provider != "openai" || got[0].ModelID != "gpt-5" {
		t.Fatalf("agent provenance = %#v, want codex/openai/gpt-5", got[0])
	}
}

func TestBuildPushUsesGXStackRefForPushBranchName(t *testing.T) {
	store := newBundleTestStore(t)
	_ = store
	ctx := context.Background()

	main := "main"
	stackRef := "feature/cli"
	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		GXStackRef:   stackRef,
		Repo: vcs.RepoInfo{
			RootPath:   "/repo",
			Backend:    "git",
			BranchName: &main,
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if bundle.Push.BranchName == nil || *bundle.Push.BranchName != stackRef {
		t.Fatalf("push branch = %#v, want %q", bundle.Push.BranchName, stackRef)
	}
}

func TestBuildPushIncludesTrailerlessCommits(t *testing.T) {
	store := newBundleTestStore(t)
	_ = store
	ctx := context.Background()

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-2",
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID: "commit-2",
			Message:  "no trailer here",
			Files:    []string{"a.txt"},
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Revisions) != 1 {
		t.Fatalf("bundle revisions = %#v, want the trailer-less commit included", bundle.Revisions)
	}
	if bundle.Revisions[0].RevisionID != "" {
		t.Fatalf("revision id = %q, want empty for trailer-less commit", bundle.Revisions[0].RevisionID)
	}
	if bundle.Revisions[0].ReviewContext != nil {
		t.Fatalf("review context = %#v, want none for trailer-less commit", bundle.Revisions[0].ReviewContext)
	}
}

// TestBuildPushMarshalsSchemaV2WireShape locks the v2 wire contract: the
// bundle must declare schema_version 2, carry `revisions` entries with exactly
// the agreed keys, and contain none of the retired v1 keys (`stack`, `change`,
// `jj_change_id`) or the retired process-shaped session keys anywhere.
//
// `command` is deliberately NOT in the retired set. The keys this test bans are
// process-shaped leftovers of the daemon (`client_pid`, `exit_code`,
// `process_name`, `parent_pid`); `command` now carries the agent that produced
// the transcript ("claude", "codex", "cursor"), which the console reads first
// and only then falls back to `source`. Dropping it rendered every session as
// `command=?`, so it is required here rather than forbidden.
func TestBuildPushMarshalsSchemaV2WireShape(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	changeID := seedChange(t, store, repoID, "gxr-3f9c", "0f4b21c9", "add retry helper\n\nCovers the timeout path.", []string{"retry.go"})
	if err := store.UpsertObservedSession(ctx, storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "claude",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, changeID, []string{"session-one"}, 2); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}
	if err := store.WriteRequest(ctx, storage.Request{
		ID:             "request-one",
		SessionID:      "session-one",
		CreatedAt:      3,
		Provider:       "anthropic",
		Endpoint:       "/v1/messages",
		Method:         "POST",
		RequestBody:    []byte(`{"input":"add retry"}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest() error = %v", err)
	}

	main := "main"
	remote := "origin"
	prURL := "https://github.com/acme/repo/pull/7"
	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "0f4b21c9",
		RemoteName:   &remote,
		Repo: vcs.RepoInfo{
			RootPath:      "/repo",
			Backend:       "jj", // stale stored value; the wire must still say git
			DefaultBranch: &main,
			BranchName:    &main,
		},
		GXStackRef:           "feature/retry",
		GitHubPullRequestURL: &prURL,
		Commits: []vcs.PushedCommit{{
			CommitID:   "0f4b21c9",
			RevisionID: "gxr-3f9c",
			Message:    "add retry helper\n\nCovers the timeout path.",
			Files:      []string{"retry.go"},
			Patch:      "diff --git a/retry.go b/retry.go\n@@ -0,0 +1 @@\n+package retry\n",
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}

	data, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	t.Logf("v2 bundle JSON:\n%s", indentJSON(t, data))

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode bundle: %v", err)
	}
	if got, ok := decoded["schema_version"].(float64); !ok || int(got) != 2 {
		t.Fatalf("schema_version = %#v, want 2", decoded["schema_version"])
	}
	if decoded["repo"].(map[string]any)["backend"] != "git" {
		t.Fatalf("repo.backend = %#v, want git", decoded["repo"])
	}
	for _, key := range []string{"event", "schema_version", "created_at", "gx_version", "repo", "push", "revisions", "sessions"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("bundle JSON missing %q: %s", key, data)
		}
	}
	revisions, ok := decoded["revisions"].([]any)
	if !ok || len(revisions) != 1 {
		t.Fatalf("bundle revisions JSON = %#v, want one entry", decoded["revisions"])
	}
	revision, ok := revisions[0].(map[string]any)
	if !ok {
		t.Fatalf("revision JSON = %#v", revisions[0])
	}
	for _, key := range []string{
		"revision_id", "commit_id", "description", "files",
		"branch_name", "base_branch_name", "patch",
		"github_pull_request_url", "review_context",
	} {
		if _, ok := revision[key]; !ok {
			t.Fatalf("revision JSON missing %q: %s", key, data)
		}
	}
	if revision["revision_id"] != "gxr-3f9c" || revision["commit_id"] != "0f4b21c9" {
		t.Fatalf("revision identity JSON = %#v", revision)
	}
	if revision["branch_name"] != "feature/retry" || revision["base_branch_name"] != "main" {
		t.Fatalf("revision branch JSON = %#v", revision)
	}

	for _, forbidden := range []string{`"stack"`, `"change"`, `"jj_change_id"`} {
		if strings.Contains(string(data), forbidden) {
			t.Fatalf("bundle JSON still contains retired key %s:\n%s", forbidden, data)
		}
	}
	sessions, ok := decoded["sessions"].([]any)
	if !ok || len(sessions) != 1 {
		t.Fatalf("bundle sessions JSON = %#v, want one entry", decoded["sessions"])
	}
	session, ok := sessions[0].(map[string]any)
	if !ok {
		t.Fatalf("session JSON = %#v", sessions[0])
	}
	for _, forbidden := range []string{"client_pid", "exit_code", "process_name", "parent_pid"} {
		if _, present := session[forbidden]; present {
			t.Fatalf("session JSON still contains retired key %q: %s", forbidden, data)
		}
	}
	for _, key := range []string{"id", "created_at", "command", "cwd", "gx_version", "requests"} {
		if _, ok := session[key]; !ok {
			t.Fatalf("session JSON missing %q: %s", key, data)
		}
	}
	// The console renders this verbatim; an empty string here is the `command=?`
	// bug, which a mere presence check would not catch.
	if session["command"] != "claude" {
		t.Fatalf("session command = %#v, want %q: %s", session["command"], "claude", data)
	}
}

func TestArtifactJSONShapeIsFlattenedForReviewIngest(t *testing.T) {
	artifact := NewArtifact(Bundle{
		Event:         "gx.pr",
		SchemaVersion: SchemaVersion,
		Repo:          RepoPayload{RootPath: "/repo", Backend: "git"},
		Push:          PushPayload{HeadCommitID: "commit-1"},
		Revisions: []RevisionPayload{{
			RevisionID:  "change-1",
			CommitID:    "commit-1",
			Description: "alpha",
			Files:       []string{"a.txt"},
			BranchName:  "feature/alpha",
			Patch:       "diff --git a/a.txt b/a.txt\n",
		}},
		Sessions: []SessionPayload{},
	})
	artifact.ReviewID = "review-one"
	artifact.ReviewURL = "http://gx.test/reviews/review-one"

	data, err := json.Marshal(artifact)
	if err != nil {
		t.Fatalf("marshal artifact: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode artifact: %v", err)
	}
	for _, key := range []string{"review_id", "review_url", "index_status", "event", "schema_version", "repo", "push", "revisions", "sessions"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("artifact JSON missing %q: %s", key, data)
		}
	}
	if _, ok := decoded["bundle"]; ok {
		t.Fatalf("artifact JSON should be flattened, got nested bundle: %s", data)
	}
	revisions, ok := decoded["revisions"].([]any)
	if !ok || len(revisions) != 1 {
		t.Fatalf("artifact revisions = %#v, want one entry", decoded["revisions"])
	}
	entry, ok := revisions[0].(map[string]any)
	if !ok || entry["patch"] == "" {
		t.Fatalf("artifact revision entry = %#v, want patch for diffs ingest", revisions[0])
	}
}

func TestBuildReviewContextComputesRiskFromEvidence(t *testing.T) {
	evidence := []DemuxEvidencePayload{
		{
			HunkIDs:          []string{"h1", "h2", "h3", "h4", "h5"},
			Confidence:       0.4,
			ProvenanceStatus: "repo_local",
			Evidence:         json.RawMessage(reviewContextFixture()),
		},
	}
	context := buildReviewContext(
		[]string{"a.go", "b.go", "c.go", "d.go", "e.go"},
		evidence,
		reviewsource.BuildGraph(evidenceStatuses(evidence), nil, nil),
		nil,
	)
	if context == nil {
		t.Fatal("buildReviewContext() = nil")
	}
	if context.Risk.Level != "high" {
		t.Fatalf("risk = %#v, want high", context.Risk)
	}
	for _, want := range []string{"many_files", "many_hunks", "repo_local_provenance", "low_demux_confidence", "structural_dependencies", "warning:structural_dependency"} {
		if !containsString(context.Risk.Signals, want) {
			t.Fatalf("risk signals = %#v, missing %q", context.Risk.Signals, want)
		}
	}
}

func TestProvenanceRiskHandlesEveryReportedStatus(t *testing.T) {
	// Every status reviewsource can report must produce a signal, so a new or
	// unrecognized status can never pass through risk scoring unnoticed.
	for _, status := range []string{
		reviewsource.StatusAbsent,
		reviewsource.StatusRepoLocal,
		reviewsource.StatusUnknown,
		reviewsource.StatusLinked,
		reviewsource.StatusExplicit,
		"some-status-we-have-not-seen",
	} {
		points, signal := provenanceRisk(status)
		if signal == "" {
			t.Fatalf("provenanceRisk(%q) returned no signal", status)
		}
		if points < 0 {
			t.Fatalf("provenanceRisk(%q) points = %d, want >= 0", status, points)
		}
	}
	if points, _ := provenanceRisk(reviewsource.StatusAbsent); points <= 0 {
		t.Fatal("absent provenance should raise risk")
	}
	if points, _ := provenanceRisk(reviewsource.StatusExplicit); points != 0 {
		t.Fatal("explicit provenance should not raise risk")
	}
	if points, signal := provenanceRisk(""); points != 0 || signal != "" {
		t.Fatalf("empty status = (%d, %q), want (0, \"\")", points, signal)
	}
}

func indentJSON(t *testing.T, data []byte) string {
	t.Helper()
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		t.Fatalf("indent bundle JSON: %v", err)
	}
	return out.String()
}

func newBundleTestStore(t *testing.T) *storage.Store {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func seedRepo(t *testing.T, store *storage.Store, root string) int64 {
	t.Helper()
	id, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  root,
		Backend:   "git",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	return id
}

func seedChange(t *testing.T, store *storage.Store, repoID int64, changeID, commitID, description string, files []string) int64 {
	t.Helper()
	id, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      changeID,
		CurrentCommitID: commitID,
		Description:     description,
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	filesJSON, err := json.Marshal(files)
	if err != nil {
		t.Fatalf("marshal files: %v", err)
	}
	if err := store.WriteChangeRevision(context.Background(), storage.ChangeRevision{
		ChangeID:      id,
		JJCommitID:    commitID,
		JJOperationID: "op-" + changeID,
		ChangedFiles:  string(filesJSON),
		CreatedAt:     1,
	}); err != nil {
		t.Fatalf("WriteChangeRevision() error = %v", err)
	}
	return id
}

func reviewContextFixture() string {
	return `{"revision":{"id":"r1"},"structural_facts":[{"file":"alpha.txt","language":"go","defined_symbols":["Run"],"symbols":[{"name":"Run","kind":"function","start_line":1,"end_line":5}]}],"structural_dependencies":[{"from_file":"alpha.txt","to_file":"beta.txt","symbol":"Beta"}],"changed_symbols":[{"hunk_id":"h1","file":"alpha.txt","symbol":"Run","kind":"function","start_line":1,"end_line":5}],"feasibility_warnings":[{"revision_id":"r1","severity":"warning","source":"structural_dependency","depends_on":"r2","symbol":"Beta","from_file":"alpha.txt","to_file":"beta.txt","message":"alpha references Beta"}]}`
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func hasEvidenceKind(evidence []ReviewEvidencePayload, kind string) bool {
	for _, item := range evidence {
		if item.Kind == kind {
			return true
		}
	}
	return false
}

// TestBuildPushFindsSessionsWhenRepoRootPathDiverges pins the bundle read path
// to the same repository identity the write path uses.
//
// storage.Store.UpsertRepo matches an existing repos row on git_common_dir
// first and its UPDATE never rewrites root_path, so the stored root_path is
// whatever worktree registered the row first — and can even be empty. The
// sessions this push wrote therefore hang off a repos row whose root_path is
// not the pushing worktree's. Reading them back by root_path found nothing and
// published `sessions: []` with no error anywhere.
func TestBuildPushFindsSessionsWhenRepoRootPathDiverges(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()

	// The row as UpsertRepo left it: identified by the git common dir, with a
	// root_path that does not match the pushing worktree.
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:     "/stale/other-worktree",
		GitCommonDir: "/repo/.git",
		Backend:      "git",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID := seedChange(t, store, repoID, "change-1", "commit-1", "alpha", []string{"alpha.txt"})
	if err := store.UpsertSession(ctx, storage.Session{
		ID: "session-one", CreatedAt: 1, Command: "claude", Cwd: "/repo", GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertSession() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, changeID, []string{"session-one"}, 1); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		Repo: vcs.RepoInfo{
			RootPath:     "/repo",
			GitCommonDir: "/repo/.git",
			Backend:      "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-1",
			RevisionID: "change-1",
			Message:    "alpha",
			Files:      []string{"alpha.txt"},
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Sessions) != 1 || bundle.Sessions[0].ID != "session-one" {
		t.Fatalf("bundle sessions = %#v, want session-one resolved via git_common_dir", bundle.Sessions)
	}
	if bundle.Sessions[0].Command != "claude" {
		t.Fatalf("bundle session command = %q, want claude", bundle.Sessions[0].Command)
	}
}

// TestBuildPushDoesNotReadAnotherReposChanges guards the replacement lookup:
// resolving by identity must still scope the revision to this repository.
func TestBuildPushDoesNotReadAnotherReposChanges(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()

	otherID := seedRepo(t, store, "/other")
	otherChange := seedChange(t, store, otherID, "change-1", "commit-1", "alpha", []string{"alpha.txt"})
	if err := store.UpsertSession(ctx, storage.Session{
		ID: "other-session", CreatedAt: 1, Command: "claude", Cwd: "/other", GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertSession() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, otherChange, []string{"other-session"}, 1); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}
	seedRepo(t, store, "/repo")

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		Repo:         vcs.RepoInfo{RootPath: "/repo", GitCommonDir: "/repo", Backend: "git"},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-1",
			RevisionID: "change-1",
			Message:    "alpha",
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Sessions) != 0 {
		t.Fatalf("bundle sessions = %#v, want none from a different repo", bundle.Sessions)
	}
}
