package codereview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/semantic"
)

const (
	defaultReviewKnowledgeNamespace = "gx-review-knowledge"
	defaultReviewResourceBaseURL    = "https://gcp-us-central1.turbopuffer.com"
	defaultReviewResourceTopK       = 8
	defaultReviewResourceDeepTopK   = 24
)

type ReviewResourceRetriever struct {
	Embedder  reviewResourceEmbedder
	Store     reviewResourceStore
	Namespace string
	Limit     int
}

type reviewResourceEmbedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float32, error)
}

type reviewResourceStore interface {
	Query(ctx context.Context, req reviewResourceQuery) ([]reviewResourceRow, error)
}

type reviewResourceQuery struct {
	Vector            []float32
	Text              string
	Limit             int
	Filters           any
	IncludeAttributes []string
}

type reviewResourceRow map[string]any

type turboPufferReviewResourceStore struct {
	apiKey     string
	baseURL    string
	namespace  string
	httpClient *http.Client
}

func reviewResourceRetrieverFromEnv() ContextRetriever {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_RESOURCES")), "0") {
		return nil
	}
	openAIKey := strings.TrimSpace(firstNonEmpty(os.Getenv("OPENAI_API_KEY"), os.Getenv("GX_OPENAI_API_KEY")))
	tpufKey := strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY"))
	if openAIKey == "" || tpufKey == "" {
		return nil
	}
	namespace := strings.TrimSpace(os.Getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE"))
	if namespace == "" {
		namespace = defaultReviewKnowledgeNamespace
	}
	cfg := semantic.Config{
		OpenAIAPIKey:         openAIKey,
		OpenAIBaseURL:        normalizeReviewOpenAIBaseURL(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL"), "https://api.openai.com")),
		OpenAIEmbeddingModel: firstNonEmpty(os.Getenv("GX_OPENAI_EMBEDDING_MODEL"), "text-embedding-3-small"),
		EmbeddingDimensions:  reviewEnvInt("GX_EMBEDDING_DIMENSIONS", 512),
	}
	return ReviewResourceRetriever{
		Embedder: semantic.NewOpenAIEmbedder(cfg),
		Store: turboPufferReviewResourceStore{
			apiKey:     tpufKey,
			baseURL:    strings.TrimRight(firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), defaultReviewResourceBaseURL), "/"),
			namespace:  strings.Trim(namespace, "/"),
			httpClient: &http.Client{Timeout: 20 * time.Second},
		},
		Namespace: namespace,
		Limit:     reviewEnvInt("GX_REVIEW_RESOURCES_TOP_K", defaultReviewResourceTopK),
	}
}

func (r ReviewResourceRetriever) Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	if r.Embedder == nil || r.Store == nil {
		return nil, nil
	}
	signals := reviewResourceSignals(ctx, repoRoot, opts, facts, hints)
	queryText := reviewResourceQueryText(opts, signals)
	if strings.TrimSpace(queryText) == "" {
		return nil, nil
	}
	vectors, err := r.Embedder.Embed(ctx, []string{queryText})
	if err != nil {
		return nil, err
	}
	if len(vectors) != 1 {
		return nil, fmt.Errorf("review resource embedding returned %d vectors", len(vectors))
	}

	limit := r.Limit
	if opts.Deep && limit < defaultReviewResourceDeepTopK {
		limit = defaultReviewResourceDeepTopK
	}
	if limit <= 0 {
		limit = defaultReviewResourceTopK
	}
	include := []string{
		"source_id",
		"url",
		"title",
		"category",
		"authority",
		"evidence_level",
		"language_tags",
		"framework_tags",
		"risk_tag_values",
		"review_tag_values",
		"chunk_kind",
		"chunk_index",
		"text",
	}
	broadRows, err := r.Store.Query(ctx, reviewResourceQuery{
		Vector:            vectors[0],
		Text:              queryText,
		Limit:             maxInt(2, limit/2),
		Filters:           reviewResourceBaseFilter(),
		IncludeAttributes: include,
	})
	if err != nil {
		return nil, err
	}
	rows := append([]reviewResourceRow{}, broadRows...)
	if filter := reviewResourceSignalFilter(signals); filter != nil {
		if filteredRows, err := r.Store.Query(ctx, reviewResourceQuery{
			Vector:            vectors[0],
			Text:              queryText,
			Limit:             limit,
			Filters:           filter,
			IncludeAttributes: include,
		}); err == nil {
			rows = append(rows, filteredRows...)
		}
	}
	return reviewResourceSnippets(rows, limit, r.Namespace), nil
}

