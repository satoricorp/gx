package semantic

import "fmt"

const (
	transcriptFieldText       = "text"
	transcriptFieldSourceKind = "source_kind"
	transcriptFieldSourceID   = "source_id"
	transcriptFieldRepoRoot   = "repo_root"
	transcriptFieldBranchName = "branch_name"
	transcriptFieldRevisionID = "revision_id"
	// transcriptFieldJJChangeID is the legacy attribute name for the lgtm
	// revision id (commit-SHA fallback) from the jj era. The lgtm-sessions
	// namespace is upserted into and never rebuilt wholesale, so rows written
	// before the rename carry only this key. Dual-write keeps those rows and
	// new ones queryable by one shared name during the transition; delete
	// this field (and the fallback expectation in any future reader) once a
	// full reindex has rewritten every row with revision_id.
	transcriptFieldJJChangeID       = "jj_change_id"
	transcriptFieldCommitID         = "commit_id"
	transcriptFieldRevisionTitle    = "revision_title"
	transcriptFieldSessionID        = "session_id"
	transcriptFieldRequestID        = "request_id"
	transcriptFieldResponseID       = "response_id"
	transcriptFieldAgentTool        = "agent_tool"
	transcriptFieldProvider         = "provider"
	transcriptFieldModel            = "model"
	transcriptFieldCreatedAt        = "created_at"
	transcriptFieldProvenanceStatus = "provenance_status"
)

type TranscriptMetadata struct {
	SourceID   string
	RepoRoot   string
	BranchName string
	// RevisionID carries the lgtm revision id, falling back to the commit SHA
	// for commits without a lgtm trailer.
	RevisionID       string
	CommitID         string
	RevisionTitle    string
	SessionID        string
	RequestID        string
	ResponseID       string
	AgentTool        string
	Provider         string
	Model            string
	CreatedAt        int64
	ProvenanceStatus string
	Text             string
}

func (m TranscriptMetadata) Attributes() map[string]any {
	return map[string]any{
		transcriptFieldSourceKind: "session_transcript",
		transcriptFieldSourceID:   m.SourceID,
		transcriptFieldRepoRoot:   m.RepoRoot,
		transcriptFieldBranchName: m.BranchName,
		transcriptFieldRevisionID: m.RevisionID,
		// Dual-write the legacy name so old rows and new rows stay queryable
		// by one attribute; see transcriptFieldJJChangeID.
		transcriptFieldJJChangeID:       m.RevisionID,
		transcriptFieldCommitID:         m.CommitID,
		transcriptFieldRevisionTitle:    m.RevisionTitle,
		transcriptFieldSessionID:        m.SessionID,
		transcriptFieldRequestID:        m.RequestID,
		transcriptFieldResponseID:       m.ResponseID,
		transcriptFieldAgentTool:        m.AgentTool,
		transcriptFieldProvider:         m.Provider,
		transcriptFieldModel:            m.Model,
		transcriptFieldCreatedAt:        m.CreatedAt,
		transcriptFieldProvenanceStatus: m.ProvenanceStatus,
		transcriptFieldText:             m.Text,
	}
}

// transcriptTurboPufferSchema is the single schema every lgtm writer declares.
// One namespace holds session transcripts, published diffs and code chunks, so
// the schema is the union of their fields; TurboPuffer only stores attributes a
// row actually sets, so unused columns cost nothing.
//
// Two full-text columns are declared on purpose. `text` carries the embedded
// body and is stemmed, which is what conceptual queries need ("retry" must
// reach "retries"). `symbol` carries identifier terms and is NOT stemmed,
// because identifier lookups must stay exact; sub-word recall is handled by
// writing pre-split word parts into the value (see CodeChunk.SymbolText)
// rather than by asking the tokenizer to split camelCase, which it does not do.
func transcriptTurboPufferSchema(dimensions int) map[string]any {
	stringFilter := map[string]any{"type": "string", "filterable": true}
	return map[string]any{
		"vector": map[string]any{
			"type": fmt.Sprintf("[%d]f32", dimensions),
			"ann":  true,
		},
		transcriptFieldText: map[string]any{
			"type": "string",
			"full_text_search": map[string]any{
				"stemming":         true,
				"remove_stopwords": false,
				"case_sensitive":   false,
			},
		},
		codeFieldSymbol: map[string]any{
			"type": "string",
			"full_text_search": map[string]any{
				"stemming":         false,
				"remove_stopwords": false,
				"case_sensitive":   false,
			},
		},
		transcriptFieldSourceKind: stringFilter,
		transcriptFieldSourceID:   stringFilter,
		transcriptFieldRepoRoot:   stringFilter,
		transcriptFieldBranchName: stringFilter,
		transcriptFieldRevisionID: stringFilter,
		// Legacy twin of revision_id; drop after a full reindex.
		transcriptFieldJJChangeID:       stringFilter,
		transcriptFieldCommitID:         stringFilter,
		transcriptFieldRevisionTitle:    map[string]any{"type": "string"},
		transcriptFieldSessionID:        stringFilter,
		transcriptFieldRequestID:        stringFilter,
		transcriptFieldResponseID:       stringFilter,
		transcriptFieldAgentTool:        stringFilter,
		transcriptFieldProvider:         stringFilter,
		transcriptFieldModel:            stringFilter,
		transcriptFieldCreatedAt:        map[string]any{"type": "uint"},
		transcriptFieldProvenanceStatus: stringFilter,
		codeFieldOrgID:                  stringFilter,
		codeFieldRepoFullName:           stringFilter,
		codeFieldFilePath:               stringFilter,
		codeFieldFile:                   stringFilter,
		codeFieldSymbolName:             stringFilter,
		codeFieldSymbolKind:             stringFilter,
		codeFieldPackage:                stringFilter,
		codeFieldStartLine:              map[string]any{"type": "uint", "filterable": true},
		codeFieldEndLine:                map[string]any{"type": "uint", "filterable": true},
		codeFieldChunkHash:              stringFilter,
		codeFieldLanguage:               stringFilter,
		codeFieldDocType:                stringFilter,
		codeFieldIndexedAt:              map[string]any{"type": "uint"},
		codeFieldIndexedReason:          stringFilter,
		codeFieldHeadSha:                stringFilter,
	}
}
