package matcher

import (
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/capture"
)

// The change this exists for: an agent that edits with a heredoc, `sed -i` or
// an inline script never emits a file-edit event, so its work reached the diff
// with nothing in attribution claiming it and was scored as hand-written. The
// transcript held the command text the whole time; nothing read it.
func TestCommandTextClaimsHunksNoEditEventCanClaim(t *testing.T) {
	added := []string{
		"func resolveBedrockJudgeModel() string {",
		"	return normalizeBedrockModelID(firstNonEmpty(",
		"		os.Getenv(\"TOTALITY_REVIEW_JUDGE_MODEL\"),",
		"		defaultBedrockJudgeModel,",
		"	))",
		"}",
	}
	hunks := []HunkRef{{
		Index:      0,
		CommitSHA:  "abc123",
		FilePath:   "internal/codereview/judge.go",
		AddedLines: added,
		CommitTime: 1_000_000,
	}}
	events := []capture.SessionEvent{{
		SessionID: "s1",
		Tool:      capture.ToolClaude,
		Kind:      capture.KindCommand,
		TS:        1_000_000,
		Command:   "python3 - <<'PY'\np='internal/codereview/judge.go'\ns=open(p).read()\ns=s.replace('old', '''" + strings.Join(added, "\n") + "''')\nopen(p,'w').write(s)\nPY",
	}}

	result := Match(events, hunks, DefaultConfig())

	if _, ok := result.CommandHunkIndexes[0]; !ok {
		t.Fatalf("CommandHunkIndexes = %#v, want the hunk claimed from the command text", result.CommandHunkIndexes)
	}
	if got := result.AuthorshipHunkCount(); got != 1 {
		t.Fatalf("AuthorshipHunkCount() = %d, want 1: a command-authored hunk is agent work", got)
	}
	links := BuildHunkLinks([]capture.CommitHunk{{
		CommitSHA:  "abc123",
		FilePath:   "internal/codereview/judge.go",
		AddedLines: added,
	}}, events, result)
	if links[0].Authorship != AuthorshipAgent {
		t.Fatalf("Authorship = %q, want %q", links[0].Authorship, AuthorshipAgent)
	}
	if links[0].Tier != TierCommand || links[0].SessionID != "s1" {
		t.Fatalf("link = %#v, want tier %d attributed to s1", links[0], TierCommand)
	}
}

// A command that merely mentions a line is not its author. Without a size floor
// the character-n-gram fallback in tier1Match makes any short hunk claimable by
// any command quoting it, which would turn `grep` into authorship.
func TestCommandTextDoesNotClaimHunksItOnlyMentions(t *testing.T) {
	hunks := []HunkRef{{
		Index:      0,
		FilePath:   "internal/app/app.go",
		AddedLines: []string{"	return nil"},
		CommitTime: 1_000_000,
	}}
	events := []capture.SessionEvent{{
		SessionID: "s1",
		Kind:      capture.KindCommand,
		TS:        1_000_000,
		Command:   `grep -rn "return nil" internal/`,
	}}

	result := Match(events, hunks, DefaultConfig())

	if len(result.CommandHunkIndexes) != 0 {
		t.Fatalf("CommandHunkIndexes = %#v, want a passing mention to claim nothing", result.CommandHunkIndexes)
	}
}

// A real edit event is better evidence than a command, and the command pass runs
// second precisely so it can never take a hunk away from one.
func TestCommandMatchNeverDowngradesAnEditEventMatch(t *testing.T) {
	added := []string{
		"func resolveBedrockJudgeModel() string {",
		"	return normalizeBedrockModelID(firstNonEmpty(",
		"		os.Getenv(\"TOTALITY_REVIEW_JUDGE_MODEL\"),",
		"		defaultBedrockJudgeModel,",
		"	))",
		"}",
	}
	text := strings.Join(added, "\n")
	hunks := []HunkRef{{
		Index:      0,
		FilePath:   "internal/codereview/judge.go",
		AddedLines: added,
		CommitTime: 1_000_000,
	}}
	events := []capture.SessionEvent{
		{SessionID: "s1", Kind: capture.KindEdit, TS: 1_000_000, FilePath: "internal/codereview/judge.go", NewText: text},
		{SessionID: "s2", Kind: capture.KindCommand, TS: 1_000_000, Command: "cat > judge.go <<'EOF'\n" + text + "\nEOF"},
	}

	result := Match(events, hunks, DefaultConfig())

	if _, ok := result.Tier1HunkIndexes[0]; !ok {
		t.Fatalf("Tier1HunkIndexes = %#v, want the edit event to keep the hunk", result.Tier1HunkIndexes)
	}
	if len(result.CommandHunkIndexes) != 0 {
		t.Fatalf("CommandHunkIndexes = %#v, want the command pass to skip an already-claimed hunk", result.CommandHunkIndexes)
	}
	if got := result.AuthorshipHunkCount(); got != 1 {
		t.Fatalf("AuthorshipHunkCount() = %d, want the hunk counted once, not twice", got)
	}
}

