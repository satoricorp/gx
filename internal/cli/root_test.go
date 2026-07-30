package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/totality/internal/cloud"
	"github.com/satoricorp/totality/internal/codereview"
	"github.com/satoricorp/totality/internal/vcs"
)

var errTestComposeRepair = errors.New("compose repair unavailable")

func TestVersionCommandPrintsLabeledVersion(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := strings.TrimSpace(out.String())
	if !strings.HasPrefix(text, "version ") {
		t.Fatalf("version output = %q, want compact version line", text)
	}
	if strings.Contains(text, "$ tx version") {
		t.Fatalf("version output should not echo command:\n%s", text)
	}
}

func TestVersionCommandPrintsJSON(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"version", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	var got struct {
		Version  string `json:"version"`
		Release  string `json:"release_version"`
		Revision string `json:"revision"`
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("version --json output is not JSON: %v\n%s", err, out.String())
	}
	if got.Version == "" {
		t.Fatalf("version --json = %#v, want non-empty version", got)
	}
}

func TestRootHelpPrintsAsciiLogoAtTop(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.HasPrefix(text, txLogoRaw+"\n") {
		t.Fatalf("root help should start with logo:\n%s", text)
	}
	if !strings.Contains(text, "$ tx help") {
		t.Fatalf("root help missing tx help invocation:\n%s", text)
	}
	if !strings.Contains(text, "Not signed in  tx auth login") {
		t.Fatalf("root help missing signed-out auth line:\n%s", text)
	}
	if !strings.Contains(text, "review (txr)") {
		t.Fatalf("root help missing alias %q:\n%s", "review (txr)", text)
	}
	if strings.Contains(text, "status") {
		t.Fatalf("root help should not offer a status command:\n%s", text)
	}
	if strings.Contains(text, "generate (tlg)") {
		t.Fatalf("root help should hide generate:\n%s", text)
	}
	if strings.Contains(text, "Shortcuts:") {
		t.Fatalf("root help should render aliases inline instead of a shortcut section:\n%s", text)
	}
}

func TestReviewAliasResolves(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"txr"})
	if err != nil || cmd == nil || cmd.Name() != "review" {
		t.Fatalf("Find(txr) = cmd=%v err=%v, want review command", cmd, err)
	}
}

func TestRootHelpShowsSignedInUser(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "octocat", GitHubAccessToken: "gho_saved", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	if text := out.String(); !strings.Contains(text, "Signed in as octocat") {
		t.Fatalf("root help missing signed-in auth line:\n%s", text)
	}
}

func TestRootHelpIgnoresLegacyLoginWithoutGitHubToken(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "api-key", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if strings.Contains(text, "Signed in as api-key") {
		t.Fatalf("root help showed stale legacy login:\n%s", text)
	}
	if !strings.Contains(text, "Not signed in  tx auth login") {
		t.Fatalf("root help missing signed-out auth line:\n%s", text)
	}
}

func TestRootHelpShowsStoredLoginWithCloudEnvPresent(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_CLOUD_URL", "http://localhost:3200/tx/pr")
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "joelachance", GitHubAccessToken: "gho_saved", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "Signed in as joelachance") {
		t.Fatalf("root help missing stored login auth line:\n%s", text)
	}
}

func TestDemoAndOpsCommandsAreRemoved(t *testing.T) {
	root := NewRoot(context.Background())
	for _, args := range [][]string{{"demo"}, {"ops"}, {"ops", "diagnose", "doctor"}, {"ops", "diag", "doctor"}} {
		cmd, _, err := root.Find(args)
		if err == nil && cmd != nil && cmd != root {
			t.Fatalf("Find(%v) resolved removed command %q", args, cmd.Name())
		}
	}
	for _, cmd := range root.Commands() {
		if cmd.Name() == "demo" || cmd.Name() == "ops" {
			t.Fatalf("root still registers removed command %q", cmd.Name())
		}
	}
	for _, group := range root.Groups() {
		if group.Title == "Advanced:" {
			t.Fatalf("root still declares the now-empty %q group", group.Title)
		}
	}
}

