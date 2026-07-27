package codereview

import (
	"fmt"
	"strings"
)

// Coverage is the answer to "how much of the subject did this review actually
// read?", and it exists so the review can stop implying the answer is "all of
// it".
//
// A review that read ten of a branch's 181 changed files and then printed "No
// material issues found in this change" was not lying about the findings; it
// was lying about the noun. The same sentence is true and useless when the
// change it names is a tenth of the change. Every number here is counted at the
// point where the brief is handed to the model, so it describes what was read
// rather than what was intended.
type Coverage struct {
	// Subject is what was being reviewed: "the change" or "the repository".
	Subject string `json:"subject,omitempty"`
	// Units is the plural noun for what the counts count.
	Units string `json:"units,omitempty"`
	// Total and Read are files in the subject, and files that reached the
	// model.
	Total int `json:"total"`
	Read  int `json:"read"`
	// Truncated is how many of the files that were read arrived only in part,
	// because one file exceeded the per-snippet budget.
	Truncated int `json:"truncated,omitempty"`
	// Shards is how many concurrent model calls the subject was split across,
	// and ShardsFailed how many of them did not answer. A failed shard is a
	// hole in the coverage, not a clean region.
	Shards       int `json:"shards,omitempty"`
	ShardsFailed int `json:"shards_failed,omitempty"`
	// ContextDropped is how many retrieved/local context snippets did not fit
	// the per-call budget. Unlike the subject itself, context is ranked and
	// genuinely optional, so dropping some is a normal outcome — it is reported
	// rather than fanned out.
	ContextDropped int `json:"context_dropped,omitempty"`
	// OutOfScope is how many tracked files were never candidates: docs,
	// testdata, vendored trees, generated files. Without it, "read all 363
	// files in the repository" is a true sentence about a denominator the
	// reader cannot see, in a repository that tracks 533.
	OutOfScope int `json:"out_of_scope,omitempty"`
	// Planned is false when no coverage was computed (nothing was reviewed, or
	// the model was never called). It keeps a zero value from reading as "read
	// nothing".
	Planned bool `json:"planned,omitempty"`
}

// Complete reports whether everything in the subject reached the model.
func (c Coverage) Complete() bool {
	if !c.Planned {
		return false
	}
	return c.Read >= c.Total && c.Truncated == 0 && c.ShardsFailed == 0
}

// Partial reports whether some but not all of the subject was read.
func (c Coverage) Partial() bool {
	return c.Planned && !c.Complete()
}

// Statement is one sentence stating what was read, in the units a reader can
// check. It is appended to the no-findings message rather than hidden behind
// --verbose, because the no-findings case is exactly where the difference
// between "clean" and "barely looked" decides what someone does next.
func (c Coverage) Statement() string {
	if !c.Planned || c.Total == 0 {
		return ""
	}
	units := strings.TrimSpace(c.Units)
	if units == "" {
		units = "files"
	}
	var b strings.Builder
	if c.Read >= c.Total {
		fmt.Fprintf(&b, "Read all %d %s", c.Total, units)
	} else {
		fmt.Fprintf(&b, "Read %d of %d %s", c.Read, c.Total, units)
	}
	if subject := strings.TrimSpace(c.Subject); subject != "" {
		fmt.Fprintf(&b, " in %s", subject)
	}
	if c.Shards > 1 {
		fmt.Fprintf(&b, " across %d parallel reviews", c.Shards)
	}
	b.WriteString(".")
	if c.OutOfScope > 0 {
		fmt.Fprintf(&b, " %d further tracked file(s) are outside the review's source scope (docs, testdata, vendored, generated).", c.OutOfScope)
	}
	if c.Truncated > 0 {
		fmt.Fprintf(&b, " %d were too large to send whole and were read in part.", c.Truncated)
	}
	if c.ShardsFailed > 0 {
		fmt.Fprintf(&b, " %d of %d parallel reviews failed, so that share of the subject was not reviewed at all.", c.ShardsFailed, c.Shards)
	}
	if c.ContextDropped > 0 {
		fmt.Fprintf(&b, " %d lower-ranked context snippet(s) did not fit the budget.", c.ContextDropped)
	}
	return b.String()
}

// noFindingsMessage is what a review with nothing to report says. The wording
// is conditioned on coverage: "no material issues found in the repository" is a
// claim about the repository, and a review that read a fraction of it has not
// earned that sentence.
func noFindingsMessage(reviewMode string, coverage Coverage) string {
	subject := "this change"
	if reviewMode == ReviewModeRepo {
		subject = "the repository"
	}
	if coverage.Partial() {
		what := "the part of this change that was reviewed"
		if reviewMode == ReviewModeRepo {
			what = "the part of the repository that was reviewed"
		}
		return fmt.Sprintf("No material issues found in %s. This is not a clean bill of health for %s: %s", what, subject, coverage.Statement())
	}
	message := fmt.Sprintf("No material issues found in %s.", subject)
	if statement := coverage.Statement(); statement != "" {
		message += " " + statement
	}
	return message
}
