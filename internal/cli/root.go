package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/capture/extract"
	"github.com/satoricorp/gx/internal/clitui"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/daemon"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/launcher"
	"github.com/satoricorp/gx/internal/postlist"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/uploadauth"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/satoricorp/gx/internal/version"
)

var publishRevisionOutputLine = regexp.MustCompile(`(?m)^(\s{2})([0-9a-f]{4,})(\s+)(.*)$`)

func NewRoot(ctx context.Context) *cobra.Command {
	engine := authoring.NewEngine()

	root := &cobra.Command{
		Use:           "gx",
		Short:         "GX CLI for capturing coding context and shipping stacked pull requests",
		Long:          gxTagline,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			printInitNoteIfNeeded(cmd)
			return printRootHelp(cmd)
		},
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	root.AddGroup(
		&cobra.Group{ID: groupSetup, Title: "Setup:"},
		&cobra.Group{ID: groupWork, Title: "Work:"},
		&cobra.Group{ID: groupShip, Title: "Ship:"},
		&cobra.Group{ID: groupAdvanced, Title: "Advanced:"},
	)

	root.AddCommand(
		newInternalDaemonCommand(ctx),
		newVersionCommand(),
		newDoctorCommand(ctx),
		newLoginCommand(ctx),
		newAuthCommand(ctx),
		newInitCommand(ctx, engine),
		newCaptureCommand(ctx),
		newPublishUploadCommand(ctx),
		newBaseCommand(ctx, engine),
		newDemoCommand(),
		newComposeCommand(ctx, engine),
		newAddCommand(ctx, engine, "add [filesets...]", "add", true),
		newEditCommand(ctx, engine),
		newStatusCommand(ctx, engine, "status", "status", false),
		newStacksCommand(ctx, engine),
		newReviewCommand(ctx),
		newPublishCommand(ctx, engine, "publish [stack]", false),
		newSyncCommand(ctx, engine),
		newOpsCommand(ctx),
	)
	assignCommandGroups(root)
	rootRef := root
	root.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c != rootRef {
			if err := renderHelpTemplate(c, c.OutOrStdout()); err != nil {
				c.PrintErrln(err)
			}
			return
		}
		if err := printRootHelp(c); err != nil {
			c.PrintErrln(err)
		}
	})
	installHelpStyling(root)

	return root
}

func assignCommandGroups(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		switch cmd.Name() {
		case "init", "auth", "login", "version", "demo", "doctor":
			cmd.GroupID = groupSetup
		case "add", "base", "compose", "edit", "review", "status", "stacks":
			cmd.GroupID = groupWork
		case "publish", "sync":
			cmd.GroupID = groupShip
		case "ops":
			cmd.GroupID = groupAdvanced
		}
	}
}

func printInitNoteIfNeeded(cmd *cobra.Command) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	cfg, err := gxconfig.LoadAt(cwd)
	if err != nil || cfg.HasIdentity() {
		return
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, danger("Run `gx init` first."))
	fmt.Fprintln(out)
}

func ShouldLaunch(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] == "help" {
		return false
	}
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help":
			return false
		}
	}
	return launcher.IsWrappedTool(args[0])
}

func newInternalDaemonCommand(ctx context.Context) *cobra.Command {
	var ambient bool
	cmd := &cobra.Command{
		Use:    "__gx-daemon",
		Short:  "Run the internal gx background service",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return daemon.Run(ctx, daemon.Options{Ambient: ambient})
		},
	}
	cmd.Flags().BoolVar(&ambient, "ambient", false, "run stable ambient proxy")
	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print gx version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), logoText("version "+version.Current()))
		},
	}
}

func newInitCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var name string
	var email string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Set up GX in the current repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), commandLine("gx init", true))
			fmt.Fprintln(cmd.OutOrStdout())
			result, err := engine.Init(ctx, authoring.InitOptions{
				Name:        name,
				Email:       email,
				Interactive: true,
				In:          cmd.InOrStdin(),
				Out:         cmd.OutOrStdout(),
			})
			if err != nil {
				return err
			}
			if result.IdentityName != "" && result.IdentityEmail != "" {
				if result.IdentityApplied {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Configured", fmt.Sprintf("gx identity as %s <%s>", result.IdentityName, result.IdentityEmail)))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Using", fmt.Sprintf("gx identity %s <%s>", result.IdentityName, result.IdentityEmail)))
				}
				if client := postlist.NewFromEnv(); client != nil {
					err := client.UpsertIdentity(ctx, postlist.Identity{
						Name:      result.IdentityName,
						Email:     result.IdentityEmail,
						GXVersion: version.Current(),
					})
					if err != nil {
						fmt.Fprintln(cmd.ErrOrStderr(), danger(fmt.Sprintf("gx signup upload failed: %v", err)))
					}
				}
			}
			if err := installCaptureHook(cmd, result.Repo.RootPath); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Pre-push hook", err.Error()))
			}
			if err := gxconfig.EnsureCaptureRetention(); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Retention", err.Error()))
			} else {
				cfg, _ := gxconfig.Load()
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Session retention", fmt.Sprintf("%d days", cfg.Capture.CleanupPeriodDays)))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "user name to store in GX and JJ config")
	cmd.Flags().StringVar(&email, "email", "", "user email to store in GX and JJ config")
	return cmd
}

func newAddCommand(ctx context.Context, engine *authoring.Engine, use string, helpName string, hidden bool) *cobra.Command {
	var message string
	var interactive bool
	var hunk bool
	var patchFile string
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   use,
		Short: "Create the next GX change from your current work",
		Long: strings.Join([]string{
			"Create the next GX change from your current work.",
			"",
			"A non-empty commit message is required; pass it with -m.",
			"",
			"From gx/<base>, gx add records the current work into a new inferred GX stack based on <base>.",
			"From gx/edit, gx add is edit-mode surgery: it records onto the active edited stack/revision and should only run after intentionally choosing that stack.",
		}, "\n"),
		Example: strings.Join([]string{
			`  gx add -m "describe this revision"`,
			`  gx add -i internal/termstyle/theme.go -m "update terminal theme"`,
			`  gx add --hunk --patch-file /tmp/selected.patch -m "record selected hunks"`,
		}, "\n"),
		Hidden: hidden,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := engine.Checkpoint(ctx, authoring.CheckpointOptions{
				Intent:      message,
				Filesets:    args,
				Interactive: interactive,
				Hunk:        hunk,
				PatchFile:   patchFile,
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, addResultJSON{
					Result: result,
					Split:  interactive || len(args) > 0 || hunk,
				})
			}
			if strings.TrimSpace(result.Output) != "" {
				fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(result.Output))
			}
			printAddSummary(cmd.OutOrStdout(), result, interactive || len(args) > 0)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message (required)")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactively choose changes to add")
	cmd.Flags().BoolVar(&hunk, "hunk", false, "select changes from a patch file (agent-friendly)")
	cmd.Flags().StringVar(&patchFile, "patch-file", "", "path to a unified patch file used with --hunk")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	setHelpName(cmd, helpName)
	return cmd
}

type addResultJSON struct {
	Result authoring.CheckpointResult `json:"result"`
	Split  bool                       `json:"split"`
}

func printAddSummary(out io.Writer, result authoring.CheckpointResult, split bool) {
	invocation := "gx add"
	if result.Change.Description != "" {
		invocation = fmt.Sprintf("gx add -m %q", result.Change.Description)
	}
	fmt.Fprintln(out, commandLine(invocation, false))
	fmt.Fprintln(out)
	fmt.Fprintln(out, success("Revision recorded"))
	fmt.Fprintln(out, labelValue("Message", result.Change.Description))
	fmt.Fprintln(out, labelValue("Revision", shortID(result.Change.ChangeID, 12)))
	if result.Change.CommitID != "" {
		fmt.Fprintln(out, labelValue("Commit", shortID(result.Change.CommitID, 8)))
	}
	if split {
		fmt.Fprintln(out, labelValue("Split", "recorded the selected changes; remaining edits stay in the current revision"))
	} else {
		fmt.Fprintln(out, labelValue("Split", "recorded the current revision and opened the next one"))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, hint("Remaining changes stay in the current revision."))
	fmt.Fprintln(out, labelValue("Edit", fmt.Sprintf("gx edit %s", result.Change.ChangeID)))
	fmt.Fprintln(out, labelValue("Next", fmt.Sprintf("gx add -m %q", result.Change.Description)))
}

func newBaseCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var setRef string
	var jsonOut bool
	cmd := &cobra.Command{
		Use:    "base",
		Short:  "Utility: show or set the GX authoring base",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				result authoring.BaseResult
				err    error
			)
			if strings.TrimSpace(setRef) != "" {
				result, err = engine.SetBase(ctx, setRef)
			} else {
				result, err = engine.Base(ctx)
			}
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			printBaseSummary(cmd.OutOrStdout(), result, strings.TrimSpace(setRef) != "")
			return nil
		},
	}
	cmd.Flags().StringVar(&setRef, "set", "", "set the GX authoring base ref")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

func printBaseSummary(out io.Writer, result authoring.BaseResult, updated bool) {
	if updated {
		fmt.Fprintln(out, success("Base updated"))
	} else {
		fmt.Fprintln(out, section("GX base"))
	}
	fmt.Fprintln(out, labelValue("Base", result.BaseRef))
	if result.DefaultBranch != "" && result.DefaultBranch != result.BaseRef {
		fmt.Fprintln(out, labelValue("Default branch", result.DefaultBranch))
	}
	current := result.CurrentRef
	if current == "" {
		current = "(detached)"
	}
	fmt.Fprintln(out, labelValue("Current", current))
	fmt.Fprintln(out, labelValue("On base", yesNo(result.OnBase)))
	if !result.OnBase {
		fmt.Fprintln(out)
		fmt.Fprintln(out, labelValue("Return", fmt.Sprintf("gx base --set %s", result.BaseRef)))
	}
}

func newDemuxCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var intent string
	var raw bool
	var planOnly bool
	var model string
	var maxWarnings int
	var excludeFilesets []string
	cmd := &cobra.Command{
		Use:    "demux [filesets...]",
		Short:  "Compose ordered GX revisions from the current working copy",
		Hidden: true,
		Args:   cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			run := func(progress io.Writer) (authoring.DemuxPlanPacket, error) {
				if err := engine.RequireAuthoringBase(ctx, "gx compose"); err != nil {
					return authoring.DemuxPlanPacket{}, err
				}
				return engine.DemuxChanges(ctx, authoring.ProposeDemuxOptions{
					Intent:          intent,
					Filesets:        args,
					ExcludeFilesets: excludeFilesets,
					PlanOnly:        planOnly,
					Model:           model,
					MaxWarnings:     maxWarnings,
					ProgressWriter:  progress,
				})
			}
			var (
				packet authoring.DemuxPlanPacket
				err    error
			)
			if !jsonOut && useDemuxLoader(cmd.InOrStdin(), cmd.ErrOrStderr()) {
				packet, err = runDemuxWithLoader(cmd.InOrStdin(), cmd.ErrOrStderr(), run)
			} else {
				var progress io.Writer
				if !jsonOut {
					progress = cmd.ErrOrStderr()
				}
				packet, err = run(progress)
			}
			if err != nil {
				if strings.TrimSpace(packet.Proposal.ID) != "" {
					if jsonOut {
						_ = writeJSON(cmd, map[string]any{
							"error":  err.Error(),
							"packet": packet,
						})
						return err
					}
					printDemuxChangesPacket(cmd.OutOrStdout(), packet, demuxPrintOptions{Raw: raw})
				}
				return err
			}
			if jsonOut {
				return writeJSON(cmd, packet)
			}
			if !raw && !planOnly && useStatusInteractive(cmd.InOrStdin(), cmd.OutOrStdout()) {
				action, err := runDemuxInteractive(cmd.InOrStdin(), cmd.OutOrStdout(), packet.Proposal)
				if err != nil {
					return err
				}
				result, applied, err := applyDemuxInteractiveAction(ctx, engine, packet, action, cmd.Name() == "compose")
				if err != nil {
					return err
				}
				if applied {
					printDemuxApplySummary(cmd.OutOrStdout(), result)
				}
				return nil
			}
			printDemuxChangesPacket(cmd.OutOrStdout(), packet, demuxPrintOptions{Raw: raw})
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().StringVar(&intent, "intent", "", "optional intent prefix for proposed revisions")
	cmd.Flags().BoolVar(&raw, "raw", false, "show full compose diagnostics in human output")
	cmd.Flags().BoolVar(&planOnly, "plan", false, "stop after local deterministic planning; do not call OpenAI")
	cmd.Flags().StringVar(&model, "model", "", "OpenAI model for compose AI repair; defaults to GX_DEMUX_REVIEW_MODEL or gpt-4.1-mini")
	cmd.Flags().IntVar(&maxWarnings, "max-warnings", 0, "maximum warning-severity diagnostics to send to the model; defaults to 50")
	cmd.Flags().StringArrayVar(&excludeFilesets, "exclude", nil, "exclude changed file or directory from compose; may be repeated")
	cmd.AddCommand(newDemuxListCommand(ctx, engine))
	cmd.AddCommand(newDemuxReviewCommand(ctx, engine))
	cmd.AddCommand(newDemuxFixCommand(ctx, engine))
	cmd.AddCommand(newDemuxShowCommand(ctx, engine))
	cmd.AddCommand(newDemuxApplyCommand(ctx, engine))
	cmd.AddCommand(newDemuxApplyPlanCommand(ctx, engine))
	cmd.AddCommand(newDemuxReviewPlanCommand(ctx, engine))
	cmd.AddCommand(newComposeDoctorCommand(ctx, engine))
	return cmd
}

