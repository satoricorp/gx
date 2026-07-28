package codereview

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// The index copy review reads through must carry the original index's
// modification time.
//
// Git decides which entries are worth re-reading from disk by comparing each
// entry's cached mtime against the mtime of the index file itself: an entry at
// or after it is "racily clean", because a file written in the same instant as
// the index cannot be trusted on stat data alone. A fresh copy has a newer
// mtime than the index it came from, which moves that boundary past exactly
// those entries. Git then skips the re-read, reports the file unchanged, and
// the review finds nothing to review — exiting 0 with empty output, so nothing
// downstream can tell a clean tree from a missed one.
//
// This surfaced as a test that failed about 2% of the time under load (2 of 90
// runs), and not at all unloaded (0 of 150) — load is what bunches a write and
// a commit into the same timestamp window. With the timestamp preserved the
// same load produced 0 failures in 240 runs.
//
// The assertion is on the timestamp rather than on a reproduced miss on
// purpose: git also compares ctime, which cannot be set from a test, so the
// full condition cannot be staged deterministically. The timestamp is the
// property the fix establishes and the one a future change would break.
func TestScratchGitIndexPreservesTheIndexTimestamp(t *testing.T) {
	root := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}
	runGitHere(t, root, "init", "-q", "-b", "main")
	runGitHere(t, root, "config", "user.email", "test@example.com")
	runGitHere(t, root, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitHere(t, root, "add", "README.md")
	runGitHere(t, root, "commit", "-q", "-m", "initial")

	realIndex, err := os.Stat(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, release := withScratchGitIndex(context.Background(), root)
	defer release()

	scratch := scratchGitIndex(ctx)
	if scratch == "" {
		t.Fatal("no scratch index was created")
	}
	copied, err := os.Stat(scratch)
	if err != nil {
		t.Fatal(err)
	}
	if !copied.ModTime().Equal(realIndex.ModTime()) {
		t.Fatalf("scratch index mtime = %v, want the real index's %v (difference: %v)",
			copied.ModTime(), realIndex.ModTime(), copied.ModTime().Sub(realIndex.ModTime()))
	}
}

// And the copy still has to be a faithful copy, or preserving its timestamp
// would just make a wrong answer look trustworthy.
func TestScratchGitIndexCopiesTheIndexContents(t *testing.T) {
	root := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}
	runGitHere(t, root, "init", "-q", "-b", "main")
	runGitHere(t, root, "config", "user.email", "test@example.com")
	runGitHere(t, root, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitHere(t, root, "add", "README.md")
	runGitHere(t, root, "commit", "-q", "-m", "initial")

	want, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, release := withScratchGitIndex(context.Background(), root)
	defer release()
	got, err := os.ReadFile(scratchGitIndex(ctx))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("scratch index differs from the real index (%d vs %d bytes)", len(got), len(want))
	}
}

func runGitHere(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
