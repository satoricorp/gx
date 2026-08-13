package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/auth"
	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/semantic"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func newDoctorCommand(ctx context.Context) *cobra.Command {
	var jsonOut bool
	var fix bool
	var sendReport bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Fix current gx state",
		RunE: func(cmd *cobra.Command, args []string) error {
			var repair vcs.RepairResult
			var repairErr error
			var staleRepair vcs.StaleStackCleanupResult
			var staleRepairErr error
			var missingBaseRepair vcs.RebaseOntoDefaultResult
			var missingBaseErr error
			if fix {
				staleRepair, staleRepairErr = vcs.NewService().CleanupStaleStacks(ctx)
				if staleRepairErr != nil {
					return staleRepairErr
				}
				missingBaseRepair, missingBaseErr = vcs.NewService().RepairMissingStackBaseRefs(ctx)
				repair, repairErr = vcs.NewService().RepairWorkflow(ctx)
			}
			stale, staleErr := doctorStaleStacks(ctx, fix, staleRepair)
			missingBase, missingDetectErr := doctorMissingStackBaseRefs(ctx, fix, missingBaseRepair)
			if jsonOut {
				payload := map[string]any{"doctor": doctorStatusJSON(ctx)}
				if outboxStatus, err := publication.QueuedUploadStatus(); err == nil {
					payload["publish_uploads"] = outboxStatus
				}
				if missingBaseErr == nil && (len(missingBaseRepair.Fixed) > 0 || len(missingBaseRepair.Actions) > 0) {
					payload["missing_base_refs"] = missingBaseRepair
				}
				if staleErr == nil {
					payload["stale_stacks"] = stale
				}
				if missingDetectErr == nil {
					payload["missing_base_refs"] = missingBase
				}
				if fix {
					payload["repair"] = repair
					if missingBaseErr == nil {
						payload["missing_base_repair"] = missingBaseRepair
					}
				}
				if !sendReport {
					return writeJSON(cmd, payload)
				}
				// Send the diagnosis that was just produced, then fold the
				// outcome into the same document so stdout stays one JSON value.
				result, reportErr := sendSupportReport(ctx, supportAttachment("gx-doctor.json", marshalDoctorDiagnosis(payload)))
				if reportErr != nil {
					payload["report_error"] = reportErr.Error()
				} else {
					payload["report"] = result
				}
				if err := writeJSON(cmd, payload); err != nil {
					return err
				}
				return reportErr
			}
			// With --report the printed diagnosis is also the thing we send, so
			// tee it into a buffer while it renders.
			out := cmd.OutOrStdout()
			var diagnosis bytes.Buffer
			if sendReport {
				out = io.MultiWriter(out, &diagnosis)
			}
			capture := captureDoctorStatus(ctx, "")
			printCaptureDoctor(out, capture)
			reportAndDrainPublishOutbox(ctx, out)
			reportCodeIndexFreshness(ctx, out, repoRootForDoctor(ctx))
			printDoctorMissingBaseRefRepair(out, missingBaseRepair, missingBaseErr)
			printDoctorStaleStacks(out, fix, stale, staleErr)
			printDoctorMissingStackBaseRefs(out, fix, missingBase, missingDetectErr, missingBaseRepair, missingBaseErr)
			if fix {
				if len(repair.Actions) == 0 {
					fmt.Fprintln(out, labelValue("Workflow repair", success("ok")+": no changes needed"))
				} else {
					fmt.Fprintln(out, section("Workflow repair"))
					for _, action := range repair.Actions {
						fmt.Fprintf(out, "  %s\n", action)
					}
				}
				for _, warning := range repair.Warnings {
					if strings.TrimSpace(warning) != "" {
						fmt.Fprintln(out, labelWarningValue("Warning", strings.TrimSpace(warning)))
					}
				}
				if repairErr != nil {
					fmt.Fprintln(out, labelWarningValue("Workflow repair", "repair failed: "+repairErr.Error()))
				}
			}
			if sendReport {
				result, reportErr := sendSupportReport(ctx, supportAttachment("gx-doctor.txt", diagnosis.String()))
				if reportErr != nil {
					return reportErr
				}
				printSupportReportResult(cmd.OutOrStdout(), result)
			}
			if staleErr != nil {
				return staleErr
			}
			if missingDetectErr != nil {
				return missingDetectErr
			}
			if missingBaseErr != nil {
				return missingBaseErr
			}
			if missingBaseErr != nil {
				return missingBaseErr
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	cmd.Flags().BoolVar(&fix, "fix", false, "repair safe gx workflow state issues")
	cmd.Flags().BoolVar(&sendReport, "report", false, "send this diagnosis and recent gx logs to support")
	return cmd
}

// marshalDoctorDiagnosis renders the JSON-mode payload for the support
// attachment, so `--json --report` sends the same diagnosis it printed.
func marshalDoctorDiagnosis(payload map[string]any) string {
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return ""
	}
	return string(raw)
}

