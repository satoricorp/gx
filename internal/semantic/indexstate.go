package semantic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// repoIndexStateVersion is the on-disk format version of the manifest itself.
const repoIndexStateVersion = 1

// FileIndexState records what was last uploaded for one file. Hash is over the
// raw file bytes, so an unchanged file is skipped before it is even chunked.
// ChunkHashes is positional: entry i corresponds to chunk index i, whose row id
// is derived from (repo, path, i). Keeping per-chunk hashes means editing one
// function in a large file re-embeds only the chunks that actually moved.
type FileIndexState struct {
	Hash        string   `json:"hash"`
	Size        int64    `json:"size"`
	ChunkHashes []string `json:"chunk_hashes"`
}

// RepoIndexState is the persisted manifest that makes re-indexing incremental.
// Every field that can change chunk content or vector shape is recorded; if any
// of them differs from the running configuration the manifest is discarded and
// the repository is rebuilt from scratch, because partial mixes of two chunk
// shapes inside one namespace are unrankable.
type RepoIndexState struct {
	Version        int                       `json:"version"`
	Namespace      string                    `json:"namespace"`
	RepoFullName   string                    `json:"repo_full_name"`
	RepoRoot       string                    `json:"repo_root"`
	EmbeddingModel string                    `json:"embedding_model"`
	Dimensions     int                       `json:"dimensions"`
	ChunkerVersion int                       `json:"chunker_version"`
	SchemaVersion  int                       `json:"schema_version"`
	CommitID       string                    `json:"commit_id,omitempty"`
	UpdatedAt      int64                     `json:"updated_at"`
	Files          map[string]FileIndexState `json:"files"`
}

// Compatible reports whether a manifest may be reused for an index run with the
// given parameters.
func (s *RepoIndexState) Compatible(namespace, model string, dimensions, chunkerVersion, schemaVersion int) bool {
	if s == nil || s.Version != repoIndexStateVersion {
		return false
	}
	return s.Namespace == namespace &&
		s.EmbeddingModel == model &&
		s.Dimensions == dimensions &&
		s.ChunkerVersion == chunkerVersion &&
		s.SchemaVersion == schemaVersion
}

// SortedFiles returns the manifest paths in deterministic order.
func (s *RepoIndexState) SortedFiles() []string {
	if s == nil || len(s.Files) == 0 {
		return nil
	}
	out := make([]string, 0, len(s.Files))
	for path := range s.Files {
		out = append(out, path)
	}
	sort.Strings(out)
	return out
}

var unsafeStatePathChars = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// RepoIndexStatePath is where the manifest for a namespace lives:
// $TOTALITY_HOME/index/<namespace>.json, defaulting to ~/.totality/index.
func RepoIndexStatePath(namespace string) (string, error) {
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		return "", fmt.Errorf("namespace is required for the index manifest path")
	}
	base := strings.TrimSpace(os.Getenv("TOTALITY_HOME"))
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".totality")
	}
	name := unsafeStatePathChars.ReplaceAllString(namespace, "-")
	if len(name) > 120 {
		name = name[:80] + "-" + shortHash(namespace)[:16]
	}
	return filepath.Join(base, "index", name+".json"), nil
}

// LoadRepoIndexState reads a manifest. A missing or unreadable manifest is not
// an error: it simply means the next run is a full index.
func LoadRepoIndexState(path string) *RepoIndexState {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var state RepoIndexState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil
	}
	if state.Files == nil {
		state.Files = map[string]FileIndexState{}
	}
	return &state
}

// SaveRepoIndexState writes the manifest atomically so an interrupted run
// cannot leave a half-written file that would be silently treated as truth.
func SaveRepoIndexState(path string, state *RepoIndexState) error {
	if state == nil {
		return fmt.Errorf("index state is nil")
	}
	state.Version = repoIndexStateVersion
	state.UpdatedAt = time.Now().Unix()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create index state dir: %w", err)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal index state: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("write index state: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("commit index state: %w", err)
	}
	return nil
}
