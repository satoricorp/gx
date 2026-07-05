package authoring

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/structural"
	"github.com/satoricorp/gx/internal/vcs"
)

func TestGroupFilesPairsSourceAndTests(t *testing.T) {
	got := groupFiles([]string{
		"internal/foo/foo_test.go",
		"internal/foo/foo.go",
		"internal/storage/schema.sql",
	})
	want := [][]string{
		{"internal/storage/schema.sql"},
		{"internal/foo/foo.go", "internal/foo/foo_test.go"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupFiles() = %#v, want %#v", got, want)
	}
}

func TestGroupFilesClustersGreenfieldTauriScaffold(t *testing.T) {
	got := groupFiles([]string{
		".gitignore",
		".vscode/extensions.json",
		"README.md",
		"index.html",
		"package-lock.json",
		"package.json",
		"src/App.css",
		"src/App.tsx",
		"src/main.tsx",
		"src/vite-env.d.ts",
		"src-tauri/.gitignore",
		"src-tauri/Cargo.lock",
		"src-tauri/Cargo.toml",
		"src-tauri/build.rs",
		"src-tauri/capabilities/default.json",
		"src-tauri/src/lib.rs",
		"src-tauri/src/main.rs",
		"src-tauri/tauri.conf.json",
		"tsconfig.json",
		"tsconfig.node.json",
		"vite.config.ts",
	})
	want := [][]string{
		{".gitignore", ".vscode/extensions.json", "package-lock.json", "package.json", "tsconfig.json", "tsconfig.node.json", "vite.config.ts"},
		{"src-tauri/.gitignore", "src-tauri/Cargo.lock", "src-tauri/Cargo.toml", "src-tauri/build.rs", "src-tauri/capabilities/default.json", "src-tauri/src/lib.rs", "src-tauri/src/main.rs", "src-tauri/tauri.conf.json"},
		{"index.html", "src/App.css", "src/App.tsx", "src/main.tsx", "src/vite-env.d.ts"},
		{"README.md"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupFiles() = %#v, want %#v", got, want)
	}
}

func TestCounterpartKeyPairsPackageManifestAndLock(t *testing.T) {
	got := groupFiles([]string{"package-lock.json", "package.json"})
	want := [][]string{{"package-lock.json", "package.json"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupFiles() = %#v, want %#v", got, want)
	}
}

func TestProposalIntentUsesPrefix(t *testing.T) {
	got := proposalIntent("add demuxr", []string{"internal/foo/foo.go", "internal/foo/foo_test.go"})
	if got != "add demuxr: internal/foo/foo" {
		t.Fatalf("proposalIntent() = %q", got)
	}
}

func TestCleanFilesSortsAndDedupes(t *testing.T) {
	got := cleanFiles([]string{"b.go", "", "a.go", "b.go"})
	want := []string{"a.go", "b.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("cleanFiles() = %#v, want %#v", got, want)
	}
}

func TestFilterFilesByFilesetsMatchesExactFile(t *testing.T) {
	got, err := filterFilesByFilesets(
		[]string{"internal/cli/root.go", "internal/authoring/demux.go"},
		[]string{"./internal/cli/root.go"},
	)
	if err != nil {
		t.Fatalf("filterFilesByFilesets() error = %v", err)
	}
	want := []string{"internal/cli/root.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterFilesByFilesets() = %#v, want %#v", got, want)
	}
}

func TestFilterFilesByFilesetsMatchesDirectoryPrefix(t *testing.T) {
	got, err := filterFilesByFilesets(
		[]string{"internal/cli/root.go", "internal/authoring/demux.go", "README.md"},
		[]string{"internal/cli"},
	)
	if err != nil {
		t.Fatalf("filterFilesByFilesets() error = %v", err)
	}
	want := []string{"internal/cli/root.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterFilesByFilesets() = %#v, want %#v", got, want)
	}
}

func TestFilterFilesByFilesetsReturnsNoMatchError(t *testing.T) {
	_, err := filterFilesByFilesets([]string{"internal/cli/root.go"}, []string{"internal/storage"})
	if err == nil || !strings.Contains(err.Error(), "no working-copy changes matched demux filesets") {
		t.Fatalf("filterFilesByFilesets() error = %v, want no-match error", err)
	}
}

func TestExcludeFilesByFilesetsMatchesDirectoryPrefix(t *testing.T) {
	got, err := excludeFilesByFilesets(
		[]string{"internal/cli/root.go", "internal/authoring/demux.go", "README.md"},
		[]string{"internal/authoring"},
	)
	if err != nil {
		t.Fatalf("excludeFilesByFilesets() error = %v", err)
	}
	want := []string{"internal/cli/root.go", "README.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("excludeFilesByFilesets() = %#v, want %#v", got, want)
	}
}

func TestExcludeFilesByFilesetsReturnsAllExcludedError(t *testing.T) {
	_, err := excludeFilesByFilesets([]string{"internal/cli/root.go"}, []string{"internal"})
	if err == nil || !strings.Contains(err.Error(), "all working-copy changes were excluded from demux") {
		t.Fatalf("excludeFilesByFilesets() error = %v, want all-excluded error", err)
	}
}

func TestFilterHunksByFiles(t *testing.T) {
	got := filterHunksByFiles(
		[]HunkRange{
			{ID: "h1", File: "internal/cli/root.go"},
			{ID: "h2", File: "internal/authoring/demux.go"},
		},
		[]string{"internal/cli/root.go"},
	)
	want := []HunkRange{{ID: "h1", File: "internal/cli/root.go"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterHunksByFiles() = %#v, want %#v", got, want)
	}
}

func TestOrderGroupsByStructuralDependencies(t *testing.T) {
	groups := [][]string{{"internal/app/app.go"}, {"internal/core/helper.go"}}
	got := orderGroupsByStructuralDependencies(groups, structural.Facts{
		Edges: []structural.DependencyEdge{{
			FromFile: "internal/app/app.go",
			ToFile:   "internal/core/helper.go",
			Symbol:   "NewThing",
		}},
	})
	want := [][]string{{"internal/core/helper.go"}, {"internal/app/app.go"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("orderGroupsByStructuralDependencies() = %#v, want %#v", got, want)
	}
}

func TestAnnotateHunksWithSymbols(t *testing.T) {
	got := annotateHunksWithSymbols(
		[]HunkRange{{ID: "h1", File: "app.go"}},
		[]structural.HunkSymbol{{HunkID: "h1", File: "app.go", Symbol: "Run", Kind: "function"}},
	)
	if got[0].Symbol != "Run" || got[0].SymbolKind != "function" {
		t.Fatalf("annotateHunksWithSymbols() = %#v", got)
	}
}

func TestExpandGroupsByChangedSymbolsSplitsSingleFileBySymbol(t *testing.T) {
	got := expandGroupsByChangedSymbols([][]string{{"app.go"}}, map[string][]HunkRange{
		"app.go": {
			{ID: "h1", File: "app.go", Symbol: "First", SymbolKind: "function"},
			{ID: "h2", File: "app.go", Symbol: "Second", SymbolKind: "function"},
		},
	})
	if len(got) != 2 {
		t.Fatalf("expandGroupsByChangedSymbols() = %#v, want 2 groups", got)
	}
	if !got[0].UseHunks || got[0].Symbol != "First" || !reflect.DeepEqual(got[0].Hunks, []HunkRange{{ID: "h1", File: "app.go", Symbol: "First", SymbolKind: "function"}}) {
		t.Fatalf("first group = %#v", got[0])
	}
	if !got[1].UseHunks || got[1].Symbol != "Second" || !reflect.DeepEqual(got[1].Hunks, []HunkRange{{ID: "h2", File: "app.go", Symbol: "Second", SymbolKind: "function"}}) {
		t.Fatalf("second group = %#v", got[1])
	}
}

func TestExpandGroupsByChangedSymbolsFallsBackWhenAnyHunkHasNoSymbol(t *testing.T) {
	got := expandGroupsByChangedSymbols([][]string{{"app.go"}}, map[string][]HunkRange{
		"app.go": {
			{ID: "h1", File: "app.go", Symbol: "First", SymbolKind: "function"},
			{ID: "h2", File: "app.go"},
		},
	})
	if len(got) != 1 || got[0].UseHunks || !reflect.DeepEqual(got[0].Files, []string{"app.go"}) || len(got[0].Hunks) != 2 {
		t.Fatalf("expandGroupsByChangedSymbols() = %#v, want file-level fallback", got)
	}
}

func TestNormalizeDemuxProposalResolvesHunkIDs(t *testing.T) {
	proposal, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.txt", Patch: "patch one"},
			{ID: "h2", File: "beta.txt", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "beta", UseHunks: true, HunkIDs: []string{"h2"}},
		},
	})
	if err != nil {
		t.Fatalf("normalizeDemuxProposal() error = %v", err)
	}
	if !reflect.DeepEqual(proposal.Revisions[0].Files, []string{"alpha.txt"}) {
		t.Fatalf("first revision files = %#v", proposal.Revisions[0].Files)
	}
	if got := proposal.Revisions[0].Hunks; len(got) != 1 || got[0].Patch != "patch one" {
		t.Fatalf("first revision hunks = %#v", got)
	}
}

func TestNormalizeDemuxProposalRejectsDuplicateHunkIDs(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "beta", UseHunks: true, HunkIDs: []string{"h1"}},
		},
	})
	if err == nil {
		t.Fatal("normalizeDemuxProposal() error = nil, want duplicate hunk error")
	}
}

func TestApplyRoutedDemuxRevisionRestoresFilesFromProposalCommit(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}

	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(context.Background(), store, RepoInfo{RootPath: repoRoot, Backend: "jj"}, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	targetHead := "target-head"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "internal hooks",
		BookmarkName: "feature/internal-hooks",
		BaseRef:      "main",
		BaseCommitID: "base-commit",
		HeadChangeID: &targetHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)

	runner := &demuxRoutingFakeRunner{repoRoot: repoRoot}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))
	revision := RevisionProposal{
		ID:          "u1",
		Intent:      "restore hooks",
		Files:       []string{"internal/hooks/run.go", "internal/hooks/run_test.go"},
		TargetStack: "feature/internal-hooks",
	}

	if _, err := engine.applyRoutedDemuxRevision(context.Background(), DemuxProposal{
		RepoRoot:         repoRoot,
		ProposedCommitID: "proposal-commit",
		Revisions:        []RevisionProposal{revision},
	}, revision); err != nil {
		t.Fatalf("applyRoutedDemuxRevision() error = %v", err)
	}

	if runner.called("git apply") {
		t.Fatalf("routed demux used git apply: %#v", runner.calls)
	}
	for _, want := range []string{
		"jj restore --from proposal-commit internal/hooks/run.go internal/hooks/run_test.go",
		"jj restore --from @- internal/hooks/run.go internal/hooks/run_test.go",
	} {
		if !runner.called(want) {
			t.Fatalf("missing call %q in %#v", want, runner.calls)
		}
	}
}