func shouldRunDemuxInteractiveOnError(packet authoring.DemuxPlanPacket, raw bool, planOnly bool, in io.Reader, out io.Writer) bool {
	return demuxErrorHasInteractiveProposal(packet, raw, planOnly) && useStatusInteractive(in, out)
}

func demuxErrorHasInteractiveProposal(packet authoring.DemuxPlanPacket, raw bool, planOnly bool) bool {
	return strings.TrimSpace(packet.Proposal.ID) != "" && !raw && !planOnly
}

const demuxPartialComposeWarning = "Partial compose proposal: these revisions do not cover all current changes. Accept them if they look right, then run gx compose again for the remaining changes."

func saveDemuxPacketProposal(ctx context.Context, engine *authoring.Engine, packet authoring.DemuxPlanPacket) (authoring.DemuxPlanPacket, error) {
	saved, err := engine.SaveDemuxProposal(ctx, packet.Proposal)
	if err != nil {
		return authoring.DemuxPlanPacket{}, err
	}
	packet.Proposal = saved
	packet.Review.Proposal = saved
	packet.StructuralFacts = saved.StructuralFacts
	packet.StructuralDeps = saved.StructuralDeps
	packet.ChangedSymbols = saved.ChangedSymbols
	packet.FeasibilityWarnings = saved.FeasibilityWarnings
	packet.Warnings = saved.Warnings
	return packet, nil
}

func newComposeDoctorCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check jj provisioning, capture parsers, and compose matcher coverage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := engine.RunComposeDoctor(ctx)
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, report)
			}
			fmt.Fprintln(cmd.OutOrStdout(), section("Compose"))
			if report.JJOK {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("jj", success("ok")+": "+report.JJVersion))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("jj", danger("fail")+": not available"))
			}
			if report.JJRepoOK {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("jj repo", success("ok")))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("jj repo", danger("fail")+": not colocated"))
			}
			if len(report.ParsersReachable) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Capture parsers", success("ok")+": "+strings.Join(report.ParsersReachable, ", ")))
			}
			if len(report.ParsersMissing) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Missing parsers", danger("warn")+": "+strings.Join(report.ParsersMissing, ", ")))
			}
			if report.LastCoverageFound {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Last hunk coverage", fmt.Sprintf("%.0f%%", report.LastHunkCoverage*100)))
			}
			if report.FeasibilityHint != "" {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Feasibility", muted(report.FeasibilityHint)))
			}
			for _, warning := range report.Warnings {
				fmt.Fprintln(cmd.OutOrStdout(), labelWarningValue("Warning", warning))
			}
			if !report.OK {
				return fmt.Errorf("compose doctor found issues")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

func newComposeCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	cmd := newDemuxCommand(ctx, engine)
	cmd.Use = "compose [filesets...]"
	cmd.Short = "Compose ordered GX revisions from the current working copy"
	cmd.Hidden = false
	configureComposeApplyDefaults(cmd)
	renameDemuxCommandSurface(cmd)
	return cmd
}

func configureComposeApplyDefaults(cmd *cobra.Command) {
	for _, sub := range cmd.Commands() {
		switch sub.Name() {
		case "apply", "apply-plan":
			sub.Annotations = map[string]string{"returnToDefaultBranch": "true"}
		}
	}
}

func renameDemuxCommandSurface(cmd *cobra.Command) {
	for _, sub := range cmd.Commands() {
		sub.Short = strings.ReplaceAll(sub.Short, "demux", "compose")
		sub.Short = strings.ReplaceAll(sub.Short, "Demux", "Compose")
	}
}

func newDemuxListCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var includeApplied bool
	var limit int
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"proposals"},
		Short:   "List saved compose proposals for the current repo",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			proposals, err := engine.ListDemuxProposals(ctx, authoring.ListDemuxProposalsOptions{
				IncludeApplied: includeApplied,
				Limit:          limit,
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, proposals)
			}
			printDemuxProposalList(cmd.OutOrStdout(), proposals)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&includeApplied, "all", false, "include applied and expired proposals")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum proposals to list; use 0 for no limit")
	return cmd
}

func newDemuxReviewCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var raw bool
	cmd := &cobra.Command{
		Use:   "review <proposal-id>",
		Short: "Review a saved compose proposal for deterministic warnings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packet, err := engine.ShowDemuxProposal(ctx, args[0])
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, packet)
			}
			printDemuxReviewPacket(cmd.OutOrStdout(), packet, demuxPrintOptions{Raw: raw})
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&raw, "raw", false, "show full review diagnostics in human output")
	return cmd
}

func newDemuxFixCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var model string
	var maxWarnings int
	var planOnly bool
	cmd := &cobra.Command{
		Use:   "fix <proposal-id>",
		Short: "Repair a compose proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !jsonOut {
				mode := "deterministic repair, then AI if still needed"
				if planOnly {
					mode = "deterministic repair only"
				}
				fmt.Fprintf(cmd.ErrOrStderr(), "Repairing compose proposal %s (%s)...\n", args[0], mode)
			}
			result, err := engine.ReviewDemuxProposalWithAI(ctx, args[0], authoring.DemuxAIReviewOptions{
				Model:       model,
				MaxWarnings: maxWarnings,
				PlanOnly:    planOnly,
			})
			if err != nil {
				if jsonOut {
					_ = writeJSON(cmd, map[string]any{
						"error":  err.Error(),
						"result": result,
					})
					return err
				}
				if result.Proposal.ID != "" {
					printDemuxAIReviewResult(cmd.OutOrStdout(), result)
				}
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			printDemuxAIReviewResult(cmd.OutOrStdout(), result)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&planOnly, "plan", false, "stop after local deterministic repair; do not call OpenAI")
	cmd.Flags().StringVar(&model, "model", "", "OpenAI model for AI compose fix; defaults to GX_DEMUX_REVIEW_MODEL or gpt-4.1-mini")
	cmd.Flags().IntVar(&maxWarnings, "max-warnings", 0, "maximum warning-severity diagnostics to send to the model; defaults to 50")
	return cmd
}

func newDemuxShowCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var raw bool
	var proposalID string
	cmd := &cobra.Command{
		Use:   "show <proposal-or-revision-id>",
		Short: "Show a saved compose proposal or proposed revision",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(proposalID) == "" && demuxShowTargetIsProposal(args[0]) {
				packet, err := engine.ShowDemuxProposal(ctx, args[0])
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(cmd, packet)
				}
				printDemuxChangesPacket(cmd.OutOrStdout(), packet, demuxPrintOptions{Raw: raw})
				return nil
			}
			view, err := engine.ShowDemuxRevision(ctx, proposalID, args[0])
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, view)
			}
			printDemuxRevisionView(cmd.OutOrStdout(), view, demuxPrintOptions{Raw: raw})
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&raw, "raw", false, "show full revision diagnostics in human output")
	cmd.Flags().StringVar(&proposalID, "proposal", "", "compose proposal id; defaults to latest pending proposal in the current repo")
	return cmd
}

func demuxShowTargetIsProposal(target string) bool {
	target = strings.TrimSpace(strings.ToLower(target))
	if strings.HasPrefix(target, "demux-") {
		return true
	}
	if !strings.HasPrefix(target, "d") || len(target) < 2 {
		return false
	}
	_, err := strconv.Atoi(strings.TrimPrefix(target, "d"))
	return err == nil
}

func newDemuxApplyCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var allowWarnings bool
	cmd := &cobra.Command{
		Use:    "apply <proposal-id>",
		Short:  "Apply a pending GX compose proposal",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := engine.ApplyDemuxProposalWithOptions(ctx, args[0], authoring.ApplyDemuxOptions{
				AllowWarnings:         allowWarnings,
				ReturnToDefaultBranch: composeReturnToDefaultBranch(cmd),
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			printDemuxApplySummary(cmd.OutOrStdout(), result)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&allowWarnings, "allow-warnings", false, "apply despite blocking feasibility warnings")
	return cmd
}

func printDemuxApplySummary(out io.Writer, result authoring.ApplyDemuxResult) {
	fmt.Fprintln(out, success("Compose applied"))
	fmt.Fprintln(out, labelValue("Proposal", result.Proposal.ID))
	fmt.Fprintln(out, labelValue("Created", fmt.Sprintf("%d revisions", len(result.Revisions))))
	for index, revision := range result.Revisions {
		fmt.Fprintf(out, "  %s %s  %s\n", command(fmt.Sprintf("r%d", index+1)), value(revision.Change.Description), muted(shortID(revision.Change.ChangeID, 12)))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Next"))
	fmt.Fprintf(out, "  %s\n", command("gx stacks"))
	fmt.Fprintf(out, "  %s\n", command("gx publish"))
	if result.RemainingChanges {
		printDemuxPartialFollowup(out)
	}
}

func printDemuxPartialFollowup(out io.Writer) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Remaining changes"))
	fmt.Fprintln(out, muted("Selected compose changes were applied. Run "+command("gx compose")+" again for the remaining changes."))
}

func newDemuxReviewPlanCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var planFile string
	var raw bool
	cmd := &cobra.Command{
		Use:    "review-plan --plan-file <path>",
		Short:  "Review an LLM-authored GX revision plan without applying it",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(planFile) == "" {
				return fmt.Errorf("--plan-file is required")
			}
			data, err := os.ReadFile(planFile)
			if err != nil {
				return fmt.Errorf("read revision plan: %w", err)
			}
			var proposal authoring.DemuxProposal
			if err := json.Unmarshal(data, &proposal); err != nil {
				return fmt.Errorf("decode revision plan: %w", err)
			}
			result, err := engine.ReviewDemuxPlan(ctx, proposal)
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, section("Revision plan review"))
			fmt.Fprintln(out, labelValue("Valid", fmt.Sprintf("%t", result.Valid)))
			if len(result.Errors) > 0 {
				fmt.Fprintln(out)
				fmt.Fprintln(out, section("Errors"))
				for _, reviewErr := range result.Errors {
					fmt.Fprintf(out, "  %s %s\n", danger("!"), reviewErr)
				}
			}
			printDemuxProposal(out, result.Proposal, demuxPrintOptions{Raw: raw})
			return nil
		},
	}
	cmd.Flags().StringVar(&planFile, "plan-file", "", "path to a JSON compose proposal containing LLM-authored revisions")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&raw, "raw", false, "show full compose diagnostics in human output")
	return cmd
}

func newDemuxApplyPlanCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var planFile string
	var allowWarnings bool
	cmd := &cobra.Command{
		Use:    "apply-plan --plan-file <path>",
		Short:  "Apply an LLM-authored GX revision plan",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(planFile) == "" {
				return fmt.Errorf("--plan-file is required")
			}
			data, err := os.ReadFile(planFile)
			if err != nil {
				return fmt.Errorf("read revision plan: %w", err)
			}
			var proposal authoring.DemuxProposal
			if err := json.Unmarshal(data, &proposal); err != nil {
				return fmt.Errorf("decode revision plan: %w", err)
			}
			result, err := engine.ApplyDemuxPlanWithOptions(ctx, proposal, authoring.ApplyDemuxOptions{
				AllowWarnings:         allowWarnings,
				ReturnToDefaultBranch: composeReturnToDefaultBranch(cmd),
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, section("Revision plan applied"))
			fmt.Fprintln(out, labelValue("Proposal", result.Proposal.ID))
			fmt.Fprintln(out, labelValue("Revisions", fmt.Sprintf("%d", len(result.Revisions))))
			return nil
		},
	}
	cmd.Flags().StringVar(&planFile, "plan-file", "", "path to a JSON compose proposal containing LLM-authored revisions")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&allowWarnings, "allow-warnings", false, "apply despite blocking feasibility warnings")
	return cmd
}

func composeReturnToDefaultBranch(cmd *cobra.Command) bool {
	return cmd != nil && cmd.Annotations["returnToDefaultBranch"] == "true"
}

type demuxPrintOptions struct {
	Raw bool
}

