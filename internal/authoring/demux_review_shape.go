package authoring

import (
	"fmt"
	"sort"
	"strings"
)

type demuxReviewShapePolicy struct {
	TinyEffectiveLOC int
	TargetMinLOC     int
	TargetMaxLOC     int
	SoftMaxLOC       int
	ClusterMaxLOC    int
}

var defaultDemuxReviewShapePolicy = demuxReviewShapePolicy{
	TinyEffectiveLOC: 5,
	TargetMinLOC:     20,
	TargetMaxLOC:     80,
	SoftMaxLOC:       150,
	ClusterMaxLOC:    220,
}

func shapeDemuxProposalForReview(proposal DemuxProposal) (DemuxProposal, bool) {
	return shapeDemuxProposalWithPolicy(proposal, defaultDemuxReviewShapePolicy)
}

func shapeDemuxProposalWithPolicy(proposal DemuxProposal, policy demuxReviewShapePolicy) (DemuxProposal, bool) {
	if len(proposal.Revisions) == 0 {
		return proposal, false
	}
	normalized, normalizedRoutes := normalizeDemuxReviewShapeRoutes(proposal)
	if normalizedRoutes {
		proposal = normalized
	}
	changed := annotateDemuxRevisionShape(proposal, policy)
	shaped, coalesced := coalesceTinyDemuxRevisions(proposal, policy)
	if coalesced {
		proposal = shaped
		changed = true
		changed = annotateDemuxRevisionShape(proposal, policy) || changed
	}
	shaped, semanticCoalesced := coalesceSemanticDemuxRevisions(proposal, policy)
	if semanticCoalesced {
		proposal = shaped
		changed = true
		changed = annotateDemuxRevisionShape(proposal, policy) || changed
	}
	return proposal, changed || normalizedRoutes
}

func normalizeDemuxReviewShapeRoutes(proposal DemuxProposal) (DemuxProposal, bool) {
	groups := reviewShapeDirectoryGroups(proposal.Revisions)
	if len(groups) == 0 {
		return proposal, false
	}
	changed := false
	for _, indexes := range groups {
		if len(indexes) < 2 || !reviewShapeGroupHasInternalEdges(proposal, indexes) {
			continue
		}
		target, base, ok := dominantReviewShapeRoute(proposal.Revisions, indexes)
		if !ok {
			continue
		}
		for _, index := range indexes {
			revision := &proposal.Revisions[index]
			if strings.TrimSpace(revision.TargetStack) == target && strings.TrimSpace(revision.BaseStack) == base {
				continue
			}
			revision.TargetStack = target
			revision.BaseStack = base
			revision.RouteSource = routeSourceHeuristic
			revision.RouteReason = "review shape: same package and structural dependency cluster"
			if revision.RouteConfidence < 0.76 {
				revision.RouteConfidence = 0.76
			}
			appendShapeReason(revision, "review shape: normalized route to package cluster "+target)
			changed = true
		}
	}
	return proposal, changed
}

func reviewShapeDirectoryGroups(revisions []RevisionProposal) map[string][]int {
	groups := map[string][]int{}
	for index, revision := range revisions {
		dir, ok := commonRevisionDirectory(revision)
		if !ok {
			continue
		}
		groups[dir] = append(groups[dir], index)
	}
	return groups
}

func commonRevisionDirectory(revision RevisionProposal) (string, bool) {
	files := cleanFiles(revision.Files)
	if len(files) == 0 {
		return "", false
	}
	dir := packageDirectory(files[0])
	if dir == "" {
		return "", false
	}
	for _, file := range files[1:] {
		if packageDirectory(file) != dir {
			return "", false
		}
	}
	return dir, true
}

func packageDirectory(file string) string {
	file = normalizeFilesetPath(file)
	if file == "" {
		return ""
	}
	index := strings.LastIndex(file, "/")
	if index < 0 {
		return ""
	}
	return file[:index]
}

func reviewShapeGroupHasInternalEdges(proposal DemuxProposal, indexes []int) bool {
	files := map[string]struct{}{}
	for _, index := range indexes {
		for _, file := range proposal.Revisions[index].Files {
			files[file] = struct{}{}
		}
	}
	for _, dep := range proposal.StructuralDeps {
		_, hasFrom := files[dep.FromFile]
		_, hasTo := files[dep.ToFile]
		if hasFrom && hasTo {
			return true
		}
	}
	return false
}

