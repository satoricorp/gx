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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create OpenAI embeddings request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request OpenAI embeddings: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request OpenAI embeddings: status %s", resp.Status)
	}

	var decoded openAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode OpenAI embeddings response: %w", err)
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
