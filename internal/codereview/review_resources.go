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

	"github.com/satoricorp/gx/internal/cloud"
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
	// CloudSearcher overrides the gx Cloud retrieval client; injected in
	// tests. When nil it is resolved from the signed-in credentials.
	CloudSearcher reviewCloudSearcher
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
	retriever := ReviewResourceRetriever{
		Limit: reviewEnvInt("GX_REVIEW_RESOURCES_TOP_K", defaultReviewResourceTopK),
	}
	openAIKey := strings.TrimSpace(firstNonEmpty(os.Getenv("OPENAI_API_KEY"), os.Getenv("GX_OPENAI_API_KEY")))
	tpufKey := strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY"))
	if openAIKey == "" || tpufKey == "" {
		// No direct keys. The retriever is still returned: Retrieve goes
		// through gx Cloud when signed in, and otherwise records the source as
		// disabled. It used to return nil here, which made the shared
		// knowledge corpus vanish from the review with no evidence line — the
		// exact silent-degradation the evidence log exists to prevent.
		return retriever
	}
	namespace := strings.TrimSpace(os.Getenv("GX_REVIEW_KNOWLEDGE_NAMESPACE"))
	if namespace == "" {
		namespace = defaultReviewKnowledgeNamespace
	}
	cfg := semantic.Config{
		OpenAIAPIKey:         openAIKey,
		OpenAIBaseURL:        normalizeReviewOpenAIBaseURL(firstNonEmpty(os.Getenv("GX_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL"), "https://api.openai.com")),
		OpenAIEmbeddingModel: firstNonEmpty(os.Getenv("GX_OPENAI_EMBEDDING_MODEL"), "text-embedding-3-small"),
		// 512 stays pinned for the shared corpus: its ~35k rows were embedded
		// at that width and TurboPuffer rejects a query at any other.
		EmbeddingDimensions: reviewEnvInt("GX_EMBEDDING_DIMENSIONS", 512),
	}
	retriever.Embedder = semantic.NewOpenAIEmbedder(cfg)
	retriever.Store = turboPufferReviewResourceStore{
		apiKey:     tpufKey,
		baseURL:    strings.TrimRight(firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), defaultReviewResourceBaseURL), "/"),
		namespace:  strings.Trim(namespace, "/"),
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
	retriever.Namespace = namespace
	return retriever
}

// reviewKnowledgeEvidenceSource is the name this source reports under.
const reviewKnowledgeEvidenceSource = "review knowledge"

// EvidenceSource implements evidenceNamer.
func (ReviewResourceRetriever) EvidenceSource() string { return reviewKnowledgeEvidenceSource }

// retrieveKnowledgeViaCloud searches the shared review-knowledge corpus
// through gx Cloud. The server embeds the query at the corpus's own width and
// applies the review_corpus filter; the tag-narrowed second query of the
// direct path is not replicated — the broad hybrid search is what the server's
// own summary broker uses for this corpus.
func (r ReviewResourceRetriever) retrieveKnowledgeViaCloud(
	ctx context.Context,
	in RetrieveInput,
	searcher reviewCloudSearcher,
	queryText string,
) ([]ContextSnippet, error) {
	limit := r.Limit
	if in.Options.Deep && limit < defaultReviewResourceDeepTopK {
		limit = defaultReviewResourceDeepTopK
	}
	if limit <= 0 {
		limit = defaultReviewResourceTopK
	}
	result, err := searcher.SearchReviewIndex(ctx, cloud.ReviewSearchRequest{
		Target: "knowledge",
		Query:  queryText,
		Limit:  limit,
	})
	if err != nil {
		in.Evidence.Record(EvidenceStatus{
			Source: reviewKnowledgeEvidenceSource,
			State:  EvidenceUnavailable,
			Detail: cloudUnavailableDetail(err),
		})
		return nil, nil
	}
	status := EvidenceStatus{Source: reviewKnowledgeEvidenceSource, Namespace: result.Namespace}
	switch {
	case !result.Available:
		status.State = EvidenceUnavailable
		status.Detail = "gx cloud retrieval is not configured server-side"
		if result.Reason != "" {
			status.Detail = "gx cloud: " + result.Reason
		}
	case !result.Exists:
		status.State = EvidenceMissing
		status.Detail = "the shared review-knowledge corpus does not exist on this cloud"
	case len(result.Rows) == 0:
		status.State = EvidenceEmpty
	}
	if status.State != "" {
		in.Evidence.Record(status)
		return nil, nil
	}
	rows := make([]reviewResourceRow, 0, len(result.Rows))
	for _, row := range cloudSearchRows(result.Rows) {
		rows = append(rows, reviewResourceRow(row))
	}
	snippets := reviewResourceSnippets(rows, limit, result.Namespace)
	status.State = EvidenceOK
	status.Snippets = len(snippets)
	status.Detail = "via gx cloud"
	in.Evidence.Record(status)
	return snippets, nil
}

