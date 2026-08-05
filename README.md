# MPP — Sistem Antrian Mall Pelayanan Publik

Monorepo untuk **Sistem Antrian Mall Pelayanan Publik (MPP)** — sistem manajemen
antrian multi-instansi untuk mal pelayanan publik: registrasi via WhatsApp AI agent,
web, dan walk-in; check-in QR di kiosk; display TV dengan pemanggilan suara (TTS
Bahasa Indonesia offline); aplikasi loket untuk petugas; verifikasi dokumen front
office; dashboard admin; dan pelaporan.

> _Monorepo for a public-service-mall queue management system: multi-agency queue
> handling with WhatsApp AI-agent / web / walk-in registration, QR kiosk check-in,
> TV displays with offline Indonesian text-to-speech voice calling, counter (loket)
> operator apps, front-office document verification, admin dashboards, and reporting._

---

## Struktur monorepo

```
mpp/
├── apps/
│   ├── api/          # Backend — Go 1.26 + Gin + pgx  (dari venturo-skeleton-go)
│   └── web/          # Frontend — Next.js 16 + MUI    (dari venturo-skeleton-next.js)
├── packages/
│   └── api-contract/ # Kontrak & tipe bersama FE↔BE (placeholder)
├── infra/
│   └── supabase/     # Catatan setup Postgres/Supabase
├── docs/             # 📚 Dokumentasi requirement lengkap (mulai di sini)
├── docker-compose.yml
├── Makefile
└── .env.example
```

Ini adalah **meta-repo polyglot ringan**: `apps/api` (Go) dan `apps/web`
(Next.js/Yarn 1) berdiri sendiri dengan toolchain masing-masing. Root hanya
mengorkestrasi keduanya lewat `Makefile` + `docker-compose.yml` — tidak ada
Turborepo/Nx, karena Go dan Next tidak berbagi dependency graph JS.

## Tech stack

| Lapisan   | Teknologi                                                              |
|-----------|-----------------------------------------------------------------------|
| Backend   | Go 1.26, Gin, pgx v5 (tanpa ORM), golang-migrate, Redis, JWT/API-key   |
| Frontend  | Next.js 16 (App Router), React 19, MUI v9, TanStack Query + ky, RHF+Zod |
| Database  | PostgreSQL 16 (Supabase-compatible — lihat `infra/supabase/`)          |
| Realtime  | WebSocket (pemanggilan antrian → TV/loket/kiosk) + Redis pub/sub        |
| Cache     | Redis (permission cache + counter nomor antrian + state aktif)         |

## Quickstart

```bash
# 1. Install deps + siapkan file .env
make bootstrap

# 2. Jalankan infra (postgres + redis)
make up

# 3. Migrasi + seed database
make db-setup

# 4. Jalankan backend & frontend (dua terminal, hot-reload)
make api-dev      # http://localhost:8080
make web-dev      # http://localhost:8002
```

Atau jalankan seluruh stack dalam container: `make up-full`.
Lihat semua perintah: `make help`.

## Testing (TDD)

Fitur dikerjakan **test-first**: tulis test yang gagal dulu, baru implementasinya.

```bash
make check        # gate lengkap: backend test + frontend test + type-check
make api-test     # Go — testify; integration test pakai Postgres/Redis dari compose
make web-test     # Vitest + Testing Library (happy-dom)
```

- Backend: helper `internal/testsupport` menyediakan `Postgres(t)`/`Redis(t)`; test
  **skip** (bukan gagal) kalau `make up` belum dijalankan. Jalankan lewat `make api-test`,
  bukan `go test ./...` telanjang, supaya `.env` ikut ter-export.
- Frontend: test berdampingan dengan kode sebagai `src/**/*.test.ts(x)`.
- Catatan mesin lokal: kalau port 5432 sudah dipakai Postgres lain, `.env` memakai
  `POSTGRES_PORT=5433` — samakan dengan `DB_PORT` di `apps/api/.env`.

## Database: PostgreSQL & Supabase

Development lokal memakai Postgres + Redis dari `docker-compose.yml`. Untuk hosting
sementara di **Supabase** (managed Postgres), arahkan `DB_*` backend ke koneksi
Supabase (`DB_SSLMODE=require`, gunakan endpoint pooler). **Catatan penting:**
Supabase tidak menyediakan Redis — sedangkan backend butuh Redis untuk permission
cache _dan_ engine antrian (counter nomor, state aktif, pub/sub). Sediakan Redis
terpisah (compose lokal atau hosted seperti Upstash). Detail: `infra/supabase/README.md`.

## Dokumentasi

Seluruh dokumentasi requirement ada di **[`docs/`](./docs/)** — mulai dari
[`docs/README.md`](./docs/README.md). Dokumen bisnis (PRD/FRD/aturan bisnis/UI)
ditulis dalam Bahasa Indonesia; dokumen teknis (arsitektur/API/domain/keamanan)
dalam Bahasa Inggris.

## Asal template

- Backend di-vendor dari [`venturo-id/venturo-skeleton-go`](https://github.com/venturo-id/venturo-skeleton-go)
  (module path di-rename ke `github.com/ndollem/mpp/apps/api`).
- Frontend di-vendor dari [`venturo-id/venturo-skeleton-next.js`](https://github.com/venturo-id/venturo-skeleton-next.js).

Keduanya disalin penuh (bukan submodule); riwayat git upstream tidak ikut.
