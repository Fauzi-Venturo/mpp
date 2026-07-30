import type { Metadata } from 'next';

import { notFound } from 'next/navigation';

import { paths, pathWithSlash } from 'src/routes/paths';

import { CONFIG } from 'src/global-config';
import { ApiError, getArticle, getArticles } from 'src/lib/api';
import { articleJsonLd, toJsonLdScript, breadcrumbListJsonLd } from 'src/lib/seo';

import { ArticleDetailsView } from 'src/sections/article/view/article-details-view';

// ----------------------------------------------------------------------
// ISR: nothing is prerendered at build time (no API dependency during build);
// each slug is statically generated on first request, then regenerated in the
// background at most every 5 minutes. The empty generateStaticParams is what
// opts the route into static generation at runtime — do not remove it.

export const revalidate = 300;

export function generateStaticParams() {
  return [];
}

type Props = {
  params: Promise<{ slug: string }>;
};

// Only a definitive 404 from the API becomes notFound() — a transient backend
// failure must throw (error boundary) so ISR keeps serving the last good page
// instead of replacing it with a cached 404.
async function fetchArticle(slug: string) {
  try {
    return await getArticle(slug);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;

  let article;
  try {
    article = await fetchArticle(slug);
  } catch {
    // Transient backend failure: fall back to minimal metadata — throwing
    // from generateMetadata bypasses error.tsx and yields a bare 500. The
    // page body re-throws and gets the branded error boundary instead.
    return { title: 'Article' };
  }

  if (!article) {
    notFound();
  }

  const canonical = pathWithSlash(paths.article.details(slug));

  return {
    title: article.title,
    description: article.excerpt,
    alternates: { canonical },
    openGraph: {
      type: 'article',
      url: canonical,
      siteName: CONFIG.appName,
      title: article.title,
      description: article.excerpt,
      publishedTime: article.published_at ?? undefined,
      modifiedTime: article.updated_at ?? article.published_at ?? undefined,
      ...(article.cover_url && { images: [{ url: article.cover_url, alt: article.title }] }),
    },
    twitter: {
      card: 'summary_large_image',
      title: article.title,
      description: article.excerpt,
      ...(article.cover_url && { images: [article.cover_url] }),
    },
  };
}

export default async function Page({ params }: Props) {
  const { slug } = await params;

  // Deduped with generateMetadata's call via Next.js request memoization.
  const article = await fetchArticle(slug);

  if (!article) {
    notFound();
  }

  // revalidate matches the segment's 300s — the route regenerates at the
  // LOWEST fetch revalidate, so the default 60s would silently win here.
  const relatedArticles = await getArticles(
    { limit: 4, category: article.category?.slug },
    { revalidate: 300 }
  )
    .then((result) => result.articles.filter((item) => item.slug !== slug).slice(0, 3))
    .catch(() => []);

  // Article + BreadcrumbList (breadcrumb visualnya dirender di view).
  const jsonLd = [
    articleJsonLd(article, `${CONFIG.siteUrl}${pathWithSlash(paths.article.details(slug))}`),
    breadcrumbListJsonLd([
      { name: 'Home', url: `${CONFIG.siteUrl}/` },
      { name: 'Article', url: `${CONFIG.siteUrl}${pathWithSlash(paths.article.root)}` },
      { name: article.title },
    ]),
  ];

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: toJsonLdScript(jsonLd) }}
      />

      <ArticleDetailsView article={article} relatedArticles={relatedArticles} />
    </>
  );
}
