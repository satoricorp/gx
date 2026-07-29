#!/usr/bin/env zsh
# Emit shell to export TOTALITY_LDFLAGS from repo .env (sourced by just via dotenv-load or caller).
set -euo pipefail

github_id="${GITHUB_CLIENT_ID:-}"
convex_site="${CONVEX_SITE_URL:-}"
cloud_url="${TOTALITY_CLOUD_URL:-}"
posthog_key="${TOTALITY_POSTHOG_KEY:-}"
posthog_host="${TOTALITY_POSTHOG_HOST:-}"
tl_version="${VERSION:-dev}"

flags=()

if [[ -n "${github_id}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/buildconfig.GitHubClientID=${github_id}")
fi
if [[ -n "${convex_site}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/buildconfig.ConvexSiteURL=${convex_site}")
fi
if [[ -n "${cloud_url}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/buildconfig.CloudURL=${cloud_url}")
fi
if [[ -n "${posthog_key}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/buildconfig.PostHogKey=${posthog_key}")
fi
if [[ -n "${posthog_host}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/buildconfig.PostHogHost=${posthog_host}")
fi
if [[ -n "${tl_version}" ]]; then
  flags+=("-X" "github.com/satoricorp/totality/internal/version.Version=${tl_version}")
fi

printf 'export TOTALITY_LDFLAGS=%q\n' "${flags[*]}"
