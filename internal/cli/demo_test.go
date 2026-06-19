package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestDemoCommandWalksThroughRequestedScreens(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetIn(strings.NewReader(strings.Join([]string{
		"",
		`gx add -m "initial change for demo"`,
		"gx stacks",
		"",
		"gx compose",
		"Shift + A",
		"gx publish",
		"",
		"",
		"",
	}, "\n") + "\n"))
	root.SetArgs([]string{"demo"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	for _, want := range []string{
		"Welcome to gx!",
		"gx creates stacks and revisions autonomously.",
		`gx add -m "initial change for demo"`,
		"Revision recorded",
		"gx stacks",
		"This command is interactive",
		"gx compose",
		"No more `gx add`!",
		"Shift + A",
		"gx publish",
		"https://gx.run/reviews/demo-stack",
		"> Cursor",
		"> Codex",
		"> Claude Code",
		"> Opencode",
		"> Antigravity",
		"Congrats for finishing the gx demo!",
		"- Joe (Founder)",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("demo output missing %q in:\n%s", want, text)
		}
	}
}

func TestDemoCommandCanBeCanceled(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer

	err := runDemo(strings.NewReader("q\n"), &out)
	if err == nil || !strings.Contains(err.Error(), "demo canceled") {
		t.Fatalf("runDemo() error = %v, want demo canceled", err)
	}
}
