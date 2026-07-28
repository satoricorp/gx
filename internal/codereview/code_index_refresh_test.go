package codereview

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestRefreshCodeIndexSkipsWithoutCredentials pins that the refresh seam is
// inert when there is nothing to refresh into. It must not fabricate evidence
// and must not attempt a network call.
func TestRefreshCodeIndexSkipsWithoutCredentials(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("TURBOPUFFER_API_KEY", "")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), t.TempDir(), log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want no evidence recorded when unconfigured", statuses)
	}
}

// TestRefreshCodeIndexHonoursKillSwitch pins the escape hatch: a review must be
// able to run without touching the index at all.
func TestRefreshCodeIndexHonoursKillSwitch(t *testing.T) {
	t.Setenv("GX_REVIEW_INDEX_REFRESH", "0")
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TURBOPUFFER_API_KEY", "test-key")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), t.TempDir(), log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want the refresh skipped entirely", statuses)
	}
}

// TestRefreshCodeIndexWithoutAGXHomeStillIndexes pins the CI case: a machine
// with no GX home must still get a fresh index.
//
// CI is where the code is newest and the index most likely to be stale, so it
// is the last place a refresh should be turned off. The manifest is only a
// cache, so a machine that keeps nothing between runs gets a disposable one
// rather than no refresh.
func TestRefreshCodeIndexWithoutAGXHomeStillIndexes(t *testing.T) {
	backend := gxtest.NewIndexBackend(t)
	backend.Use(t)
	gxHome := filepath.Join(t.TempDir(), "absent")
	t.Setenv("GX_HOME", gxHome)

	refreshCodeIndex(context.Background(), indexableRepo(t), &EvidenceLog{})

	if !backend.Upserted() {
		t.Fatalf("no rows upserted; the refresh was skipped on a machine with no GX home.\nrequests: %v", backend.Requests())
	}
	if _, err := os.Stat(gxHome); !os.IsNotExist(err) {
		t.Fatalf("refresh created %s (stat error = %v); the manifest must go somewhere disposable", gxHome, err)
	}
}

// TestRefreshCodeIndexWritesItsManifestIntoTheGXHome is the other half: when
// there is a GX home the manifest belongs in it, so the next review on this
// machine is incremental rather than re-embedding the whole checkout.
func TestRefreshCodeIndexWritesItsManifestIntoTheGXHome(t *testing.T) {
	backend := gxtest.NewIndexBackend(t)
	backend.Use(t)
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	refreshCodeIndex(context.Background(), indexableRepo(t), &EvidenceLog{})

	manifests, err := filepath.Glob(filepath.Join(gxHome, "index", "*.json"))
	if err != nil {
		t.Fatalf("Glob() error = %v", err)
	}
	if len(manifests) == 0 {
		t.Fatalf("no index manifest under %s, so the next review would re-embed everything", gxHome)
	}
}

// indexableRepo is a directory with enough source in it to produce chunks, so
// a refresh over it reaches the upsert rather than finding nothing to do.
func indexableRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", `package app

// Run starts the application and returns the first error it meets.
func Run(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}
`)
	return root
}

// TestRefreshCodeIndexSkipsEmptyRepoRoot guards the seam against being called
// before the repository root is known.
func TestRefreshCodeIndexSkipsEmptyRepoRoot(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("TURBOPUFFER_API_KEY", "test-key")

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), "  ", log)

	if statuses := log.Statuses(); len(statuses) != 0 {
		t.Fatalf("Statuses() = %#v, want no refresh without a repository root", statuses)
	}
}
