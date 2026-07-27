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
type RevisionSummary = vcs.RevisionSummary
type StackSummary = vcs.StackSummary
type GitWorkingStatus = vcs.GitWorkingStatus
type PruneEmptyStacksResult = vcs.PruneEmptyStacksResult

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
