package termstyle

import (
	"os"
	"strings"
	"testing"
)

// pinTTY substitutes the stdout-is-a-terminal answer for one test.
func pinTTY(t *testing.T, isTTY bool) {
	t.Helper()
	prev := stdoutIsTerminal
	stdoutIsTerminal = func() bool { return isTTY }
	t.Cleanup(func() { stdoutIsTerminal = prev })
}

func TestEnabledRespectsNoColor(t *testing.T) {
	pinTTY(t, true)
	t.Setenv("NO_COLOR", "1")
	if Enabled() {
		t.Fatal("expected color disabled when NO_COLOR is set")
	}
}

func TestEnabledNoColorBeatsForceColor(t *testing.T) {
	pinTTY(t, true)
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")
	if Enabled() {
		t.Fatal("NO_COLOR must win over FORCE_COLOR")
	}
}

func TestLabelValuePlainWithoutColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := LabelValue("Proxy", "http://127.0.0.1:8787")
	if !strings.Contains(got, "Proxy") || !strings.Contains(got, "http://127.0.0.1:8787") {
		t.Fatalf("LabelValue() = %q", got)
	}
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("LabelValue() should not contain ANSI when NO_COLOR is set, got %q", got)
	}
}

func TestHyperlinkWithoutColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	url := "https://github.com/example/compare/main...feature"
	if got := Hyperlink(url, url); got != url {
		t.Fatalf("Hyperlink() = %q, want plain url", got)
	}
}

func TestEnabledOnTerminal(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FORCE_COLOR")
	t.Setenv("TERM", "xterm-256color")
	pinTTY(t, true)
	if !Enabled() {
		t.Fatal("expected color enabled on a terminal")
	}
}

// The gap this fixes: `gx review | tee out.txt` used to write ANSI into the
// file because nothing checked the fd.
func TestEnabledOffWhenPiped(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("FORCE_COLOR")
	t.Setenv("TERM", "xterm-256color")
	pinTTY(t, false)
	if Enabled() {
		t.Fatal("expected color disabled when stdout is not a terminal")
	}
}

func TestEnabledForceColorThroughPipe(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("FORCE_COLOR", "1")
	pinTTY(t, false)
	if !Enabled() {
		t.Fatal("FORCE_COLOR must enable color through a pipe")
	}
	t.Setenv("FORCE_COLOR", "0")
	if Enabled() {
		t.Fatal("FORCE_COLOR=0 must not force color")
	}
}

func TestEnabledDumbTerm(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("TERM", "dumb")
	pinTTY(t, true)
	if Enabled() {
		t.Fatal("TERM=dumb must disable color even under FORCE_COLOR")
	}
}
