package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func repoLockPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".gx", "repo.lock")
}

func withRepoLock(repoRoot string, fn func() error) error {
	if err := os.MkdirAll(filepath.Join(repoRoot, ".gx"), 0o755); err != nil {
		return fmt.Errorf("create gx dir: %w", err)
	}
	lockPath := repoLockPath(repoRoot)
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("open repo lock: %w", err)
	}
	defer file.Close()

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("acquire repo lock: %w", err)
	}
	defer func() { _ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN) }()

	return fn()
}

func withGitRetry(ctx context.Context, label string, fn func() error) error {
	const maxAttempts = 5
	backoff := 80 * time.Millisecond
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			if !isGitIndexLock(err) || attempt == maxAttempts {
				break
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			continue
		}
		return nil
	}
	if lastErr == nil {
		return nil
	}
	if isGitIndexLock(lastErr) {
		return fmt.Errorf("%s failed after %d retries: %w", label, maxAttempts, lastErr)
	}
	return lastErr
}

func isGitIndexLock(err error) bool {
	if err == nil {
		return false
	}
	value := strings.ToLower(err.Error())
	return strings.Contains(value, "index.lock") || strings.Contains(value, "unable to create") && strings.Contains(value, "index")
}

func (s *Service) runGitLocked(ctx context.Context, repoRoot, label string, fn func() error) error {
	return withGitRetry(ctx, label, fn)
}
