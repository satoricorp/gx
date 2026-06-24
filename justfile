set shell := ["zsh", "-cu"]
set dotenv-load := true

repo_root := `pwd`

# `just build` / `just build-release` / `just install` bake public auth/client endpoints from .env into buildconfig.
gx_ldflags := "\
  -X github.com/satoricorp/gx/internal/buildconfig.GitHubClientID=${GITHUB_CLIENT_ID:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.ConvexSiteURL=${CONVEX_SITE_URL:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.CloudURL=${GX_CLOUD_URL:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.PostHogKey=${GX_POSTHOG_KEY:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.PostHogHost=${GX_POSTHOG_HOST:-} \
  -X github.com/satoricorp/gx/internal/version.Version=${VERSION:-dev}"

build:
  go build -ldflags "{{gx_ldflags}}" -o gx ./cmd/gx

build-release:
  go build -ldflags "{{gx_ldflags}}" -o gx ./cmd/gx

install-completions:
  mkdir -p ~/.local/share/bash-completion/completions ~/.zfunc
  go run ./cmd/gx-gen-completions ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx

install:
  mkdir -p ~/.local/bin
  just build
  cp ./gx ~/.local/bin/gx
  just install-completions
  just verify-bake ./gx
  just verify-bake ~/.local/bin/gx

verify-bake bin="gx":
  #!/usr/bin/env bash
  set -euo pipefail
  bin="{{bin}}"
  test -x "$bin" || { echo "missing executable: $bin" >&2; exit 1; }
  blob="$(strings "$bin")"
  missing=0
  check() {
    local label="$1" value="$2"
    if [[ -z "$value" ]]; then
      echo "skip $label (empty in .env)" >&2
      return 0
    fi
    if ! grep -F -- "$value" <<<"$blob" >/dev/null; then
      echo "missing baked $label in $bin" >&2
      missing=1
    fi
  }
  check GITHUB_CLIENT_ID "${GITHUB_CLIENT_ID:-}"
  check CONVEX_SITE_URL "${CONVEX_SITE_URL:-}"
  check GX_CLOUD_URL "${GX_CLOUD_URL:-}"
  check GX_POSTHOG_KEY "${GX_POSTHOG_KEY:-}"
  check GX_POSTHOG_HOST "${GX_POSTHOG_HOST:-}"
  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi
  echo "bake ok: $bin"

test:
  go test ./...

test-e2e:
  go test -tags=e2e ./test/e2e

test-all:
  go test ./...
  go test -tags=e2e ./test/e2e

run *args:
  {{repo_root}}/gx {{args}}

menubar:
  rm -f ~/.local/bin/gx
  apps/menubar/scripts/package.sh
