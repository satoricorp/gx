package vcs

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	githubapi "github.com/satoricorp/gx/internal/github"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/provenance"
	"github.com/satoricorp/gx/internal/storage"
)

var ErrNoGitBranch = errors.New("gx requires an active branch; create or checkout a branch first")

type Runner interface {
	Run(ctx context.Context, dir, name string, args ...string) (string, error)
	RunStdout(ctx context.Context, dir, name string, args ...string) (string, error)
	RunStream(ctx context.Context, dir, name string, args ...string) error
	RunWithStdin(ctx context.Context, dir, name string, stdin string, args ...string) (string, error)
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	exe, err := resolveExecutable(name)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (ExecRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	exe, err := resolveExecutable(name)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func (ExecRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	exe, err := resolveExecutable(name)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

func (ExecRunner) RunWithStdin(ctx context.Context, dir, name string, stdin string, args ...string) (string, error) {
	exe, err := resolveExecutable(name)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func resolveExecutable(name string) (string, error) {
	if strings.ContainsRune(name, filepath.Separator) {
		return name, nil
	}
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	home, _ := os.UserHomeDir()
	return resolveExecutableFrom(name, home, []string{
		"/opt/homebrew/bin",
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
	})
}

func resolveExecutableFrom(name, home string, dirs []string) (string, error) {
	if home != "" {
		dirs = append([]string{filepath.Join(home, ".local", "bin")}, dirs...)
	}
	for _, dir := range dirs {
		candidate := filepath.Join(dir, name)
		if executableFile(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%s executable not found; install %s or set PATH so %q is available", name, name, name)
}

func executableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}

type Service struct {
	runner Runner
}

func NewService() *Service {
	return &Service{runner: ExecRunner{}}
}

func NewServiceWithRunner(runner Runner) *Service {
	return &Service{runner: runner}
}

type RepoInfo struct {
	RootPath      string
	GitCommonDir  string
	GitDir        string
	Backend       string
	DefaultRemote *string
	DefaultBranch *string
	AuthoringBase *string
	RemoteURL     *string
	BranchName    *string
}

type ChangeInfo struct {
	ChangeID       string
	CommitID       string
	Description    string
	ParentChangeID *string
	Files          []string
}

type CommitResult struct {
	Repo                     RepoInfo
	Change                   ChangeInfo
	Stack                    *StackInfo
	OperationID              string
	PreferredSessionIDs      []string
	SessionContexts          []storage.SessionContext
	ProvenanceStatus         string
	SkipRepoLocalSessions    bool
	SessionEventAttributions []storage.SessionEventAttribution
}

type SplitCommitOptions struct {
	Message                string
	Filesets               []string
	Interactive            bool
	Hunk                   bool
	PatchFile              string
	BookmarkRecordedCommit bool
}

type RevisionOptions struct {
	Message                string
	Filesets               []string
	Interactive            bool
	Hunk                   bool
	PatchFile              string
	PreferredSessionIDs    []string
	SessionContexts        []storage.SessionContext
	BookmarkRecordedCommit bool
}

type ModifyResult struct {
	Repo           RepoInfo
	PreviousChange ChangeInfo
	CurrentChange  ChangeInfo
	Stack          *StackInfo
	OperationID    string
	Output         string
}

type SwitchResult struct {
	Repo          RepoInfo
	Stack         StackInfo
	CurrentChange ChangeInfo
	Output        string
}

type CreateStackResult = SwitchResult

type BaseResult struct {
	Repo          RepoInfo
	BaseRef       string
	DefaultBranch string
	CurrentRef    string
	OnBase        bool
}

type DeleteRevisionResult struct {
	Repo   RepoInfo
	Change ChangeInfo
	Output string
}

type DeleteStackResult struct {
	Repo      RepoInfo
	Stack     StackInfo
	Revisions []ChangeInfo
	Output    string
}

type PruneEmptyStacksResult struct {
	Repo    RepoInfo
	Deleted []DeleteStackResult
}

type PruneGitHubPullRequestsResult struct {
	Checked int
	Removed int
}

type RepairResult struct {
	Repo     RepoInfo `json:"repo"`
	Actions  []string `json:"actions"`
	Warnings []string `json:"warnings,omitempty"`
}

type PushResult struct {
	Repo                 RepoInfo
	CurrentChange        *ChangeInfo
	Stack                *StackInfo
	Published            []PushedChange
	GitExported          bool
	GitPushStatus        string
	GitHubPRStatus       string
	HeadCommitID         string
	RemoteName           *string
	GXStackRef           string
	GXBaseRef            string
	GitPublishedRef      string
	GitCheckoutRef       string
	Output               string
	Warnings             []string
	GitHubPullRequestURL *string
}

type PushedChange struct {
	Change               ChangeInfo
	BranchName           string
	BaseBranchName       string
	Patch                string
	GitHubPullRequestURL *string
}

type StackInfo struct {
	ID             int64
	Name           string
	Alias          string
	BookmarkName   string
	BaseRef        string
	BaseCommitID   string
	HeadChangeID   *string
	HeadCommitID   *string
	RemoteName     *string
	RemoteRef      *string
	GitHubPRURL    *string
	Status         string
	RevisionCount  int
	PublishedCount int
	Revisions      []RevisionSummary
}

type InitResult struct {
	Repo            RepoInfo
	Initialized     bool
	Output          string
	IdentityName    string
	IdentityEmail   string
	IdentityApplied bool
}

type SyncResult struct {
	Repo       RepoInfo
	RemoteName string
	Output     string
}

type RevisionSummary struct {
	Index       int
	ChangeID    string
	CommitID    string
	Description string
	Status      string
	Active      bool
	Published   bool
}

type UnitSummary = RevisionSummary

type StackSummary struct {
	Repo           RepoInfo
	ContainerName  string
	Stack          *StackInfo
	Stacks         []StackInfo
	PublishedCount int
	Units          []UnitSummary
	Revisions      []RevisionSummary

	MissingBaseRefs        []MissingStackBaseRef `json:"missing_base_refs,omitempty"`
	NeedsRebaseOntoDefault bool                  `json:"needs_rebase_onto_default,omitempty"`
	RepairCommand          string                `json:"repair_command,omitempty"`

	GitWorking GitWorkingStatus `json:"git_working,omitempty"`
	Next       []string         `json:"next,omitempty"`
}

type InitOptions struct {
	Name        string
	Email       string
	Interactive bool
	In          io.Reader
	Out         io.Writer
}

func (s *Service) InitWithOptions(ctx context.Context, opts InitOptions) (InitResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return InitResult{}, err
	}
	return s.InitAtPath(ctx, cwd, opts)
}

func (s *Service) EnsureRepoAtPath(ctx context.Context, startPath string) (InitResult, error) {
	return s.InitAtPath(ctx, startPath, InitOptions{})
}

func (s *Service) InitAtPath(ctx context.Context, startPath string, opts InitOptions) (InitResult, error) {
	repo, err := s.ResolveGXRepoAtPath(ctx, startPath)
	if err != nil {
		return InitResult{}, err
	}
	initialized := false
	if ok, initErr := s.isRepoInitializedByIdentity(ctx, repo.GitCommonDir, repo.RootPath); initErr != nil {
		return InitResult{}, initErr
	} else if !ok {
		initialized = true
	}

	name, email, applied, err := s.ensureIdentity(ctx, repo.RootPath, opts)
	if err != nil {
		return InitResult{}, err
	}
	if err := s.ensureGXInternalIgnored(ctx, repo.RootPath); err != nil {
		return InitResult{}, err
	}
	repo, err = s.ResolveGXRepoAtPath(ctx, repo.RootPath)
	if err != nil {
		return InitResult{}, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return InitResult{}, err
	}
	defer store.Close()
	if _, err := upsertRepo(ctx, store, repo); err != nil {
		return InitResult{}, err
	}
	if err := store.RecordInitializedRepo(ctx, repo.RootPath, time.Now().UnixMilli()); err != nil {
		return InitResult{}, err
	}
	return InitResult{
		Repo:            repo,
		Initialized:     initialized,
		IdentityName:    name,
		IdentityEmail:   email,
		IdentityApplied: applied,
	}, nil
}

func (s *Service) Sync(ctx context.Context, remote string) (SyncResult, error) {
	repo, err := s.ResolveGitRepo(ctx)
	if err != nil {
		return SyncResult{}, err
	}
	remoteName := strings.TrimSpace(remote)
	if remoteName == "" {
		if repo.DefaultRemote != nil && strings.TrimSpace(*repo.DefaultRemote) != "" {
			remoteName = strings.TrimSpace(*repo.DefaultRemote)
		} else {
			remoteName = "origin"
		}
	}
	output, err := s.runner.Run(ctx, repo.RootPath, "git", "fetch", "--prune", remoteName)
	if err != nil {
		return SyncResult{}, err
	}
	return SyncResult{
		Repo:       repo,
		RemoteName: remoteName,
		Output:     output,
	}, nil
}

func (s *Service) PruneEmptyStacks(ctx context.Context) (PruneEmptyStacksResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return PruneEmptyStacksResult{}, err
	}
	repo, err := s.ResolveGXRepoAtPath(ctx, cwd)
	if err != nil {
		return PruneEmptyStacksResult{}, err
	}
	result := PruneEmptyStacksResult{Repo: repo}
	err = withRepoIdentityLock(repo.GitCommonDir, func() error {
		model, err := s.loadGitStackReadModel(ctx, repo)
		if err != nil {
			return err
		}
		for _, stack := range model.stacks {
			if !emptyStackForPrune(stack) {
				continue
			}
			deleted, err := s.deleteStackUnlocked(ctx, repo, stack.BookmarkName, false)
			if err != nil {
				return err
			}
			result.Deleted = append(result.Deleted, deleted)
		}
		return nil
	})
	return result, err
}

func emptyStackForPrune(stack StackInfo) bool {
	if IsTerminalStackStatus(stack.Status) {
		return false
	}
	if len(stack.Revisions) > 0 {
		return false
	}
	return stack.RevisionCount == 0
}

func (s *Service) deleteStackUnlocked(ctx context.Context, repo RepoInfo, bookmarkName string, abandonChanges bool) (DeleteStackResult, error) {
	store, err := openStore(ctx)
	if err != nil {
		return DeleteStackResult{}, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
	if err != nil {
		return DeleteStackResult{}, err
	}
	if repoRow == nil {
		return DeleteStackResult{}, fmt.Errorf("unknown repo %s", repo.RootPath)
	}
	storedStack, err := store.FindStackByBookmark(ctx, repoRow.ID, bookmarkName)
	if err != nil {
		return DeleteStackResult{}, err
	}
	if storedStack == nil {
		return DeleteStackResult{}, fmt.Errorf("unknown stack %q", bookmarkName)
	}
	changes, err := store.ListChangesByStackID(ctx, storedStack.ID)
	if err != nil {
		return DeleteStackResult{}, err
	}
	if abandonChanges {
		for _, change := range changes {
			if err := store.MarkChangeStatus(ctx, change.ID, "abandoned", time.Now().UnixMilli()); err != nil {
				return DeleteStackResult{}, err
			}
		}
	}
	if err := store.DeleteStack(ctx, storedStack.ID); err != nil {
		return DeleteStackResult{}, err
	}
	return DeleteStackResult{
		Repo:      repo,
		Stack:     stackInfoFromStorage(*storedStack),
		Revisions: changesToChangeInfo(changes),
	}, nil
}

func (s *Service) RevisionExists(ctx context.Context, repoRoot, rev string) (bool, error) {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return false, nil
	}
	if s.refExists(ctx, repoRoot, rev) {
		return true, nil
	}
	return s.revExists(ctx, repoRoot, rev), nil
}

func (s *Service) RepairWorkflow(ctx context.Context) (RepairResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepairResult{}, err
	}
	repo, err := s.ResolveGXRepoAtPath(ctx, cwd)
	if err != nil {
		return RepairResult{}, err
	}
	return RepairResult{Repo: repo}, nil
}

// shouldReturnCleanCheckoutToBase reports whether workflow repair may move a
// clean checkout back to the authoring base. Deliberate checkouts of regular
// named branches stay put; only detached HEADs and legacy gx-internal refs
// are returned to base.

type PublishMode string

const (
	PublishModeReviewOnly   PublishMode = "review"
	PublishModeReviewAndGit PublishMode = "review+git"
)

type PushOptions struct {
	Mode PublishMode

	// UploadOnly is kept for compatibility with older callers. New callers should
	// prefer Mode so review publication and Git export are explicit.
	UploadOnly bool
}

func (opts PushOptions) gitExportEnabled() bool {
	switch opts.Mode {
	case PublishModeReviewOnly:
		return false
	case PublishModeReviewAndGit:
		return true
	}
	return !opts.UploadOnly
}

type PublishHook func(PushResult) error

func (s *Service) gitCommitIsAncestor(ctx context.Context, repoRoot, ancestor, descendant string) (bool, error) {
	ancestor = strings.TrimSpace(ancestor)
	descendant = strings.TrimSpace(descendant)
	if ancestor == "" || descendant == "" {
		return false, nil
	}
	if _, err := s.runner.Run(ctx, repoRoot, "git", "merge-base", "--is-ancestor", ancestor, descendant); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) remoteBranchHead(ctx context.Context, repoRoot, remoteName, branchName string) (string, error) {
	remoteName = strings.TrimSpace(remoteName)
	branchName = strings.TrimSpace(branchName)
	if remoteName == "" || branchName == "" {
		return "", nil
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "ls-remote", "--heads", remoteName, branchName)
	if err != nil {
		return "", err
	}
	for _, line := range splitLines(out) {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "refs/heads/"+branchName {
			return strings.TrimSpace(fields[0]), nil
		}
	}
	return "", nil
}

func (s *Service) Stack(ctx context.Context) (StackSummary, error) {
	var summary StackSummary
	err := s.preservingGitIndexForCwd(ctx, func() error {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return cwdErr
		}
		repo, err := s.ResolveGXRepoAtPath(ctx, cwd)
		if err != nil {
			return err
		}
		model, err := s.loadGitStackReadModel(ctx, repo)
		if err != nil {
			return err
		}
		summary = stackSummaryForStack(model.repo, model.currentStack, model.currentStackFound, model.stacks, model.currentRevisions, model.currentPublishedCount)
		return nil
	})
	return summary, err
}

// PreservingGitIndexForCwd runs fn with the Git index protected.
func (s *Service) PreservingGitIndexForCwd(ctx context.Context, fn func() error) error {
	return s.preservingGitIndexForCwd(ctx, fn)
}

// preservingGitIndexForCwd applies preservingGitIndex to the git repo that
// contains the current working directory; outside a git repo it runs fn as-is.
func (s *Service) preservingGitIndexForCwd(ctx context.Context, fn func() error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fn()
	}
	root, err := s.runTrimmed(ctx, cwd, "git", "rev-parse", "--show-toplevel")
	if err != nil || strings.TrimSpace(root) == "" {
		return fn()
	}
	return s.preservingGitIndex(ctx, strings.TrimSpace(root), fn)
}

