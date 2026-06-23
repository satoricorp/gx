package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	gxservice "github.com/satoricorp/gx/internal/service"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func newServiceCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the gx background capture service",
	}
	cmd.AddCommand(
		newServiceInstallCommand(ctx),
		newServiceUninstallCommand(ctx),
		newServiceStartCommand(ctx),
		newServiceStopCommand(ctx),
		newServiceStatusCommand(ctx),
		newServiceEnvCommand(ctx),
	)
	return cmd
}

func newServiceInstallCommand(ctx context.Context) *cobra.Command {
	var gxPath string
	return &cobra.Command{
		Use:   "install",
		Short: "Install and start the gx LaunchAgent with ambient capture env",
		RunE: func(cmd *cobra.Command, args []string) error {
			if gxPath == "" {
				exe, err := os.Executable()
				if err != nil {
					return err
				}
				gxPath = exe
			}
			if err := gxservice.NewManager().Install(ctx, gxPath); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Installed", gxservice.Label))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Proxy", gxservice.AnthropicBaseURL()))
			fmt.Fprintln(cmd.OutOrStdout(), muted("Open a new terminal window for environment changes to apply."))
			return nil
		},
	}
}

func newServiceUninstallCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Stop gx and remove LaunchAgent/env integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gxservice.NewManager().Uninstall(ctx); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Removed", gxservice.Label))
			return nil
		},
	}
}

func newServiceStartCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start or restart the gx capture service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gxservice.NewManager().Start(ctx); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Started", gxservice.Label))
			return nil
		},
	}
}

func newServiceStopCommand(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the gx capture service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gxservice.NewManager().Stop(ctx); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Stopped", gxservice.Label))
			return nil
		},
	}
}

func newServiceStatusCommand(ctx context.Context) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Print launchd status for the gx capture service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOut {
				status := serviceStatusJSON(ctx)
				return writeJSON(cmd, status)
			}
			text, err := gxservice.NewManager().Status(ctx)
			if text != "" {
				fmt.Fprintln(cmd.OutOrStdout(), text)
			}
			return err
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

func newServiceEnvCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage gx ambient capture environment variables",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "apply",
			Short: "Set gx base URLs with launchctl and managed shell blocks",
			RunE: func(cmd *cobra.Command, args []string) error {
				mgr := gxservice.NewManager()
				if err := mgr.ApplyEnv(ctx); err != nil {
					return err
				}
				if err := gxservice.InstallShellBlocks(); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), envSummary())
				return nil
			},
		},
		&cobra.Command{
			Use:   "clear",
			Short: "Clear gx launchctl env and managed shell blocks",
			RunE: func(cmd *cobra.Command, args []string) error {
				mgr := gxservice.NewManager()
				if err := mgr.ClearEnv(ctx); err != nil {
					return err
				}
				if err := gxservice.RemoveShellBlocks(); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), success("Cleared gx ambient capture env."))
				return nil
			},
		},
		func() *cobra.Command {
			var jsonOut bool
			statusCmd := &cobra.Command{
				Use:   "status",
				Short: "Print expected gx ambient capture environment values",
				RunE: func(cmd *cobra.Command, args []string) error {
					values, err := gxservice.NewManager().EnvStatus(ctx)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(cmd, envStatusJSON(values))
					}
					fmt.Fprintln(cmd.OutOrStdout(), labelValue(gxservice.EnvAnthropic, compareEnv(values[gxservice.EnvAnthropic], gxservice.AnthropicBaseURL())))
					fmt.Fprintln(cmd.OutOrStdout(), labelValue(gxservice.EnvOpenAI, compareEnv(values[gxservice.EnvOpenAI], gxservice.OpenAIBaseURL())))
					return nil
				},
			}
			statusCmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
			return statusCmd
		}(),
	)
	return cmd
}

