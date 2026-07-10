package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/capture/extract"
	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/uploadauth"
)

func newCaptureCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "capture",
		Short:  "Capture agent sessions for pushed commits",
		Hidden: true,
	}
	cmd.AddCommand(newCapturePushCommand(ctx))
	cmd.AddCommand(newCaptureSyncCommand(ctx))
	return cmd
}

func newCaptureSyncCommand(ctx context.Context) *cobra.Command {
	var quiet bool
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Upload shareable capture staging rows to the GX server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCaptureSync(ctx, cmd.OutOrStdout(), quiet)
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "suppress upload summary")
	return cmd
}

func runCaptureSync(ctx context.Context, out interface{ Write([]byte) (int, error) }, quiet bool) error {
	creds, ok := uploadauth.Load()
	if !ok {
		if quiet {
			return nil
		}
		return fmt.Errorf("not logged in for capture upload — run `gx auth login`")
	}
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		return err
	}
	result, err := extract.SyncPending(ctx, stager, creds, telemetry.NewFromEnv())
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Fprintf(out, "capture sync: extracts=%d sessions=%d errors=%d/%d\n",
			result.ExtractsUploaded,
			result.SessionsUploaded,
			result.ExtractErrors+result.SessionErrors,
			result.ExtractErrors+result.SessionErrors+result.ExtractsUploaded+result.SessionsUploaded,
		)
	}
	return nil
}

func newCapturePushCommand(ctx context.Context) *cobra.Command {
	var (
		repoRoot string
		remote   string
		refRange string
		base     string
		head     string
		localRef string
		headSHA  string
		tools    string
	)
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Run capture pipeline for a pushed ref range",
		RunE: func(cmd *cobra.Command, args []string) error {
			toolList := splitCaptureTools(tools)
			outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
				RepoRoot: repoRoot,
				Remote:   remote,
				RefRange: refRange,
				Base:     base,
				Head:     head,
				LocalRef: localRef,
				HeadSHA:  headSHA,
				Tools:    toolList,
			})
			if err != nil {
				return nil
			}
			result := outcome.Result
			fmt.Fprintf(cmd.OutOrStdout(), "capture staged extract=%s sessions=%d coverage=%.1f%% ref=%s shareable=%d/%d\n",
				result.StagedExtractID,
				result.StagedSessions,
				result.HunkCoverage*100,
				result.RefRange,
				outcome.ShareableExtract,
				outcome.ShareableSession,
			)
			if outcome.Publication.Queued {
				fmt.Fprintf(cmd.OutOrStdout(), "publication queued id=%s\n", outcome.Publication.QueueID)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root (default: current directory)")
	cmd.Flags().StringVar(&remote, "remote", "origin", "git remote name")
	cmd.Flags().StringVar(&refRange, "ref-range", "", "git ref range (e.g. abc..def)")
	cmd.Flags().StringVar(&base, "base", "", "base ref when ref-range omitted")
	cmd.Flags().StringVar(&head, "head", "", "head ref when ref-range omitted")
	cmd.Flags().StringVar(&localRef, "local-ref", "", "local ref name from pre-push stdin")
	cmd.Flags().StringVar(&headSHA, "head-sha", "", "local commit SHA from pre-push stdin")
	cmd.Flags().StringVar(&tools, "tools", "claude,codex,cursor", "comma-separated capture tools")
	return cmd
}

func splitCaptureTools(raw string) []string {
	var out []string
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func installCaptureHook(cmd *cobra.Command, repoRoot string) error {
	return installCaptureHookWithOutput(cmd, repoRoot, true)
}

func installCaptureHookQuiet(cmd *cobra.Command, repoRoot string) error {
	return installCaptureHookWithOutput(cmd, repoRoot, false)
}

func installCaptureHookWithOutput(cmd *cobra.Command, repoRoot string, printSuccess bool) error {
	if repoRoot == "" {
		var err error
		repoRoot, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	gxPath, err := os.Executable()
	if err != nil {
		return err
	}
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repoRoot, GXPath: gxPath}); err != nil {
		return err
	}
	if printSuccess {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Pre-push hook", success("ok")+": capture on git push"))
	}
	return nil
}
