#!/bin/bash
set -euo pipefail

# lil helper to highlight all the places there are versions of things that may need to be updated
# during release or just maintenance stuff. ignores action versions (dependabots problem) and the
# actual versions in go.mod (other than go version itself)
PURPLE=$(printf '\033[1;35m')
CYAN=$(printf '\033[1;36m')
NC=$(printf '\033[0m')

PATTERN_SEMVER='\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?'
PATTERN_CALVER='\d{4}\.\d{1,2}\.\d{1,2}'
PATTERN_HASH='[a-f0-9]{7,}'
PATTERN_GOVER='1\.\d{1,2}'

PATTERN_VERSIONS="${PATTERN_SEMVER}|${PATTERN_CALVER}|${PATTERN_HASH}|${PATTERN_GOVER}"

highlight_version() {
    echo -e "\n${CYAN}=============== $1 :: $3${NC}"
    grep -E "$2" "$1" --color=never | grep -E "${PATTERN_VERSIONS}"
}

# file :: re to match the line :: nice name to print
locations=(
    "constants/versions.go     ^var\\sVersion =                         scrapligo"
    "constants/versions.go     ^var\\sLibScrapliVersion\\s=             libscrapli"
    "constants/versions.go     ^var\\sScrapliDefinitionsVersion\\s=     definitions"
    "go.mod                    ^go\\s                                   go"
    ".github/vars.env          GO_VERSION=                              ci go"
    ".github/vars.env          GCI_VERSION=                             ci gci"
    ".github/vars.env          GOFUMPT_VERSION=                         ci gofumpt"
    ".github/vars.env          GOLANGCI_LINT_VERSION=                   ci golangci-lint"
    ".github/vars.env          GOLINES_VERSION=                         ci golines"
    ".github/vars.env          GOTESTSUM_VERSION=                       ci gotestsum"
    "Makefile                  ghcr.io/scrapli/                         local/ci clab setup"
)

for entry in "${locations[@]}"; do
    read -r file pattern label <<<"$entry"
    highlight_version "$file" "$pattern" "$label"
done
