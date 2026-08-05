'use client';

import type { FaqGroup } from 'src/lib/api';

import { HomeCta } from '../home-cta';
import { HomeHero } from '../home-hero';
import { HomeFaqs } from '../home-faqs';
import { HomeFocus } from '../home-focus';
import { HomeClients } from '../home-clients';
import { HomeProblem } from '../home-problem';
import { HomeSolution } from '../home-solution';
import { HomeResource } from '../home-resource';
import { HomeTechStack } from '../home-tech-stack';
import { HomeManagement } from '../home-management';
import { HomeFloatingCta } from '../home-floating-cta';
import { HomeSpecialOffer } from '../home-special-offer';

// ----------------------------------------------------------------------
// Data: RSC props dari src/app/(home)/page.tsx — faqGroups (FAQ API) + waLink
// (site-content API), keduanya ISR 300 dtk. Prop null = backend mati → tiap
// section jatuh ke konten statis home-data.ts; home tidak pernah pecah.

type HomeViewProps = {
  faqGroups?: FaqGroup[] | null;
  waLink?: string | null;
};

export function HomeView({ faqGroups, waLink }: HomeViewProps) {
  return (
    <>
      <HomeHero waLink={waLink} />
      <HomeClients />
      <HomeProblem />
      <HomeSolution />
      <HomeFocus />
      <HomeManagement />
      <HomeResource />
      <HomeSpecialOffer waLink={waLink} />
      <HomeTechStack />
      <HomeFaqs groups={faqGroups} />
      <HomeCta waLink={waLink} />
      <HomeFloatingCta waLink={waLink} />
    </>
  );
}
