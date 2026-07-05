package publication

import (
	"context"
	"os"
	"os/exec"
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

func renderBlastRadiusSection(artifact reviewbundle.Artifact, catalog prBodyCatalog, reach lexicalReach) string {
	if reach.ReferenceCount <= 0 {
		return ""
	}
	prURL := pullRequestURL(artifact)
	sha := headCommitSHA(artifact, catalog)
	var body strings.Builder
	body.WriteString("\n\n## Blast Radius\n\n")
	lines := 0
	for _, sym := range reach.Symbols {
		if lines >= 5 {
			break
		}
		body.WriteString("- **")
		body.WriteString(sym.Symbol)
		body.WriteString("**: ")
		body.WriteString(strconv.Itoa(len(sym.References)))
		body.WriteString(" reference")
		if len(sym.References) != 1 {
			body.WriteByte('s')
		}
		linkCount := 0
		for _, ref := range sym.References {
			if linkCount >= 2 {
				break
			}
			if link := blobPermalink(prURL, sha, ref.File, ref.Line); link != "" {
				if linkCount == 0 {
					body.WriteString(" — ")
				} else {
					body.WriteString(", ")
				}
				body.WriteString("[")
				body.WriteString(ref.File)
				body.WriteString(":")
				body.WriteString(strconv.Itoa(ref.Line))
				body.WriteString("](")
				body.WriteString(link)
				body.WriteString(")")
				linkCount++
			}
		}
		body.WriteByte('\n')
		lines++
	}
	if reach.Truncated {
		body.WriteString("\nReach is lexical (text search) and truncated; common names were skipped.\n")
	}
	return body.String()
}
