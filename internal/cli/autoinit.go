package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
)

func ensureAutoInitializedRepo(ctx context.Context, engine *authoring.Engine, cmd *cobra.Command) error {
	if shouldSkipAutoInit(cmd) {
		return nil
	}
	var result authoring.EnsureReadyResult
	// Guard the index: autoinit runs jj, and jj rewrites the git index on
	// snapshot/import, which would silently discard the user's git add state.
	err := engine.PreservingGitIndex(ctx, func() error {
		var readyErr error
		result, readyErr = engine.EnsureReadyRepo(ctx)
		return readyErr
	})
	if err != nil {
		return err
	}
	if !result.Prepared {
		return nil
	}
	if !commandRequestsJSON(cmd) {
		fmt.Fprintln(cmd.ErrOrStderr(), muted("Initializing gx for this repository..."))
	}
	if result.Repo.RootPath != "" {
		_ = installCaptureHookQuiet(cmd, result.Repo.RootPath)
	}
	return nil
}

func commandRequestsJSON(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if current.Flags().Lookup("json") == nil {
			continue
		}
		jsonOut, err := current.Flags().GetBool("json")
		if err == nil && jsonOut {
			return true
		}
	}
	return false
}

func shouldSkipAutoInit(cmd *cobra.Command) bool {
	if cmd == nil {
		return true
	}
	if strings.Contains(cmd.CommandPath(), "__") {
		return true
	}
	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		// commit skips autoinit because autoinit runs jj, and jj rewrites the
		// git index (staged entries become intent-to-add) before gx commit can
		// capture them; RecordStagedRevision self-initializes after capture.
		case "init", "version", "login", "auth", "set", "demo", "commit":
			return true
		}
	}
	return false
}
