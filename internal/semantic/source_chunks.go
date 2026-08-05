package semantic

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/lgtm/internal/reviewbundle"
)

// BuildRepositoryChunks chunks every source file in a bundle's checkout.
//
// It shares one chunker with the on-demand repository indexer
// (ChunkSourceFile), so a chunk written by the publish path and one written by
// `lgtm index` are byte-identical and share a row id, instead of being two
// incompatible shapes inside the same namespace.
func BuildRepositoryChunks(bundle reviewbundle.Bundle) []Chunk {
	root := strings.TrimSpace(bundle.Repo.RootPath)
	if root == "" {
		return nil
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil
	}
	fullName := repoFullName(bundle)
	commitID := strings.TrimSpace(bundle.Push.HeadCommitID)
	if commitID == "" {
		commitID = repositorySnapshotID(bundle)
	}
	branchName := firstNonEmptyString(deref(bundle.Push.BranchName), deref(bundle.Repo.BranchName), deref(bundle.Repo.DefaultBranch))

	files, err := listRepositoryFiles(root)
	if err != nil {
		return nil
	}
	ctx := codeRowContext{
		RepoRoot:     root,
		RepoFullName: fullName,
		CommitID:     commitID,
		BranchName:   branchName,
		Reason:       "publish",
		IndexedAt:    time.Now().Unix(),
	}

	var chunks []Chunk
	for _, file := range files {
		body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file.path)))
		if err != nil || len(body) == 0 || isBinaryContent(body) {
			continue
		}
		for index, chunk := range ChunkSourceFile(file.path, string(body)) {
			chunks = append(chunks, codeChunkRow(ctx, chunk, index))
		}
	}
	return chunks
}

func languageFromPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ts", ".tsx", ".mts", ".cts":
		return "typescript"
	case ".js", ".jsx", ".mjs", ".cjs":
		return "javascript"
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".kt", ".kts":
		return "kotlin"
	case ".rb":
		return "ruby"
	case ".swift":
		return "swift"
	case ".c", ".h":
		return "c"
	case ".cc", ".cpp", ".hpp", ".hh":
		return "cpp"
	case ".sh", ".bash", ".zsh":
		return "shell"
	case ".sql":
		return "sql"
	case ".md", ".mdx":
		return "markdown"
	case ".json":
		return "json"
	case ".toml":
		return "toml"
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "text"
	}
}

func docTypeFromPath(path string) string {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "__tests__/") || strings.Contains(lower, "/test/") ||
		strings.HasPrefix(lower, "test/") || strings.HasSuffix(lower, "_test.go") ||
		strings.Contains(lower, ".test.") || strings.Contains(lower, ".spec.") {
		return "test_file"
	}
	if strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".mdx") ||
		strings.HasPrefix(lower, "docs/") || strings.Contains(lower, "/adr/") {
		return "architecture_doc"
	}
	if strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".yaml") ||
		strings.HasSuffix(lower, ".yml") || strings.HasSuffix(lower, ".toml") {
		return "config_file"
	}
	return "code_chunk"
}

func repoFullName(bundle reviewbundle.Bundle) string {
	for _, value := range []string{deref(bundle.Repo.RemoteURL), deref(bundle.Repo.DefaultRemote)} {
		if full := repoFullNameFromRemoteURL(value); full != "" {
			return full
		}
	}
	return ""
}

func repoFullNameFromRemoteURL(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return ""
	}
	if match := strings.TrimPrefix(remoteURL, "git@github.com:"); match != remoteURL {
		parts := strings.Split(strings.TrimSuffix(match, ".git"), "/")
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}
	if strings.Contains(remoteURL, "github.com") {
		idx := strings.Index(remoteURL, "github.com/")
		if idx >= 0 {
			path := strings.TrimSuffix(remoteURL[idx+len("github.com/"):], ".git")
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				return parts[0] + "/" + parts[1]
			}
		}
	}
	return ""
}

func shortHash(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])[:40]
}

func repositorySnapshotID(bundle reviewbundle.Bundle) string {
	return shortHash(strings.Join([]string{bundle.Repo.RootPath, deref(bundle.Repo.BranchName), deref(bundle.Repo.RemoteURL)}, ":"))
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
