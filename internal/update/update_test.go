package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/version"
)

func TestBehindComparesCommitsNotVersionStrings(t *testing.T) {
	// The same commit publishes twice under two names — a 12-character sha
	// from a main push and a semver from its tag — and both are the same
	// binary. Comparing names would call that an update.
	sha := "87406ffcaaeb16bf73572f81a99de50269a212b7"
	if Behind(version.Info{Revision: sha}, Manifest{Version: "0.2.2", GitSHA: sha}) {
		t.Fatal("Behind() = true for the same commit under a different name")
	}
	if !Behind(version.Info{Revision: sha}, Manifest{Version: "0.2.3", GitSHA: strings.Repeat("a", 40)}) {
		t.Fatal("Behind() = false for a different commit")
	}
	// A binary with no VCS metadata cannot be compared, and guessing is worse
	// than staying quiet.
	if Behind(version.Info{}, Manifest{GitSHA: sha}) {
		t.Fatal("Behind() = true with no local revision")
	}
}

func TestSkipLeavesDevelopmentBuildsAlone(t *testing.T) {
	t.Setenv("CI", "")
	t.Setenv("GX_UPDATE_CHECK", "")

	released := version.Info{Release: "0.2.2", Revision: "abc"}
	if Skip(released) {
		t.Fatal("Skip() = true for a released build")
	}
	// `just build` bakes "dev", and its commit is whatever is checked out —
	// usually ahead of the published one, which a sha comparison reads as
	// behind.
	if !Skip(version.Info{Release: "dev", Revision: "abc"}) {
		t.Fatal("Skip() = false for a dev build")
	}
	if !Skip(version.Info{Release: "0.2.2", Revision: "abc", Modified: true}) {
		t.Fatal("Skip() = false for a dirty build")
	}

	t.Setenv("GX_UPDATE_CHECK", "0")
	if !Skip(released) {
		t.Fatal("Skip() = false with GX_UPDATE_CHECK=0")
	}
	t.Setenv("GX_UPDATE_CHECK", "")
	t.Setenv("CI", "true")
	if !Skip(released) {
		t.Fatal("Skip() = false under CI")
	}
}

func TestLatestServesFromCacheWithinTheInterval(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_ = json.NewEncoder(w).Encode(Manifest{Version: "0.2.2", GitSHA: strings.Repeat("b", 40)})
	}))
	defer server.Close()
	t.Setenv("GX_DOWNLOAD_BASE_URL", server.URL)

	now := time.Now()
	if _, err := Latest(context.Background(), server.Client(), now); err != nil {
		t.Fatalf("Latest() error = %v", err)
	}
	if _, err := Latest(context.Background(), server.Client(), now.Add(time.Hour)); err != nil {
		t.Fatalf("Latest() second call error = %v", err)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want the second call served from cache", requests)
	}

	// Past the interval it asks again, or a published fix would go unnoticed
	// on a machine that stays open for a week.
	if _, err := Latest(context.Background(), server.Client(), now.Add(checkInterval+time.Minute)); err != nil {
		t.Fatalf("Latest() after interval error = %v", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want a refetch after the interval", requests)
	}
}

func TestLatestSurvivesAnUnreadableCache(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	if err := os.WriteFile(filepath.Join(home, "update-check.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Manifest{Version: "0.2.2", GitSHA: strings.Repeat("c", 40)})
	}))
	defer server.Close()
	t.Setenv("GX_DOWNLOAD_BASE_URL", server.URL)

	manifest, err := Latest(context.Background(), server.Client(), time.Now())
	if err != nil {
		t.Fatalf("Latest() error = %v", err)
	}
	if manifest.Version != "0.2.2" {
		t.Fatalf("Latest() = %+v, want the fetched manifest", manifest)
	}
}

// buildArchive writes the layout scripts/package-cli.sh publishes.
func buildArchive(t *testing.T, cli, mcp string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	add := func(name, body string, mode int64) {
		if err := tw.WriteHeader(&tar.Header{
			Name: name, Mode: mode, Size: int64(len(body)), Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	add("gx/bin/gx", cli, 0o755)
	add("gx/bin/gx-mcp", mcp, 0o755)
	// Present in the real archive and deliberately not installed into the
	// binary directory.
	add("gx/README.md", "not installed", 0o644)
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func serveArchive(t *testing.T, archive []byte) (*httptest.Server, Manifest) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(archive)
	}))
	digest := sha256.Sum256(archive)
	manifest := Manifest{
		Version: "0.2.3",
		GitSHA:  strings.Repeat("d", 40),
		Assets: map[string]Asset{
			AssetKey(): {URL: server.URL + "/cli/latest/gx.tar.gz", SHA256: hex.EncodeToString(digest[:])},
		},
	}
	return server, manifest
}

