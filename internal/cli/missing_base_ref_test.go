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
		FixPrompt:              "Rebase stack feature/structural onto main? [y/N]",
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
		`"missing_base_refs"`,
		`"feature/authoring"`,
		`"fix_prompt"`,
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
		FixPrompt: "Rebase stack feature/structural onto main? [y/N]",
	}
	printMissingStackBaseRefNotice(&buf, status)
	out := buf.String()
	if !strings.Contains(out, "feature/authoring") || !strings.Contains(out, "feature/structural") {
		t.Fatalf("notice = %q, want missing base details", out)
	}
}

func TestPromptFixMissingStackBase(t *testing.T) {
	if promptFixMissingStackBase(strings.NewReader("y\n"), &bytes.Buffer{}, "Rebase? [y/N]") != true {
		t.Fatal("expected yes")
	}
	if promptFixMissingStackBase(strings.NewReader("n\n"), &bytes.Buffer{}, "Rebase? [y/N]") != false {
		t.Fatal("expected no")
	}
}
