package cli

import (
	"context"
	"strings"
	"testing"
)

func TestIndexCommandIsRegisteredAsHidden(t *testing.T) {
	root := NewRoot(context.Background())
	cmd, _, err := root.Find([]string{"index"})
	if err != nil || cmd == nil || cmd.Name() != "index" {
		t.Fatalf("Find(index) = cmd=%v err=%v, want the index command", cmd, err)
	}
	if !cmd.Hidden {
		t.Fatal("gx index is a maintenance command and must stay hidden")
	}
	if cmd.GroupID != "" {
		t.Fatalf("gx index GroupID = %q, want no group", cmd.GroupID)
	}
	for _, flag := range []string{"full", "json", "quiet", "concurrency", "namespace"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Fatalf("gx index is missing the --%s flag", flag)
		}
	}
}

func TestIndexCommandReportsMissingCredentials(t *testing.T) {
	// Indexing must fail loudly when it cannot run. The historical failure mode
	// was the opposite: a missing key nulled the indexer and every publish
	// reported success while writing nothing.
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GX_OPENAI_API_KEY", "")
	t.Setenv("TURBOPUFFER_API_KEY", "")
	t.Setenv("GX_HOME", t.TempDir())

	root := NewRoot(context.Background())
	root.SetArgs([]string{"index", "--quiet"})
	root.SetOut(&strings.Builder{})
	root.SetErr(&strings.Builder{})
	err := root.Execute()
	if err == nil {
		t.Fatal("gx index without credentials must return an error, not silently do nothing")
	}
	if !strings.Contains(err.Error(), "OPENAI_API_KEY") {
		t.Fatalf("error = %q, want it to name the missing credential", err)
	}
}
