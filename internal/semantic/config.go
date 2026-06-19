package semantic

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultOpenAIBaseURL      = "https://api.openai.com"
	defaultOpenAIEmbedModel   = "text-embedding-3-small"
	defaultEmbeddingDims      = 512
	defaultTurboPufferBaseURL = "https://gcp-us-central1.turbopuffer.com"
	defaultTurboPufferNS      = "gx-sessions"
	defaultBatchSize          = 64
	defaultMaxChunkBytes      = 12000
)

type Config struct {
	Enabled              bool
	OpenAIAPIKey         string
	OpenAIBaseURL        string
	OpenAIEmbeddingModel string
	EmbeddingDimensions  int
	TurboPufferAPIKey    string
	TurboPufferBaseURL   string
	TurboPufferNamespace string
	BatchSize            int
	MaxChunkBytes        int
}

func ConfigFromEnv() (Config, error) {
	enabled := truthy(os.Getenv("GX_SEMANTIC_INDEX")) || strings.TrimSpace(os.Getenv("GX_TPUF_NAMESPACE")) != ""
	cfg := Config{
		Enabled:              enabled,
		OpenAIAPIKey:         strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		OpenAIBaseURL:        envOrDefault("GX_OPENAI_BASE_URL", defaultOpenAIBaseURL),
		OpenAIEmbeddingModel: envOrDefault("GX_OPENAI_EMBEDDING_MODEL", defaultOpenAIEmbedModel),
		EmbeddingDimensions:  envIntOrDefault("GX_EMBEDDING_DIMENSIONS", defaultEmbeddingDims),
		TurboPufferAPIKey:    strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")),
		TurboPufferBaseURL:   envOrDefault("GX_TPUF_BASE_URL", defaultTurboPufferBaseURL),
		TurboPufferNamespace: envOrDefault("GX_TPUF_NAMESPACE", defaultTurboPufferNS),
		BatchSize:            envIntOrDefault("GX_SEMANTIC_BATCH_SIZE", defaultBatchSize),
		MaxChunkBytes:        envIntOrDefault("GX_SEMANTIC_MAX_CHUNK_BYTES", defaultMaxChunkBytes),
	}
	if !cfg.Enabled {
		return cfg, nil
	}
	if cfg.OpenAIAPIKey == "" {
		return cfg, fmt.Errorf("OPENAI_API_KEY is required when semantic indexing is enabled")
	}
	if cfg.TurboPufferAPIKey == "" {
		return cfg, fmt.Errorf("TURBOPUFFER_API_KEY is required when semantic indexing is enabled")
	}
	if cfg.EmbeddingDimensions <= 0 {
		return cfg, fmt.Errorf("GX_EMBEDDING_DIMENSIONS must be positive")
	}
	if cfg.BatchSize <= 0 {
		return cfg, fmt.Errorf("GX_SEMANTIC_BATCH_SIZE must be positive")
	}
	if cfg.MaxChunkBytes <= 0 {
		return cfg, fmt.Errorf("GX_SEMANTIC_MAX_CHUNK_BYTES must be positive")
	}
	if _, err := url.ParseRequestURI(strings.TrimRight(cfg.OpenAIBaseURL, "/")); err != nil {
		return cfg, fmt.Errorf("GX_OPENAI_BASE_URL is invalid: %w", err)
	}
	if _, err := url.ParseRequestURI(strings.TrimRight(cfg.TurboPufferBaseURL, "/")); err != nil {
		return cfg, fmt.Errorf("GX_TPUF_BASE_URL is invalid: %w", err)
	}
	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func envIntOrDefault(name string, fallback int) int {
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

func truthy(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "t", "yes", "y", "on":
		return true
	default:
		return false
	}
}
