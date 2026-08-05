# Security & Privacy — MPP _(EN)_

Public-sector queue system handling citizen PII. Security is server-enforced; privacy
follows data minimization.

## Authentication & authorization

- **TLS everywhere.** No plaintext transport.
- **JWT** for staff (short-lived access + refresh); **scoped API-keys** for devices.
  The backend rejects weak/default JWT secrets at boot in production.
- **RBAC** enforced on every mutating endpoint, with agency/loket scoping (see
  [`rbac-matrix.md`](./rbac-matrix.md)). Permissions cached in Redis, DB fallback.
- **Account protection:** login attempt limits + lockout (skeleton settings), password
  hashing (`golang.org/x/crypto`).

## QR token security

- Tokens are **single-use** and **time-bound** (valid on the booking day within a
  window). Used/expired/wrong-day tokens are rejected (`409/410`).
- Tokens are unguessable (crypto-random), stored hashed where feasible, and invalidated
  on check-in, cancel, or expiry.

## PII handling & data minimization

| Data | Collected? | Notes |
|------|-----------|-------|
| Name | yes | needed for service |
| Phone (WhatsApp) | yes | primary contact/channel |
| Email | optional | confirmations |
| NIK / ID number | **only if the service requires it** | store hashed/masked where possible; access-restricted |
| Uploaded documents | only if FO upload is enabled | restricted bucket, short retention |

- Collect the **minimum** required per service; don't demand NIK unless the service does.
- **Do not** print sensitive PII on kiosk tickets (number + service only).
- Minimize PII in LLM prompts and logs (WhatsApp agent).

## Storage & retention

- Uploaded documents in an **access-controlled bucket**, retained only as long as needed
  for the service, then purged per policy.
- Operational records (antrian, serving_session) retained for reporting; consider
  **anonymizing** PII after a defined period while keeping aggregate metrics.
- Redis holds transient queue state (rebuildable) — no long-term PII store.

## Audit & accountability

- **Audit log** (skeleton `core.audit_logs`) for sensitive actions: master-data edits,
  config changes, manual resets, transfers, role assignments.
- Logs capture actor, action, target, timestamp — not raw PII payloads.

## Abuse prevention

- **Rate limiting** on public registration (web + WhatsApp) per IP/phone to deter spam
  and quota exhaustion.
- Webhook **signature verification** for WhatsApp inbound.
- Idempotency keys to prevent duplicate bookings from retries.

## Operational security

- Secrets via env/secret manager, never committed (`.gitignore` covers `.env`,
  service-account JSON, PEM).
- Least-privilege DB and storage credentials.
- Separate device API-keys per device class; rotate on compromise.
- `TZ=UTC` in storage; consistent timestamps for audit integrity.

## Compliance posture

- Align with Indonesian personal-data protection (UU PDP) principles: lawful basis,
  minimization, purpose limitation, retention limits, and subject rights (access,
  correction, deletion where applicable). Confirm specifics with the MPP legal owner
  before go-live.