func (s *Service) stackMergedIntoBase(ctx context.Context, repoRoot, fallbackBaseRef string, stack StackInfo, bookmarkTargets map[string]string) bool {
	baseRef := strings.TrimSpace(stack.BaseRef)
	if baseRef == "" || IsTerminalStackStatus(stack.Status) {
		return false
	}
	remoteName := ""
	if stack.RemoteName != nil {
		remoteName = strings.TrimSpace(*stack.RemoteName)
	}
	baseSelectors := s.stackMergeBaseSelectors(ctx, repoRoot, baseRef, remoteName, fallbackBaseRef)
	selectors := make([]string, 0, 2)
	if bookmark := strings.TrimSpace(stack.BookmarkName); bookmark != "" {
		if _, exists := bookmarkTargets[bookmark]; exists {
			selectors = append(selectors, bookmark)
		}
	}
	if len(selectors) == 0 {
		if stack.HeadCommitID != nil {
			if head := strings.TrimSpace(*stack.HeadCommitID); head != "" {
				selectors = append(selectors, head)
			}
		}
		if stack.HeadChangeID != nil {
			if head := strings.TrimSpace(*stack.HeadChangeID); head != "" {
				selectors = append(selectors, head)
			}
		}
	}
	for _, selector := range selectors {
		for _, baseSelector := range baseSelectors {
			if merged, _ := s.gitCommitIsAncestor(ctx, repoRoot, selector, baseSelector); merged {
				return true
			}
		}
	}
	return false
}