func printDemuxChangesPacket(out io.Writer, packet authoring.DemuxPlanPacket, opts demuxPrintOptions) {
	fmt.Fprintln(out, commandLine("gx compose", true))
	fmt.Fprintln(out)
	printDemuxPartialWarning(out, packet.Proposal)
	printDemuxProposal(out, packet.Proposal, opts)
	printDemuxAcceptHint(out, packet)
}

func printDemuxPartialWarning(out io.Writer, proposal authoring.DemuxProposal) {
	if !demuxProposalIsPartial(proposal) {
		return
	}
	fmt.Fprintln(out, section("Partial proposal"))
	fmt.Fprintln(out, muted(demuxPartialComposeWarning))
	fmt.Fprintln(out)
}

func demuxProposalIsPartial(proposal authoring.DemuxProposal) bool {
	for _, warning := range proposal.Warnings {
		if warning == demuxPartialComposeWarning {
			return true
		}
	}
	return false
}

func printDemuxAcceptHint(out io.Writer, packet authoring.DemuxPlanPacket) {
	if packet.State != authoring.DemuxWorkflowReadyToApply || strings.TrimSpace(packet.Proposal.ID) == "" {
		return
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, labelValue("Accept", fmt.Sprintf("gx compose apply %s", packet.Proposal.ID)))
	if demuxProposalIsPartial(packet.Proposal) {
		fmt.Fprintln(out, labelValue("Then", "gx compose"))
	}
}

func printDemuxReviewPacket(out io.Writer, packet authoring.DemuxPlanPacket, opts demuxPrintOptions) {
	proposal := packet.Proposal
	review := packet.Review
	fmt.Fprintln(out, commandLine("gx compose review "+proposal.ID, true))
	fmt.Fprintln(out)
	printDemuxProposalSummary(out, proposal, packet.State)
	if len(review.Errors) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Blocking warnings"))
		for _, reviewErr := range review.Errors {
			fmt.Fprintf(out, "  %s\n", muted(reviewErr))
		}
	}
	if len(review.RepairHints) > 0 {
		fmt.Fprintln(out)
		if opts.Raw {
			fmt.Fprintln(out, section("Repair hints"))
			for _, repairHint := range review.RepairHints {
				fmt.Fprintf(out, "  %s %s\n", muted("-"), formatRepairHint(repairHint))
			}
		} else {
			fmt.Fprintln(out, labelValue("Repair hints", fmt.Sprintf("%d hidden; use --raw or --json", len(review.RepairHints))))
		}
	}
	blockingWarnings, diagnosticWarnings := splitFeasibilityWarnings(proposal.FeasibilityWarnings)
	if len(blockingWarnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Blocking warnings"))
		for _, warning := range blockingWarnings {
			details := warning.Source
			if warning.RevisionID != "" {
				details += " / " + warning.RevisionID
			}
			fmt.Fprintf(out, "  %s  %s\n", muted(firstNonEmptyString(warning.RevisionID, warningMarker(warning.Severity))), muted(warning.Message+"  ["+details+"]"))
		}
	}
	if opts.Raw && len(diagnosticWarnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Diagnostics"))
		for _, warning := range diagnosticWarnings {
			details := warning.Source
			if warning.RevisionID != "" {
				details += " / " + warning.RevisionID
			}
			fmt.Fprintf(out, "  %s  %s\n", muted(firstNonEmptyString(warning.RevisionID, warningMarker(warning.Severity))), muted(warning.Message+"  ["+details+"]"))
		}
	}
	if opts.Raw && len(proposal.Warnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Diagnostics"))
		for _, warning := range proposal.Warnings {
			fmt.Fprintf(out, "  %s %s\n", danger("!"), warning)
		}
	} else if !opts.Raw && len(proposal.Warnings)+len(diagnosticWarnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, labelValue("Diagnostics", fmt.Sprintf("%d hidden; use --raw or --json", len(proposal.Warnings)+len(diagnosticWarnings))))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Next"))
	if packet.State == authoring.DemuxWorkflowReadyToApply {
		fmt.Fprintf(out, "  %s\n", command("gx compose"))
	} else {
		fmt.Fprintf(out, "  %s\n", command(fmt.Sprintf("gx compose fix %s", proposal.ID)))
	}
}

func printDemuxAIReviewResult(out io.Writer, result authoring.DemuxAIReviewResult) {
	fmt.Fprintln(out, commandLine("gx compose fix "+result.Proposal.ID, true))
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Compose fix"))
	fmt.Fprintln(out, labelValue("Proposal", result.Proposal.ID))
	fmt.Fprintln(out, labelValue("Model", result.Model))
	fmt.Fprintln(out, labelValue("Updated", yesNo(result.Updated)))
	fmt.Fprintln(out, labelValue("State", string(result.State)))
	fmt.Fprintln(out, labelValue("Valid", fmt.Sprintf("%t", result.Review.Valid)))
	fmt.Fprintln(out, labelValue("Revisions", fmt.Sprintf("%d", len(result.Proposal.Revisions))))
	fmt.Fprintln(out, labelValue("Blocking issues", fmt.Sprintf("%d", demuxBlockingIssueCount(result))))
	if diagnostics := demuxHiddenDiagnosticCount(result); diagnostics > 0 {
		fmt.Fprintln(out, labelValue("Diagnostics", fmt.Sprintf("%d hidden; use --raw or --json", diagnostics)))
	}
	if result.WarningsTotal > 0 || result.RepairHintTotal > 0 {
		fmt.Fprintln(out, labelValue("Sent", fmt.Sprintf("%d/%d warnings, %d/%d repair hints", result.WarningsSent, result.WarningsTotal, result.RepairHintSent, result.RepairHintTotal)))
	}
	if len(result.Notes) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Notes"))
		for _, note := range result.Notes {
			fmt.Fprintf(out, "  %s %s\n", muted("-"), note)
		}
	}
	if len(result.Review.Errors) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Errors"))
		for _, reviewErr := range result.Review.Errors {
			fmt.Fprintf(out, "  %s %s\n", danger("!"), reviewErr)
		}
	}
	fmt.Fprintln(out)
	if result.State == authoring.DemuxWorkflowReadyToApply {
		fmt.Fprintln(out, labelValue("Accept", "gx compose"))
	} else {
		fmt.Fprintln(out, labelValue("Review", fmt.Sprintf("gx compose review %s --raw", result.Proposal.ID)))
	}
}

func demuxBlockingIssueCount(result authoring.DemuxAIReviewResult) int {
	blocking := 0
	for _, warning := range result.Proposal.FeasibilityWarnings {
		if warning.Severity == "warning" {
			blocking++
		}
	}
	return blocking + len(result.Review.Errors)
}

func demuxHiddenDiagnosticCount(result authoring.DemuxAIReviewResult) int {
	return len(result.Proposal.FeasibilityWarnings) + len(result.Proposal.Warnings) + len(result.Review.RepairHints)
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func reviewFeasibilityWarnings(warnings []authoring.FeasibilityWarning, raw bool) ([]authoring.FeasibilityWarning, int) {
	if raw {
		return warnings, 0
	}
	return nil, len(warnings)
}

func splitFeasibilityWarnings(warnings []authoring.FeasibilityWarning) ([]authoring.FeasibilityWarning, []authoring.FeasibilityWarning) {
	var blocking []authoring.FeasibilityWarning
	var diagnostics []authoring.FeasibilityWarning
	for _, warning := range warnings {
		if strings.EqualFold(strings.TrimSpace(warning.Severity), "warning") {
			blocking = append(blocking, warning)
			continue
		}
		diagnostics = append(diagnostics, warning)
	}
	return blocking, diagnostics
}

func formatRepairHint(repairHint authoring.RepairHint) string {
	parts := []string{repairHint.Kind}
	if repairHint.RevisionID != "" {
		parts = append(parts, repairHint.RevisionID)
	}
	if repairHint.HunkID != "" {
		parts = append(parts, repairHint.HunkID)
	}
	if repairHint.File != "" {
		parts = append(parts, repairHint.File)
	}
	if repairHint.DependsOn != "" {
		parts = append(parts, "depends on "+repairHint.DependsOn)
	}
	if repairHint.Symbol != "" {
		parts = append(parts, repairHint.Symbol)
	}
	if repairHint.Suggestion != "" {
		parts = append(parts, repairHint.Suggestion)
	}
	return strings.Join(parts, " - ")
}

func printDemuxProposalList(out io.Writer, proposals []authoring.DemuxProposalSummary) {
	fmt.Fprintln(out, section("Compose proposals"))
	if len(proposals) == 0 {
		fmt.Fprintln(out, muted("No compose proposals found."))
		return
	}
	for _, proposal := range proposals {
		marker := " "
		if proposal.LatestPendingForShow {
			marker = "*"
		}
		details := []string{
			string(proposal.Status),
			fmt.Sprintf("%d revisions", proposal.RevisionCount),
			fmt.Sprintf("%d files", len(proposal.Files)),
		}
		if proposal.HiddenDiagnosticCount > 0 {
			details = append(details, fmt.Sprintf("%d diagnostics", proposal.HiddenDiagnosticCount))
		}
		id := proposal.ID
		if proposal.Alias != "" {
			id = fmt.Sprintf("%s  %s", proposal.Alias, muted(proposal.ID))
		}
		fmt.Fprintf(out, "%s %s  %s  %s\n", marker, id, muted(strings.Join(details, " / ")), muted(formatMillis(proposal.CreatedAt)))
		if proposal.FirstRevisionIntent != "" {
			fmt.Fprintf(out, "    %s\n", proposal.FirstRevisionIntent)
		}
		if proposal.LatestPendingForShow {
			fmt.Fprintf(out, "    %s\n", muted("default for revisions: gx compose show <revision-id>"))
		}
	}
}

func formatMillis(ms int64) string {
	if ms <= 0 {
		return "unknown time"
	}
	return time.UnixMilli(ms).Format("2006-01-02 15:04:05")
}

func printDemuxRevisionView(out io.Writer, view authoring.DemuxRevisionView, opts demuxPrintOptions) {
	revision := view.Revision
	fmt.Fprintln(out, commandLine("gx compose show "+revision.ID, true))
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Compose revision"))
	fmt.Fprintln(out, labelValue("Proposal", view.ProposalID))
	fmt.Fprintln(out, labelValue("Proposal status", string(view.ProposalStatus)))
	fmt.Fprintln(out, labelValue("Base revision", view.ProposedChangeID))
	fmt.Fprintln(out, labelValue("Revision", revision.ID))
	fmt.Fprintln(out, labelValue("Intent", revision.Intent))
	fmt.Fprintln(out, labelValue("Mode", demuxRevisionMode(revision)))
	if revision.Confidence > 0 {
		fmt.Fprintln(out, labelValue("Confidence", fmt.Sprintf("%.2f", revision.Confidence)))
	}
	if opts.Raw && revision.ProvenanceStatus != "" {
		provenance := revision.ProvenanceStatus
		if len(revision.SessionIDs) > 0 {
			provenance += " " + strings.Join(revision.SessionIDs, ", ")
		}
		fmt.Fprintln(out, labelValue("Provenance", provenance))
	}
	if len(revision.DependsOn) > 0 {
		fmt.Fprintln(out, labelValue("Depends on", strings.Join(revision.DependsOn, ", ")))
	}
	if len(revision.Files) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Files"))
		for _, file := range revision.Files {
			fmt.Fprintf(out, "  %s\n", file)
		}
	}
	if len(view.ChangedSymbols) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Changed symbols"))
		for _, symbol := range view.ChangedSymbols {
			label := symbol.Symbol
			if symbol.Kind != "" {
				label = symbol.Kind + " " + label
			}
			fmt.Fprintf(out, "  %s %s  %s\n", muted("-"), label, muted(symbol.File))
		}
	}
	if opts.Raw {
		printDemuxRevisionDiagnostics(out, view)
	} else {
		diagnostics := len(view.FeasibilityWarnings) + len(view.StructuralDeps)
		if diagnostics > 0 {
			fmt.Fprintln(out)
			fmt.Fprintln(out, labelValue("Diagnostics", fmt.Sprintf("%d hidden; use --raw or --json", diagnostics)))
		}
	}
	if len(view.Hunks) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Diff"))
		for index, hunk := range view.Hunks {
			if index > 0 {
				fmt.Fprintln(out)
			}
			if hunk.ID != "" {
				details := hunk.ID
				if hunk.Symbol != "" {
					details += " / " + hunk.Symbol
				}
				fmt.Fprintf(out, "%s %s\n", muted("hunk"), muted(details))
			}
			fmt.Fprint(out, hunk.Patch)
			if !strings.HasSuffix(hunk.Patch, "\n") {
				fmt.Fprintln(out)
			}
		}
	}
}