func dominantReviewShapeRoute(revisions []RevisionProposal, indexes []int) (string, string, bool) {
	type routeVote struct {
		target string
		base   string
		count  int
	}
	votes := map[string]routeVote{}
	for _, index := range indexes {
		target := strings.TrimSpace(revisions[index].TargetStack)
		if target == "" {
			continue
		}
		base := strings.TrimSpace(revisions[index].BaseStack)
		key := target + "\x00" + base
		vote := votes[key]
		vote.target = target
		vote.base = base
		vote.count++
		votes[key] = vote
	}
	best := routeVote{}
	for _, vote := range votes {
		if vote.count > best.count {
			best = vote
		}
	}
	if best.target == "" || best.count < 2 {
		return "", "", false
	}
	return best.target, best.base, true
}

func annotateDemuxRevisionShape(proposal DemuxProposal, policy demuxReviewShapePolicy) bool {
	changed := false
	for index := range proposal.Revisions {
		revision := &proposal.Revisions[index]
		metrics := demuxRevisionMetrics(proposal, *revision)
		if revision.EffectiveLOC != metrics.EffectiveLOC {
			revision.EffectiveLOC = metrics.EffectiveLOC
			changed = true
		}
		reason := reviewShapeReason(metrics, policy)
		if reason != "" && appendShapeReason(revision, reason) {
			changed = true
		}
	}
	return changed
}

type demuxRevisionShapeMetrics struct {
	EffectiveLOC int
	FileCount    int
	HunkCount    int
	Symbols      []string
}

func demuxRevisionMetrics(proposal DemuxProposal, revision RevisionProposal) demuxRevisionShapeMetrics {
	hunks := revisionHunksForMetrics(proposal, revision)
	symbols := map[string]struct{}{}
	effectiveLOC := 0
	for _, hunk := range hunks {
		effectiveLOC += hunkEffectiveLOC(hunk)
		if strings.TrimSpace(hunk.Symbol) != "" {
			symbols[hunk.Symbol] = struct{}{}
		}
	}
	out := demuxRevisionShapeMetrics{
		EffectiveLOC: effectiveLOC,
		FileCount:    len(cleanFiles(revision.Files)),
		HunkCount:    len(hunks),
		Symbols:      make([]string, 0, len(symbols)),
	}
	for symbol := range symbols {
		out.Symbols = append(out.Symbols, symbol)
	}
	sort.Strings(out.Symbols)
	return out
}

func revisionHunksForMetrics(proposal DemuxProposal, revision RevisionProposal) []HunkRange {
	if len(revision.Hunks) > 0 {
		return append([]HunkRange(nil), revision.Hunks...)
	}
	if revision.UseHunks && len(revision.HunkIDs) > 0 {
		out := make([]HunkRange, 0, len(revision.HunkIDs))
		for _, id := range revision.HunkIDs {
			if hunk := findRevisionHunkByID(proposal.Hunks, id); hunk.ID != "" {
				out = append(out, hunk)
			}
		}
		return out
	}
	return hunksForFiles(hunksByFile(proposal.Hunks), revision.Files)
}

func hunkEffectiveLOC(hunk HunkRange) int {
	loc := changedLinesInPatch(hunk.Patch)
	if loc > 0 {
		return loc
	}
	return hunk.OldLines + hunk.NewLines
}

func changedLinesInPatch(patch string) int {
	loc := 0
	for _, line := range strings.Split(patch, "\n") {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			continue
		case strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-"):
			loc++
		}
	}
	return loc
}

func reviewShapeReason(metrics demuxRevisionShapeMetrics, policy demuxReviewShapePolicy) string {
	switch {
	case metrics.EffectiveLOC > 0 && metrics.EffectiveLOC < policy.TinyEffectiveLOC:
		return fmt.Sprintf("review shape: tiny revision (%d effective LOC; merge unless independently reviewable)", metrics.EffectiveLOC)
	case metrics.EffectiveLOC > policy.SoftMaxLOC:
		return fmt.Sprintf("review shape: large revision (%d effective LOC; soft max %d)", metrics.EffectiveLOC, policy.SoftMaxLOC)
	case metrics.EffectiveLOC >= policy.TargetMinLOC && metrics.EffectiveLOC <= policy.TargetMaxLOC:
		return fmt.Sprintf("review shape: target-sized revision (%d effective LOC)", metrics.EffectiveLOC)
	default:
		return ""
	}
}

