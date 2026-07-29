package codereview

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/satoricorp/totality/internal/auth"
	"github.com/satoricorp/totality/internal/semantic"
)

const (
	// defaultCodeIndexTopK is deliberately the same number as
	// maxAIContextSnippets. The two are one budget seen from opposite ends: the
	// retriever decides how much code it is willing to find and the prompt
	// builder decides how much context it is willing to carry, and while the
	// retriever asked for 12 the prompt had room for 48, so the review was
	// starved by the smaller of two numbers that nobody had compared.
	//
	// Measured on this repository against 20 real changes, retrieving the
	// changed files' own chunks: k=12 recall 0.677, k=20 0.729, k=32 0.782,
	// k=48 0.807, k=64 0.819, k=96 0.851. The curve is still climbing at 96, so
	// this is not the recall ceiling — it is the point where extra rows stop
	// fitting the prompt and start being fetched only to be evicted.
	//
	// This is a real trade, not a free win, and the retrieval cost is the small
	// half of it. The extra turbopuffer work is 78ms. The whole review measured
	// slower end to end: median 82.7s at k=12 (n=7) against 111.4s at k=48
	// (n=6), on a subject whose per-run spread is 63-163s in both arms. The
	// mechanism was not isolated; the likely cause is that the shared context
	// now fills maxShardContextBytes instead of falling short of it, so the
	// model call carries a bigger prompt. It is taken deliberately, because
	// accuracy is the first-order goal here and a review that never sees the
	// code cannot find anything in it. k=32 keeps 81% of the recall gain and
	// was not separately timed, so it is the obvious knob if this trade is ever
	// judged the wrong way round.
	defaultCodeIndexTopK = 48
	// The deep and whole-repo modes buy their extra recall from the flat part of
	// that curve, where each additional row is worth less but still positive.
	defaultCodeIndexDeepTopK      = 96
	defaultCodeIndexWholeRepoTopK = 96
	// codeIndexEvidenceSource is the name this source reports under.
	codeIndexEvidenceSource = "code index"
	// codeIndexStaleAfter is when a code index stops being described as current.
	// A month-old index of a repository under active development answers with
	// deleted code, which is worse than answering with nothing, so the review
	// says how old the answer is.
	codeIndexStaleAfter = 14 * 24 * time.Hour
)

// CodeIndexRetriever retrieves the repository's own source from whichever
// TurboPuffer namespace holds it.
//
// Retrieval is hybrid on purpose. A review query has two very different shapes
// inside it: exact identifiers lifted out of the diff, which only a lexical
// index can match (a vector index has no idea that `contextRetrieverFromEnv` is
// one specific function), and a description of the change, which only a vector
// index can match against prose-free code. Running BM25 over the identifier
// column, BM25 over the body column and ANN over the vector, then fusing the
// three by reciprocal rank, gets both without hand-tuning a score blend.
type CodeIndexRetriever struct {
	Store      indexStore
	Namespaces []codeIndexTarget
	Limit      int
	// EmbedderFor builds an embedder for a namespace width. Injected so tests
	// can drive fusion without an embeddings API.
	EmbedderFor embedderFactory
	// CloudSearcher overrides the Totality Cloud retrieval client; injected in
	// tests. When nil it is resolved from the signed-in credentials.
	CloudSearcher reviewCloudSearcher
}

// codeIndexTarget is one candidate namespace and where it came from.
type codeIndexTarget struct {
	Namespace string
	// Origin explains the namespace to a human reading the evidence line.
	Origin string
}

// EvidenceSource implements evidenceNamer.
func (CodeIndexRetriever) EvidenceSource() string { return codeIndexEvidenceSource }

func codeIndexRetrieverFromEnv() ContextRetriever {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("TOTALITY_REVIEW_CODE_INDEX")), "0") {
		return nil
	}
	return CodeIndexRetriever{
		Limit: reviewEnvInt("TOTALITY_REVIEW_CODE_INDEX_TOP_K", defaultCodeIndexTopK),
	}
}

