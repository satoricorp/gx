package codereview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/lgtm/internal/cloud"
	"github.com/satoricorp/lgtm/internal/semantic"
)

type ReviewHistoryRetriever struct {
	Client *cloud.Client
	Limit  int
}

func reviewHistoryRetrieverFromEnv() ContextRetriever {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LGTM_REVIEW_HISTORY_CONTEXT")), "0") {
		return nil
	}
	client := cloud.NewClient()
	if client == nil {
		return nil
	}
	limit := reviewEnvInt("LGTM_REVIEW_HISTORY_TOP_K", 8)
	return ReviewHistoryRetriever{Client: client, Limit: limit}
}

// EvidenceSource implements evidenceNamer.
func (ReviewHistoryRetriever) EvidenceSource() string { return "prior review findings" }

func (r ReviewHistoryRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	if r.Client == nil {
		return nil, nil
	}
	repoRoot := in.RepoRoot
	opts := in.Options
	facts := in.Facts
	hints := in.Hints
	repoFullName := reviewHistoryRepoFullName(ctx, repoRoot)
	if repoFullName == "" {
		return nil, nil
	}
	limit := r.Limit
	if limit <= 0 {
		limit = 8
	}
	if opts.Deep && limit < 15 {
		limit = 15
	}
	result, err := r.Client.SearchCodeReviewHistory(ctx, cloud.CodeReviewHistorySearchRequest{
		RepoFullName: repoFullName,
		Query:        reviewHistoryQuery(opts, facts, hints),
		FilePaths:    facts.Files,
		Limit:        limit,
	})
	if err != nil {
		return nil, err
	}
	return reviewHistorySnippets(result, limit), nil
}

// reviewHistoryRepoFullName is the repository's "owner/name".
//
// It delegates to the same resolver the index writes through. It used to have
// its own copy — `git config --get remote.origin.url` then `remote.upstream.url`,
// parsed by cloud.RepoFullNameFromRemoteURL — while indexing used
// `git remote get-url origin` parsed by a different function. Two copies of a
// derivation are two answers waiting to differ, and they did.
func reviewHistoryRepoFullName(ctx context.Context, repoRoot string) string {
	return semantic.RepoFullNameForRoot(ctx, repoRoot)
}

func reviewHistoryQuery(opts Options, facts RepoFacts, hints []ReviewHint) string {
	parts := []string{
		"lgtm code review history recall",
		strings.TrimSpace(opts.Prompt),
		strings.TrimSpace(opts.Scope),
		strings.Join(facts.DependencyFiles, " "),
		reviewPolicyQueryText(opts.ReviewPolicy),
	}
	for _, file := range facts.Files {
		if language := reviewHistoryLanguage(file); language != "" {
			parts = append(parts, language, file)
		}
	}
	for _, hint := range hints {
		parts = append(parts, hint.Kind, hint.Title, hint.Why, strings.Join(hint.Modules, " "))
	}
	return strings.Join(parts, "\n")
}

func reviewHistorySnippets(result cloud.CodeReviewHistorySearchResult, limit int) []ContextSnippet {
	var snippets []ContextSnippet
	for _, finding := range result.ExactMatches {
		text := strings.TrimSpace(strings.Join([]string{
			"Prior code review finding.",
			"Outcome: " + finding.Outcome,
			"Category: " + finding.Category,
			"Language: " + finding.Language,
			"File: " + finding.FilePath,
			"Title: " + finding.Title,
			"Summary: " + finding.Summary,
			"Recommendation: " + finding.Recommendation,
		}, "\n"))
		snippets = append(snippets, ContextSnippet{
			Kind:      "code_review_history",
			Ref:       firstNonEmpty(finding.Fingerprint, finding.ID),
			Source:    "lgtm-cloud:code-review-history",
			Publisher: "prior lgtm reviews",
			Title:     finding.Title,
			Text:      text,
			File:      finding.FilePath,
			StartLine: finding.LineStart,
			EndLine:   finding.LineEnd,
			ChunkHash: finding.Fingerprint,
		})
		if len(snippets) >= limit {
			return snippets
		}
	}
	for _, summary := range result.Summaries {
		snippets = append(snippets, ContextSnippet{
			Kind:      "code_review_summary",
			Ref:       summary.ID,
			Source:    "lgtm-cloud:code-review-history",
			Publisher: "prior lgtm reviews",
			Title:     fmt.Sprintf("Prior %s summary", firstNonEmpty(summary.SummaryKind, "review")),
			Text:      strings.TrimSpace(summary.SummaryText),
		})
		if len(snippets) >= limit {
			return snippets
		}
	}
	for _, match := range result.SimilarMatches {
		text := strings.TrimSpace(match.Text)
		if text == "" {
			continue
		}
		snippets = append(snippets, ContextSnippet{
			Kind:      "code_review_history",
			Ref:       firstNonEmpty(match.ID, stringAttribute(match.Attributes, "review_fingerprint")),
			Source:    "turbopuffer:code-review-history",
			Publisher: "prior lgtm reviews",
			Title:     firstNonEmpty(stringAttribute(match.Attributes, "review_category"), "Code review history"),
			Text:      text,
			File:      stringAttribute(match.Attributes, "file"),
			ChunkHash: stringAttribute(match.Attributes, "review_fingerprint"),
		})
		if len(snippets) >= limit {
			return snippets
		}
	}
	return snippets
}

func reviewHistoryLanguage(file string) string {
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

func stringAttribute(attrs map[string]any, key string) string {
	if attrs == nil {
		return ""
	}
	value, _ := attrs[key].(string)
	return strings.TrimSpace(value)
}
