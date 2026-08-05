---
name: feedback_mpp_tdd_gate
description: Tiap slice MPP dikerjakan TDD (test merah dulu) + 3 gate battle — plan mode, ≥1 subagent, test hijau — sebelum /grademe & commit
metadata:
  type: feedback
---

Setiap slice MPP dikerjakan dengan urutan **plan → test MERAH dulu → implementasi →
test HIJAU → `/grademe <nama>-s<N>` → commit**. Test ditulis sebelum implementasi, bukan
sesudahnya untuk mengejar formalitas.

**Why:** MPP dipakai sebagai materi Battle #1 (w04), dan "test buatanmu hijau" adalah salah
satu dari 3 gate wajib — gagal 1 gate berarti keluar dari peringkat, dan rata-rata skor
`/grademe` (plan · context · subagent · test) jadi tie-breaker. Di luar konteks lomba pun,
bug paling berbahaya di domain ini **gagal diam-diam**: overbooking karena kuota tidak atomik
dan nomor antrian dobel hanya kelihatan lewat test negatif, tidak lewat happy-path manual.
Lihat [[project_mpp_constraints]] dan `docs/00-overview/battle-w04-mpp-relay.md`.

**How to apply:**

- **3 gate tiap slice:** (1) plan mode dulu, (2) dispatch ≥1 subagent, (3) test hijau.
- **Cakupan minimal per slice: happy-path + ≥1 kasus negatif.** Negatif yang wajib ada:
  kuota penuh → `409` tanpa overbooking (slice 1), QR reuse/expired ditolak (slice 3),
  nomor antrian tidak pernah dobel di alokasi paralel (slice 4), recall ke-4 ditolak
  (slice 5), transisi state ilegal → `409` (semua slice).
- **Gate verifikasi backend:** `make api-test` (bukan `go test ./...` telanjang — target Make
  meng-`include .env` sehingga test integrasi menemukan Postgres compose; tanpa itu test
  integrasi ter-*skip* diam-diam). Helper: `internal/testsupport` (`Postgres(t)`, `Redis(t)`,
  skip bila infra mati). Assertion: `testify`.
- **Gate verifikasi frontend:** `yarn test` (Vitest + Testing Library, environment `happy-dom`)
  + `yarn tsc:check` + `yarn lint`. **Vitest dipasang 2026-08-05 atas permintaan eksplisit
  user** ("projek ini HARUS pakai TDD") — ini membatalkan catatan lama yang melarang memasang
  test framework di `apps/web`. Konfigurasi: `vitest.config.mts` + `vitest.setup.ts`.
- **Gate gabungan:** `make check` = `api-test` + `web-test` + `web-check`.
- Jangan lompat slice: selesai satu slice → `/grademe` → commit → baru slice berikutnya.
