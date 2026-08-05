package hooks

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

// Discovery deliberately casts a wide net: any transcript that mentions this
// repository's path is a candidate, which is how sessions belonging to other
// repositories end up staged alongside this one's. Content matching used to be
// the only thing throwing them back, and matching is about to stop being a
// prerequisite for sharing.
//
// A session bound to a different repository must not become shareable here.
func TestSessionsForThisRepoWithholdsForeignSessions(t *testing.T) {
	repo := initGateRepo(t, "git@github.com:satoricorp/gx.git")

	stager := &fakeOriginStager{origins: map[string]string{
		"own":       "github.com/satoricorp/gx",
		"other":     "github.com/satoricorp/console",
		"personal":  "github.com/joe/side-project",
		"unbound":   "",
		"sameOther": "github.com/satoricorp/console",
	}}

	ids := []string{"own", "other", "personal", "unbound", "sameOther"}
	kept, withheld := sessionsForThisRepo(context.Background(), stager, repo, ids)

	if withheld != 3 {
		t.Fatalf("withheld = %d, want 3 (console twice, side-project once)", withheld)
	}
	got := map[string]bool{}
	for _, id := range kept {
		got[id] = true
	}
	if !got["own"] {
		t.Error("a session bound to this repository was withheld")
	}
	// An unbound session keeps the previous behavior: this change may only
	// ever share less, never more, and it must not strand sessions whose
	// directory could not be resolved.
	if !got["unbound"] {
		t.Error("an unbound session was withheld; the gate is not purely subtractive")
	}
	for _, id := range []string{"other", "personal", "sameOther"} {
		if got[id] {
			t.Errorf("session %q belongs to another repository but was shareable", id)
		}
	}
}

// With no identity of its own, this repository has nothing to compare against,
// so it withholds nothing rather than withholding everything.
func TestSessionsForThisRepoWithholdsNothingWithoutAnOrigin(t *testing.T) {
	repo := initGateRepo(t, "")
	stager := &fakeOriginStager{origins: map[string]string{"a": "github.com/satoricorp/console"}}

	kept, withheld := sessionsForThisRepo(context.Background(), stager, repo, []string{"a"})
	if withheld != 0 || len(kept) != 1 {
		t.Fatalf("kept=%v withheld=%d, want everything kept", kept, withheld)
	}
}

// A gate that cannot read its own data narrows nothing rather than dropping
// work on the floor.
func TestSessionsForThisRepoKeepsEverythingWhenTheLookupFails(t *testing.T) {
	repo := initGateRepo(t, "git@github.com:satoricorp/gx.git")
	stager := &fakeOriginStager{err: os.ErrClosed}

	kept, withheld := sessionsForThisRepo(context.Background(), stager, repo, []string{"a", "b"})
	if withheld != 0 || len(kept) != 2 {
		t.Fatalf("kept=%v withheld=%d, want everything kept on lookup failure", kept, withheld)
	}
}

type fakeOriginStager struct {
	storage.CaptureStager
	origins map[string]string
	err     error
}

func (f *fakeOriginStager) SessionOriginsByID(_ context.Context, ids []string) (map[string]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := map[string]string{}
	for _, id := range ids {
		if origin, ok := f.origins[id]; ok {
			out[id] = origin
		}
	}
	return out, nil
}

func initGateRepo(t *testing.T, origin string) string {
	t.Helper()
	dir := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}
	runGateGit(t, dir, "init", "-q", "-b", "main")
	runGateGit(t, dir, "config", "user.email", "test@example.com")
	runGateGit(t, dir, "config", "user.name", "Test")
	if origin != "" {
		runGateGit(t, dir, "remote", "add", "origin", origin)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGateGit(t, dir, "add", "README.md")
	runGateGit(t, dir, "commit", "-q", "-m", "initial")
	return dir
}

func runGateGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
