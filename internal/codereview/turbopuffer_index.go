package codereview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/satoricorp/totality/internal/semantic"
)

// indexProbe describes a TurboPuffer namespace well enough to query it without
// guessing. Every reader in this repository has at some point assumed a field
// name or a vector width that the namespace did not have — the console code
// index is 1536-dimensional with a `content` column, the tl code index is
// 3072-dimensional with a `text` column, and the session namespaces are 512 —
// and the failures were silent: a wrong width is an HTTP 400 and a wrong field
// name is an empty result set. Reading the schema first turns both into facts.
type indexProbe struct {
	Namespace string
	// Exists is false when the namespace has never been created.
	Exists bool
	// Dimensions is the namespace's vector width, 0 when it has no vector.
	Dimensions int
	// BodyField is the full-text column holding the chunk body.
	BodyField string
	// SymbolField is the full-text column holding identifier terms, empty when
	// the namespace has none.
	SymbolField string
	// Filterable is the set of attributes that can appear in a filter. A filter
	// on an attribute the namespace does not declare is rejected outright, so
	// filters are built from this rather than from what the caller wishes were
	// there.
	Filterable map[string]bool
	// Fields is every attribute the namespace declares. TurboPuffer rejects the
	// whole query — HTTP 400, no rows — when include_attributes names one it
	// does not have, so the union of attribute names across the writers cannot
	// simply be requested from all of them.
	Fields    map[string]bool
	RowCount  int
	LastWrite time.Time
	// Err is set when the namespace could not be described at all.
	Err error
}

// IncludeAttributes narrows a wish list of attributes to the ones this
// namespace actually declares. An empty result means "ask for everything",
// which is what the caller wants when the schema could not be read.
func (p indexProbe) IncludeAttributes(candidates []string) []string {
	if len(p.Fields) == 0 {
		return nil
	}
	out := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if p.Fields[candidate] {
			out = append(out, candidate)
		}
	}
	return out
}

// Usable reports whether the namespace can answer a query.
func (p indexProbe) Usable() bool {
	return p.Exists && p.Err == nil && (p.BodyField != "" || p.SymbolField != "" || p.Dimensions > 0)
}

// Has reports whether an attribute may be used in a filter.
func (p indexProbe) Has(attribute string) bool {
	return p.Filterable[attribute]
}

// indexStore is the TurboPuffer surface these retrievers need. It is an
// interface so tests can drive the retrievers without a network.
type indexStore interface {
	Probe(ctx context.Context, namespace string) indexProbe
	Query(ctx context.Context, namespace string, req indexQuery) ([]indexRow, error)
}

// indexQuery is one hybrid query: up to three ranked legs fused by TurboPuffer
// with reciprocal rank fusion. Legs are optional individually — a namespace
// with no symbol column, or one whose vector width no embedder can match, still
// answers on the legs that remain rather than dropping out of the review.
type indexQuery struct {
	Vector            []float32
	SymbolField       string
	SymbolQuery       string
	BodyField         string
	BodyQuery         string
	Limit             int
	Filters           any
	IncludeAttributes []string
}

// Legs reports how many ranked legs this query will issue.
func (q indexQuery) Legs() int {
	legs := 0
	if len(q.Vector) > 0 {
		legs++
	}
	if strings.TrimSpace(q.SymbolField) != "" && strings.TrimSpace(q.SymbolQuery) != "" {
		legs++
	}
	if strings.TrimSpace(q.BodyField) != "" && strings.TrimSpace(q.BodyQuery) != "" {
		legs++
	}
	return legs
}

type indexRow map[string]any

