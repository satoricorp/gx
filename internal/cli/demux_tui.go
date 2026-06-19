package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/satoricorp/gx/internal/authoring"
)

type demuxInteractiveAction struct {
	Kind        string
	Target      string
	Proposal    authoring.DemuxProposal
	ProposalSet bool
}

type demuxInteractiveMode int

const (
	demuxInteractiveStacks demuxInteractiveMode = iota
	demuxInteractiveRevisions
	demuxInteractiveDiff
)

type demuxInteractiveModel struct {
	proposal  authoring.DemuxProposal
	groups    []demuxStackDisplayGroup
	mode      demuxInteractiveMode
	stack     int
	revision  int
	selected  map[string]struct{}
	notice    string
	diffView  tuiDiffView
	diffTitle string
	action    demuxInteractiveAction
}

func runDemuxInteractive(in io.Reader, out io.Writer, proposal authoring.DemuxProposal) (demuxInteractiveAction, error) {
	model := newDemuxInteractiveModel(proposal)
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	final, err := program.Run()
	if err != nil {
		return demuxInteractiveAction{}, err
	}
	if model, ok := final.(demuxInteractiveModel); ok {
		return model.action, nil
	}
	return demuxInteractiveAction{}, nil
}

func newDemuxInteractiveModel(proposal authoring.DemuxProposal) demuxInteractiveModel {
	model := demuxInteractiveModel{
		proposal: proposal,
		groups:   demuxStackDisplayGroups(proposal.Revisions),
		mode:     demuxInteractiveStacks,
		selected: map[string]struct{}{},
		diffView: newTUIDiffView(),
	}
	model.clamp()
	return model
}

func (m demuxInteractiveModel) Init() tea.Cmd {
	return nil
}

func (m demuxInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == demuxInteractiveDiff {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.diffView.resize(msg.Width, msg.Height)
			return m, nil
		case tea.KeyPressMsg:
			switch msg.String() {
			case "esc":
				m.mode = demuxInteractiveRevisions
				return m, nil
			case "ctrl+c", "q":
				m.action = demuxInteractiveAction{Kind: "quit"}
				return m, tea.Quit
			}
		}
		view, cmd := m.diffView.view.Update(msg)
		m.diffView.view = view
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.diffView.resize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			m.move(-1)
		case "down", "j":
			m.move(1)
		case "enter":
			if m.mode == demuxInteractiveStacks && len(m.currentRevisions()) > 0 {
				m.mode = demuxInteractiveRevisions
				m.revision = 0
			}
		case "esc":
			if m.mode == demuxInteractiveRevisions {
				m.mode = demuxInteractiveStacks
			} else {
				m.action = demuxInteractiveAction{Kind: "quit"}
				return m, tea.Quit
			}
		case "A":
			m.action = demuxInteractiveAction{Kind: "apply_all", Proposal: m.proposal, ProposalSet: true}
			return m, tea.Quit
		case "s", "S":
			if group := m.selectedGroup(); group != nil {
				m.action = demuxInteractiveAction{Kind: "apply_stack", Target: group.Target, Proposal: m.proposal, ProposalSet: true}
				return m, tea.Quit
			}
		case " ", "space":
			if m.mode == demuxInteractiveRevisions {
				m.toggleSelectedRevision()
			}
		case "d", "D":
			if m.mode == demuxInteractiveRevisions {
				m.openSelectedRevisionDiff()
			}
		case "c", "C":
			if m.mode == demuxInteractiveRevisions {
				m.combineSelectedRevisions()
			}
		case "r", "R":
			if m.mode == demuxInteractiveRevisions {
				m.releaseSelectedRevisions()
			}
		case "e", "E":
			if m.mode == demuxInteractiveRevisions {
				m.notice = "Edit selected revisions is not implemented yet."
			}
		case "x", "X":
			m.selected = map[string]struct{}{}
			m.notice = "Selection cleared."
		case "ctrl+c", "q":
			m.action = demuxInteractiveAction{Kind: "quit"}
			return m, tea.Quit
		}
	}
	m.clamp()
	return m, nil
}

func (m demuxInteractiveModel) selectedGroup() *demuxStackDisplayGroup {
	if m.stack < 0 || m.stack >= len(m.groups) {
		return nil
	}
	return &m.groups[m.stack]
}