// The command pass has no file path to join on, so time is the only thing
// narrowing it. A command from days before the commit is not evidence.
func TestCommandMatchRespectsTheTemporalWindow(t *testing.T) {
	added := []string{
		"func resolveBedrockJudgeModel() string {",
		"	return normalizeBedrockModelID(firstNonEmpty(",
		"		os.Getenv(\"TOTALITY_REVIEW_JUDGE_MODEL\"),",
		"		defaultBedrockJudgeModel,",
		"	))",
		"}",
	}
	hunks := []HunkRef{{Index: 0, FilePath: "a.go", AddedLines: added, CommitTime: 100_000_000}}
	events := []capture.SessionEvent{{
		SessionID: "s1",
		Kind:      capture.KindCommand,
		TS:        1_000, // days earlier
		Command:   "cat > a.go <<'EOF'\n" + strings.Join(added, "\n") + "\nEOF",
	}}

	if result := Match(events, hunks, DefaultConfig()); len(result.CommandHunkIndexes) != 0 {
		t.Fatalf("CommandHunkIndexes = %#v, want a command outside the window to claim nothing", result.CommandHunkIndexes)
	}
}

// The committed bytes are rarely byte-identical to the command that wrote them:
// gofmt runs, the file already held something, a later edit touched a line.
// Measured against a real push, demanding every n-gram claimed one hunk and
// rejected the next at 0.99. The bar is a threshold for that reason.
func TestCommandMatchToleratesTheGapBetweenCommandAndCommittedBytes(t *testing.T) {
	added := []string{
		"func resolveBedrockJudgeModel() string {",
		"	return normalizeBedrockModelID(firstNonEmpty(",
		"		os.Getenv(\"TOTALITY_REVIEW_JUDGE_MODEL\"),",
		"		defaultBedrockJudgeModel,",
		"	))",
		"}",
		"",
		"func judgeDisabledFromEnv() bool {",
		"	return strings.EqualFold(strings.TrimSpace(os.Getenv(\"TOTALITY_REVIEW_JUDGE\")), \"0\")",
		"}",
	}
	hunks := []HunkRef{{Index: 0, FilePath: "internal/codereview/judge.go", AddedLines: added, CommitTime: 1_000_000}}

	// The command carries all but the closing line: the shape of a gofmt or
	// follow-up-edit difference between what was written and what was committed.
	events := []capture.SessionEvent{{
		SessionID: "s1",
		Kind:      capture.KindCommand,
		TS:        1_000_000,
		Command:   "cat > judge.go <<'EOF'\n" + strings.Join(added[:len(added)-1], "\n") + "\nEOF",
	}}

	result := Match(events, hunks, DefaultConfig())
	if _, ok := result.CommandHunkIndexes[0]; !ok {
		t.Fatalf("CommandHunkIndexes = %#v, want a near-complete containment to still claim the hunk", result.CommandHunkIndexes)
	}

	// Far enough short and it is no longer evidence.
	sparse := []capture.SessionEvent{{
		SessionID: "s1",
		Kind:      capture.KindCommand,
		TS:        1_000_000,
		Command:   "cat > judge.go <<'EOF'\n" + strings.Join(added[:2], "\n") + "\nEOF",
	}}
	if r := Match(sparse, hunks, DefaultConfig()); len(r.CommandHunkIndexes) != 0 {
		t.Fatalf("CommandHunkIndexes = %#v, want a command carrying a fraction of the hunk to claim nothing", r.CommandHunkIndexes)
	}
}