func stackVisibleInGX(stack StackInfo) bool {
	return !IsTerminalStackStatus(stack.Status)
}

func stackSummaryForStack(repo RepoInfo, stack StackInfo, stackFound bool, stacks []StackInfo, revisions []RevisionSummary, publishedCount int) StackSummary {
	var stackPtr *StackInfo
	containerName := ""
	if stackFound {
		stackPtr = &stack
		containerName = stack.BookmarkName
	}
	return StackSummary{
		Repo:           repo,
		ContainerName:  containerName,
		Stack:          stackPtr,
		Stacks:         stacks,
		PublishedCount: publishedCount,
		Units:          revisions,
		Revisions:      revisions,
	}
}

func (s *Service) ResolveGitRepo(ctx context.Context) (RepoInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveGitRepoAtPath(ctx, cwd)
}

func (s *Service) ResolveGitRepoAtPath(ctx context.Context, startPath string) (RepoInfo, error) {
	info, err := s.resolveGitRepoWithoutStoreAtPath(ctx, startPath)
	if err != nil {
		return RepoInfo{}, err
	}
	return s.withStoredRepoConfigByIdentity(ctx, info), nil
}

// ResolveGitRepoWithoutStore resolves the repo from git alone. Opening the GX
// store creates ~/.gx and applies its schema, so a read-only command like
// `gx review` must not go through ResolveGitRepo: it would leave GX state
// behind on a machine that has never run `gx init`. Everything review needs
// (root, remote, branch) comes from git; only AuthoringBase and a stored
// backend override come from the store, and review uses neither.
func (s *Service) ResolveGitRepoWithoutStore(ctx context.Context) (RepoInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, err
	}
	return s.resolveGitRepoWithoutStoreAtPath(ctx, cwd)
}

