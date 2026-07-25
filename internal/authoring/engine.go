package authoring

import (
	"context"
	"os"

	"github.com/satoricorp/gx/internal/vcs"
)

type RepoInfo = vcs.RepoInfo
type ChangeInfo = vcs.ChangeInfo
type InitOptions = vcs.InitOptions
type InitResult = vcs.InitResult
type StackInfo = vcs.StackInfo
type SyncResult = vcs.SyncResult
type RevisionSummary = vcs.RevisionSummary
type StackSummary = vcs.StackSummary
type StatusSnapshot = vcs.StatusSnapshot
type GitWorkingStatus = vcs.GitWorkingStatus
type PruneEmptyStacksResult = vcs.PruneEmptyStacksResult
type PruneGitHubPullRequestsResult = vcs.PruneGitHubPullRequestsResult

// Engine is GX's authoring seam. CLI and MCP adapters call this module
// instead of owning Git/storage mechanics directly.
type Engine struct {
	vcs *vcs.Service
}

func NewEngine() *Engine {
	return NewEngineWithVCS(vcs.NewService())
}

func NewEngineWithVCS(service *vcs.Service) *Engine {
	return &Engine{vcs: service}
}

func (e *Engine) Init(ctx context.Context, opts InitOptions) (InitResult, error) {
	return e.vcs.InitWithOptions(ctx, opts)
}

type EnsureReadyResult = vcs.EnsureReadyResult

func (e *Engine) EnsureReadyRepo(ctx context.Context) (EnsureReadyResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return EnsureReadyResult{}, err
	}
	return e.vcs.EnsureReadyRepo(ctx, cwd)
}

func (e *Engine) Status(ctx context.Context) (StackSummary, error) {
	return e.vcs.Stack(ctx)
}

func (e *Engine) StatusSnapshot(ctx context.Context) (StatusSnapshot, error) {
	return e.vcs.StatusSnapshot(ctx)
}

func (e *Engine) PreservingGitIndex(ctx context.Context, fn func() error) error {
	return e.vcs.PreservingGitIndexForCwd(ctx, fn)
}

func (e *Engine) GitWorkingStatus(ctx context.Context, repoRoot string) (GitWorkingStatus, error) {
	return e.vcs.GitWorkingStatus(ctx, repoRoot)
}

func (e *Engine) CurrentChange(ctx context.Context, repoRoot, rev string) (ChangeInfo, error) {
	return e.vcs.CurrentChange(ctx, repoRoot, rev)
}

func (e *Engine) PruneEmptyStacks(ctx context.Context) (PruneEmptyStacksResult, error) {
	return e.vcs.PruneEmptyStacks(ctx)
}

func (e *Engine) Sync(ctx context.Context, remote string) (SyncResult, error) {
	return e.vcs.Sync(ctx, remote)
}

func (e *Engine) SyncCloudBookmarkTip(ctx context.Context, repo RepoInfo, branchName string) error {
	return e.vcs.SyncCloudBookmarkTip(ctx, repo, branchName)
}

func (e *Engine) PrunePublishedStackByRef(ctx context.Context, repo RepoInfo, publishRef string) (bool, error) {
	return e.vcs.PrunePublishedStackByRef(ctx, repo, publishRef)
}

func (e *Engine) PruneTerminalGitHubPullRequestStacks(ctx context.Context, repo RepoInfo) (PruneGitHubPullRequestsResult, error) {
	return e.vcs.PruneTerminalGitHubPullRequestStacks(ctx, repo)
}
