package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/hooks"
	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/uploadauth"
)

type captureDoctorJSON struct {
	HookInstalled        bool                 `json:"hookInstalled"`
	HookApplicable       bool                 `json:"hookApplicable"`
	RepoRoot             string               `json:"repoRoot,omitempty"`
	RepoHooks            []repoHookDoctorJSON `json:"repoHooks,omitempty"`
	RepoHooksTotal       int                  `json:"repoHooksTotal"`
	RepoHooksMissing     int                  `json:"repoHooksMissing"`
	RepoHooksUnreachable int                  `json:"repoHooksUnreachable"`
	RepoHooksOK          bool                 `json:"repoHooksOK"`
	UploadAuthed         bool                 `json:"uploadAuthed"`
	UploadAPI            string               `json:"uploadAPI,omitempty"`
	UploadAuthError      string               `json:"uploadAuthError,omitempty"`
	PendingExtracts      int                  `json:"pendingExtracts"`
	PendingSessions      int                  `json:"pendingSessions"`
	CursorReachable      bool                 `json:"cursorReachable"`
	CursorPath           string               `json:"cursorPath,omitempty"`
	DiskFreeGB           int                  `json:"diskFreeGB"`
	DiskWarn             bool                 `json:"diskWarn"`
	Issues               []captureIssueJSON   `json:"issues,omitempty"`
	OK                   bool                 `json:"ok"`
}

type captureIssueJSON struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Action   string `json:"action"`
}

type repoHookDoctorJSON struct {
	RepoRoot      string `json:"repoRoot"`
	HookInstalled bool   `json:"hookInstalled"`
	GitReachable  bool   `json:"gitReachable"`
	Error         string `json:"error,omitempty"`
}

func captureDoctorStatus(ctx context.Context, repoRoot string) captureDoctorJSON {
	status := captureDoctorJSON{}
	status.HookInstalled, status.HookApplicable, status.RepoRoot = capturePrimaryHookStatus(ctx, repoRoot)
	status.RepoHooks = captureRegisteredRepoHooks(ctx)
	status.RepoHooksTotal = len(status.RepoHooks)
	status.RepoHooksOK = true
	for _, repoHook := range status.RepoHooks {
		switch {
		case !repoHook.GitReachable:
			status.RepoHooksUnreachable++
			status.RepoHooksOK = false
		case !repoHook.HookInstalled:
			status.RepoHooksMissing++
			status.RepoHooksOK = false
		}
	}
	if creds, kind, ok := uploadauth.LoadWithKind(); ok {
		status.UploadAuthed = true
		status.UploadAPI = creds.APIURL
		switch kind {
		case "gx-cli":
			status.UploadAuthed, status.UploadAuthError = validateCaptureCloudSession(ctx, creds.Token)
		case "github":
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
	status.OK = captureDoctorOK(status)
	status.Issues = captureDoctorIssues(status)
	return status
}

func captureDoctorOK(status captureDoctorJSON) bool {
	return (!status.HookApplicable || status.HookInstalled) &&
		status.RepoHooksOK &&
		status.UploadAuthed &&
		status.CursorReachable &&
		!status.DiskWarn
}

func printCaptureDoctor(out fmtWriter, status captureDoctorJSON) {
	fmt.Fprintln(out, labelValue("Hooks installed", captureHooksDoctorValue(status)))
	if status.UploadAuthed {
		fmt.Fprintln(out, labelValue("Upload", success("ok")+": "+status.UploadAPI))
	} else {
		hint := "run `gx auth login`"
		if strings.TrimSpace(status.UploadAuthError) != "" {
			hint = status.UploadAuthError
		}
		fmt.Fprintln(out, labelValue("Upload", danger("warn")+": "+hint))
	}
	backlog := status.PendingExtracts + status.PendingSessions
	if backlog == 0 {
		fmt.Fprintln(out, labelValue("Staging backlog", success("ok")+": 0 pending"))
	} else {
		fmt.Fprintln(out, labelValue("Staging backlog", fmt.Sprintf("%d pending", backlog)))
	}
	fmt.Fprintln(out, labelValue("Disk used", formatDiskUsedGB(gxStorageDiskUsedBytes())))
}

func captureHooksDoctorValue(status captureDoctorJSON) string {
	if status.HookApplicable && !status.HookInstalled {
		return danger("warn") + ": run `gx init` in this repo"
	}
	if !status.RepoHooksOK {
		return danger("warn") + fmt.Sprintf(": %d missing, %d unreachable of %d repos", status.RepoHooksMissing, status.RepoHooksUnreachable, status.RepoHooksTotal)
	}
	if status.HookApplicable || status.RepoHooksTotal > 0 {
		return success("ok")
	}
	return "not checked: run `gx doctor` inside a git repo"
}

func formatDiskUsedGB(bytes int64) string {
	if bytes <= 0 {
		return "0GB"
	}
	return fmt.Sprintf("%dGB", bytes/(1024*1024*1024))
}

func captureRegisteredRepoHooks(ctx context.Context) []repoHookDoctorJSON {
	candidates := captureRepoRootCandidates(ctx)
	statuses := make([]repoHookDoctorJSON, 0, len(candidates))
	for _, candidate := range candidates {
		installed, reachable, resolvedRoot := captureHookStatus(ctx, candidate.root)
		if !reachable && !candidate.reportUnreachable {
			continue
		}
		status := repoHookDoctorJSON{
			RepoRoot:      firstNonEmptyCaptureString(resolvedRoot, candidate.root),
			HookInstalled: installed,
			GitReachable:  reachable,
		}
		if !reachable {
			status.Error = "git repo not reachable"
		}
		statuses = append(statuses, status)
	}
	return statuses
}

func capturePrimaryHookStatus(ctx context.Context, repoRoot string) (installed bool, applicable bool, resolvedRoot string) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot != "" {
		return captureHookStatus(ctx, repoRoot)
	}
	if cwd, err := os.Getwd(); err == nil {
		installed, applicable, resolvedRoot = captureHookStatus(ctx, cwd)
		if applicable {
			return installed, true, resolvedRoot
		}
	}
	for _, candidate := range captureRepoRootCandidates(ctx) {
		installed, applicable, resolvedRoot = captureHookStatus(ctx, candidate.root)
		if applicable {
			return installed, true, resolvedRoot
		}
	}
	return false, false, ""
}

type captureRepoRootCandidate struct {
	root              string
	reportUnreachable bool
}

func captureRepoRootCandidates(ctx context.Context) []captureRepoRootCandidate {
	db, err := storage.Open(ctx)
	if err != nil {
		return nil
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return nil
	}
	defer store.Close()

	initialized, err := store.ListInitializedRepos(ctx)
	if err != nil {
		return nil
	}
	if len(initialized) > 0 {
		return captureInitializedRepoCandidates(initialized)
	}
	repos, err := store.ListRepos(ctx)
	if err != nil {
		return nil
	}
	return captureKnownRepoCandidates(repos)
}

func captureInitializedRepoCandidates(repos []storage.InitializedRepo) []captureRepoRootCandidate {
	roots := make([]captureRepoRootCandidate, 0, len(repos))
	seen := map[string]struct{}{}
	for _, repo := range repos {
		root := strings.TrimSpace(repo.RootPath)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, captureRepoRootCandidate{root: root, reportUnreachable: true})
	}
	return roots
}

