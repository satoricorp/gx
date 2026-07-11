package cli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/clitui"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/commitcontext"
	"github.com/satoricorp/gx/internal/daemon"
	"github.com/satoricorp/gx/internal/github"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/inference"
	"github.com/satoricorp/gx/internal/launcher"
	"github.com/satoricorp/gx/internal/postlist"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.Contains(cmd.CommandPath(), "__") {
				return nil
			}
			inference.ApplyToEnvironment()
			telemetry.EmitInstallOnce(ctx)
			return ensureAutoInitializedRepo(ctx, engine, cmd)
		},
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
		&cobra.Group{ID: groupHelp, Title: "Help:"},
		&cobra.Group{ID: groupAdvanced, Title: "Advanced:"},
	)

	root.AddCommand(
		newInternalDaemonCommand(ctx),
		newVersionCommand(),
		newDoctorCommand(ctx),
		newAuthCommand(ctx),
		newSetCommand(ctx),
		newInitCommand(ctx, engine),
		newCaptureCommand(ctx),
		newPublishUploadCommand(ctx),
		newBaseCommand(ctx, engine),
		newDemoCommand(),
		newGenerateCommand(ctx, engine),
		newCommitCommand(ctx, engine),
		newEditCommand(ctx, engine),
		newStatusCommand(ctx, engine, "status", "status", false),
		newReportCommand(ctx, engine),
		newReviewCommand(ctx, engine),
		newPushCommand(ctx, engine),
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
		case "init", "auth", "set", "demo":
			cmd.GroupID = groupSetup
		case "base", "commit", "edit", "review", "status":
			cmd.GroupID = groupWork
		case "push", "sync":
			cmd.GroupID = groupShip
		case "doctor", "report", "version":
			cmd.GroupID = groupHelp
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
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "gx version",
		Run: func(cmd *cobra.Command, args []string) {
			if jsonOut {
				if err := json.NewEncoder(cmd.OutOrStdout()).Encode(version.BuildInfo()); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err)
				}
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), logoText("version "+version.Current()))
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable version metadata")
	return cmd
}

func newInitCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var name string
	var email string
	var yes bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Set up gx in the current repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout(), commandLine("gx init", true))
				fmt.Fprintln(cmd.OutOrStdout())
			}
			result, err := engine.Init(ctx, authoring.InitOptions{
				Name:        name,
				Email:       email,
				Interactive: !yes,
				In:          initInput(cmd, yes),
				Out:         initOutput(cmd, yes),
			})
			if err != nil {
				return err
			}
			if !yes && result.IdentityName != "" && result.IdentityEmail != "" {
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
			installHook := installCaptureHook
			if yes {
				installHook = installCaptureHookQuiet
			}
			if err := installHook(cmd, result.Repo.RootPath); err != nil {
				if !yes {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Pre-push hook", err.Error()))
				}
			}
			if err := installClaudeCaptureHooks(cmd, result.Repo.RootPath, yes); err != nil {
				if !yes {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Claude hooks", err.Error()))
				}
			}
			if err := gxconfig.EnsureCaptureRetention(); err != nil {
				if !yes {
					fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Retention", err.Error()))
				}
			} else {
				cfg, _ := gxconfig.Load()
				if !yes {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Session retention", fmt.Sprintf("%d days", cfg.Capture.CleanupPeriodDays)))
				}
			}
			if err := initSetupFromCommand(cmd, result.Repo.RootPath, yes); err != nil {
				return err
			}
			if !yes {
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Next", "git add <files> && gx commit -m \"...\""))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "user name to store in GX and JJ config")
	cmd.Flags().StringVar(&email, "email", "", "user email to store in GX and JJ config")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "accept defaults and suppress successful init output")
	return cmd
}

func initInput(cmd *cobra.Command, quiet bool) io.Reader {
	if quiet {
		return nil
	}
	return cmd.InOrStdin()
}

func initOutput(cmd *cobra.Command, quiet bool) io.Writer {
	if quiet {
		return nil
	}
	return cmd.OutOrStdout()
}

func newCommitCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var message string
	var branch string
	var jsonOut bool
	var all bool
	var amend bool
	var fileMessage string
	var contextFile string
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Record staged Git changes as a GX revision",
		Long: strings.Join([]string{
			"Record staged Git changes as a GX revision.",
			"",
			"Like git commit, the revision lands on the current branch and HEAD",
			"does not move. Pass -b to create and switch to a new branch first.",
			"",
			"Use git add or git add -p to choose scope, then run gx commit -m.",
			"The resulting revision is JJ-backed and editable with GX/MCP tools.",
			"",
			"Exit codes: 0 success, 1 general error, 2 no staged changes.",
		}, "\n"),
		Example: strings.Join([]string{
			`  git add internal/cli/root.go`,
			`  gx commit -m "record staged CLI change"`,
			`  git add -p`,
			`  gx commit -m "record selected hunks"`,
		}, "\n"),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("gx commit does not accept path arguments; choose scope with git add or git add -p")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				return fmt.Errorf("gx commit -a is not supported; stage changes with git add first")
			}
			if amend {
				return fmt.Errorf("gx commit --amend is not supported; use gx edit <rev> to modify a recorded revision")
			}
			if strings.TrimSpace(fileMessage) != "" {
				return fmt.Errorf("gx commit -F is not supported yet; pass a message with -m")
			}
			if strings.TrimSpace(message) == "" {
				return fmt.Errorf("commit message is required; pass -m \"describe this revision\"")
			}
			var selfReport commitcontext.SelfReport
			if strings.TrimSpace(contextFile) != "" {
				report, err := commitcontext.LoadFile(contextFile)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
				} else {
					selfReport = report
				}
			}
			result, err := engine.CommitStaged(ctx, authoring.CommitStagedOptions{
				Message:    message,
				Branch:     branch,
				SelfReport: selfReport,
			})
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(cmd, commitResultJSON{
					Result:           result,
					CreatedBranch:    result.CreatedBranch,
					ProvenanceStatus: result.ProvenanceStatus,
				})
			}
			if strings.TrimSpace(result.Output) != "" {
				fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(result.Output))
			}
			printCommitSummary(cmd.OutOrStdout(), result)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message (required)")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "create this branch at HEAD and record the revision onto it (default: current branch)")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "unsupported; stage changes with git add first")
	cmd.Flags().BoolVar(&amend, "amend", false, "unsupported; use gx edit <rev>")
	cmd.Flags().StringVarP(&fileMessage, "file", "F", "", "unsupported; pass a message with -m")
	cmd.Flags().StringVar(&contextFile, "context-file", "", "optional JSON file with agent-declared provenance (task_summary, commands_run, tests_run)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

type commitResultJSON struct {
	Result           authoring.CheckpointResult `json:"result"`
	CreatedBranch    bool                       `json:"created_branch"`
	ProvenanceStatus string                     `json:"provenance_status"`
}

func printCommitSummary(out io.Writer, result authoring.CheckpointResult) {
	invocation := "gx commit"
	if result.Change.Description != "" {
		invocation = fmt.Sprintf("gx commit -m %q", result.Change.Description)
	}
	fmt.Fprintln(out, commandLine(invocation, false))
	fmt.Fprintln(out)
	fmt.Fprintln(out, success("Revision recorded"))
	if result.CreatedBranch && result.Stack != nil {
		fmt.Fprintln(out, labelValue("Created branch", fmt.Sprintf("%s (from %s)", result.Stack.BookmarkName, result.Stack.BaseRef)))
	}
	fmt.Fprintln(out, labelValue("Message", result.Change.Description))
	fmt.Fprintln(out, labelValue("Revision", shortID(result.Change.ChangeID, 12)))
	if result.Change.CommitID != "" {
		fmt.Fprintln(out, labelValue("Commit", shortID(result.Change.CommitID, 8)))
	}
	if result.Stack != nil {
		fmt.Fprintln(out, labelValue("Stack", result.Stack.BookmarkName))
	}
	fmt.Fprintln(out, labelValue("Provenance", firstNonEmptyString(result.ProvenanceStatus, "absent")))
	fmt.Fprintln(out)
	fmt.Fprintln(out, hint("Choose the next manual scope with git add or git add -p."))
	fmt.Fprintln(out, labelValue("Next", fmt.Sprintf("gx commit -m %q", result.Change.Description)))
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

func newGenerateCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var intent string
	var raw bool
	var model string
	var maxWarnings int
	var excludeFilesets []string
	var legacy bool
	cmd := &cobra.Command{
		Use:     "generate [filesets...]",
		Short:   "Save your work in branches & commits",
		Hidden:  true,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				switch strings.TrimSpace(args[0]) {
				case "apply", "accept", "review", "fix":
					return fmt.Errorf("unknown command %q for %q", args[0], "gx generate")
				}
			}
			startedAt := time.Now()
			for round := 1; ; round++ {
				packet, result, err := runGenerateApply(ctx, engine, cmd, generateRunOptions{
					JSON:              jsonOut,
					Raw:               raw,
					Intent:            intent,
					Model:             model,
					MaxWarnings:       maxWarnings,
					Filesets:          args,
					ExcludeFilesets:   excludeFilesets,
					PreflightAttempts: generatePreflightAttempts(),
					Legacy:            legacy,
				})
				emitGenerateRunTelemetry(ctx, packet, result, err, generateRunTelemetryOptions{
					JSON:         jsonOut,
					Raw:          raw,
					Filesets:     args,
					ExcludeCount: len(excludeFilesets),
					HasIntent:    strings.TrimSpace(intent) != "",
					ModelSet:     strings.TrimSpace(model) != "",
					Duration:     time.Since(startedAt),
				})
				if err == nil {
					return nil
				}
				if jsonOut || generateIsMCP() || !useStatusInteractive(cmd.InOrStdin(), cmd.OutOrStdout()) {
					return err
				}
				if strings.TrimSpace(packet.Proposal.ID) == "" {
					return err
				}
				printDemuxChangesPacket(cmd.OutOrStdout(), packet, demuxPrintOptions{Raw: raw})
				printGenerateRoundFailure(cmd.OutOrStdout(), err)
				continued := promptContinueGenerate(cmd.InOrStdin(), cmd.OutOrStdout())
				telemetry.EmitProductEvent(ctx, telemetry.EventCLIGeneratePrompt, map[string]any{
					"round":     round,
					"continued": continued,
				})
				if !continued {
					return err
				}
			}
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().StringVar(&intent, "intent", "", "optional intent prefix for generated revisions")
	cmd.Flags().BoolVar(&raw, "raw", false, "show full generate diagnostics in human output")
	cmd.Flags().StringVar(&model, "model", "", "OpenAI model for generate repair; defaults to GX_DEMUX_REVIEW_MODEL or gpt-4.1-mini")
	cmd.Flags().IntVar(&maxWarnings, "max-warnings", 0, "maximum warning-severity diagnostics to send to the model; defaults to 50")
	cmd.Flags().StringArrayVar(&excludeFilesets, "exclude", nil, "exclude changed file or directory from generate; may be repeated")
	cmd.Flags().BoolVar(&legacy, "legacy", false, "run the legacy generate repair and preflight pipeline")
	return cmd
}

type generateRunOptions struct {
	JSON              bool
	Raw               bool
	Intent            string
	Model             string
	MaxWarnings       int
	Filesets          []string
	ExcludeFilesets   []string
	PreflightAttempts int
	Legacy            bool
}

func runGenerateApply(ctx context.Context, engine *authoring.Engine, cmd *cobra.Command, opts generateRunOptions) (authoring.DemuxPlanPacket, authoring.ApplyDemuxResult, error) {
	run := func(progress io.Writer) (authoring.DemuxPlanPacket, error) {
		if err := ensureGenerateInitialized(ctx, engine, cmd, opts); err != nil {
			return authoring.DemuxPlanPacket{}, err
		}
		if err := engine.RequireAuthoringBase(ctx, "gx generate"); err != nil {
			return authoring.DemuxPlanPacket{}, err
		}
		return engine.GenerateChanges(ctx, authoring.ProposeDemuxOptions{
			Intent:                 opts.Intent,
			Filesets:               opts.Filesets,
			ExcludeFilesets:        opts.ExcludeFilesets,
			Legacy:                 opts.Legacy,
			Model:                  opts.Model,
			MaxWarnings:            opts.MaxWarnings,
			ApplyPreflightAttempts: opts.PreflightAttempts,
			ProgressWriter:         progress,
		})
	}
	var (
		packet authoring.DemuxPlanPacket
		result authoring.ApplyDemuxResult
		err    error
	)
	if !opts.JSON && useDemuxLoader(cmd.InOrStdin(), cmd.ErrOrStderr()) {
		packet, err = runDemuxWithLoader(cmd.InOrStdin(), cmd.ErrOrStderr(), run)
	} else {
		var progress io.Writer
		if !opts.JSON {
			progress = cmd.ErrOrStderr()
		}
		packet, err = run(progress)
	}
	if err != nil {
		if opts.JSON && strings.TrimSpace(packet.Proposal.ID) != "" {
			_ = writeJSON(cmd, map[string]any{
				"error":  err.Error(),
				"packet": packet,
			})
		}
		return packet, authoring.ApplyDemuxResult{}, err
	}
	if reason := demuxAutoAcceptBlockedReason(packet); reason != "" {
		if opts.JSON {
			_ = writeJSON(cmd, map[string]any{
				"error":  reason,
				"packet": packet,
			})
		}
		return packet, authoring.ApplyDemuxResult{}, fmt.Errorf("%s", reason)
	}
	apply := func() (authoring.ApplyDemuxResult, error) {
		return engine.ApplyDemuxPlanWithOptions(ctx, packet.Proposal, authoring.ApplyDemuxOptions{
			ReturnToDefaultBranch: true,
		})
	}
	if !opts.JSON && useDemuxLoader(cmd.InOrStdin(), cmd.ErrOrStderr()) {
		result, err = runDemuxApplyWithLoader(cmd.InOrStdin(), cmd.ErrOrStderr(), "Creating revision stacks...", apply)
	} else {
		result, err = apply()
	}
	if err != nil {
		return packet, result, err
	}
	if opts.JSON {
		return packet, result, writeJSON(cmd, result)
	}
	printDemuxApplySummary(cmd.OutOrStdout(), result)
	return packet, result, nil
}

func ensureGenerateInitialized(ctx context.Context, engine *authoring.Engine, cmd *cobra.Command, opts generateRunOptions) error {
	_, err := engine.EnsureReadyRepo(ctx)
	return err
}

// printGenerateRoundFailure surfaces the error that stopped a generate round
// before the interactive continue prompt. Without it the prompt appears with
// no explanation, which reads like a hang.
func printGenerateRoundFailure(out io.Writer, err error) {
	fmt.Fprintln(out)
	if err != nil {
		fmt.Fprintln(out, labelWarningValue("Error", err.Error()))
	}
}

func generatePreflightAttempts() int {
	if generateIsMCP() {
		return 6
	}
	return 3
}

func generateIsMCP() bool {
	return strings.TrimSpace(os.Getenv("GX_MCP")) != ""
}

func printDemuxApplySummary(out io.Writer, result authoring.ApplyDemuxResult) {
	fmt.Fprintln(out, success("Generated revisions applied"))
	fmt.Fprintln(out, labelValue("Plan", result.Proposal.ID))
	fmt.Fprintln(out, labelValue("Created", fmt.Sprintf("%d revisions", len(result.Revisions))))
	for index, revision := range result.Revisions {
		fmt.Fprintf(out, "  %s %s  %s\n", command(fmt.Sprintf("r%d", index+1)), value(revision.Change.Description), muted(shortID(revision.Change.ChangeID, 12)))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Next"))
	fmt.Fprintf(out, "  %s\n", command("gx status"))
	fmt.Fprintf(out, "  %s\n", command("gx push"))
	if result.RemainingChanges {
		printDemuxPartialFollowup(out)
	}
}

func printDemuxPartialFollowup(out io.Writer) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, section("Remaining changes"))
	fmt.Fprintln(out, muted("Selected generated changes were applied. Run "+command("gx generate")+" again for the remaining changes."))
}

const demuxPartialComposeWarning = "Partial generate plan: these revisions do not cover all current changes. Run gx generate again for the remaining changes."

type demuxPrintOptions struct {
	Raw bool
}

