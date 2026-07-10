package authoring

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	generateHighLogicConfidence   = 0.85
	generateMediumLogicConfidence = 0.50
)

type generatePipeline struct {
	engine *Engine
}

type generateTiming struct {
	Plan      time.Duration
	Refine    time.Duration
	Preflight time.Duration
}

func (e *Engine) generatePipeline() generatePipeline {
	return generatePipeline{engine: e}
}

func (e *Engine) GenerateChanges(ctx context.Context, opts ProposeDemuxOptions) (DemuxPlanPacket, error) {
	if opts.Legacy {
		return e.DemuxChanges(ctx, opts)
	}
	return e.generatePipeline().proposeChanges(ctx, opts)
}

func (p generatePipeline) proposeChanges(ctx context.Context, opts ProposeDemuxOptions) (DemuxPlanPacket, error) {
	started := time.Now()
	var timing generateTiming
	planStarted := time.Now()
	demuxProgress(opts.ProgressWriter, "Grouping changes...")
	proposal, err := p.engine.ProposeDemux(ctx, opts)
	if err != nil {
		return DemuxPlanPacket{}, err
	}
	timing.Plan = time.Since(planStarted)
	proposal = polishGenerateProposal(proposal)
	proposal = annotateGenerateLogicConfidence(proposal)
	proposal = appendGeneratePipelineWarning(proposal, fmt.Sprintf("generate logic confidence %.2f", proposal.Confidence.LogicConfidence))
	if opts.PlanOnly {
		proposal = p.engine.annotateSemanticLabels(ctx, proposal)
		proposal = appendGenerateTimingWarning(proposal, timing, time.Since(started))
	}
	if saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal); saveErr == nil {
		proposal = saved
	}
	if opts.PlanOnly {
		return p.engine.demuxPipeline().packetForProposal(ctx, proposal)
	}

	refineStarted := time.Now()
	proposal, err = p.refineDraft(ctx, proposal, opts)
	timing.Refine = time.Since(refineStarted)
	if err != nil {
		if strings.TrimSpace(proposal.ID) != "" {
			packet, packetErr := p.engine.demuxPipeline().packetForProposal(ctx, proposal)
			if packetErr == nil {
				return packet, err
			}
		}
		return DemuxPlanPacket{}, err
	}
	proposal = polishGenerateProposal(proposal)
	proposal = p.engine.annotateSemanticLabels(ctx, proposal)
	preflightStarted := time.Now()
	proposal, err = p.preflightHighConfidencePlan(ctx, proposal, opts)
	timing.Preflight = time.Since(preflightStarted)
	if err != nil {
		if strings.TrimSpace(proposal.ID) != "" {
			packet, packetErr := p.engine.demuxPipeline().packetForProposal(ctx, proposal)
			if packetErr == nil {
				return packet, err
			}
		}
		return DemuxPlanPacket{}, err
	}
	proposal = appendGenerateTimingWarning(proposal, timing, time.Since(started))
	if saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal); saveErr == nil {
		proposal = saved
	}
	return p.engine.demuxPipeline().packetForProposal(ctx, proposal)
}

