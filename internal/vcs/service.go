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
	"path"
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
	runner   Runner
	progress io.Writer
}

func NewService() *Service {
	return &Service{runner: ExecRunner{}}
}

func NewServiceWithRunner(runner Runner) *Service {
	return &Service{runner: runner}
}

func (s *Service) SetProgressWriter(out io.Writer) {
	s.progress = out
}

func (s *Service) progressf(format string, args ...any) {
	if s.progress == nil {
		return
	}
	fmt.Fprintf(s.progress, format+"\n", args...)
}

type RepoInfo struct {
	RootPath      string
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
	Repo                RepoInfo
	Change              ChangeInfo
	Stack               *StackInfo
	OperationID         string
	Output              string
	PreferredSessionIDs []string
}

type SplitCommitOptions struct {
	Message     string
	Filesets    []string
	Interactive bool
	Hunk        bool
	PatchFile   string
}

type RevisionOptions struct {
	Message             string
	Filesets            []string
	Interactive         bool
	Hunk                bool
	PatchFile           string
	PreferredSessionIDs []string
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
}

type InitOptions struct {
	Name        string
	Email       string
	Interactive bool
	In          io.Reader
	Out         io.Writer
}

func (s *Service) Init(ctx context.Context) (InitResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return InitResult{}, err
	}
	return s.InitAtPath(ctx, cwd, InitOptions{})
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
	repo, err := s.ResolveJJRepoAtPath(ctx, startPath)
	initialized := false
	output := ""
	if err != nil {
		gitRepo, gitErr := s.ResolveGitRepoAtPath(ctx, startPath)
		targetPath := startPath
		if gitErr == nil {
			targetPath = gitRepo.RootPath
		}

		output, err = s.runner.Run(ctx, targetPath, "jj", "git", "init", ".")
		if err != nil {
			return InitResult{}, err
		}
		repo, err = s.ResolveJJRepoAtPath(ctx, targetPath)
		if err != nil {
			return InitResult{}, err
		}
		initialized = true
	}

	name, email, applied, err := s.ensureIdentity(ctx, repo.RootPath, opts)
	if err != nil {
		return InitResult{}, err
	}
	if err := s.ensureGXInternalIgnored(ctx, repo.RootPath); err != nil {
		return InitResult{}, err
	}
	repo, err = s.ResolveJJRepoAtPath(ctx, repo.RootPath)
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
		Output:          output,
		IdentityName:    name,
		IdentityEmail:   email,
		IdentityApplied: applied,
	}, nil
}

func (s *Service) Commit(ctx context.Context, message string) (CommitResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var commitErr error
		result, commitErr = s.commitUnlocked(ctx, repo, message)
		return commitErr
	})
	return result, err
}

func (s *Service) RecordRevision(ctx context.Context, opts RevisionOptions) (CommitResult, error) {
	var result CommitResult
	var err error
	if opts.Hunk || opts.Interactive || len(opts.Filesets) > 0 {
		result, err = s.SplitCommit(ctx, SplitCommitOptions{
			Message:     opts.Message,
			Filesets:    opts.Filesets,
			Interactive: opts.Interactive,
			Hunk:        opts.Hunk,
			PatchFile:   opts.PatchFile,
		})
	} else {
		result, err = s.Commit(ctx, opts.Message)
	}
	if err != nil {
		return CommitResult{}, err
	}
	result.PreferredSessionIDs = opts.PreferredSessionIDs
	if err := recordCommit(ctx, result); err != nil {
		return result, fmt.Errorf("record revision metadata: %w", err)
	}
	return result, nil
}

func (s *Service) RecordAuthoringRevision(ctx context.Context, opts RevisionOptions) (CommitResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	_, mode, err := s.ensureAuthoringCheckout(ctx, repo, "gx add")
	if err != nil {
		return CommitResult{}, err
	}
	if mode == "base" {
		if err := s.ensureNewDraftStackForCurrent(ctx, repo, opts.Message); err != nil {
			return CommitResult{}, err
		}
	}
	return s.RecordRevision(ctx, opts)
}

func (s *Service) ensureNewDraftStackForCurrent(ctx context.Context, repo RepoInfo, description string) error {
	return withRepoLock(repo.RootPath, func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		_, err = s.createDraftStack(ctx, store, repo, repoID, description)
		return err
	})
}

func (s *Service) RecordCurrentRevisionInStack(ctx context.Context, stackBookmark, message string, preferredSessionIDs []string) (CommitResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var commitErr error
		result, commitErr = s.commitCurrentRevisionInStackUnlocked(ctx, repo, stackBookmark, message)
		return commitErr
	})
	if err != nil {
		return CommitResult{}, err
	}
	result.PreferredSessionIDs = preferredSessionIDs
	if err := recordCommit(ctx, result); err != nil {
		return result, fmt.Errorf("record revision metadata: %w", err)
	}
	return result, nil
}

func (s *Service) RecordCurrentRevisionInNewStack(ctx context.Context, stackName, bookmarkName, baseRef, message string, preferredSessionIDs []string) (CommitResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var commitErr error
		result, commitErr = s.commitCurrentRevisionInNewStackUnlocked(ctx, repo, stackName, bookmarkName, baseRef, message)
		return commitErr
	})
	if err != nil {
		return CommitResult{}, err
	}
	result.PreferredSessionIDs = preferredSessionIDs
	if err := recordCommit(ctx, result); err != nil {
		return result, fmt.Errorf("record revision metadata: %w", err)
	}
	return result, nil
}

func (s *Service) commitCurrentRevisionInStackUnlocked(ctx context.Context, repo RepoInfo, stackBookmark, message string) (CommitResult, error) {
	if err := ValidateCommitMessage(message); err != nil {
		return CommitResult{}, err
	}
	checkoutBranch := recordedEditCheckoutBranch(repo)
	store, err := openStore(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return CommitResult{}, err
	}
	stack, err := store.FindStackByBookmark(ctx, repoID, strings.TrimSpace(stackBookmark))
	if err != nil {
		return CommitResult{}, err
	}
	if stack == nil {
		return CommitResult{}, fmt.Errorf("unknown GX stack bookmark %q", stackBookmark)
	}
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", "commit", "-m", message)
	if err != nil {
		return CommitResult{}, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	stackInfo := stackInfoFromStorage(*stack)
	repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, stackInfo.BookmarkName, checkoutBranch)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, stackInfo.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &stackInfo,
		OperationID: opID,
		Output:      output,
	}, nil
}

func (s *Service) commitCurrentRevisionInNewStackUnlocked(ctx context.Context, repo RepoInfo, stackName, bookmarkName, baseRef, message string) (CommitResult, error) {
	if err := ValidateCommitMessage(message); err != nil {
		return CommitResult{}, err
	}
	checkoutBranch := recordedEditCheckoutBranch(repo)
	store, err := openStore(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return CommitResult{}, err
	}
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", "commit", "-m", message)
	if err != nil {
		return CommitResult{}, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	name := firstNonEmpty(strings.TrimSpace(stackName), stackNameFromBookmark(bookmarkName), message)
	bookmark := stackBookmarkName(firstNonEmpty(bookmarkName, name), change.ChangeID)
	if existing, err := s.resolveAppendExistingStack(ctx, store, repo, repoID, bookmark); err != nil {
		return CommitResult{}, err
	} else if existing != nil {
		stackInfo := stackInfoFromStorage(*existing)
		repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, stackInfo.BookmarkName, checkoutBranch)
		if err != nil {
			return CommitResult{}, err
		}
		if err := s.reconcileRepoChanges(ctx, repo, stackInfo.BookmarkName); err != nil {
			return CommitResult{}, err
		}
		opID, err := s.CurrentOperation(ctx, repo.RootPath)
		if err != nil {
			return CommitResult{}, err
		}
		return CommitResult{
			Repo:        repo,
			Change:      change,
			Stack:       &stackInfo,
			OperationID: opID,
			Output:      output,
		}, nil
	}
	if err := s.relocateKnownStackBookmarksFromChange(ctx, store, repoID, repo.RootPath, bookmark, "@-"); err != nil {
		return CommitResult{}, err
	}
	if err := s.assertBookmarkTargetAvailable(ctx, store, repoID, repo.RootPath, bookmark, "@-"); err != nil {
		return CommitResult{}, err
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "bookmark", "set", bookmark, "-r", "@-"); err != nil {
		return CommitResult{}, err
	}
	base := s.publicStackBaseRef(ctx, repo, firstNonEmpty(strings.TrimSpace(baseRef), s.defaultStackBaseRef(repo)))
	body := StackInfo{
		Name:         name,
		BookmarkName: bookmark,
		BaseRef:      base,
		BaseCommitID: s.stackBaseCommitID(ctx, repo.RootPath, base),
		HeadChangeID: &change.ChangeID,
		HeadCommitID: &change.CommitID,
		Status:       "draft",
	}
	checkoutBranch = recordedEditCheckoutBranch(repo)
	repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, body.BookmarkName, checkoutBranch)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, body.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &body,
		OperationID: opID,
		Output:      output,
	}, nil
}

func (s *Service) resolveAppendExistingStack(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, bookmark string) (*storage.Stack, error) {
	bookmark = strings.TrimSpace(bookmark)
	if bookmark == "" {
		return nil, nil
	}
	existing, err := store.FindStackByBookmark(ctx, repoID, bookmark)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	jjBookmarkExists, err := s.jjBookmarkExists(ctx, repo.RootPath, bookmark)
	if err != nil {
		return nil, err
	}
	if !jjBookmarkExists {
		return nil, nil
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoID); err != nil {
		return nil, err
	}
	existing, err = store.FindStackByBookmark(ctx, repoID, bookmark)
	if err != nil || existing != nil {
		return existing, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, bookmark)
	if err != nil {
		return nil, fmt.Errorf("inspect existing bookmark %q: %w", bookmark, err)
	}
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	headChangeID := change.ChangeID
	headCommitID := change.CommitID
	now := time.Now().UnixMilli()
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         firstNonEmpty(stackNameFromBookmark(bookmark), bookmark),
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: s.stackBaseCommitID(ctx, repo.RootPath, baseRef),
		HeadChangeID: &headChangeID,
		HeadCommitID: &headCommitID,
		Status:       "draft",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}
	changeID, err := upsertChange(ctx, store, repoID, change)
	if err != nil {
		return nil, err
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, now); err != nil {
		return nil, err
	}
	return store.FindStackByBookmark(ctx, repoID, bookmark)
}

func (s *Service) commitUnlocked(ctx context.Context, repo RepoInfo, message string) (CommitResult, error) {
	if err := ValidateCommitMessage(message); err != nil {
		return CommitResult{}, err
	}
	body, err := s.resolveCurrentStackWithDescription(ctx, repo, true, message)
	if err != nil {
		return CommitResult{}, err
	}
	checkoutBranch := recordedEditCheckoutBranch(repo)
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", "commit", "-m", message)
	if err != nil {
		return CommitResult{}, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, body.BookmarkName, checkoutBranch)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, body.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &body,
		OperationID: opID,
		Output:      output,
	}, nil
}

func (s *Service) SplitCommit(ctx context.Context, opts SplitCommitOptions) (CommitResult, error) {
	if opts.Hunk {
		return s.splitCommitByHunkPatchLocked(ctx, opts)
	}
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var splitErr error
		result, splitErr = s.splitCommitUnlocked(ctx, repo, opts)
		return splitErr
	})
	return result, err
}

func (s *Service) splitCommitUnlocked(ctx context.Context, repo RepoInfo, opts SplitCommitOptions) (CommitResult, error) {
	if err := ValidateCommitMessage(opts.Message); err != nil {
		return CommitResult{}, err
	}
	body, err := s.resolveCurrentStackWithDescription(ctx, repo, true, opts.Message)
	if err != nil {
		return CommitResult{}, err
	}
	checkoutBranch := recordedEditCheckoutBranch(repo)

	args := []string{"split", "-m", opts.Message}
	if opts.Interactive {
		args = append(args, "--interactive")
	}
	args = append(args, opts.Filesets...)

	output := ""
	if opts.Interactive {
		err = s.runner.RunStream(ctx, repo.RootPath, "jj", args...)
	} else {
		output, err = s.runner.Run(ctx, repo.RootPath, "jj", args...)
	}
	if err != nil {
		return CommitResult{}, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, body.BookmarkName, checkoutBranch)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, body.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &body,
		OperationID: opID,
		Output:      output,
	}, nil
}

func (s *Service) splitCommitByHunkPatchLocked(ctx context.Context, opts SplitCommitOptions) (CommitResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	var result CommitResult
	err = withRepoLock(repo.RootPath, func() error {
		var splitErr error
		result, splitErr = s.SplitCommitByHunkPatch(ctx, opts)
		return splitErr
	})
	return result, err
}

