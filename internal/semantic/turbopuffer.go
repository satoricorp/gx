package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	return &TurboPufferClient{
		apiKey:     cfg.TurboPufferAPIKey,
		baseURL:    strings.TrimRight(cfg.TurboPufferBaseURL, "/"),
		namespace:  strings.Trim(cfg.TurboPufferNamespace, "/"),
		dimensions: cfg.EmbeddingDimensions,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/namespaces/"+c.namespace, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create turbopuffer upsert request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upsert turbopuffer rows: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("upsert turbopuffer rows: status %s", resp.Status)
	}
	return nil
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

func (c *TurboPufferClient) schema() map[string]any {
	return transcriptTurboPufferSchema(c.dimensions)
}

type turboPufferUpsertRequest struct {
	DistanceMetric string           `json:"distance_metric"`
	Schema         map[string]any   `json:"schema"`
	UpsertRows     []map[string]any `json:"upsert_rows"`
}
