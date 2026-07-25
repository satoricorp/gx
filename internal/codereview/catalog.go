package codereview

import (
	"net/url"
	"strings"
)

type Source struct {
	ID        string   `json:"id"`
	Title     string   `json:"title,omitempty"`
	URL       string   `json:"url,omitempty"`
	Publisher string   `json:"publisher,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

type StaticCatalog struct{}

var sourceCatalog = []Source{
	{
		ID:        "owasp-asvs",
		Title:     "OWASP Application Security Verification Standard",
		URL:       "https://owasp.org/www-project-application-security-verification-standard/",
		Publisher: "OWASP",
		Scopes:    []string{"security"},
	},
	{
		ID:        "owasp-secure-coding",
		Title:     "OWASP Secure Coding Practices",
		URL:       "https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/stable-en",
		Publisher: "OWASP",
		Scopes:    []string{"security"},
	},
	{
		ID:        "nist-ssdf",
		Title:     "NIST Secure Software Development Framework",
		URL:       "https://csrc.nist.gov/pubs/sp/800/218/final",
		Publisher: "NIST",
		Scopes:    []string{"security", "dependencies", "maintainability"},
	},
	{
		ID:        "openssf-scorecard",
		Title:     "OpenSSF Scorecard",
		URL:       "https://openssf.org/scorecard/",
		Publisher: "OpenSSF",
		Scopes:    []string{"dependencies", "security"},
	},
	{
		ID:        "openssf-best-practices",
		Title:     "OpenSSF Best Practices",
		URL:       "https://openssf.org/best-practices-badge/",
		Publisher: "OpenSSF",
		Scopes:    []string{"dependencies", "maintainability"},
	},
	{
		ID:        "slsa",
		Title:     "SLSA",
		URL:       "https://slsa.dev/spec/v1.2/about",
		Publisher: "SLSA",
		Scopes:    []string{"dependencies", "security"},
	},
	{
		ID:        "core-web-vitals",
		Title:     "Core Web Vitals",
		URL:       "https://web.dev/articles/vitals",
		Publisher: "web.dev",
		Scopes:    []string{"performance"},
	},
	{
		ID:        "go-diagnostics",
		Title:     "Go Diagnostics",
		URL:       "https://go.dev/doc/diagnostics.html",
		Publisher: "Go project",
		Scopes:    []string{"performance"},
	},
	{
		ID:        "google-eng-practices",
		Title:     "Google Engineering Practices",
		URL:       "https://google.github.io/eng-practices/",
		Publisher: "Google",
		Scopes:    []string{"maintainability", "testing"},
	},
	{
		ID:        "go-code-review-comments",
		Title:     "Go Code Review Comments",
		URL:       "https://go.dev/wiki/CodeReviewComments",
		Publisher: "Go project",
		Scopes:    []string{"maintainability", "architecture"},
	},
	{
		ID:        "go-package-names",
		Title:     "Go Package Names",
		URL:       "https://go.dev/blog/package-names",
		Publisher: "Go project",
		Scopes:    []string{"architecture", "maintainability"},
	},
	{
		ID:        "google-doc-style",
		Title:     "Google Developer Documentation Style Guide",
		URL:       "https://developers.google.com/style/",
		Publisher: "Google",
		Scopes:    []string{"docs", "onboarding"},
	},
	{
		ID:        "write-the-docs-guide",
		Title:     "Write the Docs Guide",
		URL:       "https://www.writethedocs.org/guide/",
		Publisher: "Write the Docs",
		Scopes:    []string{"docs", "onboarding"},
	},
	{
		ID:        "diataxis",
		Title:     "Diataxis",
		URL:       "https://diataxis.fr/",
		Publisher: "Diátaxis",
		Scopes:    []string{"docs", "onboarding", "architecture"},
	},
	{
		ID:        "fowler-test-pyramid",
		Title:     "Practical Test Pyramid",
		URL:       "https://martinfowler.com/articles/practical-test-pyramid.html",
		Publisher: "Martin Fowler",
		Scopes:    []string{"testing"},
	},
	{
		ID:        "fowler-architecture",
		Title:     "Martin Fowler Architecture Guide",
		URL:       "https://www.martinfowler.com/architecture/",
		Publisher: "Martin Fowler",
		Scopes:    []string{"architecture"},
	},
}

func (StaticCatalog) SourcesForScopes(scopes []string) []Source {
	wanted := map[string]struct{}{}
	for _, scope := range scopes {
		wanted[scope] = struct{}{}
	}
	seen := map[string]struct{}{}
	var out []Source
	for _, source := range sourceCatalog {
		for _, scope := range source.Scopes {
			if _, ok := wanted[scope]; !ok {
				continue
			}
			if _, ok := seen[source.ID]; ok {
				break
			}
			seen[source.ID] = struct{}{}
			out = append(out, source)
			break
		}
	}
	return out
}

func SourcesForScopes(scopes []string) []Source {
	return sourcesForScopes(scopes)
}

func SourceBriefs(sources []Source) []SourceBrief {
	return sourceBriefs(sources)
}

func sourcesForScopes(scopes []string) []Source {
	return StaticCatalog{}.SourcesForScopes(scopes)
}

func publisherFromURLHost(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := strings.ToLower(strings.TrimPrefix(parsed.Hostname(), "www."))
	switch {
	case host == "owasp.org" || strings.HasSuffix(host, ".owasp.org"):
		return "OWASP"
	case host == "csrc.nist.gov" || strings.HasSuffix(host, ".nist.gov"):
		return "NIST"
	case strings.Contains(host, "openssf.org") || strings.Contains(host, "scorecard.dev"):
		return "OpenSSF"
	case host == "slsa.dev" || strings.HasSuffix(host, ".slsa.dev"):
		return "SLSA"
	case host == "go.dev" || strings.HasSuffix(host, ".golang.org"):
		return "Go project"
	case strings.Contains(host, "google"):
		return "Google"
	case strings.Contains(host, "writethedocs.org"):
		return "Write the Docs"
	case strings.Contains(host, "diataxis.fr"):
		return "Diátaxis"
	case strings.Contains(host, "martinfowler.com"):
		return "Martin Fowler"
	case strings.Contains(host, "web.dev"):
		return "web.dev"
	default:
		return ""
	}
}
