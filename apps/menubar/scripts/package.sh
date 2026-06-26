#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/../../.." && pwd)"
menubar_dir="$repo_root/apps/menubar"
dist_dir="$menubar_dir/dist"
app_path="$dist_dir/GX.app"

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
export GX_LDFLAGS_PROFILE="${GX_LDFLAGS_PROFILE:-release}"
eval "$(zsh scripts/ldflags.sh)"

mkdir -p "$menubar_dir/bin" "$dist_dir"
go build -ldflags "$GX_LDFLAGS" -o "$menubar_dir/bin/gx" ./cmd/gx

if [[ -f "$repo_root/mcp/package.json" ]]; then
  (
    cd "$repo_root/mcp"
    bun install --frozen-lockfile
    bun run build
  )
fi

swift build --package-path "$menubar_dir" -c release
swift_bin_dir="$(swift build --package-path "$menubar_dir" -c release --show-bin-path)"
swift_exe="$swift_bin_dir/GXMenuBar"
test -x "$swift_exe" || { echo "missing Swift executable: $swift_exe" >&2; exit 1; }

rm -rf "$app_path"
mkdir -p "$app_path/Contents/MacOS" "$app_path/Contents/Resources"
cp "$menubar_dir/Resources/Info.plist" "$app_path/Contents/Info.plist"
cp "$swift_exe" "$app_path/Contents/MacOS/GX"
chmod 755 "$app_path/Contents/MacOS/GX"
cp "$menubar_dir/assets/trayTemplate.png" "$app_path/Contents/Resources/trayTemplate.png"
cp "$menubar_dir/assets/trayTemplate@2x.png" "$app_path/Contents/Resources/trayTemplate@2x.png"
cp "$menubar_dir/assets/appIcon.icns" "$app_path/Contents/Resources/GX.icns"

resources_dir="$app_path/Contents/Resources"
mkdir -p "$resources_dir/bin"
cp "$menubar_dir/bin/gx" "$resources_dir/bin/gx"
chmod 755 "$resources_dir/bin/gx"

if [[ ! -x "$repo_root/mcp/dist/gx-mcp" ]]; then
  echo "missing MCP stdio binary: $repo_root/mcp/dist/gx-mcp" >&2
  exit 1
fi
cp "$repo_root/mcp/dist/gx-mcp" "$resources_dir/bin/gx-mcp"
chmod 755 "$resources_dir/bin/gx-mcp"

codesign --force --deep --sign - "$app_path"
rm -f "$dist_dir/GX-macOS.zip"
(
  cd "$(dirname "$app_path")"
  ditto -c -k --keepParent "$(basename "$app_path")" "$dist_dir/GX-macOS.zip"
)

dmg_root="$(mktemp -d)"
trap 'rm -rf "$dmg_root"' EXIT
cp -R "$app_path" "$dmg_root/GX.app"
ln -s /Applications "$dmg_root/Applications"
rm -f "$dist_dir/GX-macOS.dmg"
hdiutil create -volname "GX" -srcfolder "$dmg_root" -ov -format UDZO "$dist_dir/GX-macOS.dmg" >/dev/null

echo "Built $app_path"
echo "Built $dist_dir/GX-macOS.zip"
echo "Built $dist_dir/GX-macOS.dmg"
