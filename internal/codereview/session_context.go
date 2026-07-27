package codereview

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/semantic"
)

const (
	defaultSessionContextTopK     = 6
	defaultSessionContextDeepTopK = 14
	sessionEvidenceSource         = "indexed sessions"
	// legacySessionNamespace is the global session namespace from before
	// namespaces were scoped per org and per repo. It is still queried because
	// it still holds rows, and it is filtered by repo because it does not scope
	// itself.
	legacySessionNamespace = "gx-sessions-dev"
)

// sessionSourceKinds are the row kinds that carry agent-session evidence. Three
// names exist because three writers produced them: the CLI's transcript
// indexer, the console publish intake, and the older session-context uploader.
var sessionSourceKinds = []string{"published_session_context", "session_transcript", "session_context"}

// SessionContextRetriever retrieves the agent-session evidence for this
// repository: what an agent was asked to do, what it decided, and what it
// changed, as captured on previous pushes.
//
// Nothing emitted session snippets into a review before this. The retriever
// that would have — IndexedContextRetriever — was gated behind an environment
// variable nobody sets and pointed at a namespace that was never created, so a
// review's account of "why was this written this way" came from the diff alone.
type SessionContextRetriever struct {
	Store       indexStore
	Namespaces  []codeIndexTarget
	Limit       int
	EmbedderFor embedderFactory
}

// EvidenceSource implements evidenceNamer.
func (SessionContextRetriever) EvidenceSource() string { return sessionEvidenceSource }

func sessionContextRetrieverFromEnv() ContextRetriever {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_SESSION_CONTEXT")), "0") {
		return nil
	}
	return SessionContextRetriever{
		Limit: reviewEnvInt("GX_REVIEW_SESSION_CONTEXT_TOP_K", defaultSessionContextTopK),
	}
}

func (r SessionContextRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	store := r.Store
	if store == nil {
		apiKey := strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY"))
		if apiKey == "" {
			in.Evidence.Record(EvidenceStatus{
				Source: sessionEvidenceSource,
				State:  EvidenceDisabled,
				Detail: "TURBOPUFFER_API_KEY is not set",
			})
			return nil, nil
		}
		store = newTurboPufferIndexStore(apiKey, firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), defaultReviewResourceBaseURL))
	}
	repoFullName := reviewHistoryRepoFullName(ctx, in.RepoRoot)
	targets := r.Namespaces
	if len(targets) == 0 {
		targets = sessionNamespaceTargets(in.RepoRoot, repoFullName)
	}
	if len(targets) == 0 {
		in.Evidence.Record(EvidenceStatus{
			Source: sessionEvidenceSource,
			State:  EvidenceMissing,
			Detail: "no session namespace could be resolved for this repository",
		})
		return nil, nil
	}

	probes := probeNamespaces(ctx, store, targets)
	queryText := sessionQueryText(in)
	if strings.TrimSpace(queryText) == "" {
		return nil, nil
	}
	limit := r.Limit
	if limit <= 0 {
		limit = defaultSessionContextTopK
	}
	if in.Options.Deep && limit < defaultSessionContextDeepTopK {
		limit = defaultSessionContextDeepTopK
	}
	factory := r.EmbedderFor
	if factory == nil {
		factory = defaultEmbedderFactory
	}
	vectors := embedForWidths(ctx, sessionEvidenceSource, probes, queryText, factory, in.Evidence)

	type namespaceResult struct {
		target  codeIndexTarget
		probe   indexProbe
		rows    []indexRow
		err     error
		queried bool
	}
	results := make([]namespaceResult, len(targets))
	var wg sync.WaitGroup
	for i, target := range targets {
		probe := probes[target.Namespace]
		results[i] = namespaceResult{target: target, probe: probe}
		if !probe.Usable() {
			continue
		}
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			query := indexQuery{
				Vector:            vectors[probe.Dimensions],
				BodyField:         probe.BodyField,
				BodyQuery:         queryText,
				Limit:             limit,
				Filters:           sessionFilters(probe, repoFullName),
				IncludeAttributes: probe.IncludeAttributes(sessionAttributes()),
			}
			if query.Legs() == 0 {
				return
			}
			results[i].queried = true
			rows, err := store.Query(ctx, probe.Namespace, query)
			results[i].rows = rows
			results[i].err = err
		}()
	}
	wg.Wait()

	var lists [][]indexRow
	namespaceFor := map[string]string{}
	for _, result := range results {
		if result.probe.Usable() && result.err == nil && len(result.rows) > 0 {
			lists = append(lists, result.rows)
			for _, row := range result.rows {
				namespaceFor[sessionRowKey(row)] = result.probe.Namespace
			}
		}
	}
	var snippets []ContextSnippet
	if len(lists) > 0 {
		snippets = sessionSnippets(fuseRankedRows(lists, sessionRowKey), limit, namespaceFor)
	}
	contributed := map[string]int{}
	for _, snippet := range snippets {
		contributed[strings.TrimPrefix(snippet.Source, "turbopuffer:")]++
	}
	for _, result := range results {
		status := EvidenceStatus{Source: sessionEvidenceSource, Namespace: result.target.Namespace}
		switch {
		case result.probe.Err != nil:
			status.State = EvidenceUnavailable
			status.Detail = result.probe.Err.Error()
		case !result.probe.Exists:
			status.State = EvidenceMissing
			status.Detail = "no captured sessions have been indexed here"
		case result.err != nil:
			status.State = EvidenceUnavailable
			status.Detail = result.err.Error()
		case !result.queried:
			status.State = EvidenceUnavailable
			status.Detail = "namespace exposes no queryable text or vector column this reviewer can use"
		case len(result.rows) == 0:
			status.State = EvidenceEmpty
		default:
			status.State = EvidenceOK
			status.Snippets = contributed[result.target.Namespace]
		}
		in.Evidence.Record(status)
	}
	return snippets, nil
}

