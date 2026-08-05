import { it, expect, describe } from 'vitest';

import { fDate, fDateTime } from 'src/utils/format-time';

describe('fDate', () => {
  it('returns an empty string when there is no input', () => {
    expect(fDate('')).toBe('');
  });

  it('returns Invalid for an unparsable input', () => {
    expect(fDate('not-a-date')).toBe('Invalid');
  });

  it('formats an ISO date with the default pattern', () => {
    expect(fDate('2026-04-17')).toBe('17 Apr 2026');
  });

  it('honours a custom pattern', () => {
    expect(fDate('2026-04-17', 'DD/MM/YYYY')).toBe('17/04/2026');
  });
});

describe('fDateTime', () => {
  it('includes the time part', () => {
    expect(fDateTime('2026-04-17T09:30:00')).toBe('17 Apr 2026 9:30 am');
  });
});
