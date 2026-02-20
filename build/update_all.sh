#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."

if [[ -n "$(git status --porcelain)" ]]; then
    echo "dirty dirty developer, start from a clean repo plz kthxbye"
    exit 1
fi

LATEST_LIBSCRAPLI_TAGS=(
    $(
        {
            curl -sf \
                "https://api.github.com/repos/scrapli/libscrapli/tags" 2>/dev/null ||
                echo '[]'
        } |
            jq -r '.[].name' |
            sort -Vr |
            head -n 5
    )
)
CURRENT_LIBSCRAPLI_TAG=$(
    grep -Eo 'var LibScrapliVersion = "[^"]+"' constants/versions.go | cut -d '"' -f2
)

LATEST_DEFINITION_TAGS=(
    $(
        {
            curl -sf \
                "https://api.github.com/repos/scrapli/scrapli_definitions/tags" 2>/dev/null ||
                echo '[]'
        } |
            jq -r '.[].name' |
            sort -Vr |
            head -n 5
    )
)
CURRENT_DEFINITION_TAG=$(
    grep -Eo 'var ScrapliDefinitionsVersion = "[^"]+"' constants/versions.go | cut -d '"' -f2
)

CURRENT_SCRAPLIGO_VERSION=$(
    grep -Eo 'var Version = "[^"]+"' constants/versions.go | cut -d '"' -f2
)

echo "current libscrapli tag: $CURRENT_LIBSCRAPLI_TAG"
if [[ ${#LATEST_LIBSCRAPLI_TAGS[@]} -gt 0 ]]; then
    echo "latest libscrapli tags:"
    for tag in "${LATEST_LIBSCRAPLI_TAGS[@]}"; do
        echo "  - $tag"
    done
else
    echo "no libscrapli tags found"
fi

TARGET_LIBSCRAPLI_TAG=""
while true; do
    read -p "enter new libscrapli tag (or press enter for current): " input
    if [[ -z "$input" ]]; then
        echo "keeping current libscrapli tag"
        echo
        break
    elif [[ ${#LATEST_LIBSCRAPLI_TAGS[@]} -gt 0 && " ${LATEST_LIBSCRAPLI_TAGS[*]} " =~ " $input " ]]; then
        TARGET_LIBSCRAPLI_TAG="$input"
        break
    else
        if [[ "$input" =~ ^[0-9a-f]{7,40}$ ]]; then
            echo "using commit hash $input"
            echo
            TARGET_DEFINITION_TAG="$input"
            break
        else
            echo "invalid tag or hash. pick a valid value dork"
        fi
    fi
done

echo "current scrapli-definitions tag: $CURRENT_DEFINITION_TAG"
if [[ ${#LATEST_DEFINITION_TAGS[@]} -gt 0 ]]; then
    echo "latest scrapli-definitions tags:"
    for tag in "${LATEST_DEFINITION_TAGS[@]}"; do
        echo "  - $tag"
    done
else
    echo "no scrapli-definitions tags found"
fi

TARGET_DEFINITION_TAG=""
while true; do
    read -p "enter new scrapli-definitions tag (or hash) (or press enter for current): " input
    if [[ -z "$input" ]]; then
        echo "keeping current scrapli-definitions tag"
        echo
        break
    elif [[ ${#LATEST_DEFINITION_TAGS[@]} -gt 0 && " ${LATEST_DEFINITION_TAGS[*]} " =~ " $input " ]]; then
        TARGET_DEFINITION_TAG="$input"
        break
    else
        if [[ "$input" =~ ^[0-9a-f]{7,40}$ ]]; then
            echo "using commit hash $input"
            echo
            TARGET_DEFINITION_TAG="$input"
            break
        else
            echo "invalid tag or hash, pick something valid loser"
        fi
    fi
done

if [[ -n "$TARGET_LIBSCRAPLI_TAG" ]]; then
    echo "updating libscrapli tag to: ${TARGET_LIBSCRAPLI_TAG}"
    ./build/update_libscrapli_tag.sh "$TARGET_LIBSCRAPLI_TAG"
fi

if [[ -n "$TARGET_DEFINITION_TAG" ]]; then
    echo "updating scrapli-definitions to: ${TARGET_DEFINITION_TAG}"
    ./build/update_scrapli_definitions.sh "$TARGET_DEFINITION_TAG"
fi

echo

CHANGES=$(git diff -- constants/versions.go assets/definitions)

if [[ -z "$CHANGES" ]]; then
    echo "no changes to commit, exiting..."
    exit 0
fi

echo "$CHANGES"
read -p "looks good? (y/n): " confirm

if [[ "$confirm" == [yY] ]]; then
    rm constants/versions.go.bak
else
    echo "restoring..."
    mv constants/versions.go.bak constants/versions.go || true
fi
