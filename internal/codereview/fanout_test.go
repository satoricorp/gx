package codereview

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func diffSnippetsOfSize(count, size int) []DiffSnippet {
	out := make([]DiffSnippet, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, DiffSnippet{File: fmt.Sprintf("internal/pkg%03d/file.go", i), Diff: strings.Repeat("+x\n", size/3)})
	}
	return out
}

// The reason fan-out exists: a change larger than one call used to be reviewed
// as its first ten alphabetical files. Every changed file has to reach some
// shard, or the review is reporting on a subset and calling it the change.
func TestEveryChangedFileReachesAShard(t *testing.T) {
	snippets := diffSnippetsOfSize(120, 9000)
	files := make([]string, 0, len(snippets))
	for _, snippet := range snippets {
		files = append(files, snippet.File)
	}
	brief := ReviewBrief{
		Static:  StaticSnapshot{ChangedFiles: files, DiffSnippets: snippets},
		Context: []ContextSnippet{{Kind: "repo_doc", Ref: "AGENTS.md", Text: "conventions"}},
	}

	shards, coverage := planReviewShards(brief, Options{})
	if len(shards) < 2 {
		t.Fatalf("shards = %d; a %d-byte diff does not fit one call and must fan out", len(shards), snippetBytes(snippets))
	}
	seen := map[string]int{}
	for _, shard := range shards {
		if len(shard.Brief.Context) == 0 {
			t.Fatalf("shard %s carries no context; a shard reviewed without the project's conventions is a fragment of a review", shard.Label())
		}
		for _, snippet := range shard.Brief.Static.DiffSnippets {
			seen[snippet.File]++
		}
	}
	for _, file := range files {
		if seen[file] != 1 {
			t.Fatalf("%s appears in %d shards, want exactly 1", file, seen[file])
		}
	}
	if coverage.Read != len(files) || coverage.Total != len(files) {
		t.Fatalf("coverage = %d/%d, want %d/%d", coverage.Read, coverage.Total, len(files), len(files))
	}
	if !coverage.Complete() {
		t.Fatalf("coverage reports incomplete after reading every file: %+v", coverage)
	}
}

// Fanning out a subject that fits would spend calls to make the review worse.
// It also has to stay true that the ordinary case is one call, because every
// existing behavior is pinned against a single brief.
func TestASubjectThatFitsStaysOneCall(t *testing.T) {
	brief := ReviewBrief{
		Static: StaticSnapshot{
			ChangedFiles: []string{"a.go", "b.go"},
			DiffSnippets: []DiffSnippet{{File: "a.go", Diff: "+a\n"}, {File: "b.go", Diff: "+b\n"}},
		},
		Context: []ContextSnippet{{Kind: "repo_doc", Ref: "AGENTS.md", Text: "conventions"}},
	}
	shards, coverage := planReviewShards(brief, Options{})
	if len(shards) != 1 {
		t.Fatalf("shards = %d, want 1", len(shards))
	}
	if len(shards[0].Brief.Static.DiffSnippets) != 2 || len(shards[0].Brief.Context) != 1 {
		t.Fatalf("single shard is not the whole brief: %+v", shards[0].Brief)
	}
	if !coverage.Complete() {
		t.Fatalf("coverage = %+v, want complete", coverage)
	}
}

// A review with nothing to diff is still a review — docs-only changes and
// prompt-directed questions arrive here — and it must still reach the model.
func TestAnEmptySubjectStillProducesOneCall(t *testing.T) {
	shards, _ := planReviewShards(ReviewBrief{Context: []ContextSnippet{{Kind: "domain_doc", Ref: "CONTEXT.md", Text: "domain"}}}, Options{})
	if len(shards) != 1 {
		t.Fatalf("shards = %d, want 1", len(shards))
	}
}

type countingShardReviewer struct {
	mu    sync.Mutex
	calls int
	fail  map[int]bool
}

func (c *countingShardReviewer) Review(_ context.Context, brief ReviewBrief) ([]Finding, error) {
	c.mu.Lock()
	index := c.calls
	c.calls++
	c.mu.Unlock()
	if c.fail[index] {
		return nil, fmt.Errorf("shard %d failed", index)
	}
	return []Finding{{
		ID:             "ai.review.1",
		Title:          fmt.Sprintf("Finding from %d files", len(brief.Static.DiffSnippets)),
		Summary:        "summary",
		Benefit:        "benefit",
		Recommendation: "do it",
	}}, nil
}

