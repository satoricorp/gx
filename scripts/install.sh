#!/bin/sh
set -eu

base_url="${GX_INSTALL_BASE_URL:-https://download.gx.run}"
install_dir="${GX_INSTALL_DIR:-$HOME/.local/bin}"
tmp_dir="$(mktemp -d 2>/dev/null || mktemp -d -t gx-install)"

cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

fail() {
  echo "gx install: $*" >&2
  exit 1
}

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
archive="gx_${os}_${arch}.tar.gz"
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
test -x "$tmp_dir/gx/bin/gx" || fail "archive is missing gx"
test -x "$tmp_dir/gx/bin/gx-mcp" || fail "archive is missing gx-mcp"

mkdir -p "$install_dir"
install -m 755 "$tmp_dir/gx/bin/gx" "$install_dir/gx"
install -m 755 "$tmp_dir/gx/bin/gx-mcp" "$install_dir/gx-mcp"
ln -sf gx "$install_dir/gxg"
ln -sf gx "$install_dir/gxs"

if [ -f "$tmp_dir/gx/completions/gx.bash" ]; then
  mkdir -p "$HOME/.local/share/bash-completion/completions"
  install -m 644 "$tmp_dir/gx/completions/gx.bash" "$HOME/.local/share/bash-completion/completions/gx"
fi
if [ -f "$tmp_dir/gx/completions/_gx" ]; then
  mkdir -p "$HOME/.zfunc"
  install -m 644 "$tmp_dir/gx/completions/_gx" "$HOME/.zfunc/_gx"
fi

echo "Installed gx to $install_dir/gx"
echo "Installed gxg to $install_dir/gxg"
echo "Installed gxs to $install_dir/gxs"
echo "Installed gx-mcp to $install_dir/gx-mcp"
if ! command -v gx >/dev/null 2>&1; then
  echo "Add $install_dir to PATH before running gx."
fi
"$install_dir/gx" version