func captureKnownRepoCandidates(repos []storage.Repo) []captureRepoRootCandidate {
	roots := make([]captureRepoRootCandidate, 0, len(repos))
	seen := map[string]struct{}{}
	for _, repo := range repos {
		root := strings.TrimSpace(repo.RootPath)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		roots = append(roots, captureRepoRootCandidate{root: root})
	}
	return roots
}

func captureHookStatus(ctx context.Context, repoRoot string) (installed bool, applicable bool, resolvedRoot string) {
	resolvedRoot, ok := resolveGitRoot(ctx, repoRoot)
	if !ok {
		return false, false, ""
	}
	return hooks.IsInstalled(resolvedRoot), true, resolvedRoot
}

func resolveGitRoot(ctx context.Context, startPath string) (string, bool) {
	startPath = strings.TrimSpace(startPath)
	if startPath == "" {
		return "", false
	}
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(checkCtx, "git", "-C", startPath, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", false
	}
	return root, true
}

func captureDoctorIssues(status captureDoctorJSON) []captureIssueJSON {
	issues := []captureIssueJSON{}
	if status.HookApplicable && !status.HookInstalled {
		issues = append(issues, captureIssueJSON{
			Code:     "hook_missing",
			Severity: "fail",
			Message:  "GX lifecycle hooks missing",
			Action:   "run `gx init` in this repo",
		})
	}
	if status.RepoHooksUnreachable > 0 {
		issues = append(issues, captureIssueJSON{
			Code:     "repo_hooks_unreachable",
			Severity: "fail",
			Message:  fmt.Sprintf("%d registered repo hooks unreachable", status.RepoHooksUnreachable),
			Action:   "remove stale repo registrations or restore the missing repos",
		})
	}
	if status.RepoHooksMissing > 0 {
		issues = append(issues, captureIssueJSON{
			Code:     "repo_hooks_missing",
			Severity: "fail",
			Message:  fmt.Sprintf("%d registered repo lifecycle hook sets missing", status.RepoHooksMissing),
			Action:   "run `gx init` in each registered repo",
		})
	}
	if !status.UploadAuthed {
		code := "upload_auth_missing"
		message := "upload auth missing"
		action := "run `gx auth login`"
		if strings.TrimSpace(status.UploadAuthError) != "" {
			code = "upload_auth_invalid"
			message = "upload auth invalid"
			action = status.UploadAuthError
		}
		issues = append(issues, captureIssueJSON{
			Code:     code,
			Severity: "fail",
			Message:  message,
			Action:   action,
		})
	}
	if !status.CursorReachable {
		action := "open Cursor once so GX can read state.vscdb"
		if strings.TrimSpace(status.CursorPath) != "" {
			action = "check Cursor state.vscdb at " + status.CursorPath
		}
		issues = append(issues, captureIssueJSON{
			Code:     "cursor_vscdb_missing",
			Severity: "fail",
			Message:  "Cursor state.vscdb missing",
			Action:   action,
		})
	}
	if status.DiskWarn {
		issues = append(issues, captureIssueJSON{
			Code:     "disk_low",
			Severity: "fail",
			Message:  fmt.Sprintf("disk free low: %d GB", status.DiskFreeGB),
			Action:   "free disk space before running more capture jobs",
		})
	}
	return issues
}

func firstNonEmptyCaptureString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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

func validateCaptureCloudSession(ctx context.Context, token string) (bool, string) {
	verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	validation, err := cloud.ValidateCloudAPISession(verifyCtx, nil, token)
	if err != nil {
		return false, "could not verify GX API session: " + err.Error()
	}
	if validation.Valid {
		return true, ""
	}
	message := strings.TrimSpace(validation.Error)
	if message == "" {
		message = "GX API rejected stored session"
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
