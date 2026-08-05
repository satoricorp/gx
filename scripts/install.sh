#!/bin/sh
set -eu

base_url="${LGTM_INSTALL_BASE_URL:-https://download.lgtm.cx}"
install_dir="${LGTM_INSTALL_DIR:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d 2>/dev/null || mktemp -d -t lgtm-install)"

fail() {
  echo "lgtm install: $*" >&2
  exit 1
}

usage() {
  cat <<EOF
lgtm installer

Default install: lgtm CLI, lgtm-mcp, shell completions.
Repo git hooks are installed later by lgtm init.

Usage:
  curl -fsSL https://download.lgtm.cx/install.sh | sh

Options:
  -h, --help       Show this help

Environment:
  LGTM_INSTALL_BASE_URL   Download base URL (default: https://download.lgtm.cx)
  LGTM_INSTALL_DIR        CLI install directory (default: ~/.local/bin)
EOF
}

for arg in "$@"; do
  case "$arg" in
    -h|--help)
      usage
      exit 0
      ;;
    *)
      fail "unknown option: $arg (try --help)"
      ;;
  esac
done

cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

need() {
  command -v "$1" >/dev/null 2>&1 || fail "missing required command: $1"
}

download() {
  url="$1"
  output="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
  else
    fail "missing curl or wget"
  fi
}


case "$(uname -s)" in
  Darwin) os="darwin" ;;
  Linux) os="linux" ;;
  *) fail "unsupported OS: $(uname -s)" ;;
esac

case "$(uname -m)" in
  arm64|aarch64) arch="arm64" ;;
  x86_64|amd64) arch="amd64" ;;
  *) fail "unsupported architecture: $(uname -m)" ;;
esac

need tar
need awk
need mkdir

base_url="${base_url%/}"
archive="lgtm_${os}_${arch}.tar.gz"
archive_url="$base_url/cli/latest/$archive"
checksum_url="$archive_url.sha256"
archive_path="$tmp_dir/$archive"
checksum_path="$archive_path.sha256"

echo "Downloading $archive_url"
download "$archive_url" "$archive_path"
download "$checksum_url" "$checksum_path"

expected="$(awk '{print $1}' "$checksum_path")"
if command -v shasum >/dev/null 2>&1; then
  actual="$(shasum -a 256 "$archive_path" | awk '{print $1}')"
elif command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "$archive_path" | awk '{print $1}')"
else
  fail "missing shasum or sha256sum"
fi

if [ "$expected" != "$actual" ]; then
  fail "checksum mismatch for $archive"
fi

tar -xzf "$archive_path" -C "$tmp_dir"
test -x "$tmp_dir/lgtm/bin/lgtm" || fail "archive is missing lgtm"
test -x "$tmp_dir/lgtm/bin/lgtm-mcp" || fail "archive is missing lgtm-mcp"

mkdir -p "$install_dir"
install -m 755 "$tmp_dir/lgtm/bin/lgtm" "$install_dir/lgtm"
install -m 755 "$tmp_dir/lgtm/bin/lgtm-mcp" "$install_dir/lgtm-mcp"
ln -sf lgtm "$install_dir/lgtmr"

if [ -f "$tmp_dir/lgtm/completions/lgtm.bash" ]; then
  mkdir -p "$HOME/.local/share/bash-completion/completions"
  install -m 644 "$tmp_dir/lgtm/completions/lgtm.bash" "$HOME/.local/share/bash-completion/completions/lgtm"
fi
if [ -f "$tmp_dir/lgtm/completions/_lgtm" ]; then
  mkdir -p "$HOME/.zfunc"
  install -m 644 "$tmp_dir/lgtm/completions/_lgtm" "$HOME/.zfunc/_lgtm"
fi

echo "Installed lgtm to $install_dir/lgtm"
echo "Installed MCP and aliases"
if ! command -v lgtm >/dev/null 2>&1; then
  echo "Add $install_dir to PATH before running lgtm."
fi
# Cyan ANSI 6 + bold matches lgtm version / logo (internal/cli/logo.go).
lgtm_auth_login="lgtm auth login"
lgtm_init="lgtm init"
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  lgtm_auth_login="$(printf '\033[1;36mtl auth login\033[0m')"
  lgtm_init="$(printf '\033[1;36mtl init\033[0m')"
fi
echo ""
printf '\tRun %s to login.\n' "$lgtm_auth_login"
printf '\tRun %s in each repo to initialize lgtm.\n' "$lgtm_init"
"$install_dir/lgtm" version
