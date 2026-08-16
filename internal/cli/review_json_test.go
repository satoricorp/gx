package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
)

// --json wraps the Report in an envelope: a schema tag and a fix plan next to
// every field the Report already had. Nothing existing moves or renames, so a
// reader written for the bare Report keeps working and a reader that wants the
// to-do list finds it at the top level.
func TestWriteReviewJSONEnvelopeCarriesSchemaFixPlanAndReport(t *testing.T) {
	report := codereview.Report{
		Scope:        "architecture",
		ChangedFiles: []string{"pkg/a.go", "pkg/b.go"},
		Reviewed:     true,
		Findings: []codereview.Finding{
			{ID: "ai.review.1", Title: "advisory", Recommendation: "tidy", Lane: codereview.LaneAdvisory, File: "pkg/a.go", Line: 3},
			{ID: "ai.review.2", RuleID: "review/wrong-logic", Title: "blocking", Recommendation: "fix", Lane: codereview.LaneBlocking, File: "pkg/b.go", Line: 9},
			{ID: "ai.review.3", Title: "no action", Lane: codereview.LaneBlocking},
		},
	}
	var buf bytes.Buffer
	if err := writeReviewJSON(&buf, report); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, buf.String())
	}
	if got, _ := decoded["schema"].(string); got != "gx.review/2" {
		t.Errorf("schema = %q, want gx.review/2", got)
	}
	if reviewJSONSchema != "gx.review/2" {
		t.Errorf("reviewJSONSchema = %q, want gx.review/2", reviewJSONSchema)
	}
	// Report fields still sit at the top level.
	if got, _ := decoded["scope"].(string); got != "architecture" {
		t.Errorf("scope = %q, want architecture (report fields must stay top-level)", got)
	}
	findings, ok := decoded["findings"].([]any)
	if !ok || len(findings) != 3 {
		t.Fatalf("findings missing or wrong length: %v", decoded["findings"])
	}
	if _, present := decoded["changed_files"]; !present {
		t.Errorf("changed_files should still be present")
	}
	// The fix plan is blocking-first and skips the finding with no action.
	plan, ok := decoded["fix_plan"].([]any)
	if !ok {
		t.Fatalf("fix_plan missing: %v", decoded["fix_plan"])
	}
	if len(plan) != 2 {
		t.Fatalf("fix_plan should have 2 steps, got %d: %v", len(plan), plan)
	}
	first, _ := plan[0].(map[string]any)
	second, _ := plan[1].(map[string]any)
	if first["finding_id"] != "ai.review.2" || first["flips_verdict"] != true || first["order"] != float64(1) {
		t.Errorf("first step should be the blocking finding: %v", first)
	}
	if first["rule_id"] != "review/wrong-logic" || first["action"] != "fix" || first["file"] != "pkg/b.go" || first["line"] != float64(9) {
		t.Errorf("first step lost detail: %v", first)
	}
	if second["finding_id"] != "ai.review.1" || second["flips_verdict"] != false || second["order"] != float64(2) {
		t.Errorf("second step should be the advisory finding: %v", second)
	}
	// Round-trips back into a Report: the envelope's extra keys are ignored
	// by the bare type, so older readers of saved report files still work.
	var back codereview.Report
	if err := json.Unmarshal(buf.Bytes(), &back); err != nil {
		t.Fatalf("envelope does not decode as a Report: %v", err)
	}
	if len(back.Findings) != 3 || back.Scope != "architecture" {
		t.Errorf("Report round-trip lost data: %+v", back)
	}
}

// A report with no actionable findings omits fix_plan rather than emitting an
// empty array; schema and findings are always present.
func TestWriteReviewJSONOmitsEmptyFixPlan(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReviewJSON(&buf, codereview.Report{Reviewed: true}); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if _, present := decoded["fix_plan"]; present {
		t.Errorf("fix_plan should be omitted when empty: %s", buf.String())
	}
	if _, present := decoded["schema"]; !present {
		t.Errorf("schema must always be present")
	}
	if _, present := decoded["findings"]; !present {
		t.Errorf("findings must always be present")
	}
}
