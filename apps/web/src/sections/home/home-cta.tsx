'use client';

import { m } from 'framer-motion';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { Iconify } from 'src/components/iconify';
import { varFade, MotionViewport } from 'src/components/animate';

import { CONTACT, CLOSING_CTA } from './home-data';

// ----------------------------------------------------------------------

type HomeCtaProps = {
  waLink?: string | null;
};

export function HomeCta({ waLink }: HomeCtaProps) {
  return (
    <Box
      component="section"
      sx={(theme) => ({
        overflow: 'hidden',
        color: 'common.white',
        background: `linear-gradient(120deg, ${theme.vars.palette.primary.light}, ${theme.vars.palette.primary.main} 55%, ${theme.vars.palette.primary.dark})`,
      })}
    >
      <Container>
        <Box
          sx={{
            gap: { xs: 3, md: 6 },
            display: 'flex',
            alignItems: 'flex-end',
            flexDirection: { xs: 'column', md: 'row' },
          }}
        >
          <Box
            component={m.img}
            loading="lazy"
            initial={{ opacity: 0, y: 24 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            alt={CLOSING_CTA.title}
            src={CLOSING_CTA.image}
            sx={{ width: { xs: 220, md: 320 }, height: 'auto', flexShrink: 0 }}
          />

          <MotionViewport sx={{ flex: 1, py: { xs: 0, md: 8 }, pb: { xs: 6, md: 8 } }}>
            <Typography
              component={m.h2}
              variants={varFade('inUp')}
              variant="h2"
              sx={{ mb: 3, fontStyle: 'italic' }}
            >
              {CLOSING_CTA.title}
            </Typography>

            <Typography
              component={m.p}
              variants={varFade('inUp')}
              sx={{ mb: 4, maxWidth: 520, opacity: 0.9 }}
            >
              {CLOSING_CTA.description}{' '}
              <Box component="span" sx={{ fontWeight: 'fontWeightBold' }}>
                {CLOSING_CTA.descriptionStrong}
              </Box>
              . {CLOSING_CTA.descriptionEnd}
            </Typography>

            <Box component={m.div} variants={varFade('inUp')}>
              <Button
                size="large"
                variant="contained"
                href={waLink ?? CONTACT.wa}
                target="_blank"
                rel="noopener noreferrer"
                startIcon={<Iconify icon="solar:chat-round-call-linear" />}
                sx={{
                  color: 'primary.dark',
                  bgcolor: 'common.white',
                  '&:hover': { bgcolor: 'grey.200' },
                }}
              >
                {CLOSING_CTA.cta}
              </Button>
            </Box>
          </MotionViewport>
        </Box>
      </Container>
    </Box>
  );
}
