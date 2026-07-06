package publication

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

var (
	prSummaryBareFileDiffLinkRE = regexp.MustCompile(`#diff-[0-9a-f]{64}\)`)
	prSummaryVerdictBannerRE    = regexp.MustCompile(`(?m)^> [✅👀🔴] \*\*(No review needed|Quick scan|Requires Deep Review)\*\* — `)
)

// TODO: Source future golden cases from GX Cloud review history for live-model quality tracking.

var updatePRSummaryGoldens = flag.Bool("update", false, "regenerate internal/publication/testdata/prsummary/*/expected.md")

const reachGoldenHeadSHA = "__REACH_HEAD_SHA__"

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestPRSummaryGolden(t *testing.T) {
	root := filepath.Join("testdata", "prsummary")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", root, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		caseName := entry.Name()
		t.Run(caseName, func(t *testing.T) {
			runPRSummaryGoldenCase(t, filepath.Join(root, caseName))
		})
	}
}

func TestPRSummaryGoldenCorpusInvariants(t *testing.T) {
	root := filepath.Join("testdata", "prsummary")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("ReadDir(%s): %v", root, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(root, entry.Name(), "expected.md")
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			t.Fatalf("read %s: %v", path, err)
		}
		t.Run(entry.Name(), func(t *testing.T) {
			assertPRSummaryGoldenLayout(t, string(data))
		})
	}
}

func runPRSummaryGoldenCase(t *testing.T, dir string) {
	t.Helper()
	setup := loadGoldenSetup(t, dir)
	reachRoot, headSHA := prepareGoldenRepo(t, dir, setup)

	artifact := loadGoldenArtifact(t, filepath.Join(dir, "artifact.json"))
	if reachRoot != "" {
		artifact.Bundle.Repo.RootPath = reachRoot
		if headSHA != "" {
			artifact.Bundle.Push.HeadCommitID = headSHA
			for i := range artifact.Bundle.Stack {
				artifact.Bundle.Stack[i].Change.CurrentCommitID = headSHA
			}
		}
	} else if setup.ReviewMD {
		reviewRoot := t.TempDir()
		copyGoldenReviewMD(t, dir, reviewRoot)
		artifact.Bundle.Repo.RootPath = reviewRoot
	}

	if setup.BlastRadius == "1" {
		t.Setenv("GX_PR_BLAST_RADIUS", "1")
	} else {
		t.Setenv("GX_PR_BLAST_RADIUS", "0")
	}

	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	oldWait := prSummaryReviewRetryWait
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
		prSummaryReviewRetryWait = oldWait
	}()
	prSummaryReviewRetryWait = func(context.Context, time.Duration) error { return nil }

	reviewerCfg := loadGoldenReviewerConfig(t, dir)
	aiResponse := readGoldenOptionalFile(t, filepath.Join(dir, "ai_response.json"))
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		var failErr error
		if reviewerCfg.Error != "" {
			failErr = fmt.Errorf("%s", reviewerCfg.Error)
		}
		return &goldenPRSummaryReviewer{
			aiResponse: aiResponse,
			failCount:  reviewerCfg.FailCount,
			failErr:    failErr,
		}, codereview.ReviewerInfo{Models: reviewerCfg.Models}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return loadGoldenContext(t, dir), nil
	}

	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), artifact)
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if setup.ReachRepo && headSHA != "" {
		body = strings.ReplaceAll(body, headSHA, reachGoldenHeadSHA)
	}
	assertPRSummaryGoldenLayout(t, body)

	expectedPath := filepath.Join(dir, "expected.md")
	if *updatePRSummaryGoldens {
		if err := os.WriteFile(expectedPath, []byte(body), 0o644); err != nil {
			t.Fatalf("write expected.md: %v", err)
		}
		t.Logf("updated %s", expectedPath)
		return
	}

	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("read expected.md: %v", err)
	}
	if body != string(expected) {
		t.Fatalf("golden mismatch for %s\n--- got\n%s\n--- want\n%s", dir, body, string(expected))
	}
}

type goldenCaseSetup struct {
	ReachRepo   bool   `json:"reach_repo"`
	ReviewMD    bool   `json:"review_md"`
	BlastRadius string `json:"blast_radius"`
}

type goldenReviewerConfig struct {
	FailCount int      `json:"fail_count"`
	Error     string   `json:"error"`
	Models    []string `json:"models"`
}

