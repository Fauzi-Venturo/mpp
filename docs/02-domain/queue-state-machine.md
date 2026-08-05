# Queue State Machine — `antrian` lifecycle _(EN)_

Every queue item (`antrian`, and its originating `booking`) moves through a strict set
of states. Illegal transitions are rejected server-side (NFR-DATA-03).

## States

| State          | Meaning |
|----------------|---------|
| `BOOKED`       | Scheduled via WhatsApp/Web; awaiting on-site check-in. (Lives on `booking`.) |
| `CHECKED_IN`   | QR checked-in at kiosk; about to enter the active queue. |
| `WAITING`      | In the active queue for its service, not yet called. |
| `CALLED`       | Announced (TV + TTS) to a specific loket; awaiting the citizen to approach. |
| `SERVING`      | Being served at the loket. |
| `HOLD`         | Temporarily paused mid-service (e.g. citizen fetching a document). |
| `DONE`         | Service completed. |
| `SKIPPED`      | Called up to 3× with no-show → skipped. |
| `TRANSFERRED`  | Moved to another loket/service (a new/continued queue item may result). |
| `QUEUED_NEXT`  | Second service auto-activated; re-enters `WAITING` for the next service. |
| `EXPIRED`      | Booking not checked-in within the allowed window (no-show before queue). |
| `CANCELLED`    | Cancelled by citizen/admin before serving. |

## Diagram

```mermaid
stateDiagram-v2
    [*] --> BOOKED: register (WhatsApp/Web)
    [*] --> WAITING: walk-in (kiosk) / print ticket

    BOOKED --> CHECKED_IN: QR check-in (valid token)
    BOOKED --> EXPIRED: check-in window passes (no-show)
    BOOKED --> CANCELLED: cancel (citizen/admin)

    CHECKED_IN --> WAITING: enqueue to service

    WAITING --> CALLED: operator "call next" (idle-longest loket)
    WAITING --> CANCELLED: cancel before call

    CALLED --> SERVING: citizen present, "start"
    CALLED --> CALLED: "recall" (call_count < 3)
    CALLED --> SKIPPED: no-show after 3rd call
    CALLED --> WAITING: requeue after skip (optional grace)

    SERVING --> HOLD: "hold"
    HOLD --> SERVING: "resume"
    SERVING --> DONE: "finish"
    SERVING --> TRANSFERRED: "transfer" to other loket/service
    SERVING --> QUEUED_NEXT: trigger second service

    TRANSFERRED --> WAITING: re-enqueued at target service
    QUEUED_NEXT --> WAITING: new antrian (parent_antrian_id set)

    DONE --> [*]
    SKIPPED --> [*]
    EXPIRED --> [*]
    CANCELLED --> [*]
```

## Transition table

| From | Event / Action | To | Guard / Effect |
|------|----------------|----|----------------|
| `BOOKED` | Valid QR check-in | `CHECKED_IN` | token unused & not expired & correct day; mark `booking.status=CHECKED_IN`; print ticket |
| `BOOKED` | Check-in window elapses | `EXPIRED` | scheduled job; counts as no-show; frees nothing (quota already consumed for the day) |
| `BOOKED` | Cancel | `CANCELLED` | decrement `kuota_booking.terpakai` if before cutoff |
| `CHECKED_IN` | Enqueue | `WAITING` | assign `nomor`/`nomor_seq` (atomic Redis INCR); set `queued_at`; broadcast |
| walk-in | Kiosk register | `WAITING` | create `antrian` (source `WALK_IN`); assign number; print ticket |
| `WAITING` | Call next | `CALLED` | pick idle-longest eligible loket; set `loket_id`, `called_at`, `call_count=1`; broadcast + TTS |
| `CALLED` | Recall | `CALLED` | `call_count += 1` (≤3); re-broadcast + TTS |
| `CALLED` | Start | `SERVING` | open `serving_session`; set `served_at`; update loket busy |
| `CALLED` | No-show (after 3rd) | `SKIPPED` | close as no-show; loket freed → `last_idle_at=now` |
| `CALLED` | Requeue (grace) | `WAITING` | optional: put back to end of queue instead of skipping |
| `SERVING` | Hold | `HOLD` | pause session; loket may serve others per policy |
| `HOLD` | Resume | `SERVING` | continue same session |
| `SERVING` | Finish | `DONE` | close `serving_session` (outcome DONE); set `done_at`; loket `last_idle_at=now`; broadcast |
| `SERVING` | Transfer | `TRANSFERRED` | close session (outcome TRANSFERRED); create/continue target `antrian` → `WAITING` |
| `SERVING` | Second service | `QUEUED_NEXT` | create child `antrian` (`source=SECOND_SERVICE`, `parent_antrian_id` set) → `WAITING`; no re-registration |
| `WAITING`/`CALLED` | Cancel | `CANCELLED` | admin/citizen; broadcast |

## Notes

- **`call_count` cap = 3** enforces the "call 3× then skip" rule (FR-OPR-03).
- On **`DONE`/`SKIPPED`/`TRANSFERRED`**, the loket's `last_idle_at` is refreshed so the
  idle-longest allocator (FR-QUE-02) picks fairly.
- **Daily reset (00:00)**: number counters reset; any still-`WAITING` items from the
  prior day are closed/carried per operational policy (see business rules).
- **Second service** keeps the same `pemohon`; the child item is a fresh `antrian` in
  the target service's stream, linked via `parent_antrian_id` for reporting.
