package cloud

import (
	"context"
	"strings"
)

// ReviewSearchRequest is the POST /v1/review/search body.
//
// Target "repo" searches the caller's org/repo namespace (the server resolves
// the org and the namespace; the CLI never names one). Target "knowledge"
// searches the shared review-knowledge corpus.
type ReviewSearchRequest struct {
	Target       string   `json:"target,omitempty"`
	RepoFullName string   `json:"repo_full_name,omitempty"`
	Query        string   `json:"query"`
	SymbolQuery  string   `json:"symbol_query,omitempty"`
	SourceKinds  []string `json:"source_kinds,omitempty"`
	Limit        int      `json:"limit,omitempty"`
}

// ReviewSearchRow is one retrieved chunk. Attributes carry the namespace
// schema fields (file_path, start_line, source_kind, …) exactly as stored, so
// the review's snippet builders can consume them like rows read directly.
type ReviewSearchRow struct {
	ID         string         `json:"id"`
	Score      float64        `json:"score,omitempty"`
	Text       string         `json:"text"`
	Attributes map[string]any `json:"attributes"`
}

type ReviewSearchResult struct {
	// Available is false when the server has retrieval disabled outright
	// (no provider keys configured server-side). Reason says why.
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
	Target    string `json:"target,omitempty"`
	// Namespace is the namespace the server actually searched — the review's
	// evidence lines report it so a wrong-namespace bug is visible in output
	// rather than silent.
	Namespace string `json:"namespace,omitempty"`
	// Exists distinguishes "repository never indexed" (actionable: connect it)
	// from "indexed but nothing matched".
	Exists         bool              `json:"exists"`
	ApproxRowCount int               `json:"approx_row_count,omitempty"`
	LastWriteAt    string            `json:"last_write_at,omitempty"`
	Rows           []ReviewSearchRow `json:"rows"`
}

// SearchReviewIndex retrieves review context through Totality Cloud with the
// caller's login. This is the path every onboarded user takes: the TurboPuffer
// and OpenAI keys live on the server, the same way the /tx/openai and
// /tx/bedrock proxies gate inference.
func (c *Client) SearchReviewIndex(ctx context.Context, reqBody ReviewSearchRequest) (ReviewSearchResult, error) {
	var result ReviewSearchResult
	if c == nil || c.url == "" || strings.TrimSpace(reqBody.Query) == "" {
		return result, nil
	}
	if err := c.postJSON(ctx, "/v1/review/search", reqBody, &result); err != nil {
		return ReviewSearchResult{}, err
	}
	return result, nil
}

// HasReviewSearchCredentials reports whether a cloud retrieval call could
// authenticate right now, without making one. Retrievers use it to decide
// between "call the cloud" and "record the source as disabled with a login
// remedy".
func HasReviewSearchCredentials() bool {
	if !CloudConfigured() {
		return false
	}
	_, err := CloudAPIToken()
	return err == nil
}
