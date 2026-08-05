'use client';

import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { Iconify } from 'src/components/iconify';

import { FOCUS } from './home-data';
import { HomeHeading } from './home-heading';

// ----------------------------------------------------------------------

export function HomeFocus() {
  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 } }}>
      <Container>
        <HomeHeading title={FOCUS.caption} description={FOCUS.description} />

        <Card sx={{ p: { xs: 3, md: 5 }, boxShadow: (theme) => theme.customShadows.z16 }}>
          <Typography variant="h5" sx={{ mb: 3, textAlign: 'center', color: 'primary.main' }}>
            {FOCUS.highlight}
          </Typography>

          <Box
            sx={{
              rowGap: 1.5,
              columnGap: 3,
              display: 'flex',
              flexWrap: 'wrap',
              justifyContent: 'center',
            }}
          >
            {FOCUS.roles.map((role) => (
              <Box key={role} sx={{ gap: 1, display: 'flex', alignItems: 'center' }}>
                <Iconify width={20} icon="solar:check-circle-bold" sx={{ color: 'primary.main' }} />
                <Typography variant="body2">{role}</Typography>
              </Box>
            ))}
          </Box>

          <Box
            component="img"
            loading="lazy"
            alt={FOCUS.highlight}
            src={FOCUS.image}
            sx={{
              mt: { xs: 4, md: 6 },
              mx: 'auto',
              width: 1,
              maxWidth: 900,
              height: 'auto',
              display: 'block',
            }}
          />
        </Card>
      </Container>
    </Box>
  );
}
