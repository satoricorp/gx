package publication

import "testing"

func TestFormatPRSummaryLink(t *testing.T) {
	got := formatPRSummaryLink("lib.go:3", "https://github.com/example/repo/pull/1/files#diff-abcR3")
	want := `<a href="https://github.com/example/repo/pull/1/files#diff-abcR3" target="_blank" rel="noopener noreferrer">lib.go:3</a>`
	if got != want {
		t.Fatalf("formatPRSummaryLink() = %q, want %q", got, want)
	}
}

func TestFormatPRSummaryLinkEscapesHTML(t *testing.T) {
	got := formatPRSummaryLink(`a & "b"`, `https://example.com/?q=a&b=1`)
	if got != `<a href="https://example.com/?q=a&amp;b=1" target="_blank" rel="noopener noreferrer">a &amp; &#34;b&#34;</a>` {
		t.Fatalf("formatPRSummaryLink() = %q, want escaped HTML", got)
	}
}

func TestFormatPRSummaryLinkEmptyURL(t *testing.T) {
	if got := formatPRSummaryLink("plain text", ""); got != "plain text" {
		t.Fatalf("formatPRSummaryLink() = %q, want plain text", got)
	}
}