func (p generatePipeline) refineDraft(ctx context.Context, proposal DemuxProposal, opts ProposeDemuxOptions) (DemuxProposal, error) {
	demuxProgress(opts.ProgressWriter, "Checking revisions...")
	review, err := p.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return proposal, err
	}
	proposal = reviewedProposalOrFallback(review, proposal)
	proposal = annotateGenerateLogicConfidence(proposal)
	state := demuxWorkflowState(review, proposal)
	if state != DemuxWorkflowReadyToApply ||
		proposal.Confidence.EffectiveConfidence < generateMediumLogicConfidence {
		proposal = appendGenerateConfidenceLog(proposal, "deterministic_simplify", state)
		proposal = conservativeGenerateProposal(proposal, opts.Intent, "simplified generate plan because confidence or review did not clear the default gate")
		saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal)
		if saveErr != nil {
			return proposal, saveErr
		}
		proposal = saved
		review, err = p.engine.ReviewDemuxPlan(ctx, proposal)
		if err != nil {
			return proposal, err
		}
		proposal = reviewedProposalOrFallback(review, proposal)
		return proposal, nil
	}
	if proposal.Confidence.EffectiveConfidence < generateHighLogicConfidence {
		if !generateLLMRepairEnabled() {
			proposal = appendGenerateConfidenceLog(proposal, "skip_llm_repair_medium_confidence", state)
			if saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal); saveErr == nil {
				proposal = saved
			}
			return proposal, nil
		}
		proposal = appendGenerateConfidenceLog(proposal, "attempt_llm_repair", state)
		result, repairErr := p.repairLowConfidenceDraft(ctx, proposal, review, opts)
		if repairErr == nil && strings.TrimSpace(result.Proposal.ID) != "" {
			proposal = annotateGenerateLogicConfidence(result.Proposal)
			review = result.Review
		} else if repairErr != nil {
			proposal = appendGeneratePipelineWarning(proposal, fmt.Sprintf("generate llm repair failed: %v", repairErr))
		}
	}
	if demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply ||
		proposal.Confidence.EffectiveConfidence < generateMediumLogicConfidence {
		proposal = appendGenerateConfidenceLog(proposal, "deterministic_simplify_after_repair", demuxWorkflowState(review, proposal))
		proposal = conservativeGenerateProposal(proposal, opts.Intent, "simplified generate plan because confidence or review did not clear the default gate")
		saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal)
		if saveErr != nil {
			return proposal, saveErr
		}
		proposal = saved
		review, err = p.engine.ReviewDemuxPlan(ctx, proposal)
		if err != nil {
			return proposal, err
		}
		proposal = reviewedProposalOrFallback(review, proposal)
	}
	return proposal, nil
}

func (p generatePipeline) repairLowConfidenceDraft(ctx context.Context, proposal DemuxProposal, review ReviewDemuxResult, opts ProposeDemuxOptions) (DemuxAIReviewResult, error) {
	if demuxWorkflowState(review, proposal) == DemuxWorkflowReadyToApply {
		review.Valid = false
		review.Errors = append(review.Errors, fmt.Sprintf("logic confidence %.2f is below %.2f", proposal.Confidence.EffectiveConfidence, generateHighLogicConfidence))
	}
	review.Proposal = proposal
	review.RepairHints = append(review.RepairHints, RepairHint{
		Kind:       "low_logic_confidence",
		Suggestion: "revise grouping, stack routes, hunk coverage, or merge risky revisions so the plan clears the logic confidence gate",
	})
	return p.engine.demuxRepair().repairReviewed(ctx, proposal, review, DemuxAIReviewOptions{
		Model:       opts.Model,
		MaxWarnings: opts.MaxWarnings,
	})
}