// Every shard's reviewer numbers its findings ai.review.1, ai.review.2, … and
// mergeFindings deduplicates on that ID. Without namespacing, a fan-out would
// read the whole change and then report one shard's worth of it — the exact
// failure this work exists to remove, reintroduced one layer down.
func TestShardFindingsSurviveTheMerge(t *testing.T) {
	snippets := diffSnippetsOfSize(120, 9000)
	files := make([]string, 0, len(snippets))
	for _, snippet := range snippets {
		files = append(files, snippet.File)
	}
	brief := ReviewBrief{Static: StaticSnapshot{ChangedFiles: files, DiffSnippets: snippets}}
	shards, coverage := planReviewShards(brief, Options{})
	reviewer := &countingShardReviewer{}
	result := runShardedReview(context.Background(), reviewer, shards, coverage, Options{})
	if result.Err != nil {
		t.Fatalf("runShardedReview error = %v", result.Err)
	}
	if len(result.Findings) != len(shards) {
		t.Fatalf("findings = %d across %d shards; identical IDs were merged away", len(result.Findings), len(shards))
	}
	ids := map[string]struct{}{}
	for _, finding := range result.Findings {
		ids[finding.ID] = struct{}{}
	}
	if len(ids) != len(result.Findings) {
		t.Fatalf("finding IDs collide across shards: %#v", result.Findings)
	}
}

// A shard that fails leaves a region of the change unreviewed. Reporting that
// as a clean region is the same false pass as reviewing nothing and saying "no
// material issues".
func TestAFailedShardIsAHoleInCoverageNotACleanRegion(t *testing.T) {
	snippets := diffSnippetsOfSize(120, 9000)
	files := make([]string, 0, len(snippets))
	for _, snippet := range snippets {
		files = append(files, snippet.File)
	}
	brief := ReviewBrief{Static: StaticSnapshot{ChangedFiles: files, DiffSnippets: snippets}}
	shards, coverage := planReviewShards(brief, Options{})
	reviewer := &countingShardReviewer{fail: map[int]bool{0: true}}
	result := runShardedReview(context.Background(), reviewer, shards, coverage, Options{})
	if result.Err != nil {
		t.Fatalf("one failed shard must not fail the whole review: %v", result.Err)
	}
	if result.Coverage.ShardsFailed != 1 {
		t.Fatalf("ShardsFailed = %d, want 1", result.Coverage.ShardsFailed)
	}
	if result.Coverage.Read >= result.Coverage.Total {
		t.Fatalf("coverage still claims every file was read (%d/%d) after a shard failed", result.Coverage.Read, result.Coverage.Total)
	}
	if result.Coverage.Complete() {
		t.Fatalf("coverage reports complete with a failed shard: %+v", result.Coverage)
	}
	if !strings.Contains(result.Coverage.Statement(), "failed") {
		t.Fatalf("coverage statement hides the failed shard: %q", result.Coverage.Statement())
	}
}

// When every shard fails the review failed; that is a degraded review, not a
// partial one, and the caller has to be able to tell.
func TestAllShardsFailingIsAFailedReview(t *testing.T) {
	brief := ReviewBrief{Static: StaticSnapshot{ChangedFiles: []string{"a.go"}, DiffSnippets: []DiffSnippet{{File: "a.go", Diff: "+a\n"}}}}
	shards, coverage := planReviewShards(brief, Options{})
	reviewer := &countingShardReviewer{fail: map[int]bool{0: true}}
	result := runShardedReview(context.Background(), reviewer, shards, coverage, Options{})
	if result.Err == nil {
		t.Fatalf("want an error when no shard answered")
	}
	if result.Coverage.Read != 0 {
		t.Fatalf("Read = %d after every shard failed, want 0", result.Coverage.Read)
	}
}

// Context is shared, not partitioned: a shard reviewed without the project's
// policy and its retrieved evidence would produce a fraction of a review N
// times over.
func TestSharedContextRidesInEveryShard(t *testing.T) {
	snippets := diffSnippetsOfSize(120, 9000)
	files := make([]string, 0, len(snippets))
	for _, snippet := range snippets {
		files = append(files, snippet.File)
	}
	brief := ReviewBrief{
		Static: StaticSnapshot{ChangedFiles: files, DiffSnippets: snippets},
		Context: []ContextSnippet{
			{Kind: "review_policy", Ref: "REVIEW.md", Text: "policy"},
			{Kind: "indexed_code", Ref: "internal/x/x.go", Text: "retrieved code"},
		},
	}
	shards, _ := planReviewShards(brief, Options{})
	if len(shards) < 2 {
		t.Fatalf("shards = %d, want a fan-out", len(shards))
	}
	for _, shard := range shards {
		kinds := map[string]bool{}
		for _, snippet := range shard.Brief.Context {
			kinds[snippet.Kind] = true
		}
		if !kinds["review_policy"] || !kinds["indexed_code"] {
			t.Fatalf("shard %s lost shared context: %#v", shard.Label(), shard.Brief.Context)
		}
	}
}
