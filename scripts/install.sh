#!/bin/sh
set -eu

base_url="${TOTALITY_INSTALL_BASE_URL:-https://download.totality.sh}"
install_dir="${TOTALITY_INSTALL_DIR:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d 2>/dev/null || mktemp -d -t totality-install)"

fail() {
  echo "tx install: $*" >&2
  exit 1
}

usage() {
  cat <<EOF
Totality installer

Default install: tx CLI, tx-mcp, shell completions.
Repo git hooks are installed later by tx init.

Usage:
  curl -fsSL https://download.totality.sh/install.sh | sh

Options:
  -h, --help       Show this help

Environment:
  TOTALITY_INSTALL_BASE_URL   Download base URL (default: https://download.totality.sh)
  TOTALITY_INSTALL_DIR        CLI install directory (default: ~/.local/bin)
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
archive="tx_${os}_${arch}.tar.gz"
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
test -x "$tmp_dir/tx/bin/tx" || fail "archive is missing tx"
test -x "$tmp_dir/tx/bin/tx-mcp" || fail "archive is missing tx-mcp"

mkdir -p "$install_dir"
install -m 755 "$tmp_dir/tx/bin/tx" "$install_dir/tx"
install -m 755 "$tmp_dir/tx/bin/tx-mcp" "$install_dir/tx-mcp"
ln -sf tx "$install_dir/txr"

if [ -f "$tmp_dir/tx/completions/tx.bash" ]; then
  mkdir -p "$HOME/.local/share/bash-completion/completions"
  install -m 644 "$tmp_dir/tx/completions/tx.bash" "$HOME/.local/share/bash-completion/completions/tx"
fi
if [ -f "$tmp_dir/tx/completions/_tx" ]; then
  mkdir -p "$HOME/.zfunc"
  install -m 644 "$tmp_dir/tx/completions/_tx" "$HOME/.zfunc/_tx"
fi

echo "Installed tx to $install_dir/tx"
echo "Installed MCP and aliases"
if ! command -v tx >/dev/null 2>&1; then
  echo "Add $install_dir to PATH before running tx."
fi
# Cyan ANSI 6 + bold matches tx version / logo (internal/cli/logo.go).
tx_auth_login="tx auth login"
tx_init="tx init"
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  tx_auth_login="$(printf '\033[1;36mtl auth login\033[0m')"
  tx_init="$(printf '\033[1;36mtl init\033[0m')"
fi
echo ""
printf '\tRun %s to login.\n' "$tx_auth_login"
printf '\tRun %s in each repo to initialize tx.\n' "$tx_init"
"$install_dir/tx" version
