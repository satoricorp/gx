package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/cli"
	"github.com/satoricorp/gx/internal/gxconfig"
)

func TestShouldLaunch(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "no args", args: nil, want: false},
		{name: "help flag", args: []string{"--help"}, want: false},
		{name: "help command", args: []string{"help"}, want: false},
		{name: "internal daemon command", args: []string{"__gx-daemon"}, want: false},
		{name: "version command", args: []string{"version"}, want: false},
		{name: "init command", args: []string{"init"}, want: false},
		{name: "demo command", args: []string{"demo"}, want: false},
		{name: "compose command", args: []string{"compose"}, want: false},
		{name: "add command", args: []string{"add", "-m", "hi"}, want: false},
		{name: "add interactive command", args: []string{"add", "--interactive", "-m", "hi"}, want: false},
		{name: "edit command", args: []string{"edit", "abc"}, want: false},
		{name: "status command", args: []string{"status"}, want: false},
		{name: "stacks command", args: []string{"stacks"}, want: false},
		{name: "publish command", args: []string{"publish"}, want: false},
		{name: "auth command", args: []string{"auth"}, want: false},
		{name: "ops command", args: []string{"ops"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldLaunch(tt.args); got != tt.want {
				t.Fatalf("shouldLaunch(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestRootExposesAuthCommand(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"auth", "login"})
	if err != nil || cmd == nil {
		t.Fatalf("auth login command not found: %v", err)
	}
}

func TestRootRemovesShortcutCommands(t *testing.T) {
	root := cli.NewRoot(context.Background())
	for _, name := range []string{"gxa", "gxt"} {
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

func TestRootKeepsEditUtilityHiddenAndRemovesModifyAlias(t *testing.T) {
	root := cli.NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"edit"})
	if err != nil || cmd == nil || cmd.Name() != "edit" || !cmd.Hidden {
		t.Fatalf("Find(edit) = cmd=%v hidden=%v err=%v, want hidden edit utility command", cmd, cmd != nil && cmd.Hidden, err)
	}
	if cmd, _, err := root.Find([]string{"modify"}); err == nil && cmd != nil && cmd.Name() == "modify" {
		t.Fatalf("Find(modify) resolved removed alias")
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

func TestRootKeepsAddAndBaseUtilityCommandsHidden(t *testing.T) {
	root := cli.NewRoot(context.Background())
	add, _, addErr := root.Find([]string{"add"})
	if addErr != nil || add == nil || !add.Hidden {
		t.Fatalf("Find(add) = cmd=%v hidden=%v err=%v, want hidden add utility command", add, add != nil && add.Hidden, addErr)
	}
	base, _, baseErr := root.Find([]string{"base"})
	if baseErr != nil || base == nil || !base.Hidden {
		t.Fatalf("Find(base) = cmd=%v hidden=%v err=%v, want hidden base command", base, base != nil && base.Hidden, baseErr)
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
	t.Setenv("GX_HOME", t.TempDir())
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
		"  review (gxr)",
		"  generate (gxg)",
		"  status (gxs)",
		"Ship:",
		"  push",
		"  sync",
		"Help:",
		"  doctor",
		"  report",
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
		"Ship:",
		"Help:",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("root help missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Available Commands:") {
		t.Fatalf("root help should group commands by type:\n%s", text)
	}
	for _, hidden := range []string{"  add ", "  edit ", "  demo ", "  ops ", "  login ", "  base ", "  pr ", "  switch ", "  stack ", "  demux "} {
		if strings.Contains(text, hidden) {
			t.Fatalf("root help should hide %q in:\n%s", hidden, text)
		}
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
	for _, name := range []string{"commit", "compose", "publish", "stacks", "codex", "claude"} {
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
	t.Setenv("GX_HOME", t.TempDir())
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
	if !bytes.Contains(out.Bytes(), []byte("Run `gx init` first.")) {
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
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}

	root := cli.NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() unexpected error = %v", err)
	}

	if bytes.Contains(out.Bytes(), []byte("Run `gx init` first.")) {
		t.Fatalf("unexpected init note: %q", out.String())
	}
}
