#!/usr/bin/env bash
# Downloads the prebuilt native archives this module needs for cgo: the
# zktf-sim-ffi archive (matching zktf-sim-version) plus the zktf-sdk core
# archive (matching zktf-sdk-version, same pin zktf-sdk-go's own
# scripts/fetch-native.sh uses) for the host (or requested) target triple.
#
# Usage:
#   scripts/fetch-native.sh                # fetch for GOHOSTOS/GOHOSTARCH
#   scripts/fetch-native.sh linux amd64    # fetch for an explicit GOOS/GOARCH
#
# Output: writes a sourceable .env file at the repo root (CGO_CFLAGS,
# CGO_LDFLAGS, LD_LIBRARY_PATH covering both archives) and echoes the same
# `export` lines to stdout.
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

sim_version_file="$repo_root/zktf-sim-version"
sdk_version_file="$repo_root/zktf-sdk-version"
for f in "$sim_version_file" "$sdk_version_file"; do
  if [[ ! -f "$f" ]]; then
    echo "fetch-native: missing $f" >&2
    exit 1
  fi
done
sim_version="$(tr -d '[:space:]' < "$sim_version_file")"
sdk_version="$(tr -d '[:space:]' < "$sdk_version_file")"

goos="${1:-$(go env GOHOSTOS 2>/dev/null || go env GOOS)}"
goarch="${2:-$(go env GOHOSTARCH 2>/dev/null || go env GOARCH)}"

case "$goos/$goarch" in
  linux/amd64)  triple=x86_64-unknown-linux-gnu ;;
  linux/arm64)  triple=aarch64-unknown-linux-gnu ;;
  darwin/arm64) triple=aarch64-apple-darwin ;;
  *)
    echo "fetch-native: unsupported GOOS/GOARCH combination: $goos/$goarch" >&2
    echo "fetch-native: supported: linux/amd64, linux/arm64, darwin/arm64" >&2
    exit 1
    ;;
esac

native_root="$repo_root/.zktf-native"

download_archive() {
  local gcs_bucket_path="$1" archive_name="$2" dest_dir="$3"

  if [[ -f "$dest_dir/.fetched" ]]; then
    echo "fetch-native: $archive_name already present at $dest_dir" >&2
    return
  fi

  mkdir -p "$dest_dir"
  local tarball="$native_root/$archive_name"
  local https_url="https://storage.googleapis.com/download.joinself.com/$gcs_bucket_path/$archive_name"
  local gcs_uri="gs://download.joinself.com/$gcs_bucket_path/$archive_name"

  local download_ok=0
  if command -v curl >/dev/null 2>&1; then
    echo "fetch-native: downloading via curl: $https_url" >&2
    if curl -fsSL -o "$tarball" "$https_url"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]] && command -v wget >/dev/null 2>&1; then
    echo "fetch-native: downloading via wget: $https_url" >&2
    if wget -q -O "$tarball" "$https_url"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]] && command -v gcloud >/dev/null 2>&1; then
    echo "fetch-native: downloading via gcloud storage cp: $gcs_uri" >&2
    if gcloud storage cp "$gcs_uri" "$tarball"; then
      download_ok=1
    fi
  fi
  if [[ "$download_ok" -eq 0 ]]; then
    echo "fetch-native: failed to download $archive_name via curl, wget, or gcloud storage cp" >&2
    exit 1
  fi

  tar -xzf "$tarball" -C "$dest_dir" --strip-components=1
  rm -f "$tarball"
  touch "$dest_dir/.fetched"
}

sim_dest_dir="$native_root/sim-$triple-$sim_version"
sdk_dest_dir="$native_root/sdk-$triple-$sdk_version"

download_archive "zktf-sim" "zktf-sim-$triple-$sim_version.tar.gz" "$sim_dest_dir"
download_archive "zktf-sdk" "zktf-sdk-$triple-$sdk_version.tar.gz" "$sdk_dest_dir"

env_file="$repo_root/.env"
{
  echo "CGO_CFLAGS=-I$sim_dest_dir -I$sdk_dest_dir"
  echo "CGO_LDFLAGS=-L$sim_dest_dir -L$sdk_dest_dir -lzktf_sim -lzktf_sdk"
  echo "LD_LIBRARY_PATH=$sim_dest_dir:$sdk_dest_dir"
} > "$env_file"

echo "fetch-native: wrote $env_file" >&2
echo "export CGO_CFLAGS=\"-I$sim_dest_dir -I$sdk_dest_dir\""
echo "export CGO_LDFLAGS=\"-L$sim_dest_dir -L$sdk_dest_dir -lzktf_sim -lzktf_sdk\""
echo "export LD_LIBRARY_PATH=\"$sim_dest_dir:$sdk_dest_dir\""