func TestWipePendingDemuxProposalsPreservesAnchorOnSameChange(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	engine := NewEngine()
	repo := RepoInfo{RootPath: "/repo", Backend: "jj"}

	store, err := openStore(ctx)
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(ctx, store, repo, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	if err := store.UpsertDemuxProposal(ctx, storage.DemuxProposal{
		ID:           "demux-existing",
		RepoID:       repoID,
		BaseChangeID: "change-current",
		Status:       "pending",
		PayloadJSON:  `{}`,
		CreatedAt:    123,
		UpdatedAt:    123,
	}); err != nil {
		t.Fatalf("UpsertDemuxProposal() error = %v", err)
	}
	_ = store.Close()

	anchor, err := engine.wipePendingDemuxProposals(ctx, repo, "change-current")
	if err != nil {
		t.Fatalf("wipePendingDemuxProposals() error = %v", err)
	}
	if anchor.ID != "demux-existing" || anchor.CreatedAt != 123 {
		t.Fatalf("anchor = %#v, want existing proposal identity", anchor)
	}

	store, err = openStore(ctx)
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	defer store.Close()
	list, err := store.ListDemuxProposals(ctx, repoID, "pending", 10)
	if err != nil {
		t.Fatalf("ListDemuxProposals() error = %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("pending proposals after wipe = %#v, want none", list)
	}
}

func TestWipePendingDemuxProposalsDoesNotAnchorAcrossChanges(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	engine := NewEngine()
	repo := RepoInfo{RootPath: "/repo", Backend: "jj"}

	store, err := openStore(ctx)
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	repoID, err := upsertProposalRepo(ctx, store, repo, 1)
	if err != nil {
		t.Fatalf("upsertProposalRepo() error = %v", err)
	}
	if err := store.UpsertDemuxProposal(ctx, storage.DemuxProposal{
		ID:           "demux-existing",
		RepoID:       repoID,
		BaseChangeID: "change-old",
		Status:       "pending",
		PayloadJSON:  `{}`,
		CreatedAt:    123,
		UpdatedAt:    123,
	}); err != nil {
		t.Fatalf("UpsertDemuxProposal() error = %v", err)
	}
	_ = store.Close()

	anchor, err := engine.wipePendingDemuxProposals(ctx, repo, "change-new")
	if err != nil {
		t.Fatalf("wipePendingDemuxProposals() error = %v", err)
	}
	if anchor.ID != "" || anchor.CreatedAt != 0 {
		t.Fatalf("anchor = %#v, want empty across changes", anchor)
	}
}

func TestReviewDemuxPlanReturnsRepairableErrors(t *testing.T) {
	engine := NewEngine()
	result, err := engine.ReviewDemuxPlan(context.Background(), DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.txt", Patch: "patch one"},
			{ID: "h2", File: "alpha.txt", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha first", UseHunks: true, HunkIDs: []string{"h1"}},
		},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}
	if result.Valid || len(result.Errors) != 1 || !strings.Contains(result.Errors[0], "hunk h2 in alpha.txt is not assigned") {
		t.Fatalf("ReviewDemuxPlan() = %#v, want repairable missing hunk error", result)
	}
	if len(result.RepairHints) != 1 || result.RepairHints[0].Kind != "unassigned_hunk" || result.RepairHints[0].HunkID != "h2" {
		t.Fatalf("ReviewDemuxPlan() repair hints = %#v, want unassigned hunk h2", result.RepairHints)
	}
}

func TestDemuxPipelinePacketOwnsWorkflowState(t *testing.T) {
	engine := NewEngine()
	packet, err := engine.demuxPipeline().packetForProposal(context.Background(), DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
		}},
	})
	if err != nil {
		t.Fatalf("packetForProposal() error = %v", err)
	}
	if packet.State != DemuxWorkflowReadyToApply {
		t.Fatalf("packetForProposal() state = %q, want ready_to_apply", packet.State)
	}
	if packet.NextTool != "gx_apply_revision_plan" || packet.FinalTool != "gx_apply_revision_plan" {
		t.Fatalf("packetForProposal() tools = %q/%q", packet.NextTool, packet.FinalTool)
	}
	if !packet.Review.Valid {
		t.Fatalf("packetForProposal() review valid = false: %#v", packet.Review)
	}
}

func TestAppendDemuxApplyPreflightWarningBlocksApplyWithoutInvalidatingProposal(t *testing.T) {
	proposal := appendDemuxApplyPreflightWarning(DemuxProposal{
		Warnings:            []string{"keep this warning", "Compose apply preflight failed: old failure"},
		FeasibilityWarnings: []FeasibilityWarning{{Severity: "info", Source: "inferred_dependency", Message: "keep info"}},
	}, errors.New("apply exploded"))

	if got := demuxWorkflowState(ReviewDemuxResult{Valid: true, Proposal: proposal}, proposal); got != DemuxWorkflowRepairRequired {
		t.Fatalf("demuxWorkflowState() = %q, want repair_required", got)
	}
	if len(proposal.FeasibilityWarnings) != 2 {
		t.Fatalf("feasibility warnings = %#v, want info plus apply_preflight warning", proposal.FeasibilityWarnings)
	}
	warning := proposal.FeasibilityWarnings[1]
	if warning.Source != "apply_preflight" || warning.Severity != "warning" || !strings.Contains(warning.Message, "apply exploded") {
		t.Fatalf("apply preflight warning = %#v", warning)
	}
	if len(proposal.Warnings) != 2 || proposal.Warnings[0] != "keep this warning" || !strings.Contains(proposal.Warnings[1], "apply exploded") {
		t.Fatalf("plain warnings = %#v", proposal.Warnings)
	}
}

func TestDemuxPreflightSkipsGeneratedDirsUnlessProposalTouchesThem(t *testing.T) {
	if !demuxPreflightSkippableDir("node_modules") || !demuxPreflightSkippableDir("frontend/.next") {
		t.Fatal("expected generated dependency/build dirs to be skippable")
	}
	protected := demuxPreflightProtectedPaths(DemuxProposal{
		Revisions: []RevisionProposal{{
			Files: []string{"node_modules/local-patch/index.js", ".next/required-manifest.json"},
		}},
	})
	if !hasProtectedPathUnder(protected, "node_modules") {
		t.Fatalf("node_modules should stay protected when proposal touches it: %#v", protected)
	}
	if !hasProtectedPathUnder(protected, ".next") {
		t.Fatalf(".next should stay protected when proposal touches it: %#v", protected)
	}
	if hasProtectedPathUnder(protected, "coverage") {
		t.Fatalf("coverage should be skippable when untouched: %#v", protected)
	}
}

func TestCopyTreeWithSkipsOmitsIgnoredFilesUnlessProtected(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	mustWriteFile(t, filepath.Join(src, "keep.txt"), "keep")
	mustWriteFile(t, filepath.Join(src, "src-tauri/gen/schemas/desktop-schema.json"), "{}")
	mustWriteFile(t, filepath.Join(src, ".jj/repo/store/type"), "git")

	ignored := map[string]struct{}{
		filepath.Clean("src-tauri/gen/schemas/desktop-schema.json"): {},
		filepath.Clean(".jj/repo/store/type"):                       {},
	}
	skipIgnored := func(rel string, entry fs.DirEntry) bool {
		return demuxPreflightIgnoredPath(rel, ignored)
	}
	if err := copyTreeWithSkips(src, dst, nil, skipIgnored); err != nil {
		t.Fatalf("copyTreeWithSkips() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "keep.txt")); err != nil {
		t.Fatalf("kept file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "src-tauri/gen/schemas/desktop-schema.json")); !os.IsNotExist(err) {
		t.Fatalf("ignored generated file should be omitted, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, ".jj/repo/store/type")); err != nil {
		t.Fatalf("jj internals should be preserved even when git reports them ignored: %v", err)
	}

	protectedDst := t.TempDir()
	protected := map[string]struct{}{
		filepath.Clean("src-tauri/gen/schemas/desktop-schema.json"): {},
	}
	if err := copyTreeWithSkips(src, protectedDst, protected, skipIgnored); err != nil {
		t.Fatalf("copyTreeWithSkips(protected) error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(protectedDst, "src-tauri/gen/schemas/desktop-schema.json")); err != nil {
		t.Fatalf("protected ignored file should be copied: %v", err)
	}
}

func TestDemuxTargetBlockingFilesIgnoresSourceIgnoredGeneratedOutput(t *testing.T) {
	ignored := map[string]struct{}{
		filepath.Clean("node_modules/react/index.js"): {},
		filepath.Clean("dist/index.html"):             {},
	}
	got := demuxTargetBlockingFiles(
		[]string{"node_modules/react/index.js", "dist/index.html", "package.json", "README.md"},
		[]string{"package.json"},
		ignored,
	)
	want := []string{"package.json", "README.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("demuxTargetBlockingFiles() = %#v, want %#v", got, want)
	}
}

func mustWriteFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestReviewDemuxPlanReturnsNormalizedWarnings(t *testing.T) {
	engine := NewEngine()
	result, err := engine.ReviewDemuxPlan(context.Background(), DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "helper.go", Patch: "patch helper", Symbol: "NewThing"},
			{ID: "h2", File: "app.go", Patch: "patch app", Symbol: "Run"},
		},
		StructuralDeps: []StructuralDependency{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "helper", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "app", UseHunks: true, HunkIDs: []string{"h2"}},
		},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}
	if !result.Valid {
		t.Fatalf("ReviewDemuxPlan() valid = false: %#v", result)
	}
	if len(result.Proposal.Revisions[0].Hunks) != 1 || result.Proposal.Revisions[0].Files[0] != "helper.go" {
		t.Fatalf("ReviewDemuxPlan() did not normalize hunk ids: %#v", result.Proposal.Revisions[0])
	}
	if len(result.Proposal.FeasibilityWarnings) != 1 || result.Proposal.FeasibilityWarnings[0].Source != "inferred_dependency" {
		t.Fatalf("ReviewDemuxPlan() warnings = %#v, want inferred dependency", result.Proposal.FeasibilityWarnings)
	}
	if len(result.RepairHints) != 0 {
		t.Fatalf("ReviewDemuxPlan() repair hints = %#v, want none for advisory dependency warning", result.RepairHints)
	}
}

func TestReviewDemuxPlanKeepsStructuralDependencyAdvisoryNonBlocking(t *testing.T) {
	engine := NewEngine()
	result, err := engine.ReviewDemuxPlan(context.Background(), DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "patch app", Symbol: "Run"},
			{ID: "h2", File: "helper.go", Patch: "patch helper", Symbol: "NewThing"},
		},
		StructuralDeps: []StructuralDependency{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "app", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "helper", UseHunks: true, HunkIDs: []string{"h2"}},
		},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}
	if !result.Valid {
		t.Fatalf("ReviewDemuxPlan() valid = false: %#v", result)
	}
	if len(result.Proposal.FeasibilityWarnings) != 1 || result.Proposal.FeasibilityWarnings[0].Source != "structural_dependency" {
		t.Fatalf("ReviewDemuxPlan() warnings = %#v, want structural dependency", result.Proposal.FeasibilityWarnings)
	}
	warning := result.Proposal.FeasibilityWarnings[0]
	if warning.RevisionID != "r1" || warning.DependsOn != "r2" || warning.Symbol != "NewThing" {
		t.Fatalf("ReviewDemuxPlan() warning = %#v, want typed dependency fields", warning)
	}
	if len(result.RepairHints) != 0 {
		t.Fatalf("ReviewDemuxPlan() repair hints = %#v, want none for advisory dependency warning", result.RepairHints)
	}
}

func TestShapeDemuxProposalCoalescesTinySameFileRevisions(t *testing.T) {
	proposal, changed := shapeDemuxProposalWithPolicy(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "@@\n+func First() {}\n", Symbol: "First"},
			{ID: "h2", File: "app.go", Patch: "@@\n+func Second() {}\n", Symbol: "Second"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "first", Files: []string{"app.go"}, UseHunks: true, HunkIDs: []string{"h1"}, Hunks: []HunkRange{{ID: "h1", File: "app.go", Patch: "@@\n+func First() {}\n", Symbol: "First"}}},
			{ID: "u2", Intent: "second", Files: []string{"app.go"}, UseHunks: true, HunkIDs: []string{"h2"}, Hunks: []HunkRange{{ID: "h2", File: "app.go", Patch: "@@\n+func Second() {}\n", Symbol: "Second"}}},
		},
	}, defaultDemuxReviewShapePolicy)
	if !changed {
		t.Fatalf("shapeDemuxProposalWithPolicy() changed = false, want true")
	}
	if len(proposal.Revisions) != 1 {
		t.Fatalf("shaped revisions = %#v, want one coalesced revision", proposal.Revisions)
	}
	revision := proposal.Revisions[0]
	if !reflect.DeepEqual(revision.HunkIDs, []string{"h1", "h2"}) {
		t.Fatalf("coalesced revision = %#v, want h1/h2", revision)
	}
	if revision.EffectiveLOC != 2 {
		t.Fatalf("coalesced effective LOC = %d, want 2", revision.EffectiveLOC)
	}
	if !strings.Contains(strings.Join(revision.ShapeReasons, "\n"), "merged tiny revision") {
		t.Fatalf("coalesced shape reasons = %#v, want merge reason", revision.ShapeReasons)
	}
}

