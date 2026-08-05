# RBAC Matrix — MPP _(EN)_

Maps MPP roles to the skeleton's permission model. The backend uses `resource:action`
(or `resource.subresource:action`) permission strings with levels `viewer|editor|admin`,
stored as JSONB on `core.roles`, cached in Redis. Roles are seeded (`is_system`) and
assigned via `core.user_roles`.

## Roles

| Role | Scope | Summary |
|------|-------|---------|
| `admin` | tenant-wide | full master data, config, users, all reports, monitoring |
| `supervisor` | assigned agency(ies) | monitor + limited ops (open/close loket, reset, broadcast), agency reports |
| `front_office` | tenant / assigned desk | document verification only |
| `petugas_loket` | assigned loket | queue operations at the loket |
| _device: kiosk_ | API-key | check-in + walk-in + print |
| _device: tv_ | API-key | display snapshot + subscribe |
| _citizen_ | none (public) | registration/check-in via public endpoints |

## Permission resources (MPP)

`mpp.instansi`, `mpp.layanan`, `mpp.loket`, `mpp.kuota`, `mpp.booking`, `mpp.checkin`,
`mpp.queue`, `mpp.antrian`, `mpp.fo`, `mpp.display`, `mpp.monitoring`, `mpp.report`,
`mpp.config`, `mpp.audit` (+ skeleton `core.user`, `core.role`).

## Matrix

Legend: **A** = admin/manage, **E** = editor/act, **V** = viewer/read, — = none.

| Resource → | admin | supervisor | front_office | petugas_loket | kiosk | tv |
|------------|:-----:|:----------:|:------------:|:-------------:|:-----:|:--:|
| `mpp.instansi` | A | V | V | V | — | — |
| `mpp.layanan` | A | V | V | V | V | — |
| `mpp.loket` | A | E¹ | — | V | — | — |
| `mpp.kuota` | A | V | — | — | — | — |
| `mpp.booking` | A | V | V | V | E² | — |
| `mpp.checkin` | A | — | — | — | E | — |
| `mpp.queue` | A | E¹ | V | E³ | — | V |
| `mpp.antrian` | A | E¹ | V | E³ | — | V |
| `mpp.fo` | A | V¹ | E | V | — | — |
| `mpp.display` | A | E¹ | — | — | — | V |
| `mpp.monitoring` | A | V¹ | — | — | — | — |
| `mpp.report` | A | V¹ | — | — | — | — |
| `mpp.config` | A | — | — | — | — | — |
| `mpp.audit` | A | — | — | — | — | — |
| `core.user` | A | — | — | — | — | — |
| `core.role` | A | — | — | — | — | — |

Footnotes:
1. **Scoped to assigned agency** — supervisor acts only within their instansi(s).
2. Kiosk creates walk-in bookings/numbers via device API-key (limited to registration).
3. Petugas acts only on **their loket's** queue items (call/recall/start/done/skip/hold/
   transfer/second-service).

## Enforcement

- Every mutating MPP endpoint declares the required permission; middleware checks it
  server-side (Redis-cached, DB fallback).
- **Agency/loket scoping** is enforced in addition to the permission level (a supervisor
  with `mpp.queue:editor` still cannot touch another agency's queue).
- Devices use **pre-scoped API-keys**, never user JWTs.
- Super-admin bypass follows the skeleton behavior; use sparingly.

## Seeding

Seed the four MPP roles with their JSONB permission sets under `seeders/mpp/` following
the skeleton's `core` seeding pattern (roles → users → assignments). Keep them
`is_system` so they aren't accidentally deleted.
