---
name: reference_mpp_dev_engine
description: Kondisi engine dev MPP di mesin ini — Postgres compose di port 5433 (bentrok Homebrew PG), psql keg-only, Node 25 vs .nvmrc 22, cache Yarn global korup
metadata:
  type: reference
---

Hasil penyiapan engine 2026-08-05 (mesin macOS Fauzi). Semua jebakan di bawah menyita waktu
saat `make bootstrap` / `make db-setup` pertama:

- **Postgres compose dipetakan ke host port 5433, bukan 5432.** Homebrew `postgresql@16`
  (PID native, `/opt/homebrew/opt/postgresql@16`) sudah memegang `127.0.0.1:5432`, jadi
  `localhost:5432` dari host nyasar ke Postgres native → migrate gagal dengan
  `role "postgres" does not exist`. Perbaikannya di `.env` (`POSTGRES_PORT=5433`) dan
  `apps/api/.env` (`DB_PORT=5433`) — **bukan** mematikan service Homebrew.
- **`apps/api/.env.example` masih membawa nilai skeleton `tuai`** — `DB_NAME=tuai` harus
  diganti `mpp`, dan komentar `REDIS_DB=10 # reserved for tuai-be` menyesatkan.
- **`psql` keg-only** (Postgres 16 terpasang tapi tak di PATH) sementara target `make seed*`
  memanggil `psql` langsung → sudah dibereskan: `~/.zshrc` meng-export
  `/opt/homebrew/opt/postgresql@16/bin` ke PATH.
- **Node: pakai fnm, bukan Homebrew node.** Homebrew memasang node 25.1.0 (non-LTS) padahal
  `apps/web/.nvmrc` = 22 dan `engines: >=22.12.0`; Node 25 menolak `jsdom@30`. Sejak
  2026-08-05 terpasang **fnm + Node 22.23.2 (default)** dengan `eval "$(fnm env --use-on-cd
  --shell zsh)"` di `~/.zshrc` → masuk direktori ber-`.nvmrc` otomatis pindah versi. Node
  Homebrew tetap ada sebagai `system`. Environment test FE tetap **happy-dom** (sudah jalan &
  lebih ringan) — jangan gonta-ganti ke jsdom tanpa alasan. Divalidasi di Node 22:
  `yarn build` Next.js 16 exit 0 (~33 s), `make check` exit 0.
- **Cache Yarn global (`~/Library/Caches/Yarn/v6`) korup** — beberapa paket gagal ekstraksi
  (`hls.js`, `is-async-function`). Solusi tanpa merusak project lain: `yarn install
  --cache-folder <dir-sementara>`. Perbaikan permanen `yarn cache clean` belum dijalankan.
- **Colima** = runtime Docker di mesin ini (`colima start` dulu sebelum `make up`);
  `ssh ... colima/ssh.sock` yang terlihat listen di port 5432/6379 adalah port-forward Lima,
  bukan proses asing.

Verifikasi yang sudah lolos: `make check` exit 0 (backend `go test ./...` + Vitest 6 test +
`tsc:check`), schema `core` 17 tabel & `mpp` 13 tabel terisi seeder.

Lihat [[feedback_mpp_tdd_gate]] dan [[project_mpp_overview]].