func coalesceTinyDemuxRevisions(proposal DemuxProposal, policy demuxReviewShapePolicy) (DemuxProposal, bool) {
	if len(proposal.Revisions) < 2 {
		return proposal, false
	}
	revisions := append([]RevisionProposal(nil), proposal.Revisions...)
	changed := false
	for {
		sourceIndex, targetIndex, reason, ok := nextTinyDemuxRevisionMerge(proposal, revisions, policy)
		if !ok {
			break
		}
		sourceID := revisions[sourceIndex].ID
		targetID := revisions[targetIndex].ID
		revisions[targetIndex] = mergeDemuxRevision(proposal, revisions[targetIndex], revisions[sourceIndex])
		appendShapeReason(&revisions[targetIndex], fmt.Sprintf("review shape: merged tiny revision %s into %s (%s)", sourceID, targetID, reason))
		revisions = removeRevisionAt(revisions, sourceIndex)
		revisions = rewriteDemuxRevisionDependencies(revisions, sourceID, targetID)
		changed = true
	}
	proposal.Revisions = revisions
	return proposal, changed
}

func coalesceSemanticDemuxRevisions(proposal DemuxProposal, policy demuxReviewShapePolicy) (DemuxProposal, bool) {
	if len(proposal.Revisions) < 2 {
		return proposal, false
	}
	revisions := append([]RevisionProposal(nil), proposal.Revisions...)
	changed := false
	for {
		sourceIndex, targetIndex, reason, ok := nextSemanticDemuxRevisionMerge(proposal, revisions, policy)
		if !ok {
			break
		}
		sourceID := revisions[sourceIndex].ID
		targetID := revisions[targetIndex].ID
		revisions[targetIndex] = mergeDemuxRevision(proposal, revisions[targetIndex], revisions[sourceIndex])
		appendShapeReason(&revisions[targetIndex], fmt.Sprintf("review shape: merged related revision %s into %s (%s)", sourceID, targetID, reason))
		revisions = removeRevisionAt(revisions, sourceIndex)
		revisions = rewriteDemuxRevisionDependencies(revisions, sourceID, targetID)
		changed = true
	}
	proposal.Revisions = revisions
	return proposal, changed
}

func nextSemanticDemuxRevisionMerge(proposal DemuxProposal, revisions []RevisionProposal, policy demuxReviewShapePolicy) (int, int, string, bool) {
	bestSource := -1
	bestTarget := -1
	bestScore := 0
	bestDistance := len(revisions) + 1
	bestReason := ""
	for sourceIndex, source := range revisions {
		sourceLOC := demuxRevisionMetrics(proposal, source).EffectiveLOC
		if sourceLOC <= 0 || sourceLOC > policy.SoftMaxLOC || sourceLOC > policy.ClusterMaxLOC {
			continue
		}
		for targetIndex, target := range revisions {
			if targetIndex == sourceIndex || !sameDemuxReviewRoute(source, target) {
				continue
			}
			targetLOC := demuxRevisionMetrics(proposal, target).EffectiveLOC
			if targetLOC <= 0 || targetLOC > policy.SoftMaxLOC || sourceLOC+targetLOC > policy.ClusterMaxLOC {
				continue
			}
			score, reason := semanticRevisionMergeScore(proposal, source, target)
			if score < 6 {
				continue
			}
			distance := absInt(targetIndex - sourceIndex)
			if distance == 1 {
				score++
			}
			if score > bestScore || (score == bestScore && distance < bestDistance) {
				bestSource = sourceIndex
				bestTarget = targetIndex
				bestScore = score
				bestDistance = distance
				bestReason = reason
			}
		}
	}
	if bestSource < 0 {
		return -1, -1, "", false
	}
	return bestSource, bestTarget, bestReason, true
}

func semanticRevisionMergeScore(proposal DemuxProposal, source, target RevisionProposal) (int, string) {
	if sameRevisionPackage(source, target) && revisionsHaveStructuralEdge(proposal, source, target) {
		return 7, "same package structural cluster"
	}
	score, reason := tinyRevisionMergeScore(proposal, source, target)
	if score > 0 {
		return score, reason
	}
	if sameRevisionPackage(source, target) && revisionTouchesTests(source, target) && sharedReviewShapeVocabulary(source, target) {
		return 6, "same package test/support cluster"
	}
	if sameRevisionPackage(source, target) && sharedReviewShapeVocabulary(source, target) {
		return 5, "same package vocabulary"
	}
	return 0, ""
}

func sameRevisionPackage(left, right RevisionProposal) bool {
	leftDir, leftOK := commonRevisionDirectory(left)
	rightDir, rightOK := commonRevisionDirectory(right)
	return leftOK && rightOK && leftDir == rightDir
}

