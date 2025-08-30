#!/bin/bash
set -euo pipefail

LIBSCRAPLI_TAG="${1:-}"

if [[ -z "$LIBSCRAPLI_TAG" ]]; then
    echo "error: libscrapli tag must be set"
    exit 1
fi

LIBSCRAPLI_TAG="${LIBSCRAPLI_TAG#v}"

sed -i.bak -E "s|(var LibScrapliVersion = )(.*)|\1\"${LIBSCRAPLI_TAG}\"|g" constants/versions.go
