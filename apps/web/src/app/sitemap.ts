import type { MetadataRoute } from 'next';

import { paths, pathWithSlash } from 'src/routes/paths';

import { getArticles } from 'src/lib/api';
import { CONFIG } from 'src/global-config';

// ----------------------------------------------------------------------
// sitemap.ts is cached by default (frozen at build) — revalidate keeps it
// fresh against the backend. URLs carry a trailing slash to match the
// canonical form served under `trailingSlash: true`.

export const revalidate = 3600;

const PAGE_LIMIT = 100;
// 50 pages × 100 = 5.000 artikel; Google's hard cap is 50.000 URLs per file —
// shard with generateSitemaps() long before that becomes a concern.
const MAX_PAGES = 50;

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const staticEntries: MetadataRoute.Sitemap = [
    { url: `${CONFIG.siteUrl}${pathWithSlash(paths.home)}` },
    { url: `${CONFIG.siteUrl}${pathWithSlash(paths.article.root)}` },
    { url: `${CONFIG.siteUrl}${pathWithSlash(paths.support)}` },
  ];

  const articleEntries: MetadataRoute.Sitemap = [];

  try {
    let page = 1;
    let totalPages = 1;

    do {
      // Matches the segment's revalidate — the route regenerates at the
      // lowest fetch revalidate, so the default 60s would override 3600s.
      const { articles, pagination } = await getArticles(
        { page, limit: PAGE_LIMIT },
        { revalidate: 3600 }
      );

      totalPages = pagination.total_pages;

      articleEntries.push(
        ...articles.map((article) => {
          const lastModified = article.updated_at ?? article.published_at;

          return {
            url: `${CONFIG.siteUrl}${pathWithSlash(paths.article.details(article.slug))}`,
            ...(lastModified && { lastModified }),
          };
        })
      );

      page += 1;
    } while (page <= totalPages && page <= MAX_PAGES);
  } catch {
    // Backend unreachable — still serve the static entries instead of a 500.
  }

  return [...staticEntries, ...articleEntries];
}