func (s *Service) SplitCommitByHunkPatch(ctx context.Context, opts SplitCommitOptions) (CommitResult, error) {
	if strings.TrimSpace(opts.PatchFile) == "" {
		return CommitResult{}, fmt.Errorf("patch file is required; use --patch-file with --hunk")
	}
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	body, err := s.resolveCurrentStack(ctx, repo, true)
	if err != nil {
		return CommitResult{}, err
	}
	checkoutBranch := recordedEditCheckoutBranch(repo)
	current, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return CommitResult{}, err
	}
	if len(current.Files) == 0 {
		return CommitResult{}, fmt.Errorf("no changes available to split")
	}

	remainingPatch, err := s.remainingPatchFromSelection(ctx, repo.RootPath, opts.PatchFile)
	if err != nil {
		return CommitResult{}, err
	}
	remainingPatchPath := ""
	if strings.TrimSpace(remainingPatch) != "" {
		file, err := os.CreateTemp("", "gx-remaining-*.patch")
		if err != nil {
			return CommitResult{}, fmt.Errorf("create remaining patch temp file: %w", err)
		}
		if _, err := file.WriteString(remainingPatch); err != nil {
			_ = file.Close()
			return CommitResult{}, fmt.Errorf("write remaining patch temp file: %w", err)
		}
		if err := file.Close(); err != nil {
			return CommitResult{}, fmt.Errorf("close remaining patch temp file: %w", err)
		}
		remainingPatchPath = file.Name()
		defer os.Remove(remainingPatchPath)
	}

	args := append([]string{"split", "-m", opts.Message}, current.Files...)
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", args...)
	if err != nil {
		return CommitResult{}, err
	}
	childChange, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return CommitResult{}, err
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "edit", "@-"); err != nil {
		return CommitResult{}, err
	}
	if remainingPatchPath != "" {
		if _, err := s.runner.Run(ctx, repo.RootPath, "git", "apply", "-R", "--whitespace=nowarn", remainingPatchPath); err != nil {
			return CommitResult{}, err
		}
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "edit", childChange.ChangeID); err != nil {
		return CommitResult{}, err
	}
	if remainingPatchPath != "" {
		if _, err := s.runner.Run(ctx, repo.RootPath, "git", "apply", "--whitespace=nowarn", remainingPatchPath); err != nil {
			return CommitResult{}, err
		}
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, "@-")
	if err != nil {
		return CommitResult{}, err
	}
	repo, err = s.reattachRecordedContainer(ctx, repo.RootPath, body.BookmarkName, checkoutBranch)
	if err != nil {
		return CommitResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, body.BookmarkName); err != nil {
		return CommitResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return CommitResult{}, err
	}
	return CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &body,
		OperationID: opID,
		Output:      output,
	}, nil
}

func (s *Service) remainingPatchFromSelection(ctx context.Context, repoRoot, selectedPatchPath string) (string, error) {
	if _, err := os.Stat(selectedPatchPath); err != nil {
		return "", fmt.Errorf("patch file %q: %w", selectedPatchPath, err)
	}
	if _, err := s.runner.Run(ctx, repoRoot, "git", "apply", "-R", "--whitespace=nowarn", selectedPatchPath); err != nil {
		return "", err
	}
	remaining, diffErr := s.runner.RunStdout(ctx, repoRoot, "git", "diff", "--binary", "--full-index")
	if _, err := s.runner.Run(ctx, repoRoot, "git", "apply", "--whitespace=nowarn", selectedPatchPath); err != nil {
		return "", err
	}
	if diffErr != nil {
		return "", diffErr
	}
	return remaining, nil
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

func (s *Service) Modify(ctx context.Context, rev string) (ModifyResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return ModifyResult{}, err
	}
	var result ModifyResult
	err = withRepoLock(repo.RootPath, func() error {
		var modifyErr error
		result, modifyErr = s.modifyUnlocked(ctx, repo, rev)
		return modifyErr
	})
	return result, err
}

func (s *Service) RecordModification(ctx context.Context, rev string) (ModifyResult, error) {
	result, err := s.Modify(ctx, rev)
	if err != nil {
		return ModifyResult{}, err
	}
	if err := recordModify(ctx, result); err != nil {
		return result, fmt.Errorf("record revision modification metadata: %w", err)
	}
	return result, nil
}

func (s *Service) EditRevision(ctx context.Context, repoRoot, rev string) error {
	return s.editRevision(ctx, repoRoot, rev, true)
}

func (s *Service) EditWorkingCopyRevision(ctx context.Context, repoRoot, rev string) error {
	return s.editRevision(ctx, repoRoot, rev, false)
}

func (s *Service) editRevision(ctx context.Context, repoRoot, rev string, reattachEdit bool) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	if strings.TrimSpace(rev) == "" {
		return fmt.Errorf("revision is empty")
	}
	return withRepoLock(repoRoot, func() error {
		if _, err := s.runner.Run(ctx, repoRoot, "jj", "edit", rev); err != nil {
			return err
		}
		if !reattachEdit {
			return nil
		}
		repo, err := s.ResolveJJRepoAtPath(ctx, repoRoot)
		if err != nil {
			return err
		}
		bookmark := s.publicBookmarkForRev(ctx, repoRoot, "@", repo.defaultBaseBranch())
		if bookmark == "" || bookmark == repo.defaultBaseBranch() {
			return nil
		}
		commitID, err := s.commitIDForRev(ctx, repoRoot, bookmark)
		if err != nil {
			return err
		}
		return s.attachGitBranch(ctx, repoRoot, bookmark, commitID)
	})
}

func (s *Service) ReattachGitCheckoutRef(ctx context.Context, repoRoot, checkoutRef string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	checkoutRef = strings.TrimSpace(checkoutRef)
	if checkoutRef == "" {
		return nil
	}
	return withRepoLock(repoRoot, func() error {
		repo, err := s.configuredJJRepo(ctx)
		if err != nil {
			return err
		}
		_, err = s.reattachGitHeadToBaseRef(ctx, repo, checkoutRef)
		return err
	})
}

func (s *Service) ForceGitCheckout(ctx context.Context, repoRoot, checkoutRef string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	checkoutRef = strings.TrimSpace(checkoutRef)
	if checkoutRef == "" {
		return nil
	}
	checkoutRef = strings.TrimPrefix(checkoutRef, "refs/heads/")
	return withRepoLock(repoRoot, func() error {
		commitID, err := s.commitIDForContainer(ctx, repoRoot)
		if err != nil {
			return err
		}
		return s.attachGitBranch(ctx, repoRoot, checkoutRef, commitID)
	})
}

func (s *Service) ForceGitCheckoutPreservingWorktree(ctx context.Context, repoRoot, checkoutRef string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	checkoutRef = strings.TrimSpace(checkoutRef)
	if checkoutRef == "" {
		return nil
	}
	checkoutRef = strings.TrimPrefix(checkoutRef, "refs/heads/")
	return withRepoLock(repoRoot, func() error {
		_, err := s.runner.Run(ctx, repoRoot, "git", "symbolic-ref", "HEAD", "refs/heads/"+checkoutRef)
		return err
	})
}

func (s *Service) NewRevisionChild(ctx context.Context, repoRoot string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	return withRepoLock(repoRoot, func() error {
		_, err := s.runner.Run(ctx, repoRoot, "jj", "new", "@")
		return err
	})
}

func (s *Service) NewRevisionFrom(ctx context.Context, repoRoot, rev string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return fmt.Errorf("revision is empty")
	}
	return withRepoLock(repoRoot, func() error {
		_, err := s.runner.Run(ctx, repoRoot, "jj", "new", rev)
		return err
	})
}

func (s *Service) DeleteRevision(ctx context.Context, rev string) (DeleteRevisionResult, error) {
	repo, err := s.ResolveJJRepo(ctx)
	if err != nil {
		return DeleteRevisionResult{}, err
	}
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return DeleteRevisionResult{}, fmt.Errorf("revision is required")
	}
	var result DeleteRevisionResult
	err = withRepoLock(repo.RootPath, func() error {
		targetRev, err := s.resolveStoredRevisionRef(ctx, repo, rev)
		if err != nil {
			return err
		}
		change, err := s.CurrentChange(ctx, repo.RootPath, targetRev)
		if err != nil {
			return err
		}
		output, err := s.runner.Run(ctx, repo.RootPath, "jj", "abandon", targetRev)
		if err != nil {
			return err
		}
		if err := s.markChangeAbandoned(ctx, repo.RootPath, change.ChangeID); err != nil {
			return err
		}
		result = DeleteRevisionResult{Repo: repo, Change: change, Output: output}
		return nil
	})
	return result, err
}

func (s *Service) DeleteStack(ctx context.Context, bookmarkName string) (DeleteStackResult, error) {
	repo, err := s.ResolveJJRepo(ctx)
	if err != nil {
		return DeleteStackResult{}, err
	}
	bookmarkName = strings.TrimSpace(bookmarkName)
	if bookmarkName == "" {
		return DeleteStackResult{}, fmt.Errorf("stack is required")
	}
	var result DeleteStackResult
	err = withRepoLock(repo.RootPath, func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
		if err != nil {
			return err
		}
		if repoRow == nil {
			return fmt.Errorf("unknown repo %s", repo.RootPath)
		}
		storedStack, err := store.FindStackByBookmark(ctx, repoRow.ID, bookmarkName)
		if err != nil {
			return err
		}
		if storedStack == nil {
			return fmt.Errorf("unknown stack %q", bookmarkName)
		}
		changes, err := store.ListChangesByStackID(ctx, storedStack.ID)
		if err != nil {
			return err
		}
		var output strings.Builder
		for _, change := range changes {
			changeID := strings.TrimSpace(change.JJChangeID)
			if changeID == "" {
				continue
			}
			out, err := s.runner.Run(ctx, repo.RootPath, "jj", "abandon", changeID)
			if err != nil {
				return err
			}
			output.WriteString(out)
			if err := store.MarkChangeStatus(ctx, change.ID, "abandoned", time.Now().UnixMilli()); err != nil {
				return err
			}
		}
		out, err := s.runner.Run(ctx, repo.RootPath, "jj", "bookmark", "delete", storedStack.BookmarkName)
		if err != nil {
			return err
		}
		output.WriteString(out)
		if err := store.DeleteStack(ctx, storedStack.ID); err != nil {
			return err
		}
		result = DeleteStackResult{
			Repo:      repo,
			Stack:     stackInfoFromStorage(*storedStack),
			Revisions: changesToChangeInfo(changes),
			Output:    output.String(),
		}
		return nil
	})
	return result, err
}

func (s *Service) RepairWorkflow(ctx context.Context) (RepairResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return RepairResult{}, err
	}
	result := RepairResult{Repo: repo}
	err = withRepoLock(repo.RootPath, func() error {
		logicalBase := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
		authoringRef := gxAuthoringCheckoutRef(logicalBase)
		currentChange, changeErr := s.CurrentChange(ctx, repo.RootPath, "@")
		clean := changeErr == nil && len(currentChange.Files) == 0
		if clean {
			if _, err := s.setBaseUnlocked(ctx, repo, logicalBase); err != nil {
				return err
			}
			repo.AuthoringBase = &logicalBase
			reattached, err := s.reattachGitHeadToBaseRef(ctx, repo, logicalBase)
			if err != nil {
				return err
			}
			repo = reattached
			result.Repo = reattached
			result.Actions = append(result.Actions, fmt.Sprintf("returned clean checkout to %s", authoringRef))
		} else if s.currentGitCheckoutRef(ctx, repo.RootPath) == "" {
			if stack, found, stackErr := s.resolveCurrentStackForRead(ctx, repo); stackErr != nil {
				return stackErr
			} else if found && strings.TrimSpace(stack.BookmarkName) != "" {
				reattached, err := s.reattachContainer(ctx, repo.RootPath, stack.BookmarkName)
				if err != nil {
					return err
				}
				repo = reattached
				result.Repo = reattached
				result.Actions = append(result.Actions, fmt.Sprintf("reattached detached Git HEAD to %s because the working copy has changes", stack.BookmarkName))
			} else {
				result.Warnings = append(result.Warnings, "left detached Git HEAD with changes because no stack bookmark could be identified")
			}
		} else if changeErr == nil && len(currentChange.Files) > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("left checkout on %s because %d changed %s must be recorded or moved first", firstNonEmpty(s.currentGitCheckoutRef(ctx, repo.RootPath), "(detached)"), len(currentChange.Files), vcsPluralize("file", len(currentChange.Files))))
		}

		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		if repo.AuthoringBase != nil && legacyGXInternalCheckoutRef(*repo.AuthoringBase) {
			previousBase := strings.TrimSpace(*repo.AuthoringBase)
			if err := store.SetRepoAuthoringBase(ctx, repoID, logicalBase, time.Now().UnixMilli()); err != nil {
				return err
			}
			repo.AuthoringBase = &logicalBase
			result.Repo = repo
			result.Actions = append(result.Actions, fmt.Sprintf("changed GX authoring base from %s to %s", previousBase, logicalBase))
		}
		stacks, err := store.ListStacksByRepoID(ctx, repoID)
		if err != nil {
			return err
		}
		now := time.Now().UnixMilli()
		for _, stack := range stacks {
			updated := false
			if legacyGXInternalCheckoutRef(stack.BaseRef) {
				next := s.publicBookmarkForRev(ctx, repo.RootPath, firstNonEmpty(stack.BaseCommitID, stack.BaseRef), repo.defaultBaseBranch())
				if next == "" {
					next = s.publicStackBaseRef(ctx, repo, stack.BaseRef)
				}
				if next != stack.BaseRef {
					result.Actions = append(result.Actions, fmt.Sprintf("changed GX base ref for %s from %s to %s", stack.BookmarkName, stack.BaseRef, next))
					stack.BaseRef = next
					updated = true
				}
			}
			if stack.RemoteRef != nil && isInternalGitPublishedRef(*stack.RemoteRef) {
				next := "refs/heads/" + strings.TrimSpace(stack.BookmarkName)
				result.Actions = append(result.Actions, fmt.Sprintf("changed Git published ref for %s from %s to %s", stack.BookmarkName, *stack.RemoteRef, next))
				stack.RemoteRef = &next
				updated = true
			}
			if updated {
				stack.UpdatedAt = now
				if _, err := store.UpsertStack(ctx, stack); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return result, err
}

func isInternalGitPublishedRef(value string) bool {
	value = strings.TrimPrefix(strings.TrimSpace(value), "refs/heads/")
	return legacyGXInternalCheckoutRef(value)
}

func (s *Service) ApplyPatch(ctx context.Context, repoRoot, patchFile string, reverse bool) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	if strings.TrimSpace(patchFile) == "" {
		return fmt.Errorf("patch file is empty")
	}
	return withRepoLock(repoRoot, func() error {
		args := []string{"apply"}
		if reverse {
			args = append(args, "-R")
		}
		args = append(args, "--whitespace=nowarn", patchFile)
		_, err := s.runner.Run(ctx, repoRoot, "git", args...)
		return err
	})
}

