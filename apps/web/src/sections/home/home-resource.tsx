'use client';

import { m } from 'framer-motion';
import dynamic from 'next/dynamic';
import { useBoolean } from 'minimal-shared/hooks';

import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import ButtonBase from '@mui/material/ButtonBase';

import { Iconify } from 'src/components/iconify';
import { varFade, MotionViewport } from 'src/components/animate';

import { asset, RESOURCE } from './home-data';

// react-player berat — dialog baru dimuat (dan chunk-nya baru diunduh)
// setelah user mengklik tombol play, bukan ikut bundle landing.
const PlayerDialog = dynamic(
  () => import('src/components/player').then((mod) => mod.PlayerDialog),
  { ssr: false }
);

// ----------------------------------------------------------------------

export function HomeResource() {
  const openVideo = useBoolean();

  return (
    <Box
      component="section"
      sx={{
        py: { xs: 8, md: 12 },
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundImage: `url(${asset('binary-bg-light.webp')})`,
      }}
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
          <ButtonBase
            component={m.div}
            onClick={openVideo.onTrue}
            initial={{ opacity: 0, scale: 0.94 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            sx={{
              flex: 1,
              width: 1,
              borderRadius: 2,
              overflow: 'hidden',
              position: 'relative',
              display: 'block',
              boxShadow: (theme) => theme.customShadows.z16,
            }}
          >
            <Box
              component="img"
              loading="lazy"
              alt={RESOURCE.videoLabel}
              src={RESOURCE.videoThumb}
              sx={{ width: 1, display: 'block', aspectRatio: '16/9', objectFit: 'cover' }}
            />
            <Box
              sx={{
                inset: 0,
                display: 'flex',
                position: 'absolute',
                alignItems: 'center',
                justifyContent: 'center',
                bgcolor: 'rgba(0,0,0,0.12)',
              }}
            >
              <Iconify icon="solar:play-circle-bold" width={72} sx={{ color: 'error.main' }} />
            </Box>
          </ButtonBase>

          <MotionViewport sx={{ flex: 1 }}>
            <Typography component={m.h2} variants={varFade('inUp')} variant="h2" sx={{ mb: 3 }}>
              {RESOURCE.titleLines.map((line) => (
                <Box key={line} component="span" sx={{ display: 'block' }}>
                  {line}
                </Box>
              ))}
            </Typography>

            <Typography component={m.p} variants={varFade('inUp')} sx={{ color: 'text.secondary' }}>
              {RESOURCE.description}
            </Typography>
          </MotionViewport>
        </Box>
      </Container>

      {openVideo.value && (
        <PlayerDialog controls open playing src={RESOURCE.videoUrl} onClose={openVideo.onFalse} />
      )}
    </Box>
  );
}
