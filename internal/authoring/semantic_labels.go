package authoring

import (
	"context"
	"encoding/json"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	semanticLabelSourceBM25      = "local_bm25"
	semanticLabelStatusCandidate = "candidate"
	semanticLabelStatusAccepted  = "accepted"
)

var semanticLabelWordPattern = regexp.MustCompile(`[a-z0-9]+`)

type semanticLabelDocument struct {
	Label    string
	Text     string
	Source   string
	Status   string
	Evidence []string
}

func (e *Engine) annotateSemanticLabels(ctx context.Context, proposal DemuxProposal) DemuxProposal {
	store, repoID, closeStore, err := semanticLabelStore(ctx, proposal.RepoRoot)
	if err != nil {
		return proposal
	}
	defer closeStore()
	e.pruneSemanticLabelRegistry(ctx, store, repoID)
	corpus := e.semanticLabelCorpus(ctx, proposal, store, repoID)
	for index := range proposal.Revisions {
		proposal.Revisions[index].SemanticLabels = semanticLabelsForRevision(ctx, store, proposal, proposal.Revisions[index], corpus)
	}
	return proposal
}

func (e *Engine) semanticLabelCorpus(ctx context.Context, proposal DemuxProposal, store *storage.Store, repoID int64) []semanticLabelDocument {
	docs := semanticLabelDocumentsFromProposal(proposal)
	if store != nil && repoID != 0 {
		if labels, err := store.ListSemanticLabelsByRepoID(ctx, repoID, 120); err == nil {
			for _, label := range labels {
				if label.AcceptedCount <= 0 {
					continue
				}
				text := label.Label + " " + label.AliasesJSON + " " + label.SourcesJSON
				docs = append(docs, semanticLabelDocument{
					Label:  label.Label,
					Text:   text,
					Source: "registry",
					Status: semanticLabelStatusAccepted,
					Evidence: []string{
						"registry seen " + intString(label.SeenCount),
						"registry accepted " + intString(label.AcceptedCount),
					},
				})
			}
		}
	}
	if index, _, err := e.demuxStackIndex(ctx, proposal.RepoRoot); err == nil {
		for _, stack := range index.stacks {
			if IsTerminalDemuxStackStatus(stack.Status) {
				continue
			}
			text := stack.Name + " " + stack.BookmarkName + " " + stackNameFromBookmark(stack.BookmarkName)
			for _, change := range index.stackChanges[stack.ID] {
				text += " " + change.Description
			}
			if label := canonicalSemanticLabel(stack.Name + " " + stackNameFromBookmark(stack.BookmarkName)); label != "" {
				docs = append(docs, semanticLabelDocument{
					Label:    label,
					Text:     text,
					Source:   "stack_history",
					Status:   semanticLabelStatusAccepted,
					Evidence: []string{"stack " + stack.BookmarkName},
				})
			}
		}
	}
	return dedupeSemanticLabelDocuments(docs)
}

func semanticLabelDocumentsFromProposal(proposal DemuxProposal) []semanticLabelDocument {
	var docs []semanticLabelDocument
	for _, revision := range proposal.Revisions {
		text := semanticRevisionText(proposal, revision)
		for _, label := range semanticCandidateLabels(revision) {
			docs = append(docs, semanticLabelDocument{
				Label:    label,
				Text:     text + " " + label,
				Source:   "current_diff",
				Status:   semanticLabelStatusCandidate,
				Evidence: semanticLabelEvidence(revision, text),
			})
		}
	}
	return docs
}

