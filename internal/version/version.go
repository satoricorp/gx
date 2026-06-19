package version

import (
	"runtime/debug"
	"strings"
)

var Version = "dev"

var readBuildInfo = func() (*debug.BuildInfo, bool) {
	return debug.ReadBuildInfo()
}

// Current returns the version string for display and wire metadata.
// Prefer the short VCS revision for the checked-out HEAD whenever Go embedded
// VCS metadata is available. Version remains a fallback for builds without VCS
// metadata.
func Current() string {
	if revision := vcsRevision(); revision != "" {
		return revision
	}
	return strings.TrimSpace(Version)
}

func vcsRevision() string {
	info, ok := readBuildInfo()
	if !ok || info == nil {
		return ""
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return shortRevision(setting.Value)
		}
	}
	return ""
}

func shortRevision(revision string) string {
	revision = strings.TrimSpace(revision)
	if len(revision) <= 8 {
		return revision
	}
	return revision[:8]
}
