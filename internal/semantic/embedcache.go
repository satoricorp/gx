package semantic

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// defaultEmbedCacheMB bounds the on-disk vector cache. At 1536 dimensions a
// vector is 6 KiB, so the default holds roughly 85k chunks — an order of
// magnitude more than the largest between-checkout delta observed in practice
// (a few thousand chunks per benchmark run).
const defaultEmbedCacheMB = 512

// EmbedCache is a local content-addressed store of embedding vectors, keyed by
// the exact text sent to the embedder.
//
// It exists because the sync manifest tracks one state per namespace while a
// namespace can be fed from several checkouts of the same repository. Two
// worktrees at different SHAs syncing alternately each see the other's state in
// the manifest, so every sync re-embeds the delta between the two SHAs — the
// same chunks, at four-figure counts per run, repeatedly enough to hit OpenAI
// 429s during benchmark runs. The rows in TurboPuffer must still be rewritten
// on every alternation (row ids are positional, so the namespace holds exactly
// one generation and the last writer wins), but the vectors for text that was
// embedded before can come from here instead of the API.
//
// The key is a hash of the embedded text itself, not the chunk's ContentHash:
// ContentHash omits parts of what is actually embedded (the repo name and the
// rendered header), so it under-identifies the embedder input. Hashing the
// exact input makes the cache correct by construction — same key, same vector —
// and self-invalidating: any chunker or header change alters the text and
// therefore the key, so no version stamp is needed. The model and dimensions
// scope the directory instead, because the same text embeds differently under
// a different model or width.
//
// Entries are one file per vector, sharded git-object style, written via a
// temp file and rename so concurrent index runs can share the cache safely.
type EmbedCache struct {
	root  string
	scope string
	dims  int
	// maxBytes caps the cache across all scopes; Prune evicts by mtime.
	// Zero or negative means unbounded.
	maxBytes int64
}

// NewEmbedCache builds a cache rooted at dir for one model and width. A nil
// return (bad arguments) is a valid cache handle everywhere: every method
// treats a nil receiver as a permanent miss.
func NewEmbedCache(dir, model string, dims int, maxBytes int64) *EmbedCache {
	dir = strings.TrimSpace(dir)
	if dir == "" || dims <= 0 {
		return nil
	}
	scope := unsafeStatePathChars.ReplaceAllString(strings.TrimSpace(model), "-")
	if scope == "" {
		scope = "model"
	}
	return &EmbedCache{
		root:     dir,
		scope:    fmt.Sprintf("%s-%d", scope, dims),
		dims:     dims,
		maxBytes: maxBytes,
	}
}

// OpenDefaultEmbedCache resolves the cache the indexer uses when the caller
// did not inject one: $GX_HOME/index/vector-cache, sized by GX_EMBED_CACHE_MB
// (<=0 lifts the cap), disabled entirely by GX_EMBED_CACHE=0.
func OpenDefaultEmbedCache(cfg Config) *EmbedCache {
	if raw := strings.TrimSpace(os.Getenv("GX_EMBED_CACHE")); raw != "" && !truthy(raw) {
		return nil
	}
	base, err := gxHomeDir()
	if err != nil {
		return nil
	}
	maxBytes := int64(envIntOrDefault("GX_EMBED_CACHE_MB", defaultEmbedCacheMB)) << 20
	return NewEmbedCache(filepath.Join(base, "index", "vector-cache"),
		cfg.OpenAIEmbeddingModel, cfg.EmbeddingDimensions, maxBytes)
}

func (c *EmbedCache) entryPath(key string) string {
	return filepath.Join(c.root, c.scope, key[:2], key[2:])
}

// Get returns the cached vector for text, or nil. A hit refreshes the entry's
// mtime so pruning is least-recently-used rather than least-recently-written.
func (c *EmbedCache) Get(text string) []float32 {
	if c == nil {
		return nil
	}
	path := c.entryPath(hashBytes([]byte(text)))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	if len(data) != c.dims*4 {
		// Truncated by a crash or a concurrent writer caught mid-rename: not
		// recoverable, so remove it rather than miss on it forever.
		_ = os.Remove(path)
		return nil
	}
	vector := make([]float32, c.dims)
	for i := range vector {
		vector[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:]))
	}
	now := time.Now()
	_ = os.Chtimes(path, now, now)
	return vector
}

// Put stores a vector. Failures are deliberately silent: the cache is an
// optimisation, and an index run must never fail because a disk write did.
func (c *EmbedCache) Put(text string, vector []float32) {
	if c == nil || len(vector) != c.dims {
		return
	}
	path := c.entryPath(hashBytes([]byte(text)))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data := make([]byte, len(vector)*4)
	for i, value := range vector {
		binary.LittleEndian.PutUint32(data[i*4:], math.Float32bits(value))
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		_ = os.Remove(tmp.Name())
		return
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		_ = os.Remove(tmp.Name())
	}
}

// Prune evicts least-recently-used entries across every scope until the cache
// fits inside 90% of its cap, so back-to-back runs do not each pay for a walk
// that frees a few kilobytes. It reports how many entries were removed.
func (c *EmbedCache) Prune() int {
	if c == nil || c.maxBytes <= 0 {
		return 0
	}
	type entry struct {
		path  string
		size  int64
		mtime time.Time
	}
	var entries []entry
	var total int64
	_ = filepath.WalkDir(c.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(path, ".tmp") {
			return nil
		}
		info, infoErr := d.Info()
		if infoErr != nil {
			return nil
		}
		entries = append(entries, entry{path: path, size: info.Size(), mtime: info.ModTime()})
		total += info.Size()
		return nil
	})
	if total <= c.maxBytes {
		return 0
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].mtime.Before(entries[j].mtime) })
	target := c.maxBytes / 10 * 9
	removed := 0
	for _, item := range entries {
		if total <= target {
			break
		}
		if os.Remove(item.path) == nil {
			total -= item.size
			removed++
		}
	}
	return removed
}