func (r CodeIndexRetriever) Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	// The signed-in path: retrieval through Totality Cloud, no provider keys on this
	// machine. The direct-store path below stays for development (raw
	// TURBOPUFFER_API_KEY) and for tests that inject a Store.
	if r.Store == nil && !rawTurboPufferKeyPresent() {
		searcher := r.CloudSearcher
		if searcher == nil {
			searcher = reviewCloudSearcherFromEnv()
		}
		if searcher != nil {
			return retrieveCodeIndexViaCloud(ctx, in, searcher, r.limitFor(in.Options))
		}
	}
	store, ok := r.resolveStore(in.Evidence)
	if !ok {
		return nil, nil
	}
	targets := r.Namespaces
	if len(targets) == 0 {
		// Review reads; it does not index. Keeping this repository's index
		// current is Totality Cloud's job, done from the GitHub App when a pull
		// request merges, which covers every user rather than only those
		// running the CLI with an embeddings key. Two writers filling one
		// namespace with two different chunkings would also have them
		// overwriting each other's rows for the same file on alternating runs.
		targets = codeIndexTargets(ctx, in.RepoRoot)
	}
	if len(targets) == 0 {
		in.Evidence.Record(EvidenceStatus{
			Source: codeIndexEvidenceSource,
			State:  EvidenceMissing,
			Detail: "no namespace could be resolved for this repository (no git remote and no org)",
			// No remedy names the website here on purpose: without a GitHub
			// remote there is no repository for Totality Cloud to connect to, so
			// pointing at /repositories would be an instruction that cannot be
			// followed.
			Remedy: "Add a GitHub remote, then connect this repository at https://totality.sh/repositories",
		})
		return nil, nil
	}

	probes := probeNamespaces(ctx, store, targets)
	queries := codeIndexQueryText(in)
	if queries.empty() {
		return nil, nil
	}

	limit := r.limitFor(in.Options)
	factory := r.EmbedderFor
	if factory == nil {
		factory = defaultEmbedderFactory
	}
	vectors := embedForWidths(ctx, codeIndexEvidenceSource, probes, queries.Body, factory, in.Evidence)

	type namespaceResult struct {
		target codeIndexTarget
		probe  indexProbe
		rows   []indexRow
		err    error
		// queried distinguishes "asked and got nothing" from "never asked",
		// which are different facts about how well informed the review is.
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
				SymbolField:       probe.SymbolField,
				SymbolQuery:       queries.Symbols,
				BodyField:         probe.BodyField,
				BodyQuery:         queries.Body,
				Limit:             limit,
				Filters:           codeIndexFilters(probe),
				IncludeAttributes: probe.IncludeAttributes(codeIndexAttributes()),
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
	// noteFor carries a per-namespace caveat into the snippet body. A stale
	// index serves code that has since been deleted, and a model told only "here
	// is the repository's source" will reason about it as current.
	noteFor := map[string]string{}
	for _, result := range results {
		if result.probe.Usable() && result.err == nil && len(result.rows) > 0 {
			lists = append(lists, result.rows)
			noteFor[result.probe.Namespace] = codeIndexFreshnessDetail(result.probe)
			for _, row := range result.rows {
				namespaceFor[codeIndexRowKey(row)] = result.probe.Namespace
			}
		}
	}
	var snippets []ContextSnippet
	if len(lists) > 0 {
		snippets = codeIndexSnippets(fuseRankedRows(lists, codeIndexRowKey), limit, namespaceFor, noteFor)
	}
	// Counts are taken after fusion and truncation so the evidence line reports
	// what this source actually put in front of the reviewer, not how many rows
	// it fetched and then dropped.
	contributed := map[string]int{}
	for _, snippet := range snippets {
		contributed[strings.TrimPrefix(snippet.Source, "turbopuffer:")]++
	}
	for _, result := range results {
		status := EvidenceStatus{
			Source:    codeIndexEvidenceSource,
			Namespace: result.target.Namespace,
		}
		switch {
		case result.probe.Err != nil:
			status.State = EvidenceUnavailable
			status.Detail = fmt.Sprintf("%s: %v", result.target.Origin, result.probe.Err)
		case !result.probe.Exists:
			status.State = EvidenceMissing
			status.Detail = codeIndexMissingDetail(result.target)
			status.Remedy = codeIndexMissingRemedy(result.target)
		case result.err != nil:
			status.State = EvidenceUnavailable
			status.Detail = result.err.Error()
		case !result.queried:
			// The namespace is there but nothing could be asked of it: no
			// full-text column this reader understands and no usable vector
			// width. Reporting that as "empty" would credit it with an answer it
			// never gave.
			status.State = EvidenceUnavailable
			status.Detail = "namespace exposes no queryable text or vector column this reviewer can use"
		case len(result.rows) == 0:
			status.State = EvidenceEmpty
			status.Detail = codeIndexFreshnessDetail(result.probe)
		default:
			status.State = EvidenceOK
			status.Snippets = contributed[result.target.Namespace]
			status.Detail = codeIndexFreshnessDetail(result.probe)
			if note := codeIndexVectorDetail(result.probe, vectors); note != "" {
				status.Detail = joinDetail(status.Detail, note)
			}
		}
		in.Evidence.Record(status)
	}
	return snippets, nil
}