// sessionNamespaceTargets resolves where this repository's session evidence
// could be. The per-org-per-repo namespace is the same one the code index uses
// — publish artifacts, session context and code chunks share it — and the
// legacy global namespace is queried as well because it still holds rows.
func sessionNamespaceTargets(repoRoot, repoFullName string) []codeIndexTarget {
	if override := strings.TrimSpace(os.Getenv("GX_REVIEW_SESSION_NAMESPACE")); override != "" {
		var out []codeIndexTarget
		for _, name := range strings.Split(override, ",") {
			if name = strings.TrimSpace(name); name != "" {
				out = append(out, codeIndexTarget{Namespace: name, Origin: "GX_REVIEW_SESSION_NAMESPACE"})
			}
		}
		return out
	}
	orgID := ""
	if creds, ok := auth.LoadUpload(); ok {
		orgID = strings.TrimSpace(creds.OrgID)
	}
	var out []codeIndexTarget
	seen := map[string]struct{}{}
	add := func(namespace, origin string) {
		namespace = strings.TrimSpace(namespace)
		if namespace == "" {
			return
		}
		if _, ok := seen[namespace]; ok {
			return
		}
		seen[namespace] = struct{}{}
		out = append(out, codeIndexTarget{Namespace: namespace, Origin: origin})
	}
	add(semantic.NamespaceForRepo(orgID, repoFullName, repoRoot), "org session namespace")
	add(legacySessionNamespace, "legacy session namespace")
	return out
}

// sessionFilters keeps a shared namespace's non-session rows out of the result
// and, in the un-scoped legacy namespace, keeps other repositories out.
func sessionFilters(probe indexProbe, repoFullName string) any {
	var conditions []any
	if probe.Has("source_kind") {
		conditions = append(conditions, []any{"source_kind", "In", sessionSourceKinds})
	}
	if repoFullName != "" && probe.Has("repo_full_name") {
		conditions = append(conditions, []any{"repo_full_name", "Eq", repoFullName})
	}
	switch len(conditions) {
	case 0:
		return nil
	case 1:
		return conditions[0]
	default:
		return []any{"And", conditions}
	}
}

