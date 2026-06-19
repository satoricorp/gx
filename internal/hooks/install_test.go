package hooks_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/hooks"
)

func TestInstallPrePushHook(t *testing.T) {
	repo := t.TempDir()
	gitDir := filepath.Join(repo, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := hooks.Install(hooks.InstallOptions{
		RepoRoot: repo,
		GXPath:   "/usr/local/bin/gx",
	}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(gitDir, "pre-push"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !hooks.IsInstalled(repo) {
		t.Fatal("expected hook installed")
	}
	if content == "" {
		t.Fatal("empty hook script")
	}
}

func TestRefRangeForPush(t *testing.T) {
	zero := "0000000000000000000000000000000000000000"
	if got := hooks.RefRangeForPush("abc", zero); got != "abc" {
		t.Fatalf("new branch range = %q", got)
	}
	if got := hooks.RefRangeForPush("def", "abc"); got != "abc..def" {
		t.Fatalf("update range = %q", got)
	}
}

func TestParsePrePushLine(t *testing.T) {
	localRef, localSHA, remoteRef, remoteSHA, err := hooks.ParsePrePushLine(
		"refs/heads/main abc refs/heads/main def",
	)
	if err != nil {
		t.Fatal(err)
	}
	if localRef != "refs/heads/main" || localSHA != "abc" || remoteRef != "refs/heads/main" || remoteSHA != "def" {
		t.Fatalf("unexpected parse: %q %q %q %q", localRef, localSHA, remoteRef, remoteSHA)
	}
}

func TestHuskySnippet(t *testing.T) {
	snippet := hooks.HuskyPrePushSnippet("/bin/gx")
	if snippet == "" || snippet[:6] != "remote" {
		t.Fatalf("unexpected snippet: %q", snippet)
	}
}
