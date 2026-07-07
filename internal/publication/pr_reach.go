package publication

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

const (
	lexicalReachTimeout     = 10 * time.Second
	maxLexicalReachSymbols  = 20
	maxRefsPerSymbol        = 10
	maxSymbolMatchCount     = 200
	highReachRefsPerSymbol  = 25
	maxBlastRadiusCritical  = 4
	maxBlastRadiusSymbols   = 4
	maxBlastRadiusNarrative = 500
)

var goPatchSymbolRE = regexp.MustCompile(`^\+(?:func|type)\s+([A-Za-z_][A-Za-z0-9_]*)`)

type reachSite struct {
	File string
	Line int
}

type symbolReach struct {
	Symbol     string
	References []reachSite
}

type lexicalReach struct {
	Symbols        []symbolReach
	ReferenceCount int
	DependentFiles int
	Truncated      bool
}

func lexicalReachEnabled() bool {
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("GX_PR_BLAST_RADIUS")), "0")
}

func computeLexicalReach(ctx context.Context, repoRoot string, catalog prBodyCatalog, artifact reviewbundle.Artifact) lexicalReach {
	if !lexicalReachEnabled() {
		return lexicalReach{}
	}
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return lexicalReach{}
	}
	triage := triageChangeFromCatalog(catalog)
	switch triage.Class {
	case "docs-only", "config-only", "mechanical":
		return lexicalReach{}
	}

	ctx, cancel := context.WithTimeout(ctx, lexicalReachTimeout)
	defer cancel()

	changedFiles := map[string]struct{}{}
	for _, file := range catalog.Files {
		changedFiles[file] = struct{}{}
	}

	symbols := gatherReachSymbols(artifact, catalog)
	if len(symbols) == 0 {
		return lexicalReach{}
	}
	if len(symbols) > maxLexicalReachSymbols {
		symbols = symbols[:maxLexicalReachSymbols]
	}

	var reach lexicalReach
	dependentFiles := map[string]struct{}{}
	for _, symbol := range symbols {
		if len(symbol) < 4 {
			continue
		}
		sites, totalMatches, err := grepSymbolReferences(ctx, repoRoot, symbol, changedFiles)
		if err != nil {
			return lexicalReach{}
		}
		if totalMatches >= maxSymbolMatchCount {
			reach.Truncated = true
			continue
		}
		if len(sites) == 0 {
			continue
		}
		reach.Symbols = append(reach.Symbols, symbolReach{Symbol: symbol, References: sites})
		reach.ReferenceCount += len(sites)
		for _, site := range sites {
			dependentFiles[site.File] = struct{}{}
		}
	}
	reach.DependentFiles = len(dependentFiles)
	sort.Slice(reach.Symbols, func(i, j int) bool {
		return len(reach.Symbols[i].References) > len(reach.Symbols[j].References)
	})
	return reach
}

func gatherReachSymbols(artifact reviewbundle.Artifact, catalog prBodyCatalog) []string {
	seen := map[string]struct{}{}
	var symbols []string
	add := func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}
		if _, ok := seen[name]; ok {
			return
		}
		seen[name] = struct{}{}
		symbols = append(symbols, name)
	}
	for _, entry := range artifact.Stack {
		if entry.Change.ReviewContext == nil {
			continue
		}
		for _, sym := range entry.Change.ReviewContext.ChangedSymbols {
			add(sym.Symbol)
		}
	}
	if artifact.Change != nil && artifact.Change.ReviewContext != nil {
		for _, sym := range artifact.Change.ReviewContext.ChangedSymbols {
			add(sym.Symbol)
		}
	}
	if len(symbols) > 0 {
		return symbols
	}
	for _, hunk := range catalog.Hunks {
		for _, line := range strings.Split(hunk.Patch, "\n") {
			if match := goPatchSymbolRE.FindStringSubmatch(line); len(match) == 2 {
				add(match[1])
			}
		}
	}
	return symbols
}

func grepSymbolReferences(ctx context.Context, repoRoot, symbol string, changedFiles map[string]struct{}) ([]reachSite, int, error) {
	cmd := exec.CommandContext(ctx, "git", "grep", "-nw", "--fixed-strings", symbol)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	var lines []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	totalMatches := len(lines)
	if totalMatches >= maxSymbolMatchCount {
		return nil, totalMatches, nil
	}
	var sites []reachSite
	for _, line := range lines {
		file, refLine, ok := parseGitGrepLine(line)
		if !ok {
			continue
		}
		if _, changed := changedFiles[file]; changed {
			continue
		}
		if isReachTestFile(file) {
			continue
		}
		sites = append(sites, reachSite{File: file, Line: refLine})
		if len(sites) >= maxRefsPerSymbol {
			break
		}
	}
	return sites, totalMatches, nil
}

