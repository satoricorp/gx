package reviewbundle

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/totality/internal/agentprovenance"
	"github.com/satoricorp/totality/internal/reviewsource"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/vcs"
	"github.com/satoricorp/totality/internal/version"
)

const SchemaVersion = 2

type Bundle struct {
	Event         string            `json:"event"`
	SchemaVersion int               `json:"schema_version"`
	CreatedAt     int64             `json:"created_at"`
	TLVersion     string            `json:"tl_version"`
	Repo          RepoPayload       `json:"repo"`
	Push          PushPayload       `json:"push"`
	Revisions     []RevisionPayload `json:"revisions"`
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

// RevisionPayload describes one pushed commit, oldest first in
// Bundle.Revisions. RevisionID is the Totality trailer ID from the commit message
// and is empty for commits without a trailer.
type RevisionPayload struct {
	RevisionID           string                `json:"revision_id"`
	CommitID             string                `json:"commit_id"`
	Description          string                `json:"description"`
	Files                []string              `json:"files"`
	BranchName           string                `json:"branch_name"`
	BaseBranchName       string                `json:"base_branch_name"`
	Patch                string                `json:"patch"`
	GitHubPullRequestURL *string               `json:"github_pull_request_url,omitempty"`
	ReviewContext        *ReviewContextPayload `json:"review_context,omitempty"`
}

// DemuxEvidencePayload never leaves the machine anymore: it feeds the
// review-context builder from local demux rows and stays off the wire.
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

type SessionPayload struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	EndedAt   *int64 `json:"ended_at,omitempty"`
	// Command is the agent that produced the transcript ("claude", "codex",
	// "cursor"). The console renders session.command and only falls back to
	// session.source, so dropping it rendered every session as `command=?`.
	Command          string           `json:"command"`
	Cwd              string           `json:"cwd"`
	TLVersion        string           `json:"tl_version"`
	Source           *string          `json:"source,omitempty"`
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
		Event:         "tl.pr",
		SchemaVersion: SchemaVersion,
		CreatedAt:     time.Now().UnixMilli(),
		TLVersion:     version.Current(),
		Repo: RepoPayload{
			RootPath: push.Repo.RootPath,
			// Pinned rather than copied: stored repo rows can still carry a
			// stale pre-Git backend value, and this bundle always describes a
			// git push.
			Backend:       "git",
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

	// Resolve the repos row once, by the same identity rule the write side
	// uses. Doing this per-commit by root_path alone was how a push could
	// write its change_sessions rows correctly and still publish sessions: [].
	repoID, repoFound, err := findRepoID(ctx, db, push.Repo.GitCommonDir, push.Repo.RootPath)
	if err != nil {
		return Bundle{}, err
	}

	branchName := pointerValue(bundle.Push.BranchName)
	baseBranchName := pointerValue(push.Repo.DefaultBranch)
	bundle.Revisions = make([]RevisionPayload, 0, len(push.Commits))
	seenSessions := map[string]bool{}
	for _, commit := range push.Commits {
		revision := RevisionPayload{
			RevisionID:           commit.RevisionID,
			CommitID:             commit.CommitID,
			Description:          commit.Message,
			Files:                commit.Files,
			BranchName:           branchName,
			BaseBranchName:       baseBranchName,
			Patch:                commit.Patch,
			GitHubPullRequestURL: push.GitHubPullRequestURL,
		}
		if commit.RevisionID != "" && repoFound {
			row, err := findRevisionRow(ctx, db, repoID, commit.RevisionID)
			if err != nil {
				return Bundle{}, err
			}
			if row != nil {
				if len(revision.Files) == 0 {
					revision.Files = row.files
				}
				revision.ReviewContext = row.reviewContext
				sessions, err := listChangeSessions(ctx, db, row.changeRowID)
				if err != nil {
					return Bundle{}, err
				}
				for _, session := range sessions {
					if seenSessions[session.ID] {
						continue
					}
					seenSessions[session.ID] = true
					bundle.Sessions = append(bundle.Sessions, session)
				}
			}
		}
		if revision.Files == nil {
			revision.Files = []string{}
		}
		bundle.Revisions = append(bundle.Revisions, revision)
	}
	if bundle.Sessions == nil {
		bundle.Sessions = []SessionPayload{}
	}
	bundle.Sessions = capSessionPayloads(bundle.Sessions)
	return bundle, nil
}

func pointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
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

// revisionRow is the local change row for a Totality revision ID: the DB id that
// links sessions, the recorded files, and the assembled review context. The
// local `changes` table keeps its historical column names; none of them reach
// the wire.
type revisionRow struct {
	changeRowID   int64
	files         []string
	reviewContext *ReviewContextPayload
}

// findRepoID resolves the repos row for a push exactly the way the write side
// does — storage.Store.UpsertRepo (internal/storage/writer.go) matches on
// git_common_dir first and falls back to root_path.
//
// Repository identity is the git common dir, not the worktree path: linked
// worktrees of one repository deliberately share a single repos row (see
// vcs.RepoIdentityKey), and UpsertRepo's existing-row UPDATE never rewrites
// root_path. So the stored root_path is whatever worktree happened to register
// the row first, and matching on it alone reads a different row than the one
// the change and its sessions were just written under — or no row at all. That
// silently published `sessions: []` for every push from a linked worktree and
// for any repository whose stored root_path had drifted.
func findRepoID(ctx context.Context, db *sql.DB, gitCommonDir, rootPath string) (int64, bool, error) {
	commonDir := strings.TrimSpace(gitCommonDir)
	root := strings.TrimSpace(rootPath)
	if commonDir == "" {
		commonDir = root
	}
	if commonDir == "" {
		return 0, false, nil
	}
	// A repos row can carry an empty root_path, so an absent root must never be
	// spelled as "" here or it would match that row by accident.
	rootMatch := root
	if rootMatch == "" {
		rootMatch = commonDir
	}
	var id int64
	err := db.QueryRowContext(ctx, `
		SELECT id
		FROM repos
		WHERE git_common_dir = ? OR root_path = ?
		ORDER BY CASE WHEN git_common_dir = ? THEN 0 ELSE 1 END, updated_at DESC
		LIMIT 1
	`, commonDir, rootMatch, commonDir).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("resolve pushed repo: %w", err)
	}
	return id, true, nil
}