func revisionTouchesTests(left, right RevisionProposal) bool {
	for _, revision := range []RevisionProposal{left, right} {
		for _, file := range revision.Files {
			if strings.Contains(file, "_test.") || strings.Contains(file, ".test.") || strings.Contains(file, ".spec.") {
				return true
			}
		}
	}
	return false
}

func sharedReviewShapeVocabulary(left, right RevisionProposal) bool {
	leftTokens := reviewShapeTokens(left)
	rightTokens := reviewShapeTokens(right)
	for token := range leftTokens {
		if _, ok := rightTokens[token]; ok {
			return true
		}
	}
	return false
}

func reviewShapeTokens(revision RevisionProposal) map[string]struct{} {
	tokens := map[string]struct{}{}
	text := strings.ToLower(strings.TrimSpace(revision.Intent + " " + strings.Join(revision.Files, " ")))
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	}) {
		if len(token) < 4 {
			continue
		}
		switch token {
		case "internal", "authoring", "test", "tests", "update":
			continue
		default:
			tokens[token] = struct{}{}
		}
	}
	return tokens
}

func nextTinyDemuxRevisionMerge(proposal DemuxProposal, revisions []RevisionProposal, policy demuxReviewShapePolicy) (int, int, string, bool) {
	for sourceIndex, revision := range revisions {
		metrics := demuxRevisionMetrics(proposal, revision)
		if metrics.EffectiveLOC <= 0 || metrics.EffectiveLOC >= policy.TinyEffectiveLOC {
			continue
		}
		targetIndex, reason, ok := bestTinyRevisionMergeTarget(proposal, revisions, sourceIndex, policy)
		if ok {
			return sourceIndex, targetIndex, reason, true
		}
	}
	return -1, -1, "", false
}

func bestTinyRevisionMergeTarget(proposal DemuxProposal, revisions []RevisionProposal, sourceIndex int, policy demuxReviewShapePolicy) (int, string, bool) {
	source := revisions[sourceIndex]
	bestIndex := -1
	bestScore := 0
	bestDistance := len(revisions) + 1
	bestReason := ""
	for targetIndex, target := range revisions {
		if targetIndex == sourceIndex || !sameDemuxReviewRoute(source, target) {
			continue
		}
		score, reason := tinyRevisionMergeScore(proposal, source, target)
		if score == 0 {
			continue
		}
		combinedLOC := demuxRevisionMetrics(proposal, source).EffectiveLOC + demuxRevisionMetrics(proposal, target).EffectiveLOC
		if combinedLOC > policy.SoftMaxLOC {
			continue
		}
		score += 2
		distance := absInt(targetIndex - sourceIndex)
		if distance == 1 {
			score++
		}
		if score > bestScore || (score == bestScore && distance < bestDistance) {
			bestIndex = targetIndex
			bestScore = score
			bestDistance = distance
			bestReason = reason
		}
	}
	if bestIndex < 0 || bestScore < 4 {
		return -1, "", false
	}
	return bestIndex, bestReason, true
}

func sameDemuxReviewRoute(left, right RevisionProposal) bool {
	return strings.TrimSpace(left.TargetStack) == strings.TrimSpace(right.TargetStack) &&
		strings.TrimSpace(left.BaseStack) == strings.TrimSpace(right.BaseStack)
}

func tinyRevisionMergeScore(proposal DemuxProposal, source, target RevisionProposal) (int, string) {
	sourceFiles := stringSet(source.Files)
	targetFiles := stringSet(target.Files)
	for file := range sourceFiles {
		if _, ok := targetFiles[file]; ok {
			return 6, "same file " + file
		}
	}
	if sourceSymbol, targetSymbol, ok := sharedRevisionSymbol(proposal, source, target); ok {
		return 6, "same symbol " + firstNonEmpty(sourceSymbol, targetSymbol)
	}
	for file := range sourceFiles {
		counterpart := testCounterpartSource(file)
		if counterpart != "" {
			if _, ok := targetFiles[counterpart]; ok {
				return 5, "source/test counterpart"
			}
		}
	}
	for file := range targetFiles {
		counterpart := testCounterpartSource(file)
		if counterpart != "" {
			if _, ok := sourceFiles[counterpart]; ok {
				return 5, "source/test counterpart"
			}
		}
	}
	if revisionsHaveStructuralEdge(proposal, source, target) {
		return 4, "structural dependency"
	}
	return 0, ""
}