type turboPufferIndexStore struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func newTurboPufferIndexStore(apiKey, baseURL string) turboPufferIndexStore {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = defaultReviewResourceBaseURL
	}
	return turboPufferIndexStore{
		apiKey:     strings.TrimSpace(apiKey),
		baseURL:    base,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

var vectorTypePattern = regexp.MustCompile(`^\[(\d+)\]`)

func (s turboPufferIndexStore) Probe(ctx context.Context, namespace string) indexProbe {
	probe := indexProbe{Namespace: namespace, Filterable: map[string]bool{}}
	namespace = strings.Trim(strings.TrimSpace(namespace), "/")
	if namespace == "" {
		probe.Err = fmt.Errorf("empty namespace")
		return probe
	}
	endpoint := s.baseURL + "/v2/namespaces/" + url.PathEscape(namespace) + "/metadata"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		probe.Err = err
		return probe
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	client := s.httpClient
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		probe.Err = err
		return probe
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		// Not an error: this repository has never been indexed here. The caller
		// reports it as missing evidence, which is a different and much more
		// actionable thing than a failed query.
		return probe
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		probe.Err = fmt.Errorf("status %s", resp.Status)
		return probe
	}
	var decoded struct {
		Schema         map[string]indexSchemaField `json:"schema"`
		ApproxRowCount int                         `json:"approx_row_count"`
		LastWriteAt    string                      `json:"last_write_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		probe.Err = fmt.Errorf("decode namespace metadata: %w", err)
		return probe
	}
	probe.Exists = true
	probe.RowCount = decoded.ApproxRowCount
	if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(decoded.LastWriteAt)); err == nil {
		probe.LastWrite = parsed
	}
	applyIndexSchema(&probe, decoded.Schema)
	return probe
}

type indexSchemaField struct {
	Type           string `json:"type"`
	Filterable     bool   `json:"filterable"`
	FullTextSearch any    `json:"full_text_search"`
}

// bodyFieldPreference is the order in which a full-text column is treated as
// "the chunk body". `content` is the console code index, `text` is every tl
// writer, `body` is the curated review corpus.
var bodyFieldPreference = []string{"content", "text", "body"}

func applyIndexSchema(probe *indexProbe, schema map[string]indexSchemaField) {
	if probe.Filterable == nil {
		probe.Filterable = map[string]bool{}
	}
	if probe.Fields == nil {
		probe.Fields = map[string]bool{}
	}
	fullText := map[string]bool{}
	for name, field := range schema {
		probe.Fields[name] = true
		if field.Filterable {
			probe.Filterable[name] = true
		}
		if field.FullTextSearch != nil {
			fullText[name] = true
		}
		if name == "vector" {
			if match := vectorTypePattern.FindStringSubmatch(field.Type); len(match) == 2 {
				probe.Dimensions, _ = strconv.Atoi(match[1])
			}
		}
	}
	if fullText["symbol"] {
		probe.SymbolField = "symbol"
	}
	for _, candidate := range bodyFieldPreference {
		if fullText[candidate] {
			probe.BodyField = candidate
			break
		}
	}
}

func (s turboPufferIndexStore) Query(ctx context.Context, namespace string, req indexQuery) ([]indexRow, error) {
	payload, err := indexQueryPayload(req)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal index query: %w", err)
	}
	endpoint := s.baseURL + "/v2/namespaces/" + url.PathEscape(strings.Trim(strings.TrimSpace(namespace), "/")) + "/query"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create index query: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	client := s.httpClient
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("index is still building")
	}
	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// The body carries the reason (a rejected attribute, a width mismatch).
		// Dropping it is what made these failures unreadable in the first place.
		detail := strings.TrimSpace(string(responseBody))
		if len(detail) > 300 {
			detail = detail[:300]
		}
		if detail == "" {
			return nil, fmt.Errorf("status %s", resp.Status)
		}
		return nil, fmt.Errorf("status %s: %s", resp.Status, detail)
	}
	if readErr != nil {
		return nil, readErr
	}
	var decoded struct {
		Rows    []indexRow `json:"rows"`
		Results []struct {
			Rows []indexRow `json:"rows"`
		} `json:"results"`
	}
	if err := json.Unmarshal(responseBody, &decoded); err != nil {
		return nil, fmt.Errorf("decode index query: %w", err)
	}
	if len(decoded.Rows) > 0 {
		return decoded.Rows, nil
	}
	var out []indexRow
	for _, result := range decoded.Results {
		out = append(out, result.Rows...)
	}
	return out, nil
}

// indexQueryPayload builds the request body. One leg is sent on its own; two or
// three are sent as a multi-query fused server-side with RRF, which costs one
// round trip instead of three and needs no score normalisation between a BM25
// score and a cosine distance.
func indexQueryPayload(req indexQuery) (map[string]any, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = defaultCodeIndexTopK
	}
	// An empty attribute list means the schema was not known well enough to
	// narrow it; asking for everything is correct there and, unlike naming an
	// attribute the namespace lacks, is never rejected.
	var include any = true
	if len(req.IncludeAttributes) > 0 {
		include = req.IncludeAttributes
	}
	leg := func(rankBy []any) map[string]any {
		query := map[string]any{
			"rank_by":            rankBy,
			"limit":              limit,
			"include_attributes": include,
		}
		if req.Filters != nil {
			query["filters"] = req.Filters
		}
		return query
	}
	var legs []any
	if symbolField := strings.TrimSpace(req.SymbolField); symbolField != "" && strings.TrimSpace(req.SymbolQuery) != "" {
		legs = append(legs, leg([]any{symbolField, "BM25", req.SymbolQuery}))
	}
	if bodyField := strings.TrimSpace(req.BodyField); bodyField != "" && strings.TrimSpace(req.BodyQuery) != "" {
		legs = append(legs, leg([]any{bodyField, "BM25", req.BodyQuery}))
	}
	if len(req.Vector) > 0 {
		legs = append(legs, leg([]any{"vector", "ANN", req.Vector}))
	}
	switch len(legs) {
	case 0:
		return nil, fmt.Errorf("index query has no ranked legs")
	case 1:
		return legs[0].(map[string]any), nil
	default:
		for _, entry := range legs {
			// A multi-query leg takes its limit as an object; the single-query
			// form takes a bare integer.
			entry.(map[string]any)["limit"] = map[string]any{"total": limit}
		}
		return map[string]any{
			"queries":   legs,
			"rerank_by": []any{"RRF"},
		}, nil
	}
}

// probeNamespaces describes every candidate namespace concurrently. A probe is
// one small GET, and doing them in parallel keeps namespace discovery off the
// review's critical path.
func probeNamespaces(ctx context.Context, store indexStore, targets []codeIndexTarget) map[string]indexProbe {
	out := map[string]indexProbe{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, target := range targets {
		target := target
		wg.Add(1)
		go func() {
			defer wg.Done()
			probe := store.Probe(ctx, target.Namespace)
			mu.Lock()
			out[target.Namespace] = probe
			mu.Unlock()
		}()
	}
	wg.Wait()
	return out
}

// embedderFactory resolves an embedder for one namespace vector width, together
// with a label for the evidence line. It returns false when no model can serve
// that width, which downgrades the namespace to lexical retrieval rather than
// dropping it.
type embedderFactory func(dimensions int) (reviewResourceEmbedder, string, bool)

// embedForWidths produces one query vector per distinct namespace width. Two
// namespaces at the same width share a single embedding request.
func embedForWidths(ctx context.Context, source string, probes map[string]indexProbe, text string, factory embedderFactory, log *EvidenceLog) map[int][]float32 {
	widths := map[int]struct{}{}
	for _, probe := range probes {
		if probe.Usable() && probe.Dimensions > 0 {
			widths[probe.Dimensions] = struct{}{}
		}
	}
	out := map[int][]float32{}
	if len(widths) == 0 || strings.TrimSpace(text) == "" || factory == nil {
		return out
	}
	ordered := make([]int, 0, len(widths))
	for width := range widths {
		ordered = append(ordered, width)
	}
	sort.Ints(ordered)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, width := range ordered {
		embedder, label, ok := factory(width)
		if !ok || embedder == nil {
			log.Record(EvidenceStatus{
				Source: source,
				State:  EvidenceDisabled,
				Detail: fmt.Sprintf("no embedding model configured for a %d-dimension index; using lexical retrieval only", width),
			})
			continue
		}
		width := width
		wg.Add(1)
		go func() {
			defer wg.Done()
			vectors, err := embedder.Embed(ctx, []string{text})
			if err != nil || len(vectors) != 1 {
				log.Record(EvidenceStatus{
					Source: source,
					State:  EvidenceUnavailable,
					Detail: fmt.Sprintf("embedding with %s failed: %s", label, embedErrorText(err, len(vectors))),
				})
				return
			}
			mu.Lock()
			out[width] = vectors[0]
			mu.Unlock()
		}()
	}
	wg.Wait()
	return out
}

func embedErrorText(err error, count int) string {
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("returned %d vectors", count)
}

// defaultEmbedderFactory builds an OpenAI embedder whose vectors live in the
// same space as the rows already in a namespace of that width.
//
// Matching the model is a correctness constraint, not a preference: embeddings
// from two different models are not comparable even at equal width, so querying
// a 1536-dimension index built with text-embedding-3-small using a truncated
// text-embedding-3-large vector returns confident nonsense rather than an
// error. The widths below are the ones tl and Totality Cloud have actually written:
// 3072 is the current tl code index, 1536 is the console/Convex code index, and
// 512 is the older publish and session namespaces.
func defaultEmbedderFactory(dimensions int) (reviewResourceEmbedder, string, bool) {
	model, ok := embedModelForDimensions(dimensions)
	if !ok {
		return nil, "", false
	}
	apiKey := strings.TrimSpace(firstNonEmpty(os.Getenv("OPENAI_API_KEY"), os.Getenv("TOTALITY_OPENAI_API_KEY")))
	if apiKey == "" {
		return nil, "", false
	}
	return semantic.NewOpenAIEmbedder(semantic.Config{
		OpenAIAPIKey:         apiKey,
		OpenAIBaseURL:        normalizeReviewOpenAIBaseURL(firstNonEmpty(os.Getenv("TOTALITY_OPENAI_BASE_URL"), os.Getenv("OPENAI_BASE_URL"), "https://api.openai.com")),
		OpenAIEmbeddingModel: model,
		EmbeddingDimensions:  dimensions,
	}), fmt.Sprintf("%s@%d", model, dimensions), true
}

func embedModelForDimensions(dimensions int) (string, bool) {
	if dimensions <= 0 {
		return "", false
	}
	if override := strings.TrimSpace(os.Getenv("TOTALITY_REVIEW_INDEX_EMBED_MODEL")); override != "" {
		return override, true
	}
	configured := semantic.CodeIndexConfigFromEnv()
	if dimensions == configured.EmbeddingDimensions && strings.TrimSpace(configured.OpenAIEmbeddingModel) != "" {
		return configured.OpenAIEmbeddingModel, true
	}
	switch dimensions {
	case 3072:
		return "text-embedding-3-large", true
	case 1536, 512:
		return "text-embedding-3-small", true
	default:
		return "", false
	}
}

// rrfK is the reciprocal-rank-fusion constant. 60 is the value from the
// original Cormack et al. formulation and the one TurboPuffer uses server-side,
// so cross-namespace fusion here ranks on the same curve as within-namespace
// fusion there.
const rrfK = 60

// fuseRankedRows merges several ranked result lists into one by reciprocal rank
// fusion, keyed by a caller-supplied identity. Used to combine namespaces —
// TurboPuffer fuses legs within a namespace, nothing fuses across them.
func fuseRankedRows(lists [][]indexRow, key func(indexRow) string) []indexRow {
	type entry struct {
		row   indexRow
		score float64
		order int
	}
	byKey := map[string]*entry{}
	var order []string
	for _, list := range lists {
		for rank, row := range list {
			id := key(row)
			if id == "" {
				continue
			}
			score := 1.0 / float64(rrfK+rank+1)
			existing, ok := byKey[id]
			if !ok {
				byKey[id] = &entry{row: row, score: score, order: len(order)}
				order = append(order, id)
				continue
			}
			existing.score += score
			// Keep whichever copy carries more body text: the same chunk can be
			// indexed with a different truncation in each namespace.
			if indexRowBodyLength(row) > indexRowBodyLength(existing.row) {
				existing.row = row
			}
		}
	}
	entries := make([]*entry, 0, len(order))
	for _, id := range order {
		entries = append(entries, byKey[id])
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].score != entries[j].score {
			return entries[i].score > entries[j].score
		}
		return entries[i].order < entries[j].order
	})
	out := make([]indexRow, 0, len(entries))
	for _, item := range entries {
		out = append(out, item.row)
	}
	return out
}

func indexRowBodyLength(row indexRow) int {
	best := 0
	for _, field := range bodyFieldPreference {
		if length := len(stringValue(row[field])); length > best {
			best = length
		}
	}
	return best
}
