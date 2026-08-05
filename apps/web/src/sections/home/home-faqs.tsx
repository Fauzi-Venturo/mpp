'use client';

import type { FaqGroup } from 'src/lib/api';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Accordion from '@mui/material/Accordion';
import Typography from '@mui/material/Typography';
import AccordionDetails from '@mui/material/AccordionDetails';
import AccordionSummary from '@mui/material/AccordionSummary';

import { paths } from 'src/routes/paths';
import { RouterLink } from 'src/routes/components';

import { Iconify } from 'src/components/iconify';

import { FAQS } from './home-data';
import { HomeHeading } from './home-heading';

// ----------------------------------------------------------------------

// Static copy as fallback when the FAQ API is unreachable or empty.
const FALLBACK_GROUPS: FaqGroup[] = [
  {
    key: 'fallback',
    title: '',
    entries: FAQS.items.map((item) => ({ question: item.question, answer: item.answer })),
  },
];

// Home shows a teaser only — the full list lives at /support (FAQ page).
const HOME_FAQ_LIMIT = 6;

type HomeFaqsProps = {
  groups?: FaqGroup[] | null;
};

export function HomeFaqs({ groups }: HomeFaqsProps) {
  const resolvedGroups = groups?.length ? groups : FALLBACK_GROUPS;

  const entries = resolvedGroups.flatMap((group) => group.entries).slice(0, HOME_FAQ_LIMIT);

  return (
    <Box component="section" sx={{ py: { xs: 8, md: 12 } }}>
      <Container>
        <HomeHeading caption={FAQS.caption} title={FAQS.title} />

        <Box sx={{ mx: 'auto', maxWidth: 760 }}>
          {entries.map((entry) => (
            <Accordion key={entry.question}>
              <AccordionSummary expandIcon={<Iconify icon="eva:arrow-ios-downward-fill" />}>
                <Typography variant="subtitle1">{entry.question}</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <Typography sx={{ color: 'text.secondary' }}>{entry.answer}</Typography>
              </AccordionDetails>
            </Accordion>
          ))}

          <Box sx={{ mt: 5, textAlign: 'center' }}>
            <Button
              component={RouterLink}
              href={paths.support}
              size="large"
              color="primary"
              variant="outlined"
              endIcon={<Iconify icon="eva:arrow-ios-forward-fill" />}
            >
              Lihat FAQ Selengkapnya
            </Button>
          </Box>
        </Box>
      </Container>
    </Box>
  );
}
