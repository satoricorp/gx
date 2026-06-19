package extract

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/telemetry"
)

// Client uploads capture payloads to the GX server.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

// NewClient builds an upload client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

// UploadRun posts extract + sessions for one capture run.
func (c *Client) UploadRun(
	ctx context.Context,
	repoRoot, refRange, headCommit, gxVersion string,
	hunkLinks []matcher.HunkLink,
	fileStats interface{},
	sessions []SessionPayload,
	telemetryClient telemetry.Client,
) error {
	req := ExtractRequest{
		RepoRoot:   repoRoot,
		RefRange:   refRange,
		HeadCommit: headCommit,
		GxVersion:  gxVersion,
		HunkLinks:  hunkLinks,
		FileStats:  fileStats,
	}
	if err := c.postJSON(ctx, "/v1/extracts", req); err != nil {
		return fmt.Errorf("upload extract: %w", err)
	}

	for _, session := range sessions {
		content := FormatSessionContent(session)
		body := SessionRequest{
			SessionId: session.SessionID,
			Tool:      session.Tool,
			Model:     SessionModel(session),
			Content:   content,
			CapturedAtMs: time.Now().UnixMilli(),
		}
		if err := c.postJSON(ctx, "/v1/sessions", body); err != nil {
			return fmt.Errorf("upload session %s: %w", session.SessionID, err)
		}
		if telemetryClient != nil {
			telemetryClient.EmitSessionUploaded(ctx, telemetry.SessionUploadedProps{
				SessionID: session.SessionID,
				Tool:      session.Tool,
				Bytes:     len(content),
				Repo:      repoRoot,
				RefRange:  refRange,
			})
		}
	}
	return nil
}

func (c *Client) postJSON(ctx context.Context, path string, body interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(slurp)))
	}
	return nil
}
