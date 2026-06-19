package review_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/review"
)

func TestReviewOfflineSevenSections(t *testing.T) {
	repo := initReviewGitRepo(t)
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")

	report, err := review.Review(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Offline {
		t.Fatal("expected offline report")
	}
	names := sectionNames(report)
	for _, want := range []string{
		review.SectionIntentMatch,
		review.SectionStructure,
		review.SectionRules,
		review.SectionBlastRadius,
		review.SectionCollisions,
		review.SectionDoneNess,
		review.SectionBugPass,
	} {
		if !contains(names, want) {
			t.Fatalf("missing section %q in %#v", want, names)
		}
	}
	if rules := findSection(report, review.SectionRules); rules.Grade != review.GradeSkip {
		t.Fatalf("rules grade = %q, want skip offline", rules.Grade)
	}
	text := review.RenderMarkdown(report)
	if strings.Contains(strings.ToLower(text), "brief") {
		t.Fatalf("review output must not contain 'brief'")
	}
}

func TestReviewOnlineRulesFromServer(t *testing.T) {
	repo := initReviewGitRepo(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/review/context" {
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("repoRoot"); got != repo {
			http.Error(w, "missing repoRoot", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"rules": []map[string]any{{
				"id":       "r1",
				"ruleText": "no panic in handlers",
				"strength": "binding",
				"scopeExpr": "main.go",
				"status":   "enforced",
			}},
			"collisions": []map[string]any{{
				"file":   "main.go",
				"detail": "agent and human edited same hunk",
			}},
		})
	}))
	defer server.Close()

	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", server.URL)
	t.Setenv("GX_UPLOAD_TOKEN", "test-token")

	report, err := review.Review(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	if report.Offline {
		t.Fatal("expected online report")
	}
	rules := findSection(report, review.SectionRules)
	if rules.Grade != review.GradeFail {
		t.Fatalf("rules grade = %q, want fail", rules.Grade)
	}
	collisions := findSection(report, review.SectionCollisions)
	if collisions.Grade != review.GradeWarn {
		t.Fatalf("collisions grade = %q, want warn", collisions.Grade)
	}
}

func TestReviewOnlineServerShape(t *testing.T) {
	repo := initReviewGitRepo(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/review/context" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"rules": []map[string]any{{
				"id":        "r1",
				"ruleText":  "never use var",
				"strength":  "binding",
				"scopeExpr": "server/src/",
				"status":    "enforced",
			}},
			"collisions": []map[string]any{{
				"file":   "main.go",
				"kind":   "hunk_overlap",
				"detail": "Multiple sessions touched this file",
			}},
			"indexSnippets": []map[string]any{{
				"id":   "idx1",
				"text": "GX hunk link for review context.",
			}},
		})
	}))
	defer server.Close()

	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", server.URL)
	t.Setenv("GX_UPLOAD_TOKEN", "test-token")

	report, err := review.Review(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	if report.Offline {
		t.Fatal("expected online report")
	}
	rules := findSection(report, review.SectionRules)
	if rules.Grade != review.GradeFail {
		t.Fatalf("rules grade = %q, want fail", rules.Grade)
	}
	if !strings.Contains(strings.Join(rules.Evidence, " "), "never use var") {
		t.Fatalf("rules evidence = %#v", rules.Evidence)
	}
}

func initReviewGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "main.go")
	runGit(t, dir, "commit", "-m", "add main")
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "main_test.go")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func sectionNames(report review.GradedReport) []string {
	var names []string
	for _, section := range report.Sections {
		names = append(names, section.Name)
	}
	return names
}

func findSection(report review.GradedReport, name string) review.Section {
	for _, section := range report.Sections {
		if section.Name == name {
			return section
		}
	}
	return review.Section{}
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
