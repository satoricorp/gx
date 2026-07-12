package publication

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
	"github.com/satoricorp/gx/internal/semantic"
)

const (
	githubPRBodyMarker        = "<!-- gx:pr-summary:v1 -->"
	githubPRAuthorNotesMarker = "<!-- gx:author-notes -->"
	githubPRAuthorNotesHeader = "## Author Notes"
	maxPRNotableChanges       = 7
	maxPRContextSnippetSize   = 1800
	maxPROverviewLength       = 500

	verdictNoReview  = "No review needed"
	verdictQuickScan = "Quick scan"
	verdictDeepReview = "Requires Deep Review"
)

var (
	hunkHeaderRE  = regexp.MustCompile(`^@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)
	prURLStripRE  = regexp.MustCompile(`https?://\S+`)
)

var prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
	return codereview.ReviewerFromEnvWithInfo()
}
var prSummaryReviewRetryWait = func(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
var collectPRSummaryContext = collectDefaultPRSummaryContext

const prSummaryReviewRetryDelay = 2 * time.Second

type prRevisionSummary struct {
	Description string
	CommitID    string
	BranchName  string
	Files       []string
	RiskLevel   string
	RiskScore   int
	RiskSignals []string
}

type prHunkSummary struct {
	ID          string
	Revision    string
	File        string
	Header      string
	Patch       string
	OldStart    int
	OldLines    int
	NewStart    int
	NewLines    int
	ChangedLine int
	Link        string
}

type prBodyStats struct {
	RevisionCount      int
	FileCount          int
	AreaCount          int
	AddedLines         int
	DeletedLines       int
	ChangedSymbolCount int
	WarningCount       int
	InfoCount          int
	MaxRiskLevel       string
	MaxRiskScore       int
	RiskSignals        []string
}

type prBodyCatalog struct {
	Revisions []prRevisionSummary
	Hunks     []prHunkSummary
	Stats     prBodyStats
	Files     []string
	Areas     []string
}

type prSummaryContext struct {
	Snippets []codereview.ContextSnippet
	Sources  []string
}

type prAttribution struct {
	Kind   string
	Label  string
	Ref    string
	URL    string
	Opaque bool
}

type prNotableChange struct {
	Title        string
	Detail       string
	Link         string
	Flagged      bool
	Score        int
	Attributions []prAttribution
}

func GitHubPullRequestBodyFromArtifact(ctx context.Context, artifact reviewbundle.Artifact) (string, error) {
	catalog := buildPRBodyCatalog(artifact)
	reach := computeLexicalReach(ctx, artifact.Bundle.Repo.RootPath, catalog, artifact)
	applyLexicalReachToStats(&catalog.Stats, reach)
	policy := codereview.LoadReviewPolicy(ctx, artifact.Bundle.Repo.RootPath)
	summaryContext, err := collectPRSummaryContext(ctx, artifact, catalog)
	if err != nil {
		return "", err
	}
	summary, aiSucceeded, reviewerInfo, err := reviewPRSummaryFindings(ctx, artifact, catalog, summaryContext, reach, policy)
	if err != nil {
		return "", err
	}
	return renderGitHubPullRequestBody(artifact, catalog, summaryContext, summary, aiSucceeded, reach, policy, reviewerInfo), nil
}

func renderGitHubPullRequestBody(artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext, summary codereview.PRSummaryReview, aiSucceeded bool, reach lexicalReach, policy codereview.ReviewPolicy, reviewerInfo codereview.ReviewerInfo) string {
	triage := triageChangeFromCatalog(catalog)
	verdict, reason := reviewVerdict(triage, catalog.Stats, summary.Findings, aiSucceeded, catalog.Files, catalog.Areas)
	includeHeuristics := verdict != verdictNoReview
	items := notableChanges(artifact, catalog, summary, summaryContext, includeHeuristics, policy)
	var body strings.Builder
	body.WriteString(githubPRBodyMarker)
	body.WriteString("\n\n")
	body.WriteString(renderVerdictBanner(verdict, reason))
	body.WriteString("\n\n")
	body.WriteString(openingSummary(artifact, catalog, summary.Overview, aiSucceeded))
	body.WriteString(renderBlastRadiusSection(artifact, catalog, reach, policy, triage, summary, aiSucceeded))
	linkCtx := newAttributionLinkContext(artifact, catalog, summaryContext.Snippets)
	body.WriteString(renderNotableChangesSection(items, linkCtx))
	body.WriteString("\n\n")
	body.WriteString(provenanceFooter(aiSucceeded, reviewerInfo, summaryContext, linkCtx))
	return strings.TrimRight(body.String(), "\n")
}

func renderNotableChangesSection(items []prNotableChange, linkCtx attributionLinkContext) string {
	if len(items) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n\n## Notable Changes\n\n")
	for _, item := range items {
		b.WriteString("- ")
		if item.Flagged {
			b.WriteString("⚠ ")
		}
		if item.Link != "" {
			b.WriteString(formatPRSummaryLink(item.Title, item.Link))
		} else {
			b.WriteString(item.Title)
		}
		if item.Detail != "" {
			b.WriteString(" — ")
			b.WriteString(linkCtx.linkifyAttributionText(item.Detail))
		}
		b.WriteByte('\n')
		if attribution := renderAttributions(item.Attributions, linkCtx); attribution != "" {
			b.WriteString("  Attribution: ")
			b.WriteString(attribution)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func provenanceFooter(aiSucceeded bool, info codereview.ReviewerInfo, summaryContext prSummaryContext, linkCtx attributionLinkContext) string {
	if !aiSucceeded {
		return "*Generated by GX — heuristics only (AI unavailable).*"
	}
	var footer strings.Builder
	footer.WriteString("*Generated by GX")
	if len(info.Models) > 0 {
		footer.WriteString(" — model ")
		footer.WriteString(strings.Join(info.Models, ", "))
	}
	if sources := sortedUnique(summaryContext.Sources); len(sources) > 0 {
		if len(info.Models) > 0 {
			footer.WriteString("; ")
		} else {
			footer.WriteString(" — ")
		}
		footer.WriteString("context: ")
		linked := make([]string, 0, len(sources))
		for _, source := range sources {
			linked = append(linked, linkProvenanceSource(source, linkCtx.snippets))
		}
		footer.WriteString(strings.Join(linked, ", "))
	}
	footer.WriteString(".*")
	return footer.String()
}

func openingSummary(artifact reviewbundle.Artifact, catalog prBodyCatalog, overview string, aiSucceeded bool) string {
	if paragraph := sanitizedOverviewParagraph(overview, aiSucceeded); paragraph != "" {
		return paragraph
	}
	return changeSummarySentence(artifact, catalog)
}

func sanitizedOverviewParagraph(overview string, aiSucceeded bool) string {
	if !aiSucceeded {
		return ""
	}
	overview = strings.TrimSpace(overview)
	if overview == "" {
		return ""
	}
	overview = prURLStripRE.ReplaceAllString(overview, "")
	overview = sanitizePRVisibleText(overview)
	return trimSentence(overview, maxPROverviewLength)
}

func changeSummarySentence(artifact reviewbundle.Artifact, catalog prBodyCatalog) string {
	if len(catalog.Revisions) == 0 {
		return "This PR has no recorded GX revisions."
	}
	descriptions := revisionDescriptions(catalog.Revisions, 3)
	areas := humanAreaList(catalog.Areas, 3)
	if len(descriptions) == 1 {
		target := "the stack"
		if areas != "" {
			target = areas
		}
		return fmt.Sprintf("This PR changes %s: %s.", target, lowerFirst(descriptions[0]))
	}
	if name := stackName(artifact); name != "" {
		return fmt.Sprintf("This PR publishes %d GX revisions for `%s`: %s.", len(catalog.Revisions), name, strings.Join(descriptions, "; "))
	}
	return fmt.Sprintf("This PR publishes %d GX revisions: %s.", len(catalog.Revisions), strings.Join(descriptions, "; "))
}

func readinessSentence(catalog prBodyCatalog, reach lexicalReach) string {
	blast := blastRadiusLevelLabel(catalog.Stats.MaxRiskLevel)
	detail := fmt.Sprintf("%d file(s), %d area(s), +%d/-%d lines", catalog.Stats.FileCount, catalog.Stats.AreaCount, catalog.Stats.AddedLines, catalog.Stats.DeletedLines)
	if reach.ReferenceCount > 0 {
		symbols := strings.Join(topReachSymbolNames(reach, 2), ", ")
		reachDetail := fmt.Sprintf("%d references to changed symbols (%s) across %d other file(s); %d file(s), +%d/-%d lines",
			reach.ReferenceCount, symbols, reach.DependentFiles, catalog.Stats.FileCount, catalog.Stats.AddedLines, catalog.Stats.DeletedLines)
		return fmt.Sprintf("%s: %s.", blast, reachDetail)
	}
	return fmt.Sprintf("%s (%s).", blast, detail)
}

func blastRadiusLevelLabel(level string) string {
	return strings.ToUpper(firstNonEmpty(strings.TrimSpace(level), "low"))
}

func triageChangeFromCatalog(catalog prBodyCatalog) codereview.ChangeTriage {
	snippets := make([]codereview.DiffSnippet, 0, len(catalog.Hunks))
	for _, hunk := range catalog.Hunks {
		snippets = append(snippets, codereview.DiffSnippet{File: hunk.File, Diff: hunk.Patch})
	}
	return codereview.TriageChange(catalog.Files, snippets, codereview.Options{})
}

func reviewVerdict(triage codereview.ChangeTriage, stats prBodyStats, findings []codereview.Finding, aiSucceeded bool, files, areas []string) (string, string) {
	if triage.Class == "security-sensitive" {
		return verdictDeepReview, securitySensitiveVerdictReason(files, areas)
	}
	if stats.MaxRiskLevel == "high" {
		return verdictDeepReview, highRiskVerdictReason(stats)
	}
	if hasStrongOrBlockingFinding(findings) {
		return verdictDeepReview, "review findings need human judgment"
	}
	if aiSucceeded && isLowImpactTriageClass(triage.Class) && stats.WarningCount == 0 && stats.MaxRiskLevel == "low" {
		return verdictNoReview, noReviewVerdictReason(triage.Class, stats.FileCount)
	}
	if !aiSucceeded {
		return verdictQuickScan, "standard change; AI review unavailable"
	}
	return verdictQuickScan, "standard change; skim the notable changes"
}

func securitySensitiveVerdictReason(files, areas []string) string {
	target := firstNonEmpty(firstSorted(areas), verdictReasonFileArea(files))
	if target == "" {
		target = "code"
	}
	return "security-sensitive change: " + target
}

func highRiskVerdictReason(stats prBodyStats) string {
	for _, signal := range stats.RiskSignals {
		if strings.HasPrefix(signal, "lexical_reach:") {
			return "changes are referenced widely across the codebase"
		}
	}
	return "high-risk change signals"
}

func noReviewVerdictReason(class string, fileCount int) string {
	return fmt.Sprintf("%s change (%d file(s)); nothing needs human eyes", triageClassLabel(class), fileCount)
}

func triageClassLabel(class string) string {
	switch class {
	case "docs-only":
		return "documentation-only"
	case "tests-only":
		return "tests-only"
	case "config-only":
		return "configuration-only"
	case "mechanical":
		return "mechanical"
	default:
		return class
	}
}

func verdictReasonFileArea(files []string) string {
	if len(files) == 0 {
		return ""
	}
	sorted := sortedUnique(files)
	return fileArea(sorted[0])
}

func firstSorted(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return sortedUnique(values)[0]
}

func renderVerdictBanner(verdict, reason string) string {
	return fmt.Sprintf("> %s **%s** — %s.", verdictEmoji(verdict), verdict, reason)
}

func verdictEmoji(verdict string) string {
	switch verdict {
	case verdictNoReview:
		return "✅"
	case verdictQuickScan:
		return "👀"
	case verdictDeepReview:
		return "🔴"
	default:
		return ""
	}
}

func isLowImpactTriageClass(class string) bool {
	switch class {
	case "docs-only", "tests-only", "config-only", "mechanical":
		return true
	default:
		return false
	}
}

func hasStrongOrBlockingFinding(findings []codereview.Finding) bool {
	for _, finding := range findings {
		switch strings.TrimSpace(finding.Strength) {
		case "Blocking", "Strong":
			return true
		}
	}
	return false
}

func notableChanges(artifact reviewbundle.Artifact, catalog prBodyCatalog, summary codereview.PRSummaryReview, summaryContext prSummaryContext, includeHeuristics bool, policy codereview.ReviewPolicy) []prNotableChange {
	prURL := pullRequestURL(artifact)
	var items []prNotableChange
	for _, finding := range summary.Findings {
		if isSpeculativePRFinding(finding) || !patchRelevantFinding(catalog, finding) || lowValuePRFinding(finding) {
			continue
		}
		title := sanitizePRVisibleText(strings.TrimSpace(finding.Title))
		if title == "" {
			continue
		}
		detail := firstNonEmpty(strings.TrimSpace(finding.Summary), strings.TrimSpace(finding.Recommendation))
		item := prNotableChange{
			Title:        title,
			Detail:       trimSentence(sanitizePRVisibleText(detail), 240),
			Link:         hunkLinkForFinding(prURL, catalog.Hunks, finding),
			Flagged:      true,
			Score:        1000 + findingScore(catalog.Hunks, finding),
			Attributions: findingAttributions(finding, summaryContext, artifact, catalog),
		}
		if includeNotableChangeItem(item, prURL) {
			items = append(items, item)
		}
	}
	if includeHeuristics {
		items = append(items, heuristicNotableChangeItems(artifact, catalog, policy)...)
	}
	for _, change := range summary.NotableChanges {
		title := trimSentence(sanitizePRVisibleText(strings.TrimSpace(change.Note)), 120)
		if title == "" {
			continue
		}
		item := prNotableChange{
			Title:  title,
			Link:   lineLinkForAnchor(prURL, catalog.Hunks, change.File, change.Line),
			Score:  800,
			Flagged: false,
		}
		if includeNotableChangeItem(item, prURL) {
			items = append(items, item)
		}
	}
	items = dedupeNotableChanges(items)
	sortNotableChanges(items)
	if len(catalog.Hunks) > 0 && len(items) < 3 {
		items = appendTopHunkFill(items, catalog, prURL)
		items = dedupeNotableChanges(items)
		sortNotableChanges(items)
	}
	if len(items) > maxPRNotableChanges {
		items = items[:maxPRNotableChanges]
	}
	return items
}

func includeNotableChangeItem(item prNotableChange, prURL string) bool {
	if prURL == "" {
		return true
	}
	return item.Link != ""
}

func sortNotableChanges(items []prNotableChange) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].Title < items[j].Title
		}
		return items[i].Score > items[j].Score
	})
}

