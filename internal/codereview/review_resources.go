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
	Files       []string
	Languages   []string
	Frameworks  []string
	RiskTags    []string
	Categories  []string
	Hints       []string
	Intents     []string
	PolicyQuery string
}

func reviewResourceSignals(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) reviewResourceSignalSet {
	files := reviewChangedFiles(ctx, repoRoot)
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
		Files:       files,
		Languages:   languageTagsForFiles(files),
		Frameworks:  frameworkTagsForFiles(files, facts.DependencyFiles),
		RiskTags:    riskTagsForReview(files, facts.DependencyFiles, opts),
		Categories:  categoriesForReview(opts),
		Intents:     reviewResourceIntents(opts),
		PolicyQuery: reviewPolicyQueryText(opts.ReviewPolicy),
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
		"profile: " + reviewProfile(opts),
		"scope: " + opts.Scope,
		"depth: " + depthLabel(opts.Deep),
		"focus: " + strings.TrimSpace(opts.Focus),
		"review prompt: " + strings.TrimSpace(opts.Prompt),
		"primary task: find concrete review rules and failure modes for the current code change",
		"review intents: " + strings.Join(signals.Intents, "; "),
		"languages: " + strings.Join(signals.Languages, " "),
		"frameworks: " + strings.Join(signals.Frameworks, " "),
		"risk tags: " + strings.Join(signals.RiskTags, " "),
		"categories: " + strings.Join(signals.Categories, " "),
		"hints: " + strings.Join(signals.Hints, " "),
		"changed files: " + strings.Join(limitStrings(signals.Files, 30), " "),
	}
	if strings.TrimSpace(signals.PolicyQuery) != "" {
		parts = append(parts, "review policy direction:\n"+strings.TrimSpace(signals.PolicyQuery))
	}
	return strings.Join(parts, "\n")
}

