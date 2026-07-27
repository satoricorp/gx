package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Embedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float32, error)
}

type OpenAIEmbedder struct {
	apiKey     string
	baseURL    string
	model      string
	dimensions int
	httpClient *http.Client
}

func NewOpenAIEmbedder(cfg Config) *OpenAIEmbedder {
	return &OpenAIEmbedder{
		apiKey:     cfg.OpenAIAPIKey,
		baseURL:    strings.TrimRight(cfg.OpenAIBaseURL, "/"),
		model:      cfg.OpenAIEmbeddingModel,
		dimensions: cfg.EmbeddingDimensions,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// embedMaxAttempts bounds retries of one embeddings request.
//
// A full repository index is hundreds of batches wide and runs them
// concurrently, so a single transient upstream 500 used to abort the entire
// run — the first error cancels the shared context and every other in-flight
// batch dies with it. Losing a 30-second index to one blip is not a tradeoff
// worth making when indexing time is free.
const embedMaxAttempts = 4

func (e *OpenAIEmbedder) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	body, err := json.Marshal(openAIEmbeddingRequest{
		Model:          e.model,
		Input:          inputs,
		EncodingFormat: "float",
		Dimensions:     e.dimensions,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal OpenAI embeddings request: %w", err)
	}

	var decoded openAIEmbeddingResponse
	var lastErr error
	for attempt := 0; attempt < embedMaxAttempts; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
		decoded, lastErr = e.embedOnce(ctx, body)
		if lastErr == nil {
			break
		}
		if !isRetryableEmbedError(lastErr) {
			return nil, lastErr
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	if len(decoded.Data) != len(inputs) {
		return nil, fmt.Errorf("OpenAI embeddings response returned %d embeddings for %d inputs", len(decoded.Data), len(inputs))
	}
	out := make([][]float32, len(decoded.Data))
	for _, item := range decoded.Data {
		if item.Index < 0 || item.Index >= len(out) {
			return nil, fmt.Errorf("OpenAI embeddings response has invalid index %d", item.Index)
		}
		vector := make([]float32, len(item.Embedding))
		for i, value := range item.Embedding {
			vector[i] = float32(value)
		}
		out[item.Index] = vector
	}
	return out, nil
}

// embedOnce performs a single embeddings request.
func (e *OpenAIEmbedder) embedOnce(ctx context.Context, body []byte) (openAIEmbeddingResponse, error) {
	var decoded openAIEmbeddingResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return decoded, fmt.Errorf("create OpenAI embeddings request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return decoded, retryableEmbedError{fmt.Errorf("request OpenAI embeddings: %w", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err := fmt.Errorf("request OpenAI embeddings: status %s", resp.Status)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			return decoded, retryableEmbedError{err}
		}
		return decoded, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return decoded, retryableEmbedError{fmt.Errorf("decode OpenAI embeddings response: %w", err)}
	}
	return decoded, nil
}

// retryableEmbedError marks a failure that a later attempt might not hit:
// a transport error, a 429, or any 5xx.
type retryableEmbedError struct{ error }

func (e retryableEmbedError) Unwrap() error { return e.error }

func isRetryableEmbedError(err error) bool {
	var retryable retryableEmbedError
	return errors.As(err, &retryable)
}

type openAIEmbeddingRequest struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	EncodingFormat string   `json:"encoding_format"`
	Dimensions     int      `json:"dimensions,omitempty"`
}

type openAIEmbeddingResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}
