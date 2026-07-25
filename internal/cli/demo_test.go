package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestDemoCommandWalksThroughRequestedScreens(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	demoSkipExec = true
	t.Cleanup(func() { demoSkipExec = false })

	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	// Welcome + 5 steps each wait for Enter.
	root.SetIn(strings.NewReader(strings.Repeat("\n", 6)))
	root.SetArgs([]string{"demo"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"gx — version control that captures how code was made",
		"5 steps · 2 minutes · runs in a scratch repo, your work is untouched",
		"── 1/5 · init",
		"Set up gx in a repository.",
		"$ gx init",
		"Installs hooks and starts ambient capture",
		"── 2/5 · sign in",
		"Connect to gx cloud.",
		"$ gx auth login",
		"GitHub device login",
		"── 3/5 · stage",
		"Choose what goes in — plain git.",
		"$ git add -p",
		"gx doesn't replace staging",
		"── 4/5 · commit",
		"Commit with git — nothing new to learn.",
		`$ git commit -m "add rate limiter"`,
		"The hooks gx installed record the revision as you commit",
		"── 5/5 · push",
		"Ship it.",
		"$ git push",
		"Your branch goes up like always",
		"── done",
		"That's the loop: stage and commit with git, gx records the rest.",
		"git status",
		"gx review",
		"gx doctor",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("demo output missing %q in:\n%s", want, text)
		}
	}
}

func TestDemoCommandCanBeCanceled(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	demoSkipExec = true
	t.Cleanup(func() { demoSkipExec = false })

	var out bytes.Buffer
	err := runDemo(context.Background(), strings.NewReader("q\n"), &out)
	if err == nil || !strings.Contains(err.Error(), "demo canceled") {
		t.Fatalf("runDemo() error = %v, want demo canceled", err)
	}
}

func TestDemoCommandIsVisibleInHelp(t *testing.T) {
	root := NewRoot(context.Background())
	demo, _, err := root.Find([]string{"demo"})
	if err != nil {
		t.Fatalf("Find(demo) error = %v", err)
	}
	if demo.Hidden {
		t.Fatal("demo command should be visible")
	}
	if demo.GroupID != groupSetup {
		t.Fatalf("demo GroupID = %q, want %q", demo.GroupID, groupSetup)
	}
}
