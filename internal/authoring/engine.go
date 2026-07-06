package authoring

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

type RepoInfo = vcs.RepoInfo
type ChangeInfo = vcs.ChangeInfo
type InitOptions = vcs.InitOptions
type InitResult = vcs.InitResult
type CheckpointResult = vcs.CommitResult
type ModifyResult = vcs.ModifyResult
type SwitchResult = vcs.SwitchResult
type CreateStackResult = vcs.CreateStackResult
type BaseResult = vcs.BaseResult
type DeleteRevisionResult = vcs.DeleteRevisionResult
type DeleteStackResult = vcs.DeleteStackResult
type PruneEmptyStacksResult = vcs.PruneEmptyStacksResult
type PruneGitHubPullRequestsResult = vcs.PruneGitHubPullRequestsResult
type StackInfo = vcs.StackInfo
type PushOptions = vcs.PushOptions
type PublishMode = vcs.PublishMode
type PushResult = vcs.PushResult
type PublishHook = vcs.PublishHook
type SyncResult = vcs.SyncResult
type RevisionSummary = vcs.RevisionSummary
type StackSummary = vcs.StackSummary
type StatusSnapshot = vcs.StatusSnapshot

const (
	PublishModeReviewOnly   = vcs.PublishModeReviewOnly
	PublishModeReviewAndGit = vcs.PublishModeReviewAndGit
)

type CheckpointOptions struct {
	Intent                 string
	Filesets               []string
	Interactive            bool
	Hunk                   bool
	PatchFile              string
	PreferredSessionIDs      []string
	BookmarkRecordedCommit bool
}

// Engine is GX's authoring seam. CLI and MCP adapters should call this module
// instead of owning JJ/Git/storage mechanics directly.
type Engine struct {
	vcs *vcs.Service
}

func NewEngine() *Engine {
	return NewEngineWithVCS(vcs.NewService())
}

func NewEngineWithVCS(service *vcs.Service) *Engine {
	return &Engine{vcs: service}
}

func (e *Engine) SetProgressWriter(out io.Writer) {
	e.vcs.SetProgressWriter(out)
}

func (e *Engine) Init(ctx context.Context, opts InitOptions) (InitResult, error) {
	return e.vcs.InitWithOptions(ctx, opts)
}

func (e *Engine) Checkpoint(ctx context.Context, opts CheckpointOptions) (CheckpointResult, error) {
	if err := vcs.ValidateCommitMessage(opts.Intent); err != nil {
		return CheckpointResult{}, err
	}
	if opts.Hunk && opts.Interactive {
		return CheckpointResult{}, fmt.Errorf("choose one mode: --hunk or --interactive")
	}
	if opts.Hunk && len(opts.Filesets) > 0 {
		return CheckpointResult{}, fmt.Errorf("filesets are not supported with --hunk; use --patch-file")
	}
	return e.vcs.RecordAuthoringRevision(ctx, vcs.RevisionOptions{
		Message:                opts.Intent,
		Filesets:               opts.Filesets,
		Interactive:            opts.Interactive,
		Hunk:                   opts.Hunk,
		PatchFile:              opts.PatchFile,
		PreferredSessionIDs:    opts.PreferredSessionIDs,
		BookmarkRecordedCommit: opts.BookmarkRecordedCommit,
	})
}

func (e *Engine) Modify(ctx context.Context, rev string) (ModifyResult, error) {
	return e.vcs.RecordModification(ctx, rev)
}

func (e *Engine) ModifyCandidates(ctx context.Context, limit int) ([]ChangeInfo, error) {
	return e.vcs.ListModifyCandidates(ctx, limit)
}

func (e *Engine) BaseSwitchTarget(ctx context.Context, name string) (string, bool, error) {
	return e.vcs.BaseSwitchTarget(ctx, name)
}

func (e *Engine) Base(ctx context.Context) (BaseResult, error) {
	return e.vcs.Base(ctx)
}

func (e *Engine) SetBase(ctx context.Context, baseRef string) (BaseResult, error) {
	return e.vcs.SetBase(ctx, baseRef)
}

func (e *Engine) RequireAuthoringBase(ctx context.Context, commandName string) error {
	return e.vcs.RequireAuthoringBase(ctx, commandName)
}

func (e *Engine) Switch(ctx context.Context, name string) (SwitchResult, error) {
	return e.vcs.Switch(ctx, name)
}

func (e *Engine) CreateStack(ctx context.Context, name string) (CreateStackResult, error) {
	return e.vcs.CreateStack(ctx, name)
}

