#!/usr/bin/env zsh
# Emit shell to export LGTM_LDFLAGS from repo .env (sourced by just via dotenv-load or caller).
set -euo pipefail

github_id="${GITHUB_CLIENT_ID:-}"
convex_site="${CONVEX_SITE_URL:-}"
cloud_url="${LGTM_CLOUD_URL:-}"
posthog_key="${LGTM_POSTHOG_KEY:-}"
posthog_host="${LGTM_POSTHOG_HOST:-}"
lgtm_version="${VERSION:-dev}"

flags=()

if [[ -n "${github_id}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/buildconfig.GitHubClientID=${github_id}")
fi
if [[ -n "${convex_site}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/buildconfig.ConvexSiteURL=${convex_site}")
fi
if [[ -n "${cloud_url}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/buildconfig.CloudURL=${cloud_url}")
fi
if [[ -n "${posthog_key}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/buildconfig.PostHogKey=${posthog_key}")
fi
if [[ -n "${posthog_host}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/buildconfig.PostHogHost=${posthog_host}")
fi
if [[ -n "${lgtm_version}" ]]; then
  flags+=("-X" "github.com/satoricorp/lgtm/internal/version.Version=${lgtm_version}")
fi

printf 'export LGTM_LDFLAGS=%q\n' "${flags[*]}"
