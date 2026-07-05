package codereview

import (
	"strings"
	"testing"
)

func TestSecurityQualityHintFindingsOnChangedFile(t *testing.T) {
	findings := evaluateFindings(ReviewContext{
		ActiveScopes: []string{"security", "dependencies", "testing", "maintainability"},
		Sources:      []Source{{ID: "owasp-secure-coding", Scopes: []string{"security"}}},
		Options:      Options{PatchFocused: true},
		Brief: ReviewBrief{Static: StaticSnapshot{
			ChangedFiles: []string{"internal/store/users.go"},
			CodeQuality: []CodeQualityHint{{
				Kind:   "sql_interpolation",
				File:   "internal/store/users.go",
				Line:   5,
				Text:   `query := fmt.Sprintf("SELECT id FROM users WHERE id=%s", id)`,
				Reason: "SQL assembled with interpolation or concatenation is easy to turn into injection or quoting bugs; prefer parameterized queries.",
			}},
		}},
	}, defaultRules())
	if !hasFinding(findings, "security.quality-hints.1") {
		t.Fatalf("Findings = %#v, want security quality hint finding", findings)
	}
	var securityFinding Finding
	for _, finding := range findings {
		if finding.ID == "security.quality-hints.1" {
			securityFinding = finding
			break
		}
	}
	if securityFinding.Strength != "Strong" {
		t.Fatalf("Strength = %q, want Strong", securityFinding.Strength)
	}
	if !strings.Contains(evidenceText(securityFinding.Evidence), "internal/store/users.go:5") {
		t.Fatalf("Evidence = %#v, want changed-file anchor", securityFinding.Evidence)
	}
}

func TestSecurityQualityHintFindingsIgnoreUnchangedFile(t *testing.T) {
	findings := evaluateFindings(ReviewContext{
		ActiveScopes: []string{"security", "dependencies", "testing", "maintainability"},
		Sources:      []Source{{ID: "owasp-secure-coding", Scopes: []string{"security"}}},
		Options:      Options{PatchFocused: true},
		Brief: ReviewBrief{Static: StaticSnapshot{
			ChangedFiles: []string{"README.md"},
			CodeQuality: []CodeQualityHint{{
				Kind: "sql_interpolation",
				File: "internal/store/users.go",
				Line: 5,
				Text: `query := fmt.Sprintf("SELECT id FROM users WHERE id=%s", id)`,
			}},
		}},
	}, defaultRules())
	for _, finding := range findings {
		if strings.HasPrefix(finding.ID, "security.quality-hints") {
			t.Fatalf("Findings = %#v, did not expect security finding for unchanged file", findings)
		}
	}
}

func TestCollectCodeQualityHintsDetectsSecurityPatterns(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "internal/runner/run.go", strings.Join([]string{
		"package runner",
		"",
		"import \"os/exec\"",
		"",
		"func Run(input string) error {",
		"  return exec.Command(\"bash\", \"-c\", input).Run()",
		"}",
	}, "\n")+"\n")
	writeFile(t, root, "internal/validate/pattern.go", strings.Join([]string{
		"package validate",
		"",
		"import \"regexp\"",
		"",
		"func Match(input string) {",
		"  _ = regexp.MustCompile(input)",
		"}",
	}, "\n")+"\n")
	writeFile(t, root, "internal/config/config.go", "package config\n\nconst password = \"super-secret-password\"\n")

	facts := RepoFacts{Files: []string{
		"internal/runner/run.go",
		"internal/validate/pattern.go",
		"internal/config/config.go",
	}}
	hints := collectCodeQualityHints(root, facts, Options{})
	kinds := map[string]struct{}{}
	for _, hint := range hints {
		kinds[hint.Kind] = struct{}{}
	}
	for _, want := range []string{"shell_injection", "unsafe_regex", "hardcoded_secret"} {
		if _, ok := kinds[want]; !ok {
			t.Fatalf("hint kinds = %#v, want %q", kinds, want)
		}
	}
}