func (s *Service) RestorePathsFromRevision(ctx context.Context, repoRoot, fromRev string, paths []string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root is empty")
	}
	fromRev = strings.TrimSpace(fromRev)
	if fromRev == "" {
		return fmt.Errorf("from revision is empty")
	}
	var cleaned []string
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		cleaned = append(cleaned, path)
	}
	if len(cleaned) == 0 {
		return nil
	}
	return withRepoLock(repoRoot, func() error {
		args := append([]string{"restore", "--from", fromRev}, cleaned...)
		if _, err := s.runner.Run(ctx, repoRoot, "jj", args...); err != nil {
			return fmt.Errorf("jj restore --from %s: %w", fromRev, err)
		}
		return nil
	})
}

func (s *Service) modifyUnlocked(ctx context.Context, repo RepoInfo, rev string) (ModifyResult, error) {
	if rev == "" {
		rev = "@"
	}
	targetRev, err := s.resolveStoredRevisionRef(ctx, repo, rev)
	if err != nil {
		return ModifyResult{}, err
	}
	body, found, err := s.stackForRevision(ctx, repo, targetRev)
	if err != nil {
		return ModifyResult{}, err
	}
	if !found {
		body, err = s.resolveCurrentStack(ctx, repo, true)
		if err != nil {
			return ModifyResult{}, err
		}
	}
	previous, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return ModifyResult{}, err
	}
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", "edit", targetRev)
	if err != nil {
		return ModifyResult{}, err
	}
	repo, err = s.updateContainerBookmark(ctx, repo.RootPath, body.BookmarkName)
	if err != nil {
		return ModifyResult{}, err
	}
	if err := s.reconcileRepoChanges(ctx, repo, body.BookmarkName); err != nil {
		return ModifyResult{}, err
	}
	repo, err = s.reattachContainer(ctx, repo.RootPath, body.BookmarkName)
	if err != nil {
		return ModifyResult{}, err
	}
	current, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return ModifyResult{}, err
	}
	opID, err := s.CurrentOperation(ctx, repo.RootPath)
	if err != nil {
		return ModifyResult{}, err
	}
	return ModifyResult{
		Repo:           repo,
		PreviousChange: previous,
		CurrentChange:  current,
		Stack:          &body,
		OperationID:    opID,
		Output:         output,
	}, nil
}

func (s *Service) resolveStoredRevisionRef(ctx context.Context, repo RepoInfo, rev string) (string, error) {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return rev, nil
	}
	if strings.ContainsAny(rev, "@():&| ") || strings.HasPrefix(rev, "root") {
		return rev, nil
	}
	store, err := openStore(ctx)
	if err != nil {
		return "", err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return "", err
	}
	if repoRow == nil {
		return rev, nil
	}
	changes, err := store.ListChangesByRepoID(ctx, repoRow.ID)
	if err != nil {
		return "", err
	}
	exact := make([]storage.Change, 0, 1)
	prefix := make([]storage.Change, 0, 1)
	for _, change := range changes {
		changeID := strings.TrimSpace(change.JJChangeID)
		commitID := strings.TrimSpace(change.CurrentCommitID)
		if changeID == rev || commitID == rev {
			exact = append(exact, change)
			continue
		}
		if strings.HasPrefix(changeID, rev) || strings.HasPrefix(commitID, rev) {
			prefix = append(prefix, change)
		}
	}
	matches := exact
	if len(matches) == 0 {
		matches = prefix
	}
	if len(matches) == 0 {
		return rev, nil
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("revision %q is ambiguous across %d stored GX revisions", rev, len(matches))
	}
	match := matches[0]
	return firstNonEmpty(strings.TrimSpace(match.CurrentCommitID), strings.TrimSpace(match.JJChangeID), rev), nil
}

func (s *Service) ListModifyCandidates(ctx context.Context, limit int) ([]ChangeInfo, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 20
	}
	const tmpl = `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "\n"`
	out, err := s.runStdoutTrimmed(ctx, repo.RootPath, "jj", "log", "-r", "mutable() & ~empty() & ~root()", "-n", strconv.Itoa(limit), "--no-graph", "-T", tmpl)
	if err != nil {
		return nil, err
	}
	lines := splitLines(out)
	candidates := make([]ChangeInfo, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
		desc := parts[2]
		if strings.TrimSpace(desc) == "" {
			desc = "(no description set)"
		}
		candidates = append(candidates, ChangeInfo{
			ChangeID:    parts[0],
			CommitID:    parts[1],
			Description: desc,
		})
	}
	return candidates, nil
}

func (s *Service) Switch(ctx context.Context, name string) (SwitchResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return SwitchResult{}, err
	}
	var result SwitchResult
	err = withRepoLock(repo.RootPath, func() error {
		var switchErr error
		result, switchErr = s.switchUnlocked(ctx, repo, name)
		return switchErr
	})
	return result, err
}

func (s *Service) CreateStack(ctx context.Context, name string) (CreateStackResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return CreateStackResult{}, err
	}
	var result CreateStackResult
	err = withRepoLock(repo.RootPath, func() error {
		var createErr error
		result, createErr = s.createStackUnlocked(ctx, repo, name)
		return createErr
	})
	return result, err
}

func (s *Service) Base(ctx context.Context) (BaseResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return BaseResult{}, err
	}
	return s.baseResult(ctx, repo), nil
}

func (s *Service) SetBase(ctx context.Context, baseRef string) (BaseResult, error) {
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		return BaseResult{}, fmt.Errorf("base ref is required")
	}
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return BaseResult{}, err
	}
	var result BaseResult
	err = withRepoLock(repo.RootPath, func() error {
		var setErr error
		result, setErr = s.setBaseUnlocked(ctx, repo, baseRef)
		return setErr
	})
	return result, err
}

func (s *Service) workingCopyParentMatchesRef(ctx context.Context, repoRoot string, refs ...string) bool {
	parentChange, err := s.changeIDForRev(ctx, repoRoot, "@-")
	if err != nil || strings.TrimSpace(parentChange) == "" {
		return false
	}
	seen := map[string]struct{}{}
	for _, ref := range refs {
		ref = strings.TrimSpace(ref)
		if ref == "" {
			continue
		}
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}
		refChange, err := s.changeIDForRev(ctx, repoRoot, ref)
		if err != nil || strings.TrimSpace(refChange) == "" {
			continue
		}
		if parentChange == refChange {
			return true
		}
	}
	return false
}

func (s *Service) authoringBaseParentRefs(ctx context.Context, repo RepoInfo, baseRef string) []string {
	baseRef = strings.TrimSpace(baseRef)
	refs := []string{
		gxAuthoringCheckoutRef(baseRef),
		baseRef,
		s.publicStackBaseRef(ctx, repo, baseRef),
		repo.defaultBaseBranch(),
		"main",
		"main@origin",
	}
	if resolved, err := s.resolveAuthoringBaseRef(ctx, repo.RootPath, baseRef); err == nil {
		refs = append(refs, resolved)
	}
	return refs
}

func (s *Service) workingCopyOnAuthoringBase(ctx context.Context, repo RepoInfo, baseRef string) bool {
	return s.workingCopyParentMatchesRef(ctx, repo.RootPath, s.authoringBaseParentRefs(ctx, repo, baseRef)...)
}

func (s *Service) reattachAuthoringBaseGitIfParent(ctx context.Context, repo RepoInfo, baseRef string) (RepoInfo, bool, error) {
	if !s.workingCopyOnAuthoringBase(ctx, repo, baseRef) {
		return repo, false, nil
	}
	reattached, err := s.reattachGitHeadToBaseRef(ctx, repo, baseRef)
	if err != nil {
		return repo, false, err
	}
	return reattached, true, nil
}

func (s *Service) setBaseUnlocked(ctx context.Context, repo RepoInfo, baseRef string) (BaseResult, error) {
	resolvedBaseRef, err := s.resolveAuthoringBaseRef(ctx, repo.RootPath, baseRef)
	if err != nil {
		return BaseResult{}, fmt.Errorf("resolve base ref %q: %w", baseRef, err)
	}
	targetRepo := repo
	targetRepo.AuthoringBase = &resolvedBaseRef
	target := s.baseResult(ctx, targetRepo)
	if !target.OnBase {
		change, changeErr := s.CurrentChange(ctx, repo.RootPath, "@")
		if changeErr == nil && len(change.Files) > 0 {
			if reattached, ok, reattachErr := s.reattachAuthoringBaseGitIfParent(ctx, repo, resolvedBaseRef); reattachErr != nil {
				return BaseResult{}, reattachErr
			} else if ok {
				repo = reattached
				target = s.baseResult(ctx, repo)
			}
		}
		if !target.OnBase {
			change, changeErr := s.CurrentChange(ctx, repo.RootPath, "@")
			if changeErr == nil && len(change.Files) > 0 {
				current := target.CurrentRef
				if current == "" {
					current = "(detached)"
				}
				return BaseResult{}, fmt.Errorf("gx base --set %s would switch from %q with %d changed %s; record or move the changes first", resolvedBaseRef, current, len(change.Files), vcsPluralize("file", len(change.Files)))
			}
		}
	}
	store, err := openStore(ctx)
	if err != nil {
		return BaseResult{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return BaseResult{}, err
	}
	if err := store.SetRepoAuthoringBase(ctx, repoID, resolvedBaseRef, time.Now().UnixMilli()); err != nil {
		return BaseResult{}, err
	}
	repo.AuthoringBase = &resolvedBaseRef
	if target.OnBase {
		if strings.TrimSpace(target.CurrentRef) != baseCheckoutRef(resolvedBaseRef) {
			reattached, err := s.reattachGitHeadToBaseRef(ctx, repo, resolvedBaseRef)
			if err != nil {
				return BaseResult{}, err
			}
			repo = reattached
		}
		return s.baseResult(ctx, repo), nil
	}
	switched, err := s.switchBaseUnlocked(ctx, repo, resolvedBaseRef)
	if err != nil {
		return BaseResult{}, err
	}
	return s.baseResult(ctx, switched.Repo), nil
}

func (s *Service) resolveAuthoringBaseRef(ctx context.Context, repoRoot, baseRef string) (string, error) {
	baseRef = strings.TrimPrefix(strings.TrimSpace(baseRef), "origin/")
	baseRef = strings.TrimPrefix(baseRef, "refs/heads/")
	if baseRef == "" {
		baseRef = "main"
	}
	if _, err := s.commitIDForRev(ctx, repoRoot, baseRef); err != nil {
		return "", err
	}
	return baseRef, nil
}

func (s *Service) RequireAuthoringBase(ctx context.Context, commandName string) error {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return err
	}
	_, _, err = s.ensureAuthoringCheckout(ctx, repo, commandName)
	return err
}

func (s *Service) ensureAuthoringCheckout(ctx context.Context, repo RepoInfo, commandName string) (string, string, error) {
	result := s.baseResult(ctx, repo)
	current := strings.TrimSpace(result.CurrentRef)
	if current == "" {
		current = "(detached)"
	}
	if result.OnBase {
		return result.BaseRef, "base", nil
	}
	change, changeErr := s.CurrentChange(ctx, repo.RootPath, "@")
	if changeErr == nil && len(change.Files) > 0 {
		if _, ok, reattachErr := s.reattachAuthoringBaseGitIfParent(ctx, repo, result.BaseRef); reattachErr != nil {
			return "", "", reattachErr
		} else if ok {
			return result.BaseRef, "base", nil
		}
		if s.workingCopyOnAuthoringBase(ctx, repo, result.BaseRef) {
			return result.BaseRef, "base", nil
		}
		checkoutRef := gxAuthoringCheckoutRef(result.BaseRef)
		return "", "", fmt.Errorf("%s must run from %s; currently on %q with %d changed %s. Rebase onto %s with `jj rebase -s @ -d %s`, then run `gx base --set %s`", commandName, checkoutRef, current, len(change.Files), vcsPluralize("file", len(change.Files)), checkoutRef, checkoutRef, result.BaseRef)
	}
	if changeErr == nil && len(change.Files) == 0 {
		if _, err := s.switchBaseUnlocked(ctx, repo, result.BaseRef); err == nil {
			return result.BaseRef, "base", nil
		}
	}
	return "", "", fmt.Errorf("%s must run from %s; currently on %q. Run `gx base --set %s` to return to the GX authoring checkout", commandName, gxAuthoringCheckoutRef(result.BaseRef), current, result.BaseRef)
}

func vcsPluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func (s *Service) createStackUnlocked(ctx context.Context, repo RepoInfo, name string) (CreateStackResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CreateStackResult{}, fmt.Errorf("stack name is required")
	}
	store, err := openStore(ctx)
	if err != nil {
		return CreateStackResult{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return CreateStackResult{}, err
	}
	current, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return CreateStackResult{}, err
	}
	if len(current.Files) > 0 {
		return CreateStackResult{}, fmt.Errorf("record or discard current changes before creating a stack")
	}
	baseRef := s.defaultStackBaseRef(repo)
	if currentStack, found, err := s.resolveCurrentStackForRead(ctx, repo); err != nil {
		return CreateStackResult{}, err
	} else if found && strings.TrimSpace(currentStack.BookmarkName) != "" {
		baseRef = currentStack.BookmarkName
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "new", "@"); err != nil {
		return CreateStackResult{}, err
	}
	current, err = s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return CreateStackResult{}, err
	}
	bookmark := stackBookmarkName(name, current.ChangeID)
	if existing, err := store.FindStackByBookmark(ctx, repoID, bookmark); err != nil {
		return CreateStackResult{}, err
	} else if existing != nil {
		return CreateStackResult{}, fmt.Errorf("GX stack %q already exists", bookmark)
	}
	if jjBookmarkExists, err := s.jjBookmarkExists(ctx, repo.RootPath, bookmark); err != nil {
		return CreateStackResult{}, err
	} else if jjBookmarkExists {
		return CreateStackResult{}, fmt.Errorf("JJ bookmark %q already exists", bookmark)
	}
	if err := s.assertBookmarkTargetAvailable(ctx, store, repoID, repo.RootPath, bookmark, "@"); err != nil {
		return CreateStackResult{}, err
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "bookmark", "set", bookmark, "-r", "@"); err != nil {
		return CreateStackResult{}, err
	}
	now := time.Now().UnixMilli()
	headChangeID := current.ChangeID
	headCommitID := current.CommitID
	baseCommitID := s.stackBaseCommitID(ctx, repo.RootPath, baseRef)
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         name,
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommitID,
		HeadChangeID: &headChangeID,
		HeadCommitID: &headCommitID,
		Status:       "draft",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return CreateStackResult{}, err
	}
	stack := storage.Stack{
		ID:           stackID,
		RepoID:       repoID,
		Name:         name,
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommitID,
		HeadChangeID: &headChangeID,
		HeadCommitID: &headCommitID,
		Status:       "draft",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	repo, err = s.reattachContainer(ctx, repo.RootPath, bookmark)
	if err != nil {
		return CreateStackResult{}, err
	}
	current, err = s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return CreateStackResult{}, err
	}
	return CreateStackResult{
		Repo:          repo,
		Stack:         stackInfoFromStorage(stack),
		CurrentChange: current,
	}, nil
}

func (s *Service) BaseSwitchTarget(ctx context.Context, name string) (string, bool, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return "", false, err
	}
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultSwitchBaseRef(ctx, repo))
	name = strings.TrimSpace(name)
	checkoutRef := baseCheckoutRef(baseRef)
	ok := name == "base" || name == baseRef || name == repo.defaultBaseBranch() || name == checkoutRef
	return baseRef, ok, nil
}

func (s *Service) switchUnlocked(ctx context.Context, repo RepoInfo, name string) (SwitchResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return SwitchResult{}, fmt.Errorf("stack name is required")
	}
	baseRef := s.defaultSwitchBaseRef(ctx, repo)
	if name == "base" || name == baseRef || name == repo.defaultBaseBranch() || name == baseCheckoutRef(baseRef) {
		return s.switchBaseUnlocked(ctx, repo, baseRef)
	}
	store, err := openStore(ctx)
	if err != nil {
		return SwitchResult{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return SwitchResult{}, err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return SwitchResult{}, err
	}
	var stack *storage.Stack
	for _, candidate := range stacks {
		if stackMatchesDirectName(candidate, name) {
			copy := candidate
			stack = &copy
			break
		}
	}
	if stack == nil {
		for index, candidate := range stacks {
			if stackAlias(index) == name {
				copy := candidate
				stack = &copy
				break
			}
		}
	}
	if stack == nil {
		return SwitchResult{}, fmt.Errorf("unknown GX stack %q", name)
	}
	output, err := s.editStackWorkingCopy(ctx, repo.RootPath, *stack)
	if err != nil {
		return SwitchResult{}, err
	}
	repo, err = s.reattachContainer(ctx, repo.RootPath, stack.BookmarkName)
	if err != nil {
		return SwitchResult{}, err
	}
	current, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return SwitchResult{}, err
	}
	return SwitchResult{
		Repo:          repo,
		Stack:         stackInfoFromStorage(*stack),
		CurrentChange: current,
		Output:        output,
	}, nil
}

func (s *Service) switchBaseUnlocked(ctx context.Context, repo RepoInfo, baseRef string) (SwitchResult, error) {
	resolvedBase, err := s.resolveAuthoringBaseRef(ctx, repo.RootPath, baseRef)
	if err != nil {
		return SwitchResult{}, err
	}
	if s.workingCopyOnAuthoringBase(ctx, repo, resolvedBase) {
		reattached, err := s.reattachGitHeadToBaseRef(ctx, repo, resolvedBase)
		if err != nil {
			return SwitchResult{}, err
		}
		repo = reattached
		current, err := s.CurrentChange(ctx, repo.RootPath, "@")
		if err != nil {
			return SwitchResult{}, err
		}
		commitID, err := s.commitIDForRev(ctx, repo.RootPath, resolvedBase)
		if err != nil {
			return SwitchResult{}, err
		}
		return SwitchResult{
			Repo:          repo,
			Stack:         StackInfo{Name: resolvedBase, BookmarkName: resolvedBase, BaseRef: resolvedBase, BaseCommitID: commitID},
			CurrentChange: current,
		}, nil
	}
	output, err := s.runner.Run(ctx, repo.RootPath, "jj", "new", resolvedBase)
	if err != nil {
		return SwitchResult{}, err
	}
	commitID, err := s.commitIDForRev(ctx, repo.RootPath, resolvedBase)
	if err != nil {
		return SwitchResult{}, err
	}
	reattachedRepo, err := s.reattachGitHeadToBaseRef(ctx, repo, resolvedBase)
	if err != nil {
		return SwitchResult{}, err
	}
	repo = reattachedRepo
	current, err := s.CurrentChange(ctx, repo.RootPath, "@")
	if err != nil {
		return SwitchResult{}, err
	}
	return SwitchResult{
		Repo:          repo,
		Stack:         StackInfo{Name: resolvedBase, BookmarkName: resolvedBase, BaseRef: resolvedBase, BaseCommitID: commitID},
		CurrentChange: current,
		Output:        output,
	}, nil
}

func PromptForModifySelection(in io.Reader, out io.Writer, candidates []ChangeInfo) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("no mutable changes available to edit")
	}
	fmt.Fprintln(out, "Select change to edit:")
	for i, candidate := range candidates {
		fmt.Fprintf(out, "  %d. %s  [%s %s]\n", i+1, candidate.Description, shortID(candidate.ChangeID, 12), shortID(candidate.CommitID, 8))
	}
	fmt.Fprintf(out, "Enter selection [1-%d, q]: ", len(candidates))

	reader := bufio.NewReader(in)
	raw, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	choice := strings.TrimSpace(raw)
	if choice == "" {
		return candidates[0].ChangeID, nil
	}
	if choice == "q" || choice == "quit" || choice == "exit" {
		return "", context.Canceled
	}
	if index, err := strconv.Atoi(choice); err == nil {
		if index >= 1 && index <= len(candidates) {
			return candidates[index-1].ChangeID, nil
		}
	}
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate.ChangeID, choice) || strings.HasPrefix(candidate.CommitID, choice) {
			return candidate.ChangeID, nil
		}
	}
	if _, err := strconv.Atoi(choice); err == nil {
		return "", fmt.Errorf("selection %s out of range", choice)
	}
	return "", fmt.Errorf("unknown selection %q", choice)
}

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

func (s *Service) PushWithOptions(ctx context.Context, args []string, opts PushOptions) (PushResult, error) {
	repo, err := s.ResolveGitRepo(ctx)
	if err != nil {
		return PushResult{}, err
	}
	lockRoot := repo.RootPath
	if repo.Backend == "jj" || hasJJRepo(ctx, s, repo.RootPath) {
		if jjRepo, jjErr := s.ResolveJJRepoAtPath(ctx, repo.RootPath); jjErr == nil {
			lockRoot = jjRepo.RootPath
		}
	}
	var result PushResult
	err = withRepoLock(lockRoot, func() error {
		var pushErr error
		result, pushErr = s.pushUnlocked(ctx, args, opts)
		return pushErr
	})
	return result, err
}

func (s *Service) PreparePublish(ctx context.Context, args []string, opts PushOptions) (PushResult, error) {
	return s.PushWithOptions(ctx, args, opts)
}

func (s *Service) PrepareNamedPublish(ctx context.Context, name string, args []string, opts PushOptions) (PushResult, error) {
	repo, err := s.ResolveGitRepo(ctx)
	if err != nil {
		return PushResult{}, err
	}
	lockRoot := repo.RootPath
	if repo.Backend == "jj" || hasJJRepo(ctx, s, repo.RootPath) {
		if jjRepo, jjErr := s.ResolveJJRepoAtPath(ctx, repo.RootPath); jjErr == nil {
			repo = jjRepo
			lockRoot = jjRepo.RootPath
		}
	}
	var result PushResult
	err = withRepoLock(lockRoot, func() error {
		stack, err := s.stackBySelector(ctx, repo, name)
		if err != nil {
			return err
		}
		var pushErr error
		result, pushErr = s.pushStackUnlocked(ctx, repo, args, opts, stack)
		return pushErr
	})
	return result, err
}

func (s *Service) PrepareAllPublishes(ctx context.Context, args []string, opts PushOptions) ([]PushResult, error) {
	return s.PublishAllStacks(ctx, args, opts, nil)
}

func (s *Service) RecordPublish(ctx context.Context, result PushResult) error {
	if err := recordPush(ctx, result); err != nil {
		return fmt.Errorf("record stack publication metadata: %w", err)
	}
	return nil
}

func (s *Service) PublishStack(ctx context.Context, args []string, opts PushOptions, hook PublishHook) (PushResult, error) {
	result, err := s.PushWithOptions(ctx, args, opts)
	if err != nil {
		return PushResult{}, err
	}
	return result, publishResult(ctx, result, hook)
}

func (s *Service) PublishNamedStack(ctx context.Context, name string, args []string, opts PushOptions, hook PublishHook) (PushResult, error) {
	repo, err := s.ResolveGitRepo(ctx)
	if err != nil {
		return PushResult{}, err
	}
	lockRoot := repo.RootPath
	if repo.Backend == "jj" || hasJJRepo(ctx, s, repo.RootPath) {
		if jjRepo, jjErr := s.ResolveJJRepoAtPath(ctx, repo.RootPath); jjErr == nil {
			repo = jjRepo
			lockRoot = jjRepo.RootPath
		}
	}
	var result PushResult
	err = withRepoLock(lockRoot, func() error {
		stack, err := s.stackBySelector(ctx, repo, name)
		if err != nil {
			return err
		}
		var pushErr error
		result, pushErr = s.pushStackUnlocked(ctx, repo, args, opts, stack)
		return pushErr
	})
	if err != nil {
		return PushResult{}, err
	}
	return result, publishResult(ctx, result, hook)
}

func (s *Service) PublishAllStacks(ctx context.Context, args []string, opts PushOptions, hook PublishHook) ([]PushResult, error) {
	repo, err := s.ResolveGitRepo(ctx)
	if err != nil {
		return nil, err
	}
	lockRoot := repo.RootPath
	if repo.Backend == "jj" || hasJJRepo(ctx, s, repo.RootPath) {
		if jjRepo, jjErr := s.ResolveJJRepoAtPath(ctx, repo.RootPath); jjErr == nil {
			repo = jjRepo
			lockRoot = jjRepo.RootPath
		}
	}
	stacks, err := s.ListStacks(ctx)
	if err != nil {
		return nil, err
	}
	if len(stacks) == 0 {
		return []PushResult{}, nil
	}
	results := make([]PushResult, 0, len(stacks))
	var firstErr error
	err = withRepoLock(lockRoot, func() error {
		for _, stack := range stacks {
			publishable, err := s.stackHasUnpublishedRevisions(ctx, stack)
			if err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: %w", stackPublishLabel(stack), err)
				}
				continue
			}
			if !publishable {
				continue
			}
			result, err := s.pushStackUnlocked(ctx, repo, args, opts, stack)
			if err == nil {
				err = publishResult(ctx, result, hook)
			}
			if err != nil {
				if errors.Is(err, ErrNoRecordedAdds) {
					continue
				}
				if firstErr == nil {
					firstErr = fmt.Errorf("%s: %w", stackPublishLabel(stack), err)
				}
				continue
			}
			results = append(results, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 && firstErr != nil {
		return nil, firstErr
	}
	if len(results) == 0 {
		return []PushResult{}, nil
	}
	return results, firstErr
}

func (s *Service) stackHasUnpublishedRevisions(ctx context.Context, stack StackInfo) (bool, error) {
	revisions := stack.Revisions
	if revisions == nil {
		var err error
		revisions, err = s.revisionsForStack(ctx, stack)
		if err != nil {
			return false, err
		}
	}
	for _, revision := range revisions {
		if !revision.Published {
			return true, nil
		}
	}
	return false, nil
}

func publishResult(ctx context.Context, result PushResult, hook PublishHook) error {
	if err := recordPush(ctx, result); err != nil {
		return fmt.Errorf("record stack publication metadata: %w", err)
	}
	if hook != nil {
		if err := hook(result); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ListStacks(ctx context.Context) ([]StackInfo, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return nil, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return nil, err
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoID); err != nil {
		return nil, err
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoID); err != nil {
		return nil, err
	}
	stored, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}
	return s.hydrateStoredStacks(ctx, store, repo, repoID, stored)
}

func (s *Service) stackBySelector(ctx context.Context, repo RepoInfo, selector string) (StackInfo, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return StackInfo{}, fmt.Errorf("stack name is required")
	}
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, err
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoID); err != nil {
		return StackInfo{}, err
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoID); err != nil {
		return StackInfo{}, err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return StackInfo{}, err
	}
	for _, stack := range stacks {
		if stackMatchesDirectName(stack, selector) || stackNameFromBookmark(stack.BookmarkName) == selector {
			return stackInfoFromStorage(stack), nil
		}
	}
	for index, stack := range stacks {
		if stackAlias(index) == selector {
			info := stackInfoFromStorage(stack)
			info.Alias = stackAlias(index)
			return info, nil
		}
	}
	return StackInfo{}, fmt.Errorf("unknown GX stack %q", selector)
}

