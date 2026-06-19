package reviewbundle

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/satoricorp/gx/internal/reviewsource"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func TestBuildPushIncludesDemuxEvidence(t *testing.T) {
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
		CurrentChange: &vcs.ChangeInfo{
			ChangeID: "change-1",
			CommitID: "commit-1",
		},
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "jj",
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if bundle.SchemaVersion != SchemaVersion {
		t.Fatalf("schema version = %d, want %d", bundle.SchemaVersion, SchemaVersion)
	}
	if bundle.Change == nil || len(bundle.Change.DemuxEvidence) != 1 {
		t.Fatalf("bundle change demux evidence = %#v", bundle.Change)
	}
	evidence := bundle.Change.DemuxEvidence[0]
	if evidence.DemuxProposalID != "demux-test" || evidence.RevisionProposalID != "r1" || !evidence.UseHunks {
		t.Fatalf("demux evidence = %#v", evidence)
	}
	if !reflect.DeepEqual(evidence.Files, []string{"alpha.txt"}) || !reflect.DeepEqual(evidence.HunkIDs, []string{"h1"}) {
		t.Fatalf("demux evidence files/hunks = %#v/%#v", evidence.Files, evidence.HunkIDs)
	}
	if string(evidence.Evidence) != reviewContextFixture() {
		t.Fatalf("evidence JSON = %s", evidence.Evidence)
	}
	if bundle.Change.ReviewContext == nil {
		t.Fatal("bundle change review context = nil")
	}
	if bundle.Change.ReviewContext.ProvenanceStatus != "explicit" || bundle.Change.ReviewContext.StructuralStatus != "available" {
		t.Fatalf("review context status = %#v", bundle.Change.ReviewContext)
	}
	if bundle.Change.ReviewContext.Risk.Level != "high" || bundle.Change.ReviewContext.Risk.Score == 0 {
		t.Fatalf("review context risk = %#v, want high risk", bundle.Change.ReviewContext.Risk)
	}
	if len(bundle.Change.ReviewContext.StructuralFacts) != 1 || bundle.Change.ReviewContext.StructuralFacts[0].File != "alpha.txt" {
		t.Fatalf("review context structural facts = %#v", bundle.Change.ReviewContext.StructuralFacts)
	}
	if len(bundle.Change.ReviewContext.ChangedSymbols) != 1 || bundle.Change.ReviewContext.ChangedSymbols[0].Symbol != "Run" {
		t.Fatalf("review context changed symbols = %#v", bundle.Change.ReviewContext.ChangedSymbols)
	}
	if len(bundle.Change.ReviewContext.FeasibilityWarnings) != 1 || bundle.Change.ReviewContext.FeasibilityWarnings[0].Source != "structural_dependency" {
		t.Fatalf("review context warnings = %#v", bundle.Change.ReviewContext.FeasibilityWarnings)
	}
	if !hasEvidenceKind(bundle.Change.ReviewContext.Evidence, "provenance") ||
		!hasEvidenceKind(bundle.Change.ReviewContext.Evidence, "structural") ||
		!hasEvidenceKind(bundle.Change.ReviewContext.Evidence, "risk") {
		t.Fatalf("review context evidence = %#v, want typed provenance/structural/risk evidence", bundle.Change.ReviewContext.Evidence)
	}
}

