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
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
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
		Short:         "GX CLI for Git-native capture, commits, and review",
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
		newInternalHooksCommand(ctx),
		newVersionCommand(),
		newDoctorCommand(ctx),
		newAuthCommand(ctx),
		newSetCommand(ctx),
		newInitCommand(ctx, engine),
		newCaptureCommand(ctx),
		newPublishUploadCommand(ctx),
		newDemoCommand(),
		newReportCommand(ctx, engine),
		newReviewCommand(ctx, engine),
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
		case "review":
			cmd.GroupID = groupWork
		case "sync":
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
	var global bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Set up gx in the current repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			if global {
				if !yes {
					fmt.Fprintln(cmd.OutOrStdout(), commandLine("gx init --global", true))
					fmt.Fprintln(cmd.OutOrStdout())
				}
				err := runGlobalInit(ctx, cmd, yes)
				telemetry.EmitProductEvent(ctx, telemetry.EventCLIInitRun, map[string]any{
					"status":      initRunStatus(err),
					"interactive": !yes,
					"global":      true,
				})
				return err
			}
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
			telemetry.EmitProductEvent(ctx, telemetry.EventCLIInitRun, map[string]any{
				"status":      initRunStatus(err),
				"interactive": !yes,
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
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Next", "git add <files> && git commit -m \"...\""))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "user name to store in GX config")
	cmd.Flags().StringVar(&email, "email", "", "user email to store in GX config")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "accept defaults and suppress successful init output")
	cmd.Flags().BoolVar(&global, "global", false, "install machine-wide git hooks (~/.gx/hooks) so GX works in every repo; skips per-repo setup")
	return cmd
}

func initRunStatus(err error) string {
	if err != nil {
		return "error"
	}
	return "success"
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

// Exit codes for `gx review` as an automated gate. Both are distinct from the
// generic failure exit so a CI step can tell a policy failure from a crash.
const (
	// reviewFindingsExitCode means findings at or above --fail-on survived.
	reviewFindingsExitCode = 3
	// reviewNothingToReviewExitCode means the review inspected no code at all.
	// A gate that passes without looking is the false pass --fail-on exists to
	// prevent, so an explicit gate treats "never looked" as a failure too.
	reviewNothingToReviewExitCode = 4
)

func newReviewCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var scope string
	var focus string
	var base string
	var deep bool
	var verbose bool
	var jsonOut bool
	var failOn string
	var noPublish bool
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
			failOnLevel, err := codereview.ParseFailOnLevel(failOn)
			if err != nil {
				runErr = err
				return runErr
			}
			reviewScope := ""
			scopeExplicit := cmd.Flags().Changed("scope")
			if scopeExplicit {
				reviewScope = scope
			}
			prompt := ""
			if len(args) > 0 {
				prompt = strings.TrimSpace(args[0])
			}
			// Store-free on purpose: review must leave no GX state behind in a
			// repo (or on a machine) that has never run `gx init`.
			repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
			if err != nil {
				runErr = err
				emitReviewRunTelemetry(ctx, codereview.Report{}, err, reviewScope, scopeExplicit, focus, prompt, deep, verbose, time.Since(startedAt))
				return runErr
			}
			runReview := func(progress io.Writer) (codereview.Report, error) {
				return codereview.Review(ctx, repo.RootPath, codereview.Options{
					Scope:          reviewScope,
					Deep:           deep,
					Focus:          focus,
					Base:           base,
					Prompt:         prompt,
					Verbose:        verbose,
					ProgressWriter: progress,
					Color:          !jsonOut,
				})
			}
			var report codereview.Report
			if jsonOut {
				// Machine-readable output: no spinner, no ANSI, nothing on
				// stdout but the report.
				report, err = runReview(nil)
			} else {
				report, err = runReviewWithLoader(cmd.InOrStdin(), cmd.ErrOrStderr(), runReview)
			}
			emitReviewRunTelemetry(ctx, report, err, reviewScope, scopeExplicit, focus, prompt, deep, verbose, time.Since(startedAt))
			if err != nil {
				runErr = err
				return runErr
			}
			if jsonOut {
				if err := writeReviewJSON(cmd.OutOrStdout(), report); err != nil {
					runErr = err
					return runErr
				}
			} else {
				fmt.Fprint(cmd.OutOrStdout(), codereview.RenderMarkdown(report))
			}
			// Nothing was reviewed: publishing would overwrite a real review
			// comment with a no-op, and there is no review to record.
			if !noPublish && report.Reviewed {
				postReviewSummaryComment(ctx, repo, report, cmd.ErrOrStderr())
				recordReviewHistory(ctx, repo, report, prompt, scopeExplicit, deep, cmd.ErrOrStderr())
			}
			runErr = reviewGateError(report, failOnLevel)
			return runErr
		},
	}
	cmd.Flags().StringVar(&scope, "scope", codereview.DefaultScope, "review scope when explicitly set: architecture, security, performance, onboarding, docs, dependencies, testing, maintainability")
	cmd.Flags().StringVar(&focus, "focus", "", "limit review to files under this path prefix")
	cmd.Flags().StringVar(&base, "base", "", "review the commit range <ref>...HEAD instead of the working tree, e.g. --base origin/main")
	cmd.Flags().BoolVar(&deep, "deep", false, "run full-spectrum review with more local and indexed context")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "include repo facts, docs, and changed files")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print the review report as JSON instead of markdown")
	cmd.Flags().StringVar(&failOn, "fail-on", string(codereview.FailOnNone), fmt.Sprintf("exit %d when findings at or above this level survive: %s (exit %d when there was nothing to review)", reviewFindingsExitCode, strings.Join(codereview.FailOnLevels(), ", "), reviewNothingToReviewExitCode))
	cmd.Flags().BoolVar(&noPublish, "no-publish", false, "skip posting the PR review comment and recording review history")
	return cmd
}

