package semantic

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEmbedCacheRoundTrip(t *testing.T) {
	cache := NewEmbedCache(t.TempDir(), "test-model", 4, 0)
	if got := cache.Get("hello"); got != nil {
		t.Fatalf("empty cache returned %v", got)
	}
	want := []float32{1.5, -2, 0, 42}
	cache.Put("hello", want)
	got := cache.Get("hello")
	if len(got) != len(want) {
		t.Fatalf("Get returned %d floats, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Get()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if cache.Get("other") != nil {
		t.Fatal("a different text must miss")
	}

	// A vector of the wrong width must be refused, not stored truncated.
	cache.Put("short", []float32{1, 2, 3})
	if cache.Get("short") != nil {
		t.Fatal("a wrong-width vector was stored")
	}

	// A nil cache is a valid handle that always misses.
	var disabled *EmbedCache
	disabled.Put("hello", want)
	if disabled.Get("hello") != nil {
		t.Fatal("nil cache returned a vector")
	}
	if disabled.Prune() != 0 {
		t.Fatal("nil cache pruned something")
	}
}

func TestEmbedCacheRemovesCorruptEntries(t *testing.T) {
	cache := NewEmbedCache(t.TempDir(), "test-model", 4, 0)
	cache.Put("hello", []float32{1, 2, 3, 4})
	path := cache.entryPath(hashBytes([]byte("hello")))
	if err := os.WriteFile(path, []byte{1, 2, 3}, 0o600); err != nil {
		t.Fatalf("truncate entry: %v", err)
	}
	if got := cache.Get("hello"); got != nil {
		t.Fatalf("corrupt entry returned %v", got)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("corrupt entry was not removed")
	}
}

func TestEmbedCachePruneEvictsLeastRecentlyUsed(t *testing.T) {
	// Ten 16-byte entries against a 100-byte cap: prune must walk down to the
	// 90-byte target, which takes exactly five evictions, oldest first.
	cache := NewEmbedCache(t.TempDir(), "test-model", 4, 100)
	base := time.Now().Add(-time.Hour)
	texts := []string{"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9"}
	for i, text := range texts {
		cache.Put(text, []float32{float32(i), 0, 0, 0})
		stamp := base.Add(time.Duration(i) * time.Second)
		if err := os.Chtimes(cache.entryPath(hashBytes([]byte(text))), stamp, stamp); err != nil {
			t.Fatalf("stamp entry: %v", err)
		}
	}
	if removed := cache.Prune(); removed != 5 {
		t.Fatalf("pruned %d entries, want 5", removed)
	}
	for _, text := range texts[:5] {
		if cache.Get(text) != nil {
			t.Fatalf("oldest entry %q survived the prune", text)
		}
	}
	for _, text := range texts[5:] {
		if cache.Get(text) == nil {
			t.Fatalf("recent entry %q was evicted", text)
		}
	}
	// Under the cap nothing further is evicted.
	if removed := cache.Prune(); removed != 0 {
		t.Fatalf("second prune removed %d entries from an under-cap cache", removed)
	}
}

func TestOpenDefaultEmbedCacheHonoursEnvironment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GX_EMBED_CACHE", "")
	t.Setenv("GX_EMBED_CACHE_MB", "")
	cfg := Config{OpenAIEmbeddingModel: "text-embedding-3-small", EmbeddingDimensions: 1536}

	cache := OpenDefaultEmbedCache(cfg)
	if cache == nil {
		t.Fatal("default cache is disabled with no environment override")
	}
	if want := filepath.Join(home, "index", "vector-cache"); cache.root != want {
		t.Fatalf("cache root = %q, want %q", cache.root, want)
	}
	if cache.maxBytes != defaultEmbedCacheMB<<20 {
		t.Fatalf("cache cap = %d, want %d", cache.maxBytes, int64(defaultEmbedCacheMB)<<20)
	}

	t.Setenv("GX_EMBED_CACHE", "0")
	if OpenDefaultEmbedCache(cfg) != nil {
		t.Fatal("GX_EMBED_CACHE=0 did not disable the cache")
	}
}
