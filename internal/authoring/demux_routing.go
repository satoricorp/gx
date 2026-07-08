package authoring

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	routeSourceHeuristic = "heuristic"
	routeSourceAI        = "ai"
	routeSourceUser      = "user"
)

type demuxRouteReview struct {
	Warnings []FeasibilityWarning
	Hints    []RepairHint
	Errors   []string
}

type demuxRouter struct {
	index              demuxStackIndex
	staleStackWarnings []string
}

func (e *Engine) planDemuxRoutes(ctx context.Context, proposal DemuxProposal) (DemuxProposal, error) {
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return DemuxProposal{}, err
	}
	proposal = router.plan(proposal)
	for _, warning := range router.staleStackWarnings {
		proposal.Warnings = append(proposal.Warnings, warning)
	}
	return proposal, nil
}

func (e *Engine) reviewDemuxRoutes(ctx context.Context, proposal DemuxProposal) (demuxRouteReview, error) {
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return demuxRouteReview{}, err
	}
	return router.review(proposal), nil
}

func (e *Engine) demuxRouter(ctx context.Context, repoRoot string) (demuxRouter, error) {
	index, staleWarnings, err := e.demuxStackIndex(ctx, repoRoot)
	if err != nil {
		return demuxRouter{}, err
	}
	return demuxRouter{index: index, staleStackWarnings: staleWarnings}, nil
}

func (r demuxRouter) plan(proposal DemuxProposal) DemuxProposal {
	if len(r.index.stacks) > 0 {
		for index := range proposal.Revisions {
			revision := &proposal.Revisions[index]
			if hasDemuxRoute(*revision) {
				continue
			}
			stack, score, reason := r.routeForRevision(*revision)
			if stack == nil {
				continue
			}
			revision.TargetStack = stack.BookmarkName
			revision.BaseStack = stack.BaseRef
			revision.RouteReason = reason
			revision.RouteConfidence = score
			revision.RouteSource = routeSourceHeuristic
		}
	}
	return inferNewStackRoutes(proposal)
}

func (r demuxRouter) review(proposal DemuxProposal) demuxRouteReview {
	var review demuxRouteReview
	knownRevision := map[string]RevisionProposal{}
	for _, revision := range proposal.Revisions {
		knownRevision[revision.ID] = revision
	}
	for _, revision := range proposal.Revisions {
		target := routeTargetStack(revision)
		if target == "" {
			continue
		}
		stack := r.lookup(target)
		if stack == nil {
			review.Warnings = append(review.Warnings, FeasibilityWarning{
				RevisionID: revision.ID,
				Severity:   "info",
				Source:     "demux_route_new_stack",
				Message:    fmt.Sprintf("revision %s routes to new stack %q", revision.ID, target),
			})
			continue
		}
		explicitStack := strings.TrimSpace(revision.TargetStack)
		if explicitStack != "" && explicitStack != stack.BookmarkName && explicitStack != stack.Name {
			message := fmt.Sprintf("revision %s target_stack %q does not match resolved stack %q", revision.ID, explicitStack, stack.BookmarkName)
			review.Errors = append(review.Errors, message)
			review.Hints = append(review.Hints, RepairHint{
				Kind:        "invalid_demux_route",
				RevisionID:  revision.ID,
				TargetStack: stack.BookmarkName,
				Candidates:  r.candidates(),
				Suggestion:  fmt.Sprintf("use target_stack %q for %s or choose another existing stack", stack.BookmarkName, revision.ID),
			})
			continue
		}
		if strings.TrimSpace(revision.RouteSource) == "" {
			review.Warnings = append(review.Warnings, FeasibilityWarning{
				RevisionID: revision.ID,
				Severity:   "info",
				Source:     "demux_route",
				Message:    fmt.Sprintf("revision %s has a target route without route_source; mark it heuristic, ai, or user", revision.ID),
			})
		}
		for _, dep := range revision.DependsOn {
			dependency, ok := knownRevision[dep]
			if !ok {
				continue
			}
			depTarget := routeTargetStack(dependency)
			if depTarget == "" || depTarget == target {
				continue
			}
			if routeBaseStack(revision) == depTarget {
				continue
			}
			review.Warnings = append(review.Warnings, FeasibilityWarning{
				RevisionID: revision.ID,
				Severity:   "warning",
				Source:     "demux_route_cross_stack_dependency",
				DependsOn:  dep,
				Message:    fmt.Sprintf("revision %s depends on %s but routes to %q while dependency routes to %q", revision.ID, dep, target, depTarget),
			})
			review.Hints = append(review.Hints, RepairHint{
				Kind:       "cross_stack_dependency",
				RevisionID: revision.ID,
				DependsOn:  dep,
				BaseStack:  depTarget,
				Candidates: []string{target, depTarget},
				Suggestion: fmt.Sprintf("keep %s and %s on one stack, or set %s.base_stack to %q", revision.ID, dep, revision.ID, depTarget),
			})
		}
	}
	return review
}

func hasDemuxRoute(revision RevisionProposal) bool {
	return strings.TrimSpace(revision.TargetStack) != "" ||
		strings.TrimSpace(revision.BaseStack) != ""
}