func sessionAttributes() []string {
	return []string{
		"text",
		"content",
		"source_kind",
		"source_id",
		"session_id",
		"request_id",
		"response_id",
		"repo_full_name",
		"repo_root",
		"branch_name",
		"commit_id",
		"head_sha",
		"revision_title",
		"file",
		"file_path",
		"files",
		"model",
		"provider",
		"agent_tool",
		"created_at",
		"event_id",
	}
}

func sessionRowKey(row indexRow) string {
	if id := firstNonEmpty(stringValue(row["source_id"]), stringValue(row["id"])); id != "" {
		return id
	}
	return strings.Join([]string{
		stringValue(row["session_id"]),
		stringValue(row["request_id"]),
		stringValue(row["event_id"]),
	}, "\x00")
}

func sessionSnippets(rows []indexRow, limit int, namespaceFor map[string]string) []ContextSnippet {
	if limit <= 0 {
		limit = defaultSessionContextTopK
	}
	var out []ContextSnippet
	for _, row := range rows {
		text := firstNonEmpty(stringValue(row["text"]), stringValue(row["content"]))
		if strings.TrimSpace(text) == "" {
			continue
		}
		sessionID := stringValue(row["session_id"])
		requestID := stringValue(row["request_id"])
		ref := firstNonEmpty(sessionID, stringValue(row["source_id"]), stringValue(row["event_id"]))
		if ref == "" {
			continue
		}
		if sessionID != "" && requestID != "" {
			ref = sessionID + "/" + requestID
		}
		title := firstNonEmpty(stringValue(row["revision_title"]), "Agent session "+ref)
		out = append(out, ContextSnippet{
			Kind:       "indexed_session",
			Ref:        ref,
			Source:     "turbopuffer:" + namespaceFor[sessionRowKey(row)],
			Publisher:  "session",
			Title:      title,
			Text:       sessionSnippetText(row, text),
			File:       firstNonEmpty(stringValue(row["file"]), stringValue(row["file_path"])),
			Commit:     firstNonEmpty(stringValue(row["commit_id"]), stringValue(row["head_sha"])),
			SessionID:  sessionID,
			RequestID:  requestID,
			ResponseID: stringValue(row["response_id"]),
		})
		if len(out) >= limit {
			break
		}
	}
	return out
}

func sessionSnippetText(row indexRow, text string) string {
	var b strings.Builder
	b.WriteString("Captured agent session for this repository.\n")
	for _, key := range []string{
		"source_kind",
		"repo_full_name",
		"branch_name",
		"commit_id",
		"head_sha",
		"revision_title",
		"session_id",
		"request_id",
		"agent_tool",
		"provider",
		"model",
		"file",
		"files",
	} {
		if value := displayAttribute(row[key]); value != "" {
			fmt.Fprintf(&b, "%s: %s\n", key, value)
		}
	}
	b.WriteString("\n")
	b.WriteString(text)
	return strings.TrimSpace(b.String())
}

func sessionQueryText(in RetrieveInput) string {
	signals := reviewResourceSignals(in)
	identifiers := changedIdentifiers(in.DiffSnippets, in.ChangedFiles, 24)
	return strings.Join([]string{
		"GX code review: recall the agent sessions behind this code.",
		"Wanted: the instructions an agent was given, the decisions and trade-offs it recorded, constraints it was told to respect, and known-incomplete work in these files.",
		"changed files: " + strings.Join(limitStrings(signals.Files, 30), " "),
		"changed symbols: " + strings.Join(identifiers, " "),
		"review prompt: " + strings.TrimSpace(in.Options.Prompt),
		"scope: " + strings.TrimSpace(in.Options.Scope),
	}, "\n")
}
