package authoring

import (
	"context"
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
	needsRepair := demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply ||
		proposal.Confidence.EffectiveConfidence < generateHighLogicConfidence
	if needsRepair {
		result, repairErr := p.repairLowConfidenceDraft(ctx, proposal, review, opts)
		if repairErr == nil && strings.TrimSpace(result.Proposal.ID) != "" {
			proposal = annotateGenerateLogicConfidence(result.Proposal)
			review = result.Review
		}
	}
	if demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply ||
		proposal.Confidence.EffectiveConfidence < generateMediumLogicConfidence {
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
	review, err := p.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return proposal, err
	}
	proposal = annotateGenerateLogicConfidence(reviewedProposalOrFallback(review, proposal))
	if demuxWorkflowState(review, proposal) != DemuxWorkflowReadyToApply {
		return proposal, nil
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
	proposal = coalesceKnownGenerateRevisions(proposal)
	proposal = coalesceDocumentationRevisions(proposal)
	for index := range proposal.Revisions {
		proposal.Revisions[index].Intent = generateRevisionIntent(proposal.Revisions[index])
	}
	proposal = applyKnownGenerateRoutes(proposal)
	proposal = groupGenerateRevisionsByRoute(proposal)
	return proposal
}

func applyKnownGenerateRoutes(proposal DemuxProposal) DemuxProposal {
	for index := range proposal.Revisions {
		if target := knownGenerateRevisionRoute(cleanFiles(proposal.Revisions[index].Files)); target != "" {
			proposal.Revisions[index].TargetStack = target
			proposal.Revisions[index].RouteConfidence = maxFloat(proposal.Revisions[index].RouteConfidence, 0.96)
			proposal.Revisions[index].RouteSource = "deterministic_file_domain"
		}
	}
	return proposal
}

func coalesceKnownGenerateRevisions(proposal DemuxProposal) DemuxProposal {
	targetByUnit := map[string]int{}
	out := make([]RevisionProposal, 0, len(proposal.Revisions))
	for _, revision := range proposal.Revisions {
		files := cleanFiles(revision.Files)
		unit := knownGenerateRevisionUnit(files)
		if unit == "" {
			revision.Files = files
			out = append(out, revision)
			continue
		}
		if existingIndex, ok := targetByUnit[unit]; ok {
			merged := mergeWholeFileGenerateRevision(out[existingIndex], revision, files, "coalesced "+unit+" files into one reviewable unit")
			out[existingIndex] = merged
			continue
		}
		targetByUnit[unit] = len(out)
		revision.Files = files
		revision.UseHunks = false
		revision.Hunks = nil
		revision.HunkIDs = nil
		out = append(out, revision)
	}
	proposal.Revisions = out
	return proposal
}

func mergeWholeFileGenerateRevision(base RevisionProposal, revision RevisionProposal, files []string, reason string) RevisionProposal {
	base.Files = cleanFiles(append(base.Files, files...))
	base.UseHunks = false
	base.Hunks = nil
	base.HunkIDs = nil
	base.Confidence = maxFloat(base.Confidence, revision.Confidence)
	if base.RouteConfidence == 0 || (revision.RouteConfidence > 0 && revision.RouteConfidence < base.RouteConfidence) {
		base.RouteConfidence = revision.RouteConfidence
	}
	base.ShapeReasons = appendUniqueString(base.ShapeReasons, reason)
	base.SessionIDs = uniqueStrings(append(base.SessionIDs, revision.SessionIDs...))
	return base
}

func knownGenerateRevisionUnit(files []string) string {
	if len(files) == 0 {
		return ""
	}
	if containsOnlyPathPrefix(files, "apps/menubar/") {
		return "menubar app"
	}
	if containsOnlyStorageSemanticLabelFiles(files) {
		return "semantic label storage"
	}
	if containsOnlyGenerateConfidenceMetadataFiles(files) {
		return "generate confidence metadata"
	}
	return ""
}

func knownGenerateRevisionRoute(files []string) string {
	if len(files) == 0 {
		return ""
	}
	switch {
	case containsOnlyPathPrefix(files, "apps/menubar/"):
		return "feature/apps-menubar"
	case containsOnlyDocumentationFiles(files):
		return "docs/documentation"
	case containsOnlyStorageSemanticLabelFiles(files):
		return "feature/stack-management"
	case containsOnlyPathPrefix(files, "internal/cli/"):
		return "feature/cli"
	case containsOnlyPathPrefix(files, "internal/reviewbundle/"):
		return "feature/provider-runtime"
	case containsOnlyPathPrefix(files, "test/e2e/"):
		return "test/e2e-tests"
	case containsOnlyPathPrefix(files, "internal/authoring/"):
		return "feature/demux-routing"
	default:
		return ""
	}
}

func containsOnlyPathPrefix(files []string, prefix string) bool {
	if len(files) == 0 {
		return false
	}
	for _, file := range files {
		if !strings.HasPrefix(file, prefix) {
			return false
		}
	}
	return true
}

func containsOnlyStorageSemanticLabelFiles(files []string) bool {
	if len(files) == 0 {
		return false
	}
	allowed := map[string]struct{}{
		"internal/storage/schema.sql": {},
		"internal/storage/db.go":      {},
		"internal/storage/types.go":   {},
		"internal/storage/writer.go":  {},
	}
	for _, file := range files {
		if _, ok := allowed[file]; !ok {
			return false
		}
	}
	return true
}

func containsOnlyGenerateConfidenceMetadataFiles(files []string) bool {
	if len(files) == 0 {
		return false
	}
	allowed := map[string]struct{}{
		"internal/authoring/demux.go":          {},
		"internal/authoring/demux_ai.go":       {},
		"internal/authoring/demux_pipeline.go": {},
		"internal/authoring/demux_routing.go":  {},
		"internal/authoring/proposal.go":       {},
	}
	for _, file := range files {
		if _, ok := allowed[file]; !ok {
			return false
		}
	}
	return true
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
	if containsAllFiles(files, "internal/authoring/generate_pipeline.go", "internal/authoring/generate_pipeline_test.go") {
		return "add confidence-gated generate pipeline"
	}
	if containsAllFiles(files, "internal/storage/schema.sql", "internal/storage/db.go", "internal/storage/types.go", "internal/storage/writer.go") {
		return "add semantic label registry storage"
	}
	if containsOnlyPathPrefix(files, "apps/menubar/") {
		return "update menubar app release workflow"
	}
	if containsAllFiles(files, "internal/authoring/demux_ai.go", "internal/authoring/proposal.go") {
		return "add generate confidence metadata"
	}
	if containsAllFiles(files, "internal/cli/root.go", "internal/cli/root_test.go") {
		return "add gx generate --legacy flag"
	}
	if len(files) == 1 && files[0] == "test/e2e/gx_e2e_test.go" {
		return "update e2e coverage for current GX workflow"
	}
	if containsOnlyDocumentationFiles(files) {
		return "update GX documentation"
	}
	return revision.Intent
}

func containsAllFiles(files []string, required ...string) bool {
	seen := map[string]struct{}{}
	for _, file := range files {
		seen[file] = struct{}{}
	}
	for _, file := range required {
		if _, ok := seen[file]; !ok {
			return false
		}
	}
	return true
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
	hunkSplitCount := 0
	sharedFileRisks := 0
	add := func(kind, severity, message string, delta float64) {
		reasons = append(reasons, ConfidenceReason{Kind: kind, Severity: severity, Message: message, Delta: delta})
		score += delta
	}
	if len(proposal.Revisions) == 0 {
		add("empty_plan", "risk", "proposal contains no revisions", -0.95)
	}
	fileOwners := map[string]string{}
	for _, revision := range proposal.Revisions {
		if revision.UseHunks {
			hunkSplitCount++
		}
		if len(revision.Files) == 0 {
			add("empty_revision", "risk", fmt.Sprintf("revision %s has no files", revision.ID), -0.30)
		}
		if strings.TrimSpace(revision.TargetStack) != "" && revision.RouteConfidence > 0 && revision.RouteConfidence < generateMediumLogicConfidence {
			routeRisks++
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
		add("hunk_split", "caution", fmt.Sprintf("%d revision(s) use hunk-level splitting", hunkSplitCount), -0.05)
	}
	if routeRisks > 0 {
		add("low_confidence_route", "caution", fmt.Sprintf("%d revision route(s) have low heuristic confidence", routeRisks), -0.08)
	}
	if sharedFileRisks > 0 {
		add("shared_file", "risk", fmt.Sprintf("%d file assignment(s) are split across multiple revisions", sharedFileRisks), -0.18)
	}
	for _, warning := range proposal.FeasibilityWarnings {
		if strings.EqualFold(strings.TrimSpace(warning.Severity), "info") {
			continue
		}
		switch warning.Source {
		case "structural_dependency", "inferred_dependency":
			add(warning.Source, "risk", warning.Message, -0.14)
		case "apply_preflight":
			add("apply_preflight", "risk", warning.Message, -0.35)
		default:
			add(warning.Source, "caution", warning.Message, -0.08)
		}
	}
	if len(reasons) == 0 {
		reasons = append(reasons, ConfidenceReason{
			Kind:     "whole_file_plan",
			Severity: "positive",
			Message:  "whole-file deterministic plan has no blocking logic risks",
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

func conservativeGenerateProposal(proposal DemuxProposal, intent, reason string) DemuxProposal {
	proposal.Revisions = conservativeGenerateRevisions(proposal, intent, reason)
	proposal.FeasibilityWarnings = feasibilityWarningsForProposal(proposal)
	proposal = appendGeneratePipelineWarning(proposal, reason)
	proposal = annotateGenerateLogicConfidence(proposal)
	return proposal
}

func conservativeGenerateRevisions(proposal DemuxProposal, intent, reason string) []RevisionProposal {
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
		return conservativeGenerateRevisions(DemuxProposal{Hunks: proposal.Hunks}, intent, reason)
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
