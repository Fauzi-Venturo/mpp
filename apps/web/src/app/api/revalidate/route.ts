import { z } from 'zod';
import { NextResponse } from 'next/server';
import { revalidateTag } from 'next/cache';

import { env } from 'src/lib/env';
import { FAQ_TAG, ARTICLES_TAG, SITE_CONTENT_TAG, ARTICLE_CATEGORIES_TAG } from 'src/lib/api';

// ----------------------------------------------------------------------
// Webhook invalidasi ISR on-demand — dipanggil backend/CMS saat konten
// berubah, supaya editor tidak menunggu revalidate berkala (artikel 300 dtk,
// sitemap 3600 dtk).
//
//   POST /api/revalidate/   ← trailing slash WAJIB (trailingSlash: true
//                              me-redirect 308 varian tanpa slash)
//   Authorization: Bearer <REVALIDATE_TOKEN>   (server-only env, min 16 char)
//   Body: { "tags": ["articles"], "slugs": ["judul-artikel"] }
//
// `tags`  — nama cache tag yang dikenal (lihat ALLOWED_TAGS).
// `slugs` — slug artikel; di-map ke tag per-artikel `article:<slug>`.
// Endpoint menjawab 503 bila REVALIDATE_TOKEN tidak dikonfigurasi.

const ALLOWED_TAGS = new Set([ARTICLES_TAG, ARTICLE_CATEGORIES_TAG, FAQ_TAG, SITE_CONTENT_TAG]);

const bodySchema = z
  .object({
    tags: z.array(z.string().min(1)).optional(),
    slugs: z.array(z.string().min(1)).optional(),
  })
  .refine((body) => body.tags?.length || body.slugs?.length, {
    message: 'Isi minimal satu dari `tags` atau `slugs`.',
  });

export async function POST(request: Request) {
  if (!env.REVALIDATE_TOKEN) {
    return NextResponse.json(
      { error: 'REVALIDATE_TOKEN tidak dikonfigurasi — endpoint nonaktif.' },
      { status: 503 }
    );
  }

  const authorization = request.headers.get('authorization');

  if (authorization !== `Bearer ${env.REVALIDATE_TOKEN}`) {
    return NextResponse.json({ error: 'Unauthorized.' }, { status: 401 });
  }

  const parsed = bodySchema.safeParse(await request.json().catch(() => null));

  if (!parsed.success) {
    return NextResponse.json(
      { error: 'Body tidak valid.', issues: parsed.error.issues },
      { status: 400 }
    );
  }

  const { tags = [], slugs = [] } = parsed.data;

  const unknown = tags.filter((tag) => !ALLOWED_TAGS.has(tag));
  if (unknown.length) {
    return NextResponse.json(
      { error: `Tag tidak dikenal: ${unknown.join(', ')}`, allowed: [...ALLOWED_TAGS] },
      { status: 400 }
    );
  }

  const revalidated = [...tags, ...slugs.map((slug) => `article:${slug}`)];
  // Next 16: argumen ke-2 wajib — 'max' = purge dengan perilaku SWR standar.
  revalidated.forEach((tag) => revalidateTag(tag, 'max'));

  return NextResponse.json({ revalidated, now: Date.now() });
}