func (r CodeIndexRetriever) resolveStore(log *EvidenceLog) (indexStore, bool) {
	if r.Store != nil {
		return r.Store, true
	}
	apiKey := strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY"))
	if apiKey == "" {
		log.Record(EvidenceStatus{
			Source: codeIndexEvidenceSource,
			State:  EvidenceDisabled,
			Detail: "not signed in to Totality Cloud, and no TURBOPUFFER_API_KEY for direct access",
			Remedy: signInRemedy,
		})
		return nil, false
	}
	return newTurboPufferIndexStore(apiKey, firstNonEmpty(os.Getenv("TOTALITY_TPUF_BASE_URL"), defaultReviewResourceBaseURL)), true
}

func (r CodeIndexRetriever) limitFor(opts Options) int {
	limit := r.Limit
	if limit <= 0 {
		limit = defaultCodeIndexTopK
	}
	if opts.Deep && limit < defaultCodeIndexDeepTopK {
		limit = defaultCodeIndexDeepTopK
	}
	if opts.WholeRepo && limit < defaultCodeIndexWholeRepoTopK {
		limit = defaultCodeIndexWholeRepoTopK
	}
	return limit
}

// codeIndexTargets resolves the namespaces that may hold this repository's
// source, most specific first.
//
// The names come from semantic.ResolveRepoIdentity — the same function
// `tl index` writes through — so the reader and the writer cannot address
// different namespaces. They used to derive the name separately and did not
// agree; see the comment at the top of internal/semantic/repoidentity.go for
// what that cost.
//
// Several namespaces are probed rather than one because a repository's identity
// is derived from mutable git config and from a second, independent indexer:
// the index may have been written before the checkout had a remote, and the
// console/Convex indexer writes repo-<owner>-<repo> keyed only on the GitHub
// remote. Results from every namespace that answers are fused, and the evidence
// line names which ones did.
func codeIndexTargets(ctx context.Context, repoRoot string) []codeIndexTarget {
	if override := strings.TrimSpace(os.Getenv("TOTALITY_REVIEW_CODE_INDEX_NAMESPACE")); override != "" {
		var out []codeIndexTarget
		for _, name := range strings.Split(override, ",") {
			if name = strings.TrimSpace(name); name != "" {
				out = append(out, codeIndexTarget{Namespace: name, Origin: "TOTALITY_REVIEW_CODE_INDEX_NAMESPACE"})
			}
		}
		return out
	}
	identity := semantic.ResolveRepoIdentity(ctx, repoRoot, reviewOrgID(), "")
	var out []codeIndexTarget
	for _, candidate := range identity.Candidates() {
		out = append(out, codeIndexTarget{Namespace: candidate.Namespace, Origin: candidate.Origin})
	}
	return out
}

// reviewOrgID is the signed-in Totality org, or "" when this machine has never
// logged in. Read in one place so the review's namespaces and the index refresh
// it triggers cannot disagree about which org they belong to.
func reviewOrgID() string {
	if creds, ok := auth.LoadUpload(); ok {
		return strings.TrimSpace(creds.OrgID)
	}
	return ""
}