func semanticLabelsForRevision(ctx context.Context, store *storage.Store, proposal DemuxProposal, revision RevisionProposal, docs []semanticLabelDocument) []SemanticLabel {
	query := semanticRevisionText(proposal, revision)
	ranked, err := rankSemanticLabelDocuments(ctx, store, query, docs)
	if err != nil {
		return fallbackSemanticLabelsForRevision(revision)
	}
	out := make([]SemanticLabel, 0, 3)
	seen := map[string]struct{}{}
	if label := primaryIntentSemanticLabel(revision); label != "" {
		seen[label] = struct{}{}
		out = append(out, SemanticLabel{
			Label:    label,
			Source:   semanticLabelSourceBM25,
			Status:   semanticLabelStatusCandidate,
			Score:    1,
			Evidence: semanticLabelEvidence(revision, query),
		})
	}
	for _, scored := range ranked {
		if scored.Score < 0.24 {
			continue
		}
		label := strings.TrimSpace(scored.Document.Label)
		if label == "" {
			continue
		}
		if !semanticSingleLabelAllowed(label) {
			continue
		}
		if !semanticLabelRelevantToRevision(scored.Document, revision) {
			continue
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		out = append(out, SemanticLabel{
			Label:    label,
			Source:   semanticLabelSourceBM25,
			Status:   scored.Document.Status,
			Score:    roundSemanticScore(scored.Score),
			Evidence: scored.Document.Evidence,
		})
		if len(out) >= 3 {
			break
		}
	}
	if len(out) == 0 {
		return fallbackSemanticLabelsForRevision(revision)
	}
	return out
}

func semanticLabelRelevantToRevision(doc semanticLabelDocument, revision RevisionProposal) bool {
	if doc.Source == "current_diff" {
		return true
	}
	labelTokens := semanticTokens(doc.Label)
	if len(labelTokens) == 0 {
		return false
	}
	contextTokens := semanticRevisionContextTokens(revision)
	for _, token := range labelTokens {
		if _, ok := contextTokens[token]; !ok {
			return false
		}
	}
	return true
}

func semanticRevisionContextTokens(revision RevisionProposal) map[string]struct{} {
	context := map[string]struct{}{}
	add := func(value string) {
		for _, token := range semanticTokens(value) {
			context[token] = struct{}{}
		}
	}
	add(revision.Intent)
	add(revision.TargetStack)
	for _, file := range revision.Files {
		add(file)
		add(pathSemanticLabel(file))
		add(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))
	}
	for _, hunk := range revision.Hunks {
		add(hunk.Symbol)
		add(hunk.SymbolKind)
	}
	return context
}

func primaryIntentSemanticLabel(revision RevisionProposal) string {
	label := canonicalSemanticLabel(revision.Intent)
	if label == "" || !semanticSingleLabelAllowed(label) {
		return ""
	}
	if len(semanticTokens(label)) < 2 {
		return ""
	}
	return label
}

func fallbackSemanticLabelsForRevision(revision RevisionProposal) []SemanticLabel {
	labels := semanticCandidateLabels(revision)
	if len(labels) == 0 {
		return nil
	}
	evidence := semanticLabelEvidence(revision, strings.Join(revision.Files, " "))
	out := make([]SemanticLabel, 0, 3)
	for _, label := range labels {
		if !semanticSingleLabelAllowed(label) {
			continue
		}
		out = append(out, SemanticLabel{
			Label:    label,
			Source:   semanticLabelSourceBM25,
			Status:   semanticLabelStatusCandidate,
			Score:    0.62,
			Evidence: evidence,
		})
		if len(out) >= 3 {
			break
		}
	}
	return out
}

type scoredSemanticLabelDocument struct {
	Document semanticLabelDocument
	Score    float64
}