func assertPRSummaryGoldenLayout(t *testing.T, body string) {
	t.Helper()
	if !strings.HasPrefix(body, githubPRBodyMarker) {
		t.Fatalf("body missing %q marker", githubPRBodyMarker)
	}
	if !prSummaryVerdictBannerRE.MatchString(body) {
		t.Fatalf("body missing verdict banner blockquote:\n%s", body)
	}
	if strings.Contains(body, "## Needs Review") {
		t.Fatalf("body contains deprecated ## Needs Review section:\n%s", body)
	}
	if strings.Contains(body, "No specific high-impact review targets") {
		t.Fatalf("body contains deprecated filler line:\n%s", body)
	}
	if prSummaryBareFileDiffLinkRE.MatchString(body) {
		t.Fatalf("body contains bare file-level diff link (missing R/L anchor):\n%s", body)
	}
	if !strings.Contains(body, "*Generated by GX") {
		t.Fatalf("body missing provenance footer:\n%s", body)
	}
	if idx := strings.Index(body, githubPRBodyMarker); idx >= 0 {
		footerIdx := strings.LastIndex(body, "*Generated by GX")
		if footerIdx >= 0 && strings.Contains(body[idx:footerIdx], "**Review verdict:") {
			t.Fatalf("body contains deprecated Review verdict line:\n%s", body)
		}
	}
	if idx := strings.Index(body, "## Blast Radius"); idx >= 0 {
		section := body[idx:]
		if end := strings.Index(section, "\n\n## "); end > 0 {
			section = section[:end]
		} else if end := strings.Index(section, "\n\n*Generated by GX"); end > 0 {
			section = section[:end]
		}
		lines := strings.Split(strings.TrimSpace(strings.TrimPrefix(section, "## Blast Radius")), "\n")
		if len(lines) < 2 {
			t.Fatalf("blast radius section missing narrative after quantitative line:\n%s", section)
		}
		narrative := strings.TrimSpace(lines[1])
		if narrative == "" || strings.HasPrefix(narrative, "- ") || strings.HasPrefix(narrative, "**Critical path") {
			t.Fatalf("blast radius section missing narrative prose:\n%s", section)
		}
	}
}

func loadGoldenSetup(t *testing.T, dir string) goldenCaseSetup {
	t.Helper()
	path := filepath.Join(dir, "setup.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return goldenCaseSetup{}
		}
		t.Fatalf("read setup.json: %v", err)
	}
	var setup goldenCaseSetup
	if err := json.Unmarshal(data, &setup); err != nil {
		t.Fatalf("decode setup.json: %v", err)
	}
	return setup
}

func loadGoldenReviewerConfig(t *testing.T, dir string) goldenReviewerConfig {
	t.Helper()
	path := filepath.Join(dir, "reviewer.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return goldenReviewerConfig{}
		}
		t.Fatalf("read reviewer.json: %v", err)
	}
	var cfg goldenReviewerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("decode reviewer.json: %v", err)
	}
	return cfg
}

func loadGoldenContext(t *testing.T, dir string) prSummaryContext {
	t.Helper()
	path := filepath.Join(dir, "context.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return prSummaryContext{}
		}
		t.Fatalf("read context.json: %v", err)
	}
	var ctx prSummaryContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		t.Fatalf("decode context.json: %v", err)
	}
	return ctx
}

func loadGoldenArtifact(t *testing.T, path string) reviewbundle.Artifact {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read artifact.json: %v", err)
	}
	var artifact reviewbundle.Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		t.Fatalf("decode artifact.json: %v", err)
	}
	return artifact
}

func readGoldenOptionalFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func prepareGoldenRepo(t *testing.T, dir string, setup goldenCaseSetup) (root, headSHA string) {
	t.Helper()
	if setup.ReachRepo {
		return initLexicalReachRepo(t)
	}
	return "", ""
}

func copyGoldenReviewMD(t *testing.T, dir, root string) {
	t.Helper()
	src := filepath.Join(dir, "REVIEW.md")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read REVIEW.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "REVIEW.md"), data, 0o644); err != nil {
		t.Fatalf("write REVIEW.md: %v", err)
	}
}

type goldenPRSummaryReviewer struct {
	aiResponse []byte
	failCount  int
	failErr    error
	attempts   int
}

func (r *goldenPRSummaryReviewer) Review(ctx context.Context, brief codereview.ReviewBrief) ([]codereview.Finding, error) {
	_, findings, err := r.ReviewWithOverview(ctx, brief)
	return findings, err
}

func (r *goldenPRSummaryReviewer) ReviewForSummary(ctx context.Context, brief codereview.ReviewBrief) (codereview.PRSummaryReview, error) {
	r.attempts++
	if r.attempts <= r.failCount {
		if r.failErr != nil {
			return codereview.PRSummaryReview{}, r.failErr
		}
		return codereview.PRSummaryReview{}, fmt.Errorf("model unavailable")
	}
	if len(r.aiResponse) == 0 {
		return codereview.PRSummaryReview{}, nil
	}
	return codereview.ParsePRSummaryReview(string(r.aiResponse), brief)
}

func (r *goldenPRSummaryReviewer) ReviewWithOverview(ctx context.Context, brief codereview.ReviewBrief) (string, []codereview.Finding, error) {
	r.attempts++
	if r.attempts <= r.failCount {
		if r.failErr != nil {
			return "", nil, r.failErr
		}
		return "", nil, fmt.Errorf("model unavailable")
	}
	if len(r.aiResponse) == 0 {
		return "", []codereview.Finding{}, nil
	}
	return codereview.ParseAIReviewOutput(string(r.aiResponse), brief)
}
