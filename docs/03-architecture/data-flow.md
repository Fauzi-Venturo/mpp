# Data Flow — key sequences _(EN)_

## 1. Registration via WhatsApp AI agent → booking + QR

```mermaid
sequenceDiagram
    actor C as Citizen (WhatsApp)
    participant WA as WhatsApp Business API
    participant API as Core Antrian Service
    participant AG as LLM AI agent
    participant PG as PostgreSQL
    C->>WA: "mau perpanjang KTP"
    WA->>API: inbound message webhook
    API->>AG: orchestrate (context + catalog tools)
    AG-->>API: ask agency/service/date (slot-filling)
    API->>WA: prompt + document requirements
    C->>WA: choose service + date
    WA->>API: message
    API->>PG: check quota (atomic), create pemohon+booking
    API->>API: generate single-use QR token
    API->>WA: send confirmation + QR + requirements
    WA-->>C: booking confirmed
```

## 2. On-site check-in (QR) → enqueue → ticket

```mermaid
sequenceDiagram
    actor C as Citizen
    participant K as Kiosk (browser)
    participant API as Core Antrian Service
    participant R as Redis
    participant WS as WebSocket hub
    participant TV as TV display
    C->>K: scan QR token
    K->>API: POST /checkin {token}
    API->>API: validate token (single-use, day, not expired)
    API->>R: INCR counter → nomor_seq
    API->>API: create antrian (WAITING), set nomor (A-014)
    API->>WS: publish queue.updated
    WS-->>TV: waiting list updated
    API-->>K: 200 {nomor, eta}
    K->>K: print thermal ticket
```

## 3. Call next → serve → done (idle-longest allocation)

```mermaid
sequenceDiagram
    actor O as Operator (loket)
    participant L as Loket app
    participant API as Core Antrian Service
    participant R as Redis
    participant WS as WebSocket hub
    participant TV as TV display (mini PC)
    O->>L: "Call next"
    L->>API: POST /antrian/next {loket_id}
    API->>R: pick head of service queue (mode FIFO/booking)
    API->>API: assign to this loket (idle-longest eligible), status CALLED, call_count=1
    API->>WS: publish call.created {nomor, loket}
    WS-->>TV: show number + play TTS (shared audio queue)
    WS-->>L: mark current
    O->>L: "Start" (citizen present)
    L->>API: POST /antrian/{id}/start → SERVING (open serving_session)
    O->>L: "Finish"
    L->>API: POST /antrian/{id}/done → DONE, loket last_idle_at=now
    API->>WS: publish queue.updated
```

## 4. No-show: recall up to 3× then skip

```mermaid
sequenceDiagram
    actor O as Operator
    participant API as Core Antrian Service
    participant WS as WebSocket hub
    participant TV as TV
    O->>API: recall (call_count=2) 
    API->>WS: call.recalled → TV re-plays TTS
    O->>API: recall (call_count=3)
    API->>WS: call.recalled → TV re-plays TTS
    O->>API: no-show → skip
    API->>API: status SKIPPED; loket last_idle_at=now
    API->>WS: queue.updated
```

## 5. Second service (no re-registration)

```mermaid
sequenceDiagram
    actor O as Operator (loket A)
    participant API as Core Antrian Service
    participant WS as WebSocket hub
    O->>API: while SERVING, "needs second service (service X)"
    API->>API: current → QUEUED_NEXT; create child antrian (source SECOND_SERVICE, parent set) → WAITING at X
    API->>WS: queue.updated (service X)
    Note over API: same pemohon, new number in X's stream
```

## 6. Daily reset & booking expiry (worker)

```mermaid
sequenceDiagram
    participant W as Worker
    participant R as Redis
    participant PG as PostgreSQL
    Note over W: 00:00 daily
    W->>R: reset per-service number counters
    W->>PG: close/expire leftover WAITING per policy
    Note over W: periodic
    W->>PG: find BOOKED past check-in window
    W->>PG: mark EXPIRED (no-show)
```
