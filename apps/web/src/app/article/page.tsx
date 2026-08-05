import type { Metadata } from 'next';

import { Suspense } from 'react';
import { dehydrate, HydrationBoundary } from '@tanstack/react-query';

import { paths, pathWithSlash } from 'src/routes/paths';

import { getQueryClient } from 'src/lib/query';
import {
  articleKeys,
  getArticles,
  getArticleCategories,
  parseArticleListParams,
} from 'src/lib/api';

import { ArticleListView } from 'src/sections/article/view/article-list-view';

// ----------------------------------------------------------------------

const DESCRIPTION =
  'Insights and updates on technology, software engineering, IT infrastructure, and artificial intelligence from the Venturo team.';

type Props = {
  searchParams: Promise<{
    page?: string | string[];
    search?: string | string[];
    category?: string | string[];
  }>;
};

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { page, search, category } = parseArticleListParams(await searchParams);
  const isFiltered = Boolean(search || category);

  return {
    title: 'Article',
    description: DESCRIPTION,
    // Paginated pages are self-canonical; search/category variants canonicalize
    // to the root list and are kept out of the index.
    alternates: {
      canonical:
        !isFiltered && page > 1
          ? `${pathWithSlash(paths.article.root)}?page=${page}`
          : pathWithSlash(paths.article.root),
    },
    ...(isFiltered && { robots: { index: false, follow: true } }),
  };
}

export default async function Page({ searchParams }: Props) {
  const { page, search, category } = parseArticleListParams(await searchParams);

  const queryClient = getQueryClient();

  // Prefetch into the TanStack cache (through the Next Data Cache) so the
  // client view hydrates with data — full HTML for crawlers, zero client
  // refetch on load. prefetchQuery swallows errors: if the backend is down,
  // the page still renders and the view shows its error state after retrying.
  await Promise.all([
    queryClient.prefetchQuery({
      queryKey: articleKeys.list({ page, search, category }),
      queryFn: () => getArticles({ page, search, category }),
    }),
    queryClient.prefetchQuery({
      queryKey: articleKeys.categories,
      queryFn: () => getArticleCategories(),
    }),
  ]);

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <Suspense>
        <ArticleListView />
      </Suspense>
    </HydrationBoundary>
  );
}