func parseGitGrepLine(line string) (string, int, bool) {
	parts := strings.SplitN(line, ":", 3)
	if len(parts) < 3 {
		return "", 0, false
	}
	lineNum, err := strconv.Atoi(parts[1])
	if err != nil || lineNum <= 0 {
		return "", 0, false
	}
	return parts[0], lineNum, true
}

func isReachTestFile(file string) bool {
	lower := strings.ToLower(strings.TrimSpace(file))
	if strings.HasSuffix(lower, "_test.go") {
		return true
	}
	if strings.HasPrefix(lower, "test/") || strings.Contains(lower, "/test/") {
		return true
	}
	if strings.HasPrefix(lower, "testdata/") || strings.Contains(lower, "/testdata/") {
		return true
	}
	return false
}

func applyLexicalReachToStats(stats *prBodyStats, reach lexicalReach) {
	if stats == nil || reach.ReferenceCount == 0 {
		return
	}
	signal := "lexical_reach:" + strconv.Itoa(reach.ReferenceCount) + "_refs_in_" + strconv.Itoa(reach.DependentFiles) + "_files"
	stats.RiskSignals = append(stats.RiskSignals, signal)
	stats.RiskSignals = sortedUnique(stats.RiskSignals)

	level := stats.MaxRiskLevel
	score := stats.MaxRiskScore
	if reach.DependentFiles >= 3 && scoreRank("medium", 25) > scoreRank(level, score) {
		level, score = "medium", 25
	}
	for _, sym := range reach.Symbols {
		if len(sym.References) >= highReachRefsPerSymbol && scoreRank("high", 60) > scoreRank(level, score) {
			level, score = "high", 60
		}
	}
	if reach.DependentFiles >= 10 && scoreRank("high", 60) > scoreRank(level, score) {
		level, score = "high", 60
	}
	if scoreRank(level, score) > scoreRank(stats.MaxRiskLevel, stats.MaxRiskScore) {
		stats.MaxRiskLevel = level
		stats.MaxRiskScore = score
	}
}

func lexicalReachContextSnippets(reach lexicalReach) []codereview.ContextSnippet {
	if reach.ReferenceCount == 0 {
		return nil
	}
	var out []codereview.ContextSnippet
	for i, sym := range reach.Symbols {
		if i >= 5 {
			break
		}
		files := map[string]struct{}{}
		for _, ref := range sym.References {
			files[ref.File] = struct{}{}
		}
		text := "Lexical reach for " + sym.Symbol + ": " + strconv.Itoa(len(sym.References)) + " reference(s) across " + strconv.Itoa(len(files)) + " other file(s)."
		for j, ref := range sym.References {
			if j >= 3 {
				break
			}
			text += " " + ref.File + ":" + strconv.Itoa(ref.Line)
		}
		out = append(out, codereview.ContextSnippet{
			Kind: "lexical_reach",
			Ref:  sym.Symbol,
			Text: text,
		})
	}
	return out
}

func topReachSymbolNames(reach lexicalReach, limit int) []string {
	if limit <= 0 {
		return nil
	}
	names := make([]string, 0, limit)
	for _, sym := range reach.Symbols {
		if sym.Symbol == "" {
			continue
		}
		names = append(names, sym.Symbol)
		if len(names) >= limit {
			break
		}
	}
	return names
}

func shouldRenderBlastRadiusSection(triage codereview.ChangeTriage, catalog prBodyCatalog, reach lexicalReach, policy codereview.ReviewPolicy) bool {
	switch triage.Class {
	case "code", "security-sensitive":
		return true
	}
	if reach.ReferenceCount > 0 {
		return true
	}
	for _, hunk := range catalog.Hunks {
		if hunkMatchesRiskPath(hunk.File, policy) {
			return true
		}
	}
	return false
}

func blastRadiusNarrative(catalog prBodyCatalog, reach lexicalReach, policy codereview.ReviewPolicy, triage codereview.ChangeTriage, summary codereview.PRSummaryReview, aiSucceeded bool) string {
	if aiSucceeded {
		if text := sanitizedBlastRadiusNarrative(summary.DownstreamImpact); text != "" {
			return text
		}
	}
	return heuristicBlastRadiusNarrative(catalog, reach, policy, triage)
}

func sanitizedBlastRadiusNarrative(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = prURLStripRE.ReplaceAllString(text, "")
	text = sanitizePRVisibleText(text)
	return trimSentence(text, maxBlastRadiusNarrative)
}

