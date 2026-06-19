package authoring

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// demuxPipeline owns the ordered demux passes. CLI and MCP adapters should not
// need to know whether a packet came from deterministic planning, deterministic
// repair, or AI repair.
type demuxPipeline struct {
	engine *Engine
}

const defaultDemuxApplyPreflightAttempts = 3

func (e *Engine) demuxPipeline() demuxPipeline {
	return demuxPipeline{engine: e}
}

func (p demuxPipeline) proposeChanges(ctx context.Context, opts ProposeDemuxOptions) (DemuxPlanPacket, error) {
	demuxProgress(opts.ProgressWriter, "Planning compose proposal...")
	proposal, err := p.engine.ProposeDemux(ctx, opts)
	if err != nil {
		return DemuxPlanPacket{}, err
	}
	if !opts.PlanOnly {
		demuxProgress(opts.ProgressWriter, "Fixing compose proposal %s...", proposal.ID)
		result, err := p.repairProposal(ctx, proposal, DemuxAIReviewOptions{
			Model:       opts.Model,
			MaxWarnings: opts.MaxWarnings,
		})
		if err != nil {
			if strings.TrimSpace(result.Proposal.ID) != "" {
				packet, packetErr := p.packetForProposal(ctx, result.Proposal)
				if packetErr == nil {
					return packet, err
				}
			}
			return DemuxPlanPacket{}, err
		}
		proposal = result.Proposal
		proposal, err = p.preflightApplyReadyProposal(ctx, proposal, opts)
		if err != nil {
			if strings.TrimSpace(proposal.ID) != "" {
				packet, packetErr := p.packetForProposal(ctx, proposal)
				if packetErr == nil {
					return packet, err
				}
			}
			return DemuxPlanPacket{}, err
		}
	} else {
		demuxProgress(opts.ProgressWriter, "Skipping compose fix because --plan was set.")
	}
	return p.packetForProposal(ctx, proposal)
}

func demuxProgress(out io.Writer, format string, args ...any) {
	if out == nil {
		return
	}
	fmt.Fprintf(out, format+"\n", args...)
}

func (p demuxPipeline) showProposal(ctx context.Context, selector string) (DemuxPlanPacket, error) {
	proposal, err := p.engine.LoadDemuxProposalBySelector(ctx, selector)
	if err != nil {
		return DemuxPlanPacket{}, err
	}
	return p.packetForProposal(ctx, proposal)
}

func (p demuxPipeline) repairSavedProposal(ctx context.Context, selector string, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	if opts.MaxWarnings <= 0 {
		opts.MaxWarnings = envInt("GX_DEMUX_REVIEW_MAX_WARNINGS", defaultDemuxReviewMaxWarnings)
	}
	proposal, err := p.engine.LoadDemuxProposalBySelector(ctx, selector)
	if err != nil {
		return DemuxAIReviewResult{}, err
	}
	return p.repairProposal(ctx, proposal, opts)
}

func (p demuxPipeline) repairProposal(ctx context.Context, proposal DemuxProposal, opts DemuxAIReviewOptions) (DemuxAIReviewResult, error) {
	if opts.MaxWarnings <= 0 {
		opts.MaxWarnings = envInt("GX_DEMUX_REVIEW_MAX_WARNINGS", defaultDemuxReviewMaxWarnings)
	}
	return p.engine.demuxRepair().repair(ctx, proposal, opts)
}

func (p demuxPipeline) preflightApplyReadyProposal(ctx context.Context, proposal DemuxProposal, opts ProposeDemuxOptions) (DemuxProposal, error) {
	if demuxWorkflowState(ReviewDemuxResult{Valid: true, Proposal: proposal}, proposal) != DemuxWorkflowReadyToApply {
		return proposal, nil
	}
	attempts := envInt("GX_COMPOSE_APPLY_PREFLIGHT_ATTEMPTS", defaultDemuxApplyPreflightAttempts)
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		demuxProgress(opts.ProgressWriter, "Checking compose apply in disposable attempt %d/%d...", attempt, attempts)
		if err := p.engine.PreflightDemuxApply(ctx, proposal); err == nil {
			return proposal, nil
		} else {
			lastErr = err
		}
		proposal.Warnings = appendDemuxApplyPreflightWarning(proposal.Warnings, lastErr)
		if attempt == attempts {
			saved, saveErr := p.engine.SaveDemuxProposal(ctx, proposal)
			if saveErr == nil {
				proposal = saved
			}
			return proposal, fmt.Errorf("compose apply preflight failed after %d attempt(s): %w", attempt, lastErr)
		}
		demuxProgress(opts.ProgressWriter, "Repairing compose proposal after apply preflight failure...")
		result, err := p.engine.demuxRepair().repairApplyPreflightFailure(ctx, proposal, lastErr, DemuxAIReviewOptions{
			Model:       opts.Model,
			MaxWarnings: opts.MaxWarnings,
		})
		if err != nil {
			if strings.TrimSpace(result.Proposal.ID) != "" {
				proposal = result.Proposal
			}
			return proposal, fmt.Errorf("compose apply preflight failed and repair did not produce a new apply-ready proposal: %w; preflight: %v", err, lastErr)
		}
		proposal = result.Proposal
	}
	return proposal, lastErr
}

func appendDemuxApplyPreflightWarning(warnings []string, err error) []string {
	message := "Compose apply preflight failed: " + err.Error()
	out := make([]string, 0, len(warnings)+1)
	for _, warning := range warnings {
		if !strings.HasPrefix(warning, "Compose apply preflight failed: ") {
			out = append(out, warning)
		}
	}
	return append(out, message)
}

func (p demuxPipeline) packetForProposal(ctx context.Context, proposal DemuxProposal) (DemuxPlanPacket, error) {
	review, err := p.engine.ReviewDemuxPlan(ctx, proposal)
	if err != nil {
		return DemuxPlanPacket{}, err
	}
	current := review.Proposal
	if strings.TrimSpace(current.ID) == "" {
		current = proposal
	}
	state := demuxWorkflowState(review, current)
	nextTool := "gx_review_revision_plan"
	if state == DemuxWorkflowReadyToApply {
		nextTool = "gx_apply_revision_plan"
	}
	return DemuxPlanPacket{
		Action:              "demux_changes",
		State:               state,
		NextTool:            nextTool,
		FinalTool:           "gx_apply_revision_plan",
		Workflow:            demuxWorkflowInstructions(),
		PlanningContract:    demuxPlanningContract(),
		RevisionShape:       demuxRevisionShape(),
		ReviewArgumentShape: demuxToolArgumentShape(current),
		ApplyArgumentShape:  demuxToolArgumentShape(current),
		Proposal:            current,
		Review:              review,
		StructuralFacts:     current.StructuralFacts,
		StructuralDeps:      current.StructuralDeps,
		ChangedSymbols:      current.ChangedSymbols,
		FeasibilityWarnings: current.FeasibilityWarnings,
		Warnings:            current.Warnings,
	}, nil
}
