#!/usr/bin/env bash
set -e

if [[ "$1" == "-Help" || "$1" == "--help" || "$1" == "-h" ]]; then
    echo "Usage: ./scripts/commit.sh [--skip-tests]"
    echo "Enforces go vet, go test and Conventional Commits format."
    exit 0
fi

echo "==> Verifying repository state..."
if [[ "$1" != "--skip-tests" ]]; then
    echo "--> Running go vet..."
    go vet ./...

    echo "--> Running tests..."
    go test -v ./...
fi

VALID_TYPES=("feat" "fix" "refactor" "test" "docs" "chore" "ci" "perf" "style")

echo ""
echo "=== Conventional Commit Assistant ==="
echo "Allowed types: ${VALID_TYPES[*]}"

while true; do
    read -rp "Enter commit type: " TYPE
    TYPE=$(echo "$TYPE" | tr '[:upper:]' '[:lower:]' | xargs)
    for valid in "${VALID_TYPES[@]}"; do
        if [[ "$valid" == "$TYPE" ]]; then
            break 2
        fi
    done
    echo "Invalid type. Must be one of: ${VALID_TYPES[*]}"
done

read -rp "Enter scope (e.g., core, parser, repo) [optional]: " SCOPE
SCOPE=$(echo "$SCOPE" | tr '[:upper:]' '[:lower:]' | xargs)

while true; do
    read -rp "Enter commit message (imperative, lower-case start): " MESSAGE
    MESSAGE=$(echo "$MESSAGE" | xargs)
    if [[ -n "$MESSAGE" ]]; then
        break
    fi
    echo "Message cannot be empty."
done

if [[ -n "$SCOPE" ]]; then
    COMMIT_MSG="${TYPE}(${SCOPE}): ${MESSAGE}"
else
    COMMIT_MSG="${TYPE}: ${MESSAGE}"
fi

echo ""
echo "Commit message: '$COMMIT_MSG'"
read -rp "Proceed with commit? (y/n): " CONFIRM
if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
    git commit -m "$COMMIT_MSG"
    echo "Successfully committed: $COMMIT_MSG"
else
    echo "Commit cancelled."
fi
