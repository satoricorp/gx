package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// Replay harness: label the findings of an already-completed review, so the
// labeling prompt can be graded without paying for a fresh review each time.
// Skipped unless GX_RULESORT_REPLAY names a saved report.
//
//	GX_RULESORT_REPLAY=review/review-findings-*.json go test ./internal/codereview/ -run Replay -v
func TestRuleSortReplay(t *testing.T) {
	paths := strings.TrimSpace(os.Getenv("GX_RULESORT_REPLAY"))
	if paths == "" {
		t.Skip("set GX_RULESORT_REPLAY to a saved review report to replay labeling")
	}
	sorter := ruleSorterFromEnv()
	if sorter == nil || !sorter.Available() {
		t.Fatal("no rule sorter available; set GX_REVIEW_BEDROCK_DIRECT=1 and AWS credentials")
	}

	totalFindings, totalLabeled := 0, 0
	byRule := map[string]int{}
	for _, path := range strings.Split(paths, ",") {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var report struct {
			RepoRoot string    `json:"repo_root"`
			Findings []Finding `json:"findings"`
		}
		if err := json.Unmarshal(data, &report); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
		// The saved findings already carry the labels the run produced. Clear
		// them, or this measures nothing.
		for i := range report.Findings {
			report.Findings[i].RuleID = ""
		}
		reviewContext := ReviewContext{Brief: ReviewBrief{RepoRoot: report.RepoRoot}}
		outcome := runRuleSort(context.Background(), sorter, reviewContext, report.Findings)
		if outcome.Err != nil {
			t.Errorf("%s: sort error: %v", path, outcome.Err)
		}
		fmt.Printf("\n=== %s — %d labeled / %d findings ===\n", path, outcome.Labeled, len(report.Findings))
		for _, f := range report.Findings {
			label := f.RuleID
			if label == "" {
				label = "(none)"
			} else {
				label = strings.TrimPrefix(label, recommendedPackNamespace)
			}
			byRule[label]++
			title := f.Title
			if len(title) > 96 {
				title = title[:96] + "…"
			}
			fmt.Printf("  %-26s %s\n", label, title)
		}
		totalFindings += len(report.Findings)
		totalLabeled += outcome.Labeled
	}

	fmt.Printf("\n=== TOTAL: %d labeled / %d findings (%d unlabeled) ===\n",
		totalLabeled, totalFindings, totalFindings-totalLabeled)
	rules := make([]string, 0, len(byRule))
	for rule := range byRule {
		rules = append(rules, rule)
	}
	sort.Slice(rules, func(i, j int) bool {
		if byRule[rules[i]] != byRule[rules[j]] {
			return byRule[rules[i]] > byRule[rules[j]]
		}
		return rules[i] < rules[j]
	})
	for _, rule := range rules {
		fmt.Printf("  %2d  %s\n", byRule[rule], rule)
	}
}
