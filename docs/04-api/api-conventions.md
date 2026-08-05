# API Conventions — MPP _(EN)_

Inherits the `venturo-skeleton-go` conventions; MPP endpoints follow the same shape.

## Base URL & versioning

- Base: `NEXT_PUBLIC_API_URL` (e.g. `http://localhost:8080`).
- Versioned paths, e.g. `/mpp/v1/...` for the queue domain, `/core/v1/...` for skeleton
  modules (auth/user/role/company).

## Response envelope

All responses use the skeleton's envelope (the FE `ky` client unwraps it):

```json
{
  "data": { },
  "message": "human-readable status",
  "meta": { "page": 1, "per_page": 20, "total": 137 },
  "errors": [ { "field": "tanggal", "message": "quota full" } ]
}
```

- `data`: payload (object or array); `null` on pure errors.
- `meta`: pagination/extra; present on list endpoints.
- `errors`: array of field errors on 4xx validation failures.
- Non-2xx map to a single typed `ApiError` on the frontend.

## Authentication

- **JWT** (`Authorization: Bearer <token>`) for users. Claims carry `UserID`, `Email`,
  `CompanyID` (+ super-admin flag). Obtained via `/core/v1/auth/login`.
- **API-key** (`X-API-Key: <key>`) for unattended devices (kiosk/TV) with pre-scoped
  permissions. The auth middleware checks `X-API-Key` first, else the Bearer token.

## Tenancy headers

- `X-Company-Slug`: tenant (the MPP building) — sent by the public web on `/…` calls.
- `X-Client-Slug`: whitelabel client (reserved; translation overrides).

## Authorization

Server-side RBAC via `resource:action` / `resource.subresource:action` permission
strings, levels `viewer|editor|admin`, cached in Redis. See
[`../06-security/rbac-matrix.md`](../06-security/rbac-matrix.md). Every mutating MPP
endpoint declares a required permission.

## Pagination, filtering, sorting

- Query params: `page`, `per_page`, `sort`, `order`, plus resource filters
  (`instansi_id`, `jenis_layanan_id`, `status`, `date_from`, `date_to`).
- List responses include `meta` with totals.

## Errors

| HTTP | Meaning |
|------|---------|
| 400 | validation error (see `errors[]`) |
| 401 | missing/invalid/expired credentials |
| 403 | authenticated but lacks permission |
| 404 | not found / out of tenant scope |
| 409 | conflict (e.g. quota full, illegal state transition, duplicate number) |
| 429 | rate limited (public registration) |
| 500 | server error |

## Idempotency & concurrency

- Number allocation and quota consumption are **atomic** (Redis INCR / DB constraints).
- State-changing queue actions validate the current state and reject **illegal
  transitions** with `409`.
- Check-in with an already-used/expired token returns `409`/`410`.

## Realtime

Queue changes are also delivered over WebSocket — see
[`websocket-events.md`](./websocket-events.md). REST remains the source of truth; the
socket is a change feed for TVs/loket/admin.

## Timezone

All timestamps are **UTC** in requests/responses; clients localize to WIB/WITA/WIT.
