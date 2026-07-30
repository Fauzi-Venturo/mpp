# REST Endpoint Catalog — MPP _(EN)_

Indicative catalog for the `mpp` domain (prefix `/mpp/v1`). Auth/permission per
[`api-conventions.md`](./api-conventions.md) and
[`../06-security/rbac-matrix.md`](../06-security/rbac-matrix.md). Skeleton auth/user
endpoints live under `/core/v1` and are reused as-is.

## Catalog (public read + admin write)

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| GET | `/mpp/v1/instansi` | List agencies | public |
| GET | `/mpp/v1/instansi/{id}` | Agency detail | public |
| POST/PUT/DELETE | `/mpp/v1/instansi[/{id}]` | Manage agencies | admin |
| GET | `/mpp/v1/instansi/{id}/layanan` | Services of an agency (+ requirements, duration) | public |
| POST/PUT/DELETE | `/mpp/v1/layanan[/{id}]` | Manage services | admin |
| POST/PUT/DELETE | `/mpp/v1/layanan/{id}/syarat[/{sid}]` | Manage document requirements | admin |
| GET/POST/PUT/DELETE | `/mpp/v1/loket[/{id}]` | Manage lokets + service mapping + status | admin/supervisor |

## Quota & availability

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| GET | `/mpp/v1/availability?instansi_id&layanan_id&date` | Remaining quota for a date | public |
| GET/POST/PUT | `/mpp/v1/kuota` | Manage quotas (per date/agency/service) | admin |

## Registration & booking

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| POST | `/mpp/v1/booking` | Create booking (web) → returns QR token | public (rate-limited) |
| GET | `/mpp/v1/booking/{id}` | Booking status | public (token/owner) |
| POST | `/mpp/v1/booking/{id}/cancel` | Cancel booking | public (owner)/admin |
| POST | `/mpp/v1/wa/webhook` | WhatsApp inbound webhook (AI agent) | signed webhook |
| POST | `/mpp/v1/walkin` | Walk-in registration (kiosk) → number + print payload | device (API-key) |

## Check-in

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| POST | `/mpp/v1/checkin` | Check-in with QR token → antrian (WAITING) + ticket | device (API-key) |

## Queue operations (loket app)

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| POST | `/mpp/v1/loket/{id}/session` | Open/close operator session | petugas |
| GET | `/mpp/v1/queue?layanan_id` | Current waiting stream for a service | staff |
| POST | `/mpp/v1/queue/next` `{loket_id}` | Call next (idle-longest, mode-aware) → CALLED | petugas |
| POST | `/mpp/v1/antrian/{id}/recall` | Recall (call_count ≤ 3) | petugas |
| POST | `/mpp/v1/antrian/{id}/start` | Start serving → SERVING | petugas |
| POST | `/mpp/v1/antrian/{id}/done` | Finish → DONE | petugas |
| POST | `/mpp/v1/antrian/{id}/skip` | No-show → SKIPPED | petugas |
| POST | `/mpp/v1/antrian/{id}/hold` / `/resume` | Hold / resume | petugas |
| POST | `/mpp/v1/antrian/{id}/transfer` `{target_loket|layanan}` | Transfer → TRANSFERRED | petugas |
| POST | `/mpp/v1/antrian/{id}/second-service` `{layanan_id}` | Trigger second service → QUEUED_NEXT | petugas |

## Front office

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| GET | `/mpp/v1/fo/queue` | Applicants awaiting verification (+ checklist) | front office |
| POST | `/mpp/v1/fo/{antrian_id}/verify` `{result, checklist, notes}` | Record verification | front office |

## Display (TV)

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| GET | `/mpp/v1/display?instansi_id` | Snapshot for a TV (current called + next up) | device (API-key) |
| GET | `/mpp/v1/display/config` | TV layout / running text config | device |

## Admin, monitoring & reporting

| Method | Path | Purpose | Access |
|--------|------|---------|--------|
| GET | `/mpp/v1/monitoring` | Real-time snapshot across agencies/lokets | admin/supervisor |
| POST | `/mpp/v1/admin/reset` `{instansi_id?}` | Manual reset (audited) | admin/supervisor |
| POST | `/mpp/v1/admin/broadcast` | Push info/running text to TVs | admin/supervisor |
| GET | `/mpp/v1/reports/summary?date_from&date_to&instansi_id&layanan_id` | Served, no-show, avg/p90 wait & serve | admin/supervisor |
| GET | `/mpp/v1/reports/export?...` | CSV/Excel export | admin/supervisor |
| GET | `/mpp/v1/config` / PUT | System & per-agency config (mode, hours, number format, TTS text) | admin |
| GET | `/mpp/v1/audit` | Audit log of sensitive actions | admin |

> This catalog is a design contract, not the final OpenAPI spec. Finalize request/
> response DTOs and the OpenAPI document in `packages/api-contract` during the API-design
> phase; keep it in sync with the FE `src/lib/api/endpoints.ts`.
