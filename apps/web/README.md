# Venturo Skeleton — Next.js

Skeleton **company-profile / marketing site** milik Venturo — di-clone per client sebagai titik awal setiap project profil perusahaan. Basis template: Zone UI v4.6.0 (Minimals), sudah dirapikan, di-hardening, dan diintegrasikan ke backend Go **marketplace-be**.

**Stack:** Next.js 16 (App Router, RSC + ISR) · MUI 9 (theme CSS variables) · TypeScript strict · TanStack Query 5 · ky + zod (fetch boundary tervalidasi) · react-hook-form + zod · framer-motion · embla carousel.

## Apa yang sudah termasuk

- **3 vertical siap pakai** — `home` (landing), `article` (list ber-filter + detail ISR), `support` (FAQ dua kolom dari API), plus halaman utilitas (404/500, coming-soon, maintenance, blank).
- **Lapisan API standar** ([src/lib/api/](src/lib/api/)) — fetcher ky ber-timeout/retry, envelope backend di-unwrap, semua respons divalidasi zod, error dinormalkan jadi satu kelas `ApiError`, query keys + hooks TanStack. Setiap halaman ber-API punya **fallback statis** — site tidak pernah pecah saat backend mati.
- **SEO lengkap** — title template terpusat, OG/Twitter image, canonical + trailing-slash policy, sitemap dinamis, robots, manifest, JSON-LD (Organization, WebSite, Article, BreadcrumbList), webhook revalidate ISR ber-token.
- **Security & kualitas** — security headers (HSTS, nosniff, CSP report-only, dll.), env divalidasi zod saat build (fail-fast), tsconfig ketat (`noUncheckedIndexedAccess` dkk.), husky pre-commit/pre-push, Jenkinsfile CI, aksesibilitas WCAG AA (skip-link, reduced-motion, kontras).
- **Galeri komponen `/components`** — living docs semua komponen bersama untuk programmer & AI; selalu tampil saat `yarn dev`, di production di-gate env flag (deploy client → 404).
- **Docker production** — image standalone multi-stage, jalan di port 80 (lihat [docs/deployment.md](docs/deployment.md)).

## Prerequisites

- Node.js >= 22.12 (lihat `.nvmrc`)
- Yarn 1.22 (satu-satunya package manager — jangan pakai npm)
- Backend `marketplace-be` berjalan (default `http://localhost:8080`)

## Quick start

```sh
cp .env.example .env   # lalu isi variabelnya — lihat komentar di dalam file
yarn install
yarn dev               # http://localhost:8002
```

## Commands

```sh
yarn dev          # dev server (port 8002)
yarn build        # production build
yarn start        # serve production build (port 8002)

yarn tsc:check    # type-check (tidak ada test runner; ini gate utama)
yarn lint         # eslint
yarn fix:all      # lint:fix + prettier — jalankan sebelum commit
```

Gate otomatis: pre-commit menjalankan eslint+prettier pada file staged, pre-push menjalankan `tsc:check`, dan Jenkins menjalankan gate penuh (install → lint → format → typecheck → build).

Production via Docker: `docker compose up -d --build` (env dari `.env.prod`, site di port 80).

## Struktur folder

