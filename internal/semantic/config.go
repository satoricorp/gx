package semantic

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultOpenAIBaseURL = "https://api.openai.com"
	// defaultOpenAIEmbedModel / defaultEmbeddingDims are the model's native
	// width, not a truncation of it. The previous 512 was a Matryoshka
	// truncation requested at embed time via the API's `dimensions` parameter.
	//
	// The model and the width were both chosen by measurement, not by
	// reputation. Indexing this repository (3883 chunks) four ways and running
	// a 16-query ground-truth set over each:
	//
	//	config                 vector r@10  vector MRR  hybrid r@25
	//	3-small @ 512 (old)          0.513       0.710        0.692
	//	3-small @ 1536 (native)      0.538       0.750        0.692
	//	3-large @ 1536               0.410       0.598        0.667
	//	3-large @ 3072 (native)      0.436       0.568        0.641
	//
	// Paired per query, 3-small@1536 ranks the first correct file higher than
	// 3-large@3072 on 8 of 16 queries and lower on 1 (hybrid: 7 better, 0
	// worse), so the larger model is not a close call on this corpus — it is
	// worse. Going from 512 to the native 1536 is free at review time:
	// embedding latency does not vary with width (the round trip dominates)
	// and TurboPuffer ANN stays in the low tens of milliseconds.
	defaultOpenAIEmbedModel   = "text-embedding-3-small"
	defaultEmbeddingDims      = 1536
	defaultTurboPufferBaseURL = "https://gcp-us-central1.turbopuffer.com"
	// defaultTurboPufferNS is intentionally empty. The old default,
	// "totality-sessions", is a flat global namespace with no org and no repo
	// dimension; it was never created in production and could never match the
	// per-org-per-repo namespaces every reader uses. Callers resolve a
	// namespace with NamespaceForRepo instead.
	defaultTurboPufferNS = ""
	defaultBatchSize     = 64
	defaultMaxChunkBytes = 12000
)

// embeddingModelDimensions is the native (maximum) width of each supported
// embedding model.
var embeddingModelDimensions = map[string]int{
	"text-embedding-3-small": 1536,
	"text-embedding-3-large": 3072,
	"text-embedding-ada-002": 1536,
}

// MaxEmbeddingDimensions reports the native width of a model, or 0 when it is
// unknown (a custom or self-hosted model, where no cap is enforced).
func MaxEmbeddingDimensions(model string) int {
	return embeddingModelDimensions[strings.TrimSpace(model)]
}

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
	enabled := truthy(os.Getenv("TOTALITY_SEMANTIC_INDEX")) || strings.TrimSpace(os.Getenv("TOTALITY_TPUF_NAMESPACE")) != ""
	cfg := Config{
		Enabled:              enabled,
		OpenAIAPIKey:         strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		OpenAIBaseURL:        envOrDefault("TOTALITY_OPENAI_BASE_URL", defaultOpenAIBaseURL),
		OpenAIEmbeddingModel: envOrDefault("TOTALITY_OPENAI_EMBEDDING_MODEL", defaultOpenAIEmbedModel),
		EmbeddingDimensions:  envIntOrDefault("TOTALITY_EMBEDDING_DIMENSIONS", defaultEmbeddingDims),
		TurboPufferAPIKey:    strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")),
		TurboPufferBaseURL:   envOrDefault("TOTALITY_TPUF_BASE_URL", defaultTurboPufferBaseURL),
		TurboPufferNamespace: envOrDefault("TOTALITY_TPUF_NAMESPACE", defaultTurboPufferNS),
		BatchSize:            envIntOrDefault("TOTALITY_SEMANTIC_BATCH_SIZE", defaultBatchSize),
		MaxChunkBytes:        envIntOrDefault("TOTALITY_SEMANTIC_MAX_CHUNK_BYTES", defaultMaxChunkBytes),
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
		return cfg, fmt.Errorf("TOTALITY_EMBEDDING_DIMENSIONS must be positive")
	}
	if cfg.BatchSize <= 0 {
		return cfg, fmt.Errorf("TOTALITY_SEMANTIC_BATCH_SIZE must be positive")
	}
	if cfg.MaxChunkBytes <= 0 {
		return cfg, fmt.Errorf("TOTALITY_SEMANTIC_MAX_CHUNK_BYTES must be positive")
	}
	if _, err := url.ParseRequestURI(strings.TrimRight(cfg.OpenAIBaseURL, "/")); err != nil {
		return cfg, fmt.Errorf("TOTALITY_OPENAI_BASE_URL is invalid: %w", err)
	}
	if _, err := url.ParseRequestURI(strings.TrimRight(cfg.TurboPufferBaseURL, "/")); err != nil {
		return cfg, fmt.Errorf("TOTALITY_TPUF_BASE_URL is invalid: %w", err)
	}
	// The legacy bundle indexer addresses one fixed namespace and has no repo
	// context to derive one from. Failing loudly here beats the old behaviour,
	// which silently pointed every read and write at "totality-sessions" — a
	// namespace that has never existed.
	if cfg.TurboPufferNamespace == "" {
		return cfg, fmt.Errorf("TOTALITY_TPUF_NAMESPACE is required when TOTALITY_SEMANTIC_INDEX is set; repository indexing resolves its own namespace via semantic.NamespaceForRepo")
	}
	return cfg, nil
}

