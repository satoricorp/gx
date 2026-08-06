package codereview

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
)

// Cloud-mode retrieval: the review's TurboPuffer-backed sources served through
// gx Cloud instead of raw provider keys.
//
// The direct path below each retriever (raw TURBOPUFFER_API_KEY in the
// environment) predates this and remains for development, where querying the
// store without a server in the loop is worth the extra keys. Everyone else
// authenticates with `gx auth login` and retrieves through POST
// /v1/review/search, the same trade the /gx/openai and /gx/bedrock proxies
// made for inference: the provider keys live server-side, and the server
// resolves the org and namespace with the same function every indexer writes
// through — which also retires the CLI's guess at the console's namespace
// (see semantic.ResolveRepoIdentity for what that guessing cost).

// reviewCloudSearcher is the slice of cloud.Client the retrievers use;
// injected in tests.
type reviewCloudSearcher interface {
	SearchReviewIndex(ctx context.Context, req cloud.ReviewSearchRequest) (cloud.ReviewSearchResult, error)
}

// reviewCloudSearcherFromEnv returns the cloud searcher a retriever should
// use, or nil when cloud retrieval is not available (not configured, or not
// signed in).
func reviewCloudSearcherFromEnv() reviewCloudSearcher {
	if !cloud.HasReviewSearchCredentials() {
		return nil
	}
	client := cloud.NewClient()
	if client == nil {
		return nil
	}
	return client
}

// rawTurboPufferKeyPresent reports whether the direct-store path is armed.
func rawTurboPufferKeyPresent() bool {
	return strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")) != ""
}

// cloudSearchRows converts response rows to the indexRow maps the existing
// snippet builders consume: attributes at the top level, id and text merged
// in. Text goes under "text" unless the attributes already carry a body field
// (bodyFieldPreference), so the builders' field preference keeps working.
func cloudSearchRows(rows []cloud.ReviewSearchRow) []indexRow {
	out := make([]indexRow, 0, len(rows))
	for _, row := range rows {
		m := indexRow{}
		for key, value := range row.Attributes {
			m[key] = value
		}
		if _, ok := m["id"]; !ok && row.ID != "" {
			m["id"] = row.ID
		}
		if strings.TrimSpace(row.Text) != "" {
			m["text"] = row.Text
		}
		out = append(out, m)
	}
	return out
}

// cloudFreshnessDetail is codeIndexFreshnessDetail for a cloud response: the
// caution that an old index answers with deleted code.
func cloudFreshnessDetail(lastWriteAt string) string {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(lastWriteAt))
	if err != nil {
		return ""
	}
	age := time.Since(parsed)
	if age < codeIndexStaleAfter {
		return ""
	}
	return fmt.Sprintf("index last written %d day(s) ago; results may describe code that no longer exists", int(age.Hours()/24))
}

// cloudUnavailableDetail normalizes a failed cloud search into one evidence
// line detail.
func cloudUnavailableDetail(err error) string {
	return "gx cloud retrieval failed: " + err.Error()
}

// signInRemedy is what a user does when neither cloud credentials nor raw
// keys are present. The direct-key alternative is deliberately not offered:
// raw provider keys are a development setup, not a remedy.
const signInRemedy = "Run `gx auth login` so review retrieval can use gx Cloud"

