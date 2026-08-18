#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
dist_dir="${GX_CLI_DIST_DIR:-$repo_root/dist/cli}"
goos="${GOOS:-darwin}"
goarch="${GOARCH:-$(go env GOARCH)}"
version="${VERSION:-}"

cd "$repo_root"
env_github_client_id="${GITHUB_CLIENT_ID-}"
env_convex_site_url="${CONVEX_SITE_URL-}"
env_gx_cloud_url="${GX_CLOUD_URL-}"
env_gx_posthog_key="${GX_POSTHOG_KEY-}"
env_gx_posthog_host="${GX_POSTHOG_HOST-}"
set -a
if [[ -f .env ]]; then
  # shellcheck disable=SC1091
  source .env
fi
set +a
[[ -n "$env_github_client_id" ]] && export GITHUB_CLIENT_ID="$env_github_client_id"
[[ -n "$env_convex_site_url" ]] && export CONVEX_SITE_URL="$env_convex_site_url"
[[ -n "$env_gx_cloud_url" ]] && export GX_CLOUD_URL="$env_gx_cloud_url"
[[ -n "$env_gx_posthog_key" ]] && export GX_POSTHOG_KEY="$env_gx_posthog_key"
[[ -n "$env_gx_posthog_host" ]] && export GX_POSTHOG_HOST="$env_gx_posthog_host"

if [[ -z "$version" ]]; then
  version="${VERSION:-$(git rev-parse --short=12 HEAD 2>/dev/null || echo dev)}"
fi
export VERSION="$version"
export GX_LDFLAGS_PROFILE="${GX_LDFLAGS_PROFILE:-release}"
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
stage_dir="$build_root/gx"
mkdir -p "$stage_dir/bin" "$stage_dir/completions" "$stage_dir/hooks" "$dist_dir"

echo "Building gx ${version} for ${goos}/${goarch}"
GOOS="$goos" GOARCH="$goarch" CGO_ENABLED="${CGO_ENABLED:-0}" \
  go build -trimpath -ldflags "$GX_LDFLAGS" -o "$stage_dir/bin/gx" ./cmd/gx
chmod 755 "$stage_dir/bin/gx"
ln -sf gx "$stage_dir/bin/gxr"

echo "Building gx-mcp for ${goos}/${goarch}"
# gx-mcp is a Bun binary, so GX_LDFLAGS does not reach it. It gets the same
# production endpoint through --define, which is Bun's compile-time constant
# substitution. Only passed when the value is non-empty: `--define X=` is not
# valid JS, and an empty bake is the same as no bake anyway -- gx-mcp then hands
# the child nothing and the CLI falls back to its own baked endpoint.
mcp_define=()
if [[ -n "${GX_CLOUD_URL:-}" ]]; then
  mcp_define=(--define "GX_BAKED_CLOUD_URL=\"${GX_CLOUD_URL}\"")
fi
(
  cd "$repo_root/mcp"
  bun install --frozen-lockfile
  if [[ -n "$bun_target" ]]; then
    bun build src/stdio.ts --compile --target "$bun_target" "${mcp_define[@]}" --outfile "$stage_dir/bin/gx-mcp"
  else
    bun build src/stdio.ts --compile "${mcp_define[@]}" --outfile "$stage_dir/bin/gx-mcp"
  fi
)
chmod 755 "$stage_dir/bin/gx-mcp"

# The same fail-closed rule the Go bake guard follows: a shipped gx-mcp that
# cannot name its server is the build this check exists to stop, and it is
# invisible from the outside -- the MCP starts fine and every review it runs is
# quietly degraded.
if [[ -n "${GX_CLOUD_URL:-}" ]]; then
  # grep -a on the file, not strings(1). gx-mcp is cross-compiled for the
  # target platform, and macOS strings refuses a Mach-O it did not expect --
  # "LC_CODE_SIGNATURE command extends past the end of the file" -- which reads
  # here as an absent bake and fails a build whose bake was fine. Measured: the
  # native arm64 leg passed and the x64 leg failed on the same value. grep -a
  # treats the binary as bytes and does not care what object format it is.
  if ! grep -aF -- "$GX_CLOUD_URL" "$stage_dir/bin/gx-mcp" >/dev/null; then
    echo "missing baked GX_CLOUD_URL in $stage_dir/bin/gx-mcp" >&2
    exit 1
  fi
  echo "bake ok: $stage_dir/bin/gx-mcp"
elif [[ "${GX_ALLOW_UNBAKED:-}" != "1" ]]; then
  echo "cannot verify GX_CLOUD_URL for gx-mcp: unset in the environment and .env" >&2
  echo "set GX_CLOUD_URL, or GX_ALLOW_UNBAKED=1 to build a cloud-dead gx-mcp on purpose" >&2
  exit 1
fi

env -u GOOS -u GOARCH -u CGO_ENABLED go run ./cmd/gx-gen-completions "$stage_dir/completions/gx.bash" "$stage_dir/completions/_gx"

cat > "$stage_dir/hooks/README.txt" <<'EOF'
Repo git hooks are installed by `gx init` in each repository.
gx init also registers MCP with supported agent CLIs and offers AGENTS.md instructions.
EOF

cat > "$stage_dir/README.txt" <<EOF
gx CLI package

Default install (CLI + gx-mcp + completions):
  curl -fsSL https://download.gx.run/install.sh | sh

Manual install:
  install -m 755 bin/gx ~/.local/bin/gx
  ln -sf gx ~/.local/bin/gxr
  install -m 755 bin/gx-mcp ~/.local/bin/gx-mcp
EOF

archive_name="gx_${version}_${goos}_${goarch}.tar.gz"
archive_path="$dist_dir/$archive_name"
rm -f "$archive_path" "$archive_path.sha256"
tar -czf "$archive_path" -C "$build_root" gx

if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 "$archive_path" > "$archive_path.sha256"
else
  sha256sum "$archive_path" > "$archive_path.sha256"
fi

echo "Built $archive_path"
echo "Built $archive_path.sha256"
