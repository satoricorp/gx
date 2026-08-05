package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/lgtm/internal/authoring"
	"github.com/satoricorp/lgtm/internal/telemetry"
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
		fmt.Fprintln(cmd.ErrOrStderr(), muted("Initializing lgtm for this repository..."))
	}
	if result.Repo.RootPath != "" {
		if hookErr := installCaptureHookQuiet(cmd, result.Repo.RootPath); hookErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: lgtm lifecycle hooks not installed: %v\n", hookErr)
		}
		// Existing installs may still run the retired ambient-capture
		// LaunchAgent; retire it the next time lgtm touches an initialized repo.
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

// commandTelemetryContext marks read-only commands so telemetry reports
// without writing anything to the machine.
func commandTelemetryContext(ctx context.Context, cmd *cobra.Command) context.Context {
	if commandMustNotWriteLgtmState(cmd) {
		return telemetry.WithoutStateWrites(ctx)
	}
	return ctx
}

// commandMustNotWriteLgtmState names the commands that promise to leave no lgtm
// state behind. This is narrower than shouldSkipAutoInit, which also exempts
// commands like `init` and `login` whose whole job is to write lgtm state — they
// skip auto-init because they set it up themselves, not because they must not.
func commandMustNotWriteLgtmState(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		// `review` must work as a CI gate and on someone else's checkout
		// without leaving lgtm state on the machine running it.
		case "review":
			return true
		// `version` answers one question about the binary. Dockerfiles and CI
		// steps run it to check what they installed, and minting a machine ID
		// and an install sentinel to answer it means `lgtm version` creates
		// $LGTM_HOME on a machine that has not yet decided to use lgtm. The
		// install event is not worth that; the first command that actually
		// does something records it.
		case "version":
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
		// else's checkout without installing hooks or writing lgtm state into a
		// repo the reviewer does not own.
		case "init", "version", "login", "auth", "review":
			return true
		}
	}
	return false
}