func (p generatePipeline) preflightHighConfidencePlan(ctx context.Context, proposal DemuxProposal, opts ProposeDemuxOptions) (DemuxProposal, error) {
	coverageInput := proposal
	review, err := p.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return proposal, err
	}
	proposal = annotateGenerateLogicConfidence(reviewedProposalOrFallback(review, proposal))
	if demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply {
		return proposal, nil
	}
	if checked, downgraded := coverageCheckedGenerateProposal(coverageInput, proposal, opts.Intent); downgraded {
		if saved, saveErr := p.engine.SaveDemuxProposal(ctx, checked); saveErr == nil {
			checked = saved
		}
		return checked, nil
	}
	if !generateApplyPreflightEnabled() {
		proposal = appendGeneratePipelineWarning(proposal, "skipped disposable apply preflight for fast generate path; set GX_GENERATE_VERIFY=1 to force it")
		if saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal); saveErr == nil {
			proposal = saved
		}
		return proposal, nil
	}
	if proposal.Confidence.EffectiveConfidence < generateHighLogicConfidence {
		proposal = appendGeneratePipelineWarning(proposal, fmt.Sprintf("skipped disposable apply preflight because logic confidence %.2f is below %.2f", proposal.Confidence.EffectiveConfidence, generateHighLogicConfidence))
		if saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal); saveErr == nil {
			proposal = saved
		}
		return proposal, nil
	}

	demuxProgress(opts.ProgressWriter, "Verifying revision apply in a disposable check...")
	if err := p.engine.PreflightDemuxApply(ctx, proposal); err == nil {
		return proposal, nil
	} else {
		demuxProgress(opts.ProgressWriter, "Repairing revision groups after verification failed...")
		result, repairErr := p.engine.demuxRepair().repairApplyPreflightFailure(ctx, proposal, err, DemuxAIReviewOptions{
			Model:       opts.Model,
			MaxWarnings: opts.MaxWarnings,
		})
		if repairErr == nil && strings.TrimSpace(result.Proposal.ID) != "" {
			repaired := annotateGenerateLogicConfidence(result.Proposal)
			if repaired.Confidence.EffectiveConfidence >= generateHighLogicConfidence {
				demuxProgress(opts.ProgressWriter, "Verifying repaired revision groups...")
				preflightErr := p.engine.PreflightDemuxApply(ctx, repaired)
				if preflightErr == nil {
					return repaired, nil
				}
				proposal = appendDemuxApplyPreflightWarning(repaired, preflightErr)
			} else {
				proposal = repaired
			}
		} else {
			proposal = appendDemuxApplyPreflightWarning(proposal, err)
		}
	}

	fallback := conservativeGenerateProposal(proposal, opts.Intent, "simplified generate plan after apply preflight failure")
	saved, saveErr := p.engine.SaveDemuxProposal(ctx, fallback)
	if saveErr != nil {
		return proposal, saveErr
	}
	fallback = saved
	demuxProgress(opts.ProgressWriter, "Verifying simplified revision groups...")
	if err := p.engine.PreflightDemuxApply(ctx, fallback); err != nil {
		fallback = appendDemuxApplyPreflightWarning(fallback, err)
		if saved, saveErr := p.engine.SaveDemuxProposal(ctx, fallback); saveErr == nil {
			fallback = saved
		}
	}
	return fallback, nil
}

func coverageCheckedGenerateProposal(coverageInput, readyProposal DemuxProposal, intent string) (DemuxProposal, bool) {
	if err := validateDemuxPatchCoverage(coverageInput); err != nil {
		return conservativeCoverageMismatchProposal(readyProposal, intent, err), true
	}
	if err := validateDemuxPatchCoverage(readyProposal); err != nil {
		return conservativeCoverageMismatchProposal(readyProposal, intent, err), true
	}
	return readyProposal, false
}

func conservativeCoverageMismatchProposal(proposal DemuxProposal, intent string, err error) DemuxProposal {
	fallback := conservativeGenerateProposal(proposal, intent, "simplified generate plan after patch coverage mismatch")
	fallback = appendGeneratePipelineWarning(fallback, fmt.Sprintf("patch coverage mismatch: %v; using conservative revision grouping", err))
	if coverageErr := validateDemuxPatchCoverage(fallback); coverageErr == nil {
		return fallback
	}
	fallback.Revisions = conservativeGenerateWholePlanRevisions(proposal, intent, "simplified generate plan after patch coverage mismatch")
	fallback.FeasibilityWarnings = feasibilityWarningsForProposal(fallback)
	fallback = annotateGenerateLogicConfidence(fallback)
	return fallback
}

func generateLLMRepairEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_GENERATE_LLM_REPAIR"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func generateRefineUsesLLMRepair(review ReviewDemuxResult, proposal DemuxProposal) bool {
	needsRepair := demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply ||
		proposal.Confidence.EffectiveConfidence < generateHighLogicConfidence
	return needsRepair && generateLLMRepairEnabled()
}

func generateApplyPreflightEnabled() bool {
	if strings.TrimSpace(os.Getenv("GX_COMPOSE_SKIP_APPLY_PREFLIGHT")) == "1" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_GENERATE_VERIFY"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func polishGenerateProposal(proposal DemuxProposal) DemuxProposal {
	proposal = coalesceDocumentationRevisions(proposal)
	for index := range proposal.Revisions {
		proposal.Revisions[index].Intent = generateRevisionIntent(proposal.Revisions[index])
	}
	proposal = groupGenerateRevisionsByRoute(proposal)
	return proposal
}

