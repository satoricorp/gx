package publication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/reviewbundle"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

const (
	outboxDirName       = "publish-outbox"
	outboxStatusPending = "pending"
	outboxStatusRunning = "uploading"
	outboxStatusFailed  = "failed"

	// After this many failed uploads an item is quarantined: it stays on
	// disk for diagnostics but drains stop retrying it. Without a cap a
	// single oversized or rejected artifact is re-read, re-marshalled, and
	// re-uploaded on every drain forever.
	maxUploadAttempts = 5
)

type QueueAttestation struct {
	ContentHash   string `json:"content_hash"`
	AcceptorName  string `json:"acceptor_name,omitempty"`
	AcceptorEmail string `json:"acceptor_email,omitempty"`
	CreatedAt     int64  `json:"created_at"`
}

type QueueItem struct {
	ID            string                `json:"id"`
	Status        string                `json:"status"`
	Artifact      reviewbundle.Artifact `json:"artifact"`
	Attestation   QueueAttestation      `json:"attestation,omitempty"`
	CreatedAt     int64                 `json:"created_at"`
	Attempts      int                   `json:"attempts"`
	LastAttemptAt int64                 `json:"last_attempt_at,omitempty"`
	LastError     string                `json:"last_error,omitempty"`
}

type uploadLock struct {
	StartedAt int64 `json:"started_at"`
	PID       int   `json:"pid"`
}

type QueueStatus struct {
	Pending       int    `json:"pending"`
	Failed        int    `json:"failed"`
	Uploading     int    `json:"uploading"`
	OldestAt      int64  `json:"oldest_at,omitempty"`
	LastError     string `json:"last_error,omitempty"`
	LastErrorAt   int64  `json:"last_error_at,omitempty"`
	OutboxPath    string `json:"outbox_path,omitempty"`
	UploadRunning bool   `json:"upload_running"`
}

type DrainResult struct {
	Uploaded            int
	Failed              int
	Pending             int
	PRSummariesUpdated  int
	PRSummariesFailed   int
}

func EnqueuePush(ctx context.Context, push vcs.PushResult) (Result, error) {
	bundle, err := reviewbundle.BuildPush(ctx, push)
	if err != nil {
		return Result{}, err
	}
	attestation, err := queueAttestationForPush(ctx, push.Repo.RootPath, reviewbundle.NewArtifact(bundle))
	if err != nil {
		return Result{}, err
	}
	return EnqueueArtifact(ctx, reviewbundle.NewArtifact(bundle), attestation)
}

func EnqueueArtifact(ctx context.Context, artifact reviewbundle.Artifact, attestation QueueAttestation) (Result, error) {
	artifact.IndexStatus = firstNonEmpty(artifact.IndexStatus, "queued")
	item, err := newQueueItem(artifact, attestation)
	if err != nil {
		return Result{}, err
	}
	if existing, err := loadQueueItem(item.ID); err == nil {
		// Re-enqueues for the same head carry newer data (the pre-push hook
		// runs before gx push resolves the PR URL); the latest artifact must
		// win or PR linkage is lost. Leave in-flight uploads alone.
		if existing.Status == outboxStatusRunning && uploadLockActive() && !staleUpload(existing.LastAttemptAt) {
			artifactPath, artErr := WriteLocalArtifact(artifact)
			if artErr != nil {
				return Result{}, artErr
			}
			return Result{
				Queued:       true,
				QueueID:      item.ID,
				QueuePath:    queueItemPath(item.ID),
				IndexStatus:  "queued",
				ArtifactPath: artifactPath,
				Artifact:     artifact,
			}, nil
		}
		item.CreatedAt = existing.CreatedAt
		item.Attempts = existing.Attempts
	}
	path, err := writeQueueItem(item)
	if err != nil {
		return Result{}, err
	}
	artifactPath, err := WriteLocalArtifact(artifact)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Queued:       true,
		QueueID:      item.ID,
		QueuePath:    path,
		IndexStatus:  "queued",
		ArtifactPath: artifactPath,
		Artifact:     artifact,
	}, nil
}

