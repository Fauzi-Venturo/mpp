# Rencana Pengerjaan (Delivery Plan) _(ID)_

Fase inkremental. Fase 0 = deliverable saat ini (dokumentasi + monorepo). Fase
berikutnya membangun fitur MPP di atas template yang sudah di-vendor.

## Fase 0 — Fondasi (SELESAI di branch ini)
- Monorepo (`apps/api` + `apps/web` di-vendor, module Go di-rename, orkestrasi root).
- Dokumentasi requirement lengkap (`docs/`).
- Catatan Postgres/Supabase + Redis.
- **Keluaran:** repo siap dikembangkan, dokumen acuan pembangunan.

## Fase 1 — Master data & fondasi backend MPP
- Skema `mpp` + migrasi/seeder (`migrations/mpp/`, `seeders/mpp/`).
- Modul: `instansi`, `layanan` (+ `syarat_dokumen`), `loket` (+ mapping), `kuota`.
- Seed **4 role RBAC** MPP + izin (lihat `../06-security/rbac-matrix.md`).
- Admin FE: CRUD master + auth layer (token store + hook `ky`).
- **Keluaran:** katalog & konfigurasi dapat dikelola admin.

## Fase 2 — Registrasi & check-in
- Modul `booking` (web) + `availability`/kuota atomik.
- Penerbitan **QR token sekali pakai**; notifikasi konfirmasi (email/WA dasar).
- Modul `checkin` + walk-in (kiosk) + payload cetak tiket termal.
- FE: alur booking web, kiosk (check-in & walk-in).
- **Keluaran:** warga dapat booking dan check-in; nomor terbentuk.

## Fase 3 — Engine antrian + realtime + TV/TTS
- Engine: counter Redis, aliran per layanan, **alokasi idle-terlama**, mode FIFO/booking,
  **reset harian 00:00** (worker), ETA.
- **State machine** lengkap (call/recall 3×/skip/hold/transfer).
- **WebSocket hub** + Redis pub/sub; event `call.created`/`queue.updated` dll.
- FE: **display TV** (mini PC, 3 window, antrean audio bersama) + **TTS offline**.
- **Keluaran:** antrean real-time & pemanggilan suara berjalan.

## Fase 4 — Aplikasi loket & front office
- Modul `loket_ops` (panggil/ulang/mulai/selesai/skip/hold/transfer) + `serving_session`.
- **Second service** (`QUEUED_NEXT`, tanpa daftar ulang).
- Modul `fo` (verifikasi dokumen + checklist) → tampil di loket.
- FE: aplikasi loket + front office.
- **Keluaran:** pelayanan end-to-end di loket dengan verifikasi FO.

## Fase 5 — WhatsApp AI agent
- Webhook WA + orkestrasi **LLM agent** (tools: catalog, availability, create/cancel booking).
- Slot-filling, tampil syarat dokumen, kirim QR; fallback menu terstruktur.
- Pengingat H-1/hari-H + "sebentar lagi dipanggil".
- **Keluaran:** registrasi percakapan via WhatsApp.

## Fase 6 — Admin lanjutan, monitoring & pelaporan
- Dashboard monitoring real-time lintas instansi; kontrol operasional & broadcast.
- Pelaporan (dilayani, no-show, waktu tunggu/lama layanan) + ekspor CSV/Excel.
- Audit log lengkap; penyempurnaan konfigurasi (format nomor, teks TTS, template pesan).
- **Keluaran:** operasional & evaluasi berbasis data.

## Catatan lintas fase
- **Kontrak API:** finalisasi OpenAPI di `packages/api-contract`, sinkron dengan
  `apps/web/src/lib/api/endpoints.ts`.
- **Auth FE:** ditambahkan mulai Fase 1 (skeleton FE belum punya auth).
- **Redis wajib** sejak Fase 3 (dan sudah dipakai RBAC sejak Fase 1).
- **Supabase** sebagai host Postgres sementara; sediakan Redis terpisah.
- Uji tiap fase sesuai bagian "Verification" pada masing-masing dokumen terkait.
