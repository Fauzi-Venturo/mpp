# Product Brief — Sistem Antrian MPP _(ID)_

## Latar belakang

Mal Pelayanan Publik (MPP) mengumpulkan banyak **instansi** (Dukcapil, Imigrasi,
BPJS, Samsat, kepolisian, perizinan daerah, dll.) dalam satu gedung. Warga sering
tidak tahu harus mengambil nomor di mana, menunggu tanpa kepastian, dan mengantre di
loket yang salah karena tiap instansi punya alur sendiri. Sistem antrian terpadu
menyatukan pengambilan nomor, pemanggilan, dan pemantauan lintas instansi dalam satu
platform.

## Visi

> Satu sistem antrian terpadu untuk seluruh instansi di MPP: warga bisa mendaftar dari
> mana saja (WhatsApp/web/walk-in), tahu pasti kapan dilayani, dan petugas melayani
> dengan alur yang adil, cepat, dan terukur.

## Tujuan produk

1. **Kurangi waktu tunggu & antre fisik** — booking online + estimasi waktu layanan.
2. **Transparansi antrian** — nomor & pemanggilan tampil real-time di TV dan kanal warga.
3. **Distribusi beban loket adil** — alokasi otomatis ke loket yang paling lama idle.
4. **Kelengkapan dokumen sejak awal** — syarat dokumen diinformasikan saat registrasi
   dan diverifikasi front office sebelum ke loket.
5. **Data & pelaporan** — lama layanan, jumlah dilayani, no-show, per instansi/loket/hari.

## Ruang lingkup (in scope)

- Katalog instansi & jenis layanan (dua tingkat) beserta syarat dokumen dan estimasi durasi.
- Registrasi antrian: WhatsApp AI agent, website publik, dan walk-in di kiosk.
- Kuota booking per tanggal per instansi.
- Check-in QR (token sekali pakai) di kiosk + cetak tiket termal.
- Engine antrian real-time: satu antrean per layanan, alokasi loket idle-terlama,
  mode FIFO atau Booking-prioritas, pemanggilan 3× lalu skip, reset harian 00:00.
- Display TV dengan pemanggilan suara TTS Bahasa Indonesia (offline).
- Aplikasi loket untuk petugas (panggil/lewati/mulai/selesai/transfer/hold).
- Layanan kedua otomatis (second service) tanpa daftar ulang.
- Verifikasi dokumen front office.
- Dashboard admin (master data, konfigurasi, monitoring) & pelaporan.
- RBAC 4 peran: admin, supervisor, front office, petugas loket.

## Di luar lingkup (out of scope, tahap awal)

- Pembayaran/retribusi online.
- Integrasi dua arah dengan sistem inti tiap instansi (SIAK, dsb.) — cukup katalog manual.
- Aplikasi mobile native (kanal warga cukup WhatsApp + web responsif).
- Antrian multi-gedung/multi-MPP (fokus satu gedung dulu; desain data menyiapkan perluasan).

## Stakeholder

| Stakeholder            | Kepentingan                                                        |
|------------------------|--------------------------------------------------------------------|
| Warga / pemohon        | Daftar mudah, tahu kepastian antrean, dokumen jelas                |
| Petugas loket          | Alur pemanggilan sederhana, beban adil                             |
| Front office           | Verifikasi dokumen cepat sebelum warga ke loket                    |
| Supervisor instansi    | Pantau performa loket & antrean instansinya                        |
| Admin MPP              | Kelola instansi/layanan/loket, konfigurasi, laporan lintas instansi |
| Pengelola gedung MPP   | Ketertiban, kapasitas, citra pelayanan                             |

## Kriteria sukses (indikatif)

- Rata-rata waktu tunggu turun signifikan vs sistem manual.
- >70% registrasi lewat kanal online (WA/web) dalam 3 bulan.
- Tingkat no-show terukur & menurun berkat pengingat + kuota.
- Ketersediaan sistem pada jam operasional ≥ 99,5%.
