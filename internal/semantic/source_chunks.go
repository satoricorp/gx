package semantic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

const (
	sourceChunkLines       = 100
	sourceOverlapLines     = 15
	sourceMaxFileBytes     = 100 * 1024
	sourceMaxFiles         = 5000
	sourceMaxChunksPerFile = 20
)

var sourceSkipDirPrefixes = []string{
	"node_modules/",
	".git/",
	"dist/",
	"build/",
	".next/",
	"coverage/",
	"vendor/",
	".turbo/",
	".cache/",
}

var sourceSkipExtensions = map[string]struct{}{
	".png":   {},
	".jpg":   {},
	".jpeg":  {},
	".gif":   {},
	".webp":  {},
	".ico":   {},
	".svg":   {},
	".woff":  {},
	".woff2": {},
	".ttf":   {},
	".eot":   {},
	".mp4":   {},
	".mp3":   {},
	".zip":   {},
	".tar":   {},
	".gz":    {},
	".pdf":   {},
	".exe":   {},
	".dll":   {},
	".so":    {},
	".dylib": {},
	".lock":  {},
}

var sourcePriorityPrefixes = []string{
	"src/",
	"app/",
	"lib/",
	"packages/",
	"convex/",
	"services/",
}

type sourceFile struct {
	path string
	size int64
}

type sourceRange struct {
	startLineIndex int
	endLineIndex   int
	symbol         string
}

// BuildRepositoryChunks mirrors the working Console codebase indexing shape:
// prioritize source paths, chunk by coarse symbols, and attach repo metadata so
// code context can live beside session and patch context in TurboPuffer.
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

	files := listSourceFiles(root)
	if len(files) > sourceMaxFiles {
		files = files[:sourceMaxFiles]
	}

	var chunks []Chunk
	for _, file := range files {
		body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file.path)))
		if err != nil || len(body) == 0 || strings.ContainsRune(string(body), '\x00') {
			continue
		}
		fileChunks := chunkSourceFile(fullName, root, commitID, branchName, file.path, string(body))
		if len(fileChunks) > sourceMaxChunksPerFile {
			fileChunks = fileChunks[:sourceMaxChunksPerFile]
		}
		chunks = append(chunks, fileChunks...)
	}
	return chunks
}

func listSourceFiles(root string) []sourceFile {
	var files []sourceFile
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() {
			if shouldSkipSourceDir(rel + "/") {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil || !shouldIndexSourcePath(rel, info.Size()) {
			return nil
		}
		files = append(files, sourceFile{path: rel, size: info.Size()})
		return nil
	})
	sort.Slice(files, func(i, j int) bool {
		if pi, pj := sourcePathPriority(files[i].path), sourcePathPriority(files[j].path); pi != pj {
			return pi < pj
		}
		return files[i].path < files[j].path
	})
	return files
}

func shouldSkipSourceDir(path string) bool {
	for _, prefix := range sourceSkipDirPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix) || strings.Contains(path, "/"+prefix) {
			return true
		}
	}
	return false
}

func shouldIndexSourcePath(path string, size int64) bool {
	if shouldSkipSourceDir(path) || size > sourceMaxFileBytes {
		return false
	}
	if _, ok := sourceSkipExtensions[strings.ToLower(filepath.Ext(path))]; ok {
		return false
	}
	return true
}

func sourcePathPriority(path string) int {
	lower := strings.ToLower(path)
	for index, prefix := range sourcePriorityPrefixes {
		if strings.HasPrefix(lower, prefix) || strings.Contains(lower, "/"+prefix) {
			return index
		}
	}
	return len(sourcePriorityPrefixes)
}

func chunkSourceFile(fullName, repoRoot, commitID, branchName, filePath, source string) []Chunk {
	lines := strings.Split(source, "\n")
	language := languageFromPath(filePath)
	docType := docTypeFromPath(filePath)
	ranges := sourceRanges(lines, language)
	chunks := make([]Chunk, 0, len(ranges))
	for index, r := range ranges {
		startLine := r.startLineIndex + 1
		endLine := r.endLineIndex + 1
		body := strings.Join(lines[r.startLineIndex:r.endLineIndex+1], "\n")
		header := []string{
			fmt.Sprintf("repo: %s", fullName),
			fmt.Sprintf("file: %s", filePath),
			fmt.Sprintf("lines: %d-%d", startLine, endLine),
			fmt.Sprintf("language: %s", language),
			fmt.Sprintf("doc_type: %s", docType),
		}
		if r.symbol != "" {
			header = append(header[:3], append([]string{fmt.Sprintf("symbol: %s", r.symbol)}, header[3:]...)...)
		}
		text := strings.Join(header, "\n") + "\n---\n" + body
		chunks = append(chunks, Chunk{
			ID:   sourceDocumentID(fullName, commitID, filePath, index),
			Text: text,
			Attributes: map[string]any{
				"source_kind":    "code_file",
				"repo_full_name": fullName,
				"repo_root":      repoRoot,
				"branch_name":    branchName,
				"commit_id":      commitID,
				"file_path":      filePath,
				"symbol":         r.symbol,
				"start_line":     startLine,
				"end_line":       endLine,
				"chunk_hash":     shortHash(fmt.Sprintf("%s:%d:%d:%s", filePath, startLine, endLine, text)),
				"language":       language,
				"doc_type":       docType,
			},
		})
	}
	return chunks
}

