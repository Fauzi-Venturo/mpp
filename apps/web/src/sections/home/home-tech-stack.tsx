'use client';

import Box from '@mui/material/Box';
import Divider from '@mui/material/Divider';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';

import { HomeHeading } from './home-heading';
import { asset, TECH_STACK } from './home-data';

// ----------------------------------------------------------------------

type Group = (typeof TECH_STACK.groups)[number];

function Category({ group }: { group: Group }) {
  return (
    <Box>
      <Divider sx={{ mb: { xs: 4, md: 5 } }}>
        <Typography variant="h6" sx={{ px: 1, color: 'primary.main' }}>
          {group.label}
        </Typography>
      </Divider>

      <Box
        sx={{
          px: 2,
          rowGap: { xs: 4, md: 5 },
          columnGap: { xs: 4, md: 6 },
          display: 'flex',
          flexWrap: 'wrap',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        {group.logos.map(({ name, logo }) => (
          <Box
            key={name}
            component="img"
            loading="lazy"
            alt={name}
            src={asset(logo)}
            sx={{
              maxWidth: 160,
              maxHeight: 32,
              objectFit: 'contain',
              transition: (theme) => theme.transitions.create('transform'),
              '&:hover': { transform: 'scale(1.08)' },
            }}
          />
        ))}
      </Box>
    </Box>
  );
}

export function HomeTechStack() {
  // TECH_STACK.groups selalu 3 entri (lihat home-data.ts)
  const [web, mobile, ui] = TECH_STACK.groups as [Group, Group, Group];

  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 } }}>
      <Container>
        <HomeHeading
          title={TECH_STACK.caption}
          description="Stack modern dan teruji yang kami gunakan untuk membangun produk digital Anda."
        />

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: { xs: 6, md: 9 } }}>
          <Category group={web} />

          <Box
            sx={{
              gap: { xs: 6, md: 6 },
              display: 'grid',
              gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' },
            }}
          >
            <Category group={mobile} />
            <Category group={ui} />
          </Box>
        </Box>
      </Container>
    </Box>
  );
}
