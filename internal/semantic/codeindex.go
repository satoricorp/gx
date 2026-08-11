package semantic

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	// codeMaxFileBytes is the largest file that is chunked. Generated bundles
	// and vendored blobs above this are noise, but the limit is deliberately
	// generous: indexing cost is not a constraint, missing a real source file
	// is.
	codeMaxFileBytes = 512 * 1024
	// codeMaxFiles is a runaway guard, not a budget.
	codeMaxFiles = 50000
	// codeEmbedBatchSize / codeEmbedBatchBytes bound one embeddings request.
	codeEmbedBatchSize  = 32
	codeEmbedBatchBytes = 180_000
	// codeDefaultConcurrency is how many embed+upsert batches run at once.
	codeDefaultConcurrency = 8
)

var codeSkipDirNames = map[string]struct{}{
	".git":          {},
	"node_modules":  {},
	"dist":          {},
	"build":         {},
	".next":         {},
	"coverage":      {},
	"vendor":        {},
	".turbo":        {},
	".cache":        {},
	"target":        {},
	"__pycache__":   {},
	".venv":         {},
	"venv":          {},
	".idea":         {},
	".vscode":       {},
	".gx":           {},
	".terraform":    {},
	"Pods":          {},
	".gradle":       {},
	".mypy_cache":   {},
	".pytest_cache": {},
	// Fixture blobs, not source. "testdata" is the Go toolchain's own reserved
	// name and "__snapshots__"/"__fixtures__" are the JS equivalents; all three
	// hold golden payloads that exist to be compared against, not to be read.
	// Indexed as source they answer questions with a serialized expectation
	// instead of the code that produced it.
	"testdata":      {},
	"__snapshots__": {},
	"__fixtures__":  {},
}

var codeSkipExtensions = map[string]struct{}{
	".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".webp": {}, ".ico": {},
	".svg": {}, ".bmp": {}, ".tiff": {}, ".avif": {},
	".woff": {}, ".woff2": {}, ".ttf": {}, ".otf": {}, ".eot": {},
	".mp4": {}, ".mp3": {}, ".wav": {}, ".mov": {}, ".avi": {}, ".webm": {},
	".zip": {}, ".tar": {}, ".gz": {}, ".bz2": {}, ".xz": {}, ".7z": {}, ".rar": {},
	".pdf": {}, ".exe": {}, ".dll": {}, ".so": {}, ".dylib": {}, ".a": {}, ".o": {},
	".class": {}, ".jar": {}, ".wasm": {}, ".bin": {}, ".dat": {}, ".db": {},
	".sqlite": {}, ".sqlite3": {}, ".pyc": {}, ".pdb": {}, ".lock": {},
	".map": {}, ".ipynb": {},
}

var codeSkipBaseNames = map[string]struct{}{
	"package-lock.json": {},
	"yarn.lock":         {},
	"pnpm-lock.yaml":    {},
	"bun.lockb":         {},
	"go.sum":            {},
	"cargo.lock":        {},
	"poetry.lock":       {},
	"composer.lock":     {},
	"gemfile.lock":      {},
}

type repoFile struct {
	path string
	size int64
}

// CodeVectorStore is what the repository indexer needs from a vector database:
// upsert plus exact deletion by row id.
type CodeVectorStore interface {
	VectorStore
	DeleteRows(ctx context.Context, ids []string) error
}

// RepoIndexOptions describes one on-demand indexing run over a local checkout.
type RepoIndexOptions struct {
	// RepoRoot is the checkout to index. Required.
	RepoRoot string
	// RepoFullName is "owner/repo". Derived from the git remote when empty.
	RepoFullName string
	// OrgID scopes the namespace so two orgs never share rows. Empty means a
	// local-only namespace.
	OrgID string
	// Namespace overrides the computed namespace (tests, evaluation runs).
	Namespace string
	// CommitID / BranchName are recorded on every row. Derived from git when
	// empty. They are metadata only: chunk identity does not depend on them.
	CommitID   string
	BranchName string
	// Reason is recorded as indexed_reason ("review", "manual", "push").
	Reason string
	// Full ignores the manifest and re-embeds everything.
	Full bool
	// Concurrency is the number of parallel embed+upsert batches.
	Concurrency int

	// The following are injection points for tests and for callers that
	// already hold a configured client.
	Config    *Config
	Embedder  Embedder
	Store     CodeVectorStore
	StatePath string
	Now       func() time.Time
	Logf      func(format string, args ...any)
	// EmbedCache overrides the local vector cache. When nil the run resolves
	// the default cache from the environment — but only alongside a real
	// embedder: a test that injects an Embedder must not be silently coupled
	// to a persistent cache in the developer's $GX_HOME.
	EmbedCache *EmbedCache
}

