package cli

import (
	"encoding/json"
	"io"

	"github.com/satoricorp/gx/internal/codereview"
)

// reviewJSONSchema names the shape of `gx review --json` output. Bump it when
// a field is removed or renamed; adding fields does not need a bump. Version 1
// was the bare Report; version 2 adds the envelope fields below.
const reviewJSONSchema = "gx.review/2"

// reviewJSONEnvelope is what --json (and the saved report file) actually
// encodes. The Report is embedded, so every field it has today — and every
// field it grows later — serializes at the top level exactly as before; the
// envelope only adds keys next to them. Nothing existing is renamed or moved.
type reviewJSONEnvelope struct {
	Schema string `json:"schema"`
	codereview.Report
	// FixPlan is the report's findings re-cut as an ordered to-do list for an
	// agent: blocking first, one action per finding. See codereview.BuildFixPlan.
	FixPlan []codereview.FixStep `json:"fix_plan,omitempty"`
}

func writeReviewJSON(out io.Writer, report codereview.Report) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(reviewJSONEnvelope{
		Schema:  reviewJSONSchema,
		Report:  report,
		FixPlan: codereview.BuildFixPlan(report.Findings),
	})
}