func rankSemanticLabelDocuments(ctx context.Context, store *storage.Store, query string, docs []semanticLabelDocument) ([]scoredSemanticLabelDocument, error) {
	if store == nil || len(docs) == 0 || len(semanticTokens(query)) == 0 {
		return nil, nil
	}
	searchDocs := make([]storage.SemanticLabelSearchDocument, 0, len(docs))
	for index, doc := range docs {
		evidenceJSON, _ := json.Marshal(doc.Evidence)
		searchDocs = append(searchDocs, storage.SemanticLabelSearchDocument{
			Ordinal:      index,
			Label:        doc.Label,
			Text:         doc.Text,
			Source:       doc.Source,
			Status:       doc.Status,
			EvidenceJSON: string(evidenceJSON),
		})
	}
	results, err := store.RankSemanticLabelDocuments(ctx, searchDocs, query, 12)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	best := 0.0
	for _, result := range results {
		score := math.Abs(result.RawScore)
		if score > best {
			best = score
		}
	}
	if best <= 0 {
		best = 1
	}
	ranked := make([]scoredSemanticLabelDocument, 0, len(results))
	for _, result := range results {
		if result.Ordinal < 0 || result.Ordinal >= len(docs) {
			continue
		}
		score := math.Abs(result.RawScore) / best
		doc := docs[result.Ordinal]
		switch doc.Source {
		case "current_diff":
			score *= 1.18
		case "stack_history":
			score *= 0.88
		}
		if doc.Status == semanticLabelStatusAccepted {
			score *= 1.03
		}
		if !semanticSingleLabelAllowed(doc.Label) {
			score *= 0.55
		}
		if score > 1 {
			score = 1
		}
		ranked = append(ranked, scoredSemanticLabelDocument{
			Document: doc,
			Score:    score,
		})
	}
	return ranked, nil
}

func semanticRevisionText(proposal DemuxProposal, revision RevisionProposal) string {
	parts := []string{revision.Intent}
	parts = append(parts, revision.Files...)
	for _, hunk := range revision.Hunks {
		parts = append(parts, hunk.File, hunk.Symbol, hunk.SymbolKind)
	}
	for _, hunkID := range revision.HunkIDs {
		for _, hunk := range proposal.Hunks {
			if hunk.ID == hunkID {
				parts = append(parts, hunk.File, hunk.Symbol, hunk.SymbolKind)
			}
		}
	}
	return strings.Join(parts, " ")
}

func semanticCandidateLabels(revision RevisionProposal) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(label string) {
		if len(out) >= 18 {
			return
		}
		label = canonicalSemanticLabel(label)
		if label == "" {
			return
		}
		if !semanticSingleLabelAllowed(label) {
			return
		}
		if _, ok := seen[label]; ok {
			return
		}
		seen[label] = struct{}{}
		out = append(out, label)
	}
	addPhrases := func(value string) {
		tokens := semanticTokens(value)
		for i := 0; i < len(tokens); i++ {
			if i+1 < len(tokens) {
				add(tokens[i] + " " + tokens[i+1])
			}
		}
	}
	add(revision.Intent)
	addPhrases(revision.Intent)
	add(revision.TargetStack)
	for _, file := range revision.Files {
		add(pathSemanticLabel(file))
		add(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)))
	}
	for _, hunk := range revision.Hunks {
		if semanticTestFile(hunk.File) {
			continue
		}
		add(hunk.Symbol)
		if hunk.SymbolKind != "" && hunk.Symbol != "" {
			add(hunk.SymbolKind + " " + hunk.Symbol)
		}
	}
	return out
}

func pathSemanticLabel(file string) string {
	file = strings.Trim(filepath.ToSlash(strings.TrimSpace(file)), "/")
	if file == "" {
		return ""
	}
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[:2], " ")
	}
	return parts[0]
}

func semanticLabelEvidence(revision RevisionProposal, text string) []string {
	var evidence []string
	if strings.TrimSpace(revision.Intent) != "" {
		evidence = append(evidence, "intent "+revision.Intent)
	}
	for _, file := range revision.Files {
		evidence = append(evidence, "file "+file)
		if len(evidence) >= 4 {
			break
		}
	}
	if len(evidence) == 0 && strings.TrimSpace(text) != "" {
		evidence = append(evidence, "diff text")
	}
	return evidence
}