// RepoIndexResult reports exactly what a run did. Every number is observable so
// an operator can tell "indexed nothing because nothing changed" apart from
// "indexed nothing because the walk found no files" — the previous
// implementation reported both as success with zero rows.
type RepoIndexResult struct {
	Namespace      string
	RepoFullName   string
	RepoRoot       string
	CommitID       string
	BranchName     string
	FilesScanned   int
	FilesIndexed   int
	FilesSkipped   int
	FilesRemoved   int
	ChunksTotal    int
	ChunksUpserted int
	// ChunksFromCache counts upserted chunks whose vector came from the local
	// embed cache instead of the embedding API.
	ChunksFromCache int
	ChunksDeleted   int
	BytesScanned    int64
	// BytesEmbedded and EmbedBatches count only what was actually sent to the
	// embedder; cache hits appear in neither.
	BytesEmbedded int64
	EmbedBatches  int
	Duration      time.Duration
	FullRebuild   bool
}

// UpToDate reports whether the run had nothing to do.
func (r RepoIndexResult) UpToDate() bool {
	return r.ChunksUpserted == 0 && r.ChunksDeleted == 0
}

func (r RepoIndexResult) String() string {
	return fmt.Sprintf(
		"namespace=%s files=%d indexed=%d skipped=%d removed=%d chunks=%d upserted=%d reused=%d deleted=%d bytes=%d batches=%d duration=%s",
		r.Namespace, r.FilesScanned, r.FilesIndexed, r.FilesSkipped, r.FilesRemoved,
		r.ChunksTotal, r.ChunksUpserted, r.ChunksFromCache, r.ChunksDeleted, r.BytesEmbedded, r.EmbedBatches,
		r.Duration.Round(time.Millisecond),
	)
}