```
├── docs/                    # dokumentasi tim (resep, branding, deployment, distribusi)
├── public/
│   └── assets/venturo/      # aset brand (webp terkompres) — ganti per client
├── src/
│   ├── app/                 # App Router — page.tsx TIPIS: metadata + fetch + render View
│   │   ├── (home)/          # route group landing page (URL: /)
│   │   ├── (simple)/        # group layout compact: coming-soon, maintenance
│   │   ├── article/         # list (+ [slug]/ detail ISR + generateMetadata)
│   │   ├── support/         # halaman FAQ
│   │   ├── components/      # galeri komponen (di-gate middleware di production)
│   │   ├── error/           # pratinjau 404/500
│   │   ├── api/revalidate/  # webhook ISR ber-token untuk CMS/backend
│   │   ├── layout.tsx       # root: font lokal, provider, metadata default, og image
│   │   └── sitemap.ts / robots.ts / manifest.ts
│   │
│   ├── sections/            # UI per-vertical — SEMUA isi halaman ada di sini
│   │   ├── home/            #   <nama>-section.tsx  = blok presentasional
│   │   │   └── view/        #   view/<nama>-view.tsx = perangkai section (WAJIB)
│   │   ├── article/  support/  error/  coming-soon/  maintenance/
│   │   └── _examples/       # sumber galeri /components (satu-satunya prefix _)
│   │
│   ├── lib/
│   │   ├── api/             # LAPISAN DATA — articles.ts = referensi pola endpoint
│   │   │   ├── client.ts    #   ky instance + apiFetch + ApiError + envelope
│   │   │   ├── endpoints.ts #   peta path API (cermin paths.ts untuk backend)
│   │   │   ├── articles.ts  #   schema zod + fetcher + query keys (contoh kanonik)
│   │   │   └── use-*.ts     #   hooks 'use client' (sengaja di luar barrel)
│   │   ├── query/           # TanStack Query: QueryClient + provider (pola SSR)
│   │   ├── env.ts           # SATU-SATUNYA pembaca env — zod, fail-fast saat build
│   │   └── seo.ts           # helper JSON-LD (Organization, Article, Breadcrumb…)
│   │
│   ├── components/          # komponen bersama reusable (carousel, image, iconify,
│   │                        #   hook-form Field.*, markdown, animate, …) + barrel
│   ├── layouts/             # shell halaman: main/ (header+footer), simple/,
│   │                        #   nav-config-main.tsx (menu utama)
│   ├── theme/               # theme MUI CSS-vars: theme-config.ts (palette+font),
│   │                        #   theme-overrides.ts (tweak app — JANGAN edit core/)
│   ├── routes/              # paths.ts (SEMUA string URL) + hooks + RouterLink
│   ├── locales/             # i18n dormant (default id; language switcher disembunyikan)
│   ├── _mock/               # data sample HANYA untuk galeri + fallback statis
│   ├── assets/fonts/        # woff2 self-host (DM Sans + Barlow, next/font/local)
│   ├── utils/               # helper murni (format-time, dsb.)
│   ├── global-config.ts     # CONFIG — konsumen env.ts (jangan baca process.env)
│   └── middleware.ts        # gate galeri /components di production
│
├── Dockerfile / docker-compose.yml / .dockerignore   # deploy port 80, env .env.prod
├── Jenkinsfile              # CI gate
├── CLAUDE.md                # aturan arsitektur lengkap (sumber kebenaran konvensi)
└── .env.example             # template env ter-commit — copy ke .env / .env.prod
```

## Arsitektur dalam 30 detik

Setiap halaman mengikuti lapisan **page → view → section**:

1. `src/app/<route>/page.tsx` — Server Component tipis: metadata + fetch data + render satu View.
2. `src/sections/<vertical>/view/<nama>-view.tsx` — `'use client'`, menerima data sebagai props, merangkai section.
3. `src/sections/<vertical>/<nama>-section.tsx` — blok presentasional.

Tiga mode aliran data, pilih sesuai kebutuhan halaman:

| Mode | Kapan | Contoh |
|---|---|---|
| RSC props + fallback statis | konten marketing, SEO penting | home, support |
| Prefetch + HydrationBoundary | list interaktif (filter/paginasi) | /article |
| Pure RSC + ISR | detail per-slug + `generateMetadata` | /article/[slug] |

Aturan lengkap (route group, trailing slash, jebakan openGraph, aturan LCP, dsb.) ada di [CLAUDE.md](CLAUDE.md) — file itu juga yang dibaca AI assistant saat bekerja di repo ini.

## Model branch & distribusi

| Branch | Peran |
|---|---|
| `production` | skeleton hidup — semua perbaikan masuk sini |
| `template-default-4.6.0` | template Zone UI murni, beku — untuk diff & mengambil kembali bagian yang dihapus |

Project client dibuat dari repo ini sebagai **GitHub template repo**, lalu menambah remote `skeleton` untuk cherry-pick perbaikan dua arah — alur lengkapnya di [docs/distribution.md](docs/distribution.md), checklist branding-nya di [docs/branding.md](docs/branding.md).

## Dokumentasi tim

| Dokumen | Isi |
|---|---|
| [docs/recipes/add-a-page.md](docs/recipes/add-a-page.md) | Langkah menambah halaman: paths → section → view → layout → page → nav → sitemap |
| [docs/recipes/add-an-api-endpoint.md](docs/recipes/add-an-api-endpoint.md) | Menambah endpoint API meniru pola `articles.ts` + idiom fallback |
| [docs/branding.md](docs/branding.md) | Checklist branding per-client: logo, palette, font, konten, env |
| [docs/deployment.md](docs/deployment.md) | Env production, Docker (port 80, `.env.prod`), k8s, Jenkins, checklist pasca-deploy |
| [docs/distribution.md](docs/distribution.md) | Model distribusi: template repo, remote `skeleton`, alur cherry-pick perbaikan |
