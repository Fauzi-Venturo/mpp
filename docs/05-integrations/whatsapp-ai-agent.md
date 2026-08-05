# Integration — WhatsApp AI Agent _(EN)_

Citizens register through WhatsApp, guided by an **LLM-powered AI agent** that fills the
booking slots conversationally and returns a QR check-in code.

## Components

- **WhatsApp Business API** (Cloud API or a BSP) — inbound webhook + outbound send.
- **AI agent orchestrator** (in the backend) — an LLM with **tools** bound to MPP
  read/write operations. Keeps per-conversation state.
- **Core Antrian Service** — quota check, booking creation, QR issuance.

## Inbound flow

```
WhatsApp message → POST /mpp/v1/wa/webhook (signature-verified)
  → load/create conversation state (keyed by phone)
  → AI agent step (LLM + tools)
  → outbound WhatsApp reply
```

## Agent tools (function calling)

| Tool | Maps to | Purpose |
|------|---------|---------|
| `list_instansi()` | GET catalog | offer agencies |
| `list_layanan(instansi_id)` | GET services | offer services + show `syarat_dokumen` |
| `check_availability(instansi_id, layanan_id, date)` | GET availability | remaining quota |
| `create_booking(pemohon, instansi_id, layanan_id, date)` | POST booking | make booking, get QR token |
| `cancel_booking(booking_id)` | POST cancel | cancel |
| `faq(query)` | KB lookup | answer hours/location/requirements |

The agent must **only** act through these tools (no free-form DB access) and must
confirm the summary with the citizen before `create_booking`.

## Slot-filling sequence

1. Intent detection ("mau perpanjang KTP").
2. Resolve **agency** → **service** (disambiguate with buttons/lists when unsure).
3. Present **document requirements** for the chosen service.
4. Collect/confirm applicant details (name, phone auto from WA, optional NIK per policy).
5. Offer **available dates** (respect quota + operating hours).
6. Confirm summary → `create_booking` → send **QR** + requirements + date/time.

## Outbound messages

- Booking confirmation with **QR image** (single-use token) + requirements checklist.
- Optional reminders (H-1 / day-of) and "almost your turn" nudges (see
  [`notifications.md`](./notifications.md)).

## Design constraints

- **Guardrails:** the agent never overrides quota or business rules; the service is the
  authority. Quota is consumed atomically at `create_booking`, not when the agent offers
  a date.
- **Fallback:** if the LLM is uncertain, degrade to structured WhatsApp
  buttons/list messages (FR-WA-07).
- **Idempotency:** de-dupe inbound webhook deliveries by WhatsApp message id.
- **Privacy:** minimize PII in prompts; do not log full message content with PII beyond
  retention policy (see `../06-security/security-privacy.md`).
- **Security:** verify webhook signatures; rate-limit per phone number.
- **State:** conversation state stored server-side (Redis/DB), keyed by phone, with a TTL.

## Backend placement

Implement under `apps/api/internal/modules/mpp/wa/` (`handler` = webhook,
`service` = agent orchestration + tool dispatch, `repository` = conversation state).
The LLM provider is configured via env; keep the provider behind an interface so it can
be swapped.