func (s *Service) resolveGitRepoWithoutStoreAtPath(ctx context.Context, startPath string) (RepoInfo, error) {
	paths, err := s.ResolveGitPaths(ctx, startPath)
	if err != nil {
		return RepoInfo{}, err
	}
	info, err := s.resolveGitInfo(ctx, paths.WorktreeRoot)
	if err != nil {
		return RepoInfo{}, err
	}
	info.RootPath = paths.WorktreeRoot
	info.GitCommonDir = paths.CommonDir
	if info.Backend == "" {
		info.Backend = "git"
	}
	return info, nil
}

func (s *Service) withStoredRepoConfig(ctx context.Context, info RepoInfo) RepoInfo {
	store, err := openStore(ctx)
	if err != nil {
		return info
	}
	defer store.Close()
	repo, err := store.FindRepoByIdentity(ctx, info.GitCommonDir, info.RootPath)
	if err != nil || repo == nil {
		return info
	}
	info.AuthoringBase = repo.AuthoringBase
	return info
}

func (s *Service) resolveGitInfo(ctx context.Context, root string) (RepoInfo, error) {
	info := RepoInfo{RootPath: root}
	if value, err := s.gitValue(ctx, root, "remote"); err == nil && value != "" {
		lines := splitLines(value)
		if len(lines) > 0 {
			info.DefaultRemote = &lines[0]
			if url, err := s.gitValue(ctx, root, "config", "--get", "remote."+lines[0]+".url"); err == nil && url != "" {
				info.RemoteURL = &url
			} else if url, err := s.gitValue(ctx, root, "remote", "get-url", lines[0]); err == nil && url != "" {
				info.RemoteURL = &url
			}
		}
	}
	if branch, err := s.gitValue(ctx, root, "branch", "--show-current"); err == nil && branch != "" {
		info.BranchName = &branch
	}
	if head, err := s.gitValue(ctx, root, "symbolic-ref", "refs/remotes/origin/HEAD", "--short"); err == nil && head != "" {
		if idx := strings.LastIndex(head, "/"); idx >= 0 && idx+1 < len(head) {
			info.DefaultBranch = ptr(head[idx+1:])
		}
	}
	return info, nil
}

func (s *Service) CurrentChange(ctx context.Context, repoRoot, rev string) (ChangeInfo, error) {
	return s.gitCurrentChange(ctx, repoRoot, rev)
}

func recordCommit(ctx context.Context, result CommitResult) error {
	return withBusyRetry(ctx, "record commit metadata", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		return recordChangeForStack(ctx, store, result.Repo, result.Stack, result.Change, result.OperationID, result.PreferredSessionIDs, result.SessionContexts, result.SkipRepoLocalSessions, result.SessionEventAttributions)
	})
}

func (s *Service) PrunePublishedStackByRef(ctx context.Context, repo RepoInfo, publishRef string) (bool, error) {
	publishRef = strings.TrimPrefix(strings.TrimSpace(publishRef), "refs/heads/")
	if publishRef == "" {
		return false, nil
	}
	var pruned bool
	err := withBusyRetry(ctx, "prune published stack", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()

		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		stacks, err := store.ListStacksByRepoID(ctx, repoID)
		if err != nil {
			return err
		}
		for _, stack := range stacks {
			info := stackInfoFromStorage(stack)
			if stackPublishRefMatches(info, publishRef) {
				if err := store.PrunePublishedStack(ctx, repoID, stack.ID, publishRefForStack(info), time.Now().UnixMilli()); err != nil {
					return err
				}
				pruned = true
				return nil
			}
		}
		return nil
	})
	return pruned, err
}

