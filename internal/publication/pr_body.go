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
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
	"github.com/satoricorp/gx/internal/semantic"
)

const (
	githubPRBodyMarker      = "<!-- gx:pr-summary:v1 -->"
	maxPRBodyItems          = 4
	maxPRContextSnippetSize = 1800
)

var hunkHeaderRE = regexp.MustCompile(`^@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)

var prSummaryReviewerFromEnv = codereview.ReviewerFromEnv
var collectPRSummaryContext = collectDefaultPRSummaryContext

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

type prNeedsReviewItem struct {
	Title  string
	Detail string
	Link   string
}

func GitHubPullRequestBodyFromArtifact(ctx context.Context, artifact reviewbundle.Artifact) (string, error) {
	catalog := buildPRBodyCatalog(artifact)
	summaryContext, err := collectPRSummaryContext(ctx, artifact, catalog)
	if err != nil {
		return "", err
	}
	findings := reviewPRSummaryFindings(ctx, artifact, catalog, summaryContext)
	return renderGitHubPullRequestBody(artifact, catalog, summaryContext, findings), nil
}

func renderGitHubPullRequestBody(artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext, findings []codereview.Finding) string {
	items := needsReviewItems(artifact, catalog, findings)
	var body strings.Builder
	body.WriteString(githubPRBodyMarker)
	body.WriteString("\n")
	body.WriteString(openingSummary(artifact, catalog, items, summaryContext))
	body.WriteString("\n\n## Needs Review\n")
	if len(items) == 0 {
		body.WriteString("- No specific changed hunk needs extra review from the available GX context.\n")
	} else {
		for _, item := range items {
			body.WriteString("- ")
			if item.Link != "" {
				body.WriteString("[")
				body.WriteString(item.Title)
				body.WriteString("](")
				body.WriteString(item.Link)
				body.WriteString(")")
			} else {
				body.WriteString(item.Title)
			}
			if item.Detail != "" {
				body.WriteString(": ")
				body.WriteString(item.Detail)
			}
			body.WriteByte('\n')
		}
	}
	return strings.TrimRight(body.String(), "\n")
}

func openingSummary(artifact reviewbundle.Artifact, catalog prBodyCatalog, items []prNeedsReviewItem, summaryContext prSummaryContext) string {
	change := changeSummarySentence(artifact, catalog)
	readiness := readinessSentence(catalog, len(items))
	context := contextSentence(summaryContext)
	parts := []string{change, readiness}
	if context != "" {
		parts = append(parts, context)
	}
	return strings.Join(parts, " ")
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
	name := stackName(artifact)
	if name != "" {
		return fmt.Sprintf("This PR publishes %d GX revisions for `%s`: %s.", len(catalog.Revisions), name, strings.Join(descriptions, "; "))
	}
	return fmt.Sprintf("This PR publishes %d GX revisions: %s.", len(catalog.Revisions), strings.Join(descriptions, "; "))
}

func readinessSentence(catalog prBodyCatalog, issueCount int) string {
	blast := catalog.Stats.MaxRiskLevel
	if blast == "" {
		blast = "low"
	}
	detail := fmt.Sprintf("%d file(s), %d area(s), +%d/-%d lines", catalog.Stats.FileCount, catalog.Stats.AreaCount, catalog.Stats.AddedLines, catalog.Stats.DeletedLines)
	switch {
	case issueCount > 0:
		return fmt.Sprintf("Blast radius is %s (%s), so merge after reviewing the linked hunk%s below.", blast, detail, plural(issueCount))
	case blast == "high":
		return fmt.Sprintf("Blast radius is high (%s), so this should get focused review even though GX did not surface a concrete blocker.", detail)
	case blast == "medium":
		return fmt.Sprintf("Blast radius is medium (%s), and GX did not surface a concrete blocker.", detail)
	default:
		return fmt.Sprintf("Blast radius is low (%s), and GX did not surface a concrete blocker.", detail)
	}
}

func contextSentence(summaryContext prSummaryContext) string {
	if len(summaryContext.Sources) == 0 {
		return ""
	}
	sources := limitStrings(summaryContext.Sources, 3)
	return "Context used: " + strings.Join(sources, ", ") + "."
}

func needsReviewItems(artifact reviewbundle.Artifact, catalog prBodyCatalog, findings []codereview.Finding) []prNeedsReviewItem {
	var items []prNeedsReviewItem
	for _, finding := range findings {
		title := strings.TrimSpace(finding.Title)
		if title == "" {
			continue
		}
		detail := firstNonEmpty(strings.TrimSpace(finding.Summary), strings.TrimSpace(finding.Recommendation))
		link := hunkLinkForFinding(catalog.Hunks, finding)
		items = append(items, prNeedsReviewItem{
			Title:  title,
			Detail: trimSentence(detail, 220),
			Link:   link,
		})
		if len(items) >= maxPRBodyItems {
			return items
		}
	}
	if len(items) > 0 {
		return items
	}
	return heuristicNeedsReviewItems(artifact, catalog)
}

func heuristicNeedsReviewItems(_ reviewbundle.Artifact, catalog prBodyCatalog) []prNeedsReviewItem {
	var items []prNeedsReviewItem
	for _, signal := range catalog.Stats.RiskSignals {
		switch {
		case strings.Contains(signal, "warning:") || strings.Contains(signal, "structural"):
			if hunk := firstHunk(catalog.Hunks, ""); hunk != nil {
				items = append(items, prNeedsReviewItem{
					Title:  "Review structural publish assumptions",
					Detail: "GX saw structural or feasibility risk in this stack. Check the linked changed hunk before merging.",
					Link:   hunk.Link,
				})
			}
		case strings.Contains(signal, "many_files") || strings.Contains(signal, "many_hunks"):
			if hunk := firstHunk(catalog.Hunks, ""); hunk != nil {
				items = append(items, prNeedsReviewItem{
					Title:  "Review the broadest changed hunk",
					Detail: "The PR spans enough files or hunks that a focused pass is still useful even without a precise blocker.",
					Link:   hunk.Link,
				})
			}
		}
		if len(items) >= maxPRBodyItems {
			return dedupeNeedsReviewItems(items)
		}
	}
	for _, candidate := range []struct {
		filePrefix string
		title      string
		detail     string
	}{
		{
			filePrefix: "internal/publication/",
			title:      "Verify PR body refresh timing",
			detail:     "This path runs during publish, so review whether the body update still happens when artifact upload, indexing, or GitHub calls fail.",
		},
		{
			filePrefix: "internal/vcs/",
			title:      "Verify initial PR creation body",
			detail:     "This path creates or refreshes the GitHub PR before the review bundle exists, so it should not publish visible fallback prose.",
		},
		{
			filePrefix: "internal/github/",
			title:      "Verify GitHub ownership and auth handling",
			detail:     "This path reads and patches PR bodies. Confirm it only rewrites GX-owned bodies and reports auth failures clearly.",
		},
	} {
		if hunk := firstHunk(catalog.Hunks, candidate.filePrefix); hunk != nil {
			items = append(items, prNeedsReviewItem{
				Title:  candidate.title,
				Detail: candidate.detail,
				Link:   hunk.Link,
			})
		}
		if len(items) >= maxPRBodyItems {
			break
		}
	}
	return dedupeNeedsReviewItems(items)
}

func reviewPRSummaryFindings(ctx context.Context, artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext) []codereview.Finding {
	if prSummaryReviewerFromEnv == nil {
		return nil
	}
	reviewer := prSummaryReviewerFromEnv()
	if reviewer == nil {
		return nil
	}
	brief := prReviewBrief(artifact, catalog, summaryContext)
	findings, err := reviewer.Review(ctx, brief)
	if err != nil {
		return nil
	}
	return patchRelevantFindings(catalog.Hunks, findings)
}

func prReviewBrief(artifact reviewbundle.Artifact, catalog prBodyCatalog, summaryContext prSummaryContext) codereview.ReviewBrief {
	diffSnippets := make([]codereview.DiffSnippet, 0, len(catalog.Hunks))
	for _, hunk := range catalog.Hunks {
		diffSnippets = append(diffSnippets, codereview.DiffSnippet{
			File: hunk.File,
			Diff: hunk.Patch,
		})
		if len(diffSnippets) >= 24 {
			break
		}
	}
	return codereview.ReviewBrief{
		RepoRoot:      artifact.Bundle.Repo.RootPath,
		Scope:         "maintainability",
		Depth:         "shallow",
		ReviewProfile: "patch_focused",
		Focus:         strings.Join(catalog.Files, " "),
		Static: codereview.StaticSnapshot{
			FileCount:     catalog.Stats.FileCount,
			ChangedFiles:  catalog.Files,
			DiffSnippets:  diffSnippets,
			CodeQuality:   warningQualityHints(artifact),
			TestFileCount: countTestFiles(catalog.Files),
		},
		Context: summaryContext.Snippets,
		Rubric: codereview.ArchitectureRubric{
			Goal: "Produce concrete PR review findings for changed hunks only.",
			Questions: []string{
				"Can this PR be merged without much review?",
				"Which changed hunks are most likely to create a user-visible or developer-visible regression?",
				"What should the reviewer verify first?",
			},
			Reject: []string{
				"Do not list whole files when a changed hunk is available.",
				"Do not emit generic best-practice advice.",
				"Do not repeat blast-radius facts as findings.",
			},
			Output: "At most three concrete findings tied to changed hunks.",
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
	return prSummaryContext{
		Snippets: dedupePRContextSnippets(snippets),
		Sources:  sortedUnique(sources),
	}, nil
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
	if err != nil {
		return codereview.ContextSnippet{}, false
	}
	text := string(data)
	if strings.TrimSpace(text) == "" {
		return codereview.ContextSnippet{}, false
	}
	return codereview.ContextSnippet{
		Kind:        kind,
		Ref:         rel,
		Source:      "local",
		SourceLabel: rel,
		Text:        limitText(text, maxPRContextSnippetSize),
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
			Text:        limitText(chunk.Text, maxPRContextSnippetSize),
		})
		if len(out) >= 4 {
			break
		}
	}
	return out
}

func indexedPRContextSnippets(ctx context.Context, artifact reviewbundle.Artifact, catalog prBodyCatalog) ([]codereview.ContextSnippet, []string) {
	cfg, ok := prSemanticConfigFromEnv()
	if !ok {
		return nil, nil
	}
	queryText := prContextQueryText(artifact, catalog)
	if queryText == "" {
		return nil, nil
	}
	embedder := semantic.NewOpenAIEmbedder(cfg)
	vectors, err := embedder.Embed(ctx, []string{queryText})
	if err != nil || len(vectors) != 1 {
		return nil, nil
	}
	var snippets []codereview.ContextSnippet
	var sources []string
	if cfg.TurboPufferNamespace != "" {
		rows := queryPRTurboPuffer(ctx, cfg, cfg.TurboPufferNamespace, vectors[0], []any{
			"Or",
			[]any{
				[]any{"source_kind", "Eq", "code_file"},
				[]any{"source_kind", "Eq", "session_transcript"},
			},
		}, 6)
		if len(rows) > 0 {
			snippets = append(snippets, rowsToPRContextSnippets(rows, "turbopuffer:"+cfg.TurboPufferNamespace)...)
			sources = append(sources, "indexed code/session")
		}
	}
	reviewNamespace := strings.TrimSpace(os.Getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE"))
	if reviewNamespace == "" {
		reviewNamespace = "gx-review-knowledge"
	}
	rows := queryPRTurboPuffer(ctx, cfg, reviewNamespace, vectors[0], []any{"And", []any{[]any{"source_kind", "Eq", "review_knowledge"}}}, 4)
	if len(rows) > 0 {
		snippets = append(snippets, rowsToPRContextSnippets(rows, "turbopuffer:"+reviewNamespace)...)
		sources = append(sources, "review resources")
	}
	return snippets, sortedUnique(sources)
}

func prSemanticConfigFromEnv() (semantic.Config, bool) {
	apiKey := strings.TrimSpace(firstNonEmpty(os.Getenv("OPENAI_API_KEY"), os.Getenv("GX_OPENAI_API_KEY")))
	tpufKey := strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY"))
	if apiKey == "" || tpufKey == "" {
		return semantic.Config{}, false
	}
	namespace := strings.TrimSpace(os.Getenv("GX_TPUF_NAMESPACE"))
	if namespace == "" {
		namespace = "gx-sessions"
	}
	return semantic.Config{
		Enabled:              true,
		OpenAIAPIKey:         apiKey,
		OpenAIBaseURL:        normalizeOpenAIBaseURL(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL"), "https://api.openai.com")),
		OpenAIEmbeddingModel: firstNonEmpty(os.Getenv("GX_OPENAI_EMBEDDING_MODEL"), "text-embedding-3-small"),
		EmbeddingDimensions:  envInt("GX_EMBEDDING_DIMENSIONS", 512),
		TurboPufferAPIKey:    tpufKey,
		TurboPufferBaseURL:   firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), "https://gcp-us-central1.turbopuffer.com"),
		TurboPufferNamespace: namespace,
		BatchSize:            16,
		MaxChunkBytes:        maxPRContextSnippetSize,
	}, true
}

func queryPRTurboPuffer(ctx context.Context, cfg semantic.Config, namespace string, vector []float32, filters any, limit int) []map[string]any {
	if limit <= 0 {
		limit = 4
	}
	payload := map[string]any{
		"rank_by": []any{"vector", "ANN", vector},
		"limit":   map[string]any{"total": limit},
		"include_attributes": []string{
			"text",
			"source_kind",
			"source_id",
			"file_path",
			"symbol",
			"start_line",
			"end_line",
			"title",
			"url",
			"category",
		},
	}
	if filters != nil {
		payload["filters"] = filters
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
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil
	}
	var decoded struct {
		Rows []map[string]any `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil
	}
	return decoded.Rows
}

