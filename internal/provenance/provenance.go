package provenance

import (
	"context"
	"os"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	StatusExplicit  = "explicit"
	StatusRepoLocal = "repo_local"
	StatusAbsent    = "absent"
)

type Attachment struct {
	SessionIDs []string
	Status     string
}

func Resolve(ctx context.Context, store *storage.Store, repoRoot string) (Attachment, error) {
	return resolveWithPreferred(ctx, store, repoRoot, nil)
}

func ResolveWithPreferred(ctx context.Context, store *storage.Store, repoRoot string, preferredSessionIDs []string) (Attachment, error) {
	return resolveWithPreferred(ctx, store, repoRoot, preferredSessionIDs)
}

func resolveWithPreferred(ctx context.Context, store *storage.Store, repoRoot string, preferredSessionIDs []string) (Attachment, error) {
	if len(preferredSessionIDs) > 0 {
		filtered, err := store.FilterExistingSessionIDs(ctx, preferredSessionIDs)
		if err != nil {
			return Attachment{}, err
		}
		if len(filtered) > 0 {
			return Attachment{SessionIDs: filtered, Status: StatusExplicit}, nil
		}
	}

	envIDs := ExplicitSessionIDsFromEnv()
	if len(envIDs) > 0 {
		filtered, err := store.FilterExistingSessionIDs(ctx, envIDs)
		if err != nil {
			return Attachment{}, err
		}
		if len(filtered) > 0 {
			return Attachment{SessionIDs: filtered, Status: StatusExplicit}, nil
		}
	}

	sessionIDs, err := store.FindAttachableSessionsForRepo(ctx, repoRoot, 20)
	if err != nil {
		return Attachment{}, err
	}
	if len(sessionIDs) > 0 {
		return Attachment{SessionIDs: sessionIDs, Status: StatusRepoLocal}, nil
	}
	return Attachment{Status: StatusAbsent}, nil
}

func Attach(ctx context.Context, store *storage.Store, repoRoot string, changeID int64, createdAt int64) (Attachment, error) {
	return AttachPreferred(ctx, store, repoRoot, changeID, createdAt, nil)
}

func AttachPreferred(ctx context.Context, store *storage.Store, repoRoot string, changeID int64, createdAt int64, preferredSessionIDs []string) (Attachment, error) {
	attachment, err := ResolveWithPreferred(ctx, store, repoRoot, preferredSessionIDs)
	if err != nil {
		return Attachment{}, err
	}
	if len(attachment.SessionIDs) == 0 {
		return attachment, nil
	}
	filtered, err := store.FilterExistingSessionIDs(ctx, attachment.SessionIDs)
	if err != nil {
		return Attachment{}, err
	}
	if len(filtered) == 0 {
		return Attachment{Status: StatusAbsent}, nil
	}
	if err := store.WriteChangeSessions(ctx, changeID, filtered, createdAt); err != nil {
		return Attachment{}, err
	}
	return Attachment{SessionIDs: filtered, Status: attachment.Status}, nil
}

func ExplicitSessionIDsFromEnv() []string {
	raw := strings.TrimSpace(os.Getenv("GX_SESSION_IDS"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("GX_SESSION_ID"))
	}
	return ParseSessionIDs(raw)
}

func ParseSessionIDs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t'
	})
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		out = append(out, part)
	}
	return out
}

func WithExplicitSessionEnv[T any](sessionIDs []string, fn func() (T, error)) (T, error) {
	if len(sessionIDs) == 0 {
		return fn()
	}
	oldID, hadID := os.LookupEnv("GX_SESSION_ID")
	oldIDs, hadIDs := os.LookupEnv("GX_SESSION_IDS")
	_ = os.Unsetenv("GX_SESSION_ID")
	_ = os.Setenv("GX_SESSION_IDS", strings.Join(sessionIDs, ","))
	defer func() {
		if hadID {
			_ = os.Setenv("GX_SESSION_ID", oldID)
		} else {
			_ = os.Unsetenv("GX_SESSION_ID")
		}
		if hadIDs {
			_ = os.Setenv("GX_SESSION_IDS", oldIDs)
		} else {
			_ = os.Unsetenv("GX_SESSION_IDS")
		}
	}()
	return fn()
}