// connectRepositoryRemedy is where a user goes to get their repository indexed.
//
// Every "no index" message ends here rather than at `tl index`. Indexing is Totality
// Cloud's job — it runs from the GitHub App on merge, so it stays current
// without anyone remembering to re-run anything — and `tl index` is a hidden
// maintenance command that indexes from one developer's checkout into one
// developer's namespace. Sending users to it would have them build, by hand, a
// worse copy of something the server maintains for them.
const connectRepositoryRemedy = "Connect this repository at https://totality.sh/repositories to have Totality Cloud index it"

// codeIndexMissingDetail explains an absent namespace.
func codeIndexMissingDetail(target codeIndexTarget) string {
	switch target.Origin {
	case semantic.NamespaceOriginPrimary:
		return "this repository has not been indexed"
	case semantic.NamespaceOriginPreRemote:
		return "no index under this repository's pre-remote name"
	default:
		return "namespace does not exist"
	}
}

// codeIndexMissingRemedy is what the reader can do about an absent namespace.
//
// A warning with no remedy is close to useless: it tells someone their review
// saw less than it could have and leaves them nowhere to go. This is the only
// place that gap is visible to a user, so it carries the fix.
//
// Every case points at the website rather than at `tl index`. Indexing is Totality
// Cloud's job — it runs from the GitHub App on merge, so it stays current
// without anyone remembering to re-run anything — and `tl index` is a hidden
// maintenance command that indexes one developer's checkout into one
// developer's namespace. Sending users there would have them build by hand a
// worse copy of something the server maintains for them.
//
// The pre-remote name is the exception. It is probed only because an index may
// predate this checkout's git remote, and connecting the repository will not
// create a namespace under a name the repository no longer uses, so there is
// nothing here for a user to act on.
func codeIndexMissingRemedy(target codeIndexTarget) string {
	if target.Origin == semantic.NamespaceOriginPreRemote {
		return ""
	}
	return connectRepositoryRemedy
}

func codeIndexFreshnessDetail(probe indexProbe) string {
	if probe.LastWrite.IsZero() {
		return ""
	}
	age := time.Since(probe.LastWrite)
	if age < codeIndexStaleAfter {
		return ""
	}
	return fmt.Sprintf("index last written %d day(s) ago; results may describe code that no longer exists", int(age.Hours()/24))
}

func codeIndexVectorDetail(probe indexProbe, vectors map[int][]float32) string {
	if probe.Dimensions <= 0 {
		return ""
	}
	if len(vectors[probe.Dimensions]) > 0 {
		return ""
	}
	return "semantic (vector) leg unavailable; lexical results only"
}

func joinDetail(parts ...string) string {
	var out []string
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, "; ")
}

// codeIndexFilters restricts a namespace to code rows when the namespace holds
// more than code. Filters are built from the probed schema: a filter naming an
// attribute the namespace does not declare is rejected outright, which would
// take the whole namespace out of the review.
func codeIndexFilters(probe indexProbe) any {
	if !probe.Has("source_kind") {
		return nil
	}
	return []any{"source_kind", "Eq", semantic.SourceKindCodeFile}
}

func codeIndexAttributes() []string {
	return []string{
		"content",
		"text",
		"file_path",
		"file",
		"symbol",
		"symbol_name",
		"symbol_kind",
		"package_name",
		"start_line",
		"end_line",
		"language",
		"doc_type",
		"commit_id",
		"head_sha",
		"branch",
		"branch_name",
		"repo_full_name",
		"repo_id",
		"chunk_hash",
		"source_kind",
	}
}

// codeIndexRowKey identifies a chunk across namespaces so the same code indexed
// twice fuses into one snippet instead of consuming two context slots.
func codeIndexRowKey(row indexRow) string {
	file := codeIndexRowFile(row)
	if file == "" {
		return strings.TrimSpace(stringValue(row["id"]))
	}
	return fmt.Sprintf("%s:%d-%d", file, intValue(row["start_line"]), intValue(row["end_line"]))
}

func codeIndexRowFile(row indexRow) string {
	return firstNonEmpty(stringValue(row["file_path"]), stringValue(row["file"]))
}

func codeIndexRowBody(row indexRow) string {
	for _, field := range bodyFieldPreference {
		if value := strings.TrimSpace(stringValue(row[field])); value != "" {
			return value
		}
	}
	return ""
}

