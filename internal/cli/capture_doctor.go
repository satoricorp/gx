package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/hooks"
	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/uploadauth"
)

const captureBacklogWarnThreshold = 10

type captureDoctorJSON struct {
	HookInstalled   bool   `json:"hookInstalled"`
	UploadAuthed    bool   `json:"uploadAuthed"`
	UploadAPI       string `json:"uploadAPI,omitempty"`
	UploadAuthError string `json:"uploadAuthError,omitempty"`
	PendingExtracts int    `json:"pendingExtracts"`
	PendingSessions int    `json:"pendingSessions"`
	CursorReachable bool   `json:"cursorReachable"`
	CursorPath      string `json:"cursorPath,omitempty"`
	DiskFreeGB      int    `json:"diskFreeGB"`
	DiskWarn        bool   `json:"diskWarn"`
	OK              bool   `json:"ok"`
}

func captureDoctorStatus(ctx context.Context, repoRoot string) captureDoctorJSON {
	status := captureDoctorJSON{}
	if repoRoot == "" {
		if cwd, err := os.Getwd(); err == nil {
			repoRoot = cwd
		}
	}
	status.HookInstalled = hooks.IsInstalled(repoRoot)
	if creds, kind, ok := uploadauth.LoadWithKind(); ok {
		status.UploadAuthed = true
		status.UploadAPI = creds.APIURL
		if kind == "github" {
			status.UploadAuthed, status.UploadAuthError = validateCaptureUploadToken(ctx, creds.Token)
		}
	}
	if counts, err := storage.PendingCaptureCounts(ctx); err == nil {
		status.PendingExtracts = counts.Extracts
		status.PendingSessions = counts.Sessions
	}
	if path, err := cursoringest.DefaultVSCDBPath(); err == nil {
		status.CursorPath = path
		if _, statErr := os.Stat(path); statErr == nil {
			status.CursorReachable = true
		}
	}
	if freeGB, warn := diskFreeGB(); freeGB >= 0 {
		status.DiskFreeGB = freeGB
		status.DiskWarn = warn
	}
	status.OK = status.HookInstalled &&
		status.UploadAuthed &&
		status.CursorReachable &&
		status.PendingExtracts+status.PendingSessions <= captureBacklogWarnThreshold &&
		!status.DiskWarn
	return status
}

func printCaptureDoctor(out fmtWriter, status captureDoctorJSON) {
	fmt.Fprintln(out, section("Capture"))
	if status.HookInstalled {
		fmt.Fprintln(out, labelValue("Pre-push hook", success("ok")))
	} else {
		fmt.Fprintln(out, labelValue("Pre-push hook", danger("warn")+": run `gx init` in this repo"))
	}
	if status.UploadAuthed {
		fmt.Fprintln(out, labelValue("Upload credentials", success("ok")+": "+status.UploadAPI))
	} else {
		hint := "run `gx auth login`"
		if strings.TrimSpace(status.UploadAuthError) != "" {
			hint = status.UploadAuthError
		}
		fmt.Fprintln(out, labelValue("Upload credentials", danger("warn")+": "+hint))
	}
	backlog := status.PendingExtracts + status.PendingSessions
	if backlog == 0 {
		fmt.Fprintln(out, labelValue("Staging backlog", success("ok")+": 0 pending"))
	} else if backlog > captureBacklogWarnThreshold {
		fmt.Fprintln(out, labelValue("Staging backlog", danger("warn")+fmt.Sprintf(": %d pending (run `gx capture sync`)", backlog)))
	} else {
		fmt.Fprintln(out, labelValue("Staging backlog", fmt.Sprintf("%d pending", backlog)))
	}
	if status.CursorReachable {
		fmt.Fprintln(out, labelValue("Cursor vscdb", success("ok")))
	} else {
		msg := "not found"
		if status.CursorPath != "" {
			msg = status.CursorPath
		}
		fmt.Fprintln(out, labelValue("Cursor vscdb", danger("warn")+": "+msg))
	}
	if status.DiskFreeGB >= 0 {
		label := success("ok")
		if status.DiskWarn {
			label = danger("warn")
		}
		fmt.Fprintln(out, labelValue("Disk free", label+fmt.Sprintf(": %d GB", status.DiskFreeGB)))
	}
}

func validateCaptureUploadToken(ctx context.Context, token string) (bool, string) {
	verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	validation, err := cloud.ValidateGitHubAccessToken(verifyCtx, nil, token)
	if err != nil {
		return false, "could not verify GitHub token: " + err.Error()
	}
	if validation.Valid {
		return true, ""
	}
	message := strings.TrimSpace(validation.Error)
	if message == "" {
		message = "GitHub rejected stored token"
	}
	return false, message + "; run `gx auth logout` then `gx auth login`"
}

type fmtWriter interface {
	Write(p []byte) (n int, err error)
}

func diskFreeGB() (int, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return -1, false
	}
	var stat syscallStatfs
	if err := statfs(filepath.Join(home), &stat); err != nil {
		return -1, false
	}
	freeBytes := stat.FreeBytes()
	freeGB := int(freeBytes / (1024 * 1024 * 1024))
	return freeGB, freeGB < 10
}
