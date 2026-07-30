import type { Metadata } from 'next';

import { CONFIG } from 'src/global-config';
import { getFaqGroups, getWhatsAppLink } from 'src/lib/api';
import { webSiteJsonLd, toJsonLdScript, organizationJsonLd } from 'src/lib/seo';

import { CONTACT } from 'src/sections/home/home-data';
import { HomeView } from 'src/sections/home/view/home-view';

// ----------------------------------------------------------------------

const TITLE = 'Jasa Outsourcing Programmer Malang | Software House Venturo';
const DESCRIPTION =
  'Venturo adalah jasa outsource programmer terbesar di Malang dengan 130+ talenta dedicated team. Garansi 30 hari, supervisor berpengalaman, laporan progres mingguan. Konsultasi gratis.';

export const metadata: Metadata = {
  // Full custom title — `absolute` opts out of the root '%s - Venturo' template.
  title: { absolute: TITLE },
  description: DESCRIPTION,
  alternates: { canonical: '/' },
  // A page-level openGraph REPLACES the root layout's whole openGraph object
  // (and drops the file-convention og:image), so restate everything here.
  openGraph: {
    type: 'website',
    siteName: CONFIG.appName,
    locale: 'id_ID',
    url: '/',
    title: TITLE,
    description: DESCRIPTION,
    images: ['/opengraph-image.png'],
  },
};

export default async function Page() {
  // ISR via the fetch's revalidate (300s). Backend down → null → sections
  // fall back to their static copy, so home never breaks.
  const [faqGroups, waLink] = await Promise.all([
    getFaqGroups('id').catch(() => null),
    getWhatsAppLink('id').catch(() => null),
  ]);

  // Organization + WebSite structured data (knowledge panel / brand info).
  const jsonLd = [
    organizationJsonLd({ waLink: waLink ?? CONTACT.wa, email: CONTACT.email }),
    webSiteJsonLd(),
  ];

  return (
    <>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: toJsonLdScript(jsonLd) }}
      />

      <HomeView faqGroups={faqGroups} waLink={waLink} />
    </>
  );
}
