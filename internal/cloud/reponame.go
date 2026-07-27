package cloud

import "strings"

// RepoFullNameFromRemoteURL extracts "owner/repo" from a GitHub remote URL.
func RepoFullNameFromRemoteURL(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if remoteURL == "" {
		return ""
	}
	if match := strings.TrimPrefix(remoteURL, "git@github.com:"); match != remoteURL {
		parts := strings.Split(strings.TrimSuffix(match, ".git"), "/")
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}
	if strings.Contains(remoteURL, "github.com") {
		idx := strings.Index(remoteURL, "github.com/")
		if idx >= 0 {
			path := strings.TrimPrefix(remoteURL[idx+len("github.com/"):], "/")
			path = strings.TrimSuffix(path, ".git")
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				return parts[0] + "/" + parts[1]
			}
		}
	}
	return ""
}
