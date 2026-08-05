---
name: project_mpp_overview
description: Apa itu project MPP, di mana dokumennya, dan status pengerjaan (relay 6 slice — modul booking, checkin, queue, loket_ops, serving sudah ada)
metadata:
  type: project
---

**MPP = Sistem Antrian Mal Pelayanan Publik** — antrian terpadu multi-instansi satu gedung.
Alur: warga daftar (WhatsApp AI agent / web / walk-in kiosk) → check-in QR di kiosk → dipanggil
lewat TV bersuara TTS Bahasa Indonesia **offline** → dilayani di loket. Ada verifikasi dokumen
front office, dashboard admin, supervisor per-instansi, dan pelaporan.

Monorepo polyglot ringan (bukan Turborepo/Nx): `apps/api` (Go, di-vendor dari
`venturo-skeleton-go`, module path `github.com/ndollem/mpp/apps/api`) + `apps/web` (Next.js,
di-vendor dari `venturo-skeleton-next.js`). Root cuma orkestrasi via `Makefile` +
`docker-compose.yml`.

**Dokumentasi lengkap ada di `docs/` — mulai dari `docs/README.md`** (26 file, 9 bagian:
overview, requirements, domain, architecture, api, integrations, security, ui-ux, roadmap).
Dokumen bisnis Bahasa Indonesia, dokumen teknis English. Jangan menebak requirement — semuanya
sudah tertulis di sana (18 modul FR, 28 business rule `BR-xx`, state machine antrian, ERD,
katalog REST + event WebSocket, matriks RBAC).

**Status per 2026-08-05:** migrasi + seeder `mpp` lengkap, dan pengerjaan berjalan sebagai
**relay 6 slice battle** (`docs/00-overview/battle-w04-mpp-relay.md`), bukan mengikuti urutan
`docs/08-roadmap/delivery-plan.md`. Modul Go yang sudah ada di `internal/modules/mpp/`:
`booking` (slice 1 kuota atomik + slice 2 QR), `checkin` (3), `queue` (4, alokasi nomor via
Redis), `loket_ops` (5, panggil/recall/start/skip), `serving` (6, done + display). FE punya
layar tiket QR di `apps/web/src/app/booking/[id]` + `src/sections/booking/`.

Dikerjakan **paralel oleh beberapa sesi Claude di working tree yang sama** — sebelum menyentuh
modul mana pun, cek `git status` dan pastikan tidak sedang dipegang sesi lain. Gate: `make check`.
Roadmap 6 fase produk (di luar battle) tetap di `docs/08-roadmap/delivery-plan.md`.

Lihat juga [[project_mpp_constraints]] dan [[feedback_mpp_stale_claudemd]].
