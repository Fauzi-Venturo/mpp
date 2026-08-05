import { it, expect, describe } from 'vitest';

import { qrFileName, renderQrSvg } from 'src/lib/qr/qr-svg';

describe('renderQrSvg', () => {
  it('renders an inline SVG for a token', async () => {
    const svg = await renderQrSvg('a1b2c3');

    expect(svg).toMatch(/^<svg/);
    expect(svg).toContain('</svg>');
  });

  it('renders a different code for a different token', async () => {
    const [first, second] = await Promise.all([renderQrSvg('token-one'), renderQrSvg('token-two')]);

    expect(first).not.toBe(second);
  });

  it('rejects an empty token instead of rendering a meaningless code', async () => {
    await expect(renderQrSvg('')).rejects.toThrow(/token/i);
  });
});

describe('qrFileName', () => {
  it('builds a downloadable name from the queue reference', async () => {
    expect(qrFileName('DUKCAPIL-2026-08-10')).toBe('qr-dukcapil-2026-08-10.svg');
  });

  it('strips characters that are unsafe in a file name', async () => {
    expect(qrFileName('A/014 booking')).toBe('qr-a-014-booking.svg');
  });
});
