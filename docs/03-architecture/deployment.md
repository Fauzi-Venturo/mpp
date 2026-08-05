# Deployment — MPP _(EN)_

## Environments

| Env    | Backend | Frontend | Postgres | Redis |
|--------|---------|----------|----------|-------|
| Local  | `apps/api` via `make api-dev` or compose | `apps/web` via `make web-dev` or compose | compose `postgres` | compose `redis` |
| Staging/Prod | container (k8s) | container (k8s) | managed Postgres / **Supabase** | managed Redis (e.g. Upstash) |

## Local development

```bash
make bootstrap      # go mod download + yarn install + copy .env files
make up             # postgres + redis (infra only)
make db-setup       # migrations + seeders (core, then mpp)
make api-dev        # backend :8080 (air hot-reload)
make web-dev        # frontend :8002
# or the whole stack in containers:
make up-full
```

Compose defaults (`docker-compose.yml`): `postgres:16` (`mpp` db), `redis:7`, and
optional `api`/`web` under the `full` profile. Backend env is injected by compose
(`DB_HOST=postgres`, `REDIS_HOST=redis`); it reads plain `os.Getenv`.

## Backend container

The skeleton's multi-stage `Dockerfile` builds a static binary, bundles the
golang-migrate CLI + migrations + seeders, runs as non-root, `EXPOSE 80`, healthcheck
on `/health`. Its entrypoint runs migrations then starts the API; a worker CMD variant
runs background jobs (daily reset, expiry, reminders). In k8s, run the API and the
worker as **separate deployments** off the same image.

## Frontend container

The skeleton's `Dockerfile` builds Next standalone (`BUILD_STANDALONE=true`) and bakes
`NEXT_PUBLIC_*` at build time from `.env.prod`. **Consequence:** each environment needs
its own build with the correct `NEXT_PUBLIC_API_URL` / slugs. Only server-only vars
(`API_URL`, `REVALIDATE_TOKEN`) are read at runtime.

## Kubernetes (production)

- Backend: `deployments/kubernetes/` from the skeleton (adapt names/namespace).
  Deploy **api** + **worker** + a **WebSocket-capable ingress** (sticky sessions or a
  Redis-backed hub so any pod can serve any socket).
- Frontend: standalone Next server behind ingress; rebuild per environment.
- Externalize Postgres + Redis; never run stateful queue data in-pod.

## Supabase (temporary DB host)

See [`../../infra/supabase/README.md`](../../infra/supabase/README.md). Point backend
`DB_*` at Supabase (`DB_SSLMODE=require`, direct connection preferred for pgx); run the
same migrations. **Redis is still required separately** (Supabase has none) for the RBAC
cache and the queue engine.

## On-site devices

- **Kiosk**: locked-down browser (kiosk mode) to a Next route; USB QR scanner (keyboard
  emulation) + thermal printer (browser print or a local print agent).
- **TV mini PC**: one browser, three windows to the TV display route; **offline
  Indonesian TTS** + audio assets stored locally so voice calling survives internet
  loss. See [`../05-integrations/tts-voice-calling.md`](../05-integrations/tts-voice-calling.md).
- Devices authenticate with **scoped API-keys**, not user JWTs.

## Operational notes

- **Timezone:** backend enforces `TZ=UTC` (storage in UTC); UI converts to WIB/WITA/WIT.
- **Daily reset 00:00** is a worker job — ensure exactly one worker instance owns the
  schedule (leader election or a single worker deployment).
- **Backups:** managed Postgres/Supabase automated backups; Redis is rebuildable from
  Postgres, so treat it as a cache (persistence optional).
- **Health/observability:** `/health` for api; expose queue metrics (active numbers,
  avg wait/serve, per-loket throughput).