func groupGenerateRevisionsByRoute(proposal DemuxProposal) DemuxProposal {
	if len(proposal.Revisions) < 2 {
		return proposal
	}
	firstIndexByTarget := map[string]int{}
	for index, revision := range proposal.Revisions {
		if len(revision.DependsOn) > 0 || !hasNonCurrentDemuxRoute(revision) {
			return proposal
		}
		target := routeTargetStack(revision)
		if target == "" {
			return proposal
		}
		if _, ok := firstIndexByTarget[target]; !ok {
			firstIndexByTarget[target] = index
		}
	}
	ordered := append([]RevisionProposal(nil), proposal.Revisions...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left := routeTargetStack(ordered[i])
		right := routeTargetStack(ordered[j])
		if left == right {
			return false
		}
		return firstIndexByTarget[left] < firstIndexByTarget[right]
	})
	proposal.Revisions = ordered
	return proposal
}

func coalesceDocumentationRevisions(proposal DemuxProposal) DemuxProposal {
	targetByStack := map[string]int{}
	out := make([]RevisionProposal, 0, len(proposal.Revisions))
	for _, revision := range proposal.Revisions {
		files := cleanFiles(revision.Files)
		if len(files) == 0 || !containsOnlyDocumentationFiles(files) {
			out = append(out, revision)
			continue
		}
		stack := strings.TrimSpace(revision.TargetStack)
		key := stack
		if key == "" {
			key = "_docs"
		}
		if existingIndex, ok := targetByStack[key]; ok {
			merged := out[existingIndex]
			merged.Files = cleanFiles(append(merged.Files, files...))
			merged.Hunks = orderedMergedHunks(proposal, append(revisionHunksForMetrics(proposal, merged), revisionHunksForMetrics(proposal, revision)...))
			merged.HunkIDs = hunkIDs(merged.Hunks)
			merged.UseHunks = merged.UseHunks || revision.UseHunks
			merged.Confidence = maxFloat(merged.Confidence, revision.Confidence)
			if merged.RouteConfidence == 0 || (revision.RouteConfidence > 0 && revision.RouteConfidence < merged.RouteConfidence) {
				merged.RouteConfidence = revision.RouteConfidence
			}
			merged.ShapeReasons = appendUniqueString(merged.ShapeReasons, "coalesced documentation files routed to the same stack")
			merged.SessionIDs = uniqueStrings(append(merged.SessionIDs, revision.SessionIDs...))
			out[existingIndex] = merged
			continue
		}
		targetByStack[key] = len(out)
		revision.Files = files
		out = append(out, revision)
	}
	proposal.Revisions = out
	return proposal
}

