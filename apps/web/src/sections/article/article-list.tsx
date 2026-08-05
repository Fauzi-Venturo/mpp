'use client';

import type { ArticleListItem } from 'src/lib/api';

import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import Skeleton from '@mui/material/Skeleton';
import CardContent from '@mui/material/CardContent';
import PaginationItem from '@mui/material/PaginationItem';
import Pagination, { paginationClasses } from '@mui/material/Pagination';

import { RouterLink } from 'src/routes/components';

import { ArticleItem } from './article-item';

// ----------------------------------------------------------------------

const SKELETON_COUNT = 6;
// First grid row is (partially) above the fold — render those covers eagerly.
const PRIORITY_COUNT = 3;

type ArticleListProps = {
  articles: ArticleListItem[];
  loading?: boolean;
  page: number;
  totalPages: number;
  /** Href untuk nomor halaman — pagination harus berupa <a> asli agar crawlable. */
  buildPageHref: (page: number) => string;
};

export function ArticleList({
  articles,
  loading = false,
  page,
  totalPages,
  buildPageHref,
}: ArticleListProps) {
  return (
    <>
      <Box
        sx={{
          columnGap: 4,
          rowGap: { xs: 5, md: 8 },
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        }}
      >
        {loading
          ? Array.from({ length: SKELETON_COUNT }, (_, index) => (
              <ArticleItemSkeleton key={index} />
            ))
          : articles.map((article, index) => (
              <ArticleItem
                key={article.slug}
                article={article}
                priority={index < PRIORITY_COUNT && page === 1}
              />
            ))}
      </Box>

      {totalPages > 1 && (
        <Pagination
          count={totalPages}
          page={page}
          renderItem={(item) =>
            item.page && !item.disabled ? (
              <PaginationItem component={RouterLink} href={buildPageHref(item.page)} {...item} />
            ) : (
              <PaginationItem {...item} />
            )
          }
          sx={{
            py: { xs: 8, md: 10 },
            [`& .${paginationClasses.ul}`]: { justifyContent: 'center' },
          }}
        />
      )}
    </>
  );
}

// ----------------------------------------------------------------------

function ArticleItemSkeleton() {
  return (
    <Card sx={{ display: 'flex', flexDirection: 'column' }}>
      <Skeleton variant="rectangular" sx={{ width: 1, height: 'auto', aspectRatio: '16/9' }} />

      <CardContent sx={{ gap: 1, display: 'flex', flexGrow: 1, flexDirection: 'column' }}>
        <Skeleton sx={{ width: 0.9 }} />
        <Skeleton sx={{ width: 0.6 }} />

        <Box sx={{ mt: 'auto', pt: 2, gap: 1.5, display: 'flex', alignItems: 'center' }}>
          <Skeleton variant="circular" sx={{ width: 40, height: 40, flexShrink: 0 }} />
          <Box sx={{ flexGrow: 1 }}>
            <Skeleton sx={{ width: 0.5 }} />
            <Skeleton sx={{ width: 0.3 }} />
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
}