func stackPublishLabel(stack StackInfo) string {
	if name := strings.TrimSpace(stack.Name); name != "" {
		return name
	}
	if name := stackNameFromBookmark(stack.BookmarkName); name != "" {
		return name
	}
	if bookmark := strings.TrimSpace(stack.BookmarkName); bookmark != "" {
		return bookmark
	}
	return "stack"
}

func (s *Service) pushUnlocked(ctx context.Context, args []string, opts PushOptions) (PushResult, error) {
	repo, err := s.ResolveJJRepo(ctx)
	if err != nil {
		return PushResult{}, err
	}
	body, err := s.resolveCurrentStackForPublish(ctx, repo)
	if err != nil {
		return PushResult{}, err
	}
	return s.pushStackUnlocked(ctx, repo, args, opts, body)
}

func (s *Service) pushStackUnlocked(ctx context.Context, repo RepoInfo, args []string, opts PushOptions, body StackInfo) (PushResult, error) {
	remoteName := repo.DefaultRemote
	if remoteName == nil || strings.TrimSpace(*remoteName) == "" {
		remoteName = ptr("origin")
	}
	if remoteURL, err := s.remoteURLForTarget(ctx, repo.RootPath, *remoteName); err == nil && remoteURL != "" {
		repo.RemoteURL = &remoteURL
	}
	if err := s.requireRecordedAddsForPublish(ctx, body); err != nil {
		return PushResult{}, err
	}
	refName := publishRefFromArgs(args)
	if refName == "" {
		refName = publishRefForStack(body)
	}
	repo.BranchName = &refName
	gxBaseRef := s.publicStackBaseRef(ctx, repo, body.BaseRef)

	gitExported := opts.gitExportEnabled()
	pushed, gitPushStatus, warnings, stackErr := s.pushRecordedStack(ctx, repo, *remoteName, refName, body, opts)
	output := ""
	if stackErr != nil {
		return PushResult{}, stackErr
	}
	var githubPRURL *string
	githubPRStatus := "not requested"
	if gitExported {
		prURL, prStatus, prWarnings := s.ensureGitHubPullRequest(ctx, repo, body, refName, pushed)
		githubPRURL = prURL
		githubPRStatus = prStatus
		warnings = append(warnings, prWarnings...)
		body.GitHubPRURL = prURL
		for i := range pushed {
			pushed[i].GitHubPullRequestURL = prURL
		}
	}
	if len(pushed) > 0 {
		output = stackPushOutput(pushed, gitExported, gitPushStatus)
	}

	head, err := s.stackHeadCommitID(ctx, repo.RootPath, body, pushed)
	if err != nil {
		return PushResult{}, err
	}
	var reattached RepoInfo
	reattached, err = s.reattachGitHeadToBaseRef(ctx, repo, gxBaseRef)
	if err != nil {
		return PushResult{}, err
	}
	repo = reattached
	var change *ChangeInfo
	if len(pushed) > 0 {
		current := pushed[len(pushed)-1].Change
		change = &current
	}
	return PushResult{
		Repo:                 repo,
		CurrentChange:        change,
		Stack:                &body,
		Published:            pushed,
		GitExported:          gitExported,
		GitPushStatus:        gitPushStatus,
		GitHubPRStatus:       githubPRStatus,
		HeadCommitID:         head,
		RemoteName:           remoteName,
		GXStackRef:           refName,
		GXBaseRef:            gxBaseRef,
		GitPublishedRef:      "refs/heads/" + refName,
		GitCheckoutRef:       s.currentGitCheckoutRef(ctx, repo.RootPath),
		Output:               output,
		Warnings:             warnings,
		GitHubPullRequestURL: githubPRURL,
	}, nil
}

func (s *Service) stackHeadCommitID(ctx context.Context, repoRoot string, stack StackInfo, pushed []PushedChange) (string, error) {
	headRev := strings.TrimSpace(stack.BookmarkName)
	if headRev == "" && stack.HeadChangeID != nil {
		headRev = strings.TrimSpace(*stack.HeadChangeID)
	}
	if headRev != "" {
		return s.commitIDForRev(ctx, repoRoot, headRev)
	}
	if len(pushed) > 0 {
		return pushed[len(pushed)-1].Change.CommitID, nil
	}
	return "", fmt.Errorf("cannot resolve stack head for publish")
}

func (s *Service) pushRecordedStack(ctx context.Context, repo RepoInfo, remoteName, refName string, stack StackInfo, opts PushOptions) ([]PushedChange, string, []string, error) {
	gitExported := opts.gitExportEnabled()
	revisions := stack.Revisions
	if revisions == nil {
		var err error
		revisions, err = s.revisionsForStack(ctx, stack)
		if err != nil {
			return nil, "", nil, err
		}
	}
	if len(revisions) == 0 {
		return nil, "", nil, fmt.Errorf(
			"no revisions on stack to publish; run `gx add -m \"...\"` on the stack before `gx publish`",
		)
	}
	headRev := strings.TrimSpace(stack.BookmarkName)
	if headRev == "" && stack.HeadChangeID != nil {
		headRev = strings.TrimSpace(*stack.HeadChangeID)
	}
	if headRev == "" {
		headRev = revisions[len(revisions)-1].CommitID
	}
	headCommitID, err := s.commitIDForRev(ctx, repo.RootPath, headRev)
	if err != nil {
		return nil, "", nil, err
	}
	if err := s.runGitLocked(ctx, repo.RootPath, "update stack ref ref", func() error {
		return s.setGitBranchRef(ctx, repo.RootPath, refName, headCommitID)
	}); err != nil {
		return nil, "", nil, err
	}
	gitPushStatus := "not pushed"
	warnings := []string{}
	if gitExported {
		remoteHead, remoteErr := s.remoteBranchHead(ctx, repo.RootPath, remoteName, refName)
		if remoteErr == nil && remoteHead == headCommitID {
			gitPushStatus = "already up to date"
			s.progressf("GitHub branch %s already up to date on %s.", refName, remoteName)
		} else {
			if remoteErr != nil {
				warnings = append(warnings, fmt.Sprintf("Could not check whether %s/%s already exists before pushing: %v", remoteName, refName, remoteErr))
			}
			if err := s.runGitLocked(ctx, repo.RootPath, "push stack ref", func() error {
				s.progressf("Pushing %s to %s...", refName, remoteName)
				return s.runner.RunStream(ctx, repo.RootPath, "git", "push", remoteName, refName)
			}); err != nil {
				return nil, "", nil, err
			}
			gitPushStatus = "pushed"
		}
	}

	pushed := make([]PushedChange, 0, len(revisions))
	for _, unit := range revisions {
		patch, err := s.runStdoutTrimmed(ctx, repo.RootPath, "jj", "diff", "-r", unit.CommitID, "--git")
		if err != nil {
			return nil, "", nil, err
		}
		pushed = append(pushed, PushedChange{
			Change: ChangeInfo{
				ChangeID:       unit.ChangeID,
				CommitID:       unit.CommitID,
				Description:    unit.Description,
				ParentChangeID: nil,
			},
			BranchName:     refName,
			BaseBranchName: s.publicStackBaseRef(ctx, repo, stack.BaseRef),
			Patch:          patch,
		})
	}
	if gitExported {
		if err := recordStackBookmarks(ctx, repo, remoteName, pushed); err != nil {
			return nil, "", nil, err
		}
	}
	return pushed, gitPushStatus, warnings, nil
}

func (s *Service) ensureGitHubPullRequest(ctx context.Context, repo RepoInfo, stack StackInfo, refName string, pushed []PushedChange) (*string, string, []string) {
	if stack.GitHubPRURL != nil && strings.TrimSpace(*stack.GitHubPRURL) != "" {
		return ptr(strings.TrimSpace(*stack.GitHubPRURL)), "stored", nil
	}
	if repo.RemoteURL == nil || strings.TrimSpace(*repo.RemoteURL) == "" {
		return nil, "skipped", nil
	}
	host, owner, repoName, ok := parseGitHubRemote(*repo.RemoteURL)
	if !ok {
		return nil, "skipped", nil
	}
	baseRef := s.publicStackBaseRef(ctx, repo, stack.BaseRef)
	title := strings.TrimSpace(firstNonEmpty(stack.Name, stack.BookmarkName, refName))
	body := githubPullRequestBody(pushed)

	client, err := githubapi.NewClient(host)
	if err != nil {
		return nil, "warning", []string{fmt.Sprintf("Could not prepare GitHub PR for branch %s: %v. Re-run `gx auth login` or set GH_TOKEN/GITHUB_TOKEN, then retry `gx publish %s`.", refName, err, refName)}
	}
	opts := githubapi.CreatePullRequestOptions{
		Host:       host,
		Owner:      owner,
		Repo:       repoName,
		BaseBranch: baseRef,
		HeadBranch: refName,
		Title:      title,
		Body:       body,
	}
	existing, err := client.FindPullRequest(ctx, opts)
	if err != nil {
		return nil, "warning", []string{fmt.Sprintf("Could not check for an existing GitHub PR for target branch %s (the published branch): %v. Retry `gx publish %s` after GitHub access is fixed.", refName, err, refName)}
	}
	if existing != nil && strings.TrimSpace(existing.URL) != "" {
		return ptr(strings.TrimSpace(existing.URL)), "existing", nil
	}
	remoteName := "origin"
	if repo.DefaultRemote != nil && strings.TrimSpace(*repo.DefaultRemote) != "" {
		remoteName = strings.TrimSpace(*repo.DefaultRemote)
	}
	baseExists, baseErr := s.remoteBranchExists(ctx, repo.RootPath, remoteName, baseRef)
	if baseErr == nil && !baseExists {
		compare := githubPullRequestURL(*repo.RemoteURL, baseRef, refName)
		hint := fmt.Sprintf("Could not create GitHub PR for target branch %s because its base branch %q is not on %s. ", refName, baseRef, remoteName)
		if isGXStackBookmark(baseRef) {
			hint += fmt.Sprintf("Run `gx publish %s` first, then retry `gx publish %s`.", baseRef, refName)
		} else {
			hint += fmt.Sprintf("Set a valid base with `gx base --set %s` or run `gx repair`, then retry `gx publish %s`.", baseRef, refName)
		}
		if compare != nil {
			hint += " You can also inspect " + *compare + "."
		}
		return nil, "warning", []string{hint}
	}
	warnings := []string{}
	if baseErr != nil {
		warnings = append(warnings, fmt.Sprintf("Could not verify PR base branch %q on %s before creating a PR: %v", baseRef, remoteName, baseErr))
	}
	created, err := client.CreatePullRequest(ctx, opts)
	if err != nil {
		return nil, "warning", append(warnings, fmt.Sprintf("Could not create GitHub PR for target branch %s against base %s: %v. Fix the base with `gx base --set <github-base-branch>` or create the PR manually.", refName, baseRef, err))
	}
	if created == nil || strings.TrimSpace(created.URL) == "" {
		return nil, "warning", append(warnings, fmt.Sprintf("GitHub did not return a PR URL for target branch %s. Retry `gx publish %s` or create the PR manually.", refName, refName))
	}
	return ptr(strings.TrimSpace(created.URL)), "created", warnings
}

func (s *Service) remoteBranchExists(ctx context.Context, repoRoot, remoteName, branchName string) (bool, error) {
	head, err := s.remoteBranchHead(ctx, repoRoot, remoteName, branchName)
	if err != nil {
		return false, err
	}
	return head != "", nil
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

func githubPullRequestBody(_ []PushedChange) string {
	return ""
}

func (s *Service) revisionsForStack(ctx context.Context, stack StackInfo) ([]RevisionSummary, error) {
	if stack.ID == 0 {
		return nil, nil
	}
	store, err := openStore(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	changes, err := store.ListChangesByStackID(ctx, stack.ID)
	if err != nil {
		return nil, err
	}
	bookmarks, err := store.ListChangeBookmarksByName(ctx, publishRefForStack(stack))
	if err != nil {
		return nil, err
	}
	lastPushedByChangeID := make(map[int64]string, len(bookmarks))
	for _, bookmark := range bookmarks {
		lastPushedByChangeID[bookmark.ChangeID] = bookmark.LastPushedCommitID
	}
	var publishedThrough *int64
	if len(lastPushedByChangeID) == 0 {
		repo, err := s.ResolveJJRepo(ctx)
		if err != nil {
			return nil, err
		}
		repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
		if err != nil {
			return nil, err
		}
		if repoRow != nil {
			push, err := store.LatestPushByBranchName(ctx, repoRow.ID, publishRefForStack(stack))
			if err != nil {
				return nil, err
			}
			if push != nil && push.CurrentChangeID != nil {
				publishedThrough = push.CurrentChangeID
			}
		}
	}
	publishedState := NewStackPublicationState(-1, lastPushedByChangeID)
	if publishedThrough != nil {
		publishedState.PublishedThrough = *publishedThrough
	}
	revisions := make([]RevisionSummary, 0, len(changes))
	for _, change := range changes {
		if IsPlaceholderDescription(change.Description) || change.Status == "abandoned" {
			continue
		}
		revisions = append(revisions, RevisionSummary{
			Index:       len(revisions) + 1,
			ChangeID:    change.JJChangeID,
			CommitID:    change.CurrentCommitID,
			Description: change.Description,
			Status:      change.Status,
			Published:   publishedState.RevisionPublished(change),
		})
	}
	return revisions, nil
}

func (s *Service) Revisions(ctx context.Context) ([]RevisionSummary, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return nil, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return nil, err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return nil, err
	}
	if repoRow == nil {
		return nil, nil
	}
	changes, err := store.ListChangesByRepoID(ctx, repoRow.ID)
	if err != nil {
		return nil, err
	}
	activeChangeID, err := s.currentWorkChangeID(ctx, repo.RootPath)
	if err != nil {
		activeChangeID = ""
	}
	latestPush, err := store.LatestPushByRepoID(ctx, repoRow.ID)
	if err != nil {
		return nil, err
	}

	publishedThrough := int64(-1)
	if latestPush != nil && latestPush.CurrentChangeID != nil {
		publishedThrough = *latestPush.CurrentChangeID
	}

	summaries := make([]RevisionSummary, 0, len(changes))
	for i, change := range changes {
		summaries = append(summaries, RevisionSummary{
			Index:       i + 1,
			ChangeID:    change.JJChangeID,
			CommitID:    change.CurrentCommitID,
			Description: change.Description,
			Status:      change.Status,
			Active:      activeChangeID != "" && change.JJChangeID == activeChangeID,
			Published:   ChangePublishedThrough(change, publishedThrough),
		})
	}
	return summaries, nil
}

func (s *Service) Stack(ctx context.Context) (StackSummary, error) {
	repo, err := s.ResolveJJRepo(ctx)
	if err != nil {
		return StackSummary{}, err
	}
	model, err := s.loadStackReadModel(ctx, repo)
	if err != nil {
		return StackSummary{}, err
	}
	return stackSummaryForStack(model.repo, model.currentStack, model.currentStackFound, model.stacks, model.currentRevisions, model.currentPublishedCount), nil
}

func (s *Service) stackMergedIntoBase(ctx context.Context, repoRoot string, stack StackInfo, bookmarkTargets map[string]string) bool {
	baseRef := strings.TrimSpace(stack.BaseRef)
	if baseRef == "" || IsTerminalStackStatus(stack.Status) {
		return false
	}
	remoteName := ""
	if stack.RemoteName != nil {
		remoteName = strings.TrimSpace(*stack.RemoteName)
	}
	baseSelectors := s.stackMergeBaseSelectors(ctx, repoRoot, baseRef, remoteName)
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
			revset := fmt.Sprintf("(%s) & ancestors(%s)", quoteJJRev(selector), quoteJJRev(baseSelector))
			out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", revset, "-n", "1", "--no-graph", "-T", "change_id")
			if err != nil {
				continue
			}
			if strings.TrimSpace(out) != "" {
				return true
			}
		}
	}
	return false
}

