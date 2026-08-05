import type { MetadataRoute } from 'next';

import { primary } from 'src/theme';
import { CONFIG } from 'src/global-config';

// ----------------------------------------------------------------------
// Web app manifest (file convention → /manifest.webmanifest).
// Catatan per-client: untuk PWA penuh tambahkan ikon 192x192 & 512x512.

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: CONFIG.appName,
    short_name: CONFIG.appName,
    description:
      'Venturo adalah software house di Malang dengan 130+ talenta dedicated team untuk pengembangan software, outsourcing programmer, dan konsultasi IT.',
    start_url: '/',
    display: 'browser',
    background_color: '#ffffff',
    theme_color: primary.main,
    icons: [
      { src: '/favicon.ico', sizes: 'any' },
      { src: '/apple-touch-icon.png', sizes: '180x180', type: 'image/png' },
    ],
  };
}
