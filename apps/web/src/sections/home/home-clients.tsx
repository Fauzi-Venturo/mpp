'use client';

import AutoScroll from 'embla-carousel-auto-scroll';

import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import useMediaQuery from '@mui/material/useMediaQuery';

import { Iconify } from 'src/components/iconify';
import { Carousel, useCarousel } from 'src/components/carousel';

import { HERO, CLIENTS } from './home-data';

// ----------------------------------------------------------------------

export function HomeClients() {
  // WCAG 2.3.3: marquee tidak auto-jalan bila user memilih reduced motion.
  const prefersReducedMotion = useMediaQuery('(prefers-reduced-motion: reduce)');

  const carousel = useCarousel(
    { loop: true, slidesToShow: { xs: 2.5, sm: 4, md: 6 }, slideSpacing: '24px' },
    [
      AutoScroll({
        playOnInit: !prefersReducedMotion,
        speed: 1,
        stopOnInteraction: false,
        stopOnMouseEnter: true,
      }),
    ]
  );

  if (!CLIENTS.length) {
    return null;
  }

  return (
    <Box component="section" sx={{ pb: { xs: 6, md: 10 } }}>
      <Container>
        <Box sx={{ mb: 3, gap: 1, display: 'flex', alignItems: 'center' }}>
          <Typography variant="subtitle2" sx={{ color: 'text.secondary' }}>
            {HERO.trustedLabel}
          </Typography>

          {/* WCAG 2.2.2: gerakan > 5 detik wajib punya kontrol pause */}
          <IconButton
            size="small"
            onClick={carousel.autoScroll.onTogglePlay}
            aria-label={
              carousel.autoScroll.isPlaying ? 'Jeda animasi logo klien' : 'Putar animasi logo klien'
            }
            sx={{ color: 'text.disabled' }}
          >
            <Iconify
              width={18}
              icon={carousel.autoScroll.isPlaying ? 'carbon:pause' : 'carbon:play'}
            />
          </IconButton>
        </Box>

        <Carousel
          carousel={carousel}
          sx={{
            // soft fade on both edges of the marquee
            maskImage: 'linear-gradient(to right, transparent, #000 8%, #000 92%, transparent)',
          }}
        >
          {CLIENTS.map((client) => (
            <Box
              key={client.name}
              sx={{ height: 56, display: 'flex', alignItems: 'center', justifyContent: 'center' }}
            >
              <Box
                component="img"
                loading="lazy"
                alt={client.name}
                src={client.logo}
                sx={{ maxWidth: 150, maxHeight: 36, objectFit: 'contain' }}
              />
            </Box>
          ))}
        </Carousel>
      </Container>
    </Box>
  );
}
