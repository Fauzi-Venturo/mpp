import type { MetadataRoute } from 'next';

import { CONFIG } from 'src/global-config';

// ----------------------------------------------------------------------

export default function robots(): MetadataRoute.Robots {
  return {
    rules: {
      userAgent: '*',
      allow: '/',
      disallow: ['/components/', '/blank/', '/error/', '/coming-soon/', '/maintenance/'],
    },
    sitemap: `${CONFIG.siteUrl}/sitemap.xml`,
  };
}