// IndexRepository indexes a local checkout into TurboPuffer, incrementally.
//
// Unchanged files are neither re-chunked nor re-embedded: the manifest holds a
// content hash per file and per chunk, so steady-state runs upload nothing.
// Files that disappeared, and chunks that vanished because a file shrank, are
// deleted by exact row id.
func IndexRepository(ctx context.Context, opts RepoIndexOptions) (RepoIndexResult, error) {
	started := time.Now()
	nowFn := opts.Now
	if nowFn == nil {
		nowFn = time.Now
	}
	logf := opts.Logf
	if logf == nil {
		logf = func(string, ...any) {}
	}

	root := strings.TrimSpace(opts.RepoRoot)
	if root == "" {
		return RepoIndexResult{}, fmt.Errorf("repo root is required")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return RepoIndexResult{}, fmt.Errorf("resolve repo root: %w", err)
	}
	info, err := os.Stat(absRoot)
	if err != nil || !info.IsDir() {
		return RepoIndexResult{}, fmt.Errorf("repo root %s is not a directory", absRoot)
	}

	cfg, err := resolveCodeIndexConfig(opts)
	if err != nil {
		return RepoIndexResult{}, err
	}

	// Identity comes from the shared resolver, never from a local git call.
	// This is the write side of the namespace that review reads; the two used
	// to derive it separately and did not agree. See
	// internal/semantic/repoidentity.go.
	identity := ResolveRepoIdentity(ctx, absRoot, opts.OrgID, opts.RepoFullName)
	fullName := identity.RepoFullName
	commitID := strings.TrimSpace(opts.CommitID)
	if commitID == "" {
		commitID = gitOutput(absRoot, "rev-parse", "HEAD")
	}
	branch := strings.TrimSpace(opts.BranchName)
	if branch == "" {
		branch = gitOutput(absRoot, "rev-parse", "--abbrev-ref", "HEAD")
	}
	reason := strings.TrimSpace(opts.Reason)
	if reason == "" {
		reason = "manual"
	}
	namespace := strings.TrimSpace(opts.Namespace)
	if namespace == "" {
		namespace = identity.Namespace
	}

	store := opts.Store
	if store == nil {
		store = NewTurboPufferClientForNamespace(cfg, namespace)
	}
	embedder := opts.Embedder
	if embedder == nil {
		embedder = NewOpenAIEmbedder(cfg)
	}
	cache := opts.EmbedCache
	if cache == nil && opts.Embedder == nil {
		cache = OpenDefaultEmbedCache(cfg)
	}

	statePath := strings.TrimSpace(opts.StatePath)
	if statePath == "" {
		statePath, err = RepoIndexStatePath(namespace)
		if err != nil {
			return RepoIndexResult{}, fmt.Errorf("resolve index manifest path: %w", err)
		}
	}
	state := LoadRepoIndexState(statePath)
	fullRebuild := opts.Full || !state.Compatible(namespace, cfg.OpenAIEmbeddingModel, cfg.EmbeddingDimensions, CodeChunkerVersion, IndexSchemaVersion)
	// What the namespace is believed to hold right now, which is not the same
	// thing as what may be reused. A full rebuild re-embeds every file, but the
	// rows written by the previous run still exist, and the ones belonging to
	// files this run will not index — deleted files, and files newly excluded
	// as vendored or generated — are only reachable through the old manifest.
	// Resetting the manifest before the removal sweep, as this did, orphaned
	// them permanently: they stay in the namespace and keep being retrieved,
	// with no manifest left that admits they are there.
	previousFiles := map[string]FileIndexState{}
	if state != nil {
		for path, entry := range state.Files {
			previousFiles[path] = entry
		}
	}
	if fullRebuild {
		state = &RepoIndexState{Files: map[string]FileIndexState{}}
	}

	result := RepoIndexResult{
		Namespace:    namespace,
		RepoFullName: fullName,
		RepoRoot:     absRoot,
		CommitID:     commitID,
		BranchName:   branch,
		FullRebuild:  fullRebuild,
	}

	files, err := listRepositoryFiles(absRoot)
	if err != nil {
		return result, err
	}
	result.FilesScanned = len(files)
	logf("scanned %d files under %s", len(files), absRoot)

	rowCtx := codeRowContext{
		RepoRoot:     absRoot,
		RepoFullName: fullName,
		OrgID:        strings.TrimSpace(opts.OrgID),
		CommitID:     commitID,
		BranchName:   branch,
		Reason:       reason,
		IndexedAt:    nowFn().Unix(),
	}

	nextFiles := make(map[string]FileIndexState, len(files))
	var pending []Chunk
	var deleteIDs []string
	seen := map[string]struct{}{}

	for _, file := range files {
		body, err := os.ReadFile(filepath.Join(absRoot, filepath.FromSlash(file.path)))
		if err != nil {
			continue
		}
		if len(body) == 0 || isBinaryContent(body) {
			continue
		}
		seen[file.path] = struct{}{}
		result.BytesScanned += int64(len(body))
		fileHash := hashBytes(body)
		previous, hadPrevious := state.Files[file.path]
		// prior is what the namespace is believed to hold; previous is what may
		// be reused. They differ during a full rebuild, and only prior can say
		// how many rows a file used to own.
		prior, hadPrior := previousFiles[file.path]

		if hadPrevious && previous.Hash == fileHash && !fullRebuild {
			nextFiles[file.path] = previous
			result.FilesSkipped++
			result.ChunksTotal += len(previous.ChunkHashes)
			continue
		}

		chunks := ChunkSourceFile(file.path, string(body))
		if len(chunks) == 0 {
			// The file is indexable but produced nothing (whitespace only).
			// Drop any rows it used to own.
			if hadPrior {
				deleteIDs = append(deleteIDs, codeRowIDRange(fullName, file.path, 0, len(prior.ChunkHashes))...)
			}
			continue
		}
		result.FilesIndexed++
		result.ChunksTotal += len(chunks)

		hashes := make([]string, len(chunks))
		for index, chunk := range chunks {
			hashes[index] = chunk.ContentHash
			if hadPrevious && !fullRebuild && index < len(previous.ChunkHashes) &&
				previous.ChunkHashes[index] == chunk.ContentHash {
				// Same content at the same slot: the row in TurboPuffer is
				// already correct, so neither embed nor upload it again.
				continue
			}
			pending = append(pending, codeChunkRow(rowCtx, chunk, index))
		}
		if hadPrior && len(prior.ChunkHashes) > len(chunks) {
			deleteIDs = append(deleteIDs, codeRowIDRange(fullName, file.path, len(chunks), len(prior.ChunkHashes))...)
		}
		nextFiles[file.path] = FileIndexState{Hash: fileHash, Size: file.size, ChunkHashes: hashes}
	}

	removed := make([]string, 0, len(previousFiles))
	for path := range previousFiles {
		if _, ok := seen[path]; ok {
			continue
		}
		removed = append(removed, path)
	}
	sort.Strings(removed)
	for _, path := range removed {
		result.FilesRemoved++
		deleteIDs = append(deleteIDs, codeRowIDRange(fullName, path, 0, len(previousFiles[path].ChunkHashes))...)
	}

	if len(pending) == 0 && len(deleteIDs) == 0 {
		result.Duration = time.Since(started)
		// Persist the manifest anyway: it records the configuration the index
		// was verified against.
		applyIndexManifest(state, namespace, fullName, absRoot, cfg, commitID, nextFiles)
		if err := SaveRepoIndexState(statePath, state); err != nil {
			return result, err
		}
		logf("index up to date (%d files, %d chunks)", result.FilesScanned, result.ChunksTotal)
		return result, nil
	}

	batches := batchChunks(pending, codeEmbedBatchSize, codeEmbedBatchBytes)
	logf("embedding %d chunks in %d batches", len(pending), len(batches))

	// Track completion per file so a failed run can bank what finished. Before
	// this, state was saved only after every batch succeeded, so one dead
	// connection minutes into a large first index threw away the whole run —
	// the upserted rows persisted (row IDs are content-addressed and upserts
	// idempotent) but the manifest that makes the next run incremental did not,
	// and the next run re-embedded and re-paid for everything. Three attempts
	// at indexing one benchmark repo on a flaky link cost three full embeds and
	// produced nothing.
	pendingPerFile := map[string]int{}
	for _, chunk := range pending {
		if path, _ := chunk.Attributes[codeFieldFilePath].(string); path != "" {
			pendingPerFile[path]++
		}
	}
	var doneMu sync.Mutex
	donePerFile := map[string]int{}
	completedBatches := 0

	// bankPartialState persists the files whose pending chunks have all been
	// upserted. Called under doneMu.
	bankPartialState := func() int {
		banked := 0
		partial := make(map[string]FileIndexState, len(state.Files))
		for path, fileState := range state.Files {
			partial[path] = fileState
		}
		for path, next := range nextFiles {
			if count := pendingPerFile[path]; count == 0 || donePerFile[path] == count {
				partial[path] = next
				if count > 0 {
					banked++
				}
			}
		}
		snapshot := *state
		applyIndexManifest(&snapshot, namespace, fullName, absRoot, cfg, commitID, partial)
		if saveErr := SaveRepoIndexState(statePath, &snapshot); saveErr != nil {
			logf("could not save partial index state: %v", saveErr)
		}
		return banked
	}

	// bankEveryNBatches trades a little JSON writing for crash safety. The
	// error-path save below covers a run that fails and returns; it cannot
	// cover a run that is killed — a harness timeout, a closed laptop — and a
	// first index of a large repository runs long enough to meet both. Measured
	// before this: a one-hour sentry embed killed at a timeout banked nothing.
	const bankEveryNBatches = 25

	if err := runBatches(ctx, batches, opts.Concurrency, func(ctx context.Context, batch []Chunk) error {
		stats, err := embedAndUpsert(ctx, embedder, store, cache, batch)
		if err != nil {
			return err
		}
		doneMu.Lock()
		for _, chunk := range batch {
			if path, _ := chunk.Attributes[codeFieldFilePath].(string); path != "" {
				donePerFile[path]++
			}
		}
		result.ChunksFromCache += stats.reused
		result.BytesEmbedded += stats.embeddedBytes
		if stats.embedded > 0 {
			result.EmbedBatches++
		}
		completedBatches++
		if completedBatches%bankEveryNBatches == 0 {
			bankPartialState()
		}
		doneMu.Unlock()
		return nil
	}); err != nil {
		// Bank every file whose pending chunks all made it. Seeded from
		// state.Files — what the next run may reuse — not previousFiles, so an
		// interrupted full rebuild banks only what it actually rebuilt.
		doneMu.Lock()
		banked := bankPartialState()
		doneMu.Unlock()
		if banked > 0 {
			logf("index failed with %d file(s) completed; partial state saved, the next run resumes from it", banked)
		}
		result.Duration = time.Since(started)
		return result, err
	}
	result.ChunksUpserted = len(pending)

	if len(deleteIDs) > 0 {
		if err := store.DeleteRows(ctx, deleteIDs); err != nil {
			result.Duration = time.Since(started)
			return result, fmt.Errorf("delete stale code rows: %w", err)
		}
		result.ChunksDeleted = len(deleteIDs)
	}

	applyIndexManifest(state, namespace, fullName, absRoot, cfg, commitID, nextFiles)
	if err := SaveRepoIndexState(statePath, state); err != nil {
		result.Duration = time.Since(started)
		return result, err
	}

	// Only a run that wrote new cache entries can push the cache over its cap,
	// so up-to-date and all-hit runs skip the walk entirely.
	if result.BytesEmbedded > 0 {
		if removed := cache.Prune(); removed > 0 {
			logf("evicted %d cached vectors to stay under the cache cap", removed)
		}
	}

	result.Duration = time.Since(started)
	logf("%s", result.String())
	return result, nil
}

