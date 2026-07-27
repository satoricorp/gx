package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type VectorStore interface {
	Upsert(ctx context.Context, rows []VectorRow) error
}

type StaleCodeDocumentDeleter interface {
	DeleteStaleCodeDocuments(ctx context.Context, repoFullName, commitID string) error
}

type VectorRow struct {
	ID         string
	Vector     []float32
	Attributes map[string]any
}

type TurboPufferClient struct {
	apiKey     string
	baseURL    string
	namespace  string
	dimensions int
	httpClient *http.Client
}

func NewTurboPufferClient(cfg Config) *TurboPufferClient {
	return NewTurboPufferClientForNamespace(cfg, cfg.TurboPufferNamespace)
}

// NewTurboPufferClientForNamespace targets a namespace other than the one in
// the config. Namespaces are per org per repo, so a single process routinely
// writes to more than one.
func NewTurboPufferClientForNamespace(cfg Config, namespace string) *TurboPufferClient {
	return &TurboPufferClient{
		apiKey:     cfg.TurboPufferAPIKey,
		baseURL:    strings.TrimRight(cfg.TurboPufferBaseURL, "/"),
		namespace:  strings.Trim(namespace, "/"),
		dimensions: cfg.EmbeddingDimensions,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// Namespace reports the namespace this client writes to.
func (c *TurboPufferClient) Namespace() string { return c.namespace }

func (c *TurboPufferClient) Upsert(ctx context.Context, rows []VectorRow) error {
	if len(rows) == 0 {
		return nil
	}
	upsertRows := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" {
			return fmt.Errorf("turbopuffer row has empty id")
		}
		if len(row.Vector) != c.dimensions {
			return fmt.Errorf("turbopuffer row %s has vector dimension %d, want %d", row.ID, len(row.Vector), c.dimensions)
		}
		payload := map[string]any{
			"id":     row.ID,
			"vector": row.Vector,
		}
		for key, value := range row.Attributes {
			payload[key] = value
		}
		upsertRows = append(upsertRows, payload)
	}
	body, err := json.Marshal(turboPufferUpsertRequest{
		DistanceMetric: "cosine_distance",
		Schema:         c.schema(),
		UpsertRows:     upsertRows,
	})
	if err != nil {
		return fmt.Errorf("marshal turbopuffer upsert request: %w", err)
	}
	// Upserts are idempotent — a row id is a pure function of (repo, path,
	// chunk index) — so a retry can only ever rewrite the same row with the
	// same content. A full index issues hundreds of these concurrently and the
	// first failure cancels the shared context, so without a retry one reset
	// connection discards the whole run.
	var lastErr error
	for attempt := 0; attempt < upsertMaxAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		retryable := false
		retryable, lastErr = c.upsertOnce(ctx, body)
		if lastErr == nil {
			return nil
		}
		if !retryable {
			return lastErr
		}
	}
	return lastErr
}

// upsertMaxAttempts bounds retries of one upsert request.
const upsertMaxAttempts = 4

// upsertOnce performs a single upsert and reports whether a failure is worth
// retrying.
func (c *TurboPufferClient) upsertOnce(ctx context.Context, body []byte) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/namespaces/"+c.namespace, bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("create turbopuffer upsert request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("upsert turbopuffer rows: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return retryable, fmt.Errorf("upsert turbopuffer rows: status %s", resp.Status)
	}
	return false, nil
}

