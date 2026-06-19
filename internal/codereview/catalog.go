package codereview

type Source struct {
	ID     string
	Title  string
	URL    string
	Scopes []string
}

type StaticCatalog struct{}

var sourceCatalog = []Source{
	{
		ID:     "owasp-asvs",
		Title:  "OWASP Application Security Verification Standard",
		URL:    "https://owasp.org/www-project-application-security-verification-standard/",
		Scopes: []string{"security"},
	},
	{
		ID:     "owasp-secure-coding",
		Title:  "OWASP Secure Coding Practices",
		URL:    "https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/stable-en",
		Scopes: []string{"security"},
	},
	{
		ID:     "nist-ssdf",
		Title:  "NIST Secure Software Development Framework",
		URL:    "https://csrc.nist.gov/pubs/sp/800/218/final",
		Scopes: []string{"security", "dependencies", "maintainability"},
	},
	{
		ID:     "openssf-scorecard",
		Title:  "OpenSSF Scorecard",
		URL:    "https://openssf.org/scorecard/",
		Scopes: []string{"dependencies", "security"},
	},
	{
		ID:     "openssf-best-practices",
		Title:  "OpenSSF Best Practices",
		URL:    "https://openssf.org/best-practices-badge/",
		Scopes: []string{"dependencies", "maintainability"},
	},
	{
		ID:     "slsa",
		Title:  "SLSA",
		URL:    "https://slsa.dev/spec/v1.2/about",
		Scopes: []string{"dependencies", "security"},
	},
	{
		ID:     "core-web-vitals",
		Title:  "Core Web Vitals",
		URL:    "https://web.dev/articles/vitals",
		Scopes: []string{"performance"},
	},
	{
		ID:     "go-diagnostics",
		Title:  "Go Diagnostics",
		URL:    "https://go.dev/doc/diagnostics.html",
		Scopes: []string{"performance"},
	},
	{
		ID:     "google-eng-practices",
		Title:  "Google Engineering Practices",
		URL:    "https://google.github.io/eng-practices/",
		Scopes: []string{"maintainability", "testing"},
	},
	{
		ID:     "go-code-review-comments",
		Title:  "Go Code Review Comments",
		URL:    "https://go.dev/wiki/CodeReviewComments",
		Scopes: []string{"maintainability", "architecture"},
	},
	{
		ID:     "go-package-names",
		Title:  "Go Package Names",
		URL:    "https://go.dev/blog/package-names",
		Scopes: []string{"architecture", "maintainability"},
	},
	{
		ID:     "google-doc-style",
		Title:  "Google Developer Documentation Style Guide",
		URL:    "https://developers.google.com/style/",
		Scopes: []string{"docs", "onboarding"},
	},
	{
		ID:     "write-the-docs-guide",
		Title:  "Write the Docs Guide",
		URL:    "https://www.writethedocs.org/guide/",
		Scopes: []string{"docs", "onboarding"},
	},
	{
		ID:     "diataxis",
		Title:  "Diataxis",
		URL:    "https://diataxis.fr/",
		Scopes: []string{"docs", "onboarding", "architecture"},
	},
	{
		ID:     "fowler-test-pyramid",
		Title:  "Practical Test Pyramid",
		URL:    "https://martinfowler.com/articles/practical-test-pyramid.html",
		Scopes: []string{"testing"},
	},
	{
		ID:     "fowler-architecture",
		Title:  "Martin Fowler Architecture Guide",
		URL:    "https://www.martinfowler.com/architecture/",
		Scopes: []string{"architecture"},
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

func sourcesForScopes(scopes []string) []Source {
	return StaticCatalog{}.SourcesForScopes(scopes)
}