func DrainQueuedUploads(ctx context.Context, uploader Uploader, limit int) (DrainResult, error) {
	if uploader == nil {
		return DrainResult{}, errors.New("publish upload client is not configured")
	}
	if limit <= 0 {
		limit = 20
	}
	if err := writeUploadLock(); err != nil {
		return DrainResult{}, err
	}
	defer removeUploadLock()

	items, err := loadQueueItems()
	if err != nil {
		return DrainResult{}, err
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt < items[j].CreatedAt
	})

	result := DrainResult{}
	for _, item := range items {
		if item.Attempts >= maxUploadAttempts {
			continue
		}
		if result.Uploaded+result.Failed >= limit {
			result.Pending++
			continue
		}
		item.Status = outboxStatusRunning
		item.Attempts++
		item.LastAttemptAt = time.Now().UnixMilli()
		item.LastError = ""
		_, _ = writeQueueItem(item)

		updated, err := UpdateGitHubPullRequestBody(ctx, item.Artifact)
		if err != nil {
			result.PRSummariesFailed++
		} else if updated {
			result.PRSummariesUpdated++
		}

		uploaded, err := uploader.UploadReviewArtifact(ctx, item.Artifact)
		if err != nil {
			item.Status = outboxStatusFailed
			item.LastError = err.Error()
			_, _ = writeQueueItem(item)
			result.Failed++
			continue
		}
		if _, err := WriteLocalArtifact(uploaded); err != nil {
			item.Status = outboxStatusFailed
			item.LastError = err.Error()
			_, _ = writeQueueItem(item)
			result.Failed++
			continue
		}
		if err := os.Remove(queueItemPath(item.ID)); err != nil && !os.IsNotExist(err) {
			return result, err
		}
		result.Uploaded++
	}
	return result, nil
}

func QueuedUploadStatus() (QueueStatus, error) {
	dir, err := outboxDir()
	if err != nil {
		return QueueStatus{}, err
	}
	items, err := loadQueueItems()
	if err != nil {
		return QueueStatus{}, err
	}
	lockActive := uploadLockActive()
	status := QueueStatus{
		OutboxPath:    dir,
		UploadRunning: lockActive,
	}
	for _, item := range items {
		switch item.Status {
		case outboxStatusFailed:
			status.Failed++
			if item.LastError != "" && item.LastAttemptAt >= status.LastErrorAt {
				status.LastError = item.LastError
				status.LastErrorAt = item.LastAttemptAt
			}
		case outboxStatusRunning:
			if !lockActive || staleUpload(item.LastAttemptAt) {
				status.Pending++
			} else {
				status.Uploading++
			}
		default:
			status.Pending++
		}
		if item.CreatedAt > 0 && (status.OldestAt == 0 || item.CreatedAt < status.OldestAt) {
			status.OldestAt = item.CreatedAt
		}
	}
	return status, nil
}

func newQueueItem(artifact reviewbundle.Artifact, attestation QueueAttestation) (QueueItem, error) {
	payload, err := json.Marshal(artifact)
	if err != nil {
		return QueueItem{}, err
	}
	sum := sha256.Sum256(payload)
	contentHash := hex.EncodeToString(sum[:])
	if attestation.ContentHash == "" {
		attestation.ContentHash = contentHash
	}
	if attestation.CreatedAt == 0 {
		attestation.CreatedAt = time.Now().UnixMilli()
	}
	id := queueItemID(artifact, attestation)
	return QueueItem{
		ID:          id,
		Status:      outboxStatusPending,
		Artifact:    artifact,
		Attestation: attestation,
		CreatedAt:   attestation.CreatedAt,
	}, nil
}

func writeQueueItem(item QueueItem) (string, error) {
	dir, err := outboxDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := queueItemPath(item.ID)
	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, path); err != nil {
		return "", err
	}
	return path, nil
}

func shortHashPrefix(hash string) string {
	hash = strings.TrimSpace(hash)
	if len(hash) > 24 {
		return hash[:24]
	}
	return hash
}