func (c *TurboPufferClient) DeleteStaleCodeDocuments(ctx context.Context, repoFullName, commitID string) error {
	repoFullName = strings.TrimSpace(repoFullName)
	commitID = strings.TrimSpace(commitID)
	if repoFullName == "" || commitID == "" {
		return nil
	}
	body, err := json.Marshal(map[string]any{
		"delete_by_filter": []any{
			"And",
			[]any{
				[]any{"source_kind", "Eq", "code_file"},
				[]any{"repo_full_name", "Eq", repoFullName},
				[]any{"commit_id", "NotEq", commitID},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("marshal turbopuffer stale delete request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/namespaces/"+c.namespace, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create turbopuffer stale delete request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete stale turbopuffer rows: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("delete stale turbopuffer rows: status %s", resp.Status)
	}
	return nil
}

// DeleteRows removes rows by id. Re-indexing uses this for chunks that vanished
// because a file shrank or was deleted; it is exact, unlike a delete-by-filter
// sweep, and cannot evict rows that a concurrent writer just wrote.
func (c *TurboPufferClient) DeleteRows(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	for start := 0; start < len(ids); start += 500 {
		end := start + 500
		if end > len(ids) {
			end = len(ids)
		}
		body, err := json.Marshal(map[string]any{"deletes": ids[start:end]})
		if err != nil {
			return fmt.Errorf("marshal turbopuffer delete request: %w", err)
		}
		if _, err := c.do(ctx, http.MethodPost, "", body); err != nil {
			return fmt.Errorf("delete turbopuffer rows: %w", err)
		}
	}
	return nil
}

// QueryRequest is a single ranked query. Exactly one of Vector or (TextField,
// TextQuery) is required; supplying both issues a hybrid query fused with
// reciprocal rank fusion, which is how lexical identifier matches and semantic
// matches are combined without hand-tuned score weighting.
type QueryRequest struct {
	Vector            []float32
	TextField         string
	TextQuery         string
	Limit             int
	Filters           any
	IncludeAttributes []string
}

// QueryRow is one ranked result: its id, distance/score and attributes.
type QueryRow map[string]any

// Query runs a vector, BM25 or hybrid query against the namespace.
func (c *TurboPufferClient) Query(ctx context.Context, req QueryRequest) ([]QueryRow, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	build := func(rankBy []any) map[string]any {
		query := map[string]any{
			"rank_by": rankBy,
			"limit":   limit,
		}
		if req.Filters != nil {
			query["filters"] = req.Filters
		}
		if len(req.IncludeAttributes) > 0 {
			query["include_attributes"] = req.IncludeAttributes
		} else {
			query["include_attributes"] = true
		}
		return query
	}

	var payload map[string]any
	hasVector := len(req.Vector) > 0
	hasText := strings.TrimSpace(req.TextField) != "" && strings.TrimSpace(req.TextQuery) != ""
	switch {
	case hasVector && hasText:
		vectorQuery := build([]any{"vector", "ANN", req.Vector})
		textQuery := build([]any{req.TextField, "BM25", req.TextQuery})
		vectorQuery["limit"] = map[string]any{"total": limit}
		textQuery["limit"] = map[string]any{"total": limit}
		payload = map[string]any{
			"queries":   []any{vectorQuery, textQuery},
			"rerank_by": []any{"RRF"},
		}
	case hasVector:
		payload = build([]any{"vector", "ANN", req.Vector})
	case hasText:
		payload = build([]any{req.TextField, "BM25", req.TextQuery})
	default:
		return nil, fmt.Errorf("turbopuffer query requires a vector or a text query")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal turbopuffer query: %w", err)
	}
	raw, err := c.do(ctx, http.MethodPost, "query", body)
	if err != nil {
		return nil, fmt.Errorf("query turbopuffer rows: %w", err)
	}
	var decoded struct {
		Rows    []QueryRow `json:"rows"`
		Results []struct {
			Rows []QueryRow `json:"rows"`
		} `json:"results"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return nil, fmt.Errorf("decode turbopuffer query response: %w", err)
	}
	if len(decoded.Rows) > 0 {
		return decoded.Rows, nil
	}
	var out []QueryRow
	for _, result := range decoded.Results {
		out = append(out, result.Rows...)
	}
	return out, nil
}

// DeleteNamespace drops the namespace entirely. Used by scratch-namespace
// integration tests and by a forced rebuild.
func (c *TurboPufferClient) DeleteNamespace(ctx context.Context) error {
	if _, err := c.do(ctx, http.MethodDelete, "", nil); err != nil {
		return fmt.Errorf("delete turbopuffer namespace: %w", err)
	}
	return nil
}

func (c *TurboPufferClient) do(ctx context.Context, method, suffix string, body []byte) ([]byte, error) {
	endpoint := c.baseURL + "/v2/namespaces/" + url.PathEscape(c.namespace)
	if suffix != "" {
		endpoint += "/" + suffix
	}
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	payload, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(payload))
		if len(detail) > 500 {
			detail = detail[:500]
		}
		return nil, fmt.Errorf("status %s: %s", resp.Status, detail)
	}
	if readErr != nil {
		return nil, readErr
	}
	return payload, nil
}

func (c *TurboPufferClient) schema() map[string]any {
	return transcriptTurboPufferSchema(c.dimensions)
}

type turboPufferUpsertRequest struct {
	DistanceMetric string           `json:"distance_metric"`
	Schema         map[string]any   `json:"schema"`
	UpsertRows     []map[string]any `json:"upsert_rows"`
}
