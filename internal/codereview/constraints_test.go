package codereview

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setConstraintsTestEnv makes an engine test hermetic: no model, no tools, no
// retrieval, no network.
func setConstraintsTestEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_CONSTRAINTS_RETRIEVAL", "0")
}

// stubConstraintsJudge answers with a canned response (or error) and records
// the request it saw.
type stubConstraintsJudge struct {
	response constraintsJudgeResponse
	err      error
	request  *constraintsJudgeRequest
}

func (s *stubConstraintsJudge) JudgeConstraints(_ context.Context, req constraintsJudgeRequest) (constraintsJudgeResponse, error) {
	s.request = &req
	if s.err != nil {
		return constraintsJudgeResponse{}, s.err
	}
	return s.response, nil
}

func swapConstraintsJudgeFactory(t *testing.T, judge constraintsJudge, unavailable string) {
	t.Helper()
	previous := constraintsJudgeFactory
	constraintsJudgeFactory = func() (constraintsJudge, string, string, string) {
		if judge == nil {
			return nil, "", "", unavailable
		}
		return judge, "stub-model", "stub", ""
	}
	t.Cleanup(func() { constraintsJudgeFactory = previous })
}

func newConstraintsRepo(t *testing.T) string {
	t.Helper()
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommit(t, root)
	return root
}

func commitConstraintsBranchFile(t *testing.T, root, rel, content string) {
	t.Helper()
	runGit(t, root, "checkout", "-b", "feature")
	writeFile(t, root, rel, content)
	gitAdd(t, root, rel)
	runGit(t, root, "commit", "-m", "feature change")
}

func gateByID(t *testing.T, report ConstraintsReport, id GateID) GateResult {
	t.Helper()
	for _, gate := range report.Gates {
		if gate.Gate == id {
			return gate
		}
	}
	t.Fatalf("report has no %s gate: %#v", id, report.Gates)
	return GateResult{}
}

func TestCheckConstraintsNothingToCheck(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if report.Reviewed {
		t.Fatalf("Reviewed = true, want false for a clean base checkout")
	}
	if report.Verdict != VerdictNothingToCheck {
		t.Fatalf("Verdict = %q, want %q", report.Verdict, VerdictNothingToCheck)
	}
	if len(report.Gates) != 0 {
		t.Fatalf("gates = %#v, want none evaluated when nothing was inspected", report.Gates)
	}
}

func TestCheckConstraintsShipOnOrdinaryChange(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() error { return nil }\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if report.Verdict != VerdictShip {
		t.Fatalf("Verdict = %q, want %q\n%s", report.Verdict, VerdictShip, RenderConstraintsText(report))
	}
	if len(report.Gates) != len(AllGateIDs()) {
		t.Fatalf("got %d gates, want %d", len(report.Gates), len(AllGateIDs()))
	}
	if gate := gateByID(t, report, GateSecurity); gate.Status != GatePass {
		t.Fatalf("security gate = %q (%q), want PASS", gate.Status, gate.SkipReason)
	}
	// AI explicitly disabled is an opt-out, not degradation.
	if len(report.DegradedReasons) != 0 {
		t.Fatalf("DegradedReasons = %#v, want none under GX_REVIEW_AI=0", report.DegradedReasons)
	}
	backPressure := gateByID(t, report, GateBackPressure)
	if backPressure.Status != GateSkipped || !strings.Contains(backPressure.SkipReason, "GX_REVIEW_AI=0") {
		t.Fatalf("back-pressure = %q (%q), want SKIPPED with the opt-out named", backPressure.Status, backPressure.SkipReason)
	}
	if gate := gateByID(t, report, GateAccessibility); gate.Status != GateSkipped || gate.SkipReason != "no UI files changed" {
		t.Fatalf("accessibility = %q (%q), want SKIPPED for a non-UI change", gate.Status, gate.SkipReason)
	}
}

func TestCheckConstraintsSecretFailsTheSecurityGate(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/creds.go",
		"package app\n\nconst awsKey = \""+fixtureAWSKey+"\"\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if report.Verdict != VerdictNoShip {
		t.Fatalf("Verdict = %q, want %q\n%s", report.Verdict, VerdictNoShip, RenderConstraintsText(report))
	}
	gate := gateByID(t, report, GateSecurity)
	if gate.Status != GateFail {
		t.Fatalf("security gate = %q, want FAIL", gate.Status)
	}
	if len(gate.Findings) == 0 || gate.Findings[0].File != "internal/app/creds.go" || gate.Findings[0].Line != 3 {
		t.Fatalf("secret finding = %#v, want anchored to internal/app/creds.go:3", gate.Findings)
	}
	// A finding about the change shows the diff itself: the hunk, with the
	// added line prefixed "+", not a plain file excerpt.
	if !strings.Contains(gate.Findings[0].DiffHunk, "+const awsKey") || !strings.Contains(gate.Findings[0].DiffHunk, "@@") {
		t.Fatalf("finding hunk = %q, want the unified-diff window around creds.go:3", gate.Findings[0].DiffHunk)
	}
	if gate.Findings[0].CodeExcerpt != "" {
		t.Fatalf("excerpt = %q, want empty when the hunk is present", gate.Findings[0].CodeExcerpt)
	}
}