func demuxRevisionMode(revision authoring.RevisionProposal) string {
	if revision.UseHunks {
		return fmt.Sprintf("hunk (%d)", len(revision.HunkIDs))
	}
	return "file"
}

func printDemuxRevisionDiagnostics(out io.Writer, view authoring.DemuxRevisionView) {
	if len(view.FeasibilityWarnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Feasibility warnings"))
		for _, warning := range view.FeasibilityWarnings {
			details := warning.Source
			if warning.RevisionID != "" {
				details += " / " + warning.RevisionID
			}
			fmt.Fprintf(out, "  %s %s  %s\n", warningMarker(warning.Severity), warning.Message, muted("["+details+"]"))
		}
	}
	if len(view.StructuralDeps) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Structural dependencies"))
		for _, dep := range view.StructuralDeps {
			fmt.Fprintf(out, "  %s %s -> %s  %s\n", muted("-"), dep.FromFile, dep.ToFile, muted(dep.Symbol))
		}
	}
}

func printDemuxProposal(out io.Writer, proposal authoring.DemuxProposal, opts demuxPrintOptions) {
	if len(proposal.Revisions) == 0 {
		fmt.Fprintln(out, section("Nothing to demux"))
		fmt.Fprintln(out, muted("The working revision has no changes."))
		return
	}
	printDemuxProposalSummary(out, proposal, "")
	fmt.Fprintln(out)
	printDemuxProposalStacks(out, proposal)
	if opts.Raw && len(proposal.FeasibilityWarnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Warnings"))
		for _, warning := range proposal.FeasibilityWarnings {
			details := warning.Source
			if warning.RevisionID != "" {
				details += " / " + warning.RevisionID
			}
			fmt.Fprintf(out, "  %s %s  %s\n", muted(firstNonEmptyString(warning.RevisionID, warningMarker(warning.Severity))), muted(warning.Message), muted("["+details+"]"))
		}
	}
	if opts.Raw && len(proposal.Warnings) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section("Warnings"))
		for _, warning := range proposal.Warnings {
			fmt.Fprintf(out, "  %s %s\n", danger("!"), warning)
		}
	}
	if !opts.Raw {
		diagnostics := len(proposal.FeasibilityWarnings) + len(proposal.Warnings)
		if diagnostics > 0 {
			fmt.Fprintln(out)
			fmt.Fprintln(out, labelValue("Diagnostics", fmt.Sprintf("%d hidden; use --raw or --json", diagnostics)))
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, labelValue("JSON", "gx compose --json"))
}

type demuxStackDisplayGroup struct {
	Target    string
	Label     string
	Revisions []authoring.RevisionProposal
}

func printDemuxProposalStacks(out io.Writer, proposal authoring.DemuxProposal) {
	groups := demuxStackDisplayGroups(proposal.Revisions)
	if len(groups) == 1 && groups[0].Target == "" {
		fmt.Fprintln(out, section("Revisions"))
		for _, revision := range groups[0].Revisions {
			printDemuxRevisionLine(out, revision, "  ")
		}
		return
	}
	fmt.Fprintln(out, section("Stacks"))
	for groupIndex, group := range groups {
		if groupIndex > 0 {
			fmt.Fprintln(out)
		}
		meta := fmt.Sprintf("%d %s", len(group.Revisions), pluralize("revision", len(group.Revisions)))
		fmt.Fprintf(out, "  %s %s  %s\n", muted("●"), valueText(group.Label), muted(meta))
		for _, revision := range group.Revisions {
			printDemuxRevisionLine(out, revision, "    ")
		}
	}
}

func demuxStackDisplayGroups(revisions []authoring.RevisionProposal) []demuxStackDisplayGroup {
	byTarget := map[string]int{}
	groups := []demuxStackDisplayGroup{}
	for _, revision := range revisions {
		target := strings.TrimSpace(revision.TargetStack)
		index, ok := byTarget[target]
		if !ok {
			label := target
			if label == "" {
				label = "Current stack"
			}
			index = len(groups)
			byTarget[target] = index
			groups = append(groups, demuxStackDisplayGroup{Target: target, Label: label})
		}
		groups[index].Revisions = append(groups[index].Revisions, revision)
	}
	return groups
}

func printDemuxRevisionLine(out io.Writer, revision authoring.RevisionProposal, padding string) {
	fmt.Fprintf(out, "%s%s %s\n", padding, command(revision.ID), value(revision.Intent))
	for _, file := range revision.Files {
		fmt.Fprintf(out, "%s    %s\n", padding, muted(file))
	}
	for _, hunk := range revision.Hunks {
		if hunk.File != "" && !stringInSlice(hunk.File, revision.Files) {
			fmt.Fprintf(out, "%s    %s\n", padding, muted(hunk.File))
		}
	}
}

func printDemuxProposalSummary(out io.Writer, proposal authoring.DemuxProposal, state authoring.DemuxWorkflowState) {
	fmt.Fprintln(out, section("Compose proposal"))
	fmt.Fprintln(out, labelValue("Found", demuxFoundSummary(proposal)))
	fmt.Fprintln(out, labelValue("Proposes", fmt.Sprintf("%d revisions", len(proposal.Revisions))))
	status := demuxProposalDisplayStatus(proposal, state)
	if status != "" {
		fmt.Fprintln(out, labelValue("Status", status))
	}
}

func demuxFoundSummary(proposal authoring.DemuxProposal) string {
	hunks := len(proposal.Hunks)
	files := demuxProposalFileCount(proposal)
	if hunks > 0 {
		return fmt.Sprintf("%d %s across %d %s", hunks, pluralize("hunk", hunks), files, pluralize("file", files))
	}
	return fmt.Sprintf("%d %s", files, pluralize("file", files))
}

func demuxProposalFileCount(proposal authoring.DemuxProposal) int {
	seen := map[string]bool{}
	for _, hunk := range proposal.Hunks {
		if hunk.File != "" {
			seen[hunk.File] = true
		}
	}
	for _, revision := range proposal.Revisions {
		for _, file := range revision.Files {
			if file != "" {
				seen[file] = true
			}
		}
		for _, hunk := range revision.Hunks {
			if hunk.File != "" {
				seen[hunk.File] = true
			}
		}
	}
	return len(seen)
}

func demuxProposalDisplayStatus(proposal authoring.DemuxProposal, state authoring.DemuxWorkflowState) string {
	if state == authoring.DemuxWorkflowRepairRequired || state == authoring.DemuxWorkflowRepairRecommended {
		return "needs review"
	}
	blockingWarnings, _ := splitFeasibilityWarnings(proposal.FeasibilityWarnings)
	if len(blockingWarnings) > 0 {
		return "needs review"
	}
	if proposal.Status != "" && proposal.Status != authoring.ProposalPending {
		return string(proposal.Status)
	}
	return ""
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func stringInSlice(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

func warningMarker(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "warning", "warn":
		return danger("!")
	case "info":
		return muted("i")
	default:
		return danger("!")
	}
}

func newEditCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	return &cobra.Command{
		Use:    "edit [rev]",
		Short:  "Return to an existing GX change and keep working on it",
		Hidden: true,
		Args:   cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rev := ""
			if len(args) == 1 {
				rev = args[0]
			} else {
				snapshot, err := engine.StatusSnapshot(ctx)
				if err != nil {
					return err
				}
				candidates, err := engine.ModifyCandidates(ctx, 20)
				if err != nil {
					return err
				}
				selection, err := clitui.SelectModifyRevision(
					cmd.InOrStdin(),
					cmd.OutOrStdout(),
					snapshot,
					revisionsFromModifyCandidates(snapshot, candidates),
				)
				if err != nil {
					return err
				}
				rev = selection
			}
			result, err := engine.Modify(ctx, rev)
			if err != nil {
				return err
			}
			printModifySummary(cmd.OutOrStdout(), result)
			return nil
		},
	}
}

func revisionsFromModifyCandidates(snapshot authoring.StatusSnapshot, candidates []authoring.ChangeInfo) []vcs.RevisionSnapshot {
	activeChangeID := ""
	for _, bookmark := range snapshot.Bookmarks {
		if !bookmark.Current {
			continue
		}
		for _, revision := range bookmark.Revisions {
			if revision.Active {
				activeChangeID = revision.ChangeID
				break
			}
		}
		break
	}
	revisions := make([]vcs.RevisionSnapshot, 0, len(candidates))
	for index, candidate := range candidates {
		description := strings.TrimSpace(candidate.Description)
		if description == "" || description == "(no description set)" {
			description = "(no message)"
		}
		changeID := shortID(candidate.ChangeID, 12)
		commitID := shortID(candidate.CommitID, 8)
		short := changeID
		if commitID != "" {
			short = commitID
		}
		revisions = append(revisions, vcs.RevisionSnapshot{
			Index:       index + 1,
			ChangeID:    candidate.ChangeID,
			CommitID:    candidate.CommitID,
			Description: description,
			ShortID:     short,
			Active:      activeChangeID != "" && candidate.ChangeID == activeChangeID,
			Working:     activeChangeID != "" && candidate.ChangeID == activeChangeID,
			SyncNote:    "local",
		})
	}
	return revisions
}

func newStatusCommand(ctx context.Context, engine *authoring.Engine, use string, helpName string, hidden bool) *cobra.Command {
	var jsonOut bool
	var agentOut bool
	cmd := &cobra.Command{
		Use:    use,
		Short:  "Show the current GX revision and changed files",
		Hidden: hidden,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printCurrentStatus(ctx, engine, cmd.OutOrStdout(), jsonOut, agentOut)
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&agentOut, "agent", false, "print stable agent-readable text")
	setHelpName(cmd, helpName)
	return cmd
}

func newReviewCommand(ctx context.Context) *cobra.Command {
	var scope string
	var focus string
	var deep bool
	var verbose bool
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Review current changes with local facts, indexed context, and configured AI reviewers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := vcs.NewService().ResolveGitRepo(ctx)
			if err != nil {
				return err
			}
			report, err := runReviewWithLoader(
				cmd.InOrStdin(),
				cmd.ErrOrStderr(),
				func(progress io.Writer) (codereview.Report, error) {
					return codereview.Review(ctx, repo.RootPath, codereview.Options{
						Scope:          scope,
						Deep:           deep,
						Focus:          focus,
						Verbose:        verbose,
						ProgressWriter: progress,
						Color:          true,
					})
				},
			)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), codereview.RenderMarkdown(report))
			return nil
		},
	}
	cmd.Flags().StringVar(&scope, "scope", codereview.DefaultScope, "review scope: architecture, security, performance, onboarding, docs, dependencies, testing, maintainability")
	cmd.Flags().StringVar(&focus, "focus", "", "limit review to files under this path prefix")
	cmd.Flags().BoolVar(&deep, "deep", false, "retrieve more local and indexed context")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "include repo facts, docs, and changed files")
	return cmd
}

type currentStatus struct {
	Repo           authoring.RepoInfo      `json:"repo"`
	Stack          *authoring.StackInfo    `json:"stack,omitempty"`
	Refs           currentStatusRefs       `json:"refs"`
	Current        authoring.ChangeInfo    `json:"current"`
	Parent         authoring.ChangeInfo    `json:"parent"`
	PublishUploads publication.QueueStatus `json:"publish_uploads"`
	NeedsMessage   bool                    `json:"needs_message"`
	Recorded       bool                    `json:"recorded"`
	Files          []string                `json:"files"`
	Next           []string                `json:"next"`
	GitStatusNote  string                  `json:"git_status_note"`
}

type currentStatusRefs struct {
	GXBaseRef       string `json:"gx_base_ref"`
	GXStackRef      string `json:"gx_stack_ref,omitempty"`
	GitCheckoutRef  string `json:"git_checkout_ref,omitempty"`
	GitPublishedRef string `json:"git_published_ref,omitempty"`
}

func printCurrentStatus(ctx context.Context, engine *authoring.Engine, out io.Writer, jsonOut, agentOut bool) error {
	status, err := currentStatusForEngine(ctx, engine)
	if err != nil {
		return err
	}
	if jsonOut {
		return json.NewEncoder(out).Encode(status)
	}
	if agentOut {
		printCurrentStatusAgent(out, status)
		return nil
	}
	printCurrentStatusHuman(out, status)
	return nil
}

