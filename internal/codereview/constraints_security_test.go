package codereview

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The fixture credentials are split mid-literal so secret scanners — GitHub
// push protection included — never see a contiguous key shape in this source
// file. The runtime values are whole, which is what the tests exercise.
var (
	fixtureAWSKey      = "AKIA" + "ABCDEFGHIJKLMNOP"
	fixtureGitHubToken = "ghp_" + "abcdefghijklmnopqrstuvwxyz123456"
)

func TestConstraintsSecretFindings(t *testing.T) {
	tests := []struct {
		name string
		file string
		diff string
		want int
	}{
		{
			name: "added AWS key",
			file: "internal/app/creds.go",
			diff: "@@ -0,0 +1,2 @@\n+package app\n+const key = \"" + fixtureAWSKey + "\"\n",
			want: 1,
		},
		{
			name: "removed key is the fix, not the leak",
			file: "internal/app/creds.go",
			diff: "@@ -1,2 +1,1 @@\n package app\n-const key = \"" + fixtureAWSKey + "\"\n",
			want: 0,
		},
		{
			name: "testdata is exempt",
			file: "internal/app/testdata/fixture.go",
			diff: "@@ -0,0 +1,1 @@\n+const key = \"" + fixtureAWSKey + "\"\n",
			want: 0,
		},
		{
			name: "test files plant fixture keys",
			file: "internal/app/creds_test.go",
			diff: "@@ -0,0 +1,1 @@\n+const key = \"" + fixtureAWSKey + "\"\n",
			want: 0,
		},
		{
			name: "placeholder lines are spared",
			file: "internal/app/config.go",
			diff: "@@ -0,0 +1,1 @@\n+apiKey := \"AKIAEXAMPLEEXAMPLE12\" // example only\n",
			want: 0,
		},
		{
			name: "github token",
			file: "internal/app/token.go",
			diff: "@@ -0,0 +1,1 @@\n+token := \"" + fixtureGitHubToken + "\"\n",
			want: 1,
		},
		{
			name: "untracked file content fallback is scanned",
			file: "internal/app/new.go",
			diff: diffUnavailableContentHeader + "\npackage app\n\nconst key = \"" + fixtureAWSKey + "\"\n",
			want: 1,
		},
		{
			name: "quoted password assignment",
			file: "internal/app/db.go",
			diff: "@@ -0,0 +1,1 @@\n+password = \"hunter2hunter2\"\n",
			want: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			diffs := []DiffSnippet{{File: test.file, Diff: test.diff}}
			addedByFile := map[string][]constraintsAddedLine{
				test.file: constraintsAddedLinesForDiff(test.diff),
			}
			findings := constraintsSecretFindings(diffs, addedByFile)
			if len(findings) != test.want {
				t.Fatalf("findings = %d, want %d: %#v", len(findings), test.want, findings)
			}
			if test.want > 0 {
				if findings[0].Strength != "Blocking" || findings[0].File != test.file {
					t.Fatalf("finding = %#v, want Blocking on %s", findings[0], test.file)
				}
			}
		})
	}
}

func TestConstraintsAuditEnvironmentFailure(t *testing.T) {
	tests := []struct {
		name    string
		tool    string
		output  string
		skipped bool
	}{
		{name: "offline govulncheck", tool: "govulncheck", output: "govulncheck: dial tcp: lookup vuln.go.dev: no such host", skipped: true},
		{name: "npm registry unreachable", tool: "npm audit", output: "npm ERR! code ENOTFOUND", skipped: true},
		{name: "cargo advisory db fetch", tool: "cargo audit", output: "error: couldn't fetch advisory database", skipped: true},
		{name: "real vulnerabilities stay failures", tool: "govulncheck", output: "Vulnerability #1: GO-2024-1234 in golang.org/x/net", skipped: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reason := constraintsAuditEnvironmentFailure(test.tool, test.output)
			if (reason != "") != test.skipped {
				t.Fatalf("constraintsAuditEnvironmentFailure(%q) = %q, want skipped=%v", test.output, reason, test.skipped)
			}
		})
	}
}

func TestCollectConstraintsAuditResultsReportsMissingAuditor(t *testing.T) {
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	// A PATH with no govulncheck: the auditor applies but cannot run, which
	// must surface as a skipped result rather than silence or a failure.
	tools := t.TempDir()
	t.Setenv("PATH", tools)

	results := collectConstraintsAuditResults(context.Background(), root, []string{"go.mod"})
	if len(results) != 1 {
		t.Fatalf("results = %#v, want exactly the missing govulncheck", results)
	}
	if !results[0].Skipped || !strings.Contains(results[0].Reason, "govulncheck not installed") {
		t.Fatalf("result = %#v, want skipped with the missing-tool reason", results[0])
	}
}

func TestCollectConstraintsAuditResultsRunsDetectedAuditor(t *testing.T) {
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	tools := t.TempDir()
	writeExecutable(t, filepath.Join(tools, "govulncheck"), "#!/bin/sh\necho 'Vulnerability #1: GO-2024-1234'\nexit 3\n")
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))

	results := collectConstraintsAuditResults(context.Background(), root, []string{"go.mod"})
	if len(results) != 1 {
		t.Fatalf("results = %#v, want one govulncheck run", results)
	}
	if results[0].Skipped || results[0].ExitCode == 0 {
		t.Fatalf("result = %#v, want a real failure with the vuln output", results[0])
	}
}

func TestCollectConstraintsAuditResultsHonorsKillSwitch(t *testing.T) {
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")

	if results := collectConstraintsAuditResults(context.Background(), root, []string{"go.mod"}); results != nil {
		t.Fatalf("results = %#v, want nil under GX_REVIEW_STATIC_TOOLS=0", results)
	}
}
