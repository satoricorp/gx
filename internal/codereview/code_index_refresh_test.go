package codereview

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
	"github.com/satoricorp/gx/internal/semantic"
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

// TestRefreshCodeIndexSkipsWithoutAGXHome pins that review does not index on a
// machine that has never run `gx init` — the CI shape.
//
// Keeping a repository indexed is GX Cloud's job, done on merge from the GitHub
// App, so it covers every user rather than only those whose CI happens to carry
// credentials. Review on such a machine reads that index and does not try to
// maintain one: it must neither create a GX home for the manifest nor index
// without one, because a run with no manifest re-embeds everything and can
// never delete rows for files that are gone.
func TestRefreshCodeIndexSkipsWithoutAGXHome(t *testing.T) {
	backend := gxtest.NewIndexBackend(t)
	backend.Use(t)
	gxHome := filepath.Join(t.TempDir(), "absent")
	t.Setenv("GX_HOME", gxHome)

	log := &EvidenceLog{}
	refreshCodeIndex(context.Background(), indexableRepo(t), log)

	if backend.Upserted() {
		t.Fatalf("review indexed on a machine with no GX home; that is GX Cloud's job now.\nrequests: %v", backend.Requests())
	}
	if _, err := os.Stat(gxHome); !os.IsNotExist(err) {
		t.Fatalf("refresh created %s (stat error = %v), want it untouched", gxHome, err)
	}
	statuses := log.Statuses()
	if len(statuses) != 1 || statuses[0].State != EvidenceUnavailable {
		t.Fatalf("Statuses() = %#v, want one status explaining why the index may be stale", statuses)
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

// TestMissingIndexWarningsCarryTheirRemedy pins the property that makes these
// warnings worth printing: a user told their review ran on partial evidence
// must also be told how to fix it, and told once.
//
// Every route goes to the website. `gx index` is a hidden maintenance command
// that fills one developer's namespace from one developer's checkout; the
// index a review actually reads is the one GX Cloud maintains from the GitHub
// App, and that is connected on the site.
func TestMissingIndexWarningsCarryTheirRemedy(t *testing.T) {
	// The shape a repository with a git remote and no index anywhere produces:
	// three candidate namespaces, all absent.
	statuses := []EvidenceStatus{
		{
			Source: codeIndexEvidenceSource, Namespace: "gx-001-acme-api-v2", State: EvidenceMissing,
			Detail: codeIndexMissingDetail(codeIndexTarget{Origin: semantic.NamespaceOriginPrimary}),
			Remedy: codeIndexMissingRemedy(codeIndexTarget{Origin: semantic.NamespaceOriginPrimary}),
		},
		{
			Source: codeIndexEvidenceSource, Namespace: "gx-local-api-abc-v2", State: EvidenceMissing,
			Detail: codeIndexMissingDetail(codeIndexTarget{Origin: semantic.NamespaceOriginPreRemote}),
			Remedy: codeIndexMissingRemedy(codeIndexTarget{Origin: semantic.NamespaceOriginPreRemote}),
		},
		{
			Source: codeIndexEvidenceSource, Namespace: "repo-acme-api", State: EvidenceMissing,
			Detail: codeIndexMissingDetail(codeIndexTarget{Origin: semantic.NamespaceOriginConsole}),
			Remedy: codeIndexMissingRemedy(codeIndexTarget{Origin: semantic.NamespaceOriginConsole}),
		},
	}

	warnings := EvidenceWarnings(statuses)
	if len(warnings) != 1 {
		t.Fatalf("EvidenceWarnings() = %#v, want one warning for the one source", warnings)
	}
	warning := warnings[0]
	if !strings.Contains(warning, "https://gx.run/repositories") {
		t.Fatalf("warning = %q, want it to say where to get the repository indexed", warning)
	}
	if strings.Contains(warning, "gx index") {
		t.Fatalf("warning = %q, want the website rather than the hidden `gx index` command", warning)
	}
	if count := strings.Count(warning, "https://gx.run/repositories"); count != 1 {
		t.Fatalf("warning states the remedy %d times, want once:\n%s", count, warning)
	}
}

// TestEveryMissingCodeIndexOriginHasARemedy guards the switch against a new
// namespace origin being added without one, which is how a warning ends up
// telling someone their review is degraded and nothing else.
func TestEveryMissingCodeIndexOriginHasARemedy(t *testing.T) {
	for _, origin := range []string{
		semantic.NamespaceOriginPrimary,
		semantic.NamespaceOriginConsole,
		"some-future-origin",
	} {
		target := codeIndexTarget{Origin: origin}
		if remedy := codeIndexMissingRemedy(target); !strings.Contains(remedy, "https://gx.run/repositories") {
			t.Errorf("origin %q remedy = %q, want it to point at the console", origin, remedy)
		}
		if detail := codeIndexMissingDetail(target); strings.TrimSpace(detail) == "" {
			t.Errorf("origin %q has no detail", origin)
		}
	}
}