func rowsToPRContextSnippets(rows []map[string]any, source string) []codereview.ContextSnippet {
	var snippets []codereview.ContextSnippet
	for _, row := range rows {
		text := stringAny(row["text"])
		if strings.TrimSpace(text) == "" {
			continue
		}
		ref := firstNonEmpty(stringAny(row["file_path"]), stringAny(row["source_id"]), stringAny(row["url"]), stringAny(row["title"]))
		kind := firstNonEmpty(stringAny(row["source_kind"]), "indexed_context")
		snippets = append(snippets, codereview.ContextSnippet{
			Kind:        kind,
			Ref:         ref,
			Source:      source,
			SourceLabel: ref,
			Text:        limitText(text, maxPRContextSnippetSize),
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
	hunks := pullRequestHunks(artifact)
	stats := bodyStats(artifact, revisions)
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
	return prBodyCatalog{
		Revisions: revisions,
		Hunks:     hunks,
		Stats:     stats,
		Files:     sortedKeys(filesSet),
		Areas:     sortedKeys(areasSet),
	}
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
	}
	return summary
}

func pullRequestHunks(artifact reviewbundle.Artifact) []prHunkSummary {
	prURL := pullRequestURL(artifact)
	var hunks []prHunkSummary
	for revisionIndex, entry := range artifact.Stack {
		revisionLabel := strings.TrimSpace(entry.Change.Description)
		if revisionLabel == "" {
			revisionLabel = fmt.Sprintf("revision %d", revisionIndex+1)
		}
		hunks = append(hunks, parsePatchHunks(entry.Patch, revisionLabel, prURL, len(hunks))...)
	}
	if len(artifact.Stack) == 0 && artifact.Change != nil {
		hunks = append(hunks, parsePatchHunks("", artifact.Change.Description, prURL, 0)...)
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
		current.Link = githubHunkLink(prURL, current.File, current.OldStart, current.NewStart, current.NewLines)
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
	oldStart := atoiDefault(match[1], 0)
	oldLines := atoiDefault(match[2], 1)
	newStart := atoiDefault(match[3], 0)
	newLines := atoiDefault(match[4], 1)
	return oldStart, oldLines, newStart, newLines
}

func changedLineForHunk(hunk prHunkSummary) int {
	line := hunk.NewStart
	for _, patchLine := range strings.Split(hunk.Patch, "\n") {
		if strings.HasPrefix(patchLine, "@@ ") {
			line = hunk.NewStart
			continue
		}
		if strings.HasPrefix(patchLine, "+") && !strings.HasPrefix(patchLine, "+++") {
			return line
		}
		if strings.HasPrefix(patchLine, " ") || strings.HasPrefix(patchLine, "+") {
			line++
		}
	}
	return hunk.NewStart
}

func githubHunkLink(prURL, file string, oldStart, newStart, newLines int) string {
	prURL = strings.TrimRight(strings.TrimSpace(prURL), "/")
	file = strings.TrimSpace(file)
	if prURL == "" || file == "" {
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
	return fmt.Sprintf("%s/files#diff-%s%s%d", prURL, hex.EncodeToString(sum[:]), side, line)
}

func hunkLinkForFinding(hunks []prHunkSummary, finding codereview.Finding) string {
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
	if hunk := firstHunk(hunks, ""); hunk != nil {
		return hunk.Link
	}
	return ""
}

func firstHunk(hunks []prHunkSummary, filePrefix string) *prHunkSummary {
	for i := range hunks {
		if filePrefix == "" || strings.HasPrefix(hunks[i].File, filePrefix) {
			return &hunks[i]
		}
	}
	return nil
}

func patchRelevantFindings(hunks []prHunkSummary, findings []codereview.Finding) []codereview.Finding {
	files := map[string]struct{}{}
	for _, hunk := range hunks {
		files[hunk.File] = struct{}{}
	}
	var out []codereview.Finding
	for _, finding := range findings {
		text := strings.ToLower(strings.Join([]string{finding.Title, finding.Summary, finding.Recommendation, evidenceText(finding.Evidence)}, " "))
		for file := range files {
			if strings.Contains(text, strings.ToLower(file)) || strings.Contains(text, strings.ToLower(path.Base(file))) {
				out = append(out, finding)
				break
			}
		}
	}
	if len(out) > 0 {
		return out
	}
	return findings
}

func warningQualityHints(artifact reviewbundle.Artifact) []codereview.CodeQualityHint {
	var hints []codereview.CodeQualityHint
	for _, entry := range artifact.Stack {
		if entry.Change.ReviewContext == nil {
			continue
		}
		for _, warning := range entry.Change.ReviewContext.FeasibilityWarnings {
			hints = append(hints, codereview.CodeQualityHint{
				Kind:   firstNonEmpty(warning.Source, "review_warning"),
				File:   warning.FromFile,
				Text:   warning.Message,
				Reason: warning.Severity,
			})
		}
	}
	return hints
}

func bodyStats(artifact reviewbundle.Artifact, revisions []prRevisionSummary) prBodyStats {
	stats := prBodyStats{
		RevisionCount: len(revisions),
		MaxRiskLevel:  "low",
	}
	files := map[string]struct{}{}
	areas := map[string]struct{}{}
	riskSignals := map[string]struct{}{}
	for _, entry := range artifact.Stack {
		stats.AddedLines += addedLines(entry.Patch)
		stats.DeletedLines += deletedLines(entry.Patch)
		accumulateContext(&stats, entry.Change.ReviewContext, riskSignals)
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
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		if file := fileFromDiffLine(line); file != "" {
			files = append(files, file)
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
	if len(artifact.Stack) > 0 && strings.TrimSpace(artifact.Stack[0].BranchName) != "" {
		return strings.TrimSpace(artifact.Stack[0].BranchName)
	}
	if artifact.Push.BranchName != nil && strings.TrimSpace(*artifact.Push.BranchName) != "" {
		return strings.TrimSpace(*artifact.Push.BranchName)
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

func escapePath(value string) string {
	parts := strings.Split(value, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

func normalizeOpenAIBaseURL(raw string) string {
	base := strings.TrimRight(strings.TrimSpace(raw), "/")
	if base == "" {
		return "https://api.openai.com"
	}
	base = strings.TrimSuffix(base, "/v1")
	return base
}

func envInt(name string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
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
	return strings.TrimRight(value[:limit], " .,;:") + "."
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

func dedupeNeedsReviewItems(items []prNeedsReviewItem) []prNeedsReviewItem {
	seen := map[string]struct{}{}
	var out []prNeedsReviewItem
	for _, item := range items {
		key := item.Title + "\x00" + item.Link
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
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
