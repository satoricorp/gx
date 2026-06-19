package version

import (
	"runtime/debug"
	"testing"
)

func TestCurrentPrefersShortVCSRevisionOverExplicitVersion(t *testing.T) {
	restore := stubVersion(t, "0.4.2", &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: "2f59fcaa1234567890"},
	}}, true)
	defer restore()

	if got := Current(); got != "2f59fcaa" {
		t.Fatalf("Current() = %q, want short revision", got)
	}
}

func TestCurrentFallsBackToShortVCSRevision(t *testing.T) {
	restore := stubVersion(t, "dev", &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: "2f59fcaa1234567890"},
	}}, true)
	defer restore()

	if got := Current(); got != "2f59fcaa" {
		t.Fatalf("Current() = %q, want short revision", got)
	}
}

func TestCurrentIgnoresModifiedVCSStatus(t *testing.T) {
	restore := stubVersion(t, "dev", &debug.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "vcs.revision", Value: "2f59fcaa1234567890"},
		{Key: "vcs.modified", Value: "true"},
	}}, true)
	defer restore()

	if got := Current(); got != "2f59fcaa" {
		t.Fatalf("Current() = %q, want short revision without dirty suffix", got)
	}
}

func TestCurrentKeepsDevWithoutVCSRevision(t *testing.T) {
	restore := stubVersion(t, "dev", &debug.BuildInfo{}, true)
	defer restore()

	if got := Current(); got != "dev" {
		t.Fatalf("Current() = %q, want dev fallback", got)
	}
}

func stubVersion(t *testing.T, version string, info *debug.BuildInfo, ok bool) func() {
	t.Helper()
	prevVersion := Version
	prevReadBuildInfo := readBuildInfo
	Version = version
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return info, ok
	}
	return func() {
		Version = prevVersion
		readBuildInfo = prevReadBuildInfo
	}
}
