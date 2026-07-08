package storage

type Session struct {
	ID               string
	CreatedAt        int64
	EndedAt          *int64
	Command          string
	Cwd              string
	ClientPID        *int
	ExitCode         *int
	GXVersion        string
	Source           *string
	ProcessName      *string
	ParentPID        *int
	LastSeenAt       *int64
	EndReason        *string
	RepoRoot         *string
	ModelsJSON       string
	InputTokens      *int
	OutputTokens     *int
	CacheReadTokens  *int
	CacheWriteTokens *int
}

type Request struct {
	ID             string
	SessionID      string
	CreatedAt      int64
	Provider       string
	Endpoint       string
	Method         string
	Model          *string
	RequestBody    []byte
	RequestHeaders string
}

type Response struct {
	ID                string
	RequestID         string
	CreatedAt         int64
	CompletedAt       int64
	StatusCode        int
	ResponseBody      []byte
	ResponseHeaders   string
	IsStreaming       bool
	DurationMS        int64
	ProviderRequestID *string
	FinishReason      *string
	InputTokens       *int
	OutputTokens      *int
	CacheReadTokens   *int
	CacheWriteTokens  *int
	Error             *string
}

type Repo struct {
	ID            int64
	RootPath      string
	Backend       string
	DefaultRemote *string
	DefaultBranch *string
	AuthoringBase *string
	RemoteURL     *string
	CreatedAt     int64
	UpdatedAt     int64
}

type InitializedRepo struct {
	RootPath  string
	CreatedAt int64
	UpdatedAt int64
}

type Change struct {
	ID              int64
	RepoID          int64
	JJChangeID      string
	CurrentCommitID string
	Description     string
	ParentChangeID  *string
	Status          string
	FirstSeenAt     int64
	UpdatedAt       int64
}

type Stack struct {
	ID           int64
	RepoID       int64
	Name         string
	BookmarkName string
	BaseRef      string
	BaseCommitID string
	HeadChangeID *string
	HeadCommitID *string
	RemoteName   *string
	RemoteRef    *string
	GitHubPRURL  *string
	Status       string
	CreatedAt    int64
	UpdatedAt    int64
}

type StackChange struct {
	ID       int64
	StackID  int64
	ChangeID int64
	// JJChangeID and CommitID mirror the referenced change row so a stack
	// entry carries both JJ logical identity and an exact current target.
	JJChangeID string
	CommitID   string
	Position   int
	CreatedAt  int64
}

type ChangeBookmark struct {
	ID                 int64
	ChangeID           int64
	BookmarkName       string
	RemoteName         *string
	RemoteRef          *string
	LastPushedCommitID string
	CreatedAt          int64
	UpdatedAt          int64
}

type ChangeRevision struct {
	ID            int64
	ChangeID      int64
	JJCommitID    string
	JJOperationID string
	ChangedFiles  string
	CreatedAt     int64
}

type ChangeSession struct {
	ID        int64
	ChangeID  int64
	SessionID string
	CreatedAt int64
}

type ChangeSessionProvenance struct {
	ID          int64
	ChangeID    int64
	SessionID   string
	AgentTool   string
	Provider    string
	ModelID     string
	Source      *string
	ProcessName *string
	CreatedAt   int64
}

type SessionContext struct {
	SessionID   string
	Tool        string
	Model       *string
	Format      string
	ContentJSON []byte
	CapturedAt  int64
}

type DemuxProposal struct {
	ID           string
	RepoID       int64
	BaseChangeID string
	Status       string
	PayloadJSON  string
	CreatedAt    int64
	UpdatedAt    int64
	AppliedAt    *int64
}

type ChangeDemuxEvidence struct {
	ID                 int64
	ChangeID           int64
	DemuxProposalID    string
	RevisionProposalID string
	Intent             string
	FilesJSON          string
	HunkIDsJSON        string
	UseHunks           bool
	Confidence         float64
	ProvenanceStatus   string
	EvidenceJSON       string
	CreatedAt          int64
}

type SemanticLabelSearchDocument struct {
	Ordinal      int
	Label        string
	Text         string
	Source       string
	Status       string
	EvidenceJSON string
}

type SemanticLabelSearchResult struct {
	Ordinal      int
	Label        string
	Source       string
	Status       string
	EvidenceJSON string
	RawScore     float64
}

type SemanticLabelWrite struct {
	Label SemanticLabel
	Link  SemanticLabelLink
}

type SemanticLabel struct {
	ID            int64
	RepoID        int64
	Label         string
	AliasesJSON   string
	SourcesJSON   string
	SeenCount     int
	AcceptedCount int
	Confidence    float64
	CreatedAt     int64
	UpdatedAt     int64
}

type SemanticLabelLink struct {
	ID                 int64
	RepoID             int64
	LabelID            int64
	DemuxProposalID    *string
	RevisionProposalID *string
	ChangeID           *int64
	StackBookmark      *string
	Source             string
	Score              float64
	Accepted           bool
	EvidenceJSON       string
	CreatedAt          int64
}

type ModifyEvent struct {
	ID                      int64
	RepoID                  int64
	TargetChangeID          int64
	PreviousCurrentChangeID *int64
	JJOperationID           string
	CreatedAt               int64
}

type CursorMessage struct {
	ID           string
	SessionID    string
	CreatedAt    int64
	Role         string
	Text         string
	RawJSON      []byte
	InputTokens  *int
	OutputTokens *int
}

type Push struct {
	ID              int64
	RepoID          int64
	RemoteName      *string
	BranchName      *string
	HeadCommitID    string
	CurrentChangeID *int64
	CreatedAt       int64
}

type AgentLedgerSummary struct {
	Agent      string
	Filepath   string
	Calls      int
	Tokens     int
	Files      int
	LastSeenAt *int64
}
