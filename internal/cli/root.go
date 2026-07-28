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
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/github"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/inference"
	"github.com/satoricorp/gx/internal/postlist"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/satoricorp/gx/internal/version"
)

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
			telemetry.EmitInstallOnce(commandTelemetryContext(ctx, cmd))
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
		&cobra.Group{ID: groupHelp, Title: "Help:"},
	)

	root.AddCommand(
		newInternalHooksCommand(ctx),
		newVersionCommand(),
		newDoctorCommand(ctx),
		newAuthCommand(ctx),
		newSetCommand(ctx),
		newInitCommand(ctx, engine),
		newCaptureCommand(ctx),
		newPublishUploadCommand(ctx),
		newReportCommand(ctx),
		newReviewCommand(ctx),
		newIndexCommand(ctx),
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
		case "init", "auth", "set":
			cmd.GroupID = groupSetup
		case "review":
			cmd.GroupID = groupWork
		case "doctor", "version":
			cmd.GroupID = groupHelp
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
			cleanupLegacyAmbientCaptureFromInit(ctx, cmd, yes)
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

func newReviewCommand(ctx context.Context) *cobra.Command {
	var scope string
	var focus string
	var base string
	var wholeRepo bool
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
			// Everything below runs under a context that forbids writing GX
			// state, so review's own telemetry reports without minting a
			// machine ID into a $GX_HOME that may not exist.
			ctx := telemetry.WithoutStateWrites(ctx)
			startedAt := time.Now()
			failOnLevel, err := codereview.ParseFailOnLevel(failOn)
			if err != nil {
				return err
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
				emitReviewRunTelemetry(ctx, codereview.Report{}, err, reviewScope, scopeExplicit, focus, prompt, deep, wholeRepo, verbose, time.Since(startedAt))
				return err
			}
			runReview := func(progress io.Writer) (codereview.Report, error) {
				return codereview.Review(ctx, repo.RootPath, codereview.Options{
					Scope:          reviewScope,
					Deep:           deep,
					Focus:          focus,
					Base:           base,
					WholeRepo:      wholeRepo,
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
			emitReviewRunTelemetry(ctx, report, err, reviewScope, scopeExplicit, focus, prompt, deep, wholeRepo, verbose, time.Since(startedAt))
			if err != nil {
				return err
			}
			if jsonOut {
				if err := writeReviewJSON(cmd.OutOrStdout(), report); err != nil {
					return err
				}
			} else {
				fmt.Fprint(cmd.OutOrStdout(), codereview.RenderMarkdown(report))
			}
			// Nothing was reviewed: publishing would overwrite a real review
			// comment with a no-op, and there is no review to record.
			if !noPublish && report.Reviewed {
				postReviewSummaryComment(ctx, repo, report, cmd.ErrOrStderr())
				recordReviewHistory(ctx, repo, report, prompt, scopeExplicit, deep, wholeRepo, cmd.ErrOrStderr())
			}
			return reviewGateError(report, failOnLevel)
		},
	}
	cmd.Flags().StringVar(&scope, "scope", codereview.DefaultScope, "review scope when explicitly set: architecture, security, performance, onboarding, docs, dependencies, testing, maintainability")
	cmd.Flags().StringVar(&focus, "focus", "", "limit review to files under this path prefix")
	cmd.Flags().StringVar(&base, "base", "", "review the commit range <ref>...HEAD instead of the working tree, e.g. --base origin/main")
	cmd.Flags().BoolVar(&wholeRepo, "repo", false, "review the whole repository rather than just the current change; uncommitted work stays in focus, and this wins over --base")
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

func emitReviewRunTelemetry(ctx context.Context, report codereview.Report, runErr error, reviewScope string, scopeExplicit bool, focus string, prompt string, deep bool, wholeRepo bool, verbose bool, duration time.Duration) {
	status := "success"
	if runErr != nil {
		status = "error"
	}
	mode := reviewRunMode(prompt, scopeExplicit, deep, wholeRepo)
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
		"whole_repo":            wholeRepo,
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
	// Which wire the review ran over, so a fleet-wide latency or failure spike
	// can be attributed to the GX Cloud hop or to direct AWS calls instead of
	// being averaged across both. The models go with it: the panel is two
	// competing legs plus a judge, and "review got slower" is a different
	// investigation depending on which of them changed.
	if strings.TrimSpace(report.ReviewTransport) != "" {
		props["review_transport"] = report.ReviewTransport
	}
	if len(report.ReviewModels) > 0 {
		props["review_models"] = strings.Join(report.ReviewModels, ",")
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

// reviewRunMode labels a review run for GX Cloud history and telemetry. Both
// call it so the two records of the same run cannot disagree.
//
// WholeRepo comes first because it names the subject: a run recorded as
// "patch" whose summary text reads "Reviewed the repository" is a row that
// contradicts itself, and anything downstream that filters or ranks on the
// mode would treat a whole-repo review as a patch review. Depth and scope
// travel alongside as their own fields, so nothing is lost by this ordering.
func reviewRunMode(prompt string, scopeExplicit bool, deep bool, wholeRepo bool) string {
	switch {
	case wholeRepo:
		return "repo"
	case deep:
		return "deep"
	case strings.TrimSpace(prompt) != "":
		return "prompt"
	case scopeExplicit:
		return "scope"
	default:
		return "patch"
	}
}

func recordReviewHistory(ctx context.Context, repo vcs.RepoInfo, report codereview.Report, prompt string, scopeExplicit bool, deep bool, wholeRepo bool, stderr io.Writer) {
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
	req := cloud.CodeReviewHistoryRecordRequest{
		RepoRootPath: report.RepoRoot,
		RepoFullName: repoFullName,
		BranchName:   pointerString(repo.BranchName),
		HeadCommitID: currentHeadCommit(ctx, repo.RootPath),
		SourceKind:   "session_intent",
		Prompt:       prompt,
		Scope:        firstNonEmptyString(report.Scope, codereview.DefaultScope),
		Mode:         reviewRunMode(prompt, scopeExplicit, deep, wholeRepo),
		Reviewer:     report.Reviewer,
		SummaryKind:  "pr",
		SummaryText:  codereview.RenderMarkdown(report),
		Findings:     reviewHistoryFindings(repoFullName, report),
		Payload: map[string]any{
			"changed_files":  len(report.ChangedFiles),
			"deep":           deep,
			"whole_repo":     wholeRepo,
			"scope_explicit": scopeExplicit,
			// What the run actually read, so the row can be reconciled with
			// its own summary text rather than inferred from the flags.
			"review_mode": report.ReviewMode,
			"reviewed":    report.Reviewed,
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

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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

func Execute(ctx context.Context) error {
	args := os.Args[1:]
	switch filepath.Base(os.Args[0]) {
	case "gxr":
		args = append([]string{"review"}, args...)
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

func shortID(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
