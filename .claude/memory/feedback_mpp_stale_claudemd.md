---
name: feedback_mpp_stale_claudemd
description: CLAUDE.md di apps/api dan apps/web masih milik template asalnya, bukan MPP — jangan diikuti mentah-mentah
metadata:
  type: feedback
---

`apps/api/CLAUDE.md` dan `apps/web/CLAUDE.md` **ikut ter-vendor dari skeleton asalnya dan belum
disesuaikan ke MPP** (kondisi per 2026-08-05 — cek dulu, mungkin sudah diperbaiki).

- `apps/api/CLAUDE.md` mendeskripsikan dirinya sebagai *"Tuai Backend … Go 1.25"*, dan menyuruh
  `make dev` / `make db-setup` dari dalam `apps/api`. Root repo pakai `make api-dev`, dan
  `docs/03-architecture/tech-stack.md` menyebut **Go 1.26**.
- `apps/web/CLAUDE.md` mendeskripsikan *"Venturo's company-profile/marketing skeleton"* dengan
  vertikal `home`/`article`/`support` yang datanya dari backend `marketplace-be` — nol hubungan
  dengan kiosk / display TV / aplikasi loket milik MPP.

**Why:** dua file ini otomatis masuk konteks sesi. Diikuti mentah-mentah, hasilnya perintah
build salah, arsitektur salah, dan fitur dibangun untuk produk yang keliru.

**How to apply:** perlakukan bagian *pattern* keduanya sebagai valid (layout modular monolith
`{domain,dto,handler,repository,service}` di BE; pola `page → view → section`, `src/lib/api`,
`Field.*` RHF di FE — ini memang diwarisi MPP). Tapi bagian *identitas produk, perintah, dan
daftar fitur* sudah basi. **Sumber kebenaran untuk apa yang dibangun = `docs/`**, bukan CLAUDE.md
di `apps/`. Kalau menyentuh salah satunya, tawarkan sekalian memperbaruinya.

Lihat [[project_mpp_overview]].
