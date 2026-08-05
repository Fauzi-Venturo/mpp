# Non-Functional Requirements (NFR) — MPP _(EN)_

Notation: `NFR-<area>-<n>`. Targets are indicative and should be confirmed with
stakeholders before load testing.

## 1. Performance & real-time latency

| ID | Requirement |
|----|-------------|
| NFR-PERF-01 | Queue actions (call / next / skip / serve / done) reflect on TV displays and loket clients within **≤ 1 s p95** end-to-end (action → WebSocket broadcast → UI update). |
| NFR-PERF-02 | Standard REST reads (catalog, booking availability) respond **≤ 300 ms p95** under expected load. |
| NFR-PERF-03 | QR check-in (scan → ticket printed) completes within **≤ 2 s** on the kiosk. |
| NFR-PERF-04 | Number allocation and counters are served from Redis; DB is not on the hot path of each call. |

## 2. Availability & resilience

| ID | Requirement |
|----|-------------|
| NFR-AVAIL-01 | System availability during operating hours **≥ 99.5%**. |
| NFR-AVAIL-02 | **On-site voice calling (TV/TTS) must keep working during an internet outage** — audio assets and the TTS engine run locally on the mini PC; the display degrades gracefully to last-known queue state. |
| NFR-AVAIL-03 | Loss of a single loket client must not corrupt shared queue state (state lives in Redis/DB, not the client). |
| NFR-AVAIL-04 | WebSocket clients auto-reconnect and re-sync current queue state on reconnect. |
| NFR-AVAIL-05 | Thermal printer / QR scanner failure surfaces a clear kiosk error and a manual fallback path. |

## 3. Scalability

| ID | Requirement |
|----|-------------|
| NFR-SCALE-01 | Support **tens of agencies × multiple lokets** in a single MPP building, hundreds of concurrent active numbers, and dozens of simultaneous WebSocket clients (TVs, lokets, admin). |
| NFR-SCALE-02 | The Core Antrian Service is horizontally scalable; shared state (Redis + Postgres) is external, so instances are stateless behind a load balancer. |
| NFR-SCALE-03 | Data model reserves room for multi-building expansion without redesign. |

## 4. Security & privacy

| ID | Requirement |
|----|-------------|
| NFR-SEC-01 | All traffic over TLS. JWT auth for users; scoped API-keys for devices (kiosk/TV). |
| NFR-SEC-02 | RBAC enforced server-side on every mutating endpoint (see `../06-security/rbac-matrix.md`). |
| NFR-SEC-03 | PII (names, contact, uploaded documents) is minimized, access-controlled, and retained per policy (see `../06-security/security-privacy.md`). |
| NFR-SEC-04 | QR tokens are single-use and time-bound; replay is rejected. |
| NFR-SEC-05 | Sensitive/administrative actions are captured in the audit log. |
| NFR-SEC-06 | Rate limiting on public registration endpoints to deter abuse/spam. |

## 5. Usability & accessibility

| ID | Requirement |
|----|-------------|
| NFR-UX-01 | Citizen web is responsive (mobile-first) and understandable by non-technical users. |
| NFR-UX-02 | TV displays are legible from a distance (large type, high contrast); voice calling aids the visually impaired. |
| NFR-UX-03 | Loket app is operable quickly with minimal clicks; primary actions reachable in one tap. |
| NFR-UX-04 | Aim toward WCAG 2.1 AA for citizen-facing web where feasible. |

## 6. Internationalization & localization

| ID | Requirement |
|----|-------------|
| NFR-I18N-01 | Primary language **Bahasa Indonesia**; UI copy externalized to support future languages (skeleton ships translation-overrides). |
| NFR-I18N-02 | TTS pronounces Indonesian numbers and phrases correctly (e.g. "A tiga belas"). |
| NFR-I18N-03 | Timezone fixed to **UTC in storage** (backend enforces `TZ=UTC`); presentation converts to WIB/WITA/WIT as configured. |

## 7. Observability & operability

| ID | Requirement |
|----|-------------|
| NFR-OPS-01 | Structured logging (backend uses zap) with correlation of queue events. |
| NFR-OPS-02 | Health endpoints for API, DB, and Redis; container healthchecks. |
| NFR-OPS-03 | Basic metrics: active numbers, average wait/serve time, per-loket throughput. |
| NFR-OPS-04 | Migrations are versioned and reproducible across environments. |

## 8. Maintainability & portability

| ID | Requirement |
|----|-------------|
| NFR-MNT-01 | Backend follows the skeleton's modular-monolith layout; MPP domain isolated under its own module & schema. |
| NFR-MNT-02 | Runs on PostgreSQL, incl. managed Supabase; Redis provided separately. |
| NFR-MNT-03 | Local dev reproducible via `docker-compose` + `make`. |

## 9. Data integrity & consistency

| ID | Requirement |
|----|-------------|
| NFR-DATA-01 | Queue-number allocation is race-free (atomic counters in Redis; unique constraints in Postgres). |
| NFR-DATA-02 | Quota enforcement is atomic — no overbooking under concurrency. |
| NFR-DATA-03 | State transitions follow the documented state machine; illegal transitions are rejected. |
