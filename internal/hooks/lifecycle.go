package hooks

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/vcs"
)

// PrepareCommitMsgOptions configures the prepare-commit-msg hook handler.
type PrepareCommitMsgOptions struct {
	RepoRoot    string
	MessagePath string
}

// PrepareCommitMsg preserves or adds exactly one GX revision trailer.
func PrepareCommitMsg(opts PrepareCommitMsgOptions) error {
	repoRoot := strings.TrimSpace(opts.RepoRoot)
	messagePath := strings.TrimSpace(opts.MessagePath)
	if repoRoot == "" || messagePath == "" {
		return fmt.Errorf("repo root and message path required")
	}
	if !EnabledForRepo(context.Background(), repoRoot) {
		return nil
	}
	data, err := os.ReadFile(messagePath)
	if err != nil {
		return fmt.Errorf("read commit message: %w", err)
	}
	updated, err := vcs.PrepareCommitMessageHook(string(data))
	if err != nil {
		return err
	}
	if updated == string(data) {
		return nil
	}
	return os.WriteFile(messagePath, []byte(updated), 0o644)
}

// PostCommitOptions configures the post-commit hook handler.
type PostCommitOptions struct {
	RepoRoot string
}

// RunPostCommit records the new HEAD commit in GX metadata.
func RunPostCommit(ctx context.Context, opts PostCommitOptions) error {
	repoRoot := strings.TrimSpace(opts.RepoRoot)
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repoRoot = cwd
	}
	if !EnabledForRepo(ctx, repoRoot) {
		return nil
	}
	return vcs.NewService().RunPostCommitHook(ctx, repoRoot)
}

// PostRewriteOptions configures the post-rewrite hook handler.
type PostRewriteOptions struct {
	RepoRoot string
	Input    io.Reader
}

// RunPostRewrite updates stored commit OIDs after amend/rebase.
func RunPostRewrite(ctx context.Context, opts PostRewriteOptions) error {
	repoRoot := strings.TrimSpace(opts.RepoRoot)
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		repoRoot = cwd
	}
	input := opts.Input
	if input == nil {
		input = os.Stdin
	}
	var builder strings.Builder
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	// Checked after draining stdin so git never sees a short write.
	if !EnabledForRepo(ctx, repoRoot) {
		return nil
	}
	return vcs.NewService().RunPostRewriteHook(ctx, repoRoot, vcs.ParsePostRewriteMappings(builder.String()))
}

// ResolveHooksDir returns the absolute git hooks directory for repoRoot.
func ResolveHooksDir(ctx context.Context, repoRoot string) (string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return "", fmt.Errorf("repo root required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", "--path-format=absolute", "--git-path", "hooks")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --git-path hooks: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return filepath.Clean(strings.TrimSpace(string(out))), nil
}
