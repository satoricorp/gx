package codereview

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Evidence states. A review is only as good as what it managed to read, so
// every retrieval source reports which of these it ended in. The states are
// deliberately distinct: "empty" means the source answered and had nothing,
// "missing" means there is no such source to ask, and "unavailable" means the
// source exists and the question failed. Collapsing them — which is what
// swallowing the error did — makes an uninformed review indistinguishable from
// a clean one.
const (
	EvidenceOK          = "ok"
	EvidenceEmpty       = "empty"
	EvidenceMissing     = "missing"
	EvidenceUnavailable = "unavailable"
	EvidenceDisabled    = "disabled"
	// EvidenceSkipped means this review decided not to ask a source that was
	// configured and reachable. It is distinct from disabled, which is a fact
	// about the environment rather than about this review, and the reason
	// belongs in Detail: a source the reader expects to see is more alarming
	// absent from the listing than present with "skipped — docs-only change".
	EvidenceSkipped = "skipped"
)

// EvidenceStatus is what happened to one retrieval source during one review.
type EvidenceStatus struct {
	// Source is the human name of the evidence source ("code index",
	// "indexed sessions", "review knowledge", "prior review findings").
	Source string `json:"source"`
	// Namespace is the backing store the source read, when it has one.
	Namespace string `json:"namespace,omitempty"`
	// State is one of the Evidence* constants above.
	State string `json:"state"`
	// Detail explains a non-ok state in one line.
	Detail string `json:"detail,omitempty"`
	// Remedy is what the reader can do about it, when there is something.
	//
	// It is separate from Detail because a source is probed across several
	// namespaces and the fix is usually the same for all of them: folded into
	// Detail it would be repeated once per namespace in the same sentence.
	// Warnings state it once, at the end, where an instruction belongs.
	Remedy string `json:"remedy,omitempty"`
	// Snippets is how many context snippets this source contributed.
	Snippets int `json:"snippets,omitempty"`
}

// Degraded reports whether this source left the review less informed than it
// would otherwise have been. An empty source is not degraded: it answered.
func (s EvidenceStatus) Degraded() bool {
	switch s.State {
	case EvidenceMissing, EvidenceUnavailable:
		return true
	default:
		return false
	}
}

// EvidenceLog collects statuses from retrievers that run concurrently.
//
// It is carried on RetrieveInput rather than returned from Retrieve because a
// retriever reports on several namespaces at once and can succeed on one while
// failing on another; a single error return cannot say that.
type EvidenceLog struct {
	mu       sync.Mutex
	statuses []EvidenceStatus
	// parent, when set, makes this log a scope: statuses are written through to
	// the parent and counted here. See scope.
	parent  *EvidenceLog
	written int
}

// Record files one status. Safe for concurrent use.
func (l *EvidenceLog) Record(status EvidenceStatus) {
	if l == nil {
		return
	}
	status.Source = strings.TrimSpace(status.Source)
	if status.Source == "" {
		return
	}
	if strings.TrimSpace(status.State) == "" {
		status.State = EvidenceOK
	}
	l.mu.Lock()
	if l.parent != nil {
		l.written++
		l.mu.Unlock()
		l.parent.Record(status)
		return
	}
	defer l.mu.Unlock()
	l.written++
	l.statuses = append(l.statuses, status)
}

// scope returns a log that writes through to l but counts only the statuses
// filed through it.
//
// The composite hands each sub-retriever its own scope because it has to know
// whether *that* retriever reported for itself before deciding to file a
// fallback line. Counting the shared log instead answers a different question —
// "did anyone report while it was running" — and the retrievers run
// concurrently, so for the slowest source in the set the answer is nearly
// always yes. That is how the shared review-knowledge corpus came to be
// missing from most reviews: the code index files one status per namespace
// probe, and while a corpus query (an embedding call plus two vector queries)
// was still in flight those probes made it look like it had already reported.
// Its fallback line was skipped, so a review that queried the corpus and a
// review that never reached it produced identical reports — on a different
// subset of reviews every run, since it turned on which goroutine happened to
// finish first.
func (l *EvidenceLog) scope() *EvidenceLog {
	if l == nil {
		return nil
	}
	return &EvidenceLog{parent: l}
}

