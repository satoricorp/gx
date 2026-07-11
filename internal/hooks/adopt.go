package hooks

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"

	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

// AdoptPushOptions configures auto-adopt publication for a raw git push.
type AdoptPushOptions struct {
	RepoRoot    string
	Remote      string
	LocalRef    string
	HeadSHA     string
	RevisionIDs []string
}

// EnqueueAdoptedPublication queues a PR-summary artifact for a plain git push.
func EnqueueAdoptedPublication(ctx context.Context, opts AdoptPushOptions) (publication.Result, error) {
	push, err := buildAdoptPushResult(ctx, opts)
	if err != nil {
		return publication.Result{}, err
	}
	if push.HeadCommitID == "" {
		return publication.Result{}, nil
	}
	return publication.EnqueuePush(ctx, push)
}

func buildAdoptPushResult(ctx context.Context, opts AdoptPushOptions) (vcs.PushResult, error) {
	repoRoot := strings.TrimSpace(opts.RepoRoot)
	if repoRoot == "" {
		return vcs.PushResult{}, fmt.Errorf("repo root required")
	}
	headSHA := strings.TrimSpace(opts.HeadSHA)
	if headSHA == "" {
		var err error
		headSHA, err = gitRevParse(ctx, repoRoot, "HEAD")
		if err != nil {
			return vcs.PushResult{}, err
		}
	}
	remote := strings.TrimSpace(opts.Remote)
	var remotePtr *string
	if remote != "" {
		remotePtr = &remote
	}
	branch := branchNameFromRef(opts.LocalRef)
	if branch == "" {
		// Hooks installed by older gx versions do not pass --local-ref;
		// without a branch the artifact routes to an "unknown" bookmark on
		// the server and PR linkage is lost.
		branch, _ = gitCurrentBranch(ctx, repoRoot)
	}
	var branchPtr *string
	if branch != "" {
		branchPtr = &branch
	}
	repo := vcs.RepoInfo{
		RootPath:   repoRoot,
		Backend:    "jj",
		BranchName: branchPtr,
	}
	if remotePtr != nil {
		repo.DefaultRemote = remotePtr
	}
	if url, err := gitRemoteURL(ctx, repoRoot, remote); err == nil && url != "" {
		repo.RemoteURL = &url
	}
	if base, err := gitSymbolicRef(ctx, repoRoot, "refs/remotes/"+remote+"/HEAD"); err == nil {
		base = strings.TrimPrefix(base, "refs/remotes/"+remote+"/")
		if base != "" {
			repo.DefaultBranch = &base
		}
	}

	result := vcs.PushResult{
		Repo:         repo,
		HeadCommitID: headSHA,
		RemoteName:   remotePtr,
	}
	revisionIDs := uniqueStrings(opts.RevisionIDs)
	if len(revisionIDs) == 0 {
		return result, nil
	}
	changes, err := lookupAdoptedChanges(ctx, repoRoot, revisionIDs)
	if err != nil {
		return result, err
	}
	if len(changes) == 0 {
		return result, nil
	}
	result.CurrentChange = &changes[len(changes)-1]
	result.Published = make([]vcs.PushedChange, 0, len(changes))
	for _, change := range changes {
		entry := vcs.PushedChange{Change: change}
		if branch != "" {
			entry.BranchName = branch
		}
		result.Published = append(result.Published, entry)
	}
	return result, nil
}

func lookupAdoptedChanges(ctx context.Context, repoRoot string, revisionIDs []string) ([]vcs.ChangeInfo, error) {
	db, err := storage.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var out []vcs.ChangeInfo
	for _, revisionID := range revisionIDs {
		var change vcs.ChangeInfo
		var parent sql.NullString
		err := db.QueryRowContext(ctx, `
			SELECT c.jj_change_id, c.current_commit_id, c.description, c.parent_change_id
			FROM repos r
			JOIN changes c ON c.repo_id = r.id
			WHERE r.root_path = ? AND c.jj_change_id = ?
			ORDER BY c.updated_at DESC
			LIMIT 1
		`, repoRoot, revisionID).Scan(&change.ChangeID, &change.CommitID, &change.Description, &parent)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("lookup revision %s: %w", revisionID, err)
		}
		if parent.Valid {
			change.ParentChangeID = &parent.String
		}
		out = append(out, change)
	}
	return out, nil
}

func branchNameFromRef(localRef string) string {
	localRef = strings.TrimSpace(localRef)
	switch {
	case strings.HasPrefix(localRef, "refs/heads/"):
		return strings.TrimPrefix(localRef, "refs/heads/")
	default:
		return localRef
	}
}

func gitCurrentBranch(ctx context.Context, repoRoot string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "branch", "--show-current")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch --show-current: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func gitRevParse(ctx context.Context, repoRoot, ref string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", ref)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse %s: %w", ref, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func gitRemoteURL(ctx context.Context, repoRoot, remote string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "remote", "get-url", remote)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitSymbolicRef(ctx context.Context, repoRoot, ref string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "symbolic-ref", ref)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
