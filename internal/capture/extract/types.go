package extract

import "github.com/satoricorp/lgtm/internal/capture/matcher"

// ExtractRequest is the POST /v1/extracts body.
type ExtractRequest struct {
	RepoRoot         string             `json:"repoRoot"`
	RefRange         string             `json:"refRange"`
	HeadCommit       string             `json:"headCommit"`
	LgtmVersion        string             `json:"lgtmVersion,omitempty"`
	IntentCandidates []interface{}      `json:"intentCandidates,omitempty"`
	HunkLinks        []matcher.HunkLink `json:"hunkLinks"`
	StruggleSignals  []interface{}      `json:"struggleSignals,omitempty"`
	HumanOverrides   []interface{}      `json:"humanOverrides,omitempty"`
	FileStats        interface{}        `json:"fileStats,omitempty"`
	ToolVersions     map[string]string  `json:"toolVersions,omitempty"`
}

// SessionRequest is the POST /v1/sessions body.
type SessionRequest struct {
	SessionId    string `json:"sessionId"`
	Tool         string `json:"tool"`
	Model        string `json:"model,omitempty"`
	Content      string `json:"content"`
	CapturedAtMs int64  `json:"capturedAtMs,omitempty"`
}
