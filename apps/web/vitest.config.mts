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
    // Storage/API are UTC; the UI renders WIB — pin the zone so date assertions
    // don't drift with the machine running the suite.
    env: { TZ: 'Asia/Jakarta' },
    setupFiles: ['./vitest.setup.ts'],
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
  },
});