func queueItemID(artifact reviewbundle.Artifact, attestation QueueAttestation) string {
	if head := strings.TrimSpace(artifact.Push.HeadCommitID); head != "" {
		return "gx-context-" + shortClean(head)
	}
	return "gx-context-" + shortHashPrefix(attestation.ContentHash)
}

func loadQueueItem(id string) (QueueItem, error) {
	data, err := os.ReadFile(queueItemPath(id))
	if err != nil {
		return QueueItem{}, err
	}
	var item QueueItem
	if err := json.Unmarshal(data, &item); err != nil {
		return QueueItem{}, err
	}
	if item.ID == "" {
		item.ID = id
	}
	return item, nil
}

func queueAttestationForPush(ctx context.Context, repoRoot string, artifact reviewbundle.Artifact) (QueueAttestation, error) {
	payload, err := json.Marshal(artifact)
	if err != nil {
		return QueueAttestation{}, err
	}
	sum := sha256.Sum256(payload)
	att := QueueAttestation{
		ContentHash: hex.EncodeToString(sum[:]),
		CreatedAt:   time.Now().UnixMilli(),
	}
	if repoRoot == "" {
		return att, nil
	}
	name, err := gitConfigValue(ctx, repoRoot, "user.name")
	if err == nil {
		att.AcceptorName = name
	}
	email, err := gitConfigValue(ctx, repoRoot, "user.email")
	if err == nil {
		att.AcceptorEmail = email
	}
	return att, nil
}

func gitConfigValue(ctx context.Context, repoRoot, key string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "config", key)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git config %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func loadQueueItems() ([]QueueItem, error) {
	dir, err := outboxDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var items []QueueItem
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		var item QueueItem
		if err := json.Unmarshal(data, &item); err != nil {
			return nil, fmt.Errorf("decode %s: %w", entry.Name(), err)
		}
		if item.ID == "" {
			item.ID = strings.TrimSuffix(entry.Name(), ".json")
		}
		if item.Status == "" {
			item.Status = outboxStatusPending
		}
		items = append(items, item)
	}
	return items, nil
}

func outboxDir() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, outboxDirName), nil
}

func queueItemPath(id string) string {
	dir, _ := outboxDir()
	return filepath.Join(dir, id+".json")
}

func writeUploadLock() error {
	path, err := uploadLockPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if uploadLockActive() {
		return errors.New("publish context upload already running")
	}
	_ = os.Remove(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if os.IsExist(err) {
		return errors.New("publish context upload already running")
	}
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(uploadLock{
		StartedAt: time.Now().UnixMilli(),
		PID:       os.Getpid(),
	})
}

func removeUploadLock() {
	if path, err := uploadLockPath(); err == nil {
		_ = os.Remove(path)
	}
}

func uploadLockActive() bool {
	path, err := uploadLockPath()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	lock, err := parseUploadLock(data)
	if err != nil {
		return false
	}
	if staleUpload(lock.StartedAt) {
		return false
	}
	if lock.PID <= 0 {
		return false
	}
	return processAlive(lock.PID)
}

func parseUploadLock(data []byte) (uploadLock, error) {
	var lock uploadLock
	if err := json.Unmarshal(data, &lock); err == nil && lock.StartedAt > 0 {
		return lock, nil
	}
	var ts int64
	var pid int
	fields, err := fmt.Sscanf(strings.TrimSpace(string(data)), "%d %d", &ts, &pid)
	if err != nil && fields == 0 {
		return uploadLock{}, err
	}
	return uploadLock{StartedAt: ts, PID: pid}, nil
}

func uploadLockPath() (string, error) {
	dir, err := outboxDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".upload.lock"), nil
}

func staleUpload(lastAttemptAt int64) bool {
	if lastAttemptAt == 0 {
		return true
	}
	return time.Since(time.UnixMilli(lastAttemptAt)) > 10*time.Minute
}

func shortClean(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 12 {
		value = value[:12]
	}
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		default:
			return '-'
		}
	}, value)
}
