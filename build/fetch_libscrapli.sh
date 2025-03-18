#!/usr/bin/env bash
set -euxo pipefail

LIBSCRAPLI_TAG=$1

LIBSCRAPLI_RELEASE="https://api.github.com/repos/scrapli/libscrapli/releases/tags/$LIBSCRAPLI_TAG"
ASSETS_URL=$(curl -s "$LIBSCRAPLI_RELEASE" | jq -r '.assets[].browser_download_url')

if [ -z "$ASSETS_URL" ]; then
    echo "no assets for release $LIBSCRAPLI_TAG"
    exit 1
fi

mkdir .tmp || true

for url in $ASSETS_URL; do
    filename=$(basename "$url")
    curl -L -o ".tmp/$filename" "$url"
done

# get the "lib version" (doesnt include any tag alpha/post/etc. suffix)
LIB_VERSION=$(find .tmp | grep -E -o "\d+\.\d+\.\d+" | head -1)

# revert the naming shenanigans libscrapli does to appease gh release assets
mv .tmp/libscrapli-aarch64-linux.so.* assets/lib/aarch64-linux/libscrapli.so.$LIB_VERSION
mv .tmp/libscrapli-aarch64-macos.dylib.* assets/lib/aarch64-macos/libscrapli.$LIB_VERSION.dylib
mv .tmp/libscrapli-x86_64-linux-gnu.so.* assets/lib/x86_64-linux-gnu/libscrapli.so.$LIB_VERSION
mv .tmp/libscrapli-x86_64-linux-musl.so.* assets/lib/x86_64-linux-musl/libscrapli.so.$LIB_VERSION
mv .tmp/libscrapli-x86_64-macos.dylib.* assets/lib/x86_64-macos/libscrapli.$LIB_VERSION.dylib

while IFS='  ' read -r checksum filepath; do
    new_path=$(echo "$filepath" | sed 's/^zig-out\//assets\/lib\//')
    calculated_checksum=$(sha256sum "$new_path" | awk '{print $1}')

    if [ "$checksum" != "$calculated_checksum" ]; then
        echo "Checksum mismatch for $filename!"
        echo "Expected: $checksum"
        echo "Found: $calculated_checksum"
        exit 1
    fi
done < ".tmp/checksums.txt"

sed -i.bak 's|var LibScrapliVersion = "0.0.1"|var LibScrapliVersion = "'"$LIBSCRAPLI_TAG"'"|' \
  constants/versions.go && rm -f constants/versions.go.bak

rm -rf .tmp