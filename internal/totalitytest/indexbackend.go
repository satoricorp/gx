package totalitytest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// IndexBackend is a stand-in for the embeddings API and TurboPuffer that
// succeeds at everything a repository index asks of it, and records what it was
// asked.
//
// Tests about what indexing does to the machine need the indexing to actually
// complete: the manifest is written only after the upserts succeed, so a test
// pointed at a dead endpoint passes because the network failed rather than
// because the code behaved. That is how the review read-only assertion came to
// be satisfied by accident for as long as it was.
type IndexBackend struct {
	// URL is the base URL to hand to both TOTALITY_OPENAI_BASE_URL and
	// TOTALITY_TPUF_BASE_URL.
	URL string

	mu       sync.Mutex
	requests []string
}

// NewIndexBackend starts a backend and stops it when the test ends.
func NewIndexBackend(t *testing.T) *IndexBackend {
	t.Helper()
	backend := &IndexBackend{}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/embeddings", func(w http.ResponseWriter, r *http.Request) {
		backend.record(r)
		var req struct {
			Input      []string `json:"input"`
			Dimensions int      `json:"dimensions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Dimensions <= 0 {
			req.Dimensions = 1536
		}
		type item struct {
			Index     int       `json:"index"`
			Embedding []float64 `json:"embedding"`
		}
		var resp struct {
			Data []item `json:"data"`
		}
		for i := range req.Input {
			resp.Data = append(resp.Data, item{Index: i, Embedding: make([]float64, req.Dimensions)})
		}
		writeJSON(t, w, http.StatusOK, resp)
	})
	// A namespace that does not exist yet is the state a first index writes
	// into, and the state a throwaway fixture repository is always in.
	mux.HandleFunc("/v2/namespaces/{namespace}/metadata", func(w http.ResponseWriter, r *http.Request) {
		backend.record(r)
		http.Error(w, `{"error":"namespace not found"}`, http.StatusNotFound)
	})
	mux.HandleFunc("/v2/namespaces/{namespace}", func(w http.ResponseWriter, r *http.Request) {
		backend.record(r)
		writeJSON(t, w, http.StatusOK, map[string]string{"status": "OK"})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend.record(r)
		writeJSON(t, w, http.StatusOK, map[string]string{"status": "OK"})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	backend.URL = server.URL
	return backend
}

// Use points the index configuration at this backend for the duration of the
// test, credentials included, so the refresh runs instead of short-circuiting
// on a missing key.
func (b *IndexBackend) Use(t *testing.T) {
	t.Helper()
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TURBOPUFFER_API_KEY", "test-key")
	t.Setenv("TOTALITY_OPENAI_BASE_URL", b.URL)
	t.Setenv("TOTALITY_TPUF_BASE_URL", b.URL)
}

func (b *IndexBackend) record(r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.requests = append(b.requests, r.Method+" "+r.URL.Path)
}

// Requests are the calls this backend served, in order.
func (b *IndexBackend) Requests() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]string(nil), b.requests...)
}

// Upserted reports whether rows were written to a namespace, which is the
// signal that an index refresh ran rather than being skipped.
func (b *IndexBackend) Upserted() bool {
	for _, request := range b.Requests() {
		if strings.HasPrefix(request, "POST /v2/namespaces/") && !strings.HasSuffix(request, "/metadata") {
			return true
		}
	}
	return false
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, payload any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Errorf("encode fake index backend response: %v", err)
	}
}
