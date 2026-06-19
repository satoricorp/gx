package semantic

import "fmt"

const (
	transcriptFieldText             = "text"
	transcriptFieldSourceKind       = "source_kind"
	transcriptFieldSourceID         = "source_id"
	transcriptFieldRepoRoot         = "repo_root"
	transcriptFieldBranchName       = "branch_name"
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
	SourceID         string
	RepoRoot         string
	BranchName       string
	JJChangeID       string
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
		transcriptFieldSourceKind:       "session_transcript",
		transcriptFieldSourceID:         m.SourceID,
		transcriptFieldRepoRoot:         m.RepoRoot,
		transcriptFieldBranchName:       m.BranchName,
		transcriptFieldJJChangeID:       m.JJChangeID,
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

func transcriptTurboPufferSchema(dimensions int) map[string]any {
	stringFilter := map[string]any{"type": "string", "filterable": true}
	return map[string]any{
		"vector": map[string]any{
			"type": fmt.Sprintf("[%d]f32", dimensions),
			"ann":  true,
		},
		transcriptFieldText: map[string]any{
			"type":             "string",
			"full_text_search": true,
		},
		transcriptFieldSourceKind:       stringFilter,
		transcriptFieldSourceID:         stringFilter,
		transcriptFieldRepoRoot:         stringFilter,
		transcriptFieldBranchName:       stringFilter,
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
		"repo_full_name":                stringFilter,
		"file_path":                     stringFilter,
		"symbol":                        map[string]any{"type": "string", "full_text_search": true},
		"start_line":                    map[string]any{"type": "uint", "filterable": true},
		"end_line":                      map[string]any{"type": "uint", "filterable": true},
		"chunk_hash":                    stringFilter,
		"language":                      stringFilter,
		"doc_type":                      stringFilter,
	}
}
