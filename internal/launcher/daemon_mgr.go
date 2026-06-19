package launcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	gxservice "github.com/satoricorp/gx/internal/service"
)

type DaemonManager struct {
	httpClient *http.Client
}

type createSessionRequest struct {
	Command   string `json:"command"`
	Cwd       string `json:"cwd"`
	GXVersion string `json:"gx_version"`
}

type createSessionResponse struct {
	SessionID string `json:"session_id"`
	Port      int    `json:"port"`
}

type endSessionRequest struct {
	ExitCode int `json:"exit_code"`
}

func NewDaemonManager() *DaemonManager {
	return &DaemonManager{
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
}

func (m *DaemonManager) Ensure(ctx context.Context) (string, error) {
	controlURL, err := m.controlURL()
	if err == nil && m.healthy(ctx, controlURL) {
		return controlURL, nil
	}
	if err := m.spawn(); err != nil {
		return "", err
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		controlURL, err = m.controlURL()
		if err == nil && m.healthy(ctx, controlURL) {
			return controlURL, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return "", fmt.Errorf("daemon failed healthcheck")
}

func (m *DaemonManager) CreateSession(ctx context.Context, controlURL, command, cwd, gxVersion string) (createSessionResponse, error) {
	body, _ := json.Marshal(createSessionRequest{
		Command:   command,
		Cwd:       cwd,
		GXVersion: gxVersion,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, controlURL+"/sessions", bytes.NewReader(body))
	if err != nil {
		return createSessionResponse{}, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return createSessionResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return createSessionResponse{}, fmt.Errorf("create session: %s", resp.Status)
	}
	var payload createSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return createSessionResponse{}, err
	}
	return payload, nil
}

func (m *DaemonManager) EndSession(ctx context.Context, controlURL, sessionID string, exitCode int) error {
	body, _ := json.Marshal(endSessionRequest{ExitCode: exitCode})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, controlURL+"/sessions/"+sessionID+"/end", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("end session: %s", resp.Status)
	}
	return nil
}

func (m *DaemonManager) healthy(ctx context.Context, controlURL string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, controlURL+"/healthz", nil)
	if err != nil {
		return false
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (m *DaemonManager) spawn() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve gx executable: %w", err)
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer devNull.Close()

	cmd := exec.Command(exe, "__gx-daemon")
	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start daemon: %w", err)
	}
	return cmd.Process.Release()
}

func (m *DaemonManager) controlURL() (string, error) {
	return gxservice.ControlURL(), nil
}