func sharedRevisionSymbol(proposal DemuxProposal, left, right RevisionProposal) (string, string, bool) {
	leftSymbols := revisionSymbols(proposal, left)
	rightSymbols := revisionSymbols(proposal, right)
	for symbol := range leftSymbols {
		if _, ok := rightSymbols[symbol]; ok {
			return symbol, symbol, true
		}
	}
	return "", "", false
}

func revisionSymbols(proposal DemuxProposal, revision RevisionProposal) map[string]struct{} {
	out := map[string]struct{}{}
	for _, hunk := range revisionHunksForMetrics(proposal, revision) {
		symbol := strings.TrimSpace(hunk.Symbol)
		if symbol != "" {
			out[symbol] = struct{}{}
		}
	}
	return out
}

func revisionsHaveStructuralEdge(proposal DemuxProposal, left, right RevisionProposal) bool {
	leftFiles := stringSet(left.Files)
	rightFiles := stringSet(right.Files)
	for _, dep := range proposal.StructuralDeps {
		_, leftFrom := leftFiles[dep.FromFile]
		_, leftTo := leftFiles[dep.ToFile]
		_, rightFrom := rightFiles[dep.FromFile]
		_, rightTo := rightFiles[dep.ToFile]
		if (leftFrom && rightTo) || (rightFrom && leftTo) {
			return true
		}
	}
	return false
}

func mergeDemuxRevision(proposal DemuxProposal, target, source RevisionProposal) RevisionProposal {
	merged := target
	merged.Files = cleanFiles(append(append([]string(nil), target.Files...), source.Files...))
	merged.Hunks = orderedMergedHunks(proposal, append(revisionHunksForMetrics(proposal, target), revisionHunksForMetrics(proposal, source)...))
	merged.HunkIDs = hunkIDs(merged.Hunks)
	if source.UseHunks || target.UseHunks || len(merged.HunkIDs) > 0 {
		merged.UseHunks = true
	}
	merged.ShapeReasons = appendUniqueStrings(merged.ShapeReasons, source.ShapeReasons...)
	return merged
}

func orderedMergedHunks(proposal DemuxProposal, hunks []HunkRange) []HunkRange {
	if len(hunks) == 0 {
		return nil
	}
	byID := map[string]HunkRange{}
	for _, hunk := range hunks {
		if strings.TrimSpace(hunk.ID) == "" {
			continue
		}
		byID[hunk.ID] = hunk
	}
	out := make([]HunkRange, 0, len(byID))
	for _, hunk := range proposal.Hunks {
		if matched, ok := byID[hunk.ID]; ok {
			if matched.Patch == "" {
				matched = hunk
			}
			out = append(out, matched)
			delete(byID, hunk.ID)
		}
	}
	for _, hunk := range byID {
		out = append(out, hunk)
	}
	return out
}

func removeRevisionAt(revisions []RevisionProposal, index int) []RevisionProposal {
	out := make([]RevisionProposal, 0, len(revisions)-1)
	out = append(out, revisions[:index]...)
	out = append(out, revisions[index+1:]...)
	return out
}

func rewriteDemuxRevisionDependencies(revisions []RevisionProposal, fromID, toID string) []RevisionProposal {
	out := append([]RevisionProposal(nil), revisions...)
	for index := range out {
		deps := make([]string, 0, len(out[index].DependsOn))
		seen := map[string]struct{}{}
		for _, dep := range out[index].DependsOn {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			if dep == fromID {
				dep = toID
			}
			if dep == out[index].ID {
				continue
			}
			if _, ok := seen[dep]; ok {
				continue
			}
			seen[dep] = struct{}{}
			deps = append(deps, dep)
		}
		out[index].DependsOn = deps
	}
	return out
}

func appendShapeReason(revision *RevisionProposal, reason string) bool {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return false
	}
	for _, existing := range revision.ShapeReasons {
		if strings.TrimSpace(existing) == reason {
			return false
		}
	}
	revision.ShapeReasons = append(revision.ShapeReasons, reason)
	return true
}

func appendUniqueStrings(values []string, additions ...string) []string {
	out := append([]string(nil), values...)
	for _, addition := range additions {
		addition = strings.TrimSpace(addition)
		if addition == "" {
			continue
		}
		seen := false
		for _, existing := range out {
			if strings.TrimSpace(existing) == addition {
				seen = true
				break
			}
		}
		if !seen {
			out = append(out, addition)
		}
	}
	return out
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