func TestShapeDemuxProposalKeepsTinyStandaloneRevisionExplainable(t *testing.T) {
	proposal, changed := shapeDemuxProposalWithPolicy(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "config.json", Patch: "@@\n+{\"enabled\": true}\n"},
			{ID: "h2", File: "app.go", Patch: "@@\n+func Run() {}\n"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "enable config", Files: []string{"config.json"}, UseHunks: true, HunkIDs: []string{"h1"}, Hunks: []HunkRange{{ID: "h1", File: "config.json", Patch: "@@\n+{\"enabled\": true}\n"}}},
			{ID: "u2", Intent: "add run", Files: []string{"app.go"}, UseHunks: true, HunkIDs: []string{"h2"}, Hunks: []HunkRange{{ID: "h2", File: "app.go", Patch: "@@\n+func Run() {}\n"}}},
		},
	}, defaultDemuxReviewShapePolicy)
	if !changed {
		t.Fatalf("shapeDemuxProposalWithPolicy() changed = false, want true shape annotation")
	}
	if len(proposal.Revisions) != 2 {
		t.Fatalf("shaped revisions = %#v, want tiny unrelated revisions kept separate", proposal.Revisions)
	}
	if !strings.Contains(strings.Join(proposal.Revisions[0].ShapeReasons, "\n"), "tiny revision") {
		t.Fatalf("first shape reasons = %#v, want tiny revision explanation", proposal.Revisions[0].ShapeReasons)
	}
}

func TestShapeDemuxProposalNormalizesRoutesForSamePackageCluster(t *testing.T) {
	proposal, changed := shapeDemuxProposalWithPolicy(DemuxProposal{
		StructuralDeps: []StructuralDependency{
			{FromFile: "internal/authoring/demux.go", ToFile: "internal/authoring/demux_review_shape.go", Symbol: "shapeDemuxProposalForReview"},
			{FromFile: "internal/authoring/demux_ai.go", ToFile: "internal/authoring/demux_review_shape.go", Symbol: "shapeDemuxProposalForReview"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "wire compose", Files: []string{"internal/authoring/demux.go"}, TargetStack: "feature/demux-routing", EffectiveLOC: 12},
			{ID: "u2", Intent: "wire repair", Files: []string{"internal/authoring/demux_ai.go"}, TargetStack: "feature/internal-authoring", EffectiveLOC: 12},
			{ID: "u3", Intent: "add shape module", Files: []string{"internal/authoring/demux_review_shape.go"}, TargetStack: "feature/internal-authoring", EffectiveLOC: 200},
		},
	}, defaultDemuxReviewShapePolicy)
	if !changed {
		t.Fatalf("shapeDemuxProposalWithPolicy() changed = false, want route normalization")
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != "feature/internal-authoring" {
			t.Fatalf("revision %s target stack = %q, want feature/internal-authoring in %#v", revision.ID, revision.TargetStack, proposal.Revisions)
		}
	}
}

func TestShapeDemuxProposalCoalescesRelatedPackageRevisionsUnderSoftMax(t *testing.T) {
	proposal, changed := shapeDemuxProposalWithPolicy(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "internal/authoring/proposal.go", Patch: "@@\n+EffectiveLOC int\n+ShapeReasons []string\n"},
			{ID: "h2", File: "internal/authoring/demux.go", Patch: "@@\n+proposal, _ = shapeDemuxProposalForReview(proposal)\n"},
			{ID: "h3", File: "internal/authoring/demux_ai.go", Patch: "@@\n+proposal, shaped := shapeDemuxProposalForReview(proposal)\n"},
			{ID: "h4", File: "internal/authoring/demux_test.go", Patch: "@@\n+func TestShapeDemuxProposal() {}\n"},
			{ID: "h5", File: "internal/authoring/demux_review_shape.go", Patch: strings.Repeat("+func helper() {}\n", 200)},
		},
		StructuralDeps: []StructuralDependency{
			{FromFile: "internal/authoring/demux.go", ToFile: "internal/authoring/demux_review_shape.go", Symbol: "shapeDemuxProposalForReview"},
			{FromFile: "internal/authoring/demux_ai.go", ToFile: "internal/authoring/demux_review_shape.go", Symbol: "shapeDemuxProposalForReview"},
			{FromFile: "internal/authoring/demux_test.go", ToFile: "internal/authoring/demux.go", Symbol: "ReviewDemuxPlan"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "add proposal shape fields", Files: []string{"internal/authoring/proposal.go"}, UseHunks: true, HunkIDs: []string{"h1"}, TargetStack: "feature/internal-authoring"},
			{ID: "u2", Intent: "wire review shape compose", Files: []string{"internal/authoring/demux.go"}, UseHunks: true, HunkIDs: []string{"h2"}, TargetStack: "feature/internal-authoring"},
			{ID: "u3", Intent: "wire review shape repair", Files: []string{"internal/authoring/demux_ai.go"}, UseHunks: true, HunkIDs: []string{"h3"}, TargetStack: "feature/internal-authoring"},
			{ID: "u4", Intent: "test review shape", Files: []string{"internal/authoring/demux_test.go"}, UseHunks: true, HunkIDs: []string{"h4"}, TargetStack: "feature/internal-authoring"},
			{ID: "u5", Intent: "add review shape module", Files: []string{"internal/authoring/demux_review_shape.go"}, UseHunks: true, HunkIDs: []string{"h5"}, TargetStack: "feature/internal-authoring"},
		},
	}, defaultDemuxReviewShapePolicy)
	if !changed {
		t.Fatalf("shapeDemuxProposalWithPolicy() changed = false, want semantic coalescing")
	}
	if len(proposal.Revisions) != 2 {
		t.Fatalf("shaped revisions = %#v, want wiring/test cluster plus large module", proposal.Revisions)
	}
	sawCluster := false
	sawLargeModule := false
	for _, revision := range proposal.Revisions {
		switch {
		case len(revision.HunkIDs) == 4:
			sawCluster = true
			if !strings.Contains(strings.Join(revision.ShapeReasons, "\n"), "merged related revision") {
				t.Fatalf("cluster shape reasons = %#v, want semantic merge reason", revision.ShapeReasons)
			}
		case reflect.DeepEqual(revision.HunkIDs, []string{"h5"}):
			sawLargeModule = true
		}
	}
	if !sawCluster || !sawLargeModule {
		t.Fatalf("shaped revisions = %#v, want one cluster and one large module", proposal.Revisions)
	}
}

func TestDemuxStackClusterNamesScaffoldRoutes(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		wantKey  string
		wantName string
		wantKind string
	}{
		{
			name:     "tooling",
			files:    []string{".gitignore", "package-lock.json", "package.json", "tsconfig.json", "vite.config.ts"},
			wantKey:  "project-tooling",
			wantName: "project tooling",
			wantKind: "chore",
		},
		{
			name:     "tauri",
			files:    []string{"src-tauri/Cargo.toml", "src-tauri/src/lib.rs", "src-tauri/tauri.conf.json"},
			wantKey:  "tauri-app",
			wantName: "Tauri app",
			wantKind: "feature",
		},
		{
			name:     "frontend",
			files:    []string{"index.html", "src/App.tsx", "src/main.tsx"},
			wantKey:  "frontend-app",
			wantName: "frontend app",
			wantKind: "feature",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, name, kind, _, _ := demuxStackClusterForRevision(RevisionProposal{
				Intent: "stand up Tauri app shell",
				Files:  tt.files,
			})
			if key != tt.wantKey || name != tt.wantName || kind != tt.wantKind {
				t.Fatalf("demuxStackClusterForRevision() = (%q, %q, %q), want (%q, %q, %q)", key, name, kind, tt.wantKey, tt.wantName, tt.wantKind)
			}
		})
	}
}

func TestScaffoldDemuxProposalGroupsBeforeRouting(t *testing.T) {
	files := []string{
		".gitignore",
		".vscode/extensions.json",
		"README.md",
		"index.html",
		"package-lock.json",
		"package.json",
		"src/App.css",
		"src/App.tsx",
		"src/main.tsx",
		"src/vite-env.d.ts",
		"src-tauri/.gitignore",
		"src-tauri/Cargo.lock",
		"src-tauri/Cargo.toml",
		"src-tauri/build.rs",
		"src-tauri/capabilities/default.json",
		"src-tauri/src/lib.rs",
		"src-tauri/src/main.rs",
		"src-tauri/tauri.conf.json",
		"tsconfig.json",
		"tsconfig.node.json",
		"vite.config.ts",
	}
	var hunks []HunkRange
	for index, file := range files {
		hunks = append(hunks, HunkRange{
			ID:    fmt.Sprintf("h%d", index+1),
			File:  file,
			Patch: "@@\n+" + file + "\n",
		})
	}
	groups := expandGroupsByChangedSymbols(groupFiles(files), hunksByFile(hunks))
	revisions := make([]RevisionProposal, 0, len(groups))
	for index, group := range groups {
		revisions = append(revisions, RevisionProposal{
			ID:       fmt.Sprintf("u%d", index+1),
			Intent:   proposalIntentForGroup("Stand up Tauri app shell", group),
			Files:    group.Files,
			HunkIDs:  hunkIDs(group.Hunks),
			Hunks:    group.Hunks,
			UseHunks: len(group.Hunks) > 0,
		})
	}
	proposal, _ := shapeDemuxProposalForReview(DemuxProposal{
		Hunks:     hunks,
		Revisions: revisions,
	})
	proposal = inferNewStackRoutes(proposal)

	if len(proposal.Revisions) != 4 {
		t.Fatalf("revisions = %#v, want 4 scaffold groups", proposal.Revisions)
	}
	gotRoutes := make([]string, 0, len(proposal.Revisions))
	for _, revision := range proposal.Revisions {
		gotRoutes = append(gotRoutes, revision.TargetStack)
	}
	wantRoutes := []string{"chore/project-tooling", "feature/tauri-app", "feature/frontend-app", "docs/documentation"}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Fatalf("routes = %#v, want %#v", gotRoutes, wantRoutes)
	}
}

func TestDemuxWorkflowStateOnlyRequiresBlockingFailures(t *testing.T) {
	ready := demuxWorkflowState(ReviewDemuxResult{Valid: true}, DemuxProposal{})
	if ready != DemuxWorkflowReadyToApply {
		t.Fatalf("ready state = %q, want %q", ready, DemuxWorkflowReadyToApply)
	}

	hintsOnly := demuxWorkflowState(ReviewDemuxResult{
		Valid:       true,
		RepairHints: []RepairHint{{Kind: "inferred_dependency"}},
	}, DemuxProposal{})
	if hintsOnly != DemuxWorkflowReadyToApply {
		t.Fatalf("hints-only state = %q, want %q", hintsOnly, DemuxWorkflowReadyToApply)
	}

	required := demuxWorkflowState(ReviewDemuxResult{Valid: true}, DemuxProposal{
		FeasibilityWarnings: []FeasibilityWarning{{Severity: "warning", Source: "structural_dependency"}},
	})
	if required != DemuxWorkflowRepairRequired {
		t.Fatalf("required state = %q, want %q", required, DemuxWorkflowRepairRequired)
	}
}

func TestDemuxWorkflowPacketUsesReviewedProposal(t *testing.T) {
	proposal := DemuxProposal{
		ID:               "demux-1",
		RepoRoot:         "/repo",
		ProposedChangeID: "change-1",
		ProposedCommitID: "commit-1",
		Status:           ProposalPending,
		Hunks:            []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch"}},
		Revisions:        []RevisionProposal{{ID: "r1", Intent: "alpha", UseHunks: true, HunkIDs: []string{"h1"}}},
	}
	review := ReviewDemuxResult{Valid: true, Proposal: proposal}
	state := demuxWorkflowState(review, review.Proposal)
	if state != DemuxWorkflowReadyToApply {
		t.Fatalf("demuxWorkflowState() = %q, want ready", state)
	}
	shape := demuxToolArgumentShape(review.Proposal)
	rawProposal, ok := shape["proposal"].(map[string]any)
	if !ok {
		t.Fatalf("demuxToolArgumentShape() = %#v, want proposal map", shape)
	}
	if rawProposal["id"] != "demux-1" || rawProposal["proposed_change_id"] != "change-1" {
		t.Fatalf("demuxToolArgumentShape() proposal = %#v", rawProposal)
	}
}