func heuristicBlastRadiusNarrative(catalog prBodyCatalog, reach lexicalReach, policy codereview.ReviewPolicy, triage codereview.ChangeTriage) string {
	switch triage.Class {
	case "docs-only":
		return "Documentation-only change; no customer-facing runtime behavior should change."
	case "tests-only":
		return "Tests-only change; production behavior should stay the same unless the tests uncovered a needed fix."
	case "config-only":
		return "Configuration-only change; verify deployed settings, but customer-facing code paths are unlikely to shift."
	case "mechanical":
		return "Mechanical refactor with low behavioral risk; existing behavior should be preserved, though wide refactors can still hide subtle regressions."
	case "security-sensitive":
		area := firstNonEmpty(firstSorted(catalog.Areas), verdictReasonFileArea(catalog.Files))
		if area != "" {
			return fmt.Sprintf("Security-sensitive change in %s; could affect customer trust, data handling, or access control — verify auth and data paths before merge.", humanAreaList([]string{area}, 1))
		}
		return "Security-sensitive change; could affect customer trust, data handling, or access control — verify auth and data paths before merge."
	}

	totalLines := catalog.Stats.AddedLines + catalog.Stats.DeletedLines
	areas := humanAreaList(catalog.Areas, 3)
	reachNote := ""
	if reach.ReferenceCount > 0 {
		reachNote = fmt.Sprintf(" Changed symbols are referenced across %d other file(s), so downstream callers may need re-checking.", reach.DependentFiles)
	}
	if catalogTouchesCriticalPath(catalog, policy) {
		return fmt.Sprintf("Touches configured high-risk paths; could affect customer-visible behavior or billing — review end-to-end before merge.%s", reachNote)
	}

	switch catalog.Stats.MaxRiskLevel {
	case "high":
		if areas != "" {
			return fmt.Sprintf("Large high-risk change across %d file(s) in %s; could alter customer-visible behavior and ripple through downstream callers.%s Review end-to-end flows before merge.", catalog.Stats.FileCount, areas, reachNote)
		}
		return fmt.Sprintf("Large high-risk change across %d file(s); could alter customer-visible behavior and ripple through downstream callers.%s Review end-to-end flows before merge.", catalog.Stats.FileCount, reachNote)
	case "medium":
		if areas != "" {
			return fmt.Sprintf("Moderate change across %d file(s) in %s; may shift runtime behavior — sanity-check key user flows.%s", catalog.Stats.FileCount, areas, reachNote)
		}
		return fmt.Sprintf("Moderate change across %d file(s); may shift runtime behavior — sanity-check key user flows.%s", catalog.Stats.FileCount, reachNote)
	}

	switch {
	case totalLines >= 400 || catalog.Stats.FileCount >= 8:
		if areas != "" {
			return fmt.Sprintf("Broad change across %d file(s) in %s; could alter customer-visible behavior with downstream effects.%s", catalog.Stats.FileCount, areas, reachNote)
		}
		return fmt.Sprintf("Broad change across %d file(s); could alter customer-visible behavior with downstream effects.%s", catalog.Stats.FileCount, reachNote)
	case totalLines >= 100 || catalog.Stats.FileCount >= 4 || reach.ReferenceCount >= 3:
		if areas != "" {
			return fmt.Sprintf("Touches %d file(s) in %s; localized but may affect adjacent behavior.%s", catalog.Stats.FileCount, areas, reachNote)
		}
		return fmt.Sprintf("Touches %d file(s); localized but may affect adjacent behavior.%s", catalog.Stats.FileCount, reachNote)
	default:
		base := "Small localized change with low risk to existing customer-visible behavior."
		if reachNote != "" {
			return base + reachNote
		}
		return base
	}
}

func catalogTouchesCriticalPath(catalog prBodyCatalog, policy codereview.ReviewPolicy) bool {
	for _, hunk := range catalog.Hunks {
		if hunkMatchesRiskPath(hunk.File, policy) {
			return true
		}
	}
	return false
}