func hasNonCurrentDemuxRoute(revision RevisionProposal) bool {
	return strings.TrimSpace(revision.TargetStack) != ""
}

func routeTargetStack(revision RevisionProposal) string {
	return strings.TrimSpace(revision.TargetStack)
}

func routeBaseStack(revision RevisionProposal) string {
	return strings.TrimSpace(revision.BaseStack)
}

func newStackRouteName(revision RevisionProposal) string {
	return firstNonEmpty(
		stackNameFromBookmark(revision.TargetStack),
		strings.TrimSpace(revision.TargetStack),
		strings.TrimSpace(revision.Intent),
	)
}

func newStackRouteBookmark(revision RevisionProposal) string {
	candidate := strings.TrimSpace(revision.TargetStack)
	if strings.Contains(candidate, "/") {
		return candidate
	}
	return ""
}

func (r demuxRouter) lookup(key string) *storage.Stack {
	return r.index.lookup(key)
}

func (r demuxRouter) candidates() []string {
	return r.index.candidates()
}

func (r demuxRouter) routeForRevision(revision RevisionProposal) (*storage.Stack, float64, string) {
	return stackRouteForRevision(revision, r.index)
}

type demuxStackIndex struct {
	stacks       []storage.Stack
	byKey        map[string]storage.Stack
	stackChanges map[int64][]storage.Change
	stackFiles   map[int64][]string
}

func (e *Engine) demuxStackIndex(ctx context.Context, repoRoot string) (demuxStackIndex, []string, error) {
	if strings.TrimSpace(repoRoot) == "" {
		return demuxStackIndex{byKey: map[string]storage.Stack{}, stackChanges: map[int64][]storage.Change{}, stackFiles: map[int64][]string{}}, nil, nil
	}
	db, err := storage.Open(ctx)
	if err != nil {
		return demuxStackIndex{}, nil, err
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return demuxStackIndex{}, nil, err
	}
	repo, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return demuxStackIndex{}, nil, err
	}
	if repo == nil {
		return demuxStackIndex{byKey: map[string]storage.Stack{}, stackChanges: map[int64][]storage.Change{}, stackFiles: map[int64][]string{}}, nil, nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repo.ID)
	if err != nil {
		return demuxStackIndex{}, nil, err
	}
	index := demuxStackIndex{
		stacks:       make([]storage.Stack, 0, len(stacks)),
		byKey:        map[string]storage.Stack{},
		stackChanges: map[int64][]storage.Change{},
		stackFiles:   map[int64][]string{},
	}
	staleWarnings := make([]string, 0)
	warnedBookmarks := map[string]struct{}{}
	for _, stack := range stacks {
		bookmark := strings.TrimSpace(stack.BookmarkName)
		if bookmark != "" && !IsTerminalDemuxStackStatus(stack.Status) {
			exists, existsErr := e.vcs.RevisionExists(ctx, repoRoot, bookmark)
			// The probe is best-effort: if jj itself cannot answer, keep the
			// stack and let downstream operations surface any real failure.
			if existsErr == nil && !exists {
				if _, ok := warnedBookmarks[bookmark]; !ok {
					warnedBookmarks[bookmark] = struct{}{}
					staleWarnings = append(staleWarnings, demuxDedupStackSkipWarning(bookmark))
				}
				continue
			}
		}
		stackIndex := len(index.stacks)
		index.stacks = append(index.stacks, stack)
		index.add(stack.Name, stack)
		index.add(stack.BookmarkName, stack)
		index.add(stackAlias(stackIndex), stack)
		index.add(stackNameFromBookmark(stack.BookmarkName), stack)
		changes, err := store.ListChangesByStackID(ctx, stack.ID)
		if err != nil {
			return demuxStackIndex{}, nil, err
		}
		index.stackChanges[stack.ID] = changes
		files, err := stackChangedFiles(ctx, store, changes)
		if err != nil {
			return demuxStackIndex{}, nil, err
		}
		index.stackFiles[stack.ID] = files
	}
	return index, staleWarnings, nil
}

func stackChangedFiles(ctx context.Context, store *storage.Store, changes []storage.Change) ([]string, error) {
	files := []string{}
	for _, change := range changes {
		revisions, err := store.ListChangeRevisionsByChangeID(ctx, change.ID)
		if err != nil {
			return nil, err
		}
		for _, revision := range revisions {
			var changed []string
			if err := json.Unmarshal([]byte(revision.ChangedFiles), &changed); err != nil {
				continue
			}
			files = append(files, changed...)
		}
	}
	return cleanFiles(files), nil
}

func (i demuxStackIndex) add(key string, stack storage.Stack) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	i.byKey[key] = stack
	i.byKey[strings.ToLower(key)] = stack
}

func (i demuxStackIndex) lookup(key string) *storage.Stack {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil
	}
	if stack, ok := i.byKey[key]; ok {
		return &stack
	}
	if stack, ok := i.byKey[strings.ToLower(key)]; ok {
		return &stack
	}
	return nil
}

func (i demuxStackIndex) candidates() []string {
	out := make([]string, 0, len(i.stacks))
	for index, stack := range i.stacks {
		label := stack.BookmarkName
		if strings.TrimSpace(stack.Name) != "" {
			label = stack.Name + " (" + stack.BookmarkName + ")"
		}
		out = append(out, stackAlias(index)+": "+label)
	}
	return out
}

