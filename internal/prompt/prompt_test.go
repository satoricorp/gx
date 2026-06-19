package prompt

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestUseHuhFalseForBuffer(t *testing.T) {
	if UseHuh(strings.NewReader("\n")) {
		t.Fatal("expected UseHuh false for non-terminal reader")
	}
}

func TestUseHuhFalseWhenPlainPrompts(t *testing.T) {
	t.Setenv("GX_PLAIN_PROMPTS", "1")
	if UseHuh(os.Stdin) {
		t.Fatal("expected UseHuh false when GX_PLAIN_PROMPTS is set")
	}
}

func TestModifyChangeLegacy(t *testing.T) {
	t.Setenv("GX_PLAIN_PROMPTS", "1")
	candidates := []ModifyCandidate{
		{ChangeID: "abc123change", CommitID: "11111111commit", Description: "first"},
		{ChangeID: "def456change", CommitID: "22222222commit", Description: "second"},
	}
	var out bytes.Buffer
	got, err := ModifyChange(strings.NewReader("2\n"), &out, candidates)
	if err != nil {
		t.Fatalf("ModifyChange() error = %v", err)
	}
	if got != "def456change" {
		t.Fatalf("ModifyChange() = %q, want def456change", got)
	}
	if !strings.Contains(out.String(), "Select change to edit") {
		t.Fatalf("missing header: %q", out.String())
	}
}

func TestModifyChangeLegacyCancel(t *testing.T) {
	t.Setenv("GX_PLAIN_PROMPTS", "1")
	candidates := []ModifyCandidate{{ChangeID: "abc", CommitID: "111", Description: "first"}}
	_, err := ModifyChange(strings.NewReader("q\n"), &bytes.Buffer{}, candidates)
	if err == nil || !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Fatalf("ModifyChange() error = %v, want cancel", err)
	}
}

func TestIdentityLegacy(t *testing.T) {
	t.Setenv("GX_PLAIN_PROMPTS", "1")
	var out bytes.Buffer
	name, email, err := Identity(strings.NewReader("Joe Example\njoe@example.com\n"), &out, true, true, "", "", nil)
	if err != nil {
		t.Fatalf("Identity() error = %v", err)
	}
	if name != "Joe Example" || email != "joe@example.com" {
		t.Fatalf("Identity() = %q <%s>", name, email)
	}
}

func TestRequiredValueLegacyDefault(t *testing.T) {
	t.Setenv("GX_PLAIN_PROMPTS", "1")
	var out bytes.Buffer
	value, err := requiredValueLegacy(strings.NewReader("\n"), &out, "gx name", "Joe Example")
	if err != nil {
		t.Fatalf("requiredValueLegacy() error = %v", err)
	}
	if value != "Joe Example" {
		t.Fatalf("requiredValueLegacy() = %q, want Joe Example", value)
	}
}
