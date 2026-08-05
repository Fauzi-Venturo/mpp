# Battle #1 — MPP Relay (Claude Code Unleashed · Season 1 · Minggu 4)

> Rangkuman lengkap deck presentasi:
> <https://ndollem.github.io/venturo-cc-sessions/w04/slides/index.html>
> Sumber = 10 slide, ditranskrip apa adanya lalu dirapikan. Detail teknis mengikat tetap di `soal.md` + `docs/`.

**Ringkas:** Antrian Mal Pelayanan Publik, dibangun potong demi potong.
90 menit · relay 6 slice berurutan · menang karena **PALING JAUH + BUKTI**.

Repo = fork **`ndollem/mpp`** (project REAL). Kamu yang ngoding, test kamu = buktinya.
**Tim pemenang berpeluang lanjut menggarap MPP sungguhan.**

---

## 0. Status repo ini vs deck (sync 2026-08-05)

Deck ditulis untuk fork `ndollem/mpp`. Repo lokal ini (`Venturo/Mpp`, branch `main`) sudah
melewati sebagian window setup, tapi **belum ada satu pun slice yang dikerjakan**.

| Yang deck harapkan | Status di repo ini | Bukti |
|---|---|---|
| Fork + clone repo | ✅ selesai (monorepo sudah ada lokal) | commit `ba2f0c9`, `747bd90` |
| Migration `mpp` + seed disediakan | ✅ ada — 5 migration + 7 seeder | [migrations/mpp/](../../apps/api/internal/database/migrations/mpp/), [seeders/mpp/](../../apps/api/internal/database/seeders/mpp/) |
| Go 1.26 | ✅ `go 1.26.1` | [apps/api/go.mod:3](../../apps/api/go.mod#L3) |
| `make bootstrap` / `up` / `db-setup` / `api-dev` / `web-dev` | ✅ semua target ada | [Makefile](../../Makefile) |
| Postgres + Redis | ✅ tersedia via compose | [docker-compose.yml:28-36](../../docker-compose.yml#L28-L36) |
| `.env` terisi, deps ter-install | ❌ **belum** — `.env` belum ada di root/api/web, `node_modules` kosong | perlu `make bootstrap` |
| Modul Go `mpp` (slice 1–6) | ❌ **belum ada** — `internal/modules/` baru `core/` | [internal/modules/](../../apps/api/internal/modules/) |
| Layar frontend MPP | ❌ **belum ada** — `apps/web` masih template Minimal-UI mentah | [apps/web/src/app/](../../apps/web/src/app/) |
| Test slice | ❌ belum ada | — |
| `soal.md` | ❌ tidak ada di repo — kontrak detail dipakai dari `docs/` | — |

**Tabel DB yang sudah siap dipakai** (schema `mpp`): `instansi`, `jenis_layanan`,
`syarat_dokumen`, `loket`, `loket_layanan`, `loket_session`, `kuota_booking`, `pemohon`,
`booking`, `antrian`, `serving_session`, `fo_verification`, `system_config`.
Artinya slice 1–6 **tidak perlu bikin migration baru**, tinggal bangun modul + handler.

**Kontrak endpoint battle sudah terdokumentasi** di
[04-api/rest-endpoints.md](../04-api/rest-endpoints.md) — `/mpp/v1/booking` (baris 31),
`/mpp/v1/checkin` (41), `/mpp/v1/queue/next` (49), `/mpp/v1/antrian/{id}/{recall,start,done,skip}`
(50–53), `/mpp/v1/display` (69). Header `X-API-Key` & `X-Company-Slug` dijelaskan di
[04-api/api-conventions.md:33-38](../04-api/api-conventions.md#L33-L38).

### Koreksi terhadap langkah setup di deck

- **STEP 1–2 (fork + clone)** — sudah lewat untuk repo ini, langsung ke STEP 3.
- **STEP 5** — deck menyebut `REDIS_URL`; repo ini **tidak** memakai variabel itu. Yang benar:
  `REDIS_HOST` / `REDIS_PORT` / `REDIS_PASSWORD` / `REDIS_DB` di
  [apps/api/.env.example:29-34](../../apps/api/.env.example#L29-L34), dan `POSTGRES_*` /
  `REDIS_PORT` / `API_PORT=8080` / `WEB_PORT=8002` di [.env.example](../../.env.example).
- **`make up`** menjalankan `docker compose up -d postgres redis` (infra saja).
  Ada juga `make up-full` (ikut api + web) yang tidak disebut di deck.
- Verifikasi cepat sebelum mulai: `make check` (= `api-test` + `web-check`).

### Sisa pekerjaan langsung (urut)

1. `make bootstrap` → isi 3 file `.env` → `make up` → `make db-setup`
2. `make api-dev` (:8080) + `make web-dev` (:8002) — pastikan keduanya hijau
3. Baru masuk Slice 1: modul Go `internal/modules/mpp/booking` + test

---

## 1. Kenapa battle & kenapa relay

3 minggu terakhir dipakai menambah senjata. Hari ini yang diuji: **mengubah spec REAL jadi fitur jalan — berulang, di bawah tekanan waktu.**

Ini project sungguhan; tim pemenang berpeluang melanjutkan pengerjaan MPP. Skor hari ini = baris pertama leaderboard.

Alur relay:

```
Pendaftaran → QR → Check-in → Nomor → Panggil → Selesai
```

**URUT, tidak boleh lompat.**

---

## 2. Aturan main

| | Aturan |
|---|---|
| ⏱ | **90 menit** coding (di luar window setup · bisa diperpanjang +30') |
| → | **Urut & no-skip** · selesai 1 slice → `/grademe` → lanjut |
| ✓ | Boleh baca `docs/` & repo (context ✅) |
| ✗ | **Dilarang** minta bantuan manusia lain |

### 3 Gate wajib (syarat ikut juara)

1. **Plan mode** tiap slice
2. **≥1 subagent** di-dispatch
3. **Test buatanmu HIJAU**

> Gagal 1 gate → skor tetap dihitung, tapi **tidak ikut peringkat juara**.

---

## 3. Repo & Window Setup (±10 menit, timer coding belum jalan)

Stack: **Go 1.26 · Next.js 16 · MUI v9 · `ky` · RHF/Zod**, Postgres + Redis.

- **Disediakan panitia:** migration `mpp` + seed.
- **Kamu bangun:** modul Go + layar (frontend) + TEST.
- **PATUHI** `docs/` + `CLAUDE.md`.

```bash
# STEP 1 — FORK di browser: github.com/ndollem/mpp → tombol "Fork"

# STEP 2 — CLONE fork-MU (ganti <user-kamu>) lalu masuk folder:
git clone https://github.com/<user-kamu>/mpp.git && cd mpp

# STEP 3 — INSTALL deps (Go + Yarn) + siapkan .env:
make bootstrap

# STEP 4 — NYALAKAN Postgres + Redis (Redis WAJIB — Slice 4):
make up            # A) Docker: Postgres+Redis  ·  B) Supabase+Upstash → skip

# STEP 5 — ISI .env: DB_* + REDIS_URL sesuai infra STEP 4

# STEP 6 — JALANKAN MIGRATION + SEED (disediakan panitia):
make db-setup

# STEP 7 — JALANKAN backend & frontend (2 terminal):
make api-dev       # :8080
make web-dev       # :8002
```

Gerbang mulai: `db-setup` sukses, `:8080` + `:8002` nyala → soal reveal + timer coding 90' mulai.

---

## 4. Papan 6 Slice (peta skor)

| # | Slice | Tujuan | Transisi state |
|---|---|---|---|
| 1 | **Pendaftaran** | pilih instansi+layanan, cek kuota, buat booking | `[*] → BOOKED` |
| 2 | **Terbitkan QR** | token sekali pakai + terikat waktu, tampil/unduh | `qr_token` |
| 3 | **Check-in** | scan QR, validasi, tolak reuse/expired | `BOOKED → CHECKED_IN` |
| 4 | **Nomor antrian** | alokasi `A-014` (Redis `INCR`), masuk stream | `CHECKED_IN → WAITING` |
| 5 | **Panggil** | loket idle-terlama, maks 3× lalu skip | `WAITING → CALLED → SERVING` |
| 6 | **Selesai + Display** | tutup layanan + TV tampil nomor & loket | `SERVING → DONE` |

> Kerjakan **URUT** — sejauh mana kamu sampai, itu skormu.

---

## 5. Kontrak 6 slice

| # | Inti | Endpoint utama | KOMPLIT inti |
|---|---|---|---|
| 1 | **Booking** + kuota atomik | `POST /mpp/v1/booking` (penuh→409) | `201` + **BOOKED**; penuh→`409` tanpa overbook |
| 2 | **QR** sekali pakai + expiry | `qr_token` di respons booking · `GET /booking/{id}` | token unik + expiry; layar QR bisa diunduh |
| 3 | **Check-in** (device `X-API-Key`) | `POST /mpp/v1/checkin` | valid→**CHECKED_IN**; reuse/expired→ditolak |
| 4 | **Nomor A-014** (Redis `INCR`) | alokasi di `/checkin` · `GET /queue` | urut atomik tak dobel; →**WAITING** |
| 5 | **Panggil** idle-terlama, maks 3× | `/queue/next` · `/antrian/{id}/{recall,start,skip}` | **CALLED→SERVING**; recall ke-4 ditolak |
| 6 | **Selesai** + TV display | `/antrian/{id}/done` · `GET /display` | →**DONE** catat durasi; TV nomor+loket |

Ikuti kontrak **persis** (path · status · field) — envelope `{data, message, meta, errors}`, timestamp **UTC**.

---

## 6. Definisi KOMPLIT

**KOMPLIT = jalan + hijau + grademe**

1. **Happy-path jalan** (fasilitator boleh mencoba di laptopmu)
2. **Test buatanmu HIJAU** — happy path + ≥1 kasus negatif (w03: bukti, bukan janji)
3. Sudah **`/grademe <nama>-s<N>`** + commit

### ❗ Jebakan

- Atomisitas (kuota & nomor antrian)
- Device pakai **`X-API-Key`**, **bukan** JWT
- Header **`X-Company-Slug`**
- Guard transisi state ilegal

### Contoh output yang diharapkan

```console
› go test ./... && /grademe budi-s1

# sebelum — logika Slice 1 belum benar (MERAH)
--- FAIL: TestBooking_QuotaFull (0.01s)
    want 409, got 201 (overbooking!)

# sesudah — benar (HIJAU)
ok  .../internal/modules/mpp/booking  0.21s

› /grademe budi-s1
✓ transcript dinilai · skor 84/100
POST /grademe → 201 Created
```

---

## 7. Skor & Leaderboard

1. **Jumlah slice KOMPLIT (0–6)** — paling banyak selesai menang
2. Tie-break: **rata-rata skor `/grademe`** (kualitas: plan · context · subagent · test)
3. Tie-break: **slice terjauh tercepat**

> **Pemenang = paling jauh + paling rapi.** Gate FAIL → di luar peringkat juara.

---

## 8. Saat mulai (papan tetap tayang selama battle)

- **1** Booking (kuota atomik) → **2** QR sekali-pakai → **3** Check-in (`X-API-Key`)
- **4** Nomor A-014 (Redis) → **5** Panggil idle-terlama (maks 3×) → **6** Selesai + TV
- Tiap slice: **plan → build → test HIJAU → `/grademe`** · URUT, no-skip
- Menang = **slice terbanyak KOMPLIT** + rata-rata `/grademe` terbaik

Timer: **90:00** (`T` mulai/jeda · `R` reset · `E` extend +10', maks +30').

---

## 9. Selesai → submit

- ⏹ Waktu habis = **hands-off keyboard**
- ✓ Tiap slice KOMPLIT: **`git commit`** + simpan file **`.json`** hasil grademe
- ⇧ **`git push`** ke fork-mu

Sesudah battle: bedah 2 peserta terjauh — **plan · subagent · test**, bukan sekadar layar jadi.
Teaser: **w05 MCP** — rekap skor & DB otomatis.

### Konteks seri

| Minggu | Materi |
|---|---|
| w01 | plan mode · konvensi |
| w02 | subagent · code-review |
| w03 | test matrix · TDD |
| w04 | ⚔️ relay MPP |

---

## Navigasi deck (referensi presenter)

| Aksi | Tombol |
|---|---|
| Lanjut (fragment → slide) | `→` / `Space` / klik |
| Mundur | `←` |
| Slide pertama / terakhir | `Home` / `End` |
| Timer setup 10' (slide 4) | `S` |
| Timer coding 90' (slide 9) | `T` mulai · `R` reset |
| Extend coding +10' (maks +30') | `E` |
| Mode statis (screenshot) | `?static=1` |