func stackRouteForRevision(revision RevisionProposal, index demuxStackIndex) (*storage.Stack, float64, string) {
	files := cleanFiles(revisionFilesForRouting(revision))
	intent := strings.TrimSpace(revision.Intent)
	if intent == "" && len(files) == 0 {
		return nil, 0, ""
	}
	var best *storage.Stack
	secondBestScore := 0.0
	bestScore := 0.0
	bestReason := ""
	for _, stack := range index.stacks {
		if IsTerminalDemuxStackStatus(stack.Status) {
			continue
		}
		score, reason := routeStackSimilarity(revision, files, stack, index.stackChanges[stack.ID], index.stackFiles[stack.ID])
		if score > bestScore {
			secondBestScore = bestScore
			copy := stack
			best = &copy
			bestScore = score
			bestReason = reason
			continue
		}
		if score > secondBestScore {
			secondBestScore = score
		}
	}
	if bestScore < 0.75 || bestScore-secondBestScore < 0.15 {
		return nil, 0, ""
	}
	return best, bestScore, bestReason
}

func IsTerminalDemuxStackStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "merged", "closed":
		return true
	default:
		return false
	}
}

func routeStackSimilarity(revision RevisionProposal, revisionFiles []string, stack storage.Stack, changes []storage.Change, stackFiles []string) (float64, string) {
	score := 0.0
	reasons := []string{}
	intent := strings.ToLower(strings.TrimSpace(revision.Intent))
	stackText := strings.ToLower(strings.TrimSpace(stack.Name + " " + stack.BookmarkName + " " + stackNameFromBookmark(stack.BookmarkName)))
	for _, change := range changes {
		stackText = strings.TrimSpace(stackText + " " + change.Description)
	}
	textScore := 0.0
	for _, candidate := range []string{stack.Name, stack.BookmarkName, stackNameFromBookmark(stack.BookmarkName)} {
		if candidateScore := routeTextScore(intent, candidate); candidateScore > textScore {
			textScore = candidateScore
		}
	}
	if textScore > 0 {
		score = textScore
		reasons = append(reasons, "intent matched stack name")
	}
	if descriptionScore := routeDescriptionScore(intent, stackText); descriptionScore > 0 {
		score += descriptionScore * 0.15
		reasons = append(reasons, "intent matched prior revisions")
	}
	if exact := exactFileOverlapScore(revisionFiles, stackFiles); exact > 0 {
		score += exact * 0.55
		reasons = append(reasons, "files overlap existing revisions")
	}
	if dir := directoryOverlapScore(revisionFiles, stackFiles); dir > 0 {
		score += dir * 0.25
		reasons = append(reasons, "directories overlap existing revisions")
	}
	if symbol := hunkHeaderOverlapScore(revision.Hunks, stackText); symbol > 0 {
		score += symbol * 0.15
		reasons = append(reasons, "hunks match prior revision text")
	}
	if semantic := semanticLabelStackOverlapScore(revision.SemanticLabels, stackText); semantic > 0 {
		score += semantic * 0.12
		reasons = append(reasons, "semantic labels matched stack context")
	}
	if score > 1 {
		score = 1
	}
	if len(reasons) == 0 {
		return 0, ""
	}
	return score, strings.Join(reasons, "; ")
}

func exactFileOverlapScore(left, right []string) float64 {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}
	set := map[string]bool{}
	for _, file := range right {
		set[file] = true
	}
	matched := 0
	for _, file := range left {
		if set[file] {
			matched++
		}
	}
	return float64(matched) / float64(len(left))
}

func directoryOverlapScore(left, right []string) float64 {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}
	set := map[string]bool{}
	for _, file := range right {
		for _, dir := range routeDirectories(file) {
			set[dir] = true
		}
	}
	matched := 0
	for _, file := range left {
		for _, dir := range routeDirectories(file) {
			if set[dir] {
				matched++
				break
			}
		}
	}
	return float64(matched) / float64(len(left))
}

func routeDirectories(file string) []string {
	file = strings.Trim(strings.TrimSpace(file), "/")
	if file == "" || !strings.Contains(file, "/") {
		return nil
	}
	parts := strings.Split(file, "/")
	out := []string{}
	for i := 1; i < len(parts); i++ {
		out = append(out, strings.Join(parts[:i], "/"))
	}
	return out
}

func routeDescriptionScore(intent, stackText string) float64 {
	intentWords := routeWords(intent)
	if len(intentWords) == 0 || strings.TrimSpace(stackText) == "" {
		return 0
	}
	matched := 0
	for _, word := range intentWords {
		if strings.Contains(stackText, word) {
			matched++
		}
	}
	return float64(matched) / float64(len(intentWords))
}

func hunkHeaderOverlapScore(hunks []HunkRange, stackText string) float64 {
	stackText = strings.ToLower(stackText)
	if len(hunks) == 0 || stackText == "" {
		return 0
	}
	matched := 0
	total := 0
	for _, hunk := range hunks {
		for _, word := range routeWords(hunk.Header) {
			total++
			if strings.Contains(stackText, word) {
				matched++
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(matched) / float64(total)
}

func routeWords(value string) []string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	})
	out := []string{}
	for _, part := range parts {
		if len(part) > 2 {
			out = append(out, part)
		}
	}
	return out
}