// recorded reports how many statuses were filed through this log. On a scope
// that is what one retriever reported about itself.
func (l *EvidenceLog) recorded() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.written
}

// Statuses returns the recorded statuses in a stable order. Retrievers run in
// parallel, so insertion order is not reproducible and would make the review
// output differ run to run for no reason.
func (l *EvidenceLog) Statuses() []EvidenceStatus {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	out := append([]EvidenceStatus(nil), l.statuses...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Source != out[j].Source {
			return out[i].Source < out[j].Source
		}
		return out[i].Namespace < out[j].Namespace
	})
	return out
}

// EvidenceWarnings renders the sources that failed this review, for the top of
// the report. Only missing and unavailable sources appear: those are per-review
// failures that make this review weaker than the last one from the same
// checkout. A source that is simply not configured is a constant of the
// environment; it is still reported, in the verbose evidence listing and in the
// JSON report, but it is not news about this review.
//
// A source is judged whole rather than per namespace, because a retriever probes
// several namespaces without knowing in advance which one holds the repository.
// Reporting each probe separately produced the review's most misleading line:
// "review evidence unavailable — code index (repo-satoricorp-yeet): gx Cloud has
// never indexed this repository. Findings are based on the change and the
// checkout only" — printed at the top of a review whose code index had just
// answered with 96 snippets from a different namespace.
//
// Whole does not mean lossy, and the two failure states are not interchangeable:
//
//   - Missing is the expected outcome of probing a name that does not exist. It
//     is suppressed once any namespace of that source answered at all, including
//     answering empty — the source was found, it simply had nothing to say.
//   - Unavailable means a namespace that does exist could not be read. No amount
//     of emptiness elsewhere makes that untrue, and a namespace that answered
//     with nothing is not evidence that the failed one would have. It is
//     suppressed only when another namespace of the same source actually
//     returned material, and even then it is still reported — as a partial read,
//     which is what it is.
//
// The old code counted empty as answered for both, so a source with one empty
// namespace and one 503 warned about nothing at all while contributing no
// evidence whatsoever. It also kept whichever failed status sorted first, which
// is alphabetical by namespace, so "the query failed" could be replaced in the
// banner by "never indexed".
func EvidenceWarnings(statuses []EvidenceStatus) []string {
	answered := map[string]bool{}
	returnedMaterial := map[string]bool{}
	for _, status := range statuses {
		switch status.State {
		case EvidenceOK:
			answered[status.Source] = true
			returnedMaterial[status.Source] = true
		case EvidenceEmpty:
			answered[status.Source] = true
		}
	}
	bySource := map[string][]EvidenceStatus{}
	var order []string
	for _, status := range statuses {
		if !status.Degraded() {
			continue
		}
		if status.State == EvidenceMissing && answered[status.Source] {
			continue
		}
		if _, ok := bySource[status.Source]; !ok {
			order = append(order, status.Source)
		}
		bySource[status.Source] = append(bySource[status.Source], status)
	}
	var out []string
	for _, source := range order {
		failed := bySource[source]
		// Severity order, not namespace order: a source that both failed to be
		// read somewhere and was absent somewhere else leads with the failure.
		sort.SliceStable(failed, func(i, j int) bool {
			return evidenceFailureRank(failed[i].State) > evidenceFailureRank(failed[j].State)
		})
		out = append(out, evidenceSourceWarning(source, failed, returnedMaterial[source]))
	}
	return out
}

func evidenceFailureRank(state string) int {
	if state == EvidenceUnavailable {
		return 1
	}
	return 0
}

