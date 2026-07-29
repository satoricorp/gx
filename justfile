set shell := ["zsh", "-cu"]
set dotenv-load := true

repo_root := `pwd`

# `just build` / `just build-release` / `just install` bake public auth/client endpoints from .env into buildconfig.
tl_ldflags := "\
  -X github.com/satoricorp/totality/internal/buildconfig.GitHubClientID=${GITHUB_CLIENT_ID:-} \
  -X github.com/satoricorp/totality/internal/buildconfig.ConvexSiteURL=${CONVEX_SITE_URL:-} \
  -X github.com/satoricorp/totality/internal/buildconfig.CloudURL=${TOTALITY_CLOUD_URL:-} \
  -X github.com/satoricorp/totality/internal/buildconfig.PostHogKey=${TOTALITY_POSTHOG_KEY:-} \
  -X github.com/satoricorp/totality/internal/buildconfig.PostHogHost=${TOTALITY_POSTHOG_HOST:-} \
  -X github.com/satoricorp/totality/internal/version.Version=${VERSION:-dev}"

build:
  go build -ldflags "{{tl_ldflags}}" -o tl ./cmd/tl

build-release:
  go build -ldflags "{{tl_ldflags}}" -o tl ./cmd/tl

package-cli:
  scripts/package-cli.sh

install-completions:
  mkdir -p ~/.local/share/bash-completion/completions ~/.zfunc
  go run ./cmd/tl-gen-completions ~/.local/share/bash-completion/completions/tl ~/.zfunc/_tl

install:
  mkdir -p ~/.local/bin
  just build
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ./tl || true
  cp ./tl ~/.local/bin/tl
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ~/.local/bin/tl || true
  just install-completions
  just verify-bake ./tl
  just verify-bake ~/.local/bin/tl

verify-bake bin="tl":
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
  check TOTALITY_CLOUD_URL "${TOTALITY_CLOUD_URL:-}"
  check TOTALITY_POSTHOG_KEY "${TOTALITY_POSTHOG_KEY:-}"
  check TOTALITY_POSTHOG_HOST "${TOTALITY_POSTHOG_HOST:-}"
  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi
  echo "bake ok: $bin"

test:
  go test ./...

# Build the MCP: xmcp JS output plus the standalone dist/tl-mcp binary.
mcp-build:
  cd mcp && bun install && bun run build

# Publish @satoricorp/totality to npm. prepublishOnly rebuilds and tests first.
# CI equivalent: push a tag like mcp-v0.1.0 (see .github/workflows/mcp.yml).
mcp-publish:
  cd mcp && bun install --frozen-lockfile && npm publish

run *args:
  {{repo_root}}/tl {{args}}