func currentStatusForEngine(ctx context.Context, engine *authoring.Engine) (currentStatus, error) {
	stack, err := engine.Status(ctx)
	if err != nil {
		return currentStatus{}, err
	}
	current, err := engine.CurrentChange(ctx, stack.Repo.RootPath, "@")
	if err != nil {
		return currentStatus{}, err
	}
	parent, _ := engine.CurrentChange(ctx, stack.Repo.RootPath, "@-")
	recorded := false
	for _, revision := range stack.Revisions {
		if revision.ChangeID == current.ChangeID {
			recorded = true
			break
		}
	}
	needsMessage := strings.TrimSpace(current.Description) == "" || strings.TrimSpace(current.Description) == "(no description set)"
	gitCheckoutRef := pointerString(stack.Repo.BranchName)
	next := []string{"gx compose", "gx stacks"}
	if gitCheckoutRef == "gx/edit" && len(current.Files) > 0 {
		next = []string{`gx add -m "describe this revision"`, "gx stacks"}
	}
	refs := currentStatusRefs{
		GXBaseRef:      currentStatusBaseStack(currentStatus{Repo: stack.Repo, Stack: stack.Stack}),
		GitCheckoutRef: gitCheckoutRef,
	}
	if stack.Stack != nil {
		refs.GXStackRef = stack.Stack.BookmarkName
		refs.GitPublishedRef = pointerString(stack.Stack.RemoteRef)
	}
	publishUploads, _ := publication.QueuedUploadStatus()
	return currentStatus{
		Repo:           stack.Repo,
		Stack:          stack.Stack,
		Refs:           refs,
		Current:        current,
		Parent:         parent,
		PublishUploads: publishUploads,
		NeedsMessage:   needsMessage,
		Recorded:       recorded,
		Files:          current.Files,
		Next:           next,
		GitStatusNote:  "gx stores new changes in revisions, so `git status` may be clean.",
	}, nil
}

func newStackCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var newName string
	var jsonOut bool
	var agentOut bool
	cmd := &cobra.Command{
		Use:    "stack --new <name>",
		Short:  "Create a GX stack",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			newName = strings.TrimSpace(newName)
			if newName == "" {
				return fmt.Errorf("stack name is required; use gx stack --new <name>")
			}
			result, err := engine.CreateStack(ctx, newName)
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, result)
			}
			out := cmd.OutOrStdout()
			if agentOut {
				fmt.Fprintf(out, "created=true stack=%s ref=%s base=%s current_revision=%s\n",
					quoteAgent(result.Stack.Name),
					quoteAgent(result.Stack.BookmarkName),
					quoteAgent(result.Stack.BaseRef),
					quoteAgent(shortID(result.CurrentChange.ChangeID, 12)),
				)
				return nil
			}
			fmt.Fprintln(out, labelValue("created", result.Stack.Name))
			fmt.Fprintln(out, labelValue("stack", result.Stack.BookmarkName))
			fmt.Fprintln(out, labelValue("base", firstNonEmptyString(result.Stack.BaseRef, "main")))
			fmt.Fprintln(out)
			fmt.Fprintln(out, section("Next"))
			fmt.Fprintf(out, "  %s\n", command(`gx add -m "first revision"`))
			return nil
		},
	}
	cmd.Flags().StringVarP(&newName, "new", "n", "", "create and switch to a new stack")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&agentOut, "agent", false, "print stable agent-readable text")
	return cmd
}

func printCurrentStatusHuman(out io.Writer, status currentStatus) {
	fmt.Fprintln(out, commandLine("gx status", true))
	fmt.Fprintln(out)
	revision := shortID(status.Current.ChangeID, 12)
	if revision == "" {
		revision = "(unknown)"
	}
	if len(status.Files) > 0 {
		fmt.Fprintf(out, "%d files are currently waiting to be assigned.\n", len(status.Files))
		if status.Refs.GitCheckoutRef == "gx/edit" {
			fmt.Fprintf(out, "Run %s to record them onto the active edited stack revision.\n", command(`gx add -m "describe this revision"`))
		} else {
			fmt.Fprintf(out, "Run %s to add them to the pending compose proposal.\n", command("gx compose"))
		}
	} else {
		fmt.Fprintf(out, "No files are currently assigned to revision %s.\n", valueText(revision))
	}
	fmt.Fprintln(out)

	message := strings.TrimSpace(status.Current.Description)
	if status.NeedsMessage {
		fmt.Fprintln(out, danger("●")+"  "+danger("Currently no message assigned"))
	} else {
		fmt.Fprintln(out, success("●")+"  "+valueText("Message: "+message))
	}
	refs := status.Refs
	if refs.GXBaseRef == "" {
		refs.GXBaseRef = currentStatusBaseStack(status)
	}
	if refs.GXStackRef == "" && status.Stack != nil {
		refs.GXStackRef = status.Stack.BookmarkName
	}
	if refs.GitCheckoutRef == "" {
		refs.GitCheckoutRef = pointerString(status.Repo.BranchName)
	}
	if refs.GitPublishedRef == "" && status.Stack != nil {
		refs.GitPublishedRef = pointerString(status.Stack.RemoteRef)
	}
	fmt.Fprintf(out, "   %s %s\n", muted("GX base ref:"), valueText(refs.GXBaseRef))
	if refs.GXStackRef != "" {
		fmt.Fprintf(out, "   %s %s\n", muted("GX stack ref:"), valueText(refs.GXStackRef))
	}
	if refs.GitCheckoutRef != "" {
		fmt.Fprintf(out, "   %s %s\n", muted("Git checkout ref:"), valueText(refs.GitCheckoutRef))
	}
	if refs.GitPublishedRef != "" {
		fmt.Fprintf(out, "   %s %s\n", muted("Git published ref:"), valueText(refs.GitPublishedRef))
	}
	printPublishUploadStatus(out, status.PublishUploads)
	if len(status.Files) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, section(fmt.Sprintf("Files (%d)", len(status.Files))))
		for _, file := range status.Files {
			fmt.Fprintf(out, "  %s\n", file)
		}
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, hint(status.GitStatusNote))
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Next"))
	for _, next := range status.Next {
		fmt.Fprintf(out, "  %s\n", command(next))
	}
}

func printPublishUploadStatus(out io.Writer, status publication.QueueStatus) {
	total := status.Pending + status.Uploading + status.Failed
	if total == 0 {
		return
	}
	parts := []string{}
	if status.Pending > 0 {
		parts = append(parts, fmt.Sprintf("%d pending", status.Pending))
	}
	if status.Uploading > 0 {
		parts = append(parts, fmt.Sprintf("%d uploading", status.Uploading))
	}
	if status.Failed > 0 {
		parts = append(parts, fmt.Sprintf("%d failed", status.Failed))
	}
	label := strings.Join(parts, ", ")
	if status.UploadRunning {
		label += " (worker running)"
	}
	if status.Failed > 0 && strings.TrimSpace(status.LastError) != "" {
		fmt.Fprintf(out, "   %s %s\n", muted("GX context uploads:"), danger(label))
		fmt.Fprintf(out, "   %s %s\n", muted("Upload error:"), status.LastError)
		return
	}
	fmt.Fprintf(out, "   %s %s\n", muted("GX context uploads:"), valueText(label))
}

func currentStatusBaseStack(status currentStatus) string {
	if status.Repo.AuthoringBase != nil && strings.TrimSpace(*status.Repo.AuthoringBase) != "" {
		base := strings.TrimSpace(*status.Repo.AuthoringBase)
		if !strings.HasPrefix(base, "gx/") {
			return base
		}
	}
	if status.Repo.DefaultBranch != nil && strings.TrimSpace(*status.Repo.DefaultBranch) != "" {
		return strings.TrimSpace(*status.Repo.DefaultBranch)
	}
	return "main"
}

func printCurrentStatusAgent(out io.Writer, status currentStatus) {
	stackName := ""
	if status.Stack != nil {
		stackName = status.Stack.Name
	}
	fmt.Fprintf(out, "repo=%s stack=%s gx_base_ref=%s gx_stack_ref=%s git_checkout_ref=%s git_published_ref=%s revision=%s commit=%s needs_message=%t recorded=%t files=%d\n",
		quoteAgent(repoLabel(status.Repo)),
		quoteAgent(stackName),
		quoteAgent(status.Refs.GXBaseRef),
		quoteAgent(status.Refs.GXStackRef),
		quoteAgent(status.Refs.GitCheckoutRef),
		quoteAgent(status.Refs.GitPublishedRef),
		quoteAgent(shortID(status.Current.ChangeID, 12)),
		quoteAgent(shortID(status.Current.CommitID, 8)),
		status.NeedsMessage,
		status.Recorded,
		len(status.Files),
	)
	if status.PublishUploads.Pending > 0 || status.PublishUploads.Failed > 0 || status.PublishUploads.Uploading > 0 {
		fmt.Fprintf(out, "publish_uploads pending=%d uploading=%d failed=%d running=%t\n",
			status.PublishUploads.Pending,
			status.PublishUploads.Uploading,
			status.PublishUploads.Failed,
			status.PublishUploads.UploadRunning,
		)
	}
	for _, file := range status.Files {
		fmt.Fprintf(out, "file %s\n", quoteAgent(file))
	}
}

func newStacksCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var agentOut bool
	var showAll bool
	cmd := &cobra.Command{
		Use:   "stacks [stack]",
		Short: "Browse GX stacks and revisions",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOut {
				stack, err := engine.Status(ctx)
				if err != nil {
					return err
				}
				stack = stackSummaryForStacksDisplay(stack, stackDisplayOptions{ShowAll: showAll})
				return writeJSON(cmd, stack)
			}
			if len(args) > 0 && !agentOut {
				return fmt.Errorf("stack selector is only supported with --agent")
			}
			return printStacks(ctx, engine, cmd.InOrStdin(), cmd.OutOrStdout(), agentOut, firstArg(args), stackDisplayOptions{ShowAll: showAll})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&agentOut, "agent", false, "print stable agent-readable text")
	cmd.Flags().BoolVar(&showAll, "show-all", false, "accepted for compatibility; merged stacks are not shown")
	if flag := cmd.Flags().Lookup("show-all"); flag != nil {
		flag.Hidden = true
	}
	cmd.AddCommand(
		newStacksEditCommand(ctx, engine),
		newStacksDiffCommand(),
	)
	return cmd
}

func newStacksEditCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <revision>",
		Short: "Edit a revision",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := engine.Modify(ctx, args[0])
			if err != nil {
				return err
			}
			printModifySummary(cmd.OutOrStdout(), result)
			return nil
		},
	}
}

func newStacksDiffCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <revision-or-stack>",
		Short: "Show a JJ diff for a revision or stack",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJjDiff(cmd.OutOrStdout(), args[0])
		},
	}
}

func printStacks(ctx context.Context, engine *authoring.Engine, in io.Reader, out io.Writer, agentOut bool, selector string, opts stackDisplayOptions) error {
	stack, err := engine.Status(ctx)
	if err != nil {
		return err
	}
	stack = stackSummaryForStacksDisplay(stack, opts)
	if agentOut {
		printStatusAgent(out, stack, selector)
		return nil
	}
	var unrecorded *authoring.ChangeInfo
	if stack.Repo.RootPath != "" {
		if current, currentErr := engine.CurrentChange(ctx, stack.Repo.RootPath, "@"); currentErr == nil {
			found := false
			for _, unit := range stack.Revisions {
				if unit.ChangeID == current.ChangeID {
					found = true
					break
				}
			}
			if !found && len(current.Files) > 0 {
				unrecorded = &current
			}
		}
	}
	if useStatusInteractive(in, out) {
		action, err := runStacksInteractive(in, out, stack, unrecorded)
		if err != nil {
			return err
		}
		return runStacksAction(ctx, engine, out, action)
	}
	printStatusSummary(out, stack, unrecorded)
	return nil
}

func runStacksAction(ctx context.Context, engine *authoring.Engine, out io.Writer, action stacksAction) error {
	switch action.Kind {
	case "", "quit":
		return nil
	case "edit":
		result, err := engine.Modify(ctx, action.Target)
		if err != nil {
			return err
		}
		printModifySummary(out, result)
		return nil
	case "diff":
		return runJjDiff(out, action.Target)
	case "delete-revision":
		result, err := engine.DeleteRevision(ctx, action.Target)
		if err != nil {
			return err
		}
		printDeleteRevisionSummary(out, result)
		return nil
	case "delete-stack":
		result, err := engine.DeleteStack(ctx, action.Target)
		if err != nil {
			return err
		}
		printDeleteStackSummary(out, result)
		return nil
	default:
		return fmt.Errorf("unknown stacks action %q", action.Kind)
	}
}

func printDeleteRevisionSummary(out io.Writer, result authoring.DeleteRevisionResult) {
	fmt.Fprintln(out, success("Revision deleted"))
	fmt.Fprintln(out, labelValue("Revision", shortID(result.Change.ChangeID, 12)))
	fmt.Fprintln(out, labelValue("Description", firstNonEmptyString(result.Change.Description, "(no description set)")))
}

func printDeleteStackSummary(out io.Writer, result authoring.DeleteStackResult) {
	fmt.Fprintln(out, success("Stack deleted"))
	fmt.Fprintln(out, labelValue("Stack", firstNonEmptyString(result.Stack.Name, result.Stack.BookmarkName)))
	fmt.Fprintln(out, labelValue("Revisions", fmt.Sprintf("%d", len(result.Revisions))))
}