func (e *Engine) DeleteRevision(ctx context.Context, rev string) (DeleteRevisionResult, error) {
	return e.vcs.DeleteRevision(ctx, rev)
}

func (e *Engine) DeleteStack(ctx context.Context, bookmarkName string) (DeleteStackResult, error) {
	return e.vcs.DeleteStack(ctx, bookmarkName)
}

func (e *Engine) PruneEmptyStacks(ctx context.Context) (PruneEmptyStacksResult, error) {
	return e.vcs.PruneEmptyStacks(ctx)
}

func (e *Engine) Status(ctx context.Context) (StackSummary, error) {
	return e.vcs.Stack(ctx)
}

func (e *Engine) DetectMissingStackBaseRefs(ctx context.Context) (vcs.MissingStackBaseRefStatus, error) {
	return e.vcs.DetectMissingStackBaseRefs(ctx)
}

func (e *Engine) RebaseMissingStackBaseRefs(ctx context.Context, issues []vcs.MissingStackBaseRef) (vcs.RebaseOntoDefaultResult, error) {
	return e.vcs.RebaseMissingStackBaseRefs(ctx, issues)
}

func (e *Engine) RepairMissingStackBaseRefs(ctx context.Context) (vcs.RebaseOntoDefaultResult, error) {
	return e.vcs.RepairMissingStackBaseRefs(ctx)
}

func (e *Engine) ResolveJJRepo(ctx context.Context) (RepoInfo, error) {
	return e.vcs.ResolveJJRepo(ctx)
}

func (e *Engine) StatusSnapshot(ctx context.Context) (StatusSnapshot, error) {
	return e.vcs.StatusSnapshot(ctx)
}

func (e *Engine) CurrentChange(ctx context.Context, repoRoot, rev string) (ChangeInfo, error) {
	return e.vcs.CurrentChange(ctx, repoRoot, rev)
}

func (e *Engine) DiffGit(ctx context.Context, repoRoot, rev string) (string, error) {
	return e.vcs.DiffGit(ctx, repoRoot, rev)
}

func (e *Engine) PreparePublish(ctx context.Context, args []string, opts PushOptions) (PushResult, error) {
	return e.vcs.PreparePublish(ctx, args, opts)
}

func (e *Engine) PrepareNamedPublish(ctx context.Context, name string, args []string, opts PushOptions) (PushResult, error) {
	return e.vcs.PrepareNamedPublish(ctx, name, args, opts)
}

func (e *Engine) PrepareAllPublishes(ctx context.Context, args []string, opts PushOptions) ([]PushResult, error) {
	return e.vcs.PrepareAllPublishes(ctx, args, opts)
}

func (e *Engine) RecordPublish(ctx context.Context, result PushResult) error {
	return e.vcs.RecordPublish(ctx, result)
}

func (e *Engine) Publish(ctx context.Context, args []string, opts PushOptions, hook PublishHook) (PushResult, error) {
	return e.vcs.PublishStack(ctx, args, opts, hook)
}

func (e *Engine) PublishNamed(ctx context.Context, name string, args []string, opts PushOptions, hook PublishHook) (PushResult, error) {
	return e.vcs.PublishNamedStack(ctx, name, args, opts, hook)
}

func (e *Engine) PublishAll(ctx context.Context, args []string, opts PushOptions, hook PublishHook) ([]PushResult, error) {
	return e.vcs.PublishAllStacks(ctx, args, opts, hook)
}

func (e *Engine) Sync(ctx context.Context, remote string) (SyncResult, error) {
	return e.vcs.Sync(ctx, remote)
}

func (e *Engine) SyncCloudBookmarkTip(ctx context.Context, repo RepoInfo, branchName string) error {
	return e.vcs.SyncCloudBookmarkTip(ctx, repo, branchName)
}

func (e *Engine) PrunePublishedStackByRef(ctx context.Context, repo RepoInfo, publishRef string) (bool, error) {
	return e.vcs.PrunePublishedStackByRef(ctx, repo, publishRef)
}

func (e *Engine) PruneTerminalGitHubPullRequestStacks(ctx context.Context, repo RepoInfo) (PruneGitHubPullRequestsResult, error) {
	return e.vcs.PruneTerminalGitHubPullRequestStacks(ctx, repo)
}