func dedupeSemanticLabelDocuments(docs []semanticLabelDocument) []semanticLabelDocument {
	seen := map[string]semanticLabelDocument{}
	order := []string{}
	for _, doc := range docs {
		doc.Label = canonicalSemanticLabel(doc.Label)
		if doc.Label == "" {
			continue
		}
		key := doc.Label + "\x00" + doc.Source + "\x00" + doc.Status
		existing, ok := seen[key]
		if ok {
			existing.Text = strings.TrimSpace(existing.Text + " " + doc.Text)
			existing.Evidence = appendUniqueSemanticStrings(existing.Evidence, doc.Evidence...)
			seen[key] = existing
			continue
		}
		seen[key] = doc
		order = append(order, key)
	}
	out := make([]semanticLabelDocument, 0, len(order))
	for _, key := range order {
		out = append(out, seen[key])
	}
	return out
}

func canonicalSemanticLabel(value string) string {
	tokens := semanticTokens(value)
	if len(tokens) == 0 {
		return ""
	}
	tokens = normalizeSemanticLabelTokens(tokens)
	if len(tokens) > 3 {
		tokens = tokens[:3]
	}
	return strings.Join(tokens, " ")
}

func normalizeSemanticLabelTokens(tokens []string) []string {
	if len(tokens) == 0 {
		return tokens
	}
	seen := map[string]struct{}{}
	for index, token := range tokens {
		if token == "label" {
			token = "labels"
			tokens[index] = token
		}
		seen[token] = struct{}{}
	}
	has := func(token string) bool {
		_, ok := seen[token]
		return ok
	}
	allSame := true
	for _, token := range tokens[1:] {
		if token != tokens[0] {
			allSame = false
			break
		}
	}
	if allSame {
		switch tokens[0] {
		case "e2e":
			return []string{"e2e", "coverage"}
		case "documentation", "docs", "readme":
			return []string{"gx", "documentation"}
		case "demux":
			return []string{"demux", "routing"}
		case "semantic":
			return []string{"semantic", "labels"}
		default:
			return []string{tokens[0]}
		}
	}
	switch {
	case has("e2e") && has("coverage"):
		return []string{"e2e", "coverage"}
	case has("e2e") && has("gx"):
		return []string{"e2e", "coverage"}
	case has("gx") && has("workflow"):
		return []string{"gx", "workflow"}
	case has("gx") && has("documentation"):
		return []string{"gx", "documentation"}
	case has("documentation") && has("readme"):
		return []string{"gx", "documentation"}
	case has("docs") && has("documentation"):
		return []string{"gx", "documentation"}
	case has("apps") && has("menubar"):
		return []string{"menubar", "app"}
	case has("registry") && has("storage"):
		return []string{"semantic", "label", "storage"}
	case has("generate") && has("pipeline"):
		return []string{"generate", "pipeline"}
	case has("confidence") && has("metadata"):
		return []string{"confidence", "metadata"}
	case has("semantic") && has("labels"):
		return []string{"semantic", "labels"}
	case has("review") && has("context"):
		return []string{"review", "context"}
	case has("provider") && has("runtime"):
		return []string{"provider", "runtime"}
	case has("stack") && has("management"):
		return []string{"stack", "management"}
	default:
		return dedupeSemanticTokens(tokens)
	}
}

func dedupeSemanticTokens(tokens []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		out = append(out, token)
	}
	return out
}

func semanticSingleLabelAllowed(label string) bool {
	tokens := semanticTokens(label)
	if len(tokens) != 1 {
		return true
	}
	switch tokens[0] {
	case "cli", "e2e", "schema", "sql", "documentation", "storage":
		return true
	default:
		return false
	}
}

func semanticTokens(value string) []string {
	value = splitCamelSemantic(value)
	value = strings.ToLower(value)
	raw := semanticLabelWordPattern.FindAllString(value, -1)
	out := make([]string, 0, len(raw))
	for _, token := range raw {
		if semanticStopWord(token) {
			continue
		}
		out = append(out, token)
	}
	return out
}

func semanticTermCounts(value string) map[string]int {
	counts := map[string]int{}
	for _, token := range semanticTokens(value) {
		counts[token]++
	}
	return counts
}

