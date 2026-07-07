package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/vcs"
)

func TestAttachMissingStackBaseSummaryJSONFields(t *testing.T) {
	summary := authoring.StackSummary{Repo: vcs.RepoInfo{RootPath: "/repo"}}
	status := vcs.MissingStackBaseRefStatus{
		Issues: []vcs.MissingStackBaseRef{{
			BookmarkName:   "feature/structural",
			MissingBaseRef: "feature/authoring",
			DefaultBaseRef: "main",
		}},
		NeedsRebaseOntoDefault: true,
		RepairCommand:          vcs.MissingStackBaseRepairCommand,
	}
	attachMissingStackBaseSummary(&summary, status)
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	body := string(payload)
	for _, want := range []string{
		`"needs_rebase_onto_default":true`,
		`"repair_command":"gx doctor"`,
		`"missing_base_refs"`,
		`"feature/authoring"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("json = %s, want substring %q", body, want)
		}
	}
}

func TestPrintMissingStackBaseRefNotice(t *testing.T) {
	var buf bytes.Buffer
	status := vcs.MissingStackBaseRefStatus{
		Issues: []vcs.MissingStackBaseRef{{
			BookmarkName:   "feature/structural",
			MissingBaseRef: "feature/authoring",
			DefaultBaseRef: "main",
		}},
		RepairCommand: vcs.MissingStackBaseRepairCommand,
	}
	printMissingStackBaseRefNotice(&buf, status)
	out := buf.String()
	for _, want := range []string{"feature/authoring", "feature/structural", "gx doctor"} {
		if !strings.Contains(out, want) {
			t.Fatalf("notice = %q, want substring %q", out, want)
		}
	}
}
