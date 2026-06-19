package codereview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/version"
)

type CompositeContextRetriever struct {
	Retrievers []ContextRetriever
}

type EndpointContextRetriever struct {
	URL    string
	Token  string
	Client *http.Client
}

type contextQueryRequest struct {
	RepoRoot        string       `json:"repo_root"`
	Scope           string       `json:"scope"`
	Depth           string       `json:"depth"`
	Focus           string       `json:"focus,omitempty"`
	Hints           []ReviewHint `json:"hints"`
	DependencyFiles []string     `json:"dependency_files"`
	ChangedFiles    []string     `json:"changed_files"`
}

type contextQueryResponse struct {
	Snippets []ContextSnippet `json:"snippets"`
}

func contextRetrieverFromEnv() ContextRetriever {
	retrievers := []ContextRetriever{LocalContextRetriever{}}
	if url := strings.TrimSpace(os.Getenv("GX_REVIEW_CONTEXT_URL")); url != "" {
		token := strings.TrimSpace(os.Getenv("GX_REVIEW_CONTEXT_TOKEN"))
		if token == "" {
			token, _ = cloud.CloudAPIToken()
		}
		retrievers = append(retrievers, EndpointContextRetriever{
			URL:    url,
			Token:  token,
			Client: &http.Client{Timeout: 30 * time.Second},
		})
	}
	return CompositeContextRetriever{Retrievers: retrievers}
}

func (r CompositeContextRetriever) Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	var out []ContextSnippet
	for _, retriever := range r.Retrievers {
		if retriever == nil {
			continue
		}
		snippets, err := retriever.Retrieve(ctx, repoRoot, opts, facts, hints)
		if err != nil {
			continue
		}
		out = append(out, snippets...)
	}
	return dedupeContextSnippets(out), nil
}

func (r EndpointContextRetriever) Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	if strings.TrimSpace(r.URL) == "" {
		return nil, nil
	}
	payload := contextQueryRequest{
		RepoRoot:        repoRoot,
		Scope:           opts.Scope,
		Depth:           depthLabel(opts.Deep),
		Focus:           strings.TrimSpace(opts.Focus),
		Hints:           hints,
		DependencyFiles: facts.DependencyFiles,
		ChangedFiles:    changedFiles(ctx, repoRoot),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal context query: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create context query: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gx/"+version.Current())
	if strings.TrimSpace(r.Token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(r.Token))
	}
	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("query indexed review context: %w", err)
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("indexed context status %s", resp.Status)
	}
	var parsed contextQueryResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode indexed context: %w", err)
	}
	return parsed.Snippets, nil
}

func dedupeContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	seen := map[string]struct{}{}
	var out []ContextSnippet
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
