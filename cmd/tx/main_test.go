package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/cli"
	"github.com/satoricorp/totality/internal/totalityconfig"
)

func TestRootExposesAuthCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"auth", "login"})
	if err != nil || cmd == nil {
		t.Fatalf("auth login command not found: %v", err)
	}
}

func TestRootRemovesShortcutCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"tla", "tls", "tlt"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != nil && cmd.Name() == name {
			t.Fatalf("%s command resolved after removal", name)
		}
	}
}

func TestRootRemovesStackCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"stack"})
	if err == nil && cmd != nil && cmd.Name() == "stack" {
		t.Fatalf("Find(stack) resolved removed command")
	}
}

func TestRootRemovesEditAndModifyCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"edit", "modify"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != nil && cmd.Name() == name {
			t.Fatalf("Find(%s) resolved removed command", name)
		}
	}
}

func TestRootRemovesHiddenSwitchCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"switch"})
	if err == nil && cmd != nil && cmd.Name() == "switch" {
		t.Fatalf("Find(switch) resolved removed command")
	}
}

func TestStacksRemovesSwitchSubcommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"stacks", "switch"})
	if err == nil && cmd != nil && cmd.Name() == "switch" {
		t.Fatalf("Find(stacks switch) resolved removed command")
	}
}

func TestRootRemovesAddAndBaseCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"add", "base"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != nil && cmd.Name() == name {
			t.Fatalf("Find(%s) resolved removed command", name)
		}
	}
}

func TestRootRemovesHiddenPRCompatibilityCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"pr"})
	if err == nil && cmd != nil && cmd.Name() == "pr" {
		t.Fatalf("Find(pr) resolved removed command")
	}
}

func TestRootHelpShowsHumanCommandsAndHidesAgentCommands(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	root := cli.NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() unexpected error = %v", err)
	}

	text := out.String()
	ordered := []string{
		"Setup:",
		"  init",
		"  auth",
		"Work:",
		"  review (txr)",
		"Help:",
		"  doctor",
		"  version",
		"  help",
	}
	previous := -1
	for _, want := range ordered {
		index := strings.Index(text, want)
		if index < 0 {
			t.Fatalf("root help missing %q in:\n%s", want, text)
		}
		if index <= previous {
			t.Fatalf("root help renders %q out of order:\n%s", want, text)
		}
		previous = index
	}
	for _, want := range []string{
		"Setup:",
		"Work:",
		"Help:",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("root help missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Available Commands:") {
		t.Fatalf("root help should group commands by type:\n%s", text)
	}
	if strings.Contains(text, "Ship:") {
		t.Fatalf("root help should not render a Ship group after sync's retirement:\n%s", text)
	}
	if strings.Contains(text, "Advanced:") {
		t.Fatalf("root help should not render an Advanced group after ops' removal:\n%s", text)
	}
	for _, hidden := range []string{"  add ", "  edit ", "  ops ", "  demo ", "  report ", "  login ", "  base ", "  pr ", "  switch ", "  stack ", "  demux ", "  generate ", "  sync "} {
		if strings.Contains(text, hidden) {
			t.Fatalf("root help should hide %q in:\n%s", hidden, text)
		}
	}
}

func TestRootRemovesDemoAndOpsCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"demo", "ops"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != nil && cmd.Name() == name {
			t.Fatalf("Find(%s) resolved removed command", name)
		}
	}
}

func TestRootKeepsReportAsHiddenAlias(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"report"})
	if err != nil || cmd == nil || cmd.Name() != "report" {
		t.Fatalf("Find(report) = cmd=%v err=%v, want the hidden report alias", cmd, err)
	}
	if !cmd.Hidden {
		t.Fatal("tx report should be hidden after folding into tx doctor --report")
	}
}

func TestRootRemovesSyncCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"sync"})
	if err == nil && cmd != nil && cmd.Name() == "sync" {
		t.Fatalf("Find(sync) resolved removed command")
	}
}

func TestRootRemovesHiddenDemuxCompatibilityCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"demux"})
	if err == nil && cmd != nil && cmd.Name() == "demux" {
		t.Fatalf("Find(demux) resolved removed command")
	}
}

func TestRootDoesNotExposeRemovedLegacyCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"commit", "status", "compose", "publish", "stacks", "codex", "claude"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != root {
			t.Fatalf("unexpected legacy command exposed: %s", cmd.Name())
		}
	}
}

func TestRootDoesNotExposeCompletionCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	if cmd, _, err := root.Find([]string{"completion"}); err == nil && cmd != root {
		t.Fatalf("unexpected completion command exposed: %s", cmd.Name())
	}
}

func TestRootDoesNotExposeDaemonCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	if cmd, _, err := root.Find([]string{"daemon"}); err == nil && cmd != root {
		t.Fatalf("unexpected daemon command exposed: %s", cmd.Name())
	}
}

func TestRootDoesNotExposeGitOrJJPassthroughCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"git", "jj", "diff", "log", "split"} {
		if cmd, _, err := root.Find([]string{name}); err == nil && cmd != root {
			t.Fatalf("unexpected passthrough command exposed: %s", cmd.Name())
		}
	}
}

func TestRootPrintsInitNoteWhenIdentityMissing(t *testing.T) {
	wd := t.TempDir()
	prev, _ := os.Getwd()
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	root := cli.NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() unexpected error = %v", err)
	}

	text := out.String()
	if !bytes.Contains(out.Bytes(), []byte("Run `tx init` first.")) {
		t.Fatalf("missing init note: %q", text)
	}
}

func TestRootSkipsInitNoteWhenIdentityConfigured(t *testing.T) {
	wd := t.TempDir()
	prev, _ := os.Getwd()
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	if err := totalityconfig.Save(totalityconfig.Config{
		User: totalityconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("totalityconfig.Save() error = %v", err)
	}

	root := cli.NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() unexpected error = %v", err)
	}

	if bytes.Contains(out.Bytes(), []byte("Run `tx init` first.")) {
		t.Fatalf("unexpected init note: %q", out.String())
	}
}
