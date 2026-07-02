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
env_gx_appcast_url="${GX_APPCAST_URL-}"
env_gx_sparkle_public_ed_key="${GX_SPARKLE_PUBLIC_ED_KEY-}"
env_gx_app_build_version="${GX_APP_BUILD_VERSION-}"
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
[[ -n "$env_gx_appcast_url" ]] && export GX_APPCAST_URL="$env_gx_appcast_url"
[[ -n "$env_gx_sparkle_public_ed_key" ]] && export GX_SPARKLE_PUBLIC_ED_KEY="$env_gx_sparkle_public_ed_key"
[[ -n "$env_gx_app_build_version" ]] && export GX_APP_BUILD_VERSION="$env_gx_app_build_version"
export GX_LDFLAGS_PROFILE="${GX_LDFLAGS_PROFILE:-release}"
eval "$(zsh scripts/ldflags.sh)"

mkdir -p "$menubar_dir/bin" "$dist_dir"
go build -ldflags "$GX_LDFLAGS" -o "$menubar_dir/bin/gx" ./cmd/gx
git_tag_version="${VERSION:-dev}"
app_short_version="${git_tag_version#v}"
if [[ "$app_short_version" == "dev" || -z "$app_short_version" ]]; then
  app_short_version="0.1.0"
fi
app_build_version="${GX_APP_BUILD_VERSION:-$(git rev-list --count HEAD 2>/dev/null || date +%s)}"
appcast_url="${GX_APPCAST_URL:-https://download.gx.run/appcast.xml}"
sparkle_public_ed_key="${GX_SPARKLE_PUBLIC_ED_KEY:-}"
cli_version_json="$("$menubar_dir/bin/gx" version --json)"
cli_version="$(printf '%s' "$cli_version_json" | plutil -extract version raw -o - -)"

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
set_plist_string() {
  local key="$1"
  local value="$2"
  local plist="$app_path/Contents/Info.plist"
  /usr/libexec/PlistBuddy -c "Set :$key $value" "$plist" 2>/dev/null ||
    /usr/libexec/PlistBuddy -c "Add :$key string $value" "$plist"
}
set_plist_string "GXGitTagVersion" "$git_tag_version"
set_plist_string "GXCLIVersion" "$cli_version"
set_plist_string "CFBundleShortVersionString" "$app_short_version"
set_plist_string "CFBundleVersion" "$app_build_version"
set_plist_string "SUFeedURL" "$appcast_url"
set_plist_string "SUPublicEDKey" "$sparkle_public_ed_key"
cp "$swift_exe" "$app_path/Contents/MacOS/GX"
chmod 755 "$app_path/Contents/MacOS/GX"
sparkle_framework="$(find -L "$menubar_dir/.build" -path "*/Sparkle.framework" -type d -print -quit)"
if [[ -n "$sparkle_framework" ]]; then
  mkdir -p "$app_path/Contents/Frameworks"
  ditto "$sparkle_framework" "$app_path/Contents/Frameworks/Sparkle.framework"
fi
cp "$menubar_dir/assets/trayTemplate.png" "$app_path/Contents/Resources/trayTemplate.png"
cp "$menubar_dir/assets/trayTemplate@2x.png" "$app_path/Contents/Resources/trayTemplate@2x.png"
cp "$menubar_dir/assets/appIcon.icns" "$app_path/Contents/Resources/GX.icns"

resources_dir="$app_path/Contents/Resources"
mkdir -p "$resources_dir/bin"
cp "$menubar_dir/bin/gx" "$resources_dir/bin/gx"
chmod 755 "$resources_dir/bin/gx"
ln -sf gx "$resources_dir/bin/gxg"
ln -sf gx "$resources_dir/bin/gxr"
ln -sf gx "$resources_dir/bin/gxs"

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
