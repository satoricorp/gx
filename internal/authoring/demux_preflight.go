package authoring

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/satoricorp/gx/internal/storage"
)

var demuxApplyPreflightEnvMu sync.Mutex

type demuxApplyPreflightState struct {
	Repo        *storage.Repo
	Changes     []storage.Change
	Stacks      []storage.Stack
	StackChange map[int64][]storage.Change
}

func (e *Engine) PreflightDemuxApply(ctx context.Context, proposal DemuxProposal) error {
	if strings.TrimSpace(os.Getenv("GX_COMPOSE_SKIP_APPLY_PREFLIGHT")) == "1" {
		return nil
	}
	if strings.TrimSpace(proposal.RepoRoot) == "" {
		return fmt.Errorf("proposal repo root is empty")
	}
	state, err := loadDemuxApplyPreflightState(ctx, proposal.RepoRoot)
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "gx-compose-attempt-")
	if err != nil {
		return fmt.Errorf("create compose attempt dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	attemptRoot := filepath.Join(tmp, "repo")
	if err := copyTreeForDemuxPreflight(proposal, attemptRoot); err != nil {
		return fmt.Errorf("copy repo into compose attempt dir: %w", err)
	}
	attemptGXHome := filepath.Join(tmp, "gx-home")
	attemptProposal := proposal
	attemptProposal.RepoRoot = attemptRoot
	attemptProposal.Status = ProposalPending

	return withTemporaryGXHomeAndDir(attemptGXHome, attemptRoot, func() error {
		if err := seedDemuxApplyPreflightState(ctx, state, attemptRoot); err != nil {
			return err
		}
		attemptEngine := NewEngine()
		saved, err := attemptEngine.SaveDemuxProposal(ctx, attemptProposal)
		if err != nil {
			return fmt.Errorf("save compose attempt proposal: %w", err)
		}
		if _, err := attemptEngine.ApplyDemuxProposalWithOptions(ctx, saved.ID, ApplyDemuxOptions{AllowWarnings: true}); err != nil {
			return fmt.Errorf("apply compose attempt proposal: %w", err)
		}
		return nil
	})
}

func copyTreeForDemuxPreflight(proposal DemuxProposal, dst string) error {
	protected := demuxPreflightProtectedPaths(proposal)
	ignored := demuxPreflightIgnoredPaths(proposal.RepoRoot)
	return copyTreeWithSkips(proposal.RepoRoot, dst, protected, func(rel string, entry fs.DirEntry) bool {
		if entry.IsDir() && demuxPreflightSkippableDir(rel) {
			return true
		}
		return demuxPreflightIgnoredPath(rel, ignored)
	})
}

func demuxPreflightProtectedPaths(proposal DemuxProposal) map[string]struct{} {
	protected := map[string]struct{}{}
	add := func(path string) {
		path = filepath.Clean(strings.TrimSpace(path))
		if path == "." || path == "" || strings.HasPrefix(path, "..") || filepath.IsAbs(path) {
			return
		}
		protected[path] = struct{}{}
	}
	for _, hunk := range proposal.Hunks {
		add(hunk.File)
	}
	for _, revision := range proposal.Revisions {
		for _, file := range revision.Files {
			add(file)
		}
		for _, hunk := range revision.Hunks {
			add(hunk.File)
		}
	}
	return protected
}

func demuxPreflightSkippableDir(rel string) bool {
	switch filepath.Base(rel) {
	case "node_modules", ".next", ".turbo", ".cache", ".parcel-cache", ".vite",
		"dist", "build", "coverage", "target", "tmp", "temp", ".venv":
		return true
	default:
		return false
	}
}

func demuxPreflightIgnoredPaths(repoRoot string) map[string]struct{} {
	out, err := exec.Command("git", "-C", repoRoot, "ls-files", "--others", "--ignored", "--exclude-standard", "-z").Output()
	if err != nil || len(out) == 0 {
		return nil
	}
	ignored := map[string]struct{}{}
	for _, raw := range strings.Split(string(out), "\x00") {
		rel := filepath.Clean(strings.TrimSpace(raw))
		if rel == "." || rel == "" || demuxPreflightVCSInternalPath(rel) {
			continue
		}
		ignored[rel] = struct{}{}
	}
	return ignored
}

func demuxPreflightIgnoredPath(rel string, ignored map[string]struct{}) bool {
	if len(ignored) == 0 {
		return false
	}
	rel = filepath.Clean(strings.TrimSpace(rel))
	if rel == "." || rel == "" || demuxPreflightVCSInternalPath(rel) {
		return false
	}
	_, ok := ignored[rel]
	return ok
}

func demuxPreflightVCSInternalPath(rel string) bool {
	rel = filepath.ToSlash(filepath.Clean(strings.TrimSpace(rel)))
	return rel == ".git" || strings.HasPrefix(rel, ".git/") || rel == ".jj" || strings.HasPrefix(rel, ".jj/")
}

func loadDemuxApplyPreflightState(ctx context.Context, repoRoot string) (demuxApplyPreflightState, error) {
	store, err := openStore(ctx)
	if err != nil {
		return demuxApplyPreflightState{}, err
	}
	defer store.Close()

	repo, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return demuxApplyPreflightState{}, err
	}
	if repo == nil {
		return demuxApplyPreflightState{StackChange: map[int64][]storage.Change{}}, nil
	}
	changes, err := store.ListChangesByRepoID(ctx, repo.ID)
	if err != nil {
		return demuxApplyPreflightState{}, err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repo.ID)
	if err != nil {
		return demuxApplyPreflightState{}, err
	}
	stackChanges := make(map[int64][]storage.Change, len(stacks))
	for _, stack := range stacks {
		entries, err := store.ListChangesByStackID(ctx, stack.ID)
		if err != nil {
			return demuxApplyPreflightState{}, err
		}
		stackChanges[stack.ID] = entries
	}
	return demuxApplyPreflightState{
		Repo:        repo,
		Changes:     changes,
		Stacks:      stacks,
		StackChange: stackChanges,
	}, nil
}