func runJjDiff(out io.Writer, rev string) error {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return fmt.Errorf("revision is required")
	}
	cmd := exec.Command("jj", "diff", "-r", rev)
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

func printStatus(ctx context.Context, engine *authoring.Engine, in io.Reader, out io.Writer, agentOut bool, selector string) error {
	stack, err := engine.Status(ctx)
	if err != nil {
		return err
	}
	if agentOut {
		printStatusAgent(out, stack, selector)
		return nil
	}
	var unrecorded *authoring.ChangeInfo
	if stack.Repo.RootPath != "" {
		if current, currentErr := engine.CurrentChange(ctx, stack.Repo.RootPath, "@"); currentErr == nil {
			found := false
			for _, unit := range stack.Revisions {
				if unit.ChangeID == current.ChangeID {
					found = true
					break
				}
			}
			if !found && len(current.Files) > 0 {
				unrecorded = &current
			}
		}
	}
	if useStatusInteractive(in, out) {
		action, err := runStacksInteractive(in, out, stack, unrecorded)
		if err != nil {
			return err
		}
		return runStacksAction(ctx, engine, out, action)
	}
	printStatusSummary(out, stack, unrecorded)
	return nil
}

func printStatusSummary(out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo) {
	if len(stack.Revisions) == 0 && len(stack.Stacks) == 0 && stack.Stack == nil {
		fmt.Fprintln(out, muted("No GX revisions recorded yet."))
		return
	}
	fmt.Fprint(out, renderStacksSummary(stack, unrecorded, currentStackIndex(stack), false, latestRevisionDisplayIndex(stack.Revisions, unrecorded)))
}

func printCurrentRevisions(out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, indent int) {
	padding := strings.Repeat(" ", indent)
	if len(stack.Revisions) == 0 {
		fmt.Fprintf(out, "%s%s\n", padding, muted("(no revisions)"))
	} else {
		for _, unit := range stack.Revisions {
			tags := make([]string, 0, 1)
			if unit.Status != "" && unit.Status != "draft" {
				tags = append(tags, unit.Status)
			}
			line := fmt.Sprintf("%s%s %s  %s", padding, revisionStatusPrefix(unit), accent(shortID(unit.ChangeID, 7)), valueText(firstNonEmptyString(unit.Description, "(no description set)")))
			if len(tags) > 0 {
				line += "  " + muted(strings.Join(tags, ", "))
			}
			fmt.Fprintln(out, line)
		}
	}
	if unrecorded != nil {
		line := fmt.Sprintf("%s%s %s  %s", padding, logoText("*"), valueText(firstNonEmptyString(unrecorded.Description, "(no description set)")), muted("unrecorded"))
		fmt.Fprintln(out, line)
	}
}

func renderStacksSummary(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int) string {
	stack = stackSummaryWithDisplayFallback(stack)
	stacks := orderedStacks(stack)
	if len(stacks) == 0 {
		var out strings.Builder
		fmt.Fprintln(&out, commandLine("gx stacks", true))
		fmt.Fprintln(&out)
		printCurrentRevisions(&out, stack, unrecorded, 0)
		fmt.Fprintln(&out)
		fmt.Fprintln(&out, stacksLegend(stackMode))
		return out.String()
	}
	if stackCursor < 0 {
		stackCursor = 0
	}
	if stackCursor >= len(stacks) {
		stackCursor = len(stacks) - 1
	}
	lines := []string{
		commandLine("gx stacks", true),
		"",
		stacksHeaderLine(stack, len(stacks), stackMode),
		"",
	}
	previousBucket := -1
	for index, entry := range stacks {
		bucket := stackDisplayBucket(entry)
		if bucket == 1 && previousBucket != 1 {
			if index > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, section("Published"))
			lines = append(lines, "")
		} else if index > 0 {
			lines = append(lines, muted(strings.Repeat("─", 52)))
		}
		previousBucket = bucket
		selected := index == stackCursor
		if selected {
			lines = append(lines, statusBookmarkLine("●", entry, stackStatusMeta(stack, entry), true))
			cursor := -1
			if !stackMode {
				cursor = revCursor
			}
			lines = append(lines, renderStackRevisionLines(stack, entry, unrecorded, cursor)...)
			continue
		}
		lines = append(lines, statusBookmarkLine("○", entry, stackStatusMeta(stack, entry), false))
		lines = append(lines, renderStackRevisionLines(stack, entry, unrecorded, -1)...)
	}
	lines = append(lines, "", stacksLegend(stackMode))
	return strings.Join(lines, "\n") + "\n"
}

func renderSelectableRevisionLines(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, cursor int) []string {
	items := displayRevisions(stack.Revisions, unrecorded)
	return renderRevisionLines(items, cursor)
}

func renderStackRevisionLines(stack authoring.StackSummary, entry authoring.StackInfo, unrecorded *authoring.ChangeInfo, cursor int) []string {
	current := stack.Stack != nil && entry.BookmarkName == stack.Stack.BookmarkName
	entryUnrecorded := unrecorded
	if !current {
		entryUnrecorded = nil
	}
	items := displayRevisions(stackEntryRevisions(stack, entry), entryUnrecorded)
	return renderRevisionLines(items, cursor)
}

func renderRevisionLines(items []authoring.RevisionSummary, cursor int) []string {
	if len(items) == 0 {
		return []string{"    " + muted("(no revisions)")}
	}
	if cursor >= len(items) {
		cursor = len(items) - 1
	}
	lines := make([]string, 0, len(items))
	for index, unit := range items {
		marker := " "
		id := muted(shortID(unit.ChangeID, 7))
		if index == cursor {
			marker = logoText("›")
			id = accent(shortID(unit.ChangeID, 7))
		}
		line := fmt.Sprintf("    %s %s %s  %s", marker, revisionStatusPrefix(unit), id, valueText(firstNonEmptyString(unit.Description, "(no description set)")))
		if index == 0 {
			line += "  " + muted("latest")
		}
		if unit.Status != "" && unit.Status != "draft" {
			line += "  " + muted(unit.Status)
		}
		lines = append(lines, line)
	}
	return lines
}

func revisionsWithUnrecorded(revisions []authoring.RevisionSummary, unrecorded *authoring.ChangeInfo) []authoring.RevisionSummary {
	items := make([]authoring.RevisionSummary, 0, len(revisions)+1)
	items = append(items, revisions...)
	if unrecorded != nil {
		items = append(items, authoring.RevisionSummary{
			ChangeID:    unrecorded.ChangeID,
			CommitID:    unrecorded.CommitID,
			Description: firstNonEmptyString(unrecorded.Description, "(no description set)"),
			Status:      "unrecorded",
			Active:      true,
		})
	}
	return items
}

func activeRevisionIndex(revisions []authoring.RevisionSummary, unrecorded *authoring.ChangeInfo) int {
	for index, revision := range revisionsWithUnrecorded(revisions, unrecorded) {
		if revision.Active {
			return index
		}
	}
	return 0
}

func displayRevisions(revisions []authoring.RevisionSummary, unrecorded *authoring.ChangeInfo) []authoring.RevisionSummary {
	items := revisionsWithUnrecorded(revisions, unrecorded)
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}
	return items
}

func latestRevisionDisplayIndex(revisions []authoring.RevisionSummary, unrecorded *authoring.ChangeInfo) int {
	if len(displayRevisions(revisions, unrecorded)) == 0 {
		return 0
	}
	return 0
}

func printStatusAgent(out io.Writer, stack authoring.StackSummary, selector string) {
	stack = stackSummaryWithDisplayFallback(stack)
	if strings.TrimSpace(selector) != "" {
		target, ok := findStack(stack, selector)
		if !ok {
			fmt.Fprintf(out, "error=unknown-stack selector=%s\n", quoteAgent(selector))
			return
		}
		fmt.Fprintf(out, "stack=%s repo=%s\n", quoteAgent(target.Name), quoteAgent(repoLabel(stack.Repo)))
		fmt.Fprintln(out, stackAgentLine(stack, target, false))
		printRevisionAgentLines(out, stackEntryRevisions(stack, target))
		return
	}
	stacks := orderedStacks(stack)
	fmt.Fprintf(out, "repo=%s stacks=%d current=%s\n", quoteAgent(repoLabel(stack.Repo)), len(stacks), quoteAgent(currentStackName(stack)))
	for _, entry := range stacks {
		fmt.Fprintln(out, stackAgentLine(stack, entry, true))
		printRevisionAgentLines(out, stackEntryRevisions(stack, entry))
	}
}

func printRevisionAgentLines(out io.Writer, revisions []authoring.RevisionSummary) {
	for _, revision := range revisions {
		status := revision.Status
		if status == "" {
			status = "draft"
		}
		parts := []string{
			"revision",
			"id=" + quoteAgent(shortID(revision.ChangeID, 12)),
			"message=" + quoteAgent(revision.Description),
			"status=" + quoteAgent(status),
		}
		if revision.CommitID != "" {
			parts = append(parts, "commit="+quoteAgent(shortID(revision.CommitID, 8)))
		}
		if revision.Active {
			parts = append(parts, "active=true")
		}
		if revision.Published {
			parts = append(parts, "published=true")
		}
		fmt.Fprintln(out, strings.Join(parts, " "))
	}
}

func printModifySummary(out io.Writer, result authoring.ModifyResult) {
	invocation := "gx edit"
	if result.CurrentChange.ChangeID != "" {
		invocation += " " + shortID(result.CurrentChange.ChangeID, 12)
	}
	fmt.Fprintln(out, commandLine(invocation, false))
	fmt.Fprintln(out)
	if strings.TrimSpace(result.Output) != "" {
		fmt.Fprintln(out, muted(strings.TrimSpace(result.Output)))
	}
	fmt.Fprintln(out, labelToken("revision", shortID(result.CurrentChange.ChangeID, 12)))
	if result.CurrentChange.CommitID != "" {
		fmt.Fprintln(out, labelToken("commit", shortID(result.CurrentChange.CommitID, 8)))
	}
	if result.CurrentChange.Description != "" {
		fmt.Fprintln(out, labelToken("message", result.CurrentChange.Description))
	}
	if result.Stack != nil {
		fmt.Fprintln(out, labelToken("stack", result.Stack.Name)+"  "+muted(stackMeta(authoring.StackSummary{Stack: result.Stack, Revisions: []authoring.RevisionSummary{{ChangeID: result.CurrentChange.ChangeID}}}, *result.Stack)))
	}
	if result.Repo.BranchName != nil && strings.TrimSpace(*result.Repo.BranchName) != "" {
		fmt.Fprintln(out, labelToken("git branch", strings.TrimSpace(*result.Repo.BranchName)))
	}
	fmt.Fprintln(out, muted(strings.Repeat("-", 48)))
	fmt.Fprintln(out, labelToken("next", "gx compose"))
}

func stackAgentLine(summary authoring.StackSummary, entry authoring.StackInfo, includeName bool) string {
	current := summary.Stack != nil && entry.BookmarkName == summary.Stack.BookmarkName
	changeTotal, publishedTotal := stackEntryRevisionCounts(summary, entry)
	status := stackDisplayStatus(entry.Status, changeTotal, publishedTotal)
	prefix := "stack"
	if includeName {
		prefix += " name=" + quoteAgent(entry.Name)
	}
	return fmt.Sprintf("%s alias=%s ref=%s base=%s status=%s changes=%d published=%d/%d sync=%s current=%t",
		prefix,
		quoteAgent(entry.Alias),
		quoteAgent(entry.BookmarkName),
		quoteAgent(entry.BaseRef),
		quoteAgent(status),
		changeTotal,
		publishedTotal,
		changeTotal,
		quoteAgent(stackSync(entry)),
		current,
	)
}

func stackMeta(summary authoring.StackSummary, entry authoring.StackInfo) string {
	parts := []string{}
	if entry.Alias != "" {
		parts = append(parts, "alias "+entry.Alias)
	}
	if entry.BookmarkName != "" {
		parts = append(parts, "ref "+entry.BookmarkName)
	}
	if entry.BaseRef != "" {
		parts = append(parts, "base "+entry.BaseRef)
	}
	changeTotal, publishedTotal := stackEntryRevisionCounts(summary, entry)
	if summary.Stack != nil && entry.BookmarkName == summary.Stack.BookmarkName {
		parts = append(parts, "current")
	}
	if status := stackDisplayStatus(entry.Status, changeTotal, publishedTotal); status != "" {
		parts = append(parts, "state "+status)
	}
	if changeTotal > 0 {
		parts = append(parts, fmt.Sprintf("%d changes", changeTotal))
		parts = append(parts, fmt.Sprintf("%d/%d published", publishedTotal, changeTotal))
	}
	if sync := stackSync(entry); sync != "" {
		parts = append(parts, "sync "+sync)
	}
	return strings.Join(parts, " · ")
}