func printDemuxChangesPacket(out io.Writer, packet authoring.DemuxPlanPacket, opts demuxPrintOptions) {
	fmt.Fprintln(out, commandLine("gx generate", true))
	fmt.Fprintln(out)
	printDemuxPartialWarning(out, packet.Proposal)
	printDemuxProposal(out, packet.Proposal, opts)
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

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
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
	fmt.Fprintln(out, labelValue("JSON", "gx generate --json"))
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
	fmt.Fprintln(out, section("Generate plan"))
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

type generateRunTelemetryOptions struct {
	JSON         bool
	Raw          bool
	Filesets     []string
	ExcludeCount int
	HasIntent    bool
	ModelSet     bool
	Duration     time.Duration
}

func emitGenerateRunTelemetry(ctx context.Context, packet authoring.DemuxPlanPacket, result authoring.ApplyDemuxResult, runErr error, opts generateRunTelemetryOptions) {
	blockingWarnings, diagnosticWarnings := splitFeasibilityWarnings(packet.Proposal.FeasibilityWarnings)
	status := "success"
	if runErr != nil {
		status = "error"
	}
	props := map[string]any{
		"status":                   status,
		"state":                    string(packet.State),
		"revision_count":           len(packet.Proposal.Revisions),
		"generated_revision_count": len(result.Revisions),
		"stack_count":              len(demuxStackDisplayGroups(packet.Proposal.Revisions)),
		"hunk_count":               len(packet.Proposal.Hunks),
		"file_count":               demuxProposalFileCount(packet.Proposal),
		"warning_count":            len(packet.Proposal.Warnings) + len(packet.Proposal.FeasibilityWarnings),
		"blocking_warning_count":   len(blockingWarnings),
		"diagnostic_warning_count": len(diagnosticWarnings),
		"repair_hint_count":        len(packet.Review.RepairHints),
		"review_error_count":       len(packet.Review.Errors),
		"json":                     opts.JSON,
		"raw":                      opts.Raw,
		"fileset_count":            len(opts.Filesets),
		"exclude_count":            opts.ExcludeCount,
		"has_intent":               opts.HasIntent,
		"model_set":                opts.ModelSet,
		"duration_ms":              opts.Duration.Milliseconds(),
	}
	if packet.Proposal.HunkCoverage > 0 {
		props["hunk_coverage"] = packet.Proposal.HunkCoverage
	}
	telemetry.EmitProductEvent(ctx, telemetry.EventCLIGenerateRun, props)
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

func demuxAutoAcceptBlockedReason(packet authoring.DemuxPlanPacket) string {
	if packet.State != authoring.DemuxWorkflowReadyToApply {
		return fmt.Sprintf("generated revisions are not ready: %s", packet.State)
	}
	if strings.TrimSpace(packet.Proposal.ID) == "" {
		return "generated revisions are not ready: missing proposal id"
	}
	if len(packet.Review.Errors) > 0 {
		return "generated revisions are not ready: review errors found"
	}
	if len(packet.Review.RepairHints) > 0 {
		return "generated revisions are not ready: repair hints found"
	}
	blockingWarnings, _ := splitFeasibilityWarnings(packet.Proposal.FeasibilityWarnings)
	if len(blockingWarnings) > 0 {
		return "generated revisions are not ready: blocking warnings found"
	}
	if len(demuxBlockingPlainWarnings(packet.Proposal.Warnings)) > 0 {
		return "generated revisions are not ready: blocking warnings found"
	}
	return ""
}

func demuxBlockingPlainWarnings(warnings []string) []string {
	var blocking []string
	for _, warning := range warnings {
		warning = strings.TrimSpace(warning)
		switch {
		case warning == "":
			continue
		case warning == demuxPartialComposeWarning:
			blocking = append(blocking, warning)
		case strings.HasPrefix(warning, "Compose repair failed:"):
			blocking = append(blocking, warning)
		case strings.HasPrefix(warning, "Compose apply preflight failed:"):
			blocking = append(blocking, warning)
		}
	}
	return blocking
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
	var showAll bool
	var interactive bool
	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{"gxs"},
		Short:   "Show unstaged changes, local and remote stacks",
		Hidden:  hidden,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Status is read-only: protect the user's staged git index from
			// jj's index rewriting while we inspect the repo.
			return engine.PreservingGitIndex(ctx, func() error {
				if jsonOut {
					stack, _, err := statusWithMissingBaseCheck(ctx, engine, cmd.OutOrStdout(), jsonOut)
					if enrichErr := enrichStatusStack(ctx, engine, &stack); enrichErr != nil && err == nil {
						return enrichErr
					}
					stack = stackSummaryForStacksDisplay(stack, stackDisplayOptions{ShowAll: showAll})
					if writeErr := writeJSON(cmd, stack); writeErr != nil {
						return writeErr
					}
					return err
				}
				if len(args) > 0 && !agentOut {
					return fmt.Errorf("stack selector is only supported with --agent")
				}
				return printStacks(ctx, engine, cmd.InOrStdin(), cmd.OutOrStdout(), agentOut, firstArg(args), stackDisplayOptions{ShowAll: showAll}, interactive, jsonOut)
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&agentOut, "agent", false, "print stable agent-readable text")
	cmd.Flags().BoolVar(&showAll, "all", false, "list all changed files in git working tree sections")
	cmd.Flags().BoolVar(&showAll, "show-all", false, "deprecated alias for --all")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "open the interactive status browser")
	if flag := cmd.Flags().Lookup("show-all"); flag != nil {
		flag.Hidden = true
	}
	cmd.AddCommand(
		newStatusListCommand(ctx, engine),
		newStacksEditCommand(ctx, engine),
		newStacksDiffCommand(),
	)
	setHelpName(cmd, helpName)
	return cmd
}

func newStatusListCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	var agentOut bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all non-merged GX features",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return engine.PreservingGitIndex(ctx, func() error {
				stack, prunedEmpty, err := statusWithMissingBaseCheck(ctx, engine, cmd.OutOrStdout(), jsonOut)
				if enrichErr := enrichStatusStack(ctx, engine, &stack); enrichErr != nil && err == nil {
					return enrichErr
				}
				stack = stackSummaryForStacksDisplay(stack, stackDisplayOptions{ShowEmpty: true})
				if jsonOut {
					if writeErr := writeJSON(cmd, stack); writeErr != nil {
						return writeErr
					}
					return err
				}
				if err != nil {
					return err
				}
				if agentOut {
					printStatusAgent(cmd.OutOrStdout(), stack, "")
					return nil
				}
				printDeletedEmptyStacksNotice(cmd.OutOrStdout(), prunedEmpty)
				fmt.Fprint(cmd.OutOrStdout(), renderStacksSummary(stack, nil, currentStackIndex(stack), true, 0))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&agentOut, "agent", false, "print stable agent-readable text")
	return cmd
}

func newReviewCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var scope string
	var focus string
	var deep bool
	var verbose bool
	cmd := &cobra.Command{
		Use:     "review [prompt]",
		Aliases: []string{"gxr"},
		Short:   "Review changes based on codebase & session context, along with independent resources",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()
			var runErr error
			defer func() {
				if runErr != nil {
					autoReportFailure(ctx, engine, runErr, "gx review")
				}
			}()
			reviewScope := ""
			scopeExplicit := cmd.Flags().Changed("scope")
			if scopeExplicit {
				reviewScope = scope
			}
			prompt := ""
			if len(args) > 0 {
				prompt = strings.TrimSpace(args[0])
			}
			repo, err := vcs.NewService().ResolveGitRepo(ctx)
			if err != nil {
				runErr = err
				emitReviewRunTelemetry(ctx, codereview.Report{}, err, reviewScope, scopeExplicit, focus, prompt, deep, verbose, time.Since(startedAt))
				return runErr
			}
			report, err := runReviewWithLoader(
				cmd.InOrStdin(),
				cmd.ErrOrStderr(),
				func(progress io.Writer) (codereview.Report, error) {
					return codereview.Review(ctx, repo.RootPath, codereview.Options{
						Scope:          reviewScope,
						Deep:           deep,
						Focus:          focus,
						Prompt:         prompt,
						Verbose:        verbose,
						ProgressWriter: progress,
						Color:          true,
					})
				},
			)
			emitReviewRunTelemetry(ctx, report, err, reviewScope, scopeExplicit, focus, prompt, deep, verbose, time.Since(startedAt))
			if err != nil {
				runErr = err
				return runErr
			}
			fmt.Fprint(cmd.OutOrStdout(), codereview.RenderMarkdown(report))
			postReviewSummaryComment(ctx, repo, report, cmd.ErrOrStderr())
			recordReviewHistory(ctx, repo, report, prompt, scopeExplicit, deep, cmd.ErrOrStderr())
			return nil
		},
	}
	cmd.Flags().StringVar(&scope, "scope", codereview.DefaultScope, "review scope when explicitly set: architecture, security, performance, onboarding, docs, dependencies, testing, maintainability")
	cmd.Flags().StringVar(&focus, "focus", "", "limit review to files under this path prefix")
	cmd.Flags().BoolVar(&deep, "deep", false, "run full-spectrum review with more local and indexed context")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "include repo facts, docs, and changed files")
	return cmd
}

func emitReviewRunTelemetry(ctx context.Context, report codereview.Report, runErr error, reviewScope string, scopeExplicit bool, focus string, prompt string, deep bool, verbose bool, duration time.Duration) {
	status := "success"
	if runErr != nil {
		status = "error"
	}
	mode := "patch"
	if deep {
		mode = "deep"
	} else if strings.TrimSpace(prompt) != "" {
		mode = "prompt"
	} else if scopeExplicit {
		mode = "scope"
	}
	effectiveScope := strings.TrimSpace(report.Scope)
	if effectiveScope == "" {
		effectiveScope = strings.TrimSpace(reviewScope)
	}
	if effectiveScope == "" {
		effectiveScope = codereview.DefaultScope
	}
	props := map[string]any{
		"status":                status,
		"mode":                  mode,
		"scope":                 effectiveScope,
		"scope_explicit":        scopeExplicit,
		"has_focus":             strings.TrimSpace(focus) != "",
		"has_prompt":            strings.TrimSpace(prompt) != "",
		"deep":                  deep,
		"verbose":               verbose,
		"duration_ms":           duration.Milliseconds(),
		"finding_count":         len(report.Findings),
		"changed_file_count":    len(report.ChangedFiles),
		"tracked_file_count":    report.TrackedFileCount,
		"test_file_count":       report.TestFileCount,
		"context_snippet_count": report.ContextSnippets,
	}
	if strings.TrimSpace(report.Reviewer) != "" {
		props["reviewer"] = report.Reviewer
	}
	telemetry.EmitProductEvent(ctx, telemetry.EventCLIReviewRun, props)
}

