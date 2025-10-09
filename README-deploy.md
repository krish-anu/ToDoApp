Render / generic deployment notes

If your deployment fails with `pattern client/dist/**: no matching files found` or `sendfile: file client/dist/index.html not found`, it means the frontend build artifacts weren't present when `go build` ran (this project uses `go:embed` or expects `client/dist` at runtime).

Recommended: run the frontend build before `go build`.

Example (Render):
- Build Command:
  - bash render-build.sh
- Start Command:
  - ./app

If you prefer Docker, use a multi-stage Dockerfile to build the frontend and then copy `client/dist` into the final image before building the Go binary.

Alternative: remove `dist` from `.gitignore` and commit built assets (not recommended for most workflows).