func TestCheckConstraintsSkipGatesByFlag(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/creds.go",
		"package app\n\nconst awsKey = \""+fixtureAWSKey+"\"\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{SkipGates: []GateID{GateSecurity}})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	gate := gateByID(t, report, GateSecurity)
	if gate.Status != GateSkipped || gate.SkipReason != "skipped by flag" {
		t.Fatalf("security gate = %q (%q), want SKIPPED by flag", gate.Status, gate.SkipReason)
	}
	if report.Verdict != VerdictShip {
		t.Fatalf("Verdict = %q, want %q when the failing gate was skipped", report.Verdict, VerdictShip)
	}
}

func TestCheckConstraintsFailingTestSuiteFailsCorrectness(t *testing.T) {
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_CONSTRAINTS_RETRIEVAL", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "1")
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() {}\n")
	// A fake `go` that fails every subcommand: `go test` fails correctness and
	// `go vet` fails code health, both from one shared run. A fake govulncheck
	// shadows any real one so the audit leg cannot reach a network.
	tools := t.TempDir()
	writeExecutable(t, filepath.Join(tools, "go"), "#!/bin/sh\necho 'FAIL example.com/repo/internal/app'\nexit 1\n")
	writeExecutable(t, filepath.Join(tools, "govulncheck"), "#!/bin/sh\nexit 0\n")
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if report.Verdict != VerdictNoShip {
		t.Fatalf("Verdict = %q, want %q\n%s", report.Verdict, VerdictNoShip, RenderConstraintsText(report))
	}
	correctness := gateByID(t, report, GateCorrectness)
	if correctness.Status != GateFail {
		t.Fatalf("correctness = %q, want FAIL, checks: %#v", correctness.Status, correctness.Checks)
	}
	if len(correctness.Findings) == 0 || correctness.Findings[0].Strength != "Blocking" {
		t.Fatalf("correctness findings = %#v, want a Blocking finding", correctness.Findings)
	}
	if codeHealth := gateByID(t, report, GateCodeHealth); codeHealth.Status != GateFail {
		t.Fatalf("code-health = %q, want FAIL from the vet partition", codeHealth.Status)
	}
}

func TestCheckConstraintsAccessibilityGateOnUIChange(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "web/Banner.tsx",
		"export const Banner = () => (\n  <div>\n    <img src=\"/banner.png\">\n  </div>\n)\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	gate := gateByID(t, report, GateAccessibility)
	if gate.Status != GateFail {
		t.Fatalf("accessibility = %q, want FAIL for an img without alt\n%s", gate.Status, RenderConstraintsText(report))
	}
	if len(gate.Files) != 1 || gate.Files[0] != "web/Banner.tsx" {
		t.Fatalf("accessibility files = %#v, want the changed UI file", gate.Files)
	}
	if report.Verdict != VerdictNoShip {
		t.Fatalf("Verdict = %q, want %q", report.Verdict, VerdictNoShip)
	}
}

func TestCheckConstraintsDegradedWhenJudgeUnavailable(t *testing.T) {
	setConstraintsTestEnv(t)
	swapConstraintsJudgeFactory(t, nil, "not signed in to gx Cloud: run `gx auth login`")
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() {}\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if report.Verdict != VerdictDegraded {
		t.Fatalf("Verdict = %q, want %q\n%s", report.Verdict, VerdictDegraded, RenderConstraintsText(report))
	}
	if len(report.DegradedReasons) == 0 || !strings.Contains(report.DegradedReasons[0], "gx auth login") {
		t.Fatalf("DegradedReasons = %#v, want the judge's unavailability reason", report.DegradedReasons)
	}
	gate := gateByID(t, report, GateBackPressure)
	if gate.Status != GateSkipped || !strings.Contains(gate.SkipReason, "AI judgment unavailable") {
		t.Fatalf("back-pressure = %q (%q), want SKIPPED as unavailable", gate.Status, gate.SkipReason)
	}
}

