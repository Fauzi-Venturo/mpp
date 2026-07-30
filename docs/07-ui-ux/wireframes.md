# Wireframe / Sketsa Layar _(ID)_

Sketsa kasar (low-fidelity) untuk layar utama. Ini panduan tata letak, bukan desain
final; implementasi memakai komponen MUI v9.

## Web publik — pilih layanan & booking

```
┌───────────────────────────────────────────────┐
│  MPP  •  Beranda   Instansi   Status Antrean   │
├───────────────────────────────────────────────┤
│  Instansi: [ Dukcapil ▾ ]                       │
│  Jenis layanan:                                 │
│   ( ) Perekaman KTP        ~15 mnt              │
│   (•) Perpanjang KTP       ~10 mnt              │
│                                                 │
│  Syarat dokumen:                                │
│   ☑ KK asli   ☑ KTP lama   ☐ Surat pengantar    │
│                                                 │
│  Tanggal:  [ 05 Agu ]  Sisa kuota: 12           │
│                        [  Lanjut  ]             │
└───────────────────────────────────────────────┘
```

## Web publik — konfirmasi + QR

```
┌───────────────────────────────────────────────┐
│  Booking berhasil ✓                             │
│  Dukcapil • Perpanjang KTP • 05 Agu             │
│                                                 │
│        ┌───────────┐                            │
│        │  [ QR ]   │   Tunjukkan di kiosk        │
│        └───────────┘   untuk check-in            │
│                                                 │
│  [ Unduh QR ]   [ Kirim ke Email ]              │
└───────────────────────────────────────────────┘
```

## Kiosk — idle

```
┌───────────────────────────────────────────────┐
│              SELAMAT DATANG DI MPP              │
│                                                 │
│     ┌───────────────┐   ┌───────────────┐       │
│     │  CHECK-IN QR  │   │    WALK-IN    │       │
│     │   (pindai)    │   │  (tanpa QR)   │       │
│     └───────────────┘   └───────────────┘       │
└───────────────────────────────────────────────┘
```

## Display TV — pemanggilan

```
┌───────────────────────────────────────────────┐
│  DUKCAPIL                          10:24        │
│                                                 │
│        NOMOR          ┌─────────────────┐       │
│      ┌───────┐        │ Berikutnya:     │       │
│      │ A-014 │        │  A-015  A-016   │       │
│      └───────┘        │  A-017  A-018   │       │
│      LOKET 3          └─────────────────┘       │
│   🔊 "Nomor A-014, silakan ke loket 3"          │
├───────────────────────────────────────────────┤
│  ‹ running text / informasi ›                   │
└───────────────────────────────────────────────┘
```

## Aplikasi loket — panel petugas

```
┌───────────────────────────────────────────────┐
│  Loket 3 • Dukcapil • Pak Budi        [Tutup]   │
├───────────────────────────────────────────────┤
│  SEKARANG:  A-014   (dipanggil 1×)              │
│                                                 │
│  [ Panggil berikutnya ] [ Panggil ulang ]       │
│  [ Mulai ] [ Selesai ] [ Lewati ]               │
│  [ Transfer ] [ Hold ] [ Second service ]       │
├───────────────────────────────────────────────┤
│  Menunggu (8):  A-015  A-016  A-017 …           │
│  Dokumen (FO):  ☑ Lengkap                       │
└───────────────────────────────────────────────┘
```

## Front office — verifikasi

```
┌───────────────────────────────────────────────┐
│  Verifikasi Dokumen — Dukcapil                  │
├───────────────────────────────────────────────┤
│  Pemohon: Ibu Sari • Perpanjang KTP • A-014     │
│  Checklist:                                     │
│   ☑ KK asli     ☑ KTP lama     ☐ Pengantar RT   │
│                                                 │
│  Catatan: [_______________________]            │
│  ( ) Lengkap     (•) Tidak lengkap              │
│                         [  Simpan  ]            │
└───────────────────────────────────────────────┘
```

## Dashboard admin — monitoring

```
┌───────────────────────────────────────────────┐
│  Monitoring MPP (real-time)                     │
├───────────────────────────────────────────────┤
│  Instansi   Menunggu  Dilayani  T.Tunggu  Loket │
│  Dukcapil      8         3        ~12m     3/4   │
│  Imigrasi     14         5        ~20m     2/3   │
│  BPJS          3         1        ~6m      1/2   │
├───────────────────────────────────────────────┤
│  [ Reset ] [ Broadcast ] [ Buka/Tutup loket ]   │
└───────────────────────────────────────────────┘
```

> Sumber sketsa: dokumen desain artifact yang diberikan. Detail interaksi mengacu pada
> `functional-requirements.md` dan `queue-state-machine.md`.