func (m demuxInteractiveModel) currentRevisions() []authoring.RevisionProposal {
	group := m.selectedGroup()
	if group == nil {
		return nil
	}
	return group.Revisions
}

func (m demuxInteractiveModel) selectedRevision() *authoring.RevisionProposal {
	revisions := m.currentRevisions()
	if m.revision < 0 || m.revision >= len(revisions) {
		return nil
	}
	return &revisions[m.revision]
}

func (m *demuxInteractiveModel) move(delta int) {
	switch m.mode {
	case demuxInteractiveRevisions:
		m.revision += delta
	default:
		m.stack += delta
	}
	m.clamp()
}

func (m *demuxInteractiveModel) clamp() {
	if len(m.groups) == 0 {
		m.stack = 0
		m.revision = 0
		return
	}
	if m.stack < 0 {
		m.stack = 0
	}
	if m.stack >= len(m.groups) {
		m.stack = len(m.groups) - 1
	}
	revisions := m.currentRevisions()
	if len(revisions) == 0 {
		m.revision = 0
		return
	}
	if m.revision < 0 {
		m.revision = 0
	}
	if m.revision >= len(revisions) {
		m.revision = len(revisions) - 1
	}
}

func (m *demuxInteractiveModel) refreshGroups() {
	m.groups = demuxStackDisplayGroups(m.proposal.Revisions)
	m.clamp()
}

func (m *demuxInteractiveModel) toggleSelectedRevision() {
	revision := m.selectedRevision()
	if revision == nil || revision.ID == "" {
		return
	}
	if _, ok := m.selected[revision.ID]; ok {
		delete(m.selected, revision.ID)
		m.notice = fmt.Sprintf("Unselected %s.", revision.ID)
		return
	}
	m.selected[revision.ID] = struct{}{}
	m.notice = fmt.Sprintf("Selected %s.", revision.ID)
}

func (m *demuxInteractiveModel) openSelectedRevisionDiff() {
	revision := m.selectedRevision()
	if revision == nil {
		return
	}
	m.diffTitle = fmt.Sprintf("%s  %s", revision.ID, revision.Intent)
	m.diffView.open(m.diffTitle, demuxRevisionDiff(*revision))
	m.mode = demuxInteractiveDiff
}

func (m *demuxInteractiveModel) combineSelectedRevisions() {
	ids := m.selectedRevisionIDs()
	if len(ids) < 2 {
		m.notice = "Select at least 2 revisions to combine."
		return
	}
	idSet := sliceSet(ids)
	var combined authoring.RevisionProposal
	var remaining []authoring.RevisionProposal
	for _, revision := range m.proposal.Revisions {
		if _, ok := idSet[revision.ID]; !ok {
			remaining = append(remaining, revision)
			continue
		}
		if combined.ID == "" {
			combined = revision
			continue
		}
		if revision.Intent != "" && !strings.Contains(combined.Intent, revision.Intent) {
			if combined.Intent == "" {
				combined.Intent = revision.Intent
			} else {
				combined.Intent += " + " + revision.Intent
			}
		}
		combined.Files = appendUniqueStrings(combined.Files, revision.Files)
		combined.HunkIDs = appendUniqueStrings(combined.HunkIDs, revision.HunkIDs)
		combined.DependsOn = appendUniqueStrings(combined.DependsOn, revision.DependsOn)
		combined.Hunks = appendUniqueHunks(combined.Hunks, revision.Hunks)
	}
	if combined.ID == "" {
		m.notice = "Selected revisions are no longer available."
		m.selected = map[string]struct{}{}
		return
	}
	remaining = append(remaining, combined)
	m.proposal.Revisions = remaining
	m.proposal.Hunks = demuxHunksForRevisions(m.proposal, m.proposal.Revisions)
	m.selected = map[string]struct{}{}
	m.notice = fmt.Sprintf("Combined %d revisions into %s.", len(ids), combined.ID)
	m.refreshGroups()
}