func newDoctorCommand(ctx context.Context) *cobra.Command {
	var jsonOut bool
	var fix bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check gx setup and workflow state",
		RunE: func(cmd *cobra.Command, args []string) error {
			var repair vcs.RepairResult
			var repairErr error
			if fix {
				repair, repairErr = vcs.NewService().RepairWorkflow(ctx)
				if repairErr != nil {
					return repairErr
				}
			}
			if jsonOut {
				payload := map[string]any{"doctor": doctorStatusJSON(ctx)}
				if fix {
					payload["repair"] = repair
				}
				return writeJSON(cmd, payload)
			}
			capture := captureDoctorStatus(ctx, "")
			printCaptureDoctor(cmd.OutOrStdout(), capture)
			fmt.Fprintln(cmd.OutOrStdout())
			cursor := cursorStatus(ctx)
			switch {
			case !cursor.Found:
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Cursor", danger("warn")+": state.vscdb not found"))
			default:
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Cursor", fmt.Sprintf("%s: %d sessions, %d messages", success("ok"), cursor.Sessions, cursor.Messages)))
			}
			if fix {
				if len(repair.Actions) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), labelValue("Workflow repair", success("ok")+": no changes needed"))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), section("Workflow repair"))
					for _, action := range repair.Actions {
						fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", action)
					}
				}
				for _, warning := range repair.Warnings {
					if strings.TrimSpace(warning) != "" {
						fmt.Fprintln(cmd.OutOrStdout(), labelWarningValue("Warning", strings.TrimSpace(warning)))
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&fix, "fix", false, "repair safe gx workflow state issues")
	return cmd
}

func newRepairCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair gx ambient capture integrations",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "codex",
		Short: "Set Codex openai_base_url to the gx proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := gxservice.RepairCodexConfig(ctx, gxservice.ExecRunner{})
			if err != nil {
				return err
			}
			state := "updated"
			if status.Correct {
				state = "ok"
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Codex config", state))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("Path", status.ConfigPath))
			fmt.Fprintln(cmd.OutOrStdout(), labelValue("openai_base_url", gxservice.OpenAIBaseURL()))
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "workflow",
		Short: "Repair safe gx workflow ref state",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := vcs.NewService().RepairWorkflow(ctx)
			if err != nil {
				return err
			}
			if len(result.Actions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), labelValue("Workflow repair", success("ok")+": no changes needed"))
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), section("Workflow repair"))
				for _, action := range result.Actions {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", action)
				}
			}
			for _, warning := range result.Warnings {
				if strings.TrimSpace(warning) != "" {
					fmt.Fprintln(cmd.OutOrStdout(), labelWarningValue("Warning", strings.TrimSpace(warning)))
				}
			}
			return nil
		},
	})
	return cmd
}

func envSummary() string {
	return strings.Join([]string{
		labelValue(gxservice.EnvAnthropic, gxservice.AnthropicBaseURL()),
		labelValue(gxservice.EnvOpenAI, gxservice.OpenAIBaseURL()),
	}, "\n")
}

func quoteOrEmpty(value string) string {
	if value == "" {
		return "(unset)"
	}
	return value
}

func compareEnv(got, want string) string {
	if got == want {
		return success("ok") + ": " + got
	}
	if got == "" {
		return danger("warn") + ": unset, want " + want
	}
	return danger("warn") + ": " + got + ", want " + want
}

type envVarStatus struct {
	Name     string `json:"name"`
	Current  string `json:"current"`
	Expected string `json:"expected"`
	OK       bool   `json:"ok"`
}

type launchAgentStatusJSON struct {
	Label string `json:"label"`
	Path  string `json:"path"`
	OK    bool   `json:"ok"`
}

type doctorJSON struct {
	Cursor   cursorStatusJSON  `json:"cursor"`
	Capture  captureDoctorJSON `json:"capture"`
	MCP      mcpStatusJSON     `json:"mcp"`
	Ledger   []ledgerRowJSON   `json:"ledger"`
	Diagnose diagnoseJSON      `json:"diagnose"`
	OK       bool              `json:"ok"`
}

type ledgerRowJSON struct {
	Agent    string `json:"agent"`
	Filepath string `json:"filepath"`
	Status   string `json:"status"`
	Calls    string `json:"calls"`
	Tokens   string `json:"tokens"`
	Files    string `json:"files"`
	Last     string `json:"last"`
}

type diagnoseJSON struct {
	Summary   string `json:"summary"`
	LastCheck string `json:"lastCheck"`
}

type cursorStatusJSON struct {
	VSCDBPath string `json:"vscdbPath"`
	Found     bool   `json:"found"`
	Sessions  int    `json:"sessions"`
	Messages  int    `json:"messages"`
	Error     string `json:"error,omitempty"`
}

