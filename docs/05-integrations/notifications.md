# Integration — Notifications _(EN)_

## Channels

| Channel | Transport | Use |
|---------|-----------|-----|
| WhatsApp | WhatsApp Business API | booking confirmation + QR, reminders, "almost your turn" |
| Email | SMTP (skeleton `pkg/email`, `gomail.v2`) | confirmation + QR (downloadable), receipts |
| SMS | provider (optional) | fallback reminder when WhatsApp unavailable (FR-NTF-04, could-have) |

## Notification types

| Event | Trigger | Channel(s) | Content |
|-------|---------|-----------|---------|
| Booking confirmed | `create_booking` | WhatsApp (+ email if given) | number/date, **QR token**, document requirements |
| Reminder H-1 | worker, day before | WhatsApp / SMS | date/time, requirements, cancel link |
| Reminder day-of | worker, morning | WhatsApp / SMS | check-in instructions + QR |
| Almost your turn | queue position ≤ N | WhatsApp | "±N numbers before you, please head over" |
| Cancelled/Expired | cancel or expiry sweep | WhatsApp / email | status + rebook guidance |

## Timing & throttling

- Reminders and "almost your turn" are dispatched by the **worker** (leader-owned
  schedule) so they fire once.
- Throttle per applicant to avoid spam; respect quiet hours where applicable.
- "Almost your turn" threshold `N` is configurable per agency.

## Templates

- Message templates (WhatsApp/email) and TTS phrasing are **configurable**
  (FR-CFG-03). WhatsApp requires **pre-approved templates** for business-initiated
  (outside the 24-hour session) messages — plan template approval for reminders.

## Reliability

- Retries with backoff on transient send failures; dead-letter for permanent failures.
- Idempotent dispatch keyed by (applicant, notification type, date) to prevent
  duplicates on worker restart.

## Backend placement

`apps/api/internal/modules/mpp/notification/` (service + templates), reusing
`pkg/email` for SMTP and the WhatsApp send client from the WA integration. Reminder and
expiry jobs live in the worker.
