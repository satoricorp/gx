package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"
)

// `gx update` does its own fetching, and version is the one command that must
// answer about the binary in hand.
func TestSelfUpdateSkipsItsOwnCommands(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("CI", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("update check ran for a skipped command: %s", r.URL.Path)
	}))
	defer server.Close()
	t.Setenv("GX_DOWNLOAD_BASE_URL", server.URL)

	for _, name := range []string{"update", "version", "__complete"} {
		maybeSelfUpdate(context.Background(), &cobra.Command{Use: name})
	}
}