func appendTopHunkFill(items []prNotableChange, catalog prBodyCatalog, prURL string) []prNotableChange {
	usedLinks := map[string]struct{}{}
	for _, item := range items {
		if item.Link != "" {
			usedLinks[item.Link] = struct{}{}
		}
	}
	hunks := append([]prHunkSummary(nil), catalog.Hunks...)
	sort.Slice(hunks, func(i, j int) bool {
		return hunks[i].NewLines+hunks[i].OldLines > hunks[j].NewLines+hunks[j].OldLines
	})
	for _, hunk := range hunks {
		if len(items) >= 3 {
			break
		}
		if hunk.Link == "" && prURL != "" {
			continue
		}
		if hunk.Link != "" {
			if _, ok := usedLinks[hunk.Link]; ok {
				continue
			}
		}
		revision := firstNonEmpty(strings.TrimSpace(hunk.Revision), "change")
		title := trimSentence(sanitizePRVisibleText(fmt.Sprintf("%s: %s", revision, path.Base(hunk.File))), 120)
		items = append(items, prNotableChange{
			Title:  title,
			Link:   hunk.Link,
			Score:  400,
			Flagged: false,
			Attributions: []prAttribution{{
				Kind:  "codebase",
				Label: "changed hunk",
				Ref:   hunk.File,
			}},
		})
		if hunk.Link != "" {
			usedLinks[hunk.Link] = struct{}{}
		}
	}
	return items
}

func isSpeculativePRFinding(finding codereview.Finding) bool {
	switch strings.TrimSpace(finding.Strength) {
	case "Blocking", "Strong", "Worth exploring":
		return false
	default:
		return true
	}
}