func (r ReviewResourceRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	if !in.Plan.RunReviewResources && in.Plan.Triage.Class != "" {
		return nil, nil
	}
	opts := in.Options
	signals := reviewResourceSignals(in)
	queryText := reviewResourceQueryText(opts, signals)
	if strings.TrimSpace(queryText) == "" {
		return nil, nil
	}
	if r.Embedder == nil || r.Store == nil {
		// No direct keys: the signed-in path retrieves the shared corpus
		// through gx Cloud, and with no login the source reports itself
		// disabled instead of silently contributing nothing.
		searcher := r.CloudSearcher
		if searcher == nil {
			searcher = reviewCloudSearcherFromEnv()
		}
		if searcher == nil {
			in.Evidence.Record(EvidenceStatus{
				Source: reviewKnowledgeEvidenceSource,
				State:  EvidenceDisabled,
				Detail: "not signed in to gx Cloud, and no direct keys for the shared review corpus",
				Remedy: signInRemedy,
			})
			return nil, nil
		}
		return r.retrieveKnowledgeViaCloud(ctx, in, searcher, queryText)
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
	// The attribute list records what this retriever reads. It is not sent as
	// include_attributes — see turboPufferReviewResourceStore.Query.
	include := []string{
		"body",
		"text",
		"source_id",
		"publisher",
		"url",
		"source_url",
		"title",
		"category",
		"tier",
		"authority",
		"evidence_level",
		"language_tags",
		"languages",
		"framework_tags",
		"risk_tag_values",
		"review_tag_values",
		"precedence_group",
		"superseded_by",
		"historical",
		"section_path",
		"cwe_ids",
		"content_sha256",
		"usefulness_rank",
		"chunk_id",
		"chunk_kind",
		"chunk_index",
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
	// include_attributes is deliberately `true` rather than the caller's list.
	// TurboPuffer rejects the entire query — HTTP 400, zero rows — when the list
	// names an attribute the namespace does not declare, and the list did:
	// gx-review-knowledge has no `publisher` column, so every review resource
	// query in production was failing and, because the composite retriever
	// swallowed the error, failing invisibly. Asking for all attributes cannot
	// be rejected. The caller's list is kept in the request as a record of what
	// the retriever reads.
	vectorQuery := map[string]any{
		"rank_by":            []any{"vector", "ANN", req.Vector},
		"limit":              map[string]any{"total": limit},
		"include_attributes": true,
	}
	if req.Filters != nil {
		vectorQuery["filters"] = req.Filters
	}
	payload := vectorQuery
	if strings.TrimSpace(req.Text) != "" {
		textQuery := map[string]any{
			"rank_by":            []any{"body", "BM25", req.Text},
			"limit":              map[string]any{"total": limit},
			"include_attributes": true,
		}
		if req.Filters != nil {
			textQuery["filters"] = req.Filters
		}
		payload = map[string]any{
			"queries":   []any{vectorQuery, textQuery},
			"rerank_by": []any{"RRF"},
		}
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
		Rows    []reviewResourceRow `json:"rows"`
		Results []struct {
			Rows []reviewResourceRow `json:"rows"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode review resource query: %w", err)
	}
	if len(decoded.Rows) == 0 && len(decoded.Results) > 0 {
		return decoded.Results[0].Rows, nil
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

func reviewResourceSignals(in RetrieveInput) reviewResourceSignalSet {
	opts := in.Options
	facts := in.Facts
	files := normalizedChangedFiles(in.ChangedFiles)
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
		RiskTags:    in.Plan.RiskTags,
		Categories:  categoriesForReview(opts),
		Intents:     reviewResourceIntents(opts),
		PolicyQuery: reviewPolicyQueryText(opts.ReviewPolicy),
	}
	if len(signals.RiskTags) == 0 {
		signals.RiskTags = riskTagsForReview(files, facts.DependencyFiles, opts)
	}
	for _, hint := range in.Hints {
		if title := strings.TrimSpace(hint.Title); title != "" {
			signals.Hints = append(signals.Hints, title)
		}
	}
	sort.Strings(signals.Hints)
	return signals
}

func reviewResourceQueryText(opts Options, signals reviewResourceSignalSet) string {
	parts := []string{
		"gx code review resource query",
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
	return []any{"And", []any{
		[]any{"source_kind", "Eq", "review_corpus"},
		[]any{"tier", "NotEq", "research"},
		[]any{"historical", "Eq", false},
		[]any{"superseded_by", "Eq", ""},
	}}
}

func reviewResourceSignalFilter(signals reviewResourceSignalSet) any {
	var conditions []any
	if len(signals.Categories) > 0 {
		conditions = append(conditions, []any{"tier", "In", signals.Categories})
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
		[]any{"source_kind", "Eq", "review_corpus"},
		[]any{"tier", "NotEq", "research"},
		[]any{"historical", "Eq", false},
		[]any{"superseded_by", "Eq", ""},
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
		text := firstNonEmpty(stringValue(row["body"]), stringValue(row["text"]))
		if strings.TrimSpace(text) == "" {
			continue
		}
		sourceID := firstNonEmpty(stringValue(row["source_id"]), stringValue(row["url"]))
		chunkIndex := firstNonEmpty(stringValue(row["chunk_id"]), stringValue(row["chunk_index"]))
		key := sourceID + "#" + chunkIndex
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		title := stringValue(row["title"])
		publisher := firstNonEmpty(stringValue(row["publisher"]), publisherFromURLHost(firstNonEmpty(stringValue(row["source_url"]), stringValue(row["url"]))))
		ref := sourceID
		if chunkIndex != "" {
			ref += "#" + chunkIndex
		}
		snippets = append(snippets, ContextSnippet{
			Kind:      "review_resource",
			Ref:       ref,
			Source:    "turbopuffer:" + namespace,
			Publisher: publisher,
			Title:     title,
			URL:       firstNonEmpty(stringValue(row["source_url"]), stringValue(row["url"])),
			Text:      reviewResourceSnippetText(title, row, text),
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
	for _, key := range []string{"publisher", "source_url", "url", "tier", "authority", "precedence_group", "section_path", "language_tags", "cwe_ids", "usefulness_rank"} {
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
		return []string{"process", "language", "style", "security"}
	}
	if strings.TrimSpace(opts.Prompt) != "" {
		return []string{"process", "language", "style", "security"}
	}
	if opts.PatchFocused {
		return []string{"process", "language", "style", "security"}
	}
	return categoriesForReviewScope(opts.Scope)
}

func categoriesForReviewScope(scope string) []string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "security":
		return []string{"security"}
	case "dependencies":
		return []string{"security", "style"}
	case "performance":
		return []string{"language", "style"}
	case "testing":
		return []string{"process", "style"}
	case "docs", "onboarding", "architecture", "maintainability":
		return []string{"process", "language", "style"}
	default:
		return []string{"process", "language", "style", "security"}
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
	// WholeRepo belongs with the other two: these tags drive review-resource
	// retrieval, and a whole-repo review that retrieved less security guidance
	// than the default patch review would be a narrower review bought with a
	// broader request.
	if opts.PatchFocused || opts.Deep || opts.WholeRepo {
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

func intValue(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
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