func (m *demuxInteractiveModel) releaseSelectedRevisions() {
	ids := m.selectedRevisionIDs()
	if len(ids) == 0 {
		if revision := m.selectedRevision(); revision != nil && revision.ID != "" {
			ids = []string{revision.ID}
		}
	}
	if len(ids) == 0 {
		return
	}
	idSet := sliceSet(ids)
	var remaining []authoring.RevisionProposal
	for _, revision := range m.proposal.Revisions {
		if _, ok := idSet[revision.ID]; ok {
			continue
		}
		remaining = append(remaining, revision)
	}
	m.proposal.Revisions = remaining
	m.proposal.Hunks = demuxHunksForRevisions(m.proposal, m.proposal.Revisions)
	m.selected = map[string]struct{}{}
	m.notice = fmt.Sprintf("Released %d revision(s). They can be recomposed later.", len(ids))
	m.refreshGroups()
	if len(m.proposal.Revisions) == 0 {
		m.mode = demuxInteractiveStacks
	}
}

func (m demuxInteractiveModel) selectedRevisionIDs() []string {
	var ids []string
	for _, revision := range m.proposal.Revisions {
		if _, ok := m.selected[revision.ID]; ok {
			ids = append(ids, revision.ID)
		}
	}
	return ids
}

func (m demuxInteractiveModel) View() tea.View {
	if m.mode == demuxInteractiveDiff {
		return tea.NewView(m.diffView.render("gx compose"))
	}
	return tea.NewView(renderComposeSummary(m.proposal, m.groups, m.stack, m.mode == demuxInteractiveStacks, m.revision, m.selected, m.notice))
}

func renderComposeSummary(proposal authoring.DemuxProposal, groups []demuxStackDisplayGroup, stackCursor int, stackMode bool, revCursor int, selected map[string]struct{}, notice string) string {
	var lines []string
	lines = append(lines, commandLine("gx compose", true), "")
	if demuxProposalIsPartial(proposal) {
		lines = append(lines, section("Partial proposal"))
		lines = append(lines, muted(demuxPartialComposeWarning), "")
	}
	if notice != "" {
		lines = append(lines, muted(notice), "")
	}
	header := fmt.Sprintf("%d %s", len(proposal.Revisions), pluralize("revision", len(proposal.Revisions)))
	if len(groups) > 0 {
		header = fmt.Sprintf("%d %s · %d %s", len(groups), pluralize("stack", len(groups)), len(proposal.Revisions), pluralize("revision", len(proposal.Revisions)))
	}
	lines = append(lines, section("Stacks"), muted(header), "")
	if len(groups) == 0 {
		lines = append(lines, "    "+muted("(no proposed stacks)"))
	} else {
		if stackCursor < 0 {
			stackCursor = 0
		}
		if stackCursor >= len(groups) {
			stackCursor = len(groups) - 1
		}
		for index, group := range groups {
			if index > 0 {
				lines = append(lines, muted(strings.Repeat("─", 52)))
			}
			stackSelected := index == stackCursor
			marker := "○"
			if stackSelected {
				marker = "●"
			}
			meta := fmt.Sprintf("%d %s", len(group.Revisions), pluralize("revision", len(group.Revisions)))
			lines = append(lines, composeStackLine(marker, group.Label, meta, stackSelected))
			cursor := -1
			if stackSelected && !stackMode {
				cursor = revCursor
			}
			lines = append(lines, renderComposeRevisionLines(group.Revisions, cursor, selected)...)
		}
	}
	lines = append(lines, "", composeLegend(stackMode))
	return strings.Join(lines, "\n") + "\n"
}

func composeStackLine(marker, label, meta string, selected bool) string {
	renderedMarker := muted(marker)
	name := muted(label)
	if selected {
		renderedMarker = logoText(marker)
		name = valueText(label)
	}
	line := renderedMarker + " " + name
	if meta != "" {
		line += "  " + muted(meta)
	}
	return line
}

func renderComposeRevisionLines(revisions []authoring.RevisionProposal, cursor int, selected map[string]struct{}) []string {
	if len(revisions) == 0 {
		return []string{"    " + muted("(no revisions)")}
	}
	if cursor >= len(revisions) {
		cursor = len(revisions) - 1
	}
	lines := make([]string, 0, len(revisions))
	for index, revision := range revisions {
		marker := " "
		id := command(revision.ID)
		if index == cursor {
			marker = logoText("›")
			id = accent(revision.ID)
		} else if _, ok := selected[revision.ID]; ok {
			marker = success("✓")
		}
		line := fmt.Sprintf("    %s %s %s", marker, id, value(revision.Intent))
		lines = append(lines, line)
		for _, file := range revision.Files {
			lines = append(lines, "      "+muted(file))
		}
		for _, hunk := range revision.Hunks {
			if hunk.File != "" && !stringInSlice(hunk.File, revision.Files) {
				lines = append(lines, "      "+muted(hunk.File))
			}
		}
	}
	return lines
}

