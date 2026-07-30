'use client';

import type { Article, ArticleListItem } from 'src/lib/api';

import Box from '@mui/material/Box';
import Avatar from '@mui/material/Avatar';
import Divider from '@mui/material/Divider';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import ListItemText from '@mui/material/ListItemText';

import { paths } from 'src/routes/paths';

import { fDate } from 'src/utils/format-time';

import { Label } from 'src/components/label';
import { Image } from 'src/components/image';
import { Markdown } from 'src/components/markdown';
import { CustomBreadcrumbs } from 'src/components/custom-breadcrumbs';

import { ArticleItem } from '../article-item';

// ----------------------------------------------------------------------
// Data: RSC props dari src/app/article/[slug]/page.tsx (pure RSC + ISR 300
// dtk). `article` dijamin ada — slug tak dikenal sudah notFound() di page.
// relatedArticles opsional: fetch gagal → [] → bagian terkait tak dirender.

type ArticleDetailsViewProps = {
  article: Article;
  relatedArticles?: ArticleListItem[];
};

export function ArticleDetailsView({ article, relatedArticles = [] }: ArticleDetailsViewProps) {
  return (
    <Container sx={{ pt: { xs: 3, md: 5 }, pb: { xs: 10, md: 15 } }}>
      <CustomBreadcrumbs
        links={[
          { name: 'Home', href: '/' },
          { name: 'Article', href: paths.article.root },
          { name: article.title },
        ]}
        sx={{ mb: { xs: 3, md: 5 } }}
      />

      <Box sx={{ display: 'flex', flexDirection: 'column' }}>
        {article.category && (
          <Label variant="soft" color="info" sx={{ mb: 2, alignSelf: 'flex-start' }}>
            {article.category.name}
          </Label>
        )}

        <Typography variant="h2" component="h1">
          {article.title}
        </Typography>

        {article.excerpt && (
          <Typography variant="h5" component="p" sx={{ mt: 3, color: 'text.secondary' }}>
            {article.excerpt}
          </Typography>
        )}

        <Box sx={{ my: 5, gap: 2, display: 'flex', alignItems: 'center' }}>
          <Avatar alt={article.author} sx={{ width: 48, height: 48 }}>
            {article.author.charAt(0).toUpperCase()}
          </Avatar>

          <ListItemText
            primary={article.author}
            secondary={fDate(article.published_at ?? article.created_at)}
            slotProps={{
              primary: { sx: { typography: 'subtitle2' } },
              secondary: { sx: { typography: 'caption', color: 'text.disabled' } },
            }}
          />
        </Box>

        {article.cover_url && (
          <Image
            src={article.cover_url}
            alt={article.title}
            ratio="16/9"
            visibleByDefault
            slotProps={{ img: { fetchPriority: 'high' } }}
            sx={{ borderRadius: 3 }}
          />
        )}

        <Markdown content={article.content} sx={{ mt: 5 }} />
      </Box>

      {relatedArticles.length > 0 && (
        <>
          <Divider sx={{ mt: { xs: 8, md: 10 }, mb: { xs: 5, md: 8 } }} />

          <Typography variant="h4" sx={{ mb: 5 }}>
            Related articles
          </Typography>

          <Box
            sx={{
              columnGap: 4,
              rowGap: { xs: 5, md: 8 },
              display: 'grid',
              gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
            }}
          >
            {relatedArticles.map((item) => (
              <ArticleItem key={item.slug} article={item} />
            ))}
          </Box>
        </>
      )}
    </Container>
  );
}
