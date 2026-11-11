#!/usr/bin/env bash
set -e

if [[ "$1" == "-Help" || "$1" == "--help" || "$1" == "-h" ]]; then
    echo "Usage: ./scripts/run.sh"
    echo "Runs GramatiGo in development mode."
    exit 0
fi

if command -v wails &> /dev/null; then
    echo "==> Starting Wails development server..."
    wails dev
else
    echo "==> Wails CLI not found on PATH. Running standard Go app..."
    go run main.go
fi
