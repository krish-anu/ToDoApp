#!/usr/bin/env bash
set -euo pipefail

# Root build script for Render (or any CI).
# Usage on Render: bash render-build.sh
# This builds the client (Vite) and then the Go binary so that go:embed
# has files to include.

echo "==> Building client..."
cd client
# Use npm ci if you want deterministic installs in CI; fallback to npm install if package-lock missing
if [ -f package-lock.json ]; then
  npm ci
else
  npm install
fi
npm run build

cd ..

echo "==> Building Go binary..."
# Use the same flags you used locally; adjust as needed for your provider
go build -tags netgo -ldflags "-s -w" -o app

echo "==> Build complete: ./app"
