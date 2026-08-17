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

	claude, err := os.ReadFile(filepath.Join(home, ".claude", "commands", "gx.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"description:", "$ARGUMENTS", "gx_review", slashCommandMarker} {
		if !strings.Contains(string(claude), fragment) {
			t.Fatalf("claude command missing %q:\n%s", fragment, claude)
		}
	}

	claudeGates, err := os.ReadFile(filepath.Join(home, ".claude", "commands", "gates.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"description:", "argument-hint:", "$ARGUMENTS", "gx_gates", "slash-gates", slashCommandMarker} {
		if !strings.Contains(string(claudeGates), fragment) {
			t.Fatalf("claude gates command missing %q:\n%s", fragment, claudeGates)
		}
	}

	codex, err := os.ReadFile(filepath.Join(home, ".codex", "prompts", "gx.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(codex), "---") {
		t.Fatalf("codex prompt should not have frontmatter:\n%s", codex)
	}
	if !strings.Contains(string(codex), "$ARGUMENTS") {
		t.Fatalf("codex prompt missing $ARGUMENTS:\n%s", codex)
	}

	codexGates, err := os.ReadFile(filepath.Join(home, ".codex", "prompts", "gates.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(codexGates), "---") {
		t.Fatalf("codex gates prompt should not have frontmatter:\n%s", codexGates)
	}
	for _, fragment := range []string{"$ARGUMENTS", "gx_gates"} {
		if !strings.Contains(string(codexGates), fragment) {
			t.Fatalf("codex gates prompt missing %q:\n%s", fragment, codexGates)
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
	for _, name := range []string{"gx.md", "gates.md"} {
		cursor, err := os.ReadFile(filepath.Join(home, ".cursor", "commands", name))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(cursor), "$ARGUMENTS") {
			t.Fatalf("cursor %s must not use $ARGUMENTS:\n%s", name, cursor)
		}
	}
	cursorGates, err := os.ReadFile(filepath.Join(home, ".cursor", "commands", "gates.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cursorGates), "gx_gates") {
		t.Fatalf("cursor gates command missing gx_gates:\n%s", cursorGates)
	}
}

// TestInstallSlashCommandsPreservesUserOwnedFile pins the multi-command return
// semantics: a user-owned gx.md stops that FILE only. The other commands still
// install, so the agent is still reported.
func TestInstallSlashCommandsPreservesUserOwnedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "my own review command\n"
	path := filepath.Join(dir, "gx.md")
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]: gates.md still installs", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Fatalf("user-owned file was modified:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(dir, "gates.md")); err != nil {
		t.Fatalf("gates.md missing despite user-owned gx.md: %v", err)
	}
}

// TestInstallSlashCommandsPreservesUserOwnedGatesFile is the name
// collision case: "gates" is generic enough that a user may already have
// their own. The marker logic must leave it byte-identical.
func TestInstallSlashCommandsPreservesUserOwnedGatesFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "my own gates command, predating gx\n"
	path := filepath.Join(dir, "gates.md")
	if err := os.WriteFile(path, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	installed, err := installSlashCommands()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(installed, ",") != "Claude Code" {
		t.Fatalf("installed = %v, want [Claude Code]: gx.md still installs", installed)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != custom {
		t.Fatalf("user-owned gates.md was modified:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(dir, "gx.md")); err != nil {
		t.Fatalf("gx.md missing despite user-owned gates.md: %v", err)
	}
}

func TestInstallSlashCommandsRefreshesManagedFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"gx.md", "gates.md"} {
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
	for name, fragment := range map[string]string{"gx.md": "gx_review", "gates.md": "gx_gates"} {
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
	if _, err := os.Stat(filepath.Join(dir, "gx.md")); err != nil {
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

// TestRemoveLegacySlashCommandsNeverDeletesALiveCommandName guards the
// generalized legacy check: a legacyNames entry that matches any CURRENT
// command name must never be deleted, because all commands share a directory.
func TestRemoveLegacySlashCommandsNeverDeletesALiveCommandName(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	live := filepath.Join(dir, "gates.md")
	if err := os.WriteFile(live, []byte(slashCommandMarker+"\n\nlive command\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	agent := slashCommandAgents()[0]
	removeLegacySlashCommands(home, agent, slashCommand{name: "gx", legacyNames: []string{"gates"}})
	if _, err := os.Stat(live); err != nil {
		t.Fatalf("a live command name was deleted via legacyNames: %v", err)
	}
}