func generateRevisionIntent(revision RevisionProposal) string {
	files := cleanFiles(revision.Files)
	if containsOnlyDocumentationFiles(files) {
		return "update documentation"
	}
	return revision.Intent
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
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

func reviewedProposalOrFallback(review ReviewDemuxResult, fallback DemuxProposal) DemuxProposal {
	if strings.TrimSpace(review.Proposal.ID) != "" {
		return review.Proposal
	}
	return fallback
}

func annotateGenerateLogicConfidence(proposal DemuxProposal) DemuxProposal {
	confidence := generateLogicConfidence(proposal)
	proposal.Confidence.LogicConfidence = confidence.LogicConfidence
	if proposal.Confidence.LLMConfidence <= 0 {
		proposal.Confidence.EffectiveConfidence = confidence.LogicConfidence
	} else {
		proposal.Confidence.EffectiveConfidence = math.Min(confidence.LogicConfidence, proposal.Confidence.LLMConfidence)
	}
	proposal.Confidence.LogicReasons = confidence.LogicReasons
	if proposal.Confidence.EffectiveConfidence == 0 {
		proposal.Confidence.EffectiveConfidence = proposal.Confidence.LogicConfidence
	}
	return proposal
}

func generateLogicConfidence(proposal DemuxProposal) PlanConfidence {
	score := 0.95
	var reasons []ConfidenceReason
	routeRisks := 0
	oversizedRevisions := 0
	targetSizedRevisions := 0
	sharedFileRisks := 0
	hunkSplitCount := 0
	add := func(kind, severity, message, suggestion string, delta float64) {
		reasons = append(reasons, ConfidenceReason{Kind: kind, Severity: severity, Message: message, Suggestion: suggestion, Delta: delta})
		score += delta
	}
	if len(proposal.Revisions) == 0 {
		add("empty_plan", "risk", "proposal contains no revisions", "create at least one revision that covers the changed files", -0.95)
	}
	fileOwners := map[string]string{}
	for _, revision := range proposal.Revisions {
		if len(revision.Files) == 0 {
			add("empty_revision", "risk", fmt.Sprintf("revision %s has no files", revision.ID), "drop empty revisions or assign covered files/hunks to them", -0.30)
		}
		effectiveLOC := revision.EffectiveLOC
		if effectiveLOC <= 0 {
			effectiveLOC = demuxRevisionMetrics(proposal, revision).EffectiveLOC
		}
		switch {
		case effectiveLOC > defaultDemuxReviewShapePolicy.SoftMaxLOC:
			oversizedRevisions++
		case effectiveLOC >= defaultDemuxReviewShapePolicy.TargetMinLOC && effectiveLOC <= defaultDemuxReviewShapePolicy.TargetMaxLOC:
			targetSizedRevisions++
		}
		if strings.TrimSpace(revision.TargetStack) != "" && revision.RouteConfidence > 0 && revision.RouteConfidence < generateMediumLogicConfidence {
			routeRisks++
		}
		if revision.UseHunks {
			hunkSplitCount++
		}
		for _, file := range revision.Files {
			file = strings.TrimSpace(file)
			if file == "" {
				continue
			}
			if owner, ok := fileOwners[file]; ok && owner != revision.ID {
				sharedFileRisks++
				continue
			}
			fileOwners[file] = revision.ID
		}
	}
	if hunkSplitCount > 0 {
		add("hunk_split", "caution", fmt.Sprintf("%d revision(s) use hunk-level splitting", hunkSplitCount), "prefer whole-file revisions unless the split is necessary for reviewability", -0.05)
	}
	if routeRisks > 0 {
		add("low_confidence_route", "caution", fmt.Sprintf("%d revision route(s) have low heuristic confidence", routeRisks), "choose a stronger target_stack from file ownership, labels, or prior stack history", -0.08)
	}
	if sharedFileRisks > 0 {
		add("shared_file", "risk", fmt.Sprintf("%d file assignment(s) are split across multiple revisions", sharedFileRisks), "merge shared-file edits into one revision or use explicit non-overlapping hunks", -0.18)
	}
	for _, warning := range proposal.FeasibilityWarnings {
		if strings.EqualFold(strings.TrimSpace(warning.Severity), "info") {
			continue
		}
		switch warning.Source {
		case "structural_dependency", "inferred_dependency":
			add(warning.Source, "risk", warning.Message, "reorder dependent revisions, add depends_on, or merge tightly coupled edits", -0.14)
		case "apply_preflight":
			add("apply_preflight", "risk", warning.Message, "repair the apply failure before trusting this stack plan", -0.35)
		default:
			add(warning.Source, "caution", warning.Message, "tighten file grouping, route labels, or coverage metadata for this warning", -0.08)
		}
	}
	if len(reasons) == 0 {
		reasons = append(reasons, ConfidenceReason{
			Kind:     "reviewable_plan",
			Severity: "positive",
			Message:  "deterministic plan has no blocking logic risks",
		})
	}
	if score < 0 {
		score = 0
	}
	if score > 0.99 {
		score = 0.99
	}
	return PlanConfidence{
		LogicConfidence:     roundConfidence(score),
		EffectiveConfidence: roundConfidence(score),
		LogicReasons:        reasons,
	}
}

func appendGenerateConfidenceLog(proposal DemuxProposal, action string, state DemuxWorkflowState) DemuxProposal {
	proposal = annotateGenerateLogicConfidence(proposal)
	return appendGeneratePipelineWarning(proposal, fmt.Sprintf(
		"generate confidence logic=%.2f effective=%.2f state=%s action=%s suggestions=%q",
		proposal.Confidence.LogicConfidence,
		proposal.Confidence.EffectiveConfidence,
		state,
		action,
		strings.Join(generateConfidenceSuggestions(proposal), "; "),
	))
}

func generateConfidenceSuggestions(proposal DemuxProposal) []string {
	seen := map[string]struct{}{}
	var suggestions []string
	for _, reason := range proposal.Confidence.LogicReasons {
		if strings.EqualFold(strings.TrimSpace(reason.Severity), "positive") {
			continue
		}
		suggestion := strings.TrimSpace(reason.Suggestion)
		if suggestion == "" {
			suggestion = defaultGenerateConfidenceSuggestion(reason)
		}
		if suggestion == "" {
			continue
		}
		if _, ok := seen[suggestion]; ok {
			continue
		}
		seen[suggestion] = struct{}{}
		suggestions = append(suggestions, suggestion)
		if len(suggestions) >= 5 {
			break
		}
	}
	if len(suggestions) == 0 && proposal.Confidence.EffectiveConfidence < generateHighLogicConfidence {
		suggestions = append(suggestions, "collect more deterministic route evidence before using LLM repair")
	}
	return suggestions
}

func defaultGenerateConfidenceSuggestion(reason ConfidenceReason) string {
	switch strings.TrimSpace(reason.Kind) {
	case "hunk_split":
		return "prefer whole-file revisions unless the split is necessary for reviewability"
	case "low_confidence_route":
		return "choose a stronger target_stack from file ownership, labels, or prior stack history"
	case "shared_file":
		return "merge shared-file edits into one revision or use explicit non-overlapping hunks"
	case "structural_dependency", "inferred_dependency":
		return "reorder dependent revisions, add depends_on, or merge tightly coupled edits"
	default:
		return ""
	}
}

func conservativeGenerateProposal(proposal DemuxProposal, intent, reason string) DemuxProposal {
	proposal.Revisions = conservativeGenerateRevisions(proposal, intent, reason)
	proposal.FeasibilityWarnings = feasibilityWarningsForProposal(proposal)
	proposal = appendGeneratePipelineWarning(proposal, reason)
	proposal = annotateGenerateLogicConfidence(proposal)
	return proposal
}

func conservativeGenerateRevisions(proposal DemuxProposal, intent, reason string) []RevisionProposal {
	if normalized, err := normalizeDemuxProposal(proposal); err == nil {
		return markConservativeGenerateRevisions(normalized.Revisions, reason)
	} else if revisions, ok := repairConservativeGenerateRevisions(proposal, err, reason); ok {
		return revisions
	}
	return conservativeGenerateWholePlanRevisions(proposal, intent, reason)
}

func markConservativeGenerateRevisions(revisions []RevisionProposal, reason string) []RevisionProposal {
	out := make([]RevisionProposal, 0, len(revisions))
	for _, revision := range revisions {
		revision.Confidence = maxFloat(revision.Confidence, 0.90)
		revision.ShapeReasons = appendUniqueString(revision.ShapeReasons, reason)
		out = append(out, revision)
	}
	return out
}

func repairConservativeGenerateRevisions(proposal DemuxProposal, validationErr error, reason string) ([]RevisionProposal, bool) {
	repaired := proposal
	switch e := demuxValidationRepair(validationErr).(type) {
	case duplicateOwnerRepair:
		revisions, ok := mergeConservativeImplicatedRevisions(proposal.Revisions, e.revisionIDs, reason)
		if !ok {
			return nil, false
		}
		repaired.Revisions = revisions
	case removeHunkReferenceRepair:
		revisions, ok := removeConservativeHunkReference(proposal.Revisions, e.revisionID, e.hunkID, reason)
		if !ok {
			return nil, false
		}
		repaired.Revisions = revisions
	case removeDependencyRepair:
		revisions, ok := removeConservativeDependency(proposal.Revisions, e.revisionID, e.dependsOn, reason)
		if !ok {
			return nil, false
		}
		repaired.Revisions = revisions
	default:
		return nil, false
	}
	normalized, err := normalizeDemuxProposal(repaired)
	if err != nil {
		return nil, false
	}
	return markConservativeGenerateRevisions(normalized.Revisions, reason), true
}

type duplicateOwnerRepair struct {
	revisionIDs []string
}

type removeHunkReferenceRepair struct {
	revisionID string
	hunkID     string
}

type removeDependencyRepair struct {
	revisionID string
	dependsOn  string
}

func demuxValidationRepair(err error) any {
	var duplicateHunk DuplicateHunkAssignmentError
	if errors.As(err, &duplicateHunk) {
		return duplicateOwnerRepair{revisionIDs: []string{duplicateHunk.FirstOwnerID, duplicateHunk.SecondOwnerID}}
	}
	var duplicateWhole DuplicateWholeFileOwnerError
	if errors.As(err, &duplicateWhole) {
		return duplicateOwnerRepair{revisionIDs: []string{duplicateWhole.FirstOwnerID, duplicateWhole.SecondOwnerID}}
	}
	var mixed MixedHunkWholeFileCoverageError
	if errors.As(err, &mixed) {
		return duplicateOwnerRepair{revisionIDs: []string{mixed.HunkRevisionID, mixed.WholeRevisionID}}
	}
	var unknownHunk UnknownHunkReferenceError
	if errors.As(err, &unknownHunk) && strings.TrimSpace(unknownHunk.RevisionID) != "" && strings.TrimSpace(unknownHunk.HunkID) != "" {
		return removeHunkReferenceRepair{revisionID: unknownHunk.RevisionID, hunkID: unknownHunk.HunkID}
	}
	var dangling DanglingDependsOnError
	if errors.As(err, &dangling) && strings.TrimSpace(dangling.RevisionID) != "" && strings.TrimSpace(dangling.DependsOn) != "" {
		return removeDependencyRepair{revisionID: dangling.RevisionID, dependsOn: dangling.DependsOn}
	}
	return nil
}

func mergeConservativeImplicatedRevisions(revisions []RevisionProposal, revisionIDs []string, reason string) ([]RevisionProposal, bool) {
	ids := stringSet(revisionIDs)
	if len(ids) == 0 {
		return nil, false
	}
	var merged RevisionProposal
	var files []string
	mergedAny := false
	out := make([]RevisionProposal, 0, len(revisions))
	for _, revision := range revisions {
		if _, ok := ids[revision.ID]; !ok {
			out = append(out, revision)
			continue
		}
		if !mergedAny {
			merged = revision
			mergedAny = true
		}
		files = append(files, revision.Files...)
		files = append(files, proposalFilesFromHunks(revision.Hunks)...)
	}
	if !mergedAny {
		return nil, false
	}
	files = cleanFiles(files)
	if len(files) == 0 {
		return nil, false
	}
	merged.Files = files
	merged.UseHunks = false
	merged.HunkIDs = nil
	merged.Hunks = nil
	merged.DependsOn = nil
	merged.Confidence = maxFloat(merged.Confidence, 0.90)
	merged.ShapeReasons = appendUniqueString(merged.ShapeReasons, reason)
	result := make([]RevisionProposal, 0, len(out)+1)
	inserted := false
	for _, revision := range revisions {
		if _, ok := ids[revision.ID]; ok {
			if !inserted {
				result = append(result, merged)
				inserted = true
			}
			continue
		}
		result = append(result, revision)
	}
	return result, true
}

func removeConservativeHunkReference(revisions []RevisionProposal, revisionID, hunkID, reason string) ([]RevisionProposal, bool) {
	out := make([]RevisionProposal, 0, len(revisions))
	changed := false
	for _, revision := range revisions {
		if revision.ID != revisionID {
			out = append(out, revision)
			continue
		}
		revision.HunkIDs = removeString(revision.HunkIDs, hunkID)
		var hunks []HunkRange
		for _, hunk := range revision.Hunks {
			if hunk.ID != hunkID {
				hunks = append(hunks, hunk)
			}
		}
		revision.Hunks = hunks
		if revision.UseHunks && len(revision.HunkIDs) == 0 && len(revision.Hunks) == 0 {
			return nil, false
		}
		revision.ShapeReasons = appendUniqueString(revision.ShapeReasons, reason)
		out = append(out, revision)
		changed = true
	}
	return out, changed
}

func removeConservativeDependency(revisions []RevisionProposal, revisionID, dependsOn, reason string) ([]RevisionProposal, bool) {
	out := make([]RevisionProposal, 0, len(revisions))
	changed := false
	for _, revision := range revisions {
		if revision.ID == revisionID {
			next := removeString(revision.DependsOn, dependsOn)
			if len(next) != len(revision.DependsOn) {
				revision.DependsOn = next
				revision.ShapeReasons = appendUniqueString(revision.ShapeReasons, reason)
				changed = true
			}
		}
		out = append(out, revision)
	}
	return out, changed
}

func removeString(values []string, remove string) []string {
	remove = strings.TrimSpace(remove)
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == remove {
			continue
		}
		out = append(out, value)
	}
	return out
}

