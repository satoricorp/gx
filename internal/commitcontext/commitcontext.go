package commitcontext

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	AgentDeclaredTool     = "gx_commit"
	AgentDeclaredProvider = "agent-declared"
	AgentDeclaredModelID  = "self-report"
)

// SelfReport is agent-declared provenance attached to a gx commit.
type SelfReport struct {
	TaskSummary string   `json:"task_summary,omitempty"`
	CommandsRun []string `json:"commands_run,omitempty"`
	TestsRun    []string `json:"tests_run,omitempty"`
}

func (s SelfReport) Empty() bool {
	return strings.TrimSpace(s.TaskSummary) == "" &&
		len(s.CommandsRun) == 0 &&
		len(s.TestsRun) == 0
}

func (s SelfReport) JSON() ([]byte, error) {
	return json.Marshal(s)
}

// LoadFile reads a self-report JSON file from path.
func LoadFile(path string) (SelfReport, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return SelfReport{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SelfReport{}, fmt.Errorf("commit context file not found: %s", path)
		}
		return SelfReport{}, fmt.Errorf("read commit context file %s: %w", path, err)
	}
	var report SelfReport
	if err := json.Unmarshal(data, &report); err != nil {
		return SelfReport{}, fmt.Errorf("parse commit context file %s: %w", path, err)
	}
	return report, nil
}