func stackVisibleInGX(stack StackInfo) bool {
	return !IsTerminalStackStatus(stack.Status)
}

func isTerminalStackStatus(status string) bool {
	return IsTerminalStackStatus(status)
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

func (s *Service) JJVersion(ctx context.Context) (string, error) {
	return s.runStdoutTrimmed(ctx, "", "jj", "--version")
}

func (s *Service) ResolveJJRepo(ctx context.Context) (RepoInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveJJRepoAtPath(ctx, cwd)
}

func (s *Service) configuredJJRepo(ctx context.Context) (RepoInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, err
	}
	repo, err := s.ResolveJJRepoAtPath(ctx, cwd)
	if err != nil {
		return RepoInfo{}, err
	}
	_, _, _, err = s.ensureIdentity(ctx, repo.RootPath, InitOptions{})
	if err != nil {
		return RepoInfo{}, err
	}
	if err := s.ensureGXInternalIgnored(ctx, repo.RootPath); err != nil {
		return RepoInfo{}, err
	}
	return repo, nil
}

func hasJJRepo(ctx context.Context, s *Service, startPath string) bool {
	_, err := s.ResolveJJRepoAtPath(ctx, startPath)
	return err == nil
}

func (s *Service) ResolveJJRepoAtPath(ctx context.Context, startPath string) (RepoInfo, error) {
	root, err := s.runStdoutTrimmed(ctx, startPath, "jj", "root")
	if err != nil {
		return RepoInfo{}, err
	}
	info, err := s.resolveGitInfo(ctx, root)
	if err != nil {
		return RepoInfo{}, err
	}
	info.RootPath = root
	info.Backend = "jj"
	info = s.withStoredRepoConfig(ctx, info)
	return info, nil
}

func (s *Service) ResolveGitRepo(ctx context.Context) (RepoInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveGitRepoAtPath(ctx, cwd)
}

func (s *Service) ResolveGitRepoAtPath(ctx context.Context, startPath string) (RepoInfo, error) {
	root, err := s.runTrimmed(ctx, startPath, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return RepoInfo{}, err
	}
	info, err := s.resolveGitInfo(ctx, root)
	if err != nil {
		return RepoInfo{}, err
	}
	info.RootPath = root
	if info.Backend == "" {
		info.Backend = "git"
	}
	info = s.withStoredRepoConfig(ctx, info)
	return info, nil
}

