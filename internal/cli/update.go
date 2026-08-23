package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/update"
	"github.com/satoricorp/gx/internal/version"
)

func newUpdateCommand(ctx context.Context) *cobra.Command {
	var checkOnly bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gx to the latest published build",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			// Never the cached manifest: someone typing `gx update` is asking
			// about now, and the cache exists to keep the passive notice off
			// the network, not to answer a direct question.
			manifest, err := update.FetchManifest(ctx, nil)
			if err != nil {
				return err
			}
			info := version.BuildInfo()

			if !update.Behind(info, manifest) {
				fmt.Fprintln(out, success("gx is up to date")+muted(" ("+version.Current()+")"))
				return nil
			}
			if checkOnly {
				fmt.Fprintln(out, update.Notice(manifest))
				return nil
			}
			// A local build's commit is not the published one, so Behind says
			// "behind" for work that is usually ahead. The passive notice skips
			// these builds entirely; here the user asked, so say what will
			// happen and do it.
			if update.Skip(info) {
				fmt.Fprintln(out, muted("Replacing a development build ("+version.Current()+") with the published one."))
			}

			if err := update.Apply(ctx, manifest, nil, func(line string) {
				fmt.Fprintln(out, muted(line))
			}); err != nil {
				return err
			}
			name := manifest.Version
			if name == "" {
				name = manifest.GitSHA
			}
			fmt.Fprintln(out, success("Updated gx to "+name))
			return nil
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "report whether a newer build exists without installing it")
	return cmd
}

// autoUpdateTimeout bounds the download.
//
// It runs after the command's own output, so the caller already has what they
// asked for and this delays nothing they were waiting on — but it still holds
// the shell until it finishes, and a courtesy has no business doing that for
// minutes. Past it, the notice is printed and the next command tries again.
const autoUpdateTimeout = 90 * time.Second

// maybeSelfUpdate brings the binary up to date after a command has finished.
//
// After, not before: updating on the way in would mean every stale binary
// pauses to download tens of megabytes before doing the thing the user asked
// for, and it only pays off if the process then re-executes itself. Running it
// afterwards costs the caller nothing they are waiting on. The price is that
// the update lands for the *next* command, which is what the message says.
//
// It runs from PersistentPostRun, so it is reached only when the command
// itself succeeded — a failing command does not get a download stapled to it —
// and the root hook already excludes every `__`-prefixed internal command,
// which is how the git hooks invoke gx. A pre-push hook is the last place that
// should pause to fetch 30 MB.
//
// Replacing the binary does not disturb the process doing the replacing: the
// rename leaves this one running on the old inode.
//
// Output goes to stderr: `gx auth token`, `gx version --json` and the MCP path
// all have machine-readable stdout.
func maybeSelfUpdate(ctx context.Context, cmd *cobra.Command) {
	switch cmd.Name() {
	case "update", "version", "__complete", "__completeNoDesc":
		return
	}
	info := version.BuildInfo()
	if update.Skip(info) {
		return
	}
	checkCtx, cancelCheck := context.WithTimeout(ctx, 2*time.Second)
	defer cancelCheck()
	now := time.Now()
	manifest, err := update.Latest(checkCtx, nil, now)
	if err != nil || !update.Behind(info, manifest) {
		return
	}

	stderr := cmd.ErrOrStderr()

	// Nothing to do but tell them, either because they asked for that or
	// because an attempt failed recently and retrying on every command would
	// re-download tens of megabytes each time.
	if !update.AutoEnabled() || !update.AutoUpdateDue(now) {
		fmt.Fprintln(stderr, muted(update.Notice(manifest)))
		return
	}

	update.MarkAutoUpdateAttempt(now)
	applyCtx, cancelApply := context.WithTimeout(ctx, autoUpdateTimeout)
	defer cancelApply()
	if err := update.Apply(applyCtx, manifest, nil, nil); err != nil {
		// A failed self-update is not the user's problem to debug after some
		// unrelated command, but it does change the advice: the one-liner
		// still works when a permission or network fault stopped this.
		fmt.Fprintln(stderr, muted(update.Notice(manifest)))
		return
	}
	name := strings.TrimSpace(manifest.Version)
	if name == "" {
		name = manifest.GitSHA
	}
	fmt.Fprintln(stderr, muted(fmt.Sprintf(
		"Updated gx to %s — the next command uses it (this one ran %s).",
		name, version.Current())))
}