func (e *Engine) SaveDemuxProposal(ctx context.Context, proposal DemuxProposal) (DemuxProposal, error) {
	now := time.Now().UnixMilli()
	if strings.TrimSpace(proposal.ID) == "" {
		proposal.ID = newProposalID()
	}
	if proposal.Status == "" {
		proposal.Status = ProposalPending
	}
	if proposal.CreatedAt == 0 {
		proposal.CreatedAt = now
	}
	repo, err := e.repoForProposal(ctx, proposal.RepoRoot)
	if err != nil {
		return DemuxProposal{}, err
	}
	proposal.RepoRoot = repo.RootPath
	if strings.TrimSpace(proposal.ProposedChangeID) == "" {
		if current, err := e.vcs.CurrentChange(ctx, repo.RootPath, "@"); err == nil {
			proposal.ProposedChangeID = current.ChangeID
			proposal.ProposedCommitID = current.CommitID
		}
	}

	payload, err := json.Marshal(proposal)
	if err != nil {
		return DemuxProposal{}, fmt.Errorf("marshal demux proposal: %w", err)
	}
	store, err := openStore(ctx)
	if err != nil {
		return DemuxProposal{}, err
	}
	defer store.Close()
	repoID, err := upsertProposalRepo(ctx, store, repo, now)
	if err != nil {
		return DemuxProposal{}, err
	}
	if err := store.UpsertDemuxProposal(ctx, storage.DemuxProposal{
		ID:           proposal.ID,
		RepoID:       repoID,
		BaseChangeID: proposal.ProposedChangeID,
		Status:       string(proposal.Status),
		PayloadJSON:  string(payload),
		CreatedAt:    proposal.CreatedAt,
		UpdatedAt:    now,
	}); err != nil {
		return DemuxProposal{}, err
	}
	return proposal, nil
}

func (e *Engine) LoadDemuxProposal(ctx context.Context, id string) (DemuxProposal, error) {
	store, err := openStore(ctx)
	if err != nil {
		return DemuxProposal{}, err
	}
	defer store.Close()
	row, err := store.FindDemuxProposal(ctx, id)
	if err != nil {
		return DemuxProposal{}, err
	}
	if row == nil {
		return DemuxProposal{}, fmt.Errorf("unknown demux proposal %q", id)
	}
	var proposal DemuxProposal
	if err := json.Unmarshal([]byte(row.PayloadJSON), &proposal); err != nil {
		return DemuxProposal{}, fmt.Errorf("decode demux proposal %q: %w", id, err)
	}
	proposal.ID = row.ID
	proposal.ProposedChangeID = row.BaseChangeID
	proposal.Status = ProposalStatus(row.Status)
	proposal.CreatedAt = row.CreatedAt
	return proposal, nil
}

func (e *Engine) LoadDemuxProposalBySelector(ctx context.Context, selector string) (DemuxProposal, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return DemuxProposal{}, fmt.Errorf("demux proposal selector is required")
	}
	if index, ok := parseDemuxProposalAlias(selector); ok {
		proposals, err := e.listStoredDemuxProposals(ctx, ListDemuxProposalsOptions{Limit: index})
		if err != nil {
			return DemuxProposal{}, err
		}
		if len(proposals) < index {
			return DemuxProposal{}, fmt.Errorf("unknown demux proposal alias %q", selector)
		}
		return decodeStoredDemuxProposal(proposals[index-1])
	}
	return e.LoadDemuxProposal(ctx, selector)
}

func (e *Engine) LoadLatestPendingDemuxProposal(ctx context.Context) (DemuxProposal, error) {
	repo, err := e.vcs.ResolveJJRepo(ctx)
	if err != nil {
		return DemuxProposal{}, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return DemuxProposal{}, err
	}
	defer store.Close()
	rowRepo, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return DemuxProposal{}, err
	}
	if rowRepo == nil {
		return DemuxProposal{}, fmt.Errorf("no demux proposal found for current repo")
	}
	row, err := store.FindLatestDemuxProposal(ctx, rowRepo.ID, string(ProposalPending))
	if err != nil {
		return DemuxProposal{}, err
	}
	if row == nil {
		return DemuxProposal{}, fmt.Errorf("no pending demux proposal found for current repo")
	}
	return decodeStoredDemuxProposal(*row)
}

