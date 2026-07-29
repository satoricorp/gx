package cli

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
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
	cmd.Flags().StringVar(&failOn, "fail-on", string(codereview.FailOnNone), fmt.Sprintf("exit %d when findings at or above this level survive: %s (exit %d when there was nothing to review, exit %d when the review ran degraded)", reviewFindingsExitCode, strings.Join(codereview.FailOnLevels(), ", "), reviewNothingToReviewExitCode, reviewDegradedExitCode))
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