func reviewPolicyQueryText(policy *ReviewPolicy) string {
	if policy == nil {
		return ""
	}
	return policy.QueryText()
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
			Title:  title,
			URL:    stringValue(row["url"]),
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

func categoriesForReview(opts Options) []string {
	if opts.Deep {
		return []string{"core-process", "language", "framework", "security", "database", "tooling", "supply-chain", "performance", "observability"}
	}
	if strings.TrimSpace(opts.Prompt) != "" {
		return []string{"core-process", "language", "framework", "security", "database", "tooling", "supply-chain", "performance", "observability"}
	}
	if opts.PatchFocused {
		return []string{"core-process", "language", "framework", "security", "database", "tooling", "supply-chain"}
	}
	return categoriesForReviewScope(opts.Scope)
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

func reviewResourceIntents(opts Options) []string {
	intents := []string{
		"changed-line bug and regression checks",
		"security authn authz secret token and webhook verification checks",
		"data correctness migration transaction and partial failure checks",
		"race concurrency retry and idempotency checks",
		"error handling timeout cancellation and external API failure checks",
		"focused test coverage for changed behavior",
		"production observability logs errors metrics and status checks",
	}
	if opts.Deep {
		intents = append(intents,
			"Module Interface Depth locality and adapter architecture checks",
			"performance scalability and resource usage checks",
			"supply-chain dependency CI and deployment checks",
			"documentation onboarding and REVIEW.md policy checks",
		)
	} else if strings.TrimSpace(opts.Prompt) != "" {
		intents = append(intents,
			"user review_prompt checks",
			"broader codebase impact checks for the prompted concern",
		)
	} else if !opts.PatchFocused {
		intents = append(intents, "requested scope "+strings.TrimSpace(opts.Scope)+" checks")
	}
	return intents
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
	for _, tag := range []string{"bug", "regression", "error-handling", "testing"} {
		tags[tag] = struct{}{}
	}
	if opts.PatchFocused || opts.Deep {
		for _, tag := range []string{"secure-coding", "data-correctness", "race-condition", "idempotency", "observability"} {
			tags[tag] = struct{}{}
		}
	}
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

const (
	defaultIndexedContextLimit     = 8
	defaultIndexedDeepContextLimit = 24
)

type IndexedContextRetriever struct {
	Embedder  reviewResourceEmbedder
	Store     indexedContextStore
	Namespace string
	Limit     int
}

type indexedContextStore interface {
	Query(ctx context.Context, req indexedContextQuery) ([]indexedContextRow, error)
}

type indexedContextQuery struct {
	Vector            []float32
	Limit             int
	Filters           any
	IncludeAttributes []string
}

type indexedContextRow map[string]any

type turboPufferIndexedContextStore struct {
	apiKey     string
	baseURL    string
	namespace  string
	httpClient *http.Client
}

func indexedContextRetrieverFromEnv() ContextRetriever {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_INDEXED_CONTEXT")), "0") {
		return nil
	}
	cfg, err := semantic.ConfigFromEnv()
	if err != nil || !cfg.Enabled {
		return nil
	}
	return IndexedContextRetriever{
		Embedder: semantic.NewOpenAIEmbedder(cfg),
		Store: turboPufferIndexedContextStore{
			apiKey:     cfg.TurboPufferAPIKey,
			baseURL:    strings.TrimRight(cfg.TurboPufferBaseURL, "/"),
			namespace:  strings.Trim(cfg.TurboPufferNamespace, "/"),
			httpClient: &http.Client{Timeout: 20 * time.Second},
		},
		Namespace: cfg.TurboPufferNamespace,
		Limit:     reviewEnvInt("GX_REVIEW_INDEXED_CONTEXT_TOP_K", defaultIndexedContextLimit),
	}
}

func (r IndexedContextRetriever) Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	if r.Embedder == nil || r.Store == nil {
		return nil, nil
	}
	signals := reviewResourceSignals(ctx, repoRoot, opts, facts, hints)
	queryText := strings.Join([]string{
		"GX indexed codebase and session context query",
		reviewResourceQueryText(opts, signals),
		"Find code chunks and prior session transcript chunks that explain changed behavior, related Modules, previous agent decisions, and review risks.",
	}, "\n")
	vectors, err := r.Embedder.Embed(ctx, []string{queryText})
	if err != nil {
		return nil, err
	}
	if len(vectors) != 1 {
		return nil, fmt.Errorf("indexed context embedding returned %d vectors", len(vectors))
	}
	limit := r.Limit
	if opts.Deep && limit < defaultIndexedDeepContextLimit {
		limit = defaultIndexedDeepContextLimit
	}
	if limit <= 0 {
		limit = defaultIndexedContextLimit
	}
	rows, err := r.Store.Query(ctx, indexedContextQuery{
		Vector:            vectors[0],
		Limit:             limit,
		Filters:           indexedContextFilter(repoRoot),
		IncludeAttributes: indexedContextAttributes(),
	})
	if err != nil {
		return nil, err
	}
	return indexedContextSnippets(rows, limit, r.Namespace), nil
}

func indexedContextFilter(repoRoot string) any {
	conditions := []any{
		[]any{"source_kind", "In", []string{"code_file", "session_transcript"}},
	}
	if strings.TrimSpace(repoRoot) != "" {
		conditions = append(conditions, []any{"repo_root", "Eq", strings.TrimSpace(repoRoot)})
	}
	return []any{"And", conditions}
}

func indexedContextAttributes() []string {
	return []string{
		"text",
		"source_kind",
		"source_id",
		"repo_root",
		"repo_full_name",
		"branch_name",
		"commit_id",
		"revision_title",
		"file_path",
		"symbol",
		"start_line",
		"end_line",
		"chunk_hash",
		"language",
		"doc_type",
		"session_id",
		"request_id",
		"response_id",
		"agent_tool",
		"provider",
		"model",
		"created_at",
		"provenance_status",
	}
}

func (s turboPufferIndexedContextStore) Query(ctx context.Context, req indexedContextQuery) ([]indexedContextRow, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = defaultIndexedContextLimit
	}
	payload := map[string]any{
		"rank_by":            []any{"vector", "ANN", req.Vector},
		"limit":              limit,
		"include_attributes": req.IncludeAttributes,
	}
	if req.Filters != nil {
		payload["filters"] = req.Filters
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal indexed context query: %w", err)
	}
	endpoint := strings.TrimRight(s.baseURL, "/") + "/v2/namespaces/" + url.PathEscape(strings.Trim(s.namespace, "/")) + "/query"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create indexed context query: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	client := s.httpClient
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("query indexed context: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("indexed context is still building")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("query indexed context: status %s", resp.Status)
	}
	var decoded struct {
		Rows []indexedContextRow `json:"rows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode indexed context query: %w", err)
	}
	return decoded.Rows, nil
}

func indexedContextSnippets(rows []indexedContextRow, limit int, namespace string) []ContextSnippet {
	if limit <= 0 {
		limit = defaultIndexedContextLimit
	}
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		namespace = "gx-sessions"
	}
	seen := map[string]struct{}{}
	var snippets []ContextSnippet
	for _, row := range rows {
		snippet, ok := indexedContextSnippet(row, namespace)
		if !ok {
			continue
		}
		key := snippet.Kind + "\x00" + snippet.Ref + "\x00" + snippet.ChunkHash
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		snippets = append(snippets, snippet)
		if len(snippets) >= limit {
			break
		}
	}
	return snippets
}

func indexedContextSnippet(row indexedContextRow, namespace string) (ContextSnippet, bool) {
	text := stringValue(row["text"])
	if strings.TrimSpace(text) == "" {
		return ContextSnippet{}, false
	}
	sourceKind := stringValue(row["source_kind"])
	ref := firstNonEmpty(stringValue(row["source_id"]), stringValue(row["file_path"]), stringValue(row["session_id"]))
	if ref == "" {
		return ContextSnippet{}, false
	}
	snippet := ContextSnippet{
		Kind:       "indexed_context",
		Ref:        ref,
		Source:     "turbopuffer:" + strings.TrimSpace(namespace),
		Title:      firstNonEmpty(stringValue(row["revision_title"]), stringValue(row["symbol"]), ref),
		Text:       indexedContextSnippetText(row, text),
		File:       stringValue(row["file_path"]),
		StartLine:  intValue(row["start_line"]),
		EndLine:    intValue(row["end_line"]),
		Commit:     stringValue(row["commit_id"]),
		SessionID:  stringValue(row["session_id"]),
		RequestID:  stringValue(row["request_id"]),
		ResponseID: stringValue(row["response_id"]),
		ChunkHash:  stringValue(row["chunk_hash"]),
	}
	switch sourceKind {
	case "code_file":
		snippet.Kind = "indexed_code"
		if snippet.File != "" && snippet.StartLine > 0 {
			snippet.Ref = fmt.Sprintf("%s:%d", snippet.File, snippet.StartLine)
		}
	case "session_transcript":
		snippet.Kind = "indexed_session"
		if snippet.SessionID != "" {
			snippet.Ref = firstNonEmpty(snippet.SessionID+"/"+snippet.RequestID, snippet.SessionID)
		}
	default:
		return ContextSnippet{}, false
	}
	return snippet, true
}

func indexedContextSnippetText(row indexedContextRow, text string) string {
	var b strings.Builder
	for _, key := range []string{
		"source_kind",
		"repo_full_name",
		"repo_root",
		"branch_name",
		"commit_id",
		"revision_title",
		"file_path",
		"symbol",
		"start_line",
		"end_line",
		"chunk_hash",
		"session_id",
		"request_id",
		"response_id",
		"agent_tool",
		"provider",
		"model",
		"provenance_status",
	} {
		value := displayAttribute(row[key])
		if value == "" {
			continue
		}
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(value)
		b.WriteString("\n")
	}
	if b.Len() > 0 {
		b.WriteString("\n")
	}
	b.WriteString(text)
	return strings.TrimSpace(b.String())
}

func intValue(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case json.Number:
		n, _ := typed.Int64()
		return int(n)
	default:
		return 0
	}
}