func (s *Service) PruneTerminalGitHubPullRequestStacks(ctx context.Context, repo RepoInfo) (PruneGitHubPullRequestsResult, error) {
	var result PruneGitHubPullRequestsResult
	err := withBusyRetry(ctx, "prune terminal GitHub pull request stacks", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()

		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		stacks, err := store.ListStacksByRepoID(ctx, repoID)
		if err != nil {
			return err
		}
		clients := map[string]*githubapi.Client{}
		for _, stack := range stacks {
			info := stackInfoFromStorage(stack)
			if IsTerminalStackStatus(info.Status) || info.GitHubPRURL == nil {
				continue
			}
			ref, ok := githubPullRequestRefFromURL(*info.GitHubPRURL)
			if !ok {
				continue
			}
			client := clients[ref.host]
			if client == nil {
				client, err = githubapi.NewClient(ref.host)
				if err != nil {
					return err
				}
				clients[ref.host] = client
			}
			pr, err := client.GetPullRequest(ctx, ref.owner, ref.repo, ref.number)
			if err != nil {
				return err
			}
			result.Checked++
			if pr == nil || !terminalGitHubPullRequest(*pr) {
				continue
			}
			if err := store.PrunePublishedStack(ctx, repoID, stack.ID, publishRefForStack(info), time.Now().UnixMilli()); err != nil {
				return err
			}
			result.Removed++
		}
		return nil
	})
	return result, err
}

type githubPullRequestRef struct {
	host   string
	owner  string
	repo   string
	number int
}

func githubPullRequestRefFromURL(raw string) (githubPullRequestRef, bool) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return githubPullRequestRef{}, false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[0] == "" || parts[1] == "" || parts[2] != "pull" {
		return githubPullRequestRef{}, false
	}
	number, err := strconv.Atoi(parts[3])
	if err != nil || number <= 0 {
		return githubPullRequestRef{}, false
	}
	return githubPullRequestRef{host: parsed.Host, owner: parts[0], repo: strings.TrimSuffix(parts[1], ".git"), number: number}, true
}

func terminalGitHubPullRequest(pr githubapi.PullRequest) bool {
	return pr.Merged || strings.EqualFold(strings.TrimSpace(pr.State), "closed")
}

func stackPublishRefMatches(stack StackInfo, publishRef string) bool {
	publishRef = strings.TrimPrefix(strings.TrimSpace(publishRef), "refs/heads/")
	if publishRef == "" {
		return false
	}
	if publishRefForStack(stack) == publishRef {
		return true
	}
	if strings.TrimSpace(stack.BookmarkName) == publishRef {
		return true
	}
	if stack.RemoteRef != nil && strings.TrimPrefix(strings.TrimSpace(*stack.RemoteRef), "refs/heads/") == publishRef {
		return true
	}
	return false
}

func withBusyRetry(ctx context.Context, label string, fn func() error) error {
	const maxAttempts = 5
	backoff := 80 * time.Millisecond
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err
			if !isSQLiteBusy(err) || attempt == maxAttempts {
				break
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			continue
		}
		return nil
	}
	if lastErr == nil {
		return nil
	}
	if isSQLiteBusy(lastErr) {
		return fmt.Errorf("%s failed after %d retries: %w", label, maxAttempts, lastErr)
	}
	return lastErr
}

func isSQLiteBusy(err error) bool {
	if err == nil {
		return false
	}
	value := strings.ToLower(err.Error())
	return strings.Contains(value, "sqlite_busy") || strings.Contains(value, "database is locked")
}

func recordChangeForStack(ctx context.Context, store *storage.Store, repo RepoInfo, stack *StackInfo, change ChangeInfo, opID string, preferredSessionIDs []string, sessionContexts []storage.SessionContext, skipRepoLocalSessions bool, eventAttributions []storage.SessionEventAttribution) error {
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return err
	}
	changeID, err := upsertChange(ctx, store, repoID, change)
	if err != nil {
		return err
	}
	if err := writeRevision(ctx, store, changeID, change, opID); err != nil {
		return err
	}
	if stack != nil {
		head := change.ChangeID
		headCommit := change.CommitID
		stack.HeadChangeID = &head
		stack.HeadCommitID = &headCommit
		stackID, err := upsertStackInfo(ctx, store, repoID, *stack)
		if err != nil {
			return err
		}
		if err := store.AddChangeToStack(ctx, stackID, changeID, time.Now().UnixMilli()); err != nil {
			return err
		}
	}
	if skipRepoLocalSessions {
		if err := attachExplicitSessions(ctx, store, changeID, preferredSessionIDs); err != nil {
			return err
		}
	} else {
		if err := attachSessions(ctx, store, repo.RootPath, changeID, preferredSessionIDs); err != nil {
			return err
		}
	}
	if err := writeSessionContexts(ctx, store, repo.RootPath, sessionContexts); err != nil {
		return err
	}
	if len(eventAttributions) > 0 {
		stackBookmark := ""
		if stack != nil {
			stackBookmark = stack.BookmarkName
		}
		now := time.Now().UnixMilli()
		repoRootValue := repo.RootPath
		for i := range eventAttributions {
			eventAttributions[i].RepoID = repoID
			eventAttributions[i].ChangeID = changeID
			if stackBookmark != "" && eventAttributions[i].StackBookmark == nil {
				eventAttributions[i].StackBookmark = &stackBookmark
			}
			if eventAttributions[i].CreatedAt == 0 {
				eventAttributions[i].CreatedAt = now
			}
			if strings.TrimSpace(eventAttributions[i].SessionID) == "" {
				continue
			}
			createdAt := eventAttributions[i].CreatedAt
			if createdAt == 0 {
				createdAt = now
			}
			if err := store.UpsertSession(ctx, storage.Session{
				ID:        eventAttributions[i].SessionID,
				CreatedAt: createdAt,
				Command:   eventAttributions[i].Tool,
				Cwd:       repo.RootPath,
				RepoRoot:  &repoRootValue,
			}); err != nil {
				return err
			}
		}
		if err := store.WriteSessionEventAttributions(ctx, eventAttributions); err != nil {
			return err
		}
	}
	return nil
}

