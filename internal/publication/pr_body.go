package publication

import (
	"fmt"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

const maxBodyListItems = 6

type prRevisionSummary struct {
	Description string
	CommitID    string
	BranchName  string
	Files       []string
	Symbols     []string
	RiskLevel   string
	RiskScore   int
	RiskSignals []string
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
	ProvenanceStatuses []string
	LinkedSessionCount int
	Attribution        []string
	StructuralStatus   string
}

func GitHubPullRequestBodyFromArtifact(artifact reviewbundle.Artifact) string {
	revisions := revisionSummaries(artifact)
	stats := bodyStats(artifact, revisions)
	var body strings.Builder
	body.WriteString("Published by GX.\n\n")
	body.WriteString("## Summary\n")
	body.WriteString(summaryText(artifact, stats))
	body.WriteString("\n\n## Review signals\n")
	for _, line := range reviewSignalLines(stats) {
		body.WriteString("- ")
		body.WriteString(line)
		body.WriteByte('\n')
	}
	body.WriteString("\n## Important changes\n")
	if len(revisions) == 0 {
		body.WriteString("- (none)\n")
	} else {
		for _, revision := range revisions {
			body.WriteString("- ")
			body.WriteString(revisionLine(artifact, revision))
			body.WriteByte('\n')
		}
	}
	body.WriteString("\n## Links\n")
	for _, line := range linkLines(artifact) {
		body.WriteString("- ")
		body.WriteString(line)
		body.WriteByte('\n')
	}
	return strings.TrimRight(body.String(), "\n")
}

func summaryText(artifact reviewbundle.Artifact, stats prBodyStats) string {
	name := stackName(artifact)
	if stats.RevisionCount == 0 {
		if name != "" {
			return fmt.Sprintf("Publishes `%s` with no recorded revisions.", name)
		}
		return "Publishes a GX change with no recorded revisions."
	}
	if name == "" {
		return fmt.Sprintf("Publishes %d GX revision(s) with %s blast radius.", stats.RevisionCount, stats.MaxRiskLevel)
	}
	return fmt.Sprintf("Publishes %d GX revision(s) for `%s` with %s blast radius.", stats.RevisionCount, name, stats.MaxRiskLevel)
}

func reviewSignalLines(stats prBodyStats) []string {
	blast := fmt.Sprintf("Blast radius: %s (%d revision(s), %d file(s), %d area(s), +%d/-%d lines, %d changed symbol(s))",
		stats.MaxRiskLevel,
		stats.RevisionCount,
		stats.FileCount,
		stats.AreaCount,
		stats.AddedLines,
		stats.DeletedLines,
		stats.ChangedSymbolCount,
	)
	riskSignals := "none"
	if len(stats.RiskSignals) > 0 {
		riskSignals = strings.Join(limitStrings(humanSignals(stats.RiskSignals), maxBodyListItems), ", ")
	}
	risk := fmt.Sprintf("Risk: %s (%d/100; %s)", stats.MaxRiskLevel, stats.MaxRiskScore, riskSignals)
	provenance := "Provenance: absent"
	if len(stats.ProvenanceStatuses) > 0 {
		provenance = fmt.Sprintf("Provenance: %s", strings.Join(stats.ProvenanceStatuses, ", "))
	}
	if stats.LinkedSessionCount > 0 {
		provenance += fmt.Sprintf(" (%d linked session(s))", stats.LinkedSessionCount)
	}
	if len(stats.Attribution) > 0 {
		provenance += "; attribution: " + strings.Join(limitStrings(stats.Attribution, maxBodyListItems), ", ")
	}
	structural := fmt.Sprintf("Structural context: %s (%d changed symbol(s))", firstNonEmpty(stats.StructuralStatus, "unavailable"), stats.ChangedSymbolCount)
	warnings := fmt.Sprintf("Feasibility: %d warning(s), %d info item(s)", stats.WarningCount, stats.InfoCount)
	return []string{blast, risk, provenance, structural, warnings}
}

func revisionLine(artifact reviewbundle.Artifact, revision prRevisionSummary) string {
	description := strings.TrimSpace(revision.Description)
	if description == "" {
		description = "(no description set)"
	}
	linkedDescription := description
	if link := commitLink(artifact, revision.CommitID); link != "" {
		linkedDescription = fmt.Sprintf("[%s](%s)", description, link)
	}
	parts := []string{linkedDescription}
	if len(revision.Files) > 0 {
		parts = append(parts, "files: "+linkedFileList(artifact, revision.CommitID, revision.Files))
	}
	if len(revision.Symbols) > 0 {
		parts = append(parts, "symbols: "+strings.Join(limitCodeStrings(revision.Symbols, maxBodyListItems), ", "))
	}
	if revision.RiskLevel != "" {
		parts = append(parts, fmt.Sprintf("risk: %s (%d/100)", revision.RiskLevel, revision.RiskScore))
	}
	return strings.Join(parts, " | ")
}

func linkLines(artifact reviewbundle.Artifact) []string {
	var lines []string
	if prURL := pullRequestURL(artifact); prURL != "" {
		lines = append(lines, fmt.Sprintf("[GitHub PR](%s)", prURL))
		lines = append(lines, fmt.Sprintf("[Files changed](%s/files)", strings.TrimRight(prURL, "/")))
	}
	if strings.TrimSpace(artifact.ReviewURL) != "" {
		lines = append(lines, fmt.Sprintf("[GX review](%s)", strings.TrimSpace(artifact.ReviewURL)))
	}
	if len(lines) == 0 {
		lines = append(lines, "(none)")
	}
	return lines
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
	summary := prRevisionSummary{
		Description: change.Description,
		CommitID:    change.CurrentCommitID,
		BranchName:  branchName,
		Files:       sortedUnique(files),
	}
	if change.ReviewContext != nil {
		summary.RiskLevel = change.ReviewContext.Risk.Level
		summary.RiskScore = change.ReviewContext.Risk.Score
		summary.RiskSignals = append([]string(nil), change.ReviewContext.Risk.Signals...)
		for _, symbol := range change.ReviewContext.ChangedSymbols {
			label := symbol.Symbol
			if symbol.File != "" {
				label = symbol.File + ":" + label
			}
			summary.Symbols = append(summary.Symbols, strings.TrimSpace(label))
		}
		summary.Symbols = sortedUnique(summary.Symbols)
	}
	return summary
}

func bodyStats(artifact reviewbundle.Artifact, revisions []prRevisionSummary) prBodyStats {
	stats := prBodyStats{
		RevisionCount:    len(revisions),
		MaxRiskLevel:     "low",
		StructuralStatus: "unavailable",
	}
	files := map[string]struct{}{}
	areas := map[string]struct{}{}
	riskSignals := map[string]struct{}{}
	provenanceStatuses := map[string]struct{}{}
	linkedSessions := map[string]struct{}{}
	attribution := map[string]struct{}{}
	for _, entry := range artifact.Stack {
		stats.AddedLines += addedLines(entry.Patch)
		stats.DeletedLines += deletedLines(entry.Patch)
		accumulateContext(&stats, entry.Change.ReviewContext, riskSignals, provenanceStatuses, linkedSessions, attribution)
	}
	if len(artifact.Stack) == 0 && artifact.Change != nil {
		accumulateContext(&stats, artifact.Change.ReviewContext, riskSignals, provenanceStatuses, linkedSessions, attribution)
	}
	for _, revision := range revisions {
		for _, file := range revision.Files {
			files[file] = struct{}{}
			areas[fileArea(file)] = struct{}{}
		}
		stats.ChangedSymbolCount += len(revision.Symbols)
		if scoreRank(revision.RiskLevel, revision.RiskScore) > scoreRank(stats.MaxRiskLevel, stats.MaxRiskScore) {
			stats.MaxRiskLevel = firstNonEmpty(revision.RiskLevel, stats.MaxRiskLevel)
			stats.MaxRiskScore = revision.RiskScore
		}
		for _, signal := range revision.RiskSignals {
			riskSignals[signal] = struct{}{}
		}
	}
	if len(artifact.Stack) == 0 {
		if artifact.Change != nil {
			for _, file := range artifact.Change.Files {
				files[file] = struct{}{}
				areas[fileArea(file)] = struct{}{}
			}
		}
	}
	if stats.MaxRiskScore == 0 {
		switch {
		case stats.WarningCount > 0 || len(files) >= 8 || stats.AddedLines+stats.DeletedLines >= 400:
			stats.MaxRiskLevel = "high"
			stats.MaxRiskScore = 60
		case len(files) >= 4 || stats.AddedLines+stats.DeletedLines >= 100 || stats.ChangedSymbolCount >= 5:
			stats.MaxRiskLevel = "medium"
			stats.MaxRiskScore = 25
		default:
			stats.MaxRiskLevel = "low"
		}
	}
	stats.FileCount = len(files)
	stats.AreaCount = len(areas)
	stats.RiskSignals = sortedKeys(riskSignals)
	stats.ProvenanceStatuses = sortedKeys(provenanceStatuses)
	stats.LinkedSessionCount = len(linkedSessions)
	stats.Attribution = sortedKeys(attribution)
	return stats
}

func accumulateContext(stats *prBodyStats, context *reviewbundle.ReviewContextPayload, riskSignals, provenanceStatuses, linkedSessions, attribution map[string]struct{}) {
	if context == nil {
		return
	}
	if context.StructuralStatus == "available" {
		stats.StructuralStatus = "available"
	}
	if context.ProvenanceStatus != "" {
		provenanceStatuses[context.ProvenanceStatus] = struct{}{}
	}
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
	for _, source := range context.ProvenanceSources {
		if source.SessionID != "" {
			linkedSessions[source.SessionID] = struct{}{}
		}
	}
	for _, source := range context.TranscriptSources {
		if source.SessionID != "" {
			linkedSessions[source.SessionID] = struct{}{}
		}
	}
	for _, item := range context.AgentProvenance {
		if item.SessionID != "" {
			linkedSessions[item.SessionID] = struct{}{}
		}
		label := agentLabel(item)
		if label != "" {
			attribution[label] = struct{}{}
		}
	}
}

func agentLabel(item reviewbundle.ReviewAgentProvenance) string {
	parts := []string{}
	for _, value := range []string{item.AgentTool, item.Provider, item.ModelID} {
		value = strings.TrimSpace(value)
		if value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, "/")
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

func filesFromPatch(patchText string) []string {
	var files []string
	for _, line := range strings.Split(patchText, "\n") {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		file := strings.TrimPrefix(parts[3], "b/")
		if file != "" && file != "/dev/null" {
			files = append(files, file)
		}
	}
	return sortedUnique(files)
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
	file = strings.TrimSpace(file)
	if file == "" {
		return "."
	}
	dir := path.Dir(file)
	if dir == "." || dir == "/" {
		return "."
	}
	parts := strings.Split(dir, "/")
	if len(parts) > 2 {
		return strings.Join(parts[:2], "/")
	}
	return dir
}

func linkedFileList(artifact reviewbundle.Artifact, commitID string, files []string) string {
	files = limitStrings(sortedUnique(files), maxBodyListItems)
	out := make([]string, 0, len(files))
	for _, file := range files {
		label := "`" + file + "`"
		if link := fileLink(artifact, commitID, file); link != "" {
			label = fmt.Sprintf("[`%s`](%s)", file, link)
		}
		out = append(out, label)
	}
	return strings.Join(out, ", ")
}

func commitLink(artifact reviewbundle.Artifact, commitID string) string {
	commitID = strings.TrimSpace(commitID)
	if commitID == "" {
		return ""
	}
	if base := repositoryWebURL(artifact); base != "" {
		return base + "/commit/" + url.PathEscape(commitID)
	}
	return ""
}

func fileLink(artifact reviewbundle.Artifact, commitID, file string) string {
	commitID = strings.TrimSpace(commitID)
	file = strings.TrimSpace(file)
	if commitID == "" || file == "" {
		return ""
	}
	if base := repositoryWebURL(artifact); base != "" {
		return base + "/blob/" + url.PathEscape(commitID) + "/" + escapePath(file)
	}
	return ""
}

func repositoryWebURL(artifact reviewbundle.Artifact) string {
	prURL := pullRequestURL(artifact)
	if prURL == "" {
		return ""
	}
	parsed, err := url.Parse(prURL)
	if err != nil {
		return ""
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "pull" {
		return ""
	}
	parsed.Path = "/" + escapePath(parts[0]) + "/" + escapePath(parts[1])
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return strings.TrimRight(parsed.String(), "/")
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
	if artifact.Push.BranchName != nil && strings.TrimSpace(*artifact.Push.BranchName) != "" {
		return strings.TrimSpace(*artifact.Push.BranchName)
	}
	if len(artifact.Stack) > 0 && strings.TrimSpace(artifact.Stack[0].BranchName) != "" {
		return strings.TrimSpace(artifact.Stack[0].BranchName)
	}
	return ""
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
		if value == "" {
			continue
		}
		seen[value] = struct{}{}
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

func limitCodeStrings(values []string, limit int) []string {
	values = limitStrings(values, limit)
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.HasPrefix(value, "and ") {
			out = append(out, value)
			continue
		}
		out = append(out, "`"+value+"`")
	}
	return out
}

func humanSignals(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, strings.ReplaceAll(value, "_", " "))
	}
	return out
}

func escapePath(value string) string {
	parts := strings.Split(value, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}
