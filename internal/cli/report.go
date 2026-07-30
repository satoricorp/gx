package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/totality/internal/cloud"
	"github.com/satoricorp/totality/internal/publication"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/telemetry"
	"github.com/satoricorp/totality/internal/vcs"
	"github.com/satoricorp/totality/internal/version"
)

const reportLogByteLimit = 32 * 1024

// newReportCommand keeps `tx report` resolvable for one release after the
// behavior moved to `tx doctor --report`. Hidden and ungrouped: existing muscle
// memory and doc links keep working, but the public surface has one spelling.
func newReportCommand(ctx context.Context) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:    "report",
		Short:  "Send logs to support (alias for tx doctor --report)",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := sendSupportReport(ctx, nil)
			if err != nil {
				return err
			}
			if jsonOut {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			printSupportReportResult(cmd.OutOrStdout(), result)
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

// sendSupportReport packages the local support bundle — recent tx logs, cloud
// identity, and the publish-outbox error — and posts it to Totality Cloud. Callers
// pass anything they gathered first as attachments; `tx doctor --report` sends
// its diagnosis that way, so both entry points share this one implementation.
func sendSupportReport(ctx context.Context, attachments []cloud.ReportLogFile) (cloud.ReportLogResult, error) {
	client := cloud.NewClient()
	if client == nil {
		return cloud.ReportLogResult{}, fmt.Errorf("tx cloud is not configured; set TOTALITY_CLOUD_URL or rebuild with cloud endpoints")
	}
	report := buildReportLogRequest(ctx, "")
	report.Logs = append(report.Logs, attachments...)
	result, err := client.ReportLogs(ctx, report)
	if err != nil {
		return cloud.ReportLogResult{}, err
	}
	telemetry.EmitProductEvent(ctx, telemetry.EventCLIReportSent, map[string]any{
		"log_count":        len(report.Logs),
		"has_report_id":    strings.TrimSpace(result.ID) != "",
		"has_user_id":      strings.TrimSpace(report.UserID) != "",
		"has_machine_id":   strings.TrimSpace(report.MachineID) != "",
		"attachment_count": len(attachments),
	})
	return result, nil
}

func printSupportReportResult(out io.Writer, result cloud.ReportLogResult) {
	fmt.Fprintln(out, success("Report sent"))
	if strings.TrimSpace(result.ID) != "" {
		fmt.Fprintln(out, labelValue("Report ID", result.ID))
	}
	if strings.TrimSpace(result.URL) != "" {
		fmt.Fprintln(out, labelValue("Report", result.URL))
	}
}

// supportAttachment turns text a command already produced into a report log
// file: ANSI stripped so support reads plain text, secrets redacted the same
// way real log tails are.
func supportAttachment(path, content string) []cloud.ReportLogFile {
	content = strings.TrimSpace(stripANSI(content))
	if content == "" {
		return nil
	}
	return []cloud.ReportLogFile{{Path: path, Content: redactSensitive(content)}}
}

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(text string) string {
	return ansiEscapePattern.ReplaceAllString(text, "")
}

// Reporting logs to Totality Cloud is a deliberate act: `tx report` and `tx doctor
// --report` send them because the user asked. Nothing uploads on its own.
//
// `tx review` used to auto-report its failures, and it was the only command
// that did. That put a read-only command — one built to run as a CI gate and
// on checkouts the reviewer does not own — in the position of shipping log
// tails and repo identity to Totality Cloud on every failing run, including the
// ordinary "N findings at or above high" that means the gate is working. The
// upload is gone rather than narrowed: a review that fails is the user's
// business, not telemetry.

func buildReportLogRequest(ctx context.Context, overrideError string) cloud.ReportLogRequest {
	report := cloud.ReportLogRequest{
		TLVersion: version.Current(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CloudURL:  cloud.CloudBaseURL(),
		Logs:      recentTotalityLogs(),
	}
	if creds, err := cloud.LoadCloudCredentials(); err == nil && creds != nil {
		report.UserID = strings.TrimSpace(creds.UserID)
		report.Login = strings.TrimSpace(creds.Login)
		report.MachineID = strings.TrimSpace(creds.MachineID)
	}
	if report.MachineID == "" {
		// Read the machine ID, never mint one. Minting writes machine_id.json,
		// which creates $TOTALITY_HOME as a side effect — so a command that only
		// reports a failure would leave Totality state on a machine that has never
		// run tx. A report without a machine ID is worth more than that.
		if machineID, err := cloud.ExistingMachineID(); err == nil {
			report.MachineID = strings.TrimSpace(machineID)
		}
	}
	if status, err := publication.QueuedUploadStatus(); err == nil {
		report.Error = redactSensitive(status.LastError)
	}
	repo, err := vcs.NewService().ResolveGitRepoWithoutStore(ctx)
	if err != nil {
		report.StatusError = redactSensitive(err.Error())
		if strings.TrimSpace(overrideError) != "" {
			report.Error = redactSensitive(overrideError)
		}
		return report
	}
	report.RepoRoot = repo.RootPath
	report.RepoFullName = cloud.RepoFullNameFromRemoteURL(pointerString(repo.RemoteURL))
	if strings.TrimSpace(overrideError) != "" {
		report.Error = redactSensitive(overrideError)
	}
	return report
}

func recentTotalityLogs() []cloud.ReportLogFile {
	dir, err := storage.DefaultDir()
	if err != nil {
		return nil
	}
	logDir := filepath.Join(dir, "logs")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return nil
	}
	sort.Slice(entries, func(i, j int) bool {
		left, leftErr := entries[i].Info()
		right, rightErr := entries[j].Info()
		if leftErr != nil || rightErr != nil {
			return entries[i].Name() < entries[j].Name()
		}
		return left.ModTime().After(right.ModTime())
	})
	logs := make([]cloud.ReportLogFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		path := filepath.Join(logDir, entry.Name())
		content, err := tailFile(path, reportLogByteLimit)
		if err != nil {
			continue
		}
		logs = append(logs, cloud.ReportLogFile{
			Path:    entry.Name(),
			Content: redactSensitive(content),
		})
		if len(logs) >= 5 {
			break
		}
	}
	return logs
}

func tailFile(path string, limit int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	offset := int64(0)
	if info.Size() > limit {
		offset = info.Size() - limit
	}
	if _, err := file.Seek(offset, 0); err != nil {
		return "", err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var reportRedactors = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(authorization\s*[:=]\s*bearer\s+)[^\s,]+`),
	regexp.MustCompile(`(?i)((token|api[_-]?key|password|secret)\s*[:=]\s*)[^\s,]+`),
	regexp.MustCompile(`gh[opsu]_[A-Za-z0-9_]+`),
	regexp.MustCompile(`tlcs_[A-Za-z0-9._-]+`),
}

func redactSensitive(text string) string {
	for _, re := range reportRedactors {
		text = re.ReplaceAllString(text, "${1}[REDACTED]")
	}
	return text
}
