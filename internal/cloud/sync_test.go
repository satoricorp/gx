package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/buildconfig"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestSyncPushBearerAndReviewURL(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	disableSemanticIndex(t)

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	var gotAuth string
	var gotPayload reviewbundle.Artifact
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/publish" {
			t.Fatalf("path = %s, want /v1/publish", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id":  "event-uuid",
			"url": "http://localhost:3000/reviews/event-uuid",
		})
	}))
	defer server.Close()

	client := &Client{url: server.URL, http: server.Client()}
	result, err := client.UploadReviewArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: t.TempDir(), Backend: "jj"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "abc123"},
		Sessions:      []reviewbundle.SessionPayload{},
	}))
	if err != nil {
		t.Fatalf("UploadReviewBundle() error = %v", err)
	}
	if gotAuth != "Bearer gx_abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("authorization = %q", gotAuth)
	}
	if gotPayload.SchemaVersion != reviewbundle.SchemaVersion || gotPayload.Event != "gx.pr" || gotPayload.Push.HeadCommitID != "abc123" {
		t.Fatalf("payload = %#v, want canonical artifact shape", gotPayload)
	}
	if gotPayload.ReviewID != "" || gotPayload.ReviewURL != "" {
		t.Fatalf("payload review identity = %q/%q, want empty before cloud response", gotPayload.ReviewID, gotPayload.ReviewURL)
	}
	wantURL := "http://localhost:3000/reviews/event-uuid"
	if result.ReviewID != "event-uuid" || result.ReviewURL != wantURL || result.IndexStatus != "pending" {
		t.Fatalf("upload result = %#v, want id/url/pending", result)
	}
}

func TestUploadReviewArtifactAcceptsCanonicalArtifactResponse(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	disableSemanticIndex(t)

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	response := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo", Backend: "jj"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "abc123"},
		Sessions:      []reviewbundle.SessionPayload{},
	})
	response.ReviewID = "review-canonical"
	response.ReviewURL = "http://localhost:3000/reviews/review-canonical"
	response.IndexStatus = "queued"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/publish" {
			t.Fatalf("path = %s, want /v1/publish", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{url: server.URL, http: server.Client()}
	result, err := client.UploadReviewArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo", Backend: "jj"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "abc123"},
		Sessions:      []reviewbundle.SessionPayload{},
	}))
	if err != nil {
		t.Fatalf("UploadReviewArtifact() error = %v", err)
	}
	if result.ReviewID != response.ReviewID || result.ReviewURL != response.ReviewURL || result.IndexStatus != response.IndexStatus {
		t.Fatalf("result = %#v, want canonical response %#v", result, response)
	}
	if result.SchemaVersion != reviewbundle.SchemaVersion || result.Push.HeadCommitID != "abc123" {
		t.Fatalf("result payload = %#v, want artifact fields preserved", result)
	}
}

func TestSyncPushRequiresTokenWhenCloudURLSet(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	disableSemanticIndex(t)

	client := &Client{url: "http://example.invalid", http: http.DefaultClient}
	_, err := client.UploadReviewArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{}))
	if err == nil {
		t.Fatal("expected error without credentials")
	}
	if !strings.Contains(err.Error(), "gx auth login") {
		t.Fatalf("error = %v, want login hint", err)
	}
}

func TestNewClientUsesCloudURL(t *testing.T) {
	os.Unsetenv("GX_CLOUD_URL")
	buildconfig.CloudURL = "https://api.example.com/gx/pr"
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	client := NewClient()
	if client == nil {
		t.Fatal("expected client from baked CloudURL")
	}
	if client.url != "https://api.example.com" {
		t.Fatalf("client.url = %q", client.url)
	}
}

func TestNewClientEnvOverridesBakedDefault(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "http://localhost:3201")
	buildconfig.CloudURL = "https://api.example.com/gx/pr"
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	client := NewClient()
	if client == nil {
		t.Fatal("expected client from env CloudURL")
	}
	if client.url != "http://localhost:3201" {
		t.Fatalf("client.url = %q", client.url)
	}
}

func TestNewClientUsesUploadTimeoutEnv(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "http://localhost:3201")
	t.Setenv("GX_CLOUD_UPLOAD_TIMEOUT", "250ms")

	client := NewClient()
	if client == nil {
		t.Fatal("expected client from env CloudURL")
	}
	if client.http.Timeout != 250*time.Millisecond {
		t.Fatalf("client timeout = %s, want 250ms", client.http.Timeout)
	}
}

func TestNewClientNilWhenUnset(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "")
	buildconfig.CloudURL = ""
	if client := NewClient(); client != nil {
		t.Fatalf("expected nil client, got %+v", client)
	}
}

func disableSemanticIndex(t *testing.T) {
	t.Helper()
	t.Setenv("GX_SEMANTIC_INDEX", "")
	t.Setenv("GX_TPUF_NAMESPACE", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("TURBOPUFFER_API_KEY", "")
}