func patchRelevantFinding(catalog prBodyCatalog, finding codereview.Finding) bool {
	file := strings.TrimSpace(finding.File)
	if file != "" {
		return true
	}
	if len(catalog.Hunks) == 0 {
		return true
	}
	text := strings.ToLower(strings.Join([]string{finding.Title, finding.Summary, finding.Recommendation, evidenceText(finding.Evidence)}, " "))
	for _, hunk := range catalog.Hunks {
		hunkFile := strings.ToLower(hunk.File)
		if hunkFile != "" && (strings.Contains(text, hunkFile) || strings.Contains(text, strings.ToLower(path.Base(hunkFile)))) {
			return true
		}
	}
	return false
}

func lowValuePRFinding(finding codereview.Finding) bool {
	text := strings.ToLower(strings.Join([]string{finding.Title, finding.Summary, finding.Recommendation}, " "))
	return strings.Contains(text, "this is good") || strings.Contains(text, "this is correct")
}

func findingScore(hunks []prHunkSummary, finding codereview.Finding) int {
	score := 0
	switch strings.TrimSpace(finding.Strength) {
	case "Blocking":
		score = 300
	case "Strong":
		score = 200
	case "Worth exploring":
		score = 100
	}
	if validatedHunkAnchor(hunks, finding) {
		score += 50
	}
	return score
}

func sanitizePRVisibleText(value string) string {
	replacements := []struct {
		old string
		new string
	}{
		{old: "Published by GX.", new: "the legacy GX placeholder"},
		{old: "Published by GX", new: "legacy GX placeholder"},
		{old: "## Summary", new: "summary section"},
		{old: "## Review signals", new: "review-signal section"},
		{old: "## Links", new: "links section"},
		{old: "Files changed", new: "changed files"},
	}
	for _, replacement := range replacements {
		value = strings.ReplaceAll(value, replacement.old, replacement.new)
	}
	return strings.TrimSpace(value)
}

func heuristicNotableChangeItems(artifact reviewbundle.Artifact, catalog prBodyCatalog, policy codereview.ReviewPolicy) []prNotableChange {
	prURL := pullRequestURL(artifact)
	var candidates []prNotableChange
	for _, warning := range reviewWarnings(artifact) {
		file := firstNonEmpty(warning.FromFile, warning.ToFile)
		hunk := firstHunk(catalog.Hunks, file)
		if hunk == nil {
			hunk = firstHunk(catalog.Hunks, "")
		}
		score := 760
		if warning.Severity == "warning" {
			score += 120
		}
		title := "Verify " + humanizeSignal(firstNonEmpty(warning.Source, "review warning"))
		if hunk != nil && hunk.File != "" {
			title += " in " + path.Base(hunk.File)
		}
		item := prNotableChange{
			Title:  title,
			Detail: trimSentence(firstNonEmpty(warning.Message, "GX found a feasibility or structural warning tied to this change."), 240),
			Link:   hunkLink(hunk),
			Score:  score,
			Flagged: true,
			Attributions: []prAttribution{{
				Kind:  "heuristic",
				Label: firstNonEmpty(warning.Source, "feasibility warning"),
				Ref:   firstNonEmpty(file, warning.Symbol, warning.DependsOn),
			}},
		}
		if includeNotableChangeItem(item, prURL) {
			candidates = append(candidates, item)
		}
	}
	for _, hunk := range catalog.Hunks {
		score, title, detail := hunkImpact(catalog, hunk, policy)
		if score == 0 {
			continue
		}
		flagged := hunkMatchesRiskPath(hunk.File, policy)
		item := prNotableChange{
			Title:  title,
			Detail: detail,
			Link:   hunk.Link,
			Score:  score,
			Flagged: flagged,
			Attributions: []prAttribution{{
				Kind:  "codebase",
				Label: "changed hunk",
				Ref:   hunk.File,
			}},
		}
		if includeNotableChangeItem(item, prURL) {
			candidates = append(candidates, item)
		}
	}
	if highRiskCatalog(catalog) {
		if hunk := firstHunk(catalog.Hunks, ""); hunk != nil {
			item := prNotableChange{
				Title:  "Review the broadest behavioral change",
				Detail: "GX scored this stack as high risk from breadth, structural dependencies, or warning signals. Start with this changed hunk and follow its call path.",
				Link:   hunk.Link,
				Score:  520,
				Flagged: false,
				Attributions: []prAttribution{{
					Kind:  "heuristic",
					Label: "risk signals",
					Ref:   strings.Join(catalog.Stats.RiskSignals, ", "),
				}},
			}
			if includeNotableChangeItem(item, prURL) {
				candidates = append(candidates, item)
			}
		}
	}
	candidates = dedupeNotableChanges(candidates)
	sortNotableChanges(candidates)
	return candidates
}

func hunkMatchesRiskPath(file string, policy codereview.ReviewPolicy) bool {
	for _, riskPath := range policy.RiskPaths {
		if codereview.MatchRiskPathGlob(riskPath.Glob, file) {
			return true
		}
	}
	return false
}

func hunkImpact(catalog prBodyCatalog, hunk prHunkSummary, policy codereview.ReviewPolicy) (int, string, string) {
	file := strings.TrimSpace(hunk.File)
	if file == "" {
		return 0, "", ""
	}
	base := path.Base(file)
	lowerFile := strings.ToLower(file + " " + base)

	for _, riskPath := range policy.RiskPaths {
		if !codereview.MatchRiskPathGlob(riskPath.Glob, file) {
			continue
		}
		title := riskPathTitle(riskPath.Message, base)
		detail := fmt.Sprintf("%s (matched `%s` in %s).", strings.TrimSpace(riskPath.Message), riskPath.Glob, file)
		return 700, title, trimSentence(detail, 240)
	}

	if containsAny(lowerFile, []string{"auth", "token", "secret", "credential"}) {
		pattern := "auth|token|secret|credential"
		title := fmt.Sprintf("Verify auth handling in %s", base)
		detail := fmt.Sprintf("Path or filename matches generic pattern %q in %s; auth changes can leak, drop, or misuse credentials.", pattern, file)
		return 650, title, trimSentence(detail, 240)
	}
	if containsAny(lowerFile, []string{"migration", "schema"}) {
		pattern := "migration|schema"
		title := fmt.Sprintf("Verify schema change in %s", base)
		detail := fmt.Sprintf("Path matches generic pattern %q in %s; schema changes can break persistence or migrations.", pattern, file)
		return 670, title, trimSentence(detail, 240)
	}
	if strings.Contains(lowerFile, "dockerfile") || strings.Contains(lowerFile, ".github/workflows") || strings.Contains(lowerFile, "deploy") {
		pattern := "Dockerfile|.github/workflows|deploy"
		title := fmt.Sprintf("Verify deployment change in %s", base)
		detail := fmt.Sprintf("Path matches generic pattern %q in %s; deployment changes can alter runtime or CI behavior.", pattern, file)
		return 660, title, trimSentence(detail, 240)
	}
	if isGenericLockfilePath(lowerFile) {
		pattern := "lockfiles"
		title := fmt.Sprintf("Verify lockfile change in %s", base)
		detail := fmt.Sprintf("Filename matches generic pattern %q in %s; dependency lockfile changes can alter build or runtime versions.", pattern, file)
		return 640, title, trimSentence(detail, 240)
	}
	if catalog.Stats.MaxRiskLevel == "high" && hunk.NewLines+hunk.OldLines >= 40 {
		title := fmt.Sprintf("Review large high-risk hunk in %s", base)
		detail := fmt.Sprintf("%s is one of the larger hunks in a high-risk stack and may hide cross-cutting behavior.", file)
		return 500, title, trimSentence(detail, 240)
	}
	return 0, "", ""
}

func riskPathTitle(message, base string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Sprintf("Verify change in %s", base)
	}
	if idx := strings.IndexAny(message, ".;"); idx > 0 {
		clause := strings.TrimSpace(message[:idx])
		if clause != "" {
			return "Verify " + clause
		}
	}
	return fmt.Sprintf("Verify change in %s", base)
}

func isGenericLockfilePath(file string) bool {
	switch strings.ToLower(path.Base(file)) {
	case "go.sum", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb", "cargo.lock", "gemfile.lock", "poetry.lock", "composer.lock":
		return true
	default:
		return false
	}
}

func highRiskCatalog(catalog prBodyCatalog) bool {
	return catalog.Stats.MaxRiskLevel == "high" || containsAny(strings.Join(catalog.Stats.RiskSignals, " "), []string{"structural", "warning:", "many_files", "many_hunks"})
}