func TestDoctorAcceptsReportFlag(t *testing.T) {
	root := NewRoot(context.Background())
	doctor, _, err := root.Find([]string{"doctor"})
	if err != nil || doctor == nil || doctor.Name() != "doctor" {
		t.Fatalf("Find(doctor) = cmd=%v err=%v, want doctor", doctor, err)
	}
	flag := doctor.Flags().Lookup("report")
	if flag == nil {
		t.Fatal("tx doctor is missing the --report flag")
	}
	if flag.Value.Type() != "bool" {
		t.Fatalf("tx doctor --report type = %q, want bool", flag.Value.Type())
	}
}

func TestReportCommandStaysResolvableAsHiddenAlias(t *testing.T) {
	root := NewRoot(context.Background())
	report, _, err := root.Find([]string{"report"})
	if err != nil || report == nil || report.Name() != "report" {
		t.Fatalf("Find(report) = cmd=%v err=%v, want report", report, err)
	}
	if !report.Hidden {
		t.Fatal("tx report should be hidden now that tx doctor --report is the public spelling")
	}
	if report.GroupID != "" {
		t.Fatalf("tx report GroupID = %q, want no group", report.GroupID)
	}
}

func TestDoctorReportSendsDiagnosisWithLogs(t *testing.T) {
	repoRoot := initGitRepo(t)
	t.Chdir(repoRoot)
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{Login: "octocat", GitHubAccessToken: "gho_saved", ObtainedAt: time.Now()}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	var got cloud.ReportLogRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/reported-logs" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode report body: %v", err)
		}
		_, _ = w.Write([]byte(`{"id":"report-9","url":"https://totality.sh/reports/report-9"}`))
	}))
	defer server.Close()
	t.Setenv("TOTALITY_CLOUD_URL", server.URL)
	// `tx doctor` checks the saved GitHub token against GitHub's user endpoint,
	// which sent this test to api.github.com on every run with a token it had
	// just invented. The endpoint is a separate override from TOTALITY_GITHUB_API_URL,
	// so pointing that one at a fake is not enough.
	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"login":"octocat","id":583231}`))
	}))
	defer github.Close()
	t.Setenv("TOTALITY_GITHUB_USER_URL", github.URL)

	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"doctor", "--report"})

	if err := root.Execute(); err != nil {
		t.Fatalf("tx doctor --report error = %v\n%s", err, out.String())
	}

	var diagnosis string
	for _, log := range got.Logs {
		if log.Path == "totality-doctor.txt" {
			diagnosis = log.Content
		}
	}
	if diagnosis == "" {
		t.Fatalf("report logs missing the doctor diagnosis: %#v", got.Logs)
	}
	if !strings.Contains(diagnosis, "Publish outbox") {
		t.Fatalf("doctor diagnosis attachment = %q, want the printed doctor output", diagnosis)
	}
	text := out.String()
	if !strings.Contains(text, "Publish outbox") {
		t.Fatalf("tx doctor --report stopped printing the diagnosis:\n%s", text)
	}
	if !strings.Contains(text, "Report sent") || !strings.Contains(text, "report-9") {
		t.Fatalf("tx doctor --report missing send confirmation:\n%s", text)
	}
}

func TestPRAliasIsRemoved(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"pr"})
	if err == nil && cmd != nil && cmd.Name() == "pr" {
		t.Fatalf("Find(pr) resolved removed alias")
	}
}

func TestReviewCommandUsesDefaults(t *testing.T) {
	root := initGitRepo(t)
	writeTestFile(t, root, "README.md", "# repo\n")
	writeTestFile(t, root, "AGENTS.md", "# agents\n")
	writeTestFile(t, root, "go.mod", "module example.com/repo\n")
	writeTestFile(t, root, "main_test.go", "package main\n")
	gitAddTestFiles(t, root, "README.md", "AGENTS.md", "go.mod", "main_test.go")
	t.Chdir(root)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_API_URL", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_REVIEW_AI", "0")
	t.Setenv("TOTALITY_REVIEW_JUDGE", "0")
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("TOTALITY_REVIEW_RESOURCES", "0")
	t.Setenv("TOTALITY_REVIEW_INDEXED_CONTEXT", "0")

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"review"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("tx review error = %v\n%s", err, out.String())
	}
	text := out.String()
	for _, want := range []string{
		"## Recommendations",
		"- No material issues found in this change.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("tx review output missing %q in:\n%s", want, text)
		}
	}
	for _, unwanted := range []string{"brief", "PR Summary"} {
		if strings.Contains(strings.ToLower(text), unwanted) {
			t.Fatalf("tx review output should not include %q:\n%s", unwanted, text)
		}
	}
}

func writeTestTotalityConfig(t *testing.T, totalityHome, name, email string) {
	t.Helper()
	if err := os.MkdirAll(totalityHome, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := fmt.Sprintf(`{"user":{"name":%q,"email":%q}}`, name, email)
	if err := os.WriteFile(filepath.Join(totalityHome, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

// narrowPathToGit leaves git on PATH and nothing else.
//
// `tx init` registers the Totality MCP server with every agent CLI it finds installed
// by running it (`claude mcp add tx …`, `cursor mcp add tx …`). On a developer's
// machine those resolve to the real binaries, so this test was launching the
// developer's own Claude Code — which calls home to api.anthropic.com — as a
// side effect of asserting that `tx init -y` prints nothing. What the test means
// by a default machine is one with no agent CLI installed, and this is how to
// say that rather than inherit whatever the author happened to have.
func narrowPathToGit(t *testing.T) {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("locate git: %v", err)
	}
	dir := t.TempDir()
	if err := os.Symlink(git, filepath.Join(dir, "git")); err != nil {
		t.Fatalf("link git into the test PATH: %v", err)
	}
	t.Setenv("PATH", dir)
}

func TestInitYesAcceptsDefaultsAndSuppressesOutput(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	t.Chdir(root)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	narrowPathToGit(t)

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"init", "-y"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("tx init -y error = %v\n%s", err, out.String())
	}
	if out.String() != "" {
		t.Fatalf("tx init -y output = %q, want empty", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		t.Fatalf("after tx init -y .git missing: %v", err)
	}
}

func TestReviewCommandAcceptsScopeFlag(t *testing.T) {
	root := initGitRepo(t)
	t.Chdir(root)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_API_URL", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_REVIEW_AI", "0")
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("TOTALITY_REVIEW_RESOURCES", "0")
	t.Setenv("TOTALITY_REVIEW_INDEXED_CONTEXT", "0")

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"review", "--scope", "architecture"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tx review --scope error = %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "## Recommendations") {
		t.Fatalf("tx review --scope output missing recommendations:\n%s", out.String())
	}
}

func TestPostReviewSummaryCommentUpsertsGitHubPRComment(t *testing.T) {
	remote := "https://github.com/acme/tx.git"
	branch := "feature/demo"
	var gotCommentBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/pulls":
			if r.URL.Query().Get("head") != "acme:"+branch {
				t.Fatalf("head query = %q", r.URL.Query().Get("head"))
			}
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/tx/pull/7"}]`))
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode comment body: %v", err)
			}
			gotCommentBody = payload["body"]
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/tx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	t.Setenv("TOTALITY_REVIEW_AI", "0")
	t.Setenv("TOTALITY_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("TOTALITY_REVIEW_RESOURCES", "0")
	t.Setenv("TOTALITY_REVIEW_INDEXED_CONTEXT", "0")

	report := codereview.Report{
		Findings: []codereview.Finding{{
			ID:               "docs.missing-readme",
			Title:            "Missing README",
			Summary:          "The repo lacks the standard entrypoint document new maintainers expect first.",
			Benefit:          "Improves onboarding speed by giving humans and agents one place to find setup, purpose, and common commands.",
			Recommendation:   "Add a concise README.",
			Strength:         "Strong",
			SourcePublishers: []string{"Go project"},
		}},
		Sources: []codereview.Source{{
			ID:        "go-code-review-comments",
			Title:     "Go Code Review Comments",
			URL:       "https://go.dev/wiki/CodeReviewComments",
			Publisher: "Go project",
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)

	if gotCommentBody == "" {
		t.Fatal("postReviewSummaryComment() did not send a comment")
	}
	for _, want := range []string{"<!-- tx review summary -->", "## Recommendations", "**Informed by:** Go project"} {
		if !strings.Contains(gotCommentBody, want) {
			t.Fatalf("comment body missing %q:\n%s", want, gotCommentBody)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("postReviewSummaryComment() wrote warnings:\n%s", stderr.String())
	}
}

func TestPostReviewSummaryCommentAttemptsInlineCommentForValidAnchor(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test User")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeTestFile(t, root, "main.go", "package main\nfunc main() {}\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "initial")
	writeTestFile(t, root, "main.go", "package main\nfunc main() { println(\"hi\") }\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "change")

	remote := "https://github.com/acme/tx.git"
	branch := "feature/demo"
	var inlineAttempted bool
	var summaryAttempted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/pulls":
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/tx/pull/7"}]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/tx/pulls/7/comments":
			inlineAttempted = true
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode inline comment body: %v", err)
			}
			if payload["path"] != "main.go" || int(payload["line"].(float64)) != 2 {
				t.Fatalf("inline payload = %#v", payload)
			}
			_, _ = w.Write([]byte(`{"id":99}`))
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			summaryAttempted = true
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/tx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")

	report := codereview.Report{
		RepoRoot: root,
		Findings: []codereview.Finding{{
			ID:             "ai.review.1",
			Title:          "Inline finding",
			Summary:        "main.go changed.",
			Benefit:        "Keeps comments anchored.",
			Recommendation: "Fix main.go.",
			Anchors:        []codereview.FindingAnchor{{File: "main.go", Line: 2}},
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RootPath:   root,
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)
	if !inlineAttempted {
		t.Fatal("postReviewSummaryComment() did not attempt inline comment")
	}
	if !summaryAttempted {
		t.Fatal("postReviewSummaryComment() did not post summary fallback")
	}
	if stderr.Len() != 0 {
		t.Fatalf("postReviewSummaryComment() wrote warnings:\n%s", stderr.String())
	}
}

func TestPostReviewSummaryCommentFallsBackWhenInlineCommentFails(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test User")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeTestFile(t, root, "main.go", "package main\nfunc main() {}\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "initial")
	writeTestFile(t, root, "main.go", "package main\nfunc main() { println(\"hi\") }\n")
	gitAddTestFiles(t, root, "main.go")
	runGitTest(t, root, "commit", "-m", "change")

	remote := "https://github.com/acme/tx.git"
	branch := "feature/demo"
	var summaryAttempted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/pulls":
			_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/acme/tx/pull/7"}]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/tx/pulls/7/comments":
			http.Error(w, "line cannot be commented", http.StatusUnprocessableEntity)
		case r.Method == http.MethodGet && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/acme/tx/issues/7/comments":
			summaryAttempted = true
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/acme/tx/pull/7#issuecomment-12","body":"ok"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")

	report := codereview.Report{
		RepoRoot: root,
		Findings: []codereview.Finding{{
			ID:             "ai.review.1",
			Title:          "Inline finding",
			Summary:        "main.go changed.",
			Benefit:        "Keeps comments anchored.",
			Recommendation: "Fix main.go.",
			Anchors:        []codereview.FindingAnchor{{File: "main.go", Line: 2}},
		}},
	}
	var stderr bytes.Buffer
	postReviewSummaryComment(context.Background(), vcs.RepoInfo{
		RootPath:   root,
		RemoteURL:  &remote,
		BranchName: &branch,
	}, report, &stderr)
	if !summaryAttempted {
		t.Fatal("postReviewSummaryComment() did not post summary after inline failure")
	}
	if !strings.Contains(stderr.String(), "Could not post Totality inline review comment") {
		t.Fatalf("postReviewSummaryComment() warning = %q, want inline failure warning", stderr.String())
	}
}

func TestReviewCommandRejectsMultiplePrompts(t *testing.T) {
	cmd := NewRoot(context.Background())
	cmd.SetArgs([]string{"review", "one prompt", "second prompt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("tx review accepted multiple positional prompts")
	}
	if !strings.Contains(err.Error(), "accepts at most 1 arg") {
		t.Fatalf("tx review error = %v, want maximum arg error", err)
	}
}

func initGitRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-b", "main")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init error = %v\n%s", err, out)
	}
	return root
}

func writeTestFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func gitAddTestFiles(t *testing.T, root string, files ...string) {
	t.Helper()
	args := append([]string{"add"}, files...)
	runGitTest(t, root, args...)
}

func runGitTest(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s error = %v\n%s", strings.Join(args, " "), err, out)
	}
}