func TestOpenAIDemuxReviewerRequestsJSONProposal(t *testing.T) {
	var got struct {
		Model               string            `json:"model"`
		Messages            []map[string]any  `json:"messages"`
		ResponseFormat      map[string]string `json:"response_format"`
		MaxCompletionTokens int               `json:"max_completion_tokens"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		content := `{"proposal":{"id":"demux-1","repo_root":"/repo","proposed_change_id":"change-1","proposed_commit_id":"commit-1","status":"pending","revisions":[{"id":"u1","intent":"alpha","files":["alpha.go"],"provenance_status":"absent","confidence":0.8}]},"notes":["kept single revision"]}`
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "test-model",
			"choices": []map[string]any{{
				"message": map[string]any{"role": "assistant", "content": content},
			}},
		})
	}))
	defer server.Close()

	reviewer := newTestOpenAIDemuxReviewer(server, "test-model")
	result, err := reviewer.ReviewDemuxProposal(context.Background(), demuxAIReviewRequest{
		Proposal: demuxAIProposalForReview{
			ID:               "demux-1",
			RepoRoot:         "/repo",
			ProposedChangeID: "change-1",
			ProposedCommitID: "commit-1",
			Status:           ProposalPending,
			Hunks:            []HunkRange{{ID: "h1", File: "alpha.go"}},
			Revisions:        []RevisionProposal{{ID: "u1", Intent: "alpha", Files: []string{"alpha.go"}}},
		},
		Review: demuxAIReviewSummary{Valid: true},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxProposal() error = %v", err)
	}
	if got.Model != "test-model" || got.ResponseFormat["type"] != "json_object" || got.MaxCompletionTokens != defaultDemuxMaxOutputTokens || len(got.Messages) != 2 {
		t.Fatalf("request = %#v", got)
	}
	if got.Messages[0]["role"] != "developer" {
		t.Fatalf("first message role = %q, want developer", got.Messages[0]["role"])
	}
	if result.Proposal.ID != "demux-1" || len(result.Proposal.Revisions) != 1 || result.Notes[0] != "kept single revision" {
		t.Fatalf("result = %#v", result)
	}
}

func TestNormalizeOpenAIBaseURLAddsSingleV1(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"http://127.0.0.1:43123", "http://127.0.0.1:43123/v1"},
		{"http://127.0.0.1:43123/v1", "http://127.0.0.1:43123/v1"},
		{"http://127.0.0.1:43123/v1/", "http://127.0.0.1:43123/v1"},
		{"", "https://api.openai.com/v1"},
	}
	for _, tc := range tests {
		if got := normalizeOpenAIBaseURL(tc.in); got != tc.want {
			t.Fatalf("normalizeOpenAIBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDemuxReviewerFromEnvUsesOpenAIBaseURLFallback(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "off")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("GX_OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_BASE_URL", "http://127.0.0.1:43123/v1")
	t.Setenv("GX_DEMUX_REVIEW_MODEL", "test-model")

	reviewer, model, err := demuxReviewerFromEnv(DemuxAIReviewOptions{})
	if err != nil {
		t.Fatalf("demuxReviewerFromEnv() error = %v", err)
	}
	got, ok := reviewer.(fallbackDemuxReviewer)
	if !ok {
		t.Fatalf("reviewer = %T, want fallbackDemuxReviewer", reviewer)
	}
	primary, ok := got.primary.(*openAIDemuxReviewer)
	if !ok {
		t.Fatalf("primary reviewer = %T, want *openAIDemuxReviewer", got.primary)
	}
	fallback, ok := got.fallback.(*openAIDemuxReviewer)
	if !ok {
		t.Fatalf("fallback reviewer = %T, want *openAIDemuxReviewer", got.fallback)
	}
	if primary.baseURL != "http://127.0.0.1:43123/v1" || primary.apiKey != "test-key" || model != "test-model" {
		t.Fatalf("primary reviewer config = baseURL %q apiKey %q model %q", primary.baseURL, primary.apiKey, model)
	}
	if fallback.baseURL != "https://api.openai.com/v1" || fallback.apiKey != "test-key" {
		t.Fatalf("fallback reviewer config = baseURL %q apiKey %q", fallback.baseURL, fallback.apiKey)
	}
}

func TestDemuxReviewerFromEnvPrefersUserOpenAIKey(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "off")
	t.Setenv("GX_OPENAI_API_KEY", "gx-key")
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("GX_OPENAI_BASE_URL", "http://127.0.0.1:43124")
	t.Setenv("OPENAI_BASE_URL", "http://127.0.0.1:43123")
	t.Setenv("GX_DEMUX_REVIEW_MODEL", "test-model")

	reviewer, model, err := demuxReviewerFromEnv(DemuxAIReviewOptions{})
	if err != nil {
		t.Fatalf("demuxReviewerFromEnv() error = %v", err)
	}
	fallback, ok := reviewer.(fallbackDemuxReviewer)
	if !ok {
		t.Fatalf("reviewer = %T, want fallbackDemuxReviewer", reviewer)
	}
	got, ok := fallback.primary.(*openAIDemuxReviewer)
	if !ok {
		t.Fatalf("primary reviewer = %T, want *openAIDemuxReviewer", fallback.primary)
	}
	if got.baseURL != "http://127.0.0.1:43123/v1" || got.apiKey != "openai-key" || model != "test-model" {
		t.Fatalf("reviewer config = baseURL %q apiKey %q model %q", got.baseURL, got.apiKey, model)
	}
}

func TestDemuxReviewerFromEnvUsesGXCloudAsFallbackAfterUserKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"model":"test-model","choices":[{"message":{"content":"{\"revisions\":[]}"}}]}`)
	}))
	defer server.Close()

	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_CLOUD_URL", server.URL+"/gx/pr")
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("OPENAI_BASE_URL", "http://127.0.0.1:43123")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_BASE_URL", "")
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_demux",
		CLISessionToken:   "gxcs_demux",
		UserID:            "demux-user",
		Login:             "demux",
		MachineID:         "demux-machine",
		ObtainedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	reviewer, _, err := demuxReviewerFromEnv(DemuxAIReviewOptions{Model: "test-model"})
	if err != nil {
		t.Fatalf("demuxReviewerFromEnv() error = %v", err)
	}
	got, ok := reviewer.(fallbackDemuxReviewer)
	if !ok {
		t.Fatalf("reviewer = %T, want fallbackDemuxReviewer", reviewer)
	}
	directFallback, ok := got.primary.(fallbackDemuxReviewer)
	if !ok {
		t.Fatalf("primary reviewer = %T, want fallbackDemuxReviewer", got.primary)
	}
	primary, ok := directFallback.primary.(*openAIDemuxReviewer)
	if !ok {
		t.Fatalf("direct primary reviewer = %T, want *openAIDemuxReviewer", directFallback.primary)
	}
	if primary.apiKey != "openai-key" || primary.baseURL != "http://127.0.0.1:43123/v1" {
		t.Fatalf("primary reviewer config = baseURL %q apiKey %q", primary.baseURL, primary.apiKey)
	}
	if _, ok := got.fallback.(*cloudDemuxReviewer); !ok {
		t.Fatalf("fallback reviewer = %T, want *cloudDemuxReviewer", got.fallback)
	}
}

func TestDemuxReviewerFromEnvUsesGXCloudOpenAIProxy(t *testing.T) {
	var gotAuth string
	var gotPayload cloudChatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/gx/openai/chat-completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		fmt.Fprint(w, `{"model":"test-model","choices":[{"message":{"content":"{\"revisions\":[{\"id\":\"r1\",\"intent\":\"repair\",\"use_hunks\":true,\"hunk_ids\":[\"h1\"]}]}"}}]}`)
	}))
	defer server.Close()

	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_CLOUD_URL", server.URL+"/gx/pr")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_demux",
		CLISessionToken:   "gxcs_demux",
		UserID:            "demux-user",
		Login:             "demux",
		MachineID:         "demux-machine",
		ObtainedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	reviewer, model, err := demuxReviewerFromEnv(DemuxAIReviewOptions{Model: "test-model"})
	if err != nil {
		t.Fatalf("demuxReviewerFromEnv() error = %v", err)
	}
	if model != "test-model" {
		t.Fatalf("model = %q", model)
	}
	if _, ok := reviewer.(*cloudDemuxReviewer); !ok {
		t.Fatalf("reviewer = %T, want *cloudDemuxReviewer", reviewer)
	}
	response, err := reviewer.ReviewDemuxProposal(context.Background(), demuxAIReviewRequest{
		Proposal: demuxAIProposalForReview{ID: "demux-test"},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxProposal() error = %v", err)
	}
	if gotAuth != "Bearer gxcs_demux" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotPayload.Model != "test-model" || len(gotPayload.Messages) != 2 {
		t.Fatalf("payload = %#v", gotPayload)
	}
	if gotPayload.ResponseFormat["type"] != "json_object" {
		t.Fatalf("response_format = %#v", gotPayload.ResponseFormat)
	}
	if len(response.Revisions) != 1 || response.Revisions[0].ID != "r1" {
		t.Fatalf("response = %#v", response)
	}
}

func TestProposalFromAIReviewResponseAcceptsRevisionOnlyOutput(t *testing.T) {
	got := proposalFromAIReviewResponse(DemuxProposal{ID: "demux-1"}, demuxAIReviewResponse{
		Revisions: []RevisionProposal{{
			ID:      "u1",
			Intent:  "alpha",
			HunkIDs: []string{"h1"},
		}},
	})
	if got.ID != "demux-1" || len(got.Revisions) != 1 || got.Revisions[0].ID != "u1" {
		t.Fatalf("proposalFromAIReviewResponse() = %#v", got)
	}
}

func TestProposalFromAIReviewResponsePreservesBaseMetadata(t *testing.T) {
	got := proposalFromAIReviewResponse(DemuxProposal{
		ID:               "demux-1",
		RepoRoot:         "/repo",
		ProposedChangeID: "change-1",
		ProposedCommitID: "commit-1",
		Status:           ProposalPending,
		Hunks:            []HunkRange{{ID: "h1", File: "app.go"}},
	}, demuxAIReviewResponse{
		Revisions: []RevisionProposal{{
			ID:       "u1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
		}},
	})
	if got.RepoRoot != "/repo" || got.ProposedChangeID != "change-1" || got.ProposedCommitID != "commit-1" || got.Status != ProposalPending {
		t.Fatalf("metadata = %#v", got)
	}
	if len(got.Hunks) != 1 || got.Hunks[0].ID != "h1" {
		t.Fatalf("hunks = %#v", got.Hunks)
	}
	if len(got.Revisions) != 1 || got.Revisions[0].ID != "u1" {
		t.Fatalf("revisions = %#v", got.Revisions)
	}
}

func TestAutoRepairDemuxProposalReordersAndAddsDependsOn(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "caller"},
			{ID: "u2", Intent: "helper"},
		},
	}
	got, changed := autoRepairDemuxProposal(proposal, []RepairHint{{
		Kind:       "reorder_dependency",
		RevisionID: "u1",
		DependsOn:  "u2",
	}})
	if !changed {
		t.Fatal("autoRepairDemuxProposal() changed = false, want true")
	}
	if got.Revisions[0].ID != "u2" || got.Revisions[1].ID != "u1" {
		t.Fatalf("autoRepairDemuxProposal() order = %#v", got.Revisions)
	}
	if !revisionDependsOn(got.Revisions[1], "u2") {
		t.Fatalf("autoRepairDemuxProposal() did not add depends_on: %#v", got.Revisions[1])
	}
}

func TestAutoRepairDemuxProposalRemovesInvalidDependsOn(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "caller", DependsOn: []string{"u2", "missing"}},
			{ID: "u2", Intent: "helper"},
		},
	}
	got, changed := autoRepairDemuxProposal(proposal, nil)
	if !changed {
		t.Fatal("autoRepairDemuxProposal() changed = false, want true")
	}
	if len(got.Revisions[0].DependsOn) != 0 {
		t.Fatalf("autoRepairDemuxProposal() kept invalid depends_on: %#v", got.Revisions[0])
	}
}

