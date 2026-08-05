'use client';

import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Typography from '@mui/material/Typography';

import { asset, CONTACT } from './home-data';

// ----------------------------------------------------------------------

type HomeFloatingCtaProps = {
  waLink?: string | null;
};

export function HomeFloatingCta({ waLink }: HomeFloatingCtaProps) {
  return (
    <Link
      href={waLink ?? CONTACT.wa}
      target="_blank"
      rel="noopener noreferrer"
      underline="none"
      sx={{
        px: 2,
        py: 1.25,
        gap: 1,
        display: 'flex',
        borderRadius: 2,
        position: 'fixed',
        alignItems: 'center',
        color: 'common.white',
        bgcolor: '#0E7A3F',
        right: { xs: 16, md: 24 },
        bottom: { xs: 16, md: 24 },
        zIndex: (theme) => theme.zIndex.speedDial,
        boxShadow: (theme) => theme.customShadows.z16,
        transition: (theme) => theme.transitions.create('background-color'),
        '&:hover': { bgcolor: '#0B6332' },
      }}
    >
      <Box
        component="img"
        alt="WhatsApp"
        src={asset('ikon-whatsapp-venturo.webp')}
        sx={{ width: 28, height: 28, objectFit: 'contain' }}
      />
      <Box>
        <Typography variant="subtitle2" sx={{ lineHeight: 1.1 }}>
          Hubungi Kami
        </Typography>
        <Typography variant="caption">Konsultasi 100% Gratis</Typography>
      </Box>
    </Link>
  );
}