func TestBuildPushStackIncludesPatchesAndDedupedSessions(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	alphaID := seedChange(t, store, repoID, "alpha-change", "alpha-commit", "alpha", []string{"alpha.txt"})
	betaID := seedChange(t, store, repoID, "beta-change", "beta-commit", "beta", []string{"beta.txt"})
	if err := store.WriteSession(ctx, storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("WriteSession() error = %v", err)
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

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "beta-commit",
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "jj",
		},
		Published: []vcs.PushedChange{
			{
				Change:         vcs.ChangeInfo{ChangeID: "alpha-change", CommitID: "alpha-commit"},
				BranchName:     "gx/alpha",
				BaseBranchName: "main",
				Patch:          "diff --git a/alpha.txt b/alpha.txt\n",
			},
			{
				Change:         vcs.ChangeInfo{ChangeID: "beta-change", CommitID: "beta-commit"},
				BranchName:     "gx/beta",
				BaseBranchName: "gx/alpha",
				Patch:          "diff --git a/beta.txt b/beta.txt\n",
			},
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Stack) != 2 {
		t.Fatalf("bundle stack = %#v, want two entries", bundle.Stack)
	}
	if bundle.Stack[0].Patch == "" || bundle.Stack[1].Patch == "" {
		t.Fatalf("bundle stack patches missing: %#v", bundle.Stack)
	}
	if bundle.Stack[0].Change.ReviewContext == nil || bundle.Stack[1].Change.ReviewContext == nil {
		t.Fatalf("bundle stack review context missing: %#v", bundle.Stack)
	}
	if bundle.Stack[0].Change.ReviewContext.ProvenanceStatus != "linked" || bundle.Stack[0].Change.ReviewContext.LinkedSessionCount != 1 {
		t.Fatalf("alpha review context = %#v, want linked session provenance", bundle.Stack[0].Change.ReviewContext)
	}
	if len(bundle.Stack[0].Change.ReviewContext.ProvenanceSources) != 1 || bundle.Stack[0].Change.ReviewContext.ProvenanceSources[0].SessionID != "session-one" {
		t.Fatalf("alpha provenance sources = %#v, want session-one", bundle.Stack[0].Change.ReviewContext.ProvenanceSources)
	}
	if len(bundle.Stack[0].Change.ReviewContext.TranscriptSources) != 1 {
		t.Fatalf("alpha transcript sources = %#v, want one source", bundle.Stack[0].Change.ReviewContext.TranscriptSources)
	}
	transcriptSource := bundle.Stack[0].Change.ReviewContext.TranscriptSources[0]
	if transcriptSource.SessionID != "session-one" || transcriptSource.RequestID != "request-one" || transcriptSource.ResponseID == nil || *transcriptSource.ResponseID != "response-one" || transcriptSource.Status != "linked" {
		t.Fatalf("alpha transcript source = %#v, want request/response ids", transcriptSource)
	}
	if containsString(bundle.Stack[0].Change.ReviewContext.Risk.Signals, "missing_provenance") {
		t.Fatalf("alpha risk signals = %#v, did not expect missing provenance", bundle.Stack[0].Change.ReviewContext.Risk.Signals)
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
		CurrentChange: &vcs.ChangeInfo{
			ChangeID: "change-1",
			CommitID: "commit-1",
		},
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "jj",
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if bundle.Change == nil || bundle.Change.ReviewContext == nil {
		t.Fatalf("bundle change review context = %#v, want present", bundle.Change)
	}
	context := bundle.Change.ReviewContext
	if context.ProvenanceStatus != "absent" || context.StructuralStatus != "unavailable" {
		t.Fatalf("review context = %#v, want absent/unavailable", context)
	}
	if context.Risk.Level != "low" || !containsString(context.Risk.Signals, "missing_provenance") {
		t.Fatalf("risk = %#v, want low missing-provenance signal", context.Risk)
	}
}

func TestBuildPushIncludesAgentProvenance(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	changeID := seedChange(t, store, repoID, "change-1", "commit-1", "codex change", []string{"main.go"})
	source := "ambient"
	processName := "codex"
	if err := store.WriteSession(ctx, storage.Session{
		ID:          "session-one",
		CreatedAt:   1,
		Command:     "codex exec",
		Cwd:         "/repo",
		GXVersion:   "test",
		Source:      &source,
		ProcessName: &processName,
	}); err != nil {
		t.Fatalf("WriteSession() error = %v", err)
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
		CurrentChange: &vcs.ChangeInfo{
			ChangeID: "change-1",
			CommitID: "commit-1",
		},
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "jj",
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if bundle.Change == nil || bundle.Change.ReviewContext == nil {
		t.Fatalf("bundle change review context missing: %#v", bundle.Change)
	}
	got := bundle.Change.ReviewContext.AgentProvenance
	if len(got) != 1 {
		t.Fatalf("agent provenance = %#v, want one row", got)
	}
	if got[0].SessionID != "session-one" || got[0].AgentTool != "codex" || got[0].Provider != "openai" || got[0].ModelID != "gpt-5" {
		t.Fatalf("agent provenance = %#v, want codex/openai/gpt-5", got[0])
	}
}

func TestBuildPushJSONShapeIsStableForConsoleIngest(t *testing.T) {
	store := newBundleTestStore(t)
	ctx := context.Background()
	repoID := seedRepo(t, store, "/repo")
	seedChange(t, store, repoID, "change-1", "commit-1", "alpha", []string{"alpha.txt"})

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-1",
		CurrentChange: &vcs.ChangeInfo{
			ChangeID: "change-1",
			CommitID: "commit-1",
		},
		Repo: vcs.RepoInfo{
			RootPath: "/repo",
			Backend:  "jj",
		},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	data, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("decode bundle: %v", err)
	}
	for _, key := range []string{"event", "schema_version", "created_at", "gx_version", "repo", "push", "change", "sessions"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("bundle JSON missing %q: %s", key, data)
		}
	}
	change, ok := decoded["change"].(map[string]any)
	if !ok {
		t.Fatalf("bundle change JSON = %#v", decoded["change"])
	}
	reviewContext, ok := change["review_context"].(map[string]any)
	if !ok {
		t.Fatalf("bundle change missing review_context: %s", data)
	}
	for _, key := range []string{"provenance_status", "structural_status", "risk", "evidence"} {
		if _, ok := reviewContext[key]; !ok {
			t.Fatalf("review_context JSON missing %q: %s", key, data)
		}
	}
}

func TestArtifactJSONShapeIsFlattenedForReviewIngest(t *testing.T) {
	artifact := NewArtifact(Bundle{
		Event:         "gx.pr",
		SchemaVersion: SchemaVersion,
		Repo:          RepoPayload{RootPath: "/repo", Backend: "jj"},
		Push:          PushPayload{HeadCommitID: "commit-1"},
		Stack: []StackPayload{{
			BranchName: "gx/alpha",
			Patch:      "diff --git a/a.txt b/a.txt\n",
			Change: ChangePayload{
				JJChangeID:      "change-1",
				CurrentCommitID: "commit-1",
				Description:     "alpha",
			},
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
	for _, key := range []string{"review_id", "review_url", "index_status", "event", "schema_version", "repo", "push", "stack", "sessions"} {
		if _, ok := decoded[key]; !ok {
			t.Fatalf("artifact JSON missing %q: %s", key, data)
		}
	}
	if _, ok := decoded["bundle"]; ok {
		t.Fatalf("artifact JSON should be flattened, got nested bundle: %s", data)
	}
	stack, ok := decoded["stack"].([]any)
	if !ok || len(stack) != 1 {
		t.Fatalf("artifact stack = %#v, want one entry", decoded["stack"])
	}
	entry, ok := stack[0].(map[string]any)
	if !ok || entry["patch"] == "" {
		t.Fatalf("artifact stack entry = %#v, want patch for diffs ingest", stack[0])
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
		Backend:   "jj",
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
