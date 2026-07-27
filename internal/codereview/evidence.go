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
	defer l.mu.Unlock()
	l.statuses = append(l.statuses, status)
}

// count reports how many statuses have been filed so far. The composite uses it
// to tell a retriever that reported for itself from one that did not.
func (l *EvidenceLog) count() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.statuses)
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
func EvidenceWarnings(statuses []EvidenceStatus) []string {
	var out []string
	for _, status := range statuses {
		if !status.Degraded() {
			continue
		}
		out = append(out, status.warning())
	}
	return out
}

func (s EvidenceStatus) warning() string {
	label := s.Source
	if namespace := strings.TrimSpace(s.Namespace); namespace != "" {
		label += " (" + namespace + ")"
	}
	detail := strings.TrimSpace(s.Detail)
	if detail == "" {
		switch s.State {
		case EvidenceMissing:
			detail = "not indexed"
		case EvidenceUnavailable:
			detail = "query failed"
		default:
			detail = s.State
		}
	}
	return label + ": " + detail
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