func TestCheckConstraintsAppliesAIVerdicts(t *testing.T) {
	setConstraintsTestEnv(t)
	stub := &stubConstraintsJudge{response: constraintsJudgeResponse{Gates: []constraintsGateVerdict{
		{Gate: "code-health", Status: "pass", Justification: "coherent, tested change"},
		{Gate: "back-pressure", Status: "fail", Justification: "regenerated output rides along", Findings: []constraintsAIFinding{{
			Title:          "Regenerated output unrelated to intent",
			Summary:        "scripts/gen.js changed with no relation to the stated intent.",
			Recommendation: "Revert scripts/gen.js or split it into its own change.",
			File:           "scripts/gen.js",
		}}},
		{Gate: "performance", Status: "pass", Justification: "handler change adds no IO or loops"},
	}}}
	swapConstraintsJudgeFactory(t, stub, "")
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/handler.go", "package app\n\nfunc Handle() {}\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{Intent: "fix auth timeout"})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	if stub.request == nil {
		t.Fatalf("the judge was never called")
	}
	if stub.request.Intent != "fix auth timeout" || stub.request.IntentSource != "stated" {
		t.Fatalf("judge intent = %q (%q), want the stated intent", stub.request.Intent, stub.request.IntentSource)
	}
	backPressure := gateByID(t, report, GateBackPressure)
	if backPressure.Status != GateFail {
		t.Fatalf("back-pressure = %q, want FAIL from the model verdict", backPressure.Status)
	}
	if len(backPressure.Findings) != 1 || backPressure.Findings[0].Strength != "Worth exploring" {
		t.Fatalf("back-pressure findings = %#v, want the model finding with the default strength", backPressure.Findings)
	}
	performance := gateByID(t, report, GatePerformance)
	if performance.Status != GatePass || !strings.Contains(performance.Summary, "adds no IO or loops") {
		t.Fatalf("performance = %q (%q), want PASS carrying the justification", performance.Status, performance.Summary)
	}
	if report.Verdict != VerdictNoShip {
		t.Fatalf("Verdict = %q, want %q", report.Verdict, VerdictNoShip)
	}
	if len(report.DegradedReasons) != 0 {
		t.Fatalf("DegradedReasons = %#v, want none on a healthy judged run", report.DegradedReasons)
	}
}

func TestCheckConstraintsDegradesWhenReplyOmitsAGate(t *testing.T) {
	setConstraintsTestEnv(t)
	stub := &stubConstraintsJudge{response: constraintsJudgeResponse{Gates: []constraintsGateVerdict{
		{Gate: "code-health", Status: "pass"},
		// back-pressure omitted; performance not applicable in this change.
	}}}
	swapConstraintsJudgeFactory(t, stub, "")
	root := newConstraintsRepo(t)
	commitConstraintsBranchFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() {}\n")

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	gate := gateByID(t, report, GateBackPressure)
	if gate.Status != GateSkipped {
		t.Fatalf("back-pressure = %q, want SKIPPED when the reply omitted it", gate.Status)
	}
	if report.Verdict != VerdictDegraded {
		t.Fatalf("Verdict = %q, want %q\n%s", report.Verdict, VerdictDegraded, RenderConstraintsText(report))
	}
}

func TestCheckConstraintsLargeChangeWarnsWithoutBlocking(t *testing.T) {
	setConstraintsTestEnv(t)
	root := newConstraintsRepo(t)
	big := "package app\n\n" + strings.Repeat("// filler line for the size threshold\n", 1600)
	commitConstraintsBranchFile(t, root, "internal/app/big.go", big)

	report, err := CheckConstraints(context.Background(), root, ConstraintsOptions{})
	if err != nil {
		t.Fatalf("CheckConstraints() error = %v", err)
	}
	gate := gateByID(t, report, GateCodeHealth)
	if gate.Status != GatePass {
		t.Fatalf("code-health = %q, want PASS: size warns, it does not block", gate.Status)
	}
	if len(gate.Findings) == 0 || gate.Findings[0].ID != "constraints.change-size" {
		t.Fatalf("findings = %#v, want the change-size warning", gate.Findings)
	}
	if report.Verdict != VerdictShip {
		t.Fatalf("Verdict = %q, want %q: a large change still ships on its merits\n%s", report.Verdict, VerdictShip, RenderConstraintsText(report))
	}
}

func TestParseGateIDs(t *testing.T) {
	ids, err := ParseGateIDs(" security, Back-Pressure ")
	if err != nil {
		t.Fatalf("ParseGateIDs() error = %v", err)
	}
	if len(ids) != 2 || ids[0] != GateSecurity || ids[1] != GateBackPressure {
		t.Fatalf("ParseGateIDs() = %#v", ids)
	}
	if _, err := ParseGateIDs("correctness,typo"); err == nil {
		t.Fatalf("ParseGateIDs() accepted an unknown gate")
	}
	if ids, err := ParseGateIDs(""); err != nil || ids != nil {
		t.Fatalf("ParseGateIDs(\"\") = %#v, %v", ids, err)
	}
}
