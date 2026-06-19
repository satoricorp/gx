package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterPublishUsesV1Endpoint(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_publishabcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/publish" {
			t.Fatalf("path = %s, want /v1/publish", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer gx_publishabcdefghijklmnopqrstuvwxyz" {
			t.Fatalf("authorization = %q", got)
		}
		var body PublishRegistration
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if body.RepoFullName != "satoricorp/gx" || body.HeadCommitID != "abc123" {
			t.Fatalf("request body = %+v", body)
		}
		_ = json.NewEncoder(w).Encode(PublishRegistrationResult{
			ID:           "bookmark-1",
			RepoFullName: "satoricorp/gx",
			BranchName:   "main",
		})
	}))
	defer server.Close()

	client := &Client{url: server.URL, http: server.Client()}
	result, err := client.RegisterPublish(context.Background(), PublishRegistration{
		RepoFullName: "satoricorp/gx",
		BranchName:   "main",
		HeadCommitID: "abc123",
	})
	if err != nil {
		t.Fatalf("RegisterPublish() error = %v", err)
	}
	if result.ID != "bookmark-1" {
		t.Fatalf("result = %+v", result)
	}
}