const reviewCommentMarker = "<!-- gx review summary -->"

func postReviewSummaryComment(ctx context.Context, repo vcs.RepoInfo, report codereview.Report, stderr io.Writer) {
	remoteURL := pointerString(repo.RemoteURL)
	branchName := pointerString(repo.BranchName)
	if strings.TrimSpace(remoteURL) == "" || strings.TrimSpace(branchName) == "" {
		return
	}
	host, owner, repoName, ok := githubRemoteTarget(remoteURL)
	if !ok {
		return
	}
	client, err := github.NewClient(host)
	if err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not post GX review comment: %v", err)))
		return
	}
	pr, err := client.FindPullRequest(ctx, github.CreatePullRequestOptions{
		Host:       host,
		Owner:      owner,
		Repo:       repoName,
		HeadBranch: branchName,
	})
	if err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not find pull request for %s: %v", branchName, err)))
		return
	}
	if pr == nil || pr.Number == 0 {
		return
	}
	if err := postReviewInlineComments(ctx, client, repo, owner, repoName, pr.Number, report); err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not post GX inline review comment: %v", err)))
	}
	commentBody := reviewCommentMarker + "\n" + codereview.RenderMarkdown(report)
	_, err = client.UpsertIssueComment(ctx, github.IssueCommentOptions{
		Owner:  owner,
		Repo:   repoName,
		Number: pr.Number,
		Body:   commentBody,
		Marker: reviewCommentMarker,
	})
	if err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not post GX review comment: %v", err)))
	}
}

func postReviewInlineComments(ctx context.Context, client *github.Client, repo vcs.RepoInfo, owner string, repoName string, prNumber int, report codereview.Report) error {
	if client == nil || prNumber <= 0 || len(report.Findings) == 0 {
		return nil
	}
	repoRoot := strings.TrimSpace(repo.RootPath)
	if repoRoot == "" {
		repoRoot = strings.TrimSpace(report.RepoRoot)
	}
	commitID := currentHeadCommit(ctx, repoRoot)
	if strings.TrimSpace(commitID) == "" {
		return nil
	}
	seen := map[string]struct{}{}
	var firstErr error
	for _, finding := range report.Findings {
		for _, anchor := range finding.Anchors {
			if !codereview.AnchorMapsToChangedHunk(ctx, repoRoot, anchor) {
				continue
			}
			key := fmt.Sprintf("%s:%d:%s", anchor.File, anchor.Line, finding.ID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if err := client.CreatePullRequestReviewComment(ctx, github.PullRequestReviewCommentOptions{
				Owner:    owner,
				Repo:     repoName,
				Number:   prNumber,
				Body:     renderInlineReviewComment(finding),
				CommitID: commitID,
				Path:     anchor.File,
				Line:     anchor.Line,
				Side:     "RIGHT",
			}); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func renderInlineReviewComment(finding codereview.Finding) string {
	var b strings.Builder
	title := strings.TrimSpace(finding.Title)
	if title == "" {
		title = "GX review finding"
	}
	fmt.Fprintf(&b, "**%s**\n\n", title)
	if summary := strings.TrimSpace(finding.Summary); summary != "" {
		fmt.Fprintf(&b, "%s\n\n", summary)
	}
	if recommendation := strings.TrimSpace(finding.Recommendation); recommendation != "" {
		fmt.Fprintf(&b, "**Do next:** %s\n", recommendation)
	}
	return strings.TrimSpace(b.String())
}

func recordReviewHistory(ctx context.Context, repo vcs.RepoInfo, report codereview.Report, prompt string, scopeExplicit bool, deep bool, stderr io.Writer) {
	remoteURL := pointerString(repo.RemoteURL)
	repoFullName := cloud.RepoFullNameFromRemoteURL(remoteURL)
	if strings.TrimSpace(repoFullName) == "" {
		return
	}
	client := cloud.NewClient()
	if client == nil {
		return
	}
	mode := "patch"
	if deep {
		mode = "deep"
	} else if strings.TrimSpace(prompt) != "" {
		mode = "prompt"
	} else if scopeExplicit {
		mode = "scope"
	}
	req := cloud.CodeReviewHistoryRecordRequest{
		RepoRootPath: report.RepoRoot,
		RepoFullName: repoFullName,
		BranchName:   pointerString(repo.BranchName),
		HeadCommitID: currentHeadCommit(ctx, repo.RootPath),
		SourceKind:   "session_intent",
		Prompt:       prompt,
		Scope:        firstNonEmptyString(report.Scope, codereview.DefaultScope),
		Mode:         mode,
		Reviewer:     report.Reviewer,
		SummaryKind:  "pr",
		SummaryText:  codereview.RenderMarkdown(report),
		Findings:     reviewHistoryFindings(repoFullName, report),
		Payload: map[string]any{
			"changed_files":  len(report.ChangedFiles),
			"deep":           deep,
			"scope_explicit": scopeExplicit,
		},
	}
	if _, err := client.RecordCodeReviewHistory(ctx, req); err != nil {
		fmt.Fprintln(stderr, labelWarningValue("Warning", fmt.Sprintf("Could not record code review history: %v", err)))
	}
}

func reviewHistoryFindings(repoFullName string, report codereview.Report) []cloud.CodeReviewFindingRecord {
	findings := make([]cloud.CodeReviewFindingRecord, 0, len(report.Findings))
	for _, finding := range report.Findings {
		file, line := primaryFindingLocation(finding, report.ChangedFiles)
		category := "general"
		if len(finding.Scopes) > 0 && strings.TrimSpace(finding.Scopes[0]) != "" {
			category = strings.TrimSpace(finding.Scopes[0])
		} else if before, _, ok := strings.Cut(finding.ID, "."); ok && strings.TrimSpace(before) != "" {
			category = strings.TrimSpace(before)
		}
		record := cloud.CodeReviewFindingRecord{
			Fingerprint:    reviewFindingFingerprint(repoFullName, finding, file, line, category),
			Outcome:        "valid",
			Category:       category,
			Language:       reviewHistoryLanguageForFile(file),
			FilePath:       file,
			LineStart:      line,
			LineEnd:        line,
			Title:          finding.Title,
			Summary:        finding.Summary,
			Recommendation: finding.Recommendation,
			Confidence:     reviewFindingConfidence(finding.Strength),
			Severity:       finding.Strength,
			Payload: map[string]any{
				"id":       finding.ID,
				"scopes":   finding.Scopes,
				"benefit":  finding.Benefit,
				"sources":  finding.SourceIDs,
				"evidence": finding.Evidence,
			},
		}
		findings = append(findings, record)
	}
	return findings
}

func primaryFindingLocation(finding codereview.Finding, changedFiles []string) (string, int) {
	for _, evidence := range finding.Evidence {
		for _, text := range []string{evidence.Label, evidence.Value} {
			file, line := parseReviewLocation(text)
			if file != "" {
				return file, line
			}
		}
	}
	if len(changedFiles) > 0 {
		return strings.TrimSpace(changedFiles[0]), 0
	}
	return "", 0
}

func parseReviewLocation(text string) (string, int) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", 0
	}
	fields := strings.Fields(text)
	if len(fields) > 0 {
		text = fields[0]
	}
	file := text
	line := 0
	if before, after, ok := strings.Cut(text, ":"); ok {
		if parsed, err := strconv.Atoi(strings.Trim(strings.TrimSpace(after), ":,.")); err == nil {
			file = before
			line = parsed
		}
	}
	if strings.Contains(file, "/") || strings.Contains(file, ".") {
		return strings.TrimSpace(file), line
	}
	return "", 0
}

func reviewFindingFingerprint(repoFullName string, finding codereview.Finding, file string, line int, category string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		repoFullName,
		finding.ID,
		category,
		file,
		strconv.Itoa(line),
		finding.Title,
		finding.Summary,
	}, "\x00")))
	return fmt.Sprintf("%x", sum[:16])
}

func reviewFindingConfidence(strength string) int {
	switch strings.ToLower(strings.TrimSpace(strength)) {
	case "high", "strong":
		return 8
	case "low", "weak":
		return 4
	case "medium", "moderate":
		return 6
	default:
		return 6
	}
}