func conservativeGenerateWholePlanRevisions(proposal DemuxProposal, intent, reason string) []RevisionProposal {
	if len(proposal.Revisions) == 0 {
		files := demuxProposalCoveredFiles(proposal)
		if len(files) == 0 {
			files = proposalFilesFromHunks(proposal.Hunks)
		}
		sort.Strings(files)
		sessionIDs, provenanceStatus := firstRevisionProvenance(proposal.Revisions)
		if provenanceStatus == "" {
			provenanceStatus = "absent"
		}
		return []RevisionProposal{{
			ID:               "u1",
			Intent:           proposalIntent(intent, files),
			Files:            files,
			ProvenanceStatus: provenanceStatus,
			SessionIDs:       sessionIDs,
			Confidence:       0.95,
			ShapeReasons:     []string{reason},
		}}
	}
	seenFiles := map[string]string{}
	out := make([]RevisionProposal, 0, len(proposal.Revisions))
	for _, revision := range proposal.Revisions {
		files := cleanFiles(revision.Files)
		if len(files) == 0 {
			files = proposalFilesFromHunks(revision.Hunks)
		}
		var kept []string
		for _, file := range files {
			if owner := seenFiles[file]; owner != "" && owner != revision.ID {
				continue
			}
			seenFiles[file] = revision.ID
			kept = append(kept, file)
		}
		if len(kept) == 0 {
			continue
		}
		revision.Files = kept
		revision.UseHunks = false
		revision.HunkIDs = nil
		revision.Hunks = nil
		revision.DependsOn = nil
		revision.Confidence = maxFloat(revision.Confidence, 0.90)
		revision.ShapeReasons = appendUniqueString(revision.ShapeReasons, reason)
		out = append(out, revision)
	}
	if len(out) == 0 {
		return conservativeGenerateWholePlanRevisions(DemuxProposal{Hunks: proposal.Hunks}, intent, reason)
	}
	for index := range out {
		if strings.TrimSpace(out[index].ID) == "" {
			out[index].ID = fmt.Sprintf("u%d", index+1)
		}
	}
	return out
}

