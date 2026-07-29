#!/bin/sh
set -eu

base_url="${TOTALITY_INSTALL_BASE_URL:-https://download.totality.sh}"
install_dir="${TOTALITY_INSTALL_DIR:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d 2>/dev/null || mktemp -d -t totality-install)"

fail() {
  echo "tl install: $*" >&2
  exit 1
}

usage() {
  cat <<EOF
Totality installer

Default install: tl CLI, tl-mcp, shell completions.
Repo git hooks are installed later by tl init.

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
archive="tl_${os}_${arch}.tar.gz"
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
test -x "$tmp_dir/tl/bin/tl" || fail "archive is missing tl"
test -x "$tmp_dir/tl/bin/tl-mcp" || fail "archive is missing tl-mcp"

mkdir -p "$install_dir"
install -m 755 "$tmp_dir/tl/bin/tl" "$install_dir/tl"
install -m 755 "$tmp_dir/tl/bin/tl-mcp" "$install_dir/tl-mcp"
ln -sf tl "$install_dir/tlr"

if [ -f "$tmp_dir/tl/completions/tl.bash" ]; then
  mkdir -p "$HOME/.local/share/bash-completion/completions"
  install -m 644 "$tmp_dir/tl/completions/tl.bash" "$HOME/.local/share/bash-completion/completions/tl"
fi
if [ -f "$tmp_dir/tl/completions/_tl" ]; then
  mkdir -p "$HOME/.zfunc"
  install -m 644 "$tmp_dir/tl/completions/_tl" "$HOME/.zfunc/_tl"
fi

echo "Installed tl to $install_dir/tl"
echo "Installed MCP and aliases"
if ! command -v tl >/dev/null 2>&1; then
  echo "Add $install_dir to PATH before running tl."
fi
# Cyan ANSI 6 + bold matches tl version / logo (internal/cli/logo.go).
tl_auth_login="tl auth login"
tl_init="tl init"
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  tl_auth_login="$(printf '\033[1;36mtl auth login\033[0m')"
  tl_init="$(printf '\033[1;36mtl init\033[0m')"
fi
echo ""
printf '\tRun %s to login.\n' "$tl_auth_login"
printf '\tRun %s in each repo to initialize tl.\n' "$tl_init"
"$install_dir/tl" version
