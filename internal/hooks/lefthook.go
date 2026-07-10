package hooks

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteLefthookTemplate writes .lefthook.yml with gx pre-push capture.
func WriteLefthookTemplate(repoRoot, gxPath string) (string, error) {
	if gxPath == "" {
		var err error
		gxPath, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	path := filepath.Join(repoRoot, ".lefthook.yml")
	content := fmt.Sprintf(`# gx capture pre-push — see docs/hooks.md
pre-push:
  commands:
    gx-capture:
      run: |
        remote="$1"
        while read local_ref local_sha remote_ref remote_sha; do
          if [ "$local_sha" = "0000000000000000000000000000000000000000" ]; then
            continue
          fi
          if [ "$remote_sha" = "0000000000000000000000000000000000000000" ]; then
            range="$local_sha"
          else
            range="${remote_sha}..${local_sha}"
          fi
          %q capture push --remote "$remote" --ref-range "$range" --local-ref "$local_ref" --head-sha "$local_sha" || true
        done
`, gxPath)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// HuskyPrePushSnippet returns package.json husky pre-push script content.
func HuskyPrePushSnippet(gxPath string) string {
	if gxPath == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`remote="$1"
while read local_ref local_sha remote_ref remote_sha; do
  if [ "$local_sha" = "0000000000000000000000000000000000000000" ]; then
    continue
  fi
  if [ "$remote_sha" = "0000000000000000000000000000000000000000" ]; then
    range="$local_sha"
  else
    range="${remote_sha}..${local_sha}"
  fi
  %q capture push --remote "$remote" --ref-range "$range" --local-ref "$local_ref" --head-sha "$local_sha" || true
done
`, gxPath)
}
