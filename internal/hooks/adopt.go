package hooks

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/vcs"
)

// AdoptPushOptions configures auto-adopt publication for a raw git push.
type AdoptPushOptions struct {
	RepoRoot string
	Remote   string
	LocalRef string
	HeadSHA  string
	RefRange string
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
		Backend:    "git",
		BranchName: branchPtr,
	}
	if gxRepo, err := vcs.NewService().ResolveGXRepoAtPath(ctx, repoRoot); err == nil {
		repo = gxRepo
		repo.BranchName = branchPtr
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
	commits, err := vcs.PushedCommitsInGitRange(ctx, repoRoot, opts.RefRange)
	if err != nil {
		return result, err
	}
	result.Commits = commits
	return result, nil
}

func branchNameFromRef(localRef string) string {
	localRef = strings.TrimSpace(localRef)
	switch {
	case strings.HasPrefix(localRef, "refs/heads/"):
		return strings.TrimPrefix(localRef, "refs/heads/")
	case localRef == "HEAD":
		// `git push origin HEAD` reports the local ref as literal HEAD; a
		// HEAD-named artifact never unifies with the PR's bookmark on the
		// server, so fall back to resolving the current branch.
		return ""
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
