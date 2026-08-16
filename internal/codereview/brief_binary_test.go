package codereview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileContentSnippetDescribesBinariesInsteadOfSendingThem(t *testing.T) {
	root := t.TempDir()
	// A Mach-O 64-bit magic number followed by NULs and junk — what an
	// untracked compiled Go binary looks like to os.ReadFile.
	macho := append([]byte{0xcf, 0xfa, 0xed, 0xfe}, make([]byte, 600)...)
	macho = append(macho, []byte("__PAGEZERO__TEXT")...)
	if err := os.WriteFile(filepath.Join(root, "tennis"), macho, 0o755); err != nil {
		t.Fatal(err)
	}
	got := fileContentSnippet(root, "tennis")
	if !strings.HasPrefix(got, diffUnavailableContentHeader) {
		t.Fatalf("binary snippet must keep the content header:\n%s", got)
	}
	if !strings.Contains(got, "[binary file: Mach-O executable, ") || !strings.Contains(got, "no source diff to review") {
		t.Fatalf("binary snippet should name the kind and say there is no source:\n%s", got)
	}
	if strings.Contains(got, "__PAGEZERO") {
		t.Fatalf("binary bytes must not be sent to the model:\n%s", got)
	}
	if len(got) > 400 {
		t.Fatalf("binary snippet should be one short line, got %d bytes", len(got))
	}
}

func TestFileContentSnippetPassesTextThrough(t *testing.T) {
	root := t.TempDir()
	src := "package main\n\nfunc main() {}\n"
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	got := fileContentSnippet(root, "main.go")
	if !strings.HasSuffix(got, src) {
		t.Fatalf("text file should be sent verbatim:\n%s", got)
	}
}

func TestDescribeBinaryContentKinds(t *testing.T) {
	cases := map[string][]byte{
		"Mach-O executable":     append([]byte{0xcf, 0xfa, 0xed, 0xfe}, 0, 1, 2),
		"ELF executable":        append([]byte("\x7fELF"), 0, 1),
		"Windows PE executable": append([]byte("MZ"), 0, 0, 0),
		"PNG image":             append([]byte("\x89PNG"), 0, 1),
		"zip archive":           append([]byte("PK\x03\x04"), 0, 1),
		"gzip archive":          append([]byte{0x1f, 0x8b}, 0, 1),
		"PDF":                   append([]byte("%PDF-1.7"), 0),
		"binary":                {'a', 'b', 0, 'c'},
	}
	for want, data := range cases {
		kind, ok := describeBinaryContent(data)
		if !ok || kind != want {
			t.Errorf("describeBinaryContent(%q…) = (%q,%v), want (%q,true)", data[:min(4, len(data))], kind, ok, want)
		}
	}
	for _, text := range []string{"", "hello\n", "package main\n", "héllo wörld — ünïcode\n"} {
		if kind, ok := describeBinaryContent([]byte(text)); ok {
			t.Errorf("text %q wrongly classified as %q", text, kind)
		}
	}
}
