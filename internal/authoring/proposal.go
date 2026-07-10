package authoring

import (
	"io"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
)

type ProposalStatus string

const (
	ProposalPending          ProposalStatus = "pending"
	ProposalPartiallyApplied ProposalStatus = "partially_applied"
	ProposalApplied          ProposalStatus = "applied"
	ProposalExpired          ProposalStatus = "expired"
)

type DemuxApplyCheckpoint struct {
	RevisionID string `json:"revision_id"`
	ChangeID   string `json:"change_id,omitempty"`
	CommitID   string `json:"commit_id,omitempty"`
}

type HunkRange struct {
	ID         string `json:"id"`
	File       string `json:"file"`
	Header     string `json:"header"`
	OldStart   int    `json:"old_start"`
	OldLines   int    `json:"old_lines"`
	NewStart   int    `json:"new_start"`
	NewLines   int    `json:"new_lines"`
	Symbol     string `json:"symbol,omitempty"`
	SymbolKind string `json:"symbol_kind,omitempty"`
	Patch      string `json:"patch,omitempty"`
}

type RevisionProposal struct {
	ID               string             `json:"id"`
	Intent           string             `json:"intent"`
	Files            []string           `json:"files"`
	UseHunks         bool               `json:"use_hunks,omitempty"`
	HunkIDs          []string           `json:"hunk_ids,omitempty"`
	Hunks            []HunkRange        `json:"hunks,omitempty"`
	DependsOn        []string           `json:"depends_on,omitempty"`
	TargetStack      string             `json:"target_stack,omitempty"`
	BaseStack        string             `json:"base_stack,omitempty"`
	RouteReason      string             `json:"route_reason,omitempty"`
	RouteConfidence  float64            `json:"route_confidence,omitempty"`
	RouteSource      string             `json:"route_source,omitempty"`
	ProvenanceStatus string             `json:"provenance_status"`
	SessionIDs       []string           `json:"session_ids,omitempty"`
	SessionContexts  []SessionContext   `json:"session_contexts,omitempty"`
	HunkLinks        []matcher.HunkLink `json:"hunk_links,omitempty"`
	Confidence       float64            `json:"confidence"`
	EffectiveLOC     int                `json:"effective_loc,omitempty"`
	ShapeReasons     []string           `json:"shape_reasons,omitempty"`
	SemanticLabels   []SemanticLabel    `json:"semantic_labels,omitempty"`
}

type SessionContext struct {
	SessionID       string                 `json:"session_id"`
	Tool            string                 `json:"tool"`
	Model           string                 `json:"model,omitempty"`
	Format          string                 `json:"format"`
	ContentRedacted []capture.SessionEvent `json:"content_redacted"`
	CapturedAt      int64                  `json:"captured_at"`
}

type SemanticLabel struct {
	Label    string   `json:"label"`
	Source   string   `json:"source"`
	Status   string   `json:"status,omitempty"`
	Score    float64  `json:"score"`
	Evidence []string `json:"evidence,omitempty"`
}

type StructuralFact struct {
	File              string             `json:"file"`
	Language          string             `json:"language"`
	DefinedSymbols    []string           `json:"defined_symbols,omitempty"`
	ReferencedSymbols []string           `json:"referenced_symbols,omitempty"`
	Imports           []string           `json:"imports,omitempty"`
	Symbols           []StructuralSymbol `json:"symbols,omitempty"`
}

type StructuralSymbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type StructuralDependency struct {
	FromFile string `json:"from_file"`
	ToFile   string `json:"to_file"`
	Symbol   string `json:"symbol"`
}

