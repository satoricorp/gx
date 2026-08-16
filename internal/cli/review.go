package cli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/spf13/cobra"
)

// Exit codes for `gx review` as an automated gate. Both are distinct from the
// generic failure exit so a CI step can tell a policy failure from a crash.
const (
	// reviewFindingsExitCode means findings at or above --fail-on survived.
	reviewFindingsExitCode = 3
	// reviewNothingToReviewExitCode means the review inspected no code at all.
	// A gate that passes without looking is the false pass --fail-on exists to
	// prevent, so an explicit gate treats "never looked" as a failure too.
	reviewNothingToReviewExitCode = 4
	// reviewDegradedExitCode means code was read but the review that read it was
	// incomplete — no model ran, or only part of the subject reached one. Same
	// reasoning as above: "no findings at or above X" is a claim about what was
	// inspected, and a gate must not make it on a review that did not run.
	reviewDegradedExitCode = 5
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
	var noComment bool
	var fast bool
	var maxFindings int
	var clientOverride string
	cmd := &cobra.Command{
		Use:     "review [prompt]",
		Aliases: []string{"gxr"},
		Short:   "Review changes based on codebase & session context, along with independent resources",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Everything below runs under a context that forbids writing gx
			// state, so review's own telemetry reports without minting a
			// machine ID into a $GX_HOME that may not exist.
			ctx := telemetry.WithoutStateWrites(ctx)
			startedAt := time.Now()
			// Resolved once so telemetry and the cloud history row cannot
			// disagree about which surface invoked this run.
			client := telemetry.ClientSurface(clientOverride)
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
			// Store-free on purpose: review must leave no gx state behind in a
			// repo (or on a machine) that has never run `gx init`.
			repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
			if err != nil {
				emitReviewRunTelemetry(ctx, codereview.Report{}, err, client, reviewScope, scopeExplicit, focus, prompt, deep, wholeRepo, verbose, time.Since(startedAt))
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
					Fast:           fast,
					MaxFindings:    maxFindings,
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
			emitReviewRunTelemetry(ctx, report, err, client, reviewScope, scopeExplicit, focus, prompt, deep, wholeRepo, verbose, time.Since(startedAt))
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
			// The saved-report notice goes to stderr so --json stdout stays
			// pure report.
			if reportPath, saveErr := saveReviewReport(repo.RootPath, report, startedAt); saveErr != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), labelWarningValue("Report", saveErr.Error()))
			} else if reportPath != "" {
				fmt.Fprintln(cmd.ErrOrStderr(), labelValue("Report", muted("saved to "+reportPath)))
			}
			// Nothing was reviewed: publishing would overwrite a real review
			// comment with a no-op, and there is no review to record.
			// --no-comment suppresses only the outward-facing PR comment;
			// history still records so per-surface review counts stay honest.
			// --no-publish remains full suppression for runs that must leave
			// no trace in gx Cloud.
			if !noPublish && report.Reviewed {
				if !noComment {
					postReviewSummaryComment(ctx, repo, report, cmd.ErrOrStderr())
				}
				recordReviewHistory(ctx, repo, report, client, prompt, scopeExplicit, deep, wholeRepo, cmd.ErrOrStderr())
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
	// The gate is on by default. "blocking" is lane-aware: it fails only on
	// findings that earned the blocking lane — a Blocking deterministic check,
	// or a Strong finding both reviewers raised and the judge confirmed. A
	// Strong finding one model raised is demoted to advisory and does not fail
	// it. --fail-on none restores the advisory-only exit 0.
	cmd.Flags().StringVar(&failOn, "fail-on", string(codereview.FailOnBlocking), fmt.Sprintf("exit %d when findings at or above this level survive: %s; the default %q fails only on findings in the blocking lane, and \"none\" turns the gate off (exit %d when there was nothing to review, exit %d when the review ran degraded)", reviewFindingsExitCode, strings.Join(codereview.FailOnLevels(), ", "), string(codereview.FailOnBlocking), reviewNothingToReviewExitCode, reviewDegradedExitCode))
	cmd.Flags().BoolVar(&noPublish, "no-publish", false, "skip posting the PR review comment and recording review history")
	cmd.Flags().BoolVar(&noComment, "no-comment", false, "skip posting the PR review comment but still record review history; use --no-publish to suppress both")
	cmd.Flags().StringVar(&clientOverride, "client", "", "surface invoking this review, overriding $GX_CLIENT: cli, mcp, skill, slash-gx, slash-constraints")
	cmd.Flags().IntVar(&maxFindings, "max-findings", 0, "cap how many recommendations the review reports (0 uses the default); applies to both what the model is asked for and what is reported")
	cmd.Flags().BoolVar(&fast, "fast", false, "optimize for wall clock: one reviewer instead of two, no verification pass, and findings written without code examples")
	return cmd
}

func writeReviewJSON(out io.Writer, report codereview.Report) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

// saveReviewReport persists the full JSON report to
// <repoRoot>/review/review-findings-<timestamp>.json so a review survives the
// terminal it was printed in. Returns the repo-relative path, or "" when the
// review inspected nothing (an empty report is not worth a file).
func saveReviewReport(repoRoot string, report codereview.Report, startedAt time.Time) (string, error) {
	if !report.Reviewed {
		return "", nil
	}
	dir := filepath.Join(repoRoot, "review")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	// Best-effort: without the exclude entry a saved report dirties
	// `git status` and shows up as a changed file in the NEXT review, but a
	// failure to write it should not cost the report itself.
	_ = ensureReviewDirIgnored(repoRoot)
	name := "review-findings-" + startedAt.UTC().Format("20060102-150405") + ".json"
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	if err := writeReviewJSON(file, report); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return filepath.Join("review", name), nil
}

