# Supabase (managed PostgreSQL) — setup notes

Supabase is fully managed PostgreSQL, so the Go backend runs against it with only
connection-string changes. Use this while a dedicated Postgres is not yet
provisioned ("sementara bisa pakai Supabase").

## 1. Point the backend at Supabase

In `apps/api/.env` (the backend reads plain `os.Getenv`, no `.env` autoload — export
it via your process manager / `air` / compose), set:

```env
DB_HOST=db.<project-ref>.supabase.co     # or the pooler host below
DB_PORT=5432                             # 6543 when using the transaction pooler
DB_USER=postgres                         # or postgres.<project-ref> for the pooler
DB_PASSWORD=<your-db-password>
DB_NAME=postgres
DB_SSLMODE=require                        # Supabase requires TLS
```

The backend builds its DSN in `apps/api/internal/config/config.go`
(`config.GetDSN()` → `postgres://user:pass@host:port/db?sslmode=...&timezone=UTC`);
`sslmode=require` is honored as-is.

### Direct vs pooler connection

- **Direct** (`db.<ref>.supabase.co:5432`): full session features, fewer max
  connections. Fine for a single backend instance; keep `DB_MAX_CONNS` modest.
- **Transaction pooler** (`<region>.pooler.supabase.com:6543`, user
  `postgres.<ref>`): needed for many short-lived connections / serverless. Note
  pgx prepared-statement caching can conflict with transaction pooling — prefer the
  direct connection for the backend unless you hit connection limits.

## 2. Redis is still required (Supabase does not provide it)

The backend uses Redis for the RBAC permission cache
(`apps/api/internal/shared/authz`, keys `perms:user:<id>:company:<id>`), **and** the
MPP queue engine needs Redis for per-service number counters, active-queue state,
and pub/sub fan-out to TV/loket/kiosk clients. Supabase has no Redis offering, so:

- **Local:** use the `redis` service in the root `docker-compose.yml`.
- **Hosted:** provision a separate Redis (e.g. Upstash) and set `REDIS_HOST` /
  `REDIS_PORT` / `REDIS_PASSWORD` accordingly.

## 3. Migrations & schemas

`golang-migrate` runs against Supabase unchanged (`make -C apps/api db-setup`). The
skeleton keeps master/RBAC tables under the `core` schema; the MPP domain tables live
under a dedicated `mpp` schema (see `docs/02-domain/`). Both are created by the
migrations — no manual Supabase SQL editor steps required.

> ⚠️ Do not manage these tables through the Supabase dashboard UI; keep schema changes
> in versioned migrations so environments stay reproducible.

## 4. Object storage

The skeleton ships a GCS storage adapter (`apps/api/pkg/storage`). Supabase Storage
is S3-compatible; if you prefer it over GCS for uploaded documents / QR / media,
add an S3-compatible adapter alongside the GCS one rather than replacing it.