type demuxStackCluster struct {
	Key             string
	Name            string
	Kind            string
	Bookmark        string
	Reason          string
	Score           float64
	Weak            bool
	Revisions       []int
	RevisionReasons map[int]string
}

func inferNewStackRoutes(proposal DemuxProposal) DemuxProposal {
	clusters := demuxNewStackClusters(proposal)
	if len(clusters) == 0 {
		return normalizeConventionalDemuxStackRoutes(proposal)
	}
	for _, cluster := range clusters {
		for _, revisionIndex := range cluster.Revisions {
			revision := &proposal.Revisions[revisionIndex]
			if hasDemuxRoute(*revision) {
				continue
			}
			revision.TargetStack = cluster.Bookmark
			revision.RouteSource = routeSourceHeuristic
			revision.RouteReason = cluster.Reason
			if override, ok := cluster.RevisionReasons[revisionIndex]; ok {
				revision.RouteReason = override
			}
			revision.RouteConfidence = cluster.Score
		}
	}
	return normalizeConventionalDemuxStackRoutes(proposal)
}

func demuxNewStackClusters(proposal DemuxProposal) []demuxStackCluster {
	byKey := map[string]*demuxStackCluster{}
	order := []string{}
	ensureCluster := func(key, name, kind, reason string, score float64, weak bool) *demuxStackCluster {
		cluster := byKey[key]
		if cluster == nil {
			cluster = &demuxStackCluster{
				Key:             key,
				Name:            name,
				Kind:            kind,
				Bookmark:        conventionalDemuxStackBookmark(kind, name),
				Reason:          reason,
				Score:           score,
				Weak:            weak,
				RevisionReasons: map[int]string{},
			}
			byKey[key] = cluster
			order = append(order, key)
		}
		if !weak {
			cluster.Weak = false
		}
		return cluster
	}
	deferredTests := []int{}
	for index, revision := range proposal.Revisions {
		if hasDemuxRoute(revision) {
			continue
		}
		key, name, kind, reason, score, weak := demuxStackClusterForRevision(revision)
		if key == "" {
			continue
		}
		if weak && containsOnlyTestFiles(cleanFiles(revisionFilesForRouting(revision))) {
			deferredTests = append(deferredTests, index)
			continue
		}
		cluster := ensureCluster(key, name, kind, reason, score, weak)
		cluster.Revisions = append(cluster.Revisions, index)
	}
	attachDemuxTestRevisions(proposal, deferredTests, byKey, order, ensureCluster)
	mergeFragmentDemuxClusters(proposal, byKey, order)
	clusters := make([]demuxStackCluster, 0, len(order))
	for _, key := range order {
		cluster := byKey[key]
		if cluster != nil && len(cluster.Revisions) > 0 {
			textParts := make([]string, 0, len(cluster.Revisions))
			files := []string{}
			for _, revisionIndex := range cluster.Revisions {
				revision := proposal.Revisions[revisionIndex]
				textParts = append(textParts, revision.Intent)
				files = append(files, revisionFilesForRouting(revision)...)
			}
			cluster.Kind = conventionalDemuxStackKind(strings.Join(textParts, " "), cleanFiles(files))
			cluster.Bookmark = conventionalDemuxStackBookmark(cluster.Kind, cluster.Name)
			clusters = append(clusters, *cluster)
		}
	}
	return clusters
}

// attachDemuxTestRevisions groups test-only revisions with the cluster that
// contains the code they test (matched by naming conventions), and gathers any
// unmatched test revisions into one shared tests cluster instead of one stack
// per test file.
func attachDemuxTestRevisions(
	proposal DemuxProposal,
	testRevisions []int,
	byKey map[string]*demuxStackCluster,
	order []string,
	ensureCluster func(key, name, kind, reason string, score float64, weak bool) *demuxStackCluster,
) {
	if len(testRevisions) == 0 {
		return
	}
	subjects := map[string]string{}
	for _, key := range order {
		cluster := byKey[key]
		if cluster == nil {
			continue
		}
		for _, revisionIndex := range cluster.Revisions {
			for _, file := range cleanFiles(revisionFilesForRouting(proposal.Revisions[revisionIndex])) {
				if isTestFile(file) {
					continue
				}
				for _, subject := range codeFileSubjectKeys(file) {
					if _, ok := subjects[subject]; !ok {
						subjects[subject] = key
					}
				}
			}
		}
	}
	for _, revisionIndex := range testRevisions {
		files := cleanFiles(revisionFilesForRouting(proposal.Revisions[revisionIndex]))
		targetKey := ""
		for _, file := range files {
			for _, subject := range testSubjectKeys(file) {
				if key, ok := subjects[subject]; ok {
					targetKey = key
					break
				}
			}
			if targetKey != "" {
				break
			}
		}
		if targetKey != "" {
			cluster := byKey[targetKey]
			cluster.Revisions = append(cluster.Revisions, revisionIndex)
			cluster.RevisionReasons[revisionIndex] = "test file grouped with the code it tests"
			continue
		}
		cluster := ensureCluster("tests", "test updates", "test", "test-only changes grouped together", 0.7, false)
		cluster.Revisions = append(cluster.Revisions, revisionIndex)
	}
}

