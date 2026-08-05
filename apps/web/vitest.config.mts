import path from 'node:path';

import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [react()],
  resolve: {
    // tsconfig uses baseUrl "." so imports are rooted at src/…
    alias: { src: path.resolve(import.meta.dirname, 'src') },
  },
  test: {
    environment: 'happy-dom',
    globals: true,
    env: {
      // Storage/API are UTC; the UI renders WIB — pin the zone so date assertions
      // don't drift with the machine running the suite.
      TZ: 'Asia/Jakarta',
      // src/lib/env.ts validates at import time, so any module reaching the API
      // layer needs these. Fixed dummies keep tests independent of a local .env.
      NEXT_PUBLIC_API_URL: 'http://api.test',
      NEXT_PUBLIC_COMPANY_SLUG: 'mpp-test',
    },
    setupFiles: ['./vitest.setup.ts'],
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
  },
});
