package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func repoLockPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".lgtm", "repo.lock")
}

func withLockFile(lockPath string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return fmt.Errorf("create lock dir: %w", err)
	}
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

func withRepoLock(repoRoot string, fn func() error) error {
	return withLockFile(repoLockPath(repoRoot), fn)
}
