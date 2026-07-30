'use client';

import { m } from 'framer-motion';

import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { varFade, MotionViewport } from 'src/components/animate';

import { MANAGEMENT } from './home-data';
import { HomeHeading } from './home-heading';

// ----------------------------------------------------------------------

type Feature = (typeof MANAGEMENT.items)[number];

function FeatureItem({ item }: { item: Feature }) {
  return (
    <Box component={m.div} variants={varFade('inUp')} sx={{ gap: 2, display: 'flex' }}>
      <Box
        component="img"
        loading="lazy"
        alt={item.title}
        src={item.icon}
        sx={{ width: 40, height: 40, flexShrink: 0, objectFit: 'contain' }}
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
  );
}

export function HomeManagement() {
  const left = MANAGEMENT.items.slice(0, 3);
  const right = MANAGEMENT.items.slice(3, 6);

  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 }, bgcolor: 'background.neutral' }}>
      <Container>
        <HomeHeading
          caption={MANAGEMENT.caption}
          title={MANAGEMENT.title}
          description={MANAGEMENT.description}
        />

        <Box
          sx={{
            gap: { xs: 5, md: 4 },
            display: 'flex',
            alignItems: 'center',
            flexDirection: { xs: 'column', md: 'row' },
          }}
        >
          <MotionViewport sx={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 4 }}>
            {left.map((item) => (
              <FeatureItem key={item.title} item={item} />
            ))}
          </MotionViewport>

          <Box
            component={m.img}
            loading="lazy"
            initial={{ opacity: 0, y: 24 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.6 }}
            alt="Venturo app"
            src={MANAGEMENT.phone}
            sx={{ flexShrink: 0, width: { xs: 220, md: 280 }, height: 'auto' }}
          />

          <MotionViewport sx={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 4 }}>
            {right.map((item) => (
              <FeatureItem key={item.title} item={item} />
            ))}
          </MotionViewport>
        </Box>
      </Container>
    </Box>
  );
}