func composeLegend(stackMode bool) string {
	parts := []string{"● selected", "○ other", "space select", "c combine", "r release"}
	if stackMode {
		parts = append([]string{"j/k stack", "enter revisions", "s accept stack", "Shift+A accept all", "q quit"}, parts...)
	} else {
		parts = append([]string{"j/k revision", "d diff", "esc stacks", "s accept stack", "Shift+A accept all", "q quit"}, parts...)
	}
	return muted(strings.Join(parts, " · "))
}

func demuxRevisionDiff(revision authoring.RevisionProposal) string {
	var b strings.Builder
	wrotePatch := false
	for _, hunk := range revision.Hunks {
		if strings.TrimSpace(hunk.Patch) == "" {
			continue
		}
		if hunk.File != "" {
			fmt.Fprintf(&b, "diff --git a/%s b/%s\n", hunk.File, hunk.File)
		}
		b.WriteString(strings.TrimRight(hunk.Patch, "\n"))
		b.WriteString("\n")
		wrotePatch = true
	}
	if wrotePatch {
		return b.String()
	}
	if len(revision.Files) > 0 {
		fmt.Fprintf(&b, "No hunk patch text is available for %s.\n\nFiles:\n", revision.ID)
		for _, file := range revision.Files {
			fmt.Fprintf(&b, "  %s\n", file)
		}
		return b.String()
	}
	return fmt.Sprintf("No proposed diff is available for %s.\n", revision.ID)
}

func sliceSet(values []string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, value := range values {
		out[value] = struct{}{}
	}
	return out
}

func appendUniqueStrings(base []string, next []string) []string {
	seen := sliceSet(base)
	for _, value := range next {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		base = append(base, value)
	}
	return base
}

func appendUniqueHunks(base []authoring.HunkRange, next []authoring.HunkRange) []authoring.HunkRange {
	seen := map[string]struct{}{}
	for _, hunk := range base {
		if hunk.ID != "" {
			seen[hunk.ID] = struct{}{}
		}
	}
	for _, hunk := range next {
		if hunk.ID != "" {
			if _, ok := seen[hunk.ID]; ok {
				continue
			}
			seen[hunk.ID] = struct{}{}
		}
		base = append(base, hunk)
	}
	return base
}

func applyDemuxInteractiveAction(ctx context.Context, engine *authoring.Engine, packet authoring.DemuxPlanPacket, action demuxInteractiveAction, returnToDefaultBranch bool) (authoring.ApplyDemuxResult, bool, error) {
	proposal := action.Proposal
	if !action.ProposalSet {
		proposal = packet.Proposal
	}
	applyOptions := authoring.ApplyDemuxOptions{
		AllowWarnings:         true,
		ReturnToDefaultBranch: returnToDefaultBranch,
	}
	switch action.Kind {
	case "apply_all":
		result, err := engine.ApplyDemuxPlanWithOptions(ctx, proposal, applyOptions)
		return result, true, err
	case "apply_stack":
		filtered := filterDemuxProposalForStack(proposal, action.Target)
		result, err := engine.ApplyDemuxPlanWithOptions(ctx, filtered, applyOptions)
		return result, true, err
	default:
		return authoring.ApplyDemuxResult{}, false, nil
	}
}

func filterDemuxProposalForStack(proposal authoring.DemuxProposal, target string) authoring.DemuxProposal {
	target = strings.TrimSpace(target)
	filtered := proposal
	filtered.ID = ""
	filtered.Status = authoring.ProposalPending
	filtered.Revisions = nil
	selectedIDs := map[string]struct{}{}
	for _, revision := range proposal.Revisions {
		if strings.TrimSpace(revision.TargetStack) == target {
			filtered.Revisions = append(filtered.Revisions, revision)
			selectedIDs[revision.ID] = struct{}{}
		}
	}
	filtered.Hunks = demuxHunksForRevisions(proposal, filtered.Revisions)
	filtered.FeasibilityWarnings = filterDemuxWarningsForRevisions(proposal.FeasibilityWarnings, selectedIDs)
	filtered.StructuralFacts = filterDemuxStructuralFactsForRevisions(proposal.StructuralFacts, filtered.Revisions)
	filtered.StructuralDeps = filterDemuxStructuralDepsForRevisions(proposal.StructuralDeps, filtered.Revisions)
	filtered.ChangedSymbols = filterDemuxChangedSymbolsForRevisions(proposal.ChangedSymbols, filtered.Revisions)
	return filtered
}