func TestAutoRepairDemuxProposalPrunesTransitiveDependsOn(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "base"},
			{ID: "u2", Intent: "middle", DependsOn: []string{"u1"}},
			{ID: "u3", Intent: "leaf", DependsOn: []string{"u1", "u2"}},
		},
	}
	got, changed := autoRepairDemuxProposal(proposal, []RepairHint{{
		Kind:       "inferred_dependency",
		RevisionID: "u3",
		DependsOn:  "u2",
	}})
	if !changed {
		t.Fatal("autoRepairDemuxProposal() changed = false, want true")
	}
	if !reflect.DeepEqual(got.Revisions[2].DependsOn, []string{"u2"}) {
		t.Fatalf("leaf depends_on = %#v, want only direct dependency u2", got.Revisions[2].DependsOn)
	}
}

func TestAutoRepairDemuxProposalCoalescesTinyRevisionsWithoutHints(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "@@\n+const first = true\n", Symbol: "Flags"},
			{ID: "h2", File: "app.go", Patch: "@@\n+const second = true\n", Symbol: "Flags"},
		},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "first flag", Files: []string{"app.go"}, UseHunks: true, HunkIDs: []string{"h1"}, Hunks: []HunkRange{{ID: "h1", File: "app.go", Patch: "@@\n+const first = true\n", Symbol: "Flags"}}},
			{ID: "u2", Intent: "second flag", Files: []string{"app.go"}, UseHunks: true, HunkIDs: []string{"h2"}, Hunks: []HunkRange{{ID: "h2", File: "app.go", Patch: "@@\n+const second = true\n", Symbol: "Flags"}}},
		},
	}

	got, changed := autoRepairDemuxProposal(proposal, nil)

	if !changed {
		t.Fatalf("autoRepairDemuxProposal() changed = false, want deterministic shape repair")
	}
	if len(got.Revisions) != 1 || !reflect.DeepEqual(got.Revisions[0].HunkIDs, []string{"h1", "h2"}) {
		t.Fatalf("autoRepairDemuxProposal() revisions = %#v, want one coalesced revision", got.Revisions)
	}
}

func TestProposalForAIReviewStripsHeavyDiagnostics(t *testing.T) {
	got := proposalForAIReview(DemuxProposal{
		ID: "demux-1",
		Hunks: []HunkRange{{
			ID:    "h1",
			File:  "alpha.go",
			Patch: "@@ huge patch",
		}},
		Revisions: []RevisionProposal{{
			ID:    "u1",
			Hunks: []HunkRange{{ID: "h1", Patch: "@@ huge patch"}},
		}},
		FeasibilityWarnings: []FeasibilityWarning{{Message: "hidden"}},
		Warnings:            []string{"hidden"},
		StructuralFacts:     []StructuralFact{{File: "alpha.go", DefinedSymbols: []string{"Alpha"}}},
	}, 2)
	if got.Hunks[0].Patch != "" {
		t.Fatalf("proposalForAIReview() kept hunk patch: %#v", got.Hunks[0])
	}
	if len(got.Revisions[0].Hunks) != 0 {
		t.Fatalf("proposalForAIReview() kept revision hunks: %#v", got.Revisions[0].Hunks)
	}
	if len(got.StructuralDeps) != 0 {
		t.Fatalf("proposalForAIReview() kept heavy diagnostics: %#v", got)
	}
}

func TestOpenAIDemuxReviewerIncludesErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"bad request detail"}}`, http.StatusBadRequest)
	}))
	defer server.Close()

	reviewer := newTestOpenAIDemuxReviewer(server, "test-model")
	_, err := reviewer.ReviewDemuxProposal(context.Background(), demuxAIReviewRequest{})
	if err == nil || !strings.Contains(err.Error(), "400 Bad Request") || !strings.Contains(err.Error(), "bad request detail") {
		t.Fatalf("ReviewDemuxProposal() error = %v, want status and body detail", err)
	}
}

func newTestOpenAIDemuxReviewer(server *httptest.Server, model string) *openAIDemuxReviewer {
	baseURL := normalizeOpenAIBaseURL(server.URL)
	return &openAIDemuxReviewer{
		apiKey:  "test-key",
		baseURL: baseURL,
		model:   model,
		client: openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey("test-key"),
			option.WithHTTPClient(server.Client()),
		),
	}
}

func TestRepairReviewedDemuxProposalDoesNotFailOnLightweightStructuralHints(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_OPENAI_API_KEY", "test-key")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("GX_DEMUX_REVIEW_MODEL", "test-model")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "test-model",
			"choices": []map[string]any{{
				"message": map[string]string{
					"role":    "assistant",
					"content": `{"revisions":[{"id":"u1","intent":"caller","files":["caller.go"],"provenance_status":"absent","confidence":0.5},{"id":"u2","intent":"helper","files":["helper.go"],"provenance_status":"absent","confidence":0.5}],"notes":["could not repair"]}`,
				},
			}},
		})
	}))
	defer server.Close()
	t.Setenv("GX_OPENAI_BASE_URL", server.URL)

	engine := NewEngine()
	repoRoot := repoRootForTest(t)
	proposal := DemuxProposal{
		ID:               "demux-test",
		RepoRoot:         repoRoot,
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		StructuralFacts: []StructuralFact{
			{File: "caller.go", Language: "go", DefinedSymbols: []string{"Run"}, ReferencedSymbols: []string{"Help"}},
			{File: "helper.go", Language: "go", DefinedSymbols: []string{"Help"}},
		},
		StructuralDeps: []StructuralDependency{{
			FromFile: "caller.go",
			ToFile:   "helper.go",
			Symbol:   "Help",
		}, {
			FromFile: "helper.go",
			ToFile:   "caller.go",
			Symbol:   "Run",
		}},
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "caller", Files: []string{"caller.go"}, ProvenanceStatus: "absent", Confidence: 0.5},
			{ID: "u2", Intent: "helper", Files: []string{"helper.go"}, ProvenanceStatus: "absent", Confidence: 0.5},
		},
	}
	saved, err := engine.SaveDemuxProposal(context.Background(), proposal)
	if err != nil {
		t.Fatalf("SaveDemuxProposal() error = %v", err)
	}
	review, err := engine.ReviewDemuxPlan(context.Background(), saved)
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}

	result, err := engine.repairReviewedDemuxProposal(context.Background(), saved, review, DemuxAIReviewOptions{})
	if err != nil {
		t.Fatalf("repairReviewedDemuxProposal() error = %v", err)
	}
	if result.State != DemuxWorkflowReadyToApply {
		t.Fatalf("repairReviewedDemuxProposal() state = %q, want ready_to_apply", result.State)
	}
}

func TestRepairReviewedDemuxProposalDeterministicallyFixesInvalidAIDependsOn(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_OPENAI_API_KEY", "test-key")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("GX_DEMUX_REVIEW_MODEL", "test-model")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "test-model",
			"choices": []map[string]any{{
				"message": map[string]string{
					"role": "assistant",
					"content": strings.Join([]string{
						`{"revisions":[`,
						`{"id":"r1","intent":"app","use_hunks":true,"hunk_ids":["h1"],"depends_on":["r2"]},`,
						`{"id":"r2","intent":"helper","use_hunks":true,"hunk_ids":["h2"]}`,
						`],"notes":["assigned missing hunk"]}`,
					}, ""),
				},
			}},
		})
	}))
	defer server.Close()
	t.Setenv("GX_OPENAI_BASE_URL", server.URL)

	engine := NewEngine()
	saved, err := engine.SaveDemuxProposal(context.Background(), DemuxProposal{
		ID:               "demux-test",
		RepoRoot:         repoRootForTest(t),
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "patch app"},
			{ID: "h2", File: "helper.go", Patch: "patch helper"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "app", UseHunks: true, HunkIDs: []string{"h1"}},
		},
	})
	if err != nil {
		t.Fatalf("SaveDemuxProposal() error = %v", err)
	}
	review, err := engine.ReviewDemuxPlan(context.Background(), saved)
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}

	result, err := engine.repairReviewedDemuxProposal(context.Background(), saved, review, DemuxAIReviewOptions{})
	if err != nil {
		t.Fatalf("repairReviewedDemuxProposal() error = %v", err)
	}
	if result.State != DemuxWorkflowReadyToApply {
		t.Fatalf("repairReviewedDemuxProposal() state = %q, want ready_to_apply; errors=%#v", result.State, result.Review.Errors)
	}
	for _, revision := range result.Proposal.Revisions {
		if len(revision.DependsOn) != 0 {
			t.Fatalf("revision %s kept invalid depends_on: %#v", revision.ID, revision.DependsOn)
		}
	}
}

func TestReviewDemuxProposalWithAIPlanLeavesAdvisoryDependencyReady(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	engine := NewEngine()
	saved, err := engine.SaveDemuxProposal(context.Background(), DemuxProposal{
		ID:               "demux-test",
		RepoRoot:         repoRootForTest(t),
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "patch app", Symbol: "Run"},
			{ID: "h2", File: "helper.go", Patch: "patch helper", Symbol: "NewThing"},
		},
		StructuralDeps: []StructuralDependency{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "app", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "helper", UseHunks: true, HunkIDs: []string{"h2"}},
		},
	})
	if err != nil {
		t.Fatalf("SaveDemuxProposal() error = %v", err)
	}

	result, err := engine.ReviewDemuxProposalWithAI(context.Background(), saved.ID, DemuxAIReviewOptions{PlanOnly: true})
	if err != nil {
		t.Fatalf("ReviewDemuxProposalWithAI() error = %v", err)
	}
	if result.Updated || result.State != DemuxWorkflowReadyToApply {
		t.Fatalf("ReviewDemuxProposalWithAI() = updated %t state %q, want unchanged ready_to_apply", result.Updated, result.State)
	}

	packet, err := engine.ShowDemuxProposal(context.Background(), saved.ID)
	if err != nil {
		t.Fatalf("ShowDemuxProposal() error = %v", err)
	}
	if packet.State != DemuxWorkflowReadyToApply {
		t.Fatalf("ShowDemuxProposal() state = %q, want ready_to_apply; warnings=%#v", packet.State, packet.Proposal.FeasibilityWarnings)
	}
	if got := revisionIDs(packet.Proposal.Revisions); !reflect.DeepEqual(got, []string{"r1", "r2"}) {
		t.Fatalf("persisted revision order = %#v, want original advisory order", got)
	}
}

func TestReviewDemuxProposalWithRealAIRepairsUnassignedHunk(t *testing.T) {
	if os.Getenv("GX_REAL_AI_TESTS") != "1" {
		t.Skip("set GX_REAL_AI_TESTS=1 to run real OpenAI demux repair")
	}
	if strings.TrimSpace(firstNonEmpty(os.Getenv("GX_OPENAI_API_KEY"), os.Getenv("OPENAI_API_KEY"))) == "" {
		t.Skip("OPENAI_API_KEY or GX_OPENAI_API_KEY is required for real OpenAI demux repair")
	}
	t.Setenv("GX_HOME", t.TempDir())
	engine := NewEngine()
	saved, err := engine.SaveDemuxProposal(context.Background(), DemuxProposal{
		ID:               "demux-real-ai-test",
		RepoRoot:         repoRootForTest(t),
		ProposedChangeID: "change-real-ai",
		ProposedCommitID: "commit-real-ai",
		Status:           ProposalPending,
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Patch: "@@\n+func Run() {}\n", Symbol: "Run"},
			{ID: "h2", File: "worker.go", Patch: "@@\n+func Work() {}\n", Symbol: "Work"},
		},
		Revisions: []RevisionProposal{{
			ID:               "r1",
			Intent:           "add Run entrypoint",
			Files:            []string{"app.go"},
			UseHunks:         true,
			HunkIDs:          []string{"h1"},
			ProvenanceStatus: "absent",
			Confidence:       0.7,
		}},
	})
	if err != nil {
		t.Fatalf("SaveDemuxProposal() error = %v", err)
	}
	planOnly, err := engine.ReviewDemuxProposalWithAI(context.Background(), saved.ID, DemuxAIReviewOptions{PlanOnly: true})
	if err != nil {
		t.Fatalf("plan-only ReviewDemuxProposalWithAI() error = %v", err)
	}
	if planOnly.Updated || planOnly.State == DemuxWorkflowReadyToApply || planOnly.RepairHintTotal == 0 {
		t.Fatalf("plan-only repair = %#v, want unresolved deterministic state with repair hint", planOnly)
	}

	result, err := engine.ReviewDemuxProposalWithAI(context.Background(), saved.ID, DemuxAIReviewOptions{
		Model:       strings.TrimSpace(os.Getenv("GX_REAL_AI_DEMUX_MODEL")),
		MaxWarnings: 10,
	})
	if err != nil {
		t.Fatalf("real AI ReviewDemuxProposalWithAI() error = %v\nresult=%#v", err, result)
	}
	if !result.Updated || result.State != DemuxWorkflowReadyToApply || !result.Review.Valid {
		t.Fatalf("real AI repair = %#v, want updated ready proposal", result)
	}
	if result.Model == "" || result.Model == "deterministic" {
		t.Fatalf("real AI model = %q, want non-deterministic model name", result.Model)
	}
	if len(result.Review.RepairHints) != 0 || len(result.Review.Errors) != 0 {
		t.Fatalf("real AI review hints/errors = hints %#v errors %#v, want none", result.Review.RepairHints, result.Review.Errors)
	}
	covered := map[string]bool{}
	for _, revision := range result.Proposal.Revisions {
		for _, hunkID := range revision.HunkIDs {
			covered[hunkID] = true
		}
		if !revision.UseHunks {
			for _, file := range revision.Files {
				if file == "app.go" {
					covered["h1"] = true
				}
				if file == "worker.go" {
					covered["h2"] = true
				}
			}
		}
	}
	if !covered["h1"] || !covered["h2"] {
		t.Fatalf("real AI repaired revisions = %#v, want coverage for h1 and h2", result.Proposal.Revisions)
	}
}

func TestDemuxProgressWritesFormattedLine(t *testing.T) {
	var out strings.Builder

	demuxProgress(&out, "Fixing compose proposal %s...", "demux-1")

	if got, want := out.String(), "Fixing compose proposal demux-1...\n"; got != want {
		t.Fatalf("demuxProgress() = %q, want %q", got, want)
	}
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Dir(filepath.Dir(wd))
}

func TestBlockingFeasibilityWarningsOnlyBlocksWarningSeverity(t *testing.T) {
	got := blockingFeasibilityWarnings([]FeasibilityWarning{
		{Severity: "info", Source: "inferred_dependency"},
		{Severity: "warning", Source: "structural_dependency"},
	})
	if len(got) != 1 || got[0].Source != "structural_dependency" {
		t.Fatalf("blockingFeasibilityWarnings() = %#v, want only warning severity", got)
	}
}

func TestRouteTextScoreMatchesConventionalBookmarkSlug(t *testing.T) {
	got := routeTextScore("update terminal theme status", "feature/update-terminal-theme")
	if got < 0.8 {
		t.Fatalf("routeTextScore() = %v, want strong bookmark slug match", got)
	}
}

func TestReviewDemuxRoutesAllowsUnknownTargetAsNewStack(t *testing.T) {
	repoRoot := setupDemuxRouteStore(t)
	engine := NewEngine()
	review, err := engine.ReviewDemuxPlan(context.Background(), DemuxProposal{
		RepoRoot:         repoRoot,
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks:            []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch"}},
		Revisions: []RevisionProposal{{
			ID:          "r1",
			Intent:      "alpha",
			Files:       []string{"alpha.txt"},
			TargetStack: "feature/missing",
		}},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}
	if !review.Valid || len(review.Errors) != 0 {
		t.Fatalf("ReviewDemuxPlan() = %#v, want valid new stack route", review)
	}
	if len(review.Proposal.FeasibilityWarnings) < 2 || review.Proposal.FeasibilityWarnings[1].Source != "demux_route_new_stack" {
		t.Fatalf("ReviewDemuxPlan() warnings = %#v, want new stack route warning", review.Proposal.FeasibilityWarnings)
	}
}

func TestPlanDemuxRoutesUsesExistingStackHeuristic(t *testing.T) {
	repoRoot := setupDemuxRouteStore(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{{
			ID:     "r1",
			Intent: "update terminal theme colors",
			Files:  []string{"theme.go"},
		}},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	revision := proposal.Revisions[0]
	if revision.TargetStack != "feature/update-terminal-theme" || revision.RouteSource != routeSourceHeuristic {
		t.Fatalf("planned route = %#v, want update-terminal-theme heuristic route", revision)
	}
}

func TestPlanDemuxRoutesUsesStoredFileOverlapForExistingStack(t *testing.T) {
	repoRoot, repoID, store, cleanup := setupDemuxRouteStoreWithRepo(t)
	defer cleanup()
	ctx := context.Background()
	head := "change-status"
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "status display",
		BookmarkName: "feature/status-display",
		BaseRef:      "main",
		HeadChangeID: &head,
		Status:       "draft",
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      head,
		CurrentCommitID: "commit-status",
		Description:     "render status remote state",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	files, _ := json.Marshal([]string{"internal/cli/root.go", "internal/vcs/status_snapshot.go"})
	if err := store.WriteChangeRevision(ctx, storage.ChangeRevision{
		ChangeID:      changeID,
		JJCommitID:    "commit-status",
		JJOperationID: "op-status",
		ChangedFiles:  string(files),
		CreatedAt:     1,
	}); err != nil {
		t.Fatalf("WriteChangeRevision() error = %v", err)
	}
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(ctx, DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{{
			ID:     "r1",
			Intent: "tighten remote badges",
			Files:  []string{"internal/cli/root.go"},
		}},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	revision := proposal.Revisions[0]
	if revision.TargetStack != "feature/status-display" || !strings.Contains(revision.RouteReason, "files overlap") {
		t.Fatalf("planned route = %#v, want semantic file-overlap route", revision)
	}
}

func TestPlanDemuxRoutesInfersNewStacksForMultipleClusters(t *testing.T) {
	repoRoot := setupDemuxRouteRepo(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{
			{
				ID:     "r1",
				Intent: "update demux route planning",
				Files:  []string{"internal/authoring/demux_routing.go"},
			},
			{
				ID:     "r2",
				Intent: "update demux repair tests",
				Files:  []string{"internal/authoring/demux_test.go"},
			},
			{
				ID:     "r3",
				Intent: "create explicit stack command",
				Files:  []string{"internal/cli/root.go"},
			},
			{
				ID:     "r4",
				Intent: "update stack storage rows",
				Files:  []string{"internal/storage/schema.sql"},
			},
		},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	want := map[string]string{
		"r1": "feature/demux-routing",
		"r2": "feature/demux-routing",
		"r3": "feature/stack-management",
		"r4": "feature/stack-management",
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != want[revision.ID] || revision.RouteSource != routeSourceHeuristic {
			t.Fatalf("revision %s route = %#v, want target %q from heuristic", revision.ID, revision, want[revision.ID])
		}
	}
}

func TestPlanDemuxRoutesUsesConventionalNewStackPrefixes(t *testing.T) {
	repoRoot := setupDemuxRouteRepo(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "fix daemon startup panic", Files: []string{"internal/daemon/server.go"}},
			{ID: "r2", Intent: "update README guidance", Files: []string{"README.md"}},
			{ID: "r3", Intent: "update module dependencies", Files: []string{"go.mod", "go.sum"}},
		},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	want := map[string]string{
		"r1": "bug/provider-runtime",
		"r2": "docs/documentation",
		"r3": "chore/go-mod",
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != want[revision.ID] {
			t.Fatalf("revision %s target_stack = %q, want %q", revision.ID, revision.TargetStack, want[revision.ID])
		}
		if strings.HasPrefix(revision.TargetStack, "gx/") {
			t.Fatalf("revision %s target_stack kept legacy gx prefix: %q", revision.ID, revision.TargetStack)
		}
	}
}

func TestPlanDemuxRoutesUsesFeaturePrefixForMixedSourceAndTests(t *testing.T) {
	repoRoot := setupDemuxRouteRepo(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "test demux route repair", Files: []string{"internal/authoring/demux_test.go"}},
			{ID: "r2", Intent: "update demux route repair", Files: []string{"internal/authoring/demux_routing.go"}},
			{ID: "r3", Intent: "update CLI compose output", Files: []string{"internal/cli/service.go"}},
		},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	want := map[string]string{
		"r1": "feature/demux-routing",
		"r2": "feature/demux-routing",
		"r3": "feature/cli",
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != want[revision.ID] {
			t.Fatalf("revision %s target_stack = %q, want %q", revision.ID, revision.TargetStack, want[revision.ID])
		}
	}
}

func TestPrepareNextDemuxRevisionEditsParentCommitID(t *testing.T) {
	repoRoot := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	runner := &prepareNextDemuxFakeRunner{repoRoot: repoRoot}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))

	if err := engine.prepareNextDemuxRevision(context.Background(), repoRoot, "recorded-change"); err != nil {
		t.Fatalf("prepareNextDemuxRevision() error = %v", err)
	}
	if !runner.called("jj edit parent-commit") {
		t.Fatalf("prepareNextDemuxRevision() did not edit parent commit id; calls=%#v", runner.calls)
	}
	if runner.called("jj edit parent-change") {
		t.Fatalf("prepareNextDemuxRevision() edited divergent parent change id; calls=%#v", runner.calls)
	}
}

func TestCurrentDemuxApplyChangeUsesNonEmptyWorkRev(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &demuxApplySourceFakeRunner{repoRoot: repoRoot}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))

	workRev, change, err := engine.currentDemuxApplyChange(context.Background(), repoRoot)
	if err != nil {
		t.Fatalf("currentDemuxApplyChange() error = %v", err)
	}
	if workRev != "@-" {
		t.Fatalf("workRev = %q, want @-", workRev)
	}
	if change.ChangeID != "parent-change" || change.CommitID != "parent-commit" {
		t.Fatalf("change = %#v, want parent source", change)
	}
	if !reflect.DeepEqual(change.Files, []string{"server/src/app.ts"}) {
		t.Fatalf("change files = %#v, want source files", change.Files)
	}
}

func TestEditDemuxApplySourceIfNeededEditsSourceChange(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &demuxApplySourceFakeRunner{repoRoot: repoRoot}
	engine := NewEngineWithVCS(vcs.NewServiceWithRunner(runner))
	source := demuxSourceLocation{ChangeID: "parent-change", CommitID: "parent-commit"}

	if err := engine.editDemuxApplySourceIfNeeded(context.Background(), repoRoot, "@-", source); err != nil {
		t.Fatalf("editDemuxApplySourceIfNeeded() error = %v", err)
	}
	if !runner.called("jj edit parent-change") {
		t.Fatalf("editDemuxApplySourceIfNeeded() did not edit the source change; calls=%#v", runner.calls)
	}
	runner.calls = nil
	if err := engine.editDemuxApplySourceIfNeeded(context.Background(), repoRoot, "@", source); err != nil {
		t.Fatalf("editDemuxApplySourceIfNeeded(@) error = %v", err)
	}
	if runner.called("jj edit parent-change") {
		t.Fatalf("editDemuxApplySourceIfNeeded(@) edited unexpectedly; calls=%#v", runner.calls)
	}
}

func TestPlanDemuxRoutesKeepsTestPrefixForTestOnlyStack(t *testing.T) {
	repoRoot := setupDemuxRouteRepo(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "test demux route repair", Files: []string{"internal/authoring/demux_test.go"}},
			{ID: "r2", Intent: "update README guidance", Files: []string{"README.md"}},
		},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	want := map[string]string{
		"r1": "test/demux-routing",
		"r2": "docs/documentation",
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != want[revision.ID] {
			t.Fatalf("revision %s target_stack = %q, want %q", revision.ID, revision.TargetStack, want[revision.ID])
		}
	}
}

func TestNormalizeConventionalDemuxStackRoutesUsesLegacyRouteGroupKind(t *testing.T) {
	got := normalizeConventionalDemuxStackRoutes(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "test authoring shape", Files: []string{"internal/authoring/demux_test.go"}, TargetStack: "gx/internal-authoring"},
			{ID: "r2", Intent: "update authoring shape", Files: []string{"internal/authoring/demux_review_shape.go"}, TargetStack: "gx/internal-authoring"},
			{ID: "r3", Intent: "test daemon manager", Files: []string{"internal/launcher/daemon_mgr_test.go"}, TargetStack: "gx/internal-launcher"},
			{ID: "r4", Intent: "update CLI compose output", Files: []string{"internal/cli/service.go", "internal/cli/service_test.go"}, TargetStack: "test/cli"},
		},
	})
	want := map[string]string{
		"r1": "feature/internal-authoring",
		"r2": "feature/internal-authoring",
		"r3": "test/internal-launcher",
		"r4": "feature/cli",
	}
	for _, revision := range got.Revisions {
		if revision.TargetStack != want[revision.ID] {
			t.Fatalf("revision %s target_stack = %q, want %q", revision.ID, revision.TargetStack, want[revision.ID])
		}
	}
}

func TestPlanDemuxRoutesInfersConventionalStackForSingleCluster(t *testing.T) {
	repoRoot := setupDemuxRouteRepo(t)
	engine := NewEngine()
	proposal, err := engine.planDemuxRoutes(context.Background(), DemuxProposal{
		RepoRoot: repoRoot,
		Revisions: []RevisionProposal{
			{
				ID:     "r1",
				Intent: "update demux route planning",
				Files:  []string{"internal/authoring/demux_routing.go"},
			},
			{
				ID:     "r2",
				Intent: "update demux repair tests",
				Files:  []string{"internal/authoring/demux_test.go"},
			},
		},
	})
	if err != nil {
		t.Fatalf("planDemuxRoutes() error = %v", err)
	}
	want := map[string]string{
		"r1": "feature/demux-routing",
		"r2": "feature/demux-routing",
	}
	for _, revision := range proposal.Revisions {
		if revision.TargetStack != want[revision.ID] {
			t.Fatalf("revision %s target_stack = %q, want %q", revision.ID, revision.TargetStack, want[revision.ID])
		}
	}
}

func TestReviewDemuxRoutesAcceptsBaseStackForCrossStackDependency(t *testing.T) {
	repoRoot := setupDemuxRouteStore(t)
	engine := NewEngine()
	review, err := engine.ReviewDemuxPlan(context.Background(), DemuxProposal{
		RepoRoot:         repoRoot,
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks: []HunkRange{
			{ID: "h1", File: "theme.go", Patch: "patch one"},
			{ID: "h2", File: "other.go", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "theme", Files: []string{"theme.go"}, TargetStack: "feature/update-terminal-theme"},
			{ID: "r2", Intent: "other", Files: []string{"other.go"}, TargetStack: "feature/other", BaseStack: "feature/update-terminal-theme", DependsOn: []string{"r1"}},
		},
	})
	if err != nil {
		t.Fatalf("ReviewDemuxPlan() error = %v", err)
	}
	if !review.Valid {
		t.Fatalf("ReviewDemuxPlan() valid = false: %#v", review)
	}
	for _, warning := range review.Proposal.FeasibilityWarnings {
		if warning.Source == "demux_route_cross_stack_dependency" {
			t.Fatalf("ReviewDemuxPlan() warning = %#v, want base_stack to resolve cross-stack dependency", warning)
		}
	}
}

func TestAutoRepairDemuxProposalAppliesStackRouteHints(t *testing.T) {
	proposal := DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "theme", TargetStack: "s1"},
			{ID: "r2", Intent: "other", TargetStack: "feature/other", DependsOn: []string{"r1"}},
		},
	}
	got, changed := autoRepairDemuxProposal(proposal, []RepairHint{
		{Kind: "invalid_demux_route", RevisionID: "r1", TargetStack: "feature/update-terminal-theme"},
		{Kind: "cross_stack_dependency", RevisionID: "r2", DependsOn: "r1", BaseStack: "feature/update-terminal-theme"},
	})
	if !changed {
		t.Fatal("autoRepairDemuxProposal() changed = false, want true")
	}
	if got.Revisions[0].TargetStack != "feature/update-terminal-theme" {
		t.Fatalf("first target_stack = %q, want canonical stack", got.Revisions[0].TargetStack)
	}
	if got.Revisions[1].BaseStack != "feature/update-terminal-theme" {
		t.Fatalf("second base_stack = %q, want dependency stack", got.Revisions[1].BaseStack)
	}
}

func setupDemuxRouteStore(t *testing.T) string {
	t.Helper()
	repoRoot, repoID, store, cleanup := setupDemuxRouteStoreWithRepo(t)
	defer cleanup()
	head := "change-theme"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "update terminal theme",
		BookmarkName: "feature/update-terminal-theme",
		BaseRef:      "main",
		HeadChangeID: &head,
		Status:       "draft",
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	otherHead := "change-other"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "other",
		BookmarkName: "feature/other",
		BaseRef:      "main",
		HeadChangeID: &otherHead,
		Status:       "draft",
	}); err != nil {
		t.Fatalf("UpsertStack(other) error = %v", err)
	}
	return repoRoot
}

func setupDemuxRouteRepo(t *testing.T) string {
	t.Helper()
	repoRoot, _, _, cleanup := setupDemuxRouteStoreWithRepo(t)
	defer cleanup()
	return repoRoot
}

func setupDemuxRouteStoreWithRepo(t *testing.T) (string, int64, *storage.Store, func()) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		_ = db.Close()
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoRoot := filepath.Join(t.TempDir(), "repo")
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptrString("main"),
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	return repoRoot, repoID, store, func() { _ = db.Close() }
}

func ptrString(value string) *string {
	return &value
}

func TestNormalizeDemuxProposalRejectsUnassignedHunk(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.txt", Patch: "patch one"},
			{ID: "h2", File: "alpha.txt", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha first", UseHunks: true, HunkIDs: []string{"h1"}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "hunk h2 in alpha.txt is not assigned") {
		t.Fatalf("normalizeDemuxProposal() error = %v, want unassigned hunk error", err)
	}
}

func TestNormalizeDemuxProposalAllowsWholeFileCoverage(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "alpha.txt", Patch: "patch one"},
			{ID: "h2", File: "alpha.txt", Patch: "patch two"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha whole file", Files: []string{"alpha.txt"}},
		},
	})
	if err != nil {
		t.Fatalf("normalizeDemuxProposal() error = %v", err)
	}
}

func TestNormalizeDemuxProposalAllowsWholeFileRevisionWithContextHunks(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{
				ID:       "r1",
				Intent:   "alpha whole file",
				Files:    []string{"alpha.txt"},
				HunkIDs:  []string{"h1"},
				Hunks:    []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
				UseHunks: false,
			},
		},
	})
	if err != nil {
		t.Fatalf("normalizeDemuxProposal() error = %v", err)
	}
}

func TestNormalizeDemuxProposalRejectsMixedHunkAndWholeFileCoverage(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha selected", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r2", Intent: "alpha whole file", Files: []string{"alpha.txt"}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "covered by hunk revision r1 and whole-file revision r2") {
		t.Fatalf("normalizeDemuxProposal() error = %v, want mixed coverage error", err)
	}
}

func TestNormalizeDemuxProposalRejectsDuplicateRevisionIDs(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha", UseHunks: true, HunkIDs: []string{"h1"}},
			{ID: "r1", Intent: "beta", Files: []string{"beta.txt"}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "duplicate revision id r1") {
		t.Fatalf("normalizeDemuxProposal() error = %v, want duplicate revision error", err)
	}
}

func TestNormalizeDemuxProposalRejectsLaterDependency(t *testing.T) {
	_, err := normalizeDemuxProposal(DemuxProposal{
		Hunks: []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "alpha", UseHunks: true, HunkIDs: []string{"h1"}, DependsOn: []string{"r2"}},
			{ID: "r2", Intent: "beta", Files: []string{"beta.txt"}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), "revision r1 depends on unknown or later revision r2") {
		t.Fatalf("normalizeDemuxProposal() error = %v, want later dependency error", err)
	}
}

func TestMergeDemuxPlanCarriesHunkCatalog(t *testing.T) {
	base := DemuxProposal{
		ID:               "demux-1",
		RepoRoot:         "/repo",
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks:            []HunkRange{{ID: "h1", File: "alpha.txt", Patch: "patch one"}},
	}
	plan := DemuxProposal{
		ID: "demux-1",
		Revisions: []RevisionProposal{{
			ID:       "r1",
			Intent:   "alpha",
			UseHunks: true,
			HunkIDs:  []string{"h1"},
		}},
	}
	merged := mergeDemuxPlan(base, plan)
	if merged.RepoRoot != "/repo" || merged.ProposedChangeID != "change" || merged.ProposedCommitID != "commit" {
		t.Fatalf("mergeDemuxPlan() identity = %#v", merged)
	}
	if len(merged.Hunks) != 1 || merged.Hunks[0].ID != "h1" {
		t.Fatalf("mergeDemuxPlan() hunks = %#v", merged.Hunks)
	}
}

func TestMergeDemuxPlanFiltersHunksForRevisionSubset(t *testing.T) {
	base := DemuxProposal{
		ID:               "demux-1",
		RepoRoot:         "/repo",
		ProposedChangeID: "change",
		ProposedCommitID: "commit",
		Status:           ProposalPending,
		Hunks: []HunkRange{
			{ID: "h1", File: "internal/authoring/demux.go", Patch: "patch demux"},
			{ID: "h2", File: "internal/storage/schema.sql", Patch: "patch schema"},
			{ID: "h3", File: "test/e2e/gx_e2e_test.go", Patch: "patch e2e"},
		},
		Revisions: []RevisionProposal{
			{ID: "r1", Intent: "demux", Files: []string{"internal/authoring/demux.go"}, TargetStack: "feature/demux-routing"},
			{ID: "r2", Intent: "schema", Files: []string{"internal/storage/schema.sql"}, TargetStack: "feature/stack-management"},
			{ID: "r3", Intent: "e2e", Files: []string{"test/e2e/gx_e2e_test.go"}, TargetStack: "test/e2e-tests"},
		},
	}
	plan := DemuxProposal{
		ID: "demux-1",
		Revisions: []RevisionProposal{{
			ID:          "r2",
			Intent:      "schema",
			Files:       []string{"internal/storage/schema.sql"},
			TargetStack: "feature/stack-management",
		}},
	}

	merged := mergeDemuxPlan(base, plan)
	if merged.RepoRoot != "/repo" || merged.ProposedChangeID != "change" || merged.ProposedCommitID != "commit" {
		t.Fatalf("mergeDemuxPlan() identity = %#v", merged)
	}
	if len(merged.Hunks) != 1 || merged.Hunks[0].ID != "h2" {
		t.Fatalf("mergeDemuxPlan() hunks = %#v, want only selected stack h2", merged.Hunks)
	}
}

func TestFeasibilityWarningsReportsUnmappedHunks(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{{ID: "u1", HunkIDs: []string{"h1"}}},
		[]HunkRange{{ID: "h1", File: "app.go"}},
		structural.Facts{},
	)
	if len(got) != 1 || got[0].RevisionID != "u1" || got[0].Source != "changed_symbol" {
		t.Fatalf("feasibilityWarnings() = %#v, want changed_symbol warning", got)
	}
}

func TestFeasibilityWarningsReportsStructuralOrderConflict(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{ID: "u1", Files: []string{"app.go"}},
			{ID: "u2", Files: []string{"helper.go"}},
		},
		nil,
		structural.Facts{Edges: []structural.DependencyEdge{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}}},
	)
	if len(got) != 1 || got[0].RevisionID != "u1" || got[0].Source != "structural_dependency" {
		t.Fatalf("feasibilityWarnings() = %#v, want structural_dependency warning", got)
	}
}

func TestFeasibilityWarningsUsesSymbolOwnerForDependencyTarget(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{
				ID:       "u1",
				Files:    []string{"helper.go"},
				UseHunks: true,
				Hunks:    []HunkRange{{ID: "h1", File: "helper.go", Symbol: "OldThing"}},
				HunkIDs:  []string{"h1"},
			},
			{
				ID:    "u2",
				Files: []string{"app.go"},
			},
			{
				ID:       "u3",
				Files:    []string{"helper.go"},
				UseHunks: true,
				Hunks:    []HunkRange{{ID: "h2", File: "helper.go", Symbol: "NewThing"}},
				HunkIDs:  []string{"h2"},
			},
		},
		nil,
		structural.Facts{Edges: []structural.DependencyEdge{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}}},
	)
	if len(got) != 1 || got[0].RevisionID != "u2" || !strings.Contains(got[0].Message, "u3 is proposed after it") {
		t.Fatalf("feasibilityWarnings() = %#v, want dependency on symbol owner u3", got)
	}
}

func TestFeasibilityWarningsUsesReferencingHunkAsDependencySource(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{
				ID:       "u1",
				Files:    []string{"app.go"},
				UseHunks: true,
				Hunks:    []HunkRange{{ID: "h1", File: "app.go", Symbol: "Other", Patch: "@@\n+func Other() {}\n"}},
				HunkIDs:  []string{"h1"},
			},
			{
				ID:       "u2",
				Files:    []string{"app.go"},
				UseHunks: true,
				Hunks:    []HunkRange{{ID: "h2", File: "app.go", Symbol: "Run", Patch: "@@\n+func Run() { NewThing() }\n"}},
				HunkIDs:  []string{"h2"},
			},
			{
				ID:       "u3",
				Files:    []string{"helper.go"},
				UseHunks: true,
				Hunks:    []HunkRange{{ID: "h3", File: "helper.go", Symbol: "NewThing", Patch: "@@\n+func NewThing() {}\n"}},
				HunkIDs:  []string{"h3"},
			},
		},
		nil,
		structural.Facts{Edges: []structural.DependencyEdge{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}}},
	)
	if len(got) != 1 || got[0].RevisionID != "u2" || !strings.Contains(got[0].Message, "u3 is proposed after it") {
		t.Fatalf("feasibilityWarnings() = %#v, want dependency warning on referencing hunk u2", got)
	}
}

func TestFeasibilityWarningsReportsMissingInferredDependsOn(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{ID: "u1", Files: []string{"helper.go"}},
			{ID: "u2", Files: []string{"app.go"}},
		},
		nil,
		structural.Facts{Edges: []structural.DependencyEdge{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}}},
	)
	if len(got) != 1 || got[0].RevisionID != "u2" || got[0].Source != "inferred_dependency" || got[0].Severity != "info" {
		t.Fatalf("feasibilityWarnings() = %#v, want inferred dependency info for u2", got)
	}
}

func TestFeasibilityWarningsSkipsInferredDependencyWhenDependsOnIsExplicit(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{ID: "u1", Files: []string{"helper.go"}},
			{ID: "u2", Files: []string{"app.go"}, DependsOn: []string{"u1"}},
		},
		nil,
		structural.Facts{Edges: []structural.DependencyEdge{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}}},
	)
	if len(got) != 0 {
		t.Fatalf("feasibilityWarnings() = %#v, want explicit depends_on to suppress inferred warning", got)
	}
}

func TestFeasibilityWarningsForProposalUsesFinalPlanOrder(t *testing.T) {
	proposal := DemuxProposal{
		Hunks: []HunkRange{
			{ID: "h1", File: "app.go", Symbol: "Run"},
			{ID: "h2", File: "helper.go", Symbol: "NewThing"},
		},
		StructuralFacts: []StructuralFact{
			{File: "app.go", DefinedSymbols: []string{"Run"}, ReferencedSymbols: []string{"NewThing"}},
			{File: "helper.go", DefinedSymbols: []string{"NewThing"}},
		},
		StructuralDeps: []StructuralDependency{{
			FromFile: "app.go",
			ToFile:   "helper.go",
			Symbol:   "NewThing",
		}},
		Revisions: []RevisionProposal{
			{ID: "r1", Files: []string{"helper.go"}, HunkIDs: []string{"h2"}},
			{ID: "r2", Files: []string{"app.go"}, HunkIDs: []string{"h1"}, DependsOn: []string{"r1"}},
		},
		FeasibilityWarnings: []FeasibilityWarning{{
			RevisionID: "stale",
			Severity:   "warning",
			Source:     "structural_dependency",
			Message:    "stale warning from deterministic proposal",
		}},
	}

	got := feasibilityWarningsForProposal(proposal)
	if len(got) != 0 {
		t.Fatalf("feasibilityWarningsForProposal() = %#v, want no stale/final warnings", got)
	}

	proposal.Revisions = []RevisionProposal{
		{ID: "r1", Files: []string{"app.go"}, HunkIDs: []string{"h1"}},
		{ID: "r2", Files: []string{"helper.go"}, HunkIDs: []string{"h2"}},
	}
	got = feasibilityWarningsForProposal(proposal)
	if len(got) != 1 || got[0].RevisionID != "r1" || got[0].Source != "structural_dependency" {
		t.Fatalf("feasibilityWarningsForProposal() = %#v, want final structural warning", got)
	}
}

func TestFeasibilityWarningsReportsSeparatedTestCounterpart(t *testing.T) {
	got := feasibilityWarnings(
		[]RevisionProposal{
			{ID: "u1", Files: []string{"foo_test.go"}},
			{ID: "u2", Files: []string{"foo.go"}},
		},
		nil,
		structural.Facts{},
	)
	if len(got) != 1 || got[0].RevisionID != "u1" || got[0].Source != "test_pairing" {
		t.Fatalf("feasibilityWarnings() = %#v, want test_pairing warning", got)
	}
}

type demuxRoutingFakeRunner struct {
	repoRoot          string
	calls             []string
	changeLogCalls    int
	diffNameOnlyCalls int
}

func (r *demuxRoutingFakeRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	if call == "git apply" || strings.HasPrefix(call, "git apply ") {
		return "", fmt.Errorf("git apply should not be used for routed demux revisions")
	}
	return r.output(name, args)
}

func (r *demuxRoutingFakeRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	return r.output(name, args)
}

func (r *demuxRoutingFakeRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

func (r *demuxRoutingFakeRunner) output(name string, args []string) (string, error) {
	if name == "git" {
		return r.gitOutput(args)
	}
	if name != "jj" {
		return "", nil
	}
	if len(args) == 1 && args[0] == "root" {
		return r.repoRoot + "\n", nil
	}
	if len(args) >= 3 && args[0] == "config" && args[1] == "get" {
		switch args[2] {
		case "user.name":
			return "Joe Example\n", nil
		case "user.email":
			return "joe@example.com\n", nil
		}
	}
	if len(args) > 0 {
		switch args[0] {
		case "edit", "new", "commit", "restore":
			return "", nil
		case "op":
			return "op-routed\n", nil
		case "bookmark":
			return r.jjBookmarkOutput(args), nil
		case "log":
			return r.jjLogOutput(args), nil
		case "diff":
			return r.jjDiffOutput(args), nil
		}
	}
	return "", nil
}

func (r *demuxRoutingFakeRunner) gitOutput(args []string) (string, error) {
	if len(args) >= 2 && args[0] == "config" {
		switch args[len(args)-1] {
		case "user.name":
			return "Joe Example\n", nil
		case "user.email":
			return "joe@example.com\n", nil
		}
	}
	if len(args) == 2 && args[0] == "rev-parse" && args[1] == "--git-path" {
		return filepath.Join(r.repoRoot, ".git", "info", "exclude") + "\n", nil
	}
	if len(args) == 2 && args[0] == "branch" && args[1] == "--show-current" {
		return "main\n", nil
	}
	return "", nil
}

func (r *demuxRoutingFakeRunner) jjBookmarkOutput(args []string) string {
	if len(args) >= 2 && args[1] == "list" {
		return "feature/internal-hooks|target-head\n"
	}
	return ""
}

func (r *demuxRoutingFakeRunner) jjLogOutput(args []string) string {
	template := ""
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "-T" {
			template = args[i+1]
			break
		}
	}
	switch template {
	case "empty":
		return "false\n"
	case "commit_id":
		return "target-commit\n"
	case "change_id":
		return "target-head\n"
	case "change_id ++ \"\\n\"":
		return "target-head\n"
	case `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`:
		r.changeLogCalls++
		switch r.changeLogCalls {
		case 1, 2:
			return "source-change|proposal-commit|compose source|base\n"
		case 3, 4:
			return "target-work|target-commit|target stack|target-head\n"
		default:
			return "routed-change|routed-commit|restore hooks|target-head\n"
		}
	default:
		return "routed-change|routed-commit|restore hooks|target-head\n"
	}
}

func (r *demuxRoutingFakeRunner) jjDiffOutput(args []string) string {
	if len(args) >= 2 && args[0] == "diff" && args[1] == "--from" {
		return "diff-content\n"
	}
	if len(args) >= 4 && args[0] == "diff" && args[1] == "-r" && args[2] == "@" && args[3] == "--name-only" {
		r.diffNameOnlyCalls++
		if r.diffNameOnlyCalls <= 2 {
			return "internal/hooks/run.go\ninternal/hooks/run_test.go\n"
		}
	}
	return ""
}

func (r *demuxRoutingFakeRunner) called(want string) bool {
	for _, call := range r.calls {
		if call == want || strings.HasPrefix(call, want+" ") {
			return true
		}
	}
	return false
}

type prepareNextDemuxFakeRunner struct {
	repoRoot string
	calls    []string
}

func (r *prepareNextDemuxFakeRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	if call == "jj edit parent-change" {
		return "", fmt.Errorf("Change ID `parent-change` is divergent")
	}
	return r.output(name, args)
}

func (r *prepareNextDemuxFakeRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	return r.output(name, args)
}

func (r *prepareNextDemuxFakeRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

func (r *prepareNextDemuxFakeRunner) output(name string, args []string) (string, error) {
	if name == "git" {
		if len(args) == 2 && args[0] == "branch" && args[1] == "--show-current" {
			return "main\n", nil
		}
		return "", nil
	}
	if name != "jj" || len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "root":
		return r.repoRoot + "\n", nil
	case "edit":
		return "", nil
	case "log":
		rev := ""
		template := ""
		for i := 0; i+1 < len(args); i++ {
			switch args[i] {
			case "-r":
				rev = args[i+1]
			case "-T":
				template = args[i+1]
			}
		}
		if template == "empty" {
			return "false\n", nil
		}
		if template == "commit_id" {
			if rev == "@-" {
				return "parent-commit\n", nil
			}
			return "current-commit\n", nil
		}
		if template == `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"` {
			if rev == "@-" {
				return "parent-change|parent-commit|parent|base\n", nil
			}
			return "current-change|current-commit|current|parent-change\n", nil
		}
	case "diff":
		if len(args) >= 4 && args[1] == "-r" && args[2] == "@-" && args[3] == "--name-only" {
			return "leftover.go\n", nil
		}
		return "", nil
	}
	return "", nil
}

