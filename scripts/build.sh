#!/usr/bin/env bash
set -e

if [[ "$1" == "-Help" || "$1" == "--help" || "$1" == "-h" ]]; then
    echo "Usage: ./scripts/build.sh [platform]"
    echo "Default platform: linux/amd64"
    exit 0
fi

PLATFORM="${1:-linux/amd64}"
OUTPUT_DIR="build/bin"
mkdir -p "$OUTPUT_DIR"

LDFLAGS="-s -w"

echo "==> Building GramatiGo for $PLATFORM..."
if command -v wails &> /dev/null; then
    wails build -clean -platform "$PLATFORM" -ldflags "$LDFLAGS"
else
    echo "==> Wails CLI not found. Running go build..."
    EXT=""
    if [[ "$PLATFORM" == windows* ]]; then
        EXT=".exe"
    fi
    go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/corrector$EXT" .
fi

echo "==> Build finished successfully."