func splitCamelSemantic(value string) string {
	var out []rune
	var prev rune
	for i, r := range value {
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)) {
			out = append(out, ' ')
		}
		if r == '_' || r == '-' || r == '/' || r == '.' {
			out = append(out, ' ')
		} else {
			out = append(out, r)
		}
		prev = r
	}
	return string(out)
}

func semanticStopWord(token string) bool {
	if len(token) <= 1 {
		return true
	}
	if strings.HasPrefix(token, "h") {
		rest := strings.TrimPrefix(token, "h")
		if rest != "" && strings.Trim(rest, "0123456789") == "" {
			return true
		}
	}
	if strings.Trim(token, "0123456789") == "" {
		return true
	}
	switch token {
	case "add", "adds", "added", "update", "updates", "updated", "change", "changes", "changed",
		"fix", "fixes", "fixed", "new", "old", "file", "files", "test", "tests", "go", "md", "mdx",
		"internal", "src", "app", "cmd", "root", "type", "function", "method", "with", "from", "feature",
		"into", "in", "for", "and", "the", "this", "that", "true", "false", "revision", "revisions",
		"matched", "match", "matches", "local", "remote", "current", "draft", "shape", "package",
		"structural", "dependency", "cluster", "route", "routed", "routing", "unlocked", "index":
		return true
	default:
		return false
	}
}

func semanticTestFile(file string) bool {
	file = filepath.ToSlash(strings.TrimSpace(file))
	return strings.HasPrefix(file, "test/") || strings.HasSuffix(file, "_test.go")
}

func semanticLabelStackOverlapScore(labels []SemanticLabel, stackText string) float64 {
	stackTokens := map[string]struct{}{}
	for _, token := range semanticTokens(stackText) {
		stackTokens[token] = struct{}{}
	}
	if len(stackTokens) == 0 {
		return 0
	}
	best := 0.0
	for _, label := range labels {
		tokens := semanticTokens(label.Label)
		if len(tokens) == 0 {
			continue
		}
		matched := 0
		for _, token := range tokens {
			if _, ok := stackTokens[token]; ok {
				matched++
			}
		}
		score := float64(matched) / float64(len(tokens))
		score *= label.Score
		if label.Status == semanticLabelStatusAccepted {
			score *= 1.1
		}
		if score > best {
			best = score
		}
	}
	if best > 1 {
		return 1
	}
	return best
}

func primarySemanticLabel(revision RevisionProposal) (string, float64, bool) {
	for _, label := range revision.SemanticLabels {
		if label.Score >= 0.65 && strings.TrimSpace(label.Label) != "" {
			return label.Label, label.Score, true
		}
	}
	return "", 0, false
}

func (e *Engine) recordProposalSemanticLabels(ctx context.Context, proposal DemuxProposal) {
	// Draft proposals are not label authority. Persist labels only after apply so
	// repeated generate attempts do not train BM25 on rejected or stale drafts.
}

func (e *Engine) recordAppliedSemanticLabels(ctx context.Context, store *storage.Store, repoID int64, proposal DemuxProposal, revision RevisionProposal, changeID int64) error {
	return e.recordSemanticLabels(ctx, store, repoID, proposal, revision, &changeID, true)
}

func (e *Engine) recordSemanticLabels(ctx context.Context, store *storage.Store, repoID int64, proposal DemuxProposal, revision RevisionProposal, changeID *int64, accepted bool) error {
	now := time.Now().UnixMilli()
	return store.WriteSemanticLabels(ctx, semanticLabelWrites(repoID, proposal, revision, changeID, accepted, now))
}

