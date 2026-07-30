# WebSocket Events — MPP _(EN)_

Real-time change feed for TV displays, loket apps, FO, and admin dashboards. Backed by
Redis pub/sub so any API instance can serve any socket. REST remains the source of
truth; sockets carry deltas + snapshots.

## Connection

- Endpoint: `GET /mpp/v1/ws` (WebSocket upgrade).
- Auth: JWT (staff) or scoped API-key (devices) — same as REST.
- On connect, client sends a **subscribe** frame declaring interest; server replies with
  a **snapshot** so late joiners/reconnects re-sync.

```json
// client → server
{ "type": "subscribe", "channels": ["instansi:A", "loket:L3", "display:A"] }
```

## Channels

| Channel | Subscribers | Scope |
|---------|-------------|-------|
| `instansi:<prefix>` | admin, supervisor, TV | all queue activity for an agency |
| `layanan:<id>` | loket app, admin | one service's waiting stream |
| `loket:<id>` | that loket's operator | current called/serving item |
| `display:<instansi>` | TV mini PC | what to show + what to speak |
| `monitoring` | admin dashboard | cross-agency aggregate |

## Server → client events

| `type` | Payload (key fields) | Meaning |
|--------|----------------------|---------|
| `snapshot` | `{ waiting[], called[], serving[] }` | full current state on subscribe/reconnect |
| `queue.updated` | `{ layanan_id, waiting_count, next[] }` | waiting stream changed (enqueue/skip/done) |
| `call.created` | `{ antrian_id, nomor, loket, tts_text }` | a number was called → **TV shows + speaks** |
| `call.recalled` | `{ antrian_id, nomor, loket, call_count, tts_text }` | recall → **TV re-speaks** |
| `serving.started` | `{ antrian_id, loket }` | moved to SERVING |
| `serving.ended` | `{ antrian_id, outcome }` | DONE/SKIPPED/TRANSFERRED/HOLD |
| `second_service.created` | `{ parent_antrian_id, new_antrian_id, layanan_id }` | second service enqueued |
| `loket.status` | `{ loket_id, status }` | OPEN/CLOSED/BREAK changed |
| `broadcast` | `{ text, media_url? }` | admin info/running-text to TVs |
| `reset` | `{ instansi_id?, at }` | daily/manual reset happened → clients refetch |

## TTS payload contract (for the TV)

`call.created` / `call.recalled` carry a ready-to-speak string so the TV needs no extra
lookup, e.g.:

```json
{
  "type": "call.created",
  "antrian_id": "…",
  "nomor": "A-014",
  "loket": "Loket 3",
  "tts_text": "Nomor antrian A - nol satu empat, silakan menuju loket tiga"
}
```

The TV enqueues `tts_text` into the **single shared audio queue** (one mini PC, 3 TVs) so
announcements never overlap. See
[`../05-integrations/tts-voice-calling.md`](../05-integrations/tts-voice-calling.md).

## Delivery semantics

- **At-least-once**; clients dedupe by `antrian_id` + event sequence.
- On reconnect, client requests a fresh `snapshot` (do not replay missed deltas).
- Devices tolerate transient disconnects and keep displaying last-known state
  (offline resilience, NFR-AVAIL-02).
