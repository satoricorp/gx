package update

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// The payload names inside the published archive. Only these are extracted:
// matching entries by exact name means a tampered archive cannot write outside
// the destination no matter what paths it carries, without a traversal check
// to get subtly wrong.
const (
	archiveCLI          = "gx/bin/gx"
	archiveMCP          = "gx/bin/gx-mcp"
	archiveBashComplete = "gx/completions/gx.bash"
	archiveZshComplete  = "gx/completions/_gx"
)

// maxPayloadBytes bounds one extracted file. The binaries are ~30 MB; this
// leaves room without letting a decompression bomb run the disk out.
const maxPayloadBytes = 512 << 20

// AssetKey is the manifest key for the running platform.
func AssetKey() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

// InstallDir is where the running binary lives, which is where its replacement
// goes. Derived from the executable rather than assuming ~/.local/bin: the
// installer honours GX_INSTALL_DIR, and an update must land on the binary the
// caller actually ran, not on a different copy somewhere else in PATH.
func InstallDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate running binary: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		resolved = exe
	}
	return filepath.Dir(resolved), nil
}

// Apply downloads the manifest's asset for this platform, verifies it against
// the published digest, and replaces the installed files.
//
// progress, when non-nil, receives human-readable step lines.
func Apply(ctx context.Context, manifest Manifest, client *http.Client, progress func(string)) error {
	dir, err := InstallDir()
	if err != nil {
		return err
	}
	return ApplyTo(ctx, manifest, client, dir, progress)
}

// ApplyTo is Apply against an explicit directory. Separate so a test can
// install into a temp directory instead of over the test binary.
func ApplyTo(
	ctx context.Context,
	manifest Manifest,
	client *http.Client,
	dir string,
	progress func(string),
) error {
	say := func(format string, args ...any) {
		if progress != nil {
			progress(fmt.Sprintf(format, args...))
		}
	}

	key := AssetKey()
	asset, ok := manifest.Assets[key]
	if !ok {
		return fmt.Errorf("no published build for %s", key)
	}
	if strings.TrimSpace(asset.URL) == "" {
		return fmt.Errorf("published build for %s has no url", key)
	}

	// Fail before downloading 30 MB rather than after. A system-wide install
	// (/usr/local/bin) or a read-only mount is the caller's to sort out, and
	// the installer's own directory is the remedy either way.
	if err := writable(dir); err != nil {
		return fmt.Errorf("cannot write to %s: %w; re-run the installer instead: %s", dir, err, installCommand(manifest))
	}

	say("Downloading %s", asset.URL)
	archive, err := download(ctx, client, asset.URL)
	if err != nil {
		return err
	}
	defer os.Remove(archive)

	say("Verifying checksum")
	if err := verify(archive, asset.SHA256); err != nil {
		return err
	}

	say("Installing to %s", dir)
	return extractAndInstall(archive, dir, say)
}

func installCommand(manifest Manifest) string {
	if cmd := strings.TrimSpace(manifest.InstallCommand); cmd != "" {
		return cmd
	}
	return "curl -fsSL " + BaseURL() + "/install.sh | sh"
}

// writable proves the directory accepts a new file, which is what the update
// actually needs — os.Access-style permission bits do not account for an
// immutable mount or a directory owned by root with a group-writable bit.
func writable(dir string) error {
	probe, err := os.CreateTemp(dir, ".gx-update-*")
	if err != nil {
		return err
	}
	name := probe.Name()
	probe.Close()
	return os.Remove(name)
}

func download(ctx context.Context, client *http.Client, url string) (string, error) {
	if client == nil {
		// Deliberately not fetchTimeout: that bounds a metadata request, and
		// this is tens of megabytes over whatever connection the user has.
		client = &http.Client{}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create download request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("download %s: status %s", url, resp.Status)
	}
	tmp, err := os.CreateTemp("", "gx-update-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	if _, err := io.Copy(tmp, io.LimitReader(resp.Body, maxPayloadBytes)); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", fmt.Errorf("write download: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("close download: %w", err)
	}
	return tmp.Name(), nil
}

// verify checks the archive against the digest published in the manifest,
// before anything is unpacked. The manifest and the archive come from the same
// host, so this is not a defence against that host — it is a defence against a
// truncated or corrupted transfer being installed as a working binary.
func verify(path, expected string) error {
	expected = strings.ToLower(strings.TrimSpace(expected))
	if expected == "" {
		return fmt.Errorf("manifest has no checksum for this build")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open download: %w", err)
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return fmt.Errorf("read download: %w", err)
	}
	actual := hex.EncodeToString(digest.Sum(nil))
	if actual != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

func extractAndInstall(archivePath, dir string, say func(string, ...any)) error {
	payload, err := extract(archivePath)
	if err != nil {
		return err
	}
	if payload[archiveCLI] == nil {
		return fmt.Errorf("archive is missing %s", archiveCLI)
	}
	if payload[archiveMCP] == nil {
		return fmt.Errorf("archive is missing %s", archiveMCP)
	}

	if err := replace(filepath.Join(dir, "gx"), payload[archiveCLI], 0o755); err != nil {
		return err
	}
	if err := replace(filepath.Join(dir, "gx-mcp"), payload[archiveMCP], 0o755); err != nil {
		return err
	}
	// Mirrors install.sh: gxr is a symlink to gx, and gxc is removed rather
	// than left resolving to a gx with no subcommand.
	link := filepath.Join(dir, "gxr")
	_ = os.Remove(link)
	if err := os.Symlink("gx", link); err != nil {
		say("Note: could not create the gxr symlink: %v", err)
	}
	_ = os.Remove(filepath.Join(dir, "gxc"))

	installCompletions(payload, say)
	return nil
}

func installCompletions(payload map[string][]byte, say func(string, ...any)) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	if body := payload[archiveBashComplete]; body != nil {
		target := filepath.Join(home, ".local", "share", "bash-completion", "completions", "gx")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err == nil {
			if err := replace(target, body, 0o644); err != nil {
				say("Note: could not update bash completions: %v", err)
			}
		}
	}
	if body := payload[archiveZshComplete]; body != nil {
		target := filepath.Join(home, ".zfunc", "_gx")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err == nil {
			if err := replace(target, body, 0o644); err != nil {
				say("Note: could not update zsh completions: %v", err)
			}
		}
	}
}

func extract(archivePath string) (map[string][]byte, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("read archive: %w", err)
	}
	defer gz.Close()

	wanted := map[string]bool{
		archiveCLI:          true,
		archiveMCP:          true,
		archiveBashComplete: true,
		archiveZshComplete:  true,
	}
	payload := map[string][]byte{}
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read archive entry: %w", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.ToSlash(filepath.Clean(header.Name))
		if !wanted[name] {
			continue
		}
		body, err := io.ReadAll(io.LimitReader(reader, maxPayloadBytes))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		payload[name] = body
	}
	return payload, nil
}

// replace writes body to path via a sibling temp file and a rename.
//
// The rename is what makes this work on the binary that is currently running:
// a direct write to it fails with ETXTBSY on Linux, while renaming over it
// leaves the running process on the old inode and every later invocation on
// the new one. The temp file is a sibling so the rename stays within one
// filesystem, where it is atomic.
func replace(path string, body []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".new-*")
	if err != nil {
		return fmt.Errorf("create temp file in %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(body); err != nil {
		tmp.Close()
		return fmt.Errorf("write %s: %w", tmpName, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", tmpName, err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("chmod %s: %w", tmpName, err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}