func resolveCodeIndexConfig(opts RepoIndexOptions) (Config, error) {
	if opts.Config != nil {
		return *opts.Config, nil
	}
	cfg := CodeIndexConfigFromEnv()
	if opts.Embedder != nil && opts.Store != nil {
		return cfg, nil
	}
	if cfg.OpenAIAPIKey == "" {
		return cfg, fmt.Errorf("OPENAI_API_KEY is required to index a repository")
	}
	if cfg.TurboPufferAPIKey == "" {
		return cfg, fmt.Errorf("TURBOPUFFER_API_KEY is required to index a repository")
	}
	return cfg, nil
}

type codeRowContext struct {
	RepoRoot     string
	RepoFullName string
	OrgID        string
	CommitID     string
	BranchName   string
	Reason       string
	IndexedAt    int64
}

func codeChunkRow(ctx codeRowContext, chunk CodeChunk, index int) Chunk {
	text := chunk.Text(ctx.RepoFullName)
	meta := CodeRowMetadata{
		SourceID:      codeSourceID(ctx.RepoFullName, chunk.FilePath, index),
		OrgID:         ctx.OrgID,
		RepoFullName:  ctx.RepoFullName,
		RepoRoot:      ctx.RepoRoot,
		BranchName:    ctx.BranchName,
		CommitID:      ctx.CommitID,
		FilePath:      chunk.FilePath,
		Symbol:        chunk.Symbol,
		SymbolText:    chunk.SymbolText(),
		SymbolKind:    chunk.SymbolKind,
		Package:       chunk.Package,
		StartLine:     chunk.StartLine,
		EndLine:       chunk.EndLine,
		Language:      chunk.Language,
		DocType:       chunk.DocType,
		ChunkHash:     chunk.ContentHash,
		IndexedAt:     ctx.IndexedAt,
		IndexedReason: ctx.Reason,
		Text:          text,
	}
	return Chunk{
		ID:         CodeRowID(ctx.RepoFullName, chunk.FilePath, index),
		Text:       text,
		Attributes: meta.Attributes(),
	}
}

