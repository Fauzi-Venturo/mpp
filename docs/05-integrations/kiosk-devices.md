# Integration — Kiosk & On-site Devices _(EN)_

## Kiosk

A self-service station for **QR check-in**, **walk-in registration**, and **thermal
ticket printing**.

### Hardware
- Touchscreen PC in **browser kiosk mode** (locked to the kiosk route).
- **QR scanner** — typically USB HID (keyboard-emulation): a scan arrives as keystrokes
  ending in Enter. The kiosk page captures it into a hidden input.
- **Thermal printer** (58/80 mm) — via browser print (a print-CSS ticket) or a small
  **local print agent** (raw ESC/POS) for reliability.

### Check-in flow
```
Scan QR token → POST /mpp/v1/checkin {token}
  → valid (single-use, correct day, not expired)?
     yes → antrian WAITING (number assigned) → print ticket → show number + ETA
     no  → clear error (used/expired/wrong day) + guidance
```

### Walk-in flow
```
Pick agency → pick service (show requirements) → confirm
  → POST /mpp/v1/walkin → number assigned (WAITING) → print ticket
```

### Ticket content
Number (e.g. `A-014`), agency + service, date/time, estimated wait, and a note to watch
the TV/listen for the call. Keep it minimal (no sensitive PII on paper).

### Resilience
- Scanner or printer failure → clear on-screen error + manual fallback (staff-assisted).
- Kiosk authenticates with a **scoped API-key** (not a user JWT).
- Debounce double-scans; ignore repeat within a short window.

## TV display

See [`tts-voice-calling.md`](./tts-voice-calling.md) — mini PC, three windows, shared
offline audio queue.

## Loket terminal

Operator workstation running the **loket app** (browser route) behind user JWT auth.
Primary actions reachable in one tap: call next, recall, start, finish, skip, hold,
transfer, second service. Subscribes to `loket:<id>` and `layanan:<id>` WebSocket
channels for live state.

## Device auth summary

| Device | Auth | Notes |
|--------|------|-------|
| Kiosk  | scoped API-key | check-in + walk-in + print only |
| TV     | scoped API-key | read display snapshot + subscribe display channel |
| Loket  | user JWT | operator session bound to a loket |
| FO / Admin | user JWT | RBAC per role |

## Frontend placement

Kiosk, TV, and loket are `apps/web` routes with device-appropriate layouts; the
kiosk/TV routes are gated by API-key context rather than user login.