func attachExplicitSessions(ctx context.Context, store *storage.Store, changeID int64, preferredSessionIDs []string) error {
	sessionIDs := append([]string(nil), preferredSessionIDs...)
	sessionIDs = append(sessionIDs, provenance.ExplicitSessionIDsFromEnv()...)
	sessionIDs = uniqueStrings(sessionIDs)
	if len(sessionIDs) == 0 {
		return nil
	}
	filtered, err := store.FilterExistingSessionIDs(ctx, sessionIDs)
	if err != nil {
		return err
	}
	if len(filtered) == 0 {
		return nil
	}
	return store.WriteChangeSessions(ctx, changeID, filtered, time.Now().UnixMilli())
}

func writeSessionContexts(ctx context.Context, store *storage.Store, repoRoot string, contexts []storage.SessionContext) error {
	for _, sessionContext := range contexts {
		if strings.TrimSpace(sessionContext.SessionID) == "" {
			continue
		}
		createdAt := sessionContext.CapturedAt
		if createdAt == 0 {
			createdAt = time.Now().UnixMilli()
		}
		repoRootValue := repoRoot
		session := storage.Session{
			ID:        sessionContext.SessionID,
			CreatedAt: createdAt,
			Command:   sessionContext.Tool,
			Cwd:       repoRoot,
			RepoRoot:  &repoRootValue,
		}
		if err := store.UpsertSessionContext(ctx, session, sessionContext); err != nil {
			return err
		}
	}
	return nil
}

func upsertRepo(ctx context.Context, store *storage.Store, repo RepoInfo) (int64, error) {
	now := time.Now().UnixMilli()
	return store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repo.RootPath,
		GitCommonDir:  repo.GitCommonDir,
		Backend:       repo.Backend,
		DefaultRemote: repo.DefaultRemote,
		DefaultBranch: repo.DefaultBranch,
		AuthoringBase: repo.AuthoringBase,
		RemoteURL:     repo.RemoteURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
}

func upsertChange(ctx context.Context, store *storage.Store, repoID int64, change ChangeInfo) (int64, error) {
	now := time.Now().UnixMilli()
	description := change.Description
	if description == "" {
		description = "(no description set)"
	}
	return store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      change.ChangeID,
		CurrentCommitID: change.CommitID,
		Description:     description,
		ParentChangeID:  change.ParentChangeID,
		Status:          "draft",
		FirstSeenAt:     now,
		UpdatedAt:       now,
	})
}

func changesToChangeInfo(changes []storage.Change) []ChangeInfo {
	out := make([]ChangeInfo, 0, len(changes))
	for _, change := range changes {
		out = append(out, ChangeInfo{
			ChangeID:       change.JJChangeID,
			CommitID:       change.CurrentCommitID,
			Description:    change.Description,
			ParentChangeID: change.ParentChangeID,
		})
	}
	return out
}