func reviewPRSummaryFindings(ctx context.Context, artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext, reach lexicalReach, policy codereview.ReviewPolicy) (codereview.PRSummaryReview, bool, codereview.ReviewerInfo, error) {
	reviewer, info := prSummaryReviewerFromEnvWithInfo()
	if reviewer == nil {
		return codereview.PRSummaryReview{}, false, info, nil
	}
	brief := prReviewBrief(artifact, catalog, summaryContext, reach, policy)
	attemptReview := func() (codereview.PRSummaryReview, error) {
		if withSummary, ok := reviewer.(codereview.AIReviewerWithSummary); ok {
			return withSummary.ReviewForSummary(ctx, brief)
		}
		if withOverview, ok := reviewer.(codereview.AIReviewerWithOverview); ok {
			overview, findings, err := withOverview.ReviewWithOverview(ctx, brief)
			return codereview.PRSummaryReview{Overview: overview, Findings: findings}, err
		}
		findings, err := reviewer.Review(ctx, brief)
		return codereview.PRSummaryReview{Findings: findings}, err
	}
	summary, err := attemptReview()
	if err != nil {
		if cloud.IsPaymentRequired(err) {
			return codereview.PRSummaryReview{}, false, info, err
		}
		if waitErr := prSummaryReviewRetryWait(ctx, prSummaryReviewRetryDelay); waitErr != nil {
			return codereview.PRSummaryReview{}, false, info, nil
		}
		summary, err = attemptReview()
	}
	if err != nil {
		if cloud.IsPaymentRequired(err) {
			return codereview.PRSummaryReview{}, false, info, err
		}
		return codereview.PRSummaryReview{}, false, info, nil
	}
	return summary, true, info, nil
}

func prReviewBrief(artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext, reach lexicalReach, policy codereview.ReviewPolicy) codereview.ReviewBrief {
	diffSnippets := make([]codereview.DiffSnippet, 0, len(catalog.Hunks))
	for _, hunk := range catalog.Hunks {
		diffSnippets = append(diffSnippets, codereview.DiffSnippet{File: hunk.File, Diff: hunk.Patch})
		if len(diffSnippets) >= 32 {
			break
		}
	}
	snippets := append(append([]codereview.ContextSnippet(nil), summaryContext.Snippets...), lexicalReachContextSnippets(reach)...)
	labeled := codereview.LabelContextSnippets(snippets)
	focus := strings.Join(catalog.Files, " ")
	if len(policy.RiskPaths) > 0 {
		paths := make([]string, 0, len(policy.RiskPaths))
		for _, riskPath := range policy.RiskPaths {
			paths = append(paths, riskPath.Glob)
		}
		focus = focus + "; configured high-risk paths: " + strings.Join(paths, ", ")
	}
	return codereview.ReviewBrief{
		RepoRoot:      artifact.Bundle.Repo.RootPath,
		Scope:         "maintainability",
		Depth:         "shallow",
		ReviewProfile: "pr_summary",
		Focus:         focus,
		Static: codereview.StaticSnapshot{
			FileCount:     catalog.Stats.FileCount,
			ChangedFiles:  catalog.Files,
			DiffSnippets:  diffSnippets,
			CodeQuality:   warningQualityHints(artifact),
			TestFileCount: countTestFiles(catalog.Files),
		},
		Context:       labeled,
		SourceRefs:    codereview.SourceRefsFromContextSnippets(labeled),
		SourceCatalog: codereview.SourceBriefs(codereview.SourcesForScopes([]string{"maintainability", "security", "testing", "architecture", "dependencies"})),
		Rubric: codereview.ArchitectureRubric{
			Goal: "Identify the changed hunks a human should review before merging this PR.",
			Questions: []string{
				"Which changed hunks most affect the project's behavior, persistence/state, remote side effects, security, or developer workflow?",
				"Which files could create user-visible or developer-visible regressions?",
				"What exact behavior should a reviewer verify first?",
			},
			Reject: []string{
				"Do not emit generic best-practice advice.",
				"Do not include low-impact nits or quota-filling items.",
				"Do not list whole files when a changed hunk is available.",
				"Do not comment on tests unless missing coverage hides a concrete behavior risk.",
			},
			Output: "Zero to ten concrete findings tied to changed hunks, each with evidence.",
		},
	}
}

func collectDefaultPRSummaryContext(ctx context.Context, artifact reviewbundle.Artifact, catalog prBodyCatalog) (prSummaryContext, error) {
	var snippets []codereview.ContextSnippet
	var sources []string
	if local := localPRContextSnippets(artifact, catalog); len(local) > 0 {
		snippets = append(snippets, local...)
		sources = append(sources, "codebase")
	}
	if session := sessionPRContextSnippets(artifact); len(session) > 0 {
		snippets = append(snippets, session...)
		sources = append(sources, "session")
	}
	indexed, indexedSources := indexedPRContextSnippets(ctx, artifact, catalog)
	if len(indexed) > 0 {
		snippets = append(snippets, indexed...)
		sources = append(sources, indexedSources...)
	}
	return prSummaryContext{Snippets: dedupePRContextSnippets(snippets), Sources: sortedUnique(sources)}, nil
}

func localPRContextSnippets(artifact reviewbundle.Artifact, catalog prBodyCatalog) []codereview.ContextSnippet {
	root := strings.TrimSpace(artifact.Bundle.Repo.RootPath)
	if root == "" {
		return nil
	}
	var snippets []codereview.ContextSnippet
	for _, rel := range []string{"AGENTS.md", "CONTEXT.md", "REVIEW.md", "README.md"} {
		if snippet, ok := readPRContextSnippet(root, rel, "codebase_doc"); ok {
			snippets = append(snippets, snippet)
		}
	}
	for _, file := range catalog.Files {
		if len(snippets) >= 12 {
			break
		}
		if snippet, ok := readPRContextSnippet(root, file, "changed_file"); ok {
			snippets = append(snippets, snippet)
		}
	}
	return snippets
}

func readPRContextSnippet(root, rel, kind string) (codereview.ContextSnippet, bool) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil || strings.TrimSpace(string(data)) == "" {
		return codereview.ContextSnippet{}, false
	}
	return codereview.ContextSnippet{
		Kind:        kind,
		Ref:         rel,
		Source:      "local",
		SourceLabel: rel,
		Publisher:   "this repo",
		Text:        limitText(string(data), maxPRContextSnippetSize),
	}, true
}

func sessionPRContextSnippets(artifact reviewbundle.Artifact) []codereview.ContextSnippet {
	chunks := semantic.BuildSessionChunks(artifact.Bundle, maxPRContextSnippetSize)
	out := make([]codereview.ContextSnippet, 0, len(chunks))
	for _, chunk := range chunks {
		ref := stringAttribute(chunk.Attributes, "source_id")
		if ref == "" {
			ref = chunk.ID
		}
		out = append(out, codereview.ContextSnippet{
			Kind:        "session_transcript",
			Ref:         ref,
			Source:      "artifact",
			SourceLabel: "linked session",
			Publisher:   "this repo",
			Text:        limitText(chunk.Text, maxPRContextSnippetSize),
		})
		if len(out) >= 4 {
			break
		}
	}
	return out
}

