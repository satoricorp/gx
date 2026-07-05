package authoring

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

type demuxAlreadyAppliedFakeRunner struct {
	repoRoot            string
	calls               []string
	emptyRevisions      map[string]bool
	diffEmptyForFile    map[string]bool
	diffUnavailableFrom map[string]bool
	missingRevisions    map[string]bool
}

func (r *demuxAlreadyAppliedFakeRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	r.calls = append(r.calls, name+" "+strings.Join(args, " "))
	if name == "jj" && len(args) > 0 && args[0] == "root" {
		return r.repoRoot + "\n", nil
	}
	if name == "jj" && len(args) > 0 && args[0] == "op" {
		return "op-before\n", nil
	}
	return "", nil
}

func (r *demuxAlreadyAppliedFakeRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	if name != "jj" {
		return "", nil
	}
	if len(args) >= 2 && args[0] == "diff" && args[1] == "--from" {
		from := ""
		to := ""
		path := ""
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--to":
				if i+1 < len(args) {
					to = args[i+1]
					i++
				}
			default:
				if strings.HasPrefix(args[i], "-") {
					continue
				}
				if from == "" {
					from = args[i]
					continue
				}
				path = args[i]
			}
		}
		if from == "feature/existing" && to == "proposal-commit" && r.diffEmptyForFile != nil && r.diffEmptyForFile[path] {
			return "", nil
		}
		if r.diffUnavailableFrom != nil && r.diffUnavailableFrom[from] {
			return "", fmt.Errorf("Revision `%s` doesn't exist", from)
		}
		return "diff-content\n", nil
	}
	if len(args) >= 2 && args[0] == "log" && args[1] == "-r" {
		rev := args[2]
		if r.missingRevisions != nil && r.missingRevisions[rev] {
			return "", fmt.Errorf("Revision `%s` doesn't exist", rev)
		}
		template := ""
		for i := 0; i+1 < len(args); i++ {
			if args[i] == "-T" {
				template = args[i+1]
			}
		}
		switch template {
		case "empty":
			if r.emptyRevisions[rev] {
				return "true\n", nil
			}
			return "false\n", nil
		case "commit_id":
			return "recorded-commit\n", nil
		case "change_id":
			return "recorded-change\n", nil
		}
	}
	return "", nil
}

func (r *demuxAlreadyAppliedFakeRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

func TestExcludeDemuxAlreadyAppliedContentSkipsFullyAppliedRevision(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		diffEmptyForFile: map[string]bool{
			"internal/a.go": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))

	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(context.Background(), store, RepoInfo{RootPath: repoRoot, Backend: "jj"}, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "existing",
		BookmarkName: "feature/existing",
		BaseRef:      "main",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	proposal := DemuxProposal{
		RepoRoot:         repoRoot,
		ProposedCommitID: "proposal-commit",
		Hunks: []HunkRange{
			{ID: "h1", File: "internal/a.go"},
			{ID: "h2", File: "internal/b.go"},
		},
		Revisions: []RevisionProposal{
			{
				ID:          "u1",
				Intent:      "already landed",
				Files:       []string{"internal/a.go"},
				TargetStack: "feature/existing",
			},
			{
				ID:          "u2",
				Intent:      "still pending",
				Files:       []string{"internal/b.go"},
				TargetStack: "feature/existing",
			},
		},
	}
	filtered, err := engine.excludeDemuxAlreadyAppliedContent(context.Background(), proposal)
	if err != nil {
		t.Fatalf("excludeDemuxAlreadyAppliedContent() error = %v", err)
	}
	if len(filtered.Revisions) != 1 {
		t.Fatalf("revisions = %d, want 1 after excluding already-applied content", len(filtered.Revisions))
	}
	if filtered.Revisions[0].ID != "u2" {
		t.Fatalf("remaining revision = %q, want u2", filtered.Revisions[0].ID)
	}
	if len(filtered.Warnings) == 0 {
		t.Fatal("warnings empty, want skip notice for u1")
	}
	if len(filtered.Hunks) != 1 || filtered.Hunks[0].ID != "h2" {
		t.Fatalf("hunks = %#v, want only h2 after pruning already-applied h1", filtered.Hunks)
	}
}

func TestExcludeDemuxAlreadyAppliedContentToleratesStaleStackBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		diffUnavailableFrom: map[string]bool{
			"feature/stale": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))

	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(context.Background(), store, RepoInfo{RootPath: repoRoot, Backend: "jj"}, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	for _, stack := range []storage.Stack{
		{
			RepoID: repoID, Name: "stale", BookmarkName: "feature/stale",
			BaseRef: "main", Status: "draft", CreatedAt: 1, UpdatedAt: 1,
		},
	} {
		if _, err := store.UpsertStack(context.Background(), stack); err != nil {
			t.Fatalf("UpsertStack() error = %v", err)
		}
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	proposal := DemuxProposal{
		RepoRoot:         repoRoot,
		ProposedCommitID: "proposal-commit",
		Revisions: []RevisionProposal{{
			ID:          "u1",
			Intent:      "pending on stale stack",
			Files:       []string{".env.example"},
			TargetStack: "feature/stale",
		}},
	}
	filtered, err := engine.excludeDemuxAlreadyAppliedContent(context.Background(), proposal)
	if err != nil {
		t.Fatalf("excludeDemuxAlreadyAppliedContent() error = %v", err)
	}
	if len(filtered.Revisions) != 1 || filtered.Revisions[0].ID != "u1" {
		t.Fatalf("revisions = %#v, want u1 unchanged when dedup skipped", filtered.Revisions)
	}
	found := false
	for _, warning := range filtered.Warnings {
		if strings.Contains(warning, "feature/stale") && strings.Contains(warning, "skipped overlap check") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("warnings = %#v, want stale stack dedup skip warning", filtered.Warnings)
	}
}

func TestDemuxStackIndexOmitsStaleBookmarkStacks(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		missingRevisions: map[string]bool{
			"feature/stale": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))

	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(context.Background(), store, RepoInfo{RootPath: repoRoot, Backend: "jj"}, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	for _, stack := range []storage.Stack{
		{RepoID: repoID, Name: "good", BookmarkName: "feature/good", BaseRef: "main", Status: "draft", CreatedAt: 1, UpdatedAt: 1},
		{RepoID: repoID, Name: "stale", BookmarkName: "feature/stale", BaseRef: "main", Status: "draft", CreatedAt: 1, UpdatedAt: 1},
	} {
		if _, err := store.UpsertStack(context.Background(), stack); err != nil {
			t.Fatalf("UpsertStack() error = %v", err)
		}
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	index, warnings, err := engine.demuxStackIndex(context.Background(), repoRoot)
	if err != nil {
		t.Fatalf("demuxStackIndex() error = %v", err)
	}
	if len(index.stacks) != 1 || index.stacks[0].BookmarkName != "feature/good" {
		t.Fatalf("stacks = %#v, want only resolvable feature/good", index.stacks)
	}
	found := false
	for _, warning := range warnings {
		if strings.Contains(warning, "feature/stale") && strings.Contains(warning, "skipped overlap check") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("warnings = %#v, want stale stack routing warning", warnings)
	}
}

func TestDemuxFileContentMatchesRevisionUsesDiffRange(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		diffEmptyForFile: map[string]bool{
			"internal/a.go": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))
	matches, err := engine.demuxFileContentMatchesRevision(context.Background(), repoRoot, "feature/existing", "proposal-commit", "internal/a.go")
	if err != nil {
		t.Fatalf("demuxFileContentMatchesRevision() error = %v", err)
	}
	if !matches {
		t.Fatalf("matches = false, want true; calls=%#v", runner.calls)
	}
}

func TestValidateDemuxCheckpointResultTreatsAlreadyAppliedAsSkipSignal(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		emptyRevisions: map[string]bool{
			"empty-commit": true,
		},
		diffEmptyForFile: map[string]bool{
			"internal/a.go": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))
	revision := RevisionProposal{ID: "u1", Files: []string{"internal/a.go"}}
	result := CheckpointResult{Change: ChangeInfo{ChangeID: "empty-change", CommitID: "empty-commit"}}

	err := engine.validateDemuxCheckpointResult(context.Background(), repoRoot, revision, result, "feature/existing", "proposal-commit")
	if !isDemuxRevisionAlreadyApplied(err) {
		t.Fatalf("validateDemuxCheckpointResult() error = %v, want already-applied signal", err)
	}
}

func TestValidateDemuxCheckpointResultRejectsUnexpectedEmptyCommit(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &demuxAlreadyAppliedFakeRunner{
		repoRoot: repoRoot,
		emptyRevisions: map[string]bool{
			"empty-commit": true,
		},
	}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))
	revision := RevisionProposal{ID: "u1", Files: []string{"internal/b.go"}}
	result := CheckpointResult{Change: ChangeInfo{ChangeID: "empty-change", CommitID: "empty-commit"}}

	err := engine.validateDemuxCheckpointResult(context.Background(), repoRoot, revision, result, "feature/existing", "proposal-commit")
	if err == nil || isDemuxRevisionAlreadyApplied(err) {
		t.Fatalf("validateDemuxCheckpointResult() error = %v, want hard empty-commit error", err)
	}
}
