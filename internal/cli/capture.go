package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/claudehooks"
	"github.com/satoricorp/gx/internal/capture/extract"
	"github.com/satoricorp/gx/internal/capture/orchestrator"
	"github.com/satoricorp/gx/internal/capture/reparse"
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
	cmd.AddCommand(newCaptureReparseCommand(ctx))
	cmd.AddCommand(newCaptureTranscriptCommand(ctx))
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
			if outcome.RecoveryError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: repair GX revision metadata: %s\n", outcome.RecoveryError)
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
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("GX lifecycle hooks", success("ok")+": commit metadata and push capture"))
	}
	return nil
}

func installClaudeCaptureHooks(cmd *cobra.Command, repoRoot string, quiet bool) error {
	gxPath, err := os.Executable()
	if err != nil {
		return err
	}
	if err := claudehooks.MergeSettings(repoRoot, gxPath); err != nil {
		return err
	}
	if !quiet {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Claude hooks", success("ok")+": Stop/SessionEnd transcript capture"))
	}
	return nil
}

func newCaptureReparseCommand(ctx context.Context) *cobra.Command {
	var repoRoot string
	cmd := &cobra.Command{
		Use:   "reparse",
		Short: "Re-run normalization over stored raw session blobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if repoRoot == "" {
				var err error
				repoRoot, err = os.Getwd()
				if err != nil {
					return err
				}
			}
			stager, err := storage.OpenCaptureStager(ctx)
			if err != nil {
				return err
			}
			result, err := reparse.Run(ctx, stager, repoRoot, telemetry.NewFromEnv())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "capture reparse: sessions=%d updated=%d errors=%d\n",
				result.Sessions, result.Updated, result.Errors)
			return nil
		},
	}
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root (default: current directory)")
	return cmd
}

func newCaptureTranscriptCommand(ctx context.Context) *cobra.Command {
	var (
		path     string
		repoRoot string
		tool     string
	)
	cmd := &cobra.Command{
		Use:   "transcript",
		Short: "Ingest a Claude/Codex/Cursor transcript path from hooks or --path",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				resolved, err := transcriptPathFromStdin(cmd.InOrStdin())
				if err != nil {
					return err
				}
				path = resolved
			}
			if strings.TrimSpace(path) == "" {
				return nil
			}
			if repoRoot == "" {
				var err error
				repoRoot, err = os.Getwd()
				if err != nil {
					return err
				}
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if tool == "" {
				tool = inferTranscriptTool(path)
			}
			stager, err := storage.OpenCaptureStager(ctx)
			if err != nil {
				return err
			}
			session, badLines, err := orchestrator.IngestRawSession(ctx, stager, tool, path, raw, repoRoot, nil)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "capture transcript: tool=%s session=%s bad_lines=%d events=%d\n",
				tool, session.SessionID, badLines, len(session.Events))
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "transcript JSONL path")
	cmd.Flags().StringVar(&repoRoot, "repo", "", "repository root (default: current directory)")
	cmd.Flags().StringVar(&tool, "tool", "", "agent tool (claude, codex, cursor)")
	return cmd
}

func transcriptPathFromStdin(in io.Reader) (string, error) {
	if in == nil {
		return "", nil
	}
	data, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}
	payload := strings.TrimSpace(string(data))
	if payload == "" {
		return "", nil
	}
	var hook map[string]json.RawMessage
	if err := json.Unmarshal([]byte(payload), &hook); err != nil {
		return "", nil
	}
	var path string
	if err := json.Unmarshal(hook["transcript_path"], &path); err != nil {
		return "", nil
	}
	return strings.TrimSpace(path), nil
}

func inferTranscriptTool(path string) string {
	lower := strings.ToLower(filepath.ToSlash(path))
	switch {
	case strings.Contains(lower, "/.claude/"):
		return capture.ToolClaude
	case strings.Contains(lower, "/.codex/"):
		return capture.ToolCodex
	case strings.Contains(lower, "/.cursor/"):
		return capture.ToolCursor
	default:
		return capture.ToolClaude
	}
}
