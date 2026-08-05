package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReportLogsUsesReportedLogsEndpoint(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "gho_report",
		CLISessionToken:   "tlcs_report",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	var got ReportLogRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/reported-logs" {
			t.Fatalf("path = %s, want /v1/reported-logs", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer tlcs_report" {
			t.Fatalf("authorization = %q", auth)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(ReportLogResult{ID: "report-1", URL: "https://lgtm.cx/reports/report-1"})
	}))
	defer server.Close()

	client := &Client{url: server.URL, http: server.Client()}
	result, err := client.ReportLogs(context.Background(), ReportLogRequest{
		TLVersion: "dev",
		Error:     "upload failed",
		Logs:      []ReportLogFile{{Path: "publish-upload.log", Content: "failed"}},
	})
	if err != nil {
		t.Fatalf("ReportLogs() error = %v", err)
	}
	if got.Error != "upload failed" || len(got.Logs) != 1 {
		t.Fatalf("request body = %+v", got)
	}
	if result.ID != "report-1" {
		t.Fatalf("result = %+v", result)
	}
}
