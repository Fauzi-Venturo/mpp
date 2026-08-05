'use client';

import { m } from 'framer-motion';

import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { Image } from 'src/components/image';
import { varFade, MotionViewport } from 'src/components/animate';

import { PROBLEM } from './home-data';
import { HomeHeading } from './home-heading';

// ----------------------------------------------------------------------

export function HomeProblem() {
  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 }, bgcolor: 'background.neutral' }}>
      <Container>
        <HomeHeading caption={PROBLEM.caption} title={PROBLEM.title} />

        <Box
          sx={{
            gap: { xs: 5, md: 8 },
            display: 'flex',
            alignItems: 'center',
            flexDirection: { xs: 'column', md: 'row' },
          }}
        >
          <Box sx={{ flexShrink: 0, width: { xs: 240, md: 320 } }}>
            <Image
              alt={PROBLEM.title}
              src={PROBLEM.image}
              ratio="1/1"
              sx={{ borderRadius: '50%', boxShadow: (theme) => theme.customShadows.z8 }}
            />
          </Box>

          <MotionViewport sx={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 3 }}>
            {PROBLEM.items.map((item) => (
              <Box component={m.div} variants={varFade('inUp')} key={item.title}>
                <Typography variant="h6" sx={{ mb: 0.5, color: 'primary.main' }}>
                  {item.title}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  {item.description}
                </Typography>
              </Box>
            ))}
          </MotionViewport>
        </Box>

        <Stack sx={{ mt: { xs: 5, md: 7 }, alignItems: 'center' }}>
          <Typography sx={{ maxWidth: 720, textAlign: 'center', color: 'text.secondary' }}>
            {PROBLEM.closing}
          </Typography>
        </Stack>
      </Container>
    </Box>
  );
}
