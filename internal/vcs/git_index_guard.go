package vcs

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
)

// preservingGitIndex snapshots the colocated Git index around fn and restores
// it byte-for-byte if fn changed it without moving git HEAD. jj rewrites the
// index whenever it snapshots the working copy or imports git refs, converting
// the user's staged entries into intent-to-add placeholders; read-only gx
// operations must not surface that side effect.
func (s *Service) preservingGitIndex(ctx context.Context, repoRoot string, fn func() error) error {
	indexPath := s.gitIndexPath(ctx, repoRoot)
	if indexPath == "" {
		return fn()
	}
	before, readErr := os.ReadFile(indexPath)
	if readErr != nil {
		return fn()
	}
	headBefore := s.gitHeadState(ctx, repoRoot)

	fnErr := fn()

	if headBefore != "" && headBefore == s.gitHeadState(ctx, repoRoot) {
		after, err := os.ReadFile(indexPath)
		if err == nil && !bytes.Equal(before, after) {
			_ = os.WriteFile(indexPath, before, 0o644)
		}
	}
	return fnErr
}

func (s *Service) gitIndexPath(ctx context.Context, repoRoot string) string {
	path, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--git-path", "index")
	if err != nil || strings.TrimSpace(path) == "" {
		return ""
	}
	path = strings.TrimSpace(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, path)
	}
	return path
}

// gitHeadState identifies the checkout: the symbolic ref (or detached marker)
// plus the commit HEAD resolves to. Index restoration is only safe while both
// are unchanged.
func (s *Service) gitHeadState(ctx context.Context, repoRoot string) string {
	commit, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	ref, err := s.runTrimmed(ctx, repoRoot, "git", "symbolic-ref", "-q", "HEAD")
	if err != nil {
		ref = "(detached)"
	}
	return strings.TrimSpace(ref) + "@" + strings.TrimSpace(commit)
}