// CodeIndexConfigFromEnv builds a configuration for repository indexing.
//
// Unlike ConfigFromEnv it is not behind TOTALITY_SEMANTIC_INDEX. Repository indexing
// is the thing that makes review retrieval work at all, so it runs whenever
// credentials exist; TOTALITY_SEMANTIC_INDEX=0 still turns it off explicitly. It also
// never returns a namespace, because a namespace is a function of the repo
// being indexed (see NamespaceForRepo), not of the environment.
func CodeIndexConfigFromEnv() Config {
	model := envOrDefault("TOTALITY_OPENAI_EMBEDDING_MODEL", defaultOpenAIEmbedModel)
	dimensions := envIntOrDefault("TOTALITY_EMBEDDING_DIMENSIONS", defaultEmbeddingDims)
	if max := MaxEmbeddingDimensions(model); max > 0 && (dimensions <= 0 || dimensions > max) {
		dimensions = max
	}
	return Config{
		Enabled:              !strings.EqualFold(strings.TrimSpace(os.Getenv("TOTALITY_SEMANTIC_INDEX")), "0"),
		OpenAIAPIKey:         strings.TrimSpace(firstNonEmptyString(os.Getenv("OPENAI_API_KEY"), os.Getenv("TOTALITY_OPENAI_API_KEY"))),
		OpenAIBaseURL:        envOrDefault("TOTALITY_OPENAI_BASE_URL", defaultOpenAIBaseURL),
		OpenAIEmbeddingModel: model,
		EmbeddingDimensions:  dimensions,
		TurboPufferAPIKey:    strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")),
		TurboPufferBaseURL:   envOrDefault("TOTALITY_TPUF_BASE_URL", defaultTurboPufferBaseURL),
		BatchSize:            envIntOrDefault("TOTALITY_SEMANTIC_BATCH_SIZE", defaultBatchSize),
		MaxChunkBytes:        envIntOrDefault("TOTALITY_SEMANTIC_MAX_CHUNK_BYTES", defaultMaxChunkBytes),
	}
}

// NamespaceForRepo computes the TurboPuffer namespace for one repository.
//
// The shape matches the console writer (`totality-<orgId>-<slug>`) so the CLI and the
// server address the same rows, with a schema-version suffix appended: vector
// width and full-text settings are fixed per namespace in TurboPuffer, so a
// schema change has to land in a new namespace rather than corrupting an
// existing one. A repository with no org (never pushed, no credentials) gets a
// `totality-local-` namespace so a first review still has an index to read.
func NamespaceForRepo(orgID, repoFullName, repoRoot string) string {
	orgID = namespaceSlug(orgID)
	slug := namespaceSlug(repoFullName)
	if slug == "" {
		base := namespaceSlug(filepath.Base(strings.TrimRight(repoRoot, string(filepath.Separator))))
		if base == "" {
			base = "repo"
		}
		slug = base + "-" + shortHash(repoRoot)[:12]
	}
	prefix := "totality-local"
	if orgID != "" {
		prefix = "totality-" + orgID
	}
	return fmt.Sprintf("%s-%s-v%d", prefix, slug, IndexSchemaVersion)
}

var namespaceUnsafe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func namespaceSlug(value string) string {
	value = namespaceUnsafe.ReplaceAllString(strings.TrimSpace(value), "-")
	return strings.ToLower(strings.Trim(value, "-"))
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