func semanticLabelWrites(repoID int64, proposal DemuxProposal, revision RevisionProposal, changeID *int64, accepted bool, now int64) []storage.SemanticLabelWrite {
	writes := make([]storage.SemanticLabelWrite, 0, len(revision.SemanticLabels))
	for _, label := range revision.SemanticLabels {
		canonical := canonicalSemanticLabel(label.Label)
		if canonical == "" {
			continue
		}
		aliasesJSON, _ := json.Marshal([]string{label.Label})
		sourcesJSON, _ := json.Marshal([]string{label.Source})
		acceptedCount := 0
		status := semanticLabelStatusCandidate
		if accepted {
			acceptedCount = 1
			status = semanticLabelStatusAccepted
		}
		labelRow := storage.SemanticLabel{
			RepoID:        repoID,
			Label:         canonical,
			AliasesJSON:   string(aliasesJSON),
			SourcesJSON:   string(sourcesJSON),
			SeenCount:     1,
			AcceptedCount: acceptedCount,
			Confidence:    label.Score,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		evidenceJSON, _ := json.Marshal(map[string]any{
			"label":                canonical,
			"status":               status,
			"evidence":             label.Evidence,
			"target_stack":         revision.TargetStack,
			"revision_proposal_id": revision.ID,
		})
		proposalID := proposal.ID
		revisionID := revision.ID
		stack := strings.TrimSpace(revision.TargetStack)
		var stackPtr *string
		if stack != "" {
			stackPtr = &stack
		}
		writes = append(writes, storage.SemanticLabelWrite{
			Label: labelRow,
			Link: storage.SemanticLabelLink{
				RepoID:             repoID,
				DemuxProposalID:    &proposalID,
				RevisionProposalID: &revisionID,
				ChangeID:           changeID,
				StackBookmark:      stackPtr,
				Source:             label.Source,
				Score:              label.Score,
				Accepted:           accepted,
				EvidenceJSON:       string(evidenceJSON),
				CreatedAt:          now,
			},
		})
	}
	return writes
}

func (e *Engine) pruneSemanticLabelRegistry(ctx context.Context, store *storage.Store, repoID int64) {
	labels, err := store.ListSemanticLabelsByRepoID(ctx, repoID, 1000)
	if err != nil || len(labels) == 0 {
		return
	}
	merges := map[string]string{}
	for _, label := range labels {
		from := strings.TrimSpace(label.Label)
		to := semanticLabelMergeTarget(from)
		if from == "" || to == "" || from == to {
			continue
		}
		merges[from] = to
	}
	if len(merges) == 0 {
		return
	}
	_ = store.MergeSemanticLabels(ctx, repoID, merges, time.Now().UnixMilli())
}

func semanticLabelMergeTarget(label string) string {
	canonical := canonicalSemanticLabel(label)
	lower := strings.ToLower(strings.TrimSpace(label))
	switch lower {
	case "coverage gx":
		return "e2e coverage"
	case "docs index mdx", "readme docs documentation", "docs documentation", "documentation readme":
		return "gx documentation"
	case "e2e gx", "e2e gx e2e":
		return "e2e coverage"
	case "e2e e2e":
		return "e2e coverage"
	case "documentation documentation":
		return "gx documentation"
	case "demux demux":
		return "demux routing"
	case "gxgenerate attaches repo":
		return "gx generate"
	case "gxsync pulls cloud":
		return "gx sync"
	}
	if strings.HasPrefix(lower, "h") {
		tokens := semanticTokens(lower)
		if len(tokens) == 0 {
			return ""
		}
	}
	return canonical
}

func semanticLabelStore(ctx context.Context, repoRoot string) (*storage.Store, int64, func(), error) {
	db, err := storage.Open(ctx)
	if err != nil {
		return nil, 0, func() {}, err
	}
	closeStore := func() { _ = db.Close() }
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		closeStore()
		return nil, 0, func() {}, err
	}
	repo, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil || repo == nil {
		closeStore()
		return nil, 0, func() {}, err
	}
	return store, repo.ID, closeStore, nil
}

func appendUniqueSemanticStrings(values []string, extra ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values)+len(extra))
	for _, value := range append(values, extra...) {
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

func roundSemanticScore(score float64) float64 {
	return math.Round(score*100) / 100
}

func minFloat(left, right float64) float64 {
	if left < right {
		return left
	}
	return right
}

func intString(value int) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	n := value
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
