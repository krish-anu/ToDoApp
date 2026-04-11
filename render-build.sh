#!/usr/bin/env bash
set -euo pipefail

# Root build script for Render (or any CI).
# Usage on Render: bash render-build.sh
# This builds the frontend (Vite), copies dist into backend/web/dist,
# and then builds the Go binary so go:embed can include static assets.

echo "==> Building frontend..."
cd frontend
# Use npm ci if you want deterministic installs in CI; fallback to npm install if package-lock missing
if [ -f package-lock.json ]; then
  npm ci
else
  npm install
fi
npm run build

cd ..

echo "==> Copying frontend/dist to backend/web/dist..."
rm -rf backend/web/dist
mkdir -p backend/web
cp -R frontend/dist backend/web/dist

echo "==> Building Go binary..."
# Use the same flags you used locally; adjust as needed for your provider
cd backend
go build -tags netgo -ldflags "-s -w" -o ../app
cd ..

echo "==> Build complete: ./app"
