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
	err := engine.PreservingGitIndex(ctx, func() error {
		var readyErr error
		result, readyErr = engine.EnsureReadyRepo(ctx)
		return readyErr
	})
	if err != nil {
		return err
	}
	if result.Prepared && !commandRequestsJSON(cmd) {
		fmt.Fprintln(cmd.ErrOrStderr(), muted("Initializing gx for this repository..."))
	}
	if result.Repo.RootPath != "" {
		if hookErr := installCaptureHookQuiet(cmd, result.Repo.RootPath); hookErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: GX lifecycle hooks not installed: %v\n", hookErr)
		}
		// Existing installs may still run the retired ambient-capture
		// LaunchAgent; retire it the next time gx touches an initialized repo.
		cleanupLegacyAmbientCaptureQuiet(ctx, cmd.ErrOrStderr())
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
		// `review` is read-only: it must work as a CI gate and on someone
		// else's checkout without installing hooks or writing GX state into a
		// repo the reviewer does not own.
		case "init", "version", "login", "auth", "set", "review":
			return true
		}
	}
	return false
}
