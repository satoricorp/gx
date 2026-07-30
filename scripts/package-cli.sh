#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
dist_dir="${TOTALITY_CLI_DIST_DIR:-$repo_root/dist/cli}"
goos="${GOOS:-darwin}"
goarch="${GOARCH:-$(go env GOARCH)}"
version="${VERSION:-}"

cd "$repo_root"
env_github_client_id="${GITHUB_CLIENT_ID-}"
env_convex_site_url="${CONVEX_SITE_URL-}"
env_tl_cloud_url="${TOTALITY_CLOUD_URL-}"
env_tl_posthog_key="${TOTALITY_POSTHOG_KEY-}"
env_tl_posthog_host="${TOTALITY_POSTHOG_HOST-}"
set -a
if [[ -f .env ]]; then
  # shellcheck disable=SC1091
  source .env
fi
set +a
[[ -n "$env_github_client_id" ]] && export GITHUB_CLIENT_ID="$env_github_client_id"
[[ -n "$env_convex_site_url" ]] && export CONVEX_SITE_URL="$env_convex_site_url"
[[ -n "$env_tl_cloud_url" ]] && export TOTALITY_CLOUD_URL="$env_tl_cloud_url"
[[ -n "$env_tl_posthog_key" ]] && export TOTALITY_POSTHOG_KEY="$env_tl_posthog_key"
[[ -n "$env_tl_posthog_host" ]] && export TOTALITY_POSTHOG_HOST="$env_tl_posthog_host"

if [[ -z "$version" ]]; then
  version="${VERSION:-$(git rev-parse --short=12 HEAD 2>/dev/null || echo dev)}"
fi
export VERSION="$version"
export TOTALITY_LDFLAGS_PROFILE="${TOTALITY_LDFLAGS_PROFILE:-release}"
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
stage_dir="$build_root/tx"
mkdir -p "$stage_dir/bin" "$stage_dir/completions" "$stage_dir/hooks" "$dist_dir"

echo "Building tx ${version} for ${goos}/${goarch}"
GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="${CGO_ENABLED:-0}" \
  go build -trimpath -ldflags "$TOTALITY_LDFLAGS" -o "$stage_dir/bin/tx" ./cmd/tx
chmod 755 "$stage_dir/bin/tx"
ln -sf tx "$stage_dir/bin/txr"

echo "Building tx-mcp for ${goos}/${goarch}"
(
  cd "$repo_root/mcp"
  bun install --frozen-lockfile
  if [[ -n "$bun_target" ]]; then
    bun build src/stdio.ts --compile --target "$bun_target" --outfile "$stage_dir/bin/tx-mcp"
  else
    bun build src/stdio.ts --compile --outfile "$stage_dir/bin/tx-mcp"
  fi
)
chmod 755 "$stage_dir/bin/tx-mcp"

env -u GOOS -u GOARCH -u CGO_ENABLED go run ./cmd/tx-gen-completions "$stage_dir/completions/tx.bash" "$stage_dir/completions/_tx"

cat > "$stage_dir/hooks/README.txt" <<'EOF'
Repo git hooks are installed by `tx init` in each repository.
Totality init also registers MCP with supported agent CLIs and offers AGENTS.md instructions.
EOF

cat > "$stage_dir/README.txt" <<EOF
Totality CLI package

Default install (CLI + tx-mcp + completions):
  curl -fsSL https://download.totality.sh/install.sh | sh

Manual install:
  install -m 755 bin/tx ~/.local/bin/tx
  ln -sf tx ~/.local/bin/txr
  install -m 755 bin/tx-mcp ~/.local/bin/tx-mcp
EOF

archive_name="tx_${version}_${goos}_${goarch}.tar.gz"
archive_path="$dist_dir/$archive_name"
rm -f "$archive_path" "$archive_path.sha256"
tar -czf "$archive_path" -C "$build_root" tx

if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 "$archive_path" > "$archive_path.sha256"
else
  sha256sum "$archive_path" > "$archive_path.sha256"
fi

echo "Built $archive_path"
echo "Built $archive_path.sha256"