func writeReviewJSON(out io.Writer, report codereview.Report) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// reviewGateError turns a review outcome into an exit code. "Nothing to
// review" is its own outcome: it never counts as findings-clean, and under an
// explicit gate it fails, because a gate that passes on a diff it never opened
// is exactly the silent false pass --fail-on is meant to catch. Without
// --fail-on the default stays exit 0, so interactive use is unaffected.
func reviewGateError(report codereview.Report, level codereview.FailOnLevel) error {
	if !level.Enabled() {
		return nil
	}
	if !report.Reviewed {
		target := strings.TrimSpace(report.ReviewTarget)
		if target == "" {
			target = "the working tree"
		}
		return vcs.CodedErrorf(reviewNothingToReviewExitCode, fmt.Errorf("gx review: nothing was reviewed (looked at %s); refusing to pass a gate without inspecting any code", target))
	}
	failures := report.GateFailures(level)
	if len(failures) == 0 {
		return nil
	}
	return vcs.CodedErrorf(reviewFindingsExitCode, fmt.Errorf("gx review: %d finding(s) at or above %q", len(failures), string(level)))
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
	if strings.TrimSpace(report.ReviewMode) != "" {
		// What the review actually read, so a fleet-wide "nothing to review"
		// rate is visible rather than hiding inside the success count.
		props["review_mode"] = report.ReviewMode
		props["reviewed"] = report.Reviewed
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
		// No GitHub token. Review is read-only and complete without one, so
		// posting the comment is a bonus, not a failure: skip it quietly
		// rather than nag on every review from an unauthenticated checkout.
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
	if _, err := cloud.CloudAPIToken(); err != nil {
		// Signed out: history is an extra GX Cloud records for authenticated
		// users, not something a read-only review depends on.
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
}

type currentStatusRefs struct {
	GXBaseRef       string `json:"gx_base_ref"`
	GXStackRef      string `json:"gx_stack_ref,omitempty"`
	GitCheckoutRef  string `json:"git_checkout_ref,omitempty"`
	GitPublishedRef string `json:"git_published_ref,omitempty"`
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
	}, nil
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

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
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

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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

// ExitCode maps an error to a process exit status. Commands that need to fail
// with a specific status wrap their error with vcs.CodedErrorf; everything else
// falls through to the caller's generic failure code.
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