func seedDemuxApplyPreflightState(ctx context.Context, state demuxApplyPreflightState, attemptRoot string) error {
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()

	nowRepo := storage.Repo{
		RootPath: attemptRoot,
		Backend:  "jj",
	}
	if state.Repo != nil {
		nowRepo = *state.Repo
		nowRepo.ID = 0
		nowRepo.RootPath = attemptRoot
	}
	repoID, err := store.UpsertRepo(ctx, nowRepo)
	if err != nil {
		return err
	}
	changeIDs := make(map[int64]int64, len(state.Changes))
	for _, change := range state.Changes {
		originalID := change.ID
		change.ID = 0
		change.RepoID = repoID
		id, err := store.UpsertChange(ctx, change)
		if err != nil {
			return err
		}
		changeIDs[originalID] = id
	}
	for _, stack := range state.Stacks {
		originalID := stack.ID
		stack.ID = 0
		stack.RepoID = repoID
		stackID, err := store.UpsertStack(ctx, stack)
		if err != nil {
			return err
		}
		for _, change := range state.StackChange[originalID] {
			changeID, ok := changeIDs[change.ID]
			if !ok {
				copy := change
				copy.ID = 0
				copy.RepoID = repoID
				changeID, err = store.UpsertChange(ctx, copy)
				if err != nil {
					return err
				}
				changeIDs[change.ID] = changeID
			}
			if err := store.AddChangeToStack(ctx, stackID, changeID, change.FirstSeenAt); err != nil {
				return err
			}
		}
	}
	return nil
}

func withTemporaryGXHomeAndDir(gxHome, workDir string, fn func() error) error {
	demuxApplyPreflightEnvMu.Lock()
	defer demuxApplyPreflightEnvMu.Unlock()

	previousDir, err := os.Getwd()
	if err != nil {
		return err
	}
	previous, hadPrevious := os.LookupEnv("GX_HOME")
	if err := os.Setenv("GX_HOME", gxHome); err != nil {
		return err
	}
	if err := os.Chdir(workDir); err != nil {
		if hadPrevious {
			_ = os.Setenv("GX_HOME", previous)
		} else {
			_ = os.Unsetenv("GX_HOME")
		}
		return err
	}
	defer func() {
		_ = os.Chdir(previousDir)
		if hadPrevious {
			_ = os.Setenv("GX_HOME", previous)
		} else {
			_ = os.Unsetenv("GX_HOME")
		}
	}()
	return fn()
}

func copyTree(src, dst string) error {
	return copyTreeWithSkips(src, dst, nil, nil)
}

func copyTreeWithSkips(src, dst string, protected map[string]struct{}, skipEntry func(string, fs.DirEntry) bool) error {
	src = filepath.Clean(src)
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		mode := info.Mode()
		switch {
		case rel == ".":
			return os.MkdirAll(target, mode.Perm())
		case mode.Type()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		case entry.IsDir():
			if skipEntry != nil && skipEntry(rel, entry) && !hasProtectedPathUnder(protected, rel) {
				return filepath.SkipDir
			}
			return os.MkdirAll(target, mode.Perm())
		case mode.IsRegular():
			if skipEntry != nil && skipEntry(rel, entry) && !isProtectedPath(protected, rel) {
				return nil
			}
			return copyFile(path, target, mode.Perm())
		default:
			return nil
		}
	})
}

func isProtectedPath(protected map[string]struct{}, rel string) bool {
	if len(protected) == 0 {
		return false
	}
	_, ok := protected[filepath.Clean(rel)]
	return ok
}

func hasProtectedPathUnder(protected map[string]struct{}, rel string) bool {
	if len(protected) == 0 {
		return false
	}
	rel = filepath.Clean(rel)
	for path := range protected {
		if path == rel || strings.HasPrefix(path, rel+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
