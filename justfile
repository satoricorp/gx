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

package-cli:
  scripts/package-cli.sh

install-completions:
  mkdir -p ~/.local/share/bash-completion/completions ~/.zfunc
  go run ./cmd/gx-gen-completions ~/.local/share/bash-completion/completions/gx ~/.zfunc/_gx

# The bake is verified BEFORE the copy. Verifying only afterwards still left a
# broken binary installed at ~/.local/bin/gx and merely reported it, so the
# guard named the problem while shipping it anyway.
install:
  mkdir -p ~/.local/bin
  just build
  just verify-bake ./gx
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ./gx || true
  cp ./gx ~/.local/bin/gx
  command -v codesign >/dev/null 2>&1 && codesign --force --sign - ~/.local/bin/gx || true
  just install-completions
  just verify-bake ~/.local/bin/gx

verify-bake bin="gx":
  #!/usr/bin/env bash
  set -euo pipefail
  bin="{{bin}}"
  test -x "$bin" || { echo "missing executable: $bin" >&2; exit 1; }
  blob="$(strings "$bin")"
  missing=0
  unverifiable=0
  check() {
    local label="$1" value="$2"
    if [[ -z "$value" ]]; then
      # A check with nothing to check against is not a pass. This used to
      # `return 0`, so a checkout with no .env printed "bake ok" for a binary
      # whose cloud path was entirely dead -- the guard endorsing the exact
      # build it exists to catch.
      echo "cannot verify $label: unset in the environment and .env" >&2
      unverifiable=$((unverifiable + 1))
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
  if [[ "$unverifiable" -gt 0 && "${GX_ALLOW_UNBAKED:-}" != "1" ]]; then
    echo "" >&2
    echo "$bin cannot be verified: $unverifiable endpoint(s) unset." >&2
    echo "Such a binary compiles and runs, but its cloud path is dead: every" >&2
    echo "publish queues into ~/.gx/publish-outbox and no upload is ever" >&2
    echo "attempted. Populate .env, or set GX_ALLOW_UNBAKED=1 to build a" >&2
    echo "local-only binary on purpose." >&2
    exit 1
  fi
  if [[ "$unverifiable" -gt 0 ]]; then
    echo "warning: GX_ALLOW_UNBAKED=1 -- no cloud endpoints in $bin; gx Cloud is disabled in it" >&2
  fi
  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi
  echo "bake ok: $bin"

test:
  go test ./...

# Build the MCP: xmcp JS output plus the standalone dist/gx-mcp binary.
mcp-build:
  cd mcp && bun install && bun run build

# Publish @satoricorp/gx to npm. prepublishOnly rebuilds and tests first.
# CI equivalent: push a tag like mcp-v0.1.0 (see .github/workflows/mcp.yml).
mcp-publish:
  cd mcp && bun install --frozen-lockfile && npm publish

run *args:
  {{repo_root}}/gx {{args}}

