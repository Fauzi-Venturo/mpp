# Inventaris Layar (Screens Inventory) _(ID)_

Daftar layar per aplikasi. Semua dibangun di `apps/web` (Next.js App Router).

## A. Web Publik (Warga)

| Layar | Fungsi | Akses |
|-------|--------|-------|
| Beranda | Info MPP, cara daftar, status antrean berjalan | publik |
| Pilih instansi | Daftar instansi + pencarian | publik |
| Pilih jenis layanan | Layanan + **syarat dokumen** + estimasi durasi | publik |
| Pilih tanggal (booking) | Kalender + **sisa kuota** per tanggal | publik |
| Data pemohon | Form (nama, kontak, NIK bila perlu) — RHF + Zod | publik |
| Konfirmasi & QR | Ringkasan + **QR check-in** (unduh/email) | publik |
| Cek/kelola booking | Lihat status, batalkan | publik (token) |
| Status antrean publik | Nomor sedang dilayani per instansi (real-time) | publik |

## B. Kiosk (Lokasi)

| Layar | Fungsi |
|-------|--------|
| Idle / sambutan | Tombol besar: **Check-in QR** / **Walk-in** |
| Check-in QR | Pindai QR → hasil (nomor + ETA) → **cetak tiket** |
| Walk-in | Pilih instansi → layanan (+ syarat) → konfirmasi → cetak tiket |
| Error/bantuan | Pesan jelas (token dipakai/kedaluwarsa, printer error) + fallback |

## C. Display TV (Mini PC → 3 TV)

| Layar | Fungsi |
|-------|--------|
| Display pemanggilan | Nomor dipanggil (besar) + loket + daftar berikutnya |
| Running text / media | Informasi di sela pemanggilan (broadcast admin) |
| (Audio) | **Antrean suara TTS bersama** — dipimpin satu window |

## D. Aplikasi Loket (Petugas)

| Layar | Fungsi |
|-------|--------|
| Login & pilih loket | Buka sesi loket |
| Panel antrean | Nomor sekarang, aksi: **Panggil berikutnya / Ulang / Lewati / Mulai / Selesai** |
| Aksi lanjut | **Transfer**, **Hold/Resume**, **Second Service** |
| Detail pemohon | Data + **checklist dokumen** (hasil FO) |
| Statistik saya | Jumlah dilayani & rata-rata durasi hari ini |

## E. Front Office

| Layar | Fungsi |
|-------|--------|
| Antrean verifikasi | Daftar pemohon menunggu verifikasi |
| Verifikasi dokumen | Checklist syarat → **Lengkap / Tidak lengkap** + catatan |
| (Opsional) Unggah dokumen | Scan/unggah sesuai kebijakan PII |

## F. Dashboard Admin

| Layar | Fungsi |
|-------|--------|
| Ringkasan monitoring | Status antrean seluruh instansi/loket (real-time) |
| Master instansi | CRUD instansi (prefix, jam, mode antrean) |
| Master jenis layanan | CRUD layanan + **syarat dokumen** + estimasi durasi |
| Master loket | CRUD loket + mapping layanan + status |
| Kuota booking | Atur kuota per tanggal/instansi/layanan |
| Pengguna & peran | Kelola akun + RBAC |
| Konfigurasi | Mode antrean, format nomor, jam, teks TTS, template pesan |
| Kontrol operasional | Buka/tutup loket, reset manual, broadcast ke TV |
| Laporan | Dilayani, no-show, waktu tunggu/lama layanan; filter + ekspor |
| Audit log | Riwayat aksi sensitif |

## G. Supervisor

Subset dashboard admin **terbatas pada instansi yang ditugaskan**: monitoring instansi,
kontrol operasional loket instansinya, laporan instansi.

## Catatan UI

- Komponen memakai **MUI v9** (tema Zone UI bawaan skeleton); form **react-hook-form +
  Zod**; data **TanStack Query + ky**.
- Perlu **menambah lapisan auth** di FE (skeleton belum ada auth): token store + hook
  `ky` `beforeRequest`/`afterResponse`, serta middleware proteksi rute internal.
- Layar kiosk & TV dioptimalkan layar sentuh / jarak jauh; rute di-gate via API-key,
  bukan login user.
