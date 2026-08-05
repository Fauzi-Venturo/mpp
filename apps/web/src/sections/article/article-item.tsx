import type { ArticleListItem } from 'src/lib/api';

import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Card from '@mui/material/Card';
import Avatar from '@mui/material/Avatar';
import Typography from '@mui/material/Typography';
import CardContent from '@mui/material/CardContent';
import ListItemText from '@mui/material/ListItemText';
import CardActionArea from '@mui/material/CardActionArea';

import { paths } from 'src/routes/paths';
import { RouterLink } from 'src/routes/components';

import { fDate } from 'src/utils/format-time';

import { Label } from 'src/components/label';
import { Image } from 'src/components/image';

// ----------------------------------------------------------------------

type ArticleItemProps = {
  article: ArticleListItem;
  /** Above-the-fold items: render eagerly (in the SSR HTML) with high fetch priority. */
  priority?: boolean;
};

export function ArticleItem({ article, priority = false }: ArticleItemProps) {
  const detailsHref = paths.article.details(article.slug);

  return (
    <Card sx={{ display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ position: 'relative' }}>
        {article.category && (
          <Label variant="filled" sx={{ top: 16, right: 16, zIndex: 9, position: 'absolute' }}>
            {article.category.name}
          </Label>
        )}

        <CardActionArea component={RouterLink} href={detailsHref}>
          {article.cover_url ? (
            <Image
              src={article.cover_url}
              alt={article.title}
              ratio="16/9"
              visibleByDefault={priority}
              slotProps={{ img: priority ? { fetchPriority: 'high' } : { loading: 'lazy' } }}
            />
          ) : (
            <Box sx={{ aspectRatio: '16/9', bgcolor: 'background.neutral' }} />
          )}
        </CardActionArea>
      </Box>

      <CardContent sx={{ gap: 1, display: 'flex', flexGrow: 1, flexDirection: 'column' }}>
        <Link
          component={RouterLink}
          href={detailsHref}
          color="inherit"
          variant="h6"
          sx={(theme) => ({ ...theme.mixins.maxLine({ line: 2 }) })}
        >
          {article.title}
        </Link>

        <Typography
          variant="body2"
          sx={(theme) => ({ ...theme.mixins.maxLine({ line: 2 }), color: 'text.secondary' })}
        >
          {article.excerpt}
        </Typography>

        <Box sx={{ mt: 'auto', pt: 2, gap: 1.5, display: 'flex', alignItems: 'center' }}>
          <Avatar alt={article.author}>{article.author.charAt(0).toUpperCase()}</Avatar>

          <ListItemText
            primary={article.author}
            secondary={fDate(article.published_at ?? article.created_at)}
            slotProps={{
              primary: { sx: { typography: 'subtitle2' } },
              secondary: { sx: { typography: 'caption', color: 'text.disabled' } },
            }}
          />
        </Box>
      </CardContent>
    </Card>
  );
}
