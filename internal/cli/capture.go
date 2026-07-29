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

	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/claudehooks"
	"github.com/satoricorp/totality/internal/capture/extract"
	"github.com/satoricorp/totality/internal/capture/orchestrator"
	"github.com/satoricorp/totality/internal/capture/reparse"
	"github.com/satoricorp/totality/internal/hooks"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/telemetry"
	"github.com/satoricorp/totality/internal/uploadauth"
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
		Short: "Upload shareable capture staging rows to the Totality server",
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
		return fmt.Errorf("not logged in for capture upload — run `tl auth login`")
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

// runPushHook is a seam for tests: the real hook runner needs a git repo and a
// live staging database, neither of which is the subject of the command's own
// error reporting.
var runPushHook = hooks.RunPush

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
			outcome, err := runPushHook(ctx, hooks.PushOptions{
				RepoRoot: repoRoot,
				Remote:   remote,
				RefRange: refRange,
				Base:     base,
				Head:     head,
				LocalRef: localRef,
				HeadSHA:  headSHA,
				Tools:    toolList,
			})
			// The command still exits 0 — it runs inside a git hook and must
			// never block a push — but a failure is printed rather than
			// discarded. Swallowing it left an empty `extract=` as the only
			// evidence that anything went wrong, which reads exactly like a
			// clean run that found nothing.
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: capture push failed: %v\n", err)
			}
			result := outcome.Result
			// A run that deliberately did nothing says so. Printing the
			// staging line here would emit the byte-for-byte signature of a
			// failed run (`extract=` empty) for a paused capture or a repo
			// that opted out, which is how one line came to mean four things.
			if outcome.SkipReason != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "capture skipped: %s\n", outcome.SkipReason)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "capture staged extract=%s sessions=%d coverage=%.1f%% ref=%s shareable=%d/%d\n",
					result.StagedExtractID,
					result.StagedSessions,
					result.HunkCoverage*100,
					result.RefRange,
					outcome.ShareableExtract,
					outcome.ShareableSession,
				)
				// Zero coverage with sessions on disk and work in the range is
				// not a hand-written change; it is capture reading a
				// transcript the agent had not finished writing. Say so, or
				// the staging line reads as a verdict on the change.
				if result.AttributionSuspect {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Warning", fmt.Sprintf(
						"found agent sessions but linked none of %d changed hunk(s); session data was likely still being written. Re-run once the session settles: tl capture push --repo %s --ref-range %s",
						result.EligibleHunks, repoRoot, result.RefRange,
					)))
				}
			}
			if outcome.Publication.Queued {
				fmt.Fprintf(cmd.OutOrStdout(), "publication queued id=%s\n", outcome.Publication.QueueID)
			}
			if outcome.CaptureError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: capture failed, so this push has no session context: %s\n", outcome.CaptureError)
			}
			if outcome.ShareableError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: capture rows staged but not marked uploadable: %s\n", outcome.ShareableError)
			}
			if outcome.RecoveryError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: repair Totality revision metadata: %s\n", outcome.RecoveryError)
			}
			if outcome.AttachError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: this push's review artifact will list no sessions: %s\n", outcome.AttachError)
			}
			if outcome.PublicationError != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: this push queued no Totality review artifact: %s\n", outcome.PublicationError)
			}
			// Uploads happen in a detached process whose output goes nowhere,
			// so this is the only place a failing upload can reach a human.
			if failures := outcome.UploadFailures; failures.Total() > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: %d capture upload(s) failed (%d gave up after %d attempts): %s\n",
					failures.Total(), failures.ExhaustedRows, storage.MaxUploadAttempts, failures.LastError)
			}
			// One tool being unavailable is survivable but never silent: it
			// means this push carries less session context than it could.
			for _, problem := range result.DiscoveryProblems {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: capture skipped a tool: %s\n", problem)
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
	tlPath, err := os.Executable()
	if err != nil {
		return err
	}
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repoRoot, TotalityPath: tlPath}); err != nil {
		return err
	}
	if printSuccess {
		fmt.Fprintln(cmd.OutOrStdout(), labelValue("Totality lifecycle hooks", success("ok")+": commit metadata and push capture"))
	}
	return nil
}

func installClaudeCaptureHooks(cmd *cobra.Command, repoRoot string, quiet bool) error {
	tlPath, err := os.Executable()
	if err != nil {
		return err
	}
	if err := claudehooks.MergeSettings(repoRoot, tlPath); err != nil {
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
