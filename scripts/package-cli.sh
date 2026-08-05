#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
dist_dir="${LGTM_CLI_DIST_DIR:-$repo_root/dist/cli}"
goos="${GOOS:-darwin}"
goarch="${GOARCH:-$(go env GOARCH)}"
version="${VERSION:-}"

cd "$repo_root"
env_github_client_id="${GITHUB_CLIENT_ID-}"
env_convex_site_url="${CONVEX_SITE_URL-}"
env_tl_cloud_url="${LGTM_CLOUD_URL-}"
env_tl_posthog_key="${LGTM_POSTHOG_KEY-}"
env_tl_posthog_host="${LGTM_POSTHOG_HOST-}"
set -a
if [[ -f .env ]]; then
  # shellcheck disable=SC1091
  source .env
fi
set +a
[[ -n "$env_github_client_id" ]] && export GITHUB_CLIENT_ID="$env_github_client_id"
[[ -n "$env_convex_site_url" ]] && export CONVEX_SITE_URL="$env_convex_site_url"
[[ -n "$env_tl_cloud_url" ]] && export LGTM_CLOUD_URL="$env_tl_cloud_url"
[[ -n "$env_tl_posthog_key" ]] && export LGTM_POSTHOG_KEY="$env_tl_posthog_key"
[[ -n "$env_tl_posthog_host" ]] && export LGTM_POSTHOG_HOST="$env_tl_posthog_host"

if [[ -z "$version" ]]; then
  version="${VERSION:-$(git rev-parse --short=12 HEAD 2>/dev/null || echo dev)}"
fi
export VERSION="$version"
export LGTM_LDFLAGS_PROFILE="${LGTM_LDFLAGS_PROFILE:-release}"
eval "$(zsh scripts/ldflags.sh)"

case "${goos}/${goarch}" in
  darwin/arm64) bun_target="${BUN_TARGET:-bun-darwin-arm64}" ;;
  darwin/amd64) bun_target="${BUN_TARGET:-bun-darwin-x64}" ;;
  linux/arm64) bun_target="${BUN_TARGET:-bun-linux-arm64}" ;;
  linux/amd64) bun_target="${BUN_TARGET:-bun-linux-x64}" ;;
  *) bun_target="${BUN_TARGET:-}" ;;
esac

build_root="$(mktemp -d)"
trap 'rm -rf "$build_root"' EXIT
stage_dir="$build_root/lgtm"
mkdir -p "$stage_dir/bin" "$stage_dir/completions" "$stage_dir/hooks" "$dist_dir"

echo "Building lgtm ${version} for ${goos}/${goarch}"
GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="${CGO_ENABLED:-0}" \
  go build -trimpath -ldflags "$LGTM_LDFLAGS" -o "$stage_dir/bin/lgtm" ./cmd/lgtm
chmod 755 "$stage_dir/bin/lgtm"
ln -sf lgtm "$stage_dir/bin/lgtmr"

echo "Building lgtm-mcp for ${goos}/${goarch}"
(
  cd "$repo_root/mcp"
  bun install --frozen-lockfile
  if [[ -n "$bun_target" ]]; then
    bun build src/stdio.ts --compile --target "$bun_target" --outfile "$stage_dir/bin/lgtm-mcp"
  else
    bun build src/stdio.ts --compile --outfile "$stage_dir/bin/lgtm-mcp"
  fi
)
chmod 755 "$stage_dir/bin/lgtm-mcp"

env -u GOOS -u GOARCH -u CGO_ENABLED go run ./cmd/lgtm-gen-completions "$stage_dir/completions/lgtm.bash" "$stage_dir/completions/_lgtm"

cat > "$stage_dir/hooks/README.txt" <<'EOF'
Repo git hooks are installed by `lgtm init` in each repository.
lgtm init also registers MCP with supported agent CLIs and offers AGENTS.md instructions.
EOF

cat > "$stage_dir/README.txt" <<EOF
lgtm CLI package

Default install (CLI + lgtm-mcp + completions):
  curl -fsSL https://download.lgtm.cx/install.sh | sh

Manual install:
  install -m 755 bin/lgtm ~/.local/bin/lgtm
  ln -sf lgtm ~/.local/bin/lgtmr
  install -m 755 bin/lgtm-mcp ~/.local/bin/lgtm-mcp
EOF

archive_name="lgtm_${version}_${goos}_${goarch}.tar.gz"
archive_path="$dist_dir/$archive_name"
rm -f "$archive_path" "$archive_path.sha256"
tar -czf "$archive_path" -C "$build_root" lgtm

if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 "$archive_path" > "$archive_path.sha256"
else
  sha256sum "$archive_path" > "$archive_path.sha256"
fi

echo "Built $archive_path"
echo "Built $archive_path.sha256"
