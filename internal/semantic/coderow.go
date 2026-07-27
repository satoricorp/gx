package semantic

import (
	"fmt"
	"strings"
)

// IndexSchemaVersion is the version of the TurboPuffer namespace schema (field
// set, full-text configuration and vector dimensions). It is part of the
// namespace name, so bumping it provisions a fresh namespace instead of
// colliding with rows written under the old shape — TurboPuffer pins vector
// dimensions and full-text settings per namespace, and an in-place change would
// either be rejected or silently corrupt ranking.
const IndexSchemaVersion = 2

// SourceKindCodeFile is the row kind for indexed repository source.
const SourceKindCodeFile = "code_file"

const (
	codeFieldOrgID        = "org_id"
	codeFieldRepoFullName = "repo_full_name"
	codeFieldFilePath     = "file_path"
	// codeFieldFile mirrors file_path under the attribute name the console
	// writer uses, so one filter works across rows from either producer.
	codeFieldFile          = "file"
	codeFieldSymbol        = "symbol"
	codeFieldSymbolName    = "symbol_name"
	codeFieldSymbolKind    = "symbol_kind"
	codeFieldStartLine     = "start_line"
	codeFieldEndLine       = "end_line"
	codeFieldLanguage      = "language"
	codeFieldDocType       = "doc_type"
	codeFieldPackage       = "package_name"
	codeFieldChunkHash     = "chunk_hash"
	codeFieldIndexedAt     = "indexed_at"
	codeFieldIndexedReason = "indexed_reason"
	codeFieldHeadSha       = "head_sha"
)

// CodeRowMetadata is the attribute set written for one code_file row.
type CodeRowMetadata struct {
	SourceID      string
	OrgID         string
	RepoFullName  string
	RepoRoot      string
	BranchName    string
	CommitID      string
	FilePath      string
	Symbol        string
	SymbolText    string
	SymbolKind    string
	Package       string
	StartLine     int
	EndLine       int
	Language      string
	DocType       string
	ChunkHash     string
	IndexedAt     int64
	IndexedReason string
	Text          string
}

func (m CodeRowMetadata) Attributes() map[string]any {
	return map[string]any{
		transcriptFieldSourceKind: SourceKindCodeFile,
		transcriptFieldSourceID:   m.SourceID,
		transcriptFieldRepoRoot:   m.RepoRoot,
		transcriptFieldBranchName: m.BranchName,
		transcriptFieldCommitID:   m.CommitID,
		transcriptFieldText:       m.Text,
		codeFieldOrgID:            m.OrgID,
		codeFieldRepoFullName:     m.RepoFullName,
		codeFieldFilePath:         m.FilePath,
		codeFieldFile:             m.FilePath,
		codeFieldSymbol:           m.SymbolText,
		codeFieldSymbolName:       m.Symbol,
		codeFieldSymbolKind:       m.SymbolKind,
		codeFieldPackage:          m.Package,
		codeFieldStartLine:        m.StartLine,
		codeFieldEndLine:          m.EndLine,
		codeFieldLanguage:         m.Language,
		codeFieldDocType:          m.DocType,
		codeFieldChunkHash:        m.ChunkHash,
		codeFieldIndexedAt:        m.IndexedAt,
		codeFieldIndexedReason:    m.IndexedReason,
		codeFieldHeadSha:          m.CommitID,
	}
}

// CodeRowID is the TurboPuffer row id for a chunk.
//
// It is derived from (repo, file path, chunk index) and deliberately NOT from
// the commit id. Folding the commit into the id — as the previous
// implementation did — minted a brand new row for every chunk on every commit
// and then relied on a delete-by-filter sweep to evict the old generation,
// which left the namespace half-populated whenever a batch failed partway. A
// commit-independent id means a re-index overwrites rows in place, and only
// chunks that genuinely disappeared need deleting.
func CodeRowID(repoFullName, filePath string, chunkIndex int) string {
	return "gxc-" + shortHash(fmt.Sprintf("%s\x00%s\x00%d", repoFullName, filePath, chunkIndex))
}

func codeSourceID(repoFullName, filePath string, chunkIndex int) string {
	return strings.Join([]string{"code", repoFullName, filePath, fmt.Sprintf("%d", chunkIndex)}, ":")
}
