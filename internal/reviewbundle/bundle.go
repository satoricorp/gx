package reviewbundle

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/agentprovenance"
	"github.com/satoricorp/gx/internal/reviewsource"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
	"github.com/satoricorp/gx/internal/version"
)

const SchemaVersion = 1

type Bundle struct {
	Event         string            `json:"event"`
	SchemaVersion int               `json:"schema_version"`
	CreatedAt     int64             `json:"created_at"`
	GXVersion     string            `json:"gx_version"`
	Repo          RepoPayload       `json:"repo"`
	Push          PushPayload       `json:"push"`
	Change        *ChangePayload    `json:"change,omitempty"`
	Stack         []StackPayload    `json:"stack,omitempty"`
	Sessions      []SessionPayload  `json:"sessions"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type Artifact struct {
	ReviewID    string `json:"review_id,omitempty"`
	ReviewURL   string `json:"review_url,omitempty"`
	IndexStatus string `json:"index_status,omitempty"`
	Bundle
}

func NewArtifact(bundle Bundle) Artifact {
	return Artifact{
		IndexStatus: "pending",
		Bundle:      bundle,
	}
}

type RepoPayload struct {
	RootPath      string  `json:"root_path"`
	Backend       string  `json:"backend"`
	DefaultRemote *string `json:"default_remote,omitempty"`
	DefaultBranch *string `json:"default_branch,omitempty"`
	RemoteURL     *string `json:"remote_url,omitempty"`
	BranchName    *string `json:"branch_name,omitempty"`
}

type PushPayload struct {
	RemoteName           *string `json:"remote_name,omitempty"`
	BranchName           *string `json:"branch_name,omitempty"`
	HeadCommitID         string  `json:"head_commit_id"`
	GitHubPullRequestURL *string `json:"github_pull_request_url,omitempty"`
}

type ChangePayload struct {
	ID              int64                  `json:"id"`
	JJChangeID      string                 `json:"jj_change_id"`
	CurrentCommitID string                 `json:"current_commit_id"`
	Description     string                 `json:"description"`
	ParentChangeID  *string                `json:"parent_change_id,omitempty"`
	Status          string                 `json:"status"`
	Files           []string               `json:"files"`
	DemuxEvidence   []DemuxEvidencePayload `json:"demux_evidence,omitempty"`
	ReviewContext   *ReviewContextPayload  `json:"review_context,omitempty"`
}

type DemuxEvidencePayload struct {
	ID                 int64           `json:"id"`
	DemuxProposalID    string          `json:"demux_proposal_id"`
	RevisionProposalID string          `json:"revision_proposal_id"`
	Intent             string          `json:"intent"`
	Files              []string        `json:"files"`
	HunkIDs            []string        `json:"hunk_ids"`
	UseHunks           bool            `json:"use_hunks"`
	Confidence         float64         `json:"confidence"`
	ProvenanceStatus   string          `json:"provenance_status"`
	Evidence           json.RawMessage `json:"evidence"`
	CreatedAt          int64           `json:"created_at"`
}

type ReviewContextPayload struct {
	ProvenanceStatus    string                     `json:"provenance_status,omitempty"`
	LinkedSessionCount  int                        `json:"linked_session_count,omitempty"`
	ProvenanceSources   []ReviewProvenanceSource   `json:"provenance_sources,omitempty"`
	AgentProvenance     []ReviewAgentProvenance    `json:"agent_provenance,omitempty"`
	TranscriptSources   []ReviewTranscriptSource   `json:"transcript_sources,omitempty"`
	StructuralStatus    string                     `json:"structural_status,omitempty"`
	StructuralFacts     []ReviewStructuralFact     `json:"structural_facts,omitempty"`
	StructuralDeps      []ReviewStructuralDep      `json:"structural_dependencies,omitempty"`
	ChangedSymbols      []ReviewChangedSymbol      `json:"changed_symbols,omitempty"`
	SemanticLabels      []ReviewSemanticLabel      `json:"semantic_labels,omitempty"`
	FeasibilityWarnings []ReviewFeasibilityWarning `json:"feasibility_warnings,omitempty"`
	Risk                RiskPayload                `json:"risk"`
	Evidence            []ReviewEvidencePayload    `json:"evidence,omitempty"`
}

type ReviewProvenanceSource = reviewsource.ProvenanceSource

type ReviewTranscriptSource = reviewsource.TranscriptSource

type ReviewAgentProvenance struct {
	SessionID   string  `json:"session_id"`
	AgentTool   string  `json:"agent_tool"`
	Provider    string  `json:"provider,omitempty"`
	ModelID     string  `json:"model_id,omitempty"`
	Source      *string `json:"source,omitempty"`
	ProcessName *string `json:"process_name,omitempty"`
	CreatedAt   int64   `json:"created_at"`
}

type ReviewEvidencePayload struct {
	Kind   string         `json:"kind"`
	Source string         `json:"source"`
	Status string         `json:"status,omitempty"`
	Data   map[string]any `json:"data,omitempty"`
}

type ReviewSemanticLabel struct {
	Label    string   `json:"label"`
	Source   string   `json:"source,omitempty"`
	Status   string   `json:"status,omitempty"`
	Score    float64  `json:"score,omitempty"`
	Evidence []string `json:"evidence,omitempty"`
}

type ReviewStructuralFact struct {
	File              string                   `json:"file"`
	Language          string                   `json:"language"`
	DefinedSymbols    []string                 `json:"defined_symbols,omitempty"`
	ReferencedSymbols []string                 `json:"referenced_symbols,omitempty"`
	Symbols           []ReviewStructuralSymbol `json:"symbols,omitempty"`
}

type ReviewStructuralSymbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type ReviewStructuralDep struct {
	FromFile string `json:"from_file"`
	ToFile   string `json:"to_file"`
	Symbol   string `json:"symbol"`
}

type ReviewChangedSymbol struct {
	HunkID    string `json:"hunk_id"`
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type ReviewFeasibilityWarning struct {
	RevisionID string `json:"revision_id,omitempty"`
	Severity   string `json:"severity"`
	Source     string `json:"source"`
	DependsOn  string `json:"depends_on,omitempty"`
	Symbol     string `json:"symbol,omitempty"`
	FromFile   string `json:"from_file,omitempty"`
	ToFile     string `json:"to_file,omitempty"`
	Message    string `json:"message"`
}

type RiskPayload struct {
	Level   string   `json:"level"`
	Score   int      `json:"score"`
	Signals []string `json:"signals,omitempty"`
}

type StackPayload struct {
	Change               ChangePayload `json:"change"`
	BranchName           string        `json:"branch_name"`
	BaseBranchName       string        `json:"base_branch_name"`
	Patch                string        `json:"patch"`
	GitHubPullRequestURL *string       `json:"github_pull_request_url,omitempty"`
}

type SessionPayload struct {
	ID               string           `json:"id"`
	CreatedAt        int64            `json:"created_at"`
	EndedAt          *int64           `json:"ended_at,omitempty"`
	Command          string           `json:"command"`
	Cwd              string           `json:"cwd"`
	ClientPID        *int             `json:"client_pid,omitempty"`
	ExitCode         *int             `json:"exit_code,omitempty"`
	GXVersion        string           `json:"gx_version"`
	Source           *string          `json:"source,omitempty"`
	ProcessName      *string          `json:"process_name,omitempty"`
	ParentPID        *int             `json:"parent_pid,omitempty"`
	LastSeenAt       *int64           `json:"last_seen_at,omitempty"`
	EndReason        *string          `json:"end_reason,omitempty"`
	RepoRoot         *string          `json:"repo_root,omitempty"`
	Models           []string         `json:"models,omitempty"`
	InputTokens      int              `json:"input_tokens,omitempty"`
	OutputTokens     int              `json:"output_tokens,omitempty"`
	CacheReadTokens  int              `json:"cache_read_tokens,omitempty"`
	CacheWriteTokens int              `json:"cache_write_tokens,omitempty"`
	Requests         []RequestPayload `json:"requests"`
}

type RequestPayload struct {
	ID             string            `json:"id"`
	SessionID      string            `json:"session_id"`
	CreatedAt      int64             `json:"created_at"`
	Provider       string            `json:"provider"`
	Endpoint       string            `json:"endpoint"`
	Method         string            `json:"method"`
	Model          *string           `json:"model,omitempty"`
	RequestBody    []byte            `json:"request_body"`
	RequestHeaders string            `json:"request_headers"`
	Responses      []ResponsePayload `json:"responses"`
}

type ResponsePayload struct {
	ID                string  `json:"id"`
	RequestID         string  `json:"request_id"`
	CreatedAt         int64   `json:"created_at"`
	CompletedAt       int64   `json:"completed_at"`
	StatusCode        int     `json:"status_code"`
	ResponseBody      []byte  `json:"response_body"`
	ResponseHeaders   string  `json:"response_headers"`
	IsStreaming       bool    `json:"is_streaming"`
	DurationMS        int64   `json:"duration_ms"`
	ProviderRequestID *string `json:"provider_request_id,omitempty"`
	FinishReason      *string `json:"finish_reason,omitempty"`
	InputTokens       *int    `json:"input_tokens,omitempty"`
	OutputTokens      *int    `json:"output_tokens,omitempty"`
	CacheReadTokens   *int    `json:"cache_read_tokens,omitempty"`
	CacheWriteTokens  *int    `json:"cache_write_tokens,omitempty"`
	Error             *string `json:"error,omitempty"`
}

func BuildPush(ctx context.Context, push vcs.PushResult) (Bundle, error) {
	db, err := storage.Open(ctx)
	if err != nil {
		return Bundle{}, err
	}
	defer db.Close()

	bundle := Bundle{
		Event:         "gx.pr",
		SchemaVersion: SchemaVersion,
		CreatedAt:     time.Now().UnixMilli(),
		GXVersion:     version.Current(),
		Repo: RepoPayload{
			RootPath:      push.Repo.RootPath,
			Backend:       push.Repo.Backend,
			DefaultRemote: push.Repo.DefaultRemote,
			DefaultBranch: push.Repo.DefaultBranch,
			RemoteURL:     push.Repo.RemoteURL,
			BranchName:    push.Repo.BranchName,
		},
		Push: PushPayload{
			RemoteName:           push.RemoteName,
			BranchName:           publishBranchName(push),
			HeadCommitID:         push.HeadCommitID,
			GitHubPullRequestURL: push.GitHubPullRequestURL,
		},
	}

	change, err := findPushedChange(ctx, db, push)
	if err != nil {
		return Bundle{}, err
	}
	if change != nil {
		bundle.Change = change
		sessions, err := listChangeSessions(ctx, db, change.ID)
		if err != nil {
			return Bundle{}, err
		}
		bundle.Sessions = sessions
	}
	if len(push.Published) > 0 {
		stack, err := listPushedStack(ctx, db, push)
		if err != nil {
			return Bundle{}, err
		}
		bundle.Stack = stack
		bundle.Sessions, err = listStackSessions(ctx, db, stack)
		if err != nil {
			return Bundle{}, err
		}
	}
	if bundle.Sessions == nil {
		bundle.Sessions = []SessionPayload{}
	}
	bundle.Sessions = capSessionPayloads(bundle.Sessions)
	return bundle, nil
}

const (
	// Bodies above this size are omitted from the bundle (token counts and
	// request metadata stay). Oversized bodies have produced multi-GB
	// artifacts that can never finish uploading.
	maxCapturedBodyBytes = 256 << 10
	// Total budget for all request/response bodies in one bundle. Once
	// spent, remaining bodies are omitted.
	maxSessionsPayloadBytes = 32 << 20
)

func capSessionPayloads(sessions []SessionPayload) []SessionPayload {
	budget := maxSessionsPayloadBytes
	for si := range sessions {
		for ri := range sessions[si].Requests {
			request := &sessions[si].Requests[ri]
			request.RequestBody, budget = capCapturedBody(request.RequestBody, budget)
			for pi := range request.Responses {
				response := &request.Responses[pi]
				response.ResponseBody, budget = capCapturedBody(response.ResponseBody, budget)
			}
		}
	}
	return sessions
}

func capCapturedBody(body []byte, budget int) ([]byte, int) {
	// Dropping a whole body keeps every retained body valid JSON; truncating
	// would leave unparseable fragments downstream.
	if len(body) > maxCapturedBodyBytes || len(body) > budget {
		return nil, budget
	}
	return body, budget - len(body)
}

func listPushedStack(ctx context.Context, db *sql.DB, push vcs.PushResult) ([]StackPayload, error) {
	stack := make([]StackPayload, 0, len(push.Published))
	for _, entry := range push.Published {
		change, err := findChangeByJJID(ctx, db, push.Repo.RootPath, entry.Change.ChangeID)
		if err != nil {
			return nil, err
		}
		if change == nil {
			continue
		}
		stack = append(stack, StackPayload{
			Change:               *change,
			BranchName:           entry.BranchName,
			BaseBranchName:       entry.BaseBranchName,
			Patch:                entry.Patch,
			GitHubPullRequestURL: entry.GitHubPullRequestURL,
		})
	}
	return stack, nil
}

func findPushedChange(ctx context.Context, db *sql.DB, push vcs.PushResult) (*ChangePayload, error) {
	if push.CurrentChange == nil {
		return nil, nil
	}
	return findChangeByJJID(ctx, db, push.Repo.RootPath, push.CurrentChange.ChangeID)
}

func findChangeByJJID(ctx context.Context, db *sql.DB, repoRoot, jjChangeID string) (*ChangePayload, error) {
	row := db.QueryRowContext(ctx, `
		SELECT c.id, c.jj_change_id, c.current_commit_id, c.description, c.parent_change_id, c.status,
			COALESCE(cr.changed_files_json, '[]')
		FROM repos r
		JOIN changes c ON c.repo_id = r.id
		LEFT JOIN change_revisions cr ON cr.change_id = c.id
		WHERE r.root_path = ? AND c.jj_change_id = ?
		ORDER BY cr.created_at DESC, cr.id DESC
		LIMIT 1
	`, repoRoot, jjChangeID)

	var change ChangePayload
	var parent sql.NullString
	var filesJSON string
	if err := row.Scan(&change.ID, &change.JJChangeID, &change.CurrentCommitID, &change.Description, &parent, &change.Status, &filesJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find pushed change: %w", err)
	}
	if parent.Valid {
		change.ParentChangeID = &parent.String
	}
	if err := json.Unmarshal([]byte(filesJSON), &change.Files); err != nil {
		change.Files = nil
	}
	evidence, err := listChangeDemuxEvidence(ctx, db, change.ID)
	if err != nil {
		return nil, err
	}
	change.DemuxEvidence = evidence
	linkedSessionIDs, err := listChangeSessionIDs(ctx, db, change.ID)
	if err != nil {
		return nil, err
	}
	agentProvenance, err := listChangeAgentProvenance(ctx, db, change.ID)
	if err != nil {
		return nil, err
	}
	transcriptSources, err := listTranscriptSources(ctx, db, linkedSessionIDs)
	if err != nil {
		return nil, err
	}
	change.ReviewContext = buildReviewContext(change.Files, evidence, reviewsource.BuildGraph(evidenceStatuses(evidence), linkedSessionIDs, transcriptSources), agentProvenance)
	return &change, nil
}

func listChangeSessionIDs(ctx context.Context, db *sql.DB, changeID int64) ([]string, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT session_id
		FROM change_sessions
		WHERE change_id = ?
		ORDER BY created_at ASC, session_id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("list change session ids: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var sessionID string
		if err := rows.Scan(&sessionID); err != nil {
			return nil, fmt.Errorf("scan change session id: %w", err)
		}
		out = append(out, sessionID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change session ids: %w", err)
	}
	return out, nil
}

func listTranscriptSources(ctx context.Context, db *sql.DB, sessionIDs []string) ([]reviewsource.TranscriptSource, error) {
	var out []reviewsource.TranscriptSource
	for _, sessionID := range sessionIDs {
		if sessionID == "" {
			continue
		}
		sources, err := listSessionTranscriptSources(ctx, db, sessionID)
		if err != nil {
			return nil, err
		}
		out = append(out, sources...)
	}
	return out, nil
}

func listChangeAgentProvenance(ctx context.Context, db *sql.DB, changeID int64) ([]ReviewAgentProvenance, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT session_id, agent_tool, provider, model_id, source, process_name, created_at
		FROM change_session_provenance
		WHERE change_id = ?
		ORDER BY created_at ASC, session_id ASC, agent_tool ASC, provider ASC, model_id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("list change agent provenance: %w", err)
	}
	defer rows.Close()

	var out []ReviewAgentProvenance
	for rows.Next() {
		var item ReviewAgentProvenance
		if err := rows.Scan(&item.SessionID, &item.AgentTool, &item.Provider, &item.ModelID, &item.Source, &item.ProcessName, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan change agent provenance: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change agent provenance: %w", err)
	}
	if len(out) > 0 {
		return out, nil
	}
	return deriveChangeAgentProvenance(ctx, db, changeID)
}

func deriveChangeAgentProvenance(ctx context.Context, db *sql.DB, changeID int64) ([]ReviewAgentProvenance, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT DISTINCT cs.session_id, s.command, s.source, s.process_name,
			COALESCE(r.provider, '') AS provider,
			COALESCE(r.model, '') AS model_id,
			cs.created_at
		FROM change_sessions cs
		JOIN sessions s ON s.id = cs.session_id
		LEFT JOIN requests r ON r.session_id = s.id
		WHERE cs.change_id = ?
		ORDER BY cs.created_at ASC, cs.session_id ASC, provider ASC, model_id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("derive change agent provenance: %w", err)
	}
	defer rows.Close()

	var out []ReviewAgentProvenance
	for rows.Next() {
		source, err := scanAgentProvenanceSource(rows)
		if err != nil {
			return nil, fmt.Errorf("scan derived agent provenance: %w", err)
		}
		out = append(out, reviewAgentProvenanceFromRecord(agentprovenance.Resolve(source)))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate derived agent provenance: %w", err)
	}
	return out, nil
}

func scanAgentProvenanceSource(rows *sql.Rows) (agentprovenance.Source, error) {
	var source agentprovenance.Source
	var sourceValue, processName sql.NullString
	if err := rows.Scan(&source.SessionID, &source.Command, &sourceValue, &processName, &source.Provider, &source.ModelID, &source.CreatedAt); err != nil {
		return agentprovenance.Source{}, err
	}
	if sourceValue.Valid {
		source.Source = &sourceValue.String
	}
	if processName.Valid {
		source.ProcessName = &processName.String
	}
	return source, nil
}

func reviewAgentProvenanceFromRecord(record agentprovenance.Record) ReviewAgentProvenance {
	return ReviewAgentProvenance{
		SessionID:   record.SessionID,
		AgentTool:   record.AgentTool,
		Provider:    record.Provider,
		ModelID:     record.ModelID,
		Source:      record.Source,
		ProcessName: record.ProcessName,
		CreatedAt:   record.CreatedAt,
	}
}

func listSessionTranscriptSources(ctx context.Context, db *sql.DB, sessionID string) ([]reviewsource.TranscriptSource, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT req.id, req.provider, req.model, req.created_at, resp.id
		FROM requests req
		LEFT JOIN responses resp ON resp.request_id = req.id
		WHERE req.session_id = ?
		ORDER BY req.created_at ASC, req.id ASC, resp.created_at ASC, resp.id ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list transcript sources: %w", err)
	}
	defer rows.Close()

	var out []reviewsource.TranscriptSource
	for rows.Next() {
		var source reviewsource.TranscriptSource
		var responseID sql.NullString
		source.SessionID = sessionID
		source.Source = "change_sessions"
		if err := rows.Scan(&source.RequestID, &source.Provider, &source.Model, &source.CreatedAt, &responseID); err != nil {
			return nil, fmt.Errorf("scan transcript source: %w", err)
		}
		if responseID.Valid {
			source.ResponseID = &responseID.String
		}
		out = append(out, source)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transcript sources: %w", err)
	}
	return out, nil
}

func listChangeDemuxEvidence(ctx context.Context, db *sql.DB, changeID int64) ([]DemuxEvidencePayload, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, demux_proposal_id, revision_proposal_id, intent, files_json, hunk_ids_json,
			use_hunks, confidence, provenance_status, evidence_json, created_at
		FROM change_demux_evidence
		WHERE change_id = ?
		ORDER BY id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("list change demux evidence: %w", err)
	}
	defer rows.Close()

	var out []DemuxEvidencePayload
	for rows.Next() {
		var evidence DemuxEvidencePayload
		var filesJSON string
		var hunkIDsJSON string
		var useHunks int
		var evidenceJSON string
		if err := rows.Scan(
			&evidence.ID,
			&evidence.DemuxProposalID,
			&evidence.RevisionProposalID,
			&evidence.Intent,
			&filesJSON,
			&hunkIDsJSON,
			&useHunks,
			&evidence.Confidence,
			&evidence.ProvenanceStatus,
			&evidenceJSON,
			&evidence.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change demux evidence: %w", err)
		}
		evidence.UseHunks = useHunks != 0
		if err := json.Unmarshal([]byte(filesJSON), &evidence.Files); err != nil {
			evidence.Files = nil
		}
		if err := json.Unmarshal([]byte(hunkIDsJSON), &evidence.HunkIDs); err != nil {
			evidence.HunkIDs = nil
		}
		if json.Valid([]byte(evidenceJSON)) {
			evidence.Evidence = json.RawMessage(evidenceJSON)
		}
		out = append(out, evidence)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change demux evidence: %w", err)
	}
	return out, nil
}

type demuxEvidenceContext struct {
	FeasibilityWarnings []ReviewFeasibilityWarning `json:"feasibility_warnings,omitempty"`
	StructuralFacts     []ReviewStructuralFact     `json:"structural_facts,omitempty"`
	StructuralDeps      []ReviewStructuralDep      `json:"structural_dependencies,omitempty"`
	ChangedSymbols      []ReviewChangedSymbol      `json:"changed_symbols,omitempty"`
	SemanticLabels      []ReviewSemanticLabel      `json:"semantic_labels,omitempty"`
}

func buildReviewContext(files []string, evidence []DemuxEvidencePayload, sourceGraph reviewsource.Graph, agentProvenance []ReviewAgentProvenance) *ReviewContextPayload {
	context := ReviewContextPayload{
		ProvenanceStatus:   sourceGraph.ProvenanceStatus,
		LinkedSessionCount: sourceGraph.LinkedSessionCount,
		ProvenanceSources:  sourceGraph.ProvenanceSources,
		AgentProvenance:    agentProvenance,
		TranscriptSources:  sourceGraph.TranscriptSources,
		StructuralStatus:   "unavailable",
	}
	seenWarnings := map[string]struct{}{}
	seenFacts := map[string]struct{}{}
	seenDeps := map[string]struct{}{}
	seenSymbols := map[string]struct{}{}
	seenLabels := map[string]struct{}{}
	for _, item := range evidence {
		var decoded demuxEvidenceContext
		if len(item.Evidence) > 0 && json.Valid(item.Evidence) {
			_ = json.Unmarshal(item.Evidence, &decoded)
		}
		for _, warning := range decoded.FeasibilityWarnings {
			key := warning.RevisionID + "\x00" + warning.Severity + "\x00" + warning.Source + "\x00" + warning.Message
			if _, ok := seenWarnings[key]; ok {
				continue
			}
			seenWarnings[key] = struct{}{}
			context.FeasibilityWarnings = append(context.FeasibilityWarnings, warning)
		}
		for _, fact := range decoded.StructuralFacts {
			key := fact.File
			if _, ok := seenFacts[key]; ok {
				continue
			}
			seenFacts[key] = struct{}{}
			context.StructuralFacts = append(context.StructuralFacts, fact)
		}
		for _, dep := range decoded.StructuralDeps {
			key := dep.FromFile + "\x00" + dep.ToFile + "\x00" + dep.Symbol
			if _, ok := seenDeps[key]; ok {
				continue
			}
			seenDeps[key] = struct{}{}
			context.StructuralDeps = append(context.StructuralDeps, dep)
		}
		for _, symbol := range decoded.ChangedSymbols {
			key := symbol.HunkID + "\x00" + symbol.File + "\x00" + symbol.Symbol
			if _, ok := seenSymbols[key]; ok {
				continue
			}
			seenSymbols[key] = struct{}{}
			context.ChangedSymbols = append(context.ChangedSymbols, symbol)
		}
		for _, label := range decoded.SemanticLabels {
			key := label.Label + "\x00" + label.Source + "\x00" + label.Status
			if _, ok := seenLabels[key]; ok {
				continue
			}
			seenLabels[key] = struct{}{}
			context.SemanticLabels = append(context.SemanticLabels, label)
		}
	}
	if len(context.StructuralFacts) > 0 || len(context.StructuralDeps) > 0 || len(context.ChangedSymbols) > 0 {
		context.StructuralStatus = "available"
	}
	context.Risk = computeRisk(files, evidence, context)
	context.Evidence = evidenceForContext(context)
	return &context
}

func evidenceStatuses(evidence []DemuxEvidencePayload) []string {
	out := make([]string, 0, len(evidence))
	for _, item := range evidence {
		out = append(out, item.ProvenanceStatus)
	}
	return out
}

func evidenceForContext(context ReviewContextPayload) []ReviewEvidencePayload {
	return []ReviewEvidencePayload{
		{
			Kind:   "provenance",
			Source: "change_sessions",
			Status: context.ProvenanceStatus,
			Data: map[string]any{
				"linked_session_count":    context.LinkedSessionCount,
				"transcript_source_count": len(context.TranscriptSources),
			},
		},
		{
			Kind:   "structural",
			Source: "local_static_analysis",
			Status: context.StructuralStatus,
			Data: map[string]any{
				"structural_fact_count":       len(context.StructuralFacts),
				"structural_dependency_count": len(context.StructuralDeps),
				"changed_symbol_count":        len(context.ChangedSymbols),
				"feasibility_warning_count":   len(context.FeasibilityWarnings),
			},
		},
		{
			Kind:   "risk",
			Source: "reviewbundle.computeRisk",
			Status: context.Risk.Level,
			Data: map[string]any{
				"score":   context.Risk.Score,
				"signals": context.Risk.Signals,
			},
		},
		{
			Kind:   "semantic_labels",
			Source: "local_bm25",
			Status: "available",
			Data: map[string]any{
				"label_count": len(context.SemanticLabels),
			},
		},
	}
}

// provenanceRisk scores one provenance status. Every status reviewsource can
// report is handled here, so a status never passes through without a signal:
// the weaker the link between the change and a captured session, the less a
// reviewer can lean on the recorded context.
func provenanceRisk(status string) (int, string) {
	switch status {
	case reviewsource.StatusAbsent:
		return 20, "missing_provenance"
	case reviewsource.StatusRepoLocal:
		return 10, "repo_local_provenance"
	case reviewsource.StatusUnknown:
		return 10, "unknown_provenance"
	case reviewsource.StatusLinked:
		return 0, "linked_session_provenance"
	case reviewsource.StatusExplicit:
		return 0, "explicit_session_provenance"
	case "":
		return 0, ""
	default:
		// An unrecognized status is itself a reason not to trust provenance.
		return 10, "unknown_provenance"
	}
}

func computeRisk(files []string, evidence []DemuxEvidencePayload, context ReviewContextPayload) RiskPayload {
	score := 0
	var signals []string
	if len(files) >= 5 {
		score += 20
		signals = append(signals, "many_files")
	}
	totalHunks := 0
	for _, item := range evidence {
		totalHunks += len(item.HunkIDs)
		points, signal := provenanceRisk(item.ProvenanceStatus)
		score += points
		if signal != "" {
			signals = append(signals, signal)
		}
		if item.Confidence > 0 && item.Confidence < 0.6 {
			score += 15
			signals = append(signals, "low_demux_confidence")
		}
	}
	if len(evidence) == 0 {
		points, signal := provenanceRisk(context.ProvenanceStatus)
		score += points
		if signal != "" {
			signals = append(signals, signal)
		}
	}
	if totalHunks >= 5 {
		score += 15
		signals = append(signals, "many_hunks")
	}
	if len(context.StructuralDeps) > 0 {
		score += 10
		signals = append(signals, "structural_dependencies")
	}
	for _, warning := range context.FeasibilityWarnings {
		switch warning.Severity {
		case "warning":
			score += 50
			signals = append(signals, "warning:"+warning.Source)
		default:
			score += 5
			signals = append(signals, "info:"+warning.Source)
		}
	}
	if score > 100 {
		score = 100
	}
	level := "low"
	if score >= 60 {
		level = "high"
	} else if score >= 25 {
		level = "medium"
	}
	return RiskPayload{Level: level, Score: score, Signals: dedupeStrings(signals)}
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func listStackSessions(ctx context.Context, db *sql.DB, stack []StackPayload) ([]SessionPayload, error) {
	seen := map[string]bool{}
	var sessions []SessionPayload
	for _, entry := range stack {
		changeSessions, err := listChangeSessions(ctx, db, entry.Change.ID)
		if err != nil {
			return nil, err
		}
		for _, session := range changeSessions {
			if seen[session.ID] {
				continue
			}
			seen[session.ID] = true
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func listChangeSessions(ctx context.Context, db *sql.DB, changeID int64) ([]SessionPayload, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT s.id, s.created_at, s.ended_at, s.command, s.cwd, s.client_pid, s.exit_code, s.gx_version,
			s.source, s.process_name, s.parent_pid, s.last_seen_at, s.end_reason, s.repo_root,
			s.models_json, s.input_tokens, s.output_tokens, s.cache_read_tokens, s.cache_write_tokens
		FROM change_sessions cs
		JOIN sessions s ON s.id = cs.session_id
		WHERE cs.change_id = ?
		ORDER BY cs.created_at ASC, s.created_at ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("list change sessions: %w", err)
	}
	defer rows.Close()

	var sessions []SessionPayload
	for rows.Next() {
		var s SessionPayload
		var modelsJSON string
		if err := rows.Scan(
			&s.ID, &s.CreatedAt, &s.EndedAt, &s.Command, &s.Cwd, &s.ClientPID, &s.ExitCode, &s.GXVersion,
			&s.Source, &s.ProcessName, &s.ParentPID, &s.LastSeenAt, &s.EndReason, &s.RepoRoot,
			&modelsJSON, &s.InputTokens, &s.OutputTokens, &s.CacheReadTokens, &s.CacheWriteTokens,
		); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		_ = json.Unmarshal([]byte(modelsJSON), &s.Models)
		requests, err := listSessionRequests(ctx, db, s.ID)
		if err != nil {
			return nil, err
		}
		s.Requests = requests
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}
	return sessions, nil
}

func listSessionRequests(ctx context.Context, db *sql.DB, sessionID string) ([]RequestPayload, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, session_id, created_at, provider, endpoint, method, model, request_body, request_headers
		FROM requests
		WHERE session_id = ?
		ORDER BY created_at ASC, id ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list session requests: %w", err)
	}
	defer rows.Close()

	var requests []RequestPayload
	for rows.Next() {
		var req RequestPayload
		if err := rows.Scan(&req.ID, &req.SessionID, &req.CreatedAt, &req.Provider, &req.Endpoint, &req.Method, &req.Model, &req.RequestBody, &req.RequestHeaders); err != nil {
			return nil, fmt.Errorf("scan request: %w", err)
		}
		responses, err := listRequestResponses(ctx, db, req.ID)
		if err != nil {
			return nil, err
		}
		req.Responses = responses
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate requests: %w", err)
	}
	return requests, nil
}

func listRequestResponses(ctx context.Context, db *sql.DB, requestID string) ([]ResponsePayload, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, request_id, created_at, completed_at, status_code, response_body, response_headers,
			is_streaming, duration_ms, provider_request_id, finish_reason,
			input_tokens, output_tokens, cache_read_tokens, cache_write_tokens, error
		FROM responses
		WHERE request_id = ?
		ORDER BY created_at ASC, id ASC
	`, requestID)
	if err != nil {
		return nil, fmt.Errorf("list request responses: %w", err)
	}
	defer rows.Close()

	var responses []ResponsePayload
	for rows.Next() {
		var resp ResponsePayload
		if err := rows.Scan(
			&resp.ID, &resp.RequestID, &resp.CreatedAt, &resp.CompletedAt, &resp.StatusCode, &resp.ResponseBody, &resp.ResponseHeaders,
			&resp.IsStreaming, &resp.DurationMS, &resp.ProviderRequestID, &resp.FinishReason,
			&resp.InputTokens, &resp.OutputTokens, &resp.CacheReadTokens, &resp.CacheWriteTokens, &resp.Error,
		); err != nil {
			return nil, fmt.Errorf("scan response: %w", err)
		}
		responses = append(responses, resp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate responses: %w", err)
	}
	return responses, nil
}

func publishBranchName(push vcs.PushResult) *string {
	if ref := strings.TrimSpace(push.GXStackRef); ref != "" {
		return &ref
	}
	return push.Repo.BranchName
}
