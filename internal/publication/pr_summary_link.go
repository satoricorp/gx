package publication

import (
	"fmt"
	"html"
	"strings"
)

func formatPRSummaryLink(text, url string) string {
	text = strings.TrimSpace(text)
	url = strings.TrimSpace(url)
	if url == "" {
		return text
	}
	if text == "" {
		return url
	}
	return fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`,
		html.EscapeString(url), html.EscapeString(text))
}