func renderBlastRadiusSection(artifact reviewbundle.Artifact, catalog prBodyCatalog, reach lexicalReach, policy codereview.ReviewPolicy, triage codereview.ChangeTriage, summary codereview.PRSummaryReview, aiSucceeded bool) string {
	if !shouldRenderBlastRadiusSection(triage, catalog, reach, policy) {
		return ""
	}
	prURL := pullRequestURL(artifact)
	sha := headCommitSHA(artifact, catalog)
	var body strings.Builder
	body.WriteString("\n\n## Blast Radius\n\n")
	body.WriteString(readinessSentence(catalog, reach))
	body.WriteByte('\n')
	if narrative := blastRadiusNarrative(catalog, reach, policy, triage, summary, aiSucceeded); narrative != "" {
		body.WriteString(narrative)
		body.WriteByte('\n')
	}

	criticalCount := 0
	for _, riskPath := range policy.RiskPaths {
		if criticalCount >= maxBlastRadiusCritical {
			break
		}
		hunk, ok := bestHunkMatchingRiskPath(catalog, riskPath.Glob)
		if !ok {
			continue
		}
		body.WriteString(renderCriticalPathLine(riskPath, hunk, prURL))
		criticalCount++
	}

	reachLineCount := 0
	for _, sym := range reach.Symbols {
		if reachLineCount >= maxBlastRadiusSymbols {
			break
		}
		body.WriteString(renderReachSymbolLine(sym, catalog, prURL, sha))
		reachLineCount++
	}

	if reachLineCount > 0 {
		body.WriteString("\nReach is lexical (text search), not a dependency graph.")
		if reach.Truncated {
			body.WriteString(" Common names were skipped.")
		}
		body.WriteByte('\n')
	} else if reach.Truncated {
		body.WriteString("\nReach is lexical (text search) and truncated; common names were skipped.\n")
	}
	return body.String()
}

func bestHunkMatchingRiskPath(catalog prBodyCatalog, glob string) (prHunkSummary, bool) {
	var best prHunkSummary
	bestScore := -1
	found := false
	for _, hunk := range catalog.Hunks {
		if !codereview.MatchRiskPathGlob(glob, hunk.File) {
			continue
		}
		score := hunk.NewLines + hunk.OldLines
		if !found || score > bestScore {
			best = hunk
			bestScore = score
			found = true
		}
	}
	return best, found
}

func renderCriticalPathLine(riskPath codereview.RiskPath, hunk prHunkSummary, prURL string) string {
	var b strings.Builder
	b.WriteString("- **Critical path `")
	b.WriteString(riskPath.Glob)
	b.WriteString("`** — ")
	b.WriteString(strings.TrimSpace(riskPath.Message))
	b.WriteString(" (REVIEW.md)")
	link := hunk.Link
	if link == "" && prURL != "" {
		link = githubHunkLineLink(prURL, hunk.File, hunk.ChangedLine)
	}
	if link != "" {
		b.WriteString(" — ")
		b.WriteString(formatPRSummaryLink("changed here", link))
	}
	b.WriteByte('\n')
	return b.String()
}

func findDefiningHunk(catalog prBodyCatalog, symbol string) (prHunkSummary, int) {
	symRE := regexp.MustCompile(`^\+(?:func|type)\s+` + regexp.QuoteMeta(symbol) + `\b`)
	for _, hunk := range catalog.Hunks {
		if line := lineOfPatchMatch(hunk, symRE); line > 0 {
			return hunk, line
		}
	}
	return prHunkSummary{}, 0
}

func renderReachSymbolLine(sym symbolReach, catalog prBodyCatalog, prURL, sha string) string {
	uniqueFiles := map[string]struct{}{}
	for _, ref := range sym.References {
		uniqueFiles[ref.File] = struct{}{}
	}

	var b strings.Builder
	b.WriteString("- **`")
	b.WriteString(sym.Symbol)
	b.WriteString("`** — ")

	if defHunk, defLine := findDefiningHunk(catalog, sym.Symbol); defLine > 0 {
		defLink := defHunk.Link
		if defLine != defHunk.ChangedLine || defLink == "" {
			defLink = githubHunkLineLink(prURL, defHunk.File, defLine)
		}
		if defLink != "" {
			b.WriteString("redefined in ")
			b.WriteString(formatPRSummaryLink(path.Base(defHunk.File)+":"+strconv.Itoa(defLine), defLink))
			b.WriteString(", ")
		}
	}

	b.WriteString("referenced ")
	b.WriteString(strconv.Itoa(len(sym.References)))
	b.WriteString("× across ")
	b.WriteString(strconv.Itoa(len(uniqueFiles)))
	b.WriteString(" file")
	if len(uniqueFiles) != 1 {
		b.WriteByte('s')
	}
	b.WriteString(": ")

	linkCount := 0
	for _, ref := range sym.References {
		if linkCount >= 2 {
			break
		}
		if link := blobPermalink(prURL, sha, ref.File, ref.Line); link != "" {
			if linkCount > 0 {
				b.WriteString(", ")
			}
			b.WriteString(formatPRSummaryLink(ref.File+":"+strconv.Itoa(ref.Line), link))
			linkCount++
		}
	}
	b.WriteByte('\n')
	return b.String()
}
