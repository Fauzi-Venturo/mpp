import Box from '@mui/material/Box';
import Skeleton from '@mui/material/Skeleton';
import Container from '@mui/material/Container';

// ----------------------------------------------------------------------
// Streams instantly while the server fetches the article list. Also opts the
// dynamic /article routes into Next 16 partial prefetching from <Link>.

export default function Loading() {
  return (
    <Container sx={{ pt: { xs: 5, md: 8 } }}>
      <Skeleton sx={{ width: 220, height: 64 }} />
      <Skeleton sx={{ mt: 2, mb: { xs: 5, md: 8 }, width: 0.5, maxWidth: 560 }} />

      <Box
        sx={{
          columnGap: 4,
          rowGap: { xs: 5, md: 8 },
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        }}
      >
        {Array.from({ length: 6 }, (_, index) => (
          <Box key={index}>
            <Skeleton
              variant="rectangular"
              sx={{ width: 1, height: 'auto', aspectRatio: '16/9' }}
            />
            <Skeleton sx={{ mt: 2, width: 0.9 }} />
            <Skeleton sx={{ width: 0.6 }} />
          </Box>
        ))}
      </Box>
    </Container>
  );
}
