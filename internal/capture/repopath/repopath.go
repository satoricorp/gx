// Package repopath turns an absolute path an agent edited into the path a
// commit will contain.
//
// It exists because a repository can have more than one checkout. A linked
// worktree is the same repository — same objects, same commits — reached
// through a different root, and capture runs from whichever checkout ran the
// push. Relativizing a worktree edit against the pushing checkout's root
// produces a path that is syntactically valid and semantically wrong:
//
//	edited   /repo/.claude/worktrees/feature/internal/cli/report.go
//	against  /repo
//	gives    .claude/worktrees/feature/internal/cli/report.go
//	commit   internal/cli/report.go
//
// No hunk ever matches that, so every session is dropped with zero coverage
// and no error anywhere. Worktrees nested inside the repository make it
// silent: filepath.Rel succeeds and returns no "..", so the usual guard
// against escaping the repo passes and the wrong answer looks like a right
// one.
package repopath

import (
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Rel makes path relative to the repository rooted at repoRoot, accepting any
// of that repository's checkouts. A path under a linked worktree comes back
// relative to that worktree, which is the same path the commit contains.
//
// Anything outside every known checkout is returned unchanged, so a session
// that edited a different repository stays unattributable rather than being
// folded into this one.
func Rel(path, repoRoot string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	if path == "" {
		return ""
	}
	if repoRoot == "" {
		return strings.TrimPrefix(path, "./")
	}
	absPath := path
	if abs, err := filepath.Abs(path); err == nil {
		absPath = filepath.ToSlash(abs)
	}
	// A checkout reached through a symlinked path — macOS /tmp and /var are
	// the common ones — is recorded one way in the transcript and reported the
	// other way by git, so both spellings have to be tried.
	candidates := []string{absPath}
	if resolved := resolveExistingPrefix(absPath); resolved != absPath {
		candidates = append(candidates, resolved)
	}
	// Longest root first: a worktree nested inside the main checkout matches
	// both roots, and only the worktree's own root yields the committed path.
	for _, root := range checkoutRoots(repoRoot) {
		for _, candidate := range candidates {
			if rel, ok := relativeTo(candidate, path, root); ok {
				return rel
			}
		}
	}
	return strings.TrimPrefix(path, "./")
}

func relativeTo(absPath, rawPath, root string) (string, bool) {
	if root == "" {
		return "", false
	}
	if rel, err := filepath.Rel(root, absPath); err == nil {
		rel = filepath.ToSlash(rel)
		if rel != ".." && !strings.HasPrefix(rel, "../") {
			return rel, true
		}
	}
	// filepath.Abs resolved against the wrong working directory is not a
	// reason to give up when the raw path already names the root.
	prefix := root + "/"
	if strings.HasPrefix(rawPath, prefix) {
		return strings.TrimPrefix(rawPath, prefix), true
	}
	return "", false
}

var (
	rootsMu    sync.Mutex
	rootsCache = map[string][]string{}
)

// checkoutRoots returns every checkout of repoRoot's repository, longest path
// first. The answer is cached: one capture run asks per edited path, and the
// set of worktrees does not change underneath it.
func checkoutRoots(repoRoot string) []string {
	absRoot := repoRoot
	if abs, err := filepath.Abs(repoRoot); err == nil {
		absRoot = abs
	}
	absRoot = filepath.ToSlash(absRoot)

	rootsMu.Lock()
	defer rootsMu.Unlock()
	if cached, ok := rootsCache[absRoot]; ok {
		return cached
	}
	roots := discoverCheckoutRoots(absRoot)
	rootsCache[absRoot] = roots
	return roots
}

// resolveExistingPrefix resolves symlinks in the deepest part of path that
// still exists and re-appends the rest. EvalSymlinks needs every component to
// exist, and by the time capture runs the edited file often does not: it was
// renamed, deleted, or lived in a worktree that has since been removed. The
// directories above it are what carry the /var -> /private/var indirection,
// and those are usually still there.
func resolveExistingPrefix(path string) string {
	rest := ""
	current := path
	for {
		if resolved, err := filepath.EvalSymlinks(current); err == nil {
			return filepath.ToSlash(filepath.Join(resolved, rest))
		}
		parent := filepath.Dir(current)
		if parent == current {
			return path
		}
		rest = filepath.Join(filepath.Base(current), rest)
		current = parent
	}
}

func addRoot(seen map[string]struct{}, roots *[]string, root string) {
	if root == "" {
		return
	}
	if _, ok := seen[root]; ok {
		return
	}
	seen[root] = struct{}{}
	*roots = append(*roots, root)
}

func discoverCheckoutRoots(absRoot string) []string {
	seen := map[string]struct{}{absRoot: {}}
	roots := []string{absRoot}

	// A repository without git on PATH, or a path that is not a repository at
	// all, still relativizes against the root it was given.
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = absRoot
	out, err := cmd.Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "worktree ") {
				continue
			}
			root := strings.TrimSpace(strings.TrimPrefix(line, "worktree "))
			if root == "" {
				continue
			}
			if abs, err := filepath.Abs(root); err == nil {
				root = abs
			}
			root = filepath.ToSlash(root)
			addRoot(seen, &roots, root)
		}
	}
	// Both spellings of every root, for the same symlink reason as above.
	for _, root := range append([]string(nil), roots...) {
		if resolved, err := filepath.EvalSymlinks(root); err == nil {
			addRoot(seen, &roots, filepath.ToSlash(resolved))
		}
	}
	sort.SliceStable(roots, func(i, j int) bool {
		return len(roots[i]) > len(roots[j])
	})
	return roots
}

// ResetCacheForTest clears the memoized checkout roots.
func ResetCacheForTest() {
	rootsMu.Lock()
	defer rootsMu.Unlock()
	rootsCache = map[string][]string{}
}
