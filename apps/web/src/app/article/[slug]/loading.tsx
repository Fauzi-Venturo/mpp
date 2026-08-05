import Box from '@mui/material/Box';
import Skeleton from '@mui/material/Skeleton';
import Container from '@mui/material/Container';

// ----------------------------------------------------------------------

export default function Loading() {
  return (
    <Container sx={{ pt: { xs: 3, md: 5 }, pb: { xs: 10, md: 15 } }}>
      <Skeleton sx={{ width: 280, mb: { xs: 3, md: 5 } }} />

      <Skeleton sx={{ width: 0.8, height: 64 }} />
      <Skeleton sx={{ mt: 3, width: 0.6, height: 32 }} />

      <Box sx={{ my: 5, gap: 2, display: 'flex', alignItems: 'center' }}>
        <Skeleton variant="circular" sx={{ width: 48, height: 48, flexShrink: 0 }} />
        <Box sx={{ flexGrow: 1, maxWidth: 200 }}>
          <Skeleton sx={{ width: 0.7 }} />
          <Skeleton sx={{ width: 0.4 }} />
        </Box>
      </Box>

      <Skeleton
        variant="rectangular"
        sx={{ width: 1, height: 'auto', aspectRatio: '16/9', borderRadius: 3 }}
      />

      <Box sx={{ mt: 5 }}>
        {Array.from({ length: 5 }, (_, index) => (
          <Skeleton key={index} sx={{ width: index === 4 ? 0.5 : 1 }} />
        ))}
      </Box>
    </Container>
  );
}