func (e *Engine) ListDemuxProposals(ctx context.Context, opts ListDemuxProposalsOptions) ([]DemuxProposalSummary, error) {
	rows, err := e.listStoredDemuxProposals(ctx, opts)
	if err != nil {
		return nil, err
	}
	summaries := make([]DemuxProposalSummary, 0, len(rows))
	for index, row := range rows {
		proposal, err := decodeStoredDemuxProposal(row)
		if err != nil {
			return nil, err
		}
		summary := summarizeDemuxProposal(proposal, row.UpdatedAt)
		summary.LatestPendingForShow = !opts.IncludeApplied && index == 0 && summary.Status == ProposalPending
		if !opts.IncludeApplied {
			summary.Alias = fmt.Sprintf("d%d", index+1)
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (e *Engine) listStoredDemuxProposals(ctx context.Context, opts ListDemuxProposalsOptions) ([]storage.DemuxProposal, error) {
	repo, err := e.vcs.ResolveJJRepo(ctx)
	if err != nil {
		return nil, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	rowRepo, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return nil, err
	}
	if rowRepo == nil {
		return nil, nil
	}
	status := string(ProposalPending)
	if opts.IncludeApplied {
		status = ""
	}
	return store.ListDemuxProposals(ctx, rowRepo.ID, status, opts.Limit)
}

func (e *Engine) MarkDemuxProposalStatus(ctx context.Context, id string, status ProposalStatus) error {
	now := time.Now().UnixMilli()
	var appliedAt *int64
	if status == ProposalApplied {
		appliedAt = &now
	}
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()
	return store.UpdateDemuxProposalStatus(ctx, id, string(status), now, appliedAt)
}

func decodeStoredDemuxProposal(row storage.DemuxProposal) (DemuxProposal, error) {
	var proposal DemuxProposal
	if err := json.Unmarshal([]byte(row.PayloadJSON), &proposal); err != nil {
		return DemuxProposal{}, fmt.Errorf("decode demux proposal %q: %w", row.ID, err)
	}
	proposal.ID = row.ID
	proposal.ProposedChangeID = row.BaseChangeID
	proposal.Status = ProposalStatus(row.Status)
	proposal.CreatedAt = row.CreatedAt
	return proposal, nil
}

func parseDemuxProposalAlias(selector string) (int, bool) {
	selector = strings.TrimSpace(strings.ToLower(selector))
	if !strings.HasPrefix(selector, "d") || len(selector) < 2 {
		return 0, false
	}
	index, err := strconv.Atoi(strings.TrimPrefix(selector, "d"))
	if err != nil || index < 1 {
		return 0, false
	}
	return index, true
}

func summarizeDemuxProposal(proposal DemuxProposal, updatedAt int64) DemuxProposalSummary {
	files := make([]string, 0)
	seen := map[string]struct{}{}
	for _, revision := range proposal.Revisions {
		for _, file := range revision.Files {
			file = strings.TrimSpace(file)
			if file == "" {
				continue
			}
			if _, ok := seen[file]; ok {
				continue
			}
			seen[file] = struct{}{}
			files = append(files, file)
		}
	}
	files = cleanFiles(files)
	firstIntent := ""
	if len(proposal.Revisions) > 0 {
		firstIntent = proposal.Revisions[0].Intent
	}
	return DemuxProposalSummary{
		ID:                      proposal.ID,
		Status:                  proposal.Status,
		ProposedChangeID:        proposal.ProposedChangeID,
		ProposedCommitID:        proposal.ProposedCommitID,
		CreatedAt:               proposal.CreatedAt,
		UpdatedAt:               updatedAt,
		RevisionCount:           len(proposal.Revisions),
		Files:                   files,
		FeasibilityWarningCount: len(proposal.FeasibilityWarnings),
		HiddenDiagnosticCount:   len(proposal.FeasibilityWarnings) + len(proposal.Warnings),
		FirstRevisionIntent:     firstIntent,
	}
}

func (e *Engine) repoForProposal(ctx context.Context, repoRoot string) (RepoInfo, error) {
	if strings.TrimSpace(repoRoot) == "" {
		return e.vcs.ResolveJJRepo(ctx)
	}
	repo, err := e.vcs.ResolveJJRepoAtPath(ctx, repoRoot)
	if err == nil {
		return repo, nil
	}
	gitRepo, gitErr := e.vcs.ResolveGitRepoAtPath(ctx, repoRoot)
	if gitErr != nil {
		return RepoInfo{}, err
	}
	return gitRepo, nil
}

func openStore(ctx context.Context) (*storage.Store, error) {
	db, err := storage.Open(ctx)
	if err != nil {
		return nil, err
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func upsertProposalRepo(ctx context.Context, store *storage.Store, repo RepoInfo, now int64) (int64, error) {
	return store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repo.RootPath,
		Backend:       repo.Backend,
		DefaultRemote: repo.DefaultRemote,
		DefaultBranch: repo.DefaultBranch,
		RemoteURL:     repo.RemoteURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
}

func newProposalID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("demux-%d", time.Now().UnixNano())
	}
	return "demux-" + hex.EncodeToString(buf[:])
}