func (s *Service) withStoredRepoConfig(ctx context.Context, info RepoInfo) RepoInfo {
	store, err := openStore(ctx)
	if err != nil {
		return info
	}
	defer store.Close()
	repo, err := store.FindRepoByRoot(ctx, info.RootPath)
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
	const tmpl = `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", tmpl)
	if err != nil {
		return ChangeInfo{}, err
	}
	line := lastNonEmptyLine(out)
	parts := strings.SplitN(line, "|", 4)
	if len(parts) < 4 {
		return ChangeInfo{}, fmt.Errorf("unexpected jj log output: %q", out)
	}
	var parent *string
	if parts[3] != "" {
		first := strings.SplitN(parts[3], ",", 2)[0]
		parent = &first
	}
	filesOut, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "diff", "-r", rev, "--name-only")
	if err != nil {
		return ChangeInfo{}, err
	}
	return ChangeInfo{
		ChangeID:       parts[0],
		CommitID:       parts[1],
		Description:    parts[2],
		ParentChangeID: parent,
		Files:          splitLines(filesOut),
	}, nil
}

func (s *Service) DiffGit(ctx context.Context, repoRoot, rev string) (string, error) {
	return s.runStdoutTrimmed(ctx, repoRoot, "jj", "diff", "-r", rev, "--git")
}

func (s *Service) CurrentOperation(ctx context.Context, repoRoot string) (string, error) {
	return s.runStdoutTrimmed(ctx, repoRoot, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id")
}

func (s *Service) WorkRev(ctx context.Context, repoRoot string) (string, error) {
	return s.currentWorkRev(ctx, repoRoot)
}

func (s *Service) currentWorkRev(ctx context.Context, repoRoot string) (string, error) {
	empty, err := s.isRevisionEmpty(ctx, repoRoot, "@")
	if err != nil {
		return "", err
	}
	if empty {
		return "@-", nil
	}
	return "@", nil
}

func (s *Service) currentWorkChangeID(ctx context.Context, repoRoot string) (string, error) {
	rev, err := s.currentWorkRev(ctx, repoRoot)
	if err != nil {
		return "", err
	}
	change, err := s.CurrentChange(ctx, repoRoot, rev)
	if err != nil {
		return "", err
	}
	return change.ChangeID, nil
}

func recordCommit(ctx context.Context, result CommitResult) error {
	return withBusyRetry(ctx, "record commit metadata", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		return recordChangeForStack(ctx, store, result.Repo, result.Stack, result.Change, result.OperationID, result.PreferredSessionIDs)
	})
}

func recordModify(ctx context.Context, result ModifyResult) error {
	return withBusyRetry(ctx, "record modify metadata", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()

		repoID, err := upsertRepo(ctx, store, result.Repo)
		if err != nil {
			return err
		}
		prevID, err := upsertChange(ctx, store, repoID, result.PreviousChange)
		if err != nil {
			return err
		}
		currID, err := upsertChange(ctx, store, repoID, result.CurrentChange)
		if err != nil {
			return err
		}
		if err := writeRevision(ctx, store, currID, result.CurrentChange, result.OperationID); err != nil {
			return err
		}
		if result.Stack != nil {
			head := result.CurrentChange.ChangeID
			headCommit := result.CurrentChange.CommitID
			result.Stack.HeadChangeID = &head
			result.Stack.HeadCommitID = &headCommit
			stackID, err := upsertStackInfo(ctx, store, repoID, *result.Stack)
			if err != nil {
				return err
			}
			if err := store.AddChangeToStack(ctx, stackID, currID, time.Now().UnixMilli()); err != nil {
				return err
			}
		}
		if err := attachSessions(ctx, store, result.Repo.RootPath, currID, nil); err != nil {
			return err
		}
		return store.WriteModifyEvent(ctx, storage.ModifyEvent{
			RepoID:                  repoID,
			TargetChangeID:          currID,
			PreviousCurrentChangeID: &prevID,
			JJOperationID:           result.OperationID,
			CreatedAt:               time.Now().UnixMilli(),
		})
	})
}

// ResetPublished clears push history and marks every stack draft so `gx status` shows
// the full stack as unpublished again.
func ResetPublished(ctx context.Context) (int, error) {
	var stacksUpdated int
	err := withBusyRetry(ctx, "reset published state", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		stacksUpdated, err = store.ResetAllPublished(ctx)
		return err
	})
	return stacksUpdated, err
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

func recordPush(ctx context.Context, result PushResult) error {
	return withBusyRetry(ctx, "record push metadata", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()

		repoID, err := upsertRepo(ctx, store, result.Repo)
		if err != nil {
			return err
		}
		var changeID *int64
		if result.CurrentChange != nil {
			id, err := upsertChange(ctx, store, repoID, *result.CurrentChange)
			if err != nil {
				return err
			}
			if err := writeRevision(ctx, store, id, *result.CurrentChange, ""); err != nil {
				return err
			}
			changeID = &id
		}
		if result.Stack != nil {
			if result.GitExported {
				remoteRef := strings.TrimSpace(result.GitPublishedRef)
				if remoteRef == "" && result.Repo.BranchName != nil {
					remoteRef = "refs/heads/" + strings.TrimSpace(*result.Repo.BranchName)
				}
				if remoteRef != "" {
					result.Stack.RemoteRef = &remoteRef
				}
				result.Stack.RemoteName = result.RemoteName
			}
			result.Stack.Status = "published"
			if _, err := upsertStackInfo(ctx, store, repoID, *result.Stack); err != nil {
				return err
			}
		}
		return store.WritePush(ctx, storage.Push{
			RepoID:          repoID,
			RemoteName:      result.RemoteName,
			BranchName:      ptr(firstNonEmpty(result.GXStackRef, derefString(result.Repo.BranchName))),
			HeadCommitID:    result.HeadCommitID,
			CurrentChangeID: changeID,
			CreatedAt:       time.Now().UnixMilli(),
		})
	})
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

func recordChangeForStack(ctx context.Context, store *storage.Store, repo RepoInfo, stack *StackInfo, change ChangeInfo, opID string, preferredSessionIDs []string) error {
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
	if err := attachSessions(ctx, store, repo.RootPath, changeID, preferredSessionIDs); err != nil {
		return err
	}
	return nil
}

func upsertRepo(ctx context.Context, store *storage.Store, repo RepoInfo) (int64, error) {
	now := time.Now().UnixMilli()
	return store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repo.RootPath,
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

func (s *Service) markChangeAbandoned(ctx context.Context, repoRoot, jjChangeID string) error {
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return err
	}
	if repoRow == nil {
		return nil
	}
	change, err := store.FindChangeByJJChangeID(ctx, repoRow.ID, jjChangeID)
	if err != nil {
		return err
	}
	if change == nil {
		return nil
	}
	return store.MarkChangeStatus(ctx, change.ID, "abandoned", time.Now().UnixMilli())
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
	jjName, _ := s.jjConfigGet(ctx, repoRoot, "user.name")
	jjEmail, _ := s.jjConfigGet(ctx, repoRoot, "user.email")
	cfg, err := gxconfig.Load()
	if err != nil {
		return "", "", false, err
	}
	gitName := s.gitConfigValue(ctx, repoRoot, "user.name")
	gitEmail := s.gitConfigValue(ctx, repoRoot, "user.email")

	name := firstNonEmpty(
		strings.TrimSpace(opts.Name),
		strings.TrimSpace(cfg.User.Name),
		strings.TrimSpace(jjName),
		strings.TrimSpace(gitName),
	)
	email := firstNonEmpty(
		strings.TrimSpace(opts.Email),
		strings.TrimSpace(cfg.User.Email),
		strings.TrimSpace(jjEmail),
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
	if jjName != name {
		if err := s.jjConfigSetUser(ctx, repoRoot, "user.name", name); err != nil {
			return "", "", false, err
		}
		changed = true
	}
	if jjEmail != email {
		if err := s.jjConfigSetUser(ctx, repoRoot, "user.email", email); err != nil {
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

func (s *Service) activeContainerName(ctx context.Context, repo RepoInfo) (string, error) {
	body, err := s.resolveCurrentStack(ctx, repo, true)
	if err != nil {
		return "", err
	}
	return body.BookmarkName, nil
}

func (s *Service) ensureContainerAttached(ctx context.Context, repo RepoInfo) (string, error) {
	body, err := s.resolveCurrentStack(ctx, repo, true)
	if err != nil {
		return "", err
	}
	name := body.BookmarkName
	if repo.BranchName != nil && strings.TrimSpace(*repo.BranchName) == name {
		if err := s.setBookmarkTarget(ctx, repo.RootPath, name); err != nil {
			return "", err
		}
		return name, nil
	}
	_, err = s.reattachContainer(ctx, repo.RootPath, name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *Service) resolveCurrentStack(ctx context.Context, repo RepoInfo, create bool) (StackInfo, error) {
	return s.resolveCurrentStackOptions(ctx, repo, stackResolveOptions{
		create:      create,
		description: "",
	})
}

type stackResolveOptions struct {
	create      bool
	description string
}

func (s *Service) resolveCurrentStackForPublish(ctx context.Context, repo RepoInfo) (StackInfo, error) {
	return s.resolveCurrentStackOptions(ctx, repo, stackResolveOptions{
		create: true,
	})
}

func (s *Service) resolveCurrentStackWithDescription(ctx context.Context, repo RepoInfo, create bool, description string) (StackInfo, error) {
	return s.resolveCurrentStackOptions(ctx, repo, stackResolveOptions{
		create:      create,
		description: description,
	})
}

func (s *Service) resolveCurrentStackOptions(ctx context.Context, repo RepoInfo, opts stackResolveOptions) (StackInfo, error) {
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, err
	}
	defer store.Close()

	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, err
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoID); err != nil {
		return StackInfo{}, err
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoID); err != nil {
		return StackInfo{}, err
	}
	latest, err := store.LatestStackByRepoID(ctx, repoID)
	if err != nil {
		return StackInfo{}, err
	}
	if latest == nil {
		if !opts.create {
			return StackInfo{}, ErrNoGitBranch
		}
		return s.createDraftStack(ctx, store, repo, repoID, opts.description)
	}
	if current, err := s.resolveExplicitStack(ctx, store, repoID, repo); err != nil {
		return StackInfo{}, err
	} else if current != nil {
		return stackInfoFromStorage(*current), nil
	}
	if current, err := s.stackFromCurrentBookmark(ctx, store, repoID, repo.RootPath); err == nil && current != nil {
		return stackInfoFromStorage(*current), nil
	}
	if attached, err := s.stackFromAttachedBranch(ctx, store, repoID, repo); err != nil {
		return StackInfo{}, err
	} else if attached != nil {
		return stackInfoFromStorage(*attached), nil
	}
	if headChangeID, err := s.currentWorkChangeID(ctx, repo.RootPath); err == nil && strings.TrimSpace(headChangeID) != "" {
		if stack, err := store.FindStackByHeadChange(ctx, repoID, headChangeID); err == nil && stack != nil {
			return stackInfoFromStorage(*stack), nil
		}
	}
	if parentChangeID, err := s.changeIDForRev(ctx, repo.RootPath, "@-"); err == nil && strings.TrimSpace(parentChangeID) != "" {
		if stack, err := store.FindStackByHeadChange(ctx, repoID, parentChangeID); err == nil && stack != nil {
			return stackInfoFromStorage(*stack), nil
		}
	}
	if opts.create {
		return s.createDraftStack(ctx, store, repo, repoID, opts.description)
	}
	if latest == nil {
		return StackInfo{}, ErrNoGitBranch
	}
	return stackInfoFromStorage(*latest), nil
}

func (s *Service) stackForRevision(ctx context.Context, repo RepoInfo, rev string) (StackInfo, bool, error) {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return StackInfo{}, false, nil
	}
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, false, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, false, err
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoID); err != nil {
		return StackInfo{}, false, err
	}
	changeID, err := s.changeIDForRev(ctx, repo.RootPath, rev)
	if err != nil {
		return StackInfo{}, false, err
	}
	if stack, err := store.FindStackByHeadChange(ctx, repoID, changeID); err != nil {
		return StackInfo{}, false, err
	} else if stack != nil {
		return stackInfoFromStorage(*stack), true, nil
	}
	out, err := s.runStdoutTrimmed(ctx, repo.RootPath, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return StackInfo{}, false, err
	}
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) != changeID {
			continue
		}
		stack, err := store.FindStackByBookmark(ctx, repoID, strings.TrimSpace(parts[0]))
		if err != nil {
			return StackInfo{}, false, err
		}
		if stack != nil {
			return stackInfoFromStorage(*stack), true, nil
		}
	}
	return StackInfo{}, false, nil
}

func (s *Service) resolveExplicitStack(ctx context.Context, store *storage.Store, repoID int64, repo RepoInfo) (*storage.Stack, error) {
	if current, err := s.stackFromCurrentBookmark(ctx, store, repoID, repo.RootPath); err != nil {
		return nil, err
	} else if current != nil {
		return current, nil
	}
	return s.stackFromAttachedBranch(ctx, store, repoID, repo)
}

func (s *Service) stackFromAttachedBranch(ctx context.Context, store *storage.Store, repoID int64, repo RepoInfo) (*storage.Stack, error) {
	if repo.BranchName == nil {
		return nil, nil
	}
	name := strings.TrimSpace(*repo.BranchName)
	if name == "" || legacyGXInternalCheckoutRef(name) || name == repo.defaultBaseBranch() {
		return nil, nil
	}
	return store.FindStackByBookmark(ctx, repoID, name)
}

func (s *Service) repairMissingStackRowsFromBookmarks(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64) error {
	if repoID == 0 || strings.TrimSpace(repo.RootPath) == "" {
		return nil
	}
	targets, err := s.jjBookmarkTargets(ctx, repo.RootPath)
	if err != nil {
		return nil
	}
	baseNames := map[string]struct{}{
		repo.defaultBaseBranch():                     {},
		baseCheckoutRef(s.defaultStackBaseRef(repo)): {},
		baseCheckoutRef(repo.authoringBaseRef()):     {},
	}
	for _, target := range targets {
		bookmark := strings.TrimSpace(target.Name)
		if bookmark == "" || strings.Contains(bookmark, "@") || legacyGXInternalCheckoutRef(bookmark) {
			continue
		}
		if _, isBase := baseNames[bookmark]; isBase {
			continue
		}
		if legacyStackBookmarkName(bookmark) {
			canonical := stackBookmarkName(bookmark, target.ChangeID)
			if canonical == "" {
				continue
			}
			if canonical != bookmark {
				if exists, err := s.jjBookmarkExists(ctx, repo.RootPath, canonical); err != nil {
					return err
				} else if exists {
					continue
				}
				if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "bookmark", "rename", bookmark, canonical); err != nil {
					return fmt.Errorf("rename legacy stack bookmark %q to %q: %w", bookmark, canonical, err)
				}
				bookmark = canonical
			}
		}
		existing, err := store.FindStackByBookmark(ctx, repoID, bookmark)
		if err != nil {
			return err
		}
		if existing != nil {
			continue
		}
		change, err := store.FindChangeByJJChangeID(ctx, repoID, target.ChangeID)
		if err != nil {
			return err
		}
		if change == nil {
			if !isGXStackBookmark(bookmark) {
				continue
			}
			info, err := s.CurrentChange(ctx, repo.RootPath, bookmark)
			if err != nil {
				continue
			}
			changeID, err := upsertChange(ctx, store, repoID, info)
			if err != nil {
				return err
			}
			change, err = store.FindChangeByJJChangeID(ctx, repoID, info.ChangeID)
			if err != nil {
				return err
			}
			if change == nil {
				change = &storage.Change{
					ID:              changeID,
					RepoID:          repoID,
					JJChangeID:      info.ChangeID,
					CurrentCommitID: info.CommitID,
					Description:     info.Description,
					ParentChangeID:  info.ParentChangeID,
					Status:          "draft",
				}
			}
		}
		headChangeID := strings.TrimSpace(change.JJChangeID)
		if headChangeID == "" {
			headChangeID = strings.TrimSpace(target.ChangeID)
		}
		headCommitID := strings.TrimSpace(change.CurrentCommitID)
		if headCommitID == "" {
			info, err := s.CurrentChange(ctx, repo.RootPath, bookmark)
			if err == nil {
				headCommitID = strings.TrimSpace(info.CommitID)
			}
		}
		baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
		now := time.Now().UnixMilli()
		stackID, err := store.UpsertStack(ctx, storage.Stack{
			RepoID:       repoID,
			Name:         firstNonEmpty(stackNameFromBookmark(bookmark), bookmark),
			BookmarkName: bookmark,
			BaseRef:      baseRef,
			BaseCommitID: s.stackBaseCommitID(ctx, repo.RootPath, baseRef),
			HeadChangeID: &headChangeID,
			HeadCommitID: &headCommitID,
			Status:       "draft",
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if err != nil {
			return err
		}
		if change.ID != 0 {
			if err := store.AddChangeToStack(ctx, stackID, change.ID, now); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) resolveCurrentStackForRead(ctx context.Context, repo RepoInfo) (StackInfo, bool, error) {
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, false, err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return StackInfo{}, false, err
	}
	if repoRow == nil {
		return StackInfo{}, false, nil
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoRow.ID); err != nil {
		return StackInfo{}, false, err
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoRow.ID); err != nil {
		return StackInfo{}, false, err
	}
	return s.resolveCurrentStackForReadWithStore(ctx, store, repo, repoRow.ID)
}

func (s *Service) resolveCurrentStackForReadWithStore(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64) (StackInfo, bool, error) {
	if s.isOnBaseBranch(repo) || s.onAuthoringCheckout(ctx, repo) {
		return StackInfo{}, false, nil
	}
	if current, err := s.stackFromCurrentBookmark(ctx, store, repoID, repo.RootPath); err == nil && current != nil {
		return stackInfoFromStorage(*current), true, nil
	}
	if headChangeID, err := s.currentWorkChangeID(ctx, repo.RootPath); err == nil && strings.TrimSpace(headChangeID) != "" {
		if stack, err := store.FindStackByHeadChange(ctx, repoID, headChangeID); err == nil && stack != nil {
			return stackInfoFromStorage(*stack), true, nil
		}
	}
	if parentChangeID, err := s.changeIDForRev(ctx, repo.RootPath, "@-"); err == nil && strings.TrimSpace(parentChangeID) != "" {
		if stack, err := store.FindStackByHeadChange(ctx, repoID, parentChangeID); err == nil && stack != nil {
			return stackInfoFromStorage(*stack), true, nil
		}
	}
	latest, err := store.LatestStackByRepoID(ctx, repoID)
	if err != nil {
		return StackInfo{}, false, err
	}
	if latest == nil {
		return StackInfo{}, false, nil
	}
	return stackInfoFromStorage(*latest), true, nil
}

func (s *Service) isOnBaseBranch(repo RepoInfo) bool {
	if repo.BranchName == nil {
		return false
	}
	branch := strings.TrimSpace(*repo.BranchName)
	if branch == "" || isGXStackBookmark(branch) {
		return false
	}
	return branch == repo.defaultBaseBranch()
}

func (s *Service) onAuthoringCheckout(ctx context.Context, repo RepoInfo) bool {
	result := s.baseResult(ctx, repo)
	current := strings.TrimSpace(result.CurrentRef)
	return current == gxAuthoringCheckoutRef(result.BaseRef) ||
		current == result.BaseRef
}

func (s *Service) stackFromCurrentBookmark(ctx context.Context, store *storage.Store, repoID int64, repoRoot string) (*storage.Stack, error) {
	changeID, err := s.currentWorkChangeID(ctx, repoRoot)
	if err != nil || strings.TrimSpace(changeID) == "" {
		return nil, err
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return nil, err
	}
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) != changeID {
			continue
		}
		stack, err := store.FindStackByBookmark(ctx, repoID, strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, err
		}
		if stack != nil {
			return stack, nil
		}
	}
	return nil, nil
}

func (s *Service) createDraftStack(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, description string) (StackInfo, error) {
	targetRev, err := s.currentWorkRev(ctx, repo.RootPath)
	if err != nil {
		return StackInfo{}, err
	}
	change, err := s.CurrentChange(ctx, repo.RootPath, targetRev)
	if err != nil {
		return StackInfo{}, err
	}
	headChangeID := change.ChangeID
	headCommitID := change.CommitID
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	baseCommitID := s.stackBaseCommitID(ctx, repo.RootPath, baseRef)
	name := firstNonEmpty(strings.TrimSpace(description), s.inferStackName(repo, change))
	bookmark := stackBookmarkName(name, headChangeID)
	if existing, err := store.FindStackByBookmark(ctx, repoID, bookmark); err != nil {
		return StackInfo{}, err
	} else if existing != nil {
		if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "edit", bookmark); err != nil {
			return StackInfo{}, fmt.Errorf("edit existing stack %s: %w", bookmark, err)
		}
		return stackInfoFromStorage(*existing), nil
	}
	if jjBookmarkExists, err := s.jjBookmarkExists(ctx, repo.RootPath, bookmark); err != nil {
		return StackInfo{}, err
	} else if jjBookmarkExists {
		if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "edit", bookmark); err != nil {
			return StackInfo{}, fmt.Errorf("edit existing jj bookmark %s: %w", bookmark, err)
		}
		targetRev, err = s.currentWorkRev(ctx, repo.RootPath)
		if err != nil {
			return StackInfo{}, err
		}
		change, err = s.CurrentChange(ctx, repo.RootPath, targetRev)
		if err != nil {
			return StackInfo{}, err
		}
		headChangeID = change.ChangeID
		headCommitID = change.CommitID
	} else {
		if err := s.ensureForkedForNewStack(ctx, store, repoID, repo.RootPath); err != nil {
			return StackInfo{}, err
		}
		targetRev, err = s.currentWorkRev(ctx, repo.RootPath)
		if err != nil {
			return StackInfo{}, err
		}
		change, err = s.CurrentChange(ctx, repo.RootPath, targetRev)
		if err != nil {
			return StackInfo{}, err
		}
		headChangeID = change.ChangeID
		headCommitID = change.CommitID
		if err := s.relocateKnownStackBookmarksFromChange(ctx, store, repoID, repo.RootPath, bookmark, targetRev); err != nil {
			return StackInfo{}, err
		}
		if err := s.assertBookmarkTargetAvailable(ctx, store, repoID, repo.RootPath, bookmark, targetRev); err != nil {
			return StackInfo{}, err
		}
		if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "bookmark", "set", bookmark, "-r", targetRev); err != nil {
			return StackInfo{}, err
		}
	}
	now := time.Now().UnixMilli()
	id, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         name,
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommitID,
		HeadChangeID: &headChangeID,
		HeadCommitID: &headCommitID,
		Status:       "draft",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return StackInfo{}, err
	}
	body := storage.Stack{
		ID:           id,
		RepoID:       repoID,
		Name:         name,
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommitID,
		HeadChangeID: &headChangeID,
		HeadCommitID: &headCommitID,
		Status:       "draft",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return stackInfoFromStorage(body), nil
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
		if !isGXStackBookmark(branch) {
			return branch
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

func (s *Service) currentGitCheckoutRef(ctx context.Context, repoRoot string) string {
	if branch, err := s.gitValue(ctx, repoRoot, "branch", "--show-current"); err == nil {
		return strings.TrimSpace(branch)
	}
	return ""
}

func (s *Service) publicBookmarkForRev(ctx context.Context, repoRoot, rev, defaultBranch string) string {
	changeID, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", "change_id")
	if err != nil || strings.TrimSpace(changeID) == "" {
		return ""
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return ""
	}
	bestRank := 99
	best := ""
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) != strings.TrimSpace(changeID) {
			continue
		}
		name := strings.TrimSpace(parts[0])
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

func (s *Service) defaultSwitchBaseRef(ctx context.Context, repo RepoInfo) string {
	return repo.authoringBaseRef()
}

func (s *Service) baseResult(ctx context.Context, repo RepoInfo) BaseResult {
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	current := ""
	if repo.BranchName != nil {
		current = strings.TrimSpace(*repo.BranchName)
	}
	if current == "" {
		if branch, err := s.gitValue(ctx, repo.RootPath, "branch", "--show-current"); err == nil {
			current = strings.TrimSpace(branch)
		}
	}
	checkoutRef := baseCheckoutRef(baseRef)
	onBase := current == checkoutRef
	if !onBase && current != "" && baseRef != "" {
		onBase = current == baseRef || s.refsPointToSameChange(ctx, repo.RootPath, current, baseRef)
	}
	return BaseResult{
		Repo:          repo,
		BaseRef:       baseRef,
		DefaultBranch: repo.defaultBaseBranch(),
		CurrentRef:    current,
		OnBase:        onBase,
	}
}

func (s *Service) refsPointToSameChange(ctx context.Context, repoRoot, left, right string) bool {
	leftChange, err := s.changeIDForRev(ctx, repoRoot, left)
	if err != nil || strings.TrimSpace(leftChange) == "" {
		return false
	}
	rightChange, err := s.changeIDForRev(ctx, repoRoot, right)
	if err != nil || strings.TrimSpace(rightChange) == "" {
		return false
	}
	return leftChange == rightChange
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

func (s *Service) inferStackName(repo RepoInfo, change ChangeInfo) string {
	if repo.BranchName != nil {
		name := strings.TrimSpace(*repo.BranchName)
		if name != "" && name != repo.defaultBaseBranch() && !isGXStackBookmark(name) {
			return name
		}
	}
	if strings.TrimSpace(change.Description) != "" {
		return change.Description
	}
	return "change-" + shortID(change.ChangeID, 8)
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

func stackMatchesDirectName(body storage.Stack, name string) bool {
	return body.Name == name ||
		body.BookmarkName == name
}

func stackAlias(index int) string {
	return "s" + strconv.Itoa(index+1)
}

func recordedEditCheckoutBranch(_ RepoInfo) string {
	return ""
}

func (s *Service) reattachRecordedContainer(ctx context.Context, repoRoot, stackBookmark, checkoutBranch string) (RepoInfo, error) {
	repo, err := s.updateContainerBookmark(ctx, repoRoot, stackBookmark)
	if err != nil {
		return RepoInfo{}, err
	}
	checkoutBranch = strings.TrimSpace(checkoutBranch)
	stackBookmark = strings.TrimSpace(stackBookmark)
	branch := firstNonEmpty(checkoutBranch, stackBookmark)
	if branch == "" {
		return s.ResolveJJRepoAtPath(ctx, repo.RootPath)
	}
	var commitID string
	if checkoutBranch != "" {
		commitID, err = s.commitIDForContainer(ctx, repoRoot)
	} else {
		commitID, err = s.commitIDForRev(ctx, repoRoot, stackBookmark)
	}
	if err != nil {
		return RepoInfo{}, err
	}
	if err := s.attachGitBranch(ctx, repoRoot, branch, commitID); err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveJJRepoAtPath(ctx, repo.RootPath)
}

func (s *Service) reattachContainer(ctx context.Context, repoRoot, name string) (RepoInfo, error) {
	repo, err := s.updateContainerBookmark(ctx, repoRoot, name)
	if err != nil {
		return RepoInfo{}, err
	}
	commitID, err := s.commitIDForContainer(ctx, repoRoot)
	if err != nil {
		return RepoInfo{}, err
	}
	if err := s.attachGitBranch(ctx, repoRoot, name, commitID); err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveJJRepoAtPath(ctx, repo.RootPath)
}

func (s *Service) updateContainerBookmark(ctx context.Context, repoRoot, name string) (RepoInfo, error) {
	if strings.TrimSpace(name) == "" {
		return RepoInfo{}, fmt.Errorf("container name is empty")
	}
	store, err := openStore(ctx)
	if err != nil {
		return RepoInfo{}, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return RepoInfo{}, err
	}
	if repoRow != nil {
		if err := s.repairSharedStackBookmarks(ctx, store, repoRoot, repoRow.ID, name); err != nil {
			return RepoInfo{}, err
		}
		targetRev, targetErr := s.containerTargetRev(ctx, repoRoot)
		if targetErr == nil {
			if err := s.relocateKnownStackBookmarksFromChange(ctx, store, repoRow.ID, repoRoot, name, targetRev); err != nil {
				return RepoInfo{}, err
			}
			if err := s.assertBookmarkTargetAvailable(ctx, store, repoRow.ID, repoRoot, name, targetRev); err != nil {
				return RepoInfo{}, err
			}
		}
	}
	if err := s.setBookmarkTarget(ctx, repoRoot, name); err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveJJRepoAtPath(ctx, repoRoot)
}

func (s *Service) attachGitBranch(ctx context.Context, repoRoot, name, commitID string) error {
	name = strings.TrimPrefix(strings.TrimSpace(name), "refs/heads/")
	commitID = strings.TrimSpace(commitID)
	if name == "" {
		return fmt.Errorf("git branch name is empty")
	}
	if commitID == "" {
		return fmt.Errorf("git branch %s target commit is empty", name)
	}
	if _, err := s.runner.Run(ctx, repoRoot, "git", "update-ref", "refs/heads/"+name, commitID); err != nil {
		return err
	}
	if _, err := s.runner.Run(ctx, repoRoot, "git", "symbolic-ref", "HEAD", "refs/heads/"+name); err != nil {
		return err
	}
	if _, err := s.runner.Run(ctx, repoRoot, "git", "reset", "--mixed", "HEAD"); err != nil {
		return err
	}
	return nil
}

func (s *Service) reattachGitHeadToBaseRef(ctx context.Context, repo RepoInfo, baseRef string) (RepoInfo, error) {
	checkoutRef := baseCheckoutRef(baseRef)
	commitID, err := s.commitIDForContainer(ctx, repo.RootPath)
	if err != nil {
		return RepoInfo{}, err
	}
	if err := s.attachGitBranch(ctx, repo.RootPath, checkoutRef, commitID); err != nil {
		return RepoInfo{}, err
	}
	return s.ResolveJJRepoAtPath(ctx, repo.RootPath)
}

func (s *Service) setBookmarkTarget(ctx context.Context, repoRoot, name string) error {
	targetRev, err := s.containerTargetRev(ctx, repoRoot)
	if err != nil {
		return err
	}
	_, err = s.runner.Run(ctx, repoRoot, "jj", "bookmark", "set", name, "-r", targetRev, "--allow-backwards")
	return err
}

func (s *Service) containerTargetRev(ctx context.Context, repoRoot string) (string, error) {
	empty, err := s.isRevisionEmpty(ctx, repoRoot, "@")
	if err != nil {
		return "", err
	}
	if empty {
		return "@-", nil
	}
	return "@", nil
}

func (s *Service) commitIDForContainer(ctx context.Context, repoRoot string) (string, error) {
	targetRev, err := s.containerTargetRev(ctx, repoRoot)
	if err != nil {
		return "", err
	}
	return s.commitIDForRev(ctx, repoRoot, targetRev)
}

func (s *Service) commitIDForRev(ctx context.Context, repoRoot, rev string) (string, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", "commit_id")
	if err != nil {
		return "", err
	}
	return lastNonEmptyLine(out), nil
}

func (s *Service) changeIDForRev(ctx context.Context, repoRoot, rev string) (string, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", "change_id")
	if err != nil {
		return "", err
	}
	return lastNonEmptyLine(out), nil
}

func (s *Service) isRevisionEmpty(ctx context.Context, repoRoot, rev string) (bool, error) {
	value, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", "empty")
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(lastNonEmptyLine(value)), "true"), nil
}

func (s *Service) jjConfigGet(ctx context.Context, repoRoot, name string) (string, error) {
	return s.runStdoutTrimmed(ctx, repoRoot, "jj", "config", "get", name)
}

func (s *Service) jjConfigSetUser(ctx context.Context, repoRoot, name, value string) error {
	_, err := s.runner.Run(ctx, repoRoot, "jj", "config", "set", "--user", name, value)
	return err
}

func (s *Service) remoteURLForTarget(ctx context.Context, repoRoot, target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", fmt.Errorf("remote target is empty")
	}
	if strings.Contains(target, "://") || githubSSHPattern.MatchString(target) {
		return target, nil
	}
	return s.gitValue(ctx, repoRoot, "remote", "get-url", target)
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

func nonFlagArgs(args []string) []string {
	values := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		values = append(values, arg)
	}
	return values
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

func githubPullRequestURL(remoteURL, baseBranch, headBranch string) *string {
	host, owner, repo, ok := parseGitHubRemote(remoteURL)
	if !ok {
		return nil
	}
	base := url.PathEscape(strings.TrimSpace(baseBranch))
	head := url.PathEscape(strings.TrimSpace(headBranch))
	if base == "" || head == "" {
		return nil
	}
	value := fmt.Sprintf("https://%s/%s/%s/compare/%s...%s?quick_pull=1", host, owner, repo, base, head)
	return &value
}

func parseGitHubRemote(remoteURL string) (host, owner, repo string, ok bool) {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return "", "", "", false
	}
	if match := githubSSHPattern.FindStringSubmatch(remoteURL); match != nil {
		host = match[githubSSHPattern.SubexpIndex("host")]
		owner = match[githubSSHPattern.SubexpIndex("owner")]
		repo = strings.TrimSuffix(match[githubSSHPattern.SubexpIndex("repo")], ".git")
		return normalizeGitHubTarget(host, owner, repo)
	}
	if !strings.Contains(remoteURL, "://") {
		return "", "", "", false
	}
	parsed, err := url.Parse(remoteURL)
	if err != nil {
		return "", "", "", false
	}
	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) < 2 {
		return "", "", "", false
	}
	owner = segments[0]
	repo = strings.TrimSuffix(path.Base(parsed.Path), ".git")
	host = parsed.Host
	return normalizeGitHubTarget(host, owner, repo)
}

func normalizeGitHubTarget(host, owner, repo string) (string, string, string, bool) {
	host = strings.TrimSpace(host)
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	if host == "" || owner == "" || repo == "" {
		return "", "", "", false
	}
	if !strings.Contains(strings.ToLower(host), "github") {
		return "", "", "", false
	}
	return host, owner, repo, true
}

func branchFromPushArgs(args []string) *string {
	if len(args) == 0 {
		return nil
	}
	refspec := args[0]
	if idx := strings.LastIndex(refspec, ":"); idx >= 0 && idx+1 < len(refspec) {
		target := refspec[idx+1:]
		target = strings.TrimPrefix(target, "refs/heads/")
		if target != "" && target != "HEAD" {
			return &target
		}
	}
	if refspec != "" && !strings.HasPrefix(refspec, "-") && refspec != "HEAD" {
		return &refspec
	}
	return nil
}

func publishRefFromArgs(args []string) string {
	values := nonFlagArgs(args)
	if len(values) == 0 {
		return ""
	}
	if branch := branchFromPushArgs(values); branch != nil {
		return strings.TrimSpace(*branch)
	}
	return ""
}

func publishRefForStack(body StackInfo) string {
	if bookmark := strings.TrimSpace(body.BookmarkName); bookmark != "" {
		if legacyStackBookmarkName(bookmark) {
			return stackBookmarkName(bookmark, derefString(body.HeadChangeID))
		}
		return cleanRefName(bookmark)
	}
	return stackBookmarkName(body.Name, derefString(body.HeadChangeID))
}

func recordStackBookmarks(ctx context.Context, repo RepoInfo, remoteName string, stack []PushedChange) error {
	if len(stack) == 0 {
		return nil
	}
	return withBusyRetry(ctx, "record stack bookmarks", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()

		repoID, err := upsertRepo(ctx, store, repo)
		if err != nil {
			return err
		}
		now := time.Now().UnixMilli()
		remote := strings.TrimSpace(remoteName)
		for _, entry := range stack {
			changeID, err := upsertChange(ctx, store, repoID, entry.Change)
			if err != nil {
				return err
			}
			remoteRef := "refs/heads/" + entry.BranchName
			if err := store.UpsertChangeBookmark(ctx, storage.ChangeBookmark{
				ChangeID:           changeID,
				BookmarkName:       entry.BranchName,
				RemoteName:         &remote,
				RemoteRef:          &remoteRef,
				LastPushedCommitID: entry.Change.CommitID,
				CreatedAt:          now,
				UpdatedAt:          now,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func stackPushOutput(stack []PushedChange, gitExported bool, gitPushStatus string) string {
	if len(stack) == 0 {
		return ""
	}
	ref := stack[0].BranchName
	action := "Prepared"
	if gitExported {
		action = "Pushed"
		if gitPushStatus == "already up to date" {
			action = "Found"
		}
	}
	lines := []string{fmt.Sprintf("%s %d GX revisions on GitHub branch %s:", action, len(stack), ref)}
	for _, entry := range stack {
		lines = append(lines, fmt.Sprintf("  %s %s", shortID(entry.Change.ChangeID, 12), entry.Change.Description))
	}
	return strings.Join(lines, "\n")
}
