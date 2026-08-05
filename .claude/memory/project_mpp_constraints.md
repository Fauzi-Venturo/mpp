---
name: project_mpp_constraints
description: Batasan MPP yang mudah dilanggar diam-diam — tenancy, Redis wajib, UTC, Yarn 1, no ORM, FE belum punya auth, TTS pre-rendered
metadata:
  type: project
---

Keputusan yang mengikat dan tersebar di beberapa dokumen — mudah dilanggar tanpa sadar:

- **Tenancy: 1 gedung MPP = 1 `core.companies` (tenant); `instansi` adalah entitas DI DALAM
  tenant itu, bukan company terpisah.** Salah petakan → pelaporan lintas instansi & loket/kiosk
  bersama jadi pecah. (`docs/02-domain/domain-model.md`)
- **Redis wajib, bukan opsional.** Bukan cuma cache RBAC bawaan skeleton — Redis adalah engine
  antrian (counter nomor atomik, active queue, ranking idle loket, pub/sub WebSocket). Supabase
  **tidak** menyediakan Redis, jadi hosting DB di Supabase tetap butuh Redis terpisah.
- **Postgres = source of truth, Redis = hot path yang rebuildable.** Jangan simpan state yang
  tak bisa dibangun ulang dari Postgres di Redis.
- **Semua timestamp UTC** di storage & API (`TZ=UTC` di backend); konversi ke WIB/WITA/WIT hanya
  di tampilan.
- **Backend: pgx v5 tanpa ORM** — SQL ditulis tangan di layer repository. Jangan tambah GORM/ent.
- **Frontend: Yarn 1 (Classic)**, bukan pnpm/npm. Gate verifikasi = `yarn tsc:check` + `yarn lint`
  (tidak ada test framework). `NEXT_PUBLIC_*` di-bake saat build → ganti env butuh rebuild.
- **FE skeleton dikirim TANPA auth sama sekali** — token store + hook `ky`
  `beforeRequest`/`afterResponse` + proteksi rute harus dibangun sendiri, mulai Fase 1.
- **Perangkat (kiosk & TV) autentikasi pakai API-key scoped, bukan JWT user.**
- **TTS: backend mengirim `tts_text` yang sudah siap diucap** di event `call.created`/
  `call.recalled`. TV tidak menormalkan angka sendiri. Sintesis offline di mini PC — pendekatan
  yang direkomendasikan dok adalah **potongan audio pre-rendered** (digit/huruf/frasa lalu
  digabung), bukan engine TTS; Web Speech API cuma fallback.
- **Satu mini PC → 3 TV, satu browser 3 window, SATU antrean audio bersama** (leader election via
  `BroadcastChannel`) supaya suara tak tumpang tindih (BR-18).
- **Atomisitas tidak boleh dikompromikan:** alokasi nomor antrian dan konsumsi kuota booking
  harus bebas race (Redis `INCR` + unique constraint Postgres). Ini bug yang gagal diam-diam.
- **Transisi state ilegal ditolak `409`**, mengacu ke `docs/02-domain/queue-state-machine.md`.
  Aturan `call_count` maksimal 3 lalu `SKIPPED`.

Lihat [[project_mpp_overview]].
