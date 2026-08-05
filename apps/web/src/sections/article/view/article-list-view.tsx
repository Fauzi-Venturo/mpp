'use client';

import { useCallback } from 'react';

import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { useRouter, usePathname, useSearchParams } from 'src/routes/hooks';

import { parseArticleListParams } from 'src/lib/api';
import { useArticlesQuery, useArticleCategoriesQuery } from 'src/lib/api/use-articles';

import { ArticleList } from '../article-list';
import { ArticleFilters } from '../article-filters';

// ----------------------------------------------------------------------
// Data: TanStack Query (useArticlesQuery + useArticleCategoriesQuery) — sudah
// di-prefetch server oleh src/app/article/page.tsx via HydrationBoundary, jadi
// render pertama langsung berisi. Tanpa fallback statis: status loading/error/
// kosong ditangani ArticleList & ArticleFilters.
//
// URL (?page&search&category) is the single source of truth: filter changes
// navigate (replace for keystrokes, push for discrete clicks), the dynamic
// server page re-prefetches through the Next Data Cache, and TanStack Query
// bridges the transition with the previous page's data (no blank flashes).

export function ArticleListView() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const { page, search, category } = parseArticleListParams({
    page: searchParams.get('page'),
    search: searchParams.get('search'),
    category: searchParams.get('category'),
  });

  const articlesQuery = useArticlesQuery({ page, search, category });
  const categoriesQuery = useArticleCategoriesQuery();

  const buildUrl = useCallback(
    (updates: Record<string, string>) => {
      const params = new URLSearchParams(searchParams.toString());

      for (const [key, value] of Object.entries(updates)) {
        if (value) {
          params.set(key, value);
        } else {
          params.delete(key);
        }
      }

      const query = params.toString();

      return query ? `${pathname}?${query}` : pathname;
    },
    [pathname, searchParams]
  );

  const handleSearchChange = useCallback(
    (value: string) => {
      // Reset to page 1 on a new query; replace() keeps typing out of history.
      router.replace(buildUrl({ search: value, page: '' }), { scroll: false });
    },
    [router, buildUrl]
  );

  const handleCategoryChange = useCallback(
    (slug: string) => {
      router.push(buildUrl({ category: slug, page: '' }), { scroll: false });
    },
    [router, buildUrl]
  );

  // Nomor halaman dirender sebagai <a> asli (crawlable); navigasinya tetap
  // client-side lewat RouterLink.
  const buildPageHref = useCallback(
    (newPage: number) => buildUrl({ page: newPage > 1 ? String(newPage) : '' }),
    [buildUrl]
  );

  const { data, isPending, isError, isFetching, isPlaceholderData } = articlesQuery;

  const articles = data?.articles ?? [];
  const totalPages = data?.pagination.total_pages ?? 1;
  const isEmpty = !isPending && !isError && articles.length === 0;

  return (
    <Container sx={{ pt: { xs: 5, md: 8 } }}>
      <Typography variant="h2" component="h1">
        Article
      </Typography>

      <Typography sx={{ mt: 2, mb: { xs: 5, md: 8 }, color: 'text.secondary', maxWidth: 560 }}>
        Insights and updates on technology, software engineering, IT infrastructure, and artificial
        intelligence from the Venturo team.
      </Typography>

      <ArticleFilters
        search={search}
        category={category}
        categories={categoriesQuery.data ?? []}
        loadingCategories={categoriesQuery.isPending && !categoriesQuery.isError}
        onSearchChange={handleSearchChange}
        onCategoryChange={handleCategoryChange}
      />

      <Box
        sx={{
          mt: { xs: 5, md: 8 },
          ...(isFetching && isPlaceholderData && { opacity: 0.64, pointerEvents: 'none' }),
        }}
      >
        {isError && (
          <FeedbackBlock
            title="Unable to load articles"
            description="Something went wrong while contacting the server. Please try again later."
          />
        )}

        {isEmpty && (
          <FeedbackBlock
            title="No articles found"
            description={
              page > 1
                ? 'This page is out of range. Go back to the first page.'
                : search || category
                  ? 'No articles match your search or filter. Try different keywords.'
                  : 'There are no published articles yet. Please check back soon.'
            }
          />
        )}

        {!isError && !isEmpty && (
          <ArticleList
            articles={articles}
            loading={isPending}
            page={page}
            totalPages={totalPages}
            buildPageHref={buildPageHref}
          />
        )}
      </Box>
    </Container>
  );
}

// ----------------------------------------------------------------------

type FeedbackBlockProps = {
  title: string;
  description: string;
};

function FeedbackBlock({ title, description }: FeedbackBlockProps) {
  return (
    <Box sx={{ py: { xs: 8, md: 12 }, textAlign: 'center' }}>
      <Typography variant="h5" sx={{ mb: 1 }}>
        {title}
      </Typography>

      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {description}
      </Typography>
    </Box>
  );
}