func codeIndexSnippets(rows []indexRow, limit int, namespaceFor, noteFor map[string]string) []ContextSnippet {
	if limit <= 0 {
		limit = defaultCodeIndexTopK
	}
	var out []ContextSnippet
	for _, row := range rows {
		body := codeIndexRowBody(row)
		file := codeIndexRowFile(row)
		if body == "" || file == "" {
			continue
		}
		start := intValue(row["start_line"])
		end := intValue(row["end_line"])
		ref := file
		if start > 0 {
			ref = fmt.Sprintf("%s:%d", file, start)
		}
		symbol := firstNonEmpty(stringValue(row["symbol_name"]), stringValue(row["symbol"]))
		title := ref
		if symbol != "" {
			title = symbol + " (" + ref + ")"
		}
		namespace := namespaceFor[codeIndexRowKey(row)]
		out = append(out, ContextSnippet{
			Kind:      "indexed_code",
			Ref:       ref,
			Source:    "turbopuffer:" + namespace,
			Publisher: "indexed",
			Title:     title,
			Text:      codeIndexSnippetText(row, file, symbol, start, end, body, noteFor[namespace]),
			File:      file,
			StartLine: start,
			EndLine:   end,
			Commit:    firstNonEmpty(stringValue(row["commit_id"]), stringValue(row["head_sha"])),
			ChunkHash: stringValue(row["chunk_hash"]),
		})
		if len(out) >= limit {
			break
		}
	}
	return out
}

func codeIndexSnippetText(row indexRow, file, symbol string, start, end int, body, note string) string {
	var b strings.Builder
	b.WriteString("Indexed repository source.\n")
	if note = strings.TrimSpace(note); note != "" {
		fmt.Fprintf(&b, "caution: %s\n", note)
	}
	fmt.Fprintf(&b, "file: %s\n", file)
	if start > 0 {
		fmt.Fprintf(&b, "lines: %d-%d\n", start, end)
	}
	if symbol != "" {
		fmt.Fprintf(&b, "symbol: %s\n", symbol)
	}
	if language := stringValue(row["language"]); language != "" {
		fmt.Fprintf(&b, "language: %s\n", language)
	}
	if commit := firstNonEmpty(stringValue(row["commit_id"]), stringValue(row["head_sha"])); commit != "" {
		fmt.Fprintf(&b, "commit: %s\n", commit)
	}
	b.WriteString("\n")
	b.WriteString(body)
	return strings.TrimSpace(b.String())
}

// codeIndexQueries are the two textual queries a hybrid retrieval needs.
type codeIndexQueries struct {
	// Symbols are identifiers lifted from the diff, for the unstemmed symbol
	// column.
	Symbols string
	// Body is the natural-language description of the change, for the stemmed
	// body column and for the embedding.
	Body string
}

func (q codeIndexQueries) empty() bool {
	return strings.TrimSpace(q.Symbols) == "" && strings.TrimSpace(q.Body) == ""
}

// codeIndexQueryText builds the query from the change itself.
//
// It deliberately does not reuse reviewResourceQueryText. That text is written
// for the curated review corpus and is full of review vocabulary ("risk tags",
// "review intents", "failure modes"); pointed at a repository index it retrieves
// the repository's documentation about reviewing rather than the code under
// review — measured against this repository, half the top hits were markdown
// files about the review command. What retrieves code is the vocabulary of the
// code: the changed paths, the changed identifiers, and the added lines.
func codeIndexQueryText(in RetrieveInput) codeIndexQueries {
	signals := reviewResourceSignals(in)
	identifiers := changedIdentifiers(in.DiffSnippets, in.ChangedFiles, maxCodeIndexIdentifiers)
	parts := []string{
		"Source code related to the following change.",
		"files: " + strings.Join(limitStrings(signals.Files, 40), " "),
		"symbols: " + strings.Join(identifiers, " "),
		"languages: " + strings.Join(signals.Languages, " "),
	}
	if prompt := strings.TrimSpace(in.Options.Prompt); prompt != "" {
		parts = append(parts, "question: "+prompt)
	}
	if focus := strings.TrimSpace(in.Options.Focus); focus != "" {
		parts = append(parts, "focus: "+focus)
	}
	if added := addedDiffText(in.DiffSnippets, maxCodeIndexDiffTextBytes); added != "" {
		parts = append(parts, "changed lines:\n"+added)
	}
	return codeIndexQueries{
		Symbols: strings.Join(identifiers, " "),
		Body:    strings.Join(parts, "\n"),
	}
}

