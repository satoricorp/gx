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
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Upload pending capture staging rows to the GX server",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, ok := uploadauth.Load()
			if !ok {
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
			fmt.Fprintf(cmd.OutOrStdout(), "capture sync: extracts=%d sessions=%d errors=%d/%d\n",
				result.ExtractsUploaded,
				result.SessionsUploaded,
				result.ExtractErrors+result.SessionErrors,
				result.ExtractErrors+result.SessionErrors+result.ExtractsUploaded+result.SessionsUploaded,
			)
			return nil
		},
	}
	return cmd
}

func newCapturePushCommand(ctx context.Context) *cobra.Command {
	var (
		repoRoot string
		remote   string
		refRange string
		base     string
		head     string
		tools    string
	)
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Run capture pipeline for a pushed ref range",
		RunE: func(cmd *cobra.Command, args []string) error {
			toolList := splitCaptureTools(tools)
			result, err := hooks.RunPush(ctx, hooks.PushOptions{
				RepoRoot: repoRoot,
				Remote:   remote,
				RefRange: refRange,
				Base:     base,
				Head:     head,
				Tools:    toolList,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "capture staged extract=%s sessions=%d coverage=%.1f%% ref=%s\n",
				result.StagedExtractID,
				result.StagedSessions,
				result.HunkCoverage*100,
				result.RefRange,
			)
			return nil
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root (default: current directory)")
	cmd.Flags().StringVar(&remote, "remote", "origin", "git remote name")
	cmd.Flags().StringVar(&refRange, "ref-range", "", "git ref range (e.g. abc..def)")
	cmd.Flags().StringVar(&base, "base", "", "base ref when ref-range omitted")
	cmd.Flags().StringVar(&head, "head", "", "head ref when ref-range omitted")
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
	fmt.Fprintln(cmd.OutOrStdout(), labelValue("Pre-push hook", success("ok")+": capture on git push"))
	return nil
}
