package codereview

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Claim checking.
//
// A reviewer that says "`playwright` is missing from pyproject.toml" has made
// a claim the repository can settle in one read. Four of the five wrong
// findings in the first user report we triaged were this shape — a named
// symbol asserted absent, present all along — and the reader found every one
// with a single grep. Several were even phrased as "grep to verify X", the
// reviewer naming the check it did not run.
//
// That is the expensive kind of wrong. gx's output is read by another agent
// with the repository open, so an unchecked claim is not a private mistake: it
// gets verified, refuted, and narrated to the user. Precision costs more than
// recall here, because a missed finding is invisible and a false one gets a
// paragraph written about it.
//
// The rule this file enforces: a finding may not assert that a symbol is
// absent when the symbol is right there.
//
// It is deliberately narrow. It fires only on an unambiguous negation
// governing a backticked identifier, and it drops a finding only when it can
// point at the line that disproves it. Everything it cannot parse it leaves
// alone. A false drop hides a real defect, which is worse than the noise this
// removes, so the design is "prove it wrong or keep it" rather than "look
// suspicious and go".

// absenceClaim is a finding's assertion that Symbol is not present in File.
type absenceClaim struct {
	Symbol string
	File   string
}

// maxClaimFileBytes bounds one file read. A symbol in a file larger than this
// is not worth the read; the claim stays unchecked and the finding survives.
const maxClaimFileBytes = 2 << 20

// claimSymbolFirst matches "`X` is missing", "`X` is never used".
//
// The negation must follow the symbol immediately: at most one linking verb
// between them. "`onPaid` is never passed" is a claim about onPaid;
// "`VentureFocusCard` accepts `onPaid` but never uses it" is not a claim this
// pattern should read, and the intervening words keep it from matching.
var claimSymbolFirst = regexp.MustCompile(
	"`([A-Za-z_][A-Za-z0-9_]{2,})`\\s+(?:is\\s+|are\\s+|was\\s+|were\\s+)?" +
		`(?:missing|absent|undefined|not\s+(?:present|defined|imported|declared|exported)|` +
		`never\s+(?:used|called|referenced|imported|defined|declared|exported))\b`)

// claimCueFirst matches "missing `X`", "without `X`,".
//
// The symbol must end the phrase. Without that anchor, "without a `None`
// check" reads as a claim that `None` is absent — it is a claim about the
// check, and `None` appears in every Python file that mentions it. Requiring
// punctuation or end-of-string after the symbol keeps a trailing noun from
// being silently dropped.
var claimCueFirst = regexp.MustCompile(
	`(?:missing|lacks|without|does\s+not\s+(?:have|include|import|contain|define|declare))\s+` +
		"(?:a\\s+|an\\s+|the\\s+)?`([A-Za-z_][A-Za-z0-9_]{2,})`\\s*(?:[,.;:)]|$)")

// claimKeywords are tokens that name a language construct rather than a
// symbol this repository defines. They appear in nearly every file of their
// language, so "X is missing" about one of them is never settled by finding it
// somewhere in the file.
var claimKeywords = map[string]bool{
	"none": true, "null": true, "nil": true, "true": true, "false": true,
	"try": true, "catch": true, "else": true, "err": true, "error": true,
	"return": true, "await": true, "async": true, "self": true, "this": true,
	"type": true, "class": true, "func": true, "def": true, "let": true,
	"var": true, "const": true, "int": true, "str": true, "bool": true,
}

// detectAbsenceClaim reads a finding for a claim the repository can settle.
//
// The file comes from the finding's own anchor rather than from the prose: a
// finding that names no file has nowhere to check, and a path parsed out of a
// sentence is a guess. Anchors are already required to be real.
func detectAbsenceClaim(finding Finding) (absenceClaim, bool) {
	file := strings.TrimSpace(finding.File)
	if file == "" {
		return absenceClaim{}, false
	}
	text := strings.Join([]string{
		finding.Title,
		finding.Summary,
		finding.Recommendation,
	}, "\n")
	for _, pattern := range []*regexp.Regexp{claimSymbolFirst, claimCueFirst} {
		match := pattern.FindStringSubmatch(text)
		if match == nil {
			continue
		}
		symbol := strings.TrimSpace(match[1])
		if claimKeywords[strings.ToLower(symbol)] {
			continue
		}
		return absenceClaim{Symbol: symbol, File: file}, true
	}
	return absenceClaim{}, false
}

// verifyAbsenceClaim looks for the symbol the finding says is not there.
//
// Returns the 1-based line it was found on. Zero means the claim stands, or
// could not be checked — the caller keeps the finding either way, so the two
// do not need telling apart.
func verifyAbsenceClaim(repoRoot string, claim absenceClaim) int {
	path := filepath.Join(repoRoot, filepath.Clean(claim.File))
	// Never leave the repository, whatever a finding put in its file field.
	if rel, err := filepath.Rel(repoRoot, path); err != nil || strings.HasPrefix(rel, "..") {
		return 0
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Size() > maxClaimFileBytes {
		return 0
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	// Whole-word: a claim about `pid` is not settled by finding it inside
	// `rapid`, and the reviewer meant the identifier.
	word, err := regexp.Compile(`\b` + regexp.QuoteMeta(claim.Symbol) + `\b`)
	if err != nil {
		return 0
	}
	for i, line := range strings.Split(string(body), "\n") {
		if word.MatchString(line) {
			return i + 1
		}
	}
	return 0
}

// DropContradictedFindings removes findings the repository disproves, and
// reports what it removed.
//
// The returned notes are for the run details, not the degraded reasons: a
// review that caught its own bad claim is working, not degrading. They are
// reported at all because a filter that silently eats findings is impossible
// to trust or to debug — if this drops something real, the note is how anyone
// finds out.
func DropContradictedFindings(repoRoot string, findings []Finding) ([]Finding, []string) {
	if strings.TrimSpace(repoRoot) == "" || len(findings) == 0 {
		return findings, nil
	}
	kept := make([]Finding, 0, len(findings))
	var notes []string
	for _, finding := range findings {
		claim, ok := detectAbsenceClaim(finding)
		if !ok {
			kept = append(kept, finding)
			continue
		}
		line := verifyAbsenceClaim(repoRoot, claim)
		if line == 0 {
			kept = append(kept, finding)
			continue
		}
		notes = append(notes, fmt.Sprintf(
			"dropped %q: it reports `%s` absent, and `%s` is at %s:%d",
			firstSentences(finding.Title, 1), claim.Symbol, claim.Symbol, claim.File, line))
	}
	return kept, notes
}