func findRevisionRow(ctx context.Context, db *sql.DB, repoID int64, revisionID string) (*revisionRow, error) {
	row := db.QueryRowContext(ctx, `
		SELECT c.id, COALESCE(cr.changed_files_json, '[]')
		FROM changes c
		LEFT JOIN change_revisions cr ON cr.change_id = c.id
		WHERE c.repo_id = ? AND c.jj_change_id = ?
		ORDER BY cr.created_at DESC, cr.id DESC
		LIMIT 1
	`, repoID, revisionID)

	var result revisionRow
	var filesJSON string
	if err := row.Scan(&result.changeRowID, &filesJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find pushed revision: %w", err)
	}
	if err := json.Unmarshal([]byte(filesJSON), &result.files); err != nil {
		result.files = nil
	}
	evidence, err := listChangeDemuxEvidence(ctx, db, result.changeRowID)
	if err != nil {
		return nil, err
	}
	linkedSessionIDs, err := listChangeSessionIDs(ctx, db, result.changeRowID)
	if err != nil {
		return nil, err
	}
	agentProvenance, err := listChangeAgentProvenance(ctx, db, result.changeRowID)
	if err != nil {
		return nil, err
	}
	transcriptSources, err := listTranscriptSources(ctx, db, linkedSessionIDs)
	if err != nil {
		return nil, err
	}
	result.reviewContext = buildReviewContext(result.files, evidence, reviewsource.BuildGraph(evidenceStatuses(evidence), linkedSessionIDs, transcriptSources), agentProvenance)
	return &result, nil
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

func listChangeSessions(ctx context.Context, db *sql.DB, changeID int64) ([]SessionPayload, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT s.id, s.created_at, s.ended_at, s.command, s.cwd, s.tl_version,
			s.source, s.last_seen_at, s.end_reason, s.repo_root,
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
			&s.ID, &s.CreatedAt, &s.EndedAt, &s.Command, &s.Cwd, &s.TLVersion,
			&s.Source, &s.LastSeenAt, &s.EndReason, &s.RepoRoot,
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
	if ref := strings.TrimSpace(push.TotalityStackRef); ref != "" {
		return &ref
	}
	return push.Repo.BranchName
}
