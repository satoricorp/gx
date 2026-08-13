package codereview

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// The security gate is deterministic end to end: a secrets scan over the added
// lines of the diff, plus the ecosystem's own dependency auditors. The model
// sees the results as facts; it never gets a vote on them.

// constraintsSecretPatterns are token shapes that are a leak on sight. The
// generic assignment heuristic (hardcodedSecretPattern) runs alongside them.
var constraintsSecretPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{name: "AWS access key", pattern: regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{name: "private key block", pattern: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`)},
	{name: "GitHub token", pattern: regexp.MustCompile(`\b(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{20,}`)},
	{name: "GitHub fine-grained token", pattern: regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}`)},
	{name: "OpenAI key", pattern: regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{20,}`)},
	{name: "Slack token", pattern: regexp.MustCompile(`\bxox[abps]-[A-Za-z0-9-]{10,}`)},
}

// constraintsSecretPlaceholderPattern spares the lines that exist to show the
// shape of a secret rather than a secret: docs, examples, redacted samples.
var constraintsSecretPlaceholderPattern = regexp.MustCompile(`(?i)(example|placeholder|dummy|sample|redacted|xxxx|your[_-]?(api[_-]?)?key)`)

const maxConstraintsSecretFindings = 10

// constraintsSecretFindings scans the ADDED lines only: a removed secret is
// the fix, not the leak, and pre-existing ones are the repo's history rather
// than this change's doing.
func constraintsSecretFindings(diffs []DiffSnippet, addedByFile map[string][]constraintsAddedLine) []Finding {
	var findings []Finding
	for _, snippet := range diffs {
		if constraintsSecretScanExemptPath(snippet.File) {
			continue
		}
		language := qualityLanguage(snippet.File)
		for _, added := range addedByFile[snippet.File] {
			text := added.Text
			if strings.TrimSpace(text) == "" || constraintsSecretPlaceholderPattern.MatchString(text) {
				continue
			}
			label := ""
			for _, secret := range constraintsSecretPatterns {
				if secret.pattern.MatchString(text) {
					label = secret.name
					break
				}
			}
			if label == "" && hardcodedSecretPattern(language, strings.TrimSpace(text)) {
				label = "hardcoded credential"
			}
			if label == "" {
				continue
			}
			findings = append(findings, Finding{
				ID:             "constraints.secret-in-diff",
				Scopes:         []string{"security"},
				Title:          "Possible " + label + " in the diff",
				Summary:        fmt.Sprintf("An added line in %s matches the shape of a %s. Anything that lands in a commit lands in every clone.", snippet.File, label),
				Recommendation: "Remove the credential from the change, move it to the environment or a secret store, and rotate it if it was ever real.",
				Strength:       "Blocking",
				Kind:           "defect",
				File:           snippet.File,
				Line:           added.Number,
			})
			if len(findings) >= maxConstraintsSecretFindings {
				return findings
			}
		}
	}
	return findings
}

// constraintsSecretScanExemptPath spares the files whose "secrets" are props:
// tests planting fixture keys, testdata, docs. The same stance the quality
// hints take (collectCodeQualityHints skips test files) — a real credential
// pasted into a test is not impossible, but fixture keys are near-certain, and
// a gate that cries wolf on every planted AKIA gets turned off.
func constraintsSecretScanExemptPath(rel string) bool {
	if isTestFile(rel) {
		return true
	}
	lower := strings.ToLower(filepath.ToSlash(rel))
	return strings.HasPrefix(lower, "testdata/") || strings.Contains(lower, "/testdata/") ||
		strings.Contains(lower, "/fixtures/") || strings.HasSuffix(lower, ".md")
}

// constraintsAuditRunner is one dependency auditor. Same contract as the
// static-tool registry: detect inspects the checkout and the installed
// binaries; a missing auditor is evidence, never a failure.
type constraintsAuditRunner struct {
	name   string
	detect func(repoRoot string) (staticToolCommand, bool)
	// missing explains an auditor that applies to this checkout but cannot
	// run, for the evidence line ("govulncheck not installed").
	missing func(repoRoot string) string
}

var constraintsAuditRunners = []constraintsAuditRunner{
	{
		name: "govulncheck",
		detect: func(repoRoot string) (staticToolCommand, bool) {
			if !repoFileExists(repoRoot, "go.mod") {
				return staticToolCommand{}, false
			}
			bin, ok := lookStaticTool(repoRoot, "govulncheck")
			if !ok {
				return staticToolCommand{}, false
			}
			return staticToolCommand{bin: bin, argv: []string{"govulncheck", "./..."}}, true
		},
		missing: func(repoRoot string) string {
			if !repoFileExists(repoRoot, "go.mod") {
				return ""
			}
			if _, ok := lookStaticTool(repoRoot, "govulncheck"); ok {
				return ""
			}
			return "govulncheck not installed"
		},
	},
	{
		name: "npm audit",
		detect: func(repoRoot string) (staticToolCommand, bool) {
			if !repoFileExists(repoRoot, "package-lock.json") {
				return staticToolCommand{}, false
			}
			bin, ok := lookStaticTool(repoRoot, "npm")
			if !ok {
				return staticToolCommand{}, false
			}
			return staticToolCommand{bin: bin, argv: []string{"npm", "audit", "--omit=dev", "--audit-level=high"}}, true
		},
		missing: func(repoRoot string) string {
			if repoFileExists(repoRoot, "package-lock.json") {
				return ""
			}
			if repoFileExists(repoRoot, "package.json") {
				return "npm audit skipped (no package-lock.json)"
			}
			return ""
		},
	},
	{
		name: "pip-audit",
		detect: func(repoRoot string) (staticToolCommand, bool) {
			if !repoFileExists(repoRoot, "requirements.txt") && !repoFileContains(repoRoot, "pyproject.toml", "[project]") {
				return staticToolCommand{}, false
			}
			bin, ok := lookStaticTool(repoRoot, "pip-audit", pythonLocalBinDirs...)
			if !ok {
				return staticToolCommand{}, false
			}
			return staticToolCommand{bin: bin, argv: []string{"pip-audit"}}, true
		},
		missing: func(repoRoot string) string {
			if !repoFileExists(repoRoot, "requirements.txt") && !repoFileContains(repoRoot, "pyproject.toml", "[project]") {
				return ""
			}
			if _, ok := lookStaticTool(repoRoot, "pip-audit", pythonLocalBinDirs...); ok {
				return ""
			}
			return "pip-audit not installed"
		},
	},
	{
		name: "cargo audit",
		detect: func(repoRoot string) (staticToolCommand, bool) {
			if !repoFileExists(repoRoot, "Cargo.lock") {
				return staticToolCommand{}, false
			}
			// `cargo audit` is the cargo-audit plugin; detecting the plugin
			// binary is what proves the subcommand exists.
			if _, ok := lookStaticTool(repoRoot, "cargo-audit"); !ok {
				return staticToolCommand{}, false
			}
			bin, ok := lookStaticTool(repoRoot, "cargo")
			if !ok {
				return staticToolCommand{}, false
			}
			return staticToolCommand{bin: bin, argv: []string{"cargo", "audit"}}, true
		},
		missing: func(repoRoot string) string {
			if !repoFileExists(repoRoot, "Cargo.lock") {
				return ""
			}
			if _, ok := lookStaticTool(repoRoot, "cargo-audit"); ok {
				return ""
			}
			return "cargo-audit not installed"
		},
	},
}

// collectConstraintsAuditResults runs the applicable auditors concurrently.
// It honors the same kill switch as the static tools: the auditors exec
// checkout-adjacent binaries and reach the network, which is exactly what a
// hermetic run must not do.
func collectConstraintsAuditResults(ctx context.Context, repoRoot string, changed []string) []StaticToolResult {
	if constraintsStaticToolsDisabled() {
		return nil
	}
	type plannedAudit struct {
		name    string
		command staticToolCommand
		// skipReason marks an auditor that applies to this checkout but cannot
		// run; it becomes a skipped result so the gate says what went unchecked.
		skipReason string
	}
	var planned []plannedAudit
	for _, runner := range constraintsAuditRunners {
		if command, ok := runner.detect(repoRoot); ok {
			planned = append(planned, plannedAudit{name: runner.name, command: command})
			continue
		}
		if reason := runner.missing(repoRoot); reason != "" {
			planned = append(planned, plannedAudit{name: runner.name, skipReason: reason})
		}
	}
	if len(planned) == 0 {
		return nil
	}
	out := make([]StaticToolResult, len(planned))
	var wg sync.WaitGroup
	for i, audit := range planned {
		if audit.skipReason != "" {
			out[i] = StaticToolResult{Name: audit.name, Skipped: true, Reason: audit.skipReason}
			continue
		}
		i, audit := i, audit
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := runStaticTool(ctx, repoRoot, 90*time.Second, audit.name, audit.command)
			if !result.Skipped && result.ExitCode != 0 {
				if reason := constraintsAuditEnvironmentFailure(audit.name, result.Output); reason != "" {
					result.Skipped = true
					result.Reason = reason
				}
			}
			out[i] = result
		}()
	}
	wg.Wait()
	return out
}

// constraintsAuditEnvironmentFailure classifies an audit failure the change
// cannot have caused — no network, no registry — so an offline run reports
// SKIPPED rather than a flaky FAIL. Same stance as staticToolEnvironmentFailure.
func constraintsAuditEnvironmentFailure(name, output string) string {
	networkMarkers := []string{
		"dial tcp", "no such host", "i/o timeout", "connection refused",
		"TLS handshake timeout", "ENOTFOUND", "ETIMEDOUT", "ECONNREFUSED",
		"ECONNRESET", "EAI_AGAIN", "network is unreachable",
		"Temporary failure in name resolution",
		"couldn't fetch advisory database",
		"unable to connect", "Could not fetch URL",
	}
	if outputContainsAny(output, networkMarkers...) {
		return "the vulnerability database could not be reached on this checkout (network unavailable)"
	}
	if name == "npm audit" && strings.Contains(output, "ENOLOCK") {
		return "npm audit requires a lockfile on this checkout"
	}
	return ""
}

// constraintsSecurityGate assembles the verdict: secrets in the diff or a
// confirmed vulnerable dependency tree fail; a missing or unreachable auditor
// is said out loud and fails nothing.
func constraintsSecurityGate(repoRoot string, changed []string, diffs []DiffSnippet, audits []StaticToolResult) GateResult {
	gate := GateResult{Gate: GateSecurity, Title: gateTitle(GateSecurity), Status: GatePass}
	addedByFile := map[string][]constraintsAddedLine{}
	for _, snippet := range diffs {
		addedByFile[snippet.File] = constraintsAddedLinesForDiff(snippet.Diff)
	}
	secretFindings := constraintsSecretFindings(diffs, addedByFile)
	fileSet := map[string]struct{}{}
	var summary []string
	if len(secretFindings) > 0 {
		gate.Status = GateFail
		gate.Findings = append(gate.Findings, secretFindings...)
		for _, finding := range secretFindings {
			fileSet[finding.File] = struct{}{}
		}
		summary = append(summary, fmt.Sprintf("%d possible secret(s) in added lines", len(secretFindings)))
	} else {
		summary = append(summary, "no secrets in added lines")
	}

	gate.Checks = audits
	if len(audits) == 0 {
		if constraintsStaticToolsDisabled() {
			summary = append(summary, "dependency audit disabled (GX_REVIEW_STATIC_TOOLS=0)")
		} else {
			summary = append(summary, "dependency audit: no supported auditor detected")
		}
	}
	for _, audit := range audits {
		summary = append(summary, constraintsToolLine(audit))
		if audit.Skipped {
			continue
		}
		if audit.ExitCode != 0 {
			gate.Status = GateFail
			finding := constraintsToolFinding(audit, "security")
			finding.Title = audit.Name + " found vulnerable dependencies"
			finding.Summary = fmt.Sprintf("`%s` exited %d: known vulnerabilities in this dependency tree.", audit.Command, audit.ExitCode)
			finding.Recommendation = fmt.Sprintf("Run `%s` locally, then upgrade or replace the flagged dependencies before shipping.", audit.Command)
			gate.Findings = append(gate.Findings, finding)
		}
	}

	// SAST-lite: the review engine's quality hints over changed source files.
	// Candidates for the model to confirm, never an automatic failure.
	for _, file := range changed {
		if isTestFile(file) || !qualityFileSupported(file) || constraintsSecretScanExemptPath(file) {
			continue
		}
		for _, hint := range qualityHintsForFile(repoRoot, file) {
			if qualityRank(hint.Kind) != 0 || hint.Kind == "hardcoded_secret" {
				continue
			}
			gate.Evidence = append(gate.Evidence, Evidence{
				Label: fmt.Sprintf("%s:%d", hint.File, hint.Line),
				Value: hint.Kind + " candidate: " + hint.Reason,
			})
			fileSet[hint.File] = struct{}{}
			if len(gate.Evidence) >= 12 {
				break
			}
		}
		if len(gate.Evidence) >= 12 {
			break
		}
	}

	for _, file := range changed {
		if isDependencyFile(file) {
			fileSet[file] = struct{}{}
		}
	}
	for file := range fileSet {
		gate.Files = append(gate.Files, file)
	}
	sort.Strings(gate.Files)
	gate.Summary = strings.Join(summary, "; ")
	return gate
}
