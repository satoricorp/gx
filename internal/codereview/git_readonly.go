package codereview

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Review inspects a checkout it does not own, so it must not disturb one.
// Plain `git status` and `git diff` do: refreshing the stat cache rewrites
// .git/index and takes .git/index.lock, which can collide with a concurrent
// git in CI and leaves a modified index behind in a command that promises to
// change nothing.
//
// `git --no-optional-locks` only covers status; diff refreshes the index
// regardless. Pointing GIT_INDEX_FILE at a copy covers every git call
// uniformly: a refresh lands on the copy, and the real index is never opened
// for writing. The copy lives in a read-only directory so git cannot write it
// either, which is what makes one copy safe to share across the concurrent
// git calls a review fans out.
type scratchIndexKeyType struct{}

var scratchIndexKey scratchIndexKeyType

// withScratchGitIndex copies repoRoot's index somewhere git may read but not
// write, and returns a context that points review's git calls at the copy. It
// degrades to a no-op: if the index cannot be found or copied, review still
// runs, just without the protection.
func withScratchGitIndex(ctx context.Context, repoRoot string) (context.Context, func()) {
	noop := func() {}
	indexPath := gitIndexPath(ctx, repoRoot)
	if indexPath == "" {
		return ctx, noop
	}
	data, err := os.ReadFile(indexPath)
	if err != nil {
		// No index at all (a repo with no commits yet) is not a failure:
		// there is nothing for git to rewrite, so there is nothing to guard.
		return ctx, noop
	}
	dir, err := os.MkdirTemp("", "gx-review-index-")
	if err != nil {
		return ctx, noop
	}
	cleanup := func() {
		// The directory was made read-only to keep git out; take write back
		// so the cleanup can remove it.
		_ = os.Chmod(dir, 0o700)
		_ = os.RemoveAll(dir)
	}
	scratch := filepath.Join(dir, "index")
	if err := os.WriteFile(scratch, data, 0o600); err != nil {
		cleanup()
		return ctx, noop
	}
	// The copy must carry the original index's modification time.
	//
	// Git decides which entries are worth re-hashing by comparing each entry's
	// cached mtime against the mtime of the index file itself: an entry at or
	// after it is "racily clean" and gets re-read from disk, because a file
	// written in the same instant as the index cannot be trusted on stat alone.
	// A fresh copy has a newer mtime than the index it came from, which moves
	// that boundary and makes git skip exactly those re-reads — so a file
	// modified moments before the review runs is reported as unchanged, and the
	// review sees nothing to review. Preserving the timestamp keeps the copy
	// indistinguishable from the original for that decision.
	if info, err := os.Stat(indexPath); err == nil {
		_ = os.Chtimes(scratch, info.ModTime(), info.ModTime())
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		cleanup()
		return ctx, noop
	}
	return context.WithValue(ctx, scratchIndexKey, scratch), cleanup
}

func scratchGitIndex(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	path, _ := ctx.Value(scratchIndexKey).(string)
	return path
}

// gitCommand builds a git invocation for review's read-only inspection of a
// repo. Every git call in this package goes through it so none of them can
// write the repo's real index.
func gitCommand(ctx context.Context, repoRoot string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoRoot
	if scratch := scratchGitIndex(ctx); scratch != "" {
		cmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+scratch)
	}
	return cmd
}

// gitIndexPath asks git where the index lives, which is the only answer that
// holds for worktrees and an explicit GIT_DIR. This one runs unredirected by
// design: it is what discovers the path to copy, and rev-parse does not write
// the index.
func gitIndexPath(ctx context.Context, repoRoot string) string {
	if strings.TrimSpace(repoRoot) == "" {
		return ""
	}
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-path", "index")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return ""
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, path)
	}
	return path
}
