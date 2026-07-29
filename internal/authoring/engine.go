package authoring

import (
	"context"
	"os"

	"github.com/satoricorp/totality/internal/vcs"
)

type InitOptions = vcs.InitOptions
type InitResult = vcs.InitResult

// Engine is Totality's authoring seam. CLI and MCP adapters call this module
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

func (e *Engine) PreservingGitIndex(ctx context.Context, fn func() error) error {
	return e.vcs.PreservingGitIndexForCwd(ctx, fn)
}
