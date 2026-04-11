Render / generic deployment notes

If your deployment fails with `pattern web/dist/**: no matching files found` or `sendfile: file backend/web/dist/index.html not found`, it means the frontend build artifacts weren't copied into `backend/web/dist` before `go build` ran.

Recommended: run the frontend build before `go build`.

Render Web Service (recommended) settings for this repo layout:

- Root Directory:
  - (leave empty, use repository root)
- Build Command:
  - `bash render-build.sh`
- Start Command:
  - `./app`

Required Environment Variables:

- `MONGODB_URI`

Optional Environment Variables:

- `OPENAI_API_KEY` (required only for `/api/workflows/task-from-text`)
- `ALLOW_ORIGINS` (set your frontend origin if needed)
- `OPENAI_MODEL`
- `OPENAI_BASE_URL`
- `WORKFLOW_DEFAULT_USER_ID`
- `WORKFLOW_DEFAULT_TIMEZONE`

Notes:

- Do not set `PORT` manually on Render; Render injects it automatically.
- Keep `render-build.sh` at repo root and executable via `bash`.

If you prefer Docker, use a multi-stage Dockerfile to build the frontend and then copy `frontend/dist` into `backend/web/dist` before building the Go binary.

Alternative: remove `dist` from `.gitignore` and commit built assets (not recommended for most workflows).