// retrieveCodeIndexViaCloud is CodeIndexRetriever.Retrieve for the cloud path.
func retrieveCodeIndexViaCloud(
	ctx context.Context,
	in RetrieveInput,
	searcher reviewCloudSearcher,
	limit int,
) ([]ContextSnippet, error) {
	repoFullName := reviewHistoryRepoFullName(ctx, in.RepoRoot)
	if repoFullName == "" {
		in.Evidence.Record(EvidenceStatus{
			Source: codeIndexEvidenceSource,
			State:  EvidenceMissing,
			Detail: "this checkout has no GitHub remote, so gx Cloud has no repository to search",
			Remedy: "Add a GitHub remote, then connect this repository at https://gx.run/repositories",
		})
		return nil, nil
	}
	queries := codeIndexQueryText(in)
	if queries.empty() {
		return nil, nil
	}
	result, err := searcher.SearchReviewIndex(ctx, cloud.ReviewSearchRequest{
		Target:       "repo",
		RepoFullName: repoFullName,
		Query:        queries.Body,
		SymbolQuery:  queries.Symbols,
		SourceKinds:  []string{"code_file"},
		Limit:        limit,
	})
	if err != nil {
		in.Evidence.Record(EvidenceStatus{
			Source: codeIndexEvidenceSource,
			State:  EvidenceUnavailable,
			Detail: cloudUnavailableDetail(err),
		})
		return nil, nil
	}
	status := EvidenceStatus{Source: codeIndexEvidenceSource, Namespace: result.Namespace}
	switch {
	case !result.Available:
		status.State = EvidenceUnavailable
		status.Detail = "gx cloud retrieval is not configured server-side"
		if result.Reason != "" {
			status.Detail = "gx cloud: " + result.Reason
		}
	case !result.Exists:
		status.State = EvidenceMissing
		status.Detail = "gx Cloud has not indexed this repository"
		status.Remedy = connectRepositoryRemedy
	case len(result.Rows) == 0:
		status.State = EvidenceEmpty
		status.Detail = cloudFreshnessDetail(result.LastWriteAt)
	}
	if status.State != "" {
		in.Evidence.Record(status)
		return nil, nil
	}

	note := cloudFreshnessDetail(result.LastWriteAt)
	rows := cloudSearchRows(result.Rows)
	namespaceFor := map[string]string{}
	noteFor := map[string]string{}
	for _, row := range rows {
		namespaceFor[codeIndexRowKey(row)] = result.Namespace
	}
	noteFor[result.Namespace] = note
	snippets := codeIndexSnippets(rows, limit, namespaceFor, noteFor)
	status.State = EvidenceOK
	status.Snippets = len(snippets)
	status.Detail = joinDetail(note, "via gx cloud")
	in.Evidence.Record(status)
	return snippets, nil
}

// retrieveSessionsViaCloud is SessionContextRetriever.Retrieve for the cloud
// path.
func retrieveSessionsViaCloud(
	ctx context.Context,
	in RetrieveInput,
	searcher reviewCloudSearcher,
	limit int,
) ([]ContextSnippet, error) {
	repoFullName := reviewHistoryRepoFullName(ctx, in.RepoRoot)
	if repoFullName == "" {
		in.Evidence.Record(EvidenceStatus{
			Source: sessionEvidenceSource,
			State:  EvidenceMissing,
			Detail: "this checkout has no GitHub remote, so gx Cloud has no repository to search",
		})
		return nil, nil
	}
	queryText := sessionQueryText(in)
	if strings.TrimSpace(queryText) == "" {
		return nil, nil
	}
	result, err := searcher.SearchReviewIndex(ctx, cloud.ReviewSearchRequest{
		Target:       "repo",
		RepoFullName: repoFullName,
		Query:        queryText,
		SourceKinds:  sessionSourceKinds,
		Limit:        limit,
	})
	if err != nil {
		in.Evidence.Record(EvidenceStatus{
			Source: sessionEvidenceSource,
			State:  EvidenceUnavailable,
			Detail: cloudUnavailableDetail(err),
		})
		return nil, nil
	}
	status := EvidenceStatus{Source: sessionEvidenceSource, Namespace: result.Namespace}
	switch {
	case !result.Available:
		status.State = EvidenceUnavailable
		status.Detail = "gx cloud retrieval is not configured server-side"
		if result.Reason != "" {
			status.Detail = "gx cloud: " + result.Reason
		}
	case !result.Exists:
		status.State = EvidenceMissing
		status.Detail = "no captured sessions have been indexed here"
	case len(result.Rows) == 0:
		status.State = EvidenceEmpty
	}
	if status.State != "" {
		in.Evidence.Record(status)
		return nil, nil
	}

	rows := cloudSearchRows(result.Rows)
	namespaceFor := map[string]string{}
	for _, row := range rows {
		namespaceFor[sessionRowKey(row)] = result.Namespace
	}
	snippets := sessionSnippets(rows, limit, namespaceFor)
	status.State = EvidenceOK
	status.Snippets = len(snippets)
	status.Detail = "via gx cloud"
	in.Evidence.Record(status)
	return snippets, nil
}
