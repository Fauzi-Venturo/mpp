'use client';

import { m } from 'framer-motion';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { Iconify } from 'src/components/iconify';
import { varFade, MotionViewport } from 'src/components/animate';

import { CONTACT, SPECIAL_OFFER } from './home-data';

// ----------------------------------------------------------------------

type HomeSpecialOfferProps = {
  waLink?: string | null;
};

export function HomeSpecialOffer({ waLink }: HomeSpecialOfferProps) {
  return (
    <Box
      component="section"
      sx={{ py: { xs: 6, md: 8 }, overflow: 'hidden', bgcolor: 'background.neutral' }}
    >
      <Container>
        <Box
          sx={{
            mx: 'auto',
            maxWidth: 860,
            gap: { xs: 3, md: 6 },
            display: 'flex',
            alignItems: 'center',
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
            alt={SPECIAL_OFFER.heading}
            src={SPECIAL_OFFER.image}
            sx={{ width: 'auto', height: { xs: 300, md: 420 }, flexShrink: 0 }}
          />

          <MotionViewport sx={{ flex: 1 }}>
            <Typography
              component={m.h2}
              variants={varFade('inUp')}
              variant="h2"
              sx={{ mb: 2, fontStyle: 'italic', color: 'primary.main' }}
            >
              {SPECIAL_OFFER.heading}
            </Typography>

            <Typography component={m.p} variants={varFade('inUp')} sx={{ mb: 4, maxWidth: 460 }}>
              <Box component="span" sx={{ fontWeight: 'fontWeightBold' }}>
                {SPECIAL_OFFER.promoStrong}{' '}
              </Box>
              {SPECIAL_OFFER.promo}
            </Typography>

            <Box component={m.div} variants={varFade('inUp')}>
              <Button
                size="large"
                color="primary"
                variant="contained"
                href={waLink ?? CONTACT.wa}
                target="_blank"
                rel="noopener noreferrer"
                startIcon={<Iconify icon="solar:chat-round-call-linear" />}
              >
                {SPECIAL_OFFER.cta}
              </Button>

              <Box sx={{ mt: 2, gap: 1, display: 'flex', alignItems: 'center' }}>
                <Iconify icon="solar:check-circle-bold" sx={{ color: 'primary.main' }} />
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  {SPECIAL_OFFER.note}
                </Typography>
              </Box>
            </Box>
          </MotionViewport>
        </Box>
      </Container>
    </Box>
  );
}