// ensureReviewDirIgnored appends /review/ to .git/info/exclude so saved
// reports never show as untracked files. The exclude file is repo-local and
// never committed, so the ignore does not reach collaborators — a team that
// wants reports tracked can still add them explicitly with `git add -f`.
func ensureReviewDirIgnored(repoRoot string) error {
	gitDir, err := resolveGitCommonDir(repoRoot)
	if err != nil {
		return err
	}
	excludePath := filepath.Join(gitDir, "info", "exclude")
	existing, err := os.ReadFile(excludePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(line) == "/review/" {
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(excludePath), 0o755); err != nil {
		return err
	}
	content := string(existing)
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += "/review/\n"
	return os.WriteFile(excludePath, []byte(content), 0o644)
}

// resolveGitCommonDir finds the directory whose info/exclude git actually
// reads: .git itself for a normal checkout, the pointed-to gitdir for a
// linked worktree — and, when that gitdir carries a commondir file, the
// shared common directory (per-worktree gitdirs' info/exclude is ignored by
// git).
func resolveGitCommonDir(repoRoot string) (string, error) {
	gitPath := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", err
	}
	gitDir := gitPath
	if !info.IsDir() {
		data, readErr := os.ReadFile(gitPath)
		if readErr != nil {
			return "", readErr
		}
		pointer := strings.TrimSpace(string(data))
		const prefix = "gitdir:"
		if !strings.HasPrefix(pointer, prefix) {
			return "", fmt.Errorf(".git is neither a directory nor a gitdir pointer")
		}
		gitDir = strings.TrimSpace(strings.TrimPrefix(pointer, prefix))
		if !filepath.IsAbs(gitDir) {
			// Git resolves relative gitdir pointers against the directory
			// containing the .git file, not the repo root.
			gitDir = filepath.Join(filepath.Dir(gitPath), gitDir)
		}
	}
	if data, readErr := os.ReadFile(filepath.Join(gitDir, "commondir")); readErr == nil {
		common := strings.TrimSpace(string(data))
		if !filepath.IsAbs(common) {
			common = filepath.Join(gitDir, common)
		}
		gitDir = common
	}
	return filepath.Clean(gitDir), nil
}

// reviewGateError turns a review outcome into an exit code. "Nothing to
// review" is its own outcome: it never counts as findings-clean, and under a
// gate it fails, because a gate that passes on a diff it never opened is
// exactly the silent false pass --fail-on is meant to catch. The gate is on by
// default (--fail-on blocking); --fail-on none is the advisory-only mode where
// every outcome stays exit 0.
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
	// A degraded run is not a clean run with fewer findings. When no model ran,
	// `findings` is whatever the deterministic checks produced — usually nothing
	// once patch-focus filtering is applied — so the gate exited 0 and the pull
	// request merged reporting a review that never happened. The rendered report
	// says so in a banner, but an exit code is the only thing a CI step reads.
	if reason := gateDegradedReason(report); reason != "" {
		return vcs.CodedErrorf(reviewDegradedExitCode, fmt.Errorf("gx review: %s; refusing to pass a gate on an incomplete review", reason))
	}
	failures := report.GateFailures(level)
	if len(failures) == 0 {
		return nil
	}
	return vcs.CodedErrorf(reviewFindingsExitCode, fmt.Errorf("gx review: %d finding(s) at or above %q", len(failures), string(level)))
}

// gateDegradedReason states why this review cannot answer the gate's question,
// or "" when it can. Both signals are already computed and rendered; no gate
// path read either until now.
func gateDegradedReason(report codereview.Report) string {
	if len(report.DegradedReasons) > 0 {
		return strings.Join(report.DegradedReasons, "; ")
	}
	if report.Coverage.Partial() {
		if statement := strings.TrimSpace(report.Coverage.Statement()); statement != "" {
			return statement
		}
		return "only part of the change was reviewed"
	}
	return ""
}

func emitReviewRunTelemetry(ctx context.Context, report codereview.Report, runErr error, client string, reviewScope string, scopeExplicit bool, focus string, prompt string, deep bool, wholeRepo bool, verbose bool, duration time.Duration) {
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
		"client":                client,
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
	// can be attributed to the gx Cloud hop or to direct AWS calls instead of
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

// reviewRunMode labels a review run for gx Cloud history and telemetry. Both
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

func recordReviewHistory(ctx context.Context, repo vcs.RepoInfo, report codereview.Report, clientSurface string, prompt string, scopeExplicit bool, deep bool, wholeRepo bool, stderr io.Writer) {
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
		// Signed out: history is an extra gx Cloud records for authenticated
		// users, not something a read-only review depends on.
		return
	}
	req := cloud.CodeReviewHistoryRecordRequest{
		RepoRootPath: report.RepoRoot,
		RepoFullName: repoFullName,
		BranchName:   pointerString(repo.BranchName),
		HeadCommitID: currentHeadCommit(ctx, repo.RootPath),
		SourceKind:   "session_intent",
		Client:       clientSurface,
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

// reviewFindingConfidence maps a finding's Strength onto the 1-10 confidence
// the gx Cloud history row stores. The real Strength vocabulary is "Blocking",
// "Strong", "Worth exploring", "Speculative" (codereview.Finding.Strength); the
// generic high/medium/low spellings are kept as aliases so any older producer
// still lands where it used to.
func reviewFindingConfidence(strength string) int {
	switch strings.ToLower(strings.TrimSpace(strength)) {
	case "blocking":
		return 9
	case "strong", "high":
		return 8
	case "worth exploring", "worth-exploring":
		return 5
	case "speculative":
		return 3
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
	case ".cs":
		return "csharp"
	case ".dart":
		return "dart"
	default:
		return ""
	}
}
