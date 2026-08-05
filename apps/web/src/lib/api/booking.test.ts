import { it, expect, describe } from 'vitest';

import { bookingKeys, bookingSchema } from 'src/lib/api/booking';

// Mirrors the Go response of GET /mpp/v1/booking/{id}
// (apps/api/internal/modules/mpp/booking/domain/booking.go).
const payload = {
  id: '9f1d0c3a-0000-0000-0000-000000000001',
  pemohon_id: '8e1d0c3a-0000-0000-0000-000000000002',
  instansi_id: 'a2000000-0000-0000-0000-000000000001',
  jenis_layanan_id: 'b3000000-0000-0000-0000-000000000003',
  tanggal: '2026-08-10',
  channel: 'WEB',
  status: 'BOOKED',
  qr_token: 'f1e2d3c4b5a6',
  qr_expires_at: '2026-08-11T00:00:00Z',
  created_at: '2026-08-05T03:00:00Z',
  pemohon: { id: '8e1d0c3a-0000-0000-0000-000000000002', name: 'Budi Santoso' },
};

describe('bookingSchema', () => {
  it('accepts the backend payload', () => {
    const booking = bookingSchema.parse(payload);

    expect(booking.qr_token).toBe('f1e2d3c4b5a6');
    expect(booking.qr_expires_at).toBe('2026-08-11T00:00:00Z');
    expect(booking.pemohon?.name).toBe('Budi Santoso');
  });

  it('rejects a response without a QR token, instead of rendering a blank ticket', () => {
    const { qr_token: _omitted, ...withoutToken } = payload;

    expect(() => bookingSchema.parse(withoutToken)).toThrow();
  });

  it('tolerates a missing applicant block', () => {
    const { pemohon: _omitted, ...withoutPemohon } = payload;

    expect(() => bookingSchema.parse(withoutPemohon)).not.toThrow();
  });
});

describe('bookingKeys', () => {
  it('scopes the detail key by id', () => {
    expect(bookingKeys.detail('abc')).toEqual(['booking', 'detail', 'abc']);
  });
});
