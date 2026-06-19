#!/usr/bin/env zsh
# Emit shell to export GX_LDFLAGS from repo .env (sourced by just via dotenv-load or caller).
set -euo pipefail

github_id="${GITHUB_CLIENT_ID:-}"
convex_site="${CONVEX_SITE_URL:-}"
cloud_url="${GX_CLOUD_URL:-}"
gx_version="${VERSION:-dev}"

flags=()

if [[ -n "${github_id}" ]]; then
  flags+=("-X" "github.com/satoricorp/gx/internal/buildconfig.GitHubClientID=${github_id}")
fi
if [[ -n "${convex_site}" ]]; then
  flags+=("-X" "github.com/satoricorp/gx/internal/buildconfig.ConvexSiteURL=${convex_site}")
fi
if [[ -n "${cloud_url}" ]]; then
  flags+=("-X" "github.com/satoricorp/gx/internal/buildconfig.CloudURL=${cloud_url}")
fi
if [[ -n "${gx_version}" ]]; then
  flags+=("-X" "github.com/satoricorp/gx/internal/version.Version=${gx_version}")
fi

printf 'export GX_LDFLAGS=%q\n' "${flags[*]}"