func upsertStackInfo(ctx context.Context, store *storage.Store, repoID int64, stack StackInfo) (int64, error) {
	now := time.Now().UnixMilli()
	return store.UpsertStack(ctx, storage.Stack{
		ID:           stack.ID,
		RepoID:       repoID,
		Name:         firstNonEmpty(stack.Name, stack.BookmarkName),
		BookmarkName: stack.BookmarkName,
		BaseRef:      firstNonEmpty(stack.BaseRef, "main"),
		BaseCommitID: stack.BaseCommitID,
		HeadChangeID: stack.HeadChangeID,
		HeadCommitID: stack.HeadCommitID,
		RemoteName:   stack.RemoteName,
		RemoteRef:    stack.RemoteRef,
		GitHubPRURL:  stack.GitHubPRURL,
		Status:       firstNonEmpty(stack.Status, "draft"),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func writeRevision(ctx context.Context, store *storage.Store, changeID int64, change ChangeInfo, opID string) error {
	filesJSON, err := json.Marshal(change.Files)
	if err != nil {
		return err
	}
	return store.WriteChangeRevision(ctx, storage.ChangeRevision{
		ChangeID:      changeID,
		JJCommitID:    change.CommitID,
		JJOperationID: opID,
		ChangedFiles:  string(filesJSON),
		CreatedAt:     time.Now().UnixMilli(),
	})
}

func attachSessions(ctx context.Context, store *storage.Store, repoRoot string, changeID int64, preferredSessionIDs []string) error {
	_, err := provenance.AttachPreferred(ctx, store, repoRoot, changeID, time.Now().UnixMilli(), preferredSessionIDs)
	return err
}

func openStore(ctx context.Context) (*storage.Store, error) {
	db, err := storage.Open(ctx)
	if err != nil {
		return nil, err
	}
	return storage.NewStore(ctx, db)
}

func (s *Service) runTrimmed(ctx context.Context, dir, name string, args ...string) (string, error) {
	out, err := s.runner.Run(ctx, dir, name, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (s *Service) runStdoutTrimmed(ctx context.Context, dir, name string, args ...string) (string, error) {
	out, err := s.runner.RunStdout(ctx, dir, name, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (s *Service) gitValue(ctx context.Context, dir string, args ...string) (string, error) {
	out, err := s.runTrimmed(ctx, dir, "git", args...)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (s *Service) ensureIdentity(ctx context.Context, repoRoot string, opts InitOptions) (string, string, bool, error) {
	cfg, err := gxconfig.Load()
	if err != nil {
		return "", "", false, err
	}
	gitName := s.gitConfigValue(ctx, repoRoot, "user.name")
	gitEmail := s.gitConfigValue(ctx, repoRoot, "user.email")

	name := firstNonEmpty(
		strings.TrimSpace(opts.Name),
		strings.TrimSpace(cfg.User.Name),
		strings.TrimSpace(gitName),
	)
	email := firstNonEmpty(
		strings.TrimSpace(opts.Email),
		strings.TrimSpace(cfg.User.Email),
		strings.TrimSpace(gitEmail),
	)

	if opts.Interactive {
		out := opts.Out
		if out == nil {
			out = os.Stdout
		}
		in := opts.In
		if in == nil {
			in = os.Stdin
		}
		promptName := strings.TrimSpace(opts.Name) == ""
		promptEmail := strings.TrimSpace(opts.Email) == ""
		if promptName || promptEmail {
			fmt.Fprintf(out, "Your config is stored in %s\n", gxconfig.DisplayPath())
			fmt.Fprintln(out, "and by initializing with gx you share your email with gx.")
			fmt.Fprintln(out)
		}
		if promptName {
			var promptErr error
			name, promptErr = promptRequiredValue(in, out, promptFocus("gx name"), name)
			if promptErr != nil {
				return "", "", false, promptErr
			}
		}
		if promptEmail {
			var promptErr error
			email, promptErr = promptRequiredValue(in, out, promptFocus("gx email"), email)
			if promptErr != nil {
				return "", "", false, promptErr
			}
		}
	}

	if name == "" || email == "" {
		return name, email, false, nil
	}

	changed := false
	if cfg.User.Name != name || cfg.User.Email != email {
		cfg.User.Name = name
		cfg.User.Email = email
		if err := gxconfig.Save(cfg); err != nil {
			return "", "", false, err
		}
		changed = true
	}
	return name, email, changed, nil
}

func (s *Service) gitConfigValue(ctx context.Context, dir, name string) string {
	if value, err := s.runTrimmed(ctx, dir, "git", "config", name); err == nil && value != "" {
		return value
	}
	if value, err := s.runTrimmed(ctx, dir, "git", "config", "--global", name); err == nil && value != "" {
		return value
	}
	return ""
}

func (s *Service) ensureGXInternalIgnored(ctx context.Context, repoRoot string) error {
	excludePath, err := s.gitValue(ctx, repoRoot, "rev-parse", "--git-path", "info/exclude")
	if err != nil || strings.TrimSpace(excludePath) == "" {
		return nil
	}
	excludePath = strings.TrimSpace(excludePath)
	if !filepath.IsAbs(excludePath) {
		excludePath = filepath.Join(repoRoot, excludePath)
	}
	data, err := os.ReadFile(excludePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == ".gx/" {
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(excludePath), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(excludePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}
	_, err = file.WriteString(".gx/\n")
	return err
}

type stackResolveOptions struct {
	create      bool
	description string
}

func (s *Service) defaultStackBaseRef(repo RepoInfo) string {
	if repo.AuthoringBase != nil && strings.TrimSpace(*repo.AuthoringBase) != "" {
		return strings.TrimPrefix(strings.TrimSpace(*repo.AuthoringBase), "origin/")
	}
	if repo.BranchName != nil && strings.TrimSpace(*repo.BranchName) != "" {
		branch := strings.TrimSpace(*repo.BranchName)
		if base, ok := gxAuthoringBaseFromCheckoutRef(branch); ok {
			return base
		}
	}
	if repo.DefaultBranch != nil && strings.TrimSpace(*repo.DefaultBranch) != "" {
		return strings.TrimSpace(*repo.DefaultBranch)
	}
	return "main"
}

func (s *Service) publicStackBaseRef(ctx context.Context, repo RepoInfo, baseRef string) string {
	baseRef = strings.TrimPrefix(strings.TrimSpace(baseRef), "origin/")
	baseRef = strings.TrimPrefix(baseRef, "refs/heads/")
	if baseRef == "" {
		return repo.defaultBaseBranch()
	}
	if legacyStackBookmarkName(baseRef) && !legacyGXInternalCheckoutRef(baseRef) {
		return stackBookmarkName(baseRef, "")
	}
	if base, ok := gxAuthoringBaseFromCheckoutRef(baseRef); ok {
		if base == repo.defaultBaseBranch() {
			if _, err := s.commitIDForRev(ctx, repo.RootPath, base); err == nil {
				return base
			}
		}
		configuredBase := strings.TrimSpace("")
		if repo.AuthoringBase != nil {
			configuredBase = strings.TrimSpace(*repo.AuthoringBase)
		}
		if (base == configuredBase) && s.refsPointToSameChange(ctx, repo.RootPath, baseRef, base) {
			return base
		}
	}
	if legacyGXInternalCheckoutRef(baseRef) {
		if bookmark := s.publicBookmarkForRev(ctx, repo.RootPath, baseRef, repo.defaultBaseBranch()); bookmark != "" {
			return bookmark
		}
		return repo.defaultBaseBranch()
	}
	return baseRef
}

func (s *Service) publicBookmarkForRev(ctx context.Context, repoRoot, rev, defaultBranch string) string {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "for-each-ref", "--format=%(refname:short)", "--points-at", rev, "refs/heads")
	if err != nil {
		return ""
	}
	bestRank := 99
	best := ""
	for _, line := range splitLines(out) {
		name := strings.TrimSpace(line)
		if name == "" || legacyGXInternalCheckoutRef(name) {
			continue
		}
		rank := 3
		switch {
		case name != defaultBranch && !isGXStackBookmark(name):
			rank = 0
		case name == defaultBranch:
			rank = 1
		case isGXStackBookmark(name):
			rank = 2
		}
		if rank < bestRank || (rank == bestRank && (best == "" || name < best)) {
			bestRank = rank
			best = name
		}
	}
	return best
}

func (s *Service) refsPointToSameChange(ctx context.Context, repoRoot, left, right string) bool {
	leftCommit, err := s.commitIDForRev(ctx, repoRoot, left)
	if err != nil || strings.TrimSpace(leftCommit) == "" {
		return false
	}
	rightCommit, err := s.commitIDForRev(ctx, repoRoot, right)
	if err != nil || strings.TrimSpace(rightCommit) == "" {
		return false
	}
	return leftCommit == rightCommit
}

func (s *Service) stackBaseCommitID(ctx context.Context, repoRoot, baseRef string) string {
	if value, err := s.gitValue(ctx, repoRoot, "rev-parse", baseRef); err == nil && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	if value, err := s.gitValue(ctx, repoRoot, "rev-parse", "HEAD"); err == nil && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return ""
}

func stackInfoFromStorage(body storage.Stack) StackInfo {
	return StackInfo{
		ID:           body.ID,
		Name:         body.Name,
		BookmarkName: body.BookmarkName,
		BaseRef:      body.BaseRef,
		BaseCommitID: body.BaseCommitID,
		HeadChangeID: body.HeadChangeID,
		HeadCommitID: body.HeadCommitID,
		RemoteName:   body.RemoteName,
		RemoteRef:    body.RemoteRef,
		GitHubPRURL:  body.GitHubPRURL,
		Status:       body.Status,
	}
}

func stackAlias(index int) string {
	return "s" + strconv.Itoa(index+1)
}

func (s *Service) setBookmarkTargetAtRev(ctx context.Context, repoRoot, name, targetRev string, allowBackwards bool) error {
	targetRev = strings.TrimSpace(targetRev)
	if targetRev == "" {
		return fmt.Errorf("container target revision is empty")
	}
	if err := s.ensureBookmarkMutationAllowed(ctx, repoRoot, name, targetRev); err != nil {
		return err
	}
	targetCommit, err := s.commitIDForRev(ctx, repoRoot, targetRev)
	if err != nil {
		return err
	}
	_, err = s.runner.Run(ctx, repoRoot, "git", "update-ref", "refs/heads/"+cleanRefName(name), targetCommit)
	return err
}

func (s *Service) commitIDForRev(ctx context.Context, repoRoot, rev string) (string, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", strings.TrimSpace(rev)+"^{commit}")
	if err != nil {
		return "", err
	}
	return lastNonEmptyLine(out), nil
}

func (s *Service) changeIDForRev(ctx context.Context, repoRoot, rev string) (string, error) {
	change, err := s.gitChangeBySelector(ctx, repoRoot, rev)
	if err != nil {
		return "", err
	}
	return change.ChangeID, nil
}

func (s *Service) isRevisionEmpty(ctx context.Context, repoRoot, rev string) (bool, error) {
	value, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", strings.TrimSpace(rev))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(value) == "", nil
}

func splitLines(value string) []string {
	if value == "" {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(value), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func lastNonEmptyLine(value string) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line
		}
	}
	return ""
}

func ptr[T any](value T) *T {
	return &value
}

func (r RepoInfo) defaultBaseBranch() string {
	if r.DefaultBranch != nil && strings.TrimSpace(*r.DefaultBranch) != "" {
		return *r.DefaultBranch
	}
	return "main"
}

func (r RepoInfo) authoringBaseRef() string {
	if r.AuthoringBase != nil && strings.TrimSpace(*r.AuthoringBase) != "" {
		ref := strings.TrimPrefix(strings.TrimSpace(*r.AuthoringBase), "origin/")
		if legacyStackBookmarkName(ref) && !legacyGXInternalCheckoutRef(ref) {
			return stackBookmarkName(ref, "")
		}
		if isGXStackBookmark(ref) {
			return ref
		}
		if base, ok := gxAuthoringBaseFromCheckoutRef(ref); ok {
			return base
		}
		if legacyGXInternalCheckoutRef(ref) {
			return r.defaultBaseBranch()
		}
		return ref
	}
	return r.defaultBaseBranch()
}

func shortID(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func promptRequiredValue(in io.Reader, out io.Writer, label, current string) (string, error) {
	reader := bufio.NewReader(in)
	for {
		if current != "" {
			fmt.Fprintf(out, "%s [%s]: ", label, current)
		} else {
			fmt.Fprintf(out, "%s: ", label)
		}
		raw, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		value := strings.TrimSpace(raw)
		if value == "" {
			value = strings.TrimSpace(current)
		}
		if value != "" {
			return value, nil
		}
		if errors.Is(err, io.EOF) {
			return "", fmt.Errorf("%s is required", strings.ToLower(label))
		}
	}
}

func promptFocus(text string) string {
	if os.Getenv("NO_COLOR") != "" || strings.EqualFold(os.Getenv("TERM"), "dumb") || text == "" {
		return text
	}
	return "\x1b[1m\x1b[36m" + text + "\x1b[0m"
}

var githubSSHPattern = regexp.MustCompile(`^(?P<user>[^@]+@)?(?P<host>[^:]+):(?P<owner>[^/]+)/(?P<repo>[^/]+?)(?:\.git)?$`)

func publishRefForStack(body StackInfo) string {
	if bookmark := strings.TrimSpace(body.BookmarkName); bookmark != "" {
		if legacyStackBookmarkName(bookmark) {
			return stackBookmarkName(bookmark, derefString(body.HeadChangeID))
		}
		return cleanRefName(bookmark)
	}
	return stackBookmarkName(body.Name, derefString(body.HeadChangeID))
}
