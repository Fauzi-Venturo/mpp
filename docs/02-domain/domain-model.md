# Domain Model — MPP _(EN)_

The MPP domain lives under a dedicated Postgres schema **`mpp`**, alongside the
skeleton's `core` schema (users, roles, companies). This document describes the
entities and their relationships; the ERD is in [`erd.md`](./erd.md) and the queue
lifecycle in [`queue-state-machine.md`](./queue-state-machine.md).

## Tenancy mapping

The backend skeleton is multi-tenant (`core.companies`, company-scoped RBAC). For MPP,
**one MPP building = one company (tenant)**; **agencies (`instansi`) are a first-class
entity inside that tenant**, not separate companies. This keeps cross-agency reporting
and shared lokets/kiosks within one tenant boundary while reserving `company` for
future multi-building expansion.

## Core entities

### `instansi` (Agency)
An agency operating in the MPP. Owns services and lokets.
- `id`, `company_id` (tenant), `name`, `slug`, `prefix` (single letter used in queue
  numbers, unique per tenant), `description`, `logo_url`, `operating_hours` (JSONB),
  `queue_mode` (`FIFO` | `BOOKING_PRIORITY`), `is_active`, audit columns.

### `jenis_layanan` (Service)
A specific service under an agency; carries document requirements and duration.
- `id`, `instansi_id`, `name`, `description`, `estimasi_durasi_menit` (int),
  `requires_fo_verification` (bool), `is_active`, `default_second_service_id`
  (nullable self/other-service reference), audit columns.

### `syarat_dokumen` (Document requirement)
Document requirements attached to a service.
- `id`, `jenis_layanan_id`, `name`, `is_required` (bool), `notes`, `order`.

### `loket` (Counter)
A physical service counter.
- `id`, `instansi_id`, `name`/`code`, `status` (`OPEN` | `CLOSED` | `BREAK`),
  `is_active`, `last_idle_at` (timestamp used for idle-longest allocation),
  audit columns.

### `loket_layanan` (Counter↔Service map)
Which services a loket can serve (many-to-many).
- `loket_id`, `jenis_layanan_id`.

### `loket_session` (Operator session)
Binds a logged-in operator (user) to a loket for a shift.
- `id`, `loket_id`, `user_id`, `opened_at`, `closed_at`, `is_active`.

### `kuota_booking` (Booking quota)
Quota per date per agency (optionally per service).
- `id`, `instansi_id`, `jenis_layanan_id` (nullable), `tanggal` (date), `kuota` (int),
  `terpakai` (int, atomically incremented). Unique on (`instansi_id`,
  `jenis_layanan_id`, `tanggal`).

### `pemohon` (Applicant)
The citizen. Minimal PII, deduplicated by contact where possible.
- `id`, `name`, `phone` (WhatsApp), `email` (nullable), `nik` (optional/hashed per
  policy), `created_at`.

### `booking`
A scheduled registration (before check-in).
- `id`, `pemohon_id`, `instansi_id`, `jenis_layanan_id`, `tanggal`, `channel`
  (`WHATSAPP` | `WEB`), `qr_token` (single-use), `qr_expires_at`, `status`
  (`BOOKED` | `CHECKED_IN` | `EXPIRED` | `CANCELLED`), `created_at`.

### `antrian` (Queue number / Ticket)
The active queue item once a citizen is in line (from check-in or walk-in).
- `id`, `booking_id` (nullable for walk-in), `pemohon_id`, `instansi_id`,
  `jenis_layanan_id`, `nomor` (formatted, e.g. `A-014`), `nomor_seq` (int per
  service/day), `source` (`BOOKING` | `WALK_IN` | `SECOND_SERVICE`), `status`
  (see state machine), `loket_id` (nullable, set on call), `call_count` (0..3),
  `priority` (derived from mode), `fo_verified` (nullable bool),
  `queued_at`, `called_at`, `served_at`, `done_at`, `parent_antrian_id`
  (nullable, links a second-service item to its origin).

### `serving_session` (Service record)
One record of a loket serving one antrian (for reporting & durations).
- `id`, `antrian_id`, `loket_id`, `user_id`, `started_at`, `ended_at`, `outcome`
  (`DONE` | `SKIPPED` | `TRANSFERRED` | `HOLD`), `notes`.

### `fo_verification` (Front-office check)
Document verification outcome per antrian.
- `id`, `antrian_id`, `user_id` (FO), `result` (`COMPLETE` | `INCOMPLETE`),
  `checklist` (JSONB: per-syarat pass/fail), `notes`, `verified_at`.

### `audit_log`
Reuses the skeleton audit pattern (`core.audit_logs`) for sensitive MPP actions
(master-data edits, config changes, manual resets, transfers).

## Key relationships (summary)

- `instansi` **1—N** `jenis_layanan` **1—N** `syarat_dokumen`
- `instansi` **1—N** `loket` **N—N** `jenis_layanan` (via `loket_layanan`)
- `instansi` **1—N** `kuota_booking`
- `pemohon` **1—N** `booking` **1—0..1** `antrian`
- `antrian` **1—N** `serving_session`, `antrian` **1—0..1** `fo_verification`
- `antrian` **0..1—N** `antrian` (self, via `parent_antrian_id` for second service)

## Redis-held runtime state (not in Postgres)

- Per-service daily **number counter** (`mpp:counter:<service_id>:<yyyymmdd>`), atomic `INCR`.
- **Active queue** ordered set per service (waiting numbers), and current "called" per loket.
- **Pub/sub** channels for WebSocket fan-out (see `../04-api/websocket-events.md`).
- `loket` idle ranking mirror for fast idle-longest allocation.

Postgres remains the source of truth; Redis is the hot-path cache/coordination layer
and is rebuilt from Postgres on cold start / daily reset.