type mcpStatusJSON struct {
	OK        bool   `json:"ok"`
	Running   bool   `json:"running"`
	Transport string `json:"transport,omitempty"`
	URL       string `json:"url,omitempty"`
	Port      int    `json:"port,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

type serviceJSON struct {
	Label     string `json:"label"`
	Installed bool   `json:"installed"`
	Running   bool   `json:"running"`
	Path      string `json:"path"`
	Output    string `json:"output,omitempty"`
	Error     string `json:"error,omitempty"`
}

func doctorStatusJSON(ctx context.Context) doctorJSON {
	status := doctorJSON{
		Capture: captureDoctorStatus(ctx, ""),
		Cursor:  cursorStatus(ctx),
		MCP:     mcpStatus(ctx),
	}
	status.OK = status.Capture.OK
	status.Ledger = captureLedgerRows(ctx, status)
	status.Diagnose = diagnoseSummary(status)
	return status
}

func captureLedgerRows(ctx context.Context, status doctorJSON) []ledgerRowJSON {
	db, err := storage.Open(ctx)
	if err != nil {
		return fallbackLedgerRows(status)
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return fallbackLedgerRows(status)
	}
	defer store.Close()

	rows := make([]ledgerRowJSON, 0, 4)
	for _, agent := range []string{"codex", "cursor", "claude"} {
		summary, err := store.AgentLedgerSummary(ctx, agent)
		if err != nil {
			continue
		}
		row := ledgerRowJSON{
			Agent:    agent,
			Filepath: truncateLedgerPath(summary.Filepath),
			Status:   agentLedgerStatus(agent, status, summary.Calls),
			Calls:    formatLedgerCount(summary.Calls),
			Tokens:   formatLedgerTokens(summary.Tokens),
			Files:    formatLedgerCount(summary.Files),
			Last:     formatLedgerLast(summary.LastSeenAt),
		}
		if row.Filepath == "" {
			row.Filepath = defaultLedgerPath(agent, status)
		}
		rows = append(rows, row)
	}
	rows = append(rows, syncLedgerRow(status))
	return rows
}

func fallbackLedgerRows(status doctorJSON) []ledgerRowJSON {
	return []ledgerRowJSON{
		{
			Agent:    "codex",
			Filepath: defaultLedgerPath("codex", status),
			Status:   agentLedgerStatus("codex", status, 0),
			Calls:    "0",
			Tokens:   "—",
			Files:    "—",
			Last:     "—",
		},
		{
			Agent:    "cursor",
			Filepath: defaultLedgerPath("cursor", status),
			Status:   agentLedgerStatus("cursor", status, status.Cursor.Sessions),
			Calls:    formatLedgerCount(status.Cursor.Sessions),
			Tokens:   formatLedgerTokens(status.Cursor.Messages),
			Files:    "—",
			Last:     "—",
		},
		{
			Agent:    "claude",
			Filepath: defaultLedgerPath("claude", status),
			Status:   agentLedgerStatus("claude", status, 0),
			Calls:    "0",
			Tokens:   "—",
			Files:    "—",
			Last:     "—",
		},
		syncLedgerRow(status),
	}
}

func syncLedgerRow(status doctorJSON) ledgerRowJSON {
	return ledgerRowJSON{
		Agent:    "sync",
		Filepath: "review upload queue",
		Status:   "waiting",
		Calls:    "0",
		Tokens:   "—",
		Files:    "—",
		Last:     "push",
	}
}

func agentLedgerStatus(agent string, status doctorJSON, calls int) string {
	if calls > 0 {
		return "indexed"
	}
	if agent == "cursor" {
		if !status.Cursor.Found {
			return "waiting"
		}
		if status.Cursor.Sessions > 0 {
			return "indexed"
		}
		return "attached"
	}
	return "waiting"
}

func defaultLedgerPath(agent string, status doctorJSON) string {
	switch agent {
	case "cursor":
		if status.Cursor.VSCDBPath != "" {
			return truncateLedgerPath(filepath.Dir(filepath.Dir(status.Cursor.VSCDBPath)))
		}
	case "claude":
		return "—"
	}
	return "—"
}

func truncateLedgerPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 34 {
		return value
	}
	return value[:15] + "…" + value[len(value)-16:]
}

func formatLedgerCount(value int) string {
	if value <= 0 {
		return "0"
	}
	return fmt.Sprintf("%d", value)
}

func formatLedgerTokens(value int) string {
	if value <= 0 {
		return "—"
	}
	if value >= 1000 {
		return fmt.Sprintf("%dk", value/1000)
	}
	return fmt.Sprintf("%d", value)
}

func formatLedgerLast(lastSeenAt *int64) string {
	if lastSeenAt == nil || *lastSeenAt <= 0 {
		return "—"
	}
	elapsed := time.Since(time.Unix(*lastSeenAt, 0))
	switch {
	case elapsed < time.Minute:
		return "now"
	case elapsed < time.Hour:
		return fmt.Sprintf("%dm", int(elapsed.Minutes()))
	case elapsed < 24*time.Hour:
		return fmt.Sprintf("%dh", int(elapsed.Hours()))
	default:
		return fmt.Sprintf("%dd", int(elapsed.Hours()/24))
	}
}

func diagnoseSummary(status doctorJSON) diagnoseJSON {
	parts := []string{}
	if status.Capture.OK {
		parts = append(parts, "capture ok")
	} else {
		parts = append(parts, "capture warn")
	}
	if status.Cursor.Found {
		parts = append(parts, "cursor ok")
	} else {
		parts = append(parts, "cursor warn")
	}
	return diagnoseJSON{
		Summary:   strings.Join(parts, " · "),
		LastCheck: time.Now().Format("3:04 PM"),
	}
}

func cursorStatus(ctx context.Context) cursorStatusJSON {
	out := cursorStatusJSON{}
	if path, err := cursoringest.DefaultVSCDBPath(); err == nil {
		out.VSCDBPath = path
		if _, statErr := os.Stat(path); statErr == nil {
			out.Found = true
		}
	}
	db, err := storage.Open(ctx)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		out.Error = err.Error()
		return out
	}
	defer store.Close()
	if sessions, err := store.CountCursorSessions(ctx); err != nil {
		out.Error = err.Error()
	} else {
		out.Sessions = sessions
	}
	if messages, err := store.CountCursorMessages(ctx); err != nil && out.Error == "" {
		out.Error = err.Error()
	} else if err == nil {
		out.Messages = messages
	}
	return out
}

func mcpStatus(context.Context) mcpStatusJSON {
	return mcpStatusJSON{
		OK:        true,
		Running:   false,
		Transport: "stdio",
		Message:   "MCP is stdio-only; MCP clients launch the bundled server directly.",
	}
}

func daemonControlURL() (string, error) {
	return gxservice.ControlURL(), nil
}

func daemonHealthy(controlURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, controlURL+"/healthz", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func serviceStatusJSON(ctx context.Context) serviceJSON {
	launchAgent := launchAgentStatus()
	text, err := gxservice.NewManager().Status(ctx)
	status := serviceJSON{
		Label:     gxservice.Label,
		Installed: launchAgent.OK,
		Path:      launchAgent.Path,
		Output:    text,
	}
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.Running = true
	return status
}

func launchAgentStatus() launchAgentStatusJSON {
	status := launchAgentStatusJSON{Label: gxservice.Label}
	path, err := gxservice.LaunchAgentPath()
	if err != nil {
		return status
	}
	status.Path = path
	_, statErr := os.Stat(path)
	status.OK = statErr == nil
	return status
}

func envStatusJSON(values map[string]string) []envVarStatus {
	return []envVarStatus{
		{
			Name:     gxservice.EnvAnthropic,
			Current:  values[gxservice.EnvAnthropic],
			Expected: gxservice.AnthropicBaseURL(),
			OK:       values[gxservice.EnvAnthropic] == gxservice.AnthropicBaseURL(),
		},
		{
			Name:     gxservice.EnvOpenAI,
			Current:  values[gxservice.EnvOpenAI],
			Expected: gxservice.OpenAIBaseURL(),
			OK:       values[gxservice.EnvOpenAI] == gxservice.OpenAIBaseURL(),
		},
	}
}

func envOK(values []envVarStatus) bool {
	for _, value := range values {
		if !value.OK {
			return false
		}
	}
	return true
}

func writeJSON(cmd *cobra.Command, value any) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