func indexedPRContextSnippets(ctx context.Context, artifact reviewbundle.Artifact, catalog prBodyCatalog) ([]codereview.ContextSnippet, []string) {
	cfg, err := semantic.ConfigFromEnv()
	if err != nil || !cfg.Enabled {
		return nil, nil
	}
	queryText := prContextQueryText(artifact, catalog)
	if queryText == "" {
		return nil, nil
	}
	vectors, err := semantic.NewOpenAIEmbedder(cfg).Embed(ctx, []string{queryText})
	if err != nil || len(vectors) != 1 {
		return nil, nil
	}
	var snippets []codereview.ContextSnippet
	var sources []string
	if cfg.TurboPufferNamespace != "" {
		rows := queryPRTurboPuffer(ctx, cfg, cfg.TurboPufferNamespace, vectors[0], queryText, "text", []any{"Or", []any{[]any{"source_kind", "Eq", "code_file"}, []any{"source_kind", "Eq", "session_transcript"}}}, 6)
		if len(rows) > 0 {
			snippets = append(snippets, rowsToPRContextSnippets(rows, "turbopuffer:"+cfg.TurboPufferNamespace)...)
			sources = append(sources, "indexed code/session")
		}
	}
	reviewNamespace := strings.TrimSpace(os.Getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE"))
	if reviewNamespace == "" {
		reviewNamespace = "gx-review-knowledge"
	}
	rows := queryPRTurboPuffer(ctx, cfg, reviewNamespace, vectors[0], queryText, "body", []any{"And", []any{
		[]any{"source_kind", "Eq", "review_corpus"},
		[]any{"tier", "NotEq", "research"},
		[]any{"historical", "Eq", false},
		[]any{"superseded_by", "Eq", ""},
	}}, 4)
	if len(rows) > 0 {
		snippets = append(snippets, rowsToPRContextSnippets(rows, "turbopuffer:"+reviewNamespace)...)
		sources = append(sources, "review resources")
	}
	return snippets, sortedUnique(sources)
}

func queryPRTurboPuffer(ctx context.Context, cfg semantic.Config, namespace string, vector []float32, text string, textAttribute string, filters any, limit int) []map[string]any {
	vectorQuery := map[string]any{
		"rank_by": []any{"vector", "ANN", vector},
		"limit":   map[string]any{"total": limit},
		"include_attributes": []string{
			"body", "text", "source_kind", "source_id", "file_path", "symbol", "start_line", "end_line", "title", "url", "source_url", "category", "tier",
		},
	}
	if filters != nil {
		vectorQuery["filters"] = filters
	}
	payload := vectorQuery
	if strings.TrimSpace(text) != "" {
		textQuery := map[string]any{
			"rank_by":            []any{textAttribute, "BM25", text},
			"limit":              map[string]any{"total": limit},
			"include_attributes": vectorQuery["include_attributes"],
		}
		if filters != nil {
			textQuery["filters"] = filters
		}
		payload = map[string]any{
			"queries":   []any{vectorQuery, textQuery},
			"rerank_by": []any{"RRF"},
		}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	endpoint := strings.TrimRight(cfg.TurboPufferBaseURL, "/") + "/v2/namespaces/" + url.PathEscape(strings.Trim(namespace, "/")) + "/query"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+cfg.TurboPufferAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil
	}
	var decoded struct {
		Rows    []map[string]any `json:"rows"`
		Results []struct {
			Rows []map[string]any `json:"rows"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil
	}
	if len(decoded.Rows) == 0 && len(decoded.Results) > 0 {
		return decoded.Results[0].Rows
	}
	return decoded.Rows
}

func rowsToPRContextSnippets(rows []map[string]any, source string) []codereview.ContextSnippet {
	var snippets []codereview.ContextSnippet
	for _, row := range rows {
		text := firstNonEmpty(stringAny(row["body"]), stringAny(row["text"]))
		if text == "" {
			continue
		}
		ref := firstNonEmpty(stringAny(row["file_path"]), stringAny(row["source_id"]), stringAny(row["source_url"]), stringAny(row["url"]), stringAny(row["title"]))
		kind := firstNonEmpty(stringAny(row["source_kind"]), "indexed_context")
		snippets = append(snippets, codereview.ContextSnippet{
			Kind:        kind,
			Ref:         ref,
			Source:      source,
			SourceLabel: ref,
			Publisher:   firstNonEmpty(stringAny(row["publisher"]), "this repo"),
			Text:        limitText(text, maxPRContextSnippetSize),
			URL:         firstNonEmpty(stringAny(row["source_url"]), stringAny(row["url"])),
		})
	}
	return snippets
}

func prContextQueryText(artifact reviewbundle.Artifact, catalog prBodyCatalog) string {
	var parts []string
	parts = append(parts, "GX PR summary context query")
	if name := stackName(artifact); name != "" {
		parts = append(parts, "stack: "+name)
	}
	for _, revision := range catalog.Revisions {
		if revision.Description != "" {
			parts = append(parts, "revision: "+revision.Description)
		}
	}
	parts = append(parts, "files: "+strings.Join(catalog.Files, " "))
	parts = append(parts, "risk signals: "+strings.Join(catalog.Stats.RiskSignals, " "))
	return strings.Join(parts, "\n")
}

func buildPRBodyCatalog(artifact reviewbundle.Artifact) prBodyCatalog {
	revisions := revisionSummaries(artifact)
	extraRevision, extraPatch, hasExtraRevision := headRevisionFromGit(artifact, revisions)
	if hasExtraRevision {
		revisions = append(revisions, extraRevision)
	}
	hunks := pullRequestHunks(artifact)
	if hasExtraRevision {
		hunks = append(hunks, parsePatchHunks(extraPatch, extraRevision.Description, pullRequestURL(artifact), len(hunks))...)
	}
	stats := bodyStats(artifact, revisions, []string{extraPatch})
	filesSet := map[string]struct{}{}
	areasSet := map[string]struct{}{}
	for _, revision := range revisions {
		for _, file := range revision.Files {
			filesSet[file] = struct{}{}
			areasSet[fileArea(file)] = struct{}{}
		}
	}
	for _, hunk := range hunks {
		filesSet[hunk.File] = struct{}{}
		areasSet[fileArea(hunk.File)] = struct{}{}
	}
	return prBodyCatalog{Revisions: revisions, Hunks: hunks, Stats: stats, Files: sortedKeys(filesSet), Areas: sortedKeys(areasSet)}
}

func headRevisionFromGit(artifact reviewbundle.Artifact, revisions []prRevisionSummary) (prRevisionSummary, string, bool) {
	head := strings.TrimSpace(artifact.Push.HeadCommitID)
	root := strings.TrimSpace(artifact.Bundle.Repo.RootPath)
	if head == "" || root == "" {
		return prRevisionSummary{}, "", false
	}
	for _, revision := range revisions {
		if strings.EqualFold(strings.TrimSpace(revision.CommitID), head) {
			return prRevisionSummary{}, "", false
		}
	}
	base := ""
	for i := len(revisions) - 1; i >= 0; i-- {
		if commit := strings.TrimSpace(revisions[i].CommitID); commit != "" {
			base = commit
			break
		}
	}
	if base == "" {
		return prRevisionSummary{}, "", false
	}
	patch := gitOutput(root, "diff", "--no-ext-diff", base+".."+head)
	if strings.TrimSpace(patch) == "" {
		return prRevisionSummary{}, "", false
	}
	description := strings.TrimSpace(gitOutput(root, "show", "-s", "--format=%s", head))
	if description == "" {
		description = "update PR head"
	}
	return prRevisionSummary{Description: description, CommitID: head, BranchName: stackName(artifact), Files: filesFromPatch(patch)}, patch, true
}

func gitOutput(repoRoot string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = bytes.NewBuffer(nil)
	if err := cmd.Run(); err != nil {
		return ""
	}
	return out.String()
}

func revisionSummaries(artifact reviewbundle.Artifact) []prRevisionSummary {
	if len(artifact.Stack) > 0 {
		out := make([]prRevisionSummary, 0, len(artifact.Stack))
		for _, entry := range artifact.Stack {
			out = append(out, revisionSummaryFromChange(entry.Change, entry.BranchName, entry.Patch))
		}
		return out
	}
	if artifact.Change != nil {
		return []prRevisionSummary{revisionSummaryFromChange(*artifact.Change, pointerString(artifact.Push.BranchName), "")}
	}
	return nil
}

func revisionSummaryFromChange(change reviewbundle.ChangePayload, branchName, patchText string) prRevisionSummary {
	files := append([]string(nil), change.Files...)
	if len(files) == 0 {
		files = filesFromPatch(patchText)
	}
	summary := prRevisionSummary{Description: change.Description, CommitID: change.CurrentCommitID, BranchName: branchName, Files: sortedUnique(files)}
	if change.ReviewContext != nil {
		summary.RiskLevel = change.ReviewContext.Risk.Level
		summary.RiskScore = change.ReviewContext.Risk.Score
		summary.RiskSignals = append([]string(nil), change.ReviewContext.Risk.Signals...)
	}
	return summary
}

func pullRequestHunks(artifact reviewbundle.Artifact) []prHunkSummary {
	prURL := pullRequestURL(artifact)
	var hunks []prHunkSummary
	for revisionIndex, entry := range artifact.Stack {
		revisionLabel := firstNonEmpty(strings.TrimSpace(entry.Change.Description), fmt.Sprintf("revision %d", revisionIndex+1))
		hunks = append(hunks, parsePatchHunks(entry.Patch, revisionLabel, prURL, len(hunks))...)
	}
	return hunks
}

func parsePatchHunks(patchText, revisionLabel, prURL string, offset int) []prHunkSummary {
	var currentFile string
	var current *prHunkSummary
	var hunks []prHunkSummary
	flush := func() {
		if current == nil {
			return
		}
		current.Patch = strings.TrimRight(current.Patch, "\n")
		current.ChangedLine = changedLineForHunk(*current)
		current.Link = githubHunkLink(prURL, current.File, current.OldStart, current.ChangedLine, current.NewLines)
		hunks = append(hunks, *current)
		current = nil
	}
	for _, line := range strings.Split(patchText, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			flush()
			currentFile = fileFromDiffLine(line)
		case strings.HasPrefix(line, "+++ "):
			if file := strings.TrimPrefix(strings.TrimSpace(line), "+++ b/"); file != "" && file != "/dev/null" {
				currentFile = file
			}
		case strings.HasPrefix(line, "@@ "):
			flush()
			oldStart, oldLines, newStart, newLines := parsePatchHunkHeader(line)
			current = &prHunkSummary{
				ID:       fmt.Sprintf("h%d", offset+len(hunks)+1),
				Revision: revisionLabel,
				File:     currentFile,
				Header:   line,
				OldStart: oldStart,
				OldLines: oldLines,
				NewStart: newStart,
				NewLines: newLines,
				Patch:    line + "\n",
			}
		default:
			if current != nil {
				current.Patch += line + "\n"
			}
		}
	}
	flush()
	return hunks
}

func parsePatchHunkHeader(header string) (int, int, int, int) {
	match := hunkHeaderRE.FindStringSubmatch(header)
	if len(match) == 0 {
		return 0, 0, 0, 0
	}
	return atoiDefault(match[1], 0), atoiDefault(match[2], 1), atoiDefault(match[3], 0), atoiDefault(match[4], 1)
}

var anyAddedPatchLineRE = regexp.MustCompile(`^\+[^+]`)

func lineOfPatchMatch(hunk prHunkSummary, re *regexp.Regexp) int {
	line := hunk.NewStart
	for _, patchLine := range strings.Split(hunk.Patch, "\n") {
		if strings.HasPrefix(patchLine, "@@ ") {
			line = hunk.NewStart
			continue
		}
		if strings.HasPrefix(patchLine, "+") && !strings.HasPrefix(patchLine, "+++") {
			if re == nil || re.MatchString(patchLine) {
				return line
			}
			line++
			continue
		}
		if strings.HasPrefix(patchLine, " ") || strings.HasPrefix(patchLine, "+") {
			line++
		}
	}
	return 0
}

func changedLineForHunk(hunk prHunkSummary) int {
	if line := lineOfPatchMatch(hunk, anyAddedPatchLineRE); line > 0 {
		return line
	}
	return hunk.NewStart
}

func githubPullRequestFilesURL(prURL string) string {
	prURL = strings.TrimRight(strings.TrimSpace(prURL), "/")
	if prURL == "" {
		return ""
	}
	parsed, err := url.Parse(prURL)
	if err == nil && parsed.Host != "" {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) >= 4 && parts[2] == "pull" {
			number, err := strconv.Atoi(parts[3])
			if err == nil && number > 0 {
				owner := parts[0]
				repo := strings.TrimSuffix(parts[1], ".git")
				scheme := parsed.Scheme
				if scheme == "" {
					scheme = "https"
				}
				return fmt.Sprintf("%s://%s/%s/%s/pull/%d/files", scheme, parsed.Host, owner, repo, number)
			}
		}
	}
	for _, suffix := range []string{"/files", "/changes", "/commits", "/checks", "/conversation"} {
		if strings.HasSuffix(prURL, suffix) {
			return strings.TrimSuffix(prURL, suffix) + "/files"
		}
	}
	return prURL + "/files"
}

func githubHunkLink(prURL, file string, oldStart, newStart, newLines int) string {
	filesURL := githubPullRequestFilesURL(prURL)
	file = strings.TrimSpace(file)
	if filesURL == "" || file == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(file))
	side := "R"
	line := newStart
	if newLines == 0 {
		side = "L"
		line = oldStart
	}
	if line <= 0 {
		line = 1
	}
	return fmt.Sprintf("%s#diff-%s%s%d", filesURL, hex.EncodeToString(sum[:]), side, line)
}

func githubFileDiffLink(prURL, file string) string {
	filesURL := githubPullRequestFilesURL(prURL)
	file = strings.TrimSpace(file)
	if filesURL == "" || file == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(file))
	return fmt.Sprintf("%s#diff-%s", filesURL, hex.EncodeToString(sum[:]))
}

func githubHunkLineLink(prURL, file string, line int) string {
	filesURL := githubPullRequestFilesURL(prURL)
	file = strings.TrimSpace(file)
	if filesURL == "" || file == "" || line <= 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(file))
	return fmt.Sprintf("%s#diff-%sR%d", filesURL, hex.EncodeToString(sum[:]), line)
}

func blobPermalink(prURL, sha, file string, line int) string {
	prURL = strings.TrimRight(strings.TrimSpace(prURL), "/")
	sha = strings.TrimSpace(sha)
	file = strings.TrimSpace(file)
	if prURL == "" || sha == "" || file == "" {
		return ""
	}
	owner, repo, ok := parseGitHubOwnerRepo(prURL)
	if !ok {
		return ""
	}
	link := fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", owner, repo, sha, file)
	if line > 0 {
		link += fmt.Sprintf("#L%d", line)
	}
	return link
}

func parseGitHubOwnerRepo(prURL string) (string, string, bool) {
	u, err := url.Parse(strings.TrimSpace(prURL))
	if err != nil || !strings.EqualFold(u.Host, "github.com") {
		return "", "", false
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func findingAnchorLine(finding codereview.Finding) int {
	if finding.Line > 0 {
		return finding.Line
	}
	file := strings.TrimSpace(finding.File)
	for _, anchor := range finding.Anchors {
		if strings.TrimSpace(anchor.File) == file && anchor.Line > 0 {
			return anchor.Line
		}
	}
	return 0
}

func lineInHunkNewRange(hunk prHunkSummary, line int) bool {
	if line <= 0 || hunk.NewLines <= 0 {
		return false
	}
	return line >= hunk.NewStart && line < hunk.NewStart+hunk.NewLines
}

func validatedHunkAnchor(hunks []prHunkSummary, finding codereview.Finding) bool {
	file := strings.TrimSpace(finding.File)
	line := findingAnchorLine(finding)
	if file == "" || line <= 0 {
		return false
	}
	for _, hunk := range hunks {
		if hunk.File == file && lineInHunkNewRange(hunk, line) {
			return true
		}
	}
	return false
}

func lineLinkForAnchor(prURL string, hunks []prHunkSummary, file string, line int) string {
	file = strings.TrimSpace(file)
	if file == "" {
		return ""
	}
	if line > 0 && validatedHunkAnchor(hunks, codereview.Finding{File: file, Line: line}) {
		return githubHunkLineLink(prURL, file, line)
	}
	for _, hunk := range hunks {
		if hunk.File == file {
			return hunk.Link
		}
	}
	return ""
}

func hunkLinkForFinding(prURL string, hunks []prHunkSummary, finding codereview.Finding) string {
	file := strings.TrimSpace(finding.File)
	if file == "" {
		search := strings.ToLower(strings.Join([]string{finding.Title, finding.Summary, finding.Recommendation, evidenceText(finding.Evidence)}, " "))
		for _, hunk := range hunks {
			if hunk.File != "" && strings.Contains(search, strings.ToLower(hunk.File)) {
				return hunk.Link
			}
		}
		for _, hunk := range hunks {
			if base := path.Base(hunk.File); base != "" && strings.Contains(search, strings.ToLower(base)) {
				return hunk.Link
			}
		}
		return ""
	}
	return lineLinkForAnchor(prURL, hunks, file, findingAnchorLine(finding))
}

func firstHunk(hunks []prHunkSummary, filePrefix string) *prHunkSummary {
	for i := range hunks {
		if filePrefix == "" || strings.HasPrefix(hunks[i].File, filePrefix) || strings.EqualFold(hunks[i].File, filePrefix) {
			return &hunks[i]
		}
	}
	return nil
}

func hunkLink(hunk *prHunkSummary) string {
	if hunk == nil {
		return ""
	}
	return hunk.Link
}

func warningQualityHints(artifact reviewbundle.Artifact) []codereview.CodeQualityHint {
	var hints []codereview.CodeQualityHint
	for _, warning := range reviewWarnings(artifact) {
		hints = append(hints, codereview.CodeQualityHint{
			Kind:   firstNonEmpty(warning.Source, "review_warning"),
			File:   firstNonEmpty(warning.FromFile, warning.ToFile),
			Text:   warning.Message,
			Reason: warning.Severity,
		})
	}
	return hints
}

func reviewWarnings(artifact reviewbundle.Artifact) []reviewbundle.ReviewFeasibilityWarning {
	var warnings []reviewbundle.ReviewFeasibilityWarning
	for _, entry := range artifact.Stack {
		if entry.Change.ReviewContext != nil {
			warnings = append(warnings, entry.Change.ReviewContext.FeasibilityWarnings...)
		}
	}
	if artifact.Change != nil && artifact.Change.ReviewContext != nil {
		warnings = append(warnings, artifact.Change.ReviewContext.FeasibilityWarnings...)
	}
	return warnings
}

func bodyStats(artifact reviewbundle.Artifact, revisions []prRevisionSummary, extraPatches []string) prBodyStats {
	stats := prBodyStats{RevisionCount: len(revisions), MaxRiskLevel: "low"}
	files := map[string]struct{}{}
	areas := map[string]struct{}{}
	riskSignals := map[string]struct{}{}
	for _, entry := range artifact.Stack {
		stats.AddedLines += addedLines(entry.Patch)
		stats.DeletedLines += deletedLines(entry.Patch)
		accumulateContext(&stats, entry.Change.ReviewContext, riskSignals)
	}
	for _, patch := range extraPatches {
		stats.AddedLines += addedLines(patch)
		stats.DeletedLines += deletedLines(patch)
	}
	if len(artifact.Stack) == 0 && artifact.Change != nil {
		accumulateContext(&stats, artifact.Change.ReviewContext, riskSignals)
	}
	for _, revision := range revisions {
		for _, file := range revision.Files {
			files[file] = struct{}{}
			areas[fileArea(file)] = struct{}{}
		}
		if scoreRank(revision.RiskLevel, revision.RiskScore) > scoreRank(stats.MaxRiskLevel, stats.MaxRiskScore) {
			stats.MaxRiskLevel = firstNonEmpty(revision.RiskLevel, stats.MaxRiskLevel)
			stats.MaxRiskScore = revision.RiskScore
		}
		for _, signal := range revision.RiskSignals {
			riskSignals[signal] = struct{}{}
		}
	}
	if stats.MaxRiskScore == 0 {
		switch {
		case stats.WarningCount > 0:
			stats.MaxRiskLevel, stats.MaxRiskScore = "high", 60
		case len(files) >= 4 || stats.AddedLines+stats.DeletedLines >= 100 || stats.ChangedSymbolCount >= 5 || len(files) >= 8 || stats.AddedLines+stats.DeletedLines >= 400:
			stats.MaxRiskLevel, stats.MaxRiskScore = "medium", 25
		}
	}
	stats.FileCount = len(files)
	stats.AreaCount = len(areas)
	stats.RiskSignals = sortedKeys(riskSignals)
	return stats
}

func accumulateContext(stats *prBodyStats, context *reviewbundle.ReviewContextPayload, riskSignals map[string]struct{}) {
	if context == nil {
		return
	}
	stats.ChangedSymbolCount += len(context.ChangedSymbols)
	for _, signal := range context.Risk.Signals {
		riskSignals[signal] = struct{}{}
	}
	for _, warning := range context.FeasibilityWarnings {
		if warning.Severity == "warning" {
			stats.WarningCount++
		} else {
			stats.InfoCount++
		}
	}
}

func filesFromPatch(patchText string) []string {
	var files []string
	for _, line := range strings.Split(patchText, "\n") {
		if strings.HasPrefix(line, "diff --git ") {
			if file := fileFromDiffLine(line); file != "" {
				files = append(files, file)
			}
		}
	}
	return sortedUnique(files)
}

func fileFromDiffLine(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return ""
	}
	file := strings.TrimPrefix(parts[3], "b/")
	if file == "/dev/null" {
		return ""
	}
	return strings.TrimSpace(file)
}

func addedLines(patchText string) int {
	count := 0
	for _, line := range strings.Split(patchText, "\n") {
		if strings.HasPrefix(line, "+++") || !strings.HasPrefix(line, "+") {
			continue
		}
		count++
	}
	return count
}

func deletedLines(patchText string) int {
	count := 0
	for _, line := range strings.Split(patchText, "\n") {
		if strings.HasPrefix(line, "---") || !strings.HasPrefix(line, "-") {
			continue
		}
		count++
	}
	return count
}

func fileArea(file string) string {
	dir := path.Dir(strings.TrimSpace(file))
	if dir == "." || dir == "/" || dir == "" {
		return "."
	}
	parts := strings.Split(dir, "/")
	if len(parts) > 2 {
		return strings.Join(parts[:2], "/")
	}
	return dir
}

func pullRequestURL(artifact reviewbundle.Artifact) string {
	if artifact.Push.GitHubPullRequestURL != nil && strings.TrimSpace(*artifact.Push.GitHubPullRequestURL) != "" {
		return strings.TrimSpace(*artifact.Push.GitHubPullRequestURL)
	}
	for _, entry := range artifact.Stack {
		if entry.GitHubPullRequestURL != nil && strings.TrimSpace(*entry.GitHubPullRequestURL) != "" {
			return strings.TrimSpace(*entry.GitHubPullRequestURL)
		}
	}
	return ""
}

func stackName(artifact reviewbundle.Artifact) string {
	if len(artifact.Stack) > 0 && strings.TrimSpace(artifact.Stack[0].BranchName) != "" {
		return strings.TrimSpace(artifact.Stack[0].BranchName)
	}
	return pointerString(artifact.Push.BranchName)
}

func headCommitSHA(artifact reviewbundle.Artifact, catalog prBodyCatalog) string {
	if n := len(catalog.Revisions); n > 0 {
		if sha := strings.TrimSpace(catalog.Revisions[n-1].CommitID); sha != "" {
			return sha
		}
	}
	return strings.TrimSpace(artifact.Push.HeadCommitID)
}

func fileInPRDiff(catalog prBodyCatalog, file string) bool {
	for _, hunk := range catalog.Hunks {
		if hunk.File == file {
			return true
		}
	}
	return false
}

func findingAttributions(finding codereview.Finding, summaryContext prSummaryContext, artifact reviewbundle.Artifact, catalog prBodyCatalog) []prAttribution {
	if len(finding.ResolvedSources) > 0 {
		var out []prAttribution
		for _, src := range finding.ResolvedSources {
			out = append(out, renderPRResolvedSource(src, artifact, catalog))
			if len(out) >= 3 {
				break
			}
		}
		return dedupeAttributions(out)
	}
	var out []prAttribution
	for _, sourceID := range finding.SourceIDs {
		if sourceID == "" {
			continue
		}
		if attr := attributionFromSourceID(sourceID, summaryContext.Snippets); attr != nil {
			out = append(out, *attr)
			continue
		}
		out = append(out, prAttribution{Kind: "review_resource", Label: "review resource", Ref: sourceID, Opaque: true})
	}
	for _, evidence := range finding.Evidence {
		if attr := attributionFromEvidence(evidence, summaryContext.Snippets); attr != nil {
			out = append(out, *attr)
		}
		if len(out) >= 3 {
			return dedupeAttributions(out)
		}
	}
	return dedupeAttributions(out)
}

func renderPRResolvedSource(src codereview.ResolvedSource, artifact reviewbundle.Artifact, catalog prBodyCatalog) prAttribution {
	label := codereview.ResolvedSourceLabel(src)
	kind := firstNonEmpty(strings.TrimSpace(src.Kind), "source")
	if src.Opaque {
		return prAttribution{Kind: kind, Label: label, URL: strings.TrimSpace(src.URL), Opaque: true}
	}
	if strings.EqualFold(kind, "session") {
		return prAttribution{Kind: kind, Label: label, Ref: strings.TrimSpace(src.Ref)}
	}
	file := strings.TrimSpace(src.File)
	if file == "" {
		return prAttribution{Kind: kind, Label: label, Ref: strings.TrimSpace(src.Ref), URL: strings.TrimSpace(src.URL)}
	}
	prURL := pullRequestURL(artifact)
	line := src.StartLine
	if fileInPRDiff(catalog, file) {
		if line > 0 && validatedHunkAnchor(catalog.Hunks, codereview.Finding{File: file, Line: line}) {
			return prAttribution{Kind: kind, Label: label, Ref: label, URL: githubHunkLineLink(prURL, file, line)}
		}
		for _, hunk := range catalog.Hunks {
			if hunk.File != file {
				continue
			}
			link := hunk.Link
			if link == "" {
				link = githubHunkLineLink(prURL, file, hunk.ChangedLine)
			}
			return prAttribution{Kind: kind, Label: label, Ref: label, URL: link}
		}
		return prAttribution{Kind: kind, Label: label, Ref: label}
	}
	if link := blobPermalink(prURL, headCommitSHA(artifact, catalog), file, line); link != "" {
		return prAttribution{Kind: kind, Label: label, Ref: label, URL: link}
	}
	return prAttribution{Kind: kind, Label: label, Ref: label}
}

func attributionFromSourceID(sourceID string, snippets []codereview.ContextSnippet) *prAttribution {
	for _, snippet := range snippets {
		if snippet.Ref == sourceID || snippet.SourceLabel == sourceID {
			attr := attributionFromSnippet(snippet)
			return &attr
		}
	}
	return nil
}

func attributionFromEvidence(evidence codereview.Evidence, snippets []codereview.ContextSnippet) *prAttribution {
	label := strings.TrimSpace(evidence.Label)
	value := strings.TrimSpace(evidence.Value)
	for _, snippet := range snippets {
		if snippet.Ref != "" && (snippet.Ref == label || snippet.Ref == value) {
			attr := attributionFromSnippet(snippet)
			return &attr
		}
		if snippet.SourceLabel != "" && snippet.SourceLabel == label {
			attr := attributionFromSnippet(snippet)
			return &attr
		}
	}
	if label == "" && value == "" {
		return nil
	}
	return &prAttribution{Kind: "codebase", Label: label, Ref: trimSentence(value, 120)}
}

func attributionFromSnippet(snippet codereview.ContextSnippet) prAttribution {
	kind := "codebase"
	switch {
	case strings.Contains(snippet.Kind, "session"):
		kind = "session"
	case strings.Contains(snippet.Source, "review") || strings.Contains(snippet.Kind, "review"):
		kind = "review_resource"
	case strings.Contains(snippet.Source, "turbopuffer"):
		kind = "previous_review"
	}
	return prAttribution{Kind: kind, Label: firstNonEmpty(snippet.SourceLabel, snippet.Kind), Ref: firstNonEmpty(snippet.Ref, snippet.File), URL: snippet.URL}
}

func renderAttributions(attributions []prAttribution, linkCtx attributionLinkContext) string {
	attributions = dedupeAttributions(attributions)
	if len(attributions) == 0 {
		return ""
	}
	if len(attributions) > 3 {
		attributions = attributions[:3]
	}
	parts := make([]string, 0, len(attributions))
	for _, attr := range attributions {
		parts = append(parts, renderAttributionPart(attr, linkCtx))
	}
	return strings.Join(parts, "; ")
}

func appendAuthorNotes(summary, notes string) string {
	notes = strings.TrimSpace(notes)
	if notes == "" {
		return strings.TrimSpace(summary)
	}
	return strings.TrimSpace(summary) + "\n\n" + githubPRAuthorNotesMarker + "\n" + githubPRAuthorNotesHeader + "\n\n" + notes
}

func splitGeneratedPRBody(existing string) (authorNotes string, gxOwned bool) {
	existing = strings.TrimSpace(existing)
	if existing == "" {
		return "", true
	}
	if strings.HasPrefix(existing, githubPRBodyMarker) {
		if idx := strings.Index(existing, githubPRAuthorNotesMarker); idx >= 0 {
			notes := strings.TrimSpace(existing[idx+len(githubPRAuthorNotesMarker):])
			notes = strings.TrimSpace(strings.TrimPrefix(notes, githubPRAuthorNotesHeader))
			return strings.TrimSpace(notes), true
		}
		return "", true
	}
	if strings.HasPrefix(existing, "Published by GX.") {
		return "", true
	}
	return existing, false
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func sortedUnique(values []string) []string {
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			seen[value] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

func sortedKeys(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func limitStrings(values []string, limit int) []string {
	if limit <= 0 || len(values) <= limit {
		return values
	}
	out := append([]string(nil), values[:limit]...)
	out = append(out, "and "+strconv.Itoa(len(values)-limit)+" more")
	return out
}

func lowerFirst(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func trimSentence(value string, limit int) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\n", " "))
	if limit <= 0 || len(value) <= limit {
		return value
	}
	candidate := value[:limit]
	if idx := strings.LastIndexAny(candidate, ".!?"); idx >= 80 {
		return strings.TrimSpace(candidate[:idx+1])
	}
	if idx := strings.LastIndex(candidate, ";"); idx >= 80 {
		return strings.TrimSpace(candidate[:idx]) + "."
	}
	return strings.TrimRight(candidate, " .,;:") + "."
}

func limitText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "\n[truncated]"
}

func evidenceText(evidence []codereview.Evidence) string {
	var parts []string
	for _, item := range evidence {
		parts = append(parts, item.Label, item.Value)
	}
	return strings.Join(parts, " ")
}

func revisionDescriptions(revisions []prRevisionSummary, limit int) []string {
	var descriptions []string
	for _, revision := range revisions {
		description := strings.TrimSpace(revision.Description)
		if description == "" {
			description = "(no description set)"
		}
		descriptions = append(descriptions, description)
		if limit > 0 && len(descriptions) >= limit {
			break
		}
	}
	if limit > 0 && len(revisions) > limit {
		descriptions = append(descriptions, fmt.Sprintf("%d more", len(revisions)-limit))
	}
	return descriptions
}

func humanAreaList(areas []string, limit int) string {
	areas = sortedUnique(areas)
	if len(areas) == 0 {
		return ""
	}
	areas = limitStrings(areas, limit)
	for i := range areas {
		if !strings.HasPrefix(areas[i], "and ") {
			areas[i] = "`" + areas[i] + "`"
		}
	}
	return strings.Join(areas, ", ")
}

func humanizeSignal(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "_", " "))
	value = strings.TrimPrefix(value, "warning:")
	if value == "" {
		return "review signal"
	}
	return value
}

func countTestFiles(files []string) int {
	count := 0
	for _, file := range files {
		lower := strings.ToLower(file)
		if strings.HasSuffix(lower, "_test.go") || strings.Contains(lower, "test") {
			count++
		}
	}
	return count
}

func dedupeNotableChanges(items []prNotableChange) []prNotableChange {
	seen := map[string]int{}
	var out []prNotableChange
	for _, item := range items {
		key := notableChangeDedupeKey(item)
		if idx, ok := seen[key]; ok {
			if item.Score > out[idx].Score {
				out[idx] = item
			}
			continue
		}
		seen[key] = len(out)
		out = append(out, item)
	}
	return out
}

func notableChangeDedupeKey(item prNotableChange) string {
	if title := strings.ToLower(strings.TrimSpace(item.Title)); title != "" {
		return "title:" + title
	}
	if item.Link != "" {
		return "link:" + item.Link
	}
	return ""
}

func dedupeAttributions(attributions []prAttribution) []prAttribution {
	seen := map[string]struct{}{}
	var out []prAttribution
	for _, attr := range attributions {
		key := attr.Kind + "\x00" + attr.Label + "\x00" + attr.Ref + "\x00" + attr.URL
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, attr)
	}
	return out
}

func dedupePRContextSnippets(snippets []codereview.ContextSnippet) []codereview.ContextSnippet {
	seen := map[string]struct{}{}
	var out []codereview.ContextSnippet
	for _, snippet := range snippets {
		key := snippet.Kind + "\x00" + snippet.Ref + "\x00" + snippet.Source
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, snippet)
	}
	return out
}

func stringAttribute(attributes map[string]any, key string) string {
	if attributes == nil {
		return ""
	}
	return stringAny(attributes[key])
}

func stringAny(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case fmt.Stringer:
		return strings.TrimSpace(typed.String())
	default:
		return ""
	}
}

func atoiDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func scoreRank(level string, score int) int {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "high":
		return 300 + score
	case "medium":
		return 200 + score
	case "low":
		return 100 + score
	default:
		return score
	}
}

func containsAny(text string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}
