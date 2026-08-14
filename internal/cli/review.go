package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/spf13/cobra"
)

// `gx review` is the pre-ship exit gate. Unlike enhance, where --fail-on opts a
// gate in, the coded exits are the default here — review IS the gate — and
// --report-only opts out for hosts that read the report instead of the exit
// code (the MCP tool does exactly that). The exit values are shared with
// enhance, so CI authors learn one vocabulary:
//
//	0  ship: every gate passed or was legitimately skipped
//	3  no-ship: at least one gate failed        (reviewFindingsExitCode)
//	4  nothing to check                          (reviewNothingToReviewExitCode)
//	5  degraded: the AI gates could not run      (reviewDegradedExitCode)
func newReviewCommand(ctx context.Context) *cobra.Command {
	var base string
	var jsonOut bool
	var markdownOut bool
	var verbose bool
	var skipGates string
	var timeout time.Duration
	var reportOnly bool
	var clientOverride string
	cmd := &cobra.Command{
		Use:     "review [intent]",
		Aliases: []string{"gxr"},
		Short:   "Pre-ship exit gate: check the current change against six gates",
		Long: "Check the current change against six gates — correctness, security,\n" +
			"code health, back-pressure, accessibility, performance — reporting PASS, FAIL,\n" +
			"or SKIPPED for each plus a ship / no-ship verdict. Deterministic checks (tests,\n" +
			"linters, secret scan, dependency audits) decide what they can; one AI judgment\n" +
			"call covers the rest, grounded in the same context sources gx enhance uses.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Store-free and state-write-free, like enhance: a pre-ship gate may
			// run on checkouts and machines that never ran `gx init`.
			ctx := telemetry.WithoutStateWrites(ctx)
			startedAt := time.Now()
			client := telemetry.ClientSurface(clientOverride)
			skip, err := codereview.ParseGateIDs(skipGates)
			if err != nil {
				return err
			}
			intent := ""
			if len(args) > 0 {
				intent = strings.TrimSpace(args[0])
			}
			repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
			if err != nil {
				emitReviewRunTelemetry(ctx, codereview.ReviewReport{}, err, client, intent, time.Since(startedAt))
				return err
			}
			var progress io.Writer
			if !jsonOut {
				progress = cmd.ErrOrStderr()
			}
			report, err := codereview.CheckReview(ctx, repo.RootPath, codereview.ReviewOptions{
				Intent:         intent,
				Base:           base,
				SkipGates:      skip,
				Timeout:        timeout,
				Verbose:        verbose,
				ProgressWriter: progress,
				Color:          !jsonOut && !markdownOut,
			})
			emitReviewRunTelemetry(ctx, report, err, client, intent, time.Since(startedAt))
			if err != nil {
				return err
			}
			switch {
			case jsonOut:
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(report); err != nil {
					return err
				}
			case markdownOut:
				fmt.Fprint(cmd.OutOrStdout(), codereview.RenderReviewMarkdown(report))
			default:
				fmt.Fprint(cmd.OutOrStdout(), codereview.RenderReviewText(report))
			}
			if reportOnly {
				return nil
			}
			return gateVerdictError(report)
		},
	}
	cmd.Flags().StringVar(&base, "base", "", "check the commit range <ref>...HEAD instead of the working tree, e.g. --base origin/main")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print the review report as JSON instead of text")
	cmd.Flags().BoolVar(&markdownOut, "md", false, "print the review report as markdown (what the MCP tool returns to agents)")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "list the source files behind each gate")
	cmd.Flags().StringVar(&skipGates, "skip-gates", "", "comma-separated gates to skip: "+strings.Join(reviewGateNames(), ", "))
	cmd.Flags().DurationVar(&timeout, "timeout", 0, "whole-run budget (default 5m; per-tool caps still apply)")
	cmd.Flags().BoolVar(&reportOnly, "report-only", false, fmt.Sprintf("always exit 0; the verdict lives in the report only (without this, no-ship exits %d, nothing-to-check %d, degraded %d)", reviewFindingsExitCode, reviewNothingToReviewExitCode, reviewDegradedExitCode))
	cmd.Flags().StringVar(&clientOverride, "client", "", "surface invoking this run, overriding $GX_CLIENT: cli, mcp, skill, slash-enhance, slash-review")
	return cmd
}

func reviewGateNames() []string {
	ids := codereview.AllGateIDs()
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		names = append(names, string(id))
	}
	return names
}

// gateVerdictError maps the verdict onto the exit-code contract.
func gateVerdictError(report codereview.ReviewReport) error {
	switch report.Verdict {
	case codereview.VerdictNoShip:
		failed := report.FailedGates()
		names := make([]string, 0, len(failed))
		for _, id := range failed {
			names = append(names, string(id))
		}
		return vcs.CodedErrorf(reviewFindingsExitCode, fmt.Errorf("gx review: %d gate(s) failed (%s)", len(failed), strings.Join(names, ", ")))
	case codereview.VerdictNothingToCheck:
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return vcs.CodedErrorf(reviewNothingToReviewExitCode, fmt.Errorf("gx review: nothing to check (looked at %s); refusing to pass a gate without inspecting any code", target))
	case codereview.VerdictDegraded:
		return vcs.CodedErrorf(reviewDegradedExitCode, fmt.Errorf("gx review: %s; refusing to pass a gate on an incomplete run", strings.Join(report.DegradedReasons, "; ")))
	default:
		return nil
	}
}

func emitReviewRunTelemetry(ctx context.Context, report codereview.ReviewReport, runErr error, client, intent string, duration time.Duration) {
	status := "success"
	if runErr != nil {
		status = "error"
	}
	props := map[string]any{
		"status":             status,
		"client":             client,
		"has_intent":         strings.TrimSpace(intent) != "",
		"duration_ms":        duration.Milliseconds(),
		"changed_file_count": len(report.ChangedFiles),
		"degraded":           len(report.DegradedReasons) > 0,
	}
	if strings.TrimSpace(report.Verdict) != "" {
		props["verdict"] = report.Verdict
		props["reviewed"] = report.Reviewed
	}
	for _, gate := range report.Gates {
		props["gate_"+strings.ReplaceAll(string(gate.Gate), "-", "_")] = strings.ToLower(string(gate.Status))
	}
	if strings.TrimSpace(report.AIModel) != "" {
		props["ai_model"] = report.AIModel
	}
	if strings.TrimSpace(report.AITransport) != "" {
		props["ai_transport"] = report.AITransport
	}
	telemetry.EmitProductEvent(ctx, telemetry.EventCLIReviewRun, props)
}
