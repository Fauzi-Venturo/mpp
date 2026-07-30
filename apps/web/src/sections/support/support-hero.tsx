import type { BoxProps } from '@mui/material/Box';

import { varAlpha } from 'minimal-shared/utils';

import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';

import { CONFIG } from 'src/global-config';

// ----------------------------------------------------------------------

export function SupportHero({ sx, ...other }: BoxProps) {
  return (
    <Box
      component="section"
      sx={[
        (theme) => ({
          ...theme.mixins.bgGradient({
            images: [
              `linear-gradient(to bottom, ${varAlpha(theme.vars.palette.common.blackChannel, 0.8)}, ${varAlpha(theme.vars.palette.common.blackChannel, 0.8)})`,
              `url(${CONFIG.assetsDir}/assets/background/overlay-2.webp)`,
            ],
          }),
          px: 2.5,
          display: 'flex',
          alignItems: 'center',
          flexDirection: 'column',
          py: { xs: 10, md: 15 },
        }),
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
      {...other}
    >
      <Typography component="h1" variant="h2" sx={{ textAlign: 'center', color: 'common.white' }}>
        Selamat Datang di <br />
        <Box component="span" sx={{ color: 'primary.main' }}>
          FAQ
        </Box>{' '}
        Venturo
      </Typography>

      <Typography
        sx={{ mt: 3, maxWidth: 560, textAlign: 'center', color: 'common.white', opacity: 0.72 }}
      >
        Temukan jawaban seputar layanan, cara kerja tim, dan kerja sama dengan Venturo.
      </Typography>
    </Box>
  );
}
