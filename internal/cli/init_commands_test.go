package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallSlashCommandsWritesForDetectedAgents(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, dir := range []string{".claude", ".codex"} {
		if err := os.MkdirAll(filepath.Join(home, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Claude Code", "Codex"}
	if strings.Join(installed, ",") != strings.Join(want, ",") {
		t.Fatalf("installed = %v, want %v", installed, want)
	}

	claude, err := os.ReadFile(filepath.Join(home, ".claude", "commands", slashCommandName+".md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"description:", "$ARGUMENTS", "gx_review", slashCommandMarker} {
		if !strings.Contains(string(claude), fragment) {
			t.Fatalf("claude command missing %q:\n%s", fragment, claude)
		}
	}

	codex, err := os.ReadFile(filepath.Join(home, ".codex", "prompts", slashCommandName+".md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(codex), "---") {
		t.Fatalf("codex prompt should not have frontmatter:\n%s", codex)
	}
	if !strings.Contains(string(codex), "$ARGUMENTS") {
		t.Fatalf("codex prompt missing $ARGUMENTS:\n%s", codex)
	}

	if _, statErr := os.Stat(filepath.Join(home, ".cursor")); !os.IsNotExist(statErr) {
		t.Fatal("expected no .cursor directory to be created")
	}
}

func TestInstallSlashCommandsCursorBodyHasNoPlaceholder(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".cursor"), 0o755); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Cursor" {
		t.Fatalf("installed = %v, want [Cursor]", installed)
	}
	cursor, err := os.ReadFile(filepath.Join(home, ".cursor", "commands", slashCommandName+".md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(cursor), "$ARGUMENTS") {
		t.Fatalf("cursor command must not use $ARGUMENTS:\n%s", cursor)
	}
	if !strings.Contains(string(cursor), "gx_review") {
		t.Fatalf("cursor command missing gx_review:\n%s", cursor)
	}
}

func TestInstallSlashCommandsPreservesUserOwnedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "my own review command\n"
	path := filepath.Join(dir, slashCommandName+".md")
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if len(installed) != 0 {
		t.Fatalf("installed = %v, want none (file is user-owned)", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Fatalf("user-owned file was modified:\n%s", data)
	}
}

func TestInstallSlashCommandsRefreshesManagedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := slashCommandMarker + "\n\nold template\n"
	path := filepath.Join(dir, slashCommandName+".md")
	if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "old template") {
		t.Fatalf("managed file was not refreshed:\n%s", data)
	}
	if !strings.Contains(string(data), "gx_review") {
		t.Fatalf("refreshed file missing template body:\n%s", data)
	}
}

func TestInstallSlashCommandsRemovesRenamedPredecessor(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// A file this tool installed under the old name, carrying our marker.
	legacy := filepath.Join(dir, "better-review.md")
	if err := os.WriteFile(legacy, []byte(slashCommandMarker+"\n\nold command\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := installSlashCommands(); err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(legacy); !os.IsNotExist(statErr) {
		t.Error("renamed predecessor should be removed, leaving one command not two")
	}
	if _, err := os.Stat(filepath.Join(dir, slashCommandName+".md")); err != nil {
		t.Errorf("current command missing: %v", err)
	}
}

func TestInstallSlashCommandsKeepsUserOwnedPredecessor(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Same old name, but the user removed the marker: it is their file now.
	legacy := filepath.Join(dir, "better-review.md")
	custom := "my own review command\n"
	if err := os.WriteFile(legacy, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := installSlashCommands(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(legacy)
	if err != nil {
		t.Fatalf("user-owned predecessor was deleted: %v", err)
	}
	if string(data) != custom {
		t.Errorf("user-owned predecessor was modified:\n%s", data)
	}
}