type ChangedSymbol struct {
	HunkID    string `json:"hunk_id"`
	File      string `json:"file"`
	Symbol    string `json:"symbol"`
	Kind      string `json:"kind"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type FeasibilityWarning struct {
	RevisionID string `json:"revision_id,omitempty"`
	Severity   string `json:"severity"`
	Source     string `json:"source"`
	DependsOn  string `json:"depends_on,omitempty"`
	Symbol     string `json:"symbol,omitempty"`
	FromFile   string `json:"from_file,omitempty"`
	ToFile     string `json:"to_file,omitempty"`
	Message    string `json:"message"`
}

type PlanConfidence struct {
	LogicConfidence     float64            `json:"logic_confidence,omitempty"`
	LLMConfidence       float64            `json:"llm_confidence,omitempty"`
	EffectiveConfidence float64            `json:"effective_confidence,omitempty"`
	LogicReasons        []ConfidenceReason `json:"logic_reasons,omitempty"`
	LLMReasons          []ConfidenceReason `json:"llm_reasons,omitempty"`
}

type ConfidenceReason struct {
	Kind       string  `json:"kind"`
	Severity   string  `json:"severity"`
	Message    string  `json:"message"`
	Suggestion string  `json:"suggestion,omitempty"`
	Delta      float64 `json:"delta,omitempty"`
}

type DemuxProposal struct {
	ID                          string                 `json:"id"`
	RepoRoot                    string                 `json:"repo_root"`
	ProposedChangeID            string                 `json:"proposed_change_id"`
	ProposedCommitID            string                 `json:"proposed_commit_id"`
	Status                      ProposalStatus         `json:"status"`
	PlanInstructions            []string               `json:"plan_instructions,omitempty"`
	Hunks                       []HunkRange            `json:"hunks,omitempty"`
	StructuralFacts             []StructuralFact       `json:"structural_facts,omitempty"`
	StructuralDeps              []StructuralDependency `json:"structural_dependencies,omitempty"`
	ChangedSymbols              []ChangedSymbol        `json:"changed_symbols,omitempty"`
	Revisions                   []RevisionProposal     `json:"revisions"`
	FeasibilityWarnings         []FeasibilityWarning   `json:"feasibility_warnings,omitempty"`
	Warnings                    []string               `json:"warnings,omitempty"`
	HunkLinks                   []matcher.HunkLink     `json:"hunk_links,omitempty"`
	HunkCoverage                float64                `json:"hunk_coverage,omitempty"`
	CaptureTools                []string               `json:"capture_tools,omitempty"`
	Confidence                  PlanConfidence         `json:"confidence_summary,omitempty"`
	AppliedRevisionIDs          []string               `json:"applied_revision_ids,omitempty"`
	EvidenceRecordedRevisionIDs []string               `json:"evidence_recorded_revision_ids,omitempty"`
	ApplyCheckpoints            []DemuxApplyCheckpoint `json:"apply_checkpoints,omitempty"`
	CreatedAt                   int64                  `json:"created_at"`
}

type DemuxRevisionView struct {
	ProposalID          string                 `json:"proposal_id"`
	ProposalStatus      ProposalStatus         `json:"proposal_status"`
	ProposedChangeID    string                 `json:"proposed_change_id"`
	ProposedCommitID    string                 `json:"proposed_commit_id"`
	Revision            RevisionProposal       `json:"revision"`
	Hunks               []HunkRange            `json:"hunks,omitempty"`
	FeasibilityWarnings []FeasibilityWarning   `json:"feasibility_warnings,omitempty"`
	StructuralFacts     []StructuralFact       `json:"structural_facts,omitempty"`
	StructuralDeps      []StructuralDependency `json:"structural_dependencies,omitempty"`
	ChangedSymbols      []ChangedSymbol        `json:"changed_symbols,omitempty"`
}

type DemuxProposalSummary struct {
	ID                      string         `json:"id"`
	Status                  ProposalStatus `json:"status"`
	ProposedChangeID        string         `json:"proposed_change_id"`
	ProposedCommitID        string         `json:"proposed_commit_id,omitempty"`
	CreatedAt               int64          `json:"created_at"`
	UpdatedAt               int64          `json:"updated_at"`
	RevisionCount           int            `json:"revision_count"`
	Files                   []string       `json:"files,omitempty"`
	FeasibilityWarningCount int            `json:"feasibility_warning_count"`
	HiddenDiagnosticCount   int            `json:"hidden_diagnostic_count"`
	FirstRevisionIntent     string         `json:"first_revision_intent,omitempty"`
	LatestPendingForShow    bool           `json:"latest_pending_for_show,omitempty"`
	Alias                   string         `json:"alias,omitempty"`
}

type ListDemuxProposalsOptions struct {
	IncludeApplied bool
	Limit          int
}

type ProposeDemuxOptions struct {
	Intent                 string
	Filesets               []string
	ExcludeFilesets        []string
	PlanOnly               bool
	Legacy                 bool
	Model                  string
	MaxWarnings            int
	ApplyPreflightAttempts int
	ProgressWriter         io.Writer
}

type ApplyDemuxResult struct {
	Proposal         DemuxProposal      `json:"proposal"`
	Revisions        []CheckpointResult `json:"revisions"`
	RemainingChanges bool               `json:"remaining_changes,omitempty"`
	AcceptedSubset   bool               `json:"accepted_subset,omitempty"`
	NextAction       string             `json:"next_action,omitempty"`
}

type ApplyDemuxOptions struct {
	AllowWarnings         bool
	ReturnToDefaultBranch bool
}

type DemuxWorkflowState string

const (
	DemuxWorkflowReadyToApply      DemuxWorkflowState = "ready_to_apply"
	DemuxWorkflowRepairRecommended DemuxWorkflowState = "repair_recommended"
	DemuxWorkflowRepairRequired    DemuxWorkflowState = "repair_required"
)

type DemuxPlanPacket struct {
	Action              string                 `json:"action"`
	State               DemuxWorkflowState     `json:"state"`
	NextTool            string                 `json:"next_tool"`
	FinalTool           string                 `json:"final_tool"`
	Workflow            []string               `json:"workflow"`
	PlanningContract    []string               `json:"planning_contract"`
	RevisionShape       map[string]any         `json:"revision_shape"`
	ReviewArgumentShape map[string]any         `json:"review_argument_shape"`
	ApplyArgumentShape  map[string]any         `json:"apply_argument_shape"`
	Proposal            DemuxProposal          `json:"proposal"`
	Review              ReviewDemuxResult      `json:"review"`
	StructuralFacts     []StructuralFact       `json:"structural_facts,omitempty"`
	StructuralDeps      []StructuralDependency `json:"structural_dependencies,omitempty"`
	ChangedSymbols      []ChangedSymbol        `json:"changed_symbols,omitempty"`
	FeasibilityWarnings []FeasibilityWarning   `json:"feasibility_warnings,omitempty"`
	Warnings            []string               `json:"warnings,omitempty"`
}

type ReviewDemuxResult struct {
	Valid       bool          `json:"valid"`
	Proposal    DemuxProposal `json:"proposal"`
	Errors      []string      `json:"errors,omitempty"`
	RepairHints []RepairHint  `json:"repair_hints,omitempty"`
}

type demuxAIProposalForReview struct {
	ID               string                 `json:"id"`
	RepoRoot         string                 `json:"repo_root"`
	ProposedChangeID string                 `json:"proposed_change_id"`
	ProposedCommitID string                 `json:"proposed_commit_id"`
	Status           ProposalStatus         `json:"status"`
	PlanInstructions []string               `json:"plan_instructions,omitempty"`
	Confidence       PlanConfidence         `json:"confidence_summary,omitempty"`
	Hunks            []HunkRange            `json:"hunks,omitempty"`
	StructuralDeps   []StructuralDependency `json:"structural_dependencies,omitempty"`
	ChangedSymbols   []ChangedSymbol        `json:"changed_symbols,omitempty"`
	Revisions        []RevisionProposal     `json:"revisions"`
}

type demuxAIReviewSummary struct {
	Valid               bool                 `json:"valid"`
	Errors              []string             `json:"errors,omitempty"`
	FeasibilityWarnings []FeasibilityWarning `json:"feasibility_warnings,omitempty"`
	RepairHints         []RepairHint         `json:"repair_hints,omitempty"`
}

type DemuxAIReviewOptions struct {
	Model       string
	MaxWarnings int
	PlanOnly    bool
}

type DemuxAIReviewResult struct {
	Proposal        DemuxProposal      `json:"proposal"`
	Review          ReviewDemuxResult  `json:"review"`
	State           DemuxWorkflowState `json:"state"`
	Model           string             `json:"model"`
	Updated         bool               `json:"updated"`
	Notes           []string           `json:"notes,omitempty"`
	WarningsSent    int                `json:"warnings_sent"`
	WarningsTotal   int                `json:"warnings_total"`
	RepairHintSent  int                `json:"repair_hints_sent"`
	RepairHintTotal int                `json:"repair_hints_total"`
}

type RepairHint struct {
	Kind        string   `json:"kind"`
	RevisionID  string   `json:"revision_id,omitempty"`
	HunkID      string   `json:"hunk_id,omitempty"`
	File        string   `json:"file,omitempty"`
	DependsOn   string   `json:"depends_on,omitempty"`
	Symbol      string   `json:"symbol,omitempty"`
	TargetStack string   `json:"target_stack,omitempty"`
	BaseStack   string   `json:"base_stack,omitempty"`
	Candidates  []string `json:"candidates,omitempty"`
	Suggestion  string   `json:"suggestion"`
}
