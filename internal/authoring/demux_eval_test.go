package authoring

import (
	"fmt"
	"testing"
)

// evalDemuxRoutesForFiles runs the deterministic grouping pipeline the way
// ProposeDemux does (file grouping, then new-stack clustering) and returns the
// proposed target stack per file.
func evalDemuxRoutesForFiles(t *testing.T, files []string) map[string]string {
	t.Helper()
	groups := expandGroupsByChangedSymbols(groupFiles(files), map[string][]HunkRange{})
	revisions := make([]RevisionProposal, 0, len(groups))
	for index, group := range groups {
		revisions = append(revisions, RevisionProposal{
			ID:     fmt.Sprintf("u%d", index+1),
			Intent: proposalIntentForGroup("", group),
			Files:  group.Files,
		})
	}
	proposal := inferNewStackRoutes(DemuxProposal{Revisions: revisions})
	routes := map[string]string{}
	for _, revision := range proposal.Revisions {
		for _, file := range revision.Files {
			routes[file] = revision.TargetStack
		}
	}
	return routes
}

// TestDemuxEvalMultiConcernGrouping is the grouping-accuracy eval fixture: a
// representative working copy with feature code, its tests, container/CI
// infra, project tooling, and one unrelated single-file change. The planner
// must not produce one stack per file.
func TestDemuxEvalMultiConcernGrouping(t *testing.T) {
	routes := evalDemuxRoutesForFiles(t, []string{
		"server/api/billing.ts",
		"server/api/invoice.ts",
		"tests/billing.test.ts",
		"tests/invoice.test.ts",
		"Dockerfile",
		"docker-compose.yml",
		".github/workflows/ci.yaml",
		"tsconfig.json",
		".gitignore",
		"LICENSE",
		"cmd/worker/main.go",
	})

	sameStack := func(left, right string) {
		t.Helper()
		if routes[left] == "" || routes[left] != routes[right] {
			t.Fatalf("files %s (%q) and %s (%q) should share a stack; routes=%#v", left, routes[left], right, routes[right], routes)
		}
	}
	// Tests group with the code they test.
	sameStack("server/api/billing.ts", "tests/billing.test.ts")
	sameStack("server/api/billing.ts", "tests/invoice.test.ts")
	sameStack("server/api/billing.ts", "server/api/invoice.ts")
	// Infra/CI files group into one infra stack.
	sameStack("Dockerfile", "docker-compose.yml")
	sameStack("Dockerfile", ".github/workflows/ci.yaml")
	if routes["Dockerfile"] != "chore/infra" {
		t.Fatalf("Dockerfile route = %q, want chore/infra; routes=%#v", routes["Dockerfile"], routes)
	}
	// Root tooling and meta files group into one tooling stack.
	sameStack("tsconfig.json", ".gitignore")
	sameStack("tsconfig.json", "LICENSE")
	// Unrelated subsystem keeps its own stack (strong isolation signal).
	if routes["cmd/worker/main.go"] == routes["server/api/billing.ts"] || routes["cmd/worker/main.go"] == routes["Dockerfile"] {
		t.Fatalf("cmd/worker/main.go should stay independent; routes=%#v", routes)
	}
	stacks := map[string]struct{}{}
	for _, stack := range routes {
		stacks[stack] = struct{}{}
	}
	if len(stacks) != 4 {
		t.Fatalf("stack count = %d, want 4 (code+tests, infra, tooling, worker); routes=%#v", len(stacks), routes)
	}
}

func TestDemuxEvalTestOnlyChangesShareOneStack(t *testing.T) {
	routes := evalDemuxRoutesForFiles(t, []string{
		"tests/billing.test.ts",
		"tests/invoice.test.ts",
		"test/cli.test.ts",
	})
	want := "test/test-updates"
	for file, stack := range routes {
		if stack != want {
			t.Fatalf("test-only file %s routed to %q, want %q; routes=%#v", file, stack, want, routes)
		}
	}
}

func TestIsTestFileConventions(t *testing.T) {
	tests := []struct {
		file string
		want bool
	}{
		{"foo.test.ts", true},
		{"foo.spec.tsx", true},
		{"foo_test.go", true},
		{"test_foo.py", true},
		{"__tests__/foo.tsx", true},
		{"test/helpers.ts", true},
		{"src/api/main.ts", false},
		{"scripts/test-api-local.mjs", false},
		{"Dockerfile", false},
	}
	for _, tt := range tests {
		if got := isTestFile(tt.file); got != tt.want {
			t.Fatalf("isTestFile(%q) = %v, want %v", tt.file, got, tt.want)
		}
	}
}