func demuxHunksForRevisions(proposal authoring.DemuxProposal, revisions []authoring.RevisionProposal) []authoring.HunkRange {
	seen := map[string]struct{}{}
	var hunks []authoring.HunkRange
	neededHunkIDs := map[string]struct{}{}
	neededFiles := map[string]struct{}{}
	for _, revision := range revisions {
		for _, hunkID := range revision.HunkIDs {
			if hunkID != "" {
				neededHunkIDs[hunkID] = struct{}{}
			}
		}
		for _, hunk := range revision.Hunks {
			if hunk.ID == "" {
				hunks = append(hunks, hunk)
				continue
			}
			if _, ok := seen[hunk.ID]; ok {
				continue
			}
			seen[hunk.ID] = struct{}{}
			hunks = append(hunks, hunk)
		}
		if !revision.UseHunks {
			for _, file := range revision.Files {
				if file != "" {
					neededFiles[file] = struct{}{}
				}
			}
		}
	}
	for _, hunk := range proposal.Hunks {
		if hunk.ID != "" {
			if _, ok := neededHunkIDs[hunk.ID]; ok {
				if _, duplicate := seen[hunk.ID]; !duplicate {
					seen[hunk.ID] = struct{}{}
					hunks = append(hunks, hunk)
				}
				continue
			}
		}
		if _, ok := neededFiles[hunk.File]; ok {
			if hunk.ID != "" {
				if _, duplicate := seen[hunk.ID]; duplicate {
					continue
				}
				seen[hunk.ID] = struct{}{}
			}
			hunks = append(hunks, hunk)
		}
	}
	return hunks
}

func filterDemuxWarningsForRevisions(warnings []authoring.FeasibilityWarning, selectedIDs map[string]struct{}) []authoring.FeasibilityWarning {
	var out []authoring.FeasibilityWarning
	for _, warning := range warnings {
		if warning.RevisionID == "" {
			continue
		}
		if _, ok := selectedIDs[warning.RevisionID]; ok {
			out = append(out, warning)
		}
	}
	return out
}

func filterDemuxStructuralFactsForRevisions(facts []authoring.StructuralFact, revisions []authoring.RevisionProposal) []authoring.StructuralFact {
	files := demuxRevisionFileSet(revisions)
	var out []authoring.StructuralFact
	for _, fact := range facts {
		if _, ok := files[fact.File]; ok {
			out = append(out, fact)
		}
	}
	return out
}

func filterDemuxStructuralDepsForRevisions(deps []authoring.StructuralDependency, revisions []authoring.RevisionProposal) []authoring.StructuralDependency {
	files := demuxRevisionFileSet(revisions)
	var out []authoring.StructuralDependency
	for _, dep := range deps {
		_, from := files[dep.FromFile]
		_, to := files[dep.ToFile]
		if from || to {
			out = append(out, dep)
		}
	}
	return out
}

func filterDemuxChangedSymbolsForRevisions(symbols []authoring.ChangedSymbol, revisions []authoring.RevisionProposal) []authoring.ChangedSymbol {
	files := demuxRevisionFileSet(revisions)
	hunks := map[string]struct{}{}
	for _, revision := range revisions {
		for _, hunkID := range revision.HunkIDs {
			hunks[hunkID] = struct{}{}
		}
	}
	var out []authoring.ChangedSymbol
	for _, symbol := range symbols {
		_, file := files[symbol.File]
		_, hunk := hunks[symbol.HunkID]
		if file || hunk {
			out = append(out, symbol)
		}
	}
	return out
}

func demuxRevisionFileSet(revisions []authoring.RevisionProposal) map[string]struct{} {
	files := map[string]struct{}{}
	for _, revision := range revisions {
		for _, file := range revision.Files {
			files[file] = struct{}{}
		}
		for _, hunk := range revision.Hunks {
			if hunk.File != "" {
				files[hunk.File] = struct{}{}
			}
		}
	}
	return files
}
