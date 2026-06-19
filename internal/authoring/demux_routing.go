package authoring

import (
	"context"
	"fmt"
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
	index demuxStackIndex
}

func (e *Engine) planDemuxRoutes(ctx context.Context, proposal DemuxProposal) (DemuxProposal, error) {
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return DemuxProposal{}, err
	}
	return router.plan(proposal), nil
}

func (e *Engine) reviewDemuxRoutes(ctx context.Context, proposal DemuxProposal) (demuxRouteReview, error) {
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return demuxRouteReview{}, err
	}
	return router.review(proposal), nil
}

func (e *Engine) demuxRouter(ctx context.Context, repoRoot string) (demuxRouter, error) {
	index, err := e.demuxStackIndex(ctx, repoRoot)
	if err != nil {
		return demuxRouter{}, err
	}
	return demuxRouter{index: index}, nil
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
	stacks []storage.Stack
	byKey  map[string]storage.Stack
}

func (e *Engine) demuxStackIndex(ctx context.Context, repoRoot string) (demuxStackIndex, error) {
	if strings.TrimSpace(repoRoot) == "" {
		return demuxStackIndex{byKey: map[string]storage.Stack{}}, nil
	}
	db, err := storage.Open(ctx)
	if err != nil {
		return demuxStackIndex{}, err
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return demuxStackIndex{}, err
	}
	repo, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return demuxStackIndex{}, err
	}
	if repo == nil {
		return demuxStackIndex{byKey: map[string]storage.Stack{}}, nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repo.ID)
	if err != nil {
		return demuxStackIndex{}, err
	}
	index := demuxStackIndex{
		stacks: stacks,
		byKey:  map[string]storage.Stack{},
	}
	for i, stack := range stacks {
		index.add(stack.Name, stack)
		index.add(stack.BookmarkName, stack)
		index.add(stackAlias(i), stack)
		index.add(stackNameFromBookmark(stack.BookmarkName), stack)
	}
	return index, nil
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
	intent := strings.ToLower(strings.TrimSpace(revision.Intent))
	if intent == "" {
		return nil, 0, ""
	}
	var best *storage.Stack
	bestScore := 0.0
	bestReason := ""
	for _, stack := range index.stacks {
		for _, candidate := range []string{stack.Name, stack.BookmarkName, stackNameFromBookmark(stack.BookmarkName)} {
			score := routeTextScore(intent, candidate)
			if score > bestScore {
				copy := stack
				best = &copy
				bestScore = score
				bestReason = fmt.Sprintf("intent matched stack %q", candidate)
			}
		}
	}
	if bestScore < 0.72 {
		return nil, 0, ""
	}
	return best, bestScore, bestReason
}

type demuxStackCluster struct {
	Key       string
	Name      string
	Kind      string
	Bookmark  string
	Reason    string
	Score     float64
	Revisions []int
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
			revision.RouteConfidence = cluster.Score
		}
	}
	return normalizeConventionalDemuxStackRoutes(proposal)
}

func demuxNewStackClusters(proposal DemuxProposal) []demuxStackCluster {
	byKey := map[string]*demuxStackCluster{}
	order := []string{}
	for index, revision := range proposal.Revisions {
		if hasDemuxRoute(revision) {
			continue
		}
		key, name, kind, reason, score := demuxStackClusterForRevision(revision)
		if key == "" {
			continue
		}
		cluster := byKey[key]
		if cluster == nil {
			cluster = &demuxStackCluster{
				Key:      key,
				Name:     name,
				Kind:     kind,
				Bookmark: conventionalDemuxStackBookmark(kind, name),
				Reason:   reason,
				Score:    score,
			}
			byKey[key] = cluster
			order = append(order, key)
		}
		cluster.Revisions = append(cluster.Revisions, index)
	}
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

func demuxStackClusterForRevision(revision RevisionProposal) (string, string, string, string, float64) {
	text := strings.ToLower(strings.TrimSpace(revision.Intent + " " + strings.Join(revision.Files, " ")))
	files := revisionFilesForRouting(revision)
	kind := conventionalDemuxStackKind(text, files)
	if strings.Contains(text, "demux") {
		return "demux-routing", "demux routing", kind, "revision touches demux planning, repair, or routed apply", 0.82
	}
	if strings.Contains(text, "stack") || strings.Contains(text, "bookmark") ||
		containsPathPrefix(files, "internal/vcs/") ||
		containsPathPrefix(files, "internal/storage/") ||
		containsPathPrefix(files, "internal/cli/") && strings.Contains(text, "stack") {
		return "stack-management", "stack management", kind, "revision touches stack lifecycle, storage, or CLI behavior", 0.78
	}
	if containsPathPrefix(files, "apps/desktop/") {
		return "desktop-app", "desktop app", kind, "revision touches the desktop app surface", 0.78
	}
	if containsPathPrefix(files, "internal/providers/") || containsPathPrefix(files, "internal/daemon/") || containsPathPrefix(files, "internal/reviewbundle/") {
		return "provider-runtime", "provider runtime", kind, "revision touches provider, daemon, or review bundle runtime", 0.74
	}
	if containsPathPrefix(files, "docs/") || containsPathPrefix(files, "skills/") || containsMarkdownFile(files) {
		return "documentation", "documentation", kind, "revision touches documentation", 0.72
	}
	if containsPathPrefix(files, "cmd/gx/") || containsPathPrefix(files, "internal/cli/") || containsPathPrefix(files, "internal/clitui/") {
		return "cli", "CLI", kind, "revision touches the GX command-line surface", 0.72
	}
	if containsPathPrefix(files, "test/e2e/") {
		return "e2e-tests", "e2e tests", kind, "revision touches end-to-end tests", 0.72
	}
	key := demuxClusterKeyFromFiles(files)
	if key == "" {
		return "", "", "", "", 0
	}
	name := strings.ReplaceAll(key, "-", " ")
	return key, name, kind, fmt.Sprintf("revision touches %s", name), 0.66
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
	if containsOnlyBuildOrDependencyFiles(files) {
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
		file = strings.ToLower(strings.TrimSpace(file))
		if file == "" {
			continue
		}
		sawFile = true
		if !(strings.HasSuffix(file, "_test.go") || strings.HasPrefix(file, "test/") || strings.Contains(file, "/test/") || strings.Contains(file, "/tests/")) {
			return false
		}
	}
	return sawFile
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

func demuxClusterKeyFromFiles(files []string) string {
	for _, file := range files {
		file = strings.Trim(strings.TrimSpace(file), "/")
		if file == "" {
			continue
		}
		parts := strings.Split(file, "/")
		switch {
		case len(parts) >= 2 && parts[0] == "internal":
			return demuxRouteSlug(parts[0] + " " + parts[1])
		case len(parts) >= 2:
			return demuxRouteSlug(parts[0] + " " + parts[1])
		default:
			return demuxRouteSlug(parts[0])
		}
	}
	return ""
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