func (s turboPufferReviewResourceStore) Query(ctx context.Context, req reviewResourceQuery) ([]reviewResourceRow, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = defaultReviewResourceTopK
	}
	payload := map[string]any{
		"rank_by":            []any{"vector", "ANN", req.Vector},
		"limit":              map[string]any{"total": limit, "per": map[string]any{"attributes": []string{"source_id"}, "limit": 2}},
		"include_attributes": req.IncludeAttributes,
	}
	if req.Filters != nil {
		payload["filters"] = req.Filters
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal review resource query: %w", err)
	}
	endpoint := strings.TrimRight(s.baseURL, "/") + "/v2/namespaces/" + url.PathEscape(strings.Trim(s.namespace, "/")) + "/query"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create review resource query: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	client := s.httpClient
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("query review resources: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("review resource index is still building")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("query review resources: status %s", resp.Status)
	}
	var decoded struct {
		Rows []reviewResourceRow `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode review resource query: %w", err)
	}
	return decoded.Rows, nil
}

type reviewResourceSignalSet struct {
	Files      []string
	Languages  []string
	Frameworks []string
	RiskTags   []string
	Categories []string
	Hints      []string
}

func reviewResourceSignals(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) reviewResourceSignalSet {
	files := changedFiles(ctx, repoRoot)
	if len(files) == 0 {
		files = facts.Files
	}
	if strings.TrimSpace(opts.Focus) != "" {
		var focused []string
		for _, file := range files {
			if inFocus(file, opts.Focus) {
				focused = append(focused, file)
			}
		}
		files = focused
	}
	files = limitStrings(files, 80)
	signals := reviewResourceSignalSet{
		Files:      files,
		Languages:  languageTagsForFiles(files),
		Frameworks: frameworkTagsForFiles(files, facts.DependencyFiles),
		RiskTags:   riskTagsForReview(files, facts.DependencyFiles, opts),
		Categories: categoriesForReviewScope(opts.Scope),
	}
	for _, hint := range hints {
		if title := strings.TrimSpace(hint.Title); title != "" {
			signals.Hints = append(signals.Hints, title)
		}
	}
	sort.Strings(signals.Hints)
	return signals
}

func reviewResourceQueryText(opts Options, signals reviewResourceSignalSet) string {
	parts := []string{
		"GX code review resource query",
		"scope: " + opts.Scope,
		"depth: " + depthLabel(opts.Deep),
		"focus: " + strings.TrimSpace(opts.Focus),
		"languages: " + strings.Join(signals.Languages, " "),
		"frameworks: " + strings.Join(signals.Frameworks, " "),
		"risk tags: " + strings.Join(signals.RiskTags, " "),
		"categories: " + strings.Join(signals.Categories, " "),
		"hints: " + strings.Join(signals.Hints, " "),
		"changed files: " + strings.Join(limitStrings(signals.Files, 30), " "),
	}
	return strings.Join(parts, "\n")
}

func reviewResourceBaseFilter() any {
	return []any{"And", []any{[]any{"source_kind", "Eq", "review_knowledge"}}}
}

func reviewResourceSignalFilter(signals reviewResourceSignalSet) any {
	var conditions []any
	if len(signals.Categories) > 0 {
		conditions = append(conditions, []any{"category", "In", signals.Categories})
	}
	if len(signals.Languages) > 0 {
		conditions = append(conditions, []any{"language_tags", "ContainsAny", signals.Languages})
	}
	if len(signals.Frameworks) > 0 {
		conditions = append(conditions, []any{"framework_tags", "ContainsAny", signals.Frameworks})
	}
	if len(signals.RiskTags) > 0 {
		conditions = append(conditions, []any{"risk_tag_values", "ContainsAny", signals.RiskTags})
		conditions = append(conditions, []any{"review_tag_values", "ContainsAny", signals.RiskTags})
	}
	if len(conditions) == 0 {
		return nil
	}
	return []any{"And", []any{
		[]any{"source_kind", "Eq", "review_knowledge"},
		[]any{"Or", conditions},
	}}
}

func reviewResourceSnippets(rows []reviewResourceRow, limit int, namespace string) []ContextSnippet {
	if limit <= 0 {
		limit = defaultReviewResourceTopK
	}
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		namespace = defaultReviewKnowledgeNamespace
	}
	seen := map[string]struct{}{}
	var snippets []ContextSnippet
	for _, row := range rows {
		text := stringValue(row["text"])
		if strings.TrimSpace(text) == "" {
			continue
		}
		sourceID := firstNonEmpty(stringValue(row["source_id"]), stringValue(row["url"]))
		chunkIndex := stringValue(row["chunk_index"])
		key := sourceID + "#" + chunkIndex
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		title := stringValue(row["title"])
		ref := sourceID
		if chunkIndex != "" {
			ref += "#" + chunkIndex
		}
		snippets = append(snippets, ContextSnippet{
			Kind:   "review_resource",
			Ref:    ref,
			Source: "turbopuffer:" + namespace,
			Text:   reviewResourceSnippetText(title, row, text),
		})
		if len(snippets) >= limit {
			break
		}
	}
	return snippets
}

func reviewResourceSnippetText(title string, row reviewResourceRow, text string) string {
	var b strings.Builder
	if title != "" {
		b.WriteString("Review resource: ")
		b.WriteString(title)
		b.WriteString("\n")
	}
	for _, key := range []string{"url", "category", "authority", "evidence_level", "language_tags", "framework_tags", "risk_tag_values", "review_tag_values"} {
		value := displayAttribute(row[key])
		if value == "" {
			continue
		}
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(text)
	return strings.TrimSpace(b.String())
}

func categoriesForReviewScope(scope string) []string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "security":
		return []string{"security", "supply-chain", "database"}
	case "dependencies":
		return []string{"supply-chain", "security", "tooling"}
	case "performance":
		return []string{"language", "framework", "database", "tooling"}
	case "testing":
		return []string{"core-process", "tooling"}
	case "docs", "onboarding", "architecture", "maintainability":
		return []string{"core-process", "language", "framework", "tooling"}
	default:
		return []string{"core-process", "language", "framework", "security", "database", "tooling"}
	}
}

func languageTagsForFiles(files []string) []string {
	tags := map[string]struct{}{}
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		base := strings.ToLower(filepath.Base(file))
		switch ext {
		case ".ts", ".tsx":
			tags["typescript"] = struct{}{}
			tags["javascript"] = struct{}{}
		case ".js", ".jsx", ".mjs", ".cjs":
			tags["javascript"] = struct{}{}
		case ".py":
			tags["python"] = struct{}{}
		case ".go":
			tags["go"] = struct{}{}
		case ".rs":
			tags["rust"] = struct{}{}
		case ".java":
			tags["java"] = struct{}{}
		case ".c":
			tags["c"] = struct{}{}
		case ".cc", ".cpp", ".cxx", ".hpp", ".hh", ".hxx":
			tags["cpp"] = struct{}{}
		case ".h":
			tags["c"] = struct{}{}
			tags["cpp"] = struct{}{}
		case ".cs":
			tags["csharp"] = struct{}{}
		case ".sql":
			tags["sql"] = struct{}{}
		case ".sh", ".bash", ".zsh", ".ksh":
			tags["shell"] = struct{}{}
		case ".php":
			tags["php"] = struct{}{}
		case ".kt", ".kts":
			tags["kotlin"] = struct{}{}
		case ".swift":
			tags["swift"] = struct{}{}
		}
		switch base {
		case "dockerfile", "makefile", "justfile":
			tags["shell"] = struct{}{}
		}
	}
	return sortedKeys(tags)
}

func frameworkTagsForFiles(files []string, dependencyFiles []string) []string {
	tags := map[string]struct{}{}
	for _, file := range append(append([]string{}, files...), dependencyFiles...) {
		lower := strings.ToLower(file)
		switch {
		case strings.Contains(lower, "spring"):
			tags["spring"] = struct{}{}
		case strings.Contains(lower, "android"):
			tags["android"] = struct{}{}
		case strings.Contains(lower, "ios/") || strings.Contains(lower, "macos/"):
			tags["ios"] = struct{}{}
		case strings.Contains(lower, "dbt"):
			tags["dbt"] = struct{}{}
		case strings.Contains(lower, "postgres") || strings.Contains(lower, "psql"):
			tags["postgresql"] = struct{}{}
		case strings.Contains(lower, "mysql"):
			tags["mysql"] = struct{}{}
		case strings.Contains(lower, "sqlserver") || strings.Contains(lower, "mssql"):
			tags["sql-server"] = struct{}{}
		}
	}
	return sortedKeys(tags)
}

func riskTagsForReview(files []string, dependencyFiles []string, opts Options) []string {
	tags := map[string]struct{}{}
	scope := strings.ToLower(strings.TrimSpace(opts.Scope))
	if scope != "" {
		tags[scope] = struct{}{}
	}
	for _, file := range files {
		lower := strings.ToLower(file)
		if isDependencyFile(file) {
			tags["dependencies"] = struct{}{}
			tags["supply-chain"] = struct{}{}
		}
		if strings.HasSuffix(lower, ".sql") || strings.Contains(lower, "migration") || strings.Contains(lower, "schema") || strings.Contains(lower, "db/") {
			tags["database-security"] = struct{}{}
			tags["data-correctness"] = struct{}{}
			tags["transactions"] = struct{}{}
		}
		if strings.Contains(lower, ".github/workflows/") || strings.Contains(lower, "scripts/") || strings.HasSuffix(lower, ".sh") {
			tags["ci-cd"] = struct{}{}
			tags["deployment"] = struct{}{}
			tags["command-injection"] = struct{}{}
		}
		if strings.Contains(lower, "auth") || strings.Contains(lower, "session") || strings.Contains(lower, "permission") || strings.Contains(lower, "security") {
			tags["authn"] = struct{}{}
			tags["authz"] = struct{}{}
			tags["secure-coding"] = struct{}{}
		}
	}
	for _, file := range dependencyFiles {
		if isDependencyFile(file) {
			tags["dependencies"] = struct{}{}
			tags["supply-chain"] = struct{}{}
		}
	}
	return sortedKeys(tags)
}

func sortedKeys(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func displayAttribute(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case []string:
		return strings.Join(typed, ", ")
	case []any:
		var out []string
		for _, item := range typed {
			if text := stringValue(item); text != "" {
				out = append(out, text)
			}
		}
		return strings.Join(out, ", ")
	default:
		return ""
	}
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	default:
		return ""
	}
}

func normalizeReviewOpenAIBaseURL(raw string) string {
	base := strings.TrimRight(strings.TrimSpace(raw), "/")
	if strings.HasSuffix(base, "/v1") {
		base = strings.TrimSuffix(base, "/v1")
	}
	if base == "" {
		return "https://api.openai.com"
	}
	return base
}

func reviewEnvInt(name string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