func TestTestSubjectPairingConventions(t *testing.T) {
	tests := []struct {
		test    string
		subject string
	}{
		{"foo.test.ts", "src/foo.ts"},
		{"foo.spec.ts", "lib/foo.ts"},
		{"foo_test.go", "pkg/util/foo.go"},
		{"test_foo.py", "app/foo.py"},
		{"__tests__/foo.tsx", "src/components/foo.tsx"},
		{"test/sdk.test.ts", "src/sdk/index.ts"},
		{"test/http-client.test.ts", "src/http/client.ts"},
	}
	for _, tt := range tests {
		subjectKeys := map[string]struct{}{}
		for _, key := range codeFileSubjectKeys(tt.subject) {
			subjectKeys[key] = struct{}{}
		}
		matched := false
		for _, key := range testSubjectKeys(tt.test) {
			if _, ok := subjectKeys[key]; ok {
				matched = true
			}
		}
		if !matched {
			t.Fatalf("test %q did not match subject %q (test keys %v, subject keys %v)",
				tt.test, tt.subject, testSubjectKeys(tt.test), codeFileSubjectKeys(tt.subject))
		}
	}
}

func TestIsInfraFileConventions(t *testing.T) {
	tests := []struct {
		file string
		want bool
	}{
		{"Dockerfile", true},
		{"Dockerfile.dev", true},
		{".dockerignore", true},
		{"docker-compose.yml", true},
		{".github/workflows/deploy.yaml", true},
		{"main.tf", true},
		{"Makefile", true},
		{"justfile", true},
		{"cdk.json", true},
		{"cdk.context.json", true},
		{"infra/app.ts", true},
		{"deploy/rollout.sh", true},
		{"package.json", false},
		{"tsconfig.json", false},
		{"vite.config.ts", false},
		{"src/api/main.ts", false},
	}
	for _, tt := range tests {
		if got := isInfraFile(tt.file); got != tt.want {
			t.Fatalf("isInfraFile(%q) = %v, want %v", tt.file, got, tt.want)
		}
	}
}

func TestFragmentMergeJoinsSingleFileClusterIntoRelatedCluster(t *testing.T) {
	proposal := inferNewStackRoutes(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "update tools helper", Files: []string{"tools/a.sh"}},
			{ID: "u2", Intent: "update tools helper", Files: []string{"tools/b.sh"}},
			{ID: "u3", Intent: "update nested tool", Files: []string{"tools/sub/c.sh"}},
		},
	})
	routes := map[string]string{}
	for _, revision := range proposal.Revisions {
		routes[revision.ID] = revision.TargetStack
	}
	if routes["u3"] != routes["u1"] || routes["u1"] == "" {
		t.Fatalf("single-file nested cluster should merge into tools cluster; routes=%#v", routes)
	}
}

func TestFragmentMergeRespectsDistinctSessionIsolation(t *testing.T) {
	proposal := inferNewStackRoutes(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "update tools helper", Files: []string{"tools/a.sh"}, SessionIDs: []string{"session-a"}},
			{ID: "u2", Intent: "update tools helper", Files: []string{"tools/b.sh"}, SessionIDs: []string{"session-a"}},
			{ID: "u3", Intent: "update nested tool", Files: []string{"tools/sub/c.sh"}, SessionIDs: []string{"session-b"}},
		},
	})
	routes := map[string]string{}
	for _, revision := range proposal.Revisions {
		routes[revision.ID] = revision.TargetStack
	}
	if routes["u3"] == routes["u1"] {
		t.Fatalf("distinct-session single-file cluster should stay isolated; routes=%#v", routes)
	}
}

func TestFragmentMergeGroupsLeftoverStraysIntoSharedStack(t *testing.T) {
	proposal := inferNewStackRoutes(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "update changelog", Files: []string{"CHANGES.txt"}},
			{ID: "u2", Intent: "update version", Files: []string{"VERSION"}},
		},
	})
	routes := map[string]string{}
	for _, revision := range proposal.Revisions {
		routes[revision.ID] = revision.TargetStack
	}
	if routes["u1"] != routes["u2"] || routes["u1"] == "" {
		t.Fatalf("standalone root files should share one stack; routes=%#v", routes)
	}
}

func TestFragmentMergeNeverOverridesExplicitRoutes(t *testing.T) {
	proposal := inferNewStackRoutes(DemuxProposal{
		Revisions: []RevisionProposal{
			{ID: "u1", Intent: "update tools helper", Files: []string{"tools/a.sh"}},
			{ID: "u2", Intent: "user routed", Files: []string{"tools/sub/c.sh"}, TargetStack: "feature/custom", RouteSource: routeSourceUser},
		},
	})
	for _, revision := range proposal.Revisions {
		if revision.ID == "u2" {
			if revision.TargetStack != "feature/custom" || revision.RouteSource != routeSourceUser {
				t.Fatalf("explicit user route was overridden: %#v", revision)
			}
		}
	}
}