// mergeFragmentDemuxClusters is the anti-fragmentation pass: weak single-file
// clusters merge into the nearest related cluster (same directory or top-level
// subsystem, session-compatible); remaining strays merge into one shared
// cluster instead of producing N single-file stacks. Revisions with explicit
// routes never reach this pass, and clusters whose sessions are disjoint from
// a candidate's sessions are treated as intentionally isolated.
func mergeFragmentDemuxClusters(proposal DemuxProposal, byKey map[string]*demuxStackCluster, order []string) {
	leftovers := []*demuxStackCluster{}
	for _, key := range order {
		cluster := byKey[key]
		if cluster == nil || !cluster.Weak || len(cluster.Revisions) != 1 {
			continue
		}
		files := demuxClusterFiles(proposal, cluster)
		if len(files) != 1 {
			continue
		}
		target := bestFragmentMergeTarget(proposal, byKey, order, cluster, files[0])
		if target == nil {
			leftovers = append(leftovers, cluster)
			continue
		}
		moveDemuxClusterRevisions(target, cluster, fmt.Sprintf("single-file change grouped with related %s changes", target.Name))
	}
	if len(leftovers) < 2 {
		return
	}
	shared := leftovers[0]
	merged := false
	for _, cluster := range leftovers[1:] {
		if sessionSetsDisjoint(demuxClusterSessionSet(proposal, shared), demuxClusterSessionSet(proposal, cluster)) {
			continue
		}
		moveDemuxClusterRevisions(shared, cluster, "small independent change grouped with other standalone changes")
		merged = true
	}
	if merged {
		shared.Name = "misc changes"
		shared.Reason = "small independent changes grouped to avoid single-file stacks"
	}
}

func bestFragmentMergeTarget(proposal DemuxProposal, byKey map[string]*demuxStackCluster, order []string, source *demuxStackCluster, file string) *demuxStackCluster {
	sourceSessions := demuxClusterSessionSet(proposal, source)
	fileDir := packageDirectory(file)
	top := topPathSegment(file)
	var best *demuxStackCluster
	bestScore := 0
	for _, key := range order {
		candidate := byKey[key]
		if candidate == nil || candidate == source || len(candidate.Revisions) == 0 {
			continue
		}
		candidateSessions := demuxClusterSessionSet(proposal, candidate)
		if sessionSetsDisjoint(sourceSessions, candidateSessions) {
			continue
		}
		score := 0
		for _, other := range demuxClusterFiles(proposal, candidate) {
			switch {
			case packageDirectory(other) == fileDir:
				score = maxInt(score, 3)
			case top != "" && topPathSegment(other) == top:
				score = maxInt(score, 2)
			}
		}
		if score == 0 {
			continue
		}
		if sessionSetsOverlap(sourceSessions, candidateSessions) {
			score += 2
		}
		if !candidate.Weak {
			score++
		}
		if score > bestScore {
			best = candidate
			bestScore = score
		}
	}
	if bestScore < 2 {
		return nil
	}
	return best
}

func moveDemuxClusterRevisions(target, source *demuxStackCluster, reason string) {
	for _, revisionIndex := range source.Revisions {
		target.Revisions = append(target.Revisions, revisionIndex)
		target.RevisionReasons[revisionIndex] = reason
	}
	source.Revisions = nil
}

func demuxClusterFiles(proposal DemuxProposal, cluster *demuxStackCluster) []string {
	files := []string{}
	for _, revisionIndex := range cluster.Revisions {
		files = append(files, revisionFilesForRouting(proposal.Revisions[revisionIndex])...)
	}
	return cleanFiles(files)
}

func demuxClusterSessionSet(proposal DemuxProposal, cluster *demuxStackCluster) map[string]struct{} {
	out := map[string]struct{}{}
	for _, revisionIndex := range cluster.Revisions {
		for _, session := range proposal.Revisions[revisionIndex].SessionIDs {
			session = strings.TrimSpace(session)
			if session != "" {
				out[session] = struct{}{}
			}
		}
	}
	return out
}

func sessionSetsDisjoint(left, right map[string]struct{}) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	return !sessionSetsOverlap(left, right)
}