func proposalFilesFromHunks(hunks []HunkRange) []string {
	seen := map[string]struct{}{}
	var files []string
	for _, hunk := range hunks {
		file := strings.TrimSpace(hunk.File)
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		files = append(files, file)
	}
	return files
}

func firstRevisionProvenance(revisions []RevisionProposal) ([]string, string) {
	for _, revision := range revisions {
		if strings.TrimSpace(revision.ProvenanceStatus) != "" || len(revision.SessionIDs) > 0 {
			return append([]string(nil), revision.SessionIDs...), revision.ProvenanceStatus
		}
	}
	return nil, ""
}

func appendGeneratePipelineWarning(proposal DemuxProposal, message string) DemuxProposal {
	message = strings.TrimSpace(message)
	if message == "" {
		return proposal
	}
	for _, warning := range proposal.Warnings {
		if strings.TrimSpace(warning) == message {
			return proposal
		}
	}
	proposal.Warnings = append(proposal.Warnings, message)
	return proposal
}

func appendGenerateTimingWarning(proposal DemuxProposal, timing generateTiming, total time.Duration) DemuxProposal {
	return appendGeneratePipelineWarning(proposal, fmt.Sprintf(
		"generate timing plan_ms=%d refine_ms=%d preflight_ms=%d total_ms=%d",
		timing.Plan.Milliseconds(),
		timing.Refine.Milliseconds(),
		timing.Preflight.Milliseconds(),
		total.Milliseconds(),
	))
}

func roundConfidence(value float64) float64 {
	return math.Round(value*100) / 100
}

func maxFloat(left, right float64) float64 {
	if left > right {
		return left
	}
	return right
}

func appendUniqueString(values []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return values
	}
	for _, current := range values {
		if strings.TrimSpace(current) == value {
			return values
		}
	}
	return append(values, value)
}
