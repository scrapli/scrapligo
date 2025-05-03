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

# clean old assets
rm assets/lib/aarch64-linux/*
rm assets/lib/aarch64-macos/*
rm assets/lib/x86_64-linux-gnu/*
rm assets/lib/x86_64-linux-musl/*
rm assets/lib/x86_64-macos/*

# revert the naming shenanigans libscrapli does to appease gh release assets
mv .tmp/libscrapli-aarch64-linux.so.* assets/lib/aarch64-linux/libscrapli.so.${LIBSCRAPLI_TAG#v}
mv .tmp/libscrapli-aarch64-macos.dylib.* assets/lib/aarch64-macos/libscrapli.${LIBSCRAPLI_TAG#v}.dylib
mv .tmp/libscrapli-x86_64-linux-gnu.so.* assets/lib/x86_64-linux-gnu/libscrapli.so.${LIBSCRAPLI_TAG#v}
mv .tmp/libscrapli-x86_64-linux-musl.so.* assets/lib/x86_64-linux-musl/libscrapli.so.${LIBSCRAPLI_TAG#v}
mv .tmp/libscrapli-x86_64-macos.dylib.* assets/lib/x86_64-macos/libscrapli.${LIBSCRAPLI_TAG#v}.dylib

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

sed -i.bak 's|var LibScrapliVersion = ".*"|var LibScrapliVersion = "'"${LIBSCRAPLI_TAG#v}"'"|' \
  constants/versions.go && rm -f constants/versions.go.bak

rm -rf .tmp