func sessionSetsOverlap(left, right map[string]struct{}) bool {
	for session := range left {
		if _, ok := right[session]; ok {
			return true
		}
	}
	return false
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func topPathSegment(file string) string {
	file = normalizeFilesetPath(file)
	if index := strings.Index(file, "/"); index > 0 {
		return file[:index]
	}
	return ""
}

func demuxStackClusterForRevision(revision RevisionProposal) (string, string, string, string, float64, bool) {
	text := strings.ToLower(strings.TrimSpace(revision.Intent + " " + strings.Join(revision.Files, " ")))
	files := revisionFilesForRouting(revision)
	kind := conventionalDemuxStackKind(text, files)
	if strings.Contains(text, "demux") {
		return "demux-routing", "demux routing", kind, "revision touches demux planning, repair, or routed apply", 0.82, false
	}
	if containsOnlyInfraFiles(files) {
		return "infra", "infra", kind, "revision touches infrastructure, build, or CI configuration", 0.74, false
	}
	if strings.Contains(text, "stack") || strings.Contains(text, "bookmark") ||
		containsPathPrefix(files, "internal/vcs/") ||
		containsPathPrefix(files, "internal/storage/") ||
		containsPathPrefix(files, "internal/cli/") && strings.Contains(text, "stack") {
		return "stack-management", "stack management", kind, "revision touches stack lifecycle, storage, or CLI behavior", 0.78, false
	}
	if containsPathPrefix(files, "apps/desktop/") {
		return "desktop-app", "desktop app", kind, "revision touches the desktop app surface", 0.78, false
	}
	if containsPathPrefix(files, "internal/providers/") || containsPathPrefix(files, "internal/daemon/") || containsPathPrefix(files, "internal/reviewbundle/") {
		return "provider-runtime", "provider runtime", kind, "revision touches provider, daemon, or review bundle runtime", 0.74, false
	}
	if containsPathPrefix(files, "docs/") || containsPathPrefix(files, "skills/") || containsMarkdownFile(files) {
		return "documentation", "documentation", kind, "revision touches documentation", 0.72, false
	}
	if containsPathPrefix(files, "src-tauri/") {
		return "tauri-app", "Tauri app", kind, "revision touches the Tauri application shell", 0.78, false
	}
	if containsPathPrefix(files, "src/") || containsExactFile(files, "index.html") {
		return "frontend-app", "frontend app", kind, "revision touches the frontend application shell", 0.76, false
	}
	if containsProjectToolingFile(files) {
		return "project-tooling", "project tooling", kind, "revision touches project tooling, dependencies, or build configuration", 0.74, false
	}
	if containsPathPrefix(files, "cmd/gx/") || containsPathPrefix(files, "internal/cli/") || containsPathPrefix(files, "internal/clitui/") {
		return "cli", "CLI", kind, "revision touches the GX command-line surface", 0.72, false
	}
	if containsPathPrefix(files, "test/e2e/") {
		return "e2e-tests", "e2e tests", kind, "revision touches end-to-end tests", 0.72, false
	}
	key, weak := demuxClusterKeyFromFiles(files)
	if key == "" {
		return "", "", "", "", 0, false
	}
	name := strings.ReplaceAll(key, "-", " ")
	return key, name, kind, fmt.Sprintf("revision touches %s", name), 0.66, weak
}

func conventionalDemuxStackKind(text string, files []string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	if strings.Contains(text, "fix") || strings.Contains(text, "bug") ||
		strings.Contains(text, "regression") || strings.Contains(text, "panic") ||
		strings.Contains(text, "crash") || strings.Contains(text, "failure") {
		return "bug"
	}
	if containsOnlyDocumentationFiles(files) {
		return "docs"
	}
	if containsOnlyTestFiles(files) {
		return "test"
	}
	if containsOnlyBuildOrDependencyFiles(files) || containsOnlyProjectToolingFiles(files) || containsOnlyInfraFiles(files) {
		return "chore"
	}
	return "feature"
}

func conventionalDemuxStackBookmark(kind, name string) string {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		kind = "feature"
	}
	slug := demuxRouteSlug(name)
	if slug == "" {
		slug = "change"
	}
	return kind + "/" + slug
}

func normalizeConventionalDemuxStackRoutes(proposal DemuxProposal) DemuxProposal {
	groups := conventionalDemuxRouteGroups(proposal.Revisions)
	for index := range proposal.Revisions {
		revision := &proposal.Revisions[index]
		revision.TargetStack = normalizeConventionalDemuxStackRoute(revision.TargetStack, *revision, groups)
		revision.BaseStack = normalizeConventionalDemuxStackRoute(revision.BaseStack, *revision, groups)
	}
	return proposal
}

type conventionalDemuxRouteGroup struct {
	Text  string
	Files []string
}

func conventionalDemuxRouteGroups(revisions []RevisionProposal) map[string]conventionalDemuxRouteGroup {
	groups := map[string]conventionalDemuxRouteGroup{}
	for _, revision := range revisions {
		route := strings.TrimSpace(revision.TargetStack)
		if !isNormalizableDemuxStackRoute(route) {
			continue
		}
		group := groups[route]
		group.Text = strings.TrimSpace(group.Text + " " + revision.Intent)
		group.Files = append(group.Files, revisionFilesForRouting(revision)...)
		groups[route] = group
	}
	for route, group := range groups {
		group.Files = cleanFiles(group.Files)
		groups[route] = group
	}
	return groups
}

func normalizeConventionalDemuxStackRoute(route string, revision RevisionProposal, groups map[string]conventionalDemuxRouteGroup) string {
	route = strings.TrimSpace(route)
	if !isNormalizableDemuxStackRoute(route) {
		return route
	}
	name := stackNameFromBookmark(route)
	text := strings.ToLower(revision.Intent + " " + name)
	files := revisionFilesForRouting(revision)
	if group, ok := groups[route]; ok {
		text = strings.ToLower(group.Text + " " + name)
		files = group.Files
	}
	return conventionalDemuxStackBookmark(conventionalDemuxStackKind(text, files), name)
}