func codeRowIDRange(repoFullName, filePath string, from, to int) []string {
	if to <= from {
		return nil
	}
	out := make([]string, 0, to-from)
	for index := from; index < to; index++ {
		out = append(out, CodeRowID(repoFullName, filePath, index))
	}
	return out
}

func batchChunks(chunks []Chunk, maxCount, maxBytes int) [][]Chunk {
	if len(chunks) == 0 {
		return nil
	}
	if maxCount <= 0 {
		maxCount = codeEmbedBatchSize
	}
	if maxBytes <= 0 {
		maxBytes = codeEmbedBatchBytes
	}
	var batches [][]Chunk
	var current []Chunk
	currentBytes := 0
	for _, chunk := range chunks {
		size := len(chunk.Text)
		if len(current) > 0 && (len(current) >= maxCount || currentBytes+size > maxBytes) {
			batches = append(batches, current)
			current = nil
			currentBytes = 0
		}
		current = append(current, chunk)
		currentBytes += size
	}
	if len(current) > 0 {
		batches = append(batches, current)
	}
	return batches
}

// applyIndexManifest stamps the configuration an index run was built against.
// Three save sites write these same fields — success, up-to-date, and the
// partial save on failure — and a field added to one but not the others would
// quietly produce states that disagree about their own compatibility.
func applyIndexManifest(state *RepoIndexState, namespace, fullName, absRoot string,
	cfg Config, commitID string, files map[string]FileIndexState) {
	state.Namespace = namespace
	state.RepoFullName = fullName
	state.RepoRoot = absRoot
	state.EmbeddingModel = cfg.OpenAIEmbeddingModel
	state.Dimensions = cfg.EmbeddingDimensions
	state.ChunkerVersion = CodeChunkerVersion
	state.SchemaVersion = IndexSchemaVersion
	state.CommitID = commitID
	state.Files = files
}

