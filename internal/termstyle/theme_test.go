package termstyle

import (
	"os"
	"strings"
	"testing"
)

func TestEnabledRespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if Enabled() {
		t.Fatal("expected color disabled when NO_COLOR is set")
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

func TestFormatStatusPrefixes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if got := FormatStatus("ok: ready"); got != "ok: ready" {
		t.Fatalf("FormatStatus() = %q", got)
	}
}

func TestRenderStackEmpty(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := RenderStack("main", 0, nil)
	if !strings.Contains(got, "No GX stack recorded yet") {
		t.Fatalf("RenderStack() = %q", got)
	}
}

func TestRenderDoctorTablePlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := RenderDoctorTable([]DoctorRow{{Check: "Proxy", Status: "ok: up"}})
	if !strings.Contains(got, "Proxy") || !strings.Contains(got, "ok: up") {
		t.Fatalf("RenderDoctorTable() = %q", got)
	}
}

func TestHyperlinkWithoutColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	url := "https://github.com/example/compare/main...feature"
	if got := Hyperlink(url, url); got != url {
		t.Fatalf("Hyperlink() = %q, want plain url", got)
	}
}

func TestEnabledDefault(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	t.Setenv("TERM", "xterm-256color")
	if !Enabled() {
		t.Fatal("expected color enabled in normal terminal")
	}
}
