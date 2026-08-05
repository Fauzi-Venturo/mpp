'use client';

import { m } from 'framer-motion';

import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { varFade, MotionViewport } from 'src/components/animate';

import { SOLUTION } from './home-data';

// ----------------------------------------------------------------------

export function HomeSolution() {
  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 }, bgcolor: 'background.neutral' }}>
      <Container>
        <Card sx={{ overflow: 'hidden', boxShadow: (theme) => theme.customShadows.z16 }}>
          <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' } }}>
            <Box sx={{ flex: 1, p: { xs: 4, md: 6 } }}>
              <Typography variant="h3" sx={{ mb: { xs: 3, md: 5 } }}>
                {SOLUTION.caption}
              </Typography>

              <MotionViewport sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                {SOLUTION.items.map((item) => (
                  <Box
                    component={m.div}
                    variants={varFade('inUp')}
                    key={item.title}
                    sx={{ gap: 2, display: 'flex' }}
                  >
                    <Box
                      component="img"
                      loading="lazy"
                      alt={item.title}
                      src={item.icon}
                      sx={{ width: 44, height: 44, flexShrink: 0, objectFit: 'contain' }}
                    />
                    <Box>
                      <Typography variant="subtitle1" sx={{ color: 'primary.main' }}>
                        {item.title}
                      </Typography>
                      <Typography variant="body2" sx={{ mt: 0.25, color: 'text.secondary' }}>
                        {item.description}
                      </Typography>
                    </Box>
                  </Box>
                ))}
              </MotionViewport>
            </Box>

            <Box sx={{ flex: 1, position: 'relative', minHeight: { xs: 320, md: 'auto' } }}>
              <Box
                component="img"
                loading="lazy"
                alt={SOLUTION.caption}
                src={SOLUTION.image}
                sx={{ inset: 0, width: 1, height: 1, position: 'absolute', objectFit: 'cover' }}
              />
            </Box>
          </Box>
        </Card>
      </Container>
    </Box>
  );
}