func TestApplyReplacesTheInstalledBinaries(t *testing.T) {
	dir := t.TempDir()
	// Stand-ins for the installed files, so the test proves replacement rather
	// than first-time installation.
	if err := os.WriteFile(filepath.Join(dir, "gx"), []byte("old cli"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "gx-mcp"), []byte("old mcp"), 0o755); err != nil {
		t.Fatal(err)
	}
	// gxc is removed by the installer; an update that leaves it behind leaves
	// a command that opens the root help and looks broken.
	if err := os.WriteFile(filepath.Join(dir, "gxc"), []byte("stale"), 0o755); err != nil {
		t.Fatal(err)
	}

	server, manifest := serveArchive(t, buildArchive(t, "new cli", "new mcp"))
	defer server.Close()
	t.Setenv("HOME", t.TempDir())

	if err := ApplyTo(context.Background(), manifest, server.Client(), dir, nil); err != nil {
		t.Fatalf("ApplyTo() error = %v", err)
	}

	cli, err := os.ReadFile(filepath.Join(dir, "gx"))
	if err != nil || string(cli) != "new cli" {
		t.Fatalf("gx = %q (err %v), want the downloaded binary", cli, err)
	}
	mcp, err := os.ReadFile(filepath.Join(dir, "gx-mcp"))
	if err != nil || string(mcp) != "new mcp" {
		t.Fatalf("gx-mcp = %q (err %v), want the downloaded binary", mcp, err)
	}
	info, err := os.Stat(filepath.Join(dir, "gx"))
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("gx mode = %v (err %v), want 0755", info.Mode().Perm(), err)
	}
	if _, err := os.Lstat(filepath.Join(dir, "gxr")); err != nil {
		t.Fatalf("gxr symlink missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "gxc")); !os.IsNotExist(err) {
		t.Fatal("gxc still present after update")
	}
	if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
		t.Fatal("archive entries outside gx/bin were installed")
	}
}

func TestApplyRefusesAChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gx"), []byte("old cli"), 0o755); err != nil {
		t.Fatal(err)
	}

	server, manifest := serveArchive(t, buildArchive(t, "new cli", "new mcp"))
	defer server.Close()
	asset := manifest.Assets[AssetKey()]
	asset.SHA256 = strings.Repeat("0", 64)
	manifest.Assets[AssetKey()] = asset

	err := ApplyTo(context.Background(), manifest, server.Client(), dir, nil)
	if err == nil {
		t.Fatal("ApplyTo() error = nil, want a checksum failure")
	}
	if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("error = %v, want a checksum mismatch", err)
	}
	// A truncated download must not have half-replaced the install.
	cli, readErr := os.ReadFile(filepath.Join(dir, "gx"))
	if readErr != nil || string(cli) != "old cli" {
		t.Fatalf("gx = %q (err %v), want the original left in place", cli, readErr)
	}
}

func TestApplyReportsAnUnbuiltPlatform(t *testing.T) {
	server, manifest := serveArchive(t, buildArchive(t, "new cli", "new mcp"))
	defer server.Close()
	manifest.Assets = map[string]Asset{"plan9/386": {URL: "http://example.invalid", SHA256: "x"}}

	err := ApplyTo(context.Background(), manifest, server.Client(), t.TempDir(), nil)
	if err == nil || !strings.Contains(err.Error(), "no published build") {
		t.Fatalf("ApplyTo() error = %v, want an unbuilt-platform error", err)
	}
}

func TestAutoEnabledDefaultsOnAndRespectsTheOptOut(t *testing.T) {
	t.Setenv("GX_AUTO_UPDATE", "")
	if !AutoEnabled() {
		t.Fatal("AutoEnabled() = false by default, want on")
	}
	for _, value := range []string{"0", "false", "no", "off", "OFF"} {
		t.Setenv("GX_AUTO_UPDATE", value)
		if AutoEnabled() {
			t.Fatalf("AutoEnabled() = true with GX_AUTO_UPDATE=%q", value)
		}
	}
}

func TestAutoUpdateAttemptsAreRateLimited(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	now := time.Now()
	if !AutoUpdateDue(now) {
		t.Fatal("AutoUpdateDue() = false before any attempt")
	}

	// A failed update leaves the binary behind, so without this the next
	// command finds it behind again and re-downloads tens of megabytes — on
	// every command, until whatever broke gets fixed.
	MarkAutoUpdateAttempt(now)
	if AutoUpdateDue(now.Add(time.Hour)) {
		t.Fatal("AutoUpdateDue() = true an hour after an attempt")
	}
	if !AutoUpdateDue(now.Add(checkInterval + time.Minute)) {
		t.Fatal("AutoUpdateDue() = false after the interval")
	}
}

func TestCacheWritersDoNotClobberEachOther(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Manifest{Version: "0.2.3", GitSHA: strings.Repeat("e", 40)})
	}))
	defer server.Close()
	t.Setenv("GX_DOWNLOAD_BASE_URL", server.URL)

	now := time.Now()
	// The check and the attempt timestamp live in one file and are written by
	// two different paths. Either clobbering the other reintroduces exactly
	// what it was added to prevent.
	if _, err := Latest(context.Background(), server.Client(), now); err != nil {
		t.Fatalf("Latest() error = %v", err)
	}
	MarkAutoUpdateAttempt(now)

	cached, ok := readCache()
	if !ok {
		t.Fatal("cache missing after both writes")
	}
	if cached.Manifest.Version != "0.2.3" {
		t.Fatalf("manifest = %+v, want the fetched one kept through MarkAutoUpdateAttempt", cached.Manifest)
	}
	if cached.AttemptedAt.IsZero() {
		t.Fatal("attempt timestamp lost")
	}

	// And a later check must not wipe the attempt timestamp.
	if _, err := Latest(context.Background(), server.Client(), now.Add(checkInterval+time.Minute)); err != nil {
		t.Fatalf("Latest() refetch error = %v", err)
	}
	cached, _ = readCache()
	if cached.AttemptedAt.IsZero() {
		t.Fatal("attempt timestamp lost on refetch")
	}
}
