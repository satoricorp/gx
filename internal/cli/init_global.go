package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/satoricorp/totality/internal/hooks"
)

// runGlobalInit installs the machine-wide Totality lifecycle hooks. It does not need
// to run inside a git repository.
func runGlobalInit(ctx context.Context, cmd *cobra.Command, quiet bool) error {
	result, err := hooks.InstallGlobal(ctx, hooks.GlobalInstallOptions{})
	if err != nil {
		return err
	}
	if quiet {
		return nil
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, labelValue("Global hooks", success("ok")+": "+result.HooksDir))
	if result.ConfigUpdated {
		fmt.Fprintln(out, labelValue("Git config", fmt.Sprintf("core.hooksPath = %s (global)", result.HooksDir)))
	} else {
		fmt.Fprintln(out, labelValue("Git config", muted("core.hooksPath already pointed at Totality")))
	}
	fmt.Fprintln(out, labelValue("Scope", fmt.Sprintf("every git repo on this machine; %d hooks installed, repo hooks still run", len(result.Scripts))))
	fmt.Fprintln(out)
	// Labels stay short so they fit the shared label column without wrapping.
	fmt.Fprintln(out, labelValue("Opt out", "git config tl.enabled false          (one repo)"))
	fmt.Fprintln(out, labelValue("Undo", "git config --global --unset core.hooksPath"))
	return nil
}
