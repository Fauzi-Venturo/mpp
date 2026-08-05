# Peran Pengguna & Persona _(ID)_

## Peran & hak akses (RBAC)

Empat peran internal (petugas MPP) + pengguna publik (warga). Detail matriks izin di
[`../06-security/rbac-matrix.md`](../06-security/rbac-matrix.md).

| Peran            | Ringkasan hak akses |
|------------------|---------------------|
| **Admin**        | Kelola seluruh master data (instansi, layanan, loket, kuota), konfigurasi sistem, pengguna & peran, akses semua laporan, monitoring lintas instansi. |
| **Supervisor**   | Pantau antrean & performa loket **pada instansi yang ditugaskan**; kontrol operasional terbatas (buka/tutup loket, reset, broadcast); akses laporan instansinya. |
| **Front office** | Verifikasi kelengkapan dokumen pemohon; menahan/melepas pemohon ke antrean loket. |
| **Petugas loket**| Operasikan loket: panggil/ulang/skip/mulai/selesai/transfer/hold, picu second service, lihat detail & checklist pemohon. |
| **Warga (publik)** | Registrasi (WA/web/walk-in), check-in QR, lihat status antrean; tanpa akun internal. |

## Persona

### Persona 1 — Warga/Pemohon: "Ibu Sari" (35)
- **Konteks:** ingin memperpanjang KTP; sibuk, ingin kepastian jadwal.
- **Kebutuhan:** daftar dari rumah lewat WhatsApp, tahu syarat dokumen, tidak antre lama.
- **Interaksi:** chat WA AI agent → booking tanggal → terima QR → check-in di kiosk → dipanggil.
- **Sukses jika:** datang tepat slot, dokumen lengkap, langsung terlayani.

### Persona 2 — Petugas Loket: "Pak Budi" (28)
- **Konteks:** melayani puluhan warga/hari di loket Dukcapil.
- **Kebutuhan:** alur pemanggilan cepat & jelas; tidak salah panggil; beban adil.
- **Interaksi:** login → pilih loket → panggil berikutnya → layani → selesai; skip bila no-show.
- **Sukses jika:** aksi minim klik, nomor & suara sinkron dengan TV.

### Persona 3 — Front Office: "Mbak Rina" (30)
- **Konteks:** meja depan; memastikan dokumen warga lengkap sebelum ke loket.
- **Kebutuhan:** checklist syarat per jenis layanan; cara cepat menandai lengkap/kurang.
- **Interaksi:** lihat antrean verifikasi → cek dokumen → tandai lengkap → warga lanjut ke loket.
- **Sukses jika:** mengurangi bolak-balik di loket akibat dokumen kurang.

### Persona 4 — Supervisor Instansi: "Bu Dewi" (42)
- **Konteks:** memantau kinerja loket Imigrasi.
- **Kebutuhan:** lihat antrean berjalan, beban tiap loket, waktu tunggu; intervensi bila menumpuk.
- **Interaksi:** dashboard instansi → buka loket tambahan / reset / broadcast informasi.
- **Sukses jika:** antrean seimbang, waktu tunggu terkendali.

### Persona 5 — Admin MPP: "Pak Andi" (38)
- **Konteks:** mengelola seluruh sistem MPP.
- **Kebutuhan:** atur instansi/layanan/loket/kuota, konfigurasi mode antrean, laporan lintas instansi.
- **Interaksi:** panel admin → kelola master → tetapkan kuota → tinjau laporan mingguan.
- **Sukses jika:** konfigurasi mudah, data pelaporan akurat untuk evaluasi.
