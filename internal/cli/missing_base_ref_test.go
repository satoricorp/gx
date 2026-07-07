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
		FixAction:              vcs.MissingStackBaseFixAction,
		FixPrompt:              "Rebase stack feature/structural onto main because parent base ref \"feature/authoring\" is missing (parent branch was likely merged)? [y/N]",
	}
	attachMissingStackBaseSummary(&summary, status)
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	body := string(payload)
	for _, want := range []string{
		`"needs_rebase_onto_default":true`,
		`"fix_action":"rebase_onto_default"`,
		`"fix_prompt"`,
		`"missing_base_refs"`,
		`"feature/authoring"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("json = %s, want substring %q", body, want)
		}
	}
	for _, absent := range []string{`"repair_command"`} {
		if strings.Contains(body, absent) {
			t.Fatalf("json = %s, should not contain %q", body, absent)
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
	}
	printMissingStackBaseRefNotice(&buf, status)
	out := buf.String()
	for _, want := range []string{"feature/authoring", "feature/structural", "missing parent base ref"} {
		if !strings.Contains(out, want) {
			t.Fatalf("notice = %q, want substring %q", out, want)
		}
	}
}
