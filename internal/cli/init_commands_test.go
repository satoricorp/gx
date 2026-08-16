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

	claude, err := os.ReadFile(filepath.Join(home, ".claude", "commands", "enhance.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"description:", "$ARGUMENTS", "gx_enhance", "slash-enhance", slashCommandMarker} {
		if !strings.Contains(string(claude), fragment) {
			t.Fatalf("claude command missing %q:\n%s", fragment, claude)
		}
	}

	claudeReview, err := os.ReadFile(filepath.Join(home, ".claude", "commands", "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"description:", "argument-hint:", "$ARGUMENTS", "gx_review", "slash-review", slashCommandMarker} {
		if !strings.Contains(string(claudeReview), fragment) {
			t.Fatalf("claude review command missing %q:\n%s", fragment, claudeReview)
		}
	}

	codex, err := os.ReadFile(filepath.Join(home, ".codex", "prompts", "enhance.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(codex), "---") {
		t.Fatalf("codex prompt should not have frontmatter:\n%s", codex)
	}
	if !strings.Contains(string(codex), "$ARGUMENTS") {
		t.Fatalf("codex prompt missing $ARGUMENTS:\n%s", codex)
	}

	codexReview, err := os.ReadFile(filepath.Join(home, ".codex", "prompts", "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(codexReview), "---") {
		t.Fatalf("codex review prompt should not have frontmatter:\n%s", codexReview)
	}
	for _, fragment := range []string{"$ARGUMENTS", "gx_review"} {
		if !strings.Contains(string(codexReview), fragment) {
			t.Fatalf("codex review prompt missing %q:\n%s", fragment, codexReview)
		}
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
	for _, name := range []string{"enhance.md", "review.md"} {
		cursor, err := os.ReadFile(filepath.Join(home, ".cursor", "commands", name))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(cursor), "$ARGUMENTS") {
			t.Fatalf("cursor %s must not use $ARGUMENTS:\n%s", name, cursor)
		}
	}
	cursorReview, err := os.ReadFile(filepath.Join(home, ".cursor", "commands", "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cursorReview), "gx_review") {
		t.Fatalf("cursor review command missing gx_review:\n%s", cursorReview)
	}
}

// TestInstallSlashCommandsPreservesUserOwnedFile pins the multi-command return
// semantics: a user-owned enhance.md stops that FILE only. The other commands
// still install, so the agent is still reported.
func TestInstallSlashCommandsPreservesUserOwnedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "my own enhance command\n"
	path := filepath.Join(dir, "enhance.md")
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]: review.md still installs", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Fatalf("user-owned file was modified:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(dir, "review.md")); err != nil {
		t.Fatalf("review.md missing despite user-owned enhance.md: %v", err)
	}
}

// TestInstallSlashCommandsPreservesUserOwnedReviewFile is the name
// collision case: "review" is generic enough that a user may already have
// their own. The marker logic must leave it byte-identical.
func TestInstallSlashCommandsPreservesUserOwnedReviewFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "my own review command, predating gx\n"
	path := filepath.Join(dir, "review.md")
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]: enhance.md still installs", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Fatalf("user-owned review.md was modified:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(dir, "enhance.md")); err != nil {
		t.Fatalf("enhance.md missing despite user-owned review.md: %v", err)
	}
}

func TestInstallSlashCommandsRefreshesManagedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"enhance.md", "review.md"} {
		stale := slashCommandMarker + "\n\nold template\n"
		if err := os.WriteFile(filepath.Join(dir, name), []byte(stale), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]", installed)
	}
	for name, fragment := range map[string]string{"enhance.md": "gx_enhance", "review.md": "gx_review"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "old template") {
			t.Fatalf("managed %s was not refreshed:\n%s", name, data)
		}
		if !strings.Contains(string(data), fragment) {
			t.Fatalf("refreshed %s missing template body:\n%s", name, data)
		}
	}
}

func TestInstallSlashCommandsRemovesRenamedPredecessor(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// A file this tool installed under the old name ("gx" was /enhance's name
	// until the rename), carrying our marker.
	legacy := filepath.Join(dir, "gx.md")
	if err := os.WriteFile(legacy, []byte(slashCommandMarker+"\n\nold command\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := installSlashCommands(); err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(legacy); !os.IsNotExist(statErr) {
		t.Error("renamed predecessor should be removed, leaving one command not two")
	}
	if _, err := os.Stat(filepath.Join(dir, "enhance.md")); err != nil {
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
	legacy := filepath.Join(dir, "gx.md")
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

// TestRemoveLegacySlashCommandsNeverDeletesALiveCommandName guards the
// generalized legacy check: a legacyNames entry that matches any CURRENT
// command name must never be deleted, because all commands share a directory.
func TestRemoveLegacySlashCommandsNeverDeletesALiveCommandName(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(dir, "review.md")
	if err := os.WriteFile(live, []byte(slashCommandMarker+"\n\nlive command\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	agent := slashCommandAgents()[0]
	removeLegacySlashCommands(home, agent, slashCommand{name: "enhance", legacyNames: []string{"review"}})
	if _, err := os.Stat(live); err != nil {
		t.Fatalf("a live command name was deleted via legacyNames: %v", err)
	}
}
