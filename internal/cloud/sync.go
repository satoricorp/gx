package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/lgtm/internal/reviewbundle"
	"github.com/satoricorp/lgtm/internal/version"
)

type Client struct {
	url  string
	http *http.Client
}

type syncUploadResponse struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	ReviewID    string `json:"review_id"`
	ReviewURL   string `json:"review_url"`
	IndexStatus string `json:"index_status"`
}

const defaultUploadTimeout = 2 * time.Minute

func NewClient() *Client {
	url := CloudURL()
	if url == "" {
		return nil
	}
	return &Client{
		url:  url,
		http: &http.Client{Timeout: cloudUploadTimeout()},
	}
}

func cloudUploadTimeout() time.Duration {
	value := strings.TrimSpace(os.Getenv("LGTM_CLOUD_UPLOAD_TIMEOUT"))
	if value == "" {
		return defaultUploadTimeout
	}
	timeout, err := time.ParseDuration(value)
	if err != nil || timeout <= 0 {
		return defaultUploadTimeout
	}
	return timeout
}

func (c *Client) UploadReviewArtifact(ctx context.Context, artifact reviewbundle.Artifact) (reviewbundle.Artifact, error) {
	if c == nil || c.url == "" {
		return artifact, nil
	}
	body, err := json.Marshal(artifact)
	if err != nil {
		return reviewbundle.Artifact{}, fmt.Errorf("marshal lgtm cloud payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cloudURLWithPath(c.url, "/v1/publish"), bytes.NewReader(body))
	if err != nil {
		return reviewbundle.Artifact{}, fmt.Errorf("create lgtm cloud request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "lgtm/"+version.Current())

	token, err := CloudAPIToken()
	if err != nil {
		return reviewbundle.Artifact{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			timeout := c.http.Timeout
			if timeout <= 0 {
				return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload timed out; the server may still persist it, and the queued item retries on your next push (or `lgtm doctor`): %w", err)
			}
			return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload timed out after %s; the server may still persist it, and the queued item retries on your next push (or `lgtm doctor`): %w", timeout, err)
		}
		return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		detail := strings.TrimSpace(string(msg))
		if resp.StatusCode == http.StatusUnauthorized {
			hint := "run `lgtm auth logout` then `lgtm auth login`"
			if detail != "" {
				return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload: status %s: %s (%s)", resp.Status, detail, hint)
			}
			return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload: status %s (%s)", resp.Status, hint)
		}
		if detail != "" {
			return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload: status %s: %s", resp.Status, detail)
		}
		return reviewbundle.Artifact{}, fmt.Errorf("upload lgtm cloud payload: status %s", resp.Status)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return reviewbundle.Artifact{}, fmt.Errorf("decode lgtm cloud response: %w", err)
	}
	result := artifact
	if len(raw) > 0 && json.Valid(raw) {
		var returned reviewbundle.Artifact
		if err := json.Unmarshal(raw, &returned); err == nil && returned.SchemaVersion != 0 {
			result = returned
		}
		var legacy syncUploadResponse
		if err := json.Unmarshal(raw, &legacy); err == nil {
			result.ReviewID = firstNonEmpty(legacy.ReviewID, legacy.ID, result.ReviewID)
			result.ReviewURL = firstNonEmpty(legacy.ReviewURL, legacy.URL, result.ReviewURL)
			result.IndexStatus = firstNonEmpty(legacy.IndexStatus, result.IndexStatus, "pending")
		}
	}
	return result, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