// evidenceSourceWarning states what happened to one source without pretending
// only one namespace was tried. The single-namespace phrasing is kept for the
// single-namespace case, which is most of them.
func evidenceSourceWarning(source string, failed []EvidenceStatus, partial bool) string {
	// Nothing was there, and nothing broke: say that once. Which candidate
	// names were probed on the way to finding out is a fact about how lookup
	// works, not about the reader's repository, and listing three internal
	// namespace ids buries the one sentence that matters — that the repository
	// is not indexed, and where to fix it. The names stay in the structured
	// evidence for anyone debugging.
	if !partial && allMissing(failed) {
		return withRemedy(source+": "+evidenceStatusDetail(failed[0]), evidenceRemedies(failed))
	}
	if len(failed) == 1 && !partial {
		return withRemedy(failed[0].warning(), evidenceRemedies(failed))
	}
	details := make([]string, 0, len(failed))
	for _, status := range failed {
		namespace := strings.TrimSpace(status.Namespace)
		if namespace == "" {
			namespace = status.State
		}
		details = append(details, namespace+": "+evidenceStatusDetail(status))
	}
	if partial {
		return withRemedy(fmt.Sprintf("%s: read in part — %d namespace(s) could not be read (%s)",
			source, len(failed), strings.Join(details, "; ")), evidenceRemedies(failed))
	}
	return withRemedy(fmt.Sprintf("%s: none of %d namespace(s) could be read (%s)",
		source, len(failed), strings.Join(details, "; ")), evidenceRemedies(failed))
}

// evidenceRemedies collects the distinct fixes for one source, in the order
// they were recorded. Distinct because the same source probed across three
// namespaces usually has one fix, and a warning that says it three times reads
// like three different problems.
func evidenceRemedies(failed []EvidenceStatus) []string {
	seen := map[string]bool{}
	var out []string
	for _, status := range failed {
		remedy := strings.TrimSpace(status.Remedy)
		if remedy == "" || seen[remedy] {
			continue
		}
		seen[remedy] = true
		out = append(out, remedy)
	}
	return out
}

// allMissing reports whether every failure is simply an absent namespace, as
// opposed to one that could not be read. The two need different warnings: "not
// indexed" is the reader's to fix, "the query failed" is not.
func allMissing(failed []EvidenceStatus) bool {
	for _, status := range failed {
		if status.State != EvidenceMissing {
			return false
		}
	}
	return len(failed) > 0
}

// withRemedy appends the instructions as their own sentences, because that is
// what they are.
func withRemedy(warning string, remedies []string) string {
	if len(remedies) == 0 {
		return warning
	}
	return strings.TrimRight(warning, ". ") + ". " + strings.Join(remedies, ". ")
}

func (s EvidenceStatus) warning() string {
	label := s.Source
	if namespace := strings.TrimSpace(s.Namespace); namespace != "" {
		label += " (" + namespace + ")"
	}
	return label + ": " + evidenceStatusDetail(s)
}

func evidenceStatusDetail(s EvidenceStatus) string {
	if detail := strings.TrimSpace(s.Detail); detail != "" {
		return detail
	}
	switch s.State {
	case EvidenceMissing:
		return "not indexed"
	case EvidenceUnavailable:
		return "query failed"
	default:
		return s.State
	}
}

// EvidenceSummaryLine is the one-line evidence listing for the verbose report.
func EvidenceSummaryLine(status EvidenceStatus) string {
	label := status.Source
	if namespace := strings.TrimSpace(status.Namespace); namespace != "" {
		label += " (" + namespace + ")"
	}
	line := fmt.Sprintf("%s: %s", label, status.State)
	if status.Snippets > 0 {
		line += fmt.Sprintf(", %d snippet(s)", status.Snippets)
	}
	if detail := strings.TrimSpace(status.Detail); detail != "" {
		line += " — " + detail
	}
	return line
}

// evidenceNamer lets a retriever name the evidence source it stands for, so a
// bare error from it is reported as "code index" and not as a Go type name.
type evidenceNamer interface {
	EvidenceSource() string
}

func evidenceSourceName(retriever ContextRetriever) string {
	if named, ok := retriever.(evidenceNamer); ok {
		if name := strings.TrimSpace(named.EvidenceSource()); name != "" {
			return name
		}
	}
	return fmt.Sprintf("%T", retriever)
}
