set shell := ["zsh", "-cu"]
set dotenv-load := true

repo_root := `pwd`

# `just build` / `just build-release` / `just install` bake the optional gx Cloud endpoints from .env into
# buildconfig. Leave them unset for a local build.
gx_ldflags := "\
  -X github.com/satoricorp/gx/internal/buildconfig.GitHubClientID=${GITHUB_CLIENT_ID:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.ConvexSiteURL=${CONVEX_SITE_URL:-} \
  -X github.com/satoricorp/gx/internal/buildconfig.CloudURL=${GX_CLOUD_URL:-} \
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

# The gx Cloud endpoints are all or nothing. With none set this is a local
# build, which is the default: gx runs entirely on this machine and its cloud
# code stays dormant. With any set, every one must be set and baked in, because
# a binary with only some of them looks configured and publishes nowhere -- the
# silent failure AGENTS.md describes.
verify-bake bin="gx":
  #!/usr/bin/env bash
  set -euo pipefail
  bin="{{bin}}"
  test -x "$bin" || { echo "missing executable: $bin" >&2; exit 1; }
  endpoints=(GITHUB_CLIENT_ID CONVEX_SITE_URL GX_CLOUD_URL)
  configured=0
  for name in "${endpoints[@]}"; do
    [[ -n "${!name:-}" ]] && configured=$((configured + 1))
  done
  if [[ "$configured" -eq 0 ]]; then
    echo "bake ok: $bin (local build, no gx Cloud)"
    exit 0
  fi
  blob="$(strings "$bin")"
  missing=0
  for name in "${endpoints[@]}"; do
    value="${!name:-}"
    if [[ -z "$value" ]]; then
      echo "$name is unset while other gx Cloud endpoints are set; set all of ${endpoints[*]} or none" >&2
      missing=1
    elif ! grep -F -- "$value" <<<"$blob" >/dev/null; then
      echo "missing baked $name in $bin" >&2
      missing=1
    fi
  done
  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi
  echo "bake ok: $bin (gx Cloud at $GX_CLOUD_URL)"

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