func isNormalizableDemuxStackRoute(route string) bool {
	route = strings.TrimSpace(route)
	if strings.HasPrefix(route, "gx/") {
		return true
	}
	kind, _, ok := strings.Cut(route, "/")
	return ok && isConventionalDemuxStackKind(kind)
}

func isConventionalDemuxStackKind(kind string) bool {
	switch strings.TrimSpace(kind) {
	case "feature", "bug", "docs", "test", "chore":
		return true
	default:
		return false
	}
}

func revisionFilesForRouting(revision RevisionProposal) []string {
	files := append([]string(nil), revision.Files...)
	if len(files) == 0 {
		for _, hunk := range revision.Hunks {
			if strings.TrimSpace(hunk.File) != "" {
				files = append(files, hunk.File)
			}
		}
	}
	return files
}

func containsPathPrefix(files []string, prefix string) bool {
	for _, file := range files {
		if strings.HasPrefix(strings.TrimSpace(file), prefix) {
			return true
		}
	}
	return false
}

func containsExactFile(files []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, file := range files {
		if strings.TrimSpace(file) == target {
			return true
		}
	}
	return false
}

func containsMarkdownFile(files []string) bool {
	for _, file := range files {
		file = strings.ToLower(strings.TrimSpace(file))
		if strings.HasSuffix(file, ".md") {
			return true
		}
	}
	return false
}

func containsOnlyDocumentationFiles(files []string) bool {
	sawFile := false
	for _, file := range files {
		file = strings.ToLower(strings.TrimSpace(file))
		if file == "" {
			continue
		}
		sawFile = true
		if !(strings.HasPrefix(file, "docs/") || strings.HasPrefix(file, "skills/") || strings.HasSuffix(file, ".md") || strings.HasSuffix(file, ".mdx")) {
			return false
		}
	}
	return sawFile
}

func containsOnlyTestFiles(files []string) bool {
	sawFile := false
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}
		sawFile = true
		if !isTestFile(file) {
			return false
		}
	}
	return sawFile
}

func isTestFile(file string) bool {
	file = strings.ToLower(normalizeFilesetPath(file))
	if file == "" {
		return false
	}
	base := path.Base(file)
	name := strings.TrimSuffix(base, path.Ext(base))
	switch {
	case strings.HasSuffix(name, "_test"), strings.HasSuffix(name, ".test"), strings.HasSuffix(name, ".spec"), strings.HasSuffix(name, "-test"):
		return true
	case strings.HasPrefix(name, "test_"):
		return true
	}
	for _, segment := range strings.Split(path.Dir(file), "/") {
		switch segment {
		case "test", "tests", "__tests__", "testdata":
			return true
		}
	}
	return false
}

// testSubjectKeys returns identifiers a test file's name implies for the code
// under test, e.g. test/http-client.test.ts -> "http-client", foo_test.go ->
// "foo", test_foo.py -> "foo", __tests__/foo.tsx -> "foo".
func testSubjectKeys(file string) []string {
	file = strings.ToLower(normalizeFilesetPath(file))
	base := path.Base(file)
	name := strings.TrimSuffix(base, path.Ext(base))
	// Cover double extensions such as foo.test.ts and foo.spec.tsx.
	for _, suffix := range []string{"_test", ".test", ".spec", "-test", "_spec"} {
		name = strings.TrimSuffix(name, suffix)
	}
	name = strings.TrimPrefix(name, "test_")
	slug := demuxRouteSlug(name)
	if slug == "" {
		return nil
	}
	return []string{slug}
}

// codeFileSubjectKeys returns identifiers that a changed code file answers to,
// so tests can find their subject across directories: the base name, the
// parent directory name, and their combination (src/http/client.ts ->
// "client", "http", "http-client").
func codeFileSubjectKeys(file string) []string {
	file = strings.ToLower(normalizeFilesetPath(file))
	base := path.Base(file)
	name := strings.TrimSuffix(base, path.Ext(base))
	dir := path.Base(path.Dir(file))
	keys := []string{}
	seen := map[string]struct{}{}
	add := func(value string) {
		slug := demuxRouteSlug(value)
		if slug == "" {
			return
		}
		if _, ok := seen[slug]; ok {
			return
		}
		seen[slug] = struct{}{}
		keys = append(keys, slug)
	}
	add(name)
	if dir != "." && dir != "/" {
		add(dir)
		if name != "" {
			add(dir + " " + name)
		}
	}
	return keys
}

func containsOnlyBuildOrDependencyFiles(files []string) bool {
	sawFile := false
	for _, file := range files {
		file = strings.ToLower(strings.TrimSpace(file))
		if file == "" {
			continue
		}
		sawFile = true
		if !isBuildOrDependencyFile(file) {
			return false
		}
	}
	return sawFile
}

func isBuildOrDependencyFile(file string) bool {
	switch file {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "bun.lock", "bun.lockb", "pnpm-lock.yaml", "yarn.lock":
		return true
	}
	return strings.HasSuffix(file, ".config.ts") ||
		strings.HasSuffix(file, ".config.js") ||
		strings.HasSuffix(file, ".toml") ||
		strings.HasSuffix(file, ".yaml") ||
		strings.HasSuffix(file, ".yml")
}