func stackStatusMeta(summary authoring.StackSummary, entry authoring.StackInfo) string {
	parts := []string{}
	alias := strings.TrimSpace(entry.Alias)
	if alias != "" {
		parts = append(parts, alias)
	}
	ref := compactBookmarkRef(entry.BookmarkName)
	name := strings.TrimSpace(entry.Name)
	if ref != "" && ref != alias && ref != name {
		parts = append(parts, ref)
	}
	if entry.BaseRef != "" {
		parts = append(parts, entry.BaseRef)
	}
	changeTotal, publishedTotal := stackEntryRevisionCounts(summary, entry)
	if status := compactStackStatus(stackDisplayStatus(entry.Status, changeTotal, publishedTotal)); status != "" {
		parts = append(parts, status)
	}
	if pending := changeTotal - publishedTotal; pending > 0 {
		parts = append(parts, fmt.Sprintf("↑%d", pending))
	}
	if publishedTotal > 0 {
		parts = append(parts, fmt.Sprintf("✓%d", publishedTotal))
	}
	if sync := stackSync(entry); sync != "" && sync != "local" {
		parts = append(parts, "sync "+sync)
	}
	return strings.Join(parts, " · ")
}

func stacksHeaderLine(stack authoring.StackSummary, stackCount int, stackMode bool) string {
	parts := []string{
		logoText(repoLabel(stack.Repo)),
	}
	if current := currentStackName(stack); current != "" {
		parts = append(parts, muted("● ")+valueText(current))
	}
	return strings.Join(parts, "  ")
}

func statusBookmarkLine(marker string, stack authoring.StackInfo, meta string, selected bool) string {
	renderedMarker := muted(marker)
	name := muted(stack.Name)
	if selected {
		renderedMarker = logoText(marker)
		name = valueText(stack.Name)
	}
	line := renderedMarker + " " + name
	if meta != "" {
		line += "  " + muted(meta)
	}
	return line
}

func stacksLegend(stackMode bool) string {
	parts := []string{"● selected", "○ other", "↑ local", "↓ cloud", "* active"}
	if stackMode {
		parts = append([]string{"j/k stack", "enter revisions", "d diff", "Shift+D delete", "q quit"}, parts...)
	} else {
		parts = append([]string{"j/k revision", "e edit", "d diff", "Shift+D delete", "esc stacks", "q quit"}, parts...)
	}
	return muted(strings.Join(parts, " · "))
}

func compactBookmarkRef(ref string) string {
	ref = strings.TrimSpace(ref)
	ref = strings.TrimPrefix(ref, "refs/heads/")
	ref = strings.TrimPrefix(ref, "gx/")
	return ref
}

func compactStackStatus(status string) string {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case "":
		return ""
	case "published":
		return "✓"
	case "draft":
		return "draft"
	default:
		return status
	}
}

func stackDisplayStatus(status string, revisions, published int) string {
	status = strings.TrimSpace(status)
	if isTerminalStackStatus(status) {
		return status
	}
	if revisions <= 0 {
		return status
	}
	if published >= revisions {
		return "published"
	}
	return "draft"
}

func currentStackIndex(stack authoring.StackSummary) int {
	if stack.Stack == nil {
		return 0
	}
	for index, entry := range orderedStacks(stack) {
		if entry.BookmarkName == stack.Stack.BookmarkName {
			return index
		}
	}
	return 0
}

func revisionStatusPrefix(unit authoring.RevisionSummary) string {
	arrow := logoText("↑")
	if unit.Published {
		arrow = success("↓")
	}
	if unit.Active {
		return logoText("*") + " " + arrow
	}
	return "  " + arrow
}

func stackSync(stack authoring.StackInfo) string {
	parts := []string{"local"}
	if stack.RemoteName != nil && strings.TrimSpace(*stack.RemoteName) != "" {
		parts = append(parts, strings.TrimSpace(*stack.RemoteName))
	} else if stack.RemoteRef != nil && strings.TrimSpace(*stack.RemoteRef) != "" {
		parts = append(parts, "origin")
	}
	return strings.Join(parts, ",")
}

func orderedStacks(stack authoring.StackSummary) []authoring.StackInfo {
	stacks := make([]authoring.StackInfo, 0, len(stack.Stacks)+1)
	seen := map[string]bool{}
	if stack.Stack != nil {
		current := *stack.Stack
		for _, entry := range stack.Stacks {
			if entry.BookmarkName == current.BookmarkName {
				if current.Alias == "" {
					current.Alias = entry.Alias
				}
				if current.Name == "" {
					current.Name = entry.Name
				}
				current.Status = entry.Status
				current.RevisionCount = entry.RevisionCount
				current.PublishedCount = entry.PublishedCount
				current.Revisions = entry.Revisions
				break
			}
		}
		stacks = append(stacks, current)
		seen[current.BookmarkName] = true
	}
	for _, entry := range stack.Stacks {
		if seen[entry.BookmarkName] {
			continue
		}
		stacks = append(stacks, entry)
	}
	return orderStacksForDisplay(stacks)
}

func stackEntryRevisions(summary authoring.StackSummary, entry authoring.StackInfo) []authoring.RevisionSummary {
	if summary.Stack != nil && entry.BookmarkName == summary.Stack.BookmarkName && len(summary.Revisions) > 0 {
		return summary.Revisions
	}
	return entry.Revisions
}

func stackEntryRevisionCounts(summary authoring.StackSummary, entry authoring.StackInfo) (int, int) {
	revisions := stackEntryRevisions(summary, entry)
	if len(revisions) == 0 {
		return entry.RevisionCount, entry.PublishedCount
	}
	published := 0
	for _, revision := range revisions {
		if revision.Published {
			published++
		}
	}
	return len(revisions), published
}

func orderStacksForDisplay(stacks []authoring.StackInfo) []authoring.StackInfo {
	if len(stacks) <= 1 {
		return stacks
	}
	ordered := make([]authoring.StackInfo, 0, len(stacks))
	for bucket := 0; bucket <= 2; bucket++ {
		for _, entry := range stacks {
			if stackDisplayBucket(entry) == bucket {
				ordered = append(ordered, entry)
			}
		}
	}
	return ordered
}

func stackSummaryWithDisplayFallback(stack authoring.StackSummary) authoring.StackSummary {
	if stack.Stack != nil || len(stack.Stacks) > 0 || len(stack.Revisions) == 0 {
		return stack
	}
	ref := firstNonEmptyString(stack.ContainerName, pointerString(stack.Repo.BranchName), "current")
	display := compactBookmarkRef(ref)
	if display == "" {
		display = "current"
	}
	base := "main"
	if stack.Repo.DefaultBranch != nil && strings.TrimSpace(*stack.Repo.DefaultBranch) != "" {
		base = strings.TrimSpace(*stack.Repo.DefaultBranch)
	}
	current := authoring.StackInfo{
		Name:         display,
		BookmarkName: ref,
		BaseRef:      base,
		Status:       "draft",
	}
	stack.Stack = &current
	stack.Stacks = []authoring.StackInfo{current}
	return stack
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func findStack(stack authoring.StackSummary, selector string) (authoring.StackInfo, bool) {
	selector = strings.TrimSpace(selector)
	for _, entry := range orderedStacks(stack) {
		if entry.Alias == selector || entry.Name == selector || entry.BookmarkName == selector {
			return entry, true
		}
	}
	for _, entry := range orderedStacks(stack) {
		if strings.HasPrefix(entry.Alias, selector) || strings.HasPrefix(entry.Name, selector) || strings.HasPrefix(entry.BookmarkName, selector) {
			return entry, true
		}
	}
	return authoring.StackInfo{}, false
}

func currentStackName(stack authoring.StackSummary) string {
	if stack.Stack == nil {
		return ""
	}
	return firstNonEmptyString(stack.Stack.Name, stack.Stack.BookmarkName)
}

func repoLabel(repo authoring.RepoInfo) string {
	if repo.RemoteURL != nil {
		raw := strings.TrimSpace(*repo.RemoteURL)
		if raw != "" {
			raw = strings.TrimSuffix(raw, ".git")
			if idx := strings.LastIndex(raw, ":"); idx >= 0 && !strings.Contains(raw, "://") {
				return raw[idx+1:]
			}
			parts := strings.Split(raw, "/")
			if len(parts) >= 2 {
				return parts[len(parts)-2] + "/" + parts[len(parts)-1]
			}
		}
	}
	if repo.RootPath == "" {
		return "unknown"
	}
	return filepath.Base(repo.RootPath)
}

func unpublishedCount(stack authoring.StackSummary) int {
	total := 0
	for _, unit := range stack.Revisions {
		if !unit.Published {
			total++
		}
	}
	return total
}

type stackDisplayOptions struct {
	ShowAll bool
}

func stackSummaryForStacksDisplay(stack authoring.StackSummary, opts stackDisplayOptions) authoring.StackSummary {
	stack = stackSummaryWithDisplayFallback(stack)
	stacks := orderedStacks(stack)
	filtered := stacks[:0]
	for _, entry := range stacks {
		if isMergedStack(entry) {
			continue
		}
		filtered = append(filtered, entry)
	}
	stacks = filtered
	stack.Stacks = stacks
	if stack.Stack != nil && !stackListContains(stacks, stack.Stack.BookmarkName) {
		stack.Stack = nil
		stack.Revisions = nil
		stack.Units = nil
		stack.PublishedCount = 0
	}
	return stack
}

func stackListContains(stacks []authoring.StackInfo, bookmark string) bool {
	for _, entry := range stacks {
		if entry.BookmarkName == bookmark {
			return true
		}
	}
	return false
}

func isMergedStack(stack authoring.StackInfo) bool {
	return strings.EqualFold(strings.TrimSpace(stack.Status), "merged")
}

func isTerminalStackStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "merged", "closed":
		return true
	default:
		return false
	}
}

func stackDisplayBucket(stack authoring.StackInfo) int {
	switch strings.ToLower(strings.TrimSpace(stack.Status)) {
	case "published":
		return 1
	case "merged":
		return 2
	default:
		return 0
	}
}

func labelToken(label, value string) string {
	return muted(label) + " " + valueText(value)
}

func valueText(raw string) string {
	if raw == "" {
		return muted("(none)")
	}
	return value(raw)
}

func quoteAgent(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return `""`
	}
	if strings.ContainsAny(value, " \t\n\"") {
		return strconv.Quote(value)
	}
	return value
}

func firstArg(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func newPublishCommand(ctx context.Context, engine *authoring.Engine, use string, hidden bool) *cobra.Command {
	var allowBackwards bool
	var pushToGitHub bool
	var publishAll bool
	cmd := &cobra.Command{
		Use:    use,
		Short:  "Publish accepted GX stacks for review",
		Hidden: hidden,
		Args:   cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			invocation := "gx publish"
			if hidden {
				invocation = "gx pr"
			}
			publishAllRequested := publishAll || len(args) == 0
			if publishAll {
				invocation += " --all"
			} else if len(args) > 0 {
				invocation += " " + args[0]
			}
			fmt.Fprintln(out, commandLine(invocation, false))
			fmt.Fprintln(out)
			engine.SetProgressWriter(out)

			if !pushToGitHub {
				return fmt.Errorf("--github=false is no longer supported; gx publish always pushes stack refs")
			}
			mode := authoring.PublishModeReviewAndGit

			client := cloud.NewClient()
			if client == nil {
				publishHook := func(result authoring.PushResult) error {
					printPublishPush(out, result)
					branch := firstNonEmptyString(result.GXStackRef, pointerString(result.Repo.BranchName))
					fmt.Fprintln(out, labelWarningValue("Warning", fmt.Sprintf("Could not upload GX review context for branch %s: gx cloud is not configured; set GX_CLOUD_URL or rebuild with cloud endpoints. Local publish metadata was still recorded.", branch)))
					return nil
				}

				if publishAllRequested {
					if len(args) > 0 {
						return fmt.Errorf("gx publish --all does not accept a stack name")
					}
					results, err := engine.PublishAll(ctx, nil, authoring.PushOptions{Mode: mode}, publishHook)
					if err != nil {
						return err
					}
					fmt.Fprintln(out, labelValue("Published", fmt.Sprintf("%d stacks", len(results))))
					return nil
				}

				var err error
				if len(args) > 0 {
					_, err = engine.PublishNamed(ctx, args[0], nil, authoring.PushOptions{Mode: mode}, publishHook)
				} else {
					_, err = engine.Publish(ctx, nil, authoring.PushOptions{Mode: mode}, publishHook)
				}
				return err
			}

			if publishAllRequested {
				if len(args) > 0 {
					return fmt.Errorf("gx publish --all does not accept a stack name")
				}
				prepared, err := engine.PrepareAllPublishes(ctx, nil, authoring.PushOptions{Mode: mode})
				if err != nil {
					return err
				}
				published := 0
				for _, push := range prepared {
					review, err := publication.EnqueuePush(ctx, push)
					if err != nil {
						return err
					}
					if err := engine.RecordPublish(ctx, push); err != nil {
						return err
					}
					printPublishPush(out, push)
					printPublishReview(out, review)
					published++
				}
				if published > 0 {
					startPublishUploadWorker(out)
				}
				fmt.Fprintln(out, labelValue("Published", fmt.Sprintf("%d stacks", published)))
				return nil
			}

			var push authoring.PushResult
			var err error
			if len(args) > 0 {
				push, err = engine.PrepareNamedPublish(ctx, args[0], nil, authoring.PushOptions{Mode: mode})
			} else {
				push, err = engine.PreparePublish(ctx, nil, authoring.PushOptions{Mode: mode})
			}
			if err != nil {
				return err
			}
			review, err := publication.EnqueuePush(ctx, push)
			if err != nil {
				return err
			}
			if err := engine.RecordPublish(ctx, push); err != nil {
				return err
			}
			printPublishPush(out, push)
			printPublishReview(out, review)
			startPublishUploadWorker(out)
			return nil
		},
	}
	cmd.Flags().BoolVar(&allowBackwards, "allow-backwards", false, "accepted for compatibility; gx handles required JJ bookmark moves automatically")
	cmd.Flags().BoolVar(&publishAll, "all", false, "accepted for compatibility; gx publish publishes all stacks by default")
	pushToGitHub = true
	cmd.Flags().BoolVar(&pushToGitHub, "github", true, "accepted for compatibility; gx publish always pushes stack refs")
	if flag := cmd.Flags().Lookup("allow-backwards"); flag != nil {
		flag.Hidden = true
	}
	if flag := cmd.Flags().Lookup("github"); flag != nil {
		flag.Hidden = true
	}
	return cmd
}

