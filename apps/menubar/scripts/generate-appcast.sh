#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/../../.." && pwd)"
menubar_dir="$repo_root/apps/menubar"
dist_dir="$menubar_dir/dist"
app_path="$dist_dir/GX.app"
archive_path="$dist_dir/GX-macOS.zip"
appcast_dir="${GX_APPCAST_DIR:-$dist_dir/appcast}"
download_url_prefix="${GX_APPCAST_DOWNLOAD_URL_PREFIX:-}"

find_generate_appcast() {
  if [[ -n "${SPARKLE_GENERATE_APPCAST:-}" && -x "${SPARKLE_GENERATE_APPCAST:-}" ]]; then
    printf '%s\n' "$SPARKLE_GENERATE_APPCAST"
    return 0
  fi
  local candidate
  for candidate in \
    "$menubar_dir/.local-artifacts/bin/generate_appcast" \
    "$menubar_dir/.build/artifacts/sparkle/Sparkle/bin/generate_appcast" \
    "$menubar_dir/.build/checkouts/Sparkle/bin/generate_appcast"; do
    if [[ -x "$candidate" ]]; then
      printf '%s\n' "$candidate"
      return 0
    fi
  done
  return 1
}

test -d "$app_path" || { echo "missing app bundle: $app_path" >&2; exit 1; }
test -f "$archive_path" || { echo "missing app archive: $archive_path" >&2; exit 1; }
generate_appcast="$(find_generate_appcast)" || {
  echo "missing Sparkle generate_appcast; set SPARKLE_GENERATE_APPCAST" >&2
  exit 1
}

version="$(plutil -extract CFBundleShortVersionString raw -o - "$app_path/Contents/Info.plist")"
build="$(plutil -extract CFBundleVersion raw -o - "$app_path/Contents/Info.plist")"
mkdir -p "$appcast_dir"
versioned_archive="$appcast_dir/GX-${version}-${build}-macOS.zip"
cp "$archive_path" "$versioned_archive"

args=("$appcast_dir")
if [[ -n "$download_url_prefix" ]]; then
  args=(--download-url-prefix "$download_url_prefix" "${args[@]}")
fi
if [[ -n "${GX_SPARKLE_PRIVATE_ED_KEY_FILE:-}" ]]; then
  args=(--ed-key-file "$GX_SPARKLE_PRIVATE_ED_KEY_FILE" "${args[@]}")
fi

if [[ -n "${GX_SPARKLE_PRIVATE_ED_KEY:-}" ]]; then
  printf '%s' "$GX_SPARKLE_PRIVATE_ED_KEY" | "$generate_appcast" --ed-key-file - "${args[@]}"
else
  "$generate_appcast" "${args[@]}"
fi

echo "Built $appcast_dir/appcast.xml"
