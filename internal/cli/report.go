package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/version"
)

const reportLogByteLimit = 32 * 1024

func newReportCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Send logs to support",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := buildReportLogRequest(ctx, engine, "")
			client := cloud.NewClient()
			if client == nil {
				return fmt.Errorf("gx cloud is not configured; set GX_CLOUD_URL or rebuild with cloud endpoints")
			}
			result, err := client.ReportLogs(ctx, report)
			if err != nil {
				return err
			}
			telemetry.EmitProductEvent(ctx, telemetry.EventCLIReportSent, map[string]any{
				"log_count":      len(report.Logs),
				"has_report_id":  strings.TrimSpace(result.ID) != "",
				"has_user_id":    strings.TrimSpace(report.UserID) != "",
				"has_machine_id": strings.TrimSpace(report.MachineID) != "",
			})
			if jsonOut {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, success("Report sent"))
			if strings.TrimSpace(result.ID) != "" {
				fmt.Fprintln(out, labelValue("Report ID", result.ID))
			}
			if strings.TrimSpace(result.URL) != "" {
				fmt.Fprintln(out, labelValue("Report", result.URL))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "print machine-readable JSON")
	return cmd
}

func autoReportFailure(ctx context.Context, engine *authoring.Engine, err error, commandName string) {
	if !shouldAutoReportFailure(err) {
		return
	}
	client := cloud.NewClient()
	if client == nil {
		return
	}
	report := buildReportLogRequest(ctx, engine, fmt.Sprintf("%s: %s", strings.TrimSpace(commandName), redactSensitive(err.Error())))
	if _, reportErr := client.ReportLogs(ctx, report); reportErr != nil {
		fmt.Fprintln(os.Stderr, labelWarningValue("Warning", fmt.Sprintf("Could not report %s failure: %v", strings.TrimSpace(commandName), reportErr)))
	}
}

func shouldAutoReportFailure(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) {
		return false
	}
	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
	}
	switch {
	case strings.Contains(message, "does not accept"),
		strings.Contains(message, "accepts at most"),
		strings.Contains(message, "flag provided but not defined"),
		strings.HasPrefix(message, "usage:"),
		strings.Contains(message, "unknown command"),
		strings.Contains(message, "unsupported review"),
		strings.Contains(message, "unsupported format"):
		return false
	default:
		return true
	}
}

func buildReportLogRequest(ctx context.Context, engine *authoring.Engine, overrideError string) cloud.ReportLogRequest {
	report := cloud.ReportLogRequest{
		GXVersion: version.Current(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CloudURL:  cloud.CloudBaseURL(),
		Logs:      recentGXLogs(),
	}
	if creds, err := cloud.LoadCloudCredentials(); err == nil && creds != nil {
		report.UserID = strings.TrimSpace(creds.UserID)
		report.Login = strings.TrimSpace(creds.Login)
		report.MachineID = strings.TrimSpace(creds.MachineID)
	}
	if report.MachineID == "" {
		if machineID, err := cloud.DefaultMachineID(); err == nil {
			report.MachineID = strings.TrimSpace(machineID)
		}
	}
	status, err := currentStatusForEngine(ctx, engine)
	if err != nil {
		report.StatusError = redactSensitive(err.Error())
		if strings.TrimSpace(overrideError) != "" {
			report.Error = redactSensitive(overrideError)
		}
		return report
	}
	report.RepoRoot = status.Repo.RootPath
	report.RepoFullName = cloud.RepoFullNameFromRemoteURL(pointerString(status.Repo.RemoteURL))
	report.Error = redactSensitive(status.PublishUploads.LastError)
	if strings.TrimSpace(overrideError) != "" {
		report.Error = redactSensitive(overrideError)
	}
	return report
}

func recentGXLogs() []cloud.ReportLogFile {
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
	regexp.MustCompile(`gxcs_[A-Za-z0-9._-]+`),
}

func redactSensitive(text string) string {
	for _, re := range reportRedactors {
		text = re.ReplaceAllString(text, "${1}[REDACTED]")
	}
	return text
}
