package codereview

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Reference lookup.
//
// Four of the thirteen wrong findings in the user report we triaged asked the
// same question and answered it by guessing: where is this symbol used?
//
//	"not every useBuildWithCredits caller threads onPaid"  — two callers, both do
//	"creditsBannerCopy callers may use the old string API" — one caller, reads struct fields
//	"VentureFocusCard accepts onPaid but never uses it"    — passes it at page.tsx:473
//	"HaltReason is unthreaded"                             — passed at page.tsx:739
//
// Each is one `git grep` from settled, and the reader settled all four that
// way in seconds. The reviewer could not, because nothing ever put the
// reference list in front of it.
//
// This is deliberately grep and not a call graph. A real call graph needs
// import resolution, type inference and dispatch analysis, per language, and
// breaks on anything dynamic. Whole-word grep answers "where does this name
// appear" exactly, in milliseconds, in every language at once — which is the
// question these findings actually needed.

// Reference is one place a symbol appears.
type Reference struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text,omitempty"`
}

const (
	// maxReferencesPerSymbol bounds one lookup. A symbol with hundreds of hits
	// is a common word, not the specific thing a finding is about, and a long
	// list is worse than none for the reader and for the judge's budget.
	maxReferencesPerSymbol = 8
	// referenceNoiseThreshold is where a name stops identifying anything. Past
	// it the lookup reports nothing rather than an arbitrary eight of many.
	referenceNoiseThreshold = 40
	// maxReferenceLineChars keeps a minified or generated line from filling
	// the report.
	maxReferenceLineChars = 160
)

// FindReferences returns where symbol appears in the repository, whole-word.
//
// git grep rather than a filesystem walk: it reads the index, so it honours
// .gitignore for free and never descends into node_modules — the directory
// that once turned a 25-second review into 150 seconds by being treated as
// part of the change.
func FindReferences(ctx context.Context, repoRoot, symbol string) []Reference {
	symbol = strings.TrimSpace(symbol)
	if repoRoot == "" || len(symbol) < 3 {
		return nil
	}
	// -w whole word, -F literal, -I skip binaries, -n line numbers.
	cmd := gitCommand(ctx, repoRoot, "grep", "-nIwF", "--no-color", "--", symbol)
	out, err := cmd.Output()
	if err != nil {
		// git grep exits 1 when nothing matched, which is an answer, not a
		// failure: either way there is nothing to report.
		return nil
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	if len(lines) > referenceNoiseThreshold {
		return nil
	}
	refs := make([]Reference, 0, len(lines))
	for _, line := range lines {
		ref, ok := parseGrepLine(line)
		if !ok || isVendoredOrGeneratedPath(ref.File) {
			continue
		}
		refs = append(refs, ref)
	}
	// Stable order so a report does not reshuffle between runs.
	sort.SliceStable(refs, func(i, j int) bool {
		if refs[i].File != refs[j].File {
			return refs[i].File < refs[j].File
		}
		return refs[i].Line < refs[j].Line
	})
	if len(refs) > maxReferencesPerSymbol {
		refs = refs[:maxReferencesPerSymbol]
	}
	return refs
}

// parseGrepLine reads "path:line:text". The path may itself contain colons on
// some platforms, so the split is from the right of the two numeric fields
// rather than the left.
func parseGrepLine(line string) (Reference, bool) {
	first := strings.Index(line, ":")
	if first < 0 {
		return Reference{}, false
	}
	rest := line[first+1:]
	second := strings.Index(rest, ":")
	if second < 0 {
		return Reference{}, false
	}
	number, err := strconv.Atoi(rest[:second])
	if err != nil || number <= 0 {
		return Reference{}, false
	}
	text := strings.TrimSpace(rest[second+1:])
	if len(text) > maxReferenceLineChars {
		text = text[:maxReferenceLineChars] + "…"
	}
	return Reference{File: line[:first], Line: number, Text: text}, true
}

// FormatReferences renders a reference list for a finding's evidence.
func FormatReferences(symbol string, refs []Reference) string {
	if len(refs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(refs))
	for _, ref := range refs {
		parts = append(parts, fmt.Sprintf("%s:%d", ref.File, ref.Line))
	}
	noun := "reference"
	if len(refs) != 1 {
		noun = "references"
	}
	return fmt.Sprintf("`%s` — %d %s: %s", symbol, len(refs), noun, strings.Join(parts, ", "))
}

// referencesEvidenceLabel is the label the attached reference list carries.
// Named so a reader can tell a fact gx looked up from a claim it is making.
const referencesEvidenceLabel = "Checked in the repository"

// AttachSymbolReferences gives each finding the reference list for the symbol
// it names, so the claim arrives with its receipts.
//
// This exists because of who reads the output. gx's findings are read by
// another agent with the repository open, which verifies them and tells the
// user what it found. A finding that asserts something about `onPaid` and
// shows nothing invites that agent to go look and disagree in public; the
// same finding carrying "onPaid — 2 references: page.tsx:473,
// use-lifecycle-state.ts:177" invites it to check and agree. The judge gets
// the same list, so it adjudicates from the call sites instead of guessing at
// them.
//
// Lookups are deduplicated across findings, capped, and skipped entirely for
// names common enough to be meaningless.
func AttachSymbolReferences(ctx context.Context, repoRoot string, findings []Finding) []Finding {
	if strings.TrimSpace(repoRoot) == "" || len(findings) == 0 {
		return findings
	}
	cache := map[string][]Reference{}
	for i := range findings {
		symbol, ok := primarySymbol(findings[i])
		if !ok {
			continue
		}
		refs, seen := cache[symbol]
		if !seen {
			refs = FindReferences(ctx, repoRoot, symbol)
			cache[symbol] = refs
		}
		if len(refs) == 0 {
			continue
		}
		if hasEvidenceLabel(findings[i], referencesEvidenceLabel) {
			continue
		}
		findings[i].Evidence = append(findings[i].Evidence, Evidence{
			Label: referencesEvidenceLabel,
			Value: FormatReferences(symbol, refs),
		})
	}
	return findings
}

func hasEvidenceLabel(finding Finding, label string) bool {
	for _, evidence := range finding.Evidence {
		if evidence.Label == label {
			return true
		}
	}
	return false
}

// primarySymbol is the first backticked identifier a finding names, which is
// in practice the thing the finding is about. Title before summary: the title
// names the subject, while a summary's first backtick is as likely to be a
// supporting detail.
func primarySymbol(finding Finding) (string, bool) {
	for _, text := range []string{finding.Title, finding.Summary} {
		for _, candidate := range backtickedIdentifiers(text) {
			if claimKeywords[strings.ToLower(candidate)] {
				continue
			}
			return candidate, true
		}
	}
	return "", false
}
