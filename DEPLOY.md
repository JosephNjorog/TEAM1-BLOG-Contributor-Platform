# Deploying Team1 Blog Contributor Platform

Three deployable units, one shared Postgres database (Neon):

| Service | Path | Image | Listens on |
|---|---|---|---|
| API | `backend/` | `backend/Dockerfile` | `:8080` |
| Contributor/Moderator/Designer/Publisher app | `frontend/` | `frontend/Dockerfile` | `:8080` (static) |
| Super Admin app | `super-admin/` | `super-admin/Dockerfile` | `:8080` (static) |

All three are plain Dockerfiles, so they deploy as-is to Railway, Fly.io, or any VPS that can run a container. Build the two frontend images from the **repo root** (they depend on `packages/shared`):

```bash
docker build -f backend/Dockerfile -t team1-api ./backend
docker build -f frontend/Dockerfile -t team1-frontend .
docker build -f super-admin/Dockerfile -t team1-super-admin .
```

## 1. Database

Point `DATABASE_URL` at your Neon connection string with `sslmode=require`. Migrations run automatically on every backend boot (`internal/db/migrate.go`) - there's no separate migration step to run.

## 2. Choosing a topology

The two frontend apps call the API via `apiRequest`, which uses a **relative** `/api/v1` path by default. That only works if each frontend and the API share an origin (e.g. a reverse proxy in front of all three routes `/api/*` to the backend). Two options:

- **Same-origin per app (recommended, no extra config):** put a reverse proxy (Railway/Fly's built-in routing, or Caddy/nginx) in front of each frontend that forwards `/api/*` to the backend service. Leave `VITE_API_BASE_URL` unset.
- **Separate domains:** if the API is genuinely on its own domain (e.g. `api.team1blog.com`) and each frontend can't proxy to it, build the frontend images with `--build-arg VITE_API_BASE_URL=https://api.team1blog.com/api/v1`. The WebSocket notification socket derives its host from this same value, so it doesn't need separate configuration.

Either way, the backend's `CORS_ORIGINS` must list both frontend origins exactly (scheme + host, no trailing slash) or the browser will block requests.

## 3. Required backend environment variables

```bash
# Real secret - generate one, don't reuse the dev placeholder:
openssl rand -base64 48
```

| Variable | Notes |
|---|---|
| `DATABASE_URL` | Neon connection string, `sslmode=require` |
| `JWT_SECRET` | Generate with the command above. Rotating it invalidates every active session. |
| `CORS_ORIGINS` | Both frontend origins, comma-separated |
| `FRONTEND_URL` | The contributor/moderator/designer/publisher app's URL - used in invite/notification email links |
| `ADMIN_APP_URL` | The Super Admin app's URL - Super Admin invites and admin-facing emails link here instead |
| `PUBLIC_API_URL` | This API's own public URL (only matters if Cloudinary is still in mock mode - see below) |

Everything else (`RESEND_API_KEY`, `CLOUDINARY_*`, `AVALANCHE_*`, `SUBSTACK_PUBLICATION_URL`, `RESEND_WEBHOOK_SECRET`) is optional - each integration falls back to a mock implementation when its credentials are blank, so the platform runs end-to-end before any of them are configured. See `backend/.env.example` for the full list and what each one unlocks.

**Before going live for real**, set at minimum: `RESEND_API_KEY` (real emails), `CLOUDINARY_*` (real image hosting), and decide whether `AVALANCHE_TREASURY_PRIVATE_KEY`/`AVALANCHE_USDC_CONTRACT_ADDRESS` are ready - until they are, payment release stays simulated, which is safe but not real money movement.

## 4. Frontend build args

Both `frontend/Dockerfile` and `super-admin/Dockerfile` accept build args (not runtime env vars, since Vite inlines them at build time):

| App | Build arg | Purpose |
|---|---|---|
| `frontend` | `VITE_API_BASE_URL` | Override if not same-origin with the API |
| `frontend` | `VITE_ADMIN_APP_URL` | Where a Super Admin who lands here gets redirected |
| `super-admin` | `VITE_API_BASE_URL` | Override if not same-origin with the API |
| `super-admin` | `VITE_FRONTEND_URL` | Where a non-Super-Admin who lands here gets redirected |

## 5. Health check

`GET /healthz` on the API returns `200 ok` once the DB connection and migrations succeed - use it as the platform's readiness probe.

## 6. Post-deploy bootstrap

There's no public registration - the very first Super Admin account has to be created directly, since invites can only be sent by an existing Super Admin. `cmd/seed` creates one account per role (including the Super Admin) with a known password, which is convenient for the very first deploy but not something to leave as-is:

```bash
DATABASE_URL=... go run ./cmd/seed
```

Log in as the seeded Super Admin, immediately change its password (or invite a real Super Admin and deactivate the seeded one), and deactivate the other four seeded accounts from the Users & Invites / Contributors pages once real invites have gone out.
