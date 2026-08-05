---
name: project_mpp_overview
description: Apa itu project MPP, di mana dokumennya, dan fase pengerjaan saat ini (Fase 0 selesai, kode mpp belum ada)
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

**Status per 2026-08-05: Fase 0 selesai** (dokumentasi + monorepo ter-vendor). **Kode MPP belum
ditulis sama sekali** — `internal/modules/mpp/*`, `migrations/mpp/`, `seeders/mpp/` masih
rencana. Roadmap 6 fase di `docs/08-roadmap/delivery-plan.md`: 1) master data + auth FE,
2) booking & check-in, 3) engine antrian + WS + TV/TTS, 4) loket & FO, 5) WhatsApp AI agent,
6) monitoring & pelaporan.

Lihat juga [[project_mpp_constraints]] dan [[feedback_mpp_stale_claudemd]].