func (r *prepareNextDemuxFakeRunner) called(want string) bool {
	for _, call := range r.calls {
		if call == want || strings.HasPrefix(call, want+" ") {
			return true
		}
	}
	return false
}

type demuxApplySourceFakeRunner struct {
	repoRoot string
	calls    []string
}

func (r *demuxApplySourceFakeRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	return r.output(name, args)
}

func (r *demuxApplySourceFakeRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	r.calls = append(r.calls, call)
	return r.output(name, args)
}

func (r *demuxApplySourceFakeRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

func (r *demuxApplySourceFakeRunner) output(name string, args []string) (string, error) {
	if name != "jj" || len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "edit":
		return "", nil
	case "log":
		rev := ""
		template := ""
		for i := 0; i+1 < len(args); i++ {
			switch args[i] {
			case "-r":
				rev = args[i+1]
			case "-T":
				template = args[i+1]
			}
		}
		if template == "empty" {
			if rev == "@" {
				return "true\n", nil
			}
			return "false\n", nil
		}
		if template == `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"` {
			if rev == "@-" {
				return "parent-change|parent-commit|compose source|base-change\n", nil
			}
			return "empty-child|empty-commit||parent-change\n", nil
		}
	case "diff":
		if len(args) >= 4 && args[1] == "-r" && args[2] == "@-" && args[3] == "--name-only" {
			return "server/src/app.ts\n", nil
		}
		return "", nil
	}
	return "", nil
}

func (r *demuxApplySourceFakeRunner) called(want string) bool {
	for _, call := range r.calls {
		if call == want || strings.HasPrefix(call, want+" ") {
			return true
		}
	}
	return false
}