// addedDiffText is the added side of the diff, unprefixed, for the semantic
// leg: it is the closest thing available to "what the new code says".
func addedDiffText(diffs []DiffSnippet, limit int) string {
	var b strings.Builder
	for _, diff := range diffs {
		for _, line := range strings.Split(diff.Diff, "\n") {
			if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
				continue
			}
			line = strings.TrimSpace(line[1:])
			if line == "" {
				continue
			}
			if b.Len()+len(line)+1 > limit {
				return strings.TrimSpace(b.String())
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}

const (
	maxCodeIndexIdentifiers   = 48
	maxCodeIndexDiffTextBytes = 4000
)

var identifierPattern = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]{2,}`)

// changedIdentifiers extracts the identifiers a change touches, most frequent
// first, and expands compound names into their parts.
//
// The symbol column is indexed without stemming so that identifier lookups stay
// exact, which also means `parseReviewBrief` does not match `parse` or `brief`
// on its own. Emitting the parts alongside the whole name is what makes a
// renamed or partially-matching symbol still reachable.
func changedIdentifiers(diffs []DiffSnippet, changedFiles []string, limit int) []string {
	counts := map[string]int{}
	order := map[string]int{}
	next := 0
	add := func(token string) {
		token = strings.TrimSpace(token)
		if len(token) < 3 || len(token) > 64 {
			return
		}
		if identifierStopWords[strings.ToLower(token)] {
			return
		}
		if _, seen := counts[token]; !seen {
			order[token] = next
			next++
		}
		counts[token]++
	}
	for _, diff := range diffs {
		for _, line := range strings.Split(diff.Diff, "\n") {
			if len(line) == 0 {
				continue
			}
			switch line[0] {
			case '+', '-':
				if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
					continue
				}
			default:
				continue
			}
			for _, token := range identifierPattern.FindAllString(line, -1) {
				add(token)
				for _, part := range splitIdentifier(token) {
					add(part)
				}
			}
		}
	}
	for _, file := range normalizedChangedFiles(changedFiles) {
		base := path.Base(file)
		if ext := path.Ext(base); ext != "" {
			base = strings.TrimSuffix(base, ext)
		}
		add(base)
		for _, part := range splitIdentifier(base) {
			add(part)
		}
	}
	tokens := make([]string, 0, len(counts))
	for token := range counts {
		tokens = append(tokens, token)
	}
	sort.SliceStable(tokens, func(i, j int) bool {
		if counts[tokens[i]] != counts[tokens[j]] {
			return counts[tokens[i]] > counts[tokens[j]]
		}
		return order[tokens[i]] < order[tokens[j]]
	})
	if limit > 0 && len(tokens) > limit {
		tokens = tokens[:limit]
	}
	return tokens
}

var identifierBoundary = regexp.MustCompile(`[_\-.]+|([a-z0-9])([A-Z])`)

func splitIdentifier(token string) []string {
	spaced := identifierBoundary.ReplaceAllString(token, "$1 $2")
	fields := strings.Fields(strings.ReplaceAll(spaced, "_", " "))
	if len(fields) < 2 {
		return nil
	}
	return fields
}

// identifierStopWords are tokens that appear in every diff in a language and
// carry no retrieval signal. Keeping them would make every review's symbol
// query match every file.
var identifierStopWords = map[string]bool{
	"and": true, "any": true, "are": true, "bool": true, "break": true,
	"case": true, "chan": true, "class": true, "const": true, "continue": true,
	"def": true, "default": true, "defer": true, "else": true, "err": true,
	"error": true, "export": true, "false": true, "float": true, "for": true,
	"func": true, "function": true, "goto": true, "here": true, "iface": true,
	"import": true, "int": true, "interface": true, "let": true, "map": true,
	"nil": true, "not": true, "null": true, "package": true, "public": true,
	"range": true, "return": true, "self": true, "static": true, "string": true,
	"struct": true, "switch": true, "the": true, "this": true, "true": true,
	"type": true, "undefined": true, "var": true, "void": true, "while": true,
	"with": true, "yield": true,
}