func containsBuildOrDependencyFile(files []string) bool {
	for _, file := range files {
		if isBuildOrDependencyFile(strings.ToLower(strings.TrimSpace(file))) {
			return true
		}
	}
	return false
}

func containsOnlyProjectToolingFiles(files []string) bool {
	sawFile := false
	for _, file := range files {
		file = strings.ToLower(strings.TrimSpace(file))
		if file == "" {
			continue
		}
		sawFile = true
		if !isBuildOrDependencyFile(file) && !isProjectToolingFile(file) {
			return false
		}
	}
	return sawFile
}

func containsProjectToolingFile(files []string) bool {
	for _, file := range files {
		if isProjectToolingFile(strings.ToLower(strings.TrimSpace(file))) {
			return true
		}
	}
	return false
}

func isProjectToolingFile(file string) bool {
	switch {
	case file == ".gitignore", file == ".vscode/extensions.json":
		return true
	case strings.HasPrefix(file, "tsconfig"):
		return true
	case !strings.Contains(file, "/") && isRootMetaFile(path.Base(file)):
		return true
	default:
		return false
	}
}

// isRootMetaFile recognizes repository-meta files that belong with project
// tooling rather than in their own single-file stacks. Callers pass a
// lowercased base name.
func isRootMetaFile(base string) bool {
	switch base {
	case "license", "license.md", "license.txt", "notice", "codeowners",
		".editorconfig", ".nvmrc", ".npmrc", ".prettierrc", ".prettierignore",
		".eslintignore", ".ds_store":
		return true
	default:
		return false
	}
}

func containsOnlyInfraFiles(files []string) bool {
	sawFile := false
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}
		sawFile = true
		if !isInfraFile(file) {
			return false
		}
	}
	return sawFile
}

// isInfraFile recognizes infrastructure, container, CI, and deploy
// configuration so those files group into one infra stack instead of one
// stack per file.
func isInfraFile(file string) bool {
	file = strings.ToLower(normalizeFilesetPath(file))
	if file == "" {
		return false
	}
	base := path.Base(file)
	switch {
	case strings.HasPrefix(base, "dockerfile"), strings.HasSuffix(base, ".dockerfile"),
		strings.HasPrefix(base, "docker-compose"), base == ".dockerignore":
		return true
	case base == "makefile", base == "justfile", base == "jenkinsfile", base == "procfile":
		return true
	case strings.HasSuffix(base, ".tf"), strings.HasSuffix(base, ".tfvars"):
		return true
	case base == ".gitlab-ci.yml", base == ".travis.yml", base == "cloudbuild.yaml",
		base == "cdk.json", base == "cdk.context.json",
		base == "serverless.yml", base == "serverless.yaml",
		base == "fly.toml", base == "render.yaml", base == "netlify.toml", base == "vercel.json":
		return true
	}
	for _, prefix := range []string{
		".github/", ".circleci/", ".buildkite/", ".ci/",
		"terraform/", "helm/", "k8s/", "kubernetes/", "ansible/",
		"deploy/", "deployments/", "infra/", "infrastructure/",
	} {
		if strings.HasPrefix(file, prefix) {
			return true
		}
	}
	return false
}

func demuxClusterKeyFromFiles(files []string) (string, bool) {
	if containsPathPrefix(files, "src-tauri/") {
		return "tauri-app", false
	}
	if containsPathPrefix(files, "src/") || containsExactFile(files, "index.html") {
		return "frontend-app", false
	}
	if containsProjectToolingFile(files) {
		return "project-tooling", false
	}
	for _, file := range files {
		file = strings.Trim(strings.TrimSpace(file), "/")
		if file == "" {
			continue
		}
		// Key on directories only; the last path segment is a file name and
		// keying on it produces one cluster (and one stack) per file.
		parts := strings.Split(file, "/")
		switch {
		case len(parts) >= 3:
			return demuxRouteSlug(parts[0] + " " + parts[1]), true
		case len(parts) == 2:
			return demuxRouteSlug(parts[0]), true
		default:
			return demuxRouteSlug(parts[0]), true
		}
	}
	return "", false
}

func demuxRouteSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func routeTextScore(intent, candidate string) float64 {
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	if candidate == "" {
		return 0
	}
	slug := strings.ReplaceAll(stackNameFromBookmark(candidate), "-", " ")
	switch {
	case intent == candidate || intent == slug:
		return 1
	case strings.Contains(intent, candidate):
		return 0.9
	case strings.Contains(intent, slug):
		return 0.85
	default:
		words := strings.Fields(slug)
		if len(words) == 0 {
			return 0
		}
		matched := 0
		for _, word := range words {
			if len(word) > 2 && strings.Contains(intent, word) {
				matched++
			}
		}
		if matched == 0 {
			return 0
		}
		return float64(matched) / float64(len(words)) * 0.8
	}
}

func stackAlias(index int) string {
	return fmt.Sprintf("s%d", index+1)
}

func stackNameFromBookmark(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "gx/draft/")
	name = strings.TrimPrefix(name, "gx/")
	if kind, rest, ok := strings.Cut(name, "/"); ok && isConventionalDemuxStackKind(kind) {
		return rest
	}
	return name
}
