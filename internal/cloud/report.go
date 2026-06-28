package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/satoricorp/gx/internal/version"
)

type ReportLogFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ReportLogRequest struct {
	GXVersion    string          `json:"gx_version"`
	OS           string          `json:"os"`
	Arch         string          `json:"arch"`
	UserID       string          `json:"user_id,omitempty"`
	Login        string          `json:"login,omitempty"`
	MachineID    string          `json:"machine_id,omitempty"`
	RepoRoot     string          `json:"repo_root,omitempty"`
	RepoFullName string          `json:"repo_full_name,omitempty"`
	CloudURL     string          `json:"cloud_url,omitempty"`
	Error        string          `json:"error,omitempty"`
	StatusError  string          `json:"status_error,omitempty"`
	Logs         []ReportLogFile `json:"logs,omitempty"`
}

type ReportLogResult struct {
	ID  string `json:"id"`
	URL string `json:"url,omitempty"`
}

func (c *Client) ReportLogs(ctx context.Context, report ReportLogRequest) (ReportLogResult, error) {
	if c == nil || c.url == "" {
		return ReportLogResult{}, fmt.Errorf("gx cloud base URL is not configured")
	}
	reportURL := cloudURLWithPath(c.url, "/v1/reported-logs")
	if reportURL == "" {
		return ReportLogResult{}, fmt.Errorf("gx cloud base URL is not configured")
	}
	body, err := json.Marshal(report)
	if err != nil {
		return ReportLogResult{}, fmt.Errorf("marshal gx report: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reportURL, bytes.NewReader(body))
	if err != nil {
		return ReportLogResult{}, fmt.Errorf("create gx report request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gx/"+version.Current())
	token, err := CloudAPIToken()
	if err != nil {
		return ReportLogResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return ReportLogResult{}, fmt.Errorf("send gx report: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := strings.TrimSpace(string(raw))
		if resp.StatusCode == http.StatusUnauthorized {
			if detail != "" {
				return ReportLogResult{}, fmt.Errorf("send gx report: status %s: %s (run `gx auth login`)", resp.Status, detail)
			}
			return ReportLogResult{}, fmt.Errorf("send gx report: status %s (run `gx auth login`)", resp.Status)
		}
		if detail != "" {
			return ReportLogResult{}, fmt.Errorf("send gx report: status %s: %s", resp.Status, detail)
		}
		return ReportLogResult{}, fmt.Errorf("send gx report: status %s", resp.Status)
	}
	var result ReportLogResult
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &result); err != nil {
			return ReportLogResult{}, fmt.Errorf("decode gx report response: %w", err)
		}
	}
	return result, nil
}
