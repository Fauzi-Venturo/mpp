'use client';

import { Fragment } from 'react';
import { varAlpha } from 'minimal-shared/utils';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { Image } from 'src/components/image';
import { Iconify } from 'src/components/iconify';
import { AnimateCountUp } from 'src/components/animate';

import { HERO, asset, CONTACT } from './home-data';

// Above-the-fold = LCP: no entry animations here — SSR HTML must ship the
// heading and the <img> fully visible (crawlers + LCP can't wait for JS).

// ----------------------------------------------------------------------

type HomeHeroProps = {
  waLink?: string | null;
};

export function HomeHero({ waLink }: HomeHeroProps) {
  return (
    <Box
      component="section"
      sx={(theme) => ({
        overflow: 'hidden',
        position: 'relative',
        py: { xs: 8, md: 14 },
        background: `linear-gradient(180deg, ${varAlpha(theme.vars.palette.primary.lighterChannel, 0.4)}, ${varAlpha(theme.vars.palette.primary.lighterChannel, 0)})`,
      })}
    >
      <Container>
        <Box
          sx={{
            gap: { xs: 5, md: 8 },
            display: 'flex',
            alignItems: 'center',
            flexDirection: { xs: 'column', md: 'row' },
          }}
        >
          <Box sx={{ flex: 1 }}>
            <Box
              sx={{
                px: 1.5,
                py: 0.5,
                mb: 3,
                borderRadius: 1,
                typography: 'subtitle2',
                color: 'primary.dark',
                display: 'inline-flex',
                alignItems: 'center',
                gap: 0.75,
                bgcolor: 'primary.lighter',
              }}
            >
              <Iconify width={18} icon="solar:check-circle-bold" />
              {HERO.badge}
            </Box>

            <Typography component="h1" variant="h1" sx={{ mb: 3 }}>
              {HERO.title}
            </Typography>

            <Typography sx={{ mb: 4, maxWidth: 520, color: 'text.secondary' }}>
              {HERO.description}
            </Typography>

            <Box sx={{ mb: 5 }}>
              <Button
                size="large"
                color="primary"
                variant="contained"
                href={waLink ?? CONTACT.wa}
                target="_blank"
                rel="noopener noreferrer"
                endIcon={<Iconify icon="eva:arrow-ios-forward-fill" />}
              >
                {HERO.cta}
              </Button>
            </Box>
          </Box>

          <Box sx={{ flex: 1, width: 1 }}>
            <Box sx={{ width: 1, position: 'relative' }}>
              <Box
                aria-hidden
                sx={{
                  inset: 0,
                  borderRadius: 3,
                  position: 'absolute',
                  bgcolor: 'primary.lighter',
                  transform: 'translate(16px, 16px)',
                }}
              />
              <Image
                alt={HERO.title}
                src={asset('hero-team.webp')}
                ratio="16/9"
                visibleByDefault
                slotProps={{ img: { fetchPriority: 'high' } }}
                sx={{
                  position: 'relative',
                  borderRadius: 3,
                  boxShadow: (theme) => theme.customShadows.z16,
                }}
              />
            </Box>

            <Box
              sx={{
                mt: 4,
                gap: 4,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {HERO.stats.map((stat, index) => (
                <Fragment key={stat.label}>
                  {index > 0 && (
                    <Box sx={{ width: '1px', alignSelf: 'stretch', bgcolor: 'divider' }} />
                  )}
                  <Box sx={{ textAlign: 'center' }}>
                    <Box
                      sx={{
                        color: 'primary.main',
                        display: 'flex',
                        alignItems: 'baseline',
                        justifyContent: 'center',
                      }}
                    >
                      <AnimateCountUp to={stat.value} sx={{ typography: 'h2' }} />
                      <Typography component="span" variant="h3">
                        {stat.suffix}
                      </Typography>
                    </Box>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      {stat.label}
                    </Typography>
                  </Box>
                </Fragment>
              ))}
            </Box>
          </Box>
        </Box>
      </Container>
    </Box>
  );
}