func runBatches(ctx context.Context, batches [][]Chunk, concurrency int, fn func(context.Context, []Chunk) error) error {
	if len(batches) == 0 {
		return nil
	}
	if concurrency <= 0 {
		concurrency = codeDefaultConcurrency
	}
	if concurrency > len(batches) {
		concurrency = len(batches)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	work := make(chan []Chunk)
	var mu sync.Mutex
	var firstErr error
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range work {
				if err := fn(ctx, batch); err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						cancel()
					}
					mu.Unlock()
					return
				}
			}
		}()
	}
	for _, batch := range batches {
		select {
		case work <- batch:
		case <-ctx.Done():
		}
	}
	close(work)
	wg.Wait()
	mu.Lock()
	defer mu.Unlock()
	return firstErr
}

type embedStats struct {
	embedded      int
	reused        int
	embeddedBytes int64
}

// embedAndUpsert uploads one batch, taking each vector from the local cache
// when the exact same text has been embedded before and from the embedder
// otherwise. The upsert itself is never skipped: row ids are positional, so
// when two checkouts of one repository alternate into a shared namespace, each
// sync must still rewrite the rows the other checkout overwrote — reusing the
// vector removes only the API call, which is the part that was being paid for
// repeatedly (and rate-limited) before the cache existed.
func embedAndUpsert(ctx context.Context, embedder Embedder, store CodeVectorStore, cache *EmbedCache, batch []Chunk) (embedStats, error) {
	var stats embedStats
	vectors := make([][]float32, len(batch))
	missIndexes := make([]int, 0, len(batch))
	missInputs := make([]string, 0, len(batch))
	for index, chunk := range batch {
		if vector := cache.Get(chunk.Text); vector != nil {
			vectors[index] = vector
			stats.reused++
			continue
		}
		missIndexes = append(missIndexes, index)
		missInputs = append(missInputs, chunk.Text)
	}
	if len(missInputs) > 0 {
		embeddings, err := embedder.Embed(ctx, missInputs)
		if err != nil {
			return stats, fmt.Errorf("embed code chunks: %w", err)
		}
		if len(embeddings) != len(missInputs) {
			return stats, fmt.Errorf("embedder returned %d embeddings for %d chunks", len(embeddings), len(missInputs))
		}
		for position, index := range missIndexes {
			vectors[index] = embeddings[position]
			// Cached before the upsert on purpose: a vector is valid the moment
			// the embedder returns it, and a run that dies on upload should
			// still leave its retry a warm cache.
			cache.Put(batch[index].Text, embeddings[position])
			stats.embeddedBytes += int64(len(missInputs[position]))
		}
		stats.embedded = len(missInputs)
	}
	rows := make([]VectorRow, 0, len(batch))
	for index, chunk := range batch {
		rows = append(rows, VectorRow{ID: chunk.ID, Vector: vectors[index], Attributes: chunk.Attributes})
	}
	return stats, store.Upsert(ctx, rows)
}

// listRepositoryFiles prefers `git ls-files` so .gitignore is honoured without
// reimplementing it, and falls back to a filtered walk outside a work tree.
func listRepositoryFiles(root string) ([]repoFile, error) {
	paths, ok := gitTrackedFiles(root)
	if !ok {
		return walkRepositoryFiles(root)
	}
	excluded := excludedByGitAttributes(root, paths)
	files := make([]repoFile, 0, len(paths))
	for _, rel := range paths {
		if !shouldIndexRepoPath(rel) {
			continue
		}
		if excluded[rel] {
			continue
		}
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || !info.Mode().IsRegular() || info.Size() > codeMaxFileBytes {
			continue
		}
		files = append(files, repoFile{path: rel, size: info.Size()})
	}
	sortRepoFiles(files)
	if len(files) > codeMaxFiles {
		files = files[:codeMaxFiles]
	}
	return files, nil
}