func sourceRanges(lines []string, language string) []sourceRange {
	if len(lines) == 0 {
		return nil
	}
	boundaries := detectSourceSymbols(lines, language)
	if len(boundaries) == 0 {
		return splitSourceRange(0, len(lines)-1, "")
	}
	var ranges []sourceRange
	if boundaries[0].startLineIndex > 0 {
		ranges = append(ranges, sourceRange{startLineIndex: 0, endLineIndex: boundaries[0].startLineIndex - 1})
	}
	for i, boundary := range boundaries {
		end := len(lines) - 1
		if i+1 < len(boundaries) {
			end = boundaries[i+1].startLineIndex - 1
		}
		if end >= boundary.startLineIndex {
			ranges = append(ranges, splitSourceRange(boundary.startLineIndex, end, boundary.symbol)...)
		}
	}
	return ranges
}

func splitSourceRange(start, end int, symbol string) []sourceRange {
	var ranges []sourceRange
	for current := start; current <= end; {
		chunkEnd := current + sourceChunkLines - 1
		if chunkEnd > end {
			chunkEnd = end
		}
		ranges = append(ranges, sourceRange{startLineIndex: current, endLineIndex: chunkEnd, symbol: symbol})
		if chunkEnd >= end {
			break
		}
		current = chunkEnd - sourceOverlapLines + 1
	}
	return ranges
}

func detectSourceSymbols(lines []string, language string) []sourceRange {
	var out []sourceRange
	for idx, line := range lines {
		if symbol := detectSourceSymbol(line, language); symbol != "" {
			out = append(out, sourceRange{startLineIndex: idx, symbol: symbol})
		}
	}
	return out
}

var sourceSymbolMatchers = map[string][]*regexp.Regexp{
	"go": {
		regexp.MustCompile(`^func\s+(?:\([^)]+\)\s*)?([A-Za-z_]\w*)`),
	},
	"python": {
		regexp.MustCompile(`^(?:async\s+)?def\s+([A-Za-z_]\w*)`),
		regexp.MustCompile(`^class\s+([A-Za-z_]\w*)`),
	},
	"javascript": javascriptSymbolMatchers(),
	"typescript": javascriptSymbolMatchers(),
	"rust": {
		regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?(?:async\s+)?fn\s+([A-Za-z_]\w*)`),
		regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?(?:struct|enum|trait|impl)\s+([A-Za-z_]\w*)`),
	},
}

func javascriptSymbolMatchers() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`^(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+([A-Za-z_$][\w$]*)`),
		regexp.MustCompile(`^(?:export\s+)?(?:default\s+)?class\s+([A-Za-z_$][\w$]*)`),
		regexp.MustCompile(`^(?:export\s+)?(?:interface|type|enum)\s+([A-Za-z_$][\w$]*)`),
		regexp.MustCompile(`^(?:export\s+)?(?:const|let|var)\s+([A-Za-z_$][\w$]*)\s*=`),
	}
}

func detectSourceSymbol(line, language string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
		return ""
	}
	for _, matcher := range sourceSymbolMatchers[language] {
		if match := matcher.FindStringSubmatch(trimmed); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func languageFromPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".rb":
		return "ruby"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "text"
	}
}

func docTypeFromPath(path string) string {
	lower := strings.ToLower(path)
	if strings.Contains(lower, "__tests__/") || strings.Contains(lower, "/test/") ||
		strings.Contains(lower, ".test.") || strings.Contains(lower, ".spec.") {
		return "test_file"
	}
	if lower == "readme.md" || strings.HasPrefix(lower, "docs/") ||
		strings.HasSuffix(lower, ".md") || strings.Contains(lower, "/adr/") {
		return "architecture_doc"
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

func sourceDocumentID(fullName, commitID, filePath string, chunkIndex int) string {
	return "gx-code-" + shortHash(fmt.Sprintf("%s:%s:%s:%d", fullName, commitID, filePath, chunkIndex))
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
