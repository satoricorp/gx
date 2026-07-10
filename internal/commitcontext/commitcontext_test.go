package commitcontext

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFileParsesSelfReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "context.json")
	if err := os.WriteFile(path, []byte(`{
		"task_summary": "add commit surface",
		"commands_run": ["go test ./..."],
		"tests_run": ["internal/commitcontext"]
	}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if got.TaskSummary != "add commit surface" {
		t.Fatalf("TaskSummary = %q", got.TaskSummary)
	}
	if len(got.CommandsRun) != 1 || got.CommandsRun[0] != "go test ./..." {
		t.Fatalf("CommandsRun = %#v", got.CommandsRun)
	}
	if len(got.TestsRun) != 1 || got.TestsRun[0] != "internal/commitcontext" {
		t.Fatalf("TestsRun = %#v", got.TestsRun)
	}
}

func TestLoadFileMissingReturnsClearError(t *testing.T) {
	_, err := LoadFile(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("LoadFile() error = nil, want missing file error")
	}
	if got := err.Error(); !strings.Contains(got, "commit context file not found") || !strings.Contains(got, "missing.json") {
		t.Fatalf("error = %q, want missing file message", got)
	}
}

func TestSelfReportEmpty(t *testing.T) {
	var empty SelfReport
	if !empty.Empty() {
		t.Fatal("empty SelfReport should be empty")
	}
	withSummary := SelfReport{TaskSummary: "work"}
	if withSummary.Empty() {
		t.Fatal("SelfReport with task_summary should not be empty")
	}
}