// reportAndDrainPublishOutbox surfaces the publish-outbox backlog and, when
// anything is pending or failed, retries the uploads right here. With `gx sync`
// retired, doctor is the manual retry path; the pre-push hook is the automatic
// one.
func reportAndDrainPublishOutbox(ctx context.Context, w io.Writer) {
	status, err := publication.QueuedUploadStatus()
	if err != nil {
		fmt.Fprintln(w, labelWarningValue("Publish outbox", "check skipped: "+err.Error()))
		return
	}
	if status.Pending == 0 && status.Failed == 0 {
		fmt.Fprintln(w, labelValue("Publish outbox", success("ok")+": empty"))
		return
	}
	detail := fmt.Sprintf("%d pending, %d failed", status.Pending, status.Failed)
	if strings.TrimSpace(status.LastError) != "" {
		detail += "; last error: " + strings.TrimSpace(status.LastError)
	}
	fmt.Fprintln(w, labelWarningValue("Publish outbox", detail))
	if err := drainPublishUploadOutbox(ctx, w, false, 20); err != nil {
		fmt.Fprintln(w, labelWarningValue("Publish outbox", "retry failed: "+err.Error()))
	}
}

func printDoctorMissingBaseRefRepair(w io.Writer, result vcs.RebaseOntoDefaultResult, err error) {
	if err != nil {
		fmt.Fprintln(w, labelWarningValue("Missing parent base", "repair failed: "+err.Error()))
		return
	}
	if len(result.Fixed) == 0 {
		fmt.Fprintln(w, labelValue("Missing parent base", success("ok")+": none"))
		return
	}
	fmt.Fprintln(w, section("Missing parent base repair"))
	for _, action := range result.Actions {
		fmt.Fprintf(w, "  %s\n", action)
	}
}

func doctorStaleStacks(ctx context.Context, fixed bool, repaired vcs.StaleStackCleanupResult) (vcs.StaleStackCleanupResult, error) {
	if fixed {
		return repaired, nil
	}
	return vcs.NewService().ListStaleStacks(ctx)
}

func doctorMissingStackBaseRefs(ctx context.Context, fixed bool, repaired vcs.RebaseOntoDefaultResult) (vcs.MissingStackBaseRefStatus, error) {
	if fixed {
		status, err := vcs.NewService().DetectMissingStackBaseRefs(ctx)
		if err != nil {
			return vcs.MissingStackBaseRefStatus{}, err
		}
		_ = repaired
		return status, nil
	}
	return vcs.NewService().DetectMissingStackBaseRefs(ctx)
}

