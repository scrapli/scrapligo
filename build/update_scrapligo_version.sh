#!/bin/bash
set -euo pipefail

SCRAPLIGO_VERSION="${1:-}"

if [[ -z "$SCRAPLIGO_VERSION" ]]; then
    echo "error: scrapligo version must be set"
    exit 1
fi

sed -i.bak -E "s|(var Version = )(.*)|\1\"${SCRAPLIGO_VERSION}\"|g" constants/versions.go
