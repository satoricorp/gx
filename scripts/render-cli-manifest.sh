#!/usr/bin/env bash
set -euo pipefail

dist_dir="${1:-dist/cli}"
base_url="${TOTALITY_DOWNLOAD_BASE_URL:-https://download.totality.sh}"
version="${VERSION:-}"
git_sha="${GITHUB_SHA:-}"
published_at="${PUBLISHED_AT:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
manifest_path="$dist_dir/manifest.json"

if [[ -z "$version" ]]; then
  version="$(git rev-parse --short=12 HEAD 2>/dev/null || echo dev)"
fi
if [[ -z "$git_sha" ]]; then
  git_sha="$(git rev-parse HEAD 2>/dev/null || echo "")"
fi

base_url="${base_url%/}"
mkdir -p "$dist_dir"

json_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  printf '%s' "$value"
}

archive_rows=()
for archive in "$dist_dir"/tx_"$version"_*.tar.gz; do
  [[ -f "$archive" ]] || continue
  base="$(basename "$archive")"
  rest="${base#tx_${version}_}"
  platform="${rest%.tar.gz}"
  goos="${platform%_*}"
  goarch="${platform##*_}"
  latest_name="tx_${goos}_${goarch}.tar.gz"
  sha_path="$archive.sha256"
  [[ -f "$sha_path" ]] || { echo "missing checksum for $archive" >&2; exit 1; }
  sha="$(awk '{print $1}' "$sha_path")"
  size="$(wc -c < "$archive" | tr -d ' ')"
  key="${goos}/${goarch}"
  latest_url="$base_url/cli/latest/$latest_name"
  versioned_url="$base_url/cli/releases/$git_sha/$base"
  archive_rows+=("$key|$latest_url|$versioned_url|$sha|$size")
done

if [[ "${#archive_rows[@]}" -eq 0 ]]; then
  echo "no CLI archives found for version $version in $dist_dir" >&2
  exit 1
fi

{
  printf '{\n'
  printf '  "schema_version": 1,\n'
  printf '  "version": "%s",\n' "$(json_escape "$version")"
  printf '  "git_sha": "%s",\n' "$(json_escape "$git_sha")"
  printf '  "published_at": "%s",\n' "$(json_escape "$published_at")"
  printf '  "install_url": "%s/install.sh",\n' "$(json_escape "$base_url")"
  printf '  "install_command": "curl -fsSL %s/install.sh | sh",\n' "$(json_escape "$base_url")"
  printf '  "assets": {\n'
  for i in "${!archive_rows[@]}"; do
    IFS='|' read -r key latest_url versioned_url sha size <<< "${archive_rows[$i]}"
    comma=","
    if [[ "$i" -eq $((${#archive_rows[@]} - 1)) ]]; then
      comma=""
    fi
    printf '    "%s": {\n' "$(json_escape "$key")"
    printf '      "url": "%s",\n' "$(json_escape "$latest_url")"
    printf '      "versioned_url": "%s",\n' "$(json_escape "$versioned_url")"
    printf '      "sha256": "%s",\n' "$(json_escape "$sha")"
    printf '      "size": %s\n' "$size"
    printf '    }%s\n' "$comma"
  done
  printf '  }\n'
  printf '}\n'
} > "$manifest_path"

echo "Wrote $manifest_path"
