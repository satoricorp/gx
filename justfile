set shell := ["zsh", "-cu"]
set dotenv-load := true

repo_root := `pwd`

# `just build` / `just build-release` / `just install` bake public auth/client endpoints from .env into buildconfig.
lgtm_ldflags := "\
  -X github.com/satoricorp/lgtm/internal/buildconfig.GitHubClientID=${GITHUB_CLIENT_ID:-} \
  -X github.com/satoricorp/lgtm/internal/buildconfig.ConvexSiteURL=${CONVEX_SITE_URL:-} \
  -X github.com/satoricorp/lgtm/internal/buildconfig.CloudURL=${LGTM_CLOUD_URL:-} \
  -X github.com/satoricorp/lgtm/internal/buildconfig.PostHogKey=${LGTM_POSTHOG_KEY:-} \
  -X github.com/satoricorp/lgtm/internal/buildconfig.PostHogHost=${LGTM_POSTHOG_HOST:-} \
  -X github.com/satoricorp/lgtm/internal/version.Version=${VERSION:-dev}"

build:
  go build -ldflags "{{lgtm_ldflags}}" -o lgtm ./cmd/lgtm

build-release:
  go build -ldflags "{{lgtm_ldflags}}" -o lgtm ./cmd/lgtm

package-cli:
  scripts/package-cli.sh

install-completions:
  mkdir -p ~/.local/share/bash-completion/completions ~/.zfunc
  go run ./cmd/lgtm-gen-completions ~/.local/share/bash-completion/completions/lgtm ~/.zfunc/_lgtm

install:
  mkdir -p ~/.local/bin
  just build
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ./lgtm || true
  cp ./lgtm ~/.local/bin/lgtm
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ~/.local/bin/lgtm || true
  just install-completions
  just verify-bake ./lgtm
  just verify-bake ~/.local/bin/lgtm

verify-bake bin="lgtm":
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
  check LGTM_CLOUD_URL "${LGTM_CLOUD_URL:-}"
  check LGTM_POSTHOG_KEY "${LGTM_POSTHOG_KEY:-}"
  check LGTM_POSTHOG_HOST "${LGTM_POSTHOG_HOST:-}"
  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi
  echo "bake ok: $bin"

test:
  go test ./...

# Build the MCP: xmcp JS output plus the standalone dist/lgtm-mcp binary.
mcp-build:
  cd mcp && bun install && bun run build

# Publish @satoricorp/lgtm to npm. prepublishOnly rebuilds and tests first.
# CI equivalent: push a tag like mcp-v0.1.0 (see .github/workflows/mcp.yml).
mcp-publish:
  cd mcp && bun install --frozen-lockfile && npm publish

run *args:
  {{repo_root}}/lgtm {{args}}

