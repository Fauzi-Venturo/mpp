'use client';

import type { FaqGroup } from 'src/lib/api';

import { useBoolean } from 'minimal-shared/hooks';
import { useState, useEffect, useCallback } from 'react';

import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';

import { CONFIG } from 'src/global-config';

import { Iconify } from 'src/components/iconify';

// Data-only import: fallback FAQ statis yang sama dengan section home,
// supaya halaman ini tetap berisi saat API mati.
import { FAQS } from 'src/sections/home/home-data';

import { SupportNav } from '../support-nav';
import { SupportHero } from '../support-hero';
import { SupportContent } from '../support-content';

// ----------------------------------------------------------------------
// Data: RSC props dari src/app/support/page.tsx — faqGroups (FAQ API, tiap
// sub-group jadi topik di layout dua kolom) + waLink (site-content API), ISR
// 300 dtk. faqGroups null/kosong → FALLBACK_GROUPS; waLink null → CONTACT.wa.

const FALLBACK_GROUPS: FaqGroup[] = [
  {
    key: 'fallback',
    title: 'Umum',
    entries: FAQS.items.map((item) => ({ question: item.question, answer: item.answer })),
  },
];

// Grup dari API tidak membawa ikon — daftar ikon topik dirotasi berurutan.
const TOPIC_ICONS = [
  'ic-account.svg',
  'ic-payment.svg',
  'ic-delivery.svg',
  'ic-package.svg',
  'ic-refund.svg',
  'ic-assurances.svg',
];

const iconPath = (name: string) => `${CONFIG.assetsDir}/assets/icons/support/${name}`;

type SupportViewProps = {
  faqGroups?: FaqGroup[] | null;
  waLink?: string | null;
};

export function SupportView({ faqGroups, waLink }: SupportViewProps) {
  const resolvedGroups = faqGroups?.length ? faqGroups : FALLBACK_GROUPS;

  const topics = resolvedGroups.map((group, index) => ({
    key: group.key,
    title: group.title,
    icon: iconPath(TOPIC_ICONS[index % TOPIC_ICONS.length]!),
    entries: group.entries,
  }));

  const [topic, setTopic] = useState(topics[0]?.title ?? '');

  const openNavMobile = useBoolean();

  const handleChangeTopic = useCallback((event: React.SyntheticEvent, newValue: string) => {
    setTopic(newValue);
  }, []);

  useEffect(() => {
    if (openNavMobile.value) {
      openNavMobile.onFalse();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [topic]);

  return (
    <>
      <SupportHero />
      <Box
        sx={(theme) => ({
          px: 2,
          py: 1.5,
          display: { md: 'none' },
          borderBottom: `solid 1px ${theme.vars.palette.divider}`,
        })}
      >
        <IconButton onClick={openNavMobile.onTrue} aria-label="Buka daftar topik">
          <Iconify icon="carbon:menu" />
        </IconButton>
      </Box>

      <Container component="section" sx={{ pb: { xs: 10, md: 15 } }}>
        <Typography variant="h3" component="h2" sx={{ my: { xs: 3, md: 10 } }}>
          Pertanyaan yang Sering Diajukan
        </Typography>

        <Box sx={{ gap: 10, display: 'flex' }}>
          <SupportNav
            data={topics}
            topic={topic}
            waLink={waLink}
            open={openNavMobile.value}
            onChangeTopic={handleChangeTopic}
            onClose={openNavMobile.onFalse}
          />

          {topics.map(
            (item) =>
              item.title === topic && (
                <Box key={item.key} sx={{ flex: 1, minWidth: 0 }}>
                  <SupportContent contents={item.entries} />
                </Box>
              )
          )}
        </Box>
      </Container>
    </>
  );
}
