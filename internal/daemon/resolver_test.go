package daemon

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestAmbientResolverReusesExistingSession(t *testing.T) {
	repoRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(repoRoot, ".git"), 0o755); err != nil {
		t.Fatalf("create .git: %v", err)
	}
	t.Setenv("GX_HOME", t.TempDir())

	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	defer store.Close()

	resolver := NewAmbientSessionResolver(store)
	resolver.resolveProcess = func(*http.Request) ClientProcess {
		return ClientProcess{
			PID:     1234,
			Name:    "codex",
			Command: "codex",
			Cwd:     repoRoot,
		}
	}
	req := &http.Request{RemoteAddr: "127.0.0.1:5000"}

	first, err := resolver.ResolveRequestSession(req)
	if err != nil {
		t.Fatalf("first ResolveRequestSession() error = %v", err)
	}
	if !first.Capture || first.ID == "" {
		t.Fatalf("first ResolveRequestSession() = %#v, want captured session", first)
	}

	second, err := resolver.ResolveRequestSession(req)
	if err != nil {
		t.Fatalf("second ResolveRequestSession() error = %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("second session ID = %q, want existing %q", second.ID, first.ID)
	}
}
