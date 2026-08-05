# Product Requirements Document (PRD) — MPP _(ID)_

## 1. Ringkasan

Sistem antrian terpadu untuk Mal Pelayanan Publik yang menampung banyak instansi.
Mencakup registrasi multi-kanal (WhatsApp AI agent, web, walk-in), check-in QR,
engine antrian real-time, display TV bersuara (TTS Bahasa Indonesia), aplikasi loket,
verifikasi dokumen front office, dashboard admin, dan pelaporan.

## 2. Masalah yang dipecahkan

- Antrean manual: warga tidak tahu di mana ambil nomor & berapa lama menunggu.
- Beban loket tidak merata; sebagian loket menganggur, sebagian menumpuk.
- Dokumen tidak lengkap baru ketahuan di loket → bolak-balik, antrean terbuang.
- Tidak ada data terukur untuk evaluasi pelayanan lintas instansi.

## 3. Persona pengguna

Lihat [`user-roles-personas.md`](./user-roles-personas.md). Ringkas: **Warga/pemohon**,
**Petugas loket**, **Front office**, **Supervisor**, **Admin**.

## 4. Prinsip produk

1. **Real-time & transparan** — setiap perubahan nomor tampil seketika di TV & kanal warga.
2. **Adil** — alokasi loket idle-terlama; mode antrean dapat dikonfigurasi per instansi.
3. **Anti no-show** — kuota, pengingat, dan aturan panggil-3×-lalu-skip.
4. **Dokumen jelas di depan** — syarat dokumen menempel di jenis layanan sejak registrasi.
5. **Offline-tolerant di lokasi** — pemanggilan suara TV tidak bergantung internet.

## 5. Alur utama (happy path)

```
Registrasi (WA/Web/Walk-in) → Booking/nomor + syarat dokumen
   → Datang ke MPP → Check-in QR di kiosk → cetak tiket
   → (opsional) Verifikasi dokumen Front Office
   → Masuk antrean layanan → Dipanggil di TV (suara) → Dilayani di loket
   → Selesai (atau Second Service otomatis) → Data masuk laporan
```

## 6. Kebutuhan fungsional (tingkat tinggi)

Rincian per modul di [`functional-requirements.md`](./functional-requirements.md).
18 modul: master instansi, master jenis layanan, master loket, kuota booking,
registrasi WA AI agent, registrasi web, registrasi walk-in, check-in QR, engine
antrian, display TV+TTS, aplikasi loket, second service, verifikasi FO, user & RBAC,
dashboard admin & monitoring, pelaporan, notifikasi, konfigurasi & audit.

## 7. Kebutuhan non-fungsional

Lihat [`non-functional-requirements.md`](./non-functional-requirements.md): performa,
latency real-time, ketersediaan, keamanan/PII, skalabilitas, aksesibilitas, i18n,
ketahanan offline TTS.

## 8. Asumsi & batasan

- Satu gedung MPP pada tahap awal; desain data menyiapkan multi-gedung.
- Katalog instansi/layanan dikelola manual (tanpa integrasi sistem inti instansi).
- Kanal warga: WhatsApp + web responsif (tanpa app native).
- Backend Go (venturo-skeleton-go), Frontend Next.js (venturo-skeleton-next.js),
  DB PostgreSQL/Supabase, Redis wajib untuk engine antrian.

## 9. Metrik keberhasilan

- Rata-rata & p90 waktu tunggu; lama layanan per jenis layanan.
- Proporsi registrasi online vs walk-in.
- Tingkat no-show; jumlah dilayani per loket/instansi/hari.
- Ketersediaan sistem pada jam operasional.

## 10. Rilis bertahap

Lihat [`../08-roadmap/delivery-plan.md`](../08-roadmap/delivery-plan.md).