// indexExcludingAttributes are the .gitattributes markers that take a file out
// of the code index.
//
// These are GitHub Linguist's standard markers, so a repository declares
// "this is not our source" once and every tool agrees. It matters because a
// vendored corpus is not neutral filler: gx's own checkout carries 165
// third-party documents under scripts/review-knowledge/corpus (30% of its
// files), fetched from external standards bodies and, by that directory's own
// README, "repo-independent review knowledge". Indexed as if they were source,
// prose essays about code review outrank the Go files that implement it for
// every conceptual question — measured here, they cost 2 of 6 conceptual
// queries their top-5 answer. They are also already indexed in their own
// dedicated namespace, so their presence here is duplication as well as noise.
//
// linguist-documentation is deliberately NOT honoured. A repository's own docs
// describe its behaviour and are legitimate review context; only material that
// is not the repository's own work is excluded.
var indexExcludingAttributes = []string{"linguist-vendored", "linguist-generated"}

// excludedByGitAttributes reports which paths a repository has declared to be
// vendored or generated. A repository that declares nothing excludes nothing,
// so this can never silently shrink an index that did not ask for it.
func excludedByGitAttributes(root string, paths []string) map[string]bool {
	if len(paths) == 0 {
		return nil
	}
	args := append([]string{"-C", root, "check-attr", "-z", "--stdin"}, indexExcludingAttributes...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = strings.NewReader(strings.Join(paths, "\x00"))
	out, err := cmd.Output()
	if err != nil {
		// No .gitattributes, an old git, or a broken invocation: index
		// everything, which is the behaviour that existed before.
		return nil
	}
	// check-attr -z emits NUL-separated (path, attribute, value) triples.
	fields := strings.Split(string(out), "\x00")
	excluded := map[string]bool{}
	for index := 0; index+2 < len(fields); index += 3 {
		if fields[index+2] == "set" || fields[index+2] == "true" {
			excluded[filepath.ToSlash(fields[index])] = true
		}
	}
	return excluded
}

func gitTrackedFiles(root string) ([]string, bool) {
	cmd := exec.Command("git", "-C", root, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	out, err := cmd.Output()
	if err != nil {
		return nil, false
	}
	raw := strings.Split(string(out), "\x00")
	paths := make([]string, 0, len(raw))
	seen := map[string]struct{}{}
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		paths = append(paths, filepath.ToSlash(item))
	}
	return paths, true
}

func walkRepositoryFiles(root string) ([]repoFile, error) {
	var files []repoFile
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() {
			if _, skip := codeSkipDirNames[entry.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() || !shouldIndexRepoPath(rel) {
			return nil
		}
		info, infoErr := entry.Info()
		if infoErr != nil || info.Size() > codeMaxFileBytes {
			return nil
		}
		files = append(files, repoFile{path: rel, size: info.Size()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortRepoFiles(files)
	if len(files) > codeMaxFiles {
		files = files[:codeMaxFiles]
	}
	return files, nil
}

func sortRepoFiles(files []repoFile) {
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
}

func shouldIndexRepoPath(rel string) bool {
	if rel == "" {
		return false
	}
	for _, segment := range strings.Split(rel, "/") {
		if _, skip := codeSkipDirNames[segment]; skip {
			return false
		}
	}
	base := strings.ToLower(filepath.Base(rel))
	if _, skip := codeSkipBaseNames[base]; skip {
		return false
	}
	if strings.HasSuffix(base, ".min.js") || strings.HasSuffix(base, ".min.css") {
		return false
	}
	if _, skip := codeSkipExtensions[strings.ToLower(filepath.Ext(rel))]; skip {
		return false
	}
	return true
}

// isBinaryContent rejects files with a NUL byte in the first 8 KiB, which is
// the same heuristic git uses.
func isBinaryContent(body []byte) bool {
	limit := len(body)
	if limit > 8192 {
		limit = 8192
	}
	for _, b := range body[:limit] {
		if b == 0 {
			return true
		}
	}
	return false
}

func gitOutput(root string, args ...string) string {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