func printDoctorMissingStackBaseRefs(w io.Writer, fixed bool, status vcs.MissingStackBaseRefStatus, detectErr error, repaired vcs.RebaseOntoDefaultResult, repairErr error) {
	if detectErr != nil {
		fmt.Fprintln(w, labelWarningValue("Missing parent base", "check skipped: "+detectErr.Error()))
		return
	}
	if fixed {
		if len(repaired.Actions) > 0 || len(repaired.Failed) > 0 || repairErr != nil {
			fmt.Fprintln(w, section("Missing parent base repair"))
			for _, action := range repaired.Actions {
				fmt.Fprintf(w, "  %s\n", action)
			}
			for _, failure := range repaired.Failed {
				stackName := firstNonEmptyString(failure.Issue.BookmarkName, failure.Issue.Name)
				fmt.Fprintln(w, labelWarningValue("Missing parent base", fmt.Sprintf("repair failed: rebase stack %s onto %s: %s", stackName, failure.Issue.DefaultBaseRef, failure.Error)))
			}
			if repairErr != nil {
				fmt.Fprintln(w, labelWarningValue("Missing parent base", "repair failed: "+repairErr.Error()))
			}
			return
		}
		fmt.Fprintln(w, labelValue("Missing parent base", success("ok")+": none"))
		return
	}
	if len(status.Issues) == 0 {
		fmt.Fprintln(w, labelValue("Missing parent base", success("ok")+": none"))
		return
	}
	names := make([]string, 0, len(status.Issues))
	for _, issue := range status.Issues {
		names = append(names, firstNonEmptyString(issue.BookmarkName, issue.Name))
	}
	fmt.Fprintln(w, labelWarningValue("Missing parent base", fmt.Sprintf("%d stack(s) with missing parent base refs: %s (run gx doctor)", len(names), strings.Join(names, ", "))))
}

func printDoctorStaleStacks(w io.Writer, fixed bool, result vcs.StaleStackCleanupResult, err error) {
	if err != nil {
		fmt.Fprintln(w, labelWarningValue("Stale stacks", "check skipped: "+err.Error()))
		return
	}
	if len(result.Stale) == 0 {
		fmt.Fprintln(w, labelValue("Stale stacks", success("ok")+": none"))
		return
	}
	if fixed {
		fmt.Fprintln(w, section("Stale stack repair"))
		for _, action := range result.Actions {
			fmt.Fprintf(w, "  %s\n", action)
		}
		return
	}
	names := make([]string, 0, len(result.Stale))
	for _, stack := range result.Stale {
		names = append(names, stack.BookmarkName)
	}
	fmt.Fprintln(w, labelWarningValue("Stale stacks", fmt.Sprintf("%d stack(s) with missing branches: %s (run gx doctor)", len(names), strings.Join(names, ", "))))
}

