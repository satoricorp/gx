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

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/version"
)

const reportLogByteLimit = 32 * 1024

func newReportCommand(ctx context.Context, engine *authoring.Engine) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Send recent GX logs to support",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := buildReportLogRequest(ctx, engine)
			client := cloud.NewClient()
			if client == nil {
				return fmt.Errorf("gx cloud is not configured; set GX_CLOUD_URL or rebuild with cloud endpoints")
			}
			result, err := client.ReportLogs(ctx, report)
			if err != nil {
				return err
			}
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

func buildReportLogRequest(ctx context.Context, engine *authoring.Engine) cloud.ReportLogRequest {
	report := cloud.ReportLogRequest{
		GXVersion: version.Current(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CloudURL:  cloud.CloudBaseURL(),
		Logs:      recentGXLogs(),
	}
	status, err := currentStatusForEngine(ctx, engine)
	if err != nil {
		report.StatusError = redactSensitive(err.Error())
		return report
	}
	report.RepoRoot = status.Repo.RootPath
	report.RepoFullName = cloud.RepoFullNameFromRemoteURL(pointerString(status.Repo.RemoteURL))
	report.Error = redactSensitive(status.PublishUploads.LastError)
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