func reviewHistoryLanguageForFile(file string) string {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx", ".mjs", ".cjs":
		return "javascript"
	case ".rs":
		return "rust"
	case ".sql":
		return "sql"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cc", ".cpp", ".cxx", ".hh", ".hpp", ".hxx":
		return "cpp"
	default:
		return ""
	}
}

func currentHeadCommit(ctx context.Context, repoRoot string) string {
	if strings.TrimSpace(repoRoot) == "" {
		return ""
	}
	out, err := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func githubRemoteTarget(remoteURL string) (host, owner, repo string, ok bool) {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return "", "", "", false
	}
	if idx := strings.Index(remoteURL, "github.com/"); idx >= 0 {
		host = "github.com"
		path := strings.TrimSuffix(remoteURL[idx+len("github.com/"):], ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
	}
	if strings.Contains(remoteURL, "://") {
		parsed, err := url.Parse(remoteURL)
		if err != nil {
			return "", "", "", false
		}
		host = parsed.Hostname()
		if port := parsed.Port(); port != "" {
			host = parsed.Host
		}
		path := strings.TrimPrefix(parsed.Path, "/")
		path = strings.TrimSuffix(path, ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
		return "", "", "", false
	}
	if at := strings.LastIndex(remoteURL, "@"); at >= 0 {
		remoteURL = remoteURL[at+1:]
	}
	if colon := strings.Index(remoteURL, ":"); colon >= 0 {
		host = remoteURL[:colon]
		path := strings.TrimSuffix(remoteURL[colon+1:], ".git")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return host, parts[0], parts[1], true
		}
	}
	return "", "", "", false
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
	next := []string{`git add <files>`, `gx commit -m "describe this revision"`, "gx status"}
	if stack.Stack != nil && stack.Stack.BookmarkName != "" && gitCheckoutRef == stack.Stack.BookmarkName && len(current.Files) > 0 {
		next = []string{`git add <files>`, `gx commit -m "describe this revision"`, "gx status"}
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

func printCurrentStatusHuman(out io.Writer, status currentStatus) {
	fmt.Fprintln(out, commandLine("gx status", true))
	fmt.Fprintln(out)
	if strings.TrimSpace(status.PublishUploads.LastError) != "" {
		fmt.Fprintln(out, danger("ERROR: "+strings.TrimSpace(status.PublishUploads.LastError)))
		fmt.Fprintln(out, mint("Run `gx report` to report this issue."))
		fmt.Fprintln(out)
	}
	if len(status.Files) > 0 {
		fmt.Fprintf(out, "%d files are currently waiting to be assigned.\n", len(status.Files))
		fmt.Fprintln(out)
		for _, file := range status.Files {
			fmt.Fprintln(out, danger(file))
		}
	} else {
		fmt.Fprintln(out, "No files are currently waiting to be assigned.")
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
		fmt.Fprintf(out, "   %s %s\n", muted("GX Cloud uploads:"), danger(label))
		fmt.Fprintf(out, "   %s %s\n", muted("Upload error:"), status.LastError)
		return
	}
	fmt.Fprintf(out, "   %s %s\n", muted("GX Cloud uploads:"), valueText(label))
}

func currentStatusBaseStack(status currentStatus) string {
	if status.Repo.AuthoringBase != nil && strings.TrimSpace(*status.Repo.AuthoringBase) != "" {
		base := strings.TrimSpace(*status.Repo.AuthoringBase)
		if strings.HasPrefix(base, "gx/") {
			name := strings.TrimPrefix(strings.TrimPrefix(base, "gx/draft/"), "gx/")
			if name != "" {
				return "feature/" + name
			}
		} else {
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

func printStacks(ctx context.Context, engine *authoring.Engine, in io.Reader, out io.Writer, agentOut bool, selector string, opts stackDisplayOptions, interactive bool, jsonOut bool) error {
	stack, prunedEmpty, err := statusWithMissingBaseCheck(ctx, engine, out, jsonOut)
	if err != nil {
		return err
	}
	if err := enrichStatusStack(ctx, engine, &stack); err != nil {
		return err
	}
	display := stackDisplaySummaryForStacksDisplay(stack, opts)
	stack = display.Stack
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
	if interactive && useStatusInteractive(in, out) {
		action, err := runStacksInteractive(in, out, stack, unrecorded, display.HiddenEmpty)
		if err != nil {
			return err
		}
		return runStacksAction(ctx, engine, out, action)
	}
	printDeletedEmptyStacksNotice(out, prunedEmpty)
	printDefaultStatusView(out, stack, unrecorded, display.HiddenEmpty, opts.ShowAll)
	return nil
}

func enrichStatusStack(ctx context.Context, engine *authoring.Engine, stack *authoring.StackSummary) error {
	root := strings.TrimSpace(stack.Repo.RootPath)
	if root == "" {
		return nil
	}
	gitWorking, err := engine.GitWorkingStatus(ctx, root)
	if err != nil {
		return err
	}
	stack.GitWorking = gitWorking
	stack.Next = statusNextHints(gitWorking, *stack)
	return nil
}

func statusNextHints(gitWorking vcs.GitWorkingStatus, stack authoring.StackSummary) []string {
	switch {
	case gitWorking.StagedCount() > 0:
		return []string{`gx commit -m "describe this revision"`}
	case gitWorking.UnstagedCount()+gitWorking.UntrackedCount() > 0:
		return []string{"git add"}
	case unpublishedCount(stack) > 0:
		return []string{vcs.HintPush}
	default:
		if pr := currentStackPRURL(stack); pr != "" {
			return []string{pr}
		}
		return nil
	}
}

func currentStackPRURL(stack authoring.StackSummary) string {
	if stack.Stack == nil || stack.Stack.GitHubPRURL == nil {
		return ""
	}
	return strings.TrimSpace(*stack.Stack.GitHubPRURL)
}

func printDefaultStatusView(out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, hiddenEmpty int, showAll bool) {
	fmt.Fprint(out, renderDefaultStatusSummary(stack, unrecorded, hiddenEmpty, showAll))
}

func renderDefaultStatusSummary(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, hiddenEmpty int, showAll bool) string {
	stack = stackSummaryWithDisplayFallback(stack)
	var out strings.Builder
	fmt.Fprintln(&out, commandLine("gx status", true))
	fmt.Fprintln(&out)
	if hiddenEmpty > 0 {
		fmt.Fprintln(&out, muted(emptyStacksNotice(hiddenEmpty)))
		fmt.Fprintln(&out)
	}
	gitLines := gitWorkingSectionLines(stack.GitWorking, showAll)
	for _, line := range gitLines {
		fmt.Fprintln(&out, line)
	}
	if len(gitLines) > 0 {
		fmt.Fprintln(&out)
	}
	stacks := orderedStacks(stack)
	currentIdx := currentStackIndex(stack)
	if len(stacks) == 0 {
		printCurrentRevisions(&out, stack, unrecorded, 0)
	} else if currentIdx >= 0 && currentIdx < len(stacks) {
		entry := stacks[currentIdx]
		fmt.Fprintln(&out, stacksHeaderLine(stack, 1, false))
		fmt.Fprintln(&out)
		fmt.Fprintln(&out, statusBookmarkLine("●", entry, stackStatusMeta(stack, entry), true))
		fmt.Fprint(&out, renderStackRevisionLines(stack, entry, unrecorded, latestRevisionDisplayIndex(stack.Revisions, unrecorded)))
		other := len(stacks) - 1
		if other > 0 {
			fmt.Fprintln(&out)
			fmt.Fprintf(&out, "%s\n", muted(fmt.Sprintf("%d other stacks — run gx status list", other)))
		}
	} else {
		fmt.Fprint(&out, renderStacksSummaryWithHidden(stack, unrecorded, 0, false, latestRevisionDisplayIndex(stack.Revisions, unrecorded), hiddenEmpty))
		return strings.TrimRight(out.String(), "\n") + "\n"
	}
	if len(stack.Next) > 0 {
		fmt.Fprintln(&out)
		fmt.Fprintln(&out, section("Next"))
		for _, hint := range stack.Next {
			fmt.Fprintf(&out, "  %s\n", command(hint))
		}
	}
	return strings.TrimRight(out.String(), "\n") + "\n"
}

func gitWorkingSectionLines(status vcs.GitWorkingStatus, showAll bool) []string {
	var lines []string
	lines = append(lines, gitWorkingFileLines("Staged", status.Staged, showAll)...)
	lines = append(lines, gitWorkingFileLines("Unstaged", status.Unstaged, showAll)...)
	lines = append(lines, gitWorkingFileLines("Untracked", status.Untracked, showAll)...)
	return lines
}

func gitWorkingFileLines(title string, files []string, showAll bool) []string {
	if len(files) == 0 {
		return nil
	}
	lines := []string{section(title)}
	display := files
	if !showAll && len(files) > 3 {
		display = files[:3]
	}
	for _, file := range display {
		lines = append(lines, "  "+danger(file))
	}
	if !showAll && len(files) > 3 {
		lines = append(lines, muted(fmt.Sprintf("  ... %d more (gx status --all)", len(files)-3)))
	}
	return lines
}

func statusAfterPruningEmptyStacks(ctx context.Context, engine *authoring.Engine) (authoring.StackSummary, int, error) {
	pruned, err := engine.PruneEmptyStacks(ctx)
	if err != nil {
		return authoring.StackSummary{}, 0, err
	}
	if repo, repoErr := engine.ResolveJJRepo(ctx); repoErr == nil {
		_, _ = engine.PruneTerminalGitHubPullRequestStacks(ctx, repo)
	}
	stack, err := engine.Status(ctx)
	if err != nil {
		return authoring.StackSummary{}, len(pruned.Deleted), err
	}
	if strings.TrimSpace(stack.Repo.RootPath) == "" {
		if repo, repoErr := engine.ResolveJJRepo(ctx); repoErr == nil {
			stack.Repo = repo
		}
	}
	return stack, len(pruned.Deleted), nil
}

func statusWithMissingBaseCheck(ctx context.Context, engine *authoring.Engine, out io.Writer, jsonOut bool) (authoring.StackSummary, int, error) {
	stack, prunedEmpty, err := statusAfterPruningEmptyStacks(ctx, engine)
	if err != nil {
		return stack, prunedEmpty, err
	}
	missing, err := engine.DetectMissingStackBaseRefs(ctx)
	if err != nil {
		return stack, prunedEmpty, err
	}
	if len(missing.Issues) == 0 {
		return stack, prunedEmpty, nil
	}
	attachMissingStackBaseSummary(&stack, missing)
	blockErr := &vcs.ErrMissingStackBaseRefs{Status: missing}
	if jsonOut {
		return stack, prunedEmpty, blockErr
	}
	printMissingStackBaseRefNotice(out, missing)
	return stack, prunedEmpty, blockErr
}

func attachMissingStackBaseSummary(stack *authoring.StackSummary, status vcs.MissingStackBaseRefStatus) {
	if stack == nil || len(status.Issues) == 0 {
		return
	}
	stack.MissingBaseRefs = status.Issues
	stack.NeedsRebaseOntoDefault = status.NeedsRebaseOntoDefault
	stack.RepairCommand = status.RepairCommand
}

func printMissingStackBaseRefNotice(out io.Writer, status vcs.MissingStackBaseRefStatus) {
	blockErr := &vcs.ErrMissingStackBaseRefs{Status: status}
	fmt.Fprintln(out, labelWarningValue("Missing parent base", blockErr.Error()))
	for _, issue := range status.Issues {
		fmt.Fprintf(out, "  %s base=%s default=%s\n",
			valueText(firstNonEmptyString(issue.BookmarkName, issue.Name)),
			muted(issue.MissingBaseRef),
			muted(issue.DefaultBaseRef),
		)
	}
	fmt.Fprintln(out)
}

func printDeletedEmptyStacksNotice(out io.Writer, count int) {
	if count == 0 {
		return
	}
	fmt.Fprintln(out, muted(deletedEmptyStacksNotice(count)))
	fmt.Fprintln(out)
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
	default:
		return fmt.Errorf("unknown stacks action %q", action.Kind)
	}
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
		action, err := runStacksInteractive(in, out, stack, unrecorded, 0)
		if err != nil {
			return err
		}
		return runStacksAction(ctx, engine, out, action)
	}
	printStatusSummary(out, stack, unrecorded)
	return nil
}

func printStatusSummary(out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo) {
	if len(stack.Revisions) == 0 && len(stack.Stacks) == 0 && stack.Stack == nil && !hasUnrecordedFiles(unrecorded) {
		fmt.Fprintln(out, muted("No GX revisions recorded yet."))
		return
	}
	fmt.Fprint(out, renderStacksSummary(stack, unrecorded, currentStackIndex(stack), false, latestRevisionDisplayIndex(stack.Revisions, unrecorded)))
}

func printStacksSummary(out io.Writer, stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, hiddenEmpty int) {
	if len(stack.Revisions) == 0 && len(stack.Stacks) == 0 && stack.Stack == nil && !hasUnrecordedFiles(unrecorded) {
		fmt.Fprintln(out, muted("No GX revisions recorded yet."))
		if hiddenEmpty > 0 {
			fmt.Fprintln(out)
			fmt.Fprintln(out, muted(emptyStacksNotice(hiddenEmpty)))
		}
		return
	}
	fmt.Fprint(out, renderStacksSummaryWithHidden(stack, unrecorded, currentStackIndex(stack), false, latestRevisionDisplayIndex(stack.Revisions, unrecorded), hiddenEmpty))
}

func hasUnrecordedFiles(unrecorded *authoring.ChangeInfo) bool {
	return unrecorded != nil && len(unrecorded.Files) > 0
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

func printUnstagedFiles(out io.Writer, unrecorded *authoring.ChangeInfo) {
	for _, line := range unstagedFileLines(unrecorded) {
		fmt.Fprintln(out, line)
	}
	if len(unstagedFileLines(unrecorded)) > 0 {
		fmt.Fprintln(out)
	}
}

func unstagedFileLines(unrecorded *authoring.ChangeInfo) []string {
	if unrecorded == nil || len(unrecorded.Files) == 0 {
		return nil
	}
	lines := []string{section("Unstaged")}
	for _, file := range unrecorded.Files {
		lines = append(lines, "  "+danger(file))
	}
	return lines
}

func renderStacksSummary(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int) string {
	return renderStacksSummaryWithHidden(stack, unrecorded, stackCursor, stackMode, revCursor, 0)
}

func renderStacksSummaryWithHidden(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int, hiddenEmpty int) string {
	return renderStacksSummaryWithHiddenOptions(stack, unrecorded, stackCursor, stackMode, revCursor, hiddenEmpty, false)
}

func renderInteractiveStacksSummaryWithHidden(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int, hiddenEmpty int) string {
	return renderStacksSummaryWithHiddenOptions(stack, unrecorded, stackCursor, stackMode, revCursor, hiddenEmpty, true)
}

func renderStacksSummaryWithHiddenOptions(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int, hiddenEmpty int, showLegend bool) string {
	stack = stackSummaryWithDisplayFallback(stack)
	stacks := orderedStacks(stack)
	if len(stacks) == 0 {
		var out strings.Builder
		fmt.Fprintln(&out, commandLine("gx status", true))
		fmt.Fprintln(&out)
		if hiddenEmpty > 0 {
			fmt.Fprintln(&out, muted(emptyStacksNotice(hiddenEmpty)))
			fmt.Fprintln(&out)
		}
		printUnstagedFiles(&out, unrecorded)
		printCurrentRevisions(&out, stack, unrecorded, 0)
		if showLegend {
			fmt.Fprintln(&out)
			fmt.Fprintln(&out, stacksLegend())
		}
		return out.String()
	}
	if stackCursor < 0 {
		stackCursor = 0
	}
	if stackCursor >= len(stacks) {
		stackCursor = len(stacks) - 1
	}
	lines := []string{
		commandLine("gx status", true),
		"",
	}
	if hiddenEmpty > 0 {
		lines = append(lines, muted(emptyStacksNotice(hiddenEmpty)), "")
	}
	if unstaged := unstagedFileLines(unrecorded); len(unstaged) > 0 {
		lines = append(lines, unstaged...)
		lines = append(lines, "")
	}
	lines = append(lines, stacksHeaderLine(stack, len(stacks), stackMode), "")
	previousBucket := -1
	for index, entry := range stacks {
		bucket := stackDisplayBucket(entry)
		if bucket == 1 && previousBucket != 1 {
			if index > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, section("Remote"))
			lines = append(lines, "")
		} else if index > 0 {
			lines = append(lines, muted(strings.Repeat("─", 52)))
		}
		previousBucket = bucket
		selected := index == stackCursor
		if selected {
			marker := "●"
			if stackMode {
				marker = "› ●"
			}
			lines = append(lines, statusBookmarkLine(marker, entry, stackStatusMeta(stack, entry), true))
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
	if showLegend {
		lines = append(lines, "", stacksLegend())
	}
	return strings.Join(lines, "\n") + "\n"
}

func renderStacksSummaryViewport(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int, hiddenEmpty int, height int) string {
	if height <= 0 {
		return renderInteractiveStacksSummaryWithHidden(stack, unrecorded, stackCursor, stackMode, revCursor, hiddenEmpty)
	}
	stack = stackSummaryWithDisplayFallback(stack)
	stacks := orderedStacks(stack)
	top := []string{commandLine("gx status", true), ""}
	if hiddenEmpty > 0 {
		top = append(top, muted(emptyStacksNotice(hiddenEmpty)), "")
	}
	if unstaged := unstagedFileLines(unrecorded); len(unstaged) > 0 {
		top = append(top, unstaged...)
		top = append(top, "")
	}
	if len(stacks) > 0 {
		top = append(top, stacksHeaderLine(stack, len(stacks), stackMode), "")
	}
	body, selectedLine := stackSummaryBodyLines(stack, unrecorded, stackCursor, stackMode, revCursor)
	return renderScrollableView(top, body, stacksLegend(), selectedLine, height)
}

func stackSummaryBodyLines(stack authoring.StackSummary, unrecorded *authoring.ChangeInfo, stackCursor int, stackMode bool, revCursor int) ([]string, int) {
	stacks := orderedStacks(stack)
	if len(stacks) == 0 {
		var out strings.Builder
		printCurrentRevisions(&out, stack, unrecorded, 0)
		return strings.Split(strings.TrimRight(out.String(), "\n"), "\n"), 0
	}
	if stackCursor < 0 {
		stackCursor = 0
	}
	if stackCursor >= len(stacks) {
		stackCursor = len(stacks) - 1
	}
	lines := []string{}
	selectedLine := 0
	previousBucket := -1
	for index, entry := range stacks {
		bucket := stackDisplayBucket(entry)
		if bucket == 1 && previousBucket != 1 {
			if index > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, section("Remote"))
			lines = append(lines, "")
		} else if index > 0 {
			lines = append(lines, muted(strings.Repeat("─", 52)))
		}
		previousBucket = bucket
		selected := index == stackCursor
		if selected {
			marker := "●"
			if stackMode {
				marker = "› ●"
			}
			selectedLine = len(lines)
			lines = append(lines, statusBookmarkLine(marker, entry, stackStatusMeta(stack, entry), true))
			cursor := -1
			if !stackMode {
				cursor = revCursor
			}
			revisionLines := renderStackRevisionLines(stack, entry, unrecorded, cursor)
			if !stackMode && cursor >= 0 && cursor < len(revisionLines) {
				selectedLine = len(lines) + cursor
			}
			lines = append(lines, revisionLines...)
			continue
		}
		lines = append(lines, statusBookmarkLine("○", entry, stackStatusMeta(stack, entry), false))
		lines = append(lines, renderStackRevisionLines(stack, entry, unrecorded, -1)...)
	}
	return lines, selectedLine
}

func renderScrollableView(top []string, body []string, footer string, selectedLine int, height int) string {
	if height <= 0 {
		lines := append([]string{}, top...)
		lines = append(lines, body...)
		lines = append(lines, "", footer)
		return strings.Join(lines, "\n") + "\n"
	}
	if height == 1 {
		return footer
	}
	if len(top) > height-1 {
		top = top[:height-1]
	}
	bodyHeight := height - len(top) - 1
	if bodyHeight < 0 {
		bodyHeight = 0
	}
	lines := append([]string{}, top...)
	lines = append(lines, scrollBodyLines(body, selectedLine, bodyHeight)...)
	for len(lines) < height-1 {
		lines = append(lines, "")
	}
	lines = append(lines, footer)
	return strings.Join(lines, "\n") + "\n"
}

func scrollBodyLines(body []string, selectedLine int, height int) []string {
	if height <= 0 {
		return nil
	}
	if len(body) == 0 {
		return padLines(nil, height)
	}
	if len(body) <= height {
		return padLines(body, height)
	}
	if selectedLine < 0 {
		selectedLine = 0
	}
	if selectedLine >= len(body) {
		selectedLine = len(body) - 1
	}
	start := selectedLine - height/2
	if start < 0 {
		start = 0
	}
	if maxStart := len(body) - height; start > maxStart {
		start = maxStart
	}
	end := start + height
	visible := append([]string{}, body[start:end]...)
	if start > 0 {
		visible[0] = muted("... more above")
	}
	if end < len(body) {
		visible[len(visible)-1] = muted("... more below")
	}
	return visible
}

func padLines(lines []string, height int) []string {
	out := append([]string{}, lines...)
	for len(out) < height {
		out = append(out, "")
	}
	return out
}

func emptyStacksNotice(count int) string {
	word := "branches"
	if count == 1 {
		word = "branch"
	}
	return fmt.Sprintf("%d %s without revisions. Run 'gx status list' to see a full list of features.", count, word)
}

func deletedEmptyStacksNotice(count int) string {
	word := "branches"
	if count == 1 {
		word = "branch"
	}
	return fmt.Sprintf("Deleted %d %s without revisions.", count, word)
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
	if len(stack.GitWorking.Staged) > 0 {
		fmt.Fprintf(out, "staged: %s\n", quoteAgent(strings.Join(stack.GitWorking.Staged, " ")))
	}
	if len(stack.Next) > 0 {
		fmt.Fprintf(out, "next: %s\n", quoteAgent(strings.Join(stack.Next, " ")))
	}
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
	fmt.Fprintln(out, labelToken("next", `git add <files> && gx commit -m "describe this revision"`))
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

func stacksLegend() string {
	return mint("j/k up/down · d diff · esc stacks · q quit · ● selected · ↑ cloud · ↓ local")
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
		return "remote"
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
	ShowAll   bool
	ShowEmpty bool
}

type stackDisplaySummary struct {
	Stack       authoring.StackSummary
	HiddenEmpty int
}

func stackSummaryForStacksDisplay(stack authoring.StackSummary, opts stackDisplayOptions) authoring.StackSummary {
	return stackDisplaySummaryForStacksDisplay(stack, opts).Stack
}

func stackDisplaySummaryForStacksDisplay(stack authoring.StackSummary, opts stackDisplayOptions) stackDisplaySummary {
	stack = stackSummaryWithDisplayFallback(stack)
	stacks := orderedStacks(stack)
	filtered := stacks[:0]
	hiddenEmpty := 0
	for _, entry := range stacks {
		if isMergedStack(entry) {
			continue
		}
		if !opts.ShowEmpty {
			revisionCount, _ := stackEntryRevisionCounts(stack, entry)
			if revisionCount == 0 {
				hiddenEmpty++
				continue
			}
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
	return stackDisplaySummary{Stack: stack, HiddenEmpty: hiddenEmpty}
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

func newPushCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var allowBackwards bool
	var pushToGitHub bool
	cmd := &cobra.Command{
		Use:   "push [stack]",
		Short: "Push stacks and metadata to remote",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var runErr error
			defer func() {
				if runErr != nil {
					autoReportFailure(ctx, engine, runErr, "gx push")
				}
			}()
			out := cmd.OutOrStdout()
			invocation := "gx push"
			publishAllRequested := len(args) == 0
			if len(args) > 0 {
				invocation += " " + args[0]
			}
			fmt.Fprintln(out, commandLine(invocation, false))
			fmt.Fprintln(out)
			engine.SetProgressWriter(out)

			if !pushToGitHub {
				runErr = fmt.Errorf("--github=false is no longer supported; gx push always pushes stack refs")
				return runErr
			}
			mode := authoring.PublishModeReviewAndGit
			pushOpts := authoring.PushOptions{Mode: mode}

			client := cloud.NewClient()
			if client == nil {
				publishHook := func(result authoring.PushResult) error {
					printPushResult(out, result)
					branch := firstNonEmptyString(result.GXStackRef, pointerString(result.Repo.BranchName))
					fmt.Fprintln(out, labelWarningValue("Warning", fmt.Sprintf("Could not upload GX review context for branch %s: GX Cloud is not configured; set GX_CLOUD_URL or rebuild with cloud endpoints. Local remote metadata was still recorded.", branch)))
					return nil
				}

				if publishAllRequested {
					results, err := engine.PublishAll(ctx, nil, pushOpts, publishHook)
					if err != nil {
						runErr = err
						return runErr
					}
					fmt.Fprintln(out, labelValue("Pushed", fmt.Sprintf("%d stacks", len(results))))
					emitPushRunTelemetry(ctx, results, false, false, publishAllRequested, len(args) > 0)
					return nil
				}

				var (
					push authoring.PushResult
					err  error
				)
				if len(args) > 0 {
					push, err = engine.PublishNamed(ctx, args[0], nil, pushOpts, publishHook)
				} else {
					push, err = engine.Publish(ctx, nil, pushOpts, publishHook)
				}
				if err != nil {
					runErr = err
					return runErr
				}
				emitPushRunTelemetry(ctx, []authoring.PushResult{push}, false, false, publishAllRequested, len(args) > 0)
				return nil
			}

			if publishAllRequested {
				prepared, err := engine.PrepareAllPublishes(ctx, nil, pushOpts)
				if err != nil {
					runErr = err
					return runErr
				}
				published := 0
				var firstErr error
				for _, push := range prepared {
					review, err := publication.EnqueuePush(ctx, push)
					if err != nil {
						if firstErr == nil {
							firstErr = err
						}
						fmt.Fprintln(out, labelWarningValue("GX Cloud", err.Error()))
						continue
					}
					if err := engine.RecordPublish(ctx, push); err != nil {
						if firstErr == nil {
							firstErr = err
						}
						fmt.Fprintln(out, labelWarningValue("Remote metadata", err.Error()))
						continue
					}
					printPushResult(out, push)
					printPushReview(out, review)
					published++
				}
				if published > 0 {
					startPublishUploadWorker(out)
				}
				fmt.Fprintln(out, labelValue("Pushed", fmt.Sprintf("%d stacks", published)))
				emitPushRunTelemetry(ctx, prepared, true, published > 0, publishAllRequested, len(args) > 0)
				if firstErr != nil && published == 0 {
					runErr = firstErr
					return runErr
				}
				return nil
			}

			var push authoring.PushResult
			var err error
			if len(args) > 0 {
				push, err = engine.PrepareNamedPublish(ctx, args[0], nil, pushOpts)
			} else {
				push, err = engine.PreparePublish(ctx, nil, pushOpts)
			}
			if err != nil {
				runErr = err
				return runErr
			}
			review, err := publication.EnqueuePush(ctx, push)
			if err != nil {
				runErr = err
				return runErr
			}
			if err := engine.RecordPublish(ctx, push); err != nil {
				runErr = err
				return runErr
			}
			printPushResult(out, push)
			printPushReview(out, review)
			startPublishUploadWorker(out)
			emitPushRunTelemetry(ctx, []authoring.PushResult{push}, true, true, publishAllRequested, len(args) > 0)
			return nil
		},
	}
	cmd.Flags().BoolVar(&allowBackwards, "allow-backwards", false, "accepted for compatibility; gx handles required JJ bookmark moves automatically")
	pushToGitHub = true
	cmd.Flags().BoolVar(&pushToGitHub, "github", true, "accepted for compatibility; gx push always pushes stack refs")
	if flag := cmd.Flags().Lookup("allow-backwards"); flag != nil {
		flag.Hidden = true
	}
	if flag := cmd.Flags().Lookup("github"); flag != nil {
		flag.Hidden = true
	}
	return cmd
}

func emitPushRunTelemetry(ctx context.Context, pushes []authoring.PushResult, cloudConfigured bool, uploadQueued bool, publishAllRequested bool, selectedStack bool) {
	revisionCount := 0
	githubPRCount := 0
	gitPushCount := 0
	warningCount := 0
	for _, push := range pushes {
		publishedCount := len(push.Published)
		if publishedCount == 0 && push.CurrentChange != nil {
			publishedCount = 1
		}
		revisionCount += publishedCount
		warningCount += len(push.Warnings)
		if pushWasPushedToGitHub(push.GitPushStatus) {
			gitPushCount++
		}
		githubPRCount += pushGitHubPRCount(push)
	}
	telemetry.EmitProductEvent(ctx, telemetry.EventCLIPushRun, map[string]any{
		"status":             "success",
		"cloud_configured":   cloudConfigured,
		"upload_queued":      uploadQueued,
		"push_all_requested": publishAllRequested,
		"selected_stack":     selectedStack,
		"stack_count":        len(pushes),
		"revision_count":     revisionCount,
		"github_pr_count":    githubPRCount,
		"git_push_count":     gitPushCount,
		"warning_count":      warningCount,
	})
}

func pushWasPushedToGitHub(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pushed", "force pushed", "already up to date":
		return true
	default:
		return false
	}
}

func pushGitHubPRCount(push authoring.PushResult) int {
	seen := map[string]bool{}
	if push.GitHubPullRequestURL != nil {
		if url := strings.TrimSpace(*push.GitHubPullRequestURL); url != "" {
			seen[url] = true
		}
	}
	for _, published := range push.Published {
		if published.GitHubPullRequestURL == nil {
			continue
		}
		if url := strings.TrimSpace(*published.GitHubPullRequestURL); url != "" {
			seen[url] = true
		}
	}
	return len(seen)
}

func printPushResult(out io.Writer, result authoring.PushResult) {
	if strings.TrimSpace(result.Output) != "" {
		fmt.Fprintln(out, highlightPushRevisions(strings.TrimSpace(result.Output)))
	}
	if strings.TrimSpace(result.GXStackRef) != "" || result.Repo.BranchName != nil {
		branch := firstNonEmptyString(result.GXStackRef, pointerString(result.Repo.BranchName))
		fmt.Fprintln(out, labelValue("Branch", branch))
		if status := pushGitPushStatusText(result.GitPushStatus); status != "" {
			fmt.Fprintln(out, labelStatus("GitHub", status))
		}
		if prStatus := pushGitHubPRStatusText(result.GitHubPRStatus); prStatus != "" {
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

func pushGitPushStatusText(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pushed":
		return "ok: pushed"
	case "force pushed":
		return "ok: force pushed with lease"
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

func pushGitHubPRStatusText(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "created", "existing", "stored":
		return "ok: " + strings.TrimSpace(status)
	case "warning":
		return "warn: not created"
	default:
		return strings.TrimSpace(status)
	}
}

func highlightPushRevisions(output string) string {
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

func printPushReview(out io.Writer, result publication.Result) {
	if result.Queued {
		fmt.Fprintln(out, labelValue("GX Cloud", "queued for background upload"))
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
		fmt.Fprintln(out, labelValue("Remote", "gx session context"))
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
	cmd := &cobra.Command{
		Use:   "sync [remote]",
		Short: "Sync remote with local",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
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
				telemetry.EmitProductEvent(ctx, telemetry.EventCLISyncRun, map[string]any{
					"remote":        result.RemoteName,
					"listed":        summary.Listed,
					"updated":       summary.Updated,
					"caught_up":     summary.CaughtUp,
					"merged_pruned": summary.Removed,
				})
				fmt.Fprintln(out, labelValue("Remote bookmarks", fmt.Sprintf("%d listed", summary.Listed)))
				if summary.Updated > 0 {
					fmt.Fprintln(out, labelValue("Remote revisions", fmt.Sprintf("%d updated", summary.Updated)))
				}
				if summary.CaughtUp > 0 {
					fmt.Fprintln(out, labelValue("Remote catch-up", fmt.Sprintf("%d bookmark(s) fetched from remote", summary.CaughtUp)))
				}
				if summary.Removed > 0 {
					fmt.Fprintln(out, labelValue("Remote merged", fmt.Sprintf("%d bookmark(s) removed from remote", summary.Removed)))
				}
			}
			if summary, err := engine.PruneTerminalGitHubPullRequestStacks(ctx, result.Repo); err != nil {
				fmt.Fprintln(out, labelWarningValue("GitHub PR sync", err.Error()))
			} else if summary.Removed > 0 {
				fmt.Fprintln(out, labelValue("GitHub PRs closed", fmt.Sprintf("%d stack(s) removed from local db", summary.Removed)))
			}
			if status, err := publication.QueuedUploadStatus(); err == nil && (status.Pending > 0 || status.Failed > 0) {
				if err := drainPublishUploadOutbox(ctx, out, false, 20); err != nil {
					fmt.Fprintln(out, labelWarningValue("GX Cloud uploads", err.Error()))
				}
			}
			if !cloud.CloudConfigured() {
				telemetry.EmitProductEvent(ctx, telemetry.EventCLISyncRun, map[string]any{
					"remote": result.RemoteName,
				})
			}
			return nil
		},
	}
	return cmd
}

func newPublishUploadCommand(ctx context.Context) *cobra.Command {
	var quiet bool
	var limit int
	cmd := &cobra.Command{
		Use:    "__gx-upload-outbox",
		Short:  "Upload queued GX Cloud context",
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
		fmt.Fprintln(out, labelValue("GX Cloud uploads", fmt.Sprintf("%d uploaded, %d failed, %d pending", result.Uploaded, result.Failed, result.Pending)))
	}
	return nil
}

func startPublishUploadWorker(out io.Writer) {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(out, labelWarningValue("GX Cloud upload", "queued; could not find gx executable, run `gx sync` to upload"))
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
		fmt.Fprintln(out, labelWarningValue("GX Cloud upload", "queued; run `gx sync` to upload"))
		return
	}
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	fmt.Fprintln(out, labelValue("GX Cloud upload", "background upload started"))
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
	switch filepath.Base(os.Args[0]) {
	case "gxr":
		args = append([]string{"review"}, args...)
	case "gxs":
		args = append([]string{"status"}, args...)
	}
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

// ExitCode returns a non-zero process exit code when err carries one.
func ExitCode(err error) int {
	var coded *vcs.CodedError
	if errors.As(err, &coded) && coded.Code != 0 {
		return coded.Code
	}
	return 0
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
