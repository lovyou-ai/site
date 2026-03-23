#!/bin/bash
# ship.sh — Generate, build, test, deploy, commit, push. One command.
#
# Usage:
#   ./ship.sh "iter 101: description here"
#
# Runs from the site repo root.
set -euo pipefail

MSG="${1:?Usage: ./ship.sh \"commit message\"}"

echo "=== GENERATE ==="
/c/Users/matt_/go/bin/templ generate

echo "=== BUILD ==="
go.exe build -buildvcs=false ./...

echo "=== TEST ==="
go.exe test -buildvcs=false ./...

echo "=== DEPLOY ==="
/c/Users/matt_/.fly/bin/flyctl deploy --remote-only

echo "=== COMMIT ==="
git add -A
git commit -m "$MSG

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>"

echo "=== PUSH ==="
git push origin main

echo "=== SHIPPED ==="