type doctorJSON struct {
	Cursor   cursorStatusJSON  `json:"cursor"`
	Capture  captureDoctorJSON `json:"capture"`
	MCP      mcpStatusJSON     `json:"mcp"`
	Ledger   []ledgerRowJSON   `json:"ledger"`
	Stats    statsJSON         `json:"stats"`
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

type statsJSON struct {
	ApprovedStacksWaitingForPublish int              `json:"approvedStacksWaitingForPublish"`
	Agents                          []agentStatsJSON `json:"agents"`
	DiskUsedBytes                   int64            `json:"diskUsedBytes"`
}

type agentStatsJSON struct {
	Agent    string `json:"agent"`
	Label    string `json:"label"`
	Health   string `json:"health"`
	Sessions int    `json:"sessions"`
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

func doctorStatusJSON(ctx context.Context) doctorJSON {
	status := doctorJSON{
		Capture: captureDoctorStatus(ctx, ""),
		Cursor:  cursorStatus(ctx),
		MCP:     mcpStatus(ctx),
	}
	status.OK = status.Capture.OK
	status.Ledger = captureLedgerRows(ctx, status)
	status.Stats = doctorStats(ctx, status)
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

// syncLedgerRow reports the real state of the upload queue.
//
// It used to return the literals "waiting" and "0" no matter what, so the one
// screen a user checks to ask whether capture is reaching the cloud claimed a
// healthy idle queue while every row in it was failing on every push.
func syncLedgerRow(status doctorJSON) ledgerRowJSON {
	captureStatus := status.Capture
	pending := captureStatus.PendingExtracts + captureStatus.PendingSessions
	row := ledgerRowJSON{
		Agent:    "sync",
		Filepath: "review upload queue",
		Status:   "waiting",
		Calls:    formatLedgerCount(pending),
		Tokens:   "—",
		Files:    "—",
		Last:     "push",
	}
	switch {
	case captureStatus.FailedUploads > 0:
		row.Status = "failing"
		row.Files = formatLedgerCount(captureStatus.FailedUploads)
	case pending > 0:
		row.Status = "queued"
	}
	return row
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

func doctorStats(ctx context.Context, status doctorJSON) statsJSON {
	stats := statsJSON{
		Agents:        agentStatsRows(status),
		DiskUsedBytes: gxStorageDiskUsedBytes(),
	}

	db, err := storage.Open(ctx)
	if err != nil {
		return stats
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return stats
	}
	defer store.Close()

	repos, err := doctorStatsRepos(ctx, store, status.Capture.RepoRoot)
	if err != nil {
		return stats
	}
	for _, repo := range repos {
		repoStats, err := doctorRepoStackStats(ctx, store, repo.ID)
		if err != nil {
			continue
		}
		stats.ApprovedStacksWaitingForPublish += repoStats.ApprovedStacksWaitingForPublish
	}
	return stats
}

func doctorStatsRepos(ctx context.Context, store *storage.Store, repoRoot string) ([]storage.Repo, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot != "" {
		repoInfo, resolveErr := vcs.NewService().ResolveGxRepoAtPath(ctx, repoRoot)
		var repo *storage.Repo
		var err error
		if resolveErr == nil {
			repo, err = store.FindRepoByIdentity(ctx, repoInfo.GitCommonDir, repoInfo.RootPath)
		} else {
			repo, err = store.FindRepoByRoot(ctx, repoRoot)
		}
		if err != nil {
			return nil, err
		}
		if repo != nil {
			return []storage.Repo{*repo}, nil
		}
	}
	return store.ListRepos(ctx)
}

func doctorRepoStackStats(ctx context.Context, store *storage.Store, repoID int64) (statsJSON, error) {
	var stats statsJSON
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return stats, err
	}
	for _, stack := range stacks {
		if vcs.IsTerminalStackStatus(stack.Status) {
			continue
		}
		changes, err := store.ListChangesByStackID(ctx, stack.ID)
		if err != nil || len(changes) == 0 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(stack.Status), "draft") && isGxOwnedStack(stack) {
			stats.ApprovedStacksWaitingForPublish++
		}
	}
	return stats, nil
}

func isGxOwnedStack(stack storage.Stack) bool {
	return strings.HasPrefix(strings.TrimSpace(stack.BookmarkName), "gx/")
}

func agentStatsRows(status doctorJSON) []agentStatsJSON {
	byAgent := make(map[string]ledgerRowJSON, len(status.Ledger))
	for _, row := range status.Ledger {
		byAgent[row.Agent] = row
	}
	rows := make([]agentStatsJSON, 0, 3)
	for _, agent := range []string{"cursor", "codex", "claude"} {
		ledger := byAgent[agent]
		sessions := parseLedgerSessions(ledger.Calls)
		if agent == "cursor" && sessions == 0 {
			sessions = status.Cursor.Sessions
		}
		rows = append(rows, agentStatsJSON{
			Agent:    agent,
			Label:    agentStatsLabel(agent),
			Health:   agentStatsHealth(agent, ledger, status),
			Sessions: sessions,
		})
	}
	return rows
}

func agentStatsLabel(agent string) string {
	switch agent {
	case "cursor":
		return "Cursor"
	case "codex":
		return "Codex"
	case "claude":
		return "Claude"
	default:
		return agent
	}
}

func agentStatsHealth(agent string, row ledgerRowJSON, status doctorJSON) string {
	if agent == "cursor" && !status.Cursor.Found {
		return "red"
	}
	switch strings.ToLower(strings.TrimSpace(row.Status)) {
	case "indexed", "attached":
		return "green"
	case "waiting":
		return "yellow"
	default:
		return "yellow"
	}
}

func parseLedgerSessions(value string) int {
	value = strings.TrimSpace(value)
	if value == "" || value == "—" {
		return 0
	}
	var out int
	for _, r := range value {
		if r < '0' || r > '9' {
			break
		}
		out = out*10 + int(r-'0')
	}
	return out
}

func gxStorageDiskUsedBytes() int64 {
	root, err := storage.DefaultDir()
	if err != nil {
		return 0
	}
	var total int64
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "gx.db.backup-") {
			return nil
		}
		info, err := entry.Info()
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total
}

func diagnoseSummary(status doctorJSON) diagnoseJSON {
	parts := []string{}
	parts = append(parts, captureDiagnoseSummary(status.Capture))
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

func captureDiagnoseSummary(status captureDoctorJSON) string {
	if status.OK {
		return "capture ok"
	}
	if len(status.Issues) == 0 {
		return "capture warn"
	}
	prefix := "capture warn"
	for _, issue := range status.Issues {
		if issue.Severity == "fail" {
			prefix = "capture fail"
			break
		}
	}
	messages := make([]string, 0, 2)
	for _, issue := range status.Issues {
		message := strings.TrimSpace(issue.Message)
		if message == "" {
			continue
		}
		messages = append(messages, message)
		if len(messages) == 2 {
			break
		}
	}
	if len(messages) == 0 {
		return prefix
	}
	return prefix + ": " + strings.Join(messages, ", ")
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

func writeJSON(cmd *cobra.Command, value any) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

// reportCodeIndexFreshness surfaces how far the local code index has drifted
// from the checkout.
//
// `gx enhance` retrieves from this index, and a stale one fails silently: the
// query succeeds, returns chunks for code that has since changed, and the
// review reads as fully informed. Measured on this repository, an index 30 days
// behind HEAD scored 0.000 recall on every query targeting code written after
// the indexed commit — vector, BM25 and hybrid alike. Retrieval strategy cannot
// compensate for an index that does not contain the code, so doctor reports the
// drift rather than leaving it to be inferred from a disappointing review.
func reportCodeIndexFreshness(ctx context.Context, w io.Writer, repoRoot string) {
	if strings.TrimSpace(repoRoot) == "" {
		return
	}
	orgID := ""
	if creds, ok := auth.LoadUpload(); ok {
		orgID = strings.TrimSpace(creds.OrgID)
	}
	identity := semantic.ResolveRepoIdentity(ctx, repoRoot, orgID, "")

	statePath, err := semantic.RepoIndexStatePath(identity.Namespace)
	if err != nil {
		fmt.Fprintln(w, labelWarningValue("Code index", "check skipped: "+err.Error()))
		return
	}
	state := semantic.LoadRepoIndexState(statePath)
	if state == nil {
		fmt.Fprintln(w, labelWarningValue("Code index",
			"not indexed ("+identity.Namespace+"); run `gx index` so review can retrieve from this repository"))
		return
	}

	head := currentHeadCommit(ctx, repoRoot)
	age := codeIndexAge(state.UpdatedAt)

	// A matching commit is the only state that proves the index describes the
	// checkout. Age alone is not a fault: an untouched repository stays correct
	// however long it sits.
	if head != "" && state.CommitID == head {
		fmt.Fprintln(w, labelValue("Code index", success("ok")+": current at "+shortCommit(head)+age))
		return
	}
	detail := "stale"
	if state.CommitID != "" && head != "" {
		detail += ": indexed at " + shortCommit(state.CommitID) + ", HEAD is " + shortCommit(head)
	}
	detail += age + "; run `gx index` to refresh"
	fmt.Fprintln(w, labelWarningValue("Code index", detail))
}

// codeIndexAge renders " (N days old)" and nothing at all when the timestamp is
// missing or younger than a day, so a fresh index stays quiet.
func codeIndexAge(updatedAtMS int64) string {
	if updatedAtMS <= 0 {
		return ""
	}
	days := int(time.Since(time.UnixMilli(updatedAtMS)).Hours() / 24)
	if days < 1 {
		return ""
	}
	if days == 1 {
		return " (1 day old)"
	}
	return fmt.Sprintf(" (%d days old)", days)
}

func shortCommit(id string) string {
	id = strings.TrimSpace(id)
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// repoRootForDoctor resolves the repository doctor is inspecting. Doctor runs
// outside a repository too, where an empty root simply skips the repo-scoped
// sections rather than failing the whole diagnosis.
func repoRootForDoctor(ctx context.Context) string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	repo, err := vcs.NewService().ResolveGxRepoAtPath(ctx, cwd)
	if err != nil {
		return ""
	}
	return repo.RootPath
}
