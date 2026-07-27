package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentsMDContainsSnippet(t *testing.T) {
	if !agentsMDContainsSnippet(agentsMDSnippet) {
		t.Fatal("expected canonical snippet to match marker")
	}
	if agentsMDContainsSnippet("# Agents\n\nUse git commit\n") {
		t.Fatal("unexpected match")
	}
	if !strings.Contains(agentsMDSnippet, "git push") {
		t.Fatal("expected agentsMDSnippet to require plain git push")
	}
	if !strings.Contains(agentsMDSnippet, "do not run `gx push`") {
		t.Fatal("expected agentsMDSnippet to forbid gx push")
	}
	if !strings.Contains(agentsMDSnippet, "gx capture push") {
		t.Fatal("expected agentsMDSnippet to forbid gx capture push")
	}
	for _, stale := range []string{"gx_edit", "gx base", "gx edit", "gx commit", "gx status"} {
		if strings.Contains(agentsMDSnippet, stale) {
			t.Fatalf("agentsMDSnippet contains removed command %q", stale)
		}
	}
	for _, want := range []string{
		"Version control: plain Git",
		"there is no GX save verb",
		"git commit -m",
		"`git commit --amend`",
		"preserve the GX revision trailer",
		"gx_review",
	} {
		if !strings.Contains(agentsMDSnippet, want) {
			t.Fatalf("agentsMDSnippet missing canonical guidance %q", want)
		}
	}
}

func TestAppendAgentsMDSnippetCreatesFile(t *testing.T) {
	repo := t.TempDir()
	appended, err := appendAgentsMDSnippet(repo, "")
	if err != nil {
		t.Fatal(err)
	}
	if !appended {
		t.Fatal("expected append")
	}
	data, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !agentsMDContainsSnippet(string(data)) {
		t.Fatalf("missing snippet:\n%s", data)
	}
}

func TestAppendAgentsMDSnippetPreservesExistingContent(t *testing.T) {
	repo := t.TempDir()
	path := filepath.Join(repo, "AGENTS.md")
	if err := os.WriteFile(path, []byte("# Custom\n\nKeep me.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := appendAgentsMDSnippet(repo, "# Custom\n\nKeep me.\n"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "Keep me.") {
		t.Fatalf("existing content lost:\n%s", content)
	}
	if !agentsMDContainsSnippet(content) {
		t.Fatalf("snippet missing:\n%s", content)
	}
}

func TestResolveMCPBinaryFromSibling(t *testing.T) {
	dir := t.TempDir()
	gxPath := filepath.Join(dir, "gx")
	mcpPath := filepath.Join(dir, "gx-mcp")
	if err := os.WriteFile(gxPath, []byte("gx"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mcpPath, []byte("mcp"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := resolveMCPBinary(gxPath)
	if err != nil {
		t.Fatal(err)
	}
	if got != mcpPath {
		t.Fatalf("resolveMCPBinary() = %q, want %q", got, mcpPath)
	}
}

func TestMCPLaunchCommand(t *testing.T) {
	got := mcpLaunchCommand("/bin/gx", "/bin/gx-mcp")
	want := []string{"env", "GX_BINARY=/bin/gx", "/bin/gx-mcp"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("arg[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestOfferInitAgentsMDSkipsWhenDeclined(t *testing.T) {
	repo := t.TempDir()
	var out bytes.Buffer
	err := offerInitAgentsMD(initSetupOptions{
		RepoRoot: repo,
		Quiet:    false,
		In:       strings.NewReader("n\n"),
		Out:      &out,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(filepath.Join(repo, "AGENTS.md")); statErr == nil {
		t.Fatal("expected AGENTS.md to remain absent")
	}
	if !strings.Contains(out.String(), "skipped") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestOfferInitAgentsMDQuietSkipsWrite(t *testing.T) {
	repo := t.TempDir()
	err := offerInitAgentsMD(initSetupOptions{
		RepoRoot: repo,
		Quiet:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(filepath.Join(repo, "AGENTS.md")); statErr == nil {
		t.Fatal("expected AGENTS.md to remain absent in quiet mode")
	}
}

func TestMergeCodexMCPServerAppendsBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	codexDir := filepath.Join(home, ".codex")
	if err := os.MkdirAll(codexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(codexDir, "config.toml")
	if err := os.WriteFile(configPath, []byte("# codex\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	updated, err := mergeCodexMCPServer("/bin/gx", "/bin/gx-mcp")
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Fatal("expected update")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "[mcp_servers.gx]") {
		t.Fatalf("missing codex block:\n%s", data)
	}
}
