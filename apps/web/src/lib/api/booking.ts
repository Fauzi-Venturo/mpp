import { z } from 'zod';

import { apiFetch } from 'src/lib/api/client';
import { endpoints } from 'src/lib/api/endpoints';

// ----------------------------------------------------------------------
// Contract: docs/04-api/rest-endpoints.md — GET /mpp/v1/booking/{id}
// Shape mirrors apps/api/internal/modules/mpp/booking/domain/booking.go
//
// qr_token / qr_expires_at are REQUIRED: a ticket screen without them would
// render a blank QR, so contract drift must fail loudly here.

export const pemohonSchema = z.object({
  id: z.string(),
  name: z.string(),
  phone: z.string().nullish(),
  email: z.string().nullish(),
});

export const bookingSchema = z.object({
  id: z.string(),
  pemohon_id: z.string(),
  instansi_id: z.string(),
  jenis_layanan_id: z.string(),
  tanggal: z.string(),
  channel: z.string(),
  status: z.string(),
  qr_token: z.string().min(1),
  qr_expires_at: z.string(),
  created_at: z.string().optional(),
  pemohon: pemohonSchema.optional(),
});

export type Booking = z.infer<typeof bookingSchema>;

export const bookingKeys = {
  all: ['booking'] as const,
  detail: (id: string) => [...bookingKeys.all, 'detail', id] as const,
};

/**
 * Fetches one booking with its check-in token.
 *
 * Never cached: the ticket is per-citizen and time-bound (BR-09), so a shared
 * Data Cache entry would be both stale and a privacy leak.
 */
export async function getBooking(id: string, options: { signal?: AbortSignal } = {}) {
  const { data } = await apiFetch<unknown>(endpoints.booking.details(id), {
    signal: options.signal,
    cache: 'no-store',
  });

  return bookingSchema.parse(data);
}