func printPublishPush(out io.Writer, result authoring.PushResult) {
	if strings.TrimSpace(result.Output) != "" {
		fmt.Fprintln(out, highlightPublishRevisions(strings.TrimSpace(result.Output)))
	}
	if strings.TrimSpace(result.GXStackRef) != "" || result.Repo.BranchName != nil {
		branch := firstNonEmptyString(result.GXStackRef, pointerString(result.Repo.BranchName))
		fmt.Fprintln(out, labelValue("Branch", branch))
		if status := publishGitPushStatusText(result.GitPushStatus); status != "" {
			fmt.Fprintln(out, labelStatus("GitHub", status))
		}
		if prStatus := publishGitHubPRStatusText(result.GitHubPRStatus); prStatus != "" {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(prStatus)), "warn:") {
				fmt.Fprintln(out, labelWarningValue("GitHub PR", strings.TrimSpace(prStatus)))
			} else {
				fmt.Fprintln(out, labelStatus("GitHub PR", prStatus))
			}
		}
		if result.GitHubPullRequestURL != nil && strings.TrimSpace(*result.GitHubPullRequestURL) != "" {
			fmt.Fprintln(out, labelValue("GitHub PR", strings.TrimSpace(*result.GitHubPullRequestURL)))
		}
		if result.Repo.RemoteURL != nil {
			if repoFull := githubRepoFullNameFromRemote(*result.Repo.RemoteURL); repoFull != "" {
				fmt.Fprintln(out, labelValue("Actions", fmt.Sprintf("https://github.com/%s/actions", repoFull)))
			}
		}
	}
	for _, warning := range result.Warnings {
		if strings.TrimSpace(warning) != "" {
			fmt.Fprintln(out, labelWarningValue("Warning", strings.TrimSpace(warning)))
		}
	}
}

func publishGitPushStatusText(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pushed":
		return "ok: pushed"
	case "already up to date":
		return "ok: already pushed"
	case "not pushed":
		return "not pushed"
	default:
		if strings.TrimSpace(status) == "" {
			return ""
		}
		return "warn: " + strings.TrimSpace(status)
	}
}

func publishGitHubPRStatusText(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "created", "existing", "stored":
		return "ok: " + strings.TrimSpace(status)
	case "warning":
		return "warn: not created"
	default:
		return strings.TrimSpace(status)
	}
}

func highlightPublishRevisions(output string) string {
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		matches := publishRevisionOutputLine.FindStringSubmatch(line)
		if len(matches) != 5 {
			continue
		}
		if strings.TrimSpace(matches[4]) == "" {
			continue
		}
		lines[i] = matches[1] + mint(matches[2]) + matches[3] + matches[4]
	}
	return strings.Join(lines, "\n")
}

func printPublishReview(out io.Writer, result publication.Result) {
	if result.Queued {
		fmt.Fprintln(out, labelValue("GX context", "queued for background upload"))
		if result.QueueID != "" {
			fmt.Fprintln(out, labelValue("Upload ID", result.QueueID))
		}
		if result.ArtifactPath != "" {
			fmt.Fprintln(out, labelValue("Local artifact", result.ArtifactPath))
		}
		return
	}
	if !result.Uploaded {
		return
	}
	if result.ReviewID != "" {
		fmt.Fprintln(out, labelValue("Review ID", result.ReviewID))
	}
	if result.ReviewURL != "" {
		fmt.Fprintln(out, labelValue("Review:", result.ReviewURL))
	} else {
		fmt.Fprintln(out, labelValue("Synced", "gx session context"))
	}
	if result.ArtifactPath != "" {
		fmt.Fprintln(out, labelValue("Local artifact", result.ArtifactPath))
	}
	if result.SemanticIndexed {
		if result.SemanticCodeChunks > 0 {
			fmt.Fprintln(out, labelValue("Semantic index", fmt.Sprintf("%d chunks (%d code, %d session)", result.SemanticChunks, result.SemanticCodeChunks, result.SemanticSessionChunks)))
		} else {
			fmt.Fprintln(out, labelValue("Semantic index", fmt.Sprintf("%d transcript chunks", result.SemanticChunks)))
		}
	} else if result.SemanticIndexError != "" {
		fmt.Fprintln(out, labelValue("Semantic index", danger(result.SemanticIndexError)))
	} else if result.IndexStatus != "" {
		fmt.Fprintln(out, labelValue("Semantic index", result.IndexStatus))
	}
}

func newSyncCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var captureOnly bool
	cmd := &cobra.Command{
		Use:   "sync [remote]",
		Short: "Sync remote Git state for the current GX line of work",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			if captureOnly {
				fmt.Fprintln(out, commandLine("gx sync --capture", true))
				fmt.Fprintln(out)
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
				fmt.Fprintf(out, "capture sync: extracts=%d sessions=%d\n",
					result.ExtractsUploaded, result.SessionsUploaded)
				return nil
			}
			fmt.Fprintln(out, commandLine("gx sync", true))
			fmt.Fprintln(out)
			remote := ""
			if len(args) == 1 {
				remote = args[0]
			}
			result, err := engine.Sync(ctx, remote)
			if err != nil {
				return err
			}
			if strings.TrimSpace(result.Output) != "" {
				fmt.Fprintln(out, strings.TrimSpace(result.Output))
			}
			fmt.Fprintln(out, labelValue("Synced", result.RemoteName))

			if summary, err := syncCloudMetadata(ctx, engine, result.Repo, out); err != nil {
				return err
			} else if summary != nil {
				fmt.Fprintln(out, labelValue("Cloud bookmarks", fmt.Sprintf("%d listed", summary.Listed)))
				if summary.Updated > 0 {
					fmt.Fprintln(out, labelValue("Cloud revisions", fmt.Sprintf("%d updated", summary.Updated)))
				}
				if summary.CaughtUp > 0 {
					fmt.Fprintln(out, labelValue("Cloud catch-up", fmt.Sprintf("%d bookmark(s) fetched from remote", summary.CaughtUp)))
				}
			}
			if status, err := publication.QueuedUploadStatus(); err == nil && (status.Pending > 0 || status.Failed > 0) {
				if err := drainPublishUploadOutbox(ctx, out, false, 20); err != nil {
					fmt.Fprintln(out, labelWarningValue("GX context uploads", err.Error()))
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&captureOnly, "capture", false, "drain pending capture staging uploads")
	return cmd
}

func newPublishUploadCommand(ctx context.Context) *cobra.Command {
	var quiet bool
	var limit int
	cmd := &cobra.Command{
		Use:    "__gx-upload-outbox",
		Short:  "Upload queued GX publish context",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return drainPublishUploadOutbox(ctx, cmd.OutOrStdout(), quiet, limit)
		},
	}
	cmd.Flags().BoolVar(&quiet, "quiet", false, "suppress upload summary")
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum queued artifacts to upload")
	return cmd
}

func drainPublishUploadOutbox(ctx context.Context, out io.Writer, quiet bool, limit int) error {
	client := cloud.NewClient()
	if client == nil {
		return fmt.Errorf("gx cloud is not configured; set GX_CLOUD_URL or rebuild with cloud endpoints")
	}
	result, err := publication.DrainQueuedUploads(ctx, client, limit)
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Fprintln(out, labelValue("GX context uploads", fmt.Sprintf("%d uploaded, %d failed, %d pending", result.Uploaded, result.Failed, result.Pending)))
	}
	return nil
}

func startPublishUploadWorker(out io.Writer) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(out, labelWarningValue("GX context upload", "queued; could not find gx executable, run `gx sync` to upload"))
		return
	}
	cmd := exec.Command(exe, "__gx-upload-outbox", "--quiet")
	if devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0); err == nil {
		defer devNull.Close()
		cmd.Stdin = devNull
	}
	if logFile, err := publishUploadLogFile(); err == nil {
		defer logFile.Close()
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	}
	cmd.Env = os.Environ()
	configureDetachedCommand(cmd)
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(out, labelWarningValue("GX context upload", "queued; run `gx sync` to upload"))
		return
	}
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	fmt.Fprintln(out, labelValue("GX context upload", "background upload started"))
}

func publishUploadLogFile() (*os.File, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(dir, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(filepath.Join(logDir, "publish-upload.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
}

func newOpsCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "ops",
		Short:  "Advanced commands for capture, diagnostics, and ingest",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			printOpsSummary(cmd.OutOrStdout())
		},
	}
	cmd.AddCommand(
		newOpsCaptureCommand(ctx),
		newOpsDiagCommand(ctx),
		newOpsIngestCommand(ctx),
	)
	return cmd
}

func printOpsSummary(out io.Writer) {
	fmt.Fprintln(out, commandLine("gx ops", true))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Advanced commands for capture, diagnostics, and ingest")
	fmt.Fprintln(out)
	fmt.Fprintln(out, opsMenuRow("capture", "Record a gx session for replay"))
	fmt.Fprintln(out, opsMenuRow("diagnose", "Inspect local gx state and config"))
	fmt.Fprintln(out, opsMenuRow("ingest", "Import external change metadata"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, muted(`Use "gx ops [command] --help" for details.`))
}

func opsMenuRow(name, description string) string {
	return "  " + logoText(padRight(name, 11)) + " " + muted(description)
}

func newOpsCaptureCommand(ctx context.Context) *cobra.Command {
	cmd := newServiceCommand(ctx)
	cmd.Use = "capture"
	cmd.Short = "Record a gx session for replay"
	cmd.Aliases = []string{"service"}
	return cmd
}

func newOpsDiagCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "diagnose",
		Aliases: []string{"diag"},
		Short:   "Inspect local gx state and config",
	}
	doctor := newDoctorCommand(ctx)
	doctor.Use = "doctor"
	repair := newRepairCommand(ctx)
	repair.Use = "repair"
	cmd.AddCommand(doctor, repair)
	return cmd
}

func newOpsIngestCommand(ctx context.Context) *cobra.Command {
	cmd := newIngestCommand(ctx)
	cmd.Use = "ingest"
	cmd.Short = "Import external change metadata"
	return cmd
}

func Execute(ctx context.Context) error {
	args := os.Args[1:]
	if ShouldLaunch(args) {
		if err := launcher.Run(ctx, args); err != nil {
			if !errors.Is(err, context.Canceled) {
				fmt.Fprintln(os.Stderr, err)
			}
			return err
		}
		return nil
	}

	root := NewRoot(ctx)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		if !errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, err)
		}
		return err
	}
	return nil
}

func githubRepoFullNameFromRemote(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if idx := strings.Index(remoteURL, "github.com"); idx >= 0 {
		suffix := remoteURL[idx+len("github.com"):]
		suffix = strings.TrimPrefix(suffix, ":")
		suffix = strings.TrimPrefix(suffix, "/")
		suffix = strings.TrimSuffix(suffix, ".git")
		parts := strings.Split(suffix, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1]
		}
	}
	return ""
}

func shortID(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
